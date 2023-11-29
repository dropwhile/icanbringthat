// Code generated by github.com/jmattheis/goverter, DO NOT EDIT.

package converter

import (
	model "github.com/dropwhile/icbt/internal/app/model"
	dto "github.com/dropwhile/icbt/internal/app/rpc/dto"
	rpc "github.com/dropwhile/icbt/rpc"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type DTOConverter struct{}

func (c *DTOConverter) ConvertNotification(source *model.Notification) *rpc.Notification {
	var pRpcNotification *rpc.Notification
	if source != nil {
		var rpcNotification rpc.Notification
		rpcNotification.RefId = dto.GetRefId(source)
		rpcNotification.Message = (*source).Message
		rpcNotification.Created = c.timeTimeToPTimestamppbTimestamp((*source).Created)
		pRpcNotification = &rpcNotification
	}
	return pRpcNotification
}
func (c *DTOConverter) ConvertNotifications(source []*model.Notification) []*rpc.Notification {
	var pRpcNotificationList []*rpc.Notification
	if source != nil {
		pRpcNotificationList = make([]*rpc.Notification, len(source))
		for i := 0; i < len(source); i++ {
			pRpcNotificationList[i] = c.ConvertNotification(source[i])
		}
	}
	return pRpcNotificationList
}
func (c *DTOConverter) timeTimeToPTimestamppbTimestamp(source time.Time) *timestamppb.Timestamp {
	timestamppbTimestamp := dto.TimeToTimestamp(source)
	return &timestamppbTimestamp
}
