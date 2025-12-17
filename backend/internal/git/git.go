package git

import (
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

	// Initialize bare repository (for Git protocol)
	_, err := git.PlainInit(repoPath, true)
	return err
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
