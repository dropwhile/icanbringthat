package model

import (
	"context"
	"time"

	"github.com/dropwhile/refid"
	"github.com/jackc/pgx/v5"
)

//go:generate go run ../../../cmd/refidgen -t Notification -v 8

type Notification struct {
	ID           int
	RefID        NotificationRefID `db:"ref_id"`
	UserID       int               `db:"user_id"`
	Message      string
	Read         bool
	Created      time.Time
	LastModified time.Time `db:"last_modified"`
}

func NewNotification(ctx context.Context, db PgxHandle,
	userID int, message string,
) (*Notification, error) {
	refID := refid.Must(NewNotificationRefID())
	return CreateNotification(
		ctx, db,
		refID, userID, message,
	)
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
	LIMIT @limit OFFSET @offset
	`
	args := pgx.NamedArgs{
		"userID": userID,
		"limit":  limit,
		"offset": offset,
	}
	return Query[Notification](ctx, db, q, args)
}
