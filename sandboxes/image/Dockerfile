# The sandbox image both containers of the SandboxTemplate run: busybox sh
# for the two scripts, bash for the worker's bash tool, curl and jq for the
# fetch container's calls to the dispatcher, and the official `ant` CLI. The
# template names each container's command; the image sets no entrypoint.
# Build with buildx so TARGETARCH is set:
#
#   docker buildx build --platform linux/amd64 -t cma-sandbox:dev --load .
#
# The same file builds the image that carries the source of the sandbox
# image's GPL and LGPL packages, which is published beside it:
#
#   docker buildx build --target source -t cma-sandbox-source:dev --load .
#
# Alpine is the base because it carries sh, nc and sha256sum in one small
# static busybox and adds bash, curl and jq from its own repository.
# Digest resolved 2026-09-03 for the alpine:3.22 multi-arch index; bump the
# tag and digest together. The Go image is the one dispatcher/Dockerfile
# pins; bump the two together.
ARG BASE_IMAGE=alpine:3.22@sha256:14358309a308569c32bdc37e2e0e9694be33a9d99e68afb0f5ff33cc1f695dce
ARG GO_IMAGE=golang:1.26.10@sha256:3745f1aa11822297b79113f5e84ce2bffeb039734571c39855809cfe9e5d7be6
ARG ANT_VERSION=2.5.0

FROM ${BASE_IMAGE} AS packages

# A fetch from Alpine's package repository fails now and then. The install
# is tried up to six times, the wait between tries doubling from 10 s
# (310 s in all).
#
# spdx-licenses-text is installed only to give third-party-notices.sh the
# license texts it copies, and leaves in the same step, so its 5 MB reach no
# layer. The script lists the packages from apk's database, so it runs after
# the last package is in.
RUN --mount=type=bind,source=third-party-notices.sh,target=/tmp/third-party-notices.sh \
    try=1; wait=10; \
    until apk add --no-cache bash curl jq spdx-licenses-text; do \
      if [ "$try" -ge 6 ]; then echo "giving up on apk add after $try tries" >&2; exit 1; fi; \
      echo "apk add failed; trying again in ${wait}s" >&2; \
      sleep "$wait"; \
      try=$((try + 1)); wait=$((wait * 2)); \
    done; \
    sh /tmp/third-party-notices.sh /usr/share/doc/cma-sandbox \
    && apk del --no-cache spdx-licenses-text

FROM packages AS ant

# The ant download is digest-pinned per release; both digests come from
# `gh release view v<version> --repo anthropics/anthropic-cli --json assets`.
# The download is tried up to six times, the wait between tries doubling
# from 10 s (310 s in all). GitHub's release downloads answer a 5xx now
# and then. A five-minute spell of them under way when the build starts
# ends before the last try. Every try goes through the digest check, so
# a retried download installs the same file a first-try one would.
ARG TARGETARCH
ARG ANT_VERSION
ARG ANT_SHA256_AMD64=2fd19dc5d15f77e37451703e5c9b9501915060fd704044255af1c45b3c2df97c
ARG ANT_SHA256_ARM64=7bbd5964f97f63485a28b80e826be1ded609913504dab3d95bb8c15da7fea979
ARG ANT_DOWNLOAD_BASE=https://github.com/anthropics/anthropic-cli/releases/download
RUN set -e; \
    case "$TARGETARCH" in \
      amd64) sha="$ANT_SHA256_AMD64" ;; \
      arm64) sha="$ANT_SHA256_ARM64" ;; \
      *) echo "unsupported architecture: $TARGETARCH" >&2; exit 1 ;; \
    esac; \
    url="${ANT_DOWNLOAD_BASE}/v${ANT_VERSION}/ant_${ANT_VERSION}_linux_${TARGETARCH}.tar.gz"; \
    try=1; wait=10; \
    until wget -qO /tmp/ant.tgz "$url"; do \
      rm -f /tmp/ant.tgz; \
      if [ "$try" -ge 6 ]; then echo "giving up on $url after $try tries" >&2; exit 1; fi; \
      echo "downloading $url failed; trying again in ${wait}s" >&2; \
      sleep "$wait"; \
      try=$((try + 1)); wait=$((wait * 2)); \
    done; \
    echo "$sha  /tmp/ant.tgz" | sha256sum -c -; \
    tar -xzf /tmp/ant.tgz -C /usr/local/bin ant; \
    rm /tmp/ant.tgz

# ant carries its modules' code and Go's runtime, so the image carries their
# licenses, read from the binary the stage above verified. Go reads another
# architecture's binary as well as its own, so this runs on the build
# machine's.
FROM --platform=$BUILDPLATFORM ${GO_IMAGE} AS ant-licenses
ARG ANT_VERSION
COPY --from=ant /usr/local/bin/ant /tmp/ant
RUN --mount=type=bind,source=ant-licenses.sh,target=/tmp/ant-licenses.sh \
    sh /tmp/ant-licenses.sh /tmp/ant "$ANT_VERSION" /out \
    && mkdir /out/third_party/go \
    && cp "$(go env GOROOT)/LICENSE" "$(go env GOROOT)/PATENTS" /out/third_party/go/

# The GPL and LGPL ask whoever hands out a binary to hand out its source
# from the same place, so each release publishes this stage's output as
# cma-sandbox-source beside cma-sandbox. The package list is the packages
# stage's own; the release builds this target with that stage replaced by
# the image it has just pushed (--build-context packages=docker-image://…),
# so the source is that of the published packages and no others.
FROM --platform=$BUILDPLATFORM ${BASE_IMAGE} AS gpl-sources
RUN try=1; wait=10; \
    until apk add --no-cache abuild curl git; do \
      if [ "$try" -ge 6 ]; then echo "giving up on apk add after $try tries" >&2; exit 1; fi; \
      echo "apk add failed; trying again in ${wait}s" >&2; \
      sleep "$wait"; \
      try=$((try + 1)); wait=$((wait * 2)); \
    done
COPY --from=packages /lib/apk/db/installed /tmp/installed
RUN --mount=type=bind,source=gpl-sources.sh,target=/tmp/gpl-sources.sh \
    sh /tmp/gpl-sources.sh /tmp/installed /src

FROM scratch AS source
COPY --from=gpl-sources /src /src

FROM ant

# ca-filter.sh removes four root certificates from the store and fails the build
# if one is left. No step after it may install a package, because an install can
# rebuild the store.
RUN --mount=type=bind,source=ca-filter.sh,target=/tmp/ca-filter.sh sh /tmp/ca-filter.sh / any cma-sandbox 3

COPY --from=ant-licenses /out /usr/share/doc/cma-sandbox/ant
COPY --chmod=0755 fetch.sh /usr/local/bin/cma-sandbox-fetch
COPY --chmod=0755 worker.sh /usr/local/bin/cma-sandbox-worker

# uid 1000 is the fetch container and uid 1001 the worker, as the
# SandboxTemplate runs them; both are in group 1000, the template's fsGroup,
# which is how the worker reads what fetch writes. HOME is the workspace
# volume, the one path the model's tools may write to.
RUN addgroup -g 1000 cma && adduser -D -H -u 1000 -G cma fetch && adduser -D -H -u 1001 -G cma worker \
    && mkdir -p /workspace && chown 1001:1000 /workspace
USER 1001
WORKDIR /workspace
ENV HOME=/workspace
