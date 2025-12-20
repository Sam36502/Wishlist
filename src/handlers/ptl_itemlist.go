/*
 *		Partial Item-List
 *
 *		Loads and renders *just* the item datatable
 *		for use with HTMX
 *
 */

package handlers

import (
	"net/http"
	"wishlist/src/front"
	"wishlist/src/inf"
	"wishlist/src/model"

	"github.com/labstack/echo/v4"
)

func PtlUserList(c echo.Context) error {

	// Get all the user's items
	email := c.Param("email")
	if email == "" {
		return echo.ErrNotFound
	}

	show_received := c.QueryParams().Get("cb_show_received") == "on"
	show_reserved := c.QueryParams().Get("cb_show_reserved") == "on"

	user, err := model.GetUserWithEmail(email)
	if err != nil {
		return c.Redirect(http.StatusPermanentRedirect, "/home")
	}

	// TODO: Do item filtering on db-side
	allItems, err := model.GetAllItems(user.ID)
	if err != nil {
		inf.LogError(err, "Failed to retrieve User's items:")
		return c.Redirect(http.StatusPermanentRedirect, "/err/500.html")
	}

	// For now: Filtering in-memory here
	items := make([]model.Item, 0)
	for _, item := range allItems {
		skip := (!show_received && item.Status.StatusID == model.STATUS_ID_RECEIVED) ||
			(!show_reserved && item.Status.StatusID == model.STATUS_ID_RESERVED)

		if skip {
			continue
		}

		items = append(items, *item)
	}

	// Check if currently logged in as this user
	thisIsMe := false
	liUser, _, err := front.GetLoggedInUser(c)
	if err == nil {
		thisIsMe = user.Email == liUser.Email
	}

	return c.Render(http.StatusOK, "itemlist", struct {
		User            model.User
		Items           []model.Item
		GetStatusColour func(model.Status) string
		ThisIsMe        bool
	}{
		User:            *user,
		Items:           items,
		GetStatusColour: front.GetStatusColour,
		ThisIsMe:        thisIsMe,
	})
}
