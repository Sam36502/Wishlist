package inf

// DB Config
const (
	WISHLIST_DB_FILE = "wishlist.db"
)

// Cookie Config
const (
	COOKIE_TIMEOUT    = 30 * 24 * 60 * 60 // User data cookie expires after 30 days (same as API token)
	COOKIE_FORM_DATA  = "form-data"       // Cookie name for sending error information to form pages
	COOKIE_TOKEN_DATA = "token-data"      // Cookie name for the access token
)
