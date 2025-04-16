FROM golang:1.23.7-alpine AS dependencies

RUN apk update && apk add --no-cache \
  ca-certificates \
  musl-dev \
  gcc \
  build-base \
  git

ENV GO111MODULE=on \
  CGO_ENABLED=1 \
  GOOS=linux

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

FROM dependencies AS builder

COPY . .

RUN ./build.sh

WORKDIR /dist

RUN cp /build/main .

FROM alpine:latest as dist

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /dist/main /

ENTRYPOINT ["/main"]
