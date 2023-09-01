// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/alecthomas/kong"
	"github.com/quic-go/quic-go/http3"
	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icanbringthat/internal/logger"
	"github.com/dropwhile/icanbringthat/internal/util"
	"github.com/dropwhile/icanbringthat/rpc/icbt"
)

type verboseFlag bool

func (v verboseFlag) BeforeApply() error {
	logger.SetLevel(slog.LevelDebug)
	slog.Debug("debug logging enabled")
	return nil
}

type RunArgs struct {
	cli    *CLI
	client icbt.Rpc
	ctx    context.Context
}

type CLI struct { // betteralign:ignore
	// global options
	Verbose     verboseFlag      `name:"verbose" short:"v" help:"enable verbose logging"`
	Quic        bool             `name:"quic" help:"connect with http3/quic"`
	Version     kong.VersionFlag `name:"version" short:"V" help:"Print version information and quit"`
	BaseURL     string           `name:"base-url" short:"b" env:"BASE_URL" required:""`
	TwirpPrefix string           `name:"api-prefix" short:"p" env:"API_PREFIX" default:"/api"`
	AuthToken   string           `name:"auth-token" env:"AUTH_TOKEN" required:""`

	// subcommands
	Events struct { // betteralign:ignore
		Create         EventsCreateCmd       `cmd:"" aliases:"add" help:"create new event"`
		Update         EventsUpdateCmd       `cmd:"" aliases:"update" help:"update event"`
		Delete         EventsDeleteCmd       `cmd:"" aliases:"rm" help:"delete event"`
		List           EventsListCmd         `cmd:"" aliases:"ls" help:"list events"`
		Detail         EventsGetDetailsCmd   `cmd:"" aliases:"info,details" help:"get event details"`
		ListEventItems EventsListItemsCmd    `cmd:"" aliases:"items,ls-items" help:"list event items"`
		ListEarmarks   EventsListEarmarksCmd `cmd:"" aliases:"earmarks,ls-earmarks" help:"list event earmarks"`
	} `cmd:"" help:"events"`

	EventItems struct { // betteralign:ignore
		Add    EventItemsAddCmd    `cmd:"" help:"add item to event"`
		Update EventItemsUpdateCmd `cmd:"" help:"update event item"`
		Remove EventItemsRemoveCmd `cmd:"" aliases:"rm" help:"remove event item"`
	} `cmd:"" help:"event-items"`

	Earmarks struct { // betteralign:ignore
		Create EarmarksCreateCmd     `cmd:"" help:"earmark an item"`
		Detail EarmarksGetDetailsCmd `cmd:"" aliases:"info,details" help:"get earmark details"`
		Remove EarmarksRemoveCmd     `cmd:"" help:"remove an earmark"`
		List   EarmarksListCmd       `cmd:"" help:"list earmarked items"`
	} `cmd:"" help:"earmarks"`

	Favorites struct { // betteralign:ignore
		Add    FavoritesAddCmd    `cmd:"" help:"add favorite"`
		Remove FavoritesRemoveCmd `cmd:"" aliases:"rm" help:"remove favorite"`
		List   FavoritesListCmd   `cmd:"" aliases:"ls" help:"list favorites"`
	} `cmd:"" help:"favorites"`

	Notifications struct { // betteralign:ignore
		Delete    NotificationsDeleteCmd    `cmd:"" aliases:"rm" help:"Delete a single notification."`
		DeleteAll NotificationsDeleteAllCmd `cmd:"" aliases:"clear" help:"Delete all notifications."`
		List      NotificationsListCmd      `cmd:"" aliases:"ls" help:"List notifications."`
	} `cmd:"" help:"notifications"`
}

func main() {
	logger.SetupLogging(logger.NewConsoleLogger, nil)
	vinfo, _ := util.GetVersion()

	cli := CLI{}
	ctx := kong.Parse(&cli,
		kong.Name("icanbringthat client"),
		kong.Description("An api client for icanbringthat"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Tree:         true,
			NoAppSummary: true,
			Compact:      true,
		}),
		kong.Vars{
			"version": vinfo.Version,
		},
	)

	header := http.Header{}
	header.Set("Authorization", fmt.Sprintf("Bearer %s", cli.AuthToken))
	header.Set("User-Agent", fmt.Sprintf("api-client %s", vinfo.Version))

	reqCtx := context.Background()
	reqCtx, err := twirp.WithHTTPRequestHeaders(reqCtx, header)
	if err != nil {
		ctx.FatalIfErrorf(err)
		return
	}

	hc := &http.Client{}
	if cli.Quic {
		hc.Transport = &http3.RoundTripper{}
	}

	client := icbt.NewRpcProtobufClient(
		cli.BaseURL, hc,
		twirp.WithClientPathPrefix(cli.TwirpPrefix),
	)
	err = ctx.Run(&RunArgs{
		cli:    &cli,
		ctx:    reqCtx,
		client: client,
	})
	ctx.FatalIfErrorf(err)
}
