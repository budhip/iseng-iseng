package model

import "time"

// DON'T CHANGE THIS STRUCT
type Loan struct {
	ID           string    `json:"id"`
	CustomerID   string    `json:"customer_id"`
	Name         string    `json:"name"`
	Period       int       `json:"period"`
	Amount       float64   `json:"amount"`
	InterestRate float64   `json:"interest_rate"`
	TotalAmount  float64   `json:"total_amount"`
	Outstanding  float64   `json:"outstanding"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

// DON'T CHANGE THIS STRUCT
type MakePaymentRequest struct {
	LoanID        string    `json:"loan_id"`
	PaymentAmount float64   `json:"payment_amount"`
	PaymentDate   time.Time `json:"payment_date"`
}
