package repository

import (
	"Blog/models"
	"database/sql"
	"errors"
)

func (a *SQLite) GetUserById(id int) (models.User, error) {
	var user models.User
	query := "SELECT id, name, password, is_admin FROM users WHERE id = ?"

	err := a.db.QueryRow(query, id).Scan(&user.Id, &user.Name, &user.Password, &user.IsAdmin)
	if errors.Is(err, sql.ErrNoRows) {
		return user, models.ErrUserNotFound
	}
	return user, err
}

func (a *SQLite) GetUserByName(username string) (models.User, error) {
	var user models.User
	query := "SELECT id, name, password, is_admin FROM users WHERE name = ?"
	err := a.db.QueryRow(query, username).Scan(&user.Id, &user.Name, &user.Password, &user.IsAdmin)
	if errors.Is(err, sql.ErrNoRows) {
		return user, models.ErrUserNotFound
	}
	return user, err
}

func (a *SQLite) AddUser(user *models.User) error {
	query := "INSERT INTO users (name, password, is_admin) VALUES (?, ?, ?)"
	_, err := a.db.Exec(query, user.Name, user.Password, user.IsAdmin)
	return err
}
