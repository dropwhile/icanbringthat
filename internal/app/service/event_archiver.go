package service

import (
	"context"

	"github.com/dropwhile/icbt/internal/app/model"
)

func ArchiveOldEvents(db model.PgxHandle) error {
	ctx := context.Background()
	return model.ArchiveOldEvents(ctx, db)
}
