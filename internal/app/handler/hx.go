package handler

import (
	"net/http"
	"net/url"
	"strings"
)

type HxUrlVal url.URL

func (hv *HxUrlVal) HasPathPrefix(prefix string) bool {
	u := (*url.URL)(hv)
	return strings.HasPrefix(u.Path, prefix)
}

type Hxx struct {
	http.Header
}

func (hx *Hxx) CurrentUrl() *HxUrlVal {
	h := hx.Header.Get("HX-Current-URL")
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

func Hx(r *http.Request) *Hxx {
	return &Hxx{r.Header}
}
