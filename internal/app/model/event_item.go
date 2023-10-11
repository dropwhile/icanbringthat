package model

import (
	"context"
	"fmt"
	"time"

	"github.com/dropwhile/refid"
)

var EventItemRefIDT = refid.Tagger(3)

type EventItem struct {
	ID           int
	RefID        refid.RefID `db:"ref_id"`
	EventID      int         `db:"event_id"`
	Description  string
	Created      time.Time
	LastModified time.Time `db:"last_modified"`
	Event        *Event    `db:"-"`
	Earmark      *Earmark  `db:"-"`
}

func (ei *EventItem) Insert(ctx context.Context, db PgxHandle) error {
	if ei.RefID.IsNil() {
		ei.RefID = refid.Must(EventItemRefIDT.New())
	}
	q := `INSERT INTO event_item_ (ref_id, event_id, description) VALUES ($1, $2, $3) RETURNING *`
	res, err := QueryOneTx[EventItem](ctx, db, q, ei.RefID, ei.EventID, ei.Description)
	if err != nil {
		return err
	}
	ei.ID = res.ID
	ei.RefID = res.RefID
	ei.Created = res.Created
	ei.LastModified = res.LastModified
	return nil
}

func (ei *EventItem) Save(ctx context.Context, db PgxHandle) error {
	q := `UPDATE event_item_ SET description = $1 WHERE id = $2`
	return ExecTx[EventItem](ctx, db, q, ei.Description, ei.ID)
}

func (ei *EventItem) Delete(ctx context.Context, db PgxHandle) error {
	q := `DELETE FROM event_item_ WHERE id = $1`
	return ExecTx[EventItem](ctx, db, q, ei.ID)
}

func (ei *EventItem) GetEvent(ctx context.Context, db PgxHandle) (*Event, error) {
	event, err := GetEventByID(ctx, db, ei.EventID)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func NewEventItem(
	ctx context.Context,
	db PgxHandle,
	eventID int,
	description string,
) (*EventItem, error) {
	eventItem := &EventItem{
		EventID:     eventID,
		Description: description,
	}
	err := eventItem.Insert(ctx, db)
	if err != nil {
		return nil, err
	}
	return eventItem, nil
}

func GetEventItemByID(ctx context.Context, db PgxHandle, id int) (*EventItem, error) {
	q := `SELECT * FROM event_item_ WHERE id = $1`
	return QueryOne[EventItem](ctx, db, q, id)
}

func GetEventItemByRefID(ctx context.Context, db PgxHandle, refID refid.RefID) (*EventItem, error) {
	if !EventItemRefIDT.HasCorrectTag(refID) {
		err := fmt.Errorf(
			"bad refid type: got %d expected %d",
			refID.Tag(), EventItemRefIDT.Tag(),
		)
		return nil, err
	}
	q := `SELECT * FROM event_item_ WHERE ref_id = $1`
	return QueryOne[EventItem](ctx, db, q, refID)
}

func GetEventItemsByEvent(ctx context.Context, db PgxHandle, event *Event) ([]*EventItem, error) {
	q := `SELECT * FROM event_item_ WHERE event_id = $1 ORDER BY created DESC,id DESC`
	return Query[EventItem](ctx, db, q, event.ID)
}
