{{/*
Everything .Values.platform switches. Templates read it through
include "cma-sandboxes.platformPreset" . | fromYaml and never compare
.Values.platform themselves.
*/}}
{{- define "cma-sandboxes.platformPreset" -}}
{{- if eq .Values.platform "gke" -}}
sandboxAPIVersion: v1alpha1
sandboxIdentity: nonce
sandboxSchema: addon
workspaceClaimOn: dispatcher
warmPoolReplicas: 0
podPlacement: gkeSandbox
dnsPolicies: gke
createRuntimeClass: false
egressExceptExtra: [{{ .Values.gke.serviceCIDR | quote }}]
{{- else if eq .Values.platform "eks" -}}
sandboxAPIVersion: v1beta1
sandboxIdentity: projectedToken
sandboxSchema: upstream
workspaceClaimOn: template
warmPoolReplicas: {{ include "cma-sandboxes.warmPoolSize" . }}
podPlacement: runtimeClass
dnsPolicies: baseline
createRuntimeClass: {{ .Values.runtimeClass.create }}
egressExceptExtra: []
{{- else -}}
sandboxAPIVersion: v1beta1
sandboxIdentity: projectedToken
sandboxSchema: upstream
workspaceClaimOn: template
warmPoolReplicas: {{ include "cma-sandboxes.warmPoolSize" . }}
podPlacement: runtimeClass
dnsPolicies: baseline
createRuntimeClass: false
egressExceptExtra: []
{{- end -}}
{{- end -}}

{{/* sandbox.warmPoolSize, or 2 when it is unset. 0 is a value, not unset. */}}
{{- define "cma-sandboxes.warmPoolSize" -}}
{{- if kindIs "invalid" .Values.sandbox.warmPoolSize -}}2{{- else -}}{{ int .Values.sandbox.warmPoolSize }}{{- end -}}
{{- end -}}

{{/*
The reference of one image, by digest when it has one:
(dict "root" $ "repository" .Values.image.sandbox "digest" .Values.image.sandboxDigest).
*/}}
{{- define "cma-sandboxes.image" -}}
{{- $image := .root.Values.image -}}
{{- $name := printf "%s/%s" (trimSuffix "/" $image.registry) .repository -}}
{{- if .digest -}}
{{- printf "%s@%s" $name .digest -}}
{{- else -}}
{{- printf "%s:%s" $name (default .root.Chart.AppVersion $image.tag) -}}
{{- end -}}
{{- end -}}

{{/*
The refusals whose message tells the customer what to set. values.schema.json
refuses unknown keys, wrong types and values outside a fixed list.
*/}}
{{- define "cma-sandboxes.validate" -}}
{{- $onGKE := eq .Values.platform "gke" -}}
{{- if not .Values.image.registry -}}
{{- fail (printf "image.registry is empty. Set it to us-docker.pkg.dev/anthropic-public-registry/sandboxes, where each release publishes %s and %s, or to the registry you pushed them to, for example us-docker.pkg.dev/my-project/cma." .Values.image.dispatcher .Values.image.sandbox) -}}
{{- end -}}
{{- if ne (empty .Values.image.dispatcherDigest) (empty .Values.image.sandboxDigest) -}}
{{- fail "Only one of image.dispatcherDigest and image.sandboxDigest is set, so the other image would still be pulled by a tag, which can be moved. Set both, from the release notes of the release you are installing, or neither." -}}
{{- end -}}
{{- if and $onGKE (gt (int .Values.sandbox.warmPoolSize) 0) -}}
{{- fail (printf "sandbox.warmPoolSize is %d, and platform gke keeps no warm pool: a warm sandbox carries no per-sandbox nonce, so its pod could never prove who it is to the dispatcher. Remove sandbox.warmPoolSize or set it to 0." (int .Values.sandbox.warmPoolSize)) -}}
{{- end -}}
{{- if and (eq .Values.egress.mode "fqdnGke") (not $onGKE) -}}
{{- fail (printf "egress.mode is fqdnGke, which needs the FQDNNetworkPolicy of GKE Dataplane V2, and platform is %s. Set egress.mode to publicInternet, or set platform to gke on such a cluster." .Values.platform) -}}
{{- end -}}
{{- if and (eq .Values.gke.dns "cloudDns") (not $onGKE) -}}
{{- fail (printf "gke.dns is cloudDns, which only GKE serves, and platform is %s. Set gke.dns back to kubeDns, or set platform to gke." .Values.platform) -}}
{{- end -}}
{{- $nodeLocalCache := eq .Values.gke.dns "nodeLocalCache" -}}
{{- if and $nodeLocalCache (not $onGKE) -}}
{{- fail (printf "gke.dns is nodeLocalCache, which is for a GKE cluster that runs NodeLocal DNSCache, and platform is %s. Set gke.dns back to kubeDns, or set platform to gke." .Values.platform) -}}
{{- end -}}
{{- if and $nodeLocalCache (not .Values.gke.kubeDnsUpstreamIP) -}}
{{- fail "gke.dns is nodeLocalCache and gke.kubeDnsUpstreamIP is empty. gVisor pods get no answer through NodeLocal DNSCache, so sandboxes send their DNS queries to the cluster's kube-dns-upstream Service instead. Set gke.kubeDnsUpstreamIP to the address that kubectl -n kube-system get svc kube-dns-upstream -o jsonpath='{.spec.clusterIP}' prints." -}}
{{- end -}}
{{- if and .Values.gke.kubeDnsUpstreamIP (not $nodeLocalCache) -}}
{{- fail (printf "gke.kubeDnsUpstreamIP is %s and gke.dns is %s, so no sandbox would use the address. Set gke.dns to nodeLocalCache if the cluster runs NodeLocal DNSCache, or remove gke.kubeDnsUpstreamIP." .Values.gke.kubeDnsUpstreamIP .Values.gke.dns) -}}
{{- end -}}
{{- if and (eq .Values.egress.mode "publicInternet") (empty .Values.egress.clusterCIDRs) (not .Values.egress.clusterIsPrivate) -}}
{{- fail "egress.clusterCIDRs is empty, so sandboxes could reach every address outside 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16 and 100.64.0.0/10, and a cluster's API endpoint is usually one. List the endpoint and any pod, service or node range of yours outside those ranges. On GKE, gcloud container clusters describe CLUSTER --format 'value(endpoint)' prints the endpoint's address, which goes in the list as a /32. On EKS, aws eks describe-cluster --name CLUSTER --query cluster.endpoint prints an https URL, and each address its host name resolves to goes in the list as a /32. If the endpoint and every such range are already inside them, set egress.clusterIsPrivate to true instead." -}}
{{- end -}}
{{- include "cma-sandboxes.validateServerLocation" .Values.egress.cmaServer -}}
{{- end -}}

{{/* Takes .Values.egress.cmaServer and refuses it unless it names exactly one complete location. */}}
{{- define "cma-sandboxes.validateServerLocation" -}}
{{- $set := list -}}
{{- range list "public" "pods" "cidr" -}}
{{- if get $ . -}}{{- $set = append $set . -}}{{- end -}}
{{- end -}}
{{- $locations := "egress.cmaServer.public (true for api.anthropic.com or any server on a public address), egress.cmaServer.pods (a server in this cluster) or egress.cmaServer.cidr (a server on a private address outside it)" -}}
{{- if not $set -}}
{{- fail (printf "egress.cmaServer names no location for the CMA server. Set exactly one of %s." $locations) -}}
{{- end -}}
{{- if gt (len $set) 1 -}}
{{- fail (printf "egress.cmaServer names %d locations for the CMA server: %s and %s are set. Keep exactly one of %s." (len $set) (initial $set | join ", ") (last $set) $locations) -}}
{{- end -}}
{{- if and .pods (not (and .pods.namespace .pods.matchLabels)) -}}
{{- fail "egress.cmaServer.pods is incomplete. Set both egress.cmaServer.pods.namespace, the namespace the server's pods run in, and egress.cmaServer.pods.matchLabels, the labels that select them." -}}
{{- end -}}
{{- end -}}
