# doyo-img

[中文](README.md) | English

Lightweight, self-hosted image hosting service. No login required, single binary deployment, get image sharing links instantly.

## Features

- **Multiple Upload Methods** — Drag & drop, Ctrl+V paste, click to select, batch upload support
- **Instant Sharing** — Auto-generates URL / Markdown / HTML / BBCode formats after upload
- **No Login Required** — Works out of the box, anyone can upload and access images
- **Bilingual Support** — Built-in i18n, one-click language switch between Chinese and English
- **Light/Dark Themes** — Supports light / dark / system-following modes
- **Upload History** — Browser localStorage persistence, survives page refreshes
- **Real-time Progress** — Shows upload percentage in real-time
- **Image Info** — Displays filename, size, dimensions, format
- **Token-based Deletion** — Each image gets a unique deletion token for secure removal
- **IP Rate Limiting** — Token bucket algorithm for IP-level request throttling
- **CORS Support** — Configurable cross-origin policy for third-party embedding
- **Image Compression** — Optional JPEG/PNG compression
- **Thumbnails** — Auto-generated thumbnails for faster preview loading
- **S3 Storage** — Supports AWS S3, Alibaba Cloud OSS, Tencent Cloud COS, MinIO and other S3-compatible storage
- **Image Watermark** — Optional text/image watermark with custom font support (including Chinese TTF/OTF), position, and opacity
- **EXIF Auto-stripping** — Automatically removes EXIF metadata (GPS, device info) from JPEG images for privacy protection
- **SVG Security Sandbox** — Injects CSP headers for SVG responses to block embedded JavaScript execution, preventing Stored XSS
- **Link Expiration** — Optional expiration time on upload, returns 410 Gone after expiry
- **Expiration Cleanup** — Scheduled cleanup of expired images with global retention days and per-image expiration dual policy
- **Environment Variable Override** — YAML config + `DOYO_` prefix environment variables for flexible deployment
- **Single Binary Deployment** — Go embed for frontend, compiles to single binary
- **Docker Support** — Dockerfile and docker-compose configurations provided

## Architecture

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

**Backend**: Go + Gin framework, layered architecture (handler → service → storage)

**Frontend**: React 18 + TypeScript + Vite + Tailwind CSS + Zustand state management

**Build**: Frontend compiled to static files, embedded in Go binary via `go:embed`

## Directory Structure

```
doyo-img/
├── main.go                    # Entry point: init config, storage, router, start server
├── embed.go                   # go:embed frontend static resources
├── config.example.yaml        # Configuration example
├── Dockerfile                 # Docker multi-stage build
├── docker-compose.yml         # Docker Compose orchestration
├── Makefile                   # Common build commands
├── internal/
│   ├── config/config.go       # Config definition & loading (Viper)
│   ├── model/image.go         # Data models (ImageMeta, UploadResult)
│   ├── storage/
│   │   ├── storage.go         # Storage interface definition
│   │   ├── local.go           # Local filesystem storage implementation
│   │   └── s3.go              # S3-compatible storage implementation (AWS/OSS/COS/MinIO)
│   ├── service/image.go       # Business logic (upload/get/delete/watermark/EXIF)
│   ├── handler/
│   │   ├── upload.go          # POST /api/upload handler
│   │   ├── image.go           # Image access/info/delete/expiration check
│   │   └── config_handler.go  # GET /api/config public config
│   ├── processor/
│   │   ├── info.go            # Image info extraction (width/height/format)
│   │   ├── compress.go        # Image compression & thumbnail generation
│   │   ├── watermark.go       # Text/image watermark processing
│   │   └── exif.go            # JPEG EXIF metadata stripping
│   ├── middleware/
│   │   ├── cors.go            # CORS middleware
│   │   └── ratelimit.go       # IP rate limiting middleware
│   ├── router/router.go       # Route registration
│   ├── cleanup/cleanup.go     # Expired image scheduled cleanup
│   └── util/
│       ├── id.go              # Random ID / deletion token generation
│       ├── validate.go        # MIME detection / format validation
│       └── response.go        # Unified API response format
└── web/                       # React frontend
    ├── src/
    │   ├── main.tsx           # React entry
    │   ├── App.tsx            # Root component
    │   ├── api/client.ts      # API request wrapper
    │   ├── store/uploadStore.ts # Zustand global state
    │   ├── i18n/messages.ts   # Chinese/English translations
    │   ├── components/
    │   │   ├── layout/        # Header, Footer
    │   │   ├── upload/        # UploadZone, ExpirySelector
    │   │   ├── result/        # UploadResult (with expiration badge)
    │   │   └── common/        # Toast notifications
    │   ├── hooks/             # useUpload, useTheme, useCopy
    │   ├── utils/format.ts    # Link generation, file size formatting
    │   └── types/index.ts     # TypeScript type definitions
    └── vite.config.ts         # Vite config (with dev proxy)
```

## Quick Start

### Run Directly

```bash
# Download release binary, or build from source and run
./doyo-img

# Open http://localhost:9090 in browser
```

### Docker Deployment

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

### Build from Source

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

## Local Development

Frontend and backend development are separated, with frontend dev server proxying API requests to backend automatically.

**1. Configure Development Environment**:

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

**2. Start Backend** (Terminal 1):
```bash
go run .
# Backend listens on http://localhost:9090
```

**3. Start Frontend** (Terminal 2):
```bash
cd web
npm install
npm run dev
# Frontend dev server http://localhost:5173, API requests auto-proxy to :9090
```

Frontend Vite dev server is configured with proxy rules (`vite.config.ts`), `/api` and `/i` paths are automatically forwarded to backend.

## Docker Deployment

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
- `/app/data` — Image file storage directory, must be persisted
- `/app/config.yaml` — Config file, `:ro` read-only mount

## Configuration

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
    # S3-compatible service Endpoint (without protocol prefix, use_ssl controls http/https)
    # AWS S3:       s3.amazonaws.com
    # Alibaba OSS:  oss-cn-hangzhou.aliyuncs.com
    # Tencent COS:  cos.ap-guangzhou.myqcloud.com
    # MinIO:        localhost:9000
    endpoint: ""
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
    max_width: 300          # Thumbnail max width
    max_height: 300         # Thumbnail max height
  watermark:
    enabled: false          # Enable watermark
    type: "text"            # text or image
    text: "doyo-img"        # Text watermark content
    font_path: ""           # Custom font path (TTF/OTF), CJK fonts needed for Chinese
    font_size: 24
    opacity: 0.3            # Opacity 0.0 ~ 1.0
    position: "bottom-right" # top-left, top-right, bottom-left, bottom-right, center
    image_path: ""          # Image watermark PNG path (when type is image)
    padding: 20

id:
  length: 8                 # Image ID length
  alphabet: "0123456789abcdefghijklmnopqrstuvwxyz"  # ID character set

cors:
  allowed_origins:          # Allowed CORS origins
    - "*"                   # Production: set to specific domains

rate_limit:
  enabled: true             # Enable IP rate limiting
  requests_per_minute: 30   # Requests per minute limit
  burst: 10                 # Burst request allowance

cleanup:
  enabled: false            # Enable expiration cleanup
  retention_days: 30        # Global retention days, 0 = no time-based cleanup

log:
  level: "info"             # Log level: debug | info | warn | error
```

### Environment Variable Override

All config options can be overridden via `DOYO_` prefixed environment variables, using `_` as hierarchy separator:

| Config | Environment Variable | Example |
|--------|---------------------|---------|
| `server.port` | `DOYO_SERVER_PORT` | `8080` |
| `server.base_url` | `DOYO_SERVER_BASE_URL` | `https://img.example.com` |
| `storage.type` | `DOYO_STORAGE_TYPE` | `s3` |
| `storage.local.data_dir` | `DOYO_STORAGE_LOCAL_DATA_DIR` | `/data/images` |
| `upload.max_file_size` | `DOYO_UPLOAD_MAX_FILE_SIZE` | `10485760` |
| `rate_limit.enabled` | `DOYO_RATE_LIMIT_ENABLED` | `false` |
| `log.level` | `DOYO_LOG_LEVEL` | `debug` |

## API Documentation

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

### POST /api/upload — Upload Image

**Request**: `multipart/form-data`, field name `file` (supports multiple files)

Optional field `expire_hours`: Set image expiration time (hours), 0 or omitted means never expire.

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

> `expires_at` is only returned when expiration time is set.

---

### GET /i/{id}.{ext} — Direct Image Link

Returns image binary directly, usable in `<img>` tags and Markdown.

```
http://localhost:9090/i/a1b2c3d4.jpg          # Original
http://localhost:9090/i/a1b2c3d4.jpg?t=thumb  # Thumbnail
```

Response includes `Cache-Control: public, max-age=31536000, immutable`, browsers and CDNs can cache long-term.

SVG images additionally return `Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'; img-src data:; sandbox` header to prevent XSS attacks.

Expired images return `410 Gone`.

---

### GET /api/image/{id} — Get Image Info

```bash
curl http://localhost:9090/api/image/a1b2c3d4
```

Returns image metadata (ID, name, format, dimensions, size, creation time, etc.).

---

### DELETE /api/image/{id} — Delete Image

Requires deletion token returned during upload in request header:

```bash
curl -X DELETE -H "X-Delete-Token: tok_xxxxxxxxxxxxxxxxxxxxxxxx" \
  http://localhost:9090/api/image/a1b2c3d4
```

---

### GET /api/recent — Recent Uploads List

```bash
curl http://localhost:9090/api/recent?limit=20
```

Parameter `limit` range 1-50, default 20.

---

### GET /api/config — Get Public Config

```bash
curl http://localhost:9090/api/config
```

Returns frontend-required config info (file size limits, allowed formats, watermark status, expiration settings), without sensitive information.

## Production Deployment

### HTTPS Deployment (Nginx Reverse Proxy)

doyo-img itself listens on HTTP, production environments terminate TLS via Nginx:

**1. Configure `config.yaml`**:
```yaml
server:
  base_url: "https://img.example.com"
cors:
  allowed_origins:
    - "https://img.example.com"
    - "https://your-blog.com"   # Sites that need to embed images
```

**2. Nginx Config**:
```nginx
server {
    listen 443 ssl http2;
    server_name img.example.com;

    ssl_certificate     /etc/ssl/certs/img.example.com.pem;
    ssl_certificate_key /etc/ssl/private/img.example.com.key;

    client_max_body_size 10m;  # Match upload.max_file_size

    location / {
        proxy_pass http://127.0.0.1:9090;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Forwarded-Host $host;
    }

    # Image direct link cache policy (optional, server already sets Cache-Control)
    location /i/ {
        proxy_pass http://127.0.0.1:9090;
        proxy_set_header Host $host;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_valid 200 365d;
    }
}

server {
    listen 80;
    server_name img.example.com;
    return 301 https://$host$request_uri;
}
```

**3. HTTPS Deployment Checklist**:

| Config | Description |
|--------|-------------|
| `server.base_url` | **Required** Set to full HTTPS domain, otherwise generated links are HTTP |
| `cors.allowed_origins` | Set to specific HTTPS domain list, don't use `*` in production |
| Nginx `client_max_body_size` | Must be ≥ `upload.max_file_size`, otherwise large uploads are blocked by Nginx |
| Nginx `X-Forwarded-Proto` | Ensure it's passed, backend uses this to determine protocol for correct links |

### Apache Reverse Proxy

```apache
<VirtualHost *:443>
    ServerName img.example.com
    SSLEngine on
    SSLCertificateFile /etc/ssl/certs/img.example.com.pem
    SSLCertificateKeyFile /etc/ssl/private/img.example.com.key

    ProxyPreserveHost On
    ProxyPass / http://127.0.0.1:9090/
    ProxyPassReverse / http://127.0.0.1:9090/

    RequestHeader set X-Forwarded-Proto "https"
    RequestHeader set X-Forwarded-Host "img.example.com"
</VirtualHost>
```

### Security Hardening Recommendations

- Production: Set `cors.allowed_origins` to specific domain whitelist
- Enable rate limiting `rate_limit.enabled: true`, adjust parameters based on actual traffic
- SVG files have built-in `Content-Security-Policy` sandbox headers for XSS protection (automatic)
- JPEG uploads auto-strip EXIF privacy data (`strip_exif: true` default on)
- Regular backup of `data/` directory (local storage mode)
- Use `cleanup.enabled: true` to prevent unlimited storage growth

## FAQ

### Image links are HTTP instead of HTTPS after upload

Ensure one of the following:
1. Set `server.base_url: "https://your-domain.com"` in `config.yaml`
2. Reverse proxy correctly passes `X-Forwarded-Proto: https` header

### Nginx returns 413 Request Entity Too Large

Nginx's `client_max_body_size` defaults to 1MB, increase in Nginx config:
```nginx
client_max_body_size 10m;
```

### Image data lost in Docker container

Ensure persistent volume is mounted:
```yaml
volumes:
  - ./data:/app/data
```

### CORS requests blocked

Add request origin domain to `cors.allowed_origins` in `config.yaml`:
```yaml
cors:
  allowed_origins:
    - "https://your-blog.com"
    - "https://another-site.com"
```

### Rate limiting triggered (429 Too Many Requests)

Adjust rate limiting config or disable for trusted clients:
```yaml
rate_limit:
  requests_per_minute: 60
  burst: 20
```

## License

[MIT](LICENSE)
