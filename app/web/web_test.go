package web

import (
	"coding-challenge-go/app"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type MockTransactionService struct {
	userId       string
	batch        *app.Batch
	batchHistory []*app.Batch
}

func NewMockTransactionService(userId string) *MockTransactionService {
	var t MockTransactionService

	t.batch = generateTestBatch(userId)
	t.batchHistory = []*app.Batch{t.batch}

	return &t
}

func (transactionService *MockTransactionService) GetBatch(userId string) (*app.Batch, bool) {
	return transactionService.batch, true
}

func (transactionService *MockTransactionService) GetBatchHistory(userId string) ([]*app.Batch, bool) {
	return transactionService.batchHistory, true
}

func (transactionService *MockTransactionService) SaveTransaction(transaction *app.Transaction) {

}

func TestPostTransaction(t *testing.T) {
	userId := "540480ec-a932-4060-a4ee-14687a771735"
	mockTransactionService := NewMockTransactionService(userId)
	transactionHandler := TransactionHandler{TransactionService: mockTransactionService}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST",
		"/transaction",
		strings.NewReader("{ \"UserId\": \"540480ec-a932-4060-a4ee-14687a771735\", \"Amount\": 0.75}"))

	transactionHandler.postTransaction(w, r)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if body == nil {
		t.Errorf("Expected the response to contain a transaction")
	}
}

func TestGetBatch(t *testing.T) {
	userId := "540480ec-a932-4060-a4ee-14687a771735"
	mockTransactionService := NewMockTransactionService(userId)
	transactionHandler := TransactionHandler{TransactionService: mockTransactionService}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET","/batch/540480ec-a932-4060-a4ee-14687a771735", nil)

	transactionHandler.getBatch(w, r)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if body == nil {
		t.Errorf("Expected the response to contain a batch")
	}
}

func TestGetBatchHistory(t *testing.T) {
	userId := "540480ec-a932-4060-a4ee-14687a771735"
	mockTransactionService := NewMockTransactionService(userId)
	transactionHandler := TransactionHandler{TransactionService: mockTransactionService}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET","/batch/history/540480ec-a932-4060-a4ee-14687a771735", nil)

	transactionHandler.getBatchHistory(w, r)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if body == nil {
		t.Errorf("Expected the response to contain a history with a batch")
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
