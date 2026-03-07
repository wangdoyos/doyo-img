.PHONY: dev dev-frontend dev-backend build clean docker

# Development - run frontend and backend separately
dev-frontend:
	cd web && npm run dev

dev-backend:
	go run .

# Build production binary (frontend + embedded backend)
build:
	cd web && npm install && npm run build
	go build -ldflags="-s -w" -o doyo-img .

# Build for Windows
build-windows:
	cd web && npm install && npm run build
	set CGO_ENABLED=0&& go build -ldflags="-s -w" -o doyo-img.exe .

# Docker
docker:
	docker build -t doyo-img .

docker-run:
	docker-compose up -d

# Clean build artifacts
clean:
	rm -f doyo-img doyo-img.exe
	rm -rf web/dist web/node_modules
