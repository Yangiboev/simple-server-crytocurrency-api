//go:generate mockgen -source usecase.go -destination mock/usecase_mock.go -package mock
package cryptocurrency

import (
	"context"

	"github.com/Yangiboev/simple-server-crytocurrency-api/internal/models"
)

// CryptoCurrency use case
type UseCase interface {
	GetBlockByID(ctx context.Context, network, blockID string) (*models.Block, error)
	GetTransactionByID(ctx context.Context, network, blockID string) (*models.Transaction, error)
}
