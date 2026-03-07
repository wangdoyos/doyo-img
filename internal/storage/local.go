package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/wangdoyos/doyo-img/internal/model"
)

// LocalStorage 本地文件系统存储实现
type LocalStorage struct {
	dataDir string
	mu      sync.RWMutex
	index   map[string]string // id -> 相对路径（不含扩展名）
}

// NewLocalStorage 创建本地存储实例，启动时扫描数据目录构建内存索引
func NewLocalStorage(dataDir string) (*LocalStorage, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("创建数据目录失败: %w", err)
	}

	ls := &LocalStorage{
		dataDir: dataDir,
		index:   make(map[string]string),
	}

	if err := ls.buildIndex(); err != nil {
		return nil, fmt.Errorf("构建索引失败: %w", err)
	}

	return ls, nil
}

// buildIndex 遍历数据目录，根据 .json 元数据文件构建 ID -> 路径 的内存映射
func (ls *LocalStorage) buildIndex() error {
	return filepath.Walk(ls.dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 跳过错误
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".json") {
			return nil
		}

		rel, err := filepath.Rel(ls.dataDir, path)
		if err != nil {
			return nil
		}

		// 从文件名提取 ID："a1b2c3d4.json" -> "a1b2c3d4"
		base := filepath.Base(rel)
		id := strings.TrimSuffix(base, ".json")

		// 存储不含扩展名的相对路径
		pathWithoutExt := strings.TrimSuffix(rel, ".json")
		ls.index[id] = filepath.ToSlash(pathWithoutExt)

		return nil
	})
}

// datePath 返回基于当前日期的存储子目录路径：yyyy/mm/dd
func (ls *LocalStorage) datePath() string {
	now := time.Now()
	return fmt.Sprintf("%d/%02d/%02d", now.Year(), now.Month(), now.Day())
}

// Save 保存图片文件和元数据到本地文件系统
func (ls *LocalStorage) Save(ctx context.Context, id string, data io.Reader, meta *model.ImageMeta) error {
	datePath := ls.datePath()
	dir := filepath.Join(ls.dataDir, datePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 保存图片文件
	ext := meta.Format
	if ext == "jpg" {
		ext = "jpg"
	}
	imgPath := filepath.Join(dir, id+"."+ext)
	f, err := os.Create(imgPath)
	if err != nil {
		return fmt.Errorf("创建图片文件失败: %w", err)
	}
	defer f.Close()

	written, err := io.Copy(f, data)
	if err != nil {
		os.Remove(imgPath)
		return fmt.Errorf("写入图片数据失败: %w", err)
	}
	meta.Size = written

	// 更新存储路径
	relPath := filepath.ToSlash(filepath.Join(datePath, id))
	meta.StoragePath = relPath + "." + ext

	// 保存元数据 JSON 文件
	metaPath := filepath.Join(dir, id+".json")
	metaData, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		os.Remove(imgPath)
		return fmt.Errorf("序列化元数据失败: %w", err)
	}
	if err := os.WriteFile(metaPath, metaData, 0644); err != nil {
		os.Remove(imgPath)
		return fmt.Errorf("写入元数据失败: %w", err)
	}

	// 更新内存索引
	ls.mu.Lock()
	ls.index[id] = filepath.ToSlash(filepath.Join(datePath, id))
	ls.mu.Unlock()

	return nil
}

// Get 获取图片二进制流和元数据
func (ls *LocalStorage) Get(ctx context.Context, id string) (io.ReadCloser, *model.ImageMeta, error) {
	meta, err := ls.GetMeta(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	imgPath := filepath.Join(ls.dataDir, filepath.FromSlash(meta.StoragePath))
	f, err := os.Open(imgPath)
	if err != nil {
		return nil, nil, fmt.Errorf("打开图片文件失败: %w", err)
	}

	return f, meta, nil
}

// GetMeta 从 JSON 文件读取图片元数据
func (ls *LocalStorage) GetMeta(ctx context.Context, id string) (*model.ImageMeta, error) {
	ls.mu.RLock()
	relPath, exists := ls.index[id]
	ls.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("图片不存在: %s", id)
	}

	metaPath := filepath.Join(ls.dataDir, filepath.FromSlash(relPath)+".json")
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, fmt.Errorf("读取元数据失败: %w", err)
	}

	var meta model.ImageMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("解析元数据失败: %w", err)
	}

	return &meta, nil
}

// Delete 删除图片文件、缩略图和元数据
func (ls *LocalStorage) Delete(ctx context.Context, id string) error {
	ls.mu.RLock()
	relPath, exists := ls.index[id]
	ls.mu.RUnlock()

	if !exists {
		return fmt.Errorf("图片不存在: %s", id)
	}

	// 读取元数据获取实际文件扩展名
	meta, err := ls.GetMeta(ctx, id)
	if err != nil {
		return err
	}

	// 删除图片文件
	imgPath := filepath.Join(ls.dataDir, filepath.FromSlash(meta.StoragePath))
	os.Remove(imgPath)

	// 删除缩略图（如果存在）
	dir := filepath.Dir(imgPath)
	ext := filepath.Ext(imgPath)
	thumbPath := filepath.Join(dir, id+"_thumb"+ext)
	os.Remove(thumbPath)

	// 删除元数据文件
	metaPath := filepath.Join(ls.dataDir, filepath.FromSlash(relPath)+".json")
	os.Remove(metaPath)

	// 从内存索引中移除
	ls.mu.Lock()
	delete(ls.index, id)
	ls.mu.Unlock()

	return nil
}

// List 返回最近上传的图片列表，按创建时间降序排列
func (ls *LocalStorage) List(ctx context.Context, limit int) ([]*model.ImageMeta, error) {
	ls.mu.RLock()
	ids := make([]string, 0, len(ls.index))
	for id := range ls.index {
		ids = append(ids, id)
	}
	ls.mu.RUnlock()

	// 加载所有元数据
	metas := make([]*model.ImageMeta, 0, len(ids))
	for _, id := range ids {
		meta, err := ls.GetMeta(ctx, id)
		if err != nil {
			continue
		}
		metas = append(metas, meta)
	}

	// 按创建时间降序排列
	sort.Slice(metas, func(i, j int) bool {
		return metas[i].CreatedAt.After(metas[j].CreatedAt)
	})

	if limit > 0 && limit < len(metas) {
		metas = metas[:limit]
	}

	return metas, nil
}

// Exists 检查指定 ID 的图片是否存在
func (ls *LocalStorage) Exists(ctx context.Context, id string) bool {
	ls.mu.RLock()
	_, exists := ls.index[id]
	ls.mu.RUnlock()
	return exists
}

// SaveThumbnail 保存缩略图文件
func (ls *LocalStorage) SaveThumbnail(ctx context.Context, id string, data io.Reader, format string) error {
	ls.mu.RLock()
	relPath, exists := ls.index[id]
	ls.mu.RUnlock()

	if !exists {
		return fmt.Errorf("图片不存在: %s", id)
	}

	dir := filepath.Join(ls.dataDir, filepath.Dir(filepath.FromSlash(relPath)))
	ext := format
	if ext == "jpg" {
		ext = "jpg"
	}
	thumbPath := filepath.Join(dir, id+"_thumb."+ext)

	f, err := os.Create(thumbPath)
	if err != nil {
		return fmt.Errorf("创建缩略图文件失败: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, data); err != nil {
		os.Remove(thumbPath)
		return fmt.Errorf("写入缩略图失败: %w", err)
	}

	return nil
}

// GetThumbnail 获取缩略图二进制流
func (ls *LocalStorage) GetThumbnail(ctx context.Context, id string) (io.ReadCloser, error) {
	meta, err := ls.GetMeta(ctx, id)
	if err != nil {
		return nil, err
	}

	// 构建缩略图路径
	dir := filepath.Dir(filepath.Join(ls.dataDir, filepath.FromSlash(meta.StoragePath)))
	ext := meta.Format
	thumbPath := filepath.Join(dir, id+"_thumb."+ext)

	f, err := os.Open(thumbPath)
	if err != nil {
		return nil, fmt.Errorf("缩略图不存在: %w", err)
	}

	return f, nil
}
