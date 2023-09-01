// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package convert

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/errs"
	"github.com/dropwhile/icanbringthat/rpc/icbt"
)

func ToPbList[T any, V any](converter func(*T) *V, in []*T) []*V {
	out := make([]*V, len(in))
	for i := range in {
		out[i] = converter(in[i])
	}
	return out
}

func ToPbListWithService[T any, V any](ctx context.Context, converter func(context.Context, service.Servicer, *T) (*V, error), svc service.Servicer, in []*T) ([]*V, error) {
	out := make([]*V, len(in))
	var err error
	for i := range in {
		out[i], err = converter(ctx, svc, in[i])
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

func ToPbEarmark(ctx context.Context, svc service.Servicer, src *model.Earmark) (dst *icbt.Earmark, err error) {
	dst = &icbt.Earmark{}
	dst.RefId = src.RefID.String()
	dst.Note = src.Note
	dst.Created = TimeToTimestamp(src.Created)

	eventItem, err := svc.GetEventItemByID(ctx, src.EventItemID)
	if err != nil {
		return nil, err
	}
	emUser, err := svc.GetUserByID(ctx, src.UserID)
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
