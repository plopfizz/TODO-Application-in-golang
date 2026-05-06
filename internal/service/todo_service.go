package service

import (
	"todoapp/internal/domain"
	"todoapp/internal/repository"
	"todoapp/internal/util"
)

// TodoService contains business logic for to-do operations.
type TodoService struct {
	repo repository.TodoRepository
}

// NewTodoService creates a service with repository dependency injected.
func NewTodoService(repo repository.TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

// Create validates request and creates a new to-do item.
func (s *TodoService) Create(req domain.CreateTodoRequest) (domain.Todo, error) {
	text, err := util.ValidateText(req.Text)
	if err != nil {
		return domain.Todo{}, err
	}
	dueDate, err := util.ParseDueDate(req.DueDate)
	if err != nil {
		return domain.Todo{}, err
	}

	return s.repo.Create(domain.Todo{Text: text, DueDate: dueDate, Completed: false}), nil
}

// Get returns a to-do by ID.
func (s *TodoService) Get(id int64) (domain.Todo, error) { return s.repo.GetByID(id) }

// List returns to-dos filtered by completion preference.
func (s *TodoService) List(includeCompleted bool) []domain.Todo { return s.repo.List(includeCompleted) }

// Update modifies allowed fields while enforcing validation.
func (s *TodoService) Update(id int64, req domain.UpdateTodoRequest) (domain.Todo, error) {
	return s.repo.Update(id, func(todo *domain.Todo) error {
		if req.Text != nil {
			text, err := util.ValidateText(*req.Text)
			if err != nil {
				return err
			}
			todo.Text = text
		}
		if req.DueDate != nil {
			dueDate, err := util.ParseDueDate(*req.DueDate)
			if err != nil {
				return err
			}
			todo.DueDate = dueDate
		}
		if req.Completed != nil {
			todo.Completed = *req.Completed
		}
		return nil
	})
}

// Delete removes a to-do by ID.
func (s *TodoService) Delete(id int64) error { return s.repo.Delete(id) }
