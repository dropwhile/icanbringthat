// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package logger

import (
	"log/slog"
	"regexp"
	"strings"
	"time"
)

var versionyRE = regexp.MustCompile(`^v[0-9]+@`)

func trimFilePath(file string) string {
	short := file
	counter := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			if counter > 0 {
				short = file[i+1:]
				if versionyRE.MatchString(short) {
					continue
				}
				break
			}
			counter += 1
		}
	}

	// prune from after @ to next /
	atIdx := strings.Index(short, "@")
	if atIdx >= 0 && atIdx+7 <= len(short) {
		for i := atIdx; i < len(short); i++ {
			if short[i] == '/' {
				if i-atIdx > 7 {
					short = short[:atIdx+7] + "..." + short[i:]
				}
				break
			}
		}
	}

	return short
}

func replaceAttr(opts Options) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		// Remove time if in test mode
		if opts.OmitTime && a.Key == slog.TimeKey && len(groups) == 0 {
			return slog.Attr{}
		}

		switch a.Key {
		// change time keyname, as well as use UTC when
		// the option is not set to localtime
		case slog.TimeKey:
			if len(groups) == 0 {
				if v, ok := a.Value.Any().(time.Time); ok {
					a.Key = "ts"
					if !opts.UseLocalTime {
						a.Value = slog.TimeValue(v.UTC())
					}
				}
			}
		// Remove the directory from the source's filename.
		// remove function name as well
		case slog.SourceKey:
			if len(groups) == 0 {
				a.Key = "src"
				source := a.Value.Any().(*slog.Source)
				source.File = trimFilePath(source.File)
				source.Function = ""
			}
		// Customize the name of the level key and the output string, including
		// custom level values.
		case slog.LevelKey:
			if len(groups) == 0 {
				if v, ok := a.Value.Any().(slog.Level); ok {
					switch v {
					case slog.LevelDebug:
						a.Value = slog.StringValue("DBG")
					case slog.LevelInfo:
						a.Value = slog.StringValue("INF")
					case slog.LevelWarn:
						a.Value = slog.StringValue("WRN")
					case slog.LevelError:
						a.Value = slog.StringValue("ERR")
					default:
						a.Value = slog.StringValue(v.String())
					}
				}
			}
		}

		switch a.Value.Kind() {
		case slog.KindAny:
			switch v := a.Value.Any().(type) {
			// some additiona formatting for error traces?
			case error:
				a.Value = fmtErr(v)
			}
		}

		return a
	}
}

func Err(err error) slog.Attr {
	return slog.Any("error", err)
}
