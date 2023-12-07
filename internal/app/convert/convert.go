package convert

import (
	"github.com/dropwhile/icbt/internal/app/model"
	pb "github.com/dropwhile/icbt/rpc"
)

//go:generate convergen
type Convergen interface {
	// :typecast
	// :stringer
	// :case:off
	// :conv TimeToTimestamp Created Created
	ToPbNotification(*model.Notification) *pb.Notification

	// :typecast
	// :stringer
	// :case:off
	// :conv TimeToTimestamp Created Created
	// :conv TimeToTimestampTZ When() When
	ToPbEvent(*model.Event) *pb.Event

	// :typecast
	// :stringer
	// :case:off
	// :conv TimeToTimestamp Created Created
	ToPbEventItem(*model.EventItem) *pb.EventItem
}
