package main

import (
	"log"
	"net/http"

	"github.com/Mujhtech/gh-alt/backend/internal/api"
	"github.com/Mujhtech/gh-alt/backend/internal/database"
	"github.com/Mujhtech/gh-alt/backend/internal/git"
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
	
	// Star/Watch/Fork routes
	apiRouter.HandleFunc("/repositories/{id}/star", middleware.AuthMiddleware(api.StarRepository(db))).Methods("PUT", "OPTIONS")
	apiRouter.HandleFunc("/repositories/{id}/star", middleware.AuthMiddleware(api.UnstarRepository(db))).Methods("DELETE", "OPTIONS")
	apiRouter.HandleFunc("/repositories/{id}/watch", middleware.AuthMiddleware(api.WatchRepository(db))).Methods("PUT", "OPTIONS")
	apiRouter.HandleFunc("/repositories/{id}/watch", middleware.AuthMiddleware(api.UnwatchRepository(db))).Methods("DELETE", "OPTIONS")
	apiRouter.HandleFunc("/repositories/{id}/fork", middleware.AuthMiddleware(api.ForkRepository(db))).Methods("POST", "OPTIONS")
	
	// Issue routes
	apiRouter.HandleFunc("/repositories/{id}/issues", api.GetIssues(db)).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/repositories/{id}/issues", middleware.AuthMiddleware(api.CreateIssue(db))).Methods("POST", "OPTIONS")
	apiRouter.HandleFunc("/repositories/{id}/issues/{number}", middleware.AuthMiddleware(api.UpdateIssue(db))).Methods("PATCH", "OPTIONS")
	apiRouter.HandleFunc("/repositories/{id}/issues/{number}/comments", middleware.AuthMiddleware(api.AddIssueComment(db))).Methods("POST", "OPTIONS")
	
	// Pull Request routes
	apiRouter.HandleFunc("/repositories/{id}/pulls", api.GetPullRequests(db)).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/repositories/{id}/pulls", middleware.AuthMiddleware(api.CreatePullRequest(db))).Methods("POST", "OPTIONS")
	apiRouter.HandleFunc("/repositories/{id}/pulls/{number}/merge", middleware.AuthMiddleware(api.MergePullRequest(db))).Methods("POST", "OPTIONS")
	
	// Search routes
	apiRouter.HandleFunc("/search", api.SearchAll(db)).Methods("GET", "OPTIONS")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Git HTTP smart protocol routes
	gitHandler := git.NewGitHTTPHandler(db)
	router.PathPrefix("/git/").Handler(http.StripPrefix("/git", gitHandler))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
