package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	_ "database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/dropwhile/icbt/internal/app"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/util"
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
	viper.MustBindEnv("LOG_FORMAT")
	viper.SetDefault("LOG_FORMAT", "json")
	logFormat := viper.GetString("LOG_FORMAT")
	if logFormat == "plain" {
		log.Logger = util.NewLogger(os.Stderr)
	}

	log.Info().Msg("starting up...")
	log.Info().Msgf("server version: %s", ServerVersion)

	// debug logging or not
	viper.MustBindEnv("LOG_LEVEL")
	viper.SetDefault("LOG_LEVEL", "info")
	logLevel := viper.GetString("LOG_LEVEL")
	switch strings.ToLower(logLevel) {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("setting log level: debug")
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
		log.Trace().Msg("setting log level: trace")
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		log.Info().Msg("setting log level: info")
	}

	// prod mode (secure cookies) or not
	viper.MustBindEnv("PRODUCTION")
	viper.SetDefault("PRODUCTION", "true")
	isProd := viper.GetBool("PRODUCTION")
	log.Debug().Msgf("prod mode: %t", isProd)

	// listen address/port
	viper.MustBindEnv("BIND_ADDRESS")
	viper.SetDefault("BIND_ADDRESS", "127.0.0.1")
	viper.MustBindEnv("BIND_PORT")
	viper.SetDefault("BIND_PORT", "8000")
	listenAddr := viper.GetString("BIND_ADDRESS")
	if listenAddr == "" {
		log.Fatal().Msg("listen address not specified")
	}
	listenPort := viper.GetInt("BIND_PORT")
	if listenPort == 0 {
		log.Fatal().Msg("listen port not specified")
	}
	listenHostPort := fmt.Sprintf("%s:%d", listenAddr, listenPort)

	// load templates
	viper.MustBindEnv("TPL_DIR")
	viper.SetDefault("TPL_DIR", "embed")
	tplDir := path.Clean(viper.GetString("TPL_DIR"))
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
		log.Debug().Msgf("template dir: %s", tplDir)
	}
	templates := resources.MustParseTemplates(tplDir)

	// static resources
	viper.MustBindEnv("STATIC_DIR")
	viper.SetDefault("STATIC_DIR", "embed")
	staticDir := path.Clean(viper.GetString("STATIC_DIR"))
	if staticDir == "embed" {
		log.Debug().Msgf("using embedded static")
	} else {
		log.Debug().Msgf("static dir: %s", staticDir)
	}
	staticFS := resources.NewStaticFS(staticDir)

	// database
	viper.MustBindEnv("DB_DSN")
	dbDSN := viper.GetString("DB_DSN")
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
	viper.MustBindEnv("CSRF_KEY")
	csrfKeyInput := viper.GetString("CSRF_KEY")
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
		sha256.New,
	)
	hmacKeyBytes := pbkdf2.Key(
		csrfKeyBytes,
		[]byte("i4L04cpiG6JebX5brY53sBBqCyX16IwbjagbMkytmQQ="),
		4096,
		32,
		sha256.New,
	)

	viper.MustBindEnv("SMTP_HOSTNAME")
	smtpHostname := viper.GetString("SMTP_HOSTNAME")
	if smtpHostname == "" {
		log.Fatal().Msg("smtp mail host name")
	}

	viper.MustBindEnv("SMTP_HOST")
	smtpHost := viper.GetString("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = smtpHostname
	}

	viper.MustBindEnv("SMTP_PORT")
	smtpPort := viper.GetInt("SMTP_PORT")
	if smtpPort == 0 {
		log.Fatal().Msg("smtp mail port")
	}

	viper.MustBindEnv("SMTP_USER")
	smtpUser := viper.GetString("SMTP_USER")
	if smtpUser == "" {
		log.Fatal().Msg("smtp mail username")
	}

	viper.MustBindEnv("SMTP_PASS")
	smtpPass := viper.GetString("SMTP_PASS")
	if smtpPass == "" {
		log.Fatal().Msg("smtp mail password")
	}

	mailer := util.NewMailer(smtpHost, smtpPort, smtpHostname, smtpUser, smtpPass)

	// routing/handlers
	r := app.NewAPI(model, templates, mailer, csrfKeyBytes, hmacKeyBytes, isProd)
	defer r.Close()
	// serve static files
	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.FS(staticFS))))
	// serve favicon
	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		f, err := staticFS.Open("img/favicon.ico")
		if err != nil {
			log.Debug().Err(err).Msg("cant read favicon")
			http.Error(w, "Not Found", 404)
			return
		}
		defer f.Close()
		b, err := io.ReadAll(f)
		if err != nil {
			log.Debug().Err(err).Msg("cant read favicon")
			http.Error(w, "Not Found", 404)
			return
		}
		http.ServeContent(w, r, "favicon.ico", time.Time{}, bytes.NewReader(b))
	})
	r.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		f, err := staticFS.Open("robots.txt")
		if err != nil {
			log.Debug().Err(err).Msg("cant read robots.txt")
			http.Error(w, "Not Found", 404)
			return
		}
		defer f.Close()
		b, err := io.ReadAll(f)
		if err != nil {
			log.Debug().Err(err).Msg("cant read robots.txt")
			http.Error(w, "Not Found", 404)
			return
		}
		http.ServeContent(w, r, "robots.txt", time.Time{}, bytes.NewReader(b))
	})

	log.Info().Msgf("listening: %s", listenHostPort)
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
