package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/pbkdf2"
)

type Config struct {
	LogFormat      string        `split_words:"true" default:"json"`
	LogLevel       zerolog.Level `split_words:"true" default:"info"`
	Production     bool          `split_words:"true" default:"true"`
	ListenHostPort string        `split_words:"true" default:"127.0.0.1:8000"`
	TemplateDir    string        `split_words:"true" default:"embed" envconfig:"tpl_dir"`
	StaticDir      string        `split_words:"true" default:"embed"`
	DatabaseDSN    string        `split_words:"true" required:"true" envconfig:"db_dsn"`
	RPID           string        `split_words:"true" required:"true" envconfig:"rp_id"`
	RPOrigins      []string      `split_words:"true" required:"true" envconfig:"rp_origins"`
	HMACKey        string        `split_words:"true" required:"true" envconfig:"hmac_key"`
	CSRFKeyBytes   []byte        `ignored:"true"` // do these manually
	HMACKeyBytes   []byte        `ignored:"true"` // do these manually
	SMTPHostname   string        `split_words:"true" required:"true" envconfig:"smtp_hostname"`
	SMTPHost       string        `split_words:"true" envconfig:"smtp_host"`
	SMTPPort       int           `split_words:"true" required:"true" envconfig:"smtp_port"`
	SMTPUser       string        `split_words:"true" required:"true" envconfig:"smtp_user"`
	SMTPPass       string        `split_words:"true" required:"true" envconfig:"smtp_pass"`
}

func ParseConfig() (*Config, error) {
	config := &Config{}

	//----------------//
	// parse env vars //
	//----------------//
	err := envconfig.Process("", config)
	if err != nil {
		return config, err
	}

	tplDir := path.Clean(config.TemplateDir)
	if tplDir != "embed" {
		if tplDir == "" {
			return nil, fmt.Errorf("template dir not specified")
		}
		_, err := os.Stat(tplDir)
		if err != nil && os.IsNotExist(err) {
			return nil, fmt.Errorf("template dir does not exist: %s", tplDir)
		}
	}
	config.TemplateDir = tplDir

	staticDir := path.Clean(config.StaticDir)
	if staticDir != "embed" {
		if staticDir == "" {
			return nil, fmt.Errorf("static dir not specified")
		}
		_, err := os.Stat(staticDir)
		if err != nil && os.IsNotExist(err) {
			return nil, fmt.Errorf("static dir does not exist: %s", staticDir)
		}
	}
	config.StaticDir = staticDir

	// csrf Key
	keyInput := config.HMACKey
	if keyInput == "" {
		return nil, fmt.Errorf("hmac key not supplied")
	}

	// generate hmac key based on input, using pdkdf2 to stretch/shrink
	// as needed to fit 32 byte key requirement
	config.HMACKeyBytes = pbkdf2.Key(
		[]byte(keyInput), // input
		[]byte("i4L04cpiG6JebX5brY53sBBqCyX16IwbjagbMkytmQQ="), // salt
		4096,       // iterations
		32,         // desired output size
		sha256.New, // hash to use
	)
	// continue stretching with pdkdf2 to generate a csrf key
	config.CSRFKeyBytes = pbkdf2.Key(
		config.HMACKeyBytes, // input
		[]byte("C/RWyRGBRXSCL5st5bFsPstuKQNDpRIix0vUlQ4QP80="), // salt
		4096,       // iterations
		32,         // desired output size
		sha256.New, // hash to use
	)

	return config, nil
}
