package inf

import (
	"encoding/gob"
	"wishlist/src/inf/model"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
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
	gob.Register(UserFormError{})
	gob.Register(FormUser{})
	gob.Register(ItemFormError{})
	gob.Register(FormItem{})
	gob.Register(model.Token{})
}
