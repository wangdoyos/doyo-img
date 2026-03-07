package handler

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/wangdoyos/doyo-img/internal/config"
	"github.com/wangdoyos/doyo-img/internal/service"
	"github.com/wangdoyos/doyo-img/internal/util"
)

// UploadHandler 图片上传 HTTP 处理器
type UploadHandler struct {
	svc *service.ImageService
	cfg *config.Config
}

func NewUploadHandler(svc *service.ImageService, cfg *config.Config) *UploadHandler {
	return &UploadHandler{svc: svc, cfg: cfg}
}

// Upload 处理 POST /api/upload —— 接收并处理批量图片上传
func (h *UploadHandler) Upload(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		util.Error(c, 400, 400, "invalid multipart form")
		return
	}

	files := form.File["file"]
	if len(files) == 0 {
		util.Error(c, 400, 400, "no files uploaded")
		return
	}

	if len(files) > h.cfg.Upload.MaxBatchSize {
		util.Error(c, 400, 400, fmt.Sprintf("too many files, maximum is %d", h.cfg.Upload.MaxBatchSize))
		return
	}

	baseURL := h.getBaseURL(c)

	// 解析可选的过期时间参数（小时）
	expireHours := 0
	if eh := form.Value["expire_hours"]; len(eh) > 0 {
		if parsed, err := strconv.Atoi(eh[0]); err == nil && parsed > 0 {
			expireHours = parsed
		}
	}

	var results []interface{}
	var errors []string

	for _, file := range files {
		// 校验文件大小
		if file.Size > h.cfg.Upload.MaxFileSize {
			errors = append(errors, fmt.Sprintf("%s: file too large (max %d bytes)", file.Filename, h.cfg.Upload.MaxFileSize))
			continue
		}

		f, err := file.Open()
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: failed to open file", file.Filename))
			continue
		}

		// 读入缓冲区以支持多次读取（seek）
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, f); err != nil {
			f.Close()
			errors = append(errors, fmt.Sprintf("%s: failed to read file", file.Filename))
			continue
		}
		f.Close()

		reader := bytes.NewReader(buf.Bytes())
		result, err := h.svc.Upload(c.Request.Context(), file.Filename, reader, file.Size, baseURL, expireHours)
		if err != nil {
			slog.Warn("上传失败", "filename", file.Filename, "error", err)
			errors = append(errors, fmt.Sprintf("%s: %s", file.Filename, err.Error()))
			continue
		}

		results = append(results, result)
	}

	data := gin.H{"images": results}
	if len(errors) > 0 {
		data["errors"] = errors
	}

	util.Success(c, data)
}

// getBaseURL 获取服务的外部访问基础 URL，优先使用配置值，否则从请求中推断
func (h *UploadHandler) getBaseURL(c *gin.Context) string {
	if h.cfg.Server.BaseURL != "" {
		return strings.TrimRight(h.cfg.Server.BaseURL, "/")
	}
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	if proto := c.GetHeader("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}
	host := c.Request.Host
	if fwdHost := c.GetHeader("X-Forwarded-Host"); fwdHost != "" {
		host = fwdHost
	}
	return fmt.Sprintf("%s://%s", scheme, host)
}
