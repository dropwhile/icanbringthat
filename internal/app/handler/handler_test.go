package handler

import (
	"time"

	"github.com/dropwhile/icbt/internal/session"
	"github.com/dropwhile/icbt/resources"
	"github.com/pashagolub/pgxmock/v2"
)

var tstTs time.Time

func init() {
	tstTs, _ = time.Parse(time.RFC3339, "2023-01-01T03:04:05Z")
}

func NewTestHandler(mock pgxmock.PgxConnIface) *Handler {
	smgs := session.NewMemorySessionManager()
	tplMap := resources.TemplateMap{}
	h := &Handler{
		Db:      mock,
		Tpl:     tplMap,
		SessMgr: smgs,
	}
	return h
}
