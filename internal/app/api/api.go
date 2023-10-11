package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/middleware/debug"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/xhandler"
	"github.com/dropwhile/icbt/internal/session"
	"github.com/dropwhile/icbt/internal/util"
	"github.com/dropwhile/icbt/resources"
)

type API struct {
	*chi.Mux
	handler *xhandler.XHandler
	closers []func()
}

func (api *API) Close() {
	for _, f := range api.closers {
		f()
	}
}

func (api *API) OnClose(f func()) {
	api.closers = append(api.closers, f)
}

func New(
	db *pgxpool.Pool,
	tpl resources.TemplateMap,
	mailer *util.Mailer,
	hmacKey, csrfKey []byte,
	isProd bool,
) *API {
	zh := &xhandler.XHandler{
		Db:      model.SetupFromDbPool(db),
		Tpl:     tpl,
		SessMgr: session.NewDBSessionManager(db),
		Mailer:  mailer,
		Hmac:    util.NewHmac(hmacKey),
	}

	api := &API{Mux: chi.NewRouter(), handler: zh}
	api.OnClose(zh.SessMgr.Close)

	// Router/Middleware //
	r := api.Mux
	r.NotFound(zh.NotFound)
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
			csrfKey,
			// false in development only!
			csrf.Secure(isProd),
			// setup path so csrf works _between_ pages (eg. htmx calls)
			csrf.Path("/"),
			// Must be in CORS Allowed and Exposed Headers
			csrf.RequestHeader("X-CSRF-Token"),
		))
		r.Use(auth.Load(db, zh.SessMgr))

		// Routing //
		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(auth.Require)
			r.Use(middleware.NoCache)
			// acccount/settings
			r.Get("/settings", zh.ShowSettings)
			r.Post("/settings", zh.UpdateSettings)
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
			r.Post("/verify", zh.SendVerificationEmail)
			r.Get("/verify/{uvRefID:[0-9a-z]+}-{hmac:[0-9a-z]+}", zh.VerifyEmail)
		})

		// Public routes
		r.Group(func(r chi.Router) {
			r.Get("/", zh.ShowIndex)
			// login
			r.Post("/login", zh.Login)
			r.Get("/login", zh.ShowLoginForm)
			// forgot password
			r.Get("/forgot-password", zh.ShowForgotPasswordForm)
			r.Post("/forgot-password", zh.SendResetPasswordEmail)
			r.Get("/forgot-password/{upwRefID:[0-9a-z]+}-{hmac:[0-9a-z]+}", zh.ShowPasswordResetForm)
			r.Post("/forgot-password/{upwRefID:[0-9a-z]+}-{hmac:[0-9a-z]+}", zh.ResetPassword)
			// account creation
			r.Get("/create-account", zh.ShowCreateAccount)
			r.Post("/create-account", zh.CreateAccount)
			// local only debug stuff
			if !isProd {
				r.Route("/debug", func(r chi.Router) {
				})
			}
		})
	})

	return api
}
