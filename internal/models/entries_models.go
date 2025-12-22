package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
)

// EntryDto - For database interactions
type EntryDto struct {
	Id        string
	Title     string
	Content   string
	EntryDate string
	UserId    string
}

// EntryRequest - For incoming entry objects
type EntryRequest struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	EntryDate string `json:"entry_date"`
}

// Entry Response - For outgoing entry objects
type EntryResponse struct {
	SeqId     string `json:"id"` // return as id to user, but is user-scoped id
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
	res, err := m.DB.Exec("INSERT INTO entries (id, seq_id, user_id, title, content, entry_date) VALUES (?, (SELECT COALESCE(MAX(seq_id), 0) + 1 FROM entries AS temp WHERE user_id = ?), ?, ?, ?, ?)", entry.Id, entry.UserId, entry.UserId, entry.Title, entry.Content, entry.EntryDate)
	if err != nil {
		return fmt.Errorf("Error inserting entry: %s", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error inserting entry: %s", err)
		return fmt.Errorf("Error inserting entry: %v", id)
	}
	return nil
}

func (m *EntryModel) GetEntry(id string, userId string) (*EntryResponse, error) {
	var entry EntryResponse

	row := m.DB.QueryRow("SELECT seq_id, title, content, entry_date, updated_at, created_at FROM entries WHERE seq_id = ? AND user_id = ?", id, userId)
	if err := row.Scan(&entry.SeqId, &entry.Title, &entry.Content, &entry.EntryDate, &entry.UpdatedAt, &entry.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("No entry found")
			return nil, fmt.Errorf("ID %v: no such entry", id)
		}
		return nil, fmt.Errorf("ID %d: %v", id, err)
	}
	return &entry, nil
}

func (m *EntryModel) GetEntriesByUser(userId string) ([]EntryResponse, error) {
	var entries []EntryResponse

	rows, err := m.DB.Query("SELECT entries.seq_id, entries.title, entries.content, entries.entry_date, entries.updated_at, entries.created_at FROM entries WHERE entries.user_id = ?", userId)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Unable to fetch entries")
	}
	defer rows.Close()

	for rows.Next() {
		var entry EntryResponse
		if err := rows.Scan(&entry.SeqId, &entry.Title, &entry.Content, &entry.EntryDate, &entry.UpdatedAt, &entry.CreatedAt); err != nil {
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

func (m *EntryModel) UpdateEntry(id string, entry *EntryDto) error {
	// 1. Prepare the building blocks
	query := "UPDATE entries SET "
	var args []interface{}
	var setClauses []string

	if entry.Title != "" {
		setClauses = append(setClauses, "title = ?")
		args = append(args, entry.Title)
	}

	if entry.Content != "" {
		setClauses = append(setClauses, "content = ?")
		args = append(args, entry.Content)
	}

	if entry.EntryDate != "" {
		setClauses = append(setClauses, "entry_date = ?")
		args = append(args, entry.EntryDate)
	}

	if len(setClauses) == 0 {
		return errors.New("no fields provided for update")
	}

	query += strings.Join(setClauses, ", ")
	query += " WHERE user_id = ? AND seq_id = ?"

	args = append(args, entry.UserId, id)

	result, err := m.DB.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("record not found or no changes made")
	}

	return nil
}

func (m *EntryModel) DeleteEntry(id string, userId string) error {
	res, err := m.DB.Exec("DELETE FROM entries WHERE seq_id = ? AND user_id = ?", id, userId)
	if err != nil {
		log.Printf("Error deleting entry: %s", err)
		return fmt.Errorf("error deleting entry: %v", id)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("record not found")
	}

	return nil
}
