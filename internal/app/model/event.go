package model

import (
	"context"
	"fmt"
	"time"

	"github.com/dropwhile/refid"
	"github.com/georgysavva/scany/v2/pgxscan"
)

var EventRefIdT = refid.Tagger(2)

type Event struct {
	Id           int
	RefId        refid.RefId `db:"ref_id"`
	UserId       int         `db:"user_id"`
	Name         string
	Description  string
	StartTime    time.Time `db:"start_time"`
	StartTimeTZ  string    `db:"start_time_tz"`
	Created      time.Time
	LastModified time.Time    `db:"last_modified"`
	Items        []*EventItem `db:"-"`
}

func (e *Event) Insert(ctx context.Context, db PgxHandle) error {
	if e.RefId.IsNil() {
		e.RefId = refid.Must(EventRefIdT.New())
	}
	q := `INSERT INTO event_ (user_id, ref_id, name, description, start_time, start_time_tz) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`
	res, err := QueryOneTx[Event](ctx, db, q, e.UserId, e.RefId, e.Name, e.Description, e.StartTime, e.StartTimeTZ)
	if err != nil {
		return err
	}
	e.Id = res.Id
	e.RefId = res.RefId
	e.Created = res.Created
	e.LastModified = res.LastModified
	return nil
}

func (e *Event) Save(ctx context.Context, db PgxHandle) error {
	q := `UPDATE event_ SET name = $1, description = $2, start_time = $3, start_time_tz = $4 WHERE id = $5`
	return ExecTx[Event](ctx, db, q, e.Name, e.Description, e.StartTime, e.StartTimeTZ, e.Id)
}

func (e *Event) Delete(ctx context.Context, db PgxHandle) error {
	q := `DELETE FROM event_ WHERE id = $1`
	return ExecTx[Event](ctx, db, q, e.Id)
}

func NewEvent(ctx context.Context, db PgxHandle, userId int, name, description string, startTime time.Time, startTimeTZ string) (*Event, error) {
	event := &Event{
		Name:        name,
		UserId:      userId,
		Description: description,
		StartTime:   startTime,
		StartTimeTZ: startTimeTZ,
	}
	err := event.Insert(ctx, db)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func GetEventById(ctx context.Context, db PgxHandle, id int) (*Event, error) {
	q := `SELECT * FROM event_ WHERE id = $1`
	return QueryOne[Event](ctx, db, q, id)
}

func GetEventByRefId(ctx context.Context, db PgxHandle, refId refid.RefId) (*Event, error) {
	if !EventRefIdT.HasCorrectTag(refId) {
		err := fmt.Errorf(
			"bad refid type: got %d expected %d",
			refId.Tag(), EventRefIdT.Tag(),
		)
		return nil, err
	}
	q := `SELECT * FROM event_ WHERE ref_id = $1`
	return QueryOne[Event](ctx, db, q, refId)
}

func GetEventsByUser(ctx context.Context, db PgxHandle, user *User) ([]*Event, error) {
	q := `SELECT * FROM event_ WHERE event_.user_id = $1 ORDER BY start_time DESC, id DESC`
	return Query[Event](ctx, db, q, user.Id)
}

func GetEventsByUserPaginated(ctx context.Context, db PgxHandle, user *User, limit, offset int) ([]*Event, error) {
	q := `SELECT * FROM event_ WHERE event_.user_id = $1 ORDER BY start_time DESC, id DESC LIMIT $2 OFFSET $3`
	return Query[Event](ctx, db, q, user.Id, limit, offset)
}

func GetEventsComingSoonByUserPaginated(ctx context.Context, db PgxHandle, user *User, limit, offset int) ([]*Event, error) {
	q := `
		SELECT *
		FROM event_
		WHERE
			event_.user_id = $1
			start_time > CURRENT_TIMESTAMP(0)
		ORDER BY 
			start_time ASC,
			id ASC
		LIMIT $2 OFFSET $3`
	return Query[Event](ctx, db, q, user.Id, limit, offset)
}

func GetEventCountByUser(ctx context.Context, db PgxHandle, user *User) (int, error) {
	q := `SELECT count(*) FROM event_ WHERE user_id = $1`
	var count int = 0
	err := pgxscan.Get(ctx, db, &count, q, user.Id)
	return count, err
}
