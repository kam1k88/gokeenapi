# syntax=docker/dockerfile:1.5
FROM --platform=${BUILDPLATFORM} golang:1.23-alpine AS builder
WORKDIR /workspace

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod     sh -c 'go mod download'

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod     sh -c 'CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -o goarapi ./cmd/goarapi'

FROM alpine:3.21
WORKDIR /opt/goarapi
COPY --from=builder /workspace/goarapi /usr/local/bin/goarapi
ENV GOKEENAPI_INSIDE_DOCKER=true
VOLUME ["/etc/gokeenapi"]
ENTRYPOINT ["/usr/local/bin/goarapi"]
CMD ["serve", "--addr", ":8080"]
