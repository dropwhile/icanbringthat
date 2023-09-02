package main

import (
	"fmt"
	"os"

	"github.com/cactus/mlog"
	"github.com/dropwhile/icbt/internal/util/refid"
	flags "github.com/jessevdk/go-flags"
)

func main() {
	// command line flags
	var opts struct {
		TagValue uint8 `short:"t" long:"tag-value" default:"0" description:"tag value" required:"true"`
		Verbose  bool  `short:"v" long:"verbose" description:"Show verbose (debug) log level output"`
	}
	// parse said flags
	_, err := flags.Parse(&opts)
	if err != nil {
		if e, ok := err.(*flags.Error); ok {
			if e.Type == flags.ErrHelp {
				os.Exit(0)
			}
		}
		os.Exit(1)
	}
	if opts.Verbose {
		mlog.SetFlags(mlog.Flags() | mlog.Ldebug)
		mlog.Debug("debug logging enabled")
	}

	var refId refid.RefId
	if opts.TagValue != 0 {
		refId = refid.MustNewTagged(opts.TagValue)
	} else {
		refId = refid.MustNew()
	}
	fmt.Println(refId)
}
