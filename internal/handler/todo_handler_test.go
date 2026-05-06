package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"todoapp/internal/repository"
	"todoapp/internal/service"
)

func setupTestServer() *httptest.Server {
	repo := repository.NewInMemoryTodoRepository()
	svc := service.NewTodoService(repo)
	h := NewTodoHandler(svc)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return httptest.NewServer(mux)
}

func TestCRUDFlow(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	createPayload := map[string]string{"text": "Task A", "due_date": "2026-05-10"}
	buf, _ := json.Marshal(createPayload)
	resp, err := http.Post(ts.URL+"/todos", "application/json", bytes.NewReader(buf))
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("create failed err=%v status=%v", err, resp.StatusCode)
	}

	resp, err = http.Get(ts.URL + "/todos")
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("list failed err=%v status=%v", err, resp.StatusCode)
	}

	req, _ := http.NewRequest(http.MethodPut, ts.URL+"/todos/1", bytes.NewBufferString(`{"completed":true}`))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("update failed err=%v status=%v", err, resp.StatusCode)
	}

	resp, err = http.Get(ts.URL + "/todos")
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("list post-update failed err=%v status=%v", err, resp.StatusCode)
	}
	var list []map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&list)
	if len(list) != 0 {
		t.Fatalf("expected default list to exclude completed items, got %d", len(list))
	}

	req, _ = http.NewRequest(http.MethodDelete, ts.URL+"/todos/1", nil)
	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != http.StatusNoContent {
		t.Fatalf("delete failed err=%v status=%v", err, resp.StatusCode)
	}
}
