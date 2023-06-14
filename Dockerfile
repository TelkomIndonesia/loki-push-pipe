# syntax = docker/dockerfile:1.2

FROM golang:1.20 AS builder

WORKDIR /src
COPY ./ ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o loki-push-pipe



FROM alpine:3.16

COPY --from=builder /src/loki-push-pipe /usr/local/bin/loki-push-pipe
EXPOSE 3100
ENTRYPOINT ["/usr/local/bin/loki-push-pipe"]