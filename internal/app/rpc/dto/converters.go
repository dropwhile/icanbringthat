package dto

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/dropwhile/icbt/internal/app/model"
	pb "github.com/dropwhile/icbt/rpc"
)

func ToPbList[T any, V any](converter func(*T) *V, in []*T) []*V {
	out := make([]*V, len(in))
	for i := range in {
		out[i] = converter(in[i])
	}
	return out
}

func ToPbListWithDb[T any, V any](converter func(*T, model.PgxHandle) (*V, error), db model.PgxHandle, in []*T) ([]*V, error) {
	out := make([]*V, len(in))
	var err error
	for i := range in {
		out[i], err = converter(in[i], db)
		if err != nil {
			break
		}
	}
	return out, err
}

func TimeToTimestamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

func TimeToTimestampTZ(t time.Time) *pb.TimestampTZ {
	pbtz := &pb.TimestampTZ{
		Ts: timestamppb.New(t),
		Tz: t.Location().String(),
	}
	return pbtz
}

func ToPbEarmark(src *model.Earmark, db model.PgxHandle) (dst *pb.Earmark, err error) {
	dst = &pb.Earmark{}
	dst.RefId = src.RefID.String()
	dst.Note = src.Note
	dst.Created = TimeToTimestamp(src.Created)

	ctx := context.Background()
	eventItem, err := model.GetEventItemByID(ctx, db, src.EventItemID)
	if err != nil {
		return nil, err
	}
	emUser, err := model.GetUserByID(ctx, db, src.UserID)
	if err != nil {
		return nil, err
	}

	dst.EventItemRefId = eventItem.RefID.String()
	dst.Owner = emUser.Name
	return
}
