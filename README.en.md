<div align="center">

# 🖼️ doyo-img

**Lightweight, Self-Hosted Image Hosting Service**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://golang.org)
[![React](https://img.shields.io/badge/React-18-61DAFB?logo=react)](https://reactjs.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker)](https://www.docker.com)

[中文](README.md) | 🌐 **English**

**No login required, single binary deployment, get image sharing links instantly**

</div>

---

## ✨ Features

### 🚀 Core Features
- 📤 **Multiple Upload Methods** — Drag & drop, Ctrl+V paste, click to select, batch upload support
- 🔗 **Instant Sharing** — Auto-generates URL / Markdown / HTML / BBCode formats after upload
- 🔓 **No Login Required** — Works out of the box, anyone can upload and access images
- 🌐 **Bilingual Support** — Built-in i18n, one-click language switch between Chinese and English
- 🎨 **Light/Dark Themes** — Supports light / dark / system-following modes

### 💾 Storage & Management
- 📜 **Upload History** — Browser localStorage persistence, survives page refreshes
- 📊 **Real-time Progress** — Shows upload percentage in real-time
- ℹ️ **Image Info** — Displays filename, size, dimensions, format
- 🗑️ **Token-based Deletion** — Each image gets a unique deletion token for secure removal

### 🔒 Security & Privacy
- 🛡️ **IP Rate Limiting** — Token bucket algorithm for IP-level request throttling
- 🌐 **CORS Support** — Configurable cross-origin policy for third-party embedding
- 🔍 **EXIF Auto-stripping** — Automatically removes EXIF metadata from JPEG images for privacy protection
- 🔒 **SVG Security Sandbox** — Injects CSP headers for SVG responses to prevent Stored XSS

### 🎨 Image Processing
- 🗜️ **Image Compression** — Optional JPEG/PNG compression
- 🖼️ **Thumbnails** — Auto-generated thumbnails for faster preview loading
- 💧 **Image Watermark** — Optional text/image watermark with custom font support, position, and opacity
- ⏰ **Link Expiration** — Optional expiration time on upload, returns 410 Gone after expiry
- 🧹 **Expiration Cleanup** — Scheduled cleanup of expired images with dual policy support

### ☁️ Deployment & Storage
- ☁️ **S3 Storage** — Supports AWS S3, Alibaba Cloud OSS, Tencent Cloud COS, MinIO and other S3-compatible storage
- ⚙️ **Environment Variable Override** — YAML config + `DOYO_` prefix environment variables for flexible deployment
- 📦 **Single Binary Deployment** — Go embed for frontend, compiles to single binary
- 🐳 **Docker Support** — Dockerfile and docker-compose configurations provided

---

## 🏗️ Architecture

```
┌───────────────────────────────────────────────┐
│                   Client Browser                │
│  React 18 + TypeScript + Tailwind CSS + Zustand │
└──────────────────────┬────────────────────────┘
                       │ HTTP API
┌──────────────────────▼────────────────────────┐
│                 Go Backend (Gin)                │
│  ┌──────────┐ ┌──────────┐ ┌───────────────┐  │
│  │  Router   │ │Middleware│ │   Handler     │  │
│  │          │ │CORS/Rate │ │Upload/Image/  │  │
│  │  /api/*  │ │ Limiter  │ │Config         │  │
│  └──────────┘ └──────────┘ └───────┬───────┘  │
│                                     │          │
│  ┌──────────┐ ┌──────────┐ ┌───────▼───────┐  │
│  │Processor │ │   Util   │ │   Service     │  │
│  │Compress/ │ │ID/Validate│ │  ImageService │  │
│  │Thumbnail/│ │/Response │ │               │  │
│  │Watermark/│ └──────────┘ └───────┬───────┘  │
│  │EXIF Strip│                      │          │
│  └──────────┘                      │          │
│                                     │          │
│  ┌──────────────────────────────────▼───────┐  │
│  │         Storage Interface                │  │
│  │    LocalStorage  |  S3Storage            │  │
│  └──────────────────────────────────────────┘  │
└───────────────────────────────────────────────┘
```

| Layer | Tech Stack |
|-------|------------|
| **Backend** | Go + Gin framework, layered architecture (handler → service → storage) |
| **Frontend** | React 18 + TypeScript + Vite + Tailwind CSS + Zustand |
| **Build** | Frontend compiled to static files, embedded in Go binary via `go:embed` |

---

## 📁 Directory Structure

```
doyo-img/
├── 📄 main.go                    # Entry point: init config, storage, router, start server
├── 📄 embed.go                   # go:embed frontend static resources
├── 📄 config.example.yaml        # Configuration example
├── 🐳 Dockerfile                 # Docker multi-stage build
├── 📄 docker-compose.yml         # Docker Compose orchestration
├── 📄 Makefile                   # Common build commands
├── 📁 internal/
│   ├── 📁 config/                # Config definition & loading (Viper)
│   ├── 📁 model/                 # Data models
│   ├── 📁 storage/               # Storage interface & implementations
│   ├── 📁 service/               # Business logic layer
│   ├── 📁 handler/               # HTTP request handlers
│   ├── 📁 processor/             # Image processing (compress/watermark/EXIF)
│   ├── 📁 middleware/            # Middleware (CORS/rate limiting)
│   ├── 📁 router/                # Route registration
│   ├── 📁 cleanup/               # Expired image scheduled cleanup
│   └── 📁 util/                  # Utility functions
└── 📁 web/                       # React frontend
    ├── 📁 src/
    │   ├── 📁 components/        # Components
    │   ├── 📁 hooks/             # Custom Hooks
    │   ├── 📁 i18n/              # Internationalization
    │   ├── 📁 store/             # State management
    │   └── 📁 utils/             # Utility functions
    └── ⚙️ vite.config.ts         # Vite config
```

---

## 🚀 Quick Start

### ⚡ Run Directly

```bash
# Download release binary, or build from source and run
./doyo-img

# Open http://localhost:9090 in browser
```

### 🐳 Docker Deployment

```bash
# Clone project
git clone https://github.com/wangdoyos/doyo-img.git
cd doyo-img

# Copy and edit configuration
cp config.example.yaml config.yaml

# Start
docker-compose up -d

# Open http://localhost:9090 in browser
```

### 🔨 Build from Source

**Prerequisites**: Go 1.22+, Node.js 18+

```bash
# Clone project
git clone https://github.com/wangdoyos/doyo-img.git
cd doyo-img

# Option 1: Using Makefile
make build           # Linux/macOS
make build-windows   # Windows

# Option 2: Manual build
cd web && npm install && npm run build && cd ..
go build -ldflags="-s -w" -o doyo-img .

# Run
./doyo-img
```

---

## 💻 Local Development

Frontend and backend development are separated, with frontend dev server proxying API requests to backend automatically.

### 1️⃣ Configure Development Environment

The program looks for `config.yaml` in project root and `./config` directory. Two config templates are provided:

| File | Purpose |
|------|---------|
| `config.example.yaml` | Production template (rate limiting on, info log level) |
| `config.dev.yaml` | Development template (rate limiting off, debug log, relaxed upload limits) |

```bash
# Development: copy dev config
cp config.dev.yaml config.yaml

# Production: copy production config
cp config.example.yaml config.yaml
```

If no `config.yaml` is created, the program uses built-in defaults (equivalent to `config.example.yaml`).

### 2️⃣ Start Backend (Terminal 1)
```bash
go run .
# Backend listens on http://localhost:9090
```

### 3️⃣ Start Frontend (Terminal 2)
```bash
cd web
npm install
npm run dev
# Frontend dev server http://localhost:5173, API requests auto-proxy to :9090
```

---

## 🐳 Docker Deployment

### Dockerfile

Multi-stage build, final image based on Alpine for small size:

```dockerfile
# Build frontend → Build backend → Copy only binary to runtime image
docker build -t doyo-img .
```

### docker-compose.yml

```yaml
services:
  doyo-img:
    build: .
    ports:
      - "9090:9090"
    volumes:
      - ./data:/app/data          # Image persistent storage
      - ./config.yaml:/app/config.yaml:ro  # Config file mount
    restart: unless-stopped
    environment:
      - DOYO_SERVER_PORT=9090
      # For HTTPS reverse proxy deployment:
      # - DOYO_SERVER_BASE_URL=https://img.example.com
```

**Key Mount Points**:
- 📁 `/app/data` — Image file storage directory, must be persisted
- 📄 `/app/config.yaml` — Config file, `:ro` read-only mount

---

## ⚙️ Configuration

Copy `config.example.yaml` to `config.yaml` and modify as needed:

```yaml
server:
  host: "0.0.0.0"          # Listen address
  port: 9090                # Listen port
  base_url: ""              # External access base URL (required for HTTPS deployment)
                            # Example: "https://img.example.com"

storage:
  type: "local"             # Storage type: local | s3
  local:
    data_dir: "./data"      # Local storage directory
  s3:                       # S3-compatible storage (when type is "s3")
    endpoint: ""            # S3 Endpoint
    bucket: ""
    region: ""
    access_key: ""
    secret_key: ""
    use_ssl: true
    path_prefix: "images"

upload:
  max_file_size: 5242880    # Max file size in bytes (default 5MB)
  max_batch_size: 10        # Max files per upload
  allowed_formats:          # Allowed image formats
    - "jpg"
    - "jpeg"
    - "png"
    - "gif"
    - "webp"
    - "svg"
  default_expire_hours: 0   # Default expiration (hours), 0 = never expire
  max_expire_days: 0        # Max expiration days limit, 0 = no limit

processing:
  compress_enabled: false   # Enable upload compression
  compress_quality: 85      # JPEG compression quality (1-100)
  strip_exif: true          # Auto-strip EXIF metadata from JPEG
  thumbnail:
    enabled: true           # Generate thumbnails
    max_width: 300
    max_height: 300
  watermark:
    enabled: false          # Enable watermark
    type: "text"            # text | image
    text: "doyo-img"        # Text watermark content
    font_path: ""           # Custom font path (TTF/OTF)
    font_size: 24
    opacity: 0.3            # Opacity 0.0 ~ 1.0
    position: "bottom-right" # Position options
    image_path: ""          # Image watermark PNG path
    padding: 20

id:
  length: 8                 # Image ID length
  alphabet: "0123456789abcdefghijklmnopqrstuvwxyz"

cors:
  allowed_origins:          # Allowed CORS origins
    - "*"

rate_limit:
  enabled: true             # Enable IP rate limiting
  requests_per_minute: 30
  burst: 10

cleanup:
  enabled: false            # Enable expiration cleanup
  retention_days: 30

log:
  level: "info"             # Log level: debug | info | warn | error
```

### 🔧 Environment Variable Override

All config options can be overridden via `DOYO_` prefixed environment variables:

| Config | Environment Variable | Example |
|--------|---------------------|---------|
| `server.port` | `DOYO_SERVER_PORT` | `8080` |
| `server.base_url` | `DOYO_SERVER_BASE_URL` | `https://img.example.com` |
| `storage.type` | `DOYO_STORAGE_TYPE` | `s3` |
| `upload.max_file_size` | `DOYO_UPLOAD_MAX_FILE_SIZE` | `10485760` |

---

## 📡 API Documentation

### Unified Response Format

```json
{
  "code": 0,
  "data": { ... },
  "message": "ok"
}
```

`code = 0` indicates success, non-zero indicates error.

---

### 📤 POST /api/upload — Upload Image

**Request**: `multipart/form-data`, field name `file` (supports multiple files)

Optional field `expire_hours`: Set image expiration time (hours)

```bash
# Single file upload
curl -X POST -F "file=@photo.jpg" http://localhost:9090/api/upload

# Multiple files upload
curl -X POST -F "file=@a.jpg" -F "file=@b.png" http://localhost:9090/api/upload

# Upload with expiration (24 hours)
curl -X POST -F "file=@photo.jpg" -F "expire_hours=24" http://localhost:9090/api/upload
```

**Success Response**:
```json
{
  "code": 0,
  "data": {
    "images": [
      {
        "id": "a1b2c3d4",
        "name": "photo.jpg",
        "url": "http://localhost:9090/i/a1b2c3d4.jpg",
        "thumbnail_url": "http://localhost:9090/i/a1b2c3d4.jpg?t=thumb",
        "size": 102400,
        "width": 1920,
        "height": 1080,
        "format": "jpg",
        "delete_token": "tok_xxxxxxxxxxxxxxxxxxxxxxxx",
        "created_at": "2026-03-07T10:00:00Z",
        "expires_at": "2026-03-08T10:00:00Z"
      }
    ]
  },
  "message": "ok"
}
```

---

### 🖼️ GET /i/{id}.{ext} — Direct Image Link

```
http://localhost:9090/i/a1b2c3d4.jpg          # Original
http://localhost:9090/i/a1b2c3d4.jpg?t=thumb  # Thumbnail
```

- Response includes `Cache-Control: public, max-age=31536000, immutable`
- SVG images additionally return CSP header to prevent XSS
- Expired images return `410 Gone`

---

### ℹ️ GET /api/image/{id} — Get Image Info

```bash
curl http://localhost:9090/api/image/a1b2c3d4
```

---

### 🗑️ DELETE /api/image/{id} — Delete Image

```bash
curl -X DELETE -H "X-Delete-Token: tok_xxxxxxxxxxxxxxxxxxxxxxxx" \
  http://localhost:9090/api/image/a1b2c3d4
```

---

### 📜 GET /api/recent — Recent Uploads List

```bash
curl http://localhost:9090/api/recent?limit=20
```

---

### ⚙️ GET /api/config — Get Public Config

```bash
curl http://localhost:9090/api/config
```

---

## 🚀 Production Deployment

### 🔒 HTTPS Deployment (Nginx Reverse Proxy)

**1️⃣ Configure `config.yaml`**:
```yaml
server:
  base_url: "https://img.example.com"
cors:
  allowed_origins:
    - "https://img.example.com"
    - "https://your-blog.com"
```

**2️⃣ Nginx Config**:
```nginx
server {
    listen 443 ssl http2;
    server_name img.example.com;

    ssl_certificate     /etc/ssl/certs/img.example.com.pem;
    ssl_certificate_key /etc/ssl/private/img.example.com.key;

    client_max_body_size 10m;

    location / {
        proxy_pass http://127.0.0.1:9090;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Forwarded-Host $host;
    }
}

server {
    listen 80;
    server_name img.example.com;
    return 301 https://$host$request_uri;
}
```

**3️⃣ HTTPS Deployment Checklist**:

| Config | Description |
|--------|-------------|
| `server.base_url` | **Required** Set to full HTTPS domain |
| `cors.allowed_origins` | Set to specific HTTPS domain list, don't use `*` in production |
| Nginx `client_max_body_size` | Must be ≥ `upload.max_file_size` |

---

## 🛡️ Security Hardening Recommendations

- ✅ Set `cors.allowed_origins` to specific domain whitelist in production
- ✅ Enable rate limiting `rate_limit.enabled: true`
- ✅ SVG files have built-in CSP sandbox headers for XSS protection (automatic)
- ✅ JPEG uploads auto-strip EXIF privacy data
- ✅ Regular backup of `data/` directory
- ✅ Use `cleanup.enabled: true` to prevent unlimited storage growth

---

## ❓ FAQ

### Q: Image links are HTTP instead of HTTPS after upload
**A:** Ensure one of the following:
1. Set `server.base_url: "https://your-domain.com"` in `config.yaml`
2. Reverse proxy correctly passes `X-Forwarded-Proto: https` header

### Q: Nginx returns 413 Request Entity Too Large
**A:** Increase in Nginx config:
```nginx
client_max_body_size 10m;
```

### Q: Image data lost in Docker container
**A:** Ensure persistent volume is mounted:
```yaml
volumes:
  - ./data:/app/data
```

### Q: CORS requests blocked
**A:** Add request origin domain to `cors.allowed_origins` in `config.yaml`

---

## 📋 Release Roadmap

| Version | Planned Date | Major Features |
|---------|-------------|----------------|
| v1.0.0 | Mar 2026 | 🎉 First stable release |
| v1.1.0 | Apr 2026 | 📊 Admin dashboard, analytics panel |
| v1.2.0 | May 2026 | 🔐 User system, private image hosting |
| v2.0.0 | Q3 2026 | 🚀 Distributed deployment, cluster support |

---

## 📄 License

[MIT](LICENSE) © 2026 doyo-img
