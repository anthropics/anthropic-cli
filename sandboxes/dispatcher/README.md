# cma-dispatcher

`cma-dispatcher` runs once per CMA environment per Kubernetes namespace.
It polls the environment's work queue with the environment key, gives each session its own sandbox pod through [Kubernetes Agent Sandbox](https://github.com/kubernetes-sigs/agent-sandbox), hands that pod the work item's secret payload (which embeds the session-scoped token), and acks the work item only after the pod has collected it.
Only the dispatcher ever holds the environment key, and a sandbox pod's `worker` container holds nothing the dispatcher accepts.

How a pod proves its identity depends on the API version.
On `v1beta1` (upstream) the dispatcher creates a SandboxClaim and the pod authenticates with a projected service-account token the dispatcher verifies with a Kubernetes TokenReview.
On `v1alpha1` (the GKE add-on, whose admission policy bans every kubelet-issued credential in a sandbox pod) the dispatcher creates the Sandbox directly, with a fresh per-sandbox nonce in a Secret only the fetch container reads, through an env valueFrom; the pod presents that nonce, and the dispatcher accepts it only from the address of the sandbox's own Running pod.

## Build

```sh
cd sandboxes
go build -o bin/cma-dispatcher ./dispatcher/cmd/cma-dispatcher
go test ./...
```

## Run

`cma-dispatcher` (or `cma-dispatcher serve`) runs the loop.
It exits non-zero at startup unless `work/stats` succeeds for `ANTHROPIC_ENVIRONMENT_ID` at `ANTHROPIC_BASE_URL`; it never retries against another host.
Both subcommands refuse, before any request, an `ANTHROPIC_BASE_URL` that is not a URL with a host, or that has an `@` after a `/`, `?` or `#`, and do not show its value.
The dispatcher's own requests leave out a user name and password written ahead of the host, so an error that quotes a request's URL does not hold them.
Sandboxes are handed the URL whole.
In a user name or password, write `/` as `%2F`, `?` as `%3F` and `#` as `%23`, because a raw one ends the host: `http://user:123/rest@host` reads as the host `user` and the port 123.
After the host, write `@` as `%40`: `http://host/@tenant` is refused too, because its host cannot be told from such a user name.
It then lists every SandboxClaim carrying `anthropic.com/session-id`, deletes the ones whose session the API reports terminated or missing, and starts polling.
Each claim it keeps is picked up from where the previous process left it.
Only the pod the claim's last recorded handoff went to can be running the session's work: if that pod's worker has exited since, the sandbox is suspended, and if it is still running, the pod is watched.
Any other claim is watched as awaiting the session's next pod: one whose pod is gone, has not started its worker, or is not that pod, and one with no recorded handoff, whose work item was never collected.
When that next pod appears, the exit rule of step 5 below applies, and the redelivered work item resumes the suspended sandbox.
This runs once per process, so reads are retried and a failed suspend keeps the watch: a session dropped here would stay unwatched for the life of the process.

`cma-dispatcher preflight` runs the install checks below with your own kubeconfig and exits non-zero with one sentence per failure.

`cma-dispatcher -version` prints its version and exits before it loads its configuration. It needs no environment key or cluster. `serve` logs the same version in its first log line.
Configuration comes from the environment (one Kubernetes Secret) plus flags; every configuration flag's default is the `CMA_DISPATCHER_*` variable of the same name, so the chart can set either.

| Variable | Flag | Default | Meaning |
|---|---|---|---|
| `ANTHROPIC_ENVIRONMENT_KEY` | — | required | environment key |
| `ANTHROPIC_ENVIRONMENT_ID` | — | required | environment the dispatcher serves |
| `ANTHROPIC_BASE_URL` | — | required | `https://api.anthropic.com` or a custom API base URL |
| `CMA_DISPATCHER_NAMESPACE` | `-namespace` | service-account mount | namespace of the sandboxes |
| `CMA_DISPATCHER_API_VERSION` | `-api-version` | `v1beta1` | `agents.x-k8s.io` version: `v1alpha1` for the GKE add-on, `v1beta1` for upstream |
| `CMA_DISPATCHER_SANDBOX_TEMPLATE` | `-sandbox-template` | `cma-sandbox` | SandboxTemplate the add-on path builds Sandboxes from |
| `CMA_DISPATCHER_WARM_POOL` | `-warm-pool` | `cma-sandbox` | SandboxWarmPool claims check out from |
| `CMA_DISPATCHER_SANDBOX_SERVICE_ACCOUNT` | `-sandbox-service-account` | `cma-sandbox` | ServiceAccount of the sandbox pods; `default` is refused |
| `CMA_DISPATCHER_TOKEN_AUDIENCE` | `-token-audience` | `cma-dispatcher` | audience of the pods' projected tokens |
| `CMA_DISPATCHER_LISTEN` | `-listen` | `:8443` | HTTPS listen address |
| `CMA_DISPATCHER_TLS_CERT_FILE` / `_KEY_FILE` | `-tls-cert-file` / `-tls-key-file` | `/etc/cma-dispatcher/tls/tls.{crt,key}` | serving certificate |
| `CMA_DISPATCHER_API_REQUEST_TIMEOUT` | `-api-request-timeout` | `30s` | time one CMA API request may take before it fails |
| `CMA_DISPATCHER_HANDOFF_WAIT` | `-handoff-wait` | `120s` | time for a session to reach a running pod that collects its work item; the queue's reclaim window is twice this |
| `CMA_DISPATCHER_SWEEP_INTERVAL` | `-sweep-interval` | `5m` | how often sandboxes are checked against the API |
| `CMA_DISPATCHER_RETENTION_DAYS` | `-retention-days` | `0` (off) | delete suspended sandboxes idle longer than this |
| `CMA_DISPATCHER_WORKSPACE_CLAIM_SIZE` | `-workspace-claim-size` | `10Gi` | add-on path: size of each session's `/workspace` claim; `0` backs `/workspace` with an emptyDir a suspend discards |
| `CMA_DISPATCHER_WORKSPACE_STORAGE_CLASS` | `-workspace-storage-class` | empty (cluster default) | add-on path: StorageClass of that claim |
| `CMA_DISPATCHER_MAX_CONCURRENT_HANDOFFS` | `-max-concurrent-handoffs` | `0` | sessions provisioned at once; `0` reads the warm pool's replicas, or 8 when that fails |
| `CMA_DISPATCHER_LOG_FORMAT` | `-log-format` | `text` | `text` or `json` |
| `CMA_DISPATCHER_DISPATCHER_SERVICE_ACCOUNT`, `CMA_DISPATCHER_DEPLOYMENT`, `CMA_DISPATCHER_SERVICE_HOST`, `CMA_DISPATCHER_SKIP_NETWORK_PROBE` | same names | `cma-dispatcher`, `cma-dispatcher`, `cma-dispatcher:8443`, `false` | preflight only |

Do not set `http2debug=1` or `http2debug=2` in `GODEBUG` wherever the dispatcher runs with real credentials.
With either value, Go writes to standard error each header of every HTTP/2 request the dispatcher sends, values included.
It also writes up to 256 bytes of each HTTP/2 body chunk the dispatcher receives.
The environment key, the Kubernetes token and work item secrets from the CMA API can then reach the log.
Under `preflight` the Kubernetes token is your kubeconfig's, if it uses one.
With `http2debug=2`, Go also writes each header the dispatcher receives and up to 256 bytes of each body chunk it sends.
Pod tokens and assignment nonces, and the start of each work item the dispatcher hands out, can then reach the log too.
While the setting is on, `serve` and `preflight` each log a line like this at startup:

```
level=WARN msg="GODEBUG http2debug is on: the HTTP/2 debug lines on standard error can include credentials from headers and bodies; remove http2debug from GODEBUG wherever real credentials are in use"
```

## What one work item does

1. The session handle's name is derived from the session id: lowercased, non-alphanumeric runs replaced by `-`, cut so that twelve hex digits of `sha256(session id)` still fit a 63-character label. This names the SandboxClaim on `v1beta1` and the Sandbox itself on `v1alpha1`.
2. The handle is created with labels `anthropic.com/session-id` and `app.kubernetes.io/managed-by: cma-dispatcher`, and no lifetime that would delete the session's volume. An existing handle is reused only when it carries both labels for this session; a suspended Sandbox is set back to running. On `v1alpha1` the Sandbox is created directly from the template, its fetch container reading a per-sandbox nonce Secret the Sandbox owns; an existing Sandbox is reused only while the dispatcher still holds its nonce. After a restart the dispatcher holds no nonces and may not read Secrets back. While such a Sandbox's pod still exists the handoff fails and the work item redelivers, since the pod may be running a turn. Once the pod is gone, a Sandbox with no recorded handoff and no stored workspace claim is deleted so the next handoff recreates it; one stored with a claim, or past its first collection, is kept and re-armed: the dispatcher suspends it if it is not already, writes a fresh nonce Secret and points the fetch container at it, so the resumed pod starts with the new nonce.
   The labels are only as trustworthy as the namespace's writers, since whoever can create claims can also label them.
   A handoff is recorded on the handle only after a collection (step 4), so on a Sandbox found after a restart a missing record does not show that an existing pod holds no work.
   A Sandbox stored with a workspace claim is kept even without a recorded handoff, since its volume may already hold a turn's files whose record never landed; whether it has a claim is read from the stored Sandbox, because `-workspace-claim-size` may have been set to `0` after it was created.
   A pod built before the fetch container was pointed at the new Secret presents a nonce the dispatcher does not hold and is refused.
   Any failed read while resolving such a Sandbox fails the handoff.
3. The dispatcher waits for the claim's Sandbox to be Ready (on `v1beta1`) and its pod (named by the Sandbox's `agents.x-k8s.io/pod-name` annotation or after the Sandbox, controlled by the Sandbox, created no earlier than the Sandbox's last resume, not terminating) to be Running, then records a work item keyed by the pod UID.
4. The pod's `fetch` container calls `GET /work-item`; the dispatcher acks the work item only after that call succeeds, and stamps `anthropic.com/last-active` on the claim. If the pod never asks within the handoff wait, the work item is dropped and the item is left for the queue to redeliver. The poll's reclaim window is twice the handoff wait, and a redelivered item whose pod still holds it is acked at once. The work id is the session id, and the API queues the item again on a message to a session whose worker has exited, so once the pod's worker is gone the same id is handed off afresh, resuming the suspended Sandbox.
   The handoff (the pod's UID and its worker's restart count) is recorded on the claim only after the pod collected the item, so a claim never records a handoff that did not take.
   The ack waits for that record: an acked item never comes again, and a dispatcher restarted with the claim unrecorded would judge the pod that holds the token stale and suspend it mid-turn.
   A pod that already collected a work item is never given another: its worker still runs that item or has exited without running another, so a new item is left for the pod that follows the suspend.
5. The dispatcher watches the pod's `worker` container. When it exits after the handoff (a restart, or a terminated state), the dispatcher asks the API about the session: ended or gone, and the claim is deleted with its Sandbox, pod and volume; otherwise the claim is stamped `anthropic.com/last-active` and the Sandbox suspended, which deletes the pod and keeps the volume.

At most `-max-concurrent-handoffs` sessions are provisioned at a time; the loop does not poll while that many are in flight.

## HTTPS routes

Routes are exact-path, exact-method. `GET /healthz` is unauthenticated and answers 200.
`GET /work-item` authenticates by the path's identity model.
On `v1beta1` it verifies the bearer with a Kubernetes TokenReview for audience `cma-dispatcher`: the token must verify for that audience, belong to `system:serviceaccount:<namespace>:<sandbox service account>`, and carry `authentication.kubernetes.io/pod-uid`.
On `v1alpha1` it reads the `X-CMA-Assignment-Nonce` header and accepts the request only when the nonce matches a live sandbox (constant-time compare), the request's TCP peer address is that sandbox's pod's address, and the pod is observed Running; a deleted pod or a changed address is refused, and every ambiguity fails closed.
Nothing but the resolved pod UID selects the work item.
It returns `{"session_id","work_id","secret","base_url","environment_id"}`, `secret` being the work queue's secret payload verbatim: 404 means nothing is assigned to the pod yet and the caller should retry; 410 means the pod's work item is closed and it must not expect another. The same pod, identified by its verified token, gets the same work item again until the dispatcher closes it, which happens when the pod's worker exits or the pod is gone. The work item counts as collected only once the response has been written to the connection; a failed write keeps it for the retry.

## Sweep

Every sweep interval the dispatcher asks the API about each managed claim's session and deletes the claim, running pod included, when the API answers 404 `not_found_error` with the message `Session not found: <id>`, or `terminated`.
Any other answer keeps the claim: a rejected key, an outage, or any other 404, such as the page of a proxy in front of the API or the API's answer to an organization without the beta.
With `-retention-days N`, a suspended claim whose `anthropic.com/last-active` annotation is older than N days is deleted too; a claim without a parsable annotation is never counted as old.
A session whose handoff is in flight is left alone by both rules, and a handoff that resumes a suspended sandbox stamps `anthropic.com/last-active` first.

## By API version

On `v1beta1` the dispatcher creates a SandboxClaim naming `spec.warmPoolRef`, and suspends by writing `spec.operatingMode` (`Suspended` / `Running`) on the Sandbox.
On `v1alpha1` (the GKE add-on, which never reconciles warm pools) it creates the Sandbox directly from the template with `replicas: 1` and `shutdownPolicy: Retain`, the fetch container reading a fresh nonce Secret owned by the Sandbox, and suspends by writing `spec.replicas` (`0` / `1`).
The add-on's SandboxTemplate schema stores no `volumeClaimTemplates`, so the dispatcher copies none from the template; the Sandbox schema does store them, so the dispatcher writes its own `workspace` claim (`-workspace-claim-size`, 10 GiB by default, from `-workspace-storage-class` or the cluster default) onto each Sandbox, and the controller creates a PVC owned by the Sandbox and binds it to the worker's mount by name. That PVC survives a suspend and is rebound on resume, and goes with the Sandbox on delete. A template that defines its own `workspace` volume keeps it. With `-workspace-claim-size 0` the mount is backed by an explicit emptyDir instead, and `/workspace` dies with the suspended pod.
The Sandbox carries the same `anthropic.com/session-id` and handoff annotations a claim would, so the loop treats it as its own session handle.

## Preflight checks

The `v1beta1` checks:

1. `agents.x-k8s.io/<version>` serves `sandboxes` and `extensions.agents.x-k8s.io/<version>` serves `sandboxclaims`, `sandboxtemplates`, `sandboxwarmpools`.
2. The SandboxWarmPool's `sandboxTemplateRef` names the SandboxTemplate the next checks inspect, since session pods come from the pool's template.
3. The SandboxTemplate's `runtimeClassName` resolves to a RuntimeClass whose handler is `runsc` or `gvisor` (gVisor) or a Kata handler.
4. The SandboxTemplate, as the cluster stores it, carries `networkPolicyManagement: Unmanaged`.
5. The pod leaves `shareProcessNamespace` unset or false. The pod sets `automountServiceAccountToken: false`. The SandboxTemplate has a `fetch` container that mounts the projected token volume and a `worker` container. No container reads a Secret of type `kubernetes.io/service-account-token`. No other container, init containers included, mounts that volume or the CA, under any volume name, or sets any of `CMA_DISPATCHER_URL`, `CMA_SA_TOKEN_FILE`, `CMA_CA_FILE`.
6. Network policy enforcement: not implemented in this version; the check fails with that sentence unless `-skip-network-probe` turns it into a warning. Verify by hand that a sandbox pod cannot reach another sandbox pod, `169.254.169.254` or the cluster's API server and can reach the dispatcher at the `-service-host` value.
7. Each volume claim template's StorageClass (or the cluster default) exists; `Retain` is a warning.
8. The namespace carries `pod-security.kubernetes.io/enforce: restricted`.
9. The dispatcher Deployment has `replicas: 1` and strategy `Recreate`.
10. `work/stats` succeeds with the environment key.
11. The dispatcher ServiceAccount is allowed each verb below and may create TokenReviews (checked with SubjectAccessReviews).

On `v1alpha1` the checks differ in these ways:

- The API check requires only `sandboxes` and `sandboxtemplates`.
- The warm-pool check is dropped.
- The token/fetch check (5) becomes a check of the template. The pod leaves `shareProcessNamespace` unset or false. The pod sets `automountServiceAccountToken: false`. The template has `fetch` and `worker` containers. No volume projects a service-account token, ClusterTrustBundle or PodCertificate. No volume uses a hostPath. Every volume is an emptyDir, configMap, secret or persistentVolumeClaim volume, or a projected volume with a configMap or secret source. Every container, init containers included, has a cpu and a memory limit. No container carries the environment key as an env of that name, a valueFrom, an envFrom, or a mounted Secret or ConfigMap. A Secret or ConfigMap the check cannot read counts as carrying it. No container reads a Secret of type `kubernetes.io/service-account-token` by mount, envFrom or env valueFrom. A Secret the check cannot read counts as one. No container but `fetch` sets any of `CMA_DISPATCHER_URL`, `CMA_SA_TOKEN_FILE`, `CMA_CA_FILE`.
- The storage-class check (7) becomes a check of how `/workspace` is backed. A template `workspace` volume backed by a persistentVolumeClaim passes. A template `workspace` volume with any other source warns that a suspend discards `/workspace`. That warning stands even with a workspace claim size configured. With no template `workspace` volume and a workspace claim size configured, the claim's StorageClass (or the cluster default) must exist. `Retain` is a warning. With neither, a warning says an emptyDir backs `/workspace` and a suspend discards it.
- The Sandbox the dispatcher would create from the template, nonce Secret reference included, passes a dry-run create. No field may be dropped by the served schema.
- The RBAC check drops TokenReview.

## RBAC the binary needs

Namespaced Role, `v1beta1`:

| API group | resource | verbs |
|---|---|---|
| `extensions.agents.x-k8s.io` | `sandboxclaims` | create, get, list, watch, delete, patch |
| `extensions.agents.x-k8s.io` | `sandboxwarmpools` | get |
| `agents.x-k8s.io` | `sandboxes` | get, list, watch, patch |
| `""` | `pods` | get, list, watch |

Plus a ClusterRoleBinding to `system:auth-delegator` for TokenReview.
On `v1alpha1` the dispatcher instead needs `sandboxes` create, get, list, watch, patch, delete, `sandboxtemplates` get, `secrets` create, and no SandboxClaim, warm-pool, or TokenReview access.
The chart's Role carries the `v1beta1` grants, and only `platform: gke` adds the add-on path's, so the upstream installs never hold verbs they do not use.
`create` is the only Secrets verb: the dispatcher writes each sandbox's nonce Secret once and can never read, change, or delete one, and each nonce Secret is garbage-collected through its ownerReference to its Sandbox.
Because the dispatcher cannot delete Secrets, each re-arm after a restart leaves the Sandbox's earlier nonce Secret in place until the Sandbox itself is deleted.
`patch` on the session handle is for the `anthropic.com/last-active` annotation; `patch` on `sandboxes` suspends and resumes; `watch` on `pods` is how a worker's exit is seen.
