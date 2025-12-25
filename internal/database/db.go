package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Database struct {
	conn *sql.DB
}

func NewDatabase(path string) (*Database, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("Cant open DB: %w", err)
	}

	if err := db.Ping; err != nil {
		return nil, fmt.Errorf("Cant conn to DB: %w", err)
	}

	db.SetMaxOpenConns(1)
	db.SetConnMaxIdleTime(1)

	query := `
	CREATE TABLE IF NOT EXISTS tasks(
	ID INT PRIMARY KEY AUTO_INCREMENT,
	title TEXT NOT NULL
	description TEXT,
	status TEXT NOT NULL DEFAULT 'pending'
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	completed_TIMESTAMP
	)
	`

	_, err = db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("Cant create table: %w", err)
	}
	return &Database{conn: db}, nil
}
