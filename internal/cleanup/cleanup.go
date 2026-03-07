package cleanup

import (
	"context"
	"log/slog"
	"time"

	"github.com/wangdoyos/doyo-img/internal/config"
	"github.com/wangdoyos/doyo-img/internal/storage"
)

// Start 启动定时清理协程，定期删除过期图片
func Start(ctx context.Context, store storage.Storage, cfg *config.CleanupConfig) {
	if !cfg.Enabled {
		return
	}

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	slog.Info("清理调度器已启动", "retention_days", cfg.RetentionDays)

	for {
		select {
		case <-ctx.Done():
			slog.Info("清理调度器已停止")
			return
		case <-ticker.C:
			run(ctx, store, cfg.RetentionDays)
		}
	}
}

// run 执行一次清理任务，删除过期图片和超过保留天数的图片
func run(ctx context.Context, store storage.Storage, retentionDays int) {
	slog.Info("正在执行清理任务")

	now := time.Now().UTC()
	cutoff := now.AddDate(0, 0, -retentionDays)

	metas, err := store.List(ctx, 0) // 0 表示获取全部
	if err != nil {
		slog.Error("清理: 获取图片列表失败", "error", err)
		return
	}

	deleted := 0
	for _, meta := range metas {
		shouldDelete := false
		reason := ""

		// 检查单图过期时间
		if meta.ExpiresAt != nil && now.After(*meta.ExpiresAt) {
			shouldDelete = true
			reason = "单图过期"
		}

		// 检查全局保留期限
		if !shouldDelete && retentionDays > 0 && meta.CreatedAt.Before(cutoff) {
			shouldDelete = true
			reason = "超过保留期限"
		}

		if shouldDelete {
			if err := store.Delete(ctx, meta.ID); err != nil {
				slog.Error("清理: 删除图片失败", "id", meta.ID, "reason", reason, "error", err)
				continue
			}
			deleted++
			slog.Debug("清理: 已删除图片", "id", meta.ID, "reason", reason)
		}
	}

	if deleted > 0 {
		slog.Info("清理完成", "deleted", deleted)
	}
}
