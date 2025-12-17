package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func Initialize() (*sql.DB, error) {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./gh-alt.db"
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create tables
	if err := createTables(db); err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return db, nil
}

func createTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			name TEXT,
			bio TEXT,
			avatar_url TEXT,
			location TEXT,
			website TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS repositories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			owner_id INTEGER NOT NULL,
			is_private BOOLEAN DEFAULT 0,
			is_fork BOOLEAN DEFAULT 0,
			parent_id INTEGER,
			default_branch TEXT DEFAULT 'master',
			stars_count INTEGER DEFAULT 0,
			forks_count INTEGER DEFAULT 0,
			watchers_count INTEGER DEFAULT 0,
			open_issues_count INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (owner_id) REFERENCES users(id),
			FOREIGN KEY (parent_id) REFERENCES repositories(id)
		)`,
		`CREATE TABLE IF NOT EXISTS commits (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			repository_id INTEGER NOT NULL,
			hash TEXT UNIQUE NOT NULL,
			message TEXT NOT NULL,
			author TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (repository_id) REFERENCES repositories(id)
		)`,
		`CREATE TABLE IF NOT EXISTS files (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			repository_id INTEGER NOT NULL,
			commit_id INTEGER NOT NULL,
			path TEXT NOT NULL,
			content TEXT NOT NULL,
			size INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (repository_id) REFERENCES repositories(id),
			FOREIGN KEY (commit_id) REFERENCES commits(id)
		)`,
		`CREATE TABLE IF NOT EXISTS issues (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			repository_id INTEGER NOT NULL,
			number INTEGER NOT NULL,
			title TEXT NOT NULL,
			body TEXT,
			state TEXT DEFAULT 'open',
			author_id INTEGER NOT NULL,
			assignee_id INTEGER,
			comments_count INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			closed_at DATETIME,
			FOREIGN KEY (repository_id) REFERENCES repositories(id),
			FOREIGN KEY (author_id) REFERENCES users(id),
			FOREIGN KEY (assignee_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS issue_comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			issue_id INTEGER NOT NULL,
			author_id INTEGER NOT NULL,
			body TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (issue_id) REFERENCES issues(id),
			FOREIGN KEY (author_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS pull_requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			repository_id INTEGER NOT NULL,
			number INTEGER NOT NULL,
			title TEXT NOT NULL,
			body TEXT,
			state TEXT DEFAULT 'open',
			author_id INTEGER NOT NULL,
			head_branch TEXT NOT NULL,
			base_branch TEXT NOT NULL,
			merged BOOLEAN DEFAULT 0,
			merged_at DATETIME,
			comments_count INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			closed_at DATETIME,
			FOREIGN KEY (repository_id) REFERENCES repositories(id),
			FOREIGN KEY (author_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS pr_comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			pull_request_id INTEGER NOT NULL,
			author_id INTEGER NOT NULL,
			body TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (pull_request_id) REFERENCES pull_requests(id),
			FOREIGN KEY (author_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS stars (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			repository_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(user_id, repository_id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (repository_id) REFERENCES repositories(id)
		)`,
		`CREATE TABLE IF NOT EXISTS watchers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			repository_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(user_id, repository_id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (repository_id) REFERENCES repositories(id)
		)`,
		`CREATE TABLE IF NOT EXISTS collaborators (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			repository_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			permission TEXT DEFAULT 'read',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(repository_id, user_id),
			FOREIGN KEY (repository_id) REFERENCES repositories(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS branches (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			repository_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			commit_hash TEXT NOT NULL,
			is_default BOOLEAN DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(repository_id, name),
			FOREIGN KEY (repository_id) REFERENCES repositories(id)
		)`,
		`CREATE TABLE IF NOT EXISTS labels (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			repository_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			color TEXT NOT NULL,
			description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(repository_id, name),
			FOREIGN KEY (repository_id) REFERENCES repositories(id)
		)`,
		`CREATE TABLE IF NOT EXISTS issue_labels (
			issue_id INTEGER NOT NULL,
			label_id INTEGER NOT NULL,
			PRIMARY KEY (issue_id, label_id),
			FOREIGN KEY (issue_id) REFERENCES issues(id),
			FOREIGN KEY (label_id) REFERENCES labels(id)
		)`,
		`CREATE TABLE IF NOT EXISTS notifications (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			type TEXT NOT NULL,
			repository_id INTEGER,
			issue_id INTEGER,
			pull_request_id INTEGER,
			read BOOLEAN DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (repository_id) REFERENCES repositories(id),
			FOREIGN KEY (issue_id) REFERENCES issues(id),
			FOREIGN KEY (pull_request_id) REFERENCES pull_requests(id)
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}
