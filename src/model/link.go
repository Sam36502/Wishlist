package model

type Link struct {
	LinkID uint64 `json:"link_id"`
	Text   string `json:"text"`
	URL    string `json:"hyperlink"`
}
