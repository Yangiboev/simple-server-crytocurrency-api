package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

const (
	DevelopmentMode string = "Development"
	ProductionMode  string = "Production"
)

const (
	LogConsoleEncoding string = "console"
)

// App config struct
type Config struct {
	Server  Server
	Redis   Redis
	Metrics Metrics
	Logger  Logger
	Jaeger  Jaeger
}

// Server config struct
type Server struct {
	AppVersion        string
	Port              string
	Mode              string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	SSL               bool
	CtxDefaultTimeout time.Duration
	Debug             bool
}

// Logger config
type Logger struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

// Redis config
type Redis struct {
	RedisAddr      string
	RedisPassword  string
	RedisDB        string
	RedisDefaultdb string
	MinIdleConns   int
	PoolSize       int
	PoolTimeout    int
	Password       string
	DB             int
}

// Metrics config
type Metrics struct {
	URL         string
	ServiceName string
}

type Jaeger struct {
	Host        string
	ServiceName string
	LogSpans    bool
}

func Load(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

func Parse(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}
