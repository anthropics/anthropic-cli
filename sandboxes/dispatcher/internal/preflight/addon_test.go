package preflight

import (
	"errors"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/sandbox"
)

// noCredential starts the failure line of the add-on check on the template's
// containers and volumes.
const noCredential = "FAIL  sandbox template carries no kubelet credential and limits every container: "

func TestPreflightDirectFailures(t *testing.T) {
	for name, tc := range map[string]struct {
		mutate func(*directCluster)
		want   string
	}{
		"fetch has no cpu limit": {func(c *directCluster) {
			c.fetchLimits = map[string]any{"limits": map[string]any{"memory": "64Mi"}}
		}, noCredential + `container "fetch" sets no cpu or memory limit, which the add-on's policy requires of every container`},
		"worker has no memory limit": {func(c *directCluster) {
			c.workerLimits = map[string]any{"limits": map[string]any{"cpu": "500m"}}
		}, noCredential + `container "worker" sets no cpu or memory limit, which the add-on's policy requires of every container`},
		"fetch projects a token": {func(c *directCluster) { c.fetchProjectsToken = true }, noCredential + `volume "tok" projects a service-account token, ClusterTrustBundle, or PodCertificate, which the add-on's policy forbids in a sandbox pod`},
		"hostPath volume": {func(c *directCluster) {
			c.extraVolumes = []any{map[string]any{"name": "host", "hostPath": map[string]any{"path": "/etc"}}}
		}, noCredential + `volume "host" uses a hostPath, which the restricted pod security standard forbids and this check cannot confine`},
		"volume source the scans cannot inspect": {func(c *directCluster) {
			c.extraVolumes = []any{map[string]any{"name": "opaque", "csi": map[string]any{"driver": "secrets-store.csi.k8s.io"}}}
		}, noCredential + `volume "opaque" has a source this check cannot compare; use an emptyDir, configMap, secret or persistentVolumeClaim volume, or a projected volume with a configMap, secret or service-account-token source`},
		"projected volume with only a downwardAPI source": {func(c *directCluster) {
			c.extraVolumes = []any{map[string]any{"name": "meta", "projected": map[string]any{"sources": []any{map[string]any{"downwardAPI": map[string]any{"items": []any{map[string]any{"path": "labels", "fieldRef": map[string]any{"fieldPath": "metadata.labels"}}}}}}}}}
		}, noCredential + `volume "meta" has a source this check cannot compare; use an emptyDir, configMap, secret or persistentVolumeClaim volume, or a projected volume with a configMap, secret or service-account-token source`},
		"volume projects a ClusterTrustBundle": {func(c *directCluster) {
			c.extraVolumes = []any{map[string]any{"name": "bundle", "projected": map[string]any{"sources": []any{map[string]any{"clusterTrustBundle": map[string]any{"name": "root", "path": "ca.pem"}}}}}}
		}, noCredential + `volume "bundle" projects a service-account token, ClusterTrustBundle, or PodCertificate, which the add-on's policy forbids in a sandbox pod`},
		"volume projects a PodCertificate": {func(c *directCluster) {
			c.extraVolumes = []any{map[string]any{"name": "cert", "projected": map[string]any{"sources": []any{map[string]any{"podCertificate": map[string]any{"signerName": "example.com/pod", "keyType": "ED25519", "credentialBundlePath": "bundle.pem"}}}}}}
		}, noCredential + `volume "cert" projects a service-account token, ClusterTrustBundle, or PodCertificate, which the add-on's policy forbids in a sandbox pod`},
		"worker carries the environment key": {func(c *directCluster) {
			c.workerEnv = []any{map[string]any{"name": "ANTHROPIC_ENVIRONMENT_KEY", "value": "k"}}
		}, noCredential + `container "worker" carries ANTHROPIC_ENVIRONMENT_KEY; the environment key must never reach a sandbox pod`},
		"worker mounts a service-account token secret": {func(c *directCluster) {
			c.extraVolumes = []any{map[string]any{"name": "legacy", "secret": map[string]any{"secretName": "legacy-token"}}}
			c.workerMounts = []any{map[string]any{"name": "legacy", "mountPath": "/l"}}
		}, noCredential + `container "worker" reads Secret "legacy-token" of type kubernetes.io/service-account-token, a long-lived token for the Kubernetes API`},
		"worker imports an unreadable configmap": {func(c *directCluster) {
			c.workerEnvFrom = []any{map[string]any{"configMapRef": map[string]any{"name": "absent-config"}}}
		}, noCredential + `container "worker" imports ConfigMap "absent-config" through envFrom, and the check could not read it to show it free of the environment key: configmaps "absent-config" not found`},
		"fetch carries the environment key": {func(c *directCluster) {
			c.fetchEnv = []any{map[string]any{"name": "ANTHROPIC_ENVIRONMENT_KEY", "value": "k"}}
		}, noCredential + `container "fetch" carries ANTHROPIC_ENVIRONMENT_KEY; the environment key must never reach a sandbox pod`},
		"fetch mounts a service-account token secret": {func(c *directCluster) {
			c.extraVolumes = []any{map[string]any{"name": "legacy", "secret": map[string]any{"secretName": "legacy-token"}}}
			c.fetchMounts = []any{map[string]any{"name": "legacy", "mountPath": "/l"}}
		}, noCredential + `container "fetch" reads Secret "legacy-token" of type kubernetes.io/service-account-token, a long-lived token for the Kubernetes API`},
		"shared processes": {func(c *directCluster) { c.shareProcesses = true },
			noCredential + `the SandboxTemplate sets shareProcessNamespace, which lets every container read "fetch"'s nonce`},
		"no worker": {func(c *directCluster) { c.workerName = "app" },
			noCredential + `the SandboxTemplate must have containers named "fetch" and "worker"`},
		"no fetch": {func(c *directCluster) { c.noFetch = true },
			noCredential + `the SandboxTemplate must have containers named "fetch" and "worker"`},
		"automount on": {func(c *directCluster) { c.automountToken = true },
			noCredential + "the SandboxTemplate does not set automountServiceAccountToken: false, which the add-on's policy requires"},
		"served schema drops a spec field": {func(c *directCluster) { c.dropsSpecField = "shutdownPolicy" },
			"FAIL  the add-on admits a sandbox built from the template: the served Sandbox schema dropped spec.shutdownPolicy; check that -api-version v1alpha1 matches the cluster"},
		"rbac lacks create secrets": {func(c *directCluster) { c.denyResource = "secrets" },
			"FAIL  dispatcher service account holds its role: system:serviceaccount:ns:cma-dispatcher may not create secrets"},
		"automount unset":             {func(c *directCluster) { c.automountToken = nil }, noCredential + "the SandboxTemplate does not set automountServiceAccountToken: false, which the add-on's policy requires"},
		"sandbox api missing":         {func(c *directCluster) { c.missingResource = "sandboxes" }, "FAIL  agent sandbox API is served: agents.x-k8s.io/v1alpha1 does not serve sandboxes"},
		"add-on rejects the sandbox":  {func(c *directCluster) { c.rejectSandbox = true }, "FAIL  the add-on admits a sandbox built from the template: a dry-run sandbox create was refused, so the add-on would reject a real one: Security Violation: All containers must define resource limits"},
		"rbac lacks create sandboxes": {func(c *directCluster) { c.denyVerb = "create" }, "FAIL  dispatcher service account holds its role: system:serviceaccount:ns:cma-dispatcher may not create sandboxes.agents.x-k8s.io"},
		"rbac checks the configured account": {func(c *directCluster) {
			c.deploymentSA, c.dispatcherSA, c.denyVerb = "cma-dispatcher-v2", "cma-dispatcher-v2", "create"
		}, "FAIL  dispatcher service account holds its role: system:serviceaccount:ns:cma-dispatcher-v2 may not create sandboxes.agents.x-k8s.io"},
		"rbac needs no tokenreview": {func(c *directCluster) { c.denyResource = "tokenreviews" }, "ok    dispatcher service account holds its role"},
		"init container has no limits": {func(c *directCluster) {
			c.initContainers = []any{map[string]any{"name": "setup"}}
		}, noCredential + `container "setup" sets no cpu or memory limit, which the add-on's policy requires of every container`},
		"init container carries the environment key": {func(c *directCluster) {
			c.initContainers = []any{map[string]any{"name": "setup", "resources": map[string]any{"limits": map[string]any{"cpu": "10m", "memory": "64Mi"}}, "env": []any{map[string]any{"name": "ANTHROPIC_ENVIRONMENT_KEY", "value": "k"}}}}
		}, noCredential + `container "setup" carries ANTHROPIC_ENVIRONMENT_KEY; the environment key must never reach a sandbox pod`},
		"init container mounts a service-account token secret": {func(c *directCluster) {
			c.extraVolumes = []any{map[string]any{"name": "legacy", "secret": map[string]any{"secretName": "legacy-token"}}}
			c.initContainers = []any{map[string]any{"name": "setup", "resources": map[string]any{"limits": map[string]any{"cpu": "10m", "memory": "64Mi"}}, "volumeMounts": []any{map[string]any{"name": "legacy", "mountPath": "/l"}}}}
		}, noCredential + `container "setup" reads Secret "legacy-token" of type kubernetes.io/service-account-token, a long-lived token for the Kubernetes API`},
		"worker knows the dispatcher": {func(c *directCluster) {
			c.workerEnv = []any{map[string]any{"name": "CMA_DISPATCHER_URL", "value": "https://d"}}
		}, noCredential + `container "worker" sets CMA_DISPATCHER_URL, which only "fetch" may carry`},
		"worker reads the token path through an env valueFrom": {func(c *directCluster) {
			c.workerEnv = []any{map[string]any{"name": "CMA_SA_TOKEN_FILE", "valueFrom": map[string]any{"configMapKeyRef": map[string]any{"name": "cma-dispatcher-ca", "key": "ca.crt"}}}}
		}, noCredential + `container "worker" sets CMA_SA_TOKEN_FILE, which only "fetch" may carry`},
		"init container knows the dispatcher ca": {func(c *directCluster) {
			c.initContainers = []any{map[string]any{"name": "setup", "resources": map[string]any{"limits": map[string]any{"cpu": "10m", "memory": "64Mi"}}, "env": []any{map[string]any{"name": "CMA_CA_FILE", "value": "/c"}}}}
		}, noCredential + `container "setup" sets CMA_CA_FILE, which only "fetch" may carry`},
		"api lacks sandboxtemplates":    {func(c *directCluster) { c.missingResource = "sandboxtemplates" }, "FAIL  agent sandbox API is served: extensions.agents.x-k8s.io/v1alpha1 does not serve sandboxtemplates"},
		"managed policy":                {func(c *directCluster) { c.networkPolicyMode = "" }, `FAIL  sandbox template turns the managed network policy off: the stored SandboxTemplate has networkPolicyManagement "", not Unmanaged: the template omits it, or the served schema predates the field and pruned it; either way the controller's own policy leaves 100.64.0.0/10 open`},
		"managed policy set explicitly": {func(c *directCluster) { c.networkPolicyMode = "Managed" }, `FAIL  sandbox template turns the managed network policy off: the stored SandboxTemplate has networkPolicyManagement "Managed", not Unmanaged: the template omits it, or the served schema predates the field and pruned it; either way the controller's own policy leaves 100.64.0.0/10 open`},
		"runc":                          {func(c *directCluster) { c.handler = "runc" }, `FAIL  sandbox runtime is gVisor or Kata: RuntimeClass gvisor uses handler "runc"; only runsc or gvisor (gVisor) or a Kata handler isolates sandboxes`},
		"psa":                           {func(c *directCluster) { c.psaLabel = "baseline" }, "FAIL  namespace enforces restricted pod security: namespace ns lacks the label pod-security.kubernetes.io/enforce=restricted"},
		"replicas":                      {func(c *directCluster) { c.replicas = 2 }, "FAIL  dispatcher runs as one replica with Recreate: deployment cma-dispatcher has 2 replicas; the dispatcher must run exactly one"},
		"rolling":                       {func(c *directCluster) { c.strategy = appsv1.RollingUpdateDeploymentStrategyType }, `FAIL  dispatcher runs as one replica with Recreate: deployment cma-dispatcher uses strategy "RollingUpdate"; two dispatchers would overlap during a rollout without Recreate`},
		"deployment runs as another account": {func(c *directCluster) { c.deploymentSA = "default" },
			`FAIL  dispatcher runs as one replica with Recreate: deployment cma-dispatcher runs as service account "default", but the RBAC check inspects "cma-dispatcher"; the checked identity must be the deployed one`},
		"deployment names no service account": {func(c *directCluster) { c.deploymentSA = "" },
			`FAIL  dispatcher runs as one replica with Recreate: deployment cma-dispatcher runs as service account "default", but the RBAC check inspects "cma-dispatcher"; the checked identity must be the deployed one`},
		"environment key rejected": {func(c *directCluster) { c.verifyErr = errors.New("the queue answered 401") },
			"FAIL  environment key reads the work queue: the queue answered 401"},
	} {
		c := healthyDirect()
		tc.mutate(&c)
		wantLine(t, name, runDirect(t, c), tc.want)
	}
}

// Only a claim-backed /workspace survives a suspend. A configured claim passes,
// and so does a template whose own workspace volume is a claim. With neither,
// the add-on path backs /workspace with an emptyDir, and the preflight must
// say so.
func TestPreflightDirectWorkspaceBacking(t *testing.T) {
	claim := func(class string) *sandbox.WorkspaceClaim {
		return &sandbox.WorkspaceClaim{Size: resource.MustParse("10Gi"), StorageClass: class}
	}
	for name, tc := range map[string]struct {
		mutate func(*directCluster)
		want   string
	}{
		"no claim and no workspace volume": {func(*directCluster) {}, "warn  workspace survives a suspend: warning: the stored template carries no workspace volume and no workspace claim size is configured, so the dispatcher backs /workspace with an emptyDir that a suspend discards"},
		"claim on the default class":       {func(c *directCluster) { c.workspaceClaim = claim("") }, "ok    workspace survives a suspend"},
		"template's own claim-backed volume": {func(c *directCluster) {
			c.extraVolumes = []any{map[string]any{"name": "workspace", "persistentVolumeClaim": map[string]any{"claimName": "ws"}}}
		}, "ok    workspace survives a suspend"},
		"claim on a missing class": {func(c *directCluster) { c.workspaceClaim = claim("missing") }, `FAIL  workspace survives a suspend: StorageClass missing: storageclasses.storage.k8s.io "missing" not found`},
		"claim on no class with no default": {func(c *directCluster) {
			c.workspaceClaim, c.storageClass = claim(""), nil
		}, "FAIL  workspace survives a suspend: the volume claim template names no storageClassName and the cluster has no default StorageClass"},
		"claim on a Retain class": {func(c *directCluster) {
			c.workspaceClaim, c.storageClass = claim(""), defaultStorageClass(corev1.PersistentVolumeReclaimRetain)
		}, "warn  workspace survives a suspend: warning: StorageClass standard reclaims with Retain, so deleted sessions leave volumes behind"},
		"template's own emptyDir volume": {func(c *directCluster) {
			c.extraVolumes = []any{map[string]any{"name": "workspace", "emptyDir": map[string]any{}}}
		}, "warn  workspace survives a suspend: warning: the template's workspace volume is not claim-backed, so a suspend discards /workspace"},
	} {
		c := healthyDirect()
		tc.mutate(&c)
		wantLine(t, name, runDirect(t, c), tc.want)
	}
}
