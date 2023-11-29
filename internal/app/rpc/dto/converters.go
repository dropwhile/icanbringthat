package dto

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/dropwhile/icbt/internal/app/model"
	pb "github.com/dropwhile/icbt/rpc"
)

func TimeToTimestamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

func ToPbNotificationList(notifications []*model.Notification) []*pb.Notification {
	out := make([]*pb.Notification, len(notifications))
	for i := range notifications {
		out[i] = ToPbNotification(notifications[i])
	}
	return out
}
