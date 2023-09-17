package app

import (
	ah "github.com/dropwhile/icbt/internal/app/handler"
	mw "github.com/dropwhile/icbt/internal/app/middleware"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/session"
	res "github.com/dropwhile/icbt/resources"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/rs/zerolog/log"
)

type API struct {
	*chi.Mux
	Db      *model.DB
	SessMgr *session.SessionMgr
	Tpl     res.TemplateMap
}

func (api *API) Close() {
	api.SessMgr.Close()
}

func NewAPI(db *model.DB, tpl res.TemplateMap, csrfKey []byte, isProd bool) *API {
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
	if log.Debug().Enabled() {
		r.Use(mw.NewDebubRequestLogger())
	}
	r.Use(middleware.Recoverer)
	r.Use(api.SessMgr.LoadAndSave)
	r.Use(csrf.Protect(
		csrfKey,
		// false in development only!
		csrf.Secure(isProd),
		// setup path so csrf works _between_ pages (eg. htmx calls)
		csrf.Path("/"),
		// Must be in CORS Allowed and Exposed Headers
		csrf.RequestHeader("X-CSRF-Token"),
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
		r.Post("/events", ah.CreateEvent)
		r.Get("/events/add", ah.ShowCreateEventForm)
		r.Get("/events/{eRefId:[0-9a-z]+}", ah.ShowEvent)
		r.Post("/events/{eRefId:[0-9a-z]+}", ah.UpdateEvent)
		r.Delete("/events/{eRefId:[0-9a-z]+}", ah.DeleteEvent)
		r.Get("/events/{eRefId:[0-9a-z]+}/edit", ah.ShowEditEventForm)
		r.Post("/events/{eRefId:[0-9a-z]+}/items", ah.CreateEventItem)
		r.Get("/events/{eRefId:[0-9a-z]+}/items/add", ah.ShowCreateEventItemForm)
		r.Post("/events/{eRefId:[0-9a-z]+}/items/{iRefId:[0-9a-z]+}", ah.UpdateEventItem)
		r.Delete("/events/{eRefId:[0-9a-z]+}/items/{iRefId:[0-9a-z]+}", ah.DeleteEventItem)
		r.Get("/events/{eRefId:[0-9a-z]+}/items/{iRefId:[0-9a-z]+}/edit", ah.ShowEventItemEditForm)
		r.Post("/events/{eRefId:[0-9a-z]+}/items/{iRefId:[0-9a-z]+}/earmarks", ah.CreateEarmark)
		r.Get("/events/{eRefId:[0-9a-z]+}/items/{iRefId:[0-9a-z]+}/earmarks/add", ah.ShowCreateEarmarkForm)
		r.Get("/earmarks", ah.ListEarmarks)
		r.Delete("/earmarks/{mRefId:[0-9a-z]+}", ah.DeleteEarmark)
		/*
			r.Get("/earmarks/add", ah.ShowCreateEarmarkForm)
			r.Post("/earmarks/add", ah.CreateEarmark)
			r.Get("/earmarks/{mRefId:[0-9a-z]+}", ah.ShowEarmark)
			r.Post("/earmarks/{mRefId:[0-9a-z]+}", ah.UpdateEarmark)
			r.Get("/earmarks/{mRefId:[0-9a-z]+}/edit", ah.DeleteEarmark)
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
		// r.Get("/forgot-password", ah.ShowForgotPassword)
		// r.Put("/forgot-password", ah.ResetPassword)
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
