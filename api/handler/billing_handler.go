package handler

import (
	"billing/api/response"
	"billing/internal/model"
	"billing/internal/usecase"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type BillingHandler struct {
	LoanUsecase *usecase.LoanUsecase
}

/*
REQUIRED RESPONSE :
- All field in Loan
- All field in Bills
*/
func (h *BillingHandler) CreateBills(c echo.Context) error {
	req := model.Loan{}

	err := c.Bind(&req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	if req.CustomerID == "0" || req.CustomerID == "" {
		return response.Error(c, http.StatusBadRequest, "customer_id is required")
	}

	resp, err := h.LoanUsecase.CreateBills(req)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Internal Server Error")
	}

	return response.Created(c, resp)
}

/*
REQUIRED RESPONSE :
- All field in Loan
- All field in Bills
*/
func (h *BillingHandler) GetBills(c echo.Context) error {
	loanID := strings.TrimSpace(c.Param("loan_id"))

	if loanID == "" || loanID == "0" {
		return response.Error(c, http.StatusBadRequest, "loan_id is required")
	}

	resp, err := h.LoanUsecase.GetBills(loanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(c, http.StatusBadRequest, fmt.Sprintf("loan_id %s not found", loanID))
		}
		return response.Error(c, http.StatusInternalServerError, "Internal Server Error")
	}

	return response.Success(c, resp)
}

/*
REQUIRED RESPONSE :
- LoanID
- IsDelinquent
- DelinquentAt
*/
func (h *BillingHandler) GetBillStatus(c echo.Context) error {
	loanID := strings.TrimSpace(c.Param("loan_id"))

	if loanID == "" || loanID == "0" {
		return response.Error(c, http.StatusBadRequest, "loan_id is required")
	}

	resp, err := h.LoanUsecase.GetBillStatus(loanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(c, http.StatusBadRequest, fmt.Sprintf("loan_id %s not found", loanID))
		}
		return response.Error(c, http.StatusInternalServerError, "Internal Server Error")
	}

	return response.Success(c, resp)
}

/*
REQUIRED RESPONSE :
- LoanID
- Amount
- Date
*/
func (h *BillingHandler) MakePayment(c echo.Context) error {
	loanID := strings.TrimSpace(c.Param("loan_id"))

	req := model.MakePaymentRequest{}

	err := c.Bind(&req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}
	fmt.Println("NAH", req)

	if req.PaymentAmount < 0 {
		return response.Error(c, http.StatusBadRequest, "payment amount less than 0")
	}

	req.LoanID = loanID
	fmt.Println("INU", req)

	resp, err := h.LoanUsecase.MakePayment(req)
	if err != nil {
		if errors.Is(err, usecase.ErrInsufficientAmount) || errors.Is(err, usecase.ErrNoPendingBill){
			return response.Error(c, http.StatusBadRequest, err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, "Internal Server Error")
	}

	return response.Success(c, resp)
}
