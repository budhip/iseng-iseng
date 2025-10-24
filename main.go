// Package main starts the billing engine HTTP server.
package main

import (
	"billing/api"
)

// DON'T CHANGE THIS FUNCTION, PLEASE INIT YOUR DEPENDENCY IN api.Init()
func main() {
	e := api.Init()
	api.StartServer(e)
}
