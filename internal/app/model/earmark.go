package model

import (
	"context"
	"fmt"
	"time"

	"github.com/dropwhile/refid"
	"github.com/georgysavva/scany/v2/pgxscan"
)

var EarmarkRefIDT = refid.Tagger(4)

type Earmark struct {
	ID           int
	RefID        refid.RefID `db:"ref_id"`
	EventItemID  int         `db:"event_item_id"`
	UserID       int         `db:"user_id"`
	Note         string
	Created      time.Time
	LastModified time.Time  `db:"last_modified"`
	EventItem    *EventItem `db:"-"`
	User         *User      `db:"-"`
}

func (em *Earmark) Insert(ctx context.Context, db PgxHandle) error {
	if em.RefID.IsNil() {
		em.RefID = refid.Must(EarmarkRefIDT.New())
	}
	q := `INSERT INTO earmark_ (ref_id, event_item_id, user_id, note) VALUES ($1, $2, $3, $4) RETURNING *`
	res, err := QueryOneTx[Earmark](ctx, db, q, em.RefID, em.EventItemID, em.UserID, em.Note)
	if err != nil {
		return err
	}
	em.ID = res.ID
	em.RefID = res.RefID
	em.Created = res.Created
	em.LastModified = res.LastModified
	return nil
}

func (em *Earmark) Save(ctx context.Context, db PgxHandle) error {
	q := `UPDATE earmark_ SET note = $1 WHERE id = $2`
	return ExecTx[Earmark](ctx, db, q, em.Note, em.ID)
}

func (em *Earmark) Delete(ctx context.Context, db PgxHandle) error {
	q := `DELETE FROM earmark_ WHERE id = $1`
	return ExecTx[Earmark](ctx, db, q, em.ID)
}

func (em *Earmark) GetEventItem(ctx context.Context, db PgxHandle) (*EventItem, error) {
	eventItem, err := GetEventItemByID(ctx, db, em.EventItemID)
	if err != nil {
		return nil, err
	}
	return eventItem, nil
}

func NewEarmark(ctx context.Context, db PgxHandle, eventItemID, userID int, note string) (*Earmark, error) {
	earmark := &Earmark{
		EventItemID: eventItemID,
		UserID:      userID,
		Note:        note,
	}
	err := earmark.Insert(ctx, db)
	if err != nil {
		return nil, err
	}
	return earmark, nil
}

func GetEarmarkByID(ctx context.Context, db PgxHandle, id int) (*Earmark, error) {
	q := `SELECT * FROM earmark_ WHERE id = $1`
	return QueryOne[Earmark](ctx, db, q, id)
}

func GetEarmarkByRefID(ctx context.Context, db PgxHandle, refID refid.RefID) (*Earmark, error) {
	if !EarmarkRefIDT.HasCorrectTag(refID) {
		err := fmt.Errorf(
			"bad refid type: got %d expected %d",
			refID.Tag(), EarmarkRefIDT.Tag(),
		)
		return nil, err
	}
	q := `SELECT * FROM earmark_ WHERE ref_id = $1`
	return QueryOne[Earmark](ctx, db, q, refID)
}

func GetEarmarkByEventItem(ctx context.Context, db PgxHandle, eventItem *EventItem) (*Earmark, error) {
	q := `SELECT * FROM earmark_ WHERE event_item_id = $1`
	return QueryOne[Earmark](ctx, db, q, eventItem.ID)
}

func GetEarmarksByUser(ctx context.Context, db PgxHandle, user *User) ([]*Earmark, error) {
	q := `SELECT * FROM earmark_ WHERE user_id = $1 ORDER BY created DESC, id DESC`
	return Query[Earmark](ctx, db, q, user.ID)
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
	return Query[Earmark](ctx, db, q, event.ID)
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
			created DESC,
			id DESC
	`
	return Query[Earmark](ctx, db, q, user.ID)
}

func GetEarmarksByUserPaginated(ctx context.Context, db PgxHandle, user *User, limit, offset int) ([]*Earmark, error) {
	q := `SELECT * FROM earmark_ WHERE earmark_.user_id = $1 ORDER BY created DESC, id DESC LIMIT $2 OFFSET $3`
	return Query[Earmark](ctx, db, q, user.ID, limit, offset)
}

func GetEarmarkCountByUser(ctx context.Context, db PgxHandle, user *User) (int, error) {
	q := `SELECT count(*) FROM earmark_ WHERE user_id = $1`
	var count int = 0
	err := pgxscan.Get(ctx, db, &count, q, user.ID)
	return count, err
}
