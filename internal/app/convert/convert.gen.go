// Code generated by github.com/reedom/convergen
// DO NOT EDIT.

package convert

import (
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/rpc/icbt"
)

func ToPbEvent(src *model.Event) (dst *icbt.Event) {
	dst = &icbt.Event{}
	dst.RefId = src.RefID.String()
	dst.Name = src.Name
	dst.Description = src.Description
	dst.When = TimeToTimestampTZ(src.When())
	dst.Archived = src.Archived
	dst.Created = TimeToTimestamp(src.Created)

	return
}

func ToPbEventItem(src *model.EventItem) (dst *icbt.EventItem) {
	dst = &icbt.EventItem{}
	dst.RefId = src.RefID.String()
	dst.Description = src.Description
	dst.Created = TimeToTimestamp(src.Created)

	return
}

func ToPbNotification(src *model.Notification) (dst *icbt.Notification) {
	dst = &icbt.Notification{}
	dst.RefId = src.RefID.String()
	dst.Message = src.Message
	dst.Created = TimeToTimestamp(src.Created)

	return
}

func ToPbPagination(src *service.Pagination) (dst *icbt.PaginationResult) {
	dst = &icbt.PaginationResult{}
	dst.Limit = src.Limit
	dst.Offset = src.Offset
	dst.Count = src.Count

	return
}
