package model

import (
	"context"
	"time"

	"github.com/dropwhile/refid/v2"
	"github.com/dropwhile/refid/v2/reftag"
	"github.com/jackc/pgx/v5"
	"github.com/samber/mo"
)

type EventRefID = reftag.IDt2

type EventRefIDNull struct {
	reftag.NullIDt2
}

var NewEventRefID = reftag.New[EventRefID]

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

func (ev *Event) When() time.Time {
	return ev.StartTime.In(ev.StartTimeTz.Location)
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

type EventUpdateModelValues struct {
	StartTime     mo.Option[time.Time]
	Tz            mo.Option[*TimeZone]
	Name          mo.Option[string]
	Description   mo.Option[string]
	ItemSortOrder mo.Option[[]int]
}

func UpdateEvent(ctx context.Context, db PgxHandle, eventID int,
	vals *EventUpdateModelValues,
) error {
	q := `
		UPDATE event_
		SET
			name = COALESCE(@name, name),
			description = COALESCE(@description, description),
			item_sort_order = COALESCE(@itemSortOrder, item_sort_order),
			start_time = COALESCE(@startTime, start_time),
			start_time_tz = COALESCE(@startTimeTz, start_time_tz)
		WHERE id = @eventID`
	args := pgx.NamedArgs{
		"name":          vals.Name,
		"description":   vals.Description,
		"itemSortOrder": vals.ItemSortOrder,
		"startTime":     vals.StartTime,
		"startTimeTz":   vals.Tz,
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

func GetEventByEventItemID(ctx context.Context, db PgxHandle,
	eventItemID int,
) (*Event, error) {
	q := `
		SELECT ev.*
		FROM event_ ev 
		JOIN event_item_ ei ON
			ev.id = ei.event_id
		WHERE ei.id = $1`
	return QueryOne[Event](ctx, db, q, eventItemID)
}

func GetEventsByUserFiltered(
	ctx context.Context, db PgxHandle,
	userID int, archived bool,
) ([]*Event, error) {
	q := `
		SELECT * FROM event_
		WHERE
			event_.user_id = @userID AND
			archived = @archived
		ORDER BY
			start_time DESC,
			id DESC
		`
	args := pgx.NamedArgs{
		"userID":   userID,
		"archived": archived,
	}
	return Query[Event](ctx, db, q, args)
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
