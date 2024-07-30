# environment
BUILDDIR            := ${CURDIR}/build
ARCH                := $(shell go env GOHOSTARCH)
OS                  := $(shell go env GOHOSTOS)
GOVER               := $(shell go version | awk '{print $$3}' | tr -d '.')

# app specific info
APP_VER             ?= v$(shell git describe --always --tags|sed 's/^v//')
GITHASH             ?= $(shell git rev-parse --short HEAD)
GOPATH              := $(shell go env GOPATH)
GOBIN               := ${CURDIR}/.tools
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
CC_BUILD_TARGETS     = server worker client
CC_BUILD_ARCHES      = darwin/amd64 darwin/arm64 freebsd/amd64 linux/amd64 linux/arm64 windows/amd64
CC_OUTPUT_TPL       := ${BUILDDIR}/bin/{{.Dir}}.{{.OS}}-{{.Arch}}

# misc
DOCKER_PREBUILD     ?=
DOCKER_POSTBUILD    ?=
PGDATABASE          ?= icanbringthat
PGPASSWORD          ?=
REDIS_PASS          ?=

# some exported vars (pre-configure go build behavior)
export GO111MODULE=on
#export CGO_ENABLED=0
## enable go 1.21 loopvar "experiment"
export GOEXPERIMENT=loopvar
export GOOSE_DRIVER
export GOOSE_DBSTRING
export GOOSE_MIGRATION_DIR
export GOBIN
export PATH := ${GOBIN}:${PATH}

define HELP_OUTPUT
Available targets:
* help                this help (default target)
  clean               clean up
  setup               fetch related tools/utils and prepare for build
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

## begin tools

# bench tools
${GOBIN}/benchstat:
	go install golang.org/x/perf/cmd/benchstat@latest

BENCH_TOOLS := ${GOBIN}/benchstat

# other tools
${GOBIN}/modd:
	go install github.com/cortesi/modd/cmd/modd@latest

OTHER_TOOLS := ${GOBIN}/modd

# generate tools
${GOBIN}/protoc-gen-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

${GOBIN}/protoc-gen-twirp:
	go install github.com/twitchtv/twirp/protoc-gen-twirp@latest

${GOBIN}/protoc-go-inject-tag:
	go install github.com/favadi/protoc-go-inject-tag@latest

${GOBIN}/go-licenses:
	go install github.com/google/go-licenses@latest

${GOBIN}/convergen:
	go install github.com/reedom/convergen@latest

${GOBIN}/ifacemaker:
	go install github.com/vburenin/ifacemaker@latest

${GOBIN}/mockgen:
	go install go.uber.org/mock/mockgen@latest

${GOBIN}/stringer:
	go install golang.org/x/tools/cmd/stringer@latest

GENERATE_TOOLS := ${GOBIN}/stringer ${GOBIN}/protoc-gen-twirp ${GOBIN}/protoc-gen-go
GENERATE_TOOLS += ${GOBIN}/convergen ${GOBIN}/go-licenses  ${GOBIN}/protoc-go-inject-tag
GENERATE_TOOLS += ${GOBIN}/ifacemaker ${GOBIN}/mockgen

# check tools
${GOBIN}/betteralign:
	go install github.com/dkorunic/betteralign/cmd/betteralign@latest

${GOBIN}/ineffassign:
	go install github.com/gordonklaus/ineffassign@latest

${GOBIN}/errcheck:
	go install github.com/kisielk/errcheck@latest

${GOBIN}/go-errorlint:
	go install github.com/polyfloyd/go-errorlint@latest

${GOBIN}/gosec:
	go install github.com/securego/gosec/v2/cmd/gosec@latest

${GOBIN}/nilaway:
	go install go.uber.org/nilaway/cmd/nilaway@latest

${GOBIN}/deadcode:
	go install golang.org/x/tools/cmd/deadcode@latest

${GOBIN}/govulncheck:
	go install golang.org/x/vuln/cmd/govulncheck@latest

${GOBIN}/staticcheck:
	go install honnef.co/go/tools/cmd/staticcheck@latest

${GOBIN}/nilness:
	go install golang.org/x/tools/go/analysis/passes/nilness/cmd/nilness@latest

CHECK_TOOLS := ${GOBIN}/staticcheck ${GOBIN}/gosec ${GOBIN}/govulncheck
CHECK_TOOLS += ${GOBIN}/errcheck ${GOBIN}/ineffassign ${GOBIN}/nilaway
CHECK_TOOLS += ${GOBIN}/go-errorlint ${GOBIN}/deadcode ${GOBIN}/betteralign
CHECK_TOOLS += ${GOBIN}/nilness

# migrate tools
${GOBIN}/goose:
	go install github.com/pressly/goose/v3/cmd/goose@latest

MIGRATE_TOOLS += ${GOBIN}/goose

## end tools

.PHONY: setup-build
setup-build: ${BUILD_TOOLS}

.PHONY: setup-generate
setup-generate: ${GENERATE_TOOLS}

.PHONY: setup-migrate
setup-migrate: ${MIGRATE_TOOLS}

.PHONY: setup-check
setup-check: ${CHECK_TOOLS}

.PHONY: setup-bench
setup-bench: ${BENCH_TOOLS}

.PHONY: setup-other
setup-other: ${OTHER_TOOLS}

.PHONY: setup
setup: setup-build setup-generate setup-check setup-bench setup-other

.PHONY: generate
generate: setup-build setup-generate
	@echo ">> Generating..."
	@go generate ./...

.PHONY: emit-license-deps
emit-license-deps: setup-build setup-generate
	@go-licenses report \
		--template internal/app/resources/templates/license.tpl \
		./... \
		> LICENSE-backend.md

.PHONY: build
build: setup-build
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
bench: setup-bench
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
check: setup-check
	@echo ">> Running checks and validators..."
	@echo "... staticcheck ..."
	@${GOBIN}/staticcheck ./...
	@echo "... errcheck ..."
	@${GOBIN}/errcheck -ignoretests -exclude .errcheck-excludes.txt ./...
	@echo "... go-vet ..."
	@go vet ./...
	@echo "... nilness ..."
	@${GOBIN}/nilness ./...
	@echo "... ineffassign ..."
	@${GOBIN}/ineffassign ./...
	@echo "... govulncheck ..."
	@${GOBIN}/govulncheck ./...
	@echo "... betteralign ..."
	@${GOBIN}/betteralign ./...
	@echo "... gosec ..."
	@${GOBIN}/gosec -quiet -exclude-generated -exclude-dir=cmd/refidgen -exclude-dir=tools ./...

.PHONY: nilcheck
nilcheck: setup-check
	@echo ">> Running nilcheck (will have some false positives)..."
	@echo "... nilaway ..."
	@${GOBIN}/nilaway -test=false ./...

.PHONY: deadcode
deadcode: setup-check
	@echo ">> Running deadcode (will have some false positives)..."
	@echo "... deadcode ..."
	@${GOBIN}/deadcode -test ./...

.PHONY: update-go-deps
update-go-deps:
	@echo ">> updating Go dependencies..."
	@GOPRIVATE=github.com/dropwhile go get -u all
	@go mod tidy

.PHONY: migrate
migrate: setup-migrate
	@echo ">> running migrations..."
	@goose up

.PHONY: cloc
cloc:
	@echo ">> counting stuff..."
	@cloc -v 2 --force-lang=HTML,gohtml --fullpath --not-match-d resources/static/ .

.PHONY: dev-db-create
dev-db-create:
	@echo ">> starting dev postgres,redis ..."
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
		--name icanbringthat-redis \
		--restart always \
		-p 6379:6379 \
		-d redis:7-alpine \
		redis-server --requirepass "${REDIS_PASS}"

.PHONY: dev-db-start
dev-db-start:
	@echo ">> starting dev postgres,redis ..."
	@docker start icanbringthat-db-init
	@docker start icanbringthat-redis

dev-db-stop:
	@echo ">> stopping dev postgres,redis ..."
	@docker stop icanbringthat-database
	@docker stop icanbringthat-redis

dev-db-purge:
	@echo ">> purging dev postgres,redis ..."
	@docker rm -fv icanbringthat-database
	@docker rm -fv icanbringthat-redis
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
	@exec ./build/bin/server run

.PHONY: devrun
devrun: setup-other
	@echo ">> Monitoring for change, runnging tests, and restarting..."
	@modd -f .modd.conf
