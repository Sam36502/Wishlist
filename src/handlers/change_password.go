package handlers

import (
	"net/http"
	"wishlist/src/front"
	"wishlist/src/inf"
	"wishlist/src/model"

	"github.com/labstack/echo/v4"
)

func PgChangePassword(c echo.Context) error {
	var err error
	defer inf.RecoverPanic(c, &err)

	email := c.Param("email")
	if email == "" {
		return echo.ErrNotFound
	}

	// Check for data
	data := new(struct {
		User  inf.FormUser
		Error inf.UserFormError
		Email string
	})
	data.Email = email
	session, err := front.CookieStore.Get(c.Request(), front.COOKIE_FORM_DATA)
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
			inf.LogError(err, "Failed to delete form-data:")
			return echo.ErrInternalServerError
		}
	}

	return c.Render(http.StatusOK, "change_password", data)
}

func ChangePassword(c echo.Context) error {
	email := c.Param("email")
	if email == "" {
		return echo.ErrNotFound
	}

	// Get User data from form
	formUser := new(inf.FormUser)
	err := c.Bind(formUser)
	if err != nil {
		return err
	}

	// Add user to context for if there's an error
	session, err := front.CookieStore.Get(c.Request(), front.COOKIE_FORM_DATA)
	if err != nil {
		inf.LogError(err, "Failed to get form-data cookie:")
		return echo.ErrInternalServerError
	}
	session.Values["user"] = formUser
	session.Values["error"] = new(inf.UserFormError)

	// Input Validation
	hasError := false
	formError := new(inf.UserFormError)

	// Check all fields are filled
	if formUser.Password == "" {
		hasError = true
		formError.Password = "Old Password is required"
	}

	if formUser.NewPassword == "" {
		hasError = true
		formError.NewPassword = "New Password is required"
	}

	if formUser.PasswordConfirm == "" {
		hasError = true
		formError.PasswordConfirm = "Must confirm new password"
	}

	// Check confirm matches new password
	if formUser.NewPassword != formUser.PasswordConfirm {
		hasError = true
		formError.PasswordConfirm = "Passwords don't match"
	}

	// Check the provided old password is correct
	user, err := model.GetUserWithEmail(email)
	if err != nil {
		return echo.ErrNotFound
	}
	phash := model.HashPassword(formUser.Password)
	if user.Password != phash {
		hasError = true
		formError.Password = "Old Password is incorrect"
	}

	// Check new password matches constraints
	if errMsg := inf.CheckPasswordRequirements(formUser.NewPassword); errMsg != "" {
		hasError = true
		formError.NewPassword = errMsg
	}

	// Change Password
	if !hasError {
		err := user.ChangePassword(formUser.NewPassword)
		if err != nil {
			inf.LogError(err, "Failed to change password:")
			return c.Render(http.StatusOK, "status", inf.StatusPageData{
				Colour:          "red",
				MainMessage:     "Failed to change password",
				NextPageURL:     "/user/" + email,
				NextPageMessage: "Back to user page",
			})
		}
	}

	// Check Errors
	if hasError {
		session.Values["error"] = formError
		err = session.Save(c.Request(), c.Response())
		if err != nil {
			inf.LogError(err, "Failed to save form-data:")
			return echo.ErrInternalServerError
		}
		c.Redirect(http.StatusMovedPermanently, "/user/"+email+"/chgpassword")
	}

	// Log user out
	userData, err := front.CookieStore.Get(c.Request(), front.COOKIE_TOKEN_DATA)
	if err != nil {
		inf.LogError(err, "Failed to get token cookie:")
		return echo.ErrInternalServerError
	}
	userData.Options.MaxAge = -1

	err = userData.Save(c.Request(), c.Response())
	if err != nil {
		inf.LogError(err, "Failed to save token cookie:")
		return echo.ErrInternalServerError
	}

	return c.Render(http.StatusOK, "status", inf.StatusPageData{
		Colour:          "green",
		MainMessage:     "Successfully changed password",
		NextPageURL:     "/login",
		NextPageMessage: "Continue to login",
	})
}
