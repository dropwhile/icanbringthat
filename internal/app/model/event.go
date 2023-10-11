package model

import (
	"context"
	"time"

	"github.com/dropwhile/refid"
	"github.com/georgysavva/scany/v2/pgxscan"
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
	LastModified  time.Time    `db:"last_modified"`
	Items         []*EventItem `db:"-"`
}

func NewEvent(ctx context.Context, db PgxHandle,
	userID int, name, description string,
	startTime time.Time,
	startTimeTz *TimeZone,
) (*Event, error) {
	refID := refid.Must(NewEventRefID())
	itemSortOrder := []int{}
	return CreateEvent(
		ctx, db,
		refID, userID,
		name, description,
		itemSortOrder,
		startTime, startTimeTz,
	)
}

func CreateEvent(ctx context.Context, db PgxHandle,
	refID EventRefID,
	userID int, name, description string,
	itemSortOrder []int,
	startTime time.Time,
	startTimeTz *TimeZone,
) (*Event, error) {
	q := `
		INSERT INTO event_ (
			ref_id, user_id, name, description,
			item_sort_order, start_time, start_time_tz
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING *`
	return QueryOneTx[Event](
		ctx, db, q, refID,
		userID, name, description,
		itemSortOrder, startTime, startTimeTz,
	)
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
			name = $1,
			description = $2,
			item_sort_order = $3,
			start_time = $4,
			start_time_tz = $5
		WHERE id = $6`
	return ExecTx[Event](
		ctx, db, q,
		name, description,
		itemSortOrder,
		startTime, startTimeTz,
		eventID)
}

func DeleteEvent(ctx context.Context, db PgxHandle, eventID int) error {
	q := `DELETE FROM event_ WHERE id = $1`
	return ExecTx[Event](ctx, db, q, eventID)
}

func GetEventByID(ctx context.Context, db PgxHandle, eventID int) (*Event, error) {
	q := `SELECT * FROM event_ WHERE id = $1`
	return QueryOne[Event](ctx, db, q, eventID)
}

func GetEventsByIDs(ctx context.Context, db PgxHandle, eventIDs []int) ([]*Event, error) {
	q := `SELECT * FROM event_ WHERE id = ANY($1)`
	return Query[Event](ctx, db, q, eventIDs)
}

func GetEventByRefID(ctx context.Context, db PgxHandle, refID EventRefID) (*Event, error) {
	q := `SELECT * FROM event_ WHERE ref_id = $1`
	return QueryOne[Event](ctx, db, q, refID)
}

func GetEventsByUser(ctx context.Context, db PgxHandle, user *User) ([]*Event, error) {
	q := `
		SELECT * FROM event_
		WHERE
			event_.user_id = $1
		ORDER BY
			start_time DESC,
			id DESC`
	return Query[Event](ctx, db, q, user.ID)
}

func GetEventsByUserPaginated(
	ctx context.Context, db PgxHandle,
	user *User, limit, offset int,
) ([]*Event, error) {
	q := `
		SELECT * FROM event_
		WHERE
			event_.user_id = $1
		ORDER BY
			start_time DESC,
			id DESC
		LIMIT $2 OFFSET $3`
	return Query[Event](ctx, db, q, user.ID, limit, offset)
}

func GetEventsComingSoonByUserPaginated(
	ctx context.Context, db PgxHandle,
	user *User, limit, offset int,
) ([]*Event, error) {
	q := `
		SELECT *
		FROM event_
		WHERE
			event_.user_id = $1 AND
			start_time > CURRENT_TIMESTAMP(0)
		ORDER BY 
			start_time ASC,
			id ASC
		LIMIT $2 OFFSET $3`
	return Query[Event](ctx, db, q, user.ID, limit, offset)
}

func GetEventCountByUser(ctx context.Context, db PgxHandle, user *User) (int, error) {
	q := `SELECT count(*) FROM event_ WHERE user_id = $1`
	var count int = 0
	err := pgxscan.Get(ctx, db, &count, q, user.ID)
	return count, err
}
