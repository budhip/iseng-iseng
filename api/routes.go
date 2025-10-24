package api

import (
	"billing/api/handler"

	"github.com/labstack/echo/v4"
)

// DON'T CHANGE ANY PATH & METHOD
func RegisterRoutes(e *echo.Echo, handler handler.BillingHandler) {
	e.GET("/bills/:loan_id", handler.GetBills)
	e.GET("/bills/:loan_id/status", handler.GetBillStatus)
	e.POST("/bills", handler.CreateBills)
	e.POST("/bills/:loan_id/payments", handler.MakePayment)
}
