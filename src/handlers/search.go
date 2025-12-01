/*
 *		SEARCH
 *
 *		All functions pertaining to searching for users
 *
 */

package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"wishlist/src/inf"
	"wishlist/src/inf/model"

	"github.com/labstack/echo/v4"
)

type searchResultData struct {
	Search       string
	ResultString string
	Results      []model.User
}

func PgSearch(c echo.Context) error {
	var err error
	defer inf.RecoverPanic(c, &err)

	// Get Search Query
	if !c.QueryParams().Has("s") {
		return c.Render(http.StatusOK, "search", searchResultData{
			Search:       "",
			ResultString: "Enter a name above to find friends.",
			Results:      nil,
		})
	}
	search := c.QueryParam("s")
	search = strings.ToLower(search)
	search = strings.TrimSpace(search)

	// Search Users
	// TODO: Implement search in DB / with proper index
	all_users, err := model.GetAllUsers()
	if err != nil {
		inf.LogError(err, "Failed to search users:")
		return echo.ErrInternalServerError
	}
	var users []model.User
	for _, u := range all_users {
		match := strings.Contains(strings.ToLower(u.Name), search)
		match = match || strings.Contains(strings.ToLower(u.Email), search)

		if match {
			users = append(users, *u)
		}
	}

	resString := strconv.Itoa(len(users)) + " Results:"
	switch len(users) {
	case 0:
		resString = "No Results"
	case 1:
		resString = "1 Result:"
	}

	return c.Render(http.StatusOK, "search", searchResultData{
		Search:       search,
		ResultString: resString,
		Results:      users,
	})
}
