package htmx

import (
	"net/http"
	"net/url"
	"strings"
)

/*
| Header Name                | Description                                                 |
|----------------------------|-------------------------------------------------------------|
| HX-Boosted                 | indicates that the request is via an element using hx-boost |
|----------------------------|-------------------------------------------------------------|
| HX-Current-URL             | the current URL of the browser                              |
|----------------------------|-------------------------------------------------------------|
| HX-History-Restore-Request | true if the request is for history restoration              |
|                            | after a miss in the local history cache                     |
|----------------------------|-------------------------------------------------------------|
| HX-Prompt                  | the user response to an hx-prompt                           |
|----------------------------|-------------------------------------------------------------|
| HX-Request                 | always true                                                 |
|----------------------------|-------------------------------------------------------------|
| HX-Target                  | the id of the target element if it exists                   |
|----------------------------|-------------------------------------------------------------|
| HX-Trigger-Name            | the name of the triggered element if it exists              |
|----------------------------|-------------------------------------------------------------|
| HX-Trigger                 | the id of the triggered element if it exists                |
|----------------------------|-------------------------------------------------------------|
*/

type HxUrlVal url.URL

func (hv *HxUrlVal) HasPathPrefix(prefix string) bool {
	u := (*url.URL)(hv)
	return strings.HasPrefix(u.Path, prefix)
}

type Hxx struct {
	http.Header
}

func (hxx *Hxx) Boosted() bool {
	return hxx.Header.Get("hx-boosted") != ""
}

func (hxx *Hxx) CurrentUrl() *HxUrlVal {
	h := hxx.Header.Get("HX-Current-URL")
	if h != "" {
		u, err := url.Parse(h)
		if err == nil {
			x := HxUrlVal(*u)
			return &x
		}
	}
	x := HxUrlVal(url.URL{})
	return &x
}

func (hxx *Hxx) HistoryRestoreRequest() bool {
	return hxx.Header.Get("hx-history-restore-request") != ""
}

func (hxx *Hxx) Prompt() string {
	return hxx.Header.Get("hx-prompt")
}

func (hxx *Hxx) Request() bool {
	return hxx.Header.Get("hx-request") != ""
}

func (hxx *Hxx) Target() string {
	return hxx.Header.Get("hx-target")
}

func (hxx *Hxx) TriggerName() string {
	return hxx.Header.Get("hx-trigger-name")
}

func (hxx *Hxx) Trigger() string {
	return hxx.Header.Get("hx-trigger")
}

func Hx(r *http.Request) *Hxx {
	if r == nil || r.Header.Get("hx-request") != "true" {
		return &Hxx{}
	}
	return &Hxx{r.Header}
}
