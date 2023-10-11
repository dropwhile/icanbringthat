package model

import (
	"context"
	"time"

	"github.com/dropwhile/refid"
)

//go:generate go run ../../../cmd/refidgen -t EventItem -v 3

type EventItem struct {
	ID           int
	RefID        EventItemRefID `db:"ref_id"`
	EventID      int            `db:"event_id"`
	Description  string
	Created      time.Time
	LastModified time.Time `db:"last_modified"`
	Event        *Event    `db:"-"`
	Earmark      *Earmark  `db:"-"`
}

func NewEventItem(ctx context.Context, db PgxHandle,
	eventID int, description string,
) (*EventItem, error) {
	refID := refid.Must(NewEventItemRefID())
	return CreateEventItem(ctx, db, refID, eventID, description)
}

func CreateEventItem(ctx context.Context, db PgxHandle,
	refID EventItemRefID, eventID int, description string,
) (*EventItem, error) {
	q := `
		INSERT INTO event_item_ (
			ref_id, event_id, description
		)
		VALUES ($1, $2, $3)
		RETURNING *`
	return QueryOneTx[EventItem](ctx, db, q, refID, eventID, description)
}

func UpdateEventItem(ctx context.Context, db PgxHandle,
	eventItemID int, description string,
) error {
	q := `UPDATE event_item_ SET description = $1 WHERE id = $2`
	return ExecTx[EventItem](ctx, db, q, description, eventItemID)
}

func DeleteEventItem(ctx context.Context, db PgxHandle,
	eventItemID int,
) error {
	q := `DELETE FROM event_item_ WHERE id = $1`
	return ExecTx[EventItem](ctx, db, q, eventItemID)
}

func GetEventItemByID(ctx context.Context, db PgxHandle,
	eventItemID int,
) (*EventItem, error) {
	q := `SELECT * FROM event_item_ WHERE id = $1`
	return QueryOne[EventItem](ctx, db, q, eventItemID)
}

func GetEventItemByRefID(ctx context.Context, db PgxHandle,
	refID EventItemRefID,
) (*EventItem, error) {
	q := `SELECT * FROM event_item_ WHERE ref_id = $1`
	return QueryOne[EventItem](ctx, db, q, refID)
}

func GetEventItemsByEvent(ctx context.Context, db PgxHandle,
	event *Event,
) ([]*EventItem, error) {
	q := `
		SELECT * FROM event_item_
		WHERE event_id = $1
		ORDER BY
			created DESC,
			id DESC`
	return Query[EventItem](ctx, db, q, event.ID)
}
