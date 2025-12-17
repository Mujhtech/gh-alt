package git

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// GetRepoPath returns the filesystem path for a repository
func GetRepoPath(repoID int, repoName string) string {
	reposDir := os.Getenv("REPOS_DIR")
	if reposDir == "" {
		reposDir = "./repositories"
	}
	return filepath.Join(reposDir, fmt.Sprintf("%d-%s", repoID, repoName))
}

// InitRepository initializes a bare Git repository on the filesystem
func InitRepository(repoID int, repoName string) error {
	repoPath := GetRepoPath(repoID, repoName)
	
	// Create directory
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		return err
	}

	// Initialize bare repository
	_, err := git.PlainInit(repoPath, true)
	return err
}

// CommitFiles commits files to the repository
func CommitFiles(db *sql.DB, repoID int, repoName, message, author string, files map[string]string) (string, error) {
	repoPath := GetRepoPath(repoID, repoName)
	
	// Open the repository
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", err
	}

	// Get the worktree
	w, err := repo.Worktree()
	if err != nil {
		return "", err
	}

	// Write files to worktree
	for path, content := range files {
		filePath := filepath.Join(repoPath, path)
		dir := filepath.Dir(filePath)
		
		// Create directory if needed
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", err
		}

		// Write file
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return "", err
		}

		// Add to git
		if _, err := w.Add(path); err != nil {
			return "", err
		}
	}

	// Commit
	commit, err := w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name: author,
		},
	})
	if err != nil {
		return "", err
	}

	return commit.String(), nil
}

// GetCommitHistory retrieves commit history from the Git repository
func GetCommitHistory(repoID int, repoName string) ([]object.Commit, error) {
	repoPath := GetRepoPath(repoID, repoName)
	
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	// Get HEAD reference
	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	// Get commit history
	cIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, err
	}

	var commits []object.Commit
	err = cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, *c)
		return nil
	})

	return commits, err
}
