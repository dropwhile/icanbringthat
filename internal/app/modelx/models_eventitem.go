package modelx

import "context"

type EventItemExpanded struct {
	*EventItem
	Event   *Event
	Earmark *EarmarkExpanded
}

func (q *Queries) NewEventItem(ctx context.Context, eventID int32, description string) (*EventItem, error) {
	refID, err := NewEventItemRefID()
	if err != nil {
		return nil, err
	}
	params := CreateEventItemParams{
		RefID:       refID,
		EventID:     eventID,
		Description: description,
	}
	return q.CreateEventItem(ctx, params)
}
