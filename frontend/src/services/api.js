import axios from 'axios';

const API_URL = 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add token to requests if available
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const authService = {
  register: async (username, email, password) => {
    const response = await api.post('/auth/register', { username, email, password });
    return response.data;
  },
  login: async (username, password) => {
    const response = await api.post('/auth/login', { username, password });
    return response.data;
  },
};

export const repositoryService = {
  getAll: async () => {
    const response = await api.get('/repositories');
    return response.data;
  },
  getById: async (id) => {
    const response = await api.get(`/repositories/${id}`);
    return response.data;
  },
  create: async (name, description, isPrivate = false) => {
    const response = await api.post('/repositories', { name, description, is_private: isPrivate });
    return response.data;
  },
  getFiles: async (id) => {
    const response = await api.get(`/repositories/${id}/files`);
    return response.data;
  },
  getCommits: async (id) => {
    const response = await api.get(`/repositories/${id}/commits`);
    return response.data;
  },
  createCommit: async (id, message, author, files) => {
    const response = await api.post(`/repositories/${id}/commit`, { message, author, files });
    return response.data;
  },
};

export default api;
