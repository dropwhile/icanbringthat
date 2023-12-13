package main

import (
	"context"
	_ "database/sql"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/quic-go/quic-go/http3"
	"github.com/redis/go-redis/v9"

	"github.com/dropwhile/icbt/internal/app"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/resources"
	"github.com/dropwhile/icbt/internal/envconfig"
	"github.com/dropwhile/icbt/internal/logger"
	"github.com/dropwhile/icbt/internal/mail"
)

// Version holds the server version string
var Version = "no-version"

func main() {
	ctx := context.Background()

	//--------------//
	// parse config //
	//--------------//

	config, err := envconfig.Parse()
	if err != nil {
		log.Fatalf("failed to parse config: %s", err)
		return
	}

	if config.LogFormat == "plain" {
		logger.SetupLogging(logger.NewConsoleLogger, nil)
	} else {
		logger.SetupLogging(logger.NewJsonLogger, nil)
	}

	if config.Production {
		// trace level not allowed in prod mode,
		// as it may expose private data in sql
		// queries
	}
	logger.SetLevel(config.LogLevel)
	logger.Info(ctx, "setting log level",
		slog.Any("level", config.LogLevel))

	if config.TemplateDir == "embed" {
		logger.Debug(ctx, "templates",
			slog.String("location", "embedded"))
	} else {
		logger.Debug(ctx, "templates",
			slog.String("location", config.TemplateDir))
	}
	templates, err := resources.ParseTemplates(config.TemplateDir)
	if err != nil {
		logger.Fatal(ctx, "failed to parse templates",
			logger.Err(err))
		return
	}

	if config.StaticDir == "embed" {
		logger.Debug(ctx, "static",
			slog.String("location", "embedded"))
	} else {
		logger.Debug(ctx, "static",
			slog.String("location", config.StaticDir))
	}

	logger.Info(ctx, "prod mode",
		slog.Bool("mode", config.Production))

	//--------------------//
	// configure services //
	//--------------------//

	// setup dbpool pool & models
	dbpool, err := model.SetupDBPool(config.DatabaseDSN)
	if err != nil {
		logger.Fatal(ctx, "failed to connect to database",
			logger.Err(err))
		return
	}
	defer dbpool.Close()

	redisOpt, err := redis.ParseURL(config.RedisDSN)
	if err != nil {
		logger.Fatal(ctx, "failed to connect to redis",
			logger.Err(err))
		return
	}

	rdb := redis.NewClient(redisOpt)
	defer rdb.Close()

	// configure mailer
	mailConfig := &mail.Config{
		Hostname:    config.SMTPHostname,
		Host:        config.SMTPHost,
		Port:        config.SMTPPort,
		User:        config.SMTPUser,
		Pass:        config.SMTPPass,
		DefaultFrom: config.MailFrom,
	}
	mailer := mail.NewMailer(mailConfig)

	// routing/handlers
	appConfig := &app.Config{
		WebhookCreds: config.WebhookCreds,
		CSRFKeyBytes: config.CSRFKeyBytes,
		HMACKeyBytes: config.HMACKeyBytes,
		Production:   config.Production,
		BaseURL:      config.BaseURL,
	}
	r := app.New(dbpool, rdb, templates, mailer, appConfig)
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
		logger.Info(ctx, "Server shutting down...")
		if err := server.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			logger.Error(ctx, "HTTP server shutdown error",
				logger.Err(err))
		}
		if quicServer != nil {
			if err := quicServer.CloseGracefully(time.Second * 2); err != nil {
				// Error from closing listeners, or context timeout:
				logger.Error(ctx, "HTTP/3 server shutdown error",
					logger.Err(err))
			}
		}
		close(idleConnsClosed)
	}()

	// listen
	logger.Info(ctx, "starting up...",
		slog.String("version", Version))
	if config.TLSCert != "" && config.TLSKey != "" {
		if config.WithQuic {
			// add quic headers to https/tls server
			if config.WithQuic {
				handler := server.Handler
				server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					quicServer.SetQuicHeaders(w.Header()) // #nosec G104 -- this only fails if port cant be determined
					handler.ServeHTTP(w, r)
				})
			}

			// start up http3/quic server
			go func() {
				logger.Info(ctx, "listening",
					slog.String("proto", "https/quic"),
					slog.String("listen", config.Listen))
				if err := quicServer.ListenAndServeTLS(config.TLSCert, config.TLSKey); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Fatal(ctx, "HTTP/3 server error",
						logger.Err(err))
					return
				}
			}()
		}
		// startup https3/tls server
		go func() {
			logger.Info(ctx, "listening",
				slog.String("proto", "https/tls"),
				slog.String("listen", config.Listen))
			if err := server.ListenAndServeTLS(config.TLSCert, config.TLSKey); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Fatal(ctx, "HTTP server error",
					logger.Err(err))
				return
			}
		}()
	} else {
		logger.Info(ctx, "listening",
			slog.String("proto", "http"),
			slog.String("listen", config.Listen))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(ctx, "HTTP server error",
				logger.Err(err))
			return
		}
	}
	<-idleConnsClosed
}
