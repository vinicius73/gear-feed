# syntax=docker/dockerfile:1

FROM golang:1.20-alpine as builder

RUN apk add --update --no-cache gcc g++

WORKDIR /app

COPY . ./

# args
ARG BUILD_NUMBER=unknown
ARG BUILD_DATE=unknown
ARG BUILD_REF=unknown
ARG PKG=github.com/vinicius73/gamer-feed

RUN go build \
  -a -installsuffix cgo -ldflags '-s -w -extldflags "-static"' \
  -ldflags "-X $PKG/pkg.commit=$BUILD_REF -X $PKG/pkg.version=$BUILD_NUMBER -X $PKG/pkg.buildDate=$BUILD_DATE" \
  -o ./bin/gearsfeed ./apps/cli

FROM alpine:3

RUN apk add --update --no-cache ca-certificates tzdata sqlite

# args
ARG BUILD_NUMBER=unknown
ARG BUILD_DATE=unknown
ARG BUILD_REF=unknown
ARG PKG=github.com/vinicius73/gamer-feed

# Labels.
LABEL org.opencontainers.image.title="gearsfeed" \
  org.opencontainers.image.description="" \
  org.opencontainers.image.url="https://$PKG" \
  org.opencontainers.image.source="https://$PKG" \
  org.opencontainers.image.revision="$BUILD_REF"

# Environment
ENV BUILD_NUMBER=$BUILD_NUMBER \
  APP_VERSION=$BUILD_NUMBER

ENV UID=1000
ENV GID=1000
ENV UMASK=022

WORKDIR /gfeed

COPY --from=builder /app/bin /sbin

ENTRYPOINT ["/sbin/gearsfeed"]
