// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package app

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/dropwhile/icanbringthat/internal/app/handler"
	"github.com/dropwhile/icanbringthat/internal/app/resources"
	"github.com/dropwhile/icanbringthat/internal/app/rpc"
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/mail"
	"github.com/dropwhile/icanbringthat/internal/middleware/auth"
	"github.com/dropwhile/icanbringthat/internal/middleware/debug"
	"github.com/dropwhile/icanbringthat/internal/middleware/header"
	"github.com/dropwhile/icanbringthat/internal/session"
)

const TwirpPrefix = "/api"

type App struct {
	*chi.Mux
	handler *handler.Handler
	closers []func()
}

func (app *App) Close() {
	for _, f := range app.closers {
		f()
	}
}

func (app *App) OnClose(f func()) {
	app.closers = append(app.closers, f)
}

func New(
	db *pgxpool.Pool,
	rdb *redis.Client,
	templates resources.TGetter,
	mailer mail.MailSender,
	conf *Config,
) (*App, error) {
	service := service.New(service.Options{Db: db})
	baseURL := strings.TrimSuffix(conf.BaseURL, "/")
	isProd := conf.Production
	sessMgr := session.NewRedisSessionManager(rdb, conf.Production)

	zh, err := handler.New(
		handler.Options{
			Db:           db,
			Redis:        rdb,
			Templates:    templates,
			SessMgr:      sessMgr,
			Mailer:       mailer,
			HMACKeyBytes: conf.HMACKeyBytes,
			BaseURL:      baseURL,
			IsProd:       isProd,
		},
	)
	if err != nil {
		return nil, err
	}

	app := &App{Mux: chi.NewRouter(), handler: zh}
	app.OnClose(sessMgr.Close)

	// Router/Middleware //
	r := app.Mux
	r.Use(middleware.Logger)
	r.NotFound(zh.NotFoundHandler)
	r.Use(middleware.RealIP)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.GetHead)
	r.Use(header.RequestID)
	if conf.RequestLogging {
		r.Use(debug.RequestLogger())
	}
	r.Use(middleware.Recoverer)

	// nest so session and csrf middlewares are not used by
	// any static routes added onto the handler later
	r.Group(func(r chi.Router) {
		r.Use(sessMgr.LoadAndSave)
		r.Use(csrf.Protect(
			conf.CSRFKeyBytes,
			// false in development only!
			csrf.Secure(conf.Production),
			// setup path so csrf works _between_ pages (eg. htmx calls)
			csrf.Path("/"),
			// Must be in CORS Allowed and Exposed Headers
			csrf.RequestHeader("X-CSRF-Token"),
		))
		r.Use(auth.Load(service, sessMgr))

		// Routing //
		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.NoCache)
			r.Use(auth.Require)
			r.Get("/about", zh.AboutShow)
			// acccount/settings
			r.Get("/settings", zh.SettingsShow)
			r.Post("/settings", zh.SettingsUpdate)
			r.Post("/settings/auth", zh.SettingsAuthUpdate)
			r.Post("/settings/auth/api", zh.SettingsAuthApiUpdate)
			r.Post("/settings/reminders", zh.SettingsRemindersUpdate)
			r.Delete("/settings", zh.AccountDelete)
			// logout
			r.Post("/logout", zh.Logout)
			// dashboard
			r.Get("/dashboard", zh.DashboardShow)
			// event
			r.Get("/events", zh.EventsList)
			r.Post("/events", zh.EventCreate)
			r.Get("/events/add", zh.CreateEventShowAddForm)
			r.Get("/events/{eRefID:[0-9a-z]+}", zh.EventShow)
			r.Post("/events/{eRefID:[0-9a-z]+}", zh.EventUpdate)
			r.Delete("/events/{eRefID:[0-9a-z]+}", zh.EventDelete)
			r.Get("/events/{eRefID:[0-9a-z]+}/edit", zh.EventShowEditForm)
			// favorites
			r.Get("/favorites", zh.FavoritesList)
			r.Put("/events/{eRefID:[0-9a-z]+}/favorite", zh.FavoriteAdd)
			r.Delete("/events/{eRefID:[0-9a-z]+}/favorite", zh.FavoriteDelete)
			// event item
			r.Post("/events/{eRefID:[0-9a-z]+}/items", zh.EventItemCreate)
			r.Get("/events/{eRefID:[0-9a-z]+}/items/add", zh.EventItemShowCreateForm)
			r.Post("/events/{eRefID:[0-9a-z]+}/sort", zh.EventItemSortingUpdate)
			r.Post("/events/{eRefID:[0-9a-z]+}/items/{iRefID:[0-9a-z]+}", zh.EventItemUpdate)
			r.Delete("/events/{eRefID:[0-9a-z]+}/items/{iRefID:[0-9a-z]+}", zh.EventItemDelete)
			r.Get("/events/{eRefID:[0-9a-z]+}/items/{iRefID:[0-9a-z]+}/edit", zh.EventItemShowEditForm)
			// earmarks
			r.Post("/events/{eRefID:[0-9a-z]+}/items/{iRefID:[0-9a-z]+}/earmarks", zh.EarmarkCreate)
			r.Get("/events/{eRefID:[0-9a-z]+}/items/{iRefID:[0-9a-z]+}/earmarks/add", zh.CreateEarmarkShowCreateForm)
			r.Get("/earmarks", zh.EarmarksList)
			r.Delete("/earmarks/{mRefID:[0-9a-z]+}", zh.EarmarkDelete)
			// r.Get("/earmarks/{mRefID:[0-9a-z]+}", zh.EarmarkShow)
			// r.Post("/earmarks/{mRefID:[0-9a-z]+}", zh.EarmarkUpdate)
			// r.Get("/profile/{uRefID:[a-zA-Z-]+}", zh.ProfileShow)
			// notifications
			r.Get("/notifications", zh.NotificationsList)
			r.Delete("/notifications", zh.NotificationsDeleteAll)
			r.Delete("/notifications/{nRefID:[0-9a-z]+}", zh.NotificationDelete)
			// account verification
			r.Post("/verify", zh.VerifySendEmail)
			r.Get("/verify/{uvRefID:[0-9a-z]+}-{hmac:[0-9a-z]+}", zh.VerifyEmail)
			// webauthn
			r.Get("/webauthn/register", zh.WebAuthnBeginRegistration)
			r.Post("/webauthn/register", zh.WebAuthnFinishRegistration)
			r.Delete("/webauthn/key/{cRefID:[0-9a-z]+}", zh.WebAuthnDeleteKey)
		})

		// Public routes
		r.Group(func(r chi.Router) {
			r.Get("/", zh.IndexShow)
			// login
			r.Post("/login", zh.Login)
			r.Get("/login", zh.LoginShowForm)
			r.Get("/webauthn/login", zh.WebAuthnBeginLogin)
			r.Post("/webauthn/login", zh.WebAuthnFinishLogin)
			// forgot password
			r.Get("/forgot-password", zh.ForgotPasswordShowForm)
			r.Post("/forgot-password", zh.ResetPasswordSendEmail)
			r.Get("/forgot-password/{upwRefID:[0-9a-z]+}-{hmac:[0-9a-z]+}", zh.PasswordResetShowForm)
			r.Post("/forgot-password/{upwRefID:[0-9a-z]+}-{hmac:[0-9a-z]+}", zh.PasswordReset)
			// account creation
			r.Get("/create-account", zh.AccountShowCreate)
			r.Post("/create-account", zh.AccountCreate)
			// local only debug stuff
			if !conf.Production {
				r.Route("/debug", func(r chi.Router) {
				})
			}
		})
	})

	// webhooks
	r.Group(func(r chi.Router) {
		r.Use(middleware.NoCache)
		r.Use(middleware.BasicAuth("simple", conf.WebhookCreds))
		r.Post("/webhooks/pm", zh.PostmarkCallback)
	})

	if conf.RpcApi {
		// rpc api
		rpcServer, err := rpc.New(
			rpc.Options{
				Db:           db,
				Redis:        rdb,
				Templates:    templates,
				Mailer:       mailer,
				HMACKeyBytes: conf.HMACKeyBytes,
				BaseURL:      baseURL,
				IsProd:       isProd,
			},
		)
		if err != nil {
			return nil, err
		}

		r.Route(TwirpPrefix, func(r chi.Router) {
			// add auth token middleware here instead,
			// which pulls an auth token from a header,
			// looks it up in the db, and sets the user in the context
			r.NotFound(http.NotFound)
			r.Use(middleware.NoCache)
			r.Use(auth.LoadAuthToken)
			r.Mount("/", rpcServer.GenHandler(TwirpPrefix))
		})
	}

	return app, nil
}
