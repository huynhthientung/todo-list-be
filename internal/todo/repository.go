package todo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNotFound = errors.New("todo not found")
	columns     = "id, title, completed, created_at, updated_at"
)

// Repository provides CRUD operations for todo items.
type Repository struct {
	db *sql.DB
}

// NewRepository initializes a Repository backed by the provided database.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// EnsureSchema creates the todos table if it does not already exist.
func (r *Repository) EnsureSchema(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS todos (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			completed BOOLEAN NOT NULL DEFAULT false,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	return err
}

// List returns all todo items.
func (r *Repository) List(ctx context.Context) ([]Todo, error) {
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`SELECT %s FROM todos ORDER BY id`, columns))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}

// Get returns a single todo item by ID.
func (r *Repository) Get(ctx context.Context, id int64) (Todo, error) {
	var t Todo
	err := r.db.QueryRowContext(ctx, fmt.Sprintf(`SELECT %s FROM todos WHERE id = $1`, columns), id).
		Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Todo{}, ErrNotFound
		}
		return Todo{}, err
	}
	return t, nil
}

// Create inserts a new todo item.
func (r *Repository) Create(ctx context.Context, title string, completed bool) (Todo, error) {
	var t Todo
	err := r.db.QueryRowContext(
		ctx,
		fmt.Sprintf(`INSERT INTO todos (title, completed) VALUES ($1, $2) RETURNING %s`, columns),
		title,
		completed,
	).Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return Todo{}, err
	}
	return t, nil
}

// Update modifies an existing todo item. At least one of title or completed must be provided.
func (r *Repository) Update(ctx context.Context, id int64, title *string, completed *bool) (Todo, error) {
	if title == nil && completed == nil {
		return Todo{}, errors.New("no fields to update")
	}

	setParts := make([]string, 0, 2)
	args := make([]any, 0, 3)

	if title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", len(args)+1))
		args = append(args, *title)
	}
	if completed != nil {
		setParts = append(setParts, fmt.Sprintf("completed = $%d", len(args)+1))
		args = append(args, *completed)
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = NOW()"))
	setClause := strings.Join(setParts, ", ")
	args = append(args, id)

	query := fmt.Sprintf(`UPDATE todos SET %s WHERE id = $%d RETURNING %s`, setClause, len(args), columns)

	var t Todo
	err := r.db.QueryRowContext(ctx, query, args...).
		Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Todo{}, ErrNotFound
		}
		return Todo{}, err
	}
	return t, nil
}

// Delete removes a todo item.
func (r *Repository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM todos WHERE id = $1`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}
