# To-Do Application in Go

A modular RESTful to-do application built with Go.

## Features
- Create, read, list, update, and delete to-do items.
- Due-date validation with `YYYY-MM-DD` format.
- List endpoint sorts by earliest due date.
- Completed items are excluded by default and can be included via query param.
- Thread-safe in-memory storage for concurrent access.

## Run
```bash
go run ./cmd/server
```

## API Endpoints
- `POST /todos`
- `GET /todos/{id}`
- `GET /todos?include_completed=true`
- `PUT /todos/{id}`
- `DELETE /todos/{id}`

## Validation Rules
- `text` is required and cannot be empty.
- `text` max length is 250 characters.
- `due_date` is required for create and must be `YYYY-MM-DD`.
- `due_date` is optional for update but must be `YYYY-MM-DD` if provided.
- `id` in path must be a positive integer.
