package app

import (
	"github.com/shopspring/decimal"
	"time"
)

// App errors
const (
	InvalidTransactionRequestBody   string = "An error occurred parsing from request body to a Transaction"
	InvalidTransactionResponseBody  string = "An error occurred parsing from Transaction to the response body"
	BatchNotFound                   string = "Batch not found"
	BatchHistoryNotFound            string = "Batch history not found"
	InvalidBatchResponseBody        string = "An error occurred parsing from Batch to the response body"
	InvalidBatchHistoryResponseBody string = "An error occurred parsing from Batch History to the response body"
	CannotLoadConfigFile            string = "An error occurred reading the app config file"
)

type Transaction struct {
	TransactionId string
	UserId        string
	Amount        decimal.Decimal
	CreatedAt     time.Time
}

type Batch struct {
	BatchId       string
	Transactions  []Transaction
	AccruedAmount decimal.Decimal
	IsDispatched  bool
}

// Interface to expose the service operations
type TransactionService interface {
	// Saves a transaction into a batch
	SaveTransaction(transaction *Transaction)

	// Retrieves the current batch for the userId passed by parameter
	GetBatch(userId string) (*Batch, bool)

	// Retrieves the batch history for the userId passed by parameter
	GetBatchHistory(userId string) ([]*Batch, bool)
}
