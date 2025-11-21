package model

type Item struct {
	ItemID         uint64  `json:"item_id"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	Price          float32 `json:"price"`
	User           User    `json:"user"`
	ReservedByUser *User   `json:"reserved_by"`
	Status         Status  `json:"status"`
	Links          []Link  `json:"links"`
}

type Link struct {
	LinkID uint64 `json:"link_id"`
	Text   string `json:"text"`
	URL    string `json:"hyperlink"`
}

type Status struct {
	StatusID    uint64 `json:"status_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
