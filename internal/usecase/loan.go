package usecase

import (
	"billing/internal/model"
	"billing/internal/util"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LoanUsecase struct {
	DB *gorm.DB
}

var ErrNoPendingBill = errors.New("no pending bills for spcified payment_date")
var ErrInsufficientAmount = errors.New("insuffiecient payment amount")

func (u *LoanUsecase) isLoanIDExist(loanID string) (*model.Loan, error) {
	var loan model.Loan
	if err := u.DB.Where("id = ?", loanID).First(&loan).Error; err != nil {
		return nil, err 
	}
	return &loan, nil
}

func (u *LoanUsecase) CreateBills(req model.Loan) (*model.LoanWithBills, error) {
	var err error
	var resp model.LoanWithBills
	timeNow := util.GetCurrentTime()

	interest := req.Amount * req.InterestRate / 100
	req.ID = uuid.New().String()
	req.TotalAmount = req.Amount + interest
	req.Outstanding = req.TotalAmount
	req.CreatedAt = timeNow
	req.Status = "IN_PROGRESS"

	if err = u.DB.Save(&req).Error; err != nil {
		log.Println("[CreateBills] Failed to create loan", err)
		return nil, err
	}

	billings := make([]model.Billing, 0)
	billsPerMonth := req.TotalAmount / float64(req.Period)
	currentDate := timeNow.Add(7 * 24 * time.Hour)
	for i := 0; i < req.Period; i++ {
		billings = append(billings, model.Billing{
			ID:        uuid.New().String(),
			LoanID:    req.ID,
			Sequence:  i + 1,
			Date:      timeNow,
			DueDate:   currentDate,
			Amount:    billsPerMonth,
			CreatedAt: timeNow,
		})

		currentDate = currentDate.Add(7 * 24 * time.Hour)
	}

	if err = u.DB.Save(&billings).Error; err != nil {
		log.Println("[CreateBills] Failed to create billings", err)
		return nil, err
	}

	resp.Loan = req
	resp.Bills = billings

	return &resp, nil
}

func (u *LoanUsecase) GetBills(loanID string) (*model.LoanWithBills, error) {
	var err error
	var resp model.LoanWithBills

	var loan model.Loan
	if err = u.DB.Where("id = ?", loanID).First(&loan).Error; err != nil {
		log.Println("[GetBills] Failed to get loan", err)
		return nil, err
	}

	var bills []model.Billing
	if err = u.DB.Where("loan_id = ?", loanID).Find(&bills).Error; err != nil {
		log.Println("[GetBills] Failed to get bills", err)
		return nil, err
	}

	resp.Loan = loan
	resp.Bills = bills

	return &resp, nil
}

func (u *LoanUsecase) GetBillStatus(loanID string) (*model.BillingStatus, error) {
	resp := model.BillingStatus{
		LoanID: loanID,
	}
	timeNow := util.GetCurrentTime()

	_, err := u.isLoanIDExist(loanID)
	if err != nil {
		return nil, err
	}

	var bills []model.Billing
	if err := u.DB.Where("loan_id = ? AND DATE(due_date) < ? AND payment_date IS NULL", loanID, timeNow).Find(&bills).Error; err != nil {
		log.Println("[GetBillStatus] Failed to get bills", err)
		return nil, err
	}

	if len(bills) >= 2 {
		resp.IsDelinquent = true
		resp.DelinquentAt = bills[1].DueDate
	} else {
		resp.IsDelinquent = false
	}

	return &resp, nil
}

func (u *LoanUsecase) MakePayment(req model.MakePaymentRequest) (*model.Payment, error) {
	var resp model.Payment

	loan, err := u.isLoanIDExist(req.LoanID) 
	if err != nil {
		return nil, err
	}

	var bill model.Billing
	if err := u.DB.Where("loan_id = ? AND payment_date IS NULL AND DATE(due_date) < ?", req.LoanID, req.PaymentDate).Order("due_date").First(&bill).Error; err != nil {
		return nil, ErrNoPendingBill
	}

	if req.PaymentAmount != bill.Amount {
		return nil, ErrInsufficientAmount
	}

	if err := u.DB.Model(&model.Billing{}).Where("id = ?", bill.ID).Update("payment_date", req.PaymentDate).Error; err != nil {
		log.Println("[MakePayment] Failed to update bill", err)
		return nil, err
	}

	status := "IN_PROGRESS"
	if loan.Outstanding - req.PaymentAmount <= 0 {
		status = "COMPLETED"
	}

	if err := u.DB.Model(&model.Loan{}).Where("id = ?", req.LoanID).Updates(map[string]interface{}{
		"outstanding": gorm.Expr("outstanding - ?", bill.Amount),
		"status": status,
	}).Error; err != nil {
		log.Println("[MakePayment] Failed to update loan's outstanding", err)
		return nil, err
	}

	resp.LoanID = req.LoanID
	resp.Amount = req.PaymentAmount
	resp.Date = req.PaymentDate

	return &resp, nil
}