// Code generated by github.com/reedom/convergen
// DO NOT EDIT.

package dto

import (
	"github.com/dropwhile/icbt/internal/app/model"
	pb "github.com/dropwhile/icbt/rpc"
)

func ToPbEvent(src *model.Event) (dst *pb.Event) {
	dst = &pb.Event{}
	dst.RefId = src.RefID.String()
	dst.Name = src.Name
	dst.Description = src.Description
	dst.When = TimeToTimestampTZ(src.When())
	dst.Archived = src.Archived
	dst.Created = TimeToTimestamp(src.Created)

	return
}

func ToPbEventItem(src *model.EventItem) (dst *pb.EventItem) {
	dst = &pb.EventItem{}
	dst.RefId = src.RefID.String()
	dst.Description = src.Description
	dst.Created = TimeToTimestamp(src.Created)

	return
}

func ToPbNotification(src *model.Notification) (dst *pb.Notification) {
	dst = &pb.Notification{}
	dst.RefId = src.RefID.String()
	dst.Message = src.Message
	dst.Created = TimeToTimestamp(src.Created)

	return
}