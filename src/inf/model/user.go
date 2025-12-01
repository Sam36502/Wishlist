package model

import (
	"fmt"
)

type User struct {
	ID       uint64
	Email    string
	Password string
	Name     string
}

//
//	Getters / Setters (Database Functions)
//

// Returns a list of all users
func GetAllUsers() ([]*User, error) {
	if g_database == nil {
		err := ConnectionInvalidError("No open connection")
		return nil, StackError(err, "Failed to insert new user")
	}

	rows, err := g_database.Query("SELECT * FROM tbl_user")
	if err != nil {
		return nil, StackError(err, "Failed to get all users:")
	}
	defer rows.Close()

	userArr := make([]*User, 0)
	for rows.Next() {
		parsedUser := User{}
		err = rows.Scan(&parsedUser.ID, &parsedUser.Email, &parsedUser.Password, &parsedUser.Name)
		if err != nil {
			return nil, StackError(err, "Parsing Failed:")
		}
		userArr = append(userArr, &parsedUser)
	}

	return userArr, nil
}

// Gets a user by their ID
func GetUserWithID(id uint64) (*User, error) {
	if g_database == nil {
		err := ConnectionInvalidError("No open connection")
		return nil, StackError(err, "Failed to insert new user")
	}

	rows, err := g_database.Query("SELECT * FROM tbl_user WHERE id_user = ?", id)
	if err != nil {
		return nil, StackError(err, "Query Failed:")
	}
	defer rows.Close()

	parsedUser := User{}
	rows.Next()
	err = rows.Scan(&parsedUser.ID, &parsedUser.Email, &parsedUser.Password, &parsedUser.Name)
	if err != nil {
		return nil, StackError(err, "Parsing Failed:")
	}

	return &parsedUser, nil
}

// Gets a user by their email
func GetUserWithEmail(email string) (*User, error) {
	if g_database == nil {
		err := ConnectionInvalidError("No open connection")
		return nil, StackError(err, "Failed to insert new user")
	}

	row := g_database.QueryRow("SELECT * FROM tbl_user WHERE email = ?", email)

	parsedUser := User{}
	err := row.Scan(&parsedUser.ID, &parsedUser.Email, &parsedUser.Password, &parsedUser.Name)
	if err != nil {
		return nil, StackError(err, "Failed to read user from DB")
	}

	return &parsedUser, nil
}

// Adds a new user to the database
//
// Returns a reference to the newly created user (or nil if it failed)
func InsertUser(user *User) (*User, error) {
	if g_database == nil {
		err := ConnectionInvalidError("No open connection")
		return nil, StackError(err, "Failed to insert new user")
	}

	// Check email isn't already in use
	_, err := GetUserWithEmail(user.Email)
	if err == nil {
		err = EmailExistsError(fmt.Sprintf("Email '%s' already registered.", user.Email))
		return nil, StackError(err, "Failed to insert new user")
	}

	phash := HashPassword(user.Password)
	_, err = g_database.Exec("INSERT INTO tbl_user (email, password, name) VALUES(?, ?, ?)", user.Email, phash, user.Name)
	if err != nil {
		return nil, StackError(err, "Failed to insert new user")
	}

	// Fetch newly created user to confirm
	newuser, err := GetUserWithEmail(user.Email)
	if err != nil {
		return nil, StackError(err, "Failed to retrieve newly created user")
	}

	return newuser, err
}

//
//	User Methods
//

// Permanently delete a user from the database
func (usr *User) Delete() error {
	if g_database == nil {
		err := ConnectionInvalidError("No open connection")
		return StackError(err, "Failed to insert new user")
	}

	_, err := g_database.Exec("DELETE FROM tbl_user WHERE id_user = ?", usr.ID)
	if err != nil {
		return StackError(err, "Query failed:")
	}
	return nil
}

// Updates the user's email in the database
func (usr *User) ChangeEmail(new_email string) error {
	if g_database == nil {
		err := ConnectionInvalidError("No open connection")
		return StackError(err, "Failed to insert new user")
	}

	// Check email isn't already in use
	_, err := GetUserWithEmail(usr.Email)
	if err == nil {
		return EmailExistsError(fmt.Sprintf("Email '%s' already registered.", usr.Email))
	}

	_, err = g_database.Exec("UPDATE tbl_user SET email = ? WHERE id_user = ?", new_email, usr.ID)
	if err != nil {
		return StackError(err, "Query failed:")
	}

	// Fetch user to confirm change
	newuser, err := GetUserWithID(usr.ID)
	if err != nil || newuser.Email != new_email {
		return StackError(err, "Failed to confirm user was updated:")
	}

	*usr = *newuser
	return nil
}

func (usr *User) ChangeName(new_name string) error {
	if g_database == nil {
		err := ConnectionInvalidError("No open connection")
		return StackError(err, "Failed to insert new user")
	}

	_, err := g_database.Exec("UPDATE tbl_user SET name = ? WHERE id_user = ?", new_name, usr.ID)
	if err != nil {
		return StackError(err, "Query failed:")
	}

	// Fetch user to confirm change
	newuser, err := GetUserWithID(usr.ID)
	if err != nil || newuser.Name != new_name {
		return StackError(err, "Failed to confirm user was updated:")
	}

	*usr = *newuser
	return nil
}

// Update the user's password in the database
func (usr *User) ChangePassword(new_pass string) error {
	if g_database == nil {
		err := ConnectionInvalidError("No open connection")
		return StackError(err, "Failed to insert new user")
	}

	phash := HashPassword(new_pass)
	_, err := g_database.Exec("UPDATE tbl_user SET password = ? WHERE id_user = ?", phash, usr.ID)
	if err != nil {
		return StackError(err, "Query failed:")
	}

	// Fetch newly created user to confirm
	newuser, err := GetUserWithID(usr.ID)
	if err != nil || newuser.Password != phash {
		return StackError(err, "Failed to retrieve newly created user:")
	}

	*usr = *newuser
	return nil
}
