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

	"github.com/quic-go/quic-go/http3"
	"github.com/redis/go-redis/v9"
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

	config, err := util.ParseConfig()
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

	// setup dbpool pool & models
	dbpool, err := model.SetupDBPool(config.DatabaseDSN)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer dbpool.Close()

	redisOpt, err := redis.ParseURL(config.RedisDSN)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to redis")
	}

	rdb := redis.NewClient(redisOpt)
	defer rdb.Close()

	// configure mailer
	mailer := util.NewMailer(
		config.SMTPHost,
		config.SMTPPort,
		config.SMTPHostname,
		config.SMTPUser,
		config.SMTPPass,
		config.MailFrom,
	)

	// routing/handlers
	r := api.New(
		dbpool, rdb,
		templates, mailer,
		config.HMACKeyBytes,
		config.CSRFKeyBytes,
		config.Production,
		config.BaseURL,
		config.WebhookCreds,
	)
	defer r.Close()

	// serve static files dir as /static/*
	staticFS := resources.NewStaticFS(config.StaticDir)
	r.Get("/static/*", resources.FileServer(staticFS, "/static"))
	// some other single item static files
	r.Get("/favicon.ico", resources.ServeSingle(staticFS, "img/favicon.ico"))
	r.Get("/robots.txt", resources.ServeSingle(staticFS, "robots.txt"))

	server := &http.Server{
		Addr:              config.Listen,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	var quicServer *http3.Server
	if config.TLSCert != "" && config.TLSKey != "" {
		quicServer = &http3.Server{
			Addr:    config.Listen,
			Handler: r,
		}
	}

	// signal handling && graceful shutdown
	idleConnsClosed := make(chan struct{})
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals

		// We received an interrupt signal, shut down.
		log.Info().Msg("Server shutting down...")
		if err := server.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Error().Err(err).Msg("HTTP server shutdown error")
		}
		if quicServer != nil {
			if err := quicServer.CloseGracefully(time.Second * 2); err != nil {
				// Error from closing listeners, or context timeout:
				log.Error().Err(err).Msg("HTTP/3 server shutdown error")
			}
		}
		close(idleConnsClosed)
	}()

	// listen
	log.Info().Msg("starting up...")
	log.Info().Msgf("server version: %s", ServerVersion)
	if config.TLSCert != "" && config.TLSKey != "" {
		if config.WithQuic {
			// add quic headers to https/tls server
			if config.WithQuic {
				handler := server.Handler
				server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					quicServer.SetQuicHeaders(w.Header())
					handler.ServeHTTP(w, r)
				})
			}

			// start up http3/quic server
			go func() {
				log.Info().Msgf("listening(https/quic): %s", config.Listen)
				if err := quicServer.ListenAndServeTLS(config.TLSCert, config.TLSKey); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Fatal().Err(err).Msg("HTTP/3 server error")
				}
			}()
		}
		// startup https3/tls server
		go func() {
			log.Info().Msgf("listening(https/tls):  %s", config.Listen)
			if err := server.ListenAndServeTLS(config.TLSCert, config.TLSKey); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatal().Err(err).Msg("HTTP server error")
			}
		}()
	} else {
		log.Info().Msgf("listening(http): %s", config.Listen)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("HTTP server error")
		}
	}
	<-idleConnsClosed
}
