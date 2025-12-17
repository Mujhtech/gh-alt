import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import Navbar from './components/Navbar';
import Home from './pages/Home';
import Login from './pages/Login';
import Register from './pages/Register';
import Repository from './pages/Repository';
import Issues from './pages/Issues';

function App() {
  return (
    <AuthProvider>
      <Router>
        <div style={{ minHeight: '100vh', backgroundColor: '#0d1117' }}>
          <Navbar />
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            <Route path="/repository/:id" element={<Repository />} />
            <Route path="/repository/:id/issues" element={<Issues />} />
          </Routes>
        </div>
      </Router>
    </AuthProvider>
  );
}

export default App;
