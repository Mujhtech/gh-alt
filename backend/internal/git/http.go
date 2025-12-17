package git

import (
	"compress/gzip"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	gitUploadPackRe  = regexp.MustCompile(`^/(.+)/git-upload-pack$`)
	gitReceivePackRe = regexp.MustCompile(`^/(.+)/git-receive-pack$`)
	infoRefsRe       = regexp.MustCompile(`^/(.+)/info/refs$`)
)

// GitHTTPHandler handles Git HTTP smart protocol requests
type GitHTTPHandler struct {
	DB *sql.DB
}

// ServeHTTP implements the http.Handler interface
func (h *GitHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract repository name from path
	if matches := infoRefsRe.FindStringSubmatch(r.URL.Path); matches != nil {
		h.handleInfoRefs(w, r, matches[1])
		return
	}

	if matches := gitUploadPackRe.FindStringSubmatch(r.URL.Path); matches != nil {
		h.handleGitUploadPack(w, r, matches[1])
		return
	}

	if matches := gitReceivePackRe.FindStringSubmatch(r.URL.Path); matches != nil {
		h.handleGitReceivePack(w, r, matches[1])
		return
	}

	http.Error(w, "Not found", http.StatusNotFound)
}

// handleInfoRefs handles the info/refs endpoint
func (h *GitHTTPHandler) handleInfoRefs(w http.ResponseWriter, r *http.Request, repoName string) {
	service := r.URL.Query().Get("service")
	if service != "git-upload-pack" && service != "git-receive-pack" {
		http.Error(w, "Invalid service", http.StatusBadRequest)
		return
	}

	// Get repository from database
	var repoID int
	var name string
	err := h.DB.QueryRow("SELECT id, name FROM repositories WHERE name = ?", repoName).Scan(&repoID, &name)
	if err != nil {
		http.Error(w, "Repository not found", http.StatusNotFound)
		return
	}

	repoPath := GetRepoPath(repoID, name)
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		http.Error(w, "Repository not initialized", http.StatusNotFound)
		return
	}

	// Execute git service
	cmd := exec.Command("git", strings.TrimPrefix(service, "git-"), "--stateless-rpc", "--advertise-refs", repoPath)
	output, err := cmd.Output()
	if err != nil {
		http.Error(w, "Git command failed", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", fmt.Sprintf("application/x-%s-advertisement", service))
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)

	// Write packet line with proper length
	serviceLine := fmt.Sprintf("# service=%s\n", service)
	fmt.Fprintf(w, "%04x%s", len(serviceLine)+4, serviceLine)
	fmt.Fprint(w, "0000")
	w.Write(output)
}

// handleGitUploadPack handles git-upload-pack (fetch/clone)
func (h *GitHTTPHandler) handleGitUploadPack(w http.ResponseWriter, r *http.Request, repoName string) {
	// Get repository from database
	var repoID int
	var name string
	err := h.DB.QueryRow("SELECT id, name FROM repositories WHERE name = ?", repoName).Scan(&repoID, &name)
	if err != nil {
		http.Error(w, "Repository not found", http.StatusNotFound)
		return
	}

	repoPath := GetRepoPath(repoID, name)
	h.handleGitCommand(w, r, repoPath, "upload-pack")
}

// handleGitReceivePack handles git-receive-pack (push)
func (h *GitHTTPHandler) handleGitReceivePack(w http.ResponseWriter, r *http.Request, repoName string) {
	// Get repository from database
	var repoID int
	var name string
	err := h.DB.QueryRow("SELECT id, name FROM repositories WHERE name = ?", repoName).Scan(&repoID, &name)
	if err != nil {
		http.Error(w, "Repository not found", http.StatusNotFound)
		return
	}

	repoPath := GetRepoPath(repoID, name)
	h.handleGitCommand(w, r, repoPath, "receive-pack")
}

// handleGitCommand executes a git command and streams the response
func (h *GitHTTPHandler) handleGitCommand(w http.ResponseWriter, r *http.Request, repoPath, service string) {
	w.Header().Set("Content-Type", fmt.Sprintf("application/x-git-%s-result", service))
	w.Header().Set("Cache-Control", "no-cache")

	// Handle gzip encoding
	var reqBody io.Reader = r.Body
	if r.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, "Failed to decompress", http.StatusBadRequest)
			return
		}
		defer gz.Close()
		reqBody = gz
	}

	// Execute git command
	cmd := exec.Command("git", service, "--stateless-rpc", repoPath)
	cmd.Stdin = reqBody
	cmd.Stdout = w
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// Error already written to response
		return
	}
}

// NewGitHTTPHandler creates a new Git HTTP handler
func NewGitHTTPHandler(db *sql.DB) *GitHTTPHandler {
	return &GitHTTPHandler{DB: db}
}
