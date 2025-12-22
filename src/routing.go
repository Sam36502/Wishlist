package main

import (
	"net/http"
	"wishlist/src/front"
	"wishlist/src/handlers"
	"wishlist/src/inf"
	"wishlist/src/model"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {

	e.GET("/", PgMain)
	e.Static("/", "data/static")

	// User Routes
	e.GET("/search", handlers.PgSearch)

	e.GET("/register", handlers.PgRegister)
	e.POST("/register", handlers.RegisterUser)
	e.GET("/registersuccess", handlers.PgRegisterSuccess)

	e.GET("/login", handlers.PgLogin)
	e.POST("/login", handlers.LoginUser)
	e.GET("/logout", handlers.Logout)

	e.GET("/user/:email", handlers.PgUserList)
	e.GET("/user/:email/chgpassword", handlers.PgChangePassword)
	e.POST("/user/:email/chgpassword", handlers.ChangePassword)
	e.GET("/user/:email/delete", handlers.PgDeleteUser)
	e.POST("/user/:email/delete", handlers.DeleteUser)

	// Item Routes
	e.GET("/user/:email/item/:item_id", handlers.PgItem)
	e.GET("/user/:email/newitem", handlers.PgNewItem)
	e.POST("/user/:email/newitem", handlers.NewItem)

	e.GET("/user/:email/item/:item_id/delete", handlers.PgDelItem)
	e.POST("/user/:email/item/:item_id/delete", handlers.DelItem)

	e.GET("/user/:email/item/:item_id/reserve", handlers.PgReserveItem)
	e.POST("/user/:email/item/:item_id/reserve", handlers.ReserveItem)
	e.GET("/user/:email/item/:item_id/unreserve", handlers.PgUnreserveItem)
	e.POST("/user/:email/item/:item_id/unreserve", handlers.UnreserveItem)
	e.GET("/user/:email/item/:item_id/receive", handlers.PgReceiveItem)
	e.POST("/user/:email/item/:item_id/receive", handlers.ReceiveItem)
	e.GET("/user/:email/item/:item_id/unreceive", handlers.PgUnReceiveItem)
	e.POST("/user/:email/item/:item_id/unreceive", handlers.UnReceiveItem)

	// Partial Routes
	e.GET("/ptl/userlist/:email", handlers.PtlUserList)

}

func PgMain(c echo.Context) error {
	var err error
	defer inf.RecoverPanic(c, &err)

	clog, err := model.GetChangelogEntries()
	for i, e := range clog {
		r := front.Renderable(e)
		html, _ := front.RenderHTML(r)
		clog[i].Rendered = &html
	}

	return c.Render(http.StatusOK, "main", clog)
}
