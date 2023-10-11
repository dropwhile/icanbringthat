package model

import (
	"context"
	"time"
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
	q := `INSERT INTO favorite_ (user_id, event_id) VALUES ($1, $2) RETURNING *`
	return QueryOneTx[Favorite](ctx, db, q, userID, eventID)
}

func DeleteFavorite(ctx context.Context, db PgxHandle, favID int) error {
	q := `DELETE FROM favorite_ WHERE id = $1`
	return ExecTx[Favorite](ctx, db, q, favID)
}

func GetFavoriteByID(ctx context.Context, db PgxHandle, favID int) (*Favorite, error) {
	q := `SELECT * FROM favorite_ WHERE id = $1`
	return QueryOne[Favorite](ctx, db, q, favID)
}

func GetFavoritesByUser(ctx context.Context, db PgxHandle,
	user *User,
) ([]*Favorite, error) {
	q := `
		SELECT * FROM favorite_
		WHERE user_id = $1
		ORDER BY
			created DESC,
			id DESC`
	return Query[Favorite](ctx, db, q, user.ID)
}

func GetFavoritesByEvent(ctx context.Context, db PgxHandle,
	event *Event,
) ([]*Favorite, error) {
	q := `
		SELECT * FROM favorite_
		WHERE event_id = $1
		ORDER BY
			created DESC,
			id DESC`
	return Query[Favorite](ctx, db, q, event.ID)
}

func GetFavoriteByUserEvent(ctx context.Context, db PgxHandle,
	user *User, event *Event,
) (*Favorite, error) {
	q := `
		SELECT * FROM favorite_
		WHERE
			user_id = $1 AND
			event_id = $2
		ORDER BY
			created DESC,
			id DESC`
	return QueryOne[Favorite](ctx, db, q, user.ID, event.ID)
}

func GetFavoriteCountByUser(ctx context.Context, db PgxHandle,
	user *User,
) (int, error) {
	q := `SELECT count(*) FROM favorite_ WHERE user_id = $1`
	return Get[int](ctx, db, q, user.ID)
}

func GetFavoritesByUserPaginated(ctx context.Context, db PgxHandle,
	user *User, limit, offset int,
) ([]*Favorite, error) {
	q := `
	SELECT 
		favorite_.*
	FROM favorite_ 
	JOIN event_ ON
		favorite_.event_id = event_.id
	WHERE favorite_.user_id = $1 
	ORDER BY 
		event_.start_time DESC,
		event_.id DESC
	LIMIT $2 OFFSET $3
	`
	return Query[Favorite](ctx, db, q, user.ID, limit, offset)
}
