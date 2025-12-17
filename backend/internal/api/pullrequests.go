package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Mujhtech/gh-alt/backend/internal/models"
	"github.com/gorilla/mux"
)

// GetPullRequests returns all PRs for a repository
func GetPullRequests(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		repoID := vars["id"]
		state := r.URL.Query().Get("state")
		if state == "" {
			state = "open"
		}

		query := `
			SELECT pr.id, pr.repository_id, pr.number, pr.title, pr.body, pr.state,
			       pr.author_id, u.username, pr.head_branch, pr.base_branch,
			       pr.merged, pr.merged_at, pr.comments_count, pr.created_at, pr.updated_at, pr.closed_at
			FROM pull_requests pr
			JOIN users u ON pr.author_id = u.id
			WHERE pr.repository_id = ? AND pr.state = ?
			ORDER BY pr.created_at DESC`

		rows, err := db.Query(query, repoID, state)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var prs []models.PullRequest
		for rows.Next() {
			var pr models.PullRequest
			var mergedAt, closedAt sql.NullTime
			err := rows.Scan(&pr.ID, &pr.RepositoryID, &pr.Number, &pr.Title, &pr.Body,
				&pr.State, &pr.AuthorID, &pr.AuthorName, &pr.HeadBranch, &pr.BaseBranch,
				&pr.Merged, &mergedAt, &pr.CommentsCount, &pr.CreatedAt, &pr.UpdatedAt, &closedAt)
			if err != nil {
				continue
			}
			if mergedAt.Valid {
				pr.MergedAt = &mergedAt.Time
			}
			if closedAt.Valid {
				pr.ClosedAt = &closedAt.Time
			}
			prs = append(prs, pr)
		}

		if prs == nil {
			prs = []models.PullRequest{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(prs)
	}
}

// CreatePullRequest creates a new PR
func CreatePullRequest(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		repoID := vars["id"]

		userID, ok := r.Context().Value("user_id").(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req models.CreatePullRequestRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get next PR number
		var maxNumber int
		db.QueryRow("SELECT COALESCE(MAX(number), 0) FROM pull_requests WHERE repository_id = ?", repoID).Scan(&maxNumber)
		number := maxNumber + 1

		result, err := db.Exec(`
			INSERT INTO pull_requests (repository_id, number, title, body, author_id, head_branch, base_branch)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			repoID, number, req.Title, req.Body, userID, req.HeadBranch, req.BaseBranch)
		if err != nil {
			http.Error(w, "Failed to create pull request", http.StatusInternalServerError)
			return
		}

		prID, _ := result.LastInsertId()

		pr := models.PullRequest{
			ID:           int(prID),
			RepositoryID: parseIntOrZero(repoID),
			Number:       number,
			Title:        req.Title,
			Body:         req.Body,
			State:        "open",
			AuthorID:     userID,
			HeadBranch:   req.HeadBranch,
			BaseBranch:   req.BaseBranch,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(pr)
	}
}

// MergePullRequest merges a PR
func MergePullRequest(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		repoID := vars["id"]
		prNum := vars["number"]

		_, err := db.Exec(`
			UPDATE pull_requests 
			SET merged = 1, state = 'closed', merged_at = CURRENT_TIMESTAMP, closed_at = CURRENT_TIMESTAMP
			WHERE repository_id = ? AND number = ?`,
			repoID, prNum)
		if err != nil {
			http.Error(w, "Failed to merge pull request", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "merged"})
	}
}
