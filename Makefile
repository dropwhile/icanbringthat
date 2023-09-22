# environment
BUILDDIR            := ${CURDIR}/build
ARCH                := $(shell go env GOHOSTARCH)
OS                  := $(shell go env GOHOSTOS)
GOVER               := $(shell go version | awk '{print $$3}' | tr -d '.')

# app specific info
APP_VER             := $(shell git describe --always --tags|sed 's/^v//')
GITHASH             := $(shell git rev-parse --short HEAD)
GOPATH              := $(shell go env GOPATH)
VERSION_VAR         := main.ServerVersion
DB_DSN              := $(or ${DB_DSN},"postgres://postgres:password@127.0.0.1:5432/icbt?sslmode=disable")
GOOSE_DRIVER        ?= postgres
GOOSE_DBSTRING      ?= ${DB_DSN}
GOOSE_MIGRATION_DIR ?= database/migrations

# flags and build configuration
GOBUILD_OPTIONS     := -trimpath
GOTEST_FLAGS        :=
GOTEST_BENCHFLAGS   :=
GOBUILD_DEPFLAGS    := -tags netgo,production
GOBUILD_LDFLAGS     ?= -s -w
GOBUILD_FLAGS       := ${GOBUILD_DEPFLAGS} ${GOBUILD_OPTIONS} -ldflags "${GOBUILD_LDFLAGS} -X ${VERSION_VAR}=${APP_VER}"

# cross compile defs
CC_BUILD_TARGETS     = server refgen
CC_BUILD_ARCHES      = darwin/amd64 darwin/arm64 freebsd/amd64 linux/amd64 linux/arm64 windows/amd64
CC_OUTPUT_TPL       := ${BUILDDIR}/bin/{{.Dir}}.{{.OS}}-{{.Arch}}

# some exported vars (pre-configure go build behavior)
export GO111MODULE=on
#export CGO_ENABLED=0
## enable go 1.21 loopvar "experiment"
export GOEXPERIMENT=loopvar
export GOOSE_DRIVER
export GOOSE_DBSTRING
export GOOSE_MIGRATION_DIR

define HELP_OUTPUT
Available targets:
  help                this help
  clean               clean up
  all                 build binaries and man pages
  check               run checks and validators
  test                run tests
  cover               run tests with cover output
  bench               run benchmarks
  build               build all binaries
endef
export HELP_OUTPUT

.PHONY: help
help:
	@echo "$$HELP_OUTPUT"

.PHONY: clean
clean:
	@rm -rf "${BUILDDIR}"

.PHONY: setup
setup:

.PHONY: setup-check
setup-check: ${GOPATH}/bin/staticcheck ${GOPATH}/bin/gosec ${GOPATH}/bin/govulncheck

${GOPATH}/bin/staticcheck:
	go install honnef.co/go/tools/cmd/staticcheck@latest

${GOPATH}/bin/gosec:
	go install github.com/securego/gosec/v2/cmd/gosec@latest

${GOPATH}/bin/govulncheck:
	go install golang.org/x/vuln/cmd/govulncheck@latest

${GOPATH}/bin/stringer:
	go install golang.org/x/tools/cmd/stringer@latest

.PHONY: build 
build: setup
	@echo ">> Generating..."
	@go generate ./...
	@echo ">> Building..."
	@[ -d "${BUILDDIR}/bin" ] || mkdir -p "${BUILDDIR}/bin"
	@(for x in ${CC_BUILD_TARGETS}; do \
		echo "...$${x}..."; \
		go build ${GOBUILD_FLAGS} -o "${BUILDDIR}/bin/$${x}" ./cmd/$${x}; \
	done)
	@echo "done!"

.PHONY: test 
test: setup
	@echo ">> Running tests..."
	@go test -count=1 -vet=off ${GOTEST_FLAGS} ./...

.PHONY: bench
bench: setup
	@echo ">> Running benchmarks..."
	@go test -bench="." -run="^$$" -test.benchmem=true ${GOTEST_BENCHFLAGS} ./...

.PHONY: cover
cover: setup
	@echo ">> Running tests with coverage..."
	@go test -vet=off -cover ${GOTEST_FLAGS} ./...

.PHONY: check
check: setup setup-check
	@echo ">> Running checks and validators..."
	@echo "... staticcheck ..."
	@${GOPATH}/bin/staticcheck ./...
	@echo "... go-vet ..."
	@go vet ./...
	@echo "... gosec ..."
	@${GOPATH}/bin/gosec -quiet ./...
	@echo "... govulncheck ..."
	@${GOPATH}/bin/govulncheck ./...

.PHONY: update-go-deps
update-go-deps:
	@echo ">> updating Go dependencies..."
	@go get -u all
	@go mod tidy

.PHONY: migrate
migrate:
	@echo ">> running migrations..."
	@goose up

.PHONY: dev-db-create
dev-db-create:
	@echo ">> starting dev postgres..."
	@docker volume rm -f icbt-db-init
	@docker volume create icbt-db-init
	@docker create -v icbt-db-init:/data --name icbt-db-helper busybox true
	@for f in  ./database/init/*; do docker cp -q "$${f}" icbt-db-helper:/data; done
	@docker rm -f icbt-db-helper
	@docker run \
		--name icbt-database \
		-e POSTGRES_PASSWORD=password \
		-e POSTGRES_DB=icbt \
		-p 5432:5432 \
		-v "icbt-db-init:/docker-entrypoint-initdb.d/" \
		-d postgres

.PHONY: dev-db-start
dev-db-start:
	@echo ">> starting dev postgres..."
	@docker start icbt-db-init

dev-db-stop:
	@echo ">> stopping dev postgres..."
	docker stop icbt-database

dev-db-purge:
	@echo ">> purging dev postgres..."
	@docker rm -fv icbt-database
	@docker volume rm -f icbt-db-init

.PHONY: docker-build
docker-build:
	@echo ">> Building docker image..."
	@DOCKER_BUILDKIT=1 docker build \
		--build-arg GITHASH=${(GITHASH} \
		--build-arg APP_VER=${APP_VER} \
		-t icbt \
		-f docker/Dockerfile \
		.

.PHONY: run
run: build
	@echo ">> starting dev server..."
	@./build/bin/server

.PHONY: devrun
devrun: 
	@echo ">> Monitoring for change, runnging tests, and restarting..."
	@modd -f .modd.conf

.PHONY: all
all: build
