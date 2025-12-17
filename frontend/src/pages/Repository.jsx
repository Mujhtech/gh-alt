import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { repositoryService } from '../services/api';

const Repository = () => {
  const { id } = useParams();
  const [repository, setRepository] = useState(null);
  const [commits, setCommits] = useState([]);
  const [files, setFiles] = useState([]);
  const [activeTab, setActiveTab] = useState('files');
  const [selectedFile, setSelectedFile] = useState(null);

  useEffect(() => {
    loadRepository();
    loadCommits();
    loadFiles();
  }, [id]);

  const loadRepository = async () => {
    try {
      const data = await repositoryService.getById(id);
      setRepository(data);
    } catch (error) {
      console.error('Failed to load repository:', error);
    }
  };

  const loadCommits = async () => {
    try {
      const data = await repositoryService.getCommits(id);
      setCommits(data);
    } catch (error) {
      console.error('Failed to load commits:', error);
    }
  };

  const loadFiles = async () => {
    try {
      const data = await repositoryService.getFiles(id);
      setFiles(data);
    } catch (error) {
      console.error('Failed to load files:', error);
    }
  };

  if (!repository) {
    return <div style={styles.container}>Loading...</div>;
  }

  return (
    <div style={styles.container}>
      <div style={styles.header}>
        <h1 style={styles.repoName}>{repository.name}</h1>
        <p style={styles.repoDesc}>{repository.description || 'No description'}</p>
        <div style={styles.repoMeta}>
          <span>Owner: {repository.owner_name}</span>
          <span>Created: {new Date(repository.created_at).toLocaleDateString()}</span>
        </div>
      </div>

      <div style={styles.tabs}>
        <button
          onClick={() => setActiveTab('files')}
          style={activeTab === 'files' ? styles.activeTab : styles.tab}
        >
          Files ({files.length})
        </button>
        <button
          onClick={() => setActiveTab('commits')}
          style={activeTab === 'commits' ? styles.activeTab : styles.tab}
        >
          Commits ({commits.length})
        </button>
      </div>

      {activeTab === 'files' && (
        <div style={styles.content}>
          {files.length === 0 ? (
            <p style={styles.emptyMessage}>No files yet</p>
          ) : (
            <div>
              {files.map((file) => (
                <div
                  key={file.id}
                  style={styles.fileItem}
                  onClick={() => setSelectedFile(selectedFile?.id === file.id ? null : file)}
                >
                  <div style={styles.fileName}>📄 {file.path}</div>
                  <div style={styles.fileSize}>{file.size} bytes</div>
                </div>
              ))}
              {selectedFile && (
                <div style={styles.fileContent}>
                  <h3 style={styles.fileContentTitle}>{selectedFile.path}</h3>
                  <pre style={styles.code}>{selectedFile.content}</pre>
                </div>
              )}
            </div>
          )}
        </div>
      )}

      {activeTab === 'commits' && (
        <div style={styles.content}>
          {commits.length === 0 ? (
            <p style={styles.emptyMessage}>No commits yet</p>
          ) : (
            commits.map((commit) => (
              <div key={commit.id} style={styles.commitItem}>
                <div style={styles.commitHash}>{commit.hash.substring(0, 7)}</div>
                <div style={styles.commitInfo}>
                  <div style={styles.commitMessage}>{commit.message}</div>
                  <div style={styles.commitMeta}>
                    {commit.author} • {new Date(commit.created_at).toLocaleString()}
                  </div>
                </div>
              </div>
            ))
          )}
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
    color: '#58a6ff',
    marginBottom: '0.5rem',
  },
  repoDesc: {
    color: '#8b949e',
    fontSize: '1.1rem',
    marginBottom: '1rem',
  },
  repoMeta: {
    display: 'flex',
    gap: '2rem',
    color: '#8b949e',
    fontSize: '0.875rem',
  },
  tabs: {
    display: 'flex',
    gap: '0.5rem',
    marginBottom: '1rem',
    borderBottom: '1px solid #30363d',
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
  content: {
    backgroundColor: '#161b22',
    border: '1px solid #30363d',
    borderRadius: '6px',
    padding: '1.5rem',
  },
  emptyMessage: {
    textAlign: 'center',
    color: '#8b949e',
    padding: '2rem',
  },
  fileItem: {
    display: 'flex',
    justifyContent: 'space-between',
    padding: '1rem',
    borderBottom: '1px solid #30363d',
    cursor: 'pointer',
  },
  fileName: {
    color: '#58a6ff',
  },
  fileSize: {
    color: '#8b949e',
    fontSize: '0.875rem',
  },
  fileContent: {
    marginTop: '1rem',
    padding: '1rem',
    backgroundColor: '#0d1117',
    borderRadius: '6px',
  },
  fileContentTitle: {
    color: '#fff',
    marginBottom: '1rem',
  },
  code: {
    color: '#c9d1d9',
    overflow: 'auto',
    fontSize: '0.875rem',
    margin: 0,
  },
  commitItem: {
    display: 'flex',
    gap: '1rem',
    padding: '1rem',
    borderBottom: '1px solid #30363d',
  },
  commitHash: {
    backgroundColor: '#21262d',
    color: '#58a6ff',
    padding: '0.25rem 0.5rem',
    borderRadius: '6px',
    fontFamily: 'monospace',
    fontSize: '0.875rem',
    height: 'fit-content',
  },
  commitInfo: {
    flex: 1,
  },
  commitMessage: {
    color: '#fff',
    marginBottom: '0.5rem',
  },
  commitMeta: {
    color: '#8b949e',
    fontSize: '0.875rem',
  },
};

export default Repository;
