package model

import (
	"crypto/sha512"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "modernc.org/sqlite"
)

const (
	DEFAULT_DB_FILENAME     = "data/wishlist.db"
	DEFAULT_SCHEMA_FILENAME = "data/db_schema.sql"
	SCHEMA_VERSION          = "1.1.0"
)

var (
	DB_FILENAME     string  = DEFAULT_DB_FILENAME
	SCHEMA_FILENAME string  = DEFAULT_SCHEMA_FILENAME
	g_database      *sql.DB = nil
)

// Open connection to the database
func ConnectDB() (error, string) {

	// Get environment variables
	var ok bool
	DB_FILENAME, ok = os.LookupEnv("WISHLIST_DB_FILENAME")
	if !ok {
		DB_FILENAME = DEFAULT_DB_FILENAME
	}

	SCHEMA_FILENAME, ok = os.LookupEnv("WISHLIST_SCHEMA_FILENAME")
	if !ok {
		SCHEMA_FILENAME = DEFAULT_SCHEMA_FILENAME
	}

	// Connect to the database
	var err error
	g_database, err = sql.Open("sqlite", DB_FILENAME)
	if err != nil || g_database == nil {
		return StackError(err, "Failed to open Database connection:"), ""
	}

	db_version, err := EnsureSchema(g_database)
	if err != nil {
		return StackError(err, "Failed to ensure schema:"), ""
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

// Checks if the DB schema matches the required version
//
// If the DB is missing entirely, it will create it using the script in
func EnsureSchema(db *sql.DB) (string, error) {
	if db == nil {
		err := fmt.Errorf("Database is nil")
		return "", StackError(err, "Tried to ensure schema of invalid database connection")
	}

	// Read the meta info table to check if the DB matches the expected schema
	row := db.QueryRow("SELECT VALUE FROM META_INFO WHERE FIELD = 'schema_version';")
	var db_version string
	err := row.Scan(&db_version)
	if err == nil && db_version != SCHEMA_VERSION {
		return "", StackError(nil, fmt.Sprintf(
			"Schema version doesn't match expected (expected '%s', got '%s')",
			SCHEMA_VERSION, db_version,
		))
	}

	if db_version == SCHEMA_VERSION {
		return db_version, nil
	}

	// DB probably doesn't exist; create from script
	fmt.Printf("Failed to read schema metadata; attempting to create DB from script '%s'...\n", SCHEMA_FILENAME)

	file_data, err := os.ReadFile(SCHEMA_FILENAME)
	if err != nil {
		return "", StackError(err, fmt.Sprintf("Failed to read DB Schema script '%s'", SCHEMA_FILENAME))
	}

	script := strings.ReplaceAll(string(file_data), "\n", "")
	line_nr := 0
	for stmt := range strings.SplitSeq(script, ";") {
		line_nr++

		if len(stmt) <= 0 {
			continue
		}

		stmt += ";"
		fmt.Printf("---> Executing '%s'...\n", stmt)
		_, err := db.Exec(stmt)
		if err != nil {
			return "", StackError(err, fmt.Sprintf("Failed while executing %s:%d:", SCHEMA_FILENAME, line_nr))
		}
	}

	// Add standard info to metadata table
	_, err = db.Exec("INSERT INTO META_INFO VALUES ('created_date', CURRENT_TIMESTAMP);")
	if err != nil {
		return "", StackError(err, "Failed to set metadata for schema creation")
	}

	db_version, err = MetaInfo_Get("schema_version")
	if err != nil {
		return "", StackError(err, "Failed to validate schema version metadata")
	}

	return db_version, nil
}

// Sets a value in the meta info table
func MetaInfo_Set(key, val string) error {
	errmsg := fmt.Sprintf("Failed to set meta field '%s' to value '%s':", key, val)

	if g_database == nil {
		err := fmt.Errorf("Database connection is not valid")
		return StackError(err, errmsg)
	}

	_, err := g_database.Exec("INSERT INTO META_INFO VALUES (?, ?) ON CONFLICT DO UPDATE SET VALUE = ?;", key, val, val)
	if err != nil {
		StackError(err, errmsg)
	}

	return nil
}

// Sets a value in the meta info table
func MetaInfo_Get(key string) (string, error) {
	errmsg := fmt.Sprintf("Failed to read meta field '%s':", key)

	if g_database == nil {
		err := fmt.Errorf("Database connection is not valid")
		return "", StackError(err, errmsg)
	}

	row := g_database.QueryRow("SELECT VALUE FROM META_INFO WHERE FIELD = ?", key)
	var val string
	err := row.Scan(&val)
	if err != nil {
		return "", StackError(err, errmsg)
	}

	return val, nil
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
