package scoring

import (
	"billing/internal/model"
	"billing/internal/util"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCase defines a reusable test case structure
type TestCase struct {
	Name           string
	Method         string
	URL            string
	Body           string
	ExpectedStatus int
	ExpectedResult interface{}
	ParamNames     []string
	ParamValues    []string
}

// TestCreateBills_Success tests successful bill creation
func TestCreateBills_Success(t *testing.T) {
	reqBody := model.Loan{
		ID:           "loan123",
		CustomerID:   "cust123",
		Name:         "Test Loan",
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
	assert.Equal(t, reqBody.ID, respData.Loan.ID)
	assert.Equal(t, reqBody.CustomerID, respData.Loan.CustomerID)
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

// TestCreateBills_InvalidJSON tests bill creation with invalid JSON
func TestCreateBills_InvalidJSON(t *testing.T) {
	req := mapAPI[APICreatedBill]
	req.Body = []byte(`{"id":"loan123", invalid json}`)
	rec := callAPI(req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestCreateBills_NegativeValues tests bill creation with negative values
func TestCreateBills_NegativeValues(t *testing.T) {
	reqBody := model.Loan{
		ID:           "loan123",
		CustomerID:   "cust123",
		Name:         "Test Loan",
		Period:       -50,
		Amount:       5000000,
		InterestRate: 10,
	}
	req := mapAPI[APICreatedBill]
	req.Body = reqBody
	rec := callAPI(req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestGetBills_NewLoan tests getting bills for a newly created loan
func TestGetBills_NewLoan(t *testing.T) {
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
	assert.Equal(t, loan.Loan.CustomerID, respData.Loan.CustomerID)
	assert.NotNil(t, respData.Bills, "Bills array should not be nil")
	assert.Equal(t, 50, len(respData.Bills), "Should generate 50 bills for period=50")
}

// TestGetBills_MissingLoanID tests getting bills with missing loan ID
func TestGetBills_MissingLoanID(t *testing.T) {
	req := mapAPI[APIGetBill]
	rec := callAPI(req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestGetBillStatus_NotDelinquent tests bill status when not delinquent
func TestGetBillStatus_NotDelinquent(t *testing.T) {
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

	// Validate delinquency status
	assert.False(t, respData.IsDelinquent, "Should not be delinquent")
	assert.True(t, respData.DelinquentAt.IsZero(), "DelinquentAt should be zero when not delinquent")
}

// TestGetBillStatus_Delinquent_TwoConsecutiveMissedPayments tests bill status with two consecutive missed payments
func TestGetBillStatus_Delinquent_TwoConsecutiveMissedPayments(t *testing.T) {
	loan, err := seedData()
	if err != nil {
		t.Fatalf("Failed to seed data: %v", err)
	}
	cleanup := addTimeNow(3)
	defer cleanup()

	req := mapAPI[APIGetBillStatus]
	req.Param = map[string]string{
		"loan_id": loan.Loan.ID,
	}
	rec := callAPI(req)
	assert.Equal(t, http.StatusOK, rec.Code)

	respData, err := unmarshalResponse[model.BillingStatus](rec)
	assert.NoError(t, err)
	assert.Equal(t, loan.Loan.ID, respData.LoanID)

	// Validate delinquency status
	assert.True(t, respData.IsDelinquent, "Should be delinquent with two consecutive missed payments")
}
func TestGetBillStatus_InvalidRequest(t *testing.T) {
	req := mapAPI[APIGetBillStatus]
	rec := callAPI(req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestMakePayment_ExactWeeklyAmount tests making payment with exact weekly amount
func TestMakePayment_ExactWeeklyAmount(t *testing.T) {
	loan, err := seedData()
	if err != nil {
		t.Fatalf("Failed to seed data: %v", err)
	}
	cleanup := addTimeNow(1)
	defer cleanup()

	reqBody := model.MakePaymentRequest{
		PaymentAmount: 110000,
		PaymentDate:   util.GetCurrentTime(),
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
	assert.Equal(t, reqBody.LoanID, respData.LoanID)

	// Validate payment amount
	expectedWeeklyPayment := 110000.0
	assert.InDelta(t, expectedWeeklyPayment, respData.Amount, 0.01, "Payment amount should match exact weekly payment")
}

// TestMakePayment_MultipleWeeklyPayments tests making payment for multiple weeks
func TestMakePayment_MultipleWeeklyPayments(t *testing.T) {
	loan, err := seedData()
	if err != nil {
		t.Fatalf("Failed to seed data: %v", err)
	}
	cleanup := addTimeNow(2)
	defer cleanup()

	reqBody := model.MakePaymentRequest{
		PaymentAmount: 220000,
		PaymentDate:   util.GetCurrentTime(),
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
	assert.Equal(t, reqBody.LoanID, respData.LoanID)

	// Validate payment amount
	expectedWeeklyPayment := 110000.0
	expectedMultiplePayment := expectedWeeklyPayment * 2 // 2 weeks worth
	assert.InDelta(t, expectedMultiplePayment, respData.Amount, 0.01, "Payment amount should match multiple weekly payments")
	// This payment should be applied to the first 2 unpaid bills
}

// TestMakePayment_PartialAmount tests making payment with partial amount (should fail)
func TestMakePayment_PartialAmount(t *testing.T) {
	loan, err := seedData()
	if err != nil {
		t.Fatalf("Failed to seed data: %v", err)
	}
	cleanup := addTimeNow(1)
	defer cleanup()

	reqBody := model.MakePaymentRequest{
		PaymentAmount: 55000,
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

// TestMakePayment_NonExactAmount tests making payment with non-exact amount (should fail)
func TestMakePayment_NonExactAmount(t *testing.T) {
	loan, err := seedData()
	if err != nil {
		t.Fatalf("Failed to seed data: %v", err)
	}
	cleanup := addTimeNow(1)
	defer cleanup()

	reqBody := model.MakePaymentRequest{
		PaymentAmount: 120000,
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

// TestMakePayment_NegativeAmount tests making payment with negative amount
func TestMakePayment_NegativeAmount(t *testing.T) {
	loan, err := seedData()
	if err != nil {
		t.Fatalf("Failed to seed data: %v", err)
	}
	cleanup := addTimeNow(1)
	defer cleanup()

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

// TestMakePayment_ZeroAmount tests making payment with zero amount
func TestMakePayment_ZeroAmount(t *testing.T) {
	loan, err := seedData()
	if err != nil {
		t.Fatalf("Failed to seed data: %v", err)
	}
	cleanup := addTimeNow(1)
	defer cleanup()

	reqBody := model.MakePaymentRequest{
		PaymentAmount: 0,
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
