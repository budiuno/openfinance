package repo

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/budiuno/openfinance/models"
)

type DisburseRepo struct{}

func (repo DisburseRepo) InsertDisburseToDB(db *sql.DB, req models.InsertDisbursementRequest) (int, error) {
	// Perform database insertion logic

	// Construct the SQL query for insertion
	query := `
		INSERT INTO disbursements (
			amount,
			source_bank_code,
			source_account,
			destination_bank_code,
			destination_account,
			reference_id,
			remarks,
			status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;
	`

	// Execute the SQL query with the provided parameters
	var id int
	err := db.QueryRow(
		query,
		req.Amount,
		req.SourceBankCode,
		req.SourceAccount,
		req.DestinationBankCode,
		req.DestinationAccount,
		req.ReferenceID,
		req.Status,
		req.Remarks,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo DisburseRepo) PostDisbursement(req models.DisbursementRequest) (int64, error) {

	// Convert DisbursementRequest to JSON
	reqBody, err := json.Marshal(req)
	if err != nil {
		return 0, fmt.Errorf("error marshalling request body: %v", err)
	}

	// Make POST request to external API
	resp, err := http.Post("https://66100a360640280f219c2844.mockapi.io/api/v1/disburse", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return 0, fmt.Errorf("error making POST request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status code is successful (2xx)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, fmt.Errorf("unexpected response status code: %d", resp.StatusCode)
	}

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %v", err)
	}

	// Unmarshal the response body to get the ID
	var response struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return 0, fmt.Errorf("error unmarshalling response body: %v", err)
	}

	intID, err := strconv.ParseInt(response.ID, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error convert id to int64 %v", err)
	}

	return intID, nil
}
