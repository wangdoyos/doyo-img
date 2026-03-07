package processor

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/disintegration/imaging"
)

// Compress 按指定质量压缩图片
// 仅对 JPEG 和 PNG 有效，GIF/WebP/SVG 返回原始数据
func Compress(r io.ReadSeeker, format string, quality int) (io.Reader, error) {
	if format == "svg" || format == "gif" || format == "webp" {
		if _, err := r.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}
		return r, nil
	}

	img, _, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("压缩时解码图片失败: %w", err)
	}

	var buf bytes.Buffer
	switch format {
	case "jpg", "jpeg":
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
	case "png":
		err = png.Encode(&buf, img)
	default:
		if _, err := r.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}
		return r, nil
	}

	if err != nil {
		return nil, fmt.Errorf("编码压缩图片失败: %w", err)
	}

	return &buf, nil
}

// GenerateThumbnail 生成适应 maxWidth x maxHeight 尺寸的缩略图
func GenerateThumbnail(r io.ReadSeeker, format string, maxWidth, maxHeight int) (io.Reader, error) {
	if format == "svg" {
		return nil, fmt.Errorf("SVG 不支持生成缩略图")
	}

	img, _, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("生成缩略图时解码失败: %w", err)
	}

	thumb := imaging.Fit(img, maxWidth, maxHeight, imaging.Lanczos)

	var buf bytes.Buffer
	switch format {
	case "jpg", "jpeg":
		err = jpeg.Encode(&buf, thumb, &jpeg.Options{Quality: 80})
	case "png":
		err = png.Encode(&buf, thumb)
	case "gif":
		err = gif.Encode(&buf, thumb, nil)
	default:
		err = jpeg.Encode(&buf, thumb, &jpeg.Options{Quality: 80})
	}

	if err != nil {
		return nil, fmt.Errorf("编码缩略图失败: %w", err)
	}

	return &buf, nil
}
