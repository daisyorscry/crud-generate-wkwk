package config

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

var (
	c    *Config
	once sync.Once
	mu   sync.Mutex
)

type (
	Database struct {
		Host         string
		Port         string
		Name         string
		User         string
		Password     string
		MaxConn      int
		MinConn      int
		ConnLifetime time.Duration
		IdleTimeout  time.Duration
		ConnTimeout  time.Duration
	}
	JWT struct {
		SecretKey        string
		SecretKeyRefresh string
	}
	AppSettings struct {
		Host        string
		Port        string
		CORSOrigins []string
	}
	Mail struct {
		From     string
		Host     string
		Port     int
		Username string
		Password string
	}
	Config struct {
		Database    Database
		JWT         JWT
		AppSettings AppSettings
		Mail        Mail
	}
)

func New() *Config {
	if err := gotenv.Load(".env"); err != nil {
		log.Printf("Error loading .env file, fallback to system environment. Error: %v", err)
	}

	viper.AutomaticEnv()
	conf := &Config{}

	// App Settings
	conf.AppSettings.Host = viper.GetString("DAISY_LAONDRY_APP_HOST")
	conf.AppSettings.Port = viper.GetString("DAISY_LAONDRY_APP_PORT")
	corsRaw := viper.GetString("DAISY_LAONDRY_APP_CORS")
	conf.AppSettings.CORSOrigins = strings.Split(corsRaw, ",")

	// Database
	conf.Database.Host = viper.GetString("DAISY_LAONDRY_DATABASE_HOST")
	conf.Database.Port = viper.GetString("DAISY_LAONDRY_DATABASE_PORT")
	conf.Database.Name = viper.GetString("DAISY_LAONDRY_DATABASE_NAME")
	conf.Database.User = viper.GetString("DAISY_LAONDRY_DATABASE_USER")
	conf.Database.Password = viper.GetString("DAISY_LAONDRY_DATABASE_PASSWORD")
	conf.Database.MaxConn = viper.GetInt("DAISY_LAONDRY_DATABASE_MAX_CONN")
	conf.Database.MinConn = viper.GetInt("DAISY_LAONDRY_DATABASE_MIN_CONN")
	conf.Database.ConnLifetime = viper.GetDuration("DAISY_LAONDRY_DATABASE_CONN_LIFETIME")
	conf.Database.IdleTimeout = viper.GetDuration("DAISY_LAONDRY_DATABASE_CONN_IDLE_TIMEOUT")
	conf.Database.ConnTimeout = viper.GetDuration("DAISY_LAONDRY_DATABASE_CONN_TIMEOUT")

	if conf.Database.MaxConn == 0 {
		conf.Database.MaxConn = 10
	}
	if conf.Database.MinConn == 0 {
		conf.Database.MinConn = 2
	}
	if conf.Database.ConnLifetime == 0 {
		conf.Database.ConnLifetime = 30 * time.Minute
	}
	if conf.Database.IdleTimeout == 0 {
		conf.Database.IdleTimeout = 5 * time.Minute
	}
	if conf.Database.ConnTimeout == 0 {
		conf.Database.ConnTimeout = 15 * time.Second
	}

	// Mail
	conf.Mail.From = viper.GetString("DAISY_LAONDRY_MAIL_FROM")
	conf.Mail.Host = viper.GetString("DAISY_LAONDRY_MAIL_HOST")
	conf.Mail.Port = viper.GetInt("DAISY_LAONDRY_MAIL_PORT")
	conf.Mail.Username = viper.GetString("DAISY_LAONDRY_MAIL_USERNAME")
	conf.Mail.Password = viper.GetString("DAISY_LAONDRY_MAIL_PASSWORD")

	// JWT
	conf.JWT.SecretKey = viper.GetString("DAISY_LAONDRY_JWT_KEY")
	conf.JWT.SecretKeyRefresh = viper.GetString("DAISY_LAONDRY_JWT_KEY_REFRESH")

	c = conf
	return conf
}

func Get() *Config {
	mu.Lock()
	defer mu.Unlock()

	if c == nil {
		once.Do(func() {
			c = New()
		})
	}
	return c
}
