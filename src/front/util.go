package front

//
//		Front-End Utility Functions
//
//	Collects common functionalities required
//	for the user interface.
//

import (
	"strconv"
	"wishlist/src/inf"
	"wishlist/src/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// Gets and parses the currently logged in user and token from the cookie
func GetLoggedInUser(c echo.Context) (model.User, Token, error) {
	tokenData, err := CookieStore.Get(c.Request(), COOKIE_TOKEN_DATA)
	if err != nil {
		inf.LogError(err, "Failed to read token cookie:")

		if c.Path() != "/logout" {
			c.Redirect(302, "/logout")
		}

		return model.User{}, Token{}, nil
	}

	tokenInterface, exists := tokenData.Values[COOKIE_TOKEN_DATA]
	if !exists {
		return model.User{}, Token{}, inf.StackError(nil, "No user currently logged in")
	}

	cookieToken, ok := tokenInterface.(Token)
	if !ok {
		inf.LogError(nil, "Couldn't convert cookie user. Deleting Cookie...")
		tokenData.Options.MaxAge = -1
		err = tokenData.Save(c.Request(), c.Response())
		if err != nil {
			inf.LogError(err, "Failed to delete cookie:")
		}
		return model.User{}, Token{}, inf.StackError(nil, "Failed to convert cookie token")
	}

	// Parse email from JWT
	jwtParser := jwt.Parser{
		ValidMethods:         []string{},
		UseJSONNumber:        false,
		SkipClaimsValidation: true,
	}
	token, _, err := jwtParser.ParseUnverified(cookieToken.Token, jwt.MapClaims{})
	if err != nil {
		return model.User{}, Token{}, inf.StackError(err, "failed to parse JWT Claims:")
	}
	jwtClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return model.User{}, Token{}, inf.StackError(nil, "failed to convert cookie claims to token claims")
	}
	emailIfc, exists := jwtClaims["email"]
	if !exists {
		return model.User{}, Token{}, inf.StackError(nil, "invalid JWT Claims parsed; no 'email' field present")
	}
	email, ok := emailIfc.(string)
	if !ok {
		return model.User{}, Token{}, inf.StackError(nil, "invalid JWT Claims parsed; email is not a string")
	}

	user, err := model.GetUserWithEmail(email)
	if err != nil {
		return model.User{}, Token{}, inf.StackError(err, "failed to retrieve logged-in user from API:")
	}

	return *user, cookieToken, nil
}

// Gets all item data based on the "item_id" path parameter
func GetItemFromPath(c echo.Context) (model.Item, error) {
	var item_id uint64
	var err error
	if item_id_str := c.Param("item_id"); item_id_str == "" {
		return model.Item{}, echo.ErrNotFound
	} else {
		item_id, err = strconv.ParseUint(item_id_str, 10, 64)
		if err != nil {
			return model.Item{}, echo.ErrNotFound
		}
	}
	item, err := model.GetItemWithID(item_id)
	if err != nil {
		return model.Item{}, echo.ErrNotFound
	}

	if item == nil {
		return model.Item{}, nil
	}

	return *item, nil
}

func GetStatusColour(s model.Status) string {
	switch s.StatusID {

	default:
		fallthrough
	case 1:
		return "near-white"
	case 2:
		return "gold"
	case 3:
		return "green"

	}
}
