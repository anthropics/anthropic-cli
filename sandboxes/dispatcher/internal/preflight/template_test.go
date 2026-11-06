package preflight

import (
	"errors"
	"strings"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// keepsToken starts the failure line of the check on the fetch container's token.
const keepsToken = "FAIL  sandbox template keeps the token in the fetch container: "

func TestPreflightFailures(t *testing.T) {
	for name, tc := range map[string]struct {
		mutate func(*cluster)
		want   string
	}{
		"runc":     {func(c *cluster) { c.handler = "runc" }, `FAIL  sandbox runtime is gVisor or Kata: RuntimeClass gvisor uses handler "runc"; only runsc or gvisor (gVisor) or a Kata handler isolates sandboxes`},
		"psa":      {func(c *cluster) { c.psaLabel = "baseline" }, "FAIL  namespace enforces restricted pod security: namespace ns lacks the label pod-security.kubernetes.io/enforce=restricted"},
		"replicas": {func(c *cluster) { c.replicas = 2 }, "FAIL  dispatcher runs as one replica with Recreate: deployment cma-dispatcher has 2 replicas; the dispatcher must run exactly one"},
		"rolling":  {func(c *cluster) { c.strategy = appsv1.RollingUpdateDeploymentStrategyType }, `FAIL  dispatcher runs as one replica with Recreate: deployment cma-dispatcher uses strategy "RollingUpdate"; two dispatchers would overlap during a rollout without Recreate`},
		"retain":   {func(c *cluster) { c.reclaim = corev1.PersistentVolumeReclaimRetain }, "warn  workspace storage class exists: warning: StorageClass fast reclaims with Retain, so deleted sessions leave volumes behind"},
		"no class": {func(c *cluster) { c.storageClass = "missing" }, `FAIL  workspace storage class exists: StorageClass missing: storageclasses.storage.k8s.io "missing" not found`},
		"retain then no class": {func(c *cluster) {
			c.reclaim = corev1.PersistentVolumeReclaimRetain
			c.secondStorageClass = "missing"
		}, `FAIL  workspace storage class exists: StorageClass missing: storageclasses.storage.k8s.io "missing" not found`},
		"pool template":                 {func(c *cluster) { c.poolTemplate = "other" }, `FAIL  warm pool makes its pods from the checked template: SandboxWarmPool cma-sandbox makes its pods from SandboxTemplate "other", not "cma-sandbox"; point -warm-pool and -sandbox-template at the same template`},
		"no warm pool":                  {func(c *cluster) { c.noWarmPool = true }, `FAIL  warm pool makes its pods from the checked template: SandboxWarmPool cma-sandbox: sandboxwarmpools.extensions.agents.x-k8s.io "cma-sandbox" not found`},
		"no template":                   {func(c *cluster) { c.noTemplate = true }, `FAIL  sandbox runtime is gVisor or Kata: SandboxTemplate cma-sandbox: sandboxtemplates.extensions.agents.x-k8s.io "cma-sandbox" not found`},
		"pod spec does not convert":     {func(c *cluster) { c.shareProcesses = "yes" }, keepsToken + "reading the SandboxTemplate's pod spec: unrecognized type: bool"},
		"no runtime class named":        {func(c *cluster) { c.runtimeClassName = "" }, "FAIL  sandbox runtime is gVisor or Kata: the SandboxTemplate sets no runtimeClassName; sandboxes would run on the host kernel"},
		"runtime class missing":         {func(c *cluster) { c.runtimeClassName = "kata" }, `FAIL  sandbox runtime is gVisor or Kata: RuntimeClass kata: runtimeclasses.node.k8s.io "kata" not found`},
		"managed policy":                {func(c *cluster) { c.networkPolicyMode = "" }, `FAIL  sandbox template turns the managed network policy off: the stored SandboxTemplate has networkPolicyManagement "", not Unmanaged: the template omits it, or the served schema predates the field and pruned it; either way the controller's own policy leaves 100.64.0.0/10 open`},
		"managed policy set explicitly": {func(c *cluster) { c.networkPolicyMode = "Managed" }, `FAIL  sandbox template turns the managed network policy off: the stored SandboxTemplate has networkPolicyManagement "Managed", not Unmanaged: the template omits it, or the served schema predates the field and pruned it; either way the controller's own policy leaves 100.64.0.0/10 open`},
		"worker holds token": {func(c *cluster) {
			c.worker["volumeMounts"] = append(c.worker["volumeMounts"].([]any), map[string]any{"name": "dispatcher-token", "mountPath": "/t"})
		}, keepsToken + `container "worker" mounts volume "dispatcher-token", which reads what only "fetch" may use`},
		"worker holds ca": {func(c *cluster) {
			c.worker["volumeMounts"] = append(c.worker["volumeMounts"].([]any), map[string]any{"name": "dispatcher-ca", "mountPath": "/c"})
		}, keepsToken + `container "worker" mounts volume "dispatcher-ca", which reads what only "fetch" may use`},
		"worker knows dispatcher": {func(c *cluster) {
			c.worker["env"] = []any{map[string]any{"name": "CMA_SA_TOKEN_FILE", "value": "/t"}}
		}, keepsToken + `container "worker" sets CMA_SA_TOKEN_FILE, which only "fetch" may carry`},
		"worker knows the dispatcher ca": {func(c *cluster) {
			c.worker["env"] = []any{map[string]any{"name": "CMA_CA_FILE", "value": "/c"}}
		}, keepsToken + `container "worker" sets CMA_CA_FILE, which only "fetch" may carry`},
		"init container holds token": {func(c *cluster) {
			c.initContainers = []any{map[string]any{"name": "setup", "volumeMounts": []any{map[string]any{"name": "dispatcher-token", "mountPath": "/t"}}}}
		}, keepsToken + `container "setup" mounts volume "dispatcher-token", which reads what only "fetch" may use`},
		"sidecar knows dispatcher": {func(c *cluster) {
			c.sidecars = []any{map[string]any{"name": "proxy", "env": []any{map[string]any{"name": "CMA_DISPATCHER_URL", "value": "https://d"}}}}
		}, keepsToken + `container "proxy" sets CMA_DISPATCHER_URL, which only "fetch" may carry`},
		"shared processes":    {func(c *cluster) { c.shareProcesses = true }, keepsToken + `the SandboxTemplate sets shareProcessNamespace, which lets every container read "fetch"'s token`},
		"default token unset": {func(c *cluster) { c.automountToken = nil }, keepsToken + "the SandboxTemplate does not set automountServiceAccountToken: false, so every container would get a token for the Kubernetes API"},
		"default token on":    {func(c *cluster) { c.automountToken = true }, keepsToken + "the SandboxTemplate does not set automountServiceAccountToken: false, so every container would get a token for the Kubernetes API"},
		"a hostPath volume": {func(c *cluster) {
			c.extraVolumes = []any{map[string]any{"name": "host", "hostPath": map[string]any{"path": "/var/lib/kubelet/pods"}}}
			c.sidecars = []any{map[string]any{"name": "proxy", "volumeMounts": []any{map[string]any{"name": "host", "mountPath": "/h"}}}}
		}, keepsToken + `volume "host" uses a hostPath, which the restricted pod security standard forbids and this check cannot confine`},
		"sidecar mounts an uncomparable volume": {func(c *cluster) {
			c.extraVolumes = []any{map[string]any{"name": "scratch", "csi": map[string]any{"driver": "some.driver"}}}
			c.sidecars = []any{map[string]any{"name": "proxy", "volumeMounts": []any{map[string]any{"name": "scratch", "mountPath": "/s"}}}}
		}, keepsToken + `volume "scratch" has a source this check cannot compare; use an emptyDir, configMap, secret or persistentVolumeClaim volume, or a projected volume with a configMap, secret or service-account-token source`},
		"sidecar mounts a projected volume with only a downwardAPI source": {func(c *cluster) {
			c.extraVolumes = []any{map[string]any{"name": "meta", "projected": map[string]any{"sources": []any{map[string]any{"downwardAPI": map[string]any{"items": []any{map[string]any{"path": "labels", "fieldRef": map[string]any{"fieldPath": "metadata.labels"}}}}}}}}}
			c.sidecars = []any{map[string]any{"name": "proxy", "volumeMounts": []any{map[string]any{"name": "meta", "mountPath": "/m"}}}}
		}, keepsToken + `volume "meta" has a source this check cannot compare; use an emptyDir, configMap, secret or persistentVolumeClaim volume, or a projected volume with a configMap, secret or service-account-token source`},
		"sidecar reads the ca through an env valueFrom": {func(c *cluster) {
			c.sidecars = []any{map[string]any{"name": "proxy", "env": []any{map[string]any{"name": "CA", "valueFrom": map[string]any{"configMapKeyRef": map[string]any{"name": "ca", "key": "ca.crt"}}}}}}
		}, keepsToken + `container "proxy" reads configmap ca through an env valueFrom, which only "fetch" may use`},
		"sidecar reads the ca through envFrom": {func(c *cluster) {
			c.sidecars = []any{map[string]any{"name": "proxy", "envFrom": []any{map[string]any{"configMapRef": map[string]any{"name": "ca"}}}}}
		}, keepsToken + `container "proxy" reads configmap ca through envFrom, which only "fetch" may use`},
		"sidecar mounts the ca under another name": {func(c *cluster) {
			c.extraVolumes = []any{map[string]any{"name": "ca-copy", "configMap": map[string]any{"name": "ca"}}}
			c.sidecars = []any{map[string]any{"name": "proxy", "volumeMounts": []any{map[string]any{"name": "ca-copy", "mountPath": "/c"}}}}
		}, keepsToken + `container "proxy" mounts volume "ca-copy", which reads what only "fetch" may use`},
		"no fetch": {func(c *cluster) { c.noFetch = true },
			keepsToken + `the SandboxTemplate must have containers named "fetch" and "worker"`},
		"worker is named fetch": {func(c *cluster) { c.worker["name"] = "fetch" },
			keepsToken + `the SandboxTemplate must have containers named "fetch" and "worker"`},
		"environment key rejected": {func(c *cluster) { c.verifyErr = errors.New("the queue answered 401") }, "FAIL  environment key reads the work queue: the queue answered 401"},
		"rbac":                     {func(c *cluster) { c.denyVerb = "delete" }, "FAIL  dispatcher service account holds its role: system:serviceaccount:ns:cma-dispatcher may not delete sandboxclaims.extensions.agents.x-k8s.io"},
		"rbac checks the configured account": {func(c *cluster) {
			c.deploymentSA, c.dispatcherSA, c.denyVerb = "cma-dispatcher-v2", "cma-dispatcher-v2", "delete"
		}, "FAIL  dispatcher service account holds its role: system:serviceaccount:ns:cma-dispatcher-v2 may not delete sandboxclaims.extensions.agents.x-k8s.io"},
		"reviews forbidden":                  {func(c *cluster) { c.reviewsForbidden = true }, "FAIL  dispatcher service account holds its role: your credentials may not create SubjectAccessReviews, so the dispatcher's access cannot be checked: subjectaccessreviews.authorization.k8s.io is forbidden: denied"},
		"no deployment":                      {func(c *cluster) { c.noDeployment = true }, `FAIL  dispatcher runs as one replica with Recreate: deployment cma-dispatcher: deployments.apps "cma-dispatcher" not found`},
		"deployment runs as another account": {func(c *cluster) { c.deploymentSA = "default" }, `FAIL  dispatcher runs as one replica with Recreate: deployment cma-dispatcher runs as service account "default", but the RBAC check inspects "cma-dispatcher"; the checked identity must be the deployed one`},
		"deployment names no service account": {func(c *cluster) { c.deploymentSA = "" },
			`FAIL  dispatcher runs as one replica with Recreate: deployment cma-dispatcher runs as service account "default", but the RBAC check inspects "cma-dispatcher"; the checked identity must be the deployed one`},
		"token audience unset": {func(c *cluster) { c.tokenAudience = "" }, keepsToken + `volume "dispatcher-token" projects a service-account token for audience "", not "cma-dispatcher"; an unset audience is a token for the Kubernetes API itself`},
		"worker carries the environment key": {func(c *cluster) {
			c.worker["env"] = []any{map[string]any{"name": "ANTHROPIC_ENVIRONMENT_KEY", "value": "k"}}
		}, keepsToken + `container "worker" carries ANTHROPIC_ENVIRONMENT_KEY; the environment key must never reach a sandbox pod`},
		"worker reads the environment key under another name": {func(c *cluster) {
			c.worker["env"] = []any{map[string]any{"name": "KEY", "valueFrom": map[string]any{"secretKeyRef": map[string]any{"name": "env-key", "key": "ANTHROPIC_ENVIRONMENT_KEY"}}}}
		}, keepsToken + `container "worker" reads the ANTHROPIC_ENVIRONMENT_KEY key into KEY; the environment key must never reach a sandbox pod`},
		"worker imports the environment-key secret": {func(c *cluster) {
			c.worker["envFrom"] = []any{map[string]any{"secretRef": map[string]any{"name": "env-key"}}}
		}, keepsToken + `container "worker" imports Secret env-key, which holds ANTHROPIC_ENVIRONMENT_KEY,; the environment key must never reach a sandbox pod`},
		"worker imports an unreadable secret": {func(c *cluster) {
			c.worker["envFrom"] = []any{map[string]any{"secretRef": map[string]any{"name": "absent"}}}
		}, keepsToken + `container "worker" imports Secret "absent" through envFrom, and the check could not read it to show it free of the environment key: secrets "absent" not found`},
		"worker mounts a service-account token secret": {func(c *cluster) {
			c.extraVolumes = []any{map[string]any{"name": "legacy", "secret": map[string]any{"secretName": "legacy-token"}}}
			c.worker["volumeMounts"] = append(c.worker["volumeMounts"].([]any), map[string]any{"name": "legacy", "mountPath": "/l"})
		}, keepsToken + `container "worker" reads Secret "legacy-token" of type kubernetes.io/service-account-token, a long-lived token for the Kubernetes API`},
		"worker projects a service-account token secret": {func(c *cluster) {
			c.extraVolumes = []any{map[string]any{"name": "legacy", "projected": map[string]any{"sources": []any{map[string]any{"secret": map[string]any{"name": "legacy-token"}}}}}}
			c.worker["volumeMounts"] = append(c.worker["volumeMounts"].([]any), map[string]any{"name": "legacy", "mountPath": "/l"})
		}, keepsToken + `container "worker" reads Secret "legacy-token" of type kubernetes.io/service-account-token, a long-lived token for the Kubernetes API`},
		"worker imports a service-account token secret": {func(c *cluster) {
			c.worker["envFrom"] = []any{map[string]any{"secretRef": map[string]any{"name": "legacy-token"}}}
		}, keepsToken + `container "worker" reads Secret "legacy-token" of type kubernetes.io/service-account-token, a long-lived token for the Kubernetes API`},
		"worker reads a service-account token secret through valueFrom": {func(c *cluster) {
			c.worker["env"] = []any{map[string]any{"name": "T", "valueFrom": map[string]any{"secretKeyRef": map[string]any{"name": "legacy-token", "key": "token"}}}}
		}, keepsToken + `container "worker" reads Secret "legacy-token" of type kubernetes.io/service-account-token, a long-lived token for the Kubernetes API`},
		"worker mounts an unreadable secret": {func(c *cluster) {
			c.extraVolumes = []any{map[string]any{"name": "mystery", "secret": map[string]any{"secretName": "absent"}}}
			c.worker["volumeMounts"] = append(c.worker["volumeMounts"].([]any), map[string]any{"name": "mystery", "mountPath": "/m"})
		}, keepsToken + `container "worker" mounts Secret "absent", and the check could not read it to show it free of the environment key: secrets "absent" not found`},
		"worker mounts the environment-key secret": {func(c *cluster) {
			c.extraVolumes = []any{map[string]any{"name": "creds", "secret": map[string]any{"secretName": "env-key"}}}
			c.worker["volumeMounts"] = append(c.worker["volumeMounts"].([]any), map[string]any{"name": "creds", "mountPath": "/creds"})
		}, keepsToken + `container "worker" mounts Secret env-key, which holds ANTHROPIC_ENVIRONMENT_KEY,; the environment key must never reach a sandbox pod`},
		"worker projects the environment-key secret": {func(c *cluster) {
			c.extraVolumes = []any{map[string]any{"name": "creds", "projected": map[string]any{"sources": []any{map[string]any{"secret": map[string]any{"name": "env-key"}}}}}}
			c.worker["volumeMounts"] = append(c.worker["volumeMounts"].([]any), map[string]any{"name": "creds", "mountPath": "/creds"})
		}, keepsToken + `container "worker" mounts Secret env-key, which holds ANTHROPIC_ENVIRONMENT_KEY,; the environment key must never reach a sandbox pod`},
		"worker mounts a configmap holding the environment key": {func(c *cluster) {
			c.extraVolumes = []any{map[string]any{"name": "conf", "configMap": map[string]any{"name": "key-config"}}}
			c.worker["volumeMounts"] = append(c.worker["volumeMounts"].([]any), map[string]any{"name": "conf", "mountPath": "/conf"})
		}, keepsToken + `container "worker" mounts ConfigMap key-config, which holds ANTHROPIC_ENVIRONMENT_KEY,; the environment key must never reach a sandbox pod`},
		"worker projects a configmap holding the environment key": {func(c *cluster) {
			c.extraVolumes = []any{map[string]any{"name": "conf", "projected": map[string]any{"sources": []any{map[string]any{"configMap": map[string]any{"name": "key-config"}}}}}}
			c.worker["volumeMounts"] = append(c.worker["volumeMounts"].([]any), map[string]any{"name": "conf", "mountPath": "/conf"})
		}, keepsToken + `container "worker" mounts ConfigMap key-config, which holds ANTHROPIC_ENVIRONMENT_KEY,; the environment key must never reach a sandbox pod`},
		"worker mounts an unreadable configmap": {func(c *cluster) {
			c.extraVolumes = []any{map[string]any{"name": "conf", "configMap": map[string]any{"name": "absent-config"}}}
			c.worker["volumeMounts"] = append(c.worker["volumeMounts"].([]any), map[string]any{"name": "conf", "mountPath": "/conf"})
		}, keepsToken + `container "worker" mounts ConfigMap "absent-config", and the check could not read it to show it free of the environment key: configmaps "absent-config" not found`},
		"fetch mounts a service-account token secret": {func(c *cluster) {
			c.extraVolumes = []any{map[string]any{"name": "legacy", "secret": map[string]any{"secretName": "legacy-token"}}}
			c.fetchMounts = []any{map[string]any{"name": "legacy", "mountPath": "/l"}}
		}, keepsToken + `container "fetch" reads Secret "legacy-token" of type kubernetes.io/service-account-token, a long-lived token for the Kubernetes API`},
		"fetch carries the environment key": {func(c *cluster) {
			c.fetchEnv = []any{map[string]any{"name": "ANTHROPIC_ENVIRONMENT_KEY", "value": "k"}}
		}, keepsToken + `container "fetch" carries ANTHROPIC_ENVIRONMENT_KEY; the environment key must never reach a sandbox pod`},
		"api not served": {func(c *cluster) { c.unservedGroupVersion = "agents.x-k8s.io/v1beta1" }, "FAIL  agent sandbox API is served: agents.x-k8s.io/v1beta1 is not served: the server could not find the requested resource, GroupVersion \"agents.x-k8s.io/v1beta1\" not found"},
		"worker holds the token instead of fetch": {func(c *cluster) {
			c.fetchWithoutToken = true
			c.worker["volumeMounts"] = append(c.worker["volumeMounts"].([]any), map[string]any{"name": "dispatcher-token", "mountPath": "/t"})
		}, keepsToken + `container "fetch" mounts no projected service-account token`},
		"init container carries the environment key": {func(c *cluster) {
			c.initContainers = []any{map[string]any{"name": "setup", "env": []any{map[string]any{"name": "ANTHROPIC_ENVIRONMENT_KEY", "value": "k"}}}}
		}, keepsToken + `container "setup" carries ANTHROPIC_ENVIRONMENT_KEY; the environment key must never reach a sandbox pod`},
		"init container mounts a service-account token secret": {func(c *cluster) {
			c.extraVolumes = []any{map[string]any{"name": "legacy", "secret": map[string]any{"secretName": "legacy-token"}}}
			c.initContainers = []any{map[string]any{"name": "setup", "volumeMounts": []any{map[string]any{"name": "legacy", "mountPath": "/l"}}}}
		}, keepsToken + `container "setup" reads Secret "legacy-token" of type kubernetes.io/service-account-token, a long-lived token for the Kubernetes API`},
		"worker reads the environment key from a configmap under another name": {func(c *cluster) {
			c.worker["env"] = []any{map[string]any{"name": "KEY", "valueFrom": map[string]any{"configMapKeyRef": map[string]any{"name": "key-config", "key": "ANTHROPIC_ENVIRONMENT_KEY"}}}}
		}, keepsToken + `container "worker" reads the ANTHROPIC_ENVIRONMENT_KEY key into KEY; the environment key must never reach a sandbox pod`},
		"worker imports a configmap holding the environment key": {func(c *cluster) {
			c.worker["envFrom"] = []any{map[string]any{"configMapRef": map[string]any{"name": "key-config"}}}
		}, keepsToken + `container "worker" imports ConfigMap key-config, which holds ANTHROPIC_ENVIRONMENT_KEY,; the environment key must never reach a sandbox pod`},
		"worker imports an unreadable configmap": {func(c *cluster) {
			c.worker["envFrom"] = []any{map[string]any{"configMapRef": map[string]any{"name": "absent-config"}}}
		}, keepsToken + `container "worker" imports ConfigMap "absent-config" through envFrom, and the check could not read it to show it free of the environment key: configmaps "absent-config" not found`},
		"worker mounts a configmap holding the environment key as binary data": {func(c *cluster) {
			c.extraVolumes = []any{map[string]any{"name": "conf", "configMap": map[string]any{"name": "binary-key-config"}}}
			c.worker["volumeMounts"] = append(c.worker["volumeMounts"].([]any), map[string]any{"name": "conf", "mountPath": "/conf"})
		}, keepsToken + `container "worker" mounts ConfigMap binary-key-config, which holds ANTHROPIC_ENVIRONMENT_KEY,; the environment key must never reach a sandbox pod`},
		"worker reads an unreadable secret through valueFrom": {func(c *cluster) {
			c.worker["env"] = []any{map[string]any{"name": "T", "valueFrom": map[string]any{"secretKeyRef": map[string]any{"name": "absent", "key": "token"}}}}
		}, keepsToken + `container "worker" reads Secret "absent", and the check could not read it to show it is not a service-account token: secrets "absent" not found`},
		"api lacks warm pools":             {func(c *cluster) { c.missingResource = "sandboxwarmpools" }, "FAIL  agent sandbox API is served: extensions.agents.x-k8s.io/v1beta1 does not serve sandboxwarmpools"},
		"rbac tokenreviews":                {func(c *cluster) { c.denyResource = "tokenreviews" }, "FAIL  dispatcher service account holds its role: system:serviceaccount:ns:cma-dispatcher may not create tokenreviews.authentication.k8s.io"},
		"no default class":                 {func(c *cluster) { c.storageClass = "" }, "FAIL  workspace storage class exists: the volume claim template names no storageClassName and the cluster has no default StorageClass"},
		"no claim templates":               {func(c *cluster) { c.noClaimTemplates = true }, "FAIL  workspace storage class exists: the SandboxTemplate has no volumeClaimTemplates; /workspace would not survive a suspend"},
		"token audience is the api server": {func(c *cluster) { c.tokenAudience = "https://kubernetes.default.svc" }, keepsToken + `volume "dispatcher-token" projects a service-account token for audience "https://kubernetes.default.svc", not "cma-dispatcher"; an unset audience is a token for the Kubernetes API itself`},
		"api lacks sandboxclaims":          {func(c *cluster) { c.missingResource = "sandboxclaims" }, "FAIL  agent sandbox API is served: extensions.agents.x-k8s.io/v1beta1 does not serve sandboxclaims"},
		"api lacks sandboxtemplates":       {func(c *cluster) { c.missingResource = "sandboxtemplates" }, "FAIL  agent sandbox API is served: extensions.agents.x-k8s.io/v1beta1 does not serve sandboxtemplates"},
		"api lacks sandboxes":              {func(c *cluster) { c.missingResource = "sandboxes" }, "FAIL  agent sandbox API is served: agents.x-k8s.io/v1beta1 does not serve sandboxes"},
	} {
		c := healthy()
		tc.mutate(&c)
		wantLine(t, name, run(t, c), tc.want)
	}
}

func TestPreflightAllowsOrdinarySecrets(t *testing.T) {
	c := healthy()
	c.extraVolumes = []any{
		map[string]any{"name": "ordinary", "secret": map[string]any{"secretName": "ordinary"}},
		map[string]any{"name": "settings", "configMap": map[string]any{"name": "ordinary-config"}},
	}
	c.worker["volumeMounts"] = append(c.worker["volumeMounts"].([]any),
		map[string]any{"name": "ordinary", "mountPath": "/o"},
		map[string]any{"name": "settings", "mountPath": "/s"})
	c.worker["env"] = []any{map[string]any{"name": "V", "valueFrom": map[string]any{"secretKeyRef": map[string]any{"name": "ordinary", "key": "value"}}}}
	if out := run(t, c); strings.Contains(out, "FAIL") {
		t.Errorf("an Opaque secret and an ordinary ConfigMap in the worker should pass:\n%s", out)
	}
}
