package dto

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
}

/*
//  go:generate goverter gen ./...

// goverter:converter
// goverter:name DTOConverter
// goverter:output:file ./converter/generated.go
// goverter:output:package github.com/dropwhile/icbt/rpc/converter
// goverter:matchIgnoreCase
// goverter:extend TimeToTimestamp
// goverter:extend StringerToString
type Converter interface {
	ConvertNotifications(source []*model.Notification) []*rpc.Notification

	// goverter:ignoreUnexported
	// goverter:map . RefId | GetRefId
	ConvertNotification(source *model.Notification) *rpc.Notification
}

func TimeToTimestamp(t time.Time) timestamppb.Timestamp {
	return *timestamppb.New(t)
}

func GetRefId(source *model.Notification) string {
	return source.RefID.String()
}

func StringerToString(src interface{ String() string }) string {
	return src.String()
}
*/
