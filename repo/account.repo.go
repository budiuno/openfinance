package repo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	cfg "github.com/budiuno/openfinance/config"
	"github.com/budiuno/openfinance/models"
)

type AccountRepo struct{}

func (repo AccountRepo) GetAccountByAccountNumber(accountID string) (models.AccountResponse, error) {
	url := cfg.BaseURL + "/v1/accounts/" + accountID

	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error making HTTP request: %v\n", err)
		return models.AccountResponse{}, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		err := errors.New("account not found")
		fmt.Printf("Error: %s\n", err)
		return models.AccountResponse{}, err
	}

	// Decode JSON response into struct
	var acc models.Account
	if err := json.NewDecoder(response.Body).Decode(&acc); err != nil {
		fmt.Printf("Error decoding JSON: %v\n", err)
		return models.AccountResponse{}, err
	}
	accResponse := convertToAccountResponse(acc)

	// Print account
	return accResponse, nil
}

func convertToAccountResponse(account models.Account) models.AccountResponse {
	return models.AccountResponse{
		AccountName:   account.AccountName,
		AccountNumber: account.AccountNumber,
		Status:        account.Status,
	}
}
