package cryptocurrency

import "github.com/labstack/echo/v4"

// CrytoCurrency HTTP Handlers interface
type Handlers interface {
	GetBlockByID() echo.HandlerFunc
	GetTransactionByID() echo.HandlerFunc
}
