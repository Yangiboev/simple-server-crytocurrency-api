package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	cryptocurrency "github.com/Yangiboev/simple-server-crytocurrency-api/internal/cryptocurrency"
	"github.com/Yangiboev/simple-server-crytocurrency-api/internal/models"
	"github.com/Yangiboev/simple-server-crytocurrency-api/pkg/httpErrors"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

var (
	layoutDateTime = "02-01-2006 15:04"
)

// block Repository
type crytocurrencyRepo struct {
}

// block repository constructor
func NewCrytoCurrencyRepository() cryptocurrency.Repository {
	return &crytocurrencyRepo{}
}

// Get single block by id
func (ccrr *crytocurrencyRepo) GetBlockByID(ctx context.Context, network, blockID string) (*models.Block, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "crytocurrencyRepo.GetBlockByID")
	defer span.Finish()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://sochain.com/api/v2/get_block/"+network+"/"+blockID, nil)
	if err != nil {
		return nil, errors.Wrap(err, "crytocurrencyRepo.GetBlockByID.NewRequestWithContext")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "crytocurrencyRepo.GetBlockByID.ReadAll")

	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return nil, httpErrors.NotFound
	}
	//We Read the response body on the line below.
	var blockResponse = &models.BResponse{}
	if err := json.NewDecoder(res.Body).Decode(blockResponse); err != nil {
		return nil, errors.Wrap(err, "crytocurrencyRepo.GetBlockByID.NewDecoder")
	}
	dateTime := time.Unix(blockResponse.Data.DateTime, 0)
	transactions, err := ccrr.GetBlockTenTransactions(ctx, network, blockResponse.Data.Transactions[:11])
	if err != nil {
		return nil, errors.Wrap(err, "crytocurrencyRepo.GetBlockByID.GetBlockTenTransactions")
	}
	return &models.Block{
		NetworkCode:       blockResponse.Data.NetworkCode,
		BlockHash:         blockResponse.Data.BlockHash,
		BlockNumber:       blockResponse.Data.BlockNumber,
		DateTime:          dateTime.Format(layoutDateTime),
		PreviousBlockHash: blockResponse.Data.PreviousBlockHash,
		NextBlockHash:     blockResponse.Data.NextBlockHash,
		Size:              blockResponse.Data.Size,
		Transactions:      transactions,
	}, nil
}

// Get single transaction by id
func (ccrr *crytocurrencyRepo) GetTransactionByID(ctx context.Context, network, trasactionID string) (*models.Transaction, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "crytocurrencyRepo.GetTransactionByID")
	defer span.Finish()
	fmt.Println(network)
	fmt.Println(trasactionID)
	req, err := http.NewRequestWithContext(ctx, "GET", "https://sochain.com/api/v2/tx/"+network+"/"+trasactionID, nil)
	if err != nil {
		return nil, errors.Wrap(err, "crytocurrencyRepo.GetTransactionByID.NewRequestWithContext")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "crytocurrencyRepo.GetTransactionByID.ReadAll")
	}
	defer res.Body.Close()
	//We Read the response body on the line below.
	var tResponse = &models.TResponse{}
	if err := json.NewDecoder(res.Body).Decode(tResponse); err != nil {
		return nil, errors.Wrap(err, "crytocurrencyRepo.GetTransactionByID.NewDecoder")
	}
	fmt.Println(tResponse.Data)
	fmt.Println(tResponse.Data)
	fee, err := strconv.ParseFloat(tResponse.Data.Fee, 64)
	if err != nil {
		return nil, errors.Wrap(err, "crytocurrencyRepo.GetTransactionByID.ParseFloat.FEE")
	}
	sentValue, err := strconv.ParseFloat(tResponse.Data.SentValue, 64)
	if err != nil {
		return nil, errors.Wrap(err, "crytocurrencyRepo.GetTransactionByID.ParseFloat.SentValue")
	}
	dateTime := time.Unix(tResponse.Data.DateTime, 0)

	return &models.Transaction{
		TransactionID: tResponse.Data.TransactionID,
		DateTime:      dateTime.Format(layoutDateTime),
		Fee:           fee,
		SentValue:     sentValue,
	}, nil
}

// Get single block by id
func (ccrr *crytocurrencyRepo) GetBlockTenTransactions(ctx context.Context, network string, transactionIDs []string) ([]*models.Transaction, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "crytocurrencyRepo.GetCrytoCurrencyByID")
	defer span.Finish()
	var (
		response    []*models.Transaction
		results     = make(chan *models.Transaction, len(transactionIDs))
		errGroup, _ = errgroup.WithContext(ctx)
	)
	for _, transactionID := range transactionIDs {
		tID := transactionID
		errGroup.Go(func() error {
			transaction, err := ccrr.GetTransactionByID(ctx, network, tID)
			results <- transaction
			return err
		})
	}

	if err := errGroup.Wait(); err != nil {
		return nil, errors.Wrap(err, "crytocurrencyRepo.GetBlockTenTransactions.ErrorGroup")
	}
	close(results)
	for result := range results {
		response = append(response, result)
	}
	return response, nil
}
