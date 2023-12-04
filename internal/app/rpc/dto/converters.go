package dto

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/someerr"
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

func ToTwirpError(src someerr.Error) twirp.Error {
	var twErrCode twirp.ErrorCode
	switch src.Code() {
	/*
		case someerr.BadRoute:
			errString = "bad_route"
		case someerr.Malformed:
			errString = "malformed"
	*/
	case someerr.NoError:
		twErrCode = twirp.NoError
	case someerr.Canceled:
		twErrCode = twirp.Canceled
	case someerr.Unknown:
		twErrCode = twirp.Unknown
	case someerr.InvalidArgument:
		twErrCode = twirp.InvalidArgument
	case someerr.DeadlineExceeded:
		twErrCode = twirp.DeadlineExceeded
	case someerr.NotFound:
		twErrCode = twirp.NotFound
	case someerr.AlreadyExists:
		twErrCode = twirp.AlreadyExists
	case someerr.PermissionDenied:
		twErrCode = twirp.PermissionDenied
	case someerr.Unauthenticated:
		twErrCode = twirp.Unauthenticated
	case someerr.ResourceExhausted:
		twErrCode = twirp.ResourceExhausted
	case someerr.FailedPrecondition:
		twErrCode = twirp.FailedPrecondition
	case someerr.Aborted:
		twErrCode = twirp.Aborted
	case someerr.OutOfRange:
		twErrCode = twirp.OutOfRange
	case someerr.Unimplemented:
		twErrCode = twirp.Unimplemented
	case someerr.Internal:
		twErrCode = twirp.Internal
	case someerr.Unavailable:
		twErrCode = twirp.Unavailable
	case someerr.DataLoss:
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
