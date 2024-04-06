package models

type Account struct {
	AccountName   string `json:"name"`
	AccountNumber string `json:"id"`
	Status        string `json:"status"`
}

type AccountResponse struct {
	AccountName   string `json:"account_name"`
	AccountNumber string `json:"account_number"`
	Status        string `json:"status"`
}
