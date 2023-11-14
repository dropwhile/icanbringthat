package resources

import (
	"bytes"
	"fmt"
	htmltemplate "html/template"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	txttemplate "text/template"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog/log"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	"github.com/yuin/goldmark/renderer/html"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/dropwhile/icbt/internal/app/model"
)

const linkTplTxt = `<a class="text-purple-600 dark:text-purple-400 hover:underline" href="{{.href}}">{{.name}}</a>`

var (
	linkReplaceRex = regexp.MustCompile(`\blink:[^.\s]+\b`)
	linkTpl        = htmltemplate.Must(htmltemplate.New("linkTpl").Parse(linkTplTxt))
	linksMap       = map[string]string{
		"/settings": "Account Settings",
	}
)

func linkReplaceFunc(s string) string {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return s
	}

	buf := &bytes.Buffer{}
	if name, ok := linksMap[parts[1]]; ok {
		err := linkTpl.Execute(buf, map[string]any{
			"href": parts[1],
			"name": name,
		})
		if err == nil {
			return buf.String()
		}
	}
	return s
}

func truncateString(s string, size int) string {
	asRunes := []rune(s)
	if len(asRunes) > size {
		asRunes = asRunes[:size]
		if size > 3 {
			asRunes = append(asRunes[:size-3], []rune("...")...)
		}
	}
	return string(asRunes)
}

var templateFuncMap = txttemplate.FuncMap{
	"titlecase": cases.Title(language.English).String,
	"lowercase": func(s fmt.Stringer) string {
		return cases.Lower(language.English).String(s.String())
	},
	"truncate": truncateString,
	"truncate30": func(s string) string {
		return truncateString(s, 30)
	},
	"truncate45": func(s string) string {
		return truncateString(s, 45)
	},
	"truncate60": func(s string) string {
		return truncateString(s, 60)
	},
	"formatTS": func(t time.Time) string {
		return t.UTC().Format("2006-01-02T15:04Z07:00")
	},
	"formatTSLocal": func(t time.Time, zone *model.TimeZone) string {
		return t.In(zone.Location).Format("2006-01-02T15:04")
	},
	"formatDateTime": func(t time.Time) string {
		return t.Format("2006-01-02 15:04 MST")
	},
	"paginate": func(pg *PgInput) *PaginationResult {
		size, step, current := pg.Max, pg.Step, pg.Current
		maxPage := CalculateMaxPageNum(size, step)
		if current < 1 {
			current = 1
		}
		if current > maxPage {
			current = maxPage
		}

		start := (step * (current - 1)) + 1
		if start < 0 {
			start = 0
		}

		stop := start + step - 1
		if size == 0 {
			stop = 0
		}

		if current == maxPage && size%step != 0 {
			stop = ((current - 1) * step) + (size % step)
		}

		s := &PaginationResult{
			Pages:   NewPaginator(3, 3, 3).AddPages(size, step).Paginate(current),
			Start:   start,
			Stop:    stop,
			Size:    size,
			HasPrev: current > 1,
			HasNext: current < maxPage,
		}
		return s
	},
	"mod": func(i, j int) int {
		return i % j
	},
	"add": func(i, j int) int {
		return i + j
	},
	"mult": func(i, j int) int {
		return i * j
	},
	"min": func(i, j int) int {
		return min(i, j)
	},
	"subt": func(i, j int) int {
		return i - j
	},
	"set": func(ac, k, v reflect.Value) error {
		switch ac.Kind() {
		case reflect.Map:
			if k.Type() == ac.Type().Key() {
				ac.SetMapIndex(k, v)
				return nil
			}
		}
		return fmt.Errorf("calling set with unsupported type %q (%T) -> %q (%T)", ac.Kind(), ac, k.Kind(), k)
	},
	// isset is a helper func from hugo
	"isset": func(ac, kv reflect.Value) (bool, error) {
		switch ac.Kind() {
		case reflect.Array, reflect.Slice:
			k := 0
			switch kv.Kind() {
			case reflect.Int | reflect.Int8 | reflect.Int16 | reflect.Int32 | reflect.Int64:
				k = int(kv.Int())
			case reflect.Uint | reflect.Uint8 | reflect.Uint16 | reflect.Uint32 | reflect.Uint64:
				k = int(kv.Uint())
			case reflect.String:
				v, err := strconv.ParseInt(kv.String(), 0, 0)
				if err != nil {
					return false, fmt.Errorf("unable to cast %#v of type %T to int64", kv, kv)
				}
				k = int(v)
			default:
				return false, fmt.Errorf("unable to cast %#v of type %T to int", kv, kv)
			}
			if ac.Len() > k {
				return true, nil
			}
		case reflect.Map:
			if kv.Type() == ac.Type().Key() {
				return ac.MapIndex(kv).IsValid(), nil
			}
		default:
			log.Info().
				Msgf("calling IsSet with unsupported type %q (%T) will always return false", ac.Kind(), ac)
		}

		return false, nil
	},
	"eqorempty": func(arg0, arg1 reflect.Value) (bool, error) {
		k1 := arg0.Kind()
		k2 := arg1.Kind()
		if k1 != k2 {
			return false, fmt.Errorf("non-comparable types %s: %v, %s: %v", arg0, arg0.Type(), arg1.Type(), arg1)
		}

		truth := false
		switch arg0.Kind() {
		case reflect.String:
			truth = arg0.String() == "" || arg0.String() == arg1.String()
		case reflect.Invalid:
			truth = true
		}
		return truth, nil
	},
	"replaceLinks": func(input string) (htmltemplate.HTML, error) {
		body := linkReplaceRex.ReplaceAllStringFunc(input, linkReplaceFunc)
		p := bluemonday.NewPolicy()
		p.AllowElements("p", "br", "strong", "sub", "sup", "em")
		p.AllowElements("b", "i", "pre", "small", "strike", "tt", "u")
		p.RequireParseableURLs(true)
		// allow relative urls in sanitize
		p.AllowRelativeURLs(true)
		p.AllowURLSchemes("https")
		p.RequireNoFollowOnLinks(false)
		p.RequireNoReferrerOnLinks(false)
		p.AllowAttrs("href", "class").OnElements("a")
		out := p.Sanitize(body)
		return htmltemplate.HTML(out), nil // #nosec G203 -- html sanitized by bluemonday
	},
	"markdown": func(input string) (htmltemplate.HTML, error) {
		b := []byte(input)
		var buf bytes.Buffer
		md := goldmark.New(
			goldmark.WithExtensions(
				emoji.Emoji,
			),
			goldmark.WithRendererOptions(
				html.WithHardWraps(),
			),
		)
		if err := md.Convert(b, &buf); err != nil {
			return "", err
		}
		p := bluemonday.NewPolicy()
		p.AllowElements("p", "br", "strong", "sub", "sup", "em")
		p.AllowElements("b", "i", "pre", "small", "strike", "tt", "u")
		p.AllowElements("ul", "ol", "li")
		p.RequireParseableURLs(true)
		// do not allow relative urls in markdown
		p.AllowRelativeURLs(false)
		p.AllowURLSchemes("http", "https")
		p.RequireNoFollowOnLinks(true)
		p.RequireNoReferrerOnLinks(true)
		p.AllowAttrs("href").OnElements("a")
		out := p.SanitizeReader(&buf).String()
		return htmltemplate.HTML(out), nil // #nosec G203 -- html sanitized by bluemonday
	},
}
