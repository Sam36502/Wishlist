/*
 *		USER LIST
 *
 *		All functions pertaining to displaying and
 *		editing the user's wishlist items
 *
 */

package handlers

import (
	"net/http"
	"wishlist/src/inf"
	"wishlist/src/inf/model"

	"github.com/labstack/echo/v4"
)

func PgUserList(c echo.Context) (err error) {
	defer inf.RecoverPanic(c, &err)

	// Get all the user's items
	email := c.Param("email")
	if email == "" {
		return echo.ErrNotFound
	}

	user, err := model.GetUserWithEmail(email)
	if err != nil {
		return c.Redirect(http.StatusPermanentRedirect, "/home")
	}
	all_items, err := model.GetAllItems(user.ID)
	if err != nil {
		inf.LogError(err, "Failed to retrieve User's items:")
		return c.Redirect(http.StatusPermanentRedirect, "/err/500.html")
	}
	var items = make([]model.Item, len(all_items))
	for i, item := range all_items {
		items[i] = *item
	}

	// Check if currently logged in as this user
	thisIsMe := false
	liUser, _, err := inf.GetLoggedInUser(c)
	if err == nil {
		thisIsMe = user.Email == liUser.Email
	}

	return c.Render(http.StatusOK, "user_list", struct {
		User            model.User
		Items           []model.Item
		GetStatusColour func(model.Status) string
		ThisIsMe        bool
	}{
		User:            *user,
		Items:           items,
		GetStatusColour: inf.GetStatusColour,
		ThisIsMe:        thisIsMe,
	})
}
