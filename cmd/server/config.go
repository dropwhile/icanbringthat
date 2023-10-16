package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path"

	"github.com/caarlos0/env/v9"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/pbkdf2"
)

type Config struct {
	LogFormat      string        `env:"LOG_FORMAT" envDefault:"json"`
	LogLevel       zerolog.Level `env:"LOG_LEVEL" envDefault:"info"`
	Production     bool          `env:"PRODUCTION" envDefault:"true"`
	ListenHostPort string        `env:"LISTEN_HOST_PORT" envDefault:"127.0.0.1:8000"`
	TemplateDir    string        `env:"TPL_DIR" envDefault:"embed"`
	StaticDir      string        `env:"STATIC_DIR" envDefault:"embed"`
	DatabaseDSN    string        `env:"DB_DSN,required"`
	RPID           string        `env:"RP_ID,required"`
	RPOrigins      []string      `env:"RP_ORIGINS,required"`
	SMTPHostname   string        `env:"SMTP_HOSTNAME,required"`
	SMTPHost       string        `env:"SMTP_HOST" envDefault:"$SMTP_HOSTNAME"`
	SMTPPort       int           `env:"SMTP_PORT,required"`
	SMTPUser       string        `env:"SMTP_USER,required"`
	SMTPPass       string        `env:"SMTP_PASS,required"`
	HMACKey        string        `env:"HMAC_KEY,required"`
	CSRFKeyBytes   []byte
	HMACKeyBytes   []byte
}

func ParseConfig() (*Config, error) {
	config := &Config{}

	//----------------//
	// parse env vars //
	//----------------//
	if err := env.Parse(config); err != nil {
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
