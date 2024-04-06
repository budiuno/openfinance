package accountController

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	rsp "github.com/budiuno/openfinance/middlewares/response"
	"github.com/budiuno/openfinance/models"
	"github.com/gorilla/mux"
)

type MockAccountGetter struct{}

func (m MockAccountGetter) GetAccountByAccountNumber(accountNumber string) (models.AccountResponse, error) {
	// Return a mock account for testing purposes
	if accountNumber == "00000000" {
		return models.AccountResponse{}, errors.New("account not found")
	}

	accounts := map[string]models.AccountResponse{
		"12341234": {AccountName: "Skipper", AccountNumber: "12341234", Status: "active"},
		"56785678": {AccountName: "Kowalski", AccountNumber: "56785678", Status: "active"},
	}
	return accounts[accountNumber], nil
}

func TestValidateAccount(t *testing.T) {

	tests := []struct {
		name                   string
		requestBankCode        string
		requestAccountNumber   string
		expectedAccount        *models.AccountResponse
		expectedError          *rsp.ErrorResponse
		expectedHTTPStatusCode int
	}{
		{
			name:                   "valid account",
			requestBankCode:        "tsb",
			requestAccountNumber:   "12341234",
			expectedAccount:        &models.AccountResponse{AccountName: "Skipper", AccountNumber: "12341234", Status: "active"},
			expectedError:          nil,
			expectedHTTPStatusCode: 200,
		},
		{
			name:                   "unrecognize bank code",
			requestBankCode:        "bca",
			requestAccountNumber:   "12341234",
			expectedAccount:        nil,
			expectedError:          &rsp.ErrorResponse{IsError: true, Message: "Unrecognized bank code"},
			expectedHTTPStatusCode: 404,
		},
		{
			name:                   "account not found",
			requestBankCode:        "tsb",
			requestAccountNumber:   "00000000",
			expectedAccount:        nil,
			expectedError:          &rsp.ErrorResponse{IsError: true, Message: "Account not found"},
			expectedHTTPStatusCode: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/v1/account/" + tt.requestBankCode + "/" + tt.requestAccountNumber

			req := httptest.NewRequest("GET", url, nil)

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			r := mux.NewRouter()
			mockAccountGetter := MockAccountGetter{}
			r.HandleFunc("/v1/account/{bank_code}/{account_number}", func(w http.ResponseWriter, r *http.Request) {
				ValidateAccount(w, r, mockAccountGetter)
			}).Methods("GET")

			// Simulate the request by serving the HTTP request to the router
			r.ServeHTTP(rr, req)

			if rr.Code != tt.expectedHTTPStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedHTTPStatusCode, rr.Code)
			}

			if tt.expectedAccount != nil {
				var response models.AccountResponse
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Fatal(err)
				}

				// Check if account detail is correct
				if response.AccountName != tt.expectedAccount.AccountName || response.AccountNumber != tt.expectedAccount.AccountNumber || response.Status != tt.expectedAccount.Status {
					t.Errorf("Expected Account Name : %s, got %s , Expected Account Number : %s, got %s and Expected status : %s, got %s", tt.expectedAccount.AccountName, response.AccountName, tt.expectedAccount.AccountNumber, response.AccountNumber, tt.expectedAccount.Status, response.Status)
				}
			}

			if tt.expectedError != nil {
				var response rsp.ErrorResponse
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Fatal(err)
				}

				if response.IsError != tt.expectedError.IsError || response.Message != tt.expectedError.Message {
					t.Errorf("Expected isError : %t, got %t and Expected error message : %s , got %s\n", tt.expectedError.IsError, response.IsError, tt.expectedError.Message, response.Message)
				}

			}
		})
	}

}
