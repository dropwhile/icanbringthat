package model

import (
	"context"
	"time"

	"github.com/dropwhile/icbt/internal/util/refid"
)

type EventItem struct {
	Id           uint
	RefId        refid.RefId `db:"ref_id"`
	EventId      uint        `db:"event_id"`
	Description  string
	Event        *Event `db:"-"`
	Created      time.Time
	LastModified time.Time `db:"last_modified"`
}

func (ei *EventItem) Insert(db *DB, ctx context.Context) error {
	if ei.RefId.IsNil() {
		ei.RefId = EventItemRefIdT.MustNew()
	}
	q := `INSERT INTO event_item_ (ref_id, event_id, description) VALUES ($1, $2, $3) RETURNING *`
	res, err := QueryRowTx[EventItem](db, ctx, q, ei.RefId, ei.EventId, ei.Description)
	if err != nil {
		return err
	}
	ei.Id = res.Id
	ei.RefId = res.RefId
	ei.Created = res.Created
	ei.LastModified = res.LastModified
	return nil
}

func (ei *EventItem) Save(db *DB, ctx context.Context) error {
	q := `UPDATE event_item_ SET desciption = $1 WHERE id = $2`
	return ExecTx[EventItem](db, ctx, q, ei.Description, ei.Id)
}

func (ei *EventItem) Delete(db *DB, ctx context.Context) error {
	q := `DELETE FROM event_item_ WHERE id = $1`
	return ExecTx[EventItem](db, ctx, q, ei.Id)
}

func (ei *EventItem) GetEvent(db *DB, ctx context.Context) (*Event, error) {
	event, err := GetEventById(db, ctx, ei.EventId)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func NewEventItem(db *DB, ctx context.Context, eventId uint, description string) (*EventItem, error) {
	eventItem := &EventItem{
		EventId:     eventId,
		Description: description,
	}
	err := eventItem.Insert(db, ctx)
	if err != nil {
		return nil, err
	}
	return eventItem, nil
}

func GetEventItemById(db *DB, ctx context.Context, id uint) (*EventItem, error) {
	q := `SELECT * FROM event_item_ WHERE id = $1`
	return QueryRow[EventItem](db, ctx, q, id)
}

func GetEventItemByRefId(db *DB, ctx context.Context, refId refid.RefId) (*EventItem, error) {
	q := `SELECT * FROM event_item_ WHERE ref_id = $1`
	return QueryRow[EventItem](db, ctx, q, refId)
}

func GetEventItemsByEvent(db *DB, ctx context.Context, event *Event) ([]*EventItem, error) {
	q := `SELECT * FROM event_item_ WHERE event_id = $1`
	return Query[EventItem](db, ctx, q, event.Id)
}
