// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package convert

import (
	"context"
	"errors"
	"math"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/errs"
	icbt "github.com/dropwhile/icanbringthat/rpc/icbt/rpc/v1"
)

func ToPbList[T any, V any](converter func(*T) *V, in []*T) []*V {
	out := make([]*V, len(in))
	for i := range in {
		out[i] = converter(in[i])
	}
	return out
}

func ToPbListWithService[T any, V any](
	ctx context.Context,
	converter func(context.Context, service.Servicer, *T) (*V, error),
	svc service.Servicer,
	in []*T,
) ([]*V, error) {
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
	pbtz := icbt.TimestampTZ_builder{
		Ts: timestamppb.New(t),
		Tz: t.Location().String(),
	}.Build()
	return pbtz
}

func ToConnectRpcError(src errs.Error) *connect.Error {
	var rawErr error
	if err := src.Unwrap(); err != nil {
		rawErr = err
	} else {
		rawErr = errors.New(src.Msg())
	}
	code := src.Code()
	if code == errs.NoError {
		code = errs.Unknown
	}
	e := connect.NewError(connect.Code(code), rawErr)
	return e
}

func ToPbPagination(src *service.Pagination) *icbt.PaginationResult {
	dst := icbt.PaginationResult_builder{
		Limit:  uint32(min(max(src.Limit, 0), math.MaxUint32)),  // #nosec G115 -- safe conversion
		Offset: uint32(min(max(src.Offset, 0), math.MaxUint32)), // #nosec G115 -- safe conversion
		Count:  uint32(min(max(src.Count, 0), math.MaxUint32)),  // #nosec G115 -- safe conversion
	}.Build()
	return dst
}

func ToPbEvent(src *model.Event) *icbt.Event {
	dst := icbt.Event_builder{
		RefId:       src.RefID.String(),
		Name:        src.Name,
		Description: src.Description,
		When:        TimeToTimestampTZ(src.When()),
		Archived:    src.Archived,
		Created:     TimeToTimestamp(src.Created),
	}.Build()
	return dst
}

func ToPbEventItem(src *model.EventItem) *icbt.EventItem {
	dst := icbt.EventItem_builder{
		RefId:       src.RefID.String(),
		Description: src.Description,
		Created:     TimeToTimestamp(src.Created),
	}.Build()

	return dst
}

func ToPbNotification(src *model.Notification) *icbt.Notification {
	dst := icbt.Notification_builder{
		RefId:   src.RefID.String(),
		Message: src.Message,
		Created: TimeToTimestamp(src.Created),
	}.Build()

	return dst
}

func ToPbEarmark(ctx context.Context, svc service.Servicer, src *model.Earmark) (*icbt.Earmark, error) {
	eventItem, err := svc.GetEventItemByID(ctx, src.EventItemID)
	if err != nil {
		return nil, err
	}

	emUser, err := svc.GetUserByID(ctx, src.UserID)
	if err != nil {
		return nil, err
	}

	dst := icbt.Earmark_builder{
		RefId:          src.RefID.String(),
		Note:           src.Note,
		EventItemRefId: eventItem.RefID.String(),
		Owner:          emUser.Name,
		Created:        TimeToTimestamp(src.Created),
	}.Build()
	return dst, nil
}
