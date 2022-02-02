package models

// Transaction model
type Transaction struct {
	TransactionID string  `json:"transaction_id"`
	DateTime      string  `json:"date_time" validate:"required"`
	Fee           float64 `json:"fee" validate:"required"`
	SentValue     float64 `json:"sent_value" validate:"required"`
}

// Block model
type Block struct {
	NetworkCode       string         `json:"network_code" validate:"required"`
	BlockHash         string         `json:"block_hash" validate:"required"`
	BlockNumber       int64          `json:"block_number" validate:"required"`
	DateTime          string         `json:"date_time" validate:"required"`
	PreviousBlockHash string         `json:"previous_block_hash" validate:"required"`
	NextBlockHash     string         `json:"next_block_hash" validate:"required"`
	Size              int64          `json:"size" validate:"required"`
	Transactions      []*Transaction `json:"transactions" validate:"required"`
}

// TransactionResponse model
type TransactionResponse struct {
	TransactionID string `json:"txid"`
	DateTime      int64  `json:"time"`
	Fee           string `json:"fee"`
	SentValue     string `json:"sent_value"`
}

// BlockResponse model
type BlockResponse struct {
	NetworkCode       string   `json:"network"`
	BlockHash         string   `json:"blockhash"`
	BlockNumber       int64    `json:"block_no"`
	DateTime          int64    `json:"time"`
	PreviousBlockHash string   `json:"previous_blockhash"`
	NextBlockHash     string   `json:"next_blockhash"`
	Size              int64    `json:"size"`
	Transactions      []string `json:"txs"`
}
type BResponse struct {
	Status string         `json:"status"`
	Data   *BlockResponse `json:"data"`
}
type TResponse struct {
	Status string               `json:"status"`
	Data   *TransactionResponse `json:"data"`
}
