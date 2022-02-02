package cryptocurrency

import (
	"context"

	"github.com/Yangiboev/simple-server-crytocurrency-api/internal/models"
)

// CryptoCurrency Repository
type Repository interface {
	GetBlockByID(ctx context.Context, network, blockID string) (*models.Block, error)
	GetTransactionByID(ctx context.Context, network, blockID string) (*models.Transaction, error)
}
