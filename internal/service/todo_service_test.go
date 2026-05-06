package service

import (
	"testing"

	"todoapp/internal/domain"
	"todoapp/internal/repository"
)

func TestCreateValidation(t *testing.T) {
	svc := NewTodoService(repository.NewInMemoryTodoRepository())

	if _, err := svc.Create(domain.CreateTodoRequest{Text: "", DueDate: "2026-05-10"}); err == nil {
		t.Fatal("expected text validation error")
	}
	if _, err := svc.Create(domain.CreateTodoRequest{Text: "Task", DueDate: "2026/05/10"}); err == nil {
		t.Fatal("expected date validation error")
	}
}

func TestUpdateFields(t *testing.T) {
	svc := NewTodoService(repository.NewInMemoryTodoRepository())
	created, err := svc.Create(domain.CreateTodoRequest{Text: "Task", DueDate: "2026-05-10"})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	newText := "Updated Task"
	completed := true
	updated, err := svc.Update(created.ID, domain.UpdateTodoRequest{Text: &newText, Completed: &completed})
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Text != newText || !updated.Completed {
		t.Fatalf("unexpected updated todo: %+v", updated)
	}
}
