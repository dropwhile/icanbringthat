# build specific vars
APP_VER        := env("APP_VER", "")
GITHASH        := env("GITHASH", "")
BUILD_TARGETS  := env("BUILD_TARGETS", "server client")
BUILD_OPTIONS  := env("BUILD_OPTIONS", "-trimpath")
BUILD_DEPFLAGS := env("BUILD_DEPFLAGS", '-tags netgo,production')
BUILD_LDFLAGS  := env("BUILD_LDFLAGS", '-s -w')

# build binaries
[group('build')]
build: (_banner ">> building binaries")
    #!/bin/sh
    set -euo pipefail
    APP_VER="{{ APP_VER }}"
    if [ "$APP_VER" = "" ]; then
        APP_VER="$(git describe --always --tags|sed 's/^v//')"
    fi
    [ -d "{{ build_dir }}/bin" ] || mkdir -p "{{ build_dir }}/bin"
    echo "> Building as APP_VER=${APP_VER}"
    for x in {{BUILD_TARGETS}}; do
        echo " * ${x} ...";
        go build {{ BUILD_DEPFLAGS }} {{ BUILD_OPTIONS }} \
            -ldflags="{{ BUILD_LDFLAGS }} -X {{ VERSION_VAR }}=${APP_VER}" \
            -o "{{ build_dir }}/bin/${x}" ./cmd/${x};
    done

# build container images
[group('build')]
build-images *FLAGS: (_banner ">> building container images")
    #!/bin/sh
    GITHASH="{{ GITHASH }}"
    if [ "$GITHASH" = "" ]; then
        GITHASH="$(git rev-parse --short HEAD)"
    fi
    APP_VER="{{ APP_VER }}"
    if [ "$APP_VER" = "" ]; then
        APP_VER="$(git describe --always --tags|sed 's/^v//')"
    fi
    if [ "${DOCKER_PREBUILD}" != "" ]; then
        eval ${DOCKER_PREBUILD}
    fi
    set -euo pipefail
    docker build {{FLAGS}} \
        --build-arg GITHASH=${GITHASH} \
        --build-arg APP_VER=${APP_VER} \
        --build-arg BUILDKIT_INLINE_CACHE=1 \
        --cache-from icanbringthat:latest \
        -t icanbringthat:latest \
        -f docker/Dockerfile \
        .
    set +o pipefail
    if [ "${DOCKER_POSTBUILD}" != "" ]; then
        eval ${DOCKER_POSTBUILD}
    fi
