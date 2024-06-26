// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package model

import (
	"context"
	"time"

	"github.com/dropwhile/refid/v2/reftag"
	"github.com/jackc/pgx/v5"

	"github.com/dropwhile/icanbringthat/internal/util"
)

type NotificationRefID struct {
	reftag.IDt8
}

var NewNotificationRefID = reftag.New[NotificationRefID]

type Notification struct {
	Created      time.Time
	LastModified time.Time `db:"last_modified"`
	Message      string
	UserID       int `db:"user_id"`
	ID           int
	RefID        NotificationRefID `db:"ref_id"`
	Read         bool
}

func NewNotification(ctx context.Context, db PgxHandle,
	userID int, message string,
) (*Notification, error) {
	refID := util.Must(NewNotificationRefID())
	return CreateNotification(ctx, db, refID, userID, message)
}

func CreateNotification(ctx context.Context, db PgxHandle,
	refID NotificationRefID, userID int, message string,
) (*Notification, error) {
	q := `
		INSERT INTO notification_ (
			ref_id, user_id, message
		)
		VALUES (
			@refID, @userID, @message
		)
		RETURNING *`
	args := pgx.NamedArgs{
		"refID":   refID,
		"userID":  userID,
		"message": message,
	}
	return QueryOneTx[Notification](ctx, db, q, args)
}

func UpdateNotification(ctx context.Context, db PgxHandle,
	ID int, read bool,
) (*Notification, error) {
	q := `
		UPDATE notification_
		SET read = $1
		WHERE id = $2`
	return QueryOne[Notification](ctx, db, q, read, ID)
}

func DeleteNotification(ctx context.Context, db PgxHandle, ID int) error {
	q := `DELETE FROM notification_ WHERE id = $1`
	return ExecTx[Notification](ctx, db, q, ID)
}

func DeleteNotificationsByUser(ctx context.Context, db PgxHandle, userID int) error {
	q := `DELETE FROM notification_ WHERE user_id = $1`
	return ExecTx[Notification](ctx, db, q, userID)
}

func GetNotificationByID(ctx context.Context, db PgxHandle, ID int) (*Notification, error) {
	q := `SELECT * FROM notification_ WHERE id = $1`
	return QueryOne[Notification](ctx, db, q, ID)
}

func GetNotificationByRefID(ctx context.Context, db PgxHandle,
	refID NotificationRefID,
) (*Notification, error) {
	q := `SELECT * FROM notification_ WHERE ref_id = $1`
	return QueryOne[Notification](ctx, db, q, refID)
}

func GetNotificationCountByUser(ctx context.Context, db PgxHandle,
	userID int,
) (int, error) {
	q := `SELECT count(*) FROM notification_ WHERE user_id = $1`
	return Get[int](ctx, db, q, userID)
}

func GetNotificationsByUserPaginated(ctx context.Context, db PgxHandle,
	userID int, limit, offset int,
) ([]*Notification, error) {
	q := `
	SELECT *
	FROM notification_ 
	WHERE
		user_id = @userID AND
		read = FALSE
	ORDER BY 
		created DESC
	LIMIT @limit
	OFFSET @offset
	`
	args := pgx.NamedArgs{
		"userID": userID,
		"limit":  limit,
		"offset": offset,
	}
	return Query[Notification](ctx, db, q, args)
}

func GetNotificationsByUser(ctx context.Context, db PgxHandle,
	userID int,
) ([]*Notification, error) {
	q := `
	SELECT *
	FROM notification_ 
	WHERE
		user_id = @userID AND
		read = FALSE
	ORDER BY 
		created DESC
	`
	args := pgx.NamedArgs{
		"userID": userID,
	}
	return Query[Notification](ctx, db, q, args)
}
