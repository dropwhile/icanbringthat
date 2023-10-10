package modelx

import "context"

type EarmarkExpanded struct {
	*Earmark
	EventItem *EventItemExpanded
}

func (q *Queries) NewEarmark(ctx context.Context, eventItemID int32, userID int32, note string) (*Earmark, error) {
	refID, err := NewEarmarkRefID()
	if err != nil {
		return nil, err
	}
	params := CreateEarmarkParams{
		RefID:       refID,
		EventItemID: eventItemID,
		UserID:      userID,
		Note:        note,
	}
	return q.CreateEarmark(ctx, params)
}
