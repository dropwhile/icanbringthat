package main

import (
	_ "database/sql"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/envconfig"
	"github.com/dropwhile/icbt/internal/logger"
	"github.com/dropwhile/icbt/internal/mail"
	"github.com/dropwhile/icbt/resources"
)

// ServerVersion holds the server version string
var ServerVersion = "no-version"

//go:generate stringer -type=Job
type Job int

const (
	jNotifier Job = iota + 1
	jArchiver
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
		log.Fatal().Err(err).Msg("failed to parse config")
		return
	}

	if config.LogFormat == "plain" {
		log.Logger = logger.NewLogger(os.Stderr)
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
		return
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
		return
	}
	defer dbpool.Close()

	//----------------//
	// configure jobs //
	//----------------//
	workerConfig := &WorkerConfig{}
	if err := env.Parse(workerConfig); err != nil {
		log.Fatal().Err(err).Msg("failed to parse config")
		return
	}
	jobs := make(map[Job]bool)
	for _, v := range workerConfig.Jobs {
		switch strings.ToLower(v) {
		case "notifier":
			jobs[jNotifier] = true
		case "archiver":
			jobs[jArchiver] = true
		case "all":
			jobs[jNotifier] = true
			jobs[jArchiver] = true
		default:
			log.Fatal().Msgf("unknown job: %s", v)
			return
		}
	}

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
	log.Info().Msgf("server version: %s", ServerVersion)
	log.Info().Msg("starting up...")

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case sig := <-signals:
				log.Info().
					Msgf("Got %s.", sig.String())
				log.Info().Msg("Program will terminate now.")
				return
			case <-timer.C:
				if ok := jobs[jNotifier]; ok {
					if err := service.NotifyUsersPendingEvents(
						dbpool, mailer, templates, config.BaseURL,
					); err != nil {
						log.Error().Err(err).Msg("notifier error!!")
					}
				}
				if ok := jobs[jArchiver]; ok {
					if err := service.ArchiveOldEvents(dbpool); err != nil {
						log.Error().Err(err).Msg("archiver error!!")
					}
				}
				timer.Reset(timerInterval)
			}
		}
	}()

	// block
	wg.Wait()
}
