# doyo-img

中文 | [English](README.en.md)

轻量级、自托管图床服务。无需登录，单文件部署，即刻获取图片分享链接。

## 功能特性

- **多种上传方式** — 拖拽上传、Ctrl+V 粘贴、点击选择，支持批量上传
- **即时分享** — 上传后自动生成 URL / Markdown / HTML / BBCode 多种链接格式
- **无需登录** — 开箱即用，任何人都可以上传和访问图片
- **中英文切换** — 内置国际化支持，一键切换界面语言
- **明暗主题** — 支持亮色 / 暗色 / 跟随系统三种主题模式
- **上传历史** — 基于浏览器 localStorage 持久化，刷新不丢失
- **实时进度** — 上传过程实时显示进度百分比
- **图片信息** — 展示文件名、尺寸、分辨率、格式
- **令牌删除** — 每张图片生成唯一删除令牌，凭令牌安全删除
- **IP 限流** — 基于令牌桶算法的 IP 级别请求限流
- **CORS 跨域** — 可配置的跨域策略，支持第三方网站嵌入
- **图片压缩** — 可选的 JPEG/PNG 压缩处理
- **缩略图** — 自动生成缩略图，加速预览加载
- **S3 存储** — 支持 AWS S3、阿里云 OSS、腾讯云 COS、MinIO 等 S3 兼容对象存储
- **图片水印** — 可选的文字 / 图片水印，支持自定义字体（含中文 TTF/OTF）、位置、透明度
- **EXIF 自动剥离** — 自动移除 JPEG 图片中的 EXIF 元数据（GPS、设备信息等），保护隐私
- **SVG 安全沙箱** — 为 SVG 响应注入 CSP 头，阻止嵌入的 JavaScript 执行，防范 Stored XSS
- **链接过期** — 上传时可选设置过期时间，过期后自动返回 410 Gone
- **过期清理** — 定时清理过期图片，支持全局保留天数和单图过期时间双重策略
- **环境变量覆盖** — YAML 配置 + `DOYO_` 前缀环境变量，灵活部署
- **单文件部署** — Go embed 嵌入前端，编译产出单个二进制
- **Docker 支持** — 提供 Dockerfile 和 docker-compose 配置

## 技术架构

```
┌───────────────────────────────────────────────┐
│                   客户端浏览器                    │
│  React 18 + TypeScript + Tailwind CSS + Zustand │
└──────────────────────┬────────────────────────┘
                       │ HTTP API
┌──────────────────────▼────────────────────────┐
│                 Go 后端 (Gin)                   │
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

**后端**: Go + Gin 框架，分层架构 (handler → service → storage)

**前端**: React 18 + TypeScript + Vite + Tailwind CSS + Zustand 状态管理

**构建**: 前端编译为静态文件，通过 `go:embed` 嵌入 Go 二进制

## 目录结构

```
doyo-img/
├── main.go                    # 程序入口：初始化配置、存储、路由、启动服务
├── embed.go                   # go:embed 嵌入前端静态资源
├── config.example.yaml        # 配置文件示例
├── Dockerfile                 # Docker 多阶段构建
├── docker-compose.yml         # Docker Compose 编排
├── Makefile                   # 常用构建命令
├── internal/
│   ├── config/config.go       # 配置定义与加载 (Viper)
│   ├── model/image.go         # 数据模型 (ImageMeta, UploadResult)
│   ├── storage/
│   │   ├── storage.go         # Storage 接口定义
│   │   ├── local.go           # 本地文件系统存储实现
│   │   └── s3.go              # S3 兼容对象存储实现 (AWS/OSS/COS/MinIO)
│   ├── service/image.go       # 业务逻辑层 (上传/获取/删除/水印/EXIF)
│   ├── handler/
│   │   ├── upload.go          # POST /api/upload 上传处理
│   │   ├── image.go           # 图片访问/信息/删除/过期检查
│   │   └── config_handler.go  # GET /api/config 公开配置
│   ├── processor/
│   │   ├── info.go            # 图片信息提取 (宽高/格式)
│   │   ├── compress.go        # 图片压缩与缩略图生成
│   │   ├── watermark.go       # 文字/图片水印处理
│   │   └── exif.go            # JPEG EXIF 元数据剥离
│   ├── middleware/
│   │   ├── cors.go            # CORS 跨域中间件
│   │   └── ratelimit.go       # IP 限流中间件
│   ├── router/router.go       # 路由注册
│   ├── cleanup/cleanup.go     # 过期图片定时清理
│   └── util/
│       ├── id.go              # 随机 ID / 删除令牌生成
│       ├── validate.go        # MIME 检测 / 格式校验
│       └── response.go        # 统一 API 响应格式
└── web/                       # React 前端
    ├── src/
    │   ├── main.tsx           # React 入口
    │   ├── App.tsx            # 根组件
    │   ├── api/client.ts      # API 请求封装
    │   ├── store/uploadStore.ts # Zustand 全局状态
    │   ├── i18n/messages.ts   # 中英文翻译定义
    │   ├── components/
    │   │   ├── layout/        # Header, Footer
    │   │   ├── upload/        # UploadZone 上传区域, ExpirySelector 过期选择器
    │   │   ├── result/        # UploadResult 结果展示 (含过期时间徽章)
    │   │   └── common/        # Toast 提示
    │   ├── hooks/             # useUpload, useTheme, useCopy
    │   ├── utils/format.ts    # 链接生成、文件大小格式化
    │   └── types/index.ts     # TypeScript 类型定义
    └── vite.config.ts         # Vite 配置 (含开发代理)
```

## 快速开始

### 直接运行

```bash
# 下载发布版本的二进制文件，或从源码构建后直接运行
./doyo-img

# 浏览器访问 http://localhost:9090
```

### Docker 部署

```bash
# 克隆项目
git clone https://github.com/doyo-img/doyo-img.git
cd doyo-img

# 复制并编辑配置文件
cp config.example.yaml config.yaml

# 启动
docker-compose up -d

# 浏览器访问 http://localhost:9090
```

### 从源码构建

**前提条件**: Go 1.22+, Node.js 18+

```bash
# 克隆项目
git clone https://github.com/doyo-img/doyo-img.git
cd doyo-img

# 方式一：使用 Makefile
make build           # Linux/macOS
make build-windows   # Windows

# 方式二：手动构建
cd web && npm install && npm run build && cd ..
go build -ldflags="-s -w" -o doyo-img .

# 运行
./doyo-img
```

## 本地开发

前后端分离开发，前端开发服务器自动代理 API 请求到后端。

**1. 配置开发环境**:

程序启动时在项目根目录和 `./config` 目录下查找 `config.yaml`。项目提供了两个配置模板：

| 文件 | 用途 |
|------|------|
| `config.example.yaml` | 生产环境模板（限流开启、info 日志） |
| `config.dev.yaml` | 开发环境模板（限流关闭、debug 日志、上传限制放宽） |

```bash
# 开发环境：复制开发配置
cp config.dev.yaml config.yaml

# 生产环境：复制生产配置
cp config.example.yaml config.yaml
```

如果不创建 `config.yaml`，程序使用内置默认值（等同于 `config.example.yaml`）。

**2. 启动后端**（终端 1）:
```bash
go run .
# 后端监听 http://localhost:9090
```

**3. 启动前端**（终端 2）:
```bash
cd web
npm install
npm run dev
# 前端开发服务器 http://localhost:5173，API 请求自动代理到 :9090
```

前端 Vite 开发服务器已配置代理规则（`vite.config.ts`），`/api` 和 `/i` 路径的请求会自动转发到后端。

## Docker 部署

### Dockerfile

多阶段构建，最终镜像基于 Alpine，体积小：

```dockerfile
# 构建前端 → 构建后端 → 仅复制二进制到运行时镜像
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
      - ./data:/app/data          # 图片持久化存储
      - ./config.yaml:/app/config.yaml:ro  # 配置文件挂载
    restart: unless-stopped
    environment:
      - DOYO_SERVER_PORT=9090
      # HTTPS 反向代理部署时设置:
      # - DOYO_SERVER_BASE_URL=https://img.example.com
```

**关键挂载说明**:
- `/app/data` — 图片文件存储目录，必须持久化
- `/app/config.yaml` — 配置文件，`:ro` 只读挂载

## 配置文件详解

复制 `config.example.yaml` 为 `config.yaml`，按需修改：

```yaml
server:
  host: "0.0.0.0"          # 监听地址
  port: 9090                # 监听端口
  base_url: ""              # 外部访问基础 URL（HTTPS 部署时必填）
                            # 示例: "https://img.example.com"

storage:
  type: "local"             # 存储类型: local | s3
  local:
    data_dir: "./data"      # 本地存储目录
  s3:                       # S3 兼容存储（type 设为 "s3" 时生效）
    # S3 兼容服务的 Endpoint（不含协议前缀，use_ssl 控制 http/https）
    # AWS S3:       s3.amazonaws.com
    # 阿里云 OSS:   oss-cn-hangzhou.aliyuncs.com
    # 腾讯云 COS:   cos.ap-guangzhou.myqcloud.com
    # MinIO:        localhost:9000
    endpoint: ""
    bucket: ""
    region: ""
    access_key: ""
    secret_key: ""
    use_ssl: true
    path_prefix: "images"

upload:
  max_file_size: 5242880    # 单文件最大字节数 (默认 5MB)
  max_batch_size: 10        # 单次最大上传数量
  allowed_formats:          # 允许的图片格式
    - "jpg"
    - "jpeg"
    - "png"
    - "gif"
    - "webp"
    - "svg"
  default_expire_hours: 0   # 默认过期时间（小时），0 = 永不过期
  max_expire_days: 0        # 最大过期天数上限，0 = 不限制

processing:
  compress_enabled: false   # 是否启用上传压缩
  compress_quality: 85      # JPEG 压缩质量 (1-100)
  strip_exif: true          # 自动剥离 JPEG 中的 EXIF 元数据
  thumbnail:
    enabled: true           # 是否生成缩略图
    max_width: 300          # 缩略图最大宽度
    max_height: 300         # 缩略图最大高度
  watermark:
    enabled: false          # 是否启用水印
    type: "text"            # text（文本水印）或 image（图片水印）
    text: "doyo-img"        # 文本水印内容
    font_path: ""           # 自定义字体路径（TTF/OTF），中文水印需指定 CJK 字体
    font_size: 24
    opacity: 0.3            # 不透明度 0.0 ~ 1.0
    position: "bottom-right" # top-left, top-right, bottom-left, bottom-right, center
    image_path: ""          # 图片水印 PNG 路径（type 为 image 时使用）
    padding: 20

id:
  length: 8                 # 图片 ID 长度
  alphabet: "0123456789abcdefghijklmnopqrstuvwxyz"  # ID 字符集

cors:
  allowed_origins:          # 允许的跨域来源
    - "*"                   # 生产环境建议设为具体域名

rate_limit:
  enabled: true             # 是否启用 IP 限流
  requests_per_minute: 30   # 每分钟请求上限
  burst: 10                 # 突发请求允许量

cleanup:
  enabled: false            # 是否启用过期清理
  retention_days: 30        # 全局保留天数，0 = 不按创建时间清理

log:
  level: "info"             # 日志级别: debug | info | warn | error
```

### 环境变量覆盖

所有配置项均可通过 `DOYO_` 前缀的环境变量覆盖，层级用 `_` 分隔：

| 配置项 | 环境变量 | 示例 |
|--------|---------|------|
| `server.port` | `DOYO_SERVER_PORT` | `8080` |
| `server.base_url` | `DOYO_SERVER_BASE_URL` | `https://img.example.com` |
| `storage.type` | `DOYO_STORAGE_TYPE` | `s3` |
| `storage.local.data_dir` | `DOYO_STORAGE_LOCAL_DATA_DIR` | `/data/images` |
| `upload.max_file_size` | `DOYO_UPLOAD_MAX_FILE_SIZE` | `10485760` |
| `rate_limit.enabled` | `DOYO_RATE_LIMIT_ENABLED` | `false` |
| `log.level` | `DOYO_LOG_LEVEL` | `debug` |

## API 接口文档

### 统一响应格式

```json
{
  "code": 0,
  "data": { ... },
  "message": "ok"
}
```

`code = 0` 表示成功，非零表示错误。

---

### POST /api/upload — 上传图片

**请求**: `multipart/form-data`，字段名 `file`（支持多文件）

可选字段 `expire_hours`：设置图片过期时间（小时），0 或不传表示永不过期。

```bash
# 单文件上传
curl -X POST -F "file=@photo.jpg" http://localhost:9090/api/upload

# 多文件上传
curl -X POST -F "file=@a.jpg" -F "file=@b.png" http://localhost:9090/api/upload

# 带过期时间上传（24 小时后过期）
curl -X POST -F "file=@photo.jpg" -F "expire_hours=24" http://localhost:9090/api/upload
```

**成功响应**:
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

> `expires_at` 仅在设置了过期时间时返回。

---

### GET /i/{id}.{ext} — 图片直链

直接返回图片二进制内容，可用于 `<img>` 标签和 Markdown 引用。

```
http://localhost:9090/i/a1b2c3d4.jpg          # 原图
http://localhost:9090/i/a1b2c3d4.jpg?t=thumb  # 缩略图
```

响应头包含 `Cache-Control: public, max-age=31536000, immutable`，浏览器和 CDN 可长期缓存。

SVG 图片额外返回 `Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'; img-src data:; sandbox` 头，防止 XSS 攻击。

已过期的图片返回 `410 Gone`。

---

### GET /api/image/{id} — 获取图片信息

```bash
curl http://localhost:9090/api/image/a1b2c3d4
```

返回图片元数据（ID、名称、格式、尺寸、大小、创建时间等）。

---

### DELETE /api/image/{id} — 删除图片

需在请求头传入上传时返回的删除令牌：

```bash
curl -X DELETE -H "X-Delete-Token: tok_xxxxxxxxxxxxxxxxxxxxxxxx" \
  http://localhost:9090/api/image/a1b2c3d4
```

---

### GET /api/recent — 最近上传列表

```bash
curl http://localhost:9090/api/recent?limit=20
```

参数 `limit` 范围 1-50，默认 20。

---

### GET /api/config — 获取公开配置

```bash
curl http://localhost:9090/api/config
```

返回前端所需的配置信息（文件大小限制、允许格式、水印开关、过期时间配置等），不含敏感信息。

## 生产环境部署

### HTTPS 部署（Nginx 反向代理）

doyo-img 本身监听 HTTP，生产环境通过 Nginx 终止 TLS：

**1. 配置 `config.yaml`**:
```yaml
server:
  base_url: "https://img.example.com"
cors:
  allowed_origins:
    - "https://img.example.com"
    - "https://your-blog.com"   # 需要嵌入图片的网站
```

**2. Nginx 配置**:
```nginx
server {
    listen 443 ssl http2;
    server_name img.example.com;

    ssl_certificate     /etc/ssl/certs/img.example.com.pem;
    ssl_certificate_key /etc/ssl/private/img.example.com.key;

    client_max_body_size 10m;  # 与 upload.max_file_size 匹配

    location / {
        proxy_pass http://127.0.0.1:9090;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Forwarded-Host $host;
    }

    # 图片直链的缓存策略（可选，服务端已设 Cache-Control）
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

**3. HTTPS 部署要点**:

| 配置项 | 说明 |
|--------|------|
| `server.base_url` | **必须**设为 `https://` 开头的完整域名，否则生成的图片链接为 HTTP |
| `cors.allowed_origins` | 设为具体的 HTTPS 域名列表，不要在生产环境使用 `*` |
| Nginx `client_max_body_size` | 需 ≥ `upload.max_file_size`，否则大文件上传被 Nginx 拦截 |
| Nginx `X-Forwarded-Proto` | 确保传递，后端据此判断协议生成正确链接 |

### Apache 反向代理

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

### 安全加固建议

- 生产环境将 `cors.allowed_origins` 设为具体域名白名单
- 启用限流 `rate_limit.enabled: true`，根据实际流量调整参数
- SVG 文件已内置 `Content-Security-Policy` 沙箱头防范 XSS（自动生效）
- JPEG 上传自动剥离 EXIF 隐私数据（`strip_exif: true` 默认开启）
- 定期备份 `data/` 目录（本地存储模式）
- 使用 `cleanup.enabled: true` 防止存储空间无限增长

## 常见问题

### 上传后图片链接是 HTTP 而非 HTTPS

确保以下任一条件满足：
1. 在 `config.yaml` 中设置 `server.base_url: "https://your-domain.com"`
2. 反向代理正确传递 `X-Forwarded-Proto: https` 头

### Nginx 返回 413 Request Entity Too Large

Nginx 的 `client_max_body_size` 默认 1MB，需在 Nginx 配置中增大：
```nginx
client_max_body_size 10m;
```

### Docker 容器中图片数据丢失

确保挂载了持久化卷：
```yaml
volumes:
  - ./data:/app/data
```

### 跨域请求被拦截

在 `config.yaml` 的 `cors.allowed_origins` 中添加请求来源域名：
```yaml
cors:
  allowed_origins:
    - "https://your-blog.com"
    - "https://another-site.com"
```

### 限流被触发 (429 Too Many Requests)

调整限流配置或对可信客户端禁用：
```yaml
rate_limit:
  requests_per_minute: 60
  burst: 20
```

## License

[MIT](LICENSE)
