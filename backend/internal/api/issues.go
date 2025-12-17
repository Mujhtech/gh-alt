package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Mujhtech/gh-alt/backend/internal/models"
	"github.com/gorilla/mux"
)

// GetIssues returns all issues for a repository
func GetIssues(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		repoID := vars["id"]
		state := r.URL.Query().Get("state")
		if state == "" {
			state = "open"
		}

		query := `
			SELECT i.id, i.repository_id, i.number, i.title, i.body, i.state, 
			       i.author_id, u.username, i.assignee_id, i.comments_count,
			       i.created_at, i.updated_at, i.closed_at
			FROM issues i
			JOIN users u ON i.author_id = u.id
			WHERE i.repository_id = ?`
		
		if state != "all" {
			query += " AND i.state = ?"
		}
		query += " ORDER BY i.created_at DESC"

		var rows *sql.Rows
		var err error
		if state != "all" {
			rows, err = db.Query(query, repoID, state)
		} else {
			rows, err = db.Query(query, repoID)
		}
		
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var issues []models.Issue
		for rows.Next() {
			var issue models.Issue
			var closedAt sql.NullTime
			var assigneeID sql.NullInt64
			err := rows.Scan(&issue.ID, &issue.RepositoryID, &issue.Number, &issue.Title,
				&issue.Body, &issue.State, &issue.AuthorID, &issue.AuthorName,
				&assigneeID, &issue.CommentsCount, &issue.CreatedAt, &issue.UpdatedAt, &closedAt)
			if err != nil {
				continue
			}
			if closedAt.Valid {
				issue.ClosedAt = &closedAt.Time
			}
			if assigneeID.Valid {
				id := int(assigneeID.Int64)
				issue.AssigneeID = &id
			}
			issues = append(issues, issue)
		}

		if issues == nil {
			issues = []models.Issue{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(issues)
	}
}

// CreateIssue creates a new issue
func CreateIssue(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		repoID := vars["id"]

		userID, ok := r.Context().Value("user_id").(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req models.CreateIssueRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get next issue number
		var maxNumber int
		db.QueryRow("SELECT COALESCE(MAX(number), 0) FROM issues WHERE repository_id = ?", repoID).Scan(&maxNumber)
		number := maxNumber + 1

		result, err := db.Exec(`
			INSERT INTO issues (repository_id, number, title, body, author_id, assignee_id)
			VALUES (?, ?, ?, ?, ?, ?)`,
			repoID, number, req.Title, req.Body, userID, req.AssigneeID)
		if err != nil {
			http.Error(w, "Failed to create issue", http.StatusInternalServerError)
			return
		}

		issueID, _ := result.LastInsertId()

		// Update repository open issues count
		db.Exec("UPDATE repositories SET open_issues_count = open_issues_count + 1 WHERE id = ?", repoID)

		issue := models.Issue{
			ID:           int(issueID),
			RepositoryID: parseIntOrZero(repoID),
			Number:       number,
			Title:        req.Title,
			Body:         req.Body,
			State:        "open",
			AuthorID:     userID,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(issue)
	}
}

// UpdateIssue updates an issue (close/reopen)
func UpdateIssue(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		repoID := vars["id"]
		issueNum := vars["number"]

		var req struct {
			State string `json:"state"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var currentState string
		err := db.QueryRow("SELECT state FROM issues WHERE repository_id = ? AND number = ?",
			repoID, issueNum).Scan(&currentState)
		if err != nil {
			http.Error(w, "Issue not found", http.StatusNotFound)
			return
		}

		if req.State == "closed" && currentState == "open" {
			db.Exec("UPDATE issues SET state = 'closed', closed_at = CURRENT_TIMESTAMP WHERE repository_id = ? AND number = ?",
				repoID, issueNum)
			db.Exec("UPDATE repositories SET open_issues_count = open_issues_count - 1 WHERE id = ?", repoID)
		} else if req.State == "open" && currentState == "closed" {
			db.Exec("UPDATE issues SET state = 'open', closed_at = NULL WHERE repository_id = ? AND number = ?",
				repoID, issueNum)
			db.Exec("UPDATE repositories SET open_issues_count = open_issues_count + 1 WHERE id = ?", repoID)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
	}
}

// AddIssueComment adds a comment to an issue
func AddIssueComment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		repoID := vars["id"]
		issueNum := vars["number"]

		userID, ok := r.Context().Value("user_id").(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req struct {
			Body string `json:"body"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var issueID int
		err := db.QueryRow("SELECT id FROM issues WHERE repository_id = ? AND number = ?",
			repoID, issueNum).Scan(&issueID)
		if err != nil {
			http.Error(w, "Issue not found", http.StatusNotFound)
			return
		}

		result, err := db.Exec(`INSERT INTO issue_comments (issue_id, author_id, body) VALUES (?, ?, ?)`,
			issueID, userID, req.Body)
		if err != nil {
			http.Error(w, "Failed to add comment", http.StatusInternalServerError)
			return
		}

		commentID, _ := result.LastInsertId()
		db.Exec("UPDATE issues SET comments_count = comments_count + 1 WHERE id = ?", issueID)

		comment := models.IssueComment{
			ID:        int(commentID),
			IssueID:   issueID,
			AuthorID:  userID,
			Body:      req.Body,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(comment)
	}
}

func parseIntOrZero(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}
