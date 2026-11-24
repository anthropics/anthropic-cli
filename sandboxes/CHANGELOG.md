# Changelog

## [0.1.1](https://github.com/anthropics/anthropic-cli/compare/sandboxes-v0.1.0...sandboxes-v0.1.1) (2026-11-24)

### Bug Fixes

* **sandboxes:** make the worker's start check hold up under packet loss

  The check now has to pass twice in a row and waits longer for each address, so a worker starts a few seconds later.

* **sandboxes:** use the v1beta1 Agent Sandbox API on GKE clusters that serve it

  GKE 1.37.2-gke.1325000 removes the API version that 0.1.0 uses. Upgrade before the cluster reaches it, and follow "Upgrading a GKE cluster" in `sandboxes/kubernetes/README.md`, which also covers Google's warning about upgrades to 1.36.3-gke.1767000 or later. On `gke` the 0.1.0 dispatcher image does not start under this chart, so move pinned digests or `image.tag` in the same upgrade. The chart has a new setting, `gke.sandboxAPIVersion`, and refuses an upgrade by Helm 4.0.0 to 4.1.0 where both versions are served. The dispatcher has a new flag, `-claim-mode`, and `-api-version` takes `auto`.


### Chores

* bump go sdk to v2.15.0
* **sandboxes:** add a script that installs the Agent Sandbox controller on clusters other than GKE

  `sandboxes/kubernetes/install-agent-sandbox.sh` installs upstream release v1.0.1, checks its manifest against a recorded SHA-256 and pins the controller's image by digest. It changes nothing on a cluster that has another controller, unless you pass `--replace`.

## [0.1.0](https://github.com/anthropics/anthropic-cli/compare/sandboxes-v0.0.1...sandboxes-v0.1.0) (2026-11-06)

### Features

* **sandboxes:** add the Helm chart, dispatcher and sandbox image for running Claude Managed Agents self-hosted sandboxes on Kubernetes
* **sandboxes:** publish the dispatcher, sandbox and sandbox source images and the Helm chart to Google Artifact Registry


### Chores

* bump go sdk to v2.7.0
