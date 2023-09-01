// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package session

import (
	"context"
	"encoding/gob"

	"github.com/alexedwards/scs/goredisstore"
	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func init() {
	gob.Register(map[string][]string{})
}

type SessionMgr struct {
	*scs.SessionManager
}

func (sm *SessionMgr) Close() {
	if v, ok := sm.Store.(*pgxstore.PostgresStore); ok {
		v.StopCleanup()
	}
}

func (sm *SessionMgr) GetMap(ctx context.Context, key string) map[string][]string {
	var value map[string][]string
	if v, ok := sm.Get(ctx, key).(map[string][]string); ok {
		value = v
	} else {
		value = make(map[string][]string)
	}
	return value
}

func (sm *SessionMgr) PutMap(ctx context.Context, key string, value map[string][]string) {
	sm.Put(ctx, key, value)
}

func (sm *SessionMgr) PopMap(ctx context.Context, key string) map[string][]string {
	var value map[string][]string
	if v, ok := sm.Pop(ctx, key).(map[string][]string); ok {
		value = v
	}
	return value
}

func NewDBSessionManager(pool *pgxpool.Pool, secure bool) *SessionMgr {
	manager := scs.New()
	manager.Cookie.Secure = secure
	manager.Store = pgxstore.New(pool)
	return &SessionMgr{SessionManager: manager}
}

func NewRedisSessionManager(rdb *redis.Client, secure bool) *SessionMgr {
	manager := scs.New()
	manager.Cookie.Secure = secure
	manager.Store = goredisstore.New(rdb)
	return &SessionMgr{SessionManager: manager}
}

func NewTestSessionManager() *SessionMgr {
	manager := scs.New()
	manager.Cookie.Secure = false
	manager.Store = memstore.NewWithCleanupInterval(0)
	return &SessionMgr{SessionManager: manager}
}
