package inf

import (
	"fmt"
	"net/http"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
	"wishlist/src/inf/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

const (
	PASS_MIN_LEN = 8 // Minimum allowed password length

	STATUS_AVAILABLE = 1
	STATUS_RESERVED  = 2
	STATUS_RECEIVED  = 3
)

// IsEmailValid checks if the email provided is valid by regex.
// Courtesy of Stackoverflow user "icza"
func IsEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

// Logs an error message instead of crashing
func RecoverPanic(c echo.Context, err *error) {
	rec := recover()

	if rec != nil {
		var ok bool
		if *err, ok = rec.(error); !ok {
			*err = fmt.Errorf("Recovery object: %v", rec)
		}
		*err = model.StackError(*err, fmt.Sprintf("Server Panic ocurred on page '%s':", c.Request().URL.Path))

		pcs := make([]uintptr, 16)
		pc_count := runtime.Callers(1, pcs)
		for i := pc_count - 1; i > 0; i-- {
			f := runtime.FuncForPC(pcs[i])
			filename, ln := f.FileLine(pcs[i])
			fn := path.Base(filename)
			name := f.Name()
			*err = model.ErrorStack(fmt.Sprintf("%15s:%-4d | >>> %d: %s\n%s", fn, ln, i, name, *err))
		}
		LogMessage("ERROR", "Server panic occurred:\n"+(*err).Error())

		c.Render(http.StatusOK, "error", ErrorPageContent{
			ErrorCode: fmt.Sprintf("Error %v - %v", http.StatusInternalServerError, "Internal Server Error"),
			ErrorDesc: "A fatal error has occurred on the server. Please try again later.",
		})
	}
}

// Gets the current date & time in a standard format for logging
func GetTimestamp() string {
	return time.Now().Format(time.DateTime)
}

// Write a message to the console with a timestamp
func LogMessage(level, msg string) {
	fmt.Printf(" %s [%s] %s\n", GetTimestamp(), strings.ToUpper(level), msg)
}

// Prints the error message for a stacked error
func LogError(err error, msg string) {
	err = model.StackError(err, msg)
	LogMessage("ERROR", "An error occurred:\n"+err.Error())
}

// Gets and parses the currently logged in user and token from the cookie
func GetLoggedInUser(c echo.Context) (model.User, model.Token, error) {
	tokenData, err := CookieStore.Get(c.Request(), COOKIE_TOKEN_DATA)
	if err != nil {
		LogError(err, "Failed to read token cookie:")

		if c.Path() != "/logout" {
			c.Redirect(302, "/logout")
		}
		return model.User{}, model.Token{}, nil
	}

	tokenInterface, exists := tokenData.Values[COOKIE_TOKEN_DATA]
	if !exists {
		return model.User{}, model.Token{}, model.StackError(nil, "No user currently logged in")
	}

	cookieToken, ok := tokenInterface.(model.Token)
	if !ok {
		LogError(nil, "Couldn't convert cookie user. Deleting Cookie...")
		tokenData.Options.MaxAge = -1
		err = tokenData.Save(c.Request(), c.Response())
		if err != nil {
			LogError(err, "Failed to delete cookie:")
		}
		return model.User{}, model.Token{}, model.StackError(nil, "Failed to convert cookie token")
	}

	// Parse email from JWT
	jwtParser := jwt.Parser{
		ValidMethods:         []string{},
		UseJSONNumber:        false,
		SkipClaimsValidation: true,
	}
	token, _, err := jwtParser.ParseUnverified(cookieToken.Token, jwt.MapClaims{})
	if err != nil {
		return model.User{}, model.Token{}, model.StackError(err, "failed to parse JWT Claims:")
	}
	jwtClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return model.User{}, model.Token{}, model.StackError(nil, "failed to convert cookie claims to token claims")
	}
	emailIfc, exists := jwtClaims["email"]
	if !exists {
		return model.User{}, model.Token{}, model.StackError(nil, "invalid JWT Claims parsed; no 'email' field present")
	}
	email, ok := emailIfc.(string)
	if !ok {
		return model.User{}, model.Token{}, model.StackError(nil, "invalid JWT Claims parsed; email is not a string")
	}

	user, err := model.GetUserWithEmail(email)
	if err != nil {
		return model.User{}, model.Token{}, model.StackError(err, "failed to retrieve logged-in user from API:")
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

// Checks if a given password matches some basic security criteria
// If it matches, this returns an empty string
// otherwise, it returns a string with an error message
// e.g.: "Password must be at least 8 characters"
func CheckPasswordRequirements(password string) string {

	if len(password) < PASS_MIN_LEN {
		return "Password must be at least " + strconv.Itoa(PASS_MIN_LEN) + " characters"
	}

	// TODO: More rigourous constraints

	return ""
}
