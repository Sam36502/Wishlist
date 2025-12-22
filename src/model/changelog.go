package model

import (
	"wishlist/src/inf"
)

type ChangelogEntry struct {
	ID              uint64
	Version         string
	PublicationDate string
	Title           string
	Markdown        string
	Rendered        *string
}

func (che ChangelogEntry) GetMarkdown() string {
	return che.Markdown
}

func (che ChangelogEntry) GetHTML() *string {
	return che.Rendered
}

func (che ChangelogEntry) WriteHTML(html string) error {
	if g_database == nil {
		err := inf.ConnectionInvalidError("No open connection")
		return inf.StackError(err, "Failed to write rendered HTML")
	}

	_, err := g_database.Exec("UPDATE META_CHANGELOG SET RENDERED = ? WHERE ID = ?;",
		html, che.ID,
	)
	if err != nil {
		return inf.StackError(err, "Failed to write rendered HTML")
	}

	return nil
}

func GetChangelogEntries() ([]ChangelogEntry, error) {
	if g_database == nil {
		err := inf.ConnectionInvalidError("No open connection")
		return nil, inf.StackError(err, "Failed to retrieve changelog entries")
	}

	// Get All Entries
	rows, err := g_database.Query("SELECT * FROM META_CHANGELOG " +
		"WHERE PUB_DATE <= CURRENT_DATE " +
		"ORDER BY PUB_DATE DESC;",
	)
	if err != nil {
		return nil, inf.StackError(err, "Query Failed:")
	}
	defer rows.Close()

	// Parse Entries
	entries := make([]ChangelogEntry, 0)
	for rows.Next() {
		entry := ChangelogEntry{}
		err = rows.Scan(
			&entry.ID,
			&entry.PublicationDate,
			&entry.Version,
			&entry.Title,
			&entry.Markdown,
			&entry.Rendered,
		)
		if err != nil {
			return nil, inf.StackError(err, "Parsing Failed:")
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
