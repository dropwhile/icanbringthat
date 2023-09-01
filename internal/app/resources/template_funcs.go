// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package resources

import (
	"bytes"
	"fmt"
	"html/template"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	"github.com/yuin/goldmark/renderer/html"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/util"
)

const linkTplTxt = `<a class="text-purple-600 dark:text-purple-400 hover:underline" href="{{.href}}">{{.name}}</a>`

var (
	linkReplaceRex = regexp.MustCompile(`\blink:[^.\s]+\b`)
	linkTpl        = util.Must(template.New("linkTpl").Parse(linkTplTxt))
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

var templateFuncMap = template.FuncMap{
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
	"replaceLinks": func(input string) (template.HTML, error) {
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
		return template.HTML(out), nil // #nosec G203 -- html sanitized by bluemonday
	},
	"markdown": func(input string) (template.HTML, error) {
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
		return template.HTML(out), nil // #nosec G203 -- html sanitized by bluemonday
	},
}
