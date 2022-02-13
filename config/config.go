package config

import (
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

func New(filename string) (*Config, error) {
	return parseFile(filename)
}

func getViper(filename string) (*viper.Viper, error) {
	var (
		v = viper.New()
	)

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return v, nil
}

func parseFile(filename string) (*Config, error) {
	var (
		c Config
	)

	v, err := getViper(filename)
	if err != nil {
		return nil, err
	}

	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}

	return &c, nil
}
