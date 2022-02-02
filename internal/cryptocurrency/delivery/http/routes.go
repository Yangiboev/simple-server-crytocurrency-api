package http

import (
	cryptocurrency "github.com/Yangiboev/simple-server-crytocurrency-api/internal/cryptocurrency"
	"github.com/labstack/echo/v4"
)

// Map cryptocurrency routes
func MapCryptoCurrencyRoutes(ccrrGroup *echo.Group, h cryptocurrency.Handlers) {
	ccrrGroup.GET("block/:network/:block_id", h.GetBlockByID())
	ccrrGroup.GET("transaction/:network/:transaction_id", h.GetTransactionByID())
}
