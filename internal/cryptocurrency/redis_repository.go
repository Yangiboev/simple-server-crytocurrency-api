package cryptocurrency

import (
	"context"

	"github.com/Yangiboev/simple-server-crytocurrency-api/internal/models"
)

// Block redis repository
type RedisRepository interface {
	GetBlockByIDCtx(ctx context.Context, key string) (*models.Block, error)
	SetBlockCtx(ctx context.Context, key string, seconds int, block *models.Block) error
}
