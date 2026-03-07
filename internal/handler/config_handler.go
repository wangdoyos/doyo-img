package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/wangdoyos/doyo-img/internal/config"
	"github.com/wangdoyos/doyo-img/internal/util"
)

// ConfigHandler 公开配置信息 HTTP 处理器
type ConfigHandler struct {
	cfg *config.Config
}

func NewConfigHandler(cfg *config.Config) *ConfigHandler {
	return &ConfigHandler{cfg: cfg}
}

// GetPublicConfig 处理 GET /api/config —— 返回前端所需的公开配置（不含敏感信息）
func (h *ConfigHandler) GetPublicConfig(c *gin.Context) {
	util.Success(c, gin.H{
		"max_file_size":        h.cfg.Upload.MaxFileSize,
		"max_batch_size":       h.cfg.Upload.MaxBatchSize,
		"allowed_formats":      h.cfg.Upload.AllowedFormats,
		"compress_enabled":     h.cfg.Processing.CompressEnabled,
		"base_url":             h.cfg.Server.BaseURL,
		"watermark_enabled":    h.cfg.Processing.Watermark.Enabled,
		"default_expire_hours": h.cfg.Upload.DefaultExpireHours,
		"max_expire_days":      h.cfg.Upload.MaxExpireDays,
	})
}
