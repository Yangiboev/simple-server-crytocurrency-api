package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"

	"github.com/Yangiboev/simple-server-crytocurrency-api/config"
	cryptocurrency "github.com/Yangiboev/simple-server-crytocurrency-api/internal/cryptocurrency"
	"github.com/Yangiboev/simple-server-crytocurrency-api/pkg/httpErrors"
	"github.com/Yangiboev/simple-server-crytocurrency-api/pkg/logger"
	"github.com/Yangiboev/simple-server-crytocurrency-api/pkg/utils"
	"github.com/Yangiboev/simple-server-crytocurrency-api/pkg/validator"
)

// cryptoCurrency handlers
type cryptoCurrencyHandlers struct {
	cfg              *config.Config
	cryptoCurrencyUC cryptocurrency.UseCase
	logger           logger.Logger
}

// NewCryptoCurrencyHandlers handlers constructor
func NewCryptoCurrencyHandlers(cfg *config.Config, cryptoCurrencyUC cryptocurrency.UseCase, logger logger.Logger) cryptocurrency.Handlers {
	return &cryptoCurrencyHandlers{cfg: cfg, cryptoCurrencyUC: cryptoCurrencyUC, logger: logger}
}

// @Summary Get block by block hash
// @Description Get block by block hash using handler
// @Tags cryptocurrency
// @Accept json
// @Produce json
// @Param network path string true "we only support the following cryptocurrency network codes: BTC, LTC and DOGE."
// @Param block_id path string true "block_id"
// @Success 200 {object} models.Block
// @Router /block/{network}/{block_id} [get]
func (h cryptoCurrencyHandlers) GetBlockByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "cryptoCurrencyHandlers.GetByID")
		defer span.Finish()

		blockID := c.Param("block_id")
		network := c.Param("network")
		if validator.IsValidNetwork(network) {
			utils.LogResponseError(c, h.logger, httpErrors.BadQueryParams)
			return c.JSON(httpErrors.ErrorResponse(httpErrors.BadQueryParams))
		}
		if len(blockID) != 64 {
			utils.LogResponseError(c, h.logger, httpErrors.BadQueryParams)
			return c.JSON(httpErrors.ErrorResponse(httpErrors.BadQueryParams))
		}
		resp, err := h.cryptoCurrencyUC.GetBlockByID(ctx, network, blockID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, resp)
	}
}

// @Summary Get transaction by transaction hash
// @Description Get transaction by transaction hash using handler
// @Tags cryptocurrency
// @Accept json
// @Produce json
// @Param network path string true "we only support the following cryptocurrency network codes: BTC, LTC and DOGE."
// @Param transaction_id path string true "transaction_id"
// @Success 200 {object} models.Transaction
// @Router /transaction/{network}/{transaction_id} [get]
func (h cryptoCurrencyHandlers) GetTransactionByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(utils.GetRequestCtx(c), "cryptoCurrencyHandlers.GetByID")
		defer span.Finish()

		transactionID := c.Param("transaction_id")
		network := c.Param("network")
		if validator.IsValidNetwork(network) {
			utils.LogResponseError(c, h.logger, httpErrors.BadQueryParams)
			return c.JSON(httpErrors.ErrorResponse(httpErrors.BadQueryParams))
		}
		if len(transactionID) != 64 {
			utils.LogResponseError(c, h.logger, httpErrors.BadQueryParams)
			return c.JSON(httpErrors.ErrorResponse(httpErrors.BadQueryParams))
		}
		resp, err := h.cryptoCurrencyUC.GetTransactionByID(ctx, network, transactionID)
		fmt.Println(resp)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, resp)
	}
}
