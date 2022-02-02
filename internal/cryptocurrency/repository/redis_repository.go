package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	cryptocurrency "github.com/Yangiboev/simple-server-crytocurrency-api/internal/cryptocurrency"
	"github.com/Yangiboev/simple-server-crytocurrency-api/internal/models"
)

// cryto Currency redis repository
type crytoCurrencyRedisRepo struct {
	redisClient *redis.Client
}

// NewCrytoCurrency redis repository constructor
func NewCrytoCurrencyRedisRepo(redisClient *redis.Client) cryptocurrency.RedisRepository {
	return &crytoCurrencyRedisRepo{redisClient: redisClient}
}

// Get new by id
func (ccrr crytoCurrencyRedisRepo) GetBlockByIDCtx(ctx context.Context, key string) (*models.Block, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "crytoCurrencyRedisRepo.GetBlockByID")
	defer span.Finish()

	blockBytes, err := ccrr.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	blockBase := &models.Block{}
	if err := json.Unmarshal(blockBytes, blockBase); err != nil {
		return nil, err
	}

	return blockBase, nil
}

// Cache block item
func (ccrr crytoCurrencyRedisRepo) SetBlockCtx(ctx context.Context, key string, seconds int, block *models.Block) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "crytoCurrencyRedisRepo.SetBlockByID")
	defer span.Finish()

	blockBytes, err := json.Marshal(block)
	if err != nil {
		return errors.Wrap(err, "crytoCurrencyRedisRepo.SetBlockByID.json.Marshal")
	}
	if err = ccrr.redisClient.Set(ctx, key, blockBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		return errors.Wrap(err, "crytoCurrencyRedisRepo.SetBlockByID.redisClient.Set")
	}
	return nil
}
