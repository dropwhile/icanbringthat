// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package service

import (
	"context"

	"github.com/dropwhile/icanbringthat/internal/app/model"
)

func (s *Service) ArchiveOldEvents(ctx context.Context) error {
	return model.ArchiveOldEvents(ctx, s.Db)
}
