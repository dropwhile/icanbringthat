// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"context"
	_ "database/sql"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/caarlos0/env/v10"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/app/resources"
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/envconfig"
	"github.com/dropwhile/icanbringthat/internal/logger"
	"github.com/dropwhile/icanbringthat/internal/mail"
	"github.com/dropwhile/icanbringthat/internal/util"
)

type Job string

const (
	NotifierJob Job = "notifier"
	ArchiverJob Job = "archiver"
)

type WorkerConfig struct {
	Jobs []string `env:"JOBS" envDefault:"all"`
}

type RunCmd struct{}

func (c *RunCmd) Run() error {
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
		slog.With("error", err).
			Error("failed to parse templates")
		return fmt.Errorf("failed to parse templates")
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
	).Info("startup configuration")

	//--------------------//
	// configure services //
	//--------------------//

	// setup db pool & models & servicer
	db, err := model.SetupDBPool(config.DatabaseDSN, config.LogTrace)
	if err != nil {
		slog.With("error", err).
			Error("failed to connect to database")
		return fmt.Errorf("failed to connect to database")
	}
	defer db.Close()
	service := service.New(service.Options{Db: db})

	//----------------//
	// configure jobs //
	//----------------//
	workerConfig := &WorkerConfig{}
	if err := env.Parse(workerConfig); err != nil {
		slog.With("error", err).
			Error("failed to parse config")
		return fmt.Errorf("failed to parse config")
	}

	jobList := NewJobList()
	err = jobList.AddByName(workerConfig.Jobs...)
	if err != nil {
		logger.Fatal("error adding worker jobs", "error", err)
		return fmt.Errorf("error adding worker jobs: %s", err)
	}
	slog.With("worklist", jobList).
		Info("configured workers")

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

	// signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// timer
	timerInterval := 10 * time.Minute
	timer := time.NewTimer(0)
	defer timer.Stop()

	vinfo, _ := util.GetVersion()
	slog.
		With("version", vinfo.Version).
		With("go", runtime.Version()).
		Info("starting up...")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case sig := <-signals:
				slog.With("signal", sig.String()).
					Info("Got signal")
				slog.Info("Program will terminate now.")
				return
			case <-timer.C:
				if jobList.Contains(NotifierJob) {
					if err := service.NotifyUsersPendingEvents(
						context.Background(), mailer, templates, config.BaseURL,
					); err != nil {
						slog.With("erorr", err).
							Error("notifier error!!")
					}
				}
				if jobList.Contains(ArchiverJob) {
					if err := service.ArchiveOldEvents(context.Background()); err != nil {
						slog.With("erorr", err).
							Error("archiver error!!")
					}
				}
				timer.Reset(timerInterval)
			}
		}
	}()

	// block
	wg.Wait()
	return nil
}
