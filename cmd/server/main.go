package main

import (
	"context"
	_ "database/sql"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/api"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/util"
	"github.com/dropwhile/icbt/resources"
)

// ServerVersion holds the server version string
var ServerVersion = "no-version"

func main() {
	//--------------//
	// parse config //
	//--------------//

	config, err := ParseConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse config")
	}

	if config.LogFormat == "plain" {
		log.Logger = util.NewLogger(os.Stderr)
	}
	zerolog.SetGlobalLevel(config.LogLevel)
	log.Info().Msgf("setting log level: %s", config.LogLevel.String())

	if config.TemplateDir == "embed" {
		log.Debug().Msg("templates: embedded")
	} else {
		log.Debug().Msgf("templates: dir=%s", config.TemplateDir)
	}
	templates, err := resources.ParseTemplates(config.TemplateDir)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse templates")
	}

	if config.StaticDir == "embed" {
		log.Debug().Msgf("static: embedded")
	} else {
		log.Debug().Msgf("static: dir=%s", config.StaticDir)
	}

	log.Info().Msgf("prod mode: %t", config.Production)

	//--------------------//
	// configure services //
	//--------------------//

	// setup db pool & models
	db, err := model.SetupFromDsn(config.DatabaseDSN)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	// configure mailer
	mailer := util.NewMailer(
		config.SMTPHost,
		config.SMTPPort,
		config.SMTPHostname,
		config.SMTPUser,
		config.SMTPPass,
	)

	// routing/handlers
	r := api.New(
		db, templates, mailer,
		config.HMACKeyBytes,
		config.CSRFKeyBytes,
		config.Production,
	)
	defer r.Close()

	// serve static files dir as /static/*
	staticFS := resources.NewStaticFS(config.StaticDir)
	r.Get("/static/*", resources.FileServer(staticFS, "/static"))
	// some other single item static files
	r.Get("/favicon.ico", resources.ServeSingle(staticFS, "img/favicon.ico"))
	r.Get("/robots.txt", resources.ServeSingle(staticFS, "robots.txt"))

	server := &http.Server{
		Addr:              config.ListenHostPort,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	// signal handling && graceful shutdown
	idleConnsClosed := make(chan struct{})
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals

		// We received an interrupt signal, shut down.
		log.Info().Msg("HTTP server shutting down...")
		if err := server.Shutdown(context.Background()); err != nil {
			if err != nil {
				// Error from closing listeners, or context timeout:
				log.Error().Err(err).Msg("HTTP server shutdown error")
			}
		}
		close(idleConnsClosed)
	}()

	// listen
	log.Info().Msg("starting up...")
	log.Info().Msgf("listening: %s", config.ListenHostPort)
	log.Info().Msgf("server version: %s", ServerVersion)
	err = server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		log.Info().Msg("server closed")
	} else if err != nil {
		log.Fatal().Err(err).Msg("error starting server")
	}

	<-idleConnsClosed
}
