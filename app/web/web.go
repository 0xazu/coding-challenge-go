package web

import (
	"coding-challenge-go/app"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

// Handles the service requests
type TransactionHandler struct {
	TransactionService app.TransactionService
}

// Creates the handlers for the service operations and starts the service
func (h *TransactionHandler) Handle() {
	router := mux.NewRouter()

	router.HandleFunc("/transaction", h.postTransaction)
	router.HandleFunc("/batch/{userId}", h.getBatch)
	router.HandleFunc("/batch/history/{userId}", h.getBatchHistory)

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (h *TransactionHandler) postTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	transaction := app.Transaction{TransactionId: uuid.NewString(), CreatedAt: time.Now().Local()}
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, app.InvalidTransactionRequestBody, http.StatusBadRequest)
		return
	}

	h.TransactionService.SaveTransaction(&transaction)

	transactionResponse, err := json.Marshal(transaction)
	if err != nil {
		http.Error(w, app.InvalidTransactionResponseBody, http.StatusInternalServerError)
		return
	}

	h.writeResponse(w, transactionResponse)
}

func (h *TransactionHandler) getBatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}

	vars := mux.Vars(r)
	var batch, exists = h.TransactionService.GetBatch(vars["userId"])
	if !exists {
		http.Error(w, app.BatchNotFound, http.StatusNotFound)
		return
	}

	batchResponse, err := json.Marshal(batch)
	if err != nil {
		http.Error(w, app.InvalidBatchResponseBody, http.StatusInternalServerError)
		return
	}

	h.writeResponse(w, batchResponse)
}

func (h *TransactionHandler) getBatchHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}

	vars := mux.Vars(r)
	var batchHistory, exists = h.TransactionService.GetBatchHistory(vars["userId"])
	if !exists {
		http.Error(w, app.BatchHistoryNotFound, http.StatusNotFound)
		return
	}

	batchHistoryResponse, err := json.Marshal(batchHistory)
	if err != nil {
		http.Error(w, app.InvalidBatchHistoryResponseBody, http.StatusInternalServerError)
		return
	}

	h.writeResponse(w, batchHistoryResponse)
}

func (h *TransactionHandler) writeResponse(w http.ResponseWriter, transactionResponse []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(transactionResponse)
}
