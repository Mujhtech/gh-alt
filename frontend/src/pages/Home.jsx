import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { repositoryService } from '../services/api';
import { useAuth } from '../context/AuthContext';

const Home = () => {
  const [repositories, setRepositories] = useState([]);
  const [showModal, setShowModal] = useState(false);
  const [newRepo, setNewRepo] = useState({ name: '', description: '' });
  const { isAuthenticated } = useAuth();

  useEffect(() => {
    loadRepositories();
  }, []);

  const loadRepositories = async () => {
    try {
      const data = await repositoryService.getAll();
      setRepositories(data);
    } catch (error) {
      console.error('Failed to load repositories:', error);
    }
  };

  const handleCreateRepo = async (e) => {
    e.preventDefault();
    try {
      await repositoryService.create(newRepo.name, newRepo.description);
      setNewRepo({ name: '', description: '' });
      setShowModal(false);
      loadRepositories();
    } catch (error) {
      console.error('Failed to create repository:', error);
    }
  };

  return (
    <div style={styles.container}>
      <div style={styles.header}>
        <h1 style={styles.title}>Repositories</h1>
        {isAuthenticated && (
          <button onClick={() => setShowModal(true)} style={styles.createButton}>
            New Repository
          </button>
        )}
      </div>

      <div style={styles.repoList}>
        {repositories.length === 0 ? (
          <p style={styles.emptyMessage}>No repositories yet. Create one to get started!</p>
        ) : (
          repositories.map((repo) => (
            <Link to={`/repository/${repo.id}`} key={repo.id} style={styles.repoCard}>
              <h3 style={styles.repoName}>{repo.name}</h3>
              <p style={styles.repoDesc}>{repo.description || 'No description'}</p>
              <div style={styles.repoStats}>
                <span>⭐ {repo.stars_count}</span>
                <span>🍴 {repo.forks_count}</span>
                <span>👀 {repo.watchers_count}</span>
                <span>🐛 {repo.open_issues_count} issues</span>
              </div>
              <div style={styles.repoMeta}>
                <span>by {repo.owner_name}</span>
                <span>{new Date(repo.created_at).toLocaleDateString()}</span>
              </div>
            </Link>
          ))
        )}
      </div>

      {showModal && (
        <div style={styles.modal}>
          <div style={styles.modalContent}>
            <h2 style={styles.modalTitle}>Create New Repository</h2>
            <form onSubmit={handleCreateRepo}>
              <input
                type="text"
                placeholder="Repository name"
                value={newRepo.name}
                onChange={(e) => setNewRepo({ ...newRepo, name: e.target.value })}
                style={styles.input}
                required
              />
              <textarea
                placeholder="Description (optional)"
                value={newRepo.description}
                onChange={(e) => setNewRepo({ ...newRepo, description: e.target.value })}
                style={styles.textarea}
              />
              <div style={styles.modalButtons}>
                <button type="submit" style={styles.submitButton}>Create</button>
                <button type="button" onClick={() => setShowModal(false)} style={styles.cancelButton}>
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

const styles = {
  container: {
    maxWidth: '1200px',
    margin: '0 auto',
    padding: '2rem 1rem',
  },
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '2rem',
  },
  title: {
    fontSize: '2rem',
    color: '#fff',
  },
  createButton: {
    backgroundColor: '#238636',
    color: '#fff',
    border: 'none',
    padding: '0.75rem 1.5rem',
    borderRadius: '6px',
    cursor: 'pointer',
    fontSize: '1rem',
  },
  repoList: {
    display: 'grid',
    gap: '1rem',
  },
  emptyMessage: {
    textAlign: 'center',
    color: '#8b949e',
    fontSize: '1.1rem',
    padding: '3rem',
  },
  repoCard: {
    backgroundColor: '#161b22',
    border: '1px solid #30363d',
    padding: '1.5rem',
    borderRadius: '6px',
    textDecoration: 'none',
    color: '#fff',
    transition: 'border-color 0.2s',
  },
  repoName: {
    color: '#58a6ff',
    fontSize: '1.25rem',
    marginBottom: '0.5rem',
  },
  repoDesc: {
    color: '#8b949e',
    marginBottom: '0.5rem',
  },
  repoStats: {
    display: 'flex',
    gap: '1.5rem',
    color: '#8b949e',
    fontSize: '0.875rem',
    marginBottom: '1rem',
  },
  repoMeta: {
    display: 'flex',
    justifyContent: 'space-between',
    color: '#8b949e',
    fontSize: '0.875rem',
  },
  modal: {
    position: 'fixed',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    backgroundColor: 'rgba(0, 0, 0, 0.7)',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
  },
  modalContent: {
    backgroundColor: '#161b22',
    border: '1px solid #30363d',
    borderRadius: '6px',
    padding: '2rem',
    width: '90%',
    maxWidth: '500px',
  },
  modalTitle: {
    color: '#fff',
    marginBottom: '1.5rem',
  },
  input: {
    width: '100%',
    padding: '0.75rem',
    marginBottom: '1rem',
    backgroundColor: '#0d1117',
    border: '1px solid #30363d',
    borderRadius: '6px',
    color: '#fff',
    fontSize: '1rem',
  },
  textarea: {
    width: '100%',
    padding: '0.75rem',
    marginBottom: '1rem',
    backgroundColor: '#0d1117',
    border: '1px solid #30363d',
    borderRadius: '6px',
    color: '#fff',
    fontSize: '1rem',
    minHeight: '100px',
    resize: 'vertical',
  },
  modalButtons: {
    display: 'flex',
    gap: '1rem',
  },
  submitButton: {
    flex: 1,
    backgroundColor: '#238636',
    color: '#fff',
    border: 'none',
    padding: '0.75rem',
    borderRadius: '6px',
    cursor: 'pointer',
    fontSize: '1rem',
  },
  cancelButton: {
    flex: 1,
    backgroundColor: '#21262d',
    color: '#fff',
    border: '1px solid #30363d',
    padding: '0.75rem',
    borderRadius: '6px',
    cursor: 'pointer',
    fontSize: '1rem',
  },
};

export default Home;
