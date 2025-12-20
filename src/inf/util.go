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

	"github.com/labstack/echo/v4"
)

const (
	PASS_MIN_LEN = 8 // Minimum allowed password length
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
		*err = StackError(*err, fmt.Sprintf("Server Panic ocurred on page '%s':", c.Request().URL.Path))

		pcs := make([]uintptr, 16)
		pc_count := runtime.Callers(1, pcs)
		for i := pc_count - 1; i > 0; i-- {
			f := runtime.FuncForPC(pcs[i])
			filename, ln := f.FileLine(pcs[i])
			fn := path.Base(filename)
			name := f.Name()
			*err = ErrorStack(fmt.Sprintf("%15s:%-4d | >>> %d: %s\n%s", fn, ln, i, name, *err))
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
	err = StackError(err, msg)
	LogMessage("ERROR", "An error occurred:\n"+err.Error())
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
