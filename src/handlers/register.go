/*
 *		REGISTER
 *
 *		All functions pertaining to registering new users
 *
 */

package handlers

import (
	"fmt"
	"net/http"
	"wishlist/src/inf"
	"wishlist/src/inf/model"

	"github.com/labstack/echo/v4"
)

// The register page renderer
func PgRegister(c echo.Context) error {

	// Check for data
	data := new(struct {
		User  inf.FormUser
		Error inf.UserFormError
	})
	session, err := inf.CookieStore.Get(c.Request(), inf.COOKIE_FORM_DATA)
	if err == nil {
		if e := session.Values["error"]; e != nil {
			if formErr, ok := e.(inf.UserFormError); ok {
				data.Error = formErr
			}
		}
		if u := session.Values["user"]; u != nil {
			if formUsr, ok := u.(inf.FormUser); ok {
				data.User = formUsr
			}
		}
		session.Options.MaxAge = -1
		err = session.Save(c.Request(), c.Response())
		if err != nil {
			fmt.Println("[ERROR] Failed to delete form-data:\n ", err)
			return echo.ErrInternalServerError
		}
	}

	return c.Render(http.StatusOK, "register", data)
}

// Receives user's data and registers them in the DB
func RegisterUser(c echo.Context) error {

	// Get User data from form
	formUser := new(inf.FormUser)
	err := c.Bind(formUser)
	if err != nil {
		return err
	}

	// Add user to context for if there's an error
	session, err := inf.CookieStore.Get(c.Request(), inf.COOKIE_FORM_DATA)
	if err != nil {
		fmt.Println("[ERROR] Failed to get form-data cookie:\n ", err)
		return echo.ErrInternalServerError
	}
	session.Values["user"] = formUser
	session.Values["error"] = new(inf.UserFormError)

	// Input Validation
	hasError := false
	formError := new(inf.UserFormError)

	// Check all fields were filled
	if formUser.Name == "" {
		hasError = true
		formError.Name = "Name is required"
	}
	if formUser.Password == "" || formUser.PasswordConfirm == "" {
		hasError = true
		formError.Password = "Password & confirmation are required"
	}

	// Check Email is Valid
	if !inf.IsEmailValid(formUser.Email) {
		formError.Email = "Please use a valid E-Mail address"
		hasError = true
	}

	// Check Password aligns with confirm field
	if formUser.Password != formUser.PasswordConfirm {
		formError.Password = "Passwords don't match."
		hasError = true
	}

	// Check password matches constraints
	if errMsg := inf.CheckPasswordRequirements(formUser.Password); errMsg != "" {
		hasError = true
		formError.Password = errMsg
	}

	if !hasError {
		// Register the User
		user := model.User{
			Name:     formUser.Name,
			Password: formUser.Password,
			Email:    formUser.Email,
		}

		err = model.InsertUser(&user)
		if err != nil {
			hasError = true
			if _, ok := err.(model.EmailExistsError); ok {
				formError.Email = "This Email is already in use"
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
		c.Redirect(http.StatusMovedPermanently, "/register")
	}

	return c.Redirect(http.StatusMovedPermanently, "/registersuccess")
}

// The page displayed when the user successfully registers their user
func PgRegisterSuccess(c echo.Context) error {
	return c.Render(http.StatusOK, "status", inf.StatusPageData{
		Colour:          "green",
		MainMessage:     "New user successfully registered!",
		NextPageURL:     "/login",
		NextPageMessage: "Continue to login",
	})
}
