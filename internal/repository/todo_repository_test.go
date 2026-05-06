package repository

import (
	"testing"
	"time"

	"todoapp/internal/domain"
)

func TestListSortedAndFiltered(t *testing.T) {
	repo := NewInMemoryTodoRepository()

	repo.Create(domain.Todo{Text: "b", DueDate: time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)})
	repo.Create(domain.Todo{Text: "a", DueDate: time.Date(2026, 5, 8, 0, 0, 0, 0, time.UTC), Completed: true})
	repo.Create(domain.Todo{Text: "c", DueDate: time.Date(2026, 5, 9, 0, 0, 0, 0, time.UTC)})

	items := repo.List(false)
	if len(items) != 2 {
		t.Fatalf("expected 2 items when excluding completed, got %d", len(items))
	}
	if items[0].Text != "c" || items[1].Text != "b" {
		t.Fatalf("unexpected order: %+v", items)
	}

	all := repo.List(true)
	if len(all) != 3 {
		t.Fatalf("expected 3 items with completed included, got %d", len(all))
	}
	if all[0].Text != "a" {
		t.Fatalf("expected earliest due item first, got %+v", all[0])
	}
}
