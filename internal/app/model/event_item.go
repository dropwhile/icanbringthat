package model

import (
	"context"
	"time"

	"github.com/dropwhile/icbt/internal/util/refid"
)

type EventItem struct {
	Id           int
	RefId        refid.RefId `db:"ref_id"`
	EventId      int         `db:"event_id"`
	Description  string
	Created      time.Time
	LastModified time.Time `db:"last_modified"`
	Event        *Event    `db:"-"`
	Earmark      *Earmark  `db:"-"`
}

func (ei *EventItem) Insert(ctx context.Context, db PgxHandle) error {
	if ei.RefId.IsNil() {
		ei.RefId = EventItemRefIdT.MustNew()
	}
	q := `INSERT INTO event_item_ (ref_id, event_id, description) VALUES ($1, $2, $3) RETURNING *`
	res, err := QueryOneTx[EventItem](ctx, db, q, ei.RefId, ei.EventId, ei.Description)
	if err != nil {
		return err
	}
	ei.Id = res.Id
	ei.RefId = res.RefId
	ei.Created = res.Created
	ei.LastModified = res.LastModified
	return nil
}

func (ei *EventItem) Save(ctx context.Context, db PgxHandle) error {
	q := `UPDATE event_item_ SET description = $1 WHERE id = $2`
	return ExecTx[EventItem](ctx, db, q, ei.Description, ei.Id)
}

func (ei *EventItem) Delete(ctx context.Context, db PgxHandle) error {
	q := `DELETE FROM event_item_ WHERE id = $1`
	return ExecTx[EventItem](ctx, db, q, ei.Id)
}

func (ei *EventItem) GetEvent(ctx context.Context, db PgxHandle) (*Event, error) {
	event, err := GetEventById(ctx, db, ei.EventId)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func NewEventItem(ctx context.Context, db PgxHandle, eventId int, description string) (*EventItem, error) {
	eventItem := &EventItem{
		EventId:     eventId,
		Description: description,
	}
	err := eventItem.Insert(ctx, db)
	if err != nil {
		return nil, err
	}
	return eventItem, nil
}

func GetEventItemById(ctx context.Context, db PgxHandle, id int) (*EventItem, error) {
	q := `SELECT * FROM event_item_ WHERE id = $1`
	return QueryOne[EventItem](ctx, db, q, id)
}

func GetEventItemByRefId(ctx context.Context, db PgxHandle, refId refid.RefId) (*EventItem, error) {
	q := `SELECT * FROM event_item_ WHERE ref_id = $1`
	return QueryOne[EventItem](ctx, db, q, refId)
}

func GetEventItemsByEvent(ctx context.Context, db PgxHandle, event *Event) ([]*EventItem, error) {
	q := `SELECT * FROM event_item_ WHERE event_id = $1 ORDER BY created DESC`
	return Query[EventItem](ctx, db, q, event.Id)
}
