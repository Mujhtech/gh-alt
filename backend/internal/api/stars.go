package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// StarRepository stars a repository
func StarRepository(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		repoID := vars["id"]

		userID, ok := r.Context().Value("user_id").(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		result, err := db.Exec("INSERT OR IGNORE INTO stars (user_id, repository_id) VALUES (?, ?)",
			userID, repoID)
		if err != nil {
			http.Error(w, "Failed to star repository", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			db.Exec("UPDATE repositories SET stars_count = stars_count + 1 WHERE id = ?", repoID)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "starred"})
	}
}

// UnstarRepository unstars a repository
func UnstarRepository(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		repoID := vars["id"]

		userID, ok := r.Context().Value("user_id").(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		result, err := db.Exec("DELETE FROM stars WHERE user_id = ? AND repository_id = ?",
			userID, repoID)
		if err != nil {
			http.Error(w, "Failed to unstar repository", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			db.Exec("UPDATE repositories SET stars_count = stars_count - 1 WHERE id = ?", repoID)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "unstarred"})
	}
}

// WatchRepository watches a repository
func WatchRepository(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		repoID := vars["id"]

		userID, ok := r.Context().Value("user_id").(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		result, err := db.Exec("INSERT OR IGNORE INTO watchers (user_id, repository_id) VALUES (?, ?)",
			userID, repoID)
		if err != nil {
			http.Error(w, "Failed to watch repository", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			db.Exec("UPDATE repositories SET watchers_count = watchers_count + 1 WHERE id = ?", repoID)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "watching"})
	}
}

// UnwatchRepository unwatches a repository
func UnwatchRepository(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		repoID := vars["id"]

		userID, ok := r.Context().Value("user_id").(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		result, err := db.Exec("DELETE FROM watchers WHERE user_id = ? AND repository_id = ?",
			userID, repoID)
		if err != nil {
			http.Error(w, "Failed to unwatch repository", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			db.Exec("UPDATE repositories SET watchers_count = watchers_count - 1 WHERE id = ?", repoID)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "unwatched"})
	}
}

// ForkRepository forks a repository
func ForkRepository(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		repoID := vars["id"]

		userID, ok := r.Context().Value("user_id").(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get original repository details
		var name, description string
		var parentID int
		err := db.QueryRow("SELECT id, name, description FROM repositories WHERE id = ?", repoID).
			Scan(&parentID, &name, &description)
		if err != nil {
			http.Error(w, "Repository not found", http.StatusNotFound)
			return
		}

		// Create fork
		result, err := db.Exec(`
			INSERT INTO repositories (name, description, owner_id, is_fork, parent_id, is_private)
			VALUES (?, ?, ?, 1, ?, 0)`,
			name, description, userID, parentID)
		if err != nil {
			http.Error(w, "Failed to fork repository", http.StatusInternalServerError)
			return
		}

		forkID, _ := result.LastInsertId()
		db.Exec("UPDATE repositories SET forks_count = forks_count + 1 WHERE id = ?", repoID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":     forkID,
			"status": "forked",
		})
	}
}
