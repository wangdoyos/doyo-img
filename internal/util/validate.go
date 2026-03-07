package util

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// MIME 类型到短格式名的映射
var mimeToFormat = map[string]string{
	"image/jpeg":    "jpg",
	"image/png":     "png",
	"image/gif":     "gif",
	"image/webp":    "webp",
	"image/svg+xml": "svg",
}

// 短格式名到 MIME 类型的映射
var formatToMime = map[string]string{
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"png":  "image/png",
	"gif":  "image/gif",
	"webp": "image/webp",
	"svg":  "image/svg+xml",
}

// DetectMimeType 通过读取文件头部 512 字节的 magic bytes 检测真实 MIME 类型
func DetectMimeType(r io.ReadSeeker) (string, error) {
	buf := make([]byte, 512)
	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("读取文件头失败: %w", err)
	}
	// 重置读取位置到开头
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("重置读取位置失败: %w", err)
	}

	mime := http.DetectContentType(buf[:n])

	// http.DetectContentType 对 SVG 返回 "text/xml"，需要手动检测
	if strings.Contains(mime, "text/xml") || strings.Contains(mime, "text/plain") {
		content := string(buf[:n])
		if strings.Contains(content, "<svg") || strings.Contains(content, "xmlns=\"http://www.w3.org/2000/svg\"") {
			return "image/svg+xml", nil
		}
	}

	return mime, nil
}

// MimeToFormat 将 MIME 类型转换为短格式名
func MimeToFormat(mime string) string {
	if f, ok := mimeToFormat[mime]; ok {
		return f
	}
	return ""
}

// FormatToMime 将短格式名转换为 MIME 类型
func FormatToMime(format string) string {
	if m, ok := formatToMime[strings.ToLower(format)]; ok {
		return m
	}
	return "application/octet-stream"
}

// ValidateFormat 检查检测到的 MIME 类型是否在允许的格式列表中
func ValidateFormat(mime string, allowedFormats []string) bool {
	format := MimeToFormat(mime)
	if format == "" {
		return false
	}
	for _, f := range allowedFormats {
		if strings.EqualFold(f, format) || strings.EqualFold(f, "jpeg") && format == "jpg" || strings.EqualFold(f, "jpg") && format == "jpg" {
			return true
		}
	}
	return false
}

// ValidateFileSize 检查文件大小是否在限制范围内
func ValidateFileSize(size int64, maxSize int64) error {
	if size > maxSize {
		return fmt.Errorf("文件大小 %d 超过最大限制 %d", size, maxSize)
	}
	return nil
}
