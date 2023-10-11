package model

import (
	"context"
	"time"

	"github.com/dropwhile/refid"
	"github.com/georgysavva/scany/v2/pgxscan"
)

//go:generate go run ../../../cmd/refidgen -t Earmark -v 4

type Earmark struct {
	ID           int
	RefID        EarmarkRefID `db:"ref_id"`
	EventItemID  int          `db:"event_item_id"`
	UserID       int          `db:"user_id"`
	Note         string
	Created      time.Time
	LastModified time.Time  `db:"last_modified"`
	EventItem    *EventItem `db:"-"`
	User         *User      `db:"-"`
}

func NewEarmark(ctx context.Context, db PgxHandle,
	eventItemID, userID int, note string,
) (*Earmark, error) {
	refID := refid.Must(NewEarmarkRefID())
	return CreateEarmark(ctx, db, refID, eventItemID, userID, note)
}

func CreateEarmark(ctx context.Context, db PgxHandle,
	refID EarmarkRefID, eventItemID int, userID int, note string,
) (*Earmark, error) {
	q := `
		INSERT INTO earmark_ (
			ref_id, event_item_id, user_id, note
		)
		VALUES ($1, $2, $3, $4)
		RETURNING *`
	return QueryOneTx[Earmark](ctx, db, q, refID, eventItemID, userID, note)
}

func UpdateEarmark(ctx context.Context, db PgxHandle,
	earmarkID int, note string,
) error {
	q := `UPDATE earmark_ SET note = $1 WHERE id = $2`
	return ExecTx[Earmark](ctx, db, q, note, earmarkID)
}

func DeleteEarmark(ctx context.Context, db PgxHandle,
	earmarkID int,
) error {
	q := `DELETE FROM earmark_ WHERE id = $1`
	return ExecTx[Earmark](ctx, db, q, earmarkID)
}

func GetEarmarkByID(ctx context.Context, db PgxHandle,
	earmarkID int,
) (*Earmark, error) {
	q := `SELECT * FROM earmark_ WHERE id = $1`
	return QueryOne[Earmark](ctx, db, q, earmarkID)
}

func GetEarmarkByRefID(ctx context.Context, db PgxHandle,
	refID EarmarkRefID,
) (*Earmark, error) {
	q := `SELECT * FROM earmark_ WHERE ref_id = $1`
	return QueryOne[Earmark](ctx, db, q, refID)
}

func GetEarmarkByEventItem(ctx context.Context, db PgxHandle,
	eventItem *EventItem,
) (*Earmark, error) {
	q := `SELECT * FROM earmark_ WHERE event_item_id = $1`
	return QueryOne[Earmark](ctx, db, q, eventItem.ID)
}

func GetEarmarksByUser(ctx context.Context, db PgxHandle,
	user *User,
) ([]*Earmark, error) {
	q := `
		SELECT * FROM earmark_
		WHERE user_id = $1
		ORDER BY
			created DESC,
			id DESC`
	return Query[Earmark](ctx, db, q, user.ID)
}

func GetEarmarksByEvent(ctx context.Context, db PgxHandle,
	event *Event,
) ([]*Earmark, error) {
	q := `
		SELECT earmark_.*
		FROM earmark_
		JOIN event_item_ ON 
			event_item_.id = earmark_.event_item_id
		WHERE 
			event_item_.event_id = $1
	`
	return Query[Earmark](ctx, db, q, event.ID)
}

func GetEarmarksWithEventsByUser(ctx context.Context, db PgxHandle,
	user *User,
) ([]*Earmark, error) {
	q := `
		SELECT *
		FROM earmark_
		JOIN event_ ON
			event_.id = earmark_.id 
		WHERE 
			user_id = $1
		ORDER BY
			created DESC,
			id DESC
	`
	return Query[Earmark](ctx, db, q, user.ID)
}

func GetEarmarksByUserPaginated(ctx context.Context, db PgxHandle,
	user *User, limit, offset int,
) ([]*Earmark, error) {
	q := `
		SELECT * FROM earmark_
		WHERE
			earmark_.user_id = $1
		ORDER BY
			created DESC,
			id DESC
		LIMIT $2 OFFSET $3`
	return Query[Earmark](ctx, db, q, user.ID, limit, offset)
}

func GetEarmarkCountByUser(ctx context.Context, db PgxHandle, user *User) (int, error) {
	q := `SELECT count(*) FROM earmark_ WHERE user_id = $1`
	var count int = 0
	err := pgxscan.Get(ctx, db, &count, q, user.ID)
	return count, err
}
