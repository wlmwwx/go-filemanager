# Web File Manager

A standalone web-based file manager built with Go and React. Supports file upload/download, directory operations, and basic authentication.

## Features

- 🔐 **Authentication** - Basic login system with session management
- 📁 **File Operations** - Upload, download, and delete files
- 📂 **Directory Management** - Create and delete directories, navigate folder structure
- 🎨 **Modern UI** - Beautiful, responsive interface with smooth animations
- 📦 **Standalone** - Single executable with embedded React frontend
- 🔒 **Security** - Path traversal protection, operations restricted to working directory

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Node.js 18 or higher
- npm

### Build

```bash
# Install Go dependencies
go mod download

# Install frontend dependencies and build
cd frontend && npm install && cd ..

# Build standalone executable
make build
```

### Run

```bash
# Run the built executable
./filemanager

# Or use make
make run
```

The server will start on `http://localhost:8080`

**Default credentials:** `admin` / `admin123`

## Configuration

The application can be configured using a `config.yaml` file in the working directory.

### Configuration File

Create a `config.yaml` file (or copy from `config.yaml.example`):

```yaml
server:
  port: 8080
  rootDir: ""  # Leave empty to use current working directory, or specify an absolute path

secret: your-secret-key-change-this-in-production

log:
  file: ""  # Leave empty for stdout only, or specify a file path like "filemanager.log"
  level: info  # Options: info, debug
```

### Configuration Options

- **server.port**: Port number for the HTTP server (default: 8080)
- **server.rootDir**: Root directory for file operations (empty string uses current working directory)
- **secret**: Secret key for session encryption (change this in production!)
- **log.file**: Path to log file (empty string for stdout only)
- **log.level**: Logging level - `info` or `debug` (debug includes file names and line numbers)

If `config.yaml` is not found, the application will use default values.

## Development

### Backend Development

```bash
# Run Go server (requires frontend to be built first)
make dev
```

### Frontend Development

```bash
# Run frontend dev server with hot reload
cd frontend
npm run dev
```

The frontend dev server will run on `http://localhost:5173` and proxy API requests to the Go backend on port 8080.

### Full Development Setup

1. Terminal 1: Run Go backend
   ```bash
   make dev
   ```

2. Terminal 2: Run frontend dev server
   ```bash
   make dev-frontend
   ```

## Project Structure

```
.
├── main.go              # Application entry point
├── handlers.go          # API endpoint handlers
├── middleware.go        # Authentication and CORS middleware
├── models.go            # Data structures
├── go.mod              # Go dependencies
├── Makefile            # Build automation
├── frontend/           # React frontend
│   ├── src/
│   │   ├── components/
│   │   │   ├── Login.jsx
│   │   │   ├── Login.css
│   │   │   ├── FileManager.jsx
│   │   │   └── FileManager.css
│   │   ├── App.jsx
│   │   ├── api.js
│   │   └── index.css
│   ├── package.json
│   └── vite.config.js
└── README.md
```

## API Endpoints

### Authentication

- `POST /api/login` - Login with username and password
- `POST /api/logout` - Logout and clear session

### File Operations (Authenticated)

- `GET /api/files?path=<path>` - List directory contents
- `POST /api/upload` - Upload file (multipart/form-data)
- `GET /api/download?path=<path>` - Download file
- `POST /api/mkdir` - Create directory
- `DELETE /api/delete` - Delete file or directory

## Security Considerations

### Current Implementation

- Simple username/password authentication (hardcoded)
- Session-based authentication using cookies
- Path traversal protection
- Operations restricted to working directory

### Production Recommendations

1. **Change default credentials** - Modify the hardcoded credentials in `handlers.go` or implement environment variable configuration
2. **Use HTTPS** - Deploy behind a reverse proxy with SSL/TLS
3. **Implement proper password hashing** - Use bcrypt or similar for password storage
4. **Add rate limiting** - Prevent brute force attacks
5. **Configure session secret** - Change the session secret key in `main.go`
6. **Add file type restrictions** - Limit uploadable file types if needed
7. **Implement file size limits** - Already set to 32MB, adjust as needed

## Building for Production

```bash
# Build optimized production executable
make build

# The executable will be created as 'filemanager'
# Deploy it to your server and run from the desired root directory
```

The working directory where you run the executable becomes the file manager's root directory.

## Clean Up

```bash
# Remove build artifacts
make clean
```

## License

MIT

## Author

Created with Go and React
