package dto

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/somerr"
	pb "github.com/dropwhile/icbt/rpc"
)

func ToPbList[T any, V any](converter func(*T) *V, in []*T) []*V {
	out := make([]*V, len(in))
	for i := range in {
		out[i] = converter(in[i])
	}
	return out
}

func ToPbListWithDb[T any, V any](converter func(model.PgxHandle, *T) (*V, error), db model.PgxHandle, in []*T) ([]*V, error) {
	out := make([]*V, len(in))
	var err error
	for i := range in {
		out[i], err = converter(db, in[i])
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

func ToPbEarmark(db model.PgxHandle, src *model.Earmark) (dst *pb.Earmark, err error) {
	dst = &pb.Earmark{}
	dst.RefId = src.RefID.String()
	dst.Note = src.Note
	dst.Created = TimeToTimestamp(src.Created)

	ctx := context.Background()
	eventItem, err := service.GetEventItemByID(ctx, db, src.EventItemID)
	if err != nil {
		return nil, err
	}
	emUser, err := service.GetUserByID(ctx, db, src.UserID)
	if err != nil {
		return nil, err
	}

	dst.EventItemRefId = eventItem.RefID.String()
	dst.Owner = emUser.Name
	return
}

func ToTwirpError(src somerr.Error) twirp.Error {
	var twErrCode twirp.ErrorCode
	switch src.Code() {
	/*
		case somerr.BadRoute:
			errString = "bad_route"
		case somerr.Malformed:
			errString = "malformed"
	*/
	case somerr.NoError:
		twErrCode = twirp.NoError
	case somerr.Canceled:
		twErrCode = twirp.Canceled
	case somerr.Unknown:
		twErrCode = twirp.Unknown
	case somerr.InvalidArgument:
		twErrCode = twirp.InvalidArgument
	case somerr.DeadlineExceeded:
		twErrCode = twirp.DeadlineExceeded
	case somerr.NotFound:
		twErrCode = twirp.NotFound
	case somerr.AlreadyExists:
		twErrCode = twirp.AlreadyExists
	case somerr.PermissionDenied:
		twErrCode = twirp.PermissionDenied
	case somerr.Unauthenticated:
		twErrCode = twirp.Unauthenticated
	case somerr.ResourceExhausted:
		twErrCode = twirp.ResourceExhausted
	case somerr.FailedPrecondition:
		twErrCode = twirp.FailedPrecondition
	case somerr.Aborted:
		twErrCode = twirp.Aborted
	case somerr.OutOfRange:
		twErrCode = twirp.OutOfRange
	case somerr.Unimplemented:
		twErrCode = twirp.Unimplemented
	case somerr.Internal:
		twErrCode = twirp.Internal
	case somerr.Unavailable:
		twErrCode = twirp.Unavailable
	case somerr.DataLoss:
		twErrCode = twirp.DataLoss
	}
	twerr := twirp.NewError(twErrCode, src.Msg())
	for k, v := range src.MetaMap() {
		twerr = twerr.WithMeta(k, v)
	}
	if u, ok := src.(interface {
		Unwrap() error
	}); ok {
		twerr = twirp.WrapError(twerr, u.Unwrap())
	}
	return twerr
}
