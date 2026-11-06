# Running Claude Managed Agents sandboxes on Kubernetes

This is the recipe for running the self-hosted sandboxes of Claude Managed
Agents (CMA) on a Kubernetes cluster, one gVisor pod per session, managed by
the [Agent Sandbox](https://github.com/kubernetes-sigs/agent-sandbox)
controller. [`../README.md`](../README.md) says what the dispatcher and the
sandbox pod each do; this file is how to install and run them.

The recipe installs with the Helm chart in `chart/`, named `cma-sandboxes`,
which each `sandboxes-v<version>` release publishes to
`oci://us-docker.pkg.dev/anthropic-public-registry/sandboxes/charts/cma-sandboxes`. One
release holds the sandboxes of one environment, in a namespace of its own.

## Layout

- `chart/` — the Helm chart `cma-sandboxes`: `values.yaml` documents every
  setting, `values.schema.json` refuses a misspelled one, and `examples/` has
  a complete settings file for each kind of cluster.
- `chart_test.go` — renders the chart with `helm` for each kind of cluster and
  checks what an install depends on.

## Prerequisites

You need a cluster whose nodes can run gVisor, the Agent Sandbox controller,
a CNI that enforces NetworkPolicy, and `kubectl` 1.27 or newer and `helm` 3
or newer on the machine you install from. The settings this list names go in
the chart's settings file ([Install](#install), step 3).

**GKE.** Enable the Agent Sandbox add-on
(`gcloud beta container clusters update CLUSTER --enable-agent-sandbox`) and
add a GKE Sandbox node pool (`--sandbox type=gvisor`); on Autopilot the pool
is provisioned for you. The add-on serves the older
`extensions.agents.x-k8s.io/v1alpha1` API, which `platform: gke` selects.
The cluster must enforce NetworkPolicy: Dataplane V2 (Autopilot has it) or
Calico on a Standard cluster. A cluster on Cloud DNS for GKE, which is the
Autopilot default, sets `gke.dns: cloudDns`. The optional hostname egress
mode needs Dataplane V2.

**EKS.** Install the upstream Agent Sandbox controller with its extensions
(`kubectl apply -f` the release's `sandbox.yaml` and `extensions.yaml`, or the
Helm chart with `controller.extensions=true`). Create a node group for
sandboxes with gVisor's `runsc` installed and registered as a containerd
handler named `runsc`, labelled and tainted
`cma.anthropic.com/sandbox-node=true`; with `platform: eks` the chart
registers the `gvisor` RuntimeClass against that group. Use a CNI that
enforces NetworkPolicy (the VPC CNI with network policy enabled, or Cilium).

**Anything else.** The upstream controller as above, a RuntimeClass for
gVisor or another sandboxed runtime, whose name goes in `runtimeClass.name`,
and a NetworkPolicy-enforcing CNI.

**Which controller version.** The recipe has run on upstream release v1.0.1
and on the add-on as GKE 1.35 ships it. Install the upstream controller from
one release's files, by version, and pin its image by digest. The controller
creates every sandbox pod in the cluster, in every namespace, so treat a new
version as a change to the sandbox boundary. Before you move to one:

1. Compare its ClusterRoles and bindings with those of the version you run.
2. Read its release notes for changes to `SandboxTemplate`, `SandboxClaim`,
   `Sandbox` and the network policy the controller manages. The chart turns
   that policy off and ships its own, so a release that stops honouring
   `networkPolicyManagement: Unmanaged` changes what a sandbox can reach.
3. Upgrade a cluster that serves no sessions first, run `preflight` there
   ([Install](#install), step 4), and run one session.

Google upgrades the GKE add-on with the cluster, so step 3 is what you have
there: run `preflight` after every cluster upgrade. It reads the template as
the cluster stores it and dry-run-creates a sandbox, so a changed schema or
admission policy fails a check.

## Install

The chart takes its own version, which is the release's, as the tag of both
images, and pulls them from
`us-docker.pkg.dev/anthropic-public-registry/sandboxes`
([Images](../README.md#images)). To run images you built yourself, build them
from a checkout of this repository at the release's tag,
`sandboxes-v<version>`, push them under that version and name the registry in
`image.registry`:

```sh
TAG=<version>   # the release's version, which is the chart's
docker buildx build -f sandboxes/dispatcher/Dockerfile -t "REGISTRY/cma-dispatcher:$TAG" --push sandboxes
docker buildx build -t "REGISTRY/cma-sandbox:$TAG" --push sandboxes/image
```

1. Create the environment
   (`ant beta:environments create --config 'type: self_hosted'`) and an
   environment key for it, which the Console issues. Create the namespace and
   give it the two labels that hold every pod in it to the restricted Pod
   Security Standard; the chart does not create the namespace. Then load the
   key, the environment's id and the base URL into one Secret; the dispatcher
   is the only pod that mounts it:

   ```sh
   kubectl create namespace cma
   kubectl label namespace cma pod-security.kubernetes.io/enforce=restricted pod-security.kubernetes.io/enforce-version=latest
   kubectl -n cma create secret generic cma-environment \
     --from-literal=ANTHROPIC_ENVIRONMENT_KEY="$KEY" \
     --from-literal=ANTHROPIC_ENVIRONMENT_ID="$ENV_ID" \
     --from-literal=ANTHROPIC_BASE_URL="https://api.anthropic.com"   # or your custom base URL
   ```

   The chart takes only that Secret's name (`environment.secretName`), so the
   key can live in your own vault, under its audit and rotation. Have the
   External Secrets Operator, or a tool like it, write a Secret of that name
   with those three keys into the namespace, in place of the last command. The
   dispatcher reads the key when it starts, so restart it after the key
   changes ([Rotating the environment key](#rotating-the-environment-key)).

   Whoever holds the key can read every session in its environment and post
   tool results to it. Grant these rights in this namespace to as few people and
   service accounts as you can, because each can give its holder the key:
   reading or changing Secrets (the same Secret holds `ANTHROPIC_BASE_URL`,
   which decides where the dispatcher sends the key), creating pods or any
   object that makes them, changing a running pod, and running a command in
   the dispatcher's pod or adding a debug container to it.

2. Make the dispatcher's TLS certificate. One self-signed certificate is its
   own CA: the dispatcher serves it, the sandboxes trust it. Its name must be
   the Service address inside the namespace:

   ```sh
   openssl req -x509 -newkey ec -pkeyopt ec_paramgen_curve:P-256 -nodes -days 825 \
     -subj "/CN=cma-dispatcher.cma.svc" \
     -addext "subjectAltName=DNS:cma-dispatcher.cma.svc" \
     -keyout tls.key -out tls.crt
   kubectl -n cma create secret tls cma-dispatcher-tls --cert=tls.crt --key=tls.key \
     --dry-run=client -o yaml | kubectl apply -f -
   kubectl -n cma create configmap cma-dispatcher-ca --from-file=ca.crt=tls.crt \
     --dry-run=client -o yaml | kubectl apply -f -
   ```

   Renew before day 825 by running the same commands again and restarting
   the dispatcher; running sandboxes keep the old CA until they are replaced.

3. Write the settings file. Your choices go in one YAML file,
   `cma-values.yaml` in the commands below. This one is for a GKE cluster
   against `api.anthropic.com`:

   ```yaml
   platform: gke
   egress:
     clusterCIDRs:
       - 203.0.113.10/32        # the cluster's API endpoint
     cmaServer:
       public: true
   ```

   Three settings must come from you:

   | Setting | What to put there |
   |---|---|
   | `platform` | The kind of cluster: `gke` for GKE with the Agent Sandbox add-on, `eks`, or `generic` for any other cluster with the upstream controller. |
   | `egress.cmaServer` | Where the CMA server is, as exactly one of `public: true` (`api.anthropic.com`, or a custom base URL on a public address), `pods` (the `namespace` and `matchLabels` of a server in this cluster) and `cidr` (a server on a private address outside the cluster, such as `10.20.0.5/32`). With `pods` or `cidr`, `port` is the server's port. |
   | `egress.clusterCIDRs`, or `egress.clusterIsPrivate: true` | The cluster's API endpoint, and any range of yours outside the private ranges, so that sandboxes are kept off them. [Network egress](#network-egress) says how to find them and when `egress.clusterIsPrivate: true` takes their place. |

   Five more depend on the cluster and on how you admit images:

   | Setting | What to put there |
   |---|---|
   | `image.dispatcherDigest`, `image.sandboxDigest` | The two digests the release's notes list ([Images](../README.md#images)). The cluster then pulls exactly those images. Left out, it pulls them by the release's version tag. |
   | `image.registry`, `image.tag` | Only when the cluster pulls from a mirror or runs images you built: the registry you pushed both images to, and their tag when it is not the release's version. |
   | `gke.dns`, `gke.kubeDnsUpstreamIP` | On GKE, how the cluster serves DNS: `kubeDns` unless you change it. [Network egress](#network-egress) says when to change it. |
   | `runtimeClass.name`, `runtimeClass.create` | The RuntimeClass sandbox pods run under, `gvisor` unless you change it. `kubectl get runtimeclass` lists the cluster's. On `eks` the chart creates it for the `runsc` node group unless `runtimeClass.create` is `false`. |
   | `sandbox.workspace.storageClassName` | The StorageClass of each session's `/workspace` volume. Left out, it is the cluster's default class, the one `kubectl get storageclass` marks `(default)`. |

   Every other setting has a default, among them the sandbox's resources,
   idle timeout and volume size, the warm pool's size, the retention period,
   and the names of the two Secrets and the ConfigMap made above
   (`environment.secretName`, `tls.secretName`, `tls.caConfigMap`).
   `chart/values.yaml` documents each one, and `chart/examples/` has a
   complete file for each kind of cluster.

4. Install and check. `--wait` holds the command until the dispatcher is
   ready. `preflight` then runs from your own machine with your kubeconfig,
   not inside the pod, because it reads kinds (templates, RuntimeClasses,
   StorageClasses, the namespace, the Deployment) the dispatcher's Role does
   not grant; it needs the same three variables the Secret holds, at a base
   URL your machine can reach, and `-skip-network-probe` until the network
   probe exists (`../dispatcher/README.md` says how to build the binary and
   lists the checks). On GKE it also fails when the add-on's `v1alpha1`
   schema is too old to store `networkPolicyManagement`:

   ```sh
   helm install cma oci://us-docker.pkg.dev/anthropic-public-registry/sandboxes/charts/cma-sandboxes --version <version> --namespace cma -f cma-values.yaml --wait
   ANTHROPIC_ENVIRONMENT_KEY="$KEY" ANTHROPIC_ENVIRONMENT_ID="$ENV_ID" ANTHROPIC_BASE_URL="$BASE_URL" \
     cma-dispatcher preflight -namespace cma -skip-network-probe -api-version v1alpha1   # GKE
   ANTHROPIC_ENVIRONMENT_KEY="$KEY" ANTHROPIC_ENVIRONMENT_ID="$ENV_ID" ANTHROPIC_BASE_URL="$BASE_URL" \
     cma-dispatcher preflight -namespace cma -skip-network-probe                        # eks, generic
   kubectl -n cma get sandboxwarmpool cma-sandbox        # READY should reach 2 (stays 0 on GKE)
   ```

   Nothing repeats these checks while the dispatcher runs. Run `preflight`
   again after any change to the cluster, its network plugin, the controller
   or the chart's settings.

Defaults worth knowing: on the upstream controller the warm pool keeps 2
pods ready (`sandbox.warmPoolSize`); on the GKE add-on the pool is empty and
each session's sandbox is created when its first work item arrives, which
took about eight seconds from creation to a running pod on a GKE Sandbox node
with the image already pulled. Each sandbox gets a 10 GiB `ReadWriteOnce`
volume at `/workspace` (`sandbox.workspace.size`), 2 GiB of scratch space at
`/tmp`, and 2 GiB at `/mnt/memory` for the session's memory stores. It
requests 0.5 CPU and 1 GiB and is limited to 2 CPU and 4 GiB
(`sandbox.resources`), and the worker stops 30 minutes after a session goes
idle (`sandbox.idleTimeout`). On the GKE add-on, whose template schema stores
no `volumeClaimTemplates`, the dispatcher puts the claim on each sandbox
itself, from the same two `sandbox.workspace` settings. Either way the volume
survives the idle suspend and is deleted with the sandbox, so a suspended
sandbox holds its disk until the session ends or the dispatcher deletes it
after `dispatcher.retentionDays` idle days, which the default of 0 turns off.

## What the chart refuses

Helm checks the settings file before it sends anything to the cluster, and
stops with a message that names the setting to fix. It refuses a setting it
does not know, so a misspelled name is an error and not a line Helm ignores.
It refuses an empty `egress.clusterCIDRs` unless `egress.clusterIsPrivate` is
`true`, so the API endpoint is never left reachable by omission. It refuses a
missing `platform`, an empty `image.registry`, one image's digest without the
other's, an `egress.cmaServer` that
names no location or more than one, and a GKE-only setting (`gke.dns:
cloudDns` or `nodeLocalCache`, `egress.mode: fqdnGke`) on another platform.
It refuses `gke.dns: nodeLocalCache` without `gke.kubeDnsUpstreamIP`, and a
`sandbox.warmPoolSize` above 0 on `gke`, which keeps no warm pool. To run the
same checks and read the objects without touching the cluster, use
`helm template cma oci://us-docker.pkg.dev/anthropic-public-registry/sandboxes/charts/cma-sandboxes --version <version> --namespace cma -f cma-values.yaml`.

## Network egress

With the default, `egress.mode: publicInternet`, the policies let a sandbox
reach the public IPv4 internet and cluster DNS and nothing else. The
`cma-sandbox-internet` policy allows everything outbound except the private
ranges (10/8, 172.16/12, 192.168/16, 100.64/10, and the link-local range that
holds the cloud metadata server); the `cma-sandbox-baseline` policy allows
DNS and, inbound, only the dispatcher. On `gke` the chart also denies
`gke.serviceCIDR`, by default 34.118.224.0/20, the Service range new GKE
clusters get, which is public address space. This is the whole egress
allowance on its own, since the template turns the controller's own managed
policy off.

Two entries are yours to add to the internet policy's `except` list before
the first session runs, and `egress.clusterCIDRs` is where they go: the
cluster's API endpoint as a `/32`, and any pod, service or node range of
yours that falls outside the blocked ranges.
`gcloud container clusters describe CLUSTER --format 'value(endpoint)'`
prints a GKE endpoint's address.
`aws eks describe-cluster --name CLUSTER --query cluster.endpoint` prints an
EKS endpoint's URL, and each address its host name resolves to goes in the
list. Set `egress.clusterIsPrivate: true` instead only when all of them
already sit inside the blocked ranges. Sandbox nodes must not carry public
addresses (use a private cluster or a private node pool), or their range goes
in that list too. The worker container checks at start that it cannot connect
to the cloud metadata address or the cluster's API Service address, and
refuses to run when it can; the check covers only those addresses, not the
control plane's own endpoint or any other range of yours.

Those policies do not govern traffic from a pod to the cluster's own nodes.
Kubernetes always allows traffic between a pod and the node it runs on, and
most policy engines do not apply address-range rules to cluster nodes, so the
ports a node listens on are closed on the nodes themselves, as
[Node isolation you own](#node-isolation-you-own) describes.

An API server on a private address is blocked by those same ranges, and the
sandboxes dial it directly for every session. The `cidr` or `pods` form of
`egress.cmaServer` opens it with a policy of its own,
`cma-sandbox-to-server`.

IPv6 egress is closed entirely, and the chart has no setting that opens it: a
cluster's IPv6 ranges are usually global addresses that no fixed except list
can name.

`egress.mode: fqdnGke` narrows egress to the addresses that named hosts
resolve to instead: the chart leaves out the internet policy, so
`egress.clusterCIDRs` may stay empty, and allows port 443 on each host in
`egress.fqdnGke.hosts`. The list must hold the CMA server's host, which is
`api.anthropic.com`, the default, unless you use a custom base URL, and every
host your agents' commands need, such as a package index or a git host,
because those commands run inside the sandbox. Use it for sessions that handle
sensitive data: a prompt-injected agent can send what it reads to any address
the sandbox can reach. The mode is not the default because it is GKE-only,
needs Dataplane V2 and is hostname-based: a connection to an IP the sandbox
never resolved, or to an address that a permitted name resolved to earlier, is
not matched by it. It filters where a sandbox connects, not what it sends.

The list controls where a sandbox may connect, not which names it may look up.
A sandbox can still ask the cluster's DNS to resolve any name, so code inside
it can send data out inside those lookups, to whoever runs the name server of
a domain it chooses. Unlike `publicInternet`, this mode also has no list of
excluded internal ranges: it trusts each listed name's DNS to say where that
host is. List only hosts whose DNS you trust.

On other platforms `publicInternet` is the only mode. Kubernetes lets traffic
through when any policy allows it, so a NetworkPolicy of your own cannot narrow
what `cma-sandbox-internet` allows.

`gke.dns` tells the chart how a GKE cluster serves DNS: `kubeDns`, GKE's
default, `cloudDns` on Cloud DNS for GKE, which Autopilot clusters use, or
`nodeLocalCache`. On a cluster running NodeLocal DNSCache, DNS queries from
gVisor pods to the kube-dns ClusterIP are dropped, because GKE Sandbox does
not support the cache's node-level interception, so `nodeLocalCache` points
the sandboxes at the cluster's `kube-dns-upstream` Service instead. Put the
address that
`kubectl -n kube-system get svc kube-dns-upstream -o jsonpath='{.spec.clusterIP}'`
prints in `gke.kubeDnsUpstreamIP`. The sandboxes get that nameserver and the
cluster's usual search list, which the fetch container needs because it dials
the dispatcher by its short Service name.

## A second environment and upgrades

A second environment is a second namespace, labelled the same way, with its
own Secret and its own certificate, whose name is that namespace's Service
address, and a second `helm install` into it with the same settings file.

On `eks`, leave `runtimeClass.create` out of that file. The first release
keeps creating the RuntimeClass, and every later release passes
`--set runtimeClass.create=false` on its own `helm install` and on every
`helm upgrade` of that release:

```sh
helm install cma-second oci://us-docker.pkg.dev/anthropic-public-registry/sandboxes/charts/cma-sandboxes --version <version> --namespace cma-second -f cma-values.yaml --set runtimeClass.create=false --wait
```

A RuntimeClass belongs to the whole cluster, the first release created it,
and Helm refuses to create an object another release owns. With
`runtimeClass.create: false` in the shared file instead, the first release's
next upgrade would remove the RuntimeClass that the sandboxes of every
environment run under.

To upgrade, run `helm upgrade` with the newer release's version and the same
settings file, then run `preflight` again, built from that release's tag:

```sh
helm upgrade cma oci://us-docker.pkg.dev/anthropic-public-registry/sandboxes/charts/cma-sandboxes --version <newer version> --namespace cma -f cma-values.yaml --wait
```

The newer chart names the newer images by its own version. If your settings
file holds the two digests, replace them with the newer release's, or the old
images keep running. If you mirror the
images or build your own, push the newer ones under the newer release's
version first, and if your settings file sets `image.tag`, change it to that
tag. Pushed under the tag the cluster already runs, they restart nothing, the
old images keep running, and `helm upgrade` still reports success.

## Rotating the environment key

1. Create a new key for the same environment in the Console.
2. Replace it in the Secret and restart the dispatcher, which is the only
   reader:

   ```sh
   kubectl -n cma create secret generic cma-environment \
     --from-literal=ANTHROPIC_ENVIRONMENT_KEY="$NEW_KEY" \
     --from-literal=ANTHROPIC_ENVIRONMENT_ID="$ENV_ID" \
     --from-literal=ANTHROPIC_BASE_URL="$BASE_URL" \
     --dry-run=client -o yaml | kubectl apply -f -
   kubectl -n cma rollout restart deployment/cma-dispatcher
   ```

3. Revoke the old key once the sessions that started under it have gone
   idle. Session tokens already handed to sandboxes were issued under the old
   key, and where revoking a key also ends the tokens issued under it, those
   sessions lose their work in progress.

## Removal

The chart's network policies are what confine a running sandbox, so they go
after the last sandbox. The sandboxes the dispatcher made for sessions, and
their volumes, are not part of the Helm release: `helm uninstall` on its own
takes the policies and leaves those sandboxes running, free to reach the
metadata server, the cluster's API endpoint and the private ranges. Remove in
this order:

```sh
# 1. Stop the dispatcher, so it makes no new sandboxes:
kubectl -n cma scale deployment cma-dispatcher --replicas=0
# 2. Delete the warm pool and every sandbox claim, and on gke every sandbox:
kubectl -n cma delete sandboxwarmpool cma-sandbox
kubectl -n cma delete sandboxclaims --all
kubectl -n cma delete sandboxes --all          # gke only
#    Wait for the garbage collection: repeat this until it finds nothing.
kubectl -n cma get sandboxes,pods,pvc
# 3. Uninstall the release, which takes the network policies:
helm uninstall cma --namespace cma
# 4. Delete the namespace:
kubectl delete namespace cma
```

Step 2 lets the controller release the sandboxes and their volumes, and it is
done when the check prints `No resources found in cma namespace.` A remaining
claim keeps its volume, and a `Retain` storage class keeps the backing disk.
In step 3 Helm skips the warm pool, which is already gone, and takes the rest
of the release: the dispatcher, the network policies, and what the chart
keeps outside the namespace. Outside the namespace it keeps the cluster-wide
binding that lets the dispatcher ask the cluster who a sandbox's token belongs
to and, on `eks`, the RuntimeClass when this release created it, so remove
that release last. Revoke the environment key if the environment is not
moving elsewhere.

## Node isolation you own

gVisor is the boundary between model-authored commands and the node's
kernel, and the chart can only ask for it. The rest is cluster
configuration:

- Run sandboxes on their own node pool, tainted so nothing else lands
  there, and keep the dispatcher and any other workload off it.
- Give the sandbox nodes the least cloud identity you can: a node IAM role
  or service account with no permissions beyond pulling images and writing
  logs. On EKS set the instance metadata hop limit to 1 (IMDSv2 only) so a
  pod cannot reach the node's credentials even if a policy slips; on GKE use
  metadata concealment or Workload Identity, which the sandbox nodes need
  anyway.
- A sandbox must not be able to connect to the ports a node listens on,
  and NetworkPolicy cannot enforce that. What a sandbox finds there is the
  kubelet, which asks for credentials, sshd, and above all any node agent a
  DaemonSet puts on the pool (DaemonSets usually tolerate the pool's
  taint), many of which listen without authentication. On sandbox nodes
  whose boot you control (an EKS node group with your own launch template,
  or your own Kubernetes), run this at every boot before the kubelet
  starts, for example from a systemd unit the node's user data installs, and
  make the kubelet's start depend on that unit, so that a node without the
  rules runs no sandbox:

  ```sh
  for tables in iptables ip6tables; do
    $tables -w -t mangle -I PREROUTING 1 -i eni+ -m addrtype --dst-type LOCAL -m conntrack --ctstate NEW -j DROP
    $tables -w -I INPUT 1 -i eni+ -m conntrack --ctstate NEW -j DROP
  done
  ```

  `eni+` matches the host side of every pod's network interface under the
  AWS VPC CNI; another network plugin uses another prefix, which `ip link`
  on a node shows. The first rule drops each connection a pod opens to an
  address its node holds, node ports of NodePort and LoadBalancer Services
  included. The second drops each connection that ends at the node once
  kube-proxy has rewritten its destination. Neither touches the node's own
  probes of its pods. Install the IPv6 pair even on an IPv4 cluster, because
  pod interfaces have link-local IPv6 addresses. Only sandboxes run on these
  nodes, so nothing they need is lost, with two exceptions, where the node
  holds an address sandboxes use. With NodeLocal DNSCache, insert a rule with
  `-j ACCEPT` ahead of each of the two for the cache's address on port 53,
  for UDP and for TCP. With kube-proxy in IPVS mode, do the same for the
  cluster DNS Service's address on port 53, the `cma-dispatcher` Service's
  address on TCP port 8443, and a server in the cluster on its Service's
  address and port. Have each rule name one address, one protocol and one
  port, because a rule that accepts an address on every port opens the
  kubelet too. On GKE, where the nodes' boot
  is not yours and Dataplane V2 leaves every node's ports reachable, keep
  DaemonSets whose agents listen without authentication off the cluster or
  off the sandbox pool.
- Keep gVisor and the node kernel current; a sandbox escape is a
  gVisor or kernel bug, and both projects ship fixes regularly.
- Revisit `egress.clusterCIDRs` when the cluster's IP ranges change.

## Known gaps

- The worker needs `ant` 1.32.0 or later, the first release that reads the
  work item's secret payload from `ANTHROPIC_WORK_SECRET_FILE` and runs
  without an environment key; `../image/Dockerfile` pins a release that
  does. An older `ANT_VERSION` builds an image whose worker refuses to start.
- Both images leave four root certificates out of their trust store: emSign
  ECC Root CA - G3, emSign Root CA - C1, emSign ECC Root CA - C3 and AAA
  Certificate Services. A program in them that reads the store
  cannot connect over TLS to a service whose certificate chain ends only at
  one of those roots.
- `egress.mode: fqdnGke` does not filter DNS lookups and has no list of
  excluded internal ranges ([Network egress](#network-egress)).
- A session token has a limited life. A sandbox still running a session when
  its token expires is dropped by the server, and the session's next message
  queues new work.
- The recipe has not yet been run against `https://api.anthropic.com`; every
  end-to-end run so far used a custom base URL.
