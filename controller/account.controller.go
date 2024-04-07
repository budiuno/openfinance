package Controller

import (
	"net/http"
	"strings"

	rsp "github.com/budiuno/openfinance/middlewares/response"
	"github.com/budiuno/openfinance/models"

	"github.com/gorilla/mux"
)

type AccountGetter interface {
	GetAccountByAccountNumber(accountNumber string) (models.AccountResponse, error)
}

func ValidateAccount(w http.ResponseWriter, r *http.Request, getAccount AccountGetter) {
	vars := mux.Vars(r)
	bankCode := strings.ToLower(vars["bank_code"])
	accountNumber := vars["account_number"]

	if bankCode != "tsb" {
		errorMessage := rsp.ErrorResponse{IsError: true, Message: "Unrecognized bank code"}
		rsp.RespondWithError(w, errorMessage, http.StatusNotFound)
		return
	}

	// call mockapi(test-bank) to check account number exist
	account, err := getAccount.GetAccountByAccountNumber(accountNumber)
	if err != nil {
		// If the account is not found, return "account not found" as JSON
		errorMessage := rsp.ErrorResponse{IsError: true, Message: "Account not found"}
		rsp.RespondWithError(w, errorMessage, http.StatusNotFound)
		return
	}

	// If account found, return it as JSON
	rsp.RespondWithJSON(w, account, http.StatusOK)
}
