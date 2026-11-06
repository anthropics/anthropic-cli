# The cma-dispatcher image: the static binary on a distroless base that
# carries a CA bundle (the dispatcher dials the API over TLS) and nothing
# else: no shell, no package manager. Build from the sandboxes/ directory:
#
#   docker buildx build -f dispatcher/Dockerfile -t cma-dispatcher:dev --load .
#
# Both bases are pinned by digest; bump each tag and digest together.
ARG GO_IMAGE=golang:1.26.10@sha256:3745f1aa11822297b79113f5e84ce2bffeb039734571c39855809cfe9e5d7be6
ARG BASE_IMAGE=gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab
# The image the trust-store stage runs its step in: the Alpine image
# image/Dockerfile builds on. Bump the two together.
ARG FILTER_IMAGE=alpine:3.22@sha256:14358309a308569c32bdc37e2e0e9694be33a9d99e68afb0f5ff33cc1f695dce

# The compiler runs on the build machine's own architecture and cross-compiles,
# so a multi-platform build never emulates it.
FROM --platform=$BUILDPLATFORM ${GO_IMAGE} AS build
ARG TARGETOS
ARG TARGETARCH
WORKDIR /src
# The module graph alone, so source edits don't invalidate the downloaded
# dependencies.
COPY go.mod go.sum ./
RUN go mod download
# Exactly the dispatcher's build graph, so files beside it never enter the
# build's layers.
COPY dispatcher dispatcher
# The release a published image is built from, which -version prints and
# serve logs. With none passed in, the binary reports a development build.
ARG VERSION=
RUN CGO_ENABLED=0 GOOS="$TARGETOS" GOARCH="$TARGETARCH" go build -trimpath \
      -ldflags "-X github.com/anthropics/anthropic-cli/sandboxes/dispatcher/internal/version.release=$VERSION" \
      -o /out/cma-dispatcher ./dispatcher/cmd/cma-dispatcher
# The binary carries its modules' code and Go's runtime, so the image carries
# their licenses.
RUN CGO_ENABLED=0 GOOS="$TARGETOS" GOARCH="$TARGETARCH" sh dispatcher/third-party-licenses.sh /out/third_party ./dispatcher/cmd/cma-dispatcher \
    && mkdir /out/third_party/go \
    && cp "$(go env GOROOT)/LICENSE" "$(go env GOROOT)/PATENTS" /out/third_party/go/

FROM ${BASE_IMAGE} AS base

# The base's trust store without four root certificates, which
# image/ca-filter.sh lists. The base has no shell, so the step runs on a copy of
# its /etc in an Alpine stage, on the build machine's platform, because it only
# edits text files. It fails the build if one of the four is left.
FROM --platform=$BUILDPLATFORM ${FILTER_IMAGE} AS trust-store
COPY --from=base /etc /base/etc
RUN --mount=type=bind,source=image/ca-filter.sh,target=/tmp/ca-filter.sh sh /tmp/ca-filter.sh /base bundle cma-dispatcher 4

FROM base
COPY --from=trust-store /base/etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /out/cma-dispatcher /usr/local/bin/cma-dispatcher
COPY LICENSE /usr/share/doc/cma-dispatcher/LICENSE
COPY --from=build /out/third_party /usr/share/doc/cma-dispatcher/third_party
ENTRYPOINT ["/usr/local/bin/cma-dispatcher"]
