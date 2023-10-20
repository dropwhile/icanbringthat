package service

import (
	"context"

	"github.com/dropwhile/icbt/internal/app/model"
)

func NotifyUsersPendingEvents(db model.PgxHandle) error {
	ctx := context.Background()
	notifNeeded, err := model.GetUserEventNotificationNeeded(ctx, db)
	if err != nil {
		return err
	}

	for _, elem := range notifNeeded {
		// 1. get earmarked event items (if any)
		// 	  to notify about
		// 2. determine if owner of event or not
		//    a. if owner, send info on all items and their status (as well as
		// 		 any self earmarked items)?
		//    b. if not owner, send info on items earmarked to bring.
		// 3. send appropriate notification
		_ = elem
	}
	return nil
}
