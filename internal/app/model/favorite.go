package model

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
)

type Favorite struct {
	ID      int
	UserID  int `db:"user_id"`
	EventID int `db:"event_id"`
	Created time.Time
	User    *User  `db:"-"`
	Event   *Event `db:"-"`
}

func (fav *Favorite) Insert(ctx context.Context, db PgxHandle) error {
	q := `INSERT INTO favorite_ (user_id, event_id) VALUES ($1, $2) RETURNING *`
	res, err := QueryOneTx[Favorite](ctx, db, q, fav.UserID, fav.EventID)
	if err != nil {
		return err
	}
	fav.ID = res.ID
	fav.Created = res.Created
	return nil
}

func (fav *Favorite) Delete(ctx context.Context, db PgxHandle) error {
	q := `DELETE FROM favorite_ WHERE id = $1`
	return ExecTx[Favorite](ctx, db, q, fav.ID)
}

func (fav *Favorite) GetEvent(ctx context.Context, db PgxHandle) (*Event, error) {
	event, err := GetEventByID(ctx, db, fav.EventID)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func NewFavorite(ctx context.Context, db PgxHandle, userID int, eventID int) (*Favorite, error) {
	favorite := &Favorite{
		UserID:  userID,
		EventID: eventID,
	}
	err := favorite.Insert(ctx, db)
	if err != nil {
		return nil, err
	}
	return favorite, nil
}

func GetFavoriteByID(ctx context.Context, db PgxHandle, favID int) (*Favorite, error) {
	q := `SELECT * FROM favorite_ WHERE id = $1`
	return QueryOne[Favorite](ctx, db, q, favID)
}

func GetFavoritesByUser(ctx context.Context, db PgxHandle, user *User) ([]*Favorite, error) {
	q := `SELECT * FROM favorite_ WHERE user_id = $1 ORDER BY created DESC,id DESC`
	return Query[Favorite](ctx, db, q, user.ID)
}

func GetFavoritesByEvent(ctx context.Context, db PgxHandle, event *Event) ([]*Favorite, error) {
	q := `SELECT * FROM favorite_ WHERE event_id = $1 ORDER BY created DESC,id DESC`
	return Query[Favorite](ctx, db, q, event.ID)
}

func GetFavoriteByUserEvent(ctx context.Context, db PgxHandle, user *User, event *Event) (*Favorite, error) {
	q := `SELECT * FROM favorite_ WHERE user_id = $1 AND event_id = $2 ORDER BY created DESC,id DESC`
	return QueryOne[Favorite](ctx, db, q, user.ID, event.ID)
}

func GetFavoriteCountByUser(ctx context.Context, db PgxHandle, user *User) (int, error) {
	q := `SELECT count(*) FROM favorite_ WHERE user_id = $1`
	var count int = 0
	err := pgxscan.Get(ctx, db, &count, q, user.ID)
	return count, err
}

func GetFavoritesByUserPaginated(ctx context.Context, db PgxHandle, user *User, limit, offset int) ([]*Favorite, error) {
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
