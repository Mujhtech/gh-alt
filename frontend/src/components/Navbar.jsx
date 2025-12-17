import React from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

const Navbar = () => {
  const { user, logout, isAuthenticated } = useAuth();

  return (
    <nav style={styles.nav}>
      <div style={styles.container}>
        <Link to="/" style={styles.logo}>
          GH Alt
        </Link>
        <div style={styles.menu}>
          <Link to="/" style={styles.link}>Repositories</Link>
          {isAuthenticated ? (
            <>
              <span style={styles.username}>{user?.username}</span>
              <button onClick={logout} style={styles.button}>Logout</button>
            </>
          ) : (
            <>
              <Link to="/login" style={styles.link}>Login</Link>
              <Link to="/register" style={styles.link}>Register</Link>
            </>
          )}
        </div>
      </div>
    </nav>
  );
};

const styles = {
  nav: {
    backgroundColor: '#24292e',
    padding: '1rem 0',
    borderBottom: '1px solid #444',
  },
  container: {
    maxWidth: '1200px',
    margin: '0 auto',
    padding: '0 1rem',
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  logo: {
    color: '#fff',
    fontSize: '1.5rem',
    fontWeight: 'bold',
    textDecoration: 'none',
  },
  menu: {
    display: 'flex',
    gap: '1.5rem',
    alignItems: 'center',
  },
  link: {
    color: '#fff',
    textDecoration: 'none',
  },
  username: {
    color: '#58a6ff',
  },
  button: {
    backgroundColor: '#238636',
    color: '#fff',
    border: 'none',
    padding: '0.5rem 1rem',
    borderRadius: '6px',
    cursor: 'pointer',
  },
};

export default Navbar;
