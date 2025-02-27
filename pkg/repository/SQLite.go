package repository

import (
	"database/sql"
)

type SQLite struct {
	db *sql.DB
}

func NewSQLite(driverName, source string) *SQLite {

	db, err := sql.Open(driverName, source)
	if err != nil {
		panic(err)
	}
	s := &SQLite{db: db}
	err = s.InitTables()
	if err != nil {
		panic(err)
	}
	return s
}

func (a *SQLite) InitTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		is_admin BOOLEAN DEFAULT 0
	);
	CREATE TABLE IF NOT EXISTS articles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		author_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		body TEXT NOT NULL,
		cover_url TEXT NOT NULL,
		is_posted BOOLEAN DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
	);
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		article_id INTEGER NOT NULL,
		author_name TEXT NOT NULL,
		text TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE
	);`
	_, err := a.db.Exec(query)
	return err
}
