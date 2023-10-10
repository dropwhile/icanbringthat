package modelx

import (
	"context"
	"time"
)

type EventExpanded struct {
	*Event
	Items []*EventItemExpanded
}

func (q *Queries) NewEvent(ctx context.Context, userID int32, name, description string, startTime time.Time, startTimeTZ string) (*Event, error) {
	refID, err := NewEventRefID()
	if err != nil {
		return nil, err
	}
	tz := TimeZone{}
	err = tz.Scan(startTimeTZ)
	if err != nil {
		return nil, err
	}
	params := CreateEventParams{
		UserID:        userID,
		RefID:         refID,
		Name:          name,
		Description:   description,
		ItemSortOrder: []int32{},
		StartTime:     startTime,
		StartTimeTz:   tz,
	}
	return q.CreateEvent(ctx, params)
}
