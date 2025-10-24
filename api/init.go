package api

import (
	"billing/api/handler"
	"billing/internal/model"
	"billing/internal/usecase"
	"billing/pkg/db"
	"log"

	"github.com/labstack/echo/v4"
)

// INIT YOUR DEPENDENCY HERE
func Init() *echo.Echo {
	e := echo.New()

	// IF YOU WANT TO USE DATABASE
	// YOU NEED TO ASSIGN TO VARIABLE AND PASS TO ROUTE
	// YOU NEED TO FILL STRUCT IN THIS FUNCTION TO CREATE TABLE IN DB
	// EXAMPLE : _, err := db.InitAndMigrate(&model.Loan{}, &model.InvestLoan{})
	db, err := db.InitAndMigrate(&model.Loan{}, &model.Billing{})
	if err != nil {
		log.Fatal(err)
	}

	loanUsecase := usecase.LoanUsecase {
		DB: db,
	}

	handler := handler.BillingHandler {
		LoanUsecase: &loanUsecase,
	}

	RegisterRoutes(e, handler)
	return e
}
