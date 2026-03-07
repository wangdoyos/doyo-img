package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/wangdoyos/doyo-img/internal/config"
	"github.com/wangdoyos/doyo-img/internal/model"
)

// S3Storage S3 兼容对象存储实现，支持 AWS S3、阿里云 OSS、腾讯 COS、MinIO
type S3Storage struct {
	client     *s3.Client
	bucket     string
	pathPrefix string
}

// NewS3Storage 创建 S3 存储实例并校验连接
func NewS3Storage(cfg *config.S3Config) (*S3Storage, error) {
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("S3 bucket 不能为空")
	}
	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("S3 endpoint 不能为空")
	}
	if cfg.AccessKey == "" || cfg.SecretKey == "" {
		return nil, fmt.Errorf("S3 access_key 和 secret_key 不能为空")
	}

	// 构建自定义 endpoint resolver
	endpoint := cfg.Endpoint
	if cfg.UseSSL && !strings.HasPrefix(endpoint, "https://") && !strings.HasPrefix(endpoint, "http://") {
		endpoint = "https://" + endpoint
	} else if !cfg.UseSSL && !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "http://" + endpoint
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("初始化 AWS 配置失败: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true // MinIO 等兼容服务需要路径风格
	})

	ss := &S3Storage{
		client:     client,
		bucket:     cfg.Bucket,
		pathPrefix: strings.TrimRight(cfg.PathPrefix, "/"),
	}

	// 校验连接：HeadBucket
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(cfg.Bucket),
	})
	if err != nil {
		return nil, fmt.Errorf("S3 连接校验失败（HeadBucket）: %w", err)
	}

	slog.Info("S3 存储初始化成功", "endpoint", endpoint, "bucket", cfg.Bucket, "prefix", cfg.PathPrefix)
	return ss, nil
}

// metaKey 返回元数据 JSON 在 S3 上的 key
func (ss *S3Storage) metaKey(id string) string {
	return ss.pathPrefix + "/_meta/" + id + ".json"
}

// datePath 返回基于当前日期的存储子路径
func (ss *S3Storage) datePath() string {
	now := time.Now()
	return fmt.Sprintf("%d/%02d/%02d", now.Year(), now.Month(), now.Day())
}

// Save 保存图片文件及其元数据到 S3
func (ss *S3Storage) Save(ctx context.Context, id string, data io.Reader, meta *model.ImageMeta) error {
	datePath := ss.datePath()

	// 读取全部数据到内存（S3 PutObject 需要知道 ContentLength 或使用 io.ReadSeeker）
	bodyBytes, err := io.ReadAll(data)
	if err != nil {
		return fmt.Errorf("读取图片数据失败: %w", err)
	}
	meta.Size = int64(len(bodyBytes))

	// 图片 key
	ext := meta.Format
	imgKey := fmt.Sprintf("%s/%s/%s.%s", ss.pathPrefix, datePath, id, ext)
	meta.StoragePath = imgKey

	// 上传图片文件
	_, err = ss.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(ss.bucket),
		Key:         aws.String(imgKey),
		Body:        bytes.NewReader(bodyBytes),
		ContentType: aws.String(meta.MimeType),
	})
	if err != nil {
		return fmt.Errorf("上传图片到 S3 失败: %w", err)
	}

	// 序列化并上传元数据
	metaData, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化元数据失败: %w", err)
	}

	metaKey := ss.metaKey(id)
	_, err = ss.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(ss.bucket),
		Key:         aws.String(metaKey),
		Body:        bytes.NewReader(metaData),
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		// 元数据上传失败时尝试回滚图片
		ss.client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(ss.bucket),
			Key:    aws.String(imgKey),
		})
		return fmt.Errorf("上传元数据到 S3 失败: %w", err)
	}

	return nil
}

// Get 获取图片二进制流和元数据
func (ss *S3Storage) Get(ctx context.Context, id string) (io.ReadCloser, *model.ImageMeta, error) {
	meta, err := ss.GetMeta(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	output, err := ss.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(ss.bucket),
		Key:    aws.String(meta.StoragePath),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("从 S3 获取图片失败: %w", err)
	}

	return output.Body, meta, nil
}

// GetMeta 获取图片元数据（从 _meta/ 目录直接读取）
func (ss *S3Storage) GetMeta(ctx context.Context, id string) (*model.ImageMeta, error) {
	metaKey := ss.metaKey(id)

	output, err := ss.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(ss.bucket),
		Key:    aws.String(metaKey),
	})
	if err != nil {
		return nil, fmt.Errorf("图片不存在: %s", id)
	}
	defer output.Body.Close()

	data, err := io.ReadAll(output.Body)
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
func (ss *S3Storage) Delete(ctx context.Context, id string) error {
	meta, err := ss.GetMeta(ctx, id)
	if err != nil {
		return err
	}

	// 收集需要删除的 keys
	keysToDelete := []types.ObjectIdentifier{
		{Key: aws.String(meta.StoragePath)},
		{Key: aws.String(ss.metaKey(id))},
	}

	// 缩略图 key
	thumbKey := ss.thumbnailKey(meta)
	if thumbKey != "" {
		keysToDelete = append(keysToDelete, types.ObjectIdentifier{Key: aws.String(thumbKey)})
	}

	// 批量删除
	_, err = ss.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(ss.bucket),
		Delete: &types.Delete{
			Objects: keysToDelete,
			Quiet:   aws.Bool(true),
		},
	})
	if err != nil {
		return fmt.Errorf("从 S3 删除对象失败: %w", err)
	}

	return nil
}

// List 返回最近上传的图片列表，按创建时间降序排列
func (ss *S3Storage) List(ctx context.Context, limit int) ([]*model.ImageMeta, error) {
	prefix := ss.pathPrefix + "/_meta/"

	var metas []*model.ImageMeta
	var continuationToken *string

	for {
		input := &s3.ListObjectsV2Input{
			Bucket: aws.String(ss.bucket),
			Prefix: aws.String(prefix),
		}
		if continuationToken != nil {
			input.ContinuationToken = continuationToken
		}

		output, err := ss.client.ListObjectsV2(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("列举 S3 对象失败: %w", err)
		}

		for _, obj := range output.Contents {
			key := aws.ToString(obj.Key)
			if !strings.HasSuffix(key, ".json") {
				continue
			}

			getOutput, err := ss.client.GetObject(ctx, &s3.GetObjectInput{
				Bucket: aws.String(ss.bucket),
				Key:    aws.String(key),
			})
			if err != nil {
				slog.Warn("S3 List: 读取元数据失败", "key", key, "error", err)
				continue
			}

			data, err := io.ReadAll(getOutput.Body)
			getOutput.Body.Close()
			if err != nil {
				continue
			}

			var meta model.ImageMeta
			if err := json.Unmarshal(data, &meta); err != nil {
				continue
			}
			metas = append(metas, &meta)
		}

		if !aws.ToBool(output.IsTruncated) {
			break
		}
		continuationToken = output.NextContinuationToken
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
func (ss *S3Storage) Exists(ctx context.Context, id string) bool {
	metaKey := ss.metaKey(id)
	_, err := ss.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(ss.bucket),
		Key:    aws.String(metaKey),
	})
	return err == nil
}

// SaveThumbnail 保存缩略图到 S3
func (ss *S3Storage) SaveThumbnail(ctx context.Context, id string, data io.Reader, format string) error {
	meta, err := ss.GetMeta(ctx, id)
	if err != nil {
		return err
	}

	thumbKey := ss.thumbnailKey(meta)
	if thumbKey == "" {
		return fmt.Errorf("无法计算缩略图路径")
	}

	bodyBytes, err := io.ReadAll(data)
	if err != nil {
		return fmt.Errorf("读取缩略图数据失败: %w", err)
	}

	_, err = ss.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(ss.bucket),
		Key:         aws.String(thumbKey),
		Body:        bytes.NewReader(bodyBytes),
		ContentType: aws.String("image/" + format),
	})
	if err != nil {
		return fmt.Errorf("上传缩略图到 S3 失败: %w", err)
	}

	return nil
}

// GetThumbnail 获取缩略图二进制流
func (ss *S3Storage) GetThumbnail(ctx context.Context, id string) (io.ReadCloser, error) {
	meta, err := ss.GetMeta(ctx, id)
	if err != nil {
		return nil, err
	}

	thumbKey := ss.thumbnailKey(meta)
	if thumbKey == "" {
		return nil, fmt.Errorf("缩略图不存在")
	}

	output, err := ss.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(ss.bucket),
		Key:    aws.String(thumbKey),
	})
	if err != nil {
		return nil, fmt.Errorf("缩略图不存在: %w", err)
	}

	return output.Body, nil
}

// thumbnailKey 根据元数据计算缩略图的 S3 key
func (ss *S3Storage) thumbnailKey(meta *model.ImageMeta) string {
	if meta.StoragePath == "" {
		return ""
	}
	// StoragePath 格式: {prefix}/{date}/{id}.{ext}
	// 缩略图格式: {prefix}/{date}/{id}_thumb.{ext}
	dotIdx := strings.LastIndex(meta.StoragePath, ".")
	if dotIdx < 0 {
		return ""
	}
	return meta.StoragePath[:dotIdx] + "_thumb" + meta.StoragePath[dotIdx:]
}
