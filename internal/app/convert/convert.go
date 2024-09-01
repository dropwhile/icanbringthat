// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package convert

import (
	"github.com/dropwhile/icanbringthat/internal/app/model"
	icbt "github.com/dropwhile/icanbringthat/rpc/icbt/rpc/v1"
)

//go:generate convergen
type Convergen interface {
	// :typecast
	// :stringer
	// :case:off
	// :conv TimeToTimestamp Created Created
	ToPbNotification(*model.Notification) *icbt.Notification

	// :typecast
	// :stringer
	// :case:off
	// :conv TimeToTimestamp Created Created
	// :conv TimeToTimestampTZ When() When
	ToPbEvent(*model.Event) *icbt.Event

	// :typecast
	// :stringer
	// :case:off
	// :conv TimeToTimestamp Created Created
	ToPbEventItem(*model.EventItem) *icbt.EventItem
}
