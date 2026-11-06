# cma-sandboxes

The Helm chart for the Kubernetes recipe that [`../README.md`](../README.md) describes.
That file is the install guide; this one says what the chart holds.

One release installs the sandboxes of one environment into one namespace: the `cma-dispatcher` Deployment, Service and RBAC, the `cma-sandbox` ServiceAccount, SandboxTemplate and SandboxWarmPool, and the NetworkPolicies.
Object names are fixed, not prefixed with the release name, because the dispatcher's defaults, the certificate's name and the fetch container's URL assume them.
`platform` chooses everything that differs between kinds of cluster (`templates/_helpers.tpl`), and `../chart_test.go` renders each one.

## What you provide first

The chart installs none of these, and `../README.md` says how to make each.

- The Agent Sandbox controller and its resource types: the GKE add-on, or the upstream controller with its extensions.
- Nodes that run gVisor, and on a cluster other than EKS the RuntimeClass that names them.
- A CNI that enforces NetworkPolicy.
- The namespace, labelled `pod-security.kubernetes.io/enforce=restricted`.
- In that namespace: the environment Secret (`environment.secretName`), the dispatcher's TLS Secret (`tls.secretName`) and the CA ConfigMap (`tls.caConfigMap`).
  No value of this chart holds the environment key or a private key.

The two images come from `us-docker.pkg.dev/anthropic-public-registry/sandboxes`, at the chart's own version, unless `image.registry` and `image.tag` say otherwise.
`image.dispatcherDigest` and `image.sandboxDigest` pin them by digest.

## Render and install

```sh
helm template cma oci://us-docker.pkg.dev/anthropic-public-registry/sandboxes/charts/cma-sandboxes --version <version> --namespace cma -f my-values.yaml
helm install cma oci://us-docker.pkg.dev/anthropic-public-registry/sandboxes/charts/cma-sandboxes --version <version> --namespace cma -f my-values.yaml
```

`examples/` has one worked values file per platform.
`values.yaml` documents every setting, and `values.schema.json` refuses a misspelled key.

The chart refuses to install while `egress.clusterCIDRs` is empty unless `egress.clusterIsPrivate` is true, so the API endpoint is never left reachable by omission.
IPv6 egress stays closed; no setting opens it.
