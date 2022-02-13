package main

import (
	"flag"
	"log"

	"github.com/Yangiboev/simple-server-crytocurrency-api/config"
	"github.com/Yangiboev/simple-server-crytocurrency-api/internal/server"
	"github.com/Yangiboev/simple-server-crytocurrency-api/pkg/logger"
	redis "github.com/Yangiboev/simple-server-crytocurrency-api/pkg/redis"
	"github.com/opentracing/opentracing-go"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
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

	flag.StringVar(&environment, "environment", environment, `Set -environment to use load configurations for differenet environments and different log levels. Default is "local"`)
	flag.StringVar(&configPath, "path", configPath, `Set -path to use load configurations from given path which was given without file extension. Default is "./config/local"`)

	flag.Parse()

	log.Printf("process environment \"%s\"\n", environment)
	log.Printf("process configuration path \"%s\"\n", configPath)

	log.Println("Starting api server")

	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatalf("can not load configuration file: %v", err)
	}

	logger := logger.New(cfg.Server.Mode, &cfg.Logger)
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Fatalf("could not sync logger: %v\n", err)
		}
	}()

	logger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

	redisClient := redis.NewClient(&cfg.Redis)
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Fatalf("could not close redis client: %v\n", err)
		}
	}()
	logger.Info("Redis connected")

	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: cfg.Jaeger.ServiceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           cfg.Jaeger.LogSpans,
			LocalAgentHostPort: cfg.Jaeger.Host,
		},
	}

	tracer, closer, err := jaegerCfgInstance.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)
	if err != nil {
		log.Fatal("cannot create tracer", err)
	}
	logger.Info("Jaeger connected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	logger.Info("Opentracing connected")

	s := server.NewServer(cfg, redisClient, logger)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
