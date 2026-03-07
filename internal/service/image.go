package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log/slog"
	"time"

	"github.com/wangdoyos/doyo-img/internal/config"
	"github.com/wangdoyos/doyo-img/internal/model"
	"github.com/wangdoyos/doyo-img/internal/processor"
	"github.com/wangdoyos/doyo-img/internal/storage"
	"github.com/wangdoyos/doyo-img/internal/util"

	_ "golang.org/x/image/webp" // WebP 解码支持
)

// ImageService 图片业务逻辑层，串联校验、处理、存储全流程
type ImageService struct {
	storage storage.Storage
	cfg     *config.Config
}

func NewImageService(s storage.Storage, cfg *config.Config) *ImageService {
	return &ImageService{storage: s, cfg: cfg}
}

// Upload 处理单张图片上传：校验 -> 检测类型 -> 提取信息 -> 可选压缩 -> 存储 -> 生成缩略图
func (s *ImageService) Upload(ctx context.Context, filename string, data io.ReadSeeker, size int64, baseURL string, expireHours int) (*model.UploadResult, error) {
	// 校验文件大小
	if err := util.ValidateFileSize(size, s.cfg.Upload.MaxFileSize); err != nil {
		return nil, fmt.Errorf("文件过大: %w", err)
	}

	// 通过 magic bytes 检测真实 MIME 类型
	mime, err := util.DetectMimeType(data)
	if err != nil {
		return nil, fmt.Errorf("检测文件类型失败: %w", err)
	}

	// 校验格式是否在允许列表中
	if !util.ValidateFormat(mime, s.cfg.Upload.AllowedFormats) {
		return nil, fmt.Errorf("不支持的文件格式: %s", mime)
	}

	format := util.MimeToFormat(mime)
	isSVG := format == "svg"

	// 提取图片尺寸信息
	info, err := processor.GetInfo(data, isSVG)
	if err != nil {
		slog.Warn("获取图片信息失败", "error", err)
		info = &processor.ImageInfo{Width: 0, Height: 0, Format: format}
	}

	// 生成唯一 ID
	id, err := util.GenerateID(s.cfg.ID.Length, s.cfg.ID.Alphabet)
	if err != nil {
		return nil, fmt.Errorf("生成ID失败: %w", err)
	}

	// 生成删除令牌
	deleteToken, err := util.GenerateDeleteToken()
	if err != nil {
		return nil, fmt.Errorf("生成删除令牌失败: %w", err)
	}

	// 可选 EXIF 剥离（仅当压缩未启用时需要独立剥离，压缩本身的 decode→encode 已天然丢弃 EXIF）
	if s.cfg.Processing.StripExif && !isSVG && (format == "jpg" || format == "jpeg") && !s.cfg.Processing.CompressEnabled {
		stripped, err := processor.StripExif(data, format)
		if err != nil {
			slog.Warn("EXIF 剥离失败，使用原图", "error", err)
			if _, err := data.Seek(0, io.SeekStart); err != nil {
				return nil, err
			}
		} else {
			// StripExif 返回 *bytes.Buffer，后续不再需要 seek
			// 将结果包装为 bytes.Reader 以支持后续可能的 seek 操作
			strippedBytes, _ := io.ReadAll(stripped)
			data = bytes.NewReader(strippedBytes)
		}
	}

	// 可选压缩处理
	var reader io.Reader
	if s.cfg.Processing.CompressEnabled && !isSVG {
		compressed, err := processor.Compress(data, format, s.cfg.Processing.CompressQuality)
		if err != nil {
			slog.Warn("压缩失败，使用原图", "error", err)
			if _, err := data.Seek(0, io.SeekStart); err != nil {
				return nil, err
			}
			reader = data
		} else {
			reader = compressed
		}
	} else {
		if _, err := data.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}
		reader = data
	}

	// 可选水印处理（在压缩之后、存储之前，不对 SVG 和 GIF 处理）
	if s.cfg.Processing.Watermark.Enabled && !isSVG && format != "gif" {
		watermarked, err := s.applyWatermark(reader, format)
		if err != nil {
			slog.Warn("水印处理失败，使用原图", "error", err)
		} else {
			reader = watermarked
		}
	}

	// 构建元数据
	now := time.Now().UTC()
	meta := &model.ImageMeta{
		ID:           id,
		OriginalName: filename,
		Format:       format,
		MimeType:     mime,
		Size:         size,
		Width:        info.Width,
		Height:       info.Height,
		DeleteToken:  deleteToken,
		CreatedAt:    now,
	}

	// 计算过期时间
	effectiveExpireHours := expireHours
	if effectiveExpireHours == 0 && s.cfg.Upload.DefaultExpireHours > 0 {
		effectiveExpireHours = s.cfg.Upload.DefaultExpireHours
	}
	if effectiveExpireHours > 0 {
		// 如果配置了最大过期天数，则校验不超限
		if s.cfg.Upload.MaxExpireDays > 0 {
			maxHours := s.cfg.Upload.MaxExpireDays * 24
			if effectiveExpireHours > maxHours {
				effectiveExpireHours = maxHours
			}
		}
		expiresAt := now.Add(time.Duration(effectiveExpireHours) * time.Hour)
		meta.ExpiresAt = &expiresAt
	}

	// 保存到存储后端
	// 使用 TeeReader 同时缓冲数据，以便后续生成缩略图
	var buf bytes.Buffer
	teeReader := io.TeeReader(reader, &buf)

	if err := s.storage.Save(ctx, id, teeReader, meta); err != nil {
		return nil, fmt.Errorf("保存图片失败: %w", err)
	}

	// 生成缩略图（如果启用且非 SVG）
	thumbnailURL := ""
	if s.cfg.Processing.Thumbnail.Enabled && !isSVG {
		thumbReader := bytes.NewReader(buf.Bytes())
		thumb, err := processor.GenerateThumbnail(
			thumbReader,
			format,
			s.cfg.Processing.Thumbnail.MaxWidth,
			s.cfg.Processing.Thumbnail.MaxHeight,
		)
		if err != nil {
			slog.Warn("生成缩略图失败", "id", id, "error", err)
		} else {
			if err := s.storage.SaveThumbnail(ctx, id, thumb, format); err != nil {
				slog.Warn("保存缩略图失败", "id", id, "error", err)
			} else {
				thumbnailURL = fmt.Sprintf("%s/i/%s.%s?t=thumb", baseURL, id, format)
			}
		}
	}

	result := &model.UploadResult{
		ID:           id,
		Name:         filename,
		URL:          fmt.Sprintf("%s/i/%s.%s", baseURL, id, format),
		ThumbnailURL: thumbnailURL,
		Size:         meta.Size,
		Width:        meta.Width,
		Height:       meta.Height,
		Format:       format,
		DeleteToken:  deleteToken,
		CreatedAt:    meta.CreatedAt,
		ExpiresAt:    meta.ExpiresAt,
	}

	return result, nil
}

// GetImage 获取图片二进制流和元数据
func (s *ImageService) GetImage(ctx context.Context, id string) (io.ReadCloser, *model.ImageMeta, error) {
	return s.storage.Get(ctx, id)
}

// GetMeta 获取图片元数据
func (s *ImageService) GetMeta(ctx context.Context, id string) (*model.ImageMeta, error) {
	return s.storage.GetMeta(ctx, id)
}

// GetThumbnail 获取缩略图二进制流
func (s *ImageService) GetThumbnail(ctx context.Context, id string) (io.ReadCloser, error) {
	return s.storage.GetThumbnail(ctx, id)
}

// Delete 验证删除令牌后删除图片
func (s *ImageService) Delete(ctx context.Context, id string, token string) error {
	meta, err := s.storage.GetMeta(ctx, id)
	if err != nil {
		return fmt.Errorf("图片不存在: %w", err)
	}

	if meta.DeleteToken != token {
		return fmt.Errorf("invalid delete token")
	}

	return s.storage.Delete(ctx, id)
}

// ListRecent 获取最近上传的图片列表
func (s *ImageService) ListRecent(ctx context.Context, limit int) ([]*model.ImageMeta, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 50 {
		limit = 50
	}
	return s.storage.List(ctx, limit)
}

// applyWatermark 解码图片 -> 叠加水印 -> 重新编码
func (s *ImageService) applyWatermark(reader io.Reader, format string) (io.Reader, error) {
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("水印: 解码图片失败: %w", err)
	}

	watermarked, err := processor.ApplyWatermark(img, &s.cfg.Processing.Watermark)
	if err != nil {
		return nil, fmt.Errorf("水印: 叠加失败: %w", err)
	}

	var buf bytes.Buffer
	switch format {
	case "jpg", "jpeg":
		err = jpeg.Encode(&buf, watermarked, &jpeg.Options{Quality: s.cfg.Processing.CompressQuality})
	case "png":
		err = png.Encode(&buf, watermarked)
	default:
		// WebP 等其他格式回退为 JPEG
		err = jpeg.Encode(&buf, watermarked, &jpeg.Options{Quality: s.cfg.Processing.CompressQuality})
	}
	if err != nil {
		return nil, fmt.Errorf("水印: 编码图片失败: %w", err)
	}

	return &buf, nil
}
