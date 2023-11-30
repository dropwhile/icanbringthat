package dto

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/dropwhile/icbt/rpc"
)

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

func ToPbList[T any, V any](converter func(*T) *V, in []*T) []*V {
	out := make([]*V, len(in))
	for i := range in {
		out[i] = converter(in[i])
	}
	return out
}
