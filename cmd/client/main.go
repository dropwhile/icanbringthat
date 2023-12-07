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

// Version holds the server version string
var Version = "no-version"

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
	Events struct {
		Create         EventsCreateCmd       `cmd:"" aliases:"add" help:"create new event"`
		Update         EventsUpdateCmd       `cmd:"" aliases:"update" help:"update event"`
		Delete         EventsDeleteCmd       `cmd:"" aliases:"rm" help:"delete event"`
		List           EventsListCmd         `cmd:"" aliases:"ls" help:"list events"`
		Detail         EventsGetDetailsCmd   `cmd:"" aliases:"info,details" help:"get event details"`
		ListEventItems EventsListItemsCmd    `cmd:"" aliases:"items,ls-items" help:"list event items"`
		ListEarmarks   EventsListEarmarksCmd `cmd:"" aliases:"earmarks,ls-earmarks" help:"list event earmarks"`
	} `cmd:"" help:"events"`

	EventItems struct {
		Add    EventItemsAddCmd    `cmd:"" help:"add item to event"`
		Update EventItemsUpdateCmd `cmd:"" help:"update event item"`
		Remove EventItemsRemoveCmd `cmd:"" aliases:"rm" help:"remove event item"`
	} `cmd:"" help:"event-items"`

	Earmarks struct {
		Create EarmarksCreateCmd     `cmd:"" help:"earmark an item"`
		Detail EarmarksGetDetailsCmd `cmd:"" aliases:"info,details" help:"get earmark details"`
		Remove EarmarksRemoveCmd     `cmd:"" help:"remove an earmark"`
		List   EarmarksListCmd       `cmd:"" help:"list earmarked items"`
	} `cmd:"" help:"earmarks"`

	Favorites struct {
		Add    FavoritesAddCmd    `cmd:"" help:"add favorite"`
		Remove FavoritesRemoveCmd `cmd:"" aliases:"rm" help:"remove favorite"`
		List   FavoritesListCmd   `cmd:"" aliases:"ls" help:"list favorites"`
	} `cmd:"" help:"favorites"`

	Notifications struct {
		Delete    NotificationsDeleteCmd    `cmd:"" aliases:"rm" help:"Delete a single notification."`
		DeleteAll NotificationsDeleteAllCmd `cmd:"" aliases:"clear" help:"Delete all notifications."`
		List      NotificationsListCmd      `cmd:"" aliases:"ls" help:"List notifications."`
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
			"version": Version,
		},
	)

	header := http.Header{}
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
