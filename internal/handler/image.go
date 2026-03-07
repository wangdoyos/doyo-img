package handler

import (
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/wangdoyos/doyo-img/internal/config"
	"github.com/wangdoyos/doyo-img/internal/service"
	"github.com/wangdoyos/doyo-img/internal/util"
)

// ImageHandler 图片访问相关的 HTTP 处理器
type ImageHandler struct {
	svc *service.ImageService
	cfg *config.Config
}

func NewImageHandler(svc *service.ImageService, cfg *config.Config) *ImageHandler {
	return &ImageHandler{svc: svc, cfg: cfg}
}

// ServeImage 处理 GET /i/*id —— 图片直链访问，返回图片二进制内容
func (h *ImageHandler) ServeImage(c *gin.Context) {
	idParam := c.Param("id")
	// Gin 通配符参数包含前导斜杠："/a1b2c3d4.jpg" -> "a1b2c3d4"
	idParam = strings.TrimPrefix(idParam, "/")

	// 去除扩展名："a1b2c3d4.jpg" -> "a1b2c3d4"
	id := idParam
	if dotIdx := strings.LastIndex(idParam, "."); dotIdx > 0 {
		id = idParam[:dotIdx]
	}

	// 检查是否请求缩略图
	if c.Query("t") == "thumb" {
		meta, err := h.svc.GetMeta(c.Request.Context(), id)
		if err != nil {
			util.Error(c, 404, 404, "image not found")
			return
		}

		// 检查图片是否已过期
		if meta.ExpiresAt != nil && time.Now().After(*meta.ExpiresAt) {
			util.Error(c, 410, 410, "image expired")
			return
		}

		reader, err := h.svc.GetThumbnail(c.Request.Context(), id)
		if err != nil {
			// 缩略图不存在时回退到原图
			h.serveOriginal(c, id)
			return
		}
		defer reader.Close()

		c.Header("Cache-Control", "public, max-age=31536000, immutable")
		c.Header("Content-Type", util.FormatToMime(meta.Format))
		io.Copy(c.Writer, reader)
		return
	}

	h.serveOriginal(c, id)
}

// serveOriginal 返回原始图片
func (h *ImageHandler) serveOriginal(c *gin.Context, id string) {
	reader, meta, err := h.svc.GetImage(c.Request.Context(), id)
	if err != nil {
		util.Error(c, 404, 404, "image not found")
		return
	}
	defer reader.Close()

	// 检查图片是否已过期
	if meta.ExpiresAt != nil && time.Now().After(*meta.ExpiresAt) {
		util.Error(c, 410, 410, "image expired")
		return
	}

	c.Header("Cache-Control", "public, max-age=31536000, immutable")
	c.Header("Content-Type", meta.MimeType)
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("ETag", `"`+meta.ID+`"`)
	c.Header("Content-Disposition", "inline")

	// SVG 安全沙箱：阻止嵌入的 JavaScript 执行，防止 Stored XSS
	if meta.MimeType == "image/svg+xml" {
		c.Header("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; img-src data:; sandbox")
	}

	io.Copy(c.Writer, reader)
}

// GetImageInfo 处理 GET /api/image/:id —— 返回图片元数据信息
func (h *ImageHandler) GetImageInfo(c *gin.Context) {
	id := c.Param("id")
	meta, err := h.svc.GetMeta(c.Request.Context(), id)
	if err != nil {
		util.Error(c, 404, 404, "image not found")
		return
	}

	// 检查图片是否已过期
	if meta.ExpiresAt != nil && time.Now().After(*meta.ExpiresAt) {
		util.Error(c, 410, 410, "image expired")
		return
	}

	util.Success(c, meta)
}

// DeleteImage 处理 DELETE /api/image/:id —— 通过 token 验证后删除图片
func (h *ImageHandler) DeleteImage(c *gin.Context) {
	id := c.Param("id")
	token := c.GetHeader("X-Delete-Token")
	if token == "" {
		util.Error(c, 401, 401, "delete token required")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id, token); err != nil {
		if strings.Contains(err.Error(), "invalid delete token") {
			util.Error(c, 403, 403, "invalid delete token")
			return
		}
		util.Error(c, 404, 404, "image not found")
		return
	}

	util.Success(c, gin.H{"message": "deleted"})
}

// ListRecent 处理 GET /api/recent —— 返回最近上传的图片列表
func (h *ImageHandler) ListRecent(c *gin.Context) {
	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	metas, err := h.svc.ListRecent(c.Request.Context(), limit)
	if err != nil {
		util.Error(c, 500, 500, "failed to list images")
		return
	}

	util.Success(c, gin.H{"images": metas})
}
