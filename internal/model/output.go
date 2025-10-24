package model

import "time"

// DON'T CHANGE THIS STRUCT
type LoanWithBills struct {
	Loan  Loan      `json:"loan"`
	Bills []Billing `json:"bills"`
}

// DON'T CHANGE THIS STRUCT
type Billing struct {
	ID          string     `json:"id"`
	LoanID      string     `json:"loan_id"`
	Sequence    int        `json:"sequence"`
	Date        time.Time  `json:"date"`
	DueDate     time.Time  `json:"due_date"`
	PaymentDate *time.Time `json:"payment_date"`
	Amount      float64    `json:"amount"`
	CreatedAt   time.Time  `json:"created_at"`
}

// DON'T CHANGE THIS STRUCT
type BillingStatus struct {
	LoanID       string    `json:"loan_id"`
	IsDelinquent bool      `json:"is_delinquent"`
	DelinquentAt time.Time `json:"delinquent_at"` // DelinquentAt should be the date of second consecutive missed payment
}

// DON'T CHANGE THIS STRUCT
type Payment struct {
	LoanID string    `json:"loan_id"`
	Amount float64   `json:"amount"`
	Date   time.Time `json:"date"`
}
