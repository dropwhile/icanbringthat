package session

import (
	"context"
)

const flashKey = "_flashes"

func (sm *SessionMgr) FlashAppend(ctx context.Context, key string, val ...string) {
	flashes := sm.GetMap(ctx, flashKey)
	if len(val) > 0 {
		if _, ok := flashes[key]; ok {
			flashes[key] = append(flashes[key], val...)
		} else {
			flashes[key] = val
		}
	}
	sm.PutMap(ctx, flashKey, flashes)
}

func (sm *SessionMgr) FlashPopAll(ctx context.Context) map[string][]string {
	return sm.PopMap(ctx, flashKey)
}

func (sm *SessionMgr) FlashPopKey(ctx context.Context, key string) []string {
	flashes := sm.GetMap(ctx, flashKey)
	var result []string
	if v, ok := flashes[key]; ok {
		result = v
		delete(flashes, key)
		sm.PutMap(ctx, flashKey, flashes)
	} else {
		result = make([]string, 0)
	}
	return result
}
