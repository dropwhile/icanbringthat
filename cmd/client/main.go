package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icbt/rpc"
)

type verboseFlag bool

func (v verboseFlag) BeforeApply() error {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Debug().Msg("debug logging enabled")
	return nil
}

type RunArgs struct {
	cli    *CLI
	client rpc.Rpc
	ctx    context.Context
}

type CLI struct {
	// global options
	Verbose     verboseFlag      `name:"verbose" short:"v" help:"enable verbose logging"`
	Version     kong.VersionFlag `name:"version" short:"V" help:"Print version information and quit"`
	BaseURL     string           `name:"base-url" short:"b" env:"BASE_URL" required:""`
	TwirpPrefix string           `name:"api-prefix" short:"p" env:"API_PREFIX" default:"/api"`
	AuthToken   string           `name:"auth-token" env:"AUTH_TOKEN" required:""`

	// subcommands
	Notifications struct {
		List      NotificationListCmd      `cmd:"" help:"List notifications."`
		Delete    NotificationDeleteCmd    `cmd:"" aliases:"rm" help:"Delete a single notification."`
		DeleteAll NotificationDeleteAllCmd `cmd:"" aliases:"clear" help:"Delete all notifications."`
	} `cmd:"" help:"notifications"`
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:          os.Stderr,
		PartsExclude: []string{zerolog.TimestampFieldName},
	})

	cli := CLI{}
	ctx := kong.Parse(&cli,
		kong.Name("icbt-client"),
		kong.Description("An api client for icbt"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Tree:         true,
			NoAppSummary: true,
			Compact:      true,
		}),
		kong.Vars{
			"version": "0.0.0",
		},
	)

	header := make(http.Header)
	header.Set("Authorization", fmt.Sprintf("Bearer %s", cli.AuthToken))

	reqCtx := context.Background()
	reqCtx, err := twirp.WithHTTPRequestHeaders(reqCtx, header)
	if err != nil {
		ctx.FatalIfErrorf(err)
		return
	}

	client := rpc.NewRpcProtobufClient(
		cli.BaseURL, &http.Client{},
		twirp.WithClientPathPrefix(cli.TwirpPrefix),
	)
	err = ctx.Run(&RunArgs{
		cli:    &cli,
		ctx:    reqCtx,
		client: client,
	})
	ctx.FatalIfErrorf(err)
}
