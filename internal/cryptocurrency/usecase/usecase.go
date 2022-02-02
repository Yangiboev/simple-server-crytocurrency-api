package usecase

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"

	"github.com/Yangiboev/simple-server-crytocurrency-api/config"
	cryptocurrency "github.com/Yangiboev/simple-server-crytocurrency-api/internal/cryptocurrency"
	"github.com/Yangiboev/simple-server-crytocurrency-api/internal/models"
	"github.com/Yangiboev/simple-server-crytocurrency-api/pkg/logger"
)

const (
	cacheDuration = 60
)

// cryptoCurrency UseCase
type cryptoCurrencyUC struct {
	cfg       *config.Config
	ccrrRepo  cryptocurrency.Repository
	redisRepo cryptocurrency.RedisRepository
	logger    logger.Logger
}

// cryptoCurrency UseCase constructor
func NewCryptoCurrencyUseCase(cfg *config.Config, ccrrRepo cryptocurrency.Repository, redisRepo cryptocurrency.RedisRepository, logger logger.Logger) cryptocurrency.UseCase {
	return &cryptoCurrencyUC{cfg: cfg, ccrrRepo: ccrrRepo, redisRepo: redisRepo, logger: logger}
}

// Get CrytoCurrency by id
func (u *cryptoCurrencyUC) GetBlockByID(ctx context.Context, network, blockID string) (*models.Block, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cryptoCurrencyUC.GetBlockByID")
	defer span.Finish()

	blockBase, err := u.redisRepo.GetBlockByIDCtx(ctx, network+blockID)
	fmt.Println(err != redis.Nil)
	if err != redis.Nil && err != nil {
		u.logger.Errorf("cryptoCurrencyUC.GetBlockByID.GetBlockByIDCtx: %v", err)
	}
	fmt.Println(blockBase)
	if blockBase != nil {
		return blockBase, nil
	}

	b, err := u.ccrrRepo.GetBlockByID(ctx, network, blockID)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.SetBlockCtx(ctx, network+blockID, cacheDuration, b); err != nil {
		u.logger.Errorf("cryptoCurrencyUC.GetBlockByID.SetBlockCtx: %s", err)
	}
	return b, nil
}

// Get CrytoCurrency by id
func (u *cryptoCurrencyUC) GetTransactionByID(ctx context.Context, network, trasactionID string) (*models.Transaction, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cryptoCurrencyUC.GetTransactionByID")
	defer span.Finish()

	// blockBase, err := u.redisRepo.GetTransactionByIDCtx(ctx, blockID)
	// if err.Error() != repository.ErrCacheMiss.Error() && err != nil {
	// 	u.logger.Errorf("cryptoCurrencyUC.GetTransactionByID.GetTransactionByIDCtx: %v", err)
	// }
	// if blockBase != nil {
	// 	return blockBase, nil
	// }
	// fmt.Println("asdasdasd")

	t, err := u.ccrrRepo.GetTransactionByID(ctx, network, trasactionID)
	if err != nil {
		return nil, err
	}

	// if err = u.redisRepo.SetCrytoCurrencyCtx(ctx, blockID, cacheDuration, n); err != nil {
	// 	u.logger.Errorf("cryptoCurrencyUC.GetTransactionByID.SetCrytoCurrencyCtx: %s", err)
	// }

	return t, nil
}
