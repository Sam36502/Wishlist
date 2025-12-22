package model

import (
	"database/sql"
	"wishlist/src/inf"
)

type Item struct {
	ID             uint64
	Name           string
	Description    string
	Price          float32
	User           User
	ReservedByUser *User
	Status         Status
	Links          []Link
}

type Status struct {
	StatusID    uint64
	Name        string
	Description string
}

const (
	STATUS_ID_AVAILABLE = 1
	STATUS_ID_RESERVED  = 2
	STATUS_ID_RECEIVED  = 3
)

// Gets all the items in the list for a user
func GetAllItems(userID uint64) ([]*Item, error) {
	if g_database == nil {
		err := inf.ConnectionInvalidError("No open connection")
		return nil, inf.StackError(err, "Failed to fetch user list")
	}

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
		return nil, inf.StackError(err, "Query Failed:")
	}
	defer rows.Close()

	// Parse Items
	itemArr := make([]*Item, 0)
	var reservedByID sql.NullInt64
	for rows.Next() {
		parsedItem := Item{}
		err = rows.Scan(
			&parsedItem.ID,
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
			return nil, inf.StackError(err, "Parsing Failed:")
		}

		// Get Reserved-By User if Present
		if reservedByID.Valid {
			rows, err := g_database.Query("SELECT "+
				"id_user, email, name "+
				"FROM tbl_user WHERE id_user = ?", reservedByID.Int64)
			if err != nil {
				return nil, inf.StackError(err, "Query Failed:")
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
				return nil, inf.StackError(err, "Parsing Failed:")
			}
		} else {
			parsedItem.ReservedByUser = nil
		}

		// Get Links
		linkRows, err := g_database.Query("SELECT id_link, text, hyperlink FROM tbl_link WHERE item_id = ?", parsedItem.ID)
		if err != nil {
			return nil, inf.StackError(err, "Query Failed:")
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
				return nil, inf.StackError(err, "Parsing Failed:")
			}
			parsedItem.Links = append(parsedItem.Links, parsedLink)
		}

		itemArr = append(itemArr, &parsedItem)
	}

	return itemArr, nil
}

// Gets a single item from the database by its ID
func GetItemWithID(id uint64) (*Item, error) {
	if g_database == nil {
		err := inf.ConnectionInvalidError("No open connection")
		return nil, inf.StackError(err, "Failed to insert new user")
	}

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
		return nil, inf.StackError(err, "Query Failed:")
	}
	defer rows.Close()

	// Parse Items
	parsedItem := Item{}
	var reservedByID sql.NullInt64
	rows.Next()
	err = rows.Scan(
		&parsedItem.ID,
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
		return nil, inf.StackError(err, "Parsing Failed:")
	}

	// Get Reserved-By User if Present
	if reservedByID.Valid {
		rows, err := g_database.Query("SELECT "+
			"id_user, email, name "+
			"FROM tbl_user WHERE id_user = ?", reservedByID.Int64)
		if err != nil {
			return nil, inf.StackError(err, "Query Failed:")
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
			return nil, inf.StackError(err, "Parsing Failed:")
		}
	} else {
		parsedItem.ReservedByUser = nil
	}

	// Get Links
	linkRows, err := g_database.Query("SELECT id_link, text, hyperlink FROM tbl_link WHERE item_id = ?", parsedItem.ID)
	if err != nil {
		return nil, inf.StackError(err, "Failed to get Item Links:")
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
			return nil, inf.StackError(err, "Failed parse Item Links:")
		}
		parsedItem.Links = append(parsedItem.Links, parsedLink)
	}

	return &parsedItem, nil
}

// Inserts an item into the database
func InsertItem(item *Item) (*Item, error) {
	if g_database == nil {
		err := inf.ConnectionInvalidError("No open connection")
		return nil, inf.StackError(err, "Failed to insert new item:")
	}

	// Insert Item
	res, err := g_database.Exec("INSERT INTO tbl_item (name, desc, price, status_id, user_id) VALUES (?, ?, ?, ?, ?)", item.Name, item.Description, item.Price, item.Status.StatusID, item.User.ID)
	if err != nil {
		return nil, inf.StackError(err, "Query failed inserting item:")
	}

	// Get previous insert ID
	id64, err := res.LastInsertId()
	if err != nil {
		return nil, inf.StackError(err, "Failed to get new item ID:")
	}
	id := uint64(id64)

	// Insert Links
	for _, link := range item.Links {
		_, err = g_database.Query("INSERT INTO tbl_link (text, hyperlink, item_id) VALUES (?, ?, ?)", link.Text, link.URL, id) // Use Last insert ID as itemID
		if err != nil {
			return nil, inf.StackError(err, "Query failed adding links")
		}
	}

	new_item, err := GetItemWithID(id)
	if err != nil {
		return nil, inf.StackError(err, "Failed to verify new item was created successfully:")
	}

	*item = *new_item
	return new_item, nil
}

// Changes an Item's Name
func (itm *Item) ChangeName(name string) error {
	if g_database == nil {
		err := inf.ConnectionInvalidError("No open connection")
		return inf.StackError(err, "Failed to update item name")
	}

	_, err := g_database.Exec("UPDATE tbl_item SET name = ? WHERE id_item = ?;", name, itm.ID)
	if err != nil {
		return inf.StackError(err, "Failed to update item name")
	}

	return nil
}

// Changes an Item's Description
func (itm *Item) ChangeDescription(desc string) error {
	if g_database == nil {
		err := inf.ConnectionInvalidError("No open connection")
		return inf.StackError(err, "Failed to update item description")
	}

	_, err := g_database.Exec("UPDATE tbl_item SET desc = ? WHERE id_item = ?;", desc, itm.ID)
	if err != nil {
		return inf.StackError(err, "Failed to update item description:")
	}

	return nil
}

// Changes an Item's Status
func (itm *Item) ChangeStatus(status_id int) error {
	if g_database == nil {
		err := inf.ConnectionInvalidError("No open connection")
		return inf.StackError(err, "Failed to update item status:")
	}

	_, err := g_database.Exec("UPDATE tbl_item SET status_id = ? WHERE id_item = ?;", status_id, itm.ID)
	if err != nil {
		return inf.StackError(err, "Failed to update item status:")
	}

	return nil
}

// Changes an Item's Price
func (itm *Item) ChangePrice(price float32) error {
	if g_database == nil {
		err := inf.ConnectionInvalidError("No open connection")
		return inf.StackError(err, "Failed to update item price:")
	}

	_, err := g_database.Exec("UPDATE tbl_item SET price = ? WHERE id_item = ?;", int32(price*100), itm.ID)
	if err != nil {
		return inf.StackError(err, "Failed to update item price:")
	}

	return nil
}

// Changes which user has reserved an item
func (itm *Item) ChangeReservingUser(user *User) error {
	if g_database == nil {
		err := inf.ConnectionInvalidError("No open connection")
		return inf.StackError(err, "Failed to update item reservation")
	}

	if user == nil {
		_, err := g_database.Exec("UPDATE tbl_item SET reserved_by_user_id = NULL WHERE id_item = ?;", itm.ID)
		if err != nil {
			return inf.StackError(err, "Failed to update item reservation")
		}
	} else {
		_, err := g_database.Exec("UPDATE tbl_item SET reserved_by_user_id = ? WHERE id_item = ?;", user.ID, itm.ID)
		if err != nil {
			return inf.StackError(err, "Failed to update item reservation")
		}
	}

	return nil
}

// Change an item's links
func (itm *Item) ChangeLinks(links []Link) error {
	if g_database == nil {
		err := inf.ConnectionInvalidError("No open connection")
		return inf.StackError(err, "Failed to update item reservation")
	}

	// Clear all existing links
	_, err := g_database.Exec("DELETE FROM tbl_link WHERE item_id = ?;", itm.ID)
	if err != nil {
		return inf.StackError(err, "Failed to remove existing item links")
	}

	// Insert new links
	for _, lnk := range links {
		_, err := g_database.Exec("INSERT INTO tbl_link (text, hyperlink, item_id) VALUES (?, ?, ?)", lnk.Text, lnk.URL, itm.ID)
		if err != nil {
			return inf.StackError(err, "Failed to insert new item links")
		}
	}

	return nil
}

// Permanently delete an item from the database
func (itm *Item) Delete() error {
	if g_database == nil {
		err := inf.ConnectionInvalidError("No open connection")
		return inf.StackError(err, "Failed to insert new user")
	}

	_, err := g_database.Query("DELETE FROM tbl_item WHERE id_item = ?", itm.ID)
	if err != nil {
		return inf.StackError(err, "Query failed:")
	}
	return nil
}
