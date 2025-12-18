package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type UserDto struct {
	Uid            string
	Username       string
	HashedPassword string
}

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) GetUserId(username string) (string, error) {
	var userId string

	row := u.DB.QueryRow("SELECT id FROM users WHERE username = ?", username)
	if err := row.Scan(&userId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("user not found")
		}
		return "", fmt.Errorf("Error in fetching user: %s", err)
	}

	return userId, nil
}

func (u *UserModel) GetUserPassword(username string) (string, error) {
	var hashedPassword string

	row := u.DB.QueryRow("SELECT password_hash FROM users WHERE username = ?", username)
	if err := row.Scan(&hashedPassword); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("user not found")
		}
		return "", fmt.Errorf("Error in fetching user password: %s", err)
	}

	return hashedPassword, nil
}

func (u *UserModel) AddUser(user *UserDto) error {
	res, err := u.DB.Exec("INSERT INTO users (id, username, password_hash) VALUES (?, ?, ?)", user.Uid, user.Username, user.HashedPassword)
	if err != nil {
		return fmt.Errorf("Error in adding user: %s", err)
	}
	_, err = res.LastInsertId()
	if err != nil {
		return fmt.Errorf("Error in adding user: %s", err)
	}
	log.Printf("User added: %s", user.Username)
	return nil
}
