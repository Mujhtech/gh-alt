import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { authService } from '../services/api';
import { useAuth } from '../context/AuthContext';

const Login = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();
  const { login } = useAuth();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    try {
      const response = await authService.login(username, password);
      login(response.token, response.user);
      navigate('/');
    } catch (err) {
      setError('Invalid credentials');
    }
  };

  return (
    <div style={styles.container}>
      <div style={styles.formContainer}>
        <h1 style={styles.title}>Login to GH Alt</h1>
        {error && <div style={styles.error}>{error}</div>}
        <form onSubmit={handleSubmit}>
          <input
            type="text"
            placeholder="Username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            style={styles.input}
            required
          />
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            style={styles.input}
            required
          />
          <button type="submit" style={styles.button}>
            Login
          </button>
        </form>
      </div>
    </div>
  );
};

const styles = {
  container: {
    maxWidth: '400px',
    margin: '4rem auto',
    padding: '0 1rem',
  },
  formContainer: {
    backgroundColor: '#161b22',
    border: '1px solid #30363d',
    borderRadius: '6px',
    padding: '2rem',
  },
  title: {
    color: '#fff',
    textAlign: 'center',
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
  button: {
    width: '100%',
    backgroundColor: '#238636',
    color: '#fff',
    border: 'none',
    padding: '0.75rem',
    borderRadius: '6px',
    cursor: 'pointer',
    fontSize: '1rem',
    marginTop: '0.5rem',
  },
  error: {
    backgroundColor: '#da3633',
    color: '#fff',
    padding: '0.75rem',
    borderRadius: '6px',
    marginBottom: '1rem',
  },
};

export default Login;
