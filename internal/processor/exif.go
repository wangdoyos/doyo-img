package processor

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
)

// StripExif 通过 decode→encode 流程剥离 JPEG 图片中的 EXIF 元数据（GPS、设备信息等）
// 仅对 JPEG 格式生效，其他格式直接返回原始 reader
func StripExif(r io.ReadSeeker, format string) (io.Reader, error) {
	if format != "jpg" && format != "jpeg" {
		if _, err := r.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}
		return r, nil
	}

	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	img, _, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("EXIF 剥离时解码图片失败: %w", err)
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 100}); err != nil {
		return nil, fmt.Errorf("EXIF 剥离时编码图片失败: %w", err)
	}

	return &buf, nil
}
