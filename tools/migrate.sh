#!/bin/sh

set -e

MIG_DIR="$(dirname "$(dirname "${0}")")/migrations"

case $1 in
    up)
        exec migrate \
            -path "${MIG_DIR}" \
            -database "${DB_DSN}" \
        up
        ;;
    new)
        if [ -z "$2" ]; then
            echo "ERROR: No name provided for migration"
            exit 1
        fi
        exec migrate create \
            -ext sql \
            -dir "${MIG_DIR}" \
            -seq "${2}"
        ;;
    *)
        echo "ERROR: unknown command: '${1}'"
        exit 1
        ;;
esac

