package model

import (
	"context"
	"time"

	"github.com/dropwhile/refid"
	"github.com/dropwhile/refid/reftag"
	"github.com/jackc/pgx/v5"
)

type (
	EarmarkRefID     = reftag.IDt4
	EarmarkRefIDNull = reftag.NullIDt4
)

var (
	NewEarmarkRefID       = reftag.New[EarmarkRefID]
	EarmarkRefIDMatcher   = reftag.NewMatcher[EarmarkRefID]()
	EarmarkRefIDFromBytes = reftag.FromBytes[EarmarkRefID]
	ParseEarmarkRefID     = reftag.Parse[EarmarkRefID]
)

type Earmark struct {
	ID           int
	RefID        EarmarkRefID `db:"ref_id"`
	EventItemID  int          `db:"event_item_id"`
	UserID       int          `db:"user_id"`
	Note         string
	Created      time.Time
	LastModified time.Time `db:"last_modified"`
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
		VALUES (@refID, @eventItemID, @userID, @note)
		RETURNING *`
	args := pgx.NamedArgs{
		"refID":       refID,
		"eventItemID": eventItemID,
		"userID":      userID,
		"note":        note,
	}
	return QueryOneTx[Earmark](ctx, db, q, args)
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
	eventItemID int,
) (*Earmark, error) {
	q := `SELECT * FROM earmark_ WHERE event_item_id = $1`
	return QueryOne[Earmark](ctx, db, q, eventItemID)
}

func GetEarmarksByEventItemIDs(ctx context.Context, db PgxHandle,
	eventItemIDs []int,
) ([]*Earmark, error) {
	q := `SELECT * FROM earmark_ WHERE event_item_id = ANY($1)`
	return Query[Earmark](ctx, db, q, eventItemIDs)
}

func GetEarmarksByEvent(ctx context.Context, db PgxHandle,
	eventID int,
) ([]*Earmark, error) {
	q := `
		SELECT earmark_.*
		FROM earmark_
		JOIN event_item_ ON 
			event_item_.id = earmark_.event_item_id
		WHERE 
			event_item_.event_id = $1
	`
	return Query[Earmark](ctx, db, q, eventID)
}

func GetEarmarksByUserPaginated(ctx context.Context, db PgxHandle,
	userID int, limit, offset int,
) ([]*Earmark, error) {
	q := `
		SELECT * FROM earmark_
		WHERE
			earmark_.user_id = @userID
		ORDER BY
			created DESC,
			id DESC
		LIMIT @limit OFFSET @offset`
	args := pgx.NamedArgs{
		"userID": userID,
		"limit":  limit,
		"offset": offset,
	}
	return Query[Earmark](ctx, db, q, args)
}

func GetEarmarkCountByUser(ctx context.Context, db PgxHandle,
	user *User,
) (int, error) {
	q := `SELECT count(*) FROM earmark_ WHERE user_id = $1`
	return Get[int](ctx, db, q, user.ID)
}
