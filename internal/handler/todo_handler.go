package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"todoapp/internal/domain"
	"todoapp/internal/repository"
	"todoapp/internal/service"
)

// TodoHandler wires HTTP handlers to the service.
type TodoHandler struct {
	service *service.TodoService
}

// NewTodoHandler constructs a to-do handler.
func NewTodoHandler(service *service.TodoService) *TodoHandler { return &TodoHandler{service: service} }

// RegisterRoutes binds all CRUD endpoints on default mux.
func (h *TodoHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/todos", h.handleTodos)
	mux.HandleFunc("/todos/", h.handleTodoByID)
}

// handleTodos routes collection endpoints for create and list.
func (h *TodoHandler) handleTodos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createTodo(w, r)
	case http.MethodGet:
		h.listTodos(w, r)
	default:
		h.respondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// handleTodoByID routes item endpoints for read, update, and delete.
func (h *TodoHandler) handleTodoByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getTodo(w, id)
	case http.MethodPut:
		h.updateTodo(w, r, id)
	case http.MethodDelete:
		h.deleteTodo(w, id)
	default:
		h.respondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *TodoHandler) createTodo(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	todo, err := h.service.Create(req)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	h.respondJSON(w, http.StatusCreated, todo)
}

func (h *TodoHandler) getTodo(w http.ResponseWriter, id int64) {
	todo, err := h.service.Get(id)
	if err != nil {
		h.handleDomainError(w, err)
		return
	}
	h.respondJSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) listTodos(w http.ResponseWriter, r *http.Request) {
	includeCompleted := strings.EqualFold(r.URL.Query().Get("include_completed"), "true")
	todos := h.service.List(includeCompleted)
	h.respondJSON(w, http.StatusOK, todos)
}

func (h *TodoHandler) updateTodo(w http.ResponseWriter, r *http.Request, id int64) {
	var req domain.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	todo, err := h.service.Update(id, req)
	if err != nil {
		h.handleDomainError(w, err)
		return
	}
	h.respondJSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) deleteTodo(w http.ResponseWriter, id int64) {
	if err := h.service.Delete(id); err != nil {
		h.handleDomainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func parseID(path string) (int64, error) {
	idPart := strings.TrimPrefix(path, "/todos/")
	if idPart == "" || strings.Contains(idPart, "/") {
		return 0, errors.New("invalid id path")
	}
	id, err := strconv.ParseInt(idPart, 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("id must be a positive integer")
	}
	return id, nil
}

func (h *TodoHandler) handleDomainError(w http.ResponseWriter, err error) {
	if errors.Is(err, repository.ErrNotFound) {
		h.respondError(w, http.StatusNotFound, err.Error())
		return
	}
	h.respondError(w, http.StatusBadRequest, err.Error())
}

func (h *TodoHandler) respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *TodoHandler) respondError(w http.ResponseWriter, status int, msg string) {
	h.respondJSON(w, status, map[string]string{"error": msg})
}
