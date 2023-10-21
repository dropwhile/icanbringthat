package util

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/caarlos0/env/v9"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/argon2"
)

type Config struct {
	LogFormat    string        `env:"LOG_FORMAT" envDefault:"json"`
	LogLevel     zerolog.Level `env:"LOG_LEVEL" envDefault:"info"`
	Production   bool          `env:"PRODUCTION" envDefault:"true"`
	Listen       string        `env:"LISTEN" envDefault:"127.0.0.1:8000"`
	TemplateDir  string        `env:"TPL_DIR" envDefault:"embed"`
	StaticDir    string        `env:"STATIC_DIR" envDefault:"embed"`
	DatabaseDSN  string        `env:"DB_DSN,required"`
	BaseURL      string        `env:"BASE_URL,required"`
	SMTPHostname string        `env:"SMTP_HOSTNAME,required"`
	SMTPHost     string        `env:"SMTP_HOST,expand" envDefault:"$SMTP_HOSTNAME"`
	SMTPPort     int           `env:"SMTP_PORT,required"`
	SMTPUser     string        `env:"SMTP_USER,required"`
	SMTPPass     string        `env:"SMTP_PASS,required"`
	MailFrom     string        `env:"MAIL_FROM,required"`
	HMACKey      string        `env:"HMAC_KEY,required"`
	CSRFKeyBytes []byte
	HMACKeyBytes []byte
}

func ParseConfig() (*Config, error) {
	config := &Config{}

	//----------------//
	// parse env vars //
	//----------------//
	if err := env.Parse(config); err != nil {
		return config, err
	}

	if !strings.Contains(config.Listen, ":") {
		if strings.Contains(config.Listen, ".") {
			config.Listen = strings.Join([]string{config.Listen, "8000"}, ":")
		} else {
			config.Listen = strings.Join([]string{"127.0.0.1", config.Listen}, ":")
		}
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
	config.HMACKeyBytes = argon2.IDKey(
		[]byte(keyInput), // input
		[]byte("i4L04cpiG6JebX5brY53sBBqCyX16IwbjagbMkytmQQ="), // salt
		1,       // time
		64*1024, // memory
		4,       // threads
		32,      // desired keylength
	)
	// continue stretching with pdkdf2 to generate a csrf key
	config.CSRFKeyBytes = argon2.IDKey(
		config.HMACKeyBytes, // input
		[]byte("C/RWyRGBRXSCL5st5bFsPstuKQNDpRIix0vUlQ4QP80="), // salt
		1,       // time
		64*1024, // memory
		4,       // threads
		32,      // desired keylength
	)

	return config, nil
}