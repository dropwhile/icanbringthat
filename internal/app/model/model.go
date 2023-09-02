package model

import "github.com/dropwhile/icbt/internal/util/refid"

type ModelType interface {
	User | Event | EventItem | Earmark
}

var (
	UserRefIdT      = refid.RefIdTagger(1)
	EventRefIdT     = refid.RefIdTagger(2)
	EventItemRefIdT = refid.RefIdTagger(3)
	EarmarkRefIdT   = refid.RefIdTagger(4)
)
