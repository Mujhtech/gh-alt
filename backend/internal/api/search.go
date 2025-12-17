package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Mujhtech/gh-alt/backend/internal/models"
)

// SearchAll searches across repositories, users, and issues
func SearchAll(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		searchType := r.URL.Query().Get("type") // repositories, users, issues
		
		if q == "" {
			http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
			return
		}

		result := models.SearchResult{}
		searchTerm := "%" + q + "%"

		// Search repositories
		if searchType == "" || searchType == "repositories" {
			rows, err := db.Query(`
				SELECT r.id, r.name, r.description, r.owner_id, u.username, r.is_private,
				       r.stars_count, r.forks_count, r.created_at, r.updated_at
				FROM repositories r
				JOIN users u ON r.owner_id = u.id
				WHERE (r.name LIKE ? OR r.description LIKE ?) AND r.is_private = 0
				LIMIT 20`, searchTerm, searchTerm)
			if err == nil {
				defer rows.Close()
				for rows.Next() {
					var repo models.Repository
					rows.Scan(&repo.ID, &repo.Name, &repo.Description, &repo.OwnerID,
						&repo.OwnerName, &repo.IsPrivate, &repo.StarsCount, &repo.ForksCount,
						&repo.CreatedAt, &repo.UpdatedAt)
					result.Repositories = append(result.Repositories, repo)
				}
			}
		}

		// Search users
		if searchType == "" || searchType == "users" {
			rows, err := db.Query(`
				SELECT id, username, email, name, bio, avatar_url, location, created_at
				FROM users
				WHERE username LIKE ? OR name LIKE ? OR email LIKE ?
				LIMIT 20`, searchTerm, searchTerm, searchTerm)
			if err == nil {
				defer rows.Close()
				for rows.Next() {
					var user models.User
					var name, bio, avatarURL, location sql.NullString
					rows.Scan(&user.ID, &user.Username, &user.Email, &name, &bio,
						&avatarURL, &location, &user.CreatedAt)
					if name.Valid {
						user.Name = name.String
					}
					if bio.Valid {
						user.Bio = bio.String
					}
					if avatarURL.Valid {
						user.AvatarURL = avatarURL.String
					}
					if location.Valid {
						user.Location = location.String
					}
					result.Users = append(result.Users, user)
				}
			}
		}

		// Search issues
		if searchType == "" || searchType == "issues" {
			rows, err := db.Query(`
				SELECT i.id, i.repository_id, i.number, i.title, i.body, i.state,
				       i.author_id, u.username, i.comments_count, i.created_at
				FROM issues i
				JOIN users u ON i.author_id = u.id
				JOIN repositories r ON i.repository_id = r.id
				WHERE (i.title LIKE ? OR i.body LIKE ?) AND r.is_private = 0
				LIMIT 20`, searchTerm, searchTerm)
			if err == nil {
				defer rows.Close()
				for rows.Next() {
					var issue models.Issue
					rows.Scan(&issue.ID, &issue.RepositoryID, &issue.Number, &issue.Title,
						&issue.Body, &issue.State, &issue.AuthorID, &issue.AuthorName,
						&issue.CommentsCount, &issue.CreatedAt)
					result.Issues = append(result.Issues, issue)
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}
