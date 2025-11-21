/***********************************************************

		DATABASE FUNCTIONS

		It's pretty monolithic, because I don't like
		DB integration code and wanted it to all be
		hidden away in the corner.

		I tried to keep it somewhat maintainable by
		leaving copious navigation comments, but it'd
		be better if properly cleaned up.

		If I revisit this and want to clean up the
		code, I'd suggest putting this in its own
		module and grouping the functions by object
		(item, user, etc.)

		-2021

		I'm switching to using a simple SQLite file
		and shrinking the whole project down to one
		container, to make it easier to maintain.

		Here comes that rewrite I promised 4 years ago...

		-2025

************************************************************/

package model

import (
	"crypto/sha512"
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var g_database *sql.DB = nil
var timerKillChannel = make(chan bool)

var (
	FILENAME = os.Getenv("WISHLIST_DB_FILENAME")
)

func ConnectDB() error {
	var err error
	g_database, err = sql.Open("sqlite", FILENAME)
	if err != nil || g_database == nil {
		return fmt.Errorf("Failed to open Database connection:\n  %v", err)
	}

	EnsureSchema(g_database)

	row := g_database.QueryRow("SELECT VALUE FROM META_INFO WHERE FIELD = 'schema_version';")
	var db_version string
	err = row.Scan(&db_version)
	if err != nil {
		return fmt.Errorf("Error: Failed to read database meta info:\n  %v", err)
	}

	fmt.Printf("[  DB  ] Successfully connected to the Database!\n")
	fmt.Printf("[  DB  ]   Schema Version: %s\n", db_version)

	return nil
}

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

/// USER FUNCTIONS: ///
//	GetAllUsers() -> []User, err
//	GetUserWithID( ID ) -> User, err
//	GetUserWithEmail( Email ) -> User, err
//	InsertUser( User ) -> err
//	UpdateUser( User ) -> err
//	DeleteUser( ID ) -> err

// Returns a list of all users
func GetAllUsers() ([]*User, error) {
	if g_database == nil {
		return nil, fmt.Errorf("Query Failed: Database connection is invalid!")
	}

	rows, err := g_database.Query("SELECT * FROM tbl_user")
	if err != nil {
		return nil, fmt.Errorf("Failed to get all users:\n  %v", err)
	}
	defer rows.Close()

	userArr := make([]*User, 0)
	for rows.Next() {
		parsedUser := User{}
		err = rows.Scan(&parsedUser.ID, &parsedUser.Email, &parsedUser.Password, &parsedUser.Name)
		if err != nil {
			fmt.Println(" [ERROR] Parsing Failed:", err)
			return nil, err
		}
		userArr = append(userArr, &parsedUser)
	}

	return userArr, nil
}

// Gets a user by their ID
func GetUserWithID(id uint64) (*User, error) {
	rows, err := g_database.Query("SELECT * FROM tbl_user WHERE id_user = ?", id)
	if err != nil {
		fmt.Println(" [ERROR] Query Failed:", err)
		return nil, err
	}
	defer rows.Close()

	parsedUser := User{}
	rows.Next()
	err = rows.Scan(&parsedUser.ID, &parsedUser.Email, &parsedUser.Password, &parsedUser.Name)
	if err != nil {
		fmt.Println(" [ERROR] Parsing Failed:", err)
		return nil, err
	}

	return &parsedUser, nil
}

// Gets a user by their email
func GetUserWithEmail(email string) (*User, error) {
	rows, err := g_database.Query("SELECT * FROM tbl_user WHERE email = ?", email)
	if err != nil {
		fmt.Println(" [ERROR] Query Failed:", err)
		return nil, err
	}
	defer rows.Close()

	parsedUser := User{}
	rows.Next()
	err = rows.Scan(&parsedUser.ID, &parsedUser.Email, &parsedUser.Password, &parsedUser.Name)
	if err != nil {
		fmt.Println(" [ERROR] Parsing Failed:", err)
		return nil, err
	}

	return &parsedUser, nil
}

// Adds a new user to the database
func InsertUser(user *User) error {
	// Check email isn't already in use
	_, err := GetUserWithEmail(user.Email)
	if err == nil {
		fmt.Println(" [ERROR] Email '" + user.Email + "' already registered.")
		return EmailExistsError(user.Email)
	}

	phash := HashPassword(user.Password)
	_, err = g_database.Query("INSERT INTO tbl_user (email, password, name) VALUES(?, ?, ?)", user.Email, phash, user.Name)
	if err != nil {
		fmt.Println(" [ERROR] Query failed:", err)
		return err
	}

	return nil
}

// Changes the details of a user
func UpdateUser(user *User) error {
	// Add arguments to the query if they aren't empty
	queryStr := "UPDATE tbl_user SET "
	argArr := make([]interface{}, 0)

	noArgs := true
	if user.Email != "" {
		noArgs = false
		queryStr += "email = ? ,"
		argArr = append(argArr, user.Email)

		// Check email isn't already in use
		_, err := GetUserWithEmail(user.Email)
		if err == nil {
			return EmailExistsError(fmt.Sprintf("Email '%s' is already registered", user.Email))
		}
	}

	if user.Password != "" {
		noArgs = false
		queryStr += "password = ? ,"
		argArr = append(argArr, user.Password)
	}

	if user.Name != "" {
		noArgs = false
		queryStr += "name = ? ,"
		argArr = append(argArr, user.Name)
	}

	if noArgs {
		fmt.Println("[INFO] User '" + user.Email + "' tried to update their user info with no arguments.")
		return nil
	}

	queryStr = queryStr[:len(queryStr)-1]
	queryStr += "WHERE id_user = ?"
	argArr = append(argArr, user.ID)

	_, err := g_database.Exec(queryStr, argArr...)
	if err != nil {
		fmt.Println(" [ERROR] Query failed:", err)
		return err
	}

	return nil
}

// Permanently delete a user from the database
func DeleteUser(id uint64) error {
	_, err := g_database.Query("DELETE FROM tbl_user WHERE id_user = ?", id)
	if err != nil {
		fmt.Println(" [ERROR] Query failed:", err)
		return err
	}
	return nil
}

/// ITEM FUNCTIONS ///
//	GetAllItems(userID) -> []Item & error
//  GetItemWithID(itemID) -> Item & error
//	InsertItem(Item) -> error
//  UpdateItem(Item) -> error
//  DeleteItem(itemID) -> error

// Gets all the items in the list for a user
func GetAllItems(userID uint64) ([]*Item, error) {

	// Get All Items
	rows, err := g_database.Query("SELECT "+
		"i.id_item, i.name, i.desc, i.price, i.reserved_by_user_id, "+
		"s.id_status, s.name, s.desc, "+
		"u.id_user, u.email, u.name "+
		"FROM tbl_item i "+
		"JOIN tbl_status s ON i.status_id = s.id_status "+
		"JOIN tbl_user u ON i.user_id = u.id_user "+
		"WHERE u.id_user = ? ORDER BY i.id_item", userID)
	if err != nil {
		fmt.Println(" [ERROR] Query Failed:", err)
		return nil, err
	}
	defer rows.Close()

	// Parse Items
	itemArr := make([]*Item, 0)
	var reservedByID sql.NullInt64
	for rows.Next() {
		parsedItem := Item{}
		err = rows.Scan(
			&parsedItem.ItemID,
			&parsedItem.Name,
			&parsedItem.Description,
			&parsedItem.Price,
			&reservedByID,
			&parsedItem.Status.StatusID,
			&parsedItem.Status.Name,
			&parsedItem.Status.Description,
			&parsedItem.User.ID,
			&parsedItem.User.Email,
			&parsedItem.User.Name,
		)
		if err != nil {
			fmt.Println(" [ERROR] Parsing Failed:", err)
			return nil, err
		}

		// Get Reserved-By User if Present
		if reservedByID.Valid {
			rows, err := g_database.Query("SELECT "+
				"id_user, email, name "+
				"FROM tbl_user WHERE id_user = ?", reservedByID.Int64)
			if err != nil {
				fmt.Println(" [ERROR] Query Failed:", err)
				return nil, err
			}
			defer rows.Close()
			rows.Next()
			parsedItem.ReservedByUser = &User{}
			err = rows.Scan(
				&parsedItem.ReservedByUser.ID,
				&parsedItem.ReservedByUser.Email,
				&parsedItem.ReservedByUser.Name,
			)
			if err != nil {
				fmt.Println(" [ERROR] Parsing Failed:", err)
				return nil, err
			}
		} else {
			parsedItem.ReservedByUser = nil
		}

		// Get Links
		linkRows, err := g_database.Query("SELECT id_link, text, hyperlink FROM tbl_link WHERE item_id = ?", parsedItem.ItemID)
		if err != nil {
			fmt.Println(" [ERROR] Query Failed:", err)
			return nil, err
		}
		defer linkRows.Close()

		// Parse Links
		for linkRows.Next() {
			parsedLink := Link{}
			err = linkRows.Scan(
				&parsedLink.LinkID,
				&parsedLink.Text,
				&parsedLink.URL,
			)
			if err != nil {
				fmt.Println(" [ERROR] Parsing Failed:", err)
				return nil, err
			}
			parsedItem.Links = append(parsedItem.Links, parsedLink)
		}

		itemArr = append(itemArr, &parsedItem)
	}

	return itemArr, nil
}

// Gets a single item from the database by its ID
func GetItemWithID(id uint64) (*Item, error) {
	// Get All Items
	rows, err := g_database.Query("SELECT "+
		"i.id_item, i.name, i.desc, i.price, i.reserved_by_user_id, "+
		"s.id_status, s.name, s.desc, "+
		"u.id_user, u.email, u.name "+
		"FROM tbl_item i "+
		"JOIN tbl_status s ON i.status_id = s.id_status "+
		"JOIN tbl_user u ON i.user_id = u.id_user "+
		"WHERE i.id_item = ?", id)
	if err != nil {
		fmt.Println(" [ERROR] Query Failed:", err)
		return nil, err
	}
	defer rows.Close()

	// Parse Items
	parsedItem := Item{}
	var reservedByID sql.NullInt64
	rows.Next()
	err = rows.Scan(
		&parsedItem.ItemID,
		&parsedItem.Name,
		&parsedItem.Description,
		&parsedItem.Price,
		&reservedByID,
		&parsedItem.Status.StatusID,
		&parsedItem.Status.Name,
		&parsedItem.Status.Description,
		&parsedItem.User.ID,
		&parsedItem.User.Email,
		&parsedItem.User.Name,
	)
	if err != nil {
		fmt.Println(" [ERROR] Parsing Failed:", err)
		return nil, err
	}

	// Get Reserved-By User if Present
	if reservedByID.Valid {
		rows, err := g_database.Query("SELECT "+
			"id_user, email, name "+
			"FROM tbl_user WHERE id_user = ?", reservedByID.Int64)
		if err != nil {
			fmt.Println(" [ERROR] Query Failed:", err)
			return nil, err
		}
		defer rows.Close()
		rows.Next()
		parsedItem.ReservedByUser = &User{}
		err = rows.Scan(
			&parsedItem.ReservedByUser.ID,
			&parsedItem.ReservedByUser.Email,
			&parsedItem.ReservedByUser.Name,
		)
		if err != nil {
			fmt.Println(" [ERROR] Parsing Failed:", err)
			return nil, err
		}
	} else {
		parsedItem.ReservedByUser = nil
	}

	// Get Links
	linkRows, err := g_database.Query("SELECT id_link, text, hyperlink FROM tbl_link WHERE item_id = ?", parsedItem.ItemID)
	if err != nil {
		fmt.Println(" [ERROR] Query Failed:", err)
		return nil, err
	}
	defer linkRows.Close()

	// Parse Links
	for linkRows.Next() {
		parsedLink := Link{}
		err = linkRows.Scan(
			&parsedLink.LinkID,
			&parsedLink.Text,
			&parsedLink.URL,
		)
		if err != nil {
			fmt.Println(" [ERROR] Parsing Failed:", err)
			return nil, err
		}
		parsedItem.Links = append(parsedItem.Links, parsedLink)
	}

	return &parsedItem, nil
}

// Inserts an item into the database
func InsertItem(item *Item) error {

	// Insert Item
	res, err := g_database.Exec("INSERT INTO tbl_item (name, desc, price, status_id, user_id) VALUES (?, ?, ?, ?, ?)", item.Name, item.Description, item.Price, item.Status.StatusID, item.User.ID)
	if err != nil {
		fmt.Println(" [ERROR] Query failed inserting item.", err)
		return err
	}

	// Get previous insert ID
	id64, err := res.LastInsertId()
	if err != nil {
		fmt.Println(" [ERROR] Failed to get previous insert ID.", err)
		return err
	}
	id := uint64(id64)

	// Insert Links
	for _, link := range item.Links {
		_, err = g_database.Query("INSERT INTO tbl_link (text, hyperlink, item_id) VALUES (?, ?, ?)", link.Text, link.URL, id) // Use Last insert ID as itemID
		if err != nil {
			fmt.Println(" [ERROR] Query failed adding links", err)
			return err
		}
	}

	return nil
}

// Changes the details of an item
func UpdateItem(item *Item) error {
	// Add arguments to the query if they aren't empty
	queryStr := "UPDATE tbl_item SET "
	argArr := make([]interface{}, 0)

	noArgs := true
	if item.Name != "" {
		noArgs = false
		queryStr += "name = ? ,"
		argArr = append(argArr, item.Name)
	}

	if item.Description != "" {
		noArgs = false
		queryStr += "desc = ? ,"
		argArr = append(argArr, item.Description)
	}

	if item.Status != (Status{}) {
		noArgs = false
		queryStr += "status_id = ? ,"
		argArr = append(argArr, item.Status.StatusID)
	}

	if item.Price != -1 {
		noArgs = false
		queryStr += "price = ? ,"
		argArr = append(argArr, item.Price)
	}

	if item.ReservedByUser == nil {
		noArgs = false
		queryStr += "reserved_by_user_id = NULL ,"
	} else if *item.ReservedByUser != (User{}) {
		noArgs = false
		queryStr += "reserved_by_user_id = ? ,"
		argArr = append(argArr, item.ReservedByUser.ID)
	}

	queryStr = queryStr[:len(queryStr)-1]
	queryStr += "WHERE id_item = ?"
	argArr = append(argArr, item.ItemID)

	if !noArgs {
		// Update Item
		_, err := g_database.Query(queryStr, argArr...)
		if err != nil {
			fmt.Println(" [ERROR] Query failed:", err)
			return err
		}
	} else {
		fmt.Println("[INFO] User tried to update Item '" + item.Name + "' with no arguments.")
	}

	// Drop all links for this item and add the ones provided
	if item.Links != nil {
		_, err := g_database.Query("DELETE FROM tbl_link WHERE item_id = ?", item.ItemID)
		if err != nil {
			fmt.Println(" [ERROR] Query failed:", err)
			return err
		}

		for _, link := range item.Links {
			_, err = g_database.Query("INSERT INTO tbl_link (text, hyperlink, item_id) VALUES (?, ?, ?)", link.Text, link.URL, item.ItemID)
			if err != nil {
				fmt.Println(" [ERROR] Query failed.", err)
				return err
			}
		}
	}

	return nil
}

// Permanently delete an item from the database
func DeleteItem(id uint64) error {
	_, err := g_database.Query("DELETE FROM tbl_item WHERE id_item = ?", id)
	if err != nil {
		fmt.Println(" [ERROR] Query failed:", err)
		return err
	}
	return nil
}

/// MISC FUNCTIONS
// GetUsersByNameOrEmail(name) []User, error

// Returns a list of users with the provided substring in their email or name
func SearchUsersByNameOrEmail(name string) ([]*User, error) {
	wildcardName := "%" + name + "%"
	rows, err := g_database.Query("SELECT * FROM tbl_user WHERE LOWER(email) LIKE ? OR LOWER(name) LIKE ?", wildcardName, wildcardName)
	if err != nil {
		fmt.Println(" [ERROR] Query Failed:", err)
		return nil, err
	}
	defer rows.Close()

	userArr := make([]*User, 0)
	for rows.Next() {
		parsedUser := User{}
		err = rows.Scan(&parsedUser.ID, &parsedUser.Email, &parsedUser.Password, &parsedUser.Name)
		if err != nil {
			fmt.Println(" [ERROR] Parsing Failed:", err)
			return nil, err
		}
		userArr = append(userArr, &parsedUser)
	}
	return userArr, nil
}
