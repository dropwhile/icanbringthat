package app

import (
	ah "github.com/dropwhile/icbt/internal/app/handler"
	"github.com/dropwhile/icbt/internal/app/middleware/auth"
	"github.com/dropwhile/icbt/internal/app/middleware/debug"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/session"
	"github.com/dropwhile/icbt/internal/util"
	res "github.com/dropwhile/icbt/resources"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/rs/zerolog/log"
)

type API struct {
	*chi.Mux
	handler *ah.Handler
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

func NewAPI(db *model.DB, tpl res.TemplateMap, mailer *util.Mailer, csrfKey, hmacKey []byte, isProd bool) *API {
	ah := &ah.Handler{
		Db:      db,
		Tpl:     tpl,
		SessMgr: session.NewDBSessionManager(db.GetPool()),
		Mailer:  mailer,
		Hmac:    util.NewHmac(hmacKey),
	}

	api := &API{
		Mux:     chi.NewRouter(),
		handler: ah,
	}
	api.OnClose(ah.SessMgr.Close)

	// Router/Middleware //
	r := api.Mux
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.GetHead)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	if log.Debug().Enabled() {
		r.Use(debug.RequestLogger())
	}
	r.Use(middleware.Recoverer)
	r.Use(ah.SessMgr.LoadAndSave)
	r.Use(csrf.Protect(
		csrfKey,
		// false in development only!
		csrf.Secure(isProd),
		// setup path so csrf works _between_ pages (eg. htmx calls)
		csrf.Path("/"),
		// Must be in CORS Allowed and Exposed Headers
		csrf.RequestHeader("X-CSRF-Token"),
	))
	r.Use(auth.Load(db, ah.SessMgr))

	// Routing //
	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(auth.Require)
		// acccount/settings
		r.Get("/settings", ah.ShowSettings)
		r.Post("/settings", ah.UpdateSettings)
		r.Delete("/settings", ah.DeleteAccount)
		// logout
		r.Post("/logout", ah.Logout)
		// dashboard
		r.Get("/dashboard", ah.ShowDashboard)
		// event
		r.Get("/events", ah.ListEvents)
		r.Post("/events", ah.CreateEvent)
		r.Get("/events/add", ah.ShowCreateEventForm)
		r.Get("/events/{eRefId:[0-9a-z]+}", ah.ShowEvent)
		r.Post("/events/{eRefId:[0-9a-z]+}", ah.UpdateEvent)
		r.Delete("/events/{eRefId:[0-9a-z]+}", ah.DeleteEvent)
		r.Get("/events/{eRefId:[0-9a-z]+}/edit", ah.ShowEditEventForm)
		// event item
		r.Post("/events/{eRefId:[0-9a-z]+}/items", ah.CreateEventItem)
		r.Get("/events/{eRefId:[0-9a-z]+}/items/add", ah.ShowCreateEventItemForm)
		r.Post("/events/{eRefId:[0-9a-z]+}/items/{iRefId:[0-9a-z]+}", ah.UpdateEventItem)
		r.Delete("/events/{eRefId:[0-9a-z]+}/items/{iRefId:[0-9a-z]+}", ah.DeleteEventItem)
		r.Get("/events/{eRefId:[0-9a-z]+}/items/{iRefId:[0-9a-z]+}/edit", ah.ShowEventItemEditForm)
		// earmarks
		r.Post("/events/{eRefId:[0-9a-z]+}/items/{iRefId:[0-9a-z]+}/earmarks", ah.CreateEarmark)
		r.Get("/events/{eRefId:[0-9a-z]+}/items/{iRefId:[0-9a-z]+}/earmarks/add", ah.ShowCreateEarmarkForm)
		r.Get("/earmarks", ah.ListEarmarks)
		r.Delete("/earmarks/{mRefId:[0-9a-z]+}", ah.DeleteEarmark)
		/*
			r.Get("/earmarks/{mRefId:[0-9a-z]+}", ah.ShowEarmark)
			r.Post("/earmarks/{mRefId:[0-9a-z]+}", ah.UpdateEarmark)
		*/
		// r.Get("/profile/{uRefId:[a-zA-Z-]+}", ah.ShowProfile)
	})

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/", ah.ShowIndex)
		// login
		r.Post("/login", ah.Login)
		r.Get("/login", ah.ShowLoginForm)
		// forgot password
		r.Get("/forgot-password", ah.ShowForgotPasswordForm)
		r.Post("/forgot-password", ah.SendResetPasswordEmail)
		r.Get("/forgot-password/{upwRefId:[0-9a-z]+}-{hmac:[0-9a-z]+}", ah.ShowPasswordResetForm)
		r.Post("/forgot-password/{upwRefId:[0-9a-z]+}-{hmac:[0-9a-z]+}", ah.ResetPassword)
		// account creation
		r.Get("/create-account", ah.ShowCreateAccount)
		r.Post("/create-account", ah.CreateAccount)
		// local only debug stuff
		if !isProd {
			r.Route("/debug", func(r chi.Router) {
				r.Get("/templates", ah.TestTemplates)
			})
		}
	})

	r.NotFound(ah.NotFound)

	return api
}
