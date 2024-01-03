package service

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/resources"
	"github.com/dropwhile/icbt/internal/mail"
)

func (s *Service) NotifyUsersPendingEvents(ctx context.Context,
	mailer mail.MailSender, tplContainer resources.TGetter,
	siteBaseUrl string,
) error {
	notifNeeded, err := model.GetUserEventNotificationNeeded(ctx, s.Db)
	if err != nil {
		return err
	}

	tplHtml, err := tplContainer.Get("mail_reminder.gohtml")
	if err != nil {
		return fmt.Errorf("template get error: %w", err)
	}
	tplPlain, err := tplContainer.Get("mail_reminder.gotxt")
	if err != nil {
		return fmt.Errorf("template get error: %w", err)
	}

	for _, elem := range notifNeeded {
		// get user
		user, err := model.GetUserByID(ctx, s.Db, elem.UserID)
		if err != nil {
			return err
		}

		// double check if user wants any reminders
		if !user.Settings.EnableReminders {
			continue
		}

		// check if notification threshold reached
		remT := user.Settings.ReminderThresholdHours
		if remT == 0 {
			remT = 24
		}
		notifyWhen := time.Now().Add(time.Duration(remT) * time.Hour)
		if notifyWhen.Before(elem.When) {
			continue
		}

		// get event
		event, err := model.GetEventByID(ctx, s.Db, elem.EventID)
		if err != nil {
			return err
		}

		// get eventItems and earmarks
		var eventItems []*model.EventItem
		var earmarks []*model.Earmark
		if len(elem.EventItemIDs) == 0 {
			eventItems = make([]*model.EventItem, 0)
			earmarks = make([]*model.Earmark, 0)
		} else {
			eventItems, err = model.GetEventItemsByIDs(ctx, s.Db, elem.EventItemIDs)
			if err != nil {
				return err
			}
			earmarks, err = model.GetEarmarksByEventItemIDs(ctx, s.Db, elem.EventItemIDs)
			if err != nil {
				return err
			}
		}

		owner := false
		if user.ID == event.UserID {
			owner = true
		}

		// 2. determine if owner of event or not
		//    a. if owner, send info on all items and their status (as well as
		// 		 any self earmarked items)?
		//    b. if not owner, send info on items earmarked to bring.
		// 3. send appropriate notification
		eventURL, err := url.JoinPath(
			siteBaseUrl,
			fmt.Sprintf("/events/%s", event.RefID.String()),
		)
		if err != nil {
			return fmt.Errorf("url path join error: %w", err)
		}

		eventWhen := event.StartTime.
			In(event.StartTimeTz.Location).
			Format("2006-01-02 03:04PM")

		vars := map[string]any{
			"Subject":          "Upcoming Event Reminder",
			"owner":            owner,
			"eventName":        event.Name,
			"eventDescription": event.Description,
			"eventWhen":        eventWhen,
			"eventURL":         eventURL,
			"items":            eventItems,
			"earmarks":         earmarks,
		}

		var bufHtml bytes.Buffer
		err = tplHtml.Execute(&bufHtml, vars)
		if err != nil {
			return fmt.Errorf("html template exec error: %w", err)
		}

		var bufPlain bytes.Buffer
		err = tplPlain.Execute(&bufPlain, vars)
		if err != nil {
			return fmt.Errorf("plain template exec error: %w", err)
		}

		messagePlain := bufPlain.String()
		messageHtml := bufHtml.String()
		slog.DebugContext(ctx, "email content",
			slog.String("plain", messagePlain),
			slog.String("html", messageHtml),
		)

		err = mailer.Send("", []string{user.Email},
			vars["Subject"].(string),
			messagePlain, messageHtml,
			mail.MailHeader{
				"X-PM-Message-Stream": "broadcast",
			},
		)
		if err != nil {
			return fmt.Errorf("error sending email: %w", err)
		}
		_, err = model.NewUserEventNotification(ctx, s.Db, user.ID, event.ID)
		if err != nil {
			return fmt.Errorf("error updating database: %w", err)
		}
	}
	return nil
}
