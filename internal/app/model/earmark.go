package model

import (
	"context"
	"time"

	"github.com/cactus/mlog"
	"github.com/dropwhile/icbt/internal/util/refid"
)

type Earmark struct {
	Id           uint
	RefId        refid.RefId `db:"ref_id"`
	EventItemId  uint        `db:"event_item_id"`
	UserId       uint        `db:"user_id"`
	Notes        string
	Created      time.Time
	LastModified time.Time  `db:"last_modified"`
	EventItem    *EventItem `db:"-"`
}

func (em *Earmark) Insert(db *DB, ctx context.Context) error {
	if em.RefId.IsNil() {
		em.RefId = EarmarkRefIdT.MustNew()
	}
	q := `INSERT INTO earmark_ (event_item_id, user_id, notes) VALUES ($1, $2, $3) RETURNING *`
	res, err := QueryRowTx[Earmark](db, ctx, q, em.EventItemId, em.UserId, em.Notes)
	if err != nil {
		return err
	}
	em.Id = res.Id
	em.RefId = res.RefId
	em.Created = res.Created
	em.LastModified = res.LastModified
	return nil
}

func (em *Earmark) Save(db *DB, ctx context.Context) error {
	q := `UPDATE earmark_ SET notes = $1 WHERE id = $2`
	return ExecTx[Earmark](db, ctx, q, em.Notes, em.Id)
}

func (em *Earmark) Delete(db *DB, ctx context.Context) error {
	q := `DELETE FROM earmark_ WHERE id = $1`
	return ExecTx[Earmark](db, ctx, q, em.Id)
}

func (em *Earmark) GetEventItem(db *DB, ctx context.Context) (*EventItem, error) {
	eventItem, err := GetEventItemById(db, ctx, em.EventItemId)
	if err != nil {
		return nil, err
	}
	return eventItem, nil
}

func NewEarmark(db *DB, ctx context.Context, eventItemId uint, userId uint, notes string) (*Earmark, error) {
	earmark := &Earmark{
		EventItemId: eventItemId,
		UserId:      userId,
		Notes:       notes,
	}
	err := earmark.Insert(db, ctx)
	if err != nil {
		return nil, err
	}
	return earmark, nil
}

func GetEarmarkById(db *DB, ctx context.Context, id uint) (*Earmark, error) {
	q := `SELECT * FROM earmark_ WHERE id = $1`
	return QueryRow[Earmark](db, ctx, q, id)
}

func GetEarmarkByRefId(db *DB, ctx context.Context, refId refid.RefId) (*Earmark, error) {
	q := `SELECT * FROM earmark_ WHERE ref_id = $1`
	return QueryRow[Earmark](db, ctx, q, refId)
}

func GetEarmarkByEventItem(db *DB, ctx context.Context, eventItem *EventItem) (*Earmark, error) {
	q := `SELECT * FROM earmark_ WHERE event_item_id = $1`
	return QueryRow[Earmark](db, ctx, q, eventItem.Id)
}

func GetEarmarksByUser(db *DB, ctx context.Context, user *User) ([]*Earmark, error) {
	q := `SELECT * FROM earmark_ WHERE user_id = $1 ORDER BY created DESC`
	return Query[Earmark](db, ctx, q, user.Id)
}

func GetEarmarksWithEventsByUser(db *DB, ctx context.Context, user *User) ([]*Earmark, error) {
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
	return Query[Earmark](db, ctx, q, user.Id)
}

func GetEarmarksByUserPaginated(db *DB, ctx context.Context, user *User, limit, offset uint) ([]*Earmark, error) {
	q := `SELECT * FROM earmark_ WHERE earmark_.user_id = $1 ORDER BY created DESC LIMIT $2 OFFSET $3`
	return Query[Earmark](db, ctx, q, user.Id, limit, offset)
}

func GetEarmarkCountByUser(db *DB, ctx context.Context, user *User) (uint, error) {
	q := `SELECT count(*) FROM earmark_ WHERE user_id = $1`
	if mlog.HasDebug() {
		mlog.Debugx("SQL", mlog.A("query", q), mlog.A("args", user.Id))
	}
	var count uint = 0
	err := db.GetContext(ctx, &count, q, user.Id)
	return count, err
}
