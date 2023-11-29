package model

import (
	"context"
	"time"

	"github.com/dropwhile/refid"
	"github.com/dropwhile/refid/reftag"
	"github.com/jackc/pgx/v5"
)

type (
	EventRefID     = reftag.IDt2
	EventRefIDNull = reftag.NullIDt2
)

var (
	NewEventRefID       = reftag.New[EventRefID]
	EventRefIDMatcher   = reftag.NewMatcher[EventRefID]()
	EventRefIDFromBytes = reftag.FromBytes[EventRefID]
	ParseEventRefID     = reftag.Parse[EventRefID]
)

type Event struct {
	Created       time.Time
	LastModified  time.Time `db:"last_modified"`
	StartTime     time.Time `db:"start_time"`
	StartTimeTz   *TimeZone `db:"start_time_tz"`
	Name          string
	Description   string
	ItemSortOrder []int `db:"item_sort_order"`
	Archived      bool
	UserID        int `db:"user_id"`
	ID            int
	RefID         EventRefID `db:"ref_id"`
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

func GetEventByEarmark(ctx context.Context, db PgxHandle,
	earmark *Earmark,
) (*Event, error) {
	q := `
		SELECT ev.*
		FROM event_ ev 
		JOIN event_item_ ei ON
			ev.id = ei.event_id
		JOIN earmark_ em ON
			em.event_item_id = ei.id
		WHERE em.id = $1`
	return QueryOne[Event](ctx, db, q, earmark.ID)
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

func GetEventsByUserPaginatedFiltered(
	ctx context.Context, db PgxHandle,
	userID int, limit, offset int, archived bool,
) ([]*Event, error) {
	q := `
		SELECT * FROM event_
		WHERE
			event_.user_id = @userID AND
			archived = @archived
		ORDER BY
			start_time DESC,
			id DESC
		LIMIT @limit OFFSET @offset`
	args := pgx.NamedArgs{
		"userID":   userID,
		"limit":    limit,
		"offset":   offset,
		"archived": archived,
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

func GetEventCountsByUser(ctx context.Context, db PgxHandle,
	userID int,
) (*BifurcatedRowCounts, error) {
	q := `
		SELECT
			count(*) filter (WHERE archived IS NOT TRUE) as current,
			count(*) filter (WHERE archived IS TRUE) as archived
		FROM event_
		WHERE user_id = $1`
	return QueryOne[BifurcatedRowCounts](ctx, db, q, userID)
}

func ArchiveOldEvents(ctx context.Context, db PgxHandle) error {
	q := `
		UPDATE event_
		SET archived = TRUE
		WHERE
			date_trunc('hour', start_time) AT TIME ZONE 'UTC' <
				timezone('utc', CURRENT_TIMESTAMP) - INTERVAL '1 day'
			AND archived IS FALSE
	`
	return ExecTx[Event](ctx, db, q)
}
