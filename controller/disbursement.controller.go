package Controller

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	// rsp "github.com/budiuno/openfinance/middlewares/response"
	rsp "github.com/budiuno/openfinance/middlewares/response"
	"github.com/budiuno/openfinance/models"
	"github.com/google/uuid"
)

type disburseSetter interface {
	InsertDisburseToDB(db *sql.DB, req models.InsertDisbursementRequest) (int, error)
	PostDisbursement(req models.PostDisbursementRequest) (int64, error)
	UpdateDisbursementStatus(db *sql.DB, referenceID uuid.UUID, newStatus string) error
}

func DisburseHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, repoDisburse disburseSetter) {

	var requests []models.DisbursementRequest
	invalidRequests := make([]models.InvalidDisburseRequest, 0)
	processedRequests := make([]models.ProcessedDisbursement, 0)

	if err := json.NewDecoder(r.Body).Decode(&requests); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a WaitGroup to wait for all Goroutines to finish
	var wg sync.WaitGroup

	// Create a buffered channel as a semaphore to limit parallel processing
	semaphore := make(chan struct{}, 5)

	// Process each DisbursementRequest object in parallel with limited concurrency
	for _, req := range requests {
		// Acquire a slot in the semaphore
		semaphore <- struct{}{}

		// Increment WaitGroup counter
		wg.Add(1)

		// Execute processing logic in a Goroutine
		go func(req models.DisbursementRequest) {
			// Decrement WaitGroup counter when Goroutine completes
			defer func() {
				// Release slot in the semaphore after processing is done
				<-semaphore
				wg.Done()
			}()

			// Perform validation
			err := isValid(req)
			if err != nil {
				invalidRequests = append(invalidRequests, models.InvalidDisburseRequest{Request: req, Error: "Validation failed"})
			} else {
				// Hit external endpoint
				refId := uuid.New()
				_, err := repoDisburse.PostDisbursement(populatePostDisbursementRequest(req, refId))
				if err != nil {
					message := fmt.Sprintf("failed to post disbursement, error: %v", err)
					invalidRequests = append(invalidRequests, models.InvalidDisburseRequest{Request: req, Error: message})
				} else {
					// insert to DB with refId
					insertReq := populateInsertDisbursementRequest(req, refId, "pending")
					pk, err := repoDisburse.InsertDisburseToDB(db, insertReq)
					if err != nil {
						message := fmt.Sprintf("failed to insert disbursement, error: %v", err)
						invalidRequests = append(invalidRequests, models.InvalidDisburseRequest{Request: req, Error: message})
					} else {
						processedReq := populateProcessedDisbursement(insertReq, int64(pk))
						processedRequests = append(processedRequests, processedReq)
					}

				}

			}
		}(req)
	}
	// Wait for all Goroutines to finish
	wg.Wait()

	resFinal := models.DisbursementResponse{
		Processed: processedRequests,
		Failed:    invalidRequests,
	}

	rsp.RespondWithJSON(w, resFinal, http.StatusOK)
}

func populateInsertDisbursementRequest(req models.DisbursementRequest, refID uuid.UUID, status string) models.InsertDisbursementRequest {
	return models.InsertDisbursementRequest{
		Amount:              req.Amount,
		SourceBankCode:      req.SourceBankCode,
		SourceAccount:       req.SourceAccount,
		DestinationBankCode: req.DestinationBankCode,
		DestinationAccount:  req.DestinationAccount,
		Remarks:             req.Remarks,
		ReferenceID:         refID,
		Status:              status,
	}
}

func populatePostDisbursementRequest(req models.DisbursementRequest, refID uuid.UUID) models.PostDisbursementRequest {
	return models.PostDisbursementRequest{
		Amount:              req.Amount,
		SourceBankCode:      req.SourceBankCode,
		SourceAccount:       req.SourceAccount,
		DestinationBankCode: req.DestinationBankCode,
		DestinationAccount:  req.DestinationAccount,
		Remarks:             req.Remarks,
		ReferenceID:         refID,
	}
}

func populateProcessedDisbursement(req models.InsertDisbursementRequest, disbursementID int64) models.ProcessedDisbursement {
	return models.ProcessedDisbursement{
		ID:                  disbursementID,
		Amount:              req.Amount,
		SourceBankCode:      req.SourceBankCode,
		SourceAccount:       req.SourceAccount,
		DestinationBankCode: req.DestinationBankCode,
		DestinationAccount:  req.DestinationAccount,
		Remarks:             req.Remarks,
		ReferenceID:         req.ReferenceID,
		Status:              req.Status,
	}
}

func isValid(req models.DisbursementRequest) error {
	// Perform validation logic
	if req.Amount <= 0 {
		return errors.New("amount should be greater than 0")
	}

	if req.SourceBankCode == "" {
		return errors.New("aource_bank_code cannot be empty")
	}
	if req.SourceAccount == "" {
		return errors.New("aource_account cannot be empty")
	}
	if req.DestinationBankCode == "" {
		return errors.New("destination_bank_code cannot be empty")
	}
	if req.DestinationAccount == "" {
		return errors.New("destination_account cannot be empty")
	}

	return nil
}
