import axios from 'axios';

const API_BASE_URL = import.meta.env.DEV ? 'http://localhost:8080/api' : '/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
});

export const login = async (username, password) => {
  const response = await api.post('/login', { username, password });
  return response.data;
};

export const logout = async () => {
  const response = await api.post('/logout');
  return response.data;
};

export const listFiles = async (path = '.') => {
  const response = await api.get('/files', { params: { path } });
  return response.data;
};

export const uploadFile = async (path, file) => {
  const formData = new FormData();
  formData.append('file', file);
  formData.append('path', path);

  const response = await api.post('/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
  return response.data;
};

export const downloadFile = async (path) => {
  const response = await api.get('/download', {
    params: { path },
    responseType: 'blob',
  });

  // Create download link
  const url = window.URL.createObjectURL(new Blob([response.data]));
  const link = document.createElement('a');
  link.href = url;
  link.setAttribute('download', path.split('/').pop());
  document.body.appendChild(link);
  link.click();
  link.remove();
  window.URL.revokeObjectURL(url);
};

export const createDirectory = async (path, name) => {
  const response = await api.post('/mkdir', { path, name });
  return response.data;
};

export const deleteItem = async (path) => {
  const response = await api.delete('/delete', { data: { path } });
  return response.data;
};

export const changePassword = async (oldPassword, newPassword) => {
  const response = await api.post('/change-password', { oldPassword, newPassword });
  return response.data;
};

export default api;
