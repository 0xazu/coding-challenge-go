package cache

import (
	"coding-challenge-go/app"
	"coding-challenge-go/app/config"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"testing"
	"time"
)

func TestGetBatch(t *testing.T) {
	cfg := config.Config{BatchThreshold: 100}
	transactionService := NewTransactionService(cfg)
	userId := uuid.NewString()
	expectedBatch := generateTestBatch(userId)
	transactionService.batch[userId] = expectedBatch

	batch, exists := transactionService.GetBatch(userId)

	if !exists {
		t.Errorf("Expected the batch %s to exist for the userId %s", expectedBatch.BatchId, userId)
	}

	if !cmp.Equal(batch, expectedBatch) {
		t.Errorf("Expected the returned batch %s to be equal to expected batch %s", batch.BatchId, expectedBatch.BatchId)
	}
}

func TestGetBatchDoesntExist(t *testing.T) {
	cfg := config.Config{BatchThreshold: 100}
	transactionService := NewTransactionService(cfg)
	userId := uuid.NewString()

	batch, exists := transactionService.GetBatch(userId)

	if exists {
		t.Errorf("Expected no batch for the userId %s", userId)
	}

	if batch != nil {
		t.Errorf("Expected the returned batch to be equal to nil")
	}
}

func TestGetBatchHistory(t *testing.T) {
	cfg := config.Config{BatchThreshold: 100}
	transactionService := NewTransactionService(cfg)
	userId := uuid.NewString()
	expectedBatch := generateTestBatch(userId)
	transactionService.batchHistory[userId] = append(transactionService.batchHistory[userId], expectedBatch)

	batchHistory, exists := transactionService.GetBatchHistory(userId)

	if !exists {
		t.Errorf("Expected the batch history to exist for the userId %s", userId)
	}

	if !cmp.Equal(batchHistory[0], expectedBatch) {
		t.Errorf("Expected the returned batch history to contain the batch %s", expectedBatch.BatchId)
	}
}

func TestGetBatchHistoryDoesntExist(t *testing.T) {
	cfg := config.Config{BatchThreshold: 100}
	transactionService := NewTransactionService(cfg)
	userId := uuid.NewString()

	batchHistory, exists := transactionService.GetBatchHistory(userId)

	if exists {
		t.Errorf("Expected no batch history for the userId %s", userId)
	}

	if batchHistory != nil {
		t.Errorf("Expected the returned batch history to be equal to nil")
	}
}

func TestSaveTransaction(t *testing.T) {
	cfg := config.Config{BatchThreshold: 100}
	transactionService := NewTransactionService(cfg)
	userId := uuid.NewString()
	transaction := generateTestTransaction(userId)

	transactionService.SaveTransaction(transaction)
	_, exists := transactionService.GetBatch(userId)

	if !exists {
		t.Errorf("Expected a batch to exist for the userId %s", userId)
	}
}

func TestSaveTransactionExistentBatch(t *testing.T) {
	cfg := config.Config{BatchThreshold: 100}
	transactionService := NewTransactionService(cfg)
	userId := uuid.NewString()

	transaction1 := generateTestTransaction(userId)
	transactionService.SaveTransaction(transaction1)

	transaction2 := generateTestTransaction(userId)
	transactionService.SaveTransaction(transaction2)

	batch, exists := transactionService.GetBatch(userId)

	if !exists {
		t.Errorf("Expected a batch to exist for the userId %s", userId)
	}

	if !batch.AccruedAmount.Equal(transaction1.Amount.Add(transaction2.Amount)) {
		t.Errorf("Expected the batch's accrued amount to be equal to the sum of the transactions amount")
	}
}

func TestSaveTransactionMovesBatchToHistory(t *testing.T) {
	cfg := config.Config{BatchThreshold: 100}
	transactionService := NewTransactionService(cfg)
	userId := uuid.NewString()

	transaction1 := generateTestTransaction(userId)
	transactionService.SaveTransaction(transaction1)

	transaction2 := generateTestTransaction(userId)
	transactionService.SaveTransaction(transaction2)

	transaction3 := generateTestTransaction(userId)
	transactionService.SaveTransaction(transaction3)

	transaction4 := generateTestTransaction(userId)
	transactionService.SaveTransaction(transaction4)

	_, exists := transactionService.GetBatchHistory(userId)

	if !exists {
		t.Errorf("Expected a batch history to exist for the userId %s", userId)
	}
}

func generateTestTransaction(userId string) *app.Transaction {
	return &app.Transaction{
		TransactionId: uuid.NewString(),
		UserId:        userId,
		Amount:        decimal.NewFromFloat32(25.75),
		CreatedAt:     time.Now().Local(),
	}
}

func generateTestBatch(userId string) *app.Batch {
	transaction := generateTestTransaction(userId)
	return &app.Batch{
		BatchId:       uuid.NewString(),
		Transactions:  []app.Transaction{*transaction},
		AccruedAmount: decimal.NewFromFloat32(25.75),
		IsDispatched:  false,
	}
}
