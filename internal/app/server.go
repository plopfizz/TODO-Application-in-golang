package app

import (
	"log"
	"net/http"
	"time"

	"todoapp/internal/handler"
	"todoapp/internal/repository"
	"todoapp/internal/service"
)

// NewServer creates an HTTP server with production-safe timeouts.
func NewServer(addr string) *http.Server {
	repo := repository.NewInMemoryTodoRepository()
	todoService := service.NewTodoService(repo)
	todoHandler := handler.NewTodoHandler(todoService)

	mux := http.NewServeMux()
	todoHandler.RegisterRoutes(mux)

	return &http.Server{
		Addr:              addr,
		Handler:           loggingMiddleware(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

// loggingMiddleware emits basic request logs for observability.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.String())
		next.ServeHTTP(w, r)
	})
}
