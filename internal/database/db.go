package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/doub1educk/gotasker/internal/domain"
	_ "modernc.org/sqlite"
)

type Database struct {
	conn *sql.DB
}

func NewDatabase(path string) (*Database, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("cant open DB: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cant conn to DB: %w", err)
	}

	db.SetMaxOpenConns(1)
	db.SetConnMaxIdleTime(1)

	query := `
	CREATE TABLE IF NOT EXISTS tasks(
	ID INTEGER PRIMARY KEY AUTO_INCREMENT,
	title TEXT NOT NULL
	description TEXT,
	status TEXT NOT NULL DEFAULT 'pending'
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	completed_TIMESTAMP
	)
	`

	_, err = db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("cant create table: %w", err)
	}
	return &Database{conn: db}, nil
}

func (d *Database) CreateTask(title, description string) (int, error) {
	query := `
	INSERT INTO tasks(title,description,status)
	VALUES(?,?,'pending')
	`
	result, err := d.conn.Exec(query, title, description)
	if err != nil {
		return 0, fmt.Errorf("cant create a new task : %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("cant get a id: %w", err)
	}
	return int(id), nil
}

func (d *Database) CreateTaskWithDeadline(title, description string, deadline time.Time) (int, error) {
	query := `
	INSERT INTO tasks(title,description,status,deadline)
	VALUES(?,?,'pending',?)
	`
	result, err := d.conn.Exec(query, title, description, deadline)
	if err != nil {
		return 0, fmt.Errorf("cant create a new task : %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("cant get a id: %w", err)
	}
	return int(id), nil
}

func (d *Database) GetAllTasks() ([]domain.Task, error) {
	query := `
	SELECT *
	FROM tasks
	ORDER BY created_at
	`
	rows, err := d.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("cant complete the query: %w", err)
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task
		var description sql.NullString
		var completedAt, deadline sql.NullTime
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&description,
			&task.Status,
			&task.CreatedAt,
			&completedAt,
			&deadline,
		)
		if err != nil {
			return nil, fmt.Errorf("cant scan rows:%w", err)
		}

		if description.Valid {
			task.Description = description.String
		}
		if completedAt.Valid {
			task.CompletedAt = &completedAt.Time
		}
		if deadline.Valid {
			task.DeadLine = &deadline.Time
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cant handle the rows: %w", err)
	}
	return tasks, nil
}
