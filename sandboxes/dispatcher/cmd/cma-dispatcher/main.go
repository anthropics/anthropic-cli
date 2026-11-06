// Command cma-dispatcher gives each CMA session of one environment its own
// Kubernetes Agent Sandbox pod. It runs as one replica per environment per
// namespace, polls the environment's work queue with the environment key,
// creates a SandboxClaim per session, hands the pod a session-scoped token
// over HTTPS, acks the work item only after the pod has collected it, and
// parks or deletes the sandbox when its worker exits.
//
// Subcommands:
//
//	serve      run the dispatcher (the default when no subcommand is given)
//	preflight  check an install with your own kubeconfig; non-zero on any failure
//
// ANTHROPIC_ENVIRONMENT_KEY, ANTHROPIC_ENVIRONMENT_ID and ANTHROPIC_BASE_URL
// come from the environment; every other configuration setting is a flag
// with a CMA_DISPATCHER_* variable as its default (-h lists them).
package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/cmaapi"
	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/config"
	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/dispatcher"
	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/httpapi"
	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/nonce"
	checks "github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/preflight"
	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/sandbox"
	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/version"
	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/workitem"
)

func main() {
	cmd, args := "serve", os.Args[1:]
	if len(args) > 0 && args[0] != "" && args[0][0] != '-' {
		cmd, args = args[0], args[1:]
	}
	fs := flag.NewFlagSet("cma-dispatcher "+cmd, flag.ExitOnError)
	cfg := config.Bind(fs)
	showVersion := fs.Bool("version", false, "print the version and exit")
	switch cmd {
	case "serve", "preflight":
		_ = fs.Parse(args)
	case "help", "-h", "-help", "--help":
		fmt.Fprintln(os.Stderr, "usage: cma-dispatcher [serve|preflight] [flags]")
		fs.PrintDefaults()
		return
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand %q (want serve or preflight)\n", cmd)
		os.Exit(2)
	}
	if *showVersion {
		// Ahead of cfg.Load: a version query needs no environment key,
		// namespace or cluster.
		fmt.Println("cma-dispatcher", version.String())
		return
	}
	if err := cfg.Load(); err != nil {
		fatal(err)
	}
	shape, err := sandbox.ShapeFor(cfg.K8SAPIVersion)
	if err != nil {
		fatal(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	switch cmd {
	case "serve":
		if err := serve(ctx, cfg, shape); err != nil {
			fatal(err)
		}
	case "preflight":
		if err := preflight(ctx, cfg, shape); err != nil {
			fatal(err)
		}
	}
}

// serve verifies the environment key, connects to the cluster, serves the
// handoff API to sandbox pods, and polls the work queue until ctx ends.
func serve(ctx context.Context, cfg *config.Config, shape sandbox.Shape) error {
	logger := newLogger(cfg.LogFormat)
	// Logged before anything in serve that can fail. A start that cannot reach
	// the API or the cluster still shows its version.
	logger.Info("cma-dispatcher starting", "version", version.String())
	warnHTTP2Debug(logger, os.Getenv("GODEBUG"))
	api := cmaapi.New(cfg.BaseURL, cfg.EnvironmentID, cfg.EnvironmentKey, cfg.APIRequestTimeout)
	if err := api.VerifyEnvironment(ctx); err != nil {
		return err
	}
	kube, dyn, err := clusterClients()
	if err != nil {
		return err
	}
	sandboxes, auth, err := sandboxClientAndAuth(kube, dyn, cfg, shape)
	if err != nil {
		return err
	}
	maxHandoffs := cfg.MaxConcurrentHandoffs
	if maxHandoffs == 0 {
		maxHandoffs = handoffLimit(ctx, shape, sandboxes, logger)
	}
	workItems := workitem.NewStore()
	server := &httpapi.Server{Auth: auth, WorkItems: workItems, Logger: logger}
	d := dispatcher.New(api, api, sandboxes, workItems, logger, dispatcher.Options{
		EnvironmentID:         cfg.EnvironmentID,
		BaseURL:               cfg.BaseURL,
		WorkerID:              workerID(),
		HandoffWait:           cfg.HandoffWait,
		PodPollInterval:       2 * time.Second,
		SweepInterval:         cfg.SweepInterval,
		Retention:             cfg.Retention,
		MaxConcurrentHandoffs: maxHandoffs,
	})

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if err := watchPods(ctx, kube, cfg.K8SNamespace, d); err != nil {
		return err
	}
	errc := make(chan error, 2)
	go func() {
		logger.Info("listening", "addr", cfg.Listen)
		errc <- httpapi.ListenAndServe(ctx, cfg.Listen, cfg.TLSCertFile, cfg.TLSKeyFile, server.Handler())
	}()
	go func() {
		logger.Info("dispatching", "environment_id", cfg.EnvironmentID, "namespace", cfg.K8SNamespace, "api_version", shape.Version(), "max_handoffs", maxHandoffs, "worker_id", d.Options.WorkerID)
		errc <- d.Run(ctx)
	}()
	err = <-errc
	cancel()
	<-errc
	return err
}

// watchPods feeds every pod change in the namespace to the dispatcher until
// ctx ends, so a worker's exit is seen as the kubelet reports it.
func watchPods(ctx context.Context, kube kubernetes.Interface, namespace string, d *dispatcher.Dispatcher) error {
	factory := informers.NewSharedInformerFactoryWithOptions(kube, 0, informers.WithNamespace(namespace))
	pods := factory.Core().V1().Pods().Informer()
	changed := func(obj any) {
		if pod, ok := obj.(*corev1.Pod); ok {
			d.PodChanged(ctx, pod)
		}
	}
	_, err := pods.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    changed,
		UpdateFunc: func(_, obj any) { changed(obj) },
		DeleteFunc: func(obj any) {
			if tombstone, ok := obj.(cache.DeletedFinalStateUnknown); ok {
				obj = tombstone.Obj
			}
			if pod, ok := obj.(*corev1.Pod); ok {
				d.PodGone(pod.UID)
			}
		},
	})
	if err != nil {
		return err
	}
	factory.Start(ctx.Done())
	for _, synced := range factory.WaitForCacheSync(ctx.Done()) {
		if !synced {
			return fmt.Errorf("watching pods in %s: %w", namespace, ctx.Err())
		}
	}
	return nil
}

const defaultMaxHandoffs = 8

// sandboxClientAndAuth builds the sandbox client and the pod-facing
// authenticator together, since they share an identity model. On the add-on
// path the client's nonce registry is also what authenticates the endpoint.
func sandboxClientAndAuth(kube kubernetes.Interface, dyn dynamic.Interface, cfg *config.Config, shape sandbox.Shape) (*sandbox.Client, httpapi.Authenticator, error) {
	opts := sandbox.Options{Shape: shape, Namespace: cfg.K8SNamespace, EnvironmentID: cfg.EnvironmentID, TemplateName: cfg.K8SSandboxTemplateName, WarmPoolName: cfg.K8SWarmPoolName, WorkspaceClaim: workspaceClaim(cfg)}
	if !shape.Direct() {
		reviewer, err := httpapi.NewKubeTokenReviewer(kube, cfg.K8SNamespace, cfg.K8SSandboxServiceAccount, cfg.TokenAudience)
		if err != nil {
			return nil, nil, err
		}
		sandboxes, err := sandbox.New(dyn, kube, opts)
		if err != nil {
			return nil, nil, err
		}
		return sandboxes, httpapi.TokenAuth{Reviewer: reviewer}, nil
	}
	var sandboxes *sandbox.Client
	registry := nonce.New(nonce.PodLocatorFunc(func(ctx context.Context, name string) (types.UID, string, bool, error) {
		return sandboxes.SandboxPod(ctx, name)
	}))
	opts.Nonces = registry
	sandboxes, err := sandbox.New(dyn, kube, opts)
	if err != nil {
		return nil, nil, err
	}
	return sandboxes, httpapi.NonceAuth{Registry: registry}, nil
}

// handoffLimit sizes the sessions provisioned at once. The add-on path has no
// warm pool, so it takes the default; the upstream path reads the pool.
func handoffLimit(ctx context.Context, shape sandbox.Shape, sandboxes *sandbox.Client, logger *slog.Logger) int {
	if shape.Direct() {
		return defaultMaxHandoffs
	}
	replicas, err := sandboxes.WarmPoolReplicas(ctx)
	if err != nil || replicas <= 0 {
		logger.Warn("could not size handoffs from the warm pool; using the default", "default", defaultMaxHandoffs, "error", err)
		return defaultMaxHandoffs
	}
	return replicas
}

// preflight checks the cluster and the environment key from the operator's
// machine, prints one line per check, and fails when any check does.
func preflight(ctx context.Context, cfg *config.Config, shape sandbox.Shape) error {
	warnHTTP2Debug(newLogger(cfg.LogFormat), os.Getenv("GODEBUG"))
	kube, dyn, err := clusterClients()
	if err != nil {
		return err
	}
	checker := &checks.Checker{
		Kube: kube,
		Dyn:  dyn,
		API:  cmaapi.New(cfg.BaseURL, cfg.EnvironmentID, cfg.EnvironmentKey, cfg.APIRequestTimeout),
		Options: checks.Options{
			Namespace:                cfg.K8SNamespace,
			Shape:                    shape,
			SandboxTemplateName:      cfg.K8SSandboxTemplateName,
			WarmPoolName:             cfg.K8SWarmPoolName,
			DispatcherServiceAccount: cfg.K8SDispatcherServiceAccount,
			DeploymentName:           cfg.K8SDeploymentName,
			TokenAudience:            cfg.TokenAudience,
			ServiceHost:              cfg.ServiceHost,
			SkipNetworkProbe:         cfg.SkipNetworkProbe,
			WorkspaceClaim:           workspaceClaim(cfg),
		},
	}
	if failures := checker.Run(ctx, os.Stdout); failures > 0 {
		return fmt.Errorf("%d check(s) failed", failures)
	}
	return nil
}

// workspaceClaim is the claim the add-on path gives each session's
// /workspace, or nil when the size is zero.
func workspaceClaim(cfg *config.Config) *sandbox.WorkspaceClaim {
	if cfg.WorkspaceClaimSize.IsZero() {
		return nil
	}
	return &sandbox.WorkspaceClaim{Size: cfg.WorkspaceClaimSize, StorageClass: cfg.WorkspaceStorageClass}
}

// clusterClients uses the kubeconfig the operator has (KUBECONFIG or
// ~/.kube/config) and falls back to the pod's in-cluster credentials.
func clusterClients() (kubernetes.Interface, dynamic.Interface, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	restConfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		restConfig, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, nil, fmt.Errorf("connecting to the cluster: %w", err)
	}
	kube, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, nil, err
	}
	dyn, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, nil, err
	}
	return kube, dyn, nil
}

func newLogger(format string) *slog.Logger {
	if format == "json" {
		return slog.New(slog.NewJSONHandler(os.Stderr, nil))
	}
	return slog.New(slog.NewTextHandler(os.Stderr, nil))
}

func workerID() string {
	var b [4]byte
	_, _ = rand.Read(b[:])
	host, _ := os.Hostname()
	if host == "" {
		host = "cma-dispatcher"
	}
	return host + "-" + hex.EncodeToString(b[:])
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "cma-dispatcher:", err)
	os.Exit(1)
}
