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
	slog.With("level", config.LogLevel).
		Info("setting log level")

	if config.TemplateDir == "embed" {
		slog.With("location", "embedded").
			Debug("templates")
	} else {
		slog.With("location", config.TemplateDir).
			Debug("templates")
	}
	templates, err := resources.ParseTemplates(config.TemplateDir)
	if err != nil {
		slog.With("error", err).
			Error("failed to parse templates")
		os.Exit(1)
		return
	}

	if config.StaticDir == "embed" {
		slog.With("location", "embedded").
			Debug("static")
	} else {
		slog.With("location", config.StaticDir).
			Debug("static")
	}

	slog.With("mode", config.Production).
		Info("prod mode")

	//--------------------//
	// configure services //
	//--------------------//

	// setup dbpool pool & models
	dbpool, err := model.SetupDBPool(config.DatabaseDSN, config.LogTrace)
	if err != nil {
		slog.With("error", err).
			Error("failed to connect to database")
		os.Exit(1)
		return
	}
	defer dbpool.Close()

	//----------------//
	// configure jobs //
	//----------------//
	workerConfig := &WorkerConfig{}
	if err := env.Parse(workerConfig); err != nil {
		slog.With("error", err).
			Error("failed to parse config")
		os.Exit(1)
		return
	}

	jobList := NewJobList()
	err = jobList.AddByName(workerConfig.Jobs...)
	if err != nil {
		logger.Fatal("error adding worker jobs", "error", err)
		return
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

	var wg sync.WaitGroup
	slog.With("version", Version).
		Info("starting up...")

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
						dbpool, mailer, templates, config.BaseURL,
					); err != nil {
						slog.With("erorr", err).
							Error("notifier error!!")
					}
				}
				if jobList.Contains(ArchiverJob) {
					if err := service.ArchiveOldEvents(dbpool); err != nil {
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
}
