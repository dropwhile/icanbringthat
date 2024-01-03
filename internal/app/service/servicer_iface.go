// Code generated by ifacemaker; DO NOT EDIT.

package service

import (
	"context"
	"time"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/resources"
	"github.com/dropwhile/icbt/internal/errs"
	"github.com/dropwhile/icbt/internal/mail"
)

// Servicer ...
type Servicer interface {
	GetEarmarksByEventID(ctx context.Context, eventID int) ([]*model.Earmark, errs.Error)
	GetEarmarkByEventItemID(ctx context.Context, eventItemID int) (*model.Earmark, errs.Error)
	GetEarmarksCount(ctx context.Context, userID int) (*model.BifurcatedRowCounts, errs.Error)
	GetEarmarksPaginated(ctx context.Context, userID int, limit, offset int, archived bool) ([]*model.Earmark, *Pagination, errs.Error)
	GetEarmarks(ctx context.Context, userID int, archived bool) ([]*model.Earmark, errs.Error)
	NewEarmark(ctx context.Context, user *model.User, eventItemID int, note string) (*model.Earmark, errs.Error)
	GetEarmark(ctx context.Context, refID model.EarmarkRefID) (*model.Earmark, errs.Error)
	DeleteEarmark(ctx context.Context, userID int, earmark *model.Earmark) errs.Error
	DeleteEarmarkByRefID(ctx context.Context, userID int, refID model.EarmarkRefID) errs.Error
	GetEvent(ctx context.Context, refID model.EventRefID) (*model.Event, errs.Error)
	GetEventByID(ctx context.Context, ID int) (*model.Event, errs.Error)
	GetEventsByIDs(ctx context.Context, eventIDs []int) ([]*model.Event, errs.Error)
	DeleteEvent(ctx context.Context, userID int, refID model.EventRefID) errs.Error
	UpdateEvent(ctx context.Context, userID int, refID model.EventRefID, euvs *EventUpdateValues) errs.Error
	UpdateEventItemSorting(ctx context.Context, userID int, refID model.EventRefID, itemSortOrder []int) (*model.Event, errs.Error)
	CreateEvent(ctx context.Context, user *model.User, name string, description string, when time.Time, tz string) (*model.Event, errs.Error)
	GetEventsPaginated(ctx context.Context, userID int, limit, offset int, archived bool) ([]*model.Event, *Pagination, errs.Error)
	GetEventsComingSoonPaginated(ctx context.Context, userID int, limit, offset int) ([]*model.Event, *Pagination, errs.Error)
	GetEventsCount(ctx context.Context, userID int) (*model.BifurcatedRowCounts, errs.Error)
	GetEvents(ctx context.Context, userID int, archived bool) ([]*model.Event, errs.Error)
	ArchiveOldEvents(ctx context.Context) error
	GetEventItemsCount(ctx context.Context, eventIDs []int) ([]*model.EventItemCount, errs.Error)
	GetEventItemsByEvent(ctx context.Context, refID model.EventRefID) ([]*model.EventItem, errs.Error)
	GetEventItemsByEventID(ctx context.Context, eventID int) ([]*model.EventItem, errs.Error)
	GetEventItemsByIDs(ctx context.Context, eventItemIDs []int) ([]*model.EventItem, errs.Error)
	GetEventItem(ctx context.Context, eventItemRefID model.EventItemRefID) (*model.EventItem, errs.Error)
	GetEventItemByID(ctx context.Context, eventItemID int) (*model.EventItem, errs.Error)
	RemoveEventItem(ctx context.Context, userID int, eventItemRefID model.EventItemRefID, failIfChecks func(*model.EventItem) bool) errs.Error
	AddEventItem(ctx context.Context, userID int, refID model.EventRefID, description string) (*model.EventItem, errs.Error)
	UpdateEventItem(ctx context.Context, userID int, refID model.EventItemRefID, description string, failIfChecks func(*model.EventItem) bool) (*model.EventItem, errs.Error)
	AddFavorite(ctx context.Context, userID int, refID model.EventRefID) (*model.Event, errs.Error)
	RemoveFavorite(ctx context.Context, userID int, refID model.EventRefID) errs.Error
	GetFavoriteEventsPaginated(ctx context.Context, userID int, limit, offset int, archived bool) ([]*model.Event, *Pagination, errs.Error)
	GetFavoriteEventsCount(ctx context.Context, userID int) (*model.BifurcatedRowCounts, errs.Error)
	GetFavoriteEvents(ctx context.Context, userID int, archived bool) ([]*model.Event, errs.Error)
	GetFavoriteByUserEvent(ctx context.Context, userID int, eventID int) (*model.Favorite, errs.Error)
	GetNotificationsCount(ctx context.Context, userID int) (int, errs.Error)
	GetNotificationsPaginated(ctx context.Context, userID int, limit, offset int) ([]*model.Notification, *Pagination, errs.Error)
	GetNotifications(ctx context.Context, userID int) ([]*model.Notification, errs.Error)
	DeleteNotification(ctx context.Context, userID int, refID model.NotificationRefID) errs.Error
	DeleteAllNotifications(ctx context.Context, userID int) errs.Error
	NewNotification(ctx context.Context, userID int, message string) (*model.Notification, errs.Error)
	GetUser(ctx context.Context, refID model.UserRefID) (*model.User, errs.Error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, errs.Error)
	GetUserByID(ctx context.Context, ID int) (*model.User, errs.Error)
	GetUsersByIDs(ctx context.Context, userIDs []int) ([]*model.User, errs.Error)
	NewUser(ctx context.Context, email, name string, rawPass []byte) (*model.User, errs.Error)
	UpdateUser(ctx context.Context, userID int, euvs *UserUpdateValues) errs.Error
	UpdateUserSettings(ctx context.Context, userID int, pm *model.UserSettings) errs.Error
	DeleteUser(ctx context.Context, userID int) errs.Error
	GetApiKeyByUser(ctx context.Context, userID int) (*model.ApiKey, errs.Error)
	GetUserByApiKey(ctx context.Context, token string) (*model.User, errs.Error)
	NewApiKey(ctx context.Context, userID int) (*model.ApiKey, errs.Error)
	NewApiKeyIfNotExists(ctx context.Context, userID int) (*model.ApiKey, errs.Error)
	NotifyUsersPendingEvents(ctx context.Context, mailer mail.MailSender, tplContainer resources.TGetter, siteBaseUrl string) error
	GetUserPWResetByRefID(ctx context.Context, refID model.UserPWResetRefID) (*model.UserPWReset, errs.Error)
	NewUserPWReset(ctx context.Context, userID int) (*model.UserPWReset, errs.Error)
	UpdateUserPWReset(ctx context.Context, user *model.User, upw *model.UserPWReset) errs.Error
	GetUserVerifyByRefID(ctx context.Context, refID model.UserVerifyRefID) (*model.UserVerify, errs.Error)
	NewUserVerify(ctx context.Context, userID int) (*model.UserVerify, errs.Error)
	SetUserVerified(ctx context.Context, user *model.User, verifier *model.UserVerify) errs.Error
	GetUserCredentialByRefID(ctx context.Context, refID model.CredentialRefID) (*model.UserCredential, errs.Error)
	GetUserCredentialsByUser(ctx context.Context, userID int) ([]*model.UserCredential, errs.Error)
	GetUserCredentialCountByUser(ctx context.Context, userID int) (int, errs.Error)
	DeleteUserCredential(ctx context.Context, credentialID int) errs.Error
	NewUserCredential(ctx context.Context, userID int, keyName string, credential []byte) (*model.UserCredential, errs.Error)
	WebAuthnUserFrom(user *model.User) *WebAuthnUser
	DisableRemindersWithNotification(ctx context.Context, email string, suppressionReason string) errs.Error
}
