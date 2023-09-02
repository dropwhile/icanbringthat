package session

import (
	"context"
	"database/sql"
	"encoding/gob"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
)

func init() {
	gob.Register(map[string][]string{})
}

type SessionMgr struct {
	*scs.SessionManager
}

func (sm *SessionMgr) Close() {
	if v, ok := sm.Store.(*postgresstore.PostgresStore); ok {
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

func (sm *SessionMgr) GetUint(ctx context.Context, key string) uint {
	var value uint = 0
	if v, ok := sm.Get(ctx, key).(uint); ok {
		value = v
	}
	return value
}

func (sm *SessionMgr) PutUint(ctx context.Context, key string, value uint) {
	sm.Put(ctx, key, value)
}

func (sm *SessionMgr) PopUint(ctx context.Context, key string) uint {
	var value uint = 0
	if v, ok := sm.Pop(ctx, key).(uint); ok {
		value = v
	}
	return value
}

func NewDBSessionManager(db *sql.DB) *SessionMgr {
	manager := scs.New()
	store := postgresstore.New(db)
	manager.Store = store

	return &SessionMgr{SessionManager: manager}
}

func NewMemorySessionManager() *SessionMgr {
	manager := scs.New()
	store := memstore.New()
	manager.Store = store

	return &SessionMgr{SessionManager: manager}
}
