# TODOs

*   password reset / forgot password
*   sorting for event items?
*   show earmark details outside of event view table?
    click icon to view more info?
    hover modal?
*   share button on event page?
*   send email
    *   daily notifications about new item bringing
    *   week before and/or day-before notification about event
        (depending on user settings?
    *   password reset / forgot password


## component choices

*   language: go

*   db: postgres ("prod")

## library choices

*   logging: zerolog
    src: github.com/rs/zerolog
    docs: https://godoc.org/github.com/rs/zerolog

*   router: chi
    src: https://github.com/go-chi/chi
    docs: https://go-chi.io/

*   middleware: gorilla csrf
    src: https://github.com/gorilla/csrf

*   sessions: scs
    src: https://github.com/alexedwards/scs
    docs: https://pkg.go.dev/github.com/alexedwards/scs/v2

*   db driver: pgx
    src: https://github.com/jackc/pgx
    docs: https://github.com/jackc/pgx/wiki/Getting-started-with-pgx-through-database-sql#getting-started-with-pgx-through-databasesql

*   db helper lib: scany/pgxscan
    src: https://github.com/georgysavva/scany
    docs: https://pkg.go.dev/github.com/georgysavva/scany/v2
    *   was sqlx
        reasons: ended up being overly verbose and 
            a bit funky. removed quite a few LOC when
            switching to scany
        src: https://github.com/jmoiron/sqlx
        docs: https://jmoiron.github.io/sqlx/

*   db migrations: goose
    note: was golang-migrate, but switched away.
        reasons: single migration file is nicer, transactions support,
        was src: https://github.com/golang-migrate/migrate
    refs: 
        https://pressly.github.io/goose/blog/2022/overview-sql-file
        https://github.com/pressly/goose

*   templates:
    * go html/template

*   environ parsing for settings (12 factor style):
    https://github.com/caarlos0/env
    
## other refs:

*   template fragments:
    https://gist.github.com/benpate/f92b77ea9b3a8503541eb4b9eb515d8a

*   simple go htmx middleware
    https://github.com/donseba/go-htmx

*   organizing db access
    https://www.alexedwards.net/blog/organising-database-access

*   pgx with database/sql
    https://github.com/jackc/pgx/wiki/Getting-started-with-pgx-through-database-sql#getting-started-with-pgx-through-databasesql

*   how to add files to a remove docker volume/mount
    https://stackoverflow.com/questions/51305537/how-can-i-mount-a-volume-of-files-to-a-remote-docker-daemon

*   postgresql triggers
    https://www.the-art-of-web.com/sql/trigger-update-timestamp/

*   html sanitizer
    https://github.com/microcosm-cc/bluemonday

*   argon2 in go:
    https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go

*   markdown for go:
    github.com/yuin/goldmark v1.5.6

*   use t.cleanup instead of defer in tests
    https://brandur.org/fragments/go-prefer-t-cleanup-with-parallel-subtests
    https://github.com/moricho/tparallel/issues/23#issuecomment-1694646461

*   hero icons
    https://heroicons.com/

*   chi routing docs
    https://go-chi.io/

## config notes
| QueryExecMode                | pgbouncer pool_mode | a usable config? |
| ---------------------------- | ------------------- | ---------------- |
| QueryExecModeCacheStatement  | session             | no               |
| QueryExecModeCacheStatement  | transaction         | no               |
| QueryExecModeCacheStatement  | statement           | no               |
| QueryExecModeCacheDescribe   | session             | yes [^1]         |
| QueryExecModeCacheDescribe   | transaction         | yes [^1]         |
| QueryExecModeCacheDescribe   | statement           | no [^2]          |
| QueryExecModeDescribeExec    | session             | yes              |
| QueryExecModeDescribeExec    | transaction         | yes              |
| QueryExecModeDescribeExec    | statement           | no [^2]          |
| QueryExecModeExec            | session             | yes              |
| QueryExecModeExec            | transaction         | yes              |
| QueryExecModeExec            | statement           | no [^2]          |
| QueryExecModeSimpleProtocol  | session             | yes              |
| QueryExecModeSimpleProtocol  | transaction         | yes              |
| QueryExecModeSimpleProtocol  | statement           | no [^2]          |

[^1]: assuming scheme does not change
[^2]: not transaction safe

## possible future changes
*   use memcached?
*   use pgbouncer
