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
	err, db_version := model.ConnectDB()
	if err != nil {
		inf.LogError(err, "Failed to connect to the database:")
		return
	}
	defer model.DisconnectDB()

	inf.LogMessage("info", "Successfully connected to the Database!")
	inf.LogMessage("info", fmt.Sprintf("  Schema Version: %s", db_version))

	inf.InitCookieStore()
	inf.LoadTemplates(e)
	InitRoutes(e)
	e.HTTPErrorHandler = inf.RedirectHTTPErrorHandler

	// Start the Server
	err = e.Start(":" + os.Getenv("WISHLIST_PORT"))
	if err != nil {
		inf.LogError(err, "Server exited with a fatal error:")
		return
	}
}
