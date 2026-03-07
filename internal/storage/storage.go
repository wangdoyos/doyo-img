package storage

import (
	"context"
	"io"

	"github.com/wangdoyos/doyo-img/internal/model"
)

// Storage 定义图片存储后端的统一接口
type Storage interface {
	// Save 保存图片文件及其元数据
	Save(ctx context.Context, id string, data io.Reader, meta *model.ImageMeta) error

	// Get 获取图片二进制流和元数据
	Get(ctx context.Context, id string) (io.ReadCloser, *model.ImageMeta, error)

	// GetMeta 仅获取图片的元数据
	GetMeta(ctx context.Context, id string) (*model.ImageMeta, error)

	// Delete 删除图片文件及其元数据
	Delete(ctx context.Context, id string) error

	// List 返回最近上传的图片列表，按创建时间降序排列
	List(ctx context.Context, limit int) ([]*model.ImageMeta, error)

	// Exists 检查指定 ID 的图片是否存在
	Exists(ctx context.Context, id string) bool

	// SaveThumbnail 保存指定图片的缩略图
	SaveThumbnail(ctx context.Context, id string, data io.Reader, format string) error

	// GetThumbnail 获取缩略图二进制流
	GetThumbnail(ctx context.Context, id string) (io.ReadCloser, error)
}
