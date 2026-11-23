chore(sandboxes): add a script that installs the Agent Sandbox controller on clusters other than GKE

`sandboxes/kubernetes/install-agent-sandbox.sh` installs upstream release v1.0.1, checks its manifest against a recorded SHA-256 and pins the controller's image by digest. It changes nothing on a cluster that has another controller, unless you pass `--replace`.
