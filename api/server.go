package api

import (
	"log"

	"github.com/labstack/echo/v4"
)

// DON'T CHANGE THIS FUNCTION
func StartServer(e *echo.Echo) {
	if err := e.Start(":8000"); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
