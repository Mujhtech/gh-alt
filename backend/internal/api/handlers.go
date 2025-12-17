package api

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Mujhtech/gh-alt/backend/internal/models"
	"github.com/Mujhtech/gh-alt/backend/pkg/auth"
	"github.com/gorilla/mux"
)

func Register(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hashedPassword, err := auth.HashPassword(req.Password)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		result, err := db.Exec(
			"INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
			req.Username, req.Email, hashedPassword,
		)
		if err != nil {
			http.Error(w, "Username or email already exists", http.StatusConflict)
			return
		}

		userID, _ := result.LastInsertId()
		token, _ := auth.GenerateToken(int(userID), req.Username)

		user := models.User{
			ID:       int(userID),
			Username: req.Username,
			Email:    req.Email,
		}

		response := models.AuthResponse{
			Token: token,
			User:  user,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var user models.User
		var hashedPassword string
		err := db.QueryRow(
			"SELECT id, username, email, password FROM users WHERE username = ?",
			req.Username,
		).Scan(&user.ID, &user.Username, &user.Email, &hashedPassword)

		if err != nil || !auth.CheckPasswordHash(req.Password, hashedPassword) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token, _ := auth.GenerateToken(user.ID, user.Username)

		response := models.AuthResponse{
			Token: token,
			User:  user,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func GetRepositories(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`
			SELECT r.id, r.name, r.description, r.owner_id, u.username, r.is_private, r.created_at
			FROM repositories r
			JOIN users u ON r.owner_id = u.id
			WHERE r.is_private = 0
			ORDER BY r.created_at DESC
		`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var repositories []models.Repository
		for rows.Next() {
			var repo models.Repository
			rows.Scan(&repo.ID, &repo.Name, &repo.Description, &repo.OwnerID, &repo.OwnerName, &repo.IsPrivate, &repo.CreatedAt)
			repositories = append(repositories, repo)
		}

		if repositories == nil {
			repositories = []models.Repository{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(repositories)
	}
}

func CreateRepository(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var repo models.Repository
		if err := json.NewDecoder(r.Body).Decode(&repo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// For simplicity, use a default owner_id of 1
		// In production, extract from JWT token
		repo.OwnerID = 1

		result, err := db.Exec(
			"INSERT INTO repositories (name, description, owner_id, is_private) VALUES (?, ?, ?, ?)",
			repo.Name, repo.Description, repo.OwnerID, repo.IsPrivate,
		)
		if err != nil {
			http.Error(w, "Failed to create repository", http.StatusInternalServerError)
			return
		}

		repoID, _ := result.LastInsertId()
		repo.ID = int(repoID)
		repo.CreatedAt = time.Now()

		// Create initial commit
		hash := generateHash(fmt.Sprintf("%s-%d", repo.Name, repoID))
		db.Exec(
			"INSERT INTO commits (repository_id, hash, message, author) VALUES (?, ?, ?, ?)",
			repoID, hash, "Initial commit", "System",
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(repo)
	}
}

func GetRepository(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var repo models.Repository
		err := db.QueryRow(`
			SELECT r.id, r.name, r.description, r.owner_id, u.username, r.is_private, r.created_at
			FROM repositories r
			JOIN users u ON r.owner_id = u.id
			WHERE r.id = ?
		`, id).Scan(&repo.ID, &repo.Name, &repo.Description, &repo.OwnerID, &repo.OwnerName, &repo.IsPrivate, &repo.CreatedAt)

		if err != nil {
			http.Error(w, "Repository not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(repo)
	}
}

func GetRepositoryFiles(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		rows, err := db.Query(`
			SELECT id, repository_id, commit_id, path, content, size, created_at
			FROM files
			WHERE repository_id = ?
			ORDER BY path
		`, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var files []models.File
		for rows.Next() {
			var file models.File
			rows.Scan(&file.ID, &file.RepositoryID, &file.CommitID, &file.Path, &file.Content, &file.Size, &file.CreatedAt)
			files = append(files, file)
		}

		if files == nil {
			files = []models.File{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(files)
	}
}

func GetRepositoryCommits(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		rows, err := db.Query(`
			SELECT id, repository_id, hash, message, author, created_at
			FROM commits
			WHERE repository_id = ?
			ORDER BY created_at DESC
		`, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var commits []models.Commit
		for rows.Next() {
			var commit models.Commit
			rows.Scan(&commit.ID, &commit.RepositoryID, &commit.Hash, &commit.Message, &commit.Author, &commit.CreatedAt)
			commits = append(commits, commit)
		}

		if commits == nil {
			commits = []models.Commit{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(commits)
	}
}

func CreateCommit(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var req struct {
			Message string        `json:"message"`
			Author  string        `json:"author"`
			Files   []models.File `json:"files"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		repoID, _ := strconv.Atoi(id)
		hash := generateHash(fmt.Sprintf("%s-%s-%d", req.Message, req.Author, time.Now().Unix()))

		result, err := db.Exec(
			"INSERT INTO commits (repository_id, hash, message, author) VALUES (?, ?, ?, ?)",
			repoID, hash, req.Message, req.Author,
		)
		if err != nil {
			http.Error(w, "Failed to create commit", http.StatusInternalServerError)
			return
		}

		commitID, _ := result.LastInsertId()

		// Save files
		for _, file := range req.Files {
			db.Exec(
				"INSERT INTO files (repository_id, commit_id, path, content, size) VALUES (?, ?, ?, ?, ?)",
				repoID, commitID, file.Path, file.Content, len(file.Content),
			)
		}

		commit := models.Commit{
			ID:           int(commitID),
			RepositoryID: repoID,
			Hash:         hash,
			Message:      req.Message,
			Author:       req.Author,
			CreatedAt:    time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(commit)
	}
}

func generateHash(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])[:40]
}
