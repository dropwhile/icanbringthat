package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"golang.org/x/crypto/pbkdf2"
)

type Config struct {
	LogFormat      string
	LogLevel       zerolog.Level
	Production     bool
	ListenHostPort string
	TemplateDir    string
	StaticDir      string
	DatabaseDSN    string
	CSRFKeyBytes   []byte
	HMACKeyBytes   []byte
	SMTPHostname   string
	SMTPHost       string
	SMTPPort       int
	SMTPUser       string
	SMTPPass       string
}

func ParseConfig() (*Config, error) {
	config := &Config{}

	//----------------//
	// parse env vars //
	//----------------//

	// log format
	viper.MustBindEnv("LOG_FORMAT")
	viper.SetDefault("LOG_FORMAT", "json")
	config.LogFormat = viper.GetString("LOG_FORMAT")

	// debug logging or not
	viper.MustBindEnv("LOG_LEVEL")
	viper.SetDefault("LOG_LEVEL", "info")
	logLevel := viper.GetString("LOG_LEVEL")
	switch strings.ToLower(logLevel) {
	case "debug":
		config.LogLevel = zerolog.DebugLevel
	case "trace":
		config.LogLevel = zerolog.TraceLevel
	default:
		config.LogLevel = zerolog.InfoLevel
	}

	// prod mode (secure cookies) or not
	viper.MustBindEnv("PRODUCTION")
	viper.SetDefault("PRODUCTION", "true")
	config.Production = viper.GetBool("PRODUCTION")

	// listen address/port
	viper.MustBindEnv("BIND_ADDRESS")
	viper.SetDefault("BIND_ADDRESS", "127.0.0.1")
	viper.MustBindEnv("BIND_PORT")
	viper.SetDefault("BIND_PORT", "8000")
	listenAddr := viper.GetString("BIND_ADDRESS")
	if listenAddr == "" {
		return nil, fmt.Errorf("listen address not specified")
	}
	listenPort := viper.GetInt("BIND_PORT")
	if listenPort == 0 {
		return nil, fmt.Errorf("listen address not specified")
	}
	config.ListenHostPort = fmt.Sprintf("%s:%d", listenAddr, listenPort)

	// load templates
	viper.MustBindEnv("TPL_DIR")
	viper.SetDefault("TPL_DIR", "embed")
	tplDir := path.Clean(viper.GetString("TPL_DIR"))
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

	// static resources
	viper.MustBindEnv("STATIC_DIR")
	viper.SetDefault("STATIC_DIR", "embed")
	staticDir := path.Clean(viper.GetString("STATIC_DIR"))
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

	// database
	viper.MustBindEnv("DB_DSN")
	dbDSN := viper.GetString("DB_DSN")
	if dbDSN == "" {
		return nil, fmt.Errorf("database connection info not supplied")
	}
	config.DatabaseDSN = dbDSN

	// csrf Key
	viper.MustBindEnv("HMAC_KEY")
	keyInput := viper.GetString("HMAC_KEY")
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

	// smtp stuff
	viper.MustBindEnv("SMTP_HOSTNAME")
	smtpHostname := viper.GetString("SMTP_HOSTNAME")
	if smtpHostname == "" {
		return nil, fmt.Errorf("smtp mail hostname not set")
	}
	config.SMTPHostname = smtpHostname

	viper.MustBindEnv("SMTP_HOST")
	smtpHost := viper.GetString("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = smtpHostname
	}
	config.SMTPHost = smtpHost

	viper.MustBindEnv("SMTP_PORT")
	smtpPort := viper.GetInt("SMTP_PORT")
	if smtpPort == 0 {
		return nil, fmt.Errorf("smtp mail port not set")
	}
	config.SMTPPort = smtpPort

	viper.MustBindEnv("SMTP_USER")
	smtpUser := viper.GetString("SMTP_USER")
	if smtpUser == "" {
		return nil, fmt.Errorf("smtp mail username not set")
	}
	config.SMTPUser = smtpUser

	viper.MustBindEnv("SMTP_PASS")
	smtpPass := viper.GetString("SMTP_PASS")
	if smtpPass == "" {
		return nil, fmt.Errorf("smtp mail password not set")
	}
	config.SMTPPass = smtpPass

	return config, nil
}
