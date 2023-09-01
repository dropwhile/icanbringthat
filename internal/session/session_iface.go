// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package session

import (
	"context"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

// SessionManager works with sessions
// generated from "github.com/alexedwards/scs/v2"
type SessionManager interface {
	// Load retrieves the session data for the given token from the session store,
	// and returns a new context.Context containing the session data. If no matching
	// token is found then this will create a new session.
	//
	// Most applications will use the LoadAndSave() middleware and will not need to
	// use this method.
	Load(ctx context.Context, token string) (context.Context, error)
	// Commit saves the session data to the session store and returns the session
	// token and expiry time.
	//
	// Most applications will use the LoadAndSave() middleware and will not need to
	// use this method.
	Commit(ctx context.Context) (string, time.Time, error)
	// Destroy deletes the session data from the session store and sets the session
	// status to Destroyed. Any further operations in the same request cycle will
	// result in a new session being created.
	Destroy(ctx context.Context) error
	// Put adds a key and corresponding value to the session data. Any existing
	// value for the key will be replaced. The session data status will be set to
	// Modified.
	Put(ctx context.Context, key string, val interface{})
	// Get returns the value for a given key from the session data. The return
	// value has the type interface{} so will usually need to be type asserted
	// before you can use it. For example:
	//
	//	foo, ok := session.Get(r, "foo").(string)
	//	if !ok {
	//		return errors.New("type assertion to string failed")
	//	}
	//
	// Also see the GetString(), GetInt(), GetBytes() and other helper methods which
	// wrap the type conversion for common types.
	Get(ctx context.Context, key string) interface{}
	// Pop acts like a one-time Get. It returns the value for a given key from the
	// session data and deletes the key and value from the session data. The
	// session data status will be set to Modified. The return value has the type
	// interface{} so will usually need to be type asserted before you can use it.
	Pop(ctx context.Context, key string) interface{}
	// Remove deletes the given key and corresponding value from the session data.
	// The session data status will be set to Modified. If the key is not present
	// this operation is a no-op.
	Remove(ctx context.Context, key string)
	// Clear removes all data for the current session. The session token and
	// lifetime are unaffected. If there is no data in the current session this is
	// a no-op.
	Clear(ctx context.Context) error
	// Exists returns true if the given key is present in the session data.
	Exists(ctx context.Context, key string) bool
	// Keys returns a slice of all key names present in the session data, sorted
	// alphabetically. If the data contains no data then an empty slice will be
	// returned.
	Keys(ctx context.Context) []string
	// RenewToken updates the session data to have a new session token while
	// retaining the current session data. The session lifetime is also reset and
	// the session data status will be set to Modified.
	//
	// The old session token and accompanying data are deleted from the session store.
	//
	// To mitigate the risk of session fixation attacks, it's important that you call
	// RenewToken before making any changes to privilege levels (e.g. login and
	// logout operations). See https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/Session_Management_Cheat_Sheet.md#renew-the-session-id-after-any-privilege-level-change
	// for additional information.
	RenewToken(ctx context.Context) error
	// MergeSession is used to merge in data from a different session in case strict
	// session tokens are lost across an oauth or similar redirect flows. Use Clear()
	// if no values of the new session are to be used.
	MergeSession(ctx context.Context, token string) error
	// Status returns the current status of the session data.
	Status(ctx context.Context) scs.Status
	// GetString returns the string value for a given key from the session data.
	// The zero value for a string ("") is returned if the key does not exist or the
	// value could not be type asserted to a string.
	GetString(ctx context.Context, key string) string
	// GetBool returns the bool value for a given key from the session data. The
	// zero value for a bool (false) is returned if the key does not exist or the
	// value could not be type asserted to a bool.
	GetBool(ctx context.Context, key string) bool
	// GetInt returns the int value for a given key from the session data. The
	// zero value for an int (0) is returned if the key does not exist or the
	// value could not be type asserted to an int.
	GetInt(ctx context.Context, key string) int
	// GetInt64 returns the int64 value for a given key from the session data. The
	// zero value for an int64 (0) is returned if the key does not exist or the
	// value could not be type asserted to an int64.
	GetInt64(ctx context.Context, key string) int64
	// GetInt32 returns the int value for a given key from the session data. The
	// zero value for an int32 (0) is returned if the key does not exist or the
	// value could not be type asserted to an int32.
	GetInt32(ctx context.Context, key string) int32
	// GetFloat returns the float64 value for a given key from the session data. The
	// zero value for an float64 (0) is returned if the key does not exist or the
	// value could not be type asserted to a float64.
	GetFloat(ctx context.Context, key string) float64
	// GetBytes returns the byte slice ([]byte) value for a given key from the session
	// data. The zero value for a slice (nil) is returned if the key does not exist
	// or could not be type asserted to []byte.
	GetBytes(ctx context.Context, key string) []byte
	// GetTime returns the time.Time value for a given key from the session data. The
	// zero value for a time.Time object is returned if the key does not exist or the
	// value could not be type asserted to a time.Time. This can be tested with the
	// time.IsZero() method.
	GetTime(ctx context.Context, key string) time.Time
	// PopString returns the string value for a given key and then deletes it from the
	// session data. The session data status will be set to Modified. The zero
	// value for a string ("") is returned if the key does not exist or the value
	// could not be type asserted to a string.
	PopString(ctx context.Context, key string) string
	// PopBool returns the bool value for a given key and then deletes it from the
	// session data. The session data status will be set to Modified. The zero
	// value for a bool (false) is returned if the key does not exist or the value
	// could not be type asserted to a bool.
	PopBool(ctx context.Context, key string) bool
	// PopInt returns the int value for a given key and then deletes it from the
	// session data. The session data status will be set to Modified. The zero
	// value for an int (0) is returned if the key does not exist or the value could
	// not be type asserted to an int.
	PopInt(ctx context.Context, key string) int
	// PopFloat returns the float64 value for a given key and then deletes it from the
	// session data. The session data status will be set to Modified. The zero
	// value for an float64 (0) is returned if the key does not exist or the value
	// could not be type asserted to a float64.
	PopFloat(ctx context.Context, key string) float64
	// PopBytes returns the byte slice ([]byte) value for a given key and then
	// deletes it from the from the session data. The session data status will be
	// set to Modified. The zero value for a slice (nil) is returned if the key does
	// not exist or could not be type asserted to []byte.
	PopBytes(ctx context.Context, key string) []byte
	// PopTime returns the time.Time value for a given key and then deletes it from
	// the session data. The session data status will be set to Modified. The zero
	// value for a time.Time object is returned if the key does not exist or the
	// value could not be type asserted to a time.Time.
	PopTime(ctx context.Context, key string) time.Time
	// RememberMe controls whether the session cookie is persistent (i.e  whether it
	// is retained after a user closes their browser). RememberMe only has an effect
	// if you have set SessionManager.Cookie.Persist = false (the default is true) and
	// you are using the standard LoadAndSave() middleware.
	RememberMe(ctx context.Context, val bool)
	// Iterate retrieves all active (i.e. not expired) sessions from the store and
	// executes the provided function fn for each session. If the session store
	// being used does not support iteration then Iterate will panic.
	Iterate(ctx context.Context, fn func(context.Context) error) error
	// Deadline returns the 'absolute' expiry time for the session. Please note
	// that if you are using an idle timeout, it is possible that a session will
	// expire due to non-use before the returned deadline.
	Deadline(ctx context.Context) time.Time
	// SetDeadline updates the 'absolute' expiry time for the session. Please note
	// that if you are using an idle timeout, it is possible that a session will
	// expire due to non-use before the set deadline.
	SetDeadline(ctx context.Context, expire time.Time)
	// Token returns the session token. Please note that this will return the
	// empty string "" if it is called before the session has been committed to
	// the store.
	Token(ctx context.Context) string
	// LoadAndSave provides middleware which automatically loads and saves session
	// data for the current request, and communicates the session token to and from
	// the client in a cookie.
	LoadAndSave(next http.Handler) http.Handler
	// WriteSessionCookie writes a cookie to the HTTP response with the provided
	// token as the cookie value and expiry as the cookie expiry time. The expiry
	// time will be included in the cookie only if the session is set to persist
	// or has had RememberMe(true) called on it. If expiry is an empty time.Time
	// struct (so that it's IsZero() method returns true) the cookie will be
	// marked with a historical expiry time and negative max-age (so the browser
	// deletes it).
	//
	// Most applications will use the LoadAndSave() middleware and will not need to
	// use this method.
	WriteSessionCookie(ctx context.Context, w http.ResponseWriter, token string, expiry time.Time)
	// wrapper methods not present in original
	GetMap(ctx context.Context, key string) map[string][]string
	PutMap(ctx context.Context, key string, value map[string][]string)
	PopMap(ctx context.Context, key string) map[string][]string
	FlashAppend(ctx context.Context, key string, val ...string)
	FlashPopAll(ctx context.Context) map[string][]string
	FlashPopKey(ctx context.Context, key string) []string
	Close()
}
