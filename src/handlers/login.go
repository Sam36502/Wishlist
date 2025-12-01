/*
 *		LOGIN
 *
 *		All functions pertaining to logging users in/out
 *
 */

package handlers

import (
	"net/http"
	"wishlist/src/inf"
	"wishlist/src/inf/model"

	"github.com/labstack/echo/v4"
)

/// LOGIN

func PgLogin(c echo.Context) error {
	var err error
	defer inf.RecoverPanic(c, &err)

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
			inf.LogError(err, "Failed to delete form-data")
			return echo.ErrInternalServerError
		}
	}

	return c.Render(http.StatusOK, "login", data)
}

func LoginUser(c echo.Context) error {
	// Get User data from form
	formUser := new(inf.FormUser)
	err := c.Bind(formUser)
	if err != nil {
		return err
	}

	// Add user to context for if there's an error
	session, err := inf.CookieStore.Get(c.Request(), inf.COOKIE_FORM_DATA)
	if err != nil {
		inf.LogError(err, "Failed to get form-data cookie:")
		return echo.ErrInternalServerError
	}
	session.Values["user"] = formUser
	session.Values["error"] = new(inf.UserFormError)

	// Input Validation
	hasError := false
	formError := new(inf.UserFormError)

	// Check all fields were filled
	if formUser.Password == "" {
		hasError = true
		formError.Password = "Password is required"
	}

	tok, err := model.Authenticate(formUser.Email, formUser.Password)
	if err != nil {
		inf.LogError(err, "Login Failed!")
		hasError = true
		formError.Password = "Invalid Email/Password"
	}

	// Check Errors
	if hasError {
		session.Values["error"] = formError
		err = session.Save(c.Request(), c.Response())
		if err != nil {
			inf.LogError(err, "Failed to save form-data:")
			return echo.ErrInternalServerError
		}
		return c.Redirect(http.StatusMovedPermanently, "/login")
	}

	// Add Token to cookie
	tokenData, err := inf.CookieStore.Get(c.Request(), inf.COOKIE_TOKEN_DATA)
	if err != nil {
		inf.LogError(err, "Failed to get token cookie:")
		return echo.ErrInternalServerError
	}
	tokenData.Values[inf.COOKIE_TOKEN_DATA] = tok

	// Set cookie to only expire after a month if the "remember me" option is set
	if formUser.StayLoggedIn == "on" {
		tokenData.Options.MaxAge = inf.COOKIE_TIMEOUT
	}

	err = tokenData.Save(c.Request(), c.Response())
	if err != nil {
		inf.LogError(err, "Failed to save token cookie:")
		return echo.ErrInternalServerError
	}

	return c.Redirect(http.StatusMovedPermanently, "/user/"+formUser.Email)
}

func Logout(c echo.Context) error {

	// Delete Token Cookie
	logout_message := inf.StatusPageData{
		Colour:          "green",
		MainMessage:     "Successfully logged out!",
		NextPageURL:     "/",
		NextPageMessage: "Back to main page",
	}
	userData, err := inf.CookieStore.Get(c.Request(), inf.COOKIE_TOKEN_DATA)
	if err != nil {
		logout_message.MainMessage = "Your session expired and you were logged out"
		logout_message.Colour = "red"
	}

	// Delete the cookie
	userData.Options.MaxAge = -1

	err = userData.Save(c.Request(), c.Response())
	if err != nil {
		inf.LogError(err, "Failed to save token cookie:")
		return echo.ErrInternalServerError
	}

	return c.Render(http.StatusOK, "status", logout_message)
}
