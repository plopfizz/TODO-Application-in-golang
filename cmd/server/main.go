package main

import (
	"log"

	"todoapp/internal/app"
)

// main boots the REST API server.
func main() {
	server := app.NewServer(":8080")
	log.Printf("to-do API server started at %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
