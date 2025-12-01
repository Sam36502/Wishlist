package main

import (
	"encoding/csv"
	"net/http"
	"os"
	"strings"
	"wishlist/src/handlers"
	"wishlist/src/inf"

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

type Post struct {
	Title   string
	Version string
	Text    string
	Changes []string
}

func PgMain(c echo.Context) error {
	var err error
	defer inf.RecoverPanic(c, &err)

	// Load Blog Posts from CSV
	var posts []Post
	f, err := os.Open("data/changelog.csv")
	if err == nil {
		csvReader := csv.NewReader(f)
		records, err := csvReader.ReadAll()
		if err == nil {
			for _, post := range records {
				posts = append(posts, Post{
					Title:   post[0],
					Version: post[1],
					Text:    post[2],
					Changes: strings.Split(post[3], ";"),
				})
			}
		}
	}
	defer f.Close()

	return c.Render(http.StatusOK, "main", posts)
}
