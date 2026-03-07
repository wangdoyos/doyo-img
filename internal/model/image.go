package model

import "time"

// ImageMeta 图片元数据，以 JSON 文件形式存储在图片同目录下
type ImageMeta struct {
	ID           string     `json:"id"`
	OriginalName string     `json:"name"`
	Format       string     `json:"format"`
	MimeType     string     `json:"mime_type"`
	Size         int64      `json:"size"`
	Width        int        `json:"width"`
	Height       int        `json:"height"`
	StoragePath  string     `json:"storage_path"`
	DeleteToken  string     `json:"delete_token"`
	CreatedAt    time.Time  `json:"created_at"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
}

// UploadResult 上传成功后返回给前端的结果
type UploadResult struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	URL          string     `json:"url"`
	ThumbnailURL string     `json:"thumbnail_url,omitempty"`
	Size         int64      `json:"size"`
	Width        int        `json:"width"`
	Height       int        `json:"height"`
	Format       string     `json:"format"`
	DeleteToken  string     `json:"delete_token"`
	CreatedAt    time.Time  `json:"created_at"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
}
