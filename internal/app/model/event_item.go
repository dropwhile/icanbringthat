package model

import (
	"context"
	"time"

	"github.com/dropwhile/refid/v2"
	"github.com/dropwhile/refid/v2/reftag"
	"github.com/jackc/pgx/v5"
)

type EventItemRefID = reftag.IDt3

type EventItemRefIDNull struct {
	reftag.NullIDt3
}

var NewEventItemRefID = reftag.New[EventItemRefID]

type EventItem struct {
	Created      time.Time
	LastModified time.Time `db:"last_modified"`
	Description  string
	EventID      int `db:"event_id"`
	ID           int
	RefID        EventItemRefID `db:"ref_id"`
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
		VALUES (@refID, @eventID, @description)
		RETURNING *`
	args := pgx.NamedArgs{
		"refID":       refID,
		"eventID":     eventID,
		"description": description,
	}
	return QueryOneTx[EventItem](ctx, db, q, args)
}

func UpdateEventItem(ctx context.Context, db PgxHandle,
	eventItemID int, description string,
) error {
	q := `
		UPDATE event_item_
		SET description = @description
		WHERE id = @eventItemID`
	args := pgx.NamedArgs{
		"description": description,
		"eventItemID": eventItemID,
	}
	return ExecTx[EventItem](ctx, db, q, args)
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

func GetEventItemsByIDs(ctx context.Context, db PgxHandle,
	eventItemIDs []int,
) ([]*EventItem, error) {
	q := `SELECT * FROM event_item_ WHERE id = ANY($1)`
	return Query[EventItem](ctx, db, q, eventItemIDs)
}

func GetEventItemByRefID(ctx context.Context, db PgxHandle,
	refID EventItemRefID,
) (*EventItem, error) {
	q := `SELECT * FROM event_item_ WHERE ref_id = $1`
	return QueryOne[EventItem](ctx, db, q, refID)
}

func GetEventItemsByEvent(ctx context.Context, db PgxHandle,
	eventID int,
) ([]*EventItem, error) {
	q := `
		SELECT * FROM event_item_
		WHERE event_id = $1
		ORDER BY
			created DESC,
			id DESC`
	return Query[EventItem](ctx, db, q, eventID)
}

type EventItemCount struct {
	EventID int `db:"event_id"`
	Count   int
}

func GetEventItemsCountByEventIDs(ctx context.Context, db PgxHandle,
	eventIDs []int,
) ([]*EventItemCount, error) {
	q := `
	SELECT
		e.id as event_id,
		count(ei.id)
	FROM event_ e
	LEFT JOIN event_item_ ei ON
		e.id = ei.event_id
	WHERE
		e.id = ANY ($1)
	GROUP BY e.id
	ORDER BY e.id`
	return Query[EventItemCount](ctx, db, q, eventIDs)
}
