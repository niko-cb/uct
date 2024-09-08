package entity

import "time"

// Invoice represents the invoice data
type Invoice struct {
	ID            int64     `json:"id"`
	CompanyID     int64     `json:"company_id"`
	ClientID      int64     `json:"client_id"`
	IssueDate     time.Time `json:"issue_date"`
	DueDate       time.Time `json:"due_date"`
	PaymentAmount float64   `json:"payment_amount"`
	FeeAmount     float64   `json:"fee_amount"`
	TaxAmount     float64   `json:"tax_amount"`
	TotalAmount   float64   `json:"total_amount"`
	Status        string    `json:"status"`
}
