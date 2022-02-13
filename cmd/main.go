package main

import (
	"flag"
	"log"

	"github.com/Yangiboev/simple-server-crytocurrency-api/config"
	"github.com/Yangiboev/simple-server-crytocurrency-api/internal/server"
	"github.com/Yangiboev/simple-server-crytocurrency-api/pkg/jaeger"
	"github.com/Yangiboev/simple-server-crytocurrency-api/pkg/logger"
	redis "github.com/Yangiboev/simple-server-crytocurrency-api/pkg/redis"
	"github.com/Yangiboev/simple-server-crytocurrency-api/pkg/validator"
	"github.com/labstack/echo/v4"
)

// @title Go REST API
// @version 1.0
// @description Golang REST API
// @contact.name Dilmurod Yangiboev
// @contact.url https://github.com/Yangiboev/simple-server-crytocurrency-api
// @contact.email dilmurod.yangiboev
// @BasePath /api/v1
func main() {
	var (
		environment string = "local"
		configPath  string = "./config/local"
	)

	// flag sets for getting arguments from command line
	//
	// environment used for setting application environment and logging
	flag.StringVar(&environment, "environment", environment, `Set -environment to use load configurations for differenet environments and different log levels. Default is "local"`)
	//
	// configPath used for setting application configuration file path
	flag.StringVar(&configPath, "path", configPath, `Set -path to use load configurations from given path which was given without file extension. Default is "./config/local"`)
	// parse flags
	flag.Parse()

	log.Printf("process environment \"%s\"\n", environment)
	log.Printf("process configuration path \"%s\"\n", configPath)

	log.Println("Starting api server")

	// load config from configPath
	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatalf("could not load configuration file: %v", err)
	}

	// get logger instance
	logger := logger.New(cfg.Server.Mode, &cfg.Logger)
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Errorf("could not sync logger: %v\n", err)
		}
	}()

	logger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

	// connect to redis
	redisClient := redis.NewClient(&cfg.Redis)
	defer func() {
		if err := redisClient.Close(); err != nil {
			logger.Errorf("could not close redis client: %v\n", err)
		}
	}()
	logger.Info("Redis connected")

	// connect to jaeger and get closer function
	jaegerCloser, err := jaeger.New(&cfg.Jaeger)
	if err != nil {
		logger.Fatalf("could not connect to jaeger: %v\n", err)
	}
	defer func() {
		if err := jaegerCloser(); err != nil {
			logger.Errorf("could not close jaeger client: %v\n", err)
		}
	}()
	logger.Info("Jaeger connected")

	// create validator singleton instance
	validator.New()

	// create echo router instance
	echo := echo.New()

	// create new server instance and start server
	s := server.New(cfg, echo, redisClient, logger)
	if err := s.Run(); err != nil {
		logger.Fatalf("could not run server: %v\n", err)
	}
}
