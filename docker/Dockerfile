# syntax=docker/dockerfile:1.6

# do the building
FROM golang:1.21-alpine as builder
RUN apk add --no-cache ca-certificates tzdata make git
WORKDIR /workdir
ENV GOEXPERIMENT=loopvar
RUN go install github.com/pressly/goose/v3/cmd/goose@latest; rm -rf /go/pkg/ /root/.cache/
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
ARG GITHASH
ARG APP_VER
RUN make build APP_VER="${APP_VER}" GITHASH="${GITHASH}"; rm -rf /root/.cache/

# make runnable image
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder --link /go/bin/goose /app/goose
COPY --from=builder --link /workdir/build/bin/* /app/
COPY --from=builder --link /workdir/database/migrations /app/database/migrations

ENV DEBUG=false
ENV BIND_ADDRESS=0.0.0.0
ENV BIND_PORT=8000
ENV TPL_DIR=embed
ENV STATIC_DIR=embed
ENV TPL_DIR="embed"
ENV STATIC_DIR="embed"
ENV LOG_LEVEL=trace
ENV LOG_FORMAT=plain
# production (sets cookies to secure, etc)
ENV PRODUCTION=false
# csrf key
ENV HMAC_KEY=""
# mail hostname
ENV SMTP_HOSTNAME=""
ENV SMTP_PORT=""
ENV SMTP_USER=""
ENV SMTP_PASS=""
# db data
ENV DB_DSN=""
# migrations
ENV GOOSE_DRIVER=postgres
ENV GOOSE_DBSTRING=""
ENV GOOSE_MIGRATION_DIR=database/migrations
ENV PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/app"

USER nobody
EXPOSE 8000
CMD ["/app/server"]
