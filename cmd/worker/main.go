package main

import (
	"context"
	_ "database/sql"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/caarlos0/env/v10"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/resources"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/envconfig"
	"github.com/dropwhile/icbt/internal/logger"
	"github.com/dropwhile/icbt/internal/mail"
)

// Version holds the server version string
var Version = "no-version"

type Job string

const (
	NotifierJob Job = "notifier"
	ArchiverJob Job = "archiver"
)

type WorkerConfig struct {
	Jobs []string `env:"JOBS" envDefault:"all"`
}

func main() {
	ctx := context.Background()

	//--------------//
	// parse config //
	//--------------//

	config, err := envconfig.Parse()
	if err != nil {
		log.Fatal("failed to parse config")
		return
	}

	if config.LogFormat == "plain" {
		logger.SetupLogging(logger.NewConsoleLogger, nil)
	} else {
		logger.SetupLogging(logger.NewJsonLogger, nil)
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
		logger.Fatal(ctx, "failed to parse templates")
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

	//----------------//
	// configure jobs //
	//----------------//
	workerConfig := &WorkerConfig{}
	if err := env.Parse(workerConfig); err != nil {
		logger.Fatal(ctx, "failed to parse config",
			logger.Err(err))
		return
	}

	jobList := NewJobList()
	err = jobList.AddByName(workerConfig.Jobs...)
	if err != nil {
		logger.Fatal(ctx, "error adding worker jobs",
			logger.Err(err))
		return
	}
	logger.Info(ctx, "configured workers",
		slog.String("worklist", jobList.String()))

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

	var wg sync.WaitGroup
	logger.Info(ctx, "starting up...",
		slog.String("version", Version))

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case sig := <-signals:
				logger.Info(ctx, "Got",
					slog.String("signal", sig.String()))
				logger.Info(ctx, "Program will terminate now.")
				return
			case <-timer.C:
				if jobList.Contains(NotifierJob) {
					if err := service.NotifyUsersPendingEvents(
						dbpool, mailer, templates, config.BaseURL,
					); err != nil {
						logger.Error(ctx, "notifier error!!",
							logger.Err(err))
					}
				}
				if jobList.Contains(ArchiverJob) {
					if err := service.ArchiveOldEvents(dbpool); err != nil {
						logger.Error(ctx, "archiver error!!",
							logger.Err(err))
					}
				}
				timer.Reset(timerInterval)
			}
		}
	}()

	// block
	wg.Wait()
}
