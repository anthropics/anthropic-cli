#!/bin/sh
# install-agent-sandbox.sh installs the upstream Agent Sandbox controller with
# its extensions, which the sandboxes chart needs (README.md, Prerequisites). Do not
# run it on GKE, which has a controller of its own, its Agent Sandbox add-on.
#
# A release asset and an image tag can both be replaced after publication, so
# neither name alone fixes the bytes. The script refuses a manifest whose
# SHA-256 is not the one recorded below, and runs the controller's image by the
# digest recorded below. A mirror changes where the two come from, and neither
# check.
#
# Before it applies anything, it looks for the controller's Deployment and for
# the Agent Sandbox CustomResourceDefinitions, which a controller in another
# namespace has too. It installs where it finds neither, or where the Deployment
# names the digest recorded below, from any registry, so it can run again on what
# it installed. Otherwise it changes nothing, says what it found and stops.
#
# With --replace it does not look. It then puts the version below in the place
# of a controller that it or kubectl apply installed, whether newer or older.
#
# Either way, it can stop part way on a controller that another tool manages.
#
# It needs the rights to get, create and patch each kind of object in the
# manifest, and to list and watch Deployments. Kubernetes also refuses a
# ClusterRole that grants rights the user lacks.
set -eu

version=v1.0.1
asset_sha256=460c1e0272c793c98de82482aed57db1bddb0afcb979c19325dc6cc52289c44c
image=registry.k8s.io/agent-sandbox/agent-sandbox-controller
digest=sha256:e1787f95cda406e2d81322e8637b6dfa2d77514fbbd0536a872b8e5d409b7154
url="https://github.com/kubernetes-sigs/agent-sandbox/releases/download/${version}/sandbox-with-extensions.yaml"
token=/var/run/secrets/kubernetes.io/serviceaccount/token

usage() {
  cat <<USAGE
Usage: sh install-agent-sandbox.sh (--context NAME | --kubeconfig FILE | --in-cluster) [--manifest FILE] [--image REFERENCE] [--replace]

Installs the Agent Sandbox controller ${version} with its extensions.

  --context NAME      the kubectl context of the cluster to install on
  --kubeconfig FILE   a kubeconfig file whose current context is the cluster to
                      install on
  --in-cluster        the cluster whose pod this runs in, under the pod's
                      service account
  --manifest FILE     a copy of the release's manifest, read in place of the
                      download
  --image REFERENCE   your registry's copy of the controller's image
  --replace           install without looking for a controller or its
                      CustomResourceDefinitions on the cluster, which stop the
                      script otherwise

The manifest:  ${url}
  its SHA-256: ${asset_sha256}
The image:     ${image}:${version}@${digest}
USAGE
}
die() {
  printf 'install-agent-sandbox: %s\n' "$1" >&2
  exit "${2:-1}"
}

context='' kubeconfig='' manifest='' mirror='' replace='' clusters=0
while [ $# -gt 0 ]; do
  case $1 in
    -h | --help)
      usage
      exit 0
      ;;
    --in-cluster)
      clusters=$((clusters + 1))
      shift
      continue
      ;;
    --replace)
      replace=yes
      shift
      continue
      ;;
    --*=*)
      option=${1%%=*} value=${1#*=}
      shift
      ;;
    --context | --kubeconfig | --manifest | --image)
      [ $# -ge 2 ] || die "$1 needs a value; --help lists the options" 2
      option=$1 value=$2
      shift 2
      ;;
    *) die "$1 is no option; --help lists the options" 2 ;;
  esac
  case $value in '' | -*) die "$option needs a value; --help lists the options" 2 ;; esac
  case $option in
    --context) context=$value clusters=$((clusters + 1)) ;;
    --kubeconfig) kubeconfig=$value clusters=$((clusters + 1)) ;;
    --manifest) manifest=$value ;;
    --image) mirror=$value ;;
    *) die "$option is no option; --help lists the options" 2 ;;
  esac
done

# Where nothing names a cluster, kubectl takes its current context.
[ "$clusters" -eq 1 ] ||
  die 'name the cluster to install on, with exactly one of --context NAME, --kubeconfig FILE and --in-cluster' 2

runs=${mirror:-$image:$version}
case $runs in
  *@*)
    [ "${runs#*@}" = "$digest" ] ||
      die "--image names the digest ${runs#*@}, and the controller ${version} is ${digest}"
    runs=${runs%%@*}
    ;;
esac
# A registry may be named by an IPv6 address in brackets, which has characters
# no other part of a reference has.
plain=$runs
case $runs in
  \[*\]:*/* | \[*\]/*)
    address=${runs%%\]*}
    case ${address#\[} in
      *[!0-9A-Fa-f:]*) ;;
      *:*) plain=registry${runs#*\]} ;;
    esac
    ;;
esac
case $plain in
  '' | *//* | *[!A-Za-z0-9._/:-]*) die "--image $mirror is not an image reference" 2 ;;
esac
runs="${runs}@${digest}"

missing=''
need() { command -v "$1" >/dev/null 2>&1 || missing="$missing $1"; }
need awk
need kubectl
need mktemp
need tr
[ -n "$manifest" ] || need curl
command -v sha256sum >/dev/null 2>&1 || command -v shasum >/dev/null 2>&1 ||
  missing="$missing sha256sum (or shasum)"
[ -z "$missing" ] || die "not found on this machine:${missing}"

if [ -n "$context" ]; then
  set -- --context "$context"
elif [ -n "$kubeconfig" ]; then
  # For a file with no current context, kubectl takes the cluster whose pod it
  # runs in.
  context=$(kubectl --kubeconfig "$kubeconfig" config current-context) ||
    die "the kubeconfig file $kubeconfig names no current context; set one with kubectl config use-context"
  set -- --kubeconfig "$kubeconfig" --context "$context"
else
  # kubectl takes the pod's service account only where no kubeconfig file names
  # a cluster, so it reads an empty one in place of the user's.
  if [ -z "${KUBERNETES_SERVICE_HOST:-}" ] || [ -z "${KUBERNETES_SERVICE_PORT:-}" ] || [ ! -r "$token" ]; then
    die "--in-cluster needs the service account of a pod, and there is none here: KUBERNETES_SERVICE_HOST, KUBERNETES_SERVICE_PORT and the file $token"
  fi
  export KUBECONFIG=/dev/null
  set --
fi

# The signal traps end the script, so that the EXIT trap runs in every shell.
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
trap 'exit 1' HUP INT TERM

asset="$work/asset.yaml"
if [ -n "$manifest" ]; then
  cat "$manifest" >"$asset" || die "the manifest $manifest cannot be read"
else
  tries=1
  # The speed limit ends a try that has received nothing for 30 seconds.
  until curl -fsSL --connect-timeout 30 --speed-limit 1 --speed-time 30 -o "$asset" "$url"; do
    [ "$tries" -lt 5 ] ||
      die "could not download ${url}; fetch it another way and pass the file with --manifest"
    tries=$((tries + 1))
    sleep 3
  done
fi

# sha256_of prints the SHA-256 of FILE. macOS has shasum, and most Linux systems
# have sha256sum.
sha256_of() {
  if command -v sha256sum >/dev/null 2>&1; then
    sum=$(sha256sum <"$1")
  else
    sum=$(shasum -a 256 <"$1")
  fi
  printf '%s\n' "${sum%% *}"
}
[ "$(sha256_of "$asset")" = "$asset_sha256" ] ||
  die "agent-sandbox ${version}: the manifest does not match the recorded SHA-256 ${asset_sha256}"

# The verified manifest runs one image, the controller, named by tag. The quotes
# keep YAML from reading an IPv6 address in brackets as a list.
awk -v named="image: ${image}:${version}" -v runs="image: \"${runs}\"" '
  substr($0, length($0) - length(named) + 1) == named {
    $0 = substr($0, 1, length($0) - length(named)) runs
    changed = 1
  }
  { print }
  END { exit !changed }
' "$asset" >"$work/pinned.yaml" ||
  die "agent-sandbox ${version}: the manifest does not run ${image}:${version}"

# The list "$@" is empty under --in-cluster. bash 4.0 calls an empty one unset.
if [ -z "$replace" ]; then
  unread="kubectl could not say whether the cluster has a controller already, and nothing was changed. \
Its message above says why. This user needs to read the Deployment agent-sandbox-controller in the \
namespace agent-sandbox-system, and CustomResourceDefinitions. Run this script again once the cause \
is gone, or read Install the controller in sandboxes/kubernetes/README.md and run it with --replace, to install without looking."
  kubectl ${1+"$@"} -n agent-sandbox-system get deployment agent-sandbox-controller \
    --ignore-not-found -o 'jsonpath={.spec.template.spec.containers[*].image}' >"$work/there" || die "$unread"
  # Text from the cluster must not steer the terminal. A question mark takes the
  # place of each byte that is not printable ASCII, and a digest holds none.
  there=$(LC_ALL=C tr -c '\040-\176' '?' <"$work/there")
  case "$there " in
    *"@$digest "*) ;;
    *)
      served='{", which serves "}{.spec.versions[?(@.served==true)].name}'
      stored='{" and has stored "}{.status.storedVersions[*]}'
      found=$(kubectl ${1+"$@"} get customresourcedefinitions.apiextensions.k8s.io \
        sandboxes.agents.x-k8s.io sandboxclaims.extensions.agents.x-k8s.io \
        sandboxtemplates.extensions.agents.x-k8s.io sandboxwarmpools.extensions.agents.x-k8s.io \
        --ignore-not-found -o 'jsonpath={range .items[*]}{"  definition "}{.metadata.name}'"$served$stored"'{"\n"}{end}') ||
        die "$unread"
      found=$(printf %s "$found" | LC_ALL=C tr -c '\012\040-\176' '?')
      [ -z "$there" ] || found="  deployment agent-sandbox-controller in the namespace agent-sandbox-system, which runs $there${found:+
$found}"
      [ -z "$found" ] || die "the cluster already has an Agent Sandbox controller or its CustomResourceDefinitions, and nothing was changed. Found:
$found
Without --replace, this script installs ${version} over a controller only when that controller's image names the digest ${digest}.
- To keep what is there, do not run this script.
- If another tool manages it, such as Helm or GKE's add-on, perhaps in another namespace, change it with that tool alone.
- To put ${version} in its place, or to finish a first run that stopped part way, read Install the controller in sandboxes/kubernetes/README.md and run this script again with --replace."
      ;;
  esac
fi

kubectl ${1+"$@"} apply --server-side -f "$work/pinned.yaml" ||
  die "kubectl could not apply the whole manifest, and its message above says why. Part of it may be \
on the cluster. Once the cause is gone, run this script again, with --replace if it stops at what this run applied."
kubectl ${1+"$@"} -n agent-sandbox-system rollout status deployment/agent-sandbox-controller --timeout=180s ||
  die "the controller did not start in time; the events of its pod in the namespace agent-sandbox-system say why"
