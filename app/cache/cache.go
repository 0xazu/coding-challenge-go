package cache

import (
	"coding-challenge-go/app"
	"coding-challenge-go/app/config"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"sync"
)

type TransactionService struct {
	mu           sync.RWMutex
	config       *config.Config
	batch        map[string]*app.Batch
	batchHistory map[string][]*app.Batch
}

func NewTransactionService(config config.Config) *TransactionService {
	var t TransactionService

	t.config = &config
	t.batch = make(map[string]*app.Batch)
	t.batchHistory = make(map[string][]*app.Batch)

	return &t
}

func (transactionService *TransactionService) GetBatch(userId string) (*app.Batch, bool) {
	transactionService.mu.RLock()
	defer transactionService.mu.RUnlock()

	batch, exists := transactionService.batch[userId]
	return batch, exists
}

func (transactionService *TransactionService) GetBatchHistory(userId string) ([]*app.Batch, bool) {
	transactionService.mu.RLock()
	defer transactionService.mu.RUnlock()

	batchHistory, exists := transactionService.batchHistory[userId]
	return batchHistory, exists
}

func (transactionService *TransactionService) SaveTransaction(transaction *app.Transaction) {
	transactionService.mu.Lock()
	defer transactionService.mu.Unlock()

	batch, exists := transactionService.batch[transaction.UserId]

	if exists {
		batch.Transactions = append(batch.Transactions, *transaction)
		batch.AccruedAmount = batch.AccruedAmount.Add(transaction.Amount)
	} else {
		var batch = &app.Batch{
			BatchId:       uuid.NewString(),
			Transactions:  []app.Transaction{*transaction},
			AccruedAmount: transaction.Amount,
			IsDispatched:  false,
		}

		transactionService.batch[transaction.UserId] = batch
	}

	transactionService.dispatchBatch(transaction.UserId)
}

func (transactionService *TransactionService) dispatchBatch(userId string) {
	var batch = transactionService.batch[userId]

	if batch.AccruedAmount.GreaterThan(decimal.NewFromInt32(transactionService.config.BatchThreshold)) {
		batch.IsDispatched = true
		transactionService.batchHistory[userId] = append(transactionService.batchHistory[userId], batch)

		delete(transactionService.batch, userId)
	}
}
