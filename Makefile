# environment
BUILDDIR            := ${CURDIR}/build
CACHEDIR            := ${CURDIR}/.cache
TOOLBIN             := ${CURDIR}/tools
TOOLEXE             := ${TOOLBIN}/tool
GOBIN               := ${CACHEDIR}/tools

# app specific info
APP_VER             ?= v$(shell git describe --always --tags|sed 's/^v//')
GITHASH             ?= $(shell git rev-parse --short HEAD)
VERSION_VAR         := github.com/dropwhile/icanbringthat/internal/util.Version
DB_DSN              ?=
GOOSE_DRIVER        ?= postgres
GOOSE_DBSTRING      ?= ${DB_DSN}
GOOSE_MIGRATION_DIR ?= database/migrations

# flags and build configuration
GOBUILD_OPTIONS     ?= -trimpath
GOTEST_FLAGS        ?=
GOTEST_BENCHFLAGS   ?=
GOBUILD_DEPFLAGS    ?= -tags netgo,production
GOBUILD_LDFLAGS     ?= -s -w
GOBUILD_FLAGS       := ${GOBUILD_DEPFLAGS} ${GOBUILD_OPTIONS} -ldflags "${GOBUILD_LDFLAGS} -X ${VERSION_VAR}=${APP_VER}"

# cross compile defs
CC_BUILD_TARGETS     = server client
CC_BUILD_ARCHES      = darwin/amd64 darwin/arm64 freebsd/amd64 linux/amd64 linux/arm64 windows/amd64
CC_OUTPUT_TPL       := ${BUILDDIR}/bin/{{.Dir}}.{{.OS}}-{{.Arch}}

# misc
DOCKER_PREBUILD     ?=
DOCKER_POSTBUILD    ?=
PGDATABASE          ?= icanbringthat
PGPASSWORD          ?=
REDIS_PASS          ?=

# some exported vars (pre-configure go build behavior)
export GOTOOLCHAIN=local
#export CGO_ENABLED=0
export GOOSE_DRIVER
export GOOSE_DBSTRING
export GOOSE_MIGRATION_DIR
export GOBIN
export PATH := ${GOBIN}:${TOOLBIN}:${PATH}

define HELP_OUTPUT
Available targets:
* help                this help (default target)
  clean               clean up
  build               build binaries
  docker-build        build a deployable docker image
  check               run checks and validators
  nilcheck            run nilcheck; noisy/false positives, so not enabled by default
  deadcode            run deadcode; noisy/false positives, so not enabled by default
  cover               run tests with cover output
  bench               run benchmarks
  test                run tests
  generate            run code generators (go:generate, etc)
  emit-license-deps   run go-licenses
  clean-generated     remove generated files
  migrate             runs db migrations
  update-go-deps      updates go.mod and go.sum files
  cloc                counts lines of code
  dev-db-create       creates a docker postgres for development
  dev-db-start        starts a previously created docker postgres for development
  dev-db-stop         stops a previously created docker postgres for development
  dev-db-purge        deletes/destroys a previously created docker postgres for development
  run                 run local server
  devrun              run local server with modd, monitoring/restarting with changes
endef
export HELP_OUTPUT

.PHONY: help
help:
	@echo "$$HELP_OUTPUT"

.PHONY: clean
clean:
	@rm -rf "${BUILDDIR}"

clean-cache:
	@rm -rf "${CACHEDIR}"

.PHONY: generate
generate:
	@echo ">> Generating..."
	@go generate ./...

.PHONY: emit-license-deps
emit-license-deps:
	@${TOOLEXE} go-licenses report \
		--template internal/app/resources/templates/license.tpl \
		./... \
		> LICENSE-backend.md

.PHONY: build
build:
	@echo ">> Building..."
	@[ -d "${BUILDDIR}/bin" ] || mkdir -p "${BUILDDIR}/bin"
	@(for x in ${CC_BUILD_TARGETS}; do \
		echo "...$${x}..."; \
		go build ${GOBUILD_FLAGS} -o "${BUILDDIR}/bin/$${x}" ./cmd/$${x}; \
	done)
	@echo "done!"

.PHONY: test
test:
	@echo ">> Running tests..."
	@go test -count=1 -vet=off ${GOTEST_FLAGS} ./...

.PHONY: bench
bench:
	@echo ">> Running benchmarks..."
	@go test -bench="." -run="^$$" -test.benchmem=true ${GOTEST_BENCHFLAGS} ./...

.PHONY: cover
cover:
	@echo ">> Running tests with coverage..."
	@go test -vet=off -cover ${GOTEST_FLAGS} ./...

.PHONY: clean-generated
clean-generated:
	@echo ">> Purging generated files..."
	@grep -lRE '^// Code generated by (.+) DO NOT EDIT' internal rpc | xargs rm -v

.PHONY: check
check:
	@echo ">> Running checks and validators..."
	@echo "... staticcheck ..."
	@${TOOLEXE} staticcheck ./...
	@echo "... errcheck ..."
	@${TOOLEXE} errcheck -ignoretests -exclude .errcheck-excludes.txt ./...
	@echo "... go-vet ..."
	@go vet $$(go list ./... | grep -v "github.com/dropwhile/icanbringthat/rpc")
	@echo "... nilness ..."
	@${TOOLEXE} nilness ./...
	@echo "... ineffassign ..."
	@${TOOLEXE} ineffassign ./...
	@echo "... govulncheck ..."
	@${TOOLEXE} govulncheck ./...
	@echo "... betteralign ..."
	@${TOOLEXE} betteralign ./...
	@echo "... gosec ..."
	@${TOOLEXE} gosec -quiet -exclude-generated -exclude-dir=cmd/refidgen -exclude-dir=tools ./...

.PHONY: nilcheck
nilcheck:
	@echo ">> Running nilcheck (will have some false positives)..."
	@echo "... nilaway ..."
	@${TOOLEXE} nilaway -test=false \
		-include-pkgs "github.com/dropwhile/icanbringthat" \
		-exclude-file-docstrings "@generated,Code generated by,Autogenerated by" \
		./...

.PHONY: deadcode
deadcode:
	@echo ">> Running deadcode (will have some false positives)..."
	@echo "... deadcode ..."
	@${TOOLEXE} deadcode -test ./...

.PHONY: update-go-deps
update-go-deps:
	@echo ">> updating Go dependencies..."
	@GOPRIVATE=github.com/dropwhile go get -u all
	@go mod tidy

.PHONY: migrate
migrate:
	@echo ">> running migrations..."
	@${TOOLEXE} goose up

.PHONY: cloc
cloc:
	@echo ">> counting stuff..."
	@cloc -v 2 --force-lang=HTML,gohtml --fullpath --not-match-d resources/static/ .

.PHONY: dev-db-create
dev-db-create:
	@echo ">> starting dev postgres,valkey ..."
	@docker volume rm -f icanbringthat-db-init
	@docker volume create icanbringthat-db-init
	@docker create -v icanbringthat-db-init:/data --name icanbringthat-db-helper busybox true
	@for f in  ./database/init/*; do docker cp -q "$${f}" icanbringthat-db-helper:/data; done
	@docker rm -f icanbringthat-db-helper
	@docker run \
		--name icanbringthat-database \
		--restart always \
		-e POSTGRES_PASSWORD=${PGPASSWORD} \
		-e POSTGRES_DB=${PGDATABASE} \
		-p 5432:5432 \
		-v "icanbringthat-db-init:/docker-entrypoint-initdb.d/" \
		-d postgres \
		postgres -c jit=off
	@docker run \
		--name icanbringthat-valkey \
		--restart always \
		-p 6379:6379 \
		-d valkey:8-alpine \
		valkey-server --requirepass "${REDIS_PASS}"

.PHONY: dev-db-start
dev-db-start:
	@echo ">> starting dev postgres,valkey ..."
	@docker start icanbringthat-db-init
	@docker start icanbringthat-valkey

dev-db-stop:
	@echo ">> stopping dev postgres,valkey ..."
	@docker stop icanbringthat-database
	@docker stop icanbringthat-valkey

dev-db-purge:
	@echo ">> purging dev postgres,valkey ..."
	@docker rm -fv icanbringthat-database
	@docker rm -fv icanbringthat-valkey
	@docker volume rm -f icanbringthat-db-init

.PHONY: docker-build
docker-build:
	@echo ">> Building docker image..."
	@${DOCKER_PREBUILD}
	@DOCKER_BUILDKIT=1 docker build \
		--build-arg GITHASH=${GITHASH} \
		--build-arg APP_VER=${APP_VER} \
		--build-arg BUILDKIT_INLINE_CACHE=1 \
		--cache-from icanbringthat:latest \
		-t icanbringthat:latest \
		-f docker/Dockerfile \
		.
	@eval ${DOCKER_POSTBUILD}

.PHONY: run
run: build
	@echo ">> starting dev server..."
	@exec ./build/bin/server start-webserver

.PHONY: devrun
devrun:
	@echo ">> Monitoring for change, runnging tests, and restarting..."
	@${TOOLEXE} modd -f .modd.conf
