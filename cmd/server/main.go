// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"github.com/alecthomas/kong"

	"github.com/dropwhile/icanbringthat/internal/util"
)

var cli struct { // betteralign:ignore
	Version        kong.VersionFlag `name:"version" short:"V" help:"Print version information and quit"`
	StartWebserver ServerCmd        `cmd:"" help:"run web/api server"`
	StartWorker    WorkerCmd        `cmd:"" help:"run job worker"`
}

func main() {
	vinfo, _ := util.GetVersion()
	cliCtx := kong.Parse(&cli,
		kong.Description("icanbringthat server"),
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
