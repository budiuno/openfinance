package models

type DisbursementRequest struct {
	Amount              int64  `json:"amount"`
	SourceBankCode      string `json:"source_bank_code"`
	SourceAccount       string `json:"source_account"`
	DestinationBankCode string `json:"destination_bank_code"`
	DestinationAccount  string `json:"destination_account"`
	Remarks             string `json:"remarks"`
}

type InsertDisbursementRequest struct {
	Amount              int64  `json:"amount"`
	SourceBankCode      string `json:"source_bank_code"`
	SourceAccount       string `json:"source_account"`
	DestinationBankCode string `json:"destination_bank_code"`
	DestinationAccount  string `json:"destination_account"`
	Remarks             string `json:"remarks"`
	ReferenceID         string `json:"reference_id"`
	Status              string `json:"status"`
}

type ProcessedDisbursement struct {
	ID                  int64  `json:"disbursement_id"`
	Amount              int64  `json:"amount"`
	SourceBankCode      string `json:"source_bank_code"`
	SourceAccount       string `json:"source_account"`
	DestinationBankCode string `json:"destination_bank_code"`
	DestinationAccount  string `json:"destination_account"`
	Remarks             string `json:"remarks"`
	ReferenceID         string `json:"reference_id"`
	Status              string `json:"status"`
}

// DisbursementResponse represents the response for a wallet disbursement.
type PostDisbursementResponse struct {
	TransactionId string `json:"id"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}

type DisbursementResponse struct {
	Processed []ProcessedDisbursement  `json:"success_to_processed"`
	Failed    []InvalidDisburseRequest `json:"fail_to_processed"`
}

type InvalidDisburseRequest struct {
	Request DisbursementRequest
	Error   string
}

type CallbackDisbursementRequest struct {
	ReferenceID int64  `json:"id"`
	Status      string `json:"status"`
}
