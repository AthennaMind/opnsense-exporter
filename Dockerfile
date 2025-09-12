FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.25 AS build

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH
ARG Version

WORKDIR /go/src/github.com/AthennaMind/opnsense-exporter
COPY . .

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 \
  go build \
  -tags osusergo,netgo \
  -ldflags "-s -w -X main.version=${Version}" \
  -o /usr/bin/opnsense-exporter .

FROM --platform=${BUILDPLATFORM:-linux/amd64} gcr.io/distroless/static-debian12:latest

LABEL org.opencontainers.image.source=https://github.com/AthennaMind/opnsense-exporter
LABEL org.opencontainers.image.version=${Version}
LABEL org.opencontainers.image.authors="The AthennaMind Authors admins@athennamind.com"
LABEL org.opencontainers.image.title="OPNsense Prometheus Exporter"
LABEL org.opencontainers.image.description="Prometheus exporter for OPNsense"

COPY --from=build /usr/bin/opnsense-exporter /
EXPOSE 8080
ENTRYPOINT ["/opnsense-exporter"]
