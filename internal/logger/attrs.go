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
		case slog.SourceKey:
			if len(groups) == 0 {
				a.Key = "src"
				source := a.Value.Any().(*slog.Source)
				source.File = trimFilePath(source.File)
			}
		// Customize the name of the level key and the output string, including
		// custom level values.
		case slog.LevelKey:
			if len(groups) == 0 {
				if _, ok := a.Value.Any().(slog.Level); ok {
					a.Key = "lvl"
				}
			}
		}

		switch a.Value.Kind() {
		// other cases
		case slog.KindAny:
			switch v := a.Value.Any().(type) {
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
