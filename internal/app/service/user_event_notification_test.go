package service

import (
	"context"
	"html/template"
	"testing"
	"time"

	"github.com/dropwhile/refid/v2"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"go.uber.org/mock/gomock"
	"gotest.tools/v3/assert"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/resources"
	"github.com/dropwhile/icbt/internal/mail"
	"github.com/dropwhile/icbt/internal/util"
)

func TestService_NotifyUsersPendingEvents(t *testing.T) {
	t.Parallel()

	user := &model.User{
		ID:           1,
		RefID:        refid.Must(model.NewUserRefID()),
		Email:        "user@example.com",
		Name:         "user",
		PWHash:       []byte("00x00"),
		Verified:     true,
		Created:      tstTs,
		LastModified: tstTs,
		Settings: model.UserSettings{
			ReminderThresholdHours: 0,
			EnableReminders:        true,
		},
	}
	event := &model.Event{
		ID:            2,
		RefID:         refid.Must(model.NewEventRefID()),
		UserID:        user.ID,
		Name:          "event",
		Description:   "description",
		Archived:      false,
		ItemSortOrder: []int{1, 2, 3},
		StartTime:     tstTs,
		StartTimeTz:   util.Must(ParseTimeZone("Etc/UTC")),
	}
	eventItem := &model.EventItem{
		ID:          3,
		RefID:       refid.Must(model.NewEventItemRefID()),
		EventID:     event.ID,
		Description: "eventitem",
	}
	earmark := &model.Earmark{
		ID:          4,
		RefID:       refid.Must(model.NewEarmarkRefID()),
		EventItemID: eventItem.ID,
		UserID:      user.ID,
		Note:        "earmark",
	}

	t.Run("notify pending should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})
		mailer := SetupMailerMock(t)
		templates := &resources.TemplateMap{
			"mail_reminder.gohtml": template.Must(
				template.New("mail_reminder.gohtml").
					ParseFiles("../resources/templates/html/view/mail_reminder.gohtml"),
			),
		}

		now := time.Now()

		mock.ExpectQuery("WITH subt").
			WithArgs().
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"user_id", "event_id", "when", "owner", "items",
				}).
				AddRow(user.ID, event.ID, now, true, []int{eventItem.ID}),
			)
		mock.ExpectQuery("^SELECT (.+) FROM user_").
			WithArgs(user.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "email", "name", "verified", "settings",
				}).
				AddRow(
					user.ID, user.RefID, user.Email, user.Name,
					user.Verified, user.Settings,
				),
			)
		mock.ExpectQuery("^SELECT (.+) FROM event_").
			WithArgs(event.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "name", "description", "archived",
					"item_sort_order", "start_time", "start_time_tz",
				}).
				AddRow(
					event.ID, event.RefID, event.UserID, event.Name,
					event.Description, event.Archived, event.ItemSortOrder,
					event.StartTime, event.StartTimeTz,
				),
			)
		mock.ExpectQuery("^SELECT (.+) FROM event_item_ (.+)").
			WithArgs([]int{eventItem.ID}).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "event_id", "description",
				}).
				AddRow(
					eventItem.ID, eventItem.RefID, eventItem.EventID,
					eventItem.Description,
				))
		mock.ExpectQuery("^SELECT (.+) FROM earmark_ (.+)").
			WithArgs([]int{eventItem.ID}).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "user_id", "event_item_id", "note",
				}).
				AddRow(
					earmark.ID, earmark.RefID, earmark.UserID,
					earmark.EventItemID, earmark.Note,
				))
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO user_event_notification_").
			WithArgs(pgx.NamedArgs{
				"userID":  user.ID,
				"eventID": event.ID,
			}).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"user_id", "event_id",
				}).
				AddRow(user.ID, event.ID),
			)
		mock.ExpectCommit()
		mock.ExpectRollback()

		mailer.EXPECT().
			Send("", []string{user.Email},
				"Upcoming Event Reminder",
				gomock.AssignableToTypeOf("string"),
				gomock.AssignableToTypeOf("string"),
				mail.MailHeader{
					"X-PM-Message-Stream": "broadcast",
				},
			).
			Return(nil)

		err := svc.NotifyUsersPendingEvents(
			ctx, mailer, templates, "http://example.org",
		)
		assert.NilError(t, err)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("notify pending with user reminders disabled should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})
		mailer := SetupMailerMock(t)
		templates := &resources.TemplateMap{
			"mail_reminder.gohtml": template.Must(
				template.New("mail_reminder.gohtml").
					ParseFiles("../resources/templates/html/view/mail_reminder.gohtml"),
			),
		}

		now := time.Now()

		mock.ExpectQuery("WITH subt").
			WithArgs().
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"user_id", "event_id", "when", "owner", "items",
				}).
				AddRow(user.ID, event.ID, now, true, []int{eventItem.ID}),
			)
		mock.ExpectQuery("^SELECT (.+) FROM user_").
			WithArgs(user.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "email", "name", "verified", "settings",
				}).
				AddRow(
					user.ID, user.RefID, user.Email, user.Name,
					user.Verified, model.UserSettings{
						ReminderThresholdHours: 0,
						EnableReminders:        false,
					},
				),
			)

		err := svc.NotifyUsersPendingEvents(
			ctx, mailer, templates, "http://example.org",
		)
		assert.NilError(t, err)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("notify pending with time threshold not met should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})
		mailer := SetupMailerMock(t)
		templates := &resources.TemplateMap{
			"mail_reminder.gohtml": template.Must(
				template.New("mail_reminder.gohtml").
					ParseFiles("../resources/templates/html/view/mail_reminder.gohtml"),
			),
		}

		when := time.Now().Add(time.Duration(25) * time.Hour)

		mock.ExpectQuery("WITH subt").
			WithArgs().
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"user_id", "event_id", "when", "owner", "items",
				}).
				AddRow(user.ID, event.ID, when, true, []int{eventItem.ID}),
			)
		mock.ExpectQuery("^SELECT (.+) FROM user_").
			WithArgs(user.ID).
			WillReturnRows(pgxmock.NewRows(
				[]string{
					"id", "ref_id", "email", "name", "verified", "settings",
				}).
				AddRow(
					user.ID, user.RefID, user.Email, user.Name,
					user.Verified, user.Settings,
				),
			)

		err := svc.NotifyUsersPendingEvents(
			ctx, mailer, templates, "http://example.org",
		)
		assert.NilError(t, err)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})

	t.Run("notify pending with empty results should succeed", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		mock := SetupDBMock(t, ctx)
		svc := New(Options{Db: mock})
		mailer := SetupMailerMock(t)
		templates := &resources.TemplateMap{
			"mail_reminder.gohtml": template.Must(
				template.New("mail_reminder.gohtml").
					ParseFiles("../resources/templates/html/view/mail_reminder.gohtml"),
			),
		}

		mock.ExpectQuery("WITH subt").
			WithArgs().
			WillReturnRows(pgxmock.NewRows([]string{
				"user_id", "event_id", "when", "owner", "items",
			}))

		err := svc.NotifyUsersPendingEvents(
			ctx, mailer, templates, "http://example.org",
		)
		assert.NilError(t, err)
		// we make sure that all expectations were met
		assert.Assert(t, mock.ExpectationsWereMet(),
			"there were unfulfilled expectations")
	})
}
