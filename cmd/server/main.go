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
	"time"

	"github.com/cactus/mlog"
	"github.com/dropwhile/icbt/internal/app"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/resources"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/spf13/viper"
	"golang.org/x/crypto/pbkdf2"
)

// ServerVersion holds the server version string
var ServerVersion = "no-version"

func main() {
	// parse env vars //

	// debug logging or not
	viper.MustBindEnv("debug")
	viper.SetDefault("debug", "false")
	logDebug := viper.GetBool("debug")
	if logDebug {
		mlog.SetFlags(mlog.Flags() | mlog.Ldebug)
		mlog.Debug("debug logging enabled")
	}

	// prod mode (secure cookies) or not
	viper.MustBindEnv("production")
	viper.SetDefault("production", "true")
	isProd := viper.GetBool("production")
	mlog.Debugf("prod mode is %t", isProd)

	// listen address/port
	viper.MustBindEnv("bind_address")
	viper.SetDefault("bind_address", "127.0.0.1")
	viper.MustBindEnv("bind_port")
	viper.SetDefault("bind_port", "8000")
	listenAddr := viper.GetString("bind_address")
	if listenAddr == "" {
		mlog.Fatal("listen address not specified")
	}
	listenPort := viper.GetInt("bind_port")
	if listenPort == 0 {
		mlog.Fatal("listen port not specified")
	}
	listenHostPort := fmt.Sprintf("%s:%d", listenAddr, listenPort)

	// load templates
	viper.MustBindEnv("tpl_dir")
	viper.SetDefault("tpl_dir", "embed")
	tplDir := path.Clean(viper.GetString("tpl_dir"))
	if tplDir == "embed" {
		mlog.Debugf("using embedded templates")
	} else {
		if tplDir == "" {
			mlog.Fatal("template dir not specified")
		}
		_, err := os.Stat(tplDir)
		if err != nil && os.IsNotExist(err) {
			mlog.Fatalf("template dir does not exist: %s", tplDir)
		}
		mlog.Debugf("template dir set to: %s", tplDir)
	}
	templates := resources.MustParseTemplates(tplDir)

	// static resources
	viper.MustBindEnv("static_dir")
	viper.SetDefault("static_dir", "embed")
	staticDir := path.Clean(viper.GetString("static_dir"))
	if staticDir == "embed" {
		mlog.Debugf("using embedded static")
	} else {
		mlog.Debugf("static dir set to: %s", staticDir)
	}
	staticFS := resources.NewStaticFS(staticDir)

	// database
	viper.MustBindEnv("db_dsn")
	dbDSN := viper.GetString("db_dsn")
	if dbDSN == "" {
		mlog.Fatal("database connection info not supplied")
	}
	var dbpool *pgxpool.Pool
	if logDebug {
		dbpool.Close()
		config, err := pgxpool.ParseConfig(dbDSN)
		if err != nil {
			mlog.Fatalf("failed to connect to database: %s", err)
		}
		config.ConnConfig.Tracer = &tracelog.TraceLog{
			Logger: tracelog.LoggerFunc(func(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
				mlog.Debugm("DB: "+msg, data)
			}),
			LogLevel: tracelog.LogLevelTrace,
		}
		dbpool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			mlog.Fatalf("failed to connect to database: %s", err)
		}
	} else {
		var err error
		dbpool, err = pgxpool.New(context.Background(), dbDSN)
		if err != nil {
			mlog.Fatalf("failed to connect to database: %s", err)
		}
	}
	defer dbpool.Close()

	model := model.SetupFromDb(dbpool)
	defer model.Close()

	// csrf Key
	viper.MustBindEnv("csrf_key")
	csrfKeyInput := viper.GetString("csrf_key")
	if csrfKeyInput == "" {
		mlog.Fatal("csrf key not supplied")
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

	mlog.Printf("listening on %s", listenHostPort)
	server := &http.Server{
		Addr:              listenHostPort,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		mlog.Print("server closed")
	} else if err != nil {
		mlog.Fatalf("error listening for server one: %s\n", err)
	}
}
