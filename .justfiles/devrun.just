# monitoring for change, runnging tests, and restarting...
[group('dev:run')]
devrun:
    just _banner ">> starting modd"
    modd -f .modd.conf

# run db migrations
[group('dev:run')]
migrate:
	just _banner ">> running migrations..."
	goose up

# starting dev server...
[group('dev:run')]
run: build
    just _banner ">> starting webserver"
    exec ./build/bin/server start-webserver

# create local development db
[group('dev:run')]
dev-db-create:
    just _banner ">> starting dev postgres,valkey"
    docker run \
        --name icanbringthat-database \
        --restart always \
        -e POSTGRES_PASSWORD=${PGPASSWORD} \
        -e POSTGRES_DB=${PGDATABASE} \
        -p 5432:5432 \
        -v "./database/init/load-extensions.sql:/docker-entrypoint-initdb.d/load-extensions.sql" \
        -d postgres \
        postgres -c jit=off
    docker run \
        --name icanbringthat-valkey \
        --restart always \
        -p 6379:6379 \
        -d valkey/valkey:8-alpine \
        valkey-server --requirepass "${REDIS_PASS}"

# start dev db
[group('dev:run')]
dev-db-start:
    just _banner ">> starting dev postgres,valkey"
    docker start icanbringthat-database
    docker start icanbringthat-valkey

# stop dev db
[group('dev:run')]
dev-db-stop:
    just _banner ">> stopping dev postgres,valkey"
    docker stop icanbringthat-database
    docker stop icanbringthat-valkey

# purge dev db
[group('dev:run')]
dev-db-purge:
    just _banner ">> purging dev postgres,valkey"
    docker rm -fv icanbringthat-database
    docker rm -fv icanbringthat-valkey