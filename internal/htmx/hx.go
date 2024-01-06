package htmx

import (
	"net/http"
	"net/url"
	"strings"
)

/* Request Headers
|----------------------------|-------------------------------------------------------------|
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

/* Response Headers
|----------------------------|-------------------------------------------------------------|
| Header Name                | Description                                                 |
|----------------------------|-------------------------------------------------------------|
| HX-Location                | allows you to do a client-side redirect that does           |
|                            | not do a full page reload                                   |
|----------------------------|-------------------------------------------------------------|
| HX-Push-Url                | pushes a new url into the history stack                     |
|----------------------------|-------------------------------------------------------------|
| HX-Redirect                | can be used to do a client-side redirect to a new location  |
|----------------------------|-------------------------------------------------------------|
| HX-Refresh                 | if set to “true” the client-side will do a full             |
|                            | refresh of the page                                         |
|----------------------------|-------------------------------------------------------------|
| HX-Replace-Url             | replaces the current URL in the location bar                |
|----------------------------|-------------------------------------------------------------|
| HX-Reswap                  | allows you to specify how the response will be swapped.     |
|                            | See hx-swap for possible values                             |
|----------------------------|-------------------------------------------------------------|
| HX-Retarget                | a CSS selector that updates the target of the content       |
|                            | update to a different element on the page                   |
|----------------------------|-------------------------------------------------------------|
| HX-Reselect                | a CSS selector that allows you to choose which part of the  |
|                            | response is used to be swapped in. Overrides an existing    |
|                            | hx-select on the triggering element                         |
|----------------------------|-------------------------------------------------------------|
| HX-Trigger                 | allows you to trigger client-side events                    |
|----------------------------|-------------------------------------------------------------|
| HX-Trigger-After-Settle    | allows you to trigger client-side events after the settle   |
|                            | step                                                        |
|----------------------------|-------------------------------------------------------------|
| HX-Trigger-After-Swap      | allows you to trigger client-side events after the swap     |
|                            | step                                                        |
|----------------------------|-------------------------------------------------------------|
*/

const (
	/* request headers */

	// indicates that the request is via an element using hx-boost
	HxBoosted = "HX-Boosted"

	// the current URL of the browser
	HxCurrentURL = "HX-Current-URL"

	// true if the request is for history restoration after a miss in the local
	// history cache
	HxHistoryRestoreRequest = "HX-History-Restore-Request"

	// the user response to an hx-prompt
	HxPrompt = "HX-Prompt"

	// always true if an htmx request
	HxRequest = "HX-Request"

	// the id of the target element if it exists
	HxTarget = "HX-Target"

	// the name of the triggered element if it exists
	HxTriggerName = "HX-Trigger-Name"

	// the id of the triggered element if it exists
	HxTrigger = "HX-Trigger"

	/* response headers */

	// allows you to do a client-side redirect that does not do a full page reload
	HxLocation = "HX-Location"

	// pushes a new url into the history stack
	HxPushUrl = "HX-Push-Url"

	// can be used to do a client-side redirect to a new location
	HxRedirect = "HX-Redirect"

	// if set to “true” the client-side will do a full refresh of the page
	HxRefresh = "HX-Refresh"

	// replaces the current URL in the location bar
	HxReplaceUrl = "HX-Replace-Url"

	// allows you to specify how the response will be swapped. See hx-swap for possible values
	HxReswap = "HX-Reswap"

	// a CSS selector that updates the target of the content update to a different element on the page
	HxRetarget = "HX-Retarget"

	// a CSS selector that allows you to choose which part of the response is used to be swapped in. Overrides an existing hx-select on the triggering element
	HxReselect = "HX-Reselect"

	// allows you to trigger client-side events
	// HxTrigger = "HX-Trigger" // note: repeated header, use previous definition for both Req & Resp

	// allows you to trigger client-side events after the settle step
	HxTriggerAfterSettle = "HX-Trigger-After-Settle"

	// allows you to trigger client-side events after the swap step
	HxTriggerAfterSwap = "HX-Trigger-After-Swap"
)

type HxUrlVal url.URL

func (hv *HxUrlVal) HasPathPrefix(prefix string) bool {
	u := (*url.URL)(hv)
	return strings.HasPrefix(u.Path, prefix)
}

type Hxr struct {
	http.Header
}

func (hxr *Hxr) Boosted() bool {
	return hxr.Header.Get(HxBoosted) != ""
}

func (hxr *Hxr) CurrentUrl() *HxUrlVal {
	h := hxr.Header.Get(HxCurrentURL)
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

func (hxr *Hxr) HistoryRestoreRequest() bool {
	return hxr.Header.Get(HxHistoryRestoreRequest) != ""
}

func (hxr *Hxr) Prompt() string {
	return hxr.Header.Get(HxPrompt)
}

func (hxr *Hxr) IsRequest() bool {
	return hxr.Header.Get(HxRequest) != ""
}

func (hxr *Hxr) Target() string {
	return hxr.Header.Get(HxTarget)
}

func (hxr *Hxr) TriggerName() string {
	return hxr.Header.Get(HxTriggerName)
}

func (hxr *Hxr) Trigger() string {
	return hxr.Header.Get(HxTrigger)
}

type Hxw struct {
	http.ResponseWriter
}

func (hxw *Hxw) HxLocation(location string) {
	hxw.Header().Add(HxLocation, location)
}

func (hxw *Hxw) HxTriggerAfterSwap(triggerName string) {
	hxw.Header().Add(HxTriggerAfterSwap, triggerName)
}

func (hxw *Hxw) HxRedirect(target string) {
	hxw.Header().Add(HxRedirect, target)
}

func (hxw *Hxw) HxRefesh() {
	hxw.Header().Add(HxRefresh, "true")
}

func Request(r *http.Request) *Hxr {
	if r == nil || r.Header.Get(HxRequest) != "true" {
		return &Hxr{}
	}
	return &Hxr{r.Header}
}

func Response(w http.ResponseWriter) *Hxw {
	return &Hxw{w}
}
