package dto

import (
	"time"

	timestamppb "google.golang.org/protobuf/types/known/timestamppb"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/rpc"
)

//go:generate goverter gen ./...

// goverter:converter
// goverter:name DTOConverter
// goverter:output:file ./converter/generated.go
// goverter:output:package github.com/dropwhile/icbt/rpc/converter
// goverter:extend TimeToTimestamp
type Converter interface {
	ConvertNotifications(source []*model.Notification) []*rpc.Notification

	// goverter:ignoreUnexported
	// goverter:map . RefId | GetRefId
	ConvertNotification(source *model.Notification) *rpc.Notification
}

func TimeToTimestamp(t time.Time) timestamppb.Timestamp {
	return *timestamppb.New(t)
}

func GetRefId(source model.Notification) string {
	return source.RefID.String()
}
