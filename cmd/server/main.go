package main

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/cactus/mlog"
	"github.com/dropwhile/icbt/internal/api"
	"github.com/dropwhile/icbt/internal/model"
	"github.com/dropwhile/icbt/resources"
	"github.com/spf13/viper"
)

// ServerVersion holds the server version string
var ServerVersion = "no-version"

func main() {
	// parse env vars

	// debug logging or not
	viper.BindEnv("debug")
	viper.SetDefault("debug", "false")
	logDebug := viper.GetBool("debug")
	if logDebug {
		mlog.SetFlags(mlog.Flags() | mlog.Ldebug)
		mlog.Debug("debug logging enabled")
	}

	// listen address/port
	viper.BindEnv("bind_address")
	viper.SetDefault("bind_address", "127.0.0.1")
	viper.BindEnv("bind_port")
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
	viper.BindEnv("tpl_dir")
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
	viper.BindEnv("static_dir")
	viper.SetDefault("static_dir", "embed")
	staticDir := path.Clean(viper.GetString("static_dir"))
	if staticDir == "embed" {
		mlog.Debugf("using embedded static")
	} else {
		mlog.Debugf("static dir set to: %s", staticDir)
	}
	staticFS := resources.NewStaticFS(staticDir)

	// database
	viper.BindEnv("db_dsn")
	dbDSN := viper.GetString("db_dsn")
	if dbDSN == "" {
		mlog.Fatal("database connection info not supplied")
	}
	db, err := model.NewDatabase(dbDSN)
	if err != nil {
		mlog.Fatal("error connecting to database")
	}
	defer db.Close()

	// routing/handlers
	r := api.New(db, templates)
	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.FS(staticFS))))

	mlog.Printf("listening on %s", listenHostPort)
	http.ListenAndServe(listenHostPort, r)
}
