package repository

import (
	"errors"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"todoapp/internal/domain"
)

var (
	// ErrNotFound indicates that a to-do item does not exist.
	ErrNotFound = errors.New("todo item not found")
)

// TodoRepository abstracts persistence concerns.
type TodoRepository interface {
	Create(todo domain.Todo) domain.Todo
	GetByID(id int64) (domain.Todo, error)
	List(includeCompleted bool) []domain.Todo
	Update(id int64, updateFn func(*domain.Todo) error) (domain.Todo, error)
	Delete(id int64) error
}

// InMemoryTodoRepository stores to-dos with thread-safe access.
type InMemoryTodoRepository struct {
	mu    sync.RWMutex
	seq   int64
	store map[int64]domain.Todo
}

// NewInMemoryTodoRepository constructs a repository instance.
func NewInMemoryTodoRepository() *InMemoryTodoRepository {
	return &InMemoryTodoRepository{store: make(map[int64]domain.Todo)}
}

// Create persists a to-do and assigns unique metadata.
func (r *InMemoryTodoRepository) Create(todo domain.Todo) domain.Todo {
	now := time.Now().UTC()
	todo.ID = atomic.AddInt64(&r.seq, 1)
	todo.CreatedAt = now
	todo.UpdatedAt = now

	r.mu.Lock()
	r.store[todo.ID] = todo
	r.mu.Unlock()

	return todo
}

// GetByID fetches a single to-do by identifier.
func (r *InMemoryTodoRepository) GetByID(id int64) (domain.Todo, error) {
	r.mu.RLock()
	todo, ok := r.store[id]
	r.mu.RUnlock()
	if !ok {
		return domain.Todo{}, ErrNotFound
	}
	return todo, nil
}

// List returns to-dos sorted by earliest due date.
func (r *InMemoryTodoRepository) List(includeCompleted bool) []domain.Todo {
	r.mu.RLock()
	items := make([]domain.Todo, 0, len(r.store))
	for _, item := range r.store {
		if !includeCompleted && item.Completed {
			continue
		}
		items = append(items, item)
	}
	r.mu.RUnlock()

	sort.Slice(items, func(i, j int) bool {
		if items[i].DueDate.Equal(items[j].DueDate) {
			return items[i].ID < items[j].ID
		}
		return items[i].DueDate.Before(items[j].DueDate)
	})
	return items
}

// Update atomically mutates an existing to-do.
func (r *InMemoryTodoRepository) Update(id int64, updateFn func(*domain.Todo) error) (domain.Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	todo, ok := r.store[id]
	if !ok {
		return domain.Todo{}, ErrNotFound
	}

	if err := updateFn(&todo); err != nil {
		return domain.Todo{}, err
	}
	todo.UpdatedAt = time.Now().UTC()
	r.store[id] = todo
	return todo, nil
}

// Delete removes a to-do by identifier.
func (r *InMemoryTodoRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[id]; !ok {
		return ErrNotFound
	}
	delete(r.store, id)
	return nil
}
