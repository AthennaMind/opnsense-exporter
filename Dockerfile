FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.21 as build

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH
ARG Version

WORKDIR /go/src/github.com/st3ga/opnsense-exporter
COPY . .

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 \
    go build \
	    -tags osusergo,netgo \
	    -ldflags "-s -w -X main.version=${Version}" \
        -o /usr/bin/opnsense-exporter .

FROM --platform=${BUILDPLATFORM:-linux/amd64} gcr.io/distroless/static-debian12:latest

LABEL org.opencontainers.image.source=https://github.com/st3ga/opnsense-exporter
LABEL org.opencontainers.image.version=${Version}
LABEL org.opencontainers.image.authors="the st3ga Authors admins@st3ga.com"
LABEL org.opencontainers.image.title="opnsense-exporter"
LABEL org.opencontainers.image.description="Prometheus exporter for OPNsense metrics"

COPY --from=build /usr/bin/opnsense-exporter /
CMD ["/opnsense-exporter"]