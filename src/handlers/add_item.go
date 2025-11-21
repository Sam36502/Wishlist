package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"wishlist/src/inf"
	"wishlist/src/inf/model"

	"github.com/labstack/echo/v4"
)

func PgNewItem(c echo.Context) error {

	email := c.Param("email")
	if email == "" {
		return echo.ErrNotFound
	}

	// Check for form errors and/or retained values
	data := new(struct {
		User  struct{ Email string }
		Item  inf.FormItem
		Error inf.ItemFormError
	})
	data.User.Email = email
	session, err := inf.CookieStore.Get(c.Request(), inf.COOKIE_FORM_DATA)
	if err == nil {
		if e := session.Values["error"]; e != nil {
			if formErr, ok := e.(inf.ItemFormError); ok {
				data.Error = formErr
			}
		}
		if i := session.Values["item"]; i != nil {
			if formItem, ok := i.(inf.FormItem); ok {
				data.Item = formItem
			}
		}
		session.Options.MaxAge = -1
		err = session.Save(c.Request(), c.Response())
		if err != nil {
			fmt.Println("[ERROR] Failed to delete form-data:\n ", err)
			return echo.ErrInternalServerError
		}
	}

	return c.Render(http.StatusOK, "add_item", data)
}

func NewItem(c echo.Context) error {

	email := c.Param("email")
	if email == "" {
		return echo.ErrNotFound
	}

	// Get Item data from form
	formItem := new(inf.FormItem)
	err := c.Bind(formItem)
	if err != nil {
		return err
	}

	// Add item to context for if there's an error
	session, err := inf.CookieStore.Get(c.Request(), inf.COOKIE_FORM_DATA)
	if err != nil {
		fmt.Println("[ERROR] Failed to get form-data cookie:\n ", err)
		return echo.ErrInternalServerError
	}
	session.Values["item"] = formItem
	session.Values["error"] = new(inf.UserFormError)

	// Input Validation
	hasError := false
	formError := new(inf.ItemFormError)

	// Check all required fields were filled
	if formItem.Name == "" {
		hasError = true
		formError.Name = "Item Name is Required"
	}

	var price float64
	if formItem.Price == "" {
		hasError = true
		formError.Price = "Item Price is Required"
	} else {
		price, err = strconv.ParseFloat(formItem.Price, 32)
		if err != nil {
			hasError = true
			formError.Price = "Invalid Price entered"
		}
	}

	// Add URL if not already included
	if formItem.LinkText != "" && formItem.LinkURL == "" {
		formItem.LinkURL = formItem.LinkText
	}
	if formItem.LinkURL != "" && formItem.LinkText == "" {
		formItem.LinkText = formItem.LinkURL
	}

	// Check if link URL is valid
	_, err = url.ParseRequestURI(formItem.LinkURL)
	if formItem.LinkURL != "" && err != nil {
		hasError = true
		formError.Link = "Please use a valid URL"
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

	// Add item
	if !hasError {
		err = model.InsertItem(&model.Item{
			Name:        formItem.Name,
			Description: formItem.Description,
			Price:       float32(price), // parsed above
			Status:      model.Status{StatusID: 1},
			User:        liUser,
			Links: []model.Link{
				{
					Text: formItem.LinkText,
					URL:  formItem.LinkURL,
				},
			},
		})
		if err != nil {
			hasError = true
			fmt.Println("[ERROR] Failed to add the item to the database:\n ", err)

			if _, ok := err.(model.PriceOutOfRangeError); ok {
				formError.Price = "Price was out of range"
			}
		}
	}

	if hasError {
		session.Values["error"] = formError
		err = session.Save(c.Request(), c.Response())
		if err != nil {
			fmt.Println("[ERROR] Failed to save form-data:\n ", err)
			return echo.ErrInternalServerError
		}
		c.Redirect(http.StatusMovedPermanently, "/user/"+email+"/newitem")
	}

	return c.Redirect(http.StatusMovedPermanently, "/user/"+email)
}
