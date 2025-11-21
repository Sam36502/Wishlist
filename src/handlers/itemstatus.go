package handlers

import (
	"net/http"
	"strconv"
	"wishlist/src/inf"
	"wishlist/src/inf/model"

	"github.com/labstack/echo/v4"
)

func PgReserveItem(c echo.Context) error {
	email := c.Param("email")
	if email == "" {
		return echo.ErrNotFound
	}

	item, err := inf.GetItemFromPath(c)
	if err != nil {
		return echo.ErrNotFound
	}

	return c.Render(http.StatusOK, "confirm", inf.ConfirmPageData{
		MainTitle:       "Reserve to purchase " + item.Name,
		MainDescription: "Are you sure you want to mark that you're planning to get this item?",
		YesColour:       "gold",
		YesURL:          "/user/" + email + "/item/" + strconv.FormatUint(item.ItemID, 10) + "/reserve",
		YesText:         "Yes, I'm planning to buy this as a gift",
		NoColour:        "gray",
		NoURL:           "/user/" + email + "/item/" + strconv.FormatUint(item.ItemID, 10),
		NoText:          "No, I'm not getting this item",
	})
}

func ReserveItem(c echo.Context) error {
	email := c.Param("email")
	if email == "" {
		return echo.ErrNotFound
	}

	item, err := inf.GetItemFromPath(c)
	if err != nil {
		return echo.ErrNotFound
	}

	// Check if currently logged in as this user
	// Only logged-in users other than yourself can reserve items
	liUser, _, err := inf.GetLoggedInUser(c)
	if err == nil {
		if email == liUser.Email {
			return echo.ErrForbidden
		}
	} else {
		return echo.ErrUnauthorized
	}

	// Mark item as reserved
	item.Status = model.Status{StatusID: 2} // Status 2 --> Reserved
	err = model.UpdateItem(&item)
	if err != nil {
		return c.Render(http.StatusOK, "status", inf.StatusPageData{
			Colour:          "red",
			MainMessage:     "Failed to reserve item. Please try again later.",
			NextPageURL:     "/user/" + email + "/item/" + strconv.FormatUint(item.ItemID, 10),
			NextPageMessage: "Back to " + item.Name + " item page",
		})
	}

	return c.Render(http.StatusOK, "status", inf.StatusPageData{
		Colour:          "green",
		MainMessage:     "Item successfully reserved!",
		NextPageURL:     "/user/" + email + "/item/" + strconv.FormatUint(item.ItemID, 10),
		NextPageMessage: "Back to " + item.Name + " item page",
	})
}

func PgUnreserveItem(c echo.Context) error {
	email := c.Param("email")
	if email == "" {
		return echo.ErrNotFound
	}

	item, err := inf.GetItemFromPath(c)
	if err != nil {
		return echo.ErrNotFound
	}

	return c.Render(http.StatusOK, "confirm", inf.ConfirmPageData{
		MainTitle:       "Remove reservation for " + item.Name,
		MainDescription: "Are you sure you want to mark that you're no longer planning to purchase this item?",
		YesColour:       "red",
		YesURL:          "/user/" + email + "/item/" + strconv.FormatUint(item.ItemID, 10) + "/unreserve",
		YesText:         "Yes, I'm no longer planning to get this gift",
		NoColour:        "gray",
		NoURL:           "/user/" + email + "/item/" + strconv.FormatUint(item.ItemID, 10),
		NoText:          "No, I'm still going to get it",
	})
}

func UnreserveItem(c echo.Context) error {
	email := c.Param("email")
	if email == "" {
		return echo.ErrNotFound
	}

	item, err := inf.GetItemFromPath(c)
	if err != nil {
		return echo.ErrNotFound
	}

	// Check if currently logged in as the user who reserved it
	liUser, _, err := inf.GetLoggedInUser(c)
	if err == nil {
		if item.ReservedByUser.Email != liUser.Email {
			return echo.ErrForbidden
		}
	} else {
		return echo.ErrUnauthorized
	}

	// Mark item as available
	item.Status = model.Status{StatusID: 1} // Status 1 --> Available
	item.ReservedByUser = nil
	err = model.UpdateItem(&item)
	if err != nil {
		return c.Render(http.StatusOK, "status", inf.StatusPageData{
			Colour:          "red",
			MainMessage:     "Failed to remove reservation. Please try again later.",
			NextPageURL:     "/user/" + email + "/item/" + strconv.FormatUint(item.ItemID, 10),
			NextPageMessage: "Back to " + item.Name + " item page",
		})
	}

	return c.Render(http.StatusOK, "status", inf.StatusPageData{
		Colour:          "green",
		MainMessage:     "Reservation successfully removed!",
		NextPageURL:     "/user/" + email + "/item/" + strconv.FormatUint(item.ItemID, 10),
		NextPageMessage: "Back to " + item.Name + " item page",
	})
}

func PgReceiveItem(c echo.Context) error {
	email := c.Param("email")
	if email == "" {
		return echo.ErrNotFound
	}

	item, err := inf.GetItemFromPath(c)
	if err != nil {
		return echo.ErrNotFound
	}

	return c.Render(http.StatusOK, "confirm", inf.ConfirmPageData{
		MainTitle:       "Marking " + item.Name + " as received",
		MainDescription: "Are you sure you want to mark that you received this gift from someone?",
		YesColour:       "green",
		YesURL:          "/user/" + email + "/item/" + strconv.FormatUint(item.ItemID, 10) + "/receive",
		YesText:         "Yes, I've received this item",
		NoColour:        "gray",
		NoURL:           "/user/" + email + "/item/" + strconv.FormatUint(item.ItemID, 10),
		NoText:          "No, I haven't received this item yet",
	})
}

func ReceiveItem(c echo.Context) error {
	email := c.Param("email")
	if email == "" {
		return echo.ErrNotFound
	}

	item, err := inf.GetItemFromPath(c)
	if err != nil {
		return echo.ErrNotFound
	}

	// Check if currently logged in as this user
	liUser, _, err := inf.GetLoggedInUser(c)
	if err == nil {
		if email != liUser.Email {
			return echo.ErrForbidden
		}
	} else {
		return echo.ErrUnauthorized
	}

	// Mark item as received
	item.Status = model.Status{StatusID: 3} // Status 3 --> Received
	model.UpdateItem(&item)

	return c.Render(http.StatusOK, "status", inf.StatusPageData{
		Colour:          "green",
		MainMessage:     "Item successfully marked as received!",
		NextPageURL:     "/user/" + email + "/item/" + strconv.FormatUint(item.ItemID, 10),
		NextPageMessage: "Back to " + item.Name + " item page",
	})
}

func PgUnReceiveItem(c echo.Context) error {
	email := c.Param("email")
	if email == "" {
		return echo.ErrNotFound
	}

	item, err := inf.GetItemFromPath(c)
	if err != nil {
		return echo.ErrNotFound
	}

	return c.Render(http.StatusOK, "confirm", inf.ConfirmPageData{
		MainTitle:       "Marking " + item.Name + " as not received",
		MainDescription: "Are you sure you want to mark that you haven't received this?",
		YesColour:       "red",
		YesURL:          "/user/" + email + "/item/" + strconv.FormatUint(item.ItemID, 10) + "/unreceive",
		YesText:         "Yes, I didn't receive this item",
		NoColour:        "gray",
		NoURL:           "/user/" + email + "/item/" + strconv.FormatUint(item.ItemID, 10),
		NoText:          "No, I did receive the item",
	})
}

func UnReceiveItem(c echo.Context) error {
	email := c.Param("email")
	if email == "" {
		return echo.ErrNotFound
	}

	item, err := inf.GetItemFromPath(c)
	if err != nil {
		return echo.ErrNotFound
	}

	// Check if currently logged in as this user
	liUser, _, err := inf.GetLoggedInUser(c)
	if err == nil {
		if email != liUser.Email {
			return echo.ErrForbidden
		}
	} else {
		return echo.ErrUnauthorized
	}

	// Mark item as available
	item.Status = model.Status{StatusID: 1} // Status 1 --> Available
	err = model.UpdateItem(&item)
	if err != nil {
		return c.Render(http.StatusOK, "status", inf.StatusPageData{
			Colour:          "red",
			MainMessage:     "Failed to mark item as available. Please try again later.",
			NextPageURL:     "/user/" + email + "/item/" + strconv.FormatUint(item.ItemID, 10),
			NextPageMessage: "Back to " + item.Name + " item page",
		})
	}

	return c.Render(http.StatusOK, "status", inf.StatusPageData{
		Colour:          "green",
		MainMessage:     "Item successfully marked as available again!",
		NextPageURL:     "/user/" + email + "/item/" + strconv.FormatUint(item.ItemID, 10),
		NextPageMessage: "Back to " + item.Name + " item page",
	})
}
