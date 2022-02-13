package server

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/Yangiboev/simple-server-crytocurrency-api/internal/docs"

	cryptoCurrencyHttp "github.com/Yangiboev/simple-server-crytocurrency-api/internal/cryptocurrency/delivery/http"
	crytoCurrencyRepository "github.com/Yangiboev/simple-server-crytocurrency-api/internal/cryptocurrency/repository"
	cryptoCurrencyUseCase "github.com/Yangiboev/simple-server-crytocurrency-api/internal/cryptocurrency/usecase"
	"github.com/Yangiboev/simple-server-crytocurrency-api/pkg/utils"
)

// Map Server Handlers
func (s *Server) MapHandlers() error {
	// metrics, err := metric.CreateMetrics(s.cfg.Metrics.URL, s.cfg.Metrics.ServiceName)
	// if err != nil {
	// 	s.logger.Errorf("CreateMetrics Error: %s", err)
	// }
	// s.logger.Info(
	// 	"Metrics available URL: %s, ServiceName: %s",
	// 	s.cfg.Metrics.URL,
	// 	s.cfg.Metrics.ServiceName,
	// )

	// Init repositories
	ccrrRepo := crytoCurrencyRepository.NewCrytoCurrencyRepository()
	ccrrRedisRepo := crytoCurrencyRepository.NewCrytoCurrencyRedisRepo(s.redisClient)

	// Init useCases
	ccrrUC := cryptoCurrencyUseCase.NewCryptoCurrencyUseCase(s.cfg, ccrrRepo, ccrrRedisRepo, s.logger)

	// Init handlers
	ccrrHandlers := cryptoCurrencyHttp.NewCryptoCurrencyHandlers(s.cfg, ccrrUC, s.logger)

	docs.SwaggerInfo.Title = "Go REST API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Description = "Go API for getting details block and transaction."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/v1"
	s.echo.GET("/swagger/*", echoSwagger.WrapHandler)

	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID},
	}))
	s.echo.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	s.echo.Use(middleware.RequestID())

	s.echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))
	s.echo.Use(middleware.Secure())
	s.echo.Use(middleware.BodyLimit("2M"))

	v1 := s.echo.Group("/v1")

	health := v1.Group("/health")
	ccrrGroup := v1.Group("/")

	cryptoCurrencyHttp.MapCryptoCurrencyRoutes(ccrrGroup, ccrrHandlers)

	health.GET("", func(c echo.Context) error {
		s.logger.Infof("Health check RequestID: %s", utils.GetRequestID(c))
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	return nil
}
