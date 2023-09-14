package app

import (
	ah "github.com/dropwhile/icbt/internal/app/handler"
	mw "github.com/dropwhile/icbt/internal/app/middleware"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/session"
	"github.com/dropwhile/icbt/resources"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
)

type API struct {
	*chi.Mux
	Db      *model.DB
	SessMgr *session.SessionMgr
	Tpl     resources.TemplateMap
}

func (api *API) Close() {
	api.SessMgr.Close()
}

func NewAPI(db *model.DB, tpl resources.TemplateMap, csrfKey []byte, isProd bool) *API {
	api := &API{
		SessMgr: session.NewDBSessionManager(db.GetPool()),
		Mux:     chi.NewRouter(),
		Db:      db,
		Tpl:     tpl,
	}

	// Router/Middleware //
	r := api.Mux
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.GetHead)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(api.SessMgr.LoadAndSave)
	r.Use(csrf.Protect(
		csrfKey,
		csrf.Secure(isProd),                // false in development only!
		csrf.Path("/"),                     // setup path so csrf works _between_ pages (eg. htmx calls)
		csrf.RequestHeader("X-CSRF-Token"), // Must be in CORS Allowed and Exposed Headers
	))
	r.Use(mw.LoadAuth(db, api.SessMgr))

	ah := &ah.Handler{Db: api.Db, Tpl: api.Tpl, SessMgr: api.SessMgr}

	// Routing //
	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(mw.RequireAuth)
		// acccount/settings
		r.Get("/settings", ah.ShowSettings)
		r.Post("/settings", ah.UpdateSettings)
		r.Delete("/settings", ah.DeleteAccount)
		// logout
		r.Post("/logout", ah.Logout)
		// dashboard/events/earmarks/etc
		r.Get("/dashboard", ah.ShowDashboard)
		r.Get("/events", ah.ListEvents)
		r.Get("/events/{refId:[0-9a-hjkmnp-tv-z]+}", ah.ShowEvent)
		r.Delete("/events/{refId:[0-9a-hjkmnp-tv-z]+}", ah.DeleteEvent)
		//r.Get("/create-event", ah.ShowCreateEvent)
		//r.Post("/create-event", ah.CreateEvent)
		r.Get("/earmarks", ah.ListEarmarks)
		r.Delete("/earmarks/{refId:[0-9a-hjkmnp-tv-]+}", ah.DeleteEarmark)
		//r.Get("/profile/{userRefId:[a-zA-Z-]+}", ah.ShowProfile)
	})

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/", ah.ShowIndex)
		// login
		r.Post("/login", ah.Login)
		r.Get("/login", ah.ShowLoginForm)
		// forgot password
		//r.Get("/forgot-password", ah.ShowForgotPassword)
		//r.Put("/forgot-password", ah.ResetPassword)
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

	return api
}
