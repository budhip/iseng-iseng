package tests

import (
	"billing/internal/model"
	"billing/internal/util"
	"net/http"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCase defines basic test case structure for API contract testing
type TestCase struct {
	Name           string
	Method         string
	Path           string
	Body           string
	ExpectedStatus int
	HasLoanID      bool
	LoanID         string
}

// TestMainLintCheck runs the linter and marks all subsequent tests to fail if linting fails
func TestMainLintCheck(t *testing.T) {
	cmd := exec.Command("golangci-lint", "run")
	err := cmd.Run()

	if err != nil {
		t.Fail()
	}
}

func TestCreateBills_Success(t *testing.T) {
	reqBody := model.Loan{
		CustomerID:   "cust123",
		Period:       50,
		Amount:       5000000,
		InterestRate: 10,
	}
	req := mapAPI[APICreatedBill]
	req.Body = reqBody
	rec := callAPI(req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	respData, err := unmarshalResponse[model.LoanWithBills](rec)
	assert.NoError(t, err)
	assert.NotEmpty(t, respData.Loan.ID)
	assert.NotEmpty(t, respData.Loan.CustomerID)
	assert.Equal(t, 50, respData.Loan.Period)
	assert.InDelta(t, 5000000.0, respData.Loan.Amount, 0.01)
	assert.Equal(t, 50, len(respData.Bills), "Should generate 50 bills for period=50")
	assert.InDelta(t, 5500000.0, respData.Loan.TotalAmount, 0.01)
	assert.InDelta(t, 5500000.0, respData.Loan.Outstanding, 0.01)

	// Validate total amounts across all bills
	var totalAmount float64
	for _, bill := range respData.Bills {
		totalAmount += bill.Amount
	}
	assert.InDelta(t, 5500000.0, totalAmount, 0.01, "Total amount should equal loan total amount (principal + interest)")
}

// TestCreateBills_InvalidJSON tests POST /bills with invalid JSON
func TestCreateBills_InvalidJSON(t *testing.T) {
	req := mapAPI[APICreatedBill]
	req.Body = []byte(`{"id":"loan123", invalid json}`)
	rec := callAPI(req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestCreateBills_EmptyBody tests POST /bills with empty body
func TestCreateBills_EmptyBody(t *testing.T) {
	req := mapAPI[APICreatedBill]
	rec := callAPI(req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestGetBills_ValidLoanID tests GET /bills/:loan_id with valid loan ID
func TestGetBills_ValidLoanID(t *testing.T) {
	loan, err := seedData()
	if err != nil {
		t.Fatalf("Failed to seed data: %v", err)
	}
	req := mapAPI[APIGetBill]
	req.Param = map[string]string{
		"loan_id": loan.Loan.ID,
	}

	rec := callAPI(req)
	assert.Equal(t, http.StatusOK, rec.Code)

	respData, err := unmarshalResponse[model.LoanWithBills](rec)
	assert.NoError(t, err)
	assert.Equal(t, loan.Loan.ID, respData.Loan.ID)
	assert.Equal(t, 50, len(respData.Bills), "Should generate 50 bills for period=50")
}

// TestGetBills_EmptyLoanID tests GET /bills/:loan_id with empty loan ID
func TestGetBills_EmptyLoanID(t *testing.T) {
	req := mapAPI[APIGetBill]
	rec := callAPI(req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestGetBillStatus_ValidRequest tests GET /bills/:loan_id/status with valid request
func TestGetBillStatus_ValidRequest(t *testing.T) {
	loan, err := seedData()
	if err != nil {
		t.Fatalf("Failed to seed data: %v", err)
	}
	req := mapAPI[APIGetBillStatus]
	req.Param = map[string]string{
		"loan_id": loan.Loan.ID,
	}
	rec := callAPI(req)
	assert.Equal(t, http.StatusOK, rec.Code)
	respData, err := unmarshalResponse[model.BillingStatus](rec)
	assert.NoError(t, err)
	assert.Equal(t, loan.Loan.ID, respData.LoanID)
	assert.Equal(t, false, respData.IsDelinquent)
}

// TestGetBillStatus_ValidRequest tests GET /bills/:loan_id/status with valid request and delinquent
func TestGetBillStatus_ValidRequestDelinquent(t *testing.T) {
	loan, err := seedData()
	if err != nil {
		t.Fatalf("Failed to seed data: %v", err)
	}
	resetTime := addTimeNow(3)
	defer resetTime()
	req := mapAPI[APIGetBillStatus]
	req.Param = map[string]string{
		"loan_id": loan.Loan.ID,
	}
	rec := callAPI(req)
	assert.Equal(t, http.StatusOK, rec.Code)
	respData, err := unmarshalResponse[model.BillingStatus](rec)
	assert.NoError(t, err)
	assert.Equal(t, loan.Loan.ID, respData.LoanID)
	assert.Equal(t, true, respData.IsDelinquent)
}

// TestGetBillStatus_InvalidRequest tests GET /bills/:loan_id/status with invalid loan ID
func TestGetBillStatus_InvalidRequest(t *testing.T) {
	req := mapAPI[APIGetBillStatus]
	rec := callAPI(req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestMakePayment_ValidRequest tests POST /bills/:loan_id/payments with valid request
func TestMakePayment_ValidRequest(t *testing.T) {
	loan, err := seedData()
	if err != nil {
		t.Fatalf("Failed to seed data: %v", err)
	}
	reqBody := model.MakePaymentRequest{
		PaymentAmount: 110000,
		PaymentDate:   loan.Bills[0].DueDate,
	}
	req := mapAPI[APIMakePayment]
	req.Param = map[string]string{
		"loan_id": loan.Loan.ID,
	}
	req.Body = reqBody
	rec := callAPI(req)
	assert.Equal(t, http.StatusOK, rec.Code)
	respData, err := unmarshalResponse[model.Payment](rec)
	assert.NoError(t, err)
	assert.Equal(t, reqBody.PaymentAmount, respData.Amount)
	assert.Equal(t, loan.Loan.ID, respData.LoanID)
}

// TestMakePayment_InvalidJSON tests POST /bills/:loan_id/payments with invalid JSON
func TestMakePayment_InvalidJSON(t *testing.T) {
	loan, err := seedData()
	if err != nil {
		t.Fatalf("Failed to seed data: %v", err)
	}
	req := mapAPI[APIMakePayment]
	req.Param = map[string]string{
		"loan_id": loan.Loan.ID,
	}
	req.Body = []byte(`{"payment_amount":550000, invalid json}`)
	rec := callAPI(req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestMakePayment_EmptyBody tests POST /bills/:loan_id/payments with empty body
func TestMakePayment_EmptyBody(t *testing.T) {
	loan, err := seedData()
	if err != nil {
		t.Fatalf("Failed to seed data: %v", err)
	}
	req := mapAPI[APIMakePayment]
	req.Param = map[string]string{
		"loan_id": loan.Loan.ID,
	}
	rec := callAPI(req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestMakePayment_NegativeAmount tests POST /bills/:loan_id/payments with negative amount
func TestMakePayment_NegativeAmount(t *testing.T) {
	loan, err := seedData()
	if err != nil {
		t.Fatalf("Failed to seed data: %v", err)
	}
	reqBody := model.MakePaymentRequest{
		PaymentAmount: -110000,
		PaymentDate:   util.GetCurrentTime(),
	}
	req := mapAPI[APIMakePayment]
	req.Param = map[string]string{
		"loan_id": loan.Loan.ID,
	}
	req.Body = reqBody
	rec := callAPI(req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
