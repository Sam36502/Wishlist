package main

import (
	"fmt"
	"os"

	"wishlist/src/inf"
	"wishlist/src/inf/model"

	"github.com/labstack/echo/v4"
)

func main() {

	// Initialise Echo Framework
	e := echo.New()

	// Initialise Infrastructure
	err := model.ConnectDB()
	if err != nil {
		fmt.Printf("Failed to connect to the database:\n  %v", err)
		return
	}
	defer model.DisconnectDB()

	inf.InitCookieStore()
	inf.LoadTemplates(e)
	InitRoutes(e)
	e.HTTPErrorHandler = inf.RedirectHTTPErrorHandler

	// Start the Server
	err = e.Start(":" + os.Getenv("WISHLIST_PORT"))
	if err != nil {
		fmt.Printf("Server exited with a fatal error:\n  %v", err)
		return
	}
}
