package convert

import (
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/rpc/icbt"
)

//go:generate convergen
type Convergen interface {
	// :typecast
	// :stringer
	// :case:off
	// :conv TimeToTimestamp Created Created
	ToPbNotification(*model.Notification) *icbt.Notification

	// :typecast
	// :stringer
	// :case:off
	// :conv TimeToTimestamp Created Created
	// :conv TimeToTimestampTZ When() When
	ToPbEvent(*model.Event) *icbt.Event

	// :typecast
	// :stringer
	// :case:off
	// :conv TimeToTimestamp Created Created
	ToPbEventItem(*model.EventItem) *icbt.EventItem

	// :typecast
	// :stringer
	// :case:off
	ToPbPagination(*service.Pagination) *icbt.PaginationResult
}
