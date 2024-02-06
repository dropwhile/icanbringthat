package main

import (
	"github.com/alecthomas/kong"

	"github.com/dropwhile/icbt/internal/util"
)

var cli struct {
	Version kong.VersionFlag `name:"version" short:"V" help:"Print version information and quit"`
	Run     RunCmd           `cmd:"" help:"run server"`
}

func main() {
	vinfo, _ := util.GetVersion()
	cliCtx := kong.Parse(&cli,
		kong.Description("icbt server"),
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
	err := cliCtx.Run()
	cliCtx.FatalIfErrorf(err)
}
