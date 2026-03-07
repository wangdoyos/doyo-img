# 🚀 GitHub Releases 发布规划

## 版本策略

采用 **语义化版本控制 (Semantic Versioning)**：`MAJOR.MINOR.PATCH`

- **MAJOR**: 不兼容的 API 修改
- **MINOR**: 向下兼容的功能新增
- **PATCH**: 向下兼容的问题修复

---

## 📋 发布路线图

### v1.0.0 - 首个稳定版本 🎉
**计划日期**: 2026年3月

**核心功能**:
- ✅ 基础图床功能（上传、删除、访问）
- ✅ 多格式支持（JPG/PNG/GIF/WebP/SVG）
- ✅ 缩略图生成
- ✅ 图片水印（文字/图片）
- ✅ EXIF 自动剥离
- ✅ 链接过期与清理
- ✅ S3 兼容存储支持
- ✅ 中英文国际化
- ✅ 明暗主题切换
- ✅ Docker 部署支持

**发布资产**:
- `doyo-img-linux-amd64`
- `doyo-img-linux-arm64`
- `doyo-img-darwin-amd64`
- `doyo-img-darwin-arm64`
- `doyo-img-windows-amd64.exe`
- Source code (zip/tar.gz)

---

### v1.1.0 - 管理增强 📊
**计划日期**: 2026年4月

**新增功能**:
- 📊 管理后台仪表板
- 📈 上传统计面板（日/周/月）
- 🔍 图片搜索功能
- 📁 文件夹/标签分类
- 📱 PWA 支持
- 🔔 Webhook 通知

**优化**:
- 图片加载性能优化
- 缓存策略改进

---

### v1.2.0 - 用户系统 🔐
**计划日期**: 2026年5月

**新增功能**:
- 👤 用户注册/登录系统
- 🔐 JWT 认证
- 🏠 个人图床空间
- 🔒 私有/公开图片切换
- 👥 多用户配额管理
- 📋 用户操作日志

**安全增强**:
- API 密钥管理
- 更细粒度的权限控制

---

### v2.0.0 - 企业级 🚀
**计划日期**: 2026年Q3

**架构升级**:
- 🚀 分布式部署支持
- 🔄 集群模式（多实例负载均衡）
- 💾 Redis 缓存层
- 📨 消息队列异步处理
- 🗄️ 数据库支持（MySQL/PostgreSQL）

**高级功能**:
- 🎨 图片编辑（裁剪、旋转、滤镜）
- 📝 OCR 文字识别
- 🔍 以图搜图
- ☁️ CDN 集成
- 📊 高级分析报告

---

## 🏷️ 标签分类

| 标签 | 用途 |
|------|------|
| `latest` | 始终指向最新稳定版 |
| `stable` | 经过充分测试的稳定版本 |
| `beta` | 公开测试版本 |
| `alpha` | 内部测试版本 |

---

## 📝 Release Notes 模板

```markdown
## 🎉 What's New in vX.Y.Z

### ✨ Features
- 新增功能 A
- 新增功能 B

### 🐛 Bug Fixes
- 修复问题 C
- 修复问题 D

### 🚀 Improvements
- 性能优化 E
- 体验优化 F

### 🔒 Security
- 安全修复 G

### 📦 Dependencies
- 升级依赖 H

---

## 📥 Downloads

| Platform | Architecture | Download |
|----------|-------------|----------|
| Linux | amd64 | [doyo-img-linux-amd64](link) |
| Linux | arm64 | [doyo-img-linux-arm64](link) |
| macOS | amd64 | [doyo-img-darwin-amd64](link) |
| macOS | arm64 | [doyo-img-darwin-arm64](link) |
| Windows | amd64 | [doyo-img-windows-amd64.exe](link) |

### 🐳 Docker
```bash
docker pull ghcr.io/wangdoyos/doyo-img:vX.Y.Z
```

### 📋 Checksums
```
SHA256(doyo-img-linux-amd64)=xxxxxxxx...
SHA256(doyo-img-darwin-arm64)=xxxxxxxx...
```

---

## 🙏 Thanks
感谢所有贡献者！
```

---

## 🔧 自动化发布流程

### GitHub Actions 工作流

```yaml
# .github/workflows/release.yml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        goos: [linux, darwin, windows]
        goarch: [amd64, arm64]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: Build
        run: |
          cd web && npm ci && npm run build && cd ..
          GOOS=${{ matrix.goos }} GOARCH=${{ matrix.goarch }} go build -ldflags="-s -w" -o doyo-img-${{ matrix.goos }}-${{ matrix.goarch }}${{ matrix.goos == 'windows' && '.exe' || '' }} .
      - name: Upload
        uses: actions/upload-artifact@v4
        with:
          name: doyo-img-${{ matrix.goos }}-${{ matrix.goarch }}
          path: doyo-img-*

  release:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/download-artifact@v4
      - uses: softprops/action-gh-release@v1
        with:
          files: doyo-img-*
          generate_release_notes: true
```

---

## 🎯 发布检查清单

- [ ] 版本号更新（`main.go`, `web/package.json`）
- [ ] 更新 `CHANGELOG.md`
- [ ] 所有测试通过
- [ ] 文档已更新
- [ ] Docker 镜像构建成功
- [ ] 二进制文件多平台编译成功
- [ ] Release Notes 已准备
- [ ] Git Tag 已打 (`git tag -a v1.0.0 -m "Release v1.0.0"`)
- [ ] GitHub Release 已创建
- [ ] 公告已发布（可选）
