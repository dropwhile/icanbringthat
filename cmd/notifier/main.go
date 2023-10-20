package main

import (
	_ "database/sql"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/service"
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

	// configure mailer
	mailer := util.NewMailer(
		config.SMTPHost,
		config.SMTPPort,
		config.SMTPHostname,
		config.SMTPUser,
		config.SMTPPass,
	)

	// signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// timer
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	var wg sync.WaitGroup
	log.Info().Msg("starting up...")

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case sig := <-signals:
			switch sig {
			case syscall.SIGTERM:
				log.Info().Msg("Got kill signal.")
				log.Info().Msg("Program will terminate now.")
			case syscall.SIGINT:
				log.Info().Msg("Got CTRL+C signal.")
				log.Info().Msg("Program will terminate now.")
			default:
				log.Info().Stringer("signal", sig).Msg("Ignoring signal")
			}
		case <-ticker.C:
			err := service.NotifyUsersPendingEvents(
				dbpool, mailer, templates, config.BaseURL,
			)
			if err != nil {
				log.Error().Err(err).Msg("error!!")
			}
		}
	}()

	// block
	wg.Wait()
}
