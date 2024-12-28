// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"context"
	_ "database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/quic-go/quic-go/http3"
	"github.com/redis/go-redis/v9"

	"github.com/dropwhile/icanbringthat/internal/app"
	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/app/resources"
	"github.com/dropwhile/icanbringthat/internal/envconfig"
	"github.com/dropwhile/icanbringthat/internal/logger"
	"github.com/dropwhile/icanbringthat/internal/mail"
	"github.com/dropwhile/icanbringthat/internal/util"
)

type ServerCmd struct{}

func (c *ServerCmd) Run() error {
	//--------------//
	// parse config //
	//--------------//

	config, err := envconfig.Parse()
	if err != nil {
		return fmt.Errorf("failed to parse config: %s", err)
	}

	// omit src line in production logs
	opts := &logger.Options{
		OmitSource: config.Production,
	}

	switch config.LogFormat {
	case "plain":
		logger.SetupLogging(logger.NewConsoleLogger, opts)
	default:
		logger.SetupLogging(logger.NewJsonLogger, opts)
	}

	logger.SetLevel(config.LogLevel)

	templateLoc := resources.Embed
	if config.TemplateDir == "fs" {
		templateLoc = resources.Filesystem
	}
	templates, err := resources.ParseTemplates(templateLoc)
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}

	slog.With(
		slog.Group("logging",
			"level", config.LogLevel,
			"trace", config.LogTrace,
		),
		slog.Group("templates",
			"location", config.TemplateDir,
		),
		slog.Group("static",
			"location", config.StaticDir,
		),
		slog.Bool("production", config.Production),
		slog.Bool("rpc_api", config.RpcApi),
	).Info("startup configuration")

	//--------------------//
	// configure services //
	//--------------------//

	// setup db pool & models
	db, err := model.SetupDBPool(config.DatabaseDSN, config.LogTrace)
	if err != nil {
		slog.With("error", err).
			Error("failed to connect to database")
		return fmt.Errorf("failed to connect to database")
	}
	defer db.Close()

	redisOpt, err := redis.ParseURL(config.RedisDSN)
	if err != nil {
		slog.With("error", err).
			Error("failed to connect to redis")
		return fmt.Errorf("failed to connect to redis")
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
		WebhookCreds:   config.WebhookCreds,
		CSRFKeyBytes:   config.CSRFKeyBytes,
		HMACKeyBytes:   config.HMACKeyBytes,
		Production:     config.Production,
		BaseURL:        config.BaseURL,
		RequestLogging: config.LogTrace,
		RpcApi:         config.RpcApi,
	}
	r, err := app.New(db, rdb, templates, mailer, appConfig)
	if err != nil {
		slog.With("error", err).
			Error("failed to create server")
		return fmt.Errorf("failed to create server")
	}
	defer r.Close()

	// serve static files dir as /static/*
	staticLoc := resources.Embed
	if config.StaticDir == "fs" {
		staticLoc = resources.Filesystem
	}
	staticFS := resources.NewStaticFS(staticLoc)
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
		slog.Info("Server shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			slog.With("error", err).Error("HTTP server shutdown error")
		}

		if quicServer != nil {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
			defer cancel()
			if err := quicServer.Shutdown(ctx); err != nil {
				// Error from closing listeners, or context timeout:
				slog.With("error", err).Error("HTTP/3 server shutdown error")
			}
		}
		close(idleConnsClosed)
	}()

	vinfo, _ := util.GetVersion()

	// listen
	slog.
		With("version", vinfo.Version).
		With("go", vinfo.GoVersion).
		Info("starting up...")
	if config.TLSCert != "" && config.TLSKey != "" {
		if config.WithQuic {
			// add quic headers to https/tls server
			if config.WithQuic {
				handler := server.Handler
				server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					quicServer.SetQUICHeaders(w.Header()) // #nosec G104 -- this only fails if port cant be determined
					handler.ServeHTTP(w, r)
				})
			}

			// start up http3/quic server
			go func() {
				slog.With(
					slog.String("proto", "https/quic"),
					slog.String("listen", config.Listen),
				).Info("listening")
				if err := quicServer.ListenAndServeTLS(
					config.TLSCert, config.TLSKey,
				); err != nil && !errors.Is(err, http.ErrServerClosed) {
					slog.With("error", err).
						Error("HTTP/3 server error")
					os.Exit(1)
					return
				}
			}()
		}
		// startup https3/tls server
		go func() {
			slog.With(
				slog.String("proto", "https/tls"),
				slog.String("listen", config.Listen),
			).Info("listening")
			if err := server.ListenAndServeTLS(
				config.TLSCert, config.TLSKey,
			); err != nil && !errors.Is(err, http.ErrServerClosed) {
				slog.With("error", err).
					Error("HTTPS server error")
				os.Exit(1)
				return
			}
		}()
	} else {
		slog.With(
			slog.String("proto", "http"),
			slog.String("listen", config.Listen),
		).Info("listening")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.With("error", err).
				Error("HTTP server error")
			return fmt.Errorf("HTTP server error")
		}
	}

	<-idleConnsClosed
	return nil
}
