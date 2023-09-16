package main

import (
	"context"
	"crypto/sha1"
	_ "database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/dropwhile/icbt/internal/app"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/resources"
	pgxz "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"golang.org/x/crypto/pbkdf2"
)

// ServerVersion holds the server version string
var ServerVersion = "no-version"

func main() {
	// parse env vars //

	// log format
	viper.MustBindEnv("log_format")
	viper.SetDefault("log_format", "json")
	logFormat := viper.GetString("log_format")
	if logFormat == "plain" {
		log.Logger = log.Output(
			zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
		)
	}

	// debug logging or not
	viper.MustBindEnv("log_level")
	viper.SetDefault("log_level", "info")
	logLevel := viper.GetString("log_level")
	switch strings.ToLower(logLevel) {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("setting log level to debug")
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
		log.Trace().Msg("setting log level to trace")
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		log.Info().Msg("setting log level to info")
	}

	// prod mode (secure cookies) or not
	viper.MustBindEnv("production")
	viper.SetDefault("production", "true")
	isProd := viper.GetBool("production")
	log.Debug().Msgf("prod mode is %t", isProd)

	// listen address/port
	viper.MustBindEnv("bind_address")
	viper.SetDefault("bind_address", "127.0.0.1")
	viper.MustBindEnv("bind_port")
	viper.SetDefault("bind_port", "8000")
	listenAddr := viper.GetString("bind_address")
	if listenAddr == "" {
		log.Fatal().Msg("listen address not specified")
	}
	listenPort := viper.GetInt("bind_port")
	if listenPort == 0 {
		log.Fatal().Msg("listen port not specified")
	}
	listenHostPort := fmt.Sprintf("%s:%d", listenAddr, listenPort)

	// load templates
	viper.MustBindEnv("tpl_dir")
	viper.SetDefault("tpl_dir", "embed")
	tplDir := path.Clean(viper.GetString("tpl_dir"))
	if tplDir == "embed" {
		log.Debug().Msg("using embedded templates")
	} else {
		if tplDir == "" {
			log.Fatal().Msg("template dir not specified")
		}
		_, err := os.Stat(tplDir)
		if err != nil && os.IsNotExist(err) {
			log.Fatal().Msgf("template dir does not exist: %s", tplDir)
		}
		log.Debug().Msgf("template dir set to: %s", tplDir)
	}
	templates := resources.MustParseTemplates(tplDir)

	// static resources
	viper.MustBindEnv("static_dir")
	viper.SetDefault("static_dir", "embed")
	staticDir := path.Clean(viper.GetString("static_dir"))
	if staticDir == "embed" {
		log.Debug().Msgf("using embedded static")
	} else {
		log.Debug().Msgf("static dir set to: %s", staticDir)
	}
	staticFS := resources.NewStaticFS(staticDir)

	// database
	viper.MustBindEnv("db_dsn")
	dbDSN := viper.GetString("db_dsn")
	if dbDSN == "" {
		log.Fatal().Msg("database connection info not supplied")
	}
	var dbpool *pgxpool.Pool
	if zerolog.GlobalLevel() == zerolog.TraceLevel {
		config, err := pgxpool.ParseConfig(dbDSN)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to connect to database")
		}
		config.ConnConfig.Tracer = &tracelog.TraceLog{
			Logger:   pgxz.NewLogger(log.Logger),
			LogLevel: tracelog.LogLevelTrace,
		}
		dbpool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to connect to database")
		}
	} else {
		var err error
		dbpool, err = pgxpool.New(context.Background(), dbDSN)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to connect to database")
		}
	}
	defer dbpool.Close()

	model := model.SetupFromDb(dbpool)
	defer model.Close()

	// csrf Key
	viper.MustBindEnv("csrf_key")
	csrfKeyInput := viper.GetString("csrf_key")
	if csrfKeyInput == "" {
		log.Fatal().Msg("csrf key not supplied")
	}

	// generate csrfKey based on input, using pdkdf2 to stretch/shrink
	// as needed to fit 32 byte key requirement
	csrfKeyBytes := pbkdf2.Key(
		[]byte(csrfKeyInput),
		[]byte("C/RWyRGBRXSCL5st5bFsPstuKQNDpRIix0vUlQ4QP80="),
		4096,
		32,
		sha1.New,
	)

	// routing/handlers
	r := app.NewAPI(model, templates, csrfKeyBytes, isProd)
	defer r.Close()
	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.FS(staticFS))))

	log.Info().Msgf("listening on %s", listenHostPort)
	server := &http.Server{
		Addr:              listenHostPort,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		log.Info().Msg("server closed")
	} else if err != nil {
		log.Fatal().Err(err).Msg("error starting server")
	}
}
