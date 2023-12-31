package service

import (
	"context"

	"github.com/dropwhile/icbt/internal/app/model"
)

func (s *Service) ArchiveOldEvents(ctx context.Context) error {
	return model.ArchiveOldEvents(ctx, s.Db)
}
