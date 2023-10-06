package model

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
)

type Favorite struct {
	Id      int
	UserId  int `db:"user_id"`
	EventId int `db:"event_id"`
	Created time.Time
	User    *User  `db:"-"`
	Event   *Event `db:"-"`
}

func (fav *Favorite) Insert(ctx context.Context, db PgxHandle) error {
	q := `INSERT INTO favorite_ (user_id, event_id) VALUES ($1, $2) RETURNING *`
	res, err := QueryOneTx[Favorite](ctx, db, q, fav.UserId, fav.EventId)
	if err != nil {
		return err
	}
	fav.Id = res.Id
	fav.Created = res.Created
	return nil
}

func (fav *Favorite) Delete(ctx context.Context, db PgxHandle) error {
	q := `DELETE FROM favorite_ WHERE id = $1`
	return ExecTx[Favorite](ctx, db, q, fav.Id)
}

func (fav *Favorite) GetEvent(ctx context.Context, db PgxHandle) (*Event, error) {
	event, err := GetEventById(ctx, db, fav.EventId)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func NewFavorite(
	ctx context.Context,
	db PgxHandle,
	userId int,
	eventId int,
) (*Favorite, error) {
	favorite := &Favorite{
		UserId:  userId,
		EventId: eventId,
	}
	err := favorite.Insert(ctx, db)
	if err != nil {
		return nil, err
	}
	return favorite, nil
}

func GetFavoriteById(ctx context.Context, db PgxHandle, id int) (*Favorite, error) {
	q := `SELECT * FROM favorite_ WHERE id = $1`
	return QueryOne[Favorite](ctx, db, q, id)
}

func GetFavoritesByUser(ctx context.Context, db PgxHandle, user *User) ([]*Favorite, error) {
	q := `SELECT * FROM favorite_ WHERE user_id = $1 ORDER BY created DESC,id DESC`
	return Query[Favorite](ctx, db, q, user.Id)
}

func GetFavoritesByEvent(ctx context.Context, db PgxHandle, event *Event) ([]*Favorite, error) {
	q := `SELECT * FROM favorite_ WHERE event_id = $1 ORDER BY created DESC,id DESC`
	return Query[Favorite](ctx, db, q, event.Id)
}

func GetFavoriteByUserEvent(ctx context.Context, db PgxHandle, user *User, event *Event) (*Favorite, error) {
	q := `SELECT * FROM favorite_ WHERE user_id = $1 AND event_id = $2 ORDER BY created DESC,id DESC`
	return QueryOne[Favorite](ctx, db, q, user.Id, event.Id)
}

func GetFavoriteCountByUser(ctx context.Context, db PgxHandle, user *User) (int, error) {
	q := `SELECT count(*) FROM favorite_ WHERE user_id = $1`
	var count int = 0
	err := pgxscan.Get(ctx, db, &count, q, user.Id)
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
	return Query[Favorite](ctx, db, q, user.Id, limit, offset)
}
