package handlers

import (
	"net/http"
	"strconv"
	"wishlist/src/front"
	"wishlist/src/inf"

	"github.com/labstack/echo/v4"
)

func PgDelItem(c echo.Context) error {
	var err error
	defer inf.RecoverPanic(c, &err)

	email := c.Param("email")
	if email == "" {
		return echo.ErrNotFound
	}

	// Get item info
	item, err := front.GetItemFromPath(c)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "confirm", inf.ConfirmPageData{
		MainTitle:       "Deleting " + item.Name + "...",
		MainDescription: "Are you sure you want to delete this item?",
		YesColour:       "red",
		YesText:         "Yes, delete this item permanently",
		YesURL:          "/user/" + email + "/item/" + strconv.FormatUint(item.ID, 10) + "/delete",
		NoColour:        "gray",
		NoText:          "No, don't delete it",
		NoURL:           "/user/" + email + "/item/" + strconv.FormatUint(item.ID, 10),
	})
}

func DelItem(c echo.Context) error {

	email := c.Param("email")
	if email == "" {
		return echo.ErrNotFound
	}

	item, err := front.GetItemFromPath(c)
	if err != nil {
		return nil
	}

	// Check if currently logged in as this user
	liUser, _, err := front.GetLoggedInUser(c)
	if err == nil {
		if email != liUser.Email {
			return echo.ErrForbidden
		}
	} else {
		return echo.ErrUnauthorized
	}

	// Delete item
	err = item.Delete()
	if err != nil {
		return c.Render(http.StatusOK, "status", inf.StatusPageData{
			Colour:          "red",
			MainMessage:     "Failed to delete item. Please try again later",
			NextPageURL:     "/user/" + email,
			NextPageMessage: "Back to user page",
		})
	}

	return c.Redirect(http.StatusMovedPermanently, "/user/"+email)
}
