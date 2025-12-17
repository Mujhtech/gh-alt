import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { repositoryService } from '../services/api';
import { useAuth } from '../context/AuthContext';

const Issues = () => {
  const { id } = useParams();
  const [issues, setIssues] = useState([]);
  const [repository, setRepository] = useState(null);
  const [showModal, setShowModal] = useState(false);
  const [newIssue, setNewIssue] = useState({ title: '', body: '' });
  const { isAuthenticated } = useAuth();
  const [state, setState] = useState('open');

  useEffect(() => {
    loadRepository();
    loadIssues();
  }, [id, state]);

  const loadRepository = async () => {
    try {
      const data = await repositoryService.getById(id);
      setRepository(data);
    } catch (error) {
      console.error('Failed to load repository:', error);
    }
  };

  const loadIssues = async () => {
    try {
      const data = await repositoryService.getIssues(id, state);
      setIssues(data);
    } catch (error) {
      console.error('Failed to load issues:', error);
    }
  };

  const handleCreateIssue = async (e) => {
    e.preventDefault();
    try {
      await repositoryService.createIssue(id, newIssue.title, newIssue.body);
      setNewIssue({ title: '', body: '' });
      setShowModal(false);
      loadIssues();
    } catch (error) {
      console.error('Failed to create issue:', error);
    }
  };

  if (!repository) return <div style={styles.container}>Loading...</div>;

  return (
    <div style={styles.container}>
      <div style={styles.header}>
        <h1 style={styles.repoName}>
          <Link to={`/repository/${id}`} style={styles.link}>{repository.name}</Link> / Issues
        </h1>
      </div>

      <div style={styles.tabs}>
        <button
          onClick={() => setState('open')}
          style={state === 'open' ? styles.activeTab : styles.tab}
        >
          Open ({issues.filter(i => i.state === 'open').length})
        </button>
        <button
          onClick={() => setState('closed')}
          style={state === 'closed' ? styles.activeTab : styles.tab}
        >
          Closed ({issues.filter(i => i.state === 'closed').length})
        </button>
        {isAuthenticated && (
          <button onClick={() => setShowModal(true)} style={styles.newButton}>
            New Issue
          </button>
        )}
      </div>

      <div style={styles.issueList}>
        {issues.length === 0 ? (
          <p style={styles.emptyMessage}>No {state} issues</p>
        ) : (
          issues.map((issue) => (
            <div key={issue.id} style={styles.issueCard}>
              <div style={styles.issueHeader}>
                <h3 style={styles.issueTitle}>{issue.title}</h3>
                <span style={issue.state === 'open' ? styles.badgeOpen : styles.badgeClosed}>
                  {issue.state}
                </span>
              </div>
              <div style={styles.issueMeta}>
                #{issue.number} opened by {issue.author_name} • {issue.comments_count} comments
              </div>
            </div>
          ))
        )}
      </div>

      {showModal && (
        <div style={styles.modal}>
          <div style={styles.modalContent}>
            <h2 style={styles.modalTitle}>New Issue</h2>
            <form onSubmit={handleCreateIssue}>
              <input
                type="text"
                placeholder="Issue title"
                value={newIssue.title}
                onChange={(e) => setNewIssue({ ...newIssue, title: e.target.value })}
                style={styles.input}
                required
              />
              <textarea
                placeholder="Issue description"
                value={newIssue.body}
                onChange={(e) => setNewIssue({ ...newIssue, body: e.target.value })}
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
    marginBottom: '2rem',
  },
  repoName: {
    fontSize: '2rem',
    color: '#fff',
  },
  link: {
    color: '#58a6ff',
    textDecoration: 'none',
  },
  tabs: {
    display: 'flex',
    gap: '0.5rem',
    marginBottom: '1rem',
    borderBottom: '1px solid #30363d',
    alignItems: 'center',
  },
  tab: {
    backgroundColor: 'transparent',
    color: '#8b949e',
    border: 'none',
    padding: '0.75rem 1rem',
    cursor: 'pointer',
    fontSize: '1rem',
    borderBottom: '2px solid transparent',
  },
  activeTab: {
    backgroundColor: 'transparent',
    color: '#fff',
    border: 'none',
    padding: '0.75rem 1rem',
    cursor: 'pointer',
    fontSize: '1rem',
    borderBottom: '2px solid #f78166',
  },
  newButton: {
    marginLeft: 'auto',
    backgroundColor: '#238636',
    color: '#fff',
    border: 'none',
    padding: '0.5rem 1rem',
    borderRadius: '6px',
    cursor: 'pointer',
  },
  issueList: {
    display: 'grid',
    gap: '0.5rem',
  },
  emptyMessage: {
    textAlign: 'center',
    color: '#8b949e',
    padding: '3rem',
  },
  issueCard: {
    backgroundColor: '#161b22',
    border: '1px solid #30363d',
    padding: '1rem',
    borderRadius: '6px',
    cursor: 'pointer',
  },
  issueHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '0.5rem',
  },
  issueTitle: {
    color: '#fff',
    fontSize: '1rem',
    margin: 0,
  },
  badgeOpen: {
    backgroundColor: '#238636',
    color: '#fff',
    padding: '0.25rem 0.5rem',
    borderRadius: '12px',
    fontSize: '0.75rem',
    fontWeight: 'bold',
  },
  badgeClosed: {
    backgroundColor: '#8250df',
    color: '#fff',
    padding: '0.25rem 0.5rem',
    borderRadius: '12px',
    fontSize: '0.75rem',
    fontWeight: 'bold',
  },
  issueMeta: {
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
    maxWidth: '600px',
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
    minHeight: '150px',
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

export default Issues;
