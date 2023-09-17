package model

import (
	"context"
	"time"

	"github.com/dropwhile/icbt/internal/util/refid"
	"github.com/georgysavva/scany/v2/pgxscan"
)

type Earmark struct {
	Id           int
	RefId        refid.RefId `db:"ref_id"`
	EventItemId  int         `db:"event_item_id"`
	UserId       int         `db:"user_id"`
	Note         string
	Created      time.Time
	LastModified time.Time  `db:"last_modified"`
	EventItem    *EventItem `db:"-"`
	User         *User      `db:"-"`
}

func (em *Earmark) Insert(ctx context.Context, db PgxHandle) error {
	if em.RefId.IsNil() {
		em.RefId = EarmarkRefIdT.MustNew()
	}
	q := `INSERT INTO earmark_ (ref_id, event_item_id, user_id, note) VALUES ($1, $2, $3, $4) RETURNING *`
	res, err := QueryOneTx[Earmark](ctx, db, q, em.RefId, em.EventItemId, em.UserId, em.Note)
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
	q := `UPDATE earmark_ SET note = $1 WHERE id = $2`
	return ExecTx[Earmark](ctx, db, q, em.Note, em.Id)
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

func NewEarmark(ctx context.Context, db PgxHandle, eventItemId int, userId int, note string) (*Earmark, error) {
	earmark := &Earmark{
		EventItemId: eventItemId,
		UserId:      userId,
		Note:        note,
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

func GetEarmarksByEvent(ctx context.Context, db PgxHandle, event *Event) ([]*Earmark, error) {
	q := `
		SELECT earmark_.*
		FROM earmark_
		JOIN event_item_ ON 
			event_item_.id = earmark_.event_item_id
		WHERE 
			event_item_.event_id = $1
	`
	return Query[Earmark](ctx, db, q, event.Id)
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
	var count int = 0
	err := pgxscan.Get(ctx, db, &count, q, user.Id)
	return count, err
}
