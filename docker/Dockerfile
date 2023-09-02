# syntax=docker/dockerfile:1.6

# do the building
FROM golang:1.21-alpine as builder
RUN apk add --no-cache ca-certificates make git
WORKDIR /workdir
COPY . /workdir/
RUN make build
ARG GITHASH
ARG APP_VER
RUN make build

# make runnable image
FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /workdir/build/bin/* /app/

ENV DEBUG false
ENV BIND_ADDRESS 0.0.0.0
ENV BIND_PORT 8000
ENV TPL_DIR embed
ENV STATIC_DIR embed

EXPOSE 8000
ENTRYPOINT ["/app/server"]