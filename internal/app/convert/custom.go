package convert

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/rpc/icbt"
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

func TimeToTimestampTZ(t time.Time) *icbt.TimestampTZ {
	pbtz := &icbt.TimestampTZ{
		Ts: timestamppb.New(t),
		Tz: t.Location().String(),
	}
	return pbtz
}

func ToPbEarmark(db model.PgxHandle, src *model.Earmark) (dst *icbt.Earmark, err error) {
	dst = &icbt.Earmark{}
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

func ToTwirpError(src errs.Error) twirp.Error {
	var twErrCode twirp.ErrorCode
	switch src.Code() {
	/*
		case somerr.BadRoute:
			errString = "bad_route"
		case somerr.Malformed:
			errString = "malformed"
	*/
	case errs.NoError:
		twErrCode = twirp.NoError
	case errs.Canceled:
		twErrCode = twirp.Canceled
	case errs.Unknown:
		twErrCode = twirp.Unknown
	case errs.InvalidArgument:
		twErrCode = twirp.InvalidArgument
	case errs.DeadlineExceeded:
		twErrCode = twirp.DeadlineExceeded
	case errs.NotFound:
		twErrCode = twirp.NotFound
	case errs.AlreadyExists:
		twErrCode = twirp.AlreadyExists
	case errs.PermissionDenied:
		twErrCode = twirp.PermissionDenied
	case errs.Unauthenticated:
		twErrCode = twirp.Unauthenticated
	case errs.ResourceExhausted:
		twErrCode = twirp.ResourceExhausted
	case errs.FailedPrecondition:
		twErrCode = twirp.FailedPrecondition
	case errs.Aborted:
		twErrCode = twirp.Aborted
	case errs.OutOfRange:
		twErrCode = twirp.OutOfRange
	case errs.Unimplemented:
		twErrCode = twirp.Unimplemented
	case errs.Internal:
		twErrCode = twirp.Internal
	case errs.Unavailable:
		twErrCode = twirp.Unavailable
	case errs.DataLoss:
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
