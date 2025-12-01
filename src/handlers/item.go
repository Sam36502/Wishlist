package handlers

import (
	"net/http"
	"wishlist/src/inf"
	"wishlist/src/inf/model"

	"github.com/labstack/echo/v4"
)

func PgItem(c echo.Context) error {
	var err error
	defer inf.RecoverPanic(c, &err)

	email := c.Param("email")
	if email == "" {
		return echo.ErrNotFound
	}

	// Get Item information
	item, err := inf.GetItemFromPath(c)
	if err != nil {
		return err
	}

	// Check if currently logged in as this user
	loggedInHere := false
	loggedIn := false
	liUser, _, err := inf.GetLoggedInUser(c)
	if err == nil {
		loggedIn = true
		loggedInHere = liUser.Email == email
	}

	// Check if user can reserve this item
	canReserve := loggedIn &&
		!loggedInHere &&
		item.Status.StatusID == model.STATUS_ID_AVAILABLE

	// Check if user can unreserve this item
	canUnreserve := loggedIn &&
		item.Status.StatusID == model.STATUS_ID_RESERVED &&
		item.ReservedByUser != nil &&
		item.ReservedByUser.Email == liUser.Email

	// Check if user can see the status of this item
	canSeeStatus := !loggedInHere || item.Status.StatusID == model.STATUS_ID_RECEIVED

	// Check if user can mark item as received
	canReceive := loggedInHere && item.Status.StatusID != model.STATUS_ID_RECEIVED

	// Check if user can unreceive item
	canUnreceive := loggedInHere && item.Status.StatusID == model.STATUS_ID_RECEIVED

	return c.Render(http.StatusOK, "item", struct {
		Item         model.Item
		LoggedIn     bool
		CanSeeStatus bool
		CanReserve   bool
		CanUnreserve bool
		CanReceive   bool
		CanUnreceive bool
		StatusColour string
		Email        string
	}{
		Item:         item,
		LoggedIn:     loggedInHere,
		CanSeeStatus: canSeeStatus,
		CanReserve:   canReserve,
		CanUnreserve: canUnreserve,
		CanReceive:   canReceive,
		CanUnreceive: canUnreceive,
		StatusColour: inf.GetStatusColour(item.Status),
		Email:        email,
	})
}
