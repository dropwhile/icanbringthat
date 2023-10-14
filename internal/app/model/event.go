package model

import (
	"context"
	"time"

	"github.com/dropwhile/refid"
	"github.com/jackc/pgx/v5"
)

//go:generate go run ../../../cmd/refidgen -t Event -v 2

type Event struct {
	ID            int
	RefID         EventRefID `db:"ref_id"`
	UserID        int        `db:"user_id"`
	Name          string
	Description   string
	StartTime     time.Time `db:"start_time"`
	StartTimeTz   *TimeZone `db:"start_time_tz"`
	ItemSortOrder []int     `db:"item_sort_order"`
	Created       time.Time
	LastModified  time.Time `db:"last_modified"`
}

func NewEvent(ctx context.Context, db PgxHandle,
	userID int, name, description string,
	startTime time.Time,
	startTimeTz *TimeZone,
) (*Event, error) {
	refID := refid.Must(NewEventRefID())
	return CreateEvent(
		ctx, db,
		refID, userID,
		name, description,
		startTime, startTimeTz,
	)
}

func CreateEvent(ctx context.Context, db PgxHandle,
	refID EventRefID,
	userID int, name, description string,
	startTime time.Time,
	startTimeTz *TimeZone,
) (*Event, error) {
	q := `
		INSERT INTO event_ (
			ref_id, user_id, name, description,
			start_time, start_time_tz
		)
		VALUES (
			@refID, @userID, @name, @description,
			@startTime, @startTimeTz
		)
		RETURNING *`
	args := pgx.NamedArgs{
		"refID":       refID,
		"userID":      userID,
		"name":        name,
		"description": description,
		"startTime":   startTime,
		"startTimeTz": startTimeTz,
	}
	return QueryOneTx[Event](ctx, db, q, args)
}

func UpdateEvent(ctx context.Context, db PgxHandle,
	eventID int, name, description string,
	itemSortOrder []int,
	startTime time.Time,
	startTimeTz *TimeZone,
) error {
	q := `
		UPDATE event_
		SET
			name = @name,
			description = @description,
			item_sort_order = @itemSortOrder,
			start_time = @startTime,
			start_time_tz = @startTimeTz
		WHERE id = @eventID`
	args := pgx.NamedArgs{
		"name":          name,
		"description":   description,
		"itemSortOrder": itemSortOrder,
		"startTime":     startTime,
		"startTimeTz":   startTimeTz,
		"eventID":       eventID,
	}
	return ExecTx[Event](ctx, db, q, args)
}

func DeleteEvent(ctx context.Context, db PgxHandle,
	eventID int,
) error {
	q := `DELETE FROM event_ WHERE id = $1`
	return ExecTx[Event](ctx, db, q, eventID)
}

func GetEventByID(ctx context.Context, db PgxHandle,
	eventID int,
) (*Event, error) {
	q := `SELECT * FROM event_ WHERE id = $1`
	return QueryOne[Event](ctx, db, q, eventID)
}

func GetEventsByIDs(ctx context.Context, db PgxHandle,
	eventIDs []int,
) ([]*Event, error) {
	q := `SELECT * FROM event_ WHERE id = ANY($1)`
	return Query[Event](ctx, db, q, eventIDs)
}

func GetEventByRefID(ctx context.Context, db PgxHandle,
	refID EventRefID,
) (*Event, error) {
	q := `SELECT * FROM event_ WHERE ref_id = $1`
	return QueryOne[Event](ctx, db, q, refID)
}

func GetEventsByUserPaginated(
	ctx context.Context, db PgxHandle,
	userID int, limit, offset int,
) ([]*Event, error) {
	q := `
		SELECT * FROM event_
		WHERE
			event_.user_id = @userID
		ORDER BY
			start_time DESC,
			id DESC
		LIMIT @limit OFFSET @offset`
	args := pgx.NamedArgs{
		"userID": userID,
		"limit":  limit,
		"offset": offset,
	}
	return Query[Event](ctx, db, q, args)
}

func GetEventsComingSoonByUserPaginated(
	ctx context.Context, db PgxHandle,
	userID int, limit, offset int,
) ([]*Event, error) {
	q := `
		SELECT *
		FROM event_
		WHERE
			event_.user_id = @userID AND
			start_time > CURRENT_TIMESTAMP(0)
		ORDER BY 
			start_time ASC,
			id ASC
		LIMIT @limit OFFSET @offset`
	args := pgx.NamedArgs{
		"userID": userID,
		"limit":  limit,
		"offset": offset,
	}
	return Query[Event](ctx, db, q, args)
}

func GetEventCountByUser(ctx context.Context, db PgxHandle,
	userID int,
) (int, error) {
	q := `SELECT count(*) FROM event_ WHERE user_id = $1`
	return Get[int](ctx, db, q, userID)
}
