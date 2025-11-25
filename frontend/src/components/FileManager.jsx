import { useState, useEffect } from 'react';
import * as api from '../api';
import './FileManager.css';

function FileManager({ onLogout }) {
    const [currentPath, setCurrentPath] = useState('.');
    const [files, setFiles] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [uploading, setUploading] = useState(false);

    useEffect(() => {
        loadFiles(currentPath);
    }, [currentPath]);

    const loadFiles = async (path) => {
        setLoading(true);
        setError('');
        try {
            const data = await api.listFiles(path);
            setFiles(data.files || []);
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to load files');
            if (err.response?.status === 401) {
                onLogout();
            }
        } finally {
            setLoading(false);
        }
    };

    const handleNavigate = (name) => {
        const newPath = currentPath === '.' ? name : `${currentPath}/${name}`;
        setCurrentPath(newPath);
    };

    const handleBreadcrumbClick = (index) => {
        if (index === -1) {
            setCurrentPath('.');
            return;
        }
        const parts = currentPath.split('/');
        const newPath = parts.slice(0, index + 1).join('/');
        setCurrentPath(newPath);
    };

    const handleUpload = async (e) => {
        const file = e.target.files[0];
        if (!file) return;

        setUploading(true);
        setError('');
        try {
            await api.uploadFile(currentPath, file);
            await loadFiles(currentPath);
            e.target.value = '';
        } catch (err) {
            setError(err.response?.data?.error || 'Upload failed');
        } finally {
            setUploading(false);
        }
    };

    const handleDownload = async (name) => {
        const filePath = currentPath === '.' ? name : `${currentPath}/${name}`;
        try {
            await api.downloadFile(filePath);
        } catch (err) {
            setError(err.response?.data?.error || 'Download failed');
        }
    };

    const handleCreateFolder = async () => {
        const name = prompt('Enter folder name:');
        if (!name) return;

        setError('');
        try {
            await api.createDirectory(currentPath, name);
            await loadFiles(currentPath);
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to create folder');
        }
    };

    const handleDelete = async (name) => {
        if (!confirm(`Are you sure you want to delete "${name}"?`)) return;

        const itemPath = currentPath === '.' ? name : `${currentPath}/${name}`;
        setError('');
        try {
            await api.deleteItem(itemPath);
            await loadFiles(currentPath);
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to delete');
        }
    };

    const handleLogout = async () => {
        try {
            await api.logout();
        } catch (err) {
            console.error('Logout error:', err);
        }
        onLogout();
    };

    const getBreadcrumbs = () => {
        if (currentPath === '.') return [];
        return currentPath.split('/');
    };

    const formatSize = (bytes) => {
        if (bytes === 0) return '0 B';
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
    };

    const formatDate = (dateString) => {
        const date = new Date(dateString);
        return date.toLocaleString();
    };

    return (
        <div className="file-manager">
            <header className="fm-header">
                <div className="fm-header-content">
                    <h1>📁 File Manager</h1>
                    <button onClick={handleLogout} className="logout-button">
                        Logout
                    </button>
                </div>
            </header>

            <div className="fm-toolbar">
                <div className="breadcrumb">
                    <span onClick={() => handleBreadcrumbClick(-1)} className="breadcrumb-item">
                        🏠 Home
                    </span>
                    {getBreadcrumbs().map((part, index) => (
                        <span key={index}>
                            <span className="breadcrumb-separator">/</span>
                            <span onClick={() => handleBreadcrumbClick(index)} className="breadcrumb-item">
                                {part}
                            </span>
                        </span>
                    ))}
                </div>

                <div className="toolbar-actions">
                    <button onClick={handleCreateFolder} className="action-button">
                        ➕ New Folder
                    </button>
                    <label className="action-button upload-button">
                        {uploading ? '⏳ Uploading...' : '⬆️ Upload'}
                        <input
                            type="file"
                            onChange={handleUpload}
                            disabled={uploading}
                            style={{ display: 'none' }}
                        />
                    </label>
                </div>
            </div>

            {error && (
                <div className="error-banner">
                    {error}
                    <button onClick={() => setError('')} className="close-error">×</button>
                </div>
            )}

            <div className="fm-content">
                {loading ? (
                    <div className="loading">Loading...</div>
                ) : files.length === 0 ? (
                    <div className="empty-state">
                        <p>📂 This folder is empty</p>
                    </div>
                ) : (
                    <div className="file-grid">
                        {files.map((file, index) => (
                            <div key={index} className="file-item">
                                <div
                                    className="file-icon"
                                    onClick={() => file.isDir ? handleNavigate(file.name) : handleDownload(file.name)}
                                >
                                    {file.isDir ? '📁' : '📄'}
                                </div>
                                <div className="file-info">
                                    <div
                                        className="file-name"
                                        onClick={() => file.isDir ? handleNavigate(file.name) : handleDownload(file.name)}
                                        title={file.name}
                                    >
                                        {file.name}
                                    </div>
                                    <div className="file-meta">
                                        {!file.isDir && <span>{formatSize(file.size)}</span>}
                                        <span>{formatDate(file.modTime)}</span>
                                    </div>
                                </div>
                                <button
                                    onClick={() => handleDelete(file.name)}
                                    className="delete-button"
                                    title="Delete"
                                >
                                    🗑️
                                </button>
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
}

export default FileManager;
