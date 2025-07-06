# syntax=docker/dockerfile:1.6

# do the building
# optimize for caching, and not necessarily for layers,
# since the runnable image will be the thin/flatter one anyway.
#FROM golang:alpine AS builder
FROM alpine:latest AS builder
RUN apk add --no-cache ca-certificates tzdata make git mise
WORKDIR /workdir
ENV MISE_DATA_DIR="/mise"
ENV MISE_CONFIG_DIR="/mise"
ENV MISE_CACHE_DIR="/mise/cache"
ENV PATH="/mise/shims:$PATH"
COPY .mise.toml .mise.lock ./
# install go
RUN mise trust && \
    mise install go
# install goose
RUN mise install $(grep -o 'go:github.com/pressly/goose/[^"]*' .mise.toml) && \
    install -D "$(mise which goose)" build/bin/goose
# install go mod deps
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
ARG GITHASH
ARG APP_VER
# build
RUN env APP_VER="${APP_VER}" GITHASH="${GITHASH}" mise build

# make runnable image
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder --link /workdir/build/bin/* /app/
COPY --from=builder --link /workdir/database/migrations /app/database/migrations

ENV DEBUG=false
ENV BIND_ADDRESS=0.0.0.0
ENV BIND_PORT=8000
ENV TPL_DIR=embed
ENV STATIC_DIR=embed
ENV LOG_LEVEL=debug
ENV LOG_TRACE=false
ENV LOG_FORMAT=plain
# production (sets cookies to secure, etc)
ENV PRODUCTION=false
ENV ENABLE_RPC=false
# hmac key
ENV HMAC_KEY=""
# mail hostname
ENV SMTP_HOSTNAME=""
ENV SMTP_PORT=""
ENV SMTP_USER=""
ENV SMTP_PASS=""
# db data
ENV DB_DSN=""
# worker data
ENV JOBS="all"
# migrations
ENV GOOSE_DRIVER=postgres
ENV GOOSE_DBSTRING=""
ENV GOOSE_MIGRATION_DIR=database/migrations
ENV PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/app"

USER nobody
EXPOSE 8000
CMD ["/app/server", "webserver"]
