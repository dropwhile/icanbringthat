package model

import (
	"context"
	"time"

	"github.com/cactus/mlog"
	"github.com/dropwhile/icbt/internal/util/refid"
	"github.com/georgysavva/scany/v2/pgxscan"
)

type Earmark struct {
	Id           int
	RefId        refid.RefId `db:"ref_id"`
	EventItemId  int         `db:"event_item_id"`
	UserId       int         `db:"user_id"`
	Notes        string
	Created      time.Time
	LastModified time.Time  `db:"last_modified"`
	EventItem    *EventItem `db:"-"`
}

func (em *Earmark) Insert(ctx context.Context, db PgxHandle) error {
	if em.RefId.IsNil() {
		em.RefId = EarmarkRefIdT.MustNew()
	}
	q := `INSERT INTO earmark_ (event_item_id, user_id, notes) VALUES ($1, $2, $3) RETURNING *`
	res, err := QueryOneTx[Earmark](ctx, db, q, em.EventItemId, em.UserId, em.Notes)
	if err != nil {
		return err
	}
	em.Id = res.Id
	em.RefId = res.RefId
	em.Created = res.Created
	em.LastModified = res.LastModified
	return nil
}

func (em *Earmark) Save(ctx context.Context, db PgxHandle) error {
	q := `UPDATE earmark_ SET notes = $1 WHERE id = $2`
	return ExecTx[Earmark](ctx, db, q, em.Notes, em.Id)
}

func (em *Earmark) Delete(ctx context.Context, db PgxHandle) error {
	q := `DELETE FROM earmark_ WHERE id = $1`
	return ExecTx[Earmark](ctx, db, q, em.Id)
}

func (em *Earmark) GetEventItem(ctx context.Context, db PgxHandle) (*EventItem, error) {
	eventItem, err := GetEventItemById(ctx, db, em.EventItemId)
	if err != nil {
		return nil, err
	}
	return eventItem, nil
}

func NewEarmark(ctx context.Context, db PgxHandle, eventItemId int, userId int, notes string) (*Earmark, error) {
	earmark := &Earmark{
		EventItemId: eventItemId,
		UserId:      userId,
		Notes:       notes,
	}
	err := earmark.Insert(ctx, db)
	if err != nil {
		return nil, err
	}
	return earmark, nil
}

func GetEarmarkById(ctx context.Context, db PgxHandle, id int) (*Earmark, error) {
	q := `SELECT * FROM earmark_ WHERE id = $1`
	return QueryOne[Earmark](ctx, db, q, id)
}

func GetEarmarkByRefId(ctx context.Context, db PgxHandle, refId refid.RefId) (*Earmark, error) {
	q := `SELECT * FROM earmark_ WHERE ref_id = $1`
	return QueryOne[Earmark](ctx, db, q, refId)
}

func GetEarmarkByEventItem(ctx context.Context, db PgxHandle, eventItem *EventItem) (*Earmark, error) {
	q := `SELECT * FROM earmark_ WHERE event_item_id = $1`
	return QueryOne[Earmark](ctx, db, q, eventItem.Id)
}

func GetEarmarksByUser(ctx context.Context, db PgxHandle, user *User) ([]*Earmark, error) {
	q := `SELECT * FROM earmark_ WHERE user_id = $1 ORDER BY created DESC`
	return Query[Earmark](ctx, db, q, user.Id)
}

func GetEarmarksWithEventsByUser(ctx context.Context, db PgxHandle, user *User) ([]*Earmark, error) {
	q := `
		SELECT *
		FROM earmark_
		JOIN event_ ON
			event_.id = earmark_.id 
		WHERE 
			user_id = $1
		ORDER BY
			created DESC
	`
	return Query[Earmark](ctx, db, q, user.Id)
}

func GetEarmarksByUserPaginated(ctx context.Context, db PgxHandle, user *User, limit, offset int) ([]*Earmark, error) {
	q := `SELECT * FROM earmark_ WHERE earmark_.user_id = $1 ORDER BY created DESC LIMIT $2 OFFSET $3`
	return Query[Earmark](ctx, db, q, user.Id, limit, offset)
}

func GetEarmarkCountByUser(ctx context.Context, db PgxHandle, user *User) (int, error) {
	q := `SELECT count(*) FROM earmark_ WHERE user_id = $1`
	if mlog.HasDebug() {
		mlog.Debugx("SQL", mlog.A("query", q), mlog.A("args", user.Id))
	}
	var count int = 0
	err := pgxscan.Get(ctx, db, &count, q, user.Id)
	return count, err
}
