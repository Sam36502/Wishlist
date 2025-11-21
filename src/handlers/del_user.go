package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func PgDeleteUser(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func DeleteUser(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/")
}
