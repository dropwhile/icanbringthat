package model

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type Favorite struct {
	Created time.Time
	EventID int `db:"event_id"`
	UserID  int `db:"user_id"`
	ID      int
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

func GetFavoriteEventsByUserFiltered(
	ctx context.Context, db PgxHandle,
	userID int, archived bool,
) ([]*Event, error) {
	q := `
	SELECT event_.*
	FROM event_ 
	JOIN favorite_ ON
		favorite_.event_id = event_.id
	WHERE
	 	favorite_.user_id = @userID AND
		event_.archived = @archived
	ORDER BY 
		event_.start_time DESC,
		event_.id DESC
	`
	args := pgx.NamedArgs{
		"userID":   userID,
		"archived": archived,
	}
	return Query[Event](ctx, db, q, args)
}

func GetFavoriteEventsByUserPaginatedFiltered(
	ctx context.Context, db PgxHandle,
	userID int, limit, offset int, archived bool,
) ([]*Event, error) {
	q := `
	SELECT event_.*
	FROM event_ 
	JOIN favorite_ ON
		favorite_.event_id = event_.id
	WHERE
		favorite_.user_id = @userID AND
		event_.archived = @archived
	ORDER BY 
		event_.start_time DESC,
		event_.id DESC
	LIMIT @limit OFFSET @offset
	`
	args := pgx.NamedArgs{
		"userID":   userID,
		"limit":    limit,
		"offset":   offset,
		"archived": archived,
	}
	return Query[Event](ctx, db, q, args)
}

func GetFavoriteCountByUser(ctx context.Context, db PgxHandle,
	userID int,
) (*BifurcatedRowCounts, error) {
	if userID == 0 {
		return nil, errors.New("nil user supplied")
	}
	q := `
		SELECT
			count(*) filter (WHERE event_.archived IS NOT TRUE) as current,
			count(*) filter (WHERE event_.archived IS TRUE) as archived
		FROM favorite_ fav
		JOIN event_ ON 
			event_.id = fav.event_id
		WHERE fav.user_id = $1`
	return QueryOne[BifurcatedRowCounts](ctx, db, q, userID)
}
