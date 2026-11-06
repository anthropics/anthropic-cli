package sandbox

import (
	"context"
	"errors"
	"time"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

// ErrNotBound is returned while a claim's status names no Sandbox yet.
var ErrNotBound = errors.New("sandbox claim is not bound to a sandbox yet")

// fieldManager names the dispatcher's writes in the objects' managed fields.
const fieldManager = "cma-dispatcher"

// SandboxClaim is the dispatcher's view of one SandboxClaim.
type SandboxClaim struct {
	Name        string
	SessionID   string
	SandboxName string
	LastActive  time.Time
	// HandoffRestarts is the worker's restart count recorded at the last
	// handoff; nil while no handoff has recorded one.
	HandoffRestarts *int32
	// HandoffPodUID is the pod the last recorded handoff went to; empty
	// while no handoff has recorded one.
	HandoffPodUID types.UID
}

// Sandbox is the dispatcher's view of one Sandbox.
type Sandbox struct {
	Name      string
	UID       types.UID
	PodName   string
	Suspended bool
	Ready     bool
}

// NonceSink records the per-sandbox secret the dispatcher writes into each
// direct Sandbox's nonce Secret, so the pod-facing endpoint can authenticate
// the pod by it. The direct (add-on) path alone uses it.
type NonceSink interface {
	// Mint returns a fresh nonce for sandboxName and registers it.
	Mint(sandboxName string) (string, error)
	// Known reports whether sandboxName currently has a registered nonce.
	Known(sandboxName string) bool
	// Clear forgets sandboxName's nonce once its Sandbox is gone.
	Clear(sandboxName string)
}

// Client reads and writes the Agent Sandbox objects of one namespace that
// belong to one environment.
type Client struct {
	dyn           dynamic.Interface
	kube          kubernetes.Interface
	shape         Shape
	namespace     string
	environmentID string
	templateName  string
	warmPoolName  string
	// nonces registers each direct Sandbox's nonce Secret value; nil on the
	// upstream path, which authenticates pods by projected token instead.
	nonces NonceSink
	// workspaceClaim backs each direct Sandbox's /workspace; nil means an
	// emptyDir that a suspend discards.
	workspaceClaim *WorkspaceClaim
}

// WorkspaceClaim is the per-session volume the direct path gives /workspace:
// a ReadWriteOnce claim of Size from StorageClass (the cluster default when
// empty), owned by the Sandbox so it survives a suspend and goes with a
// delete. The upstream path takes its claim from the SandboxTemplate.
type WorkspaceClaim struct {
	Size         resource.Quantity
	StorageClass string
}

// Options configures a Client. TemplateName names the SandboxTemplate and
// WarmPoolName the SandboxWarmPool a new claim checks out from. Nonces is
// required when Shape.Direct() and refused otherwise. WorkspaceClaim, direct
// path only, backs /workspace with a claim; nil backs it with an emptyDir.
type Options struct {
	Shape          Shape
	Namespace      string
	EnvironmentID  string
	TemplateName   string
	WarmPoolName   string
	Nonces         NonceSink
	WorkspaceClaim *WorkspaceClaim
}

func New(dyn dynamic.Interface, kube kubernetes.Interface, opts Options) (*Client, error) {
	if dyn == nil || kube == nil {
		panic("sandbox.New: nil Kubernetes client")
	}
	if opts.EnvironmentID == "" {
		return nil, errors.New("the sandbox client needs an environment id")
	}
	if opts.Shape.Direct() && opts.Nonces == nil {
		return nil, errors.New("the add-on path needs a nonce sink to authenticate sandbox pods")
	}
	return &Client{dyn: dyn, kube: kube, shape: opts.Shape, namespace: opts.Namespace, environmentID: opts.EnvironmentID, templateName: opts.TemplateName, warmPoolName: opts.WarmPoolName, nonces: opts.Nonces, workspaceClaim: opts.WorkspaceClaim}, nil
}

// WarmPoolReplicas reads the configured warm pool's desired size.
func (c *Client) WarmPoolReplicas(ctx context.Context) (int, error) {
	pool, err := c.dyn.Resource(c.shape.WarmPools()).Namespace(c.namespace).Get(ctx, c.warmPoolName, metav1.GetOptions{})
	if err != nil {
		return 0, err
	}
	replicas, _, err := unstructured.NestedInt64(pool.Object, "spec", "replicas")
	return int(replicas), err
}

// managedSelector selects the claims of this dispatcher's environment.
func (c *Client) managedSelector() string {
	return LabelManagedBy + "=" + ManagedByValue + "," + LabelEnvironmentID + "=" + c.environmentID
}

func (c *Client) claims() dynamic.ResourceInterface {
	return c.dyn.Resource(c.shape.Claims()).Namespace(c.namespace)
}

func (c *Client) sandboxes() dynamic.ResourceInterface {
	return c.dyn.Resource(c.shape.Sandboxes()).Namespace(c.namespace)
}
