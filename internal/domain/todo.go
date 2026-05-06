package domain

import "time"

// Todo represents a single to-do item.
type Todo struct {
	ID        int64     `json:"id"`
	Text      string    `json:"text"`
	DueDate   time.Time `json:"due_date"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateTodoRequest defines payload for creating a to-do.
type CreateTodoRequest struct {
	Text    string `json:"text"`
	DueDate string `json:"due_date"`
}

// UpdateTodoRequest defines payload for updating a to-do.
type UpdateTodoRequest struct {
	Text      *string `json:"text,omitempty"`
	DueDate   *string `json:"due_date,omitempty"`
	Completed *bool   `json:"completed,omitempty"`
}
