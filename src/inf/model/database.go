package model

import (
	"crypto/sha512"
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var (
	FILENAME   string  = os.Getenv("WISHLIST_DB_FILENAME")
	g_database *sql.DB = nil
)

// Open connection to the database
func ConnectDB() (error, string) {
	var err error
	g_database, err = sql.Open("sqlite", FILENAME)
	if err != nil || g_database == nil {
		return StackError(err, "Failed to open Database connection:"), ""
	}

	EnsureSchema(g_database)

	row := g_database.QueryRow("SELECT VALUE FROM META_INFO WHERE FIELD = 'schema_version';")
	var db_version string
	err = row.Scan(&db_version)
	if err != nil {
		return StackError(err, "Error: Failed to read database meta info:"), ""
	}

	return nil, db_version
}

// Disconnect from the Database
func DisconnectDB() error {
	return g_database.Close()
}

// Pings database and returns true if it's online and can be connected to; otherwise false
func IsDatabaseOnline() bool {
	err := g_database.Ping()
	return err == nil
}

// Converts an input password string into a hex hash string
func HashPassword(pwd string) string {
	return fmt.Sprintf("%x", sha512.Sum512(([]byte)(pwd)))
}

// Returns a list of users with the provided substring in their email or name
func SearchUsersByNameOrEmail(name string) ([]*User, error) {
	if g_database == nil {
		err := ConnectionInvalidError("No open connection")
		return nil, StackError(err, "Failed to search users:")
	}

	wildcardName := "%" + name + "%"
	rows, err := g_database.Query("SELECT * FROM tbl_user WHERE LOWER(email) LIKE ? OR LOWER(name) LIKE ?", wildcardName, wildcardName)
	if err != nil {
		return nil, StackError(err, "Failed to search users:")
	}
	defer rows.Close()

	userArr := make([]*User, 0)
	for rows.Next() {
		parsedUser := User{}
		err = rows.Scan(&parsedUser.ID, &parsedUser.Email, &parsedUser.Password, &parsedUser.Name)
		if err != nil {
			return nil, StackError(err, "Failed to search users:")
		}
		userArr = append(userArr, &parsedUser)
	}
	return userArr, nil
}
