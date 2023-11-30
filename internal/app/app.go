package app

import (
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icbt/internal/app/handler"
	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/middleware/debug"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/rpc"
	"github.com/dropwhile/icbt/internal/crypto"
	"github.com/dropwhile/icbt/internal/mail"
	"github.com/dropwhile/icbt/internal/session"
	"github.com/dropwhile/icbt/resources"
	rpcdef "github.com/dropwhile/icbt/rpc"
)

type App struct {
	*chi.Mux
	handler *handler.Handler
	closers []func()
}

func (api *App) Close() {
	for _, f := range api.closers {
		f()
	}
}

func (api *App) OnClose(f func()) {
	api.closers = append(api.closers, f)
}

func New(
	db *pgxpool.Pool,
	rdb *redis.Client,
	templateMap resources.TemplateMap,
	mailer *mail.Mailer,
	conf *Config,
) *App {
	zh := &handler.Handler{
		Db:          model.SetupFromDbPool(db),
		Redis:       rdb,
		TemplateMap: templateMap,
		SessMgr:     session.NewRedisSessionManager(rdb, conf.Production),
		Mailer:      mailer,
		MAC:         crypto.NewMAC(conf.HMACKeyBytes),
		BaseURL:     strings.TrimSuffix(conf.BaseURL, "/"),
		IsProd:      conf.Production,
	}

	api := &App{Mux: chi.NewRouter(), handler: zh}
	api.OnClose(zh.SessMgr.Close)

	// Router/Middleware //
	r := api.Mux
	r.NotFound(zh.NotFoundHandler)
	r.Use(middleware.RealIP)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.GetHead)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	if log.Trace().Enabled() {
		r.Use(debug.RequestLogger())
	}
	r.Use(middleware.Recoverer)

	// nest so session and csrf middlewares are not used by
	// any static routes added onto the handler later
	r.Group(func(r chi.Router) {
		r.Use(zh.SessMgr.LoadAndSave)
		r.Use(csrf.Protect(
			conf.CSRFKeyBytes,
			// false in development only!
			csrf.Secure(conf.Production),
			// setup path so csrf works _between_ pages (eg. htmx calls)
			csrf.Path("/"),
			// Must be in CORS Allowed and Exposed Headers
			csrf.RequestHeader("X-CSRF-Token"),
		))
		r.Use(auth.Load(db, zh.SessMgr))

		// Routing //
		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.NoCache)
			r.Use(auth.Require)
			r.Get("/about", zh.ShowAbout)
			// acccount/settings
			r.Get("/settings", zh.ShowSettings)
			r.Post("/settings", zh.UpdateSettings)
			r.Post("/settings/auth", zh.UpdateAuthSettings)
			r.Post("/settings/reminders", zh.UpdateRemindersSettings)
			r.Delete("/settings", zh.DeleteAccount)
			// logout
			r.Post("/logout", zh.Logout)
			// dashboard
			r.Get("/dashboard", zh.ShowDashboard)
			// event
			r.Get("/events", zh.ListEvents)
			r.Post("/events", zh.CreateEvent)
			r.Get("/events/add", zh.ShowCreateEventForm)
			r.Get("/events/{eRefID:[0-9a-z]+}", zh.ShowEvent)
			r.Post("/events/{eRefID:[0-9a-z]+}", zh.UpdateEvent)
			r.Delete("/events/{eRefID:[0-9a-z]+}", zh.DeleteEvent)
			r.Get("/events/{eRefID:[0-9a-z]+}/edit", zh.ShowEditEventForm)
			// favorites
			r.Get("/favorites", zh.ListFavorites)
			r.Put("/events/{eRefID:[0-9a-z]+}/favorite", zh.AddFavorite)
			r.Delete("/events/{eRefID:[0-9a-z]+}/favorite", zh.DeleteFavorite)
			// event item
			r.Post("/events/{eRefID:[0-9a-z]+}/items", zh.CreateEventItem)
			r.Get("/events/{eRefID:[0-9a-z]+}/items/add", zh.ShowCreateEventItemForm)
			r.Post("/events/{eRefID:[0-9a-z]+}/sort", zh.UpdateEventItemSorting)
			r.Post("/events/{eRefID:[0-9a-z]+}/items/{iRefID:[0-9a-z]+}", zh.UpdateEventItem)
			r.Delete("/events/{eRefID:[0-9a-z]+}/items/{iRefID:[0-9a-z]+}", zh.DeleteEventItem)
			r.Get("/events/{eRefID:[0-9a-z]+}/items/{iRefID:[0-9a-z]+}/edit", zh.ShowEventItemEditForm)
			// earmarks
			r.Post("/events/{eRefID:[0-9a-z]+}/items/{iRefID:[0-9a-z]+}/earmarks", zh.CreateEarmark)
			r.Get("/events/{eRefID:[0-9a-z]+}/items/{iRefID:[0-9a-z]+}/earmarks/add", zh.ShowCreateEarmarkForm)
			r.Get("/earmarks", zh.ListEarmarks)
			r.Delete("/earmarks/{mRefID:[0-9a-z]+}", zh.DeleteEarmark)
			// r.Get("/earmarks/{mRefID:[0-9a-z]+}", zh.ShowEarmark)
			// r.Post("/earmarks/{mRefID:[0-9a-z]+}", zh.UpdateEarmark)
			// r.Get("/profile/{uRefID:[a-zA-Z-]+}", zh.ShowProfile)
			// notifications
			r.Get("/notifications", zh.ListNotifications)
			r.Delete("/notifications", zh.DeleteAllNotifications)
			r.Delete("/notifications/{nRefID:[0-9a-z]+}", zh.DeleteNotification)
			// account verification
			r.Post("/verify", zh.SendVerificationEmail)
			r.Get("/verify/{uvRefID:[0-9a-z]+}-{hmac:[0-9a-z]+}", zh.VerifyEmail)
			// webauthn
			r.Get("/webauthn/register", zh.WebAuthnBeginRegistration)
			r.Post("/webauthn/register", zh.WebAuthnFinishRegistration)
			r.Delete("/webauthn/key/{cRefID:[0-9a-z]+}", zh.DeleteWebAuthnKey)
		})

		// Public routes
		r.Group(func(r chi.Router) {
			r.Get("/", zh.ShowIndex)
			// login
			r.Post("/login", zh.Login)
			r.Get("/login", zh.ShowLoginForm)
			r.Get("/webauthn/login", zh.WebAuthnBeginLogin)
			r.Post("/webauthn/login", zh.WebAuthnFinishLogin)
			// forgot password
			r.Get("/forgot-password", zh.ShowForgotPasswordForm)
			r.Post("/forgot-password", zh.SendResetPasswordEmail)
			r.Get("/forgot-password/{upwRefID:[0-9a-z]+}-{hmac:[0-9a-z]+}", zh.ShowPasswordResetForm)
			r.Post("/forgot-password/{upwRefID:[0-9a-z]+}-{hmac:[0-9a-z]+}", zh.ResetPassword)
			// account creation
			r.Get("/create-account", zh.ShowCreateAccount)
			r.Post("/create-account", zh.CreateAccount)
			// local only debug stuff
			if !conf.Production {
				r.Route("/debug", func(r chi.Router) {
				})
			}
		})
	})

	// webhooks
	r.Group(func(r chi.Router) {
		r.Use(middleware.BasicAuth("simple", conf.WebhookCreds))
		r.Use(middleware.NoCache)
		r.Post("/webhooks/pm", zh.PostmarkCallback)
	})

	// rpc api
	rpcServer := &rpc.Server{
		Db:          zh.Db,
		Redis:       zh.Redis,
		TemplateMap: zh.TemplateMap,
		Mailer:      zh.Mailer,
		MAC:         zh.MAC,
		BaseURL:     zh.BaseURL,
		IsProd:      zh.IsProd,
	}
	twirpHooks := &twirp.ServerHooks{
		RequestReceived: rpc.AuthHook(rpcServer.Db),
	}
	twirpHandler := rpcdef.NewRpcServer(
		rpcServer,
		twirp.WithServerPathPrefix("/api"),
		twirp.WithServerHooks(twirpHooks),
	)
	r.Group(func(r chi.Router) {
		// add auth token middleware here instead,
		// which pulls an auth token from a header,
		// looks it up in the db, and sets the user in the context
		r.Use(middleware.NoCache)
		r.Use(auth.LoadToken)
		r.Mount(twirpHandler.PathPrefix(), twirpHandler)
	})

	return api
}
