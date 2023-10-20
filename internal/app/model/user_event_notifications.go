package model

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type UserEventNotification struct {
	UserID  int `db:"user_id"`
	EventID int `db:"event_id"`
	Created time.Time
}

func NewUserEventNotification(ctx context.Context, db PgxHandle,
	userID int, eventID int,
) (*UserEventNotification, error) {
	return CreateUserEventNotification(ctx, db, userID, eventID)
}

func CreateUserEventNotification(ctx context.Context, db PgxHandle,
	userID int, eventID int,
) (*UserEventNotification, error) {
	q := `
		INSERT INTO user_event_notification_ (
			user_id, event_id
		)
		VALUES (@userID, @eventID)
		RETURNING *`
	args := pgx.NamedArgs{"userID": userID, "eventID": eventID}
	return QueryOneTx[UserEventNotification](ctx, db, q, args)
}

func DeleteUserEventNotification(ctx context.Context, db PgxHandle,
	userID int, eventID int,
) error {
	q := `
		DELETE FROM user_event_notification_
		WHERE 
			user_id = @userID AND
			event_id = @eventID`
	args := pgx.NamedArgs{"userID": userID, "eventID": eventID}
	return ExecTx[UserEventNotification](ctx, db, q, args)
}

func GetUserEventNotification(ctx context.Context, db PgxHandle,
	userID int, eventID int,
) (*UserEventNotification, error) {
	q := `
		SELECT * FROM user_event_notification_
		WHERE 
			user_id = @userID AND
			event_id = @eventID`
	args := pgx.NamedArgs{"userID": userID, "eventID": eventID}
	return QueryOne[UserEventNotification](ctx, db, q, args)
}

type UserEventNotificationNeeded struct {
	UserID       int `db:"user_id"`
	EventID      int `db:"event_id"`
	When         time.Time
	Owner        bool
	EventItemIDs []int `db:"items"`
}

func GetUserEventNotificationNeeded(
	ctx context.Context, db PgxHandle,
) ([]*UserEventNotificationNeeded, error) {
	q := `
		WITH subt as (
			SELECT
				user_id,
				id as event_id,
				TRUE as owner,
				date_trunc('hour', start_time) AT TIME ZONE 'UTC' as when,
				NULL as item
			FROM event_
			WHERE
				date_trunc('hour', start_time) AT TIME ZONE 'UTC' > timezone('utc', CURRENT_TIMESTAMP)

			UNION

			SELECT DISTINCT
				u.id as user_id,
				ev.id as event_id,
				FALSE as owner,
				date_trunc('hour', ev.start_time) AT TIME ZONE 'UTC' as when,
				ei.id as item
			FROM user_ u
			JOIN earmark_ em
				ON u.id = em.user_id
			JOIN event_item_ ei
				ON em.event_item_id = ei.id
			JOIN event_ ev
				ON ei.event_id = ev.id
			WHERE
				date_trunc('hour', ev.start_time) AT TIME ZONE 'UTC' > timezone('utc', CURRENT_TIMESTAMP)
		)
		SELECT
			subt.user_id,
			subt.event_id,
			subt.when,
			bool_or(subt.owner) as owner,
			ARRAY_AGG(subt.item ORDER BY subt.item) FILTER(WHERE subt.item IS NOT NULL) as items
		FROM subt
		LEFT JOIN user_event_notification_ uen ON
			uen.user_id = subt.user_id AND
			uen.event_id = subt.event_id
		JOIN user_ u ON
			u.id = subt.user_id
		WHERE
			uen.user_id is NULL AND
			u.verified = TRUE AND
			(u.settings->>'disable_reminders')::boolean IS NOT TRUE
		GROUP BY (subt.user_id, subt.event_id, subt.when)
	`
	return Query[UserEventNotificationNeeded](ctx, db, q)
}
