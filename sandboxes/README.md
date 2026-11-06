# Self-hosted sandboxes for Claude Managed Agents

This directory holds the pieces for running the self-hosted sandboxes of
Claude Managed Agents (CMA) for many sessions at once, one isolated sandbox
per session: a dispatcher, the image a sandbox runs, and a recipe that deploys
both on a Kubernetes cluster as one gVisor pod per session, managed by the
[Agent Sandbox](https://github.com/kubernetes-sigs/agent-sandbox) controller.
They work against `https://api.anthropic.com` and against a custom API base
URL alike; the only difference is the base URL in one Secret.

Nothing here changes the `ant` CLI. The directory is its own Go module
(`sandboxes/go.mod`), so its dependencies stay out of the CLI's
module, and the dispatcher is a separate binary.

A single `cma-dispatcher` pod holds the environment key. It polls the
environment's work queue, gives each session its own sandbox pod, and
hands that pod its work secret, which embeds a short-lived session token,
over TLS. The sandbox pod never sees the environment key. Its `fetch`
container proves who the pod is and receives the work secret: on the
upstream controller it presents a Kubernetes service-account token scoped
to the dispatcher, taken from a warm pool; on the GKE add-on, whose
admission policy bans every kubelet-issued credential in a sandbox pod, the
dispatcher creates the sandbox itself, with a nonce (a random value made for
that sandbox alone) in a Secret only that container reads. The dispatcher
accepts that nonce as often as it is presented while the sandbox exists, and
only from the address of the sandbox's running pod. It keeps nonces in memory
only, so after a restart it accepts no earlier nonce, and it gives each
existing sandbox a new one, or replaces the sandbox, before handing it work.

The pod's `worker` container runs `ant beta:worker run` for that one session
on a `/workspace` volume as another user, and holds neither the token nor the
nonce. A NetworkPolicy selects whole pods, so a command the agent runs in the
`worker` container can connect to the dispatcher's port, as the `fetch`
container does. The dispatcher gives a work secret only to a request that
carries the token or the nonce. When the worker exits, the dispatcher deletes
the sandbox if the session is over and suspends it otherwise.

The recipe has been run end to end on a GKE Standard cluster with the Agent
Sandbox add-on, against an API server reached through a custom base URL: a
session was handed to a sandbox pod, the worker started with no environment
key, ran the tool call and finished the turn, and the sandbox was deleted
when the session ended. The same objects have run the same way under upstream
controller v1.0.1, on EKS and on a kind cluster with gVisor
([Known gaps](kubernetes/README.md#known-gaps)).

## Layout

- `dispatcher/` — the Go dispatcher binary and its own README, including
  the HTTP contract the sandbox's `fetch` container speaks.
- `image/` — the sandbox image: `Dockerfile`, `ca-filter.sh`, `fetch.sh` and
  `worker.sh`.
- `kubernetes/` — the Helm chart that deploys both on a cluster, and the
  [install guide](kubernetes/README.md).

## Images

Each `sandboxes-vX.Y.Z` release publishes two multi-platform images
(`linux/amd64`, `linux/arm64`), tagged with the release's version:

- `us-docker.pkg.dev/anthropic-public-registry/sandboxes/cma-dispatcher`: the dispatcher
  binary on a distroless base with a CA bundle and nothing else.
- `us-docker.pkg.dev/anthropic-public-registry/sandboxes/cma-sandbox`: Alpine with `bash`,
  `curl`, `jq`, the two scripts, and the `ant` binary that `image/Dockerfile`
  pins and the release notes name.

Beside the sandbox image, and under the same tag, each release publishes
`us-docker.pkg.dev/anthropic-public-registry/sandboxes/cma-sandbox-source`: the source of
that image's GPL and LGPL packages, as Alpine's build recipes and patches and
the upstream archives they name. It holds files only; copy them out to read
them:

```sh
docker create --name cma-sandbox-source us-docker.pkg.dev/anthropic-public-registry/sandboxes/cma-sandbox-source:<version> -
docker cp cma-sandbox-source:/src . && docker rm cma-sandbox-source
```

Each release also publishes the Helm chart, under the same version, as
`oci://us-docker.pkg.dev/anthropic-public-registry/sandboxes/charts/cma-sandboxes`.

In the dispatcher image, `/usr/share/doc/cma-dispatcher` holds this
directory's license and, under `third_party/`, the license of every Go module
compiled into the binary. In the sandbox image,
`/usr/share/doc/cma-sandbox/THIRD_PARTY_NOTICES` lists every Alpine package
with its license and says where the source of the GPL and LGPL ones is,
`licenses/` beside it holds the text of each license named, and `ant/` holds
the license of `ant` and those of Go and every module compiled into it.

The release notes list the digests, and so does the release's
`sandboxes_X.Y.Z_checksums.txt`. Put the two image digests in the chart's
`image.dispatcherDigest` and `image.sandboxDigest`, so the cluster pulls exactly
the images that the release names: a registry tag can be moved, and a digest
cannot. A cluster that admits images by policy can allow exactly those two
digests.

## License

Everything under `sandboxes/` is licensed under the Apache License 2.0; see
[LICENSE](LICENSE).
