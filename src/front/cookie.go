package front

import (
	"encoding/gob"
	"wishlist/src/inf"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// Cookie Config
const (
	COOKIE_TIMEOUT    = 30 * 24 * 60 * 60 // User data cookie expires after 30 days (same as API token)
	COOKIE_FORM_DATA  = "form-data"       // Cookie name for sending error information to form pages
	COOKIE_TOKEN_DATA = "token-data"      // Cookie name for the access token
)

var CookieStore *sessions.CookieStore

func InitCookieStore() {
	// Initialise Cookie Store
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	CookieStore = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	CookieStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   0,
		HttpOnly: false,
	}

	// Register Session Types
	gob.Register(inf.UserFormError{})
	gob.Register(inf.FormUser{})
	gob.Register(inf.ItemFormError{})
	gob.Register(inf.FormItem{})
	gob.Register(Token{})
}
