package handlers

import (
	"net/http"
	"wishlist/src/inf"

	"github.com/labstack/echo/v4"
)

func PgDeleteUser(c echo.Context) error {
	var err error
	defer inf.RecoverPanic(c, &err)

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func DeleteUser(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/")
}
