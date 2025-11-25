.PHONY: all build build-frontend clean run dev install-deps

all: build

# Install Go dependencies
install-deps:
	go mod download

# Build frontend production bundle
build-frontend:
	cd frontend && npm run build

# Build standalone executable with embedded frontend
build: build-frontend
	go build -o filemanager .

# Run in development mode (requires separate frontend dev server)
dev:
	go run .

# Run frontend dev server
dev-frontend:
	cd frontend && npm run dev

# Clean build artifacts
clean:
	rm -f filemanager
	rm -rf frontend/dist

# Run the built executable
run: build
	./filemanager
