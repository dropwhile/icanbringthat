// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package session

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dropwhile/assert"
)

func TestFlashPopAll(t *testing.T) {
	t.Parallel()

	sm := NewTestSessionManager()
	t.Cleanup(sm.Close)

	// Create a handler and wrap it using sessionManager.LoadAndSave
	handler := sm.LoadAndSave(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		ctx := r.Context()
		if r.Form.Get("read") != "" {
			messages := sm.FlashPopAll(ctx)
			fmt.Fprint(w, messages)
		} else {
			sm.FlashAppend(ctx, "test", "hello")
			sm.FlashAppend(ctx, "test", "goodbye")
			fmt.Fprintln(w, "Hello, client")
		}
	}))

	// Create a test server with the handler
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	jar, err := cookiejar.New(nil)
	assert.Nil(t, err)

	c := http.Client{Timeout: time.Duration(1) * time.Second, Jar: jar}

	// Make a request to the server
	res, err := c.Get(ts.URL)
	if res != nil {
		io.Copy(io.Discard, res.Body)
		res.Body.Close()
	}
	assert.Nil(t, err)

	// Make a request to the server
	res, err = c.Get(ts.URL + "?read=true")
	var b []byte
	if res != nil {
		b, _ = io.ReadAll(res.Body)
		res.Body.Close()
	}
	assert.Nil(t, err)

	assert.Equal(t, string(b), "map[test:[hello goodbye]]")
}

func TestFlashPopOne(t *testing.T) {
	t.Parallel()

	sm := NewTestSessionManager()
	t.Cleanup(sm.Close)

	// Create a handler and wrap it using sessionManager.LoadAndSave
	handler := sm.LoadAndSave(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		ctx := r.Context()
		switch r.Form.Get("read") {
		case "1":
			messages := sm.FlashPopKey(ctx, "test1")
			fmt.Fprint(w, messages)
		case "2":
			messages := sm.FlashPopAll(ctx)
			fmt.Fprint(w, messages)
		default:
			sm.FlashAppend(ctx, "test1", "hello")
			sm.FlashAppend(ctx, "test2", "goodbye")
			fmt.Fprintln(w, "Hello, client")
		}
	}))

	// Create a test server with the handler
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	jar, err := cookiejar.New(nil)
	assert.Nil(t, err)

	c := http.Client{Timeout: time.Duration(1) * time.Second, Jar: jar}

	// Make a request to the server
	res, err := c.Get(ts.URL)
	if res != nil {
		io.Copy(io.Discard, res.Body)
		res.Body.Close()
	}
	assert.Nil(t, err)

	// Make a request to the server
	res, err = c.Get(ts.URL + "?read=1")
	var b []byte
	if res != nil {
		b, _ = io.ReadAll(res.Body)
		res.Body.Close()
	}
	assert.Nil(t, err)
	assert.Equal(t, string(b), "[hello]")

	// Make a request to the server
	res, err = c.Get(ts.URL + "?read=2")
	if res != nil {
		b, _ = io.ReadAll(res.Body)
		res.Body.Close()
	}
	assert.Nil(t, err)
	assert.Equal(t, string(b), "map[test2:[goodbye]]")
}
