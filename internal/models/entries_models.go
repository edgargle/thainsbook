package models

import (
	"database/sql"
	"fmt"
	"log"
)

// TODO: Implement sequential ID for entries so that user can easily identify which for update and deletion
type EntryDto struct {
	Id        string
	Title     string
	Content   string
	EntryDate string
	UserId    string
}

type EntryRequest struct {
	Id        string `json:"Id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	EntryDate string `json:"entry_date"`
}

type EntryResponse struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	EntryDate string `json:"entry_date"`
	UpdatedAt string `json:"date_updated"`
	CreatedAt string `json:"date_created"`
}

type EntryModel struct {
	DB *sql.DB
}

func (m *EntryModel) AddEntry(entry *EntryDto) error {
	res, err := m.DB.Exec("INSERT INTO entries (id, user_id, title, content, entry_date) VALUES (?, ?, ?, ?, ?)", entry.Id, entry.UserId, entry.Title, entry.Content, entry.EntryDate)
	if err != nil {
		return fmt.Errorf("Error inserting entry: %s", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error with inserting entry: %v", id)
		return fmt.Errorf("Error inserting entry: %s", err)
	}
	return nil
}

func (m *EntryModel) GetEntriesByUser(userId string) ([]EntryResponse, error) {
	var entries []EntryResponse

	rows, err := m.DB.Query("SELECT entries.title, entries.content, entries.entry_date, entries.updated_at, entries.created_at FROM entries JOIN users ON entries.user_id = users.id WHERE users.id = ?", userId)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Unable to fetch entries")
	}
	defer rows.Close()

	for rows.Next() {
		var entry EntryResponse
		if err := rows.Scan(&entry.Title, &entry.Content, &entry.EntryDate, &entry.UpdatedAt, &entry.CreatedAt); err != nil {
			log.Println(err)
			return nil, fmt.Errorf("Unable to fetch entries")
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Unable to fetch entries")
	}
	return entries, nil
}
