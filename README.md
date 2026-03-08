<div align="center">

# 🖼️ doyo-img

**轻量级、自托管图床服务**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://golang.org)
[![React](https://img.shields.io/badge/React-18-61DAFB?logo=react)](https://reactjs.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker)](https://www.docker.com)

🌐 **中文** | [English](README.en.md)

**无需登录，单文件部署，即刻获取图片分享链接**

</div>

---

## ✨ 功能特性

### 🚀 核心功能
- 📤 **多种上传方式** — 拖拽上传、Ctrl+V 粘贴、点击选择，支持批量上传
- 🔗 **即时分享** — 自动生成 URL / Markdown / HTML / BBCode 多种链接格式
- 🔓 **无需登录** — 开箱即用，任何人都可以上传和访问图片
- 🌐 **中英文切换** — 内置国际化支持，一键切换界面语言
- 🎨 **明暗主题** — 支持亮色 / 暗色 / 跟随系统三种主题模式

### 💾 存储与管理
- 📜 **上传历史** — 基于浏览器 localStorage 持久化，刷新不丢失
- 📊 **实时进度** — 上传过程实时显示进度百分比
- ℹ️ **图片信息** — 展示文件名、尺寸、分辨率、格式
- 🗑️ **令牌删除** — 每张图片生成唯一删除令牌，凭令牌安全删除

### 🔒 安全与隐私
- 🛡️ **IP 限流** — 基于令牌桶算法的 IP 级别请求限流
- 🌐 **CORS 跨域** — 可配置的跨域策略，支持第三方网站嵌入
- 🔍 **EXIF 自动剥离** — 自动移除 JPEG 图片中的 EXIF 元数据，保护隐私
- 🔒 **SVG 安全沙箱** — 为 SVG 响应注入 CSP 头，防范 Stored XSS

### 🎨 图片处理
- 🗜️ **图片压缩** — 可选的 JPEG/PNG 压缩处理
- 🖼️ **缩略图** — 自动生成缩略图，加速预览加载
- 💧 **图片水印** — 可选的文字/图片水印，支持自定义字体、位置、透明度
- ⏰ **链接过期** — 上传时可选设置过期时间，过期后自动返回 410 Gone
- 🧹 **过期清理** — 定时清理过期图片，支持全局保留天数和单图过期时间双重策略

### ☁️ 部署与存储
- ☁️ **S3 存储** — 支持 AWS S3、阿里云 OSS、腾讯云 COS、MinIO 等 S3 兼容对象存储
- ⚙️ **环境变量覆盖** — YAML 配置 + `DOYO_` 前缀环境变量，灵活部署
- 📦 **单文件部署** — Go embed 嵌入前端，编译产出单个二进制
- 🐳 **Docker 支持** — 提供 Dockerfile 和 docker-compose 配置

---

## 🏗️ 技术架构

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

| 层级 | 技术栈 |
|------|--------|
| **后端** | Go + Gin 框架，分层架构 (handler → service → storage) |
| **前端** | React 18 + TypeScript + Vite + Tailwind CSS + Zustand |
| **构建** | 前端编译为静态文件，通过 `go:embed` 嵌入 Go 二进制 |

---

## 📁 目录结构

```
doyo-img/
├── 📄 main.go                    # 程序入口：初始化配置、存储、路由、启动服务
├── 📄 embed.go                   # go:embed 嵌入前端静态资源
├── 📄 config.example.yaml        # 配置文件示例
├── 🐳 Dockerfile                 # Docker 多阶段构建
├── 📄 docker-compose.yml         # Docker Compose 编排
├── 📄 Makefile                   # 常用构建命令
├── 📁 internal/
│   ├── 📁 config/                # 配置定义与加载 (Viper)
│   ├── 📁 model/                 # 数据模型
│   ├── 📁 storage/               # 存储接口与实现
│   ├── 📁 service/               # 业务逻辑层
│   ├── 📁 handler/               # HTTP 请求处理
│   ├── 📁 processor/             # 图片处理（压缩/水印/EXIF）
│   ├── 📁 middleware/            # 中间件（CORS/限流）
│   ├── 📁 router/                # 路由注册
│   ├── 📁 cleanup/               # 过期图片定时清理
│   └── 📁 util/                  # 工具函数
└── 📁 web/                       # React 前端
    ├── 📁 src/
    │   ├── 📁 components/        # 组件
    │   ├── 📁 hooks/             # 自定义 Hooks
    │   ├── 📁 i18n/              # 国际化
    │   ├── 📁 store/             # 状态管理
    │   └── 📁 utils/             # 工具函数
    └── ⚙️ vite.config.ts         # Vite 配置
```

---

## 🚀 快速开始

### 📋 系统要求

| 部署方式 | 最低要求 | 推荐配置 |
|---------|---------|---------|
| **直接运行** | Go 1.22+ | Go 1.22+, 64MB 内存 |
| **Docker** | Docker 20.10+ | Docker 24.0+, 256MB 内存 |
| **源码构建** | Go 1.22+, Node.js 18+ | Go 1.22+, Node.js 20+, 512MB 内存 |

**磁盘空间**: 根据图片存储需求决定，建议预留 1GB+ 用于系统和日志。

---

### ⚡ 方式一：直接运行（推荐快速体验）

从 [GitHub Releases](https://github.com/wangdoyos/doyo-img/releases) 下载适合您系统的预编译二进制文件：

```bash
# Linux/macOS
wget https://github.com/wangdoyos/doyo-img/releases/latest/download/doyo-img-linux-amd64 -O doyo-img
chmod +x doyo-img
./doyo-img

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/wangdoyos/doyo-img/releases/latest/download/doyo-img-windows-amd64.exe" -OutFile "doyo-img.exe"
.\doyo-img.exe
```

浏览器访问 http://localhost:9090

---

### 🐳 方式二：Docker Hub 部署（推荐生产环境）

使用官方 Docker Hub 镜像快速部署，无需克隆源码：

#### 快速启动

```bash
# 创建数据目录
mkdir -p ./data

# 下载默认配置文件
wget https://raw.githubusercontent.com/wangdoyos/doyo-img/main/config.example.yaml -O config.yaml

# 运行容器
docker run -d \
  --name doyo-img \
  -p 9090:9090 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/config.yaml:/app/config.yaml:ro \
  --restart unless-stopped \
  wangdoyo/doyo-img:latest
```

#### Docker Compose 部署（推荐）

创建 `docker-compose.yml`：

```yaml
name: doyo
services:
  doyo-img:
    image: wangdoyo/doyo-img:latest
    container_name: doyo-img
    ports:
      - "9090:9090"
    volumes:
      - ./data:/app/data          # 图片持久化存储
      - ./config.yaml:/app/config.yaml:ro  # 配置文件（只读）
    restart: unless-stopped
    environment:
      - DOYO_SERVER_PORT=9090
      # HTTPS 反向代理部署时，设置外部可访问的基础 URL
      # - DOYO_SERVER_BASE_URL=https://img.example.com
```

启动服务：

```bash
docker-compose up -d
```

#### 镜像标签说明

| 标签 | 说明 | 适用场景 |
|------|------|---------|
| `latest` | 最新稳定版本 | 开发测试 |
| `v1.x.x` | 特定版本 | 生产环境（推荐） |
| `v1` | 最新 v1 版本 | 平衡稳定和更新 |

查看所有可用标签：[Docker Hub Tags](https://hub.docker.com/r/wangdoyos/doyo-img/tags)

#### 更新镜像

```bash
# 拉取最新镜像
docker-compose pull

# 重启容器
docker-compose up -d
```

---

### 🏗️ 方式三：本地构建 Docker 镜像

如需自定义镜像或修改源码后构建：

```bash
# 克隆项目
git clone https://github.com/wangdoyos/doyo-img.git
cd doyo-img

# 复制配置文件
cp config.example.yaml config.yaml

# 构建并启动
docker-compose up -d --build
```

---

### 🔨 方式四：从源码构建

适合需要二次开发或自定义构建的用户。

**前提条件**: Go 1.22+, Node.js 18+

```bash
# 克隆项目
git clone https://github.com/wangdoyos/doyo-img.git
cd doyo-img

# 方式一：使用 Makefile（推荐）
make build           # Linux/macOS
make build-windows   # Windows

# 方式二：手动构建
cd web && npm install && npm run build && cd ..
go build -ldflags="-s -w" -o doyo-img .

# 运行
./doyo-img
```

浏览器访问 http://localhost:9090

---

## 💻 本地开发

前后端分离开发，前端开发服务器自动代理 API 请求到后端。

### 1️⃣ 配置开发环境

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

### 2️⃣ 启动后端（终端 1）
```bash
go run .
# 后端监听 http://localhost:9090
```

### 3️⃣ 启动前端（终端 2）
```bash
cd web
npm install
npm run dev
# 前端开发服务器 http://localhost:5173，API 请求自动代理到 :9090
```

---

## 🐳 Docker 部署

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
- 📁 `/app/data` — 图片文件存储目录，必须持久化
- 📄 `/app/config.yaml` — 配置文件，`:ro` 只读挂载

---

## ⚙️ 配置文件详解

配置文件采用 YAML 格式，程序启动时按以下顺序查找：
1. `./config.yaml`（当前目录）
2. `./config/config.yaml`（config 子目录）
3. 使用内置默认值

### 快速配置

```bash
# 复制默认配置
cp config.example.yaml config.yaml

# 编辑配置
vim config.yaml
```

### 配置项详细说明

#### 1. 服务器配置 (`server`)

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `host` | string | `"0.0.0.0"` | 监听地址，`0.0.0.0` 表示监听所有接口 |
| `port` | int | `9090` | 监听端口，范围 1-65535 |
| `base_url` | string | `""` | 外部访问基础 URL，**HTTPS 部署时必填** |

**重要**: `base_url` 用于生成图片链接，生产环境必须设置为完整的 HTTPS 地址，如 `https://img.example.com`。若使用反向代理，确保配置正确的协议头。

#### 2. 存储配置 (`storage`)

| 配置项 | 类型 | 默认值 | 可选值 | 说明 |
|--------|------|--------|--------|------|
| `type` | string | `"local"` | `local`, `s3` | 存储后端类型 |

**本地存储** (`storage.local`):

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `data_dir` | string | `"./data"` | 图片存储目录，需确保有写入权限 |

**S3 存储** (`storage.s3`):

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `endpoint` | string | `""` | S3 Endpoint（不含协议前缀） |
| `bucket` | string | `""` | 存储桶名称 |
| `region` | string | `""` | 区域代码，如 `us-east-1` |
| `access_key` | string | `""` | Access Key ID |
| `secret_key` | string | `""` | Secret Access Key |
| `use_ssl` | bool | `true` | 是否使用 HTTPS |
| `path_prefix` | string | `"images"` | 对象键前缀 |

**常见 S3 服务商 Endpoint 示例**:
- AWS S3: `s3.amazonaws.com`
- 阿里云 OSS: `oss-cn-hangzhou.aliyuncs.com`
- 腾讯云 COS: `cos.ap-guangzhou.myqcloud.com`
- MinIO: `localhost:9000`

#### 3. 上传配置 (`upload`)

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `max_file_size` | int | `5242880` | 单文件最大字节数（默认 5MB） |
| `max_batch_size` | int | `10` | 单次最大上传数量 |
| `allowed_formats` | array | `["jpg", "jpeg", "png", "gif", "webp", "svg"]` | 允许的图片格式 |
| `default_expire_hours` | int | `0` | 默认过期时间（小时），0 = 永不过期 |
| `max_expire_days` | int | `0` | 最大过期天数上限，0 = 不限制 |

#### 4. 图片处理配置 (`processing`)

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `compress_enabled` | bool | `false` | 是否启用上传压缩 |
| `compress_quality` | int | `85` | JPEG 压缩质量（1-100，越高质量越好） |
| `strip_exif` | bool | `true` | 自动剥离 JPEG 中的 EXIF 元数据（GPS、设备信息） |

**缩略图配置** (`processing.thumbnail`):

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `enabled` | bool | `true` | 是否生成缩略图 |
| `max_width` | int | `300` | 缩略图最大宽度（像素） |
| `max_height` | int | `300` | 缩略图最大高度（像素） |

**水印配置** (`processing.watermark`):

| 配置项 | 类型 | 默认值 | 可选值 | 说明 |
|--------|------|--------|--------|------|
| `enabled` | bool | `false` | - | 是否启用水印 |
| `type` | string | `"text"` | `text`, `image` | 水印类型 |
| `text` | string | `"doyo-img"` | - | 文本水印内容 |
| `font_path` | string | `""` | - | 自定义字体路径（TTF/OTF），中文水印需指定 CJK 字体 |
| `font_size` | int | `24` | - | 字体大小（像素） |
| `opacity` | float | `0.3` | `0.0` - `1.0` | 不透明度 |
| `position` | string | `"bottom-right"` | `top-left`, `top-right`, `bottom-left`, `bottom-right`, `center` | 水印位置 |
| `image_path` | string | `""` | - | 图片水印 PNG 路径（`type: image` 时使用） |
| `padding` | int | `20` | - | 水印边距（像素） |

#### 5. ID 生成配置 (`id`)

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `length` | int | `8` | 图片 ID 长度，建议 6-12 |
| `alphabet` | string | `"0123456789abcdefghijklmnopqrstuvwxyz"` | ID 字符集 |

#### 6. 跨域配置 (`cors`)

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `allowed_origins` | array | `["*"]` | 允许的跨域来源，生产环境建议设为具体域名 |

#### 7. 限流配置 (`rate_limit`)

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `enabled` | bool | `true` | 是否启用 IP 限流 |
| `requests_per_minute` | int | `30` | 每分钟请求数限制 |
| `burst` | int | `10` | 突发请求容量 |

#### 8. 清理配置 (`cleanup`)

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `enabled` | bool | `false` | 是否启用过期清理 |
| `retention_days` | int | `30` | 全局保留天数，0 = 不按创建时间清理 |

#### 9. 日志配置 (`log`)

| 配置项 | 类型 | 默认值 | 可选值 | 说明 |
|--------|------|--------|--------|------|
| `level` | string | `"info"` | `debug`, `info`, `warn`, `error` | 日志级别 |

### 完整配置示例

```yaml
server:
  host: "0.0.0.0"
  port: 9090
  base_url: ""

storage:
  type: "local"
  local:
    data_dir: "./data"
  s3:
    endpoint: ""
    bucket: ""
    region: ""
    access_key: ""
    secret_key: ""
    use_ssl: true
    path_prefix: "images"

upload:
  max_file_size: 5242880
  max_batch_size: 10
  allowed_formats:
    - "jpg"
    - "jpeg"
    - "png"
    - "gif"
    - "webp"
    - "svg"
  default_expire_hours: 0
  max_expire_days: 0

processing:
  compress_enabled: false
  compress_quality: 85
  strip_exif: true
  thumbnail:
    enabled: true
    max_width: 300
    max_height: 300
  watermark:
    enabled: false
    type: "text"
    text: "doyo-img"
    font_path: ""
    font_size: 24
    opacity: 0.3
    position: "bottom-right"
    image_path: ""
    padding: 20

id:
  length: 8
  alphabet: "0123456789abcdefghijklmnopqrstuvwxyz"

cors:
  allowed_origins:
    - "*"

rate_limit:
  enabled: true
  requests_per_minute: 30
  burst: 10

cleanup:
  enabled: false
  retention_days: 30

log:
  level: "info"
```

### 🔧 环境变量覆盖

所有配置项均可通过 `DOYO_` 前缀的环境变量覆盖：

| 配置项 | 环境变量 | 示例 |
|--------|---------|------|
| `server.port` | `DOYO_SERVER_PORT` | `8080` |
| `server.base_url` | `DOYO_SERVER_BASE_URL` | `https://img.example.com` |
| `storage.type` | `DOYO_STORAGE_TYPE` | `s3` |
| `upload.max_file_size` | `DOYO_UPLOAD_MAX_FILE_SIZE` | `10485760` |

---

## 📡 API 接口文档

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

### 📤 POST /api/upload — 上传图片

**请求**: `multipart/form-data`，字段名 `file`（支持多文件）

可选字段 `expire_hours`：设置图片过期时间（小时）

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

---

### 🖼️ GET /i/{id}.{ext} — 图片直链

```
http://localhost:9090/i/a1b2c3d4.jpg          # 原图
http://localhost:9090/i/a1b2c3d4.jpg?t=thumb  # 缩略图
```

- 响应头包含 `Cache-Control: public, max-age=31536000, immutable`
- SVG 图片额外返回 CSP 头防止 XSS
- 已过期的图片返回 `410 Gone`

---

### ℹ️ GET /api/image/{id} — 获取图片信息

```bash
curl http://localhost:9090/api/image/a1b2c3d4
```

---

### 🗑️ DELETE /api/image/{id} — 删除图片

```bash
curl -X DELETE -H "X-Delete-Token: tok_xxxxxxxxxxxxxxxxxxxxxxxx" \
  http://localhost:9090/api/image/a1b2c3d4
```

---

### 📜 GET /api/recent — 最近上传列表

```bash
curl http://localhost:9090/api/recent?limit=20
```

---

### ⚙️ GET /api/config — 获取公开配置

```bash
curl http://localhost:9090/api/config
```

---

## 🚀 生产环境部署

### 🔒 HTTPS 部署（Nginx 反向代理）

**1️⃣ 配置 `config.yaml`**:
```yaml
server:
  base_url: "https://img.example.com"
cors:
  allowed_origins:
    - "https://img.example.com"
    - "https://your-blog.com"
```

**2️⃣ Nginx 配置**:
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

**3️⃣ HTTPS 部署要点**:

| 配置项 | 说明 |
|--------|------|
| `server.base_url` | **必须**设为 `https://` 开头的完整域名 |
| `cors.allowed_origins` | 设为具体的 HTTPS 域名列表，不要生产环境使用 `*` |
| Nginx `client_max_body_size` | 需 ≥ `upload.max_file_size` |

---

## 🏭 生产环境最佳实践

### 1. 安全配置清单

生产环境部署前，请确保完成以下安全配置：

```yaml
# config.yaml 生产环境推荐配置
server:
  base_url: "https://img.example.com"  # 必须设置为 HTTPS 域名

cors:
  allowed_origins:
    - "https://img.example.com"        # 只允许特定域名
    - "https://your-blog.com"

rate_limit:
  enabled: true                         # 启用限流防止滥用
  requests_per_minute: 30
  burst: 10

upload:
  max_file_size: 10485760               # 根据需求调整（10MB）
  max_batch_size: 5                     # 限制批量上传数量

cleanup:
  enabled: true                         # 启用过期清理
  retention_days: 90                    # 设置合理的保留期限

log:
  level: "info"                         # 生产环境使用 info 级别
```

### 2. 性能优化建议

| 优化项 | 建议配置 | 说明 |
|--------|---------|------|
| **图片压缩** | `compress_enabled: true` | 减少存储和带宽 |
| **缩略图** | `thumbnail.enabled: true` | 加速预览加载 |
| **CDN 加速** | 配合 CloudFlare 等 CDN | 提升全球访问速度 |
| **存储选择** | 大容量场景使用 S3 | 支持分布式存储 |

### 3. 监控与日志

**查看容器日志**:
```bash
docker logs -f doyo-img
```

**日志轮转配置**（Docker）:
```yaml
services:
  doyo-img:
    image: wangdoyo/doyo-img:latest
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

### 4. 备份策略

**本地存储备份**:
```bash
# 定期备份数据目录
tar -czf backup-$(date +%Y%m%d).tar.gz ./data

# 使用 rsync 同步到远程
rsync -avz ./data backup-server:/backups/doyo-img/
```

**S3 存储**: 启用存储桶的版本控制和跨区域复制。

### 5. 高可用部署

对于高流量场景，建议：
- 使用负载均衡（Nginx/HAProxy）部署多个实例
- 使用共享存储（S3/NFS）确保多实例数据一致
- 配置健康检查接口：`GET /api/config`

---

## 🛡️ 安全加固建议

- ✅ 生产环境将 `cors.allowed_origins` 设为具体域名白名单
- ✅ 启用限流 `rate_limit.enabled: true`
- ✅ SVG 文件已内置 CSP 沙箱头防范 XSS（自动生效）
- ✅ JPEG 上传自动剥离 EXIF 隐私数据
- ✅ 定期备份 `data/` 目录
- ✅ 使用 `cleanup.enabled: true` 防止存储空间无限增长

---

## ❓ 常见问题与故障排除

### 部署问题

#### Q: 上传后图片链接是 HTTP 而非 HTTPS
**A:** 确保以下任一条件满足：
1. 在 `config.yaml` 中设置 `server.base_url: "https://your-domain.com"`
2. 反向代理正确传递 `X-Forwarded-Proto: https` 头

#### Q: Nginx 返回 413 Request Entity Too Large
**A:** 在 Nginx 配置中增大：
```nginx
client_max_body_size 10m;
```

#### Q: Docker 容器中图片数据丢失
**A:** 确保挂载了持久化卷：
```yaml
volumes:
  - ./data:/app/data
```

#### Q: 端口 9090 被占用
**A:** 修改 `config.yaml` 中的端口或使用环境变量覆盖：
```bash
export DOYO_SERVER_PORT=8080
./doyo-img
```

#### Q: 权限 denied 错误
**A:** 确保数据目录有写入权限：
```bash
chmod 755 ./data
chown -R $(whoami):$(whoami) ./data
```

### 使用问题

#### Q: 跨域请求被拦截
**A:** 在 `config.yaml` 的 `cors.allowed_origins` 中添加请求来源域名：
```yaml
cors:
  allowed_origins:
    - "https://your-website.com"
```

#### Q: 中文水印显示为方块
**A:** 需要指定支持中文的字体文件：
```yaml
processing:
  watermark:
    font_path: "/path/to/NotoSansCJK-Regular.ttc"
```

#### Q: 上传大文件失败
**A:** 检查以下配置：
1. `upload.max_file_size` 是否足够大
2. Nginx/反向代理的 `client_max_body_size`
3. 浏览器开发者工具查看具体错误

#### Q: 如何迁移数据到另一台服务器
**A:** 
```bash
# 本地存储：复制 data 目录
rsync -avz ./data new-server:/path/to/data

# S3 存储：无需迁移，直接在新服务器配置相同 S3 参数
```

#### Q: 图片链接失效或 410 错误
**A:** 检查图片是否已过期：
- 查看 `expires_at` 字段
- 检查 `cleanup.enabled` 是否误删了未过期图片
- 确认服务器时间和时区设置正确

---

## 📋 发布计划

| 版本 | 计划日期 | 主要特性 |
|------|---------|---------|
| v1.0.0 | 2026-03 | 🎉 首个稳定版本 |
| v1.1.0 | 2026-04 | 📊 管理后台、统计面板 |
| v1.2.0 | 2026-05 | 🔐 用户系统、私有图床 |
| v2.0.0 | 2026-Q3 | 🚀 分布式部署、集群支持 |

---

## 📄 License

[MIT](LICENSE) © 2026 doyo-img
