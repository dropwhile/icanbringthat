package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cactus/mlog"
	"github.com/dropwhile/icbt/internal/util/refid"
	flags "github.com/jessevdk/go-flags"
)

type GenerateCommand struct {
	TagValue uint8  `short:"t" long:"tag-value" default:"0" description:"tag value" required:"true"`
	Encoding string `short:"b" long:"base" default:"native" description:"Encode/Decode base. Either native (modified base32), hex, or base64"`
}

// Execute runs the encode command
func (c *GenerateCommand) Execute(args []string) error {
	if opts.Verbose {
		mlog.SetFlags(mlog.Flags() | mlog.Ldebug)
		mlog.Debug("debug logging enabled")
	}
	var refId refid.RefId
	if c.TagValue != 0 {
		refId = refid.MustNewTagged(c.TagValue)
	} else {
		refId = refid.MustNew()
	}

	switch c.Encoding {
	case "base64":
		fmt.Println(refId.ToBase64String())
	case "hex":
		fmt.Println(refId.ToHexString())
	default:
		fmt.Println(refId.String())
	}
	return nil
}

// DecodeCommand holds command options for the decode command
type DecodeCommand struct {
	Positional struct {
		Refid string `positional-arg-name:"refid"`
	} `positional-args:"yes" required:"true"`
}

// Execute runs the decode command
func (c *DecodeCommand) Execute(args []string) error {
	if opts.Verbose {
		mlog.SetFlags(mlog.Flags() | mlog.Ldebug)
		mlog.Debug("debug logging enabled")
	}
	refIdTxt := strings.Trim(c.Positional.Refid, "=")
	refIdTxtLen := len(refIdTxt)

	var refId refid.RefId
	var err error
	switch refIdTxtLen {
	case 0:
		return errors.New("no refid argument provided")
	case 26: // native
		refId, err = refid.Parse(refIdTxt)
		if err != nil {
			return err
		}
	case 32: // hex
		refId, err = refid.FromHexString(refIdTxt)
		if err != nil {
			return err
		}
	case 22: // base64
		refId, err = refid.FromBase64String(refIdTxt)
		if err != nil {
			return err
		}
	}

	ts := refId.Time()
	fmt.Printf("native enc:   %s\n", refId.String())
	fmt.Printf("hex enc:      %s\n", refId.ToHexString())
	fmt.Printf("base64 enc:   %s\n", refId.ToBase64String())
	fmt.Printf("tag value:    %d\n", refId.Tag())
	fmt.Printf("time(string): %s\n", ts.Format(time.RFC3339))
	fmt.Printf("time(micros): %d\n", ts.UnixMicro())

	return nil
}

var opts struct {
	Verbose bool `short:"v" long:"verbose" description:"verbose logging"`
}

// #nosec G104
func main() {
	parser := flags.NewParser(&opts, flags.Default)
	parser.AddCommand("generate", "generate a refid",
		"generate a refid", &GenerateCommand{})
	parser.AddCommand("decode", "Decode a refid",
		"Decode a refid", &DecodeCommand{})

	// parse said flags
	_, err := parser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); ok {
			if e.Type == flags.ErrHelp {
				os.Exit(0)
			}
		}
		os.Exit(1)
	}
}
