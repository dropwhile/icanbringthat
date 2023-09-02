package model

import (
	"context"
	"time"

	"github.com/cactus/mlog"
	"github.com/dropwhile/icbt/internal/util/refid"
)

type Event struct {
	Id           uint
	RefId        refid.RefId `db:"ref_id"`
	UserId       uint        `db:"user_id"`
	Name         string
	Description  string
	StartTime    time.Time `db:"start_time"`
	Created      time.Time
	LastModified time.Time    `db:"last_modified"`
	Items        []*EventItem `db:"-"`
}

func (e *Event) Insert(db *DB, ctx context.Context) error {
	if e.RefId.IsNil() {
		e.RefId = EventRefIdT.MustNew()
	}
	q := `INSERT INTO event_ (user_id, ref_id, name, description, start_time) VALUES ($1, $2, $3, $4, $5) RETURNING *`
	res, err := QueryRowTx[Event](db, ctx, q, e.UserId, e.RefId, e.Name, e.Description, e.StartTime)
	if err != nil {
		return err
	}
	e.Id = res.Id
	e.RefId = res.RefId
	e.Created = res.Created
	e.LastModified = res.LastModified
	return nil
}

func (e *Event) Save(db *DB, ctx context.Context) error {
	q := `UPDATE event_ SET name = $1, description = $2, start_time = $3 WHERE id = $4`
	return ExecTx[Event](db, ctx, q, e.Name, e.Description, e.StartTime, e.Id)
}

func (e *Event) Delete(db *DB, ctx context.Context) error {
	q := `DELETE FROM event_ WHERE id = $1`
	return ExecTx[Event](db, ctx, q, e.Id)
}

func NewEvent(db *DB, ctx context.Context, userId uint, name, description string, startTime time.Time) (*Event, error) {
	event := &Event{
		Name:        name,
		UserId:      userId,
		Description: description,
		StartTime:   startTime,
	}
	err := event.Insert(db, ctx)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func GetEventById(db *DB, ctx context.Context, id uint) (*Event, error) {
	q := `SELECT * FROM event_ WHERE id = $1`
	return QueryRow[Event](db, ctx, q, id)
}

func GetEventByRefId(db *DB, ctx context.Context, refId refid.RefId) (*Event, error) {
	q := `SELECT * FROM event_ WHERE ref_id = $1`
	return QueryRow[Event](db, ctx, q, refId)
}

func GetEventsByUser(db *DB, ctx context.Context, user *User) ([]*Event, error) {
	q := `SELECT * FROM event_ WHERE event_.user_id = $1 ORDER BY created DESC`
	return Query[Event](db, ctx, q, user.Id)
}

func GetEventsByUserPaginated(db *DB, ctx context.Context, user *User, limit, offset uint) ([]*Event, error) {
	q := `SELECT * FROM event_ WHERE event_.user_id = $1 ORDER BY start_time DESC LIMIT $2 OFFSET $3`
	return Query[Event](db, ctx, q, user.Id, limit, offset)
}

func GetEventCountByUser(db *DB, ctx context.Context, user *User) (uint, error) {
	q := `SELECT count(*) FROM event_ WHERE user_id = $1`
	if mlog.HasDebug() {
		mlog.Debugx("SQL", mlog.A("query", q), mlog.A("args", user.Id))
	}
	var count uint = 0
	err := db.GetContext(ctx, &count, q, user.Id)
	return count, err
}
