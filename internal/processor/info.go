package processor

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"

	_ "golang.org/x/image/webp"
)

// ImageInfo 图片基本信息
type ImageInfo struct {
	Width  int
	Height int
	Format string
}

// GetInfo 提取图片的宽高和格式信息
// SVG 文件无法作为光栅图像解码，需传入 isSVG=true
func GetInfo(r io.ReadSeeker, isSVG bool) (*ImageInfo, error) {
	if isSVG {
		return &ImageInfo{Width: 0, Height: 0, Format: "svg"}, nil
	}

	img, format, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("解码图片失败: %w", err)
	}

	// 将读取位置重置到开头，以便后续复用
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("重置读取位置失败: %w", err)
	}

	bounds := img.Bounds()
	return &ImageInfo{
		Width:  bounds.Dx(),
		Height: bounds.Dy(),
		Format: format,
	}, nil
}
