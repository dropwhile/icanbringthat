package model

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Favorite struct {
	ID      int
	UserID  int `db:"user_id"`
	EventID int `db:"event_id"`
	Created time.Time
	User    *User  `db:"-"`
	Event   *Event `db:"-"`
}

func CreateFavorite(ctx context.Context, db PgxHandle,
	userID int, eventID int,
) (*Favorite, error) {
	q := `
		INSERT INTO favorite_ (
			user_id, event_id
		)
		VALUES (
			@userID, @eventID
		)
		RETURNING *`
	args := pgx.NamedArgs{"userID": userID, "eventID": eventID}
	return QueryOneTx[Favorite](ctx, db, q, args)
}

func DeleteFavorite(ctx context.Context, db PgxHandle, favID int) error {
	q := `DELETE FROM favorite_ WHERE id = $1`
	return ExecTx[Favorite](ctx, db, q, favID)
}

func GetFavoriteByID(ctx context.Context, db PgxHandle, favID int) (*Favorite, error) {
	q := `SELECT * FROM favorite_ WHERE id = $1`
	return QueryOne[Favorite](ctx, db, q, favID)
}

func GetFavoriteByUserEvent(ctx context.Context, db PgxHandle,
	userID int, eventID int,
) (*Favorite, error) {
	q := `
		SELECT * FROM favorite_
		WHERE
			user_id = @userID AND
			event_id = @eventID
		ORDER BY
			created DESC,
			id DESC`
	args := pgx.NamedArgs{"userID": userID, "eventID": eventID}
	return QueryOne[Favorite](ctx, db, q, args)
}

func GetFavoriteCountByUser(ctx context.Context, db PgxHandle,
	user *User,
) (int, error) {
	q := `SELECT count(*) FROM favorite_ WHERE user_id = $1`
	return Get[int](ctx, db, q, user.ID)
}

func GetFavoritesByUserPaginated(ctx context.Context, db PgxHandle,
	userID int, limit, offset int,
) ([]*Favorite, error) {
	q := `
	SELECT favorite_.*
	FROM favorite_ 
	JOIN event_ ON
		favorite_.event_id = event_.id
	WHERE favorite_.user_id = @userID 
	ORDER BY 
		event_.start_time DESC,
		event_.id DESC
	LIMIT @limit OFFSET @offset
	`
	args := pgx.NamedArgs{
		"userID": userID,
		"limit":  limit,
		"offset": offset,
	}
	return Query[Favorite](ctx, db, q, userID, limit, args)
}
