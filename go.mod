module github.com/dropwhile/icanbringthat

go 1.23.1

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.35.1-20240920164238-5a7b106cbb87.1
	connectrpc.com/connect v1.17.0
	connectrpc.com/validate v0.1.0
	github.com/Masterminds/sprig/v3 v3.3.0
	github.com/alecthomas/kong v1.2.1
	github.com/alexedwards/scs/goredisstore v0.0.0-20240316134038-7e11d57e8885
	github.com/alexedwards/scs/pgxstore v0.0.0-20240316134038-7e11d57e8885
	github.com/alexedwards/scs/v2 v2.8.0
	github.com/caarlos0/env/v11 v11.2.2
	github.com/dropwhile/refid/v2 v2.0.2
	github.com/go-chi/chi/v5 v5.1.0
	github.com/go-chi/httplog/v2 v2.1.1
	github.com/go-playground/validator/v10 v10.22.1
	github.com/go-webauthn/webauthn v0.11.2
	github.com/gorilla/csrf v1.7.2
	github.com/jackc/pgx/v5 v5.7.1
	github.com/jaytaylor/html2text v0.0.0-20230321000545-74c2419ad056
	github.com/microcosm-cc/bluemonday v1.0.27
	github.com/muesli/reflow v0.3.0
	github.com/pashagolub/pgxmock/v4 v4.3.0
	github.com/pganalyze/pg_query_go/v5 v5.1.0
	github.com/pkg/errors v0.9.1
	github.com/quic-go/quic-go v0.47.0
	github.com/redis/go-redis/v9 v9.6.1
	github.com/samber/mo v1.13.0
	github.com/yuin/goldmark v1.7.6
	github.com/yuin/goldmark-emoji v1.0.4
	github.com/zeebo/blake3 v0.2.4
	go.uber.org/mock v0.4.0
	golang.org/x/crypto v0.28.0
	golang.org/x/exp v0.0.0-20241009180824-f66d83c29e7c
	google.golang.org/protobuf v1.35.1
	gotest.tools/v3 v3.5.1
)

require (
	dario.cat/mergo v1.0.1 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.3.0 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/bufbuild/protovalidate-go v0.7.2 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.5 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/go-webauthn/x v0.1.15 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/google/cel-go v0.21.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/go-tpm v0.9.1 // indirect
	github.com/google/pprof v0.0.0-20241009165004-a3522334989c // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/huandu/xstrings v1.5.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/onsi/ginkgo/v2 v2.20.2 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/quic-go/qpack v0.5.1 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/spf13/cast v1.7.0 // indirect
	github.com/ssor/bom v0.0.0-20170718123548-6386211fdfcf // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/mod v0.21.0 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	golang.org/x/tools v0.26.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241007155032-5fefd90f89a9 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241007155032-5fefd90f89a9 // indirect
)

replace github.com/imdario/mergo => github.com/imdario/mergo v0.3.16
