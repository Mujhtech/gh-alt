package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type Repository struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     int       `json:"owner_id"`
	OwnerName   string    `json:"owner_name"`
	IsPrivate   bool      `json:"is_private"`
	CreatedAt   time.Time `json:"created_at"`
}

type Commit struct {
	ID           int       `json:"id"`
	RepositoryID int       `json:"repository_id"`
	Hash         string    `json:"hash"`
	Message      string    `json:"message"`
	Author       string    `json:"author"`
	CreatedAt    time.Time `json:"created_at"`
}

type File struct {
	ID           int       `json:"id"`
	RepositoryID int       `json:"repository_id"`
	CommitID     int       `json:"commit_id"`
	Path         string    `json:"path"`
	Content      string    `json:"content"`
	Size         int       `json:"size"`
	CreatedAt    time.Time `json:"created_at"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
