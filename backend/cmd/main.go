package main

import (
	"log"
	"net/http"

	"github.com/Mujhtech/gh-alt/backend/internal/api"
	"github.com/Mujhtech/gh-alt/backend/internal/database"
	"github.com/Mujhtech/gh-alt/backend/internal/middleware"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize database
	db, err := database.Initialize()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Create router
	router := mux.NewRouter()

	// Apply CORS middleware
	router.Use(middleware.CORS)

	// API routes
	apiRouter := router.PathPrefix("/api").Subrouter()
	
	// User routes
	apiRouter.HandleFunc("/auth/register", api.Register(db)).Methods("POST", "OPTIONS")
	apiRouter.HandleFunc("/auth/login", api.Login(db)).Methods("POST", "OPTIONS")
	
	// Repository routes
	apiRouter.HandleFunc("/repositories", api.GetRepositories(db)).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/repositories", middleware.AuthMiddleware(api.CreateRepository(db))).Methods("POST", "OPTIONS")
	apiRouter.HandleFunc("/repositories/{id}", api.GetRepository(db)).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/repositories/{id}/files", api.GetRepositoryFiles(db)).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/repositories/{id}/commits", api.GetRepositoryCommits(db)).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/repositories/{id}/commit", middleware.AuthMiddleware(api.CreateCommit(db))).Methods("POST", "OPTIONS")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
