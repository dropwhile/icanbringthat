package main

import (
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
	slog.Info("setting log level", "level", config.LogLevel)

	if config.TemplateDir == "embed" {
		slog.Debug("templates", "location", "embedded")
	} else {
		slog.Debug("templates", "location", config.TemplateDir)
	}
	templates, err := resources.ParseTemplates(config.TemplateDir)
	if err != nil {
		logger.Fatal("failed to parse templates")
		return
	}

	if config.StaticDir == "embed" {
		slog.Debug("static", "location", "embedded")
	} else {
		slog.Debug("static", "location", config.StaticDir)
	}

	slog.Info("prod mode", "mode", config.Production)

	//--------------------//
	// configure services //
	//--------------------//

	// setup dbpool pool & models
	dbpool, err := model.SetupDBPool(config.DatabaseDSN, config.LogTrace)
	if err != nil {
		logger.Fatal("failed to connect to database",
			"error", err)
		return
	}
	defer dbpool.Close()

	//----------------//
	// configure jobs //
	//----------------//
	workerConfig := &WorkerConfig{}
	if err := env.Parse(workerConfig); err != nil {
		logger.Fatal("failed to parse config", "error", err)
		return
	}

	jobList := NewJobList()
	err = jobList.AddByName(workerConfig.Jobs...)
	if err != nil {
		logger.Fatal("error adding worker jobs", "error", err)
		return
	}
	slog.Info("configured workers", "worklist", jobList.String())

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
	slog.Info("starting up...", "version", Version)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case sig := <-signals:
				slog.Info("Got", "signal", sig.String())
				slog.Info("Program will terminate now.")
				return
			case <-timer.C:
				if jobList.Contains(NotifierJob) {
					if err := service.NotifyUsersPendingEvents(
						dbpool, mailer, templates, config.BaseURL,
					); err != nil {
						slog.Error("notifier error!!", "error", err)
					}
				}
				if jobList.Contains(ArchiverJob) {
					if err := service.ArchiveOldEvents(dbpool); err != nil {
						slog.Error("archiver error!!", "error", err)
					}
				}
				timer.Reset(timerInterval)
			}
		}
	}()

	// block
	wg.Wait()
}
