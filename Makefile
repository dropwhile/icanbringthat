# environment
BUILDDIR            := ${CURDIR}/build
ARCH                := $(shell go env GOHOSTARCH)
OS                  := $(shell go env GOHOSTOS)
GOVER               := $(shell go version | awk '{print $$3}' | tr -d '.')

# app specific info
APP_VER             ?= v$(shell git describe --always --tags|sed 's/^v//')
GITHASH             ?= $(shell git rev-parse --short HEAD)
GOPATH              := $(shell go env GOPATH)
GOBIN               := ${GOPATH}/bin
VERSION_VAR         := main.Version
DB_DSN              ?= "postgres://postgres:password@127.0.0.1:5432/icbt?sslmode=disable"
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
PGDATABASE          ?= icbt
PGPASSWORD          ?= password
REDIS_PASS          ?= password

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
  cover               run tests with cover output
  generate            run go:generate
  build               build all binaries
  test                run tests
  bench               run benchmarks
  nilcheck            run nilcheck; noisy/false positives, so not enabled by default
  update-go-deps      updates go.mod and go.sum files
  migrate             runs db migrations
  cloc                counts lines of code
  dev-db-create       creates a docker postgres for development
  dev-db-start        starts a previously created docker postgres for development
  dev-db-stop         stops a previously created docker postgres for development
  dev-db-purge        deletes/destroys a previously created docker postgres for development
  docker-build        build a deployable docker image
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

${GOBIN}/benchstat:
	go install golang.org/x/perf/cmd/benchstat@latest

${GOBIN}/stringer:
	go install golang.org/x/tools/cmd/stringer@latest

${GOBIN}/staticcheck:
	go install honnef.co/go/tools/cmd/staticcheck@latest

${GOBIN}/gosec:
	go install github.com/securego/gosec/v2/cmd/gosec@latest

${GOBIN}/govulncheck:
	go install golang.org/x/vuln/cmd/govulncheck@latest

${GOBIN}/errcheck:
	go install github.com/kisielk/errcheck@latest

${GOBIN}/ineffassign:
	go install github.com/gordonklaus/ineffassign@latest

${GOBIN}/nilaway:
	go install go.uber.org/nilaway/cmd/nilaway@latest

${GOBIN}/protoc-gen-twirp:
	go install github.com/twitchtv/twirp/protoc-gen-twirp@latest

${GOBIN}/protoc-gen-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

${GOBIN}/go-errorlint:
	go install github.com/polyfloyd/go-errorlint@latest

${GOBIN}/modd:
	go install github.com/cortesi/modd/cmd/modd@latest

${GOBIN}/convergen:
	go install github.com/reedom/convergen@latest

${GOBIN}/deadcode:
	go install golang.org/x/tools/cmd/deadcode@latest

${GOBIN}/betteralign:
	go install github.com/dkorunic/betteralign/cmd/betteralign@latest

${GOBIN}/go-licenses:
	go install github.com/google/go-licenses@latest

${GOBIN}/goose:
	go install github.com/pressly/goose/v3/cmd/goose@latest

BENCH_TOOLS := ${GOBIN}/benchstat
OTHER_TOOLS := ${GOBIN}/modd
GENERATE_TOOLS := ${GOBIN}/stringer ${GOBIN}/protoc-gen-twirp ${GOBIN}/protoc-gen-go
GENERATE_TOOLS += ${GOBIN}/convergen ${GOBIN}/go-licenses
CHECK_TOOLS := ${GOBIN}/staticcheck ${GOBIN}/gosec ${GOBIN}/govulncheck
CHECK_TOOLS += ${GOBIN}/errcheck ${GOBIN}/ineffassign ${GOBIN}/nilaway
CHECK_TOOLS += ${GOBIN}/go-errorlint ${GOBIN}/ineffassign ${GOBIN}/deadcode
CHECK_TOOLS += ${GOBIN}/betteralign
MIGRATE_TOOLS += ${GOBIN}/goose

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
	@grep -lRE '^// Code generated by (.+) DO NOT EDIT' | xargs rm -v

.PHONY: check
check: setup-check
	@echo ">> Running checks and validators..."
	@echo "... staticcheck ..."
	@${GOBIN}/staticcheck ./...
	@echo "... errcheck ..."
	@${GOBIN}/errcheck -ignoretests -exclude .errcheck-excludes.txt ./...
	@echo "... go-vet ..."
	@go vet ./...
	@echo "... gosec ..."
	@${GOBIN}/gosec -quiet -exclude-generated -exclude-dir=cmd/refidgen -exclude-dir=tools ./...
	@echo "... ineffassign ..."
	@${GOBIN}/ineffassign ./...
	@echo "... govulncheck ..."
	@${GOBIN}/govulncheck ./...

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

.PHONY: betteralign
betteralign: setup-check
	@echo ">> Running betteralign (will have some false positives)..."
	@echo "... betteralign ..."
	@${GOBIN}/betteralign ./...

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
	@docker volume rm -f icbt-db-init
	@docker volume create icbt-db-init
	@docker create -v icbt-db-init:/data --name icbt-db-helper busybox true
	@for f in  ./database/init/*; do docker cp -q "$${f}" icbt-db-helper:/data; done
	@docker rm -f icbt-db-helper
	@docker run \
		--name icbt-database \
		--restart always \
		-e POSTGRES_PASSWORD=${PGPASSWORD} \
		-e POSTGRES_DB=${PGDATABASE} \
		-p 5432:5432 \
		-v "icbt-db-init:/docker-entrypoint-initdb.d/" \
		-d postgres \
		postgres -c jit=off
	@docker run \
		--name icbt-redis \
		--restart always \
		-p 6379:6379 \
		-d redis:7-alpine \
		redis-server --requirepass "${REDIS_PASS}"

.PHONY: dev-db-start
dev-db-start:
	@echo ">> starting dev postgres,redis ..."
	@docker start icbt-db-init
	@docker start icbt-redis

dev-db-stop:
	@echo ">> stopping dev postgres,redis ..."
	@docker stop icbt-database
	@docker stop icbt-redis

dev-db-purge:
	@echo ">> purging dev postgres,redis ..."
	@docker rm -fv icbt-database
	@docker rm -fv icbt-redis
	@docker volume rm -f icbt-db-init

.PHONY: docker-build
docker-build:
	@echo ">> Building docker image..."
	@${DOCKER_PREBUILD}
	@DOCKER_BUILDKIT=1 docker build \
		--build-arg GITHASH=${GITHASH} \
		--build-arg APP_VER=${APP_VER} \
		--build-arg BUILDKIT_INLINE_CACHE=1 \
		--cache-from icbt:latest \
		-t icbt:latest \
		-f docker/Dockerfile \
		.
	@eval ${DOCKER_POSTBUILD}

.PHONY: run
run: build
	@echo ">> starting dev server..."
	@exec ./build/bin/server

.PHONY: devrun
devrun: setup-other
	@echo ">> Monitoring for change, runnging tests, and restarting..."
	@modd -f .modd.conf

.PHONY: all
all: build
