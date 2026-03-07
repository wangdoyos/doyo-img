package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config 应用全局配置结构体
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Storage    StorageConfig    `mapstructure:"storage"`
	Upload     UploadConfig     `mapstructure:"upload"`
	Processing ProcessingConfig `mapstructure:"processing"`
	ID         IDConfig         `mapstructure:"id"`
	CORS       CORSConfig       `mapstructure:"cors"`
	RateLimit  RateLimitConfig  `mapstructure:"rate_limit"`
	Cleanup    CleanupConfig    `mapstructure:"cleanup"`
	Log        LogConfig        `mapstructure:"log"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	BaseURL string `mapstructure:"base_url"`
}

// StorageConfig 存储配置
type StorageConfig struct {
	Type  string      `mapstructure:"type"`
	Local LocalConfig `mapstructure:"local"`
	S3    S3Config    `mapstructure:"s3"`
}

// LocalConfig 本地文件系统存储配置
type LocalConfig struct {
	DataDir string `mapstructure:"data_dir"`
}

// S3Config S3 兼容对象存储配置
type S3Config struct {
	Endpoint   string `mapstructure:"endpoint"`
	Bucket     string `mapstructure:"bucket"`
	Region     string `mapstructure:"region"`
	AccessKey  string `mapstructure:"access_key"`
	SecretKey  string `mapstructure:"secret_key"`
	UseSSL     bool   `mapstructure:"use_ssl"`
	PathPrefix string `mapstructure:"path_prefix"`
}

// UploadConfig 上传限制配置
type UploadConfig struct {
	MaxFileSize        int64    `mapstructure:"max_file_size"`
	MaxBatchSize       int      `mapstructure:"max_batch_size"`
	AllowedFormats     []string `mapstructure:"allowed_formats"`
	DefaultExpireHours int      `mapstructure:"default_expire_hours"`
	MaxExpireDays      int      `mapstructure:"max_expire_days"`
}

// ProcessingConfig 图片处理配置
type ProcessingConfig struct {
	CompressEnabled bool            `mapstructure:"compress_enabled"`
	CompressQuality int             `mapstructure:"compress_quality"`
	StripExif       bool            `mapstructure:"strip_exif"`
	Thumbnail       ThumbnailConfig `mapstructure:"thumbnail"`
	Watermark       WatermarkConfig `mapstructure:"watermark"`
}

// ThumbnailConfig 缩略图配置
type ThumbnailConfig struct {
	Enabled   bool `mapstructure:"enabled"`
	MaxWidth  int  `mapstructure:"max_width"`
	MaxHeight int  `mapstructure:"max_height"`
}

// WatermarkConfig 水印配置
type WatermarkConfig struct {
	Enabled   bool    `mapstructure:"enabled"`
	Type      string  `mapstructure:"type"`       // "text" or "image"
	Text      string  `mapstructure:"text"`       // 文本水印内容
	FontPath  string  `mapstructure:"font_path"`  // 自定义字体路径（TTF/OTF），支持中文水印
	FontSize  float64 `mapstructure:"font_size"`  // 字体大小
	Opacity   float64 `mapstructure:"opacity"`    // 不透明度 0.0~1.0
	Position  string  `mapstructure:"position"`   // top-left, top-right, bottom-left, bottom-right, center
	ImagePath string  `mapstructure:"image_path"` // 图片水印路径
	Padding   int     `mapstructure:"padding"`    // 距离边缘的内边距
}

// IDConfig 图片ID生成配置
type IDConfig struct {
	Length   int    `mapstructure:"length"`
	Alphabet string `mapstructure:"alphabet"`
}

// CORSConfig 跨域配置
type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

// RateLimitConfig IP限流配置
type RateLimitConfig struct {
	Enabled           bool `mapstructure:"enabled"`
	RequestsPerMinute int  `mapstructure:"requests_per_minute"`
	Burst             int  `mapstructure:"burst"`
}

// CleanupConfig 过期图片清理配置
type CleanupConfig struct {
	Enabled       bool `mapstructure:"enabled"`
	RetentionDays int  `mapstructure:"retention_days"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level string `mapstructure:"level"`
}

// Load 加载配置文件，支持 YAML 文件 + 环境变量覆盖
func Load() (*Config, error) {
	v := viper.New()

	// 设置默认值
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 9090)
	v.SetDefault("server.base_url", "")
	v.SetDefault("storage.type", "local")
	v.SetDefault("storage.local.data_dir", "./data")
	v.SetDefault("storage.s3.use_ssl", true)
	v.SetDefault("storage.s3.path_prefix", "images")
	v.SetDefault("upload.max_file_size", 5242880)
	v.SetDefault("upload.max_batch_size", 10)
	v.SetDefault("upload.allowed_formats", []string{"jpg", "jpeg", "png", "gif", "webp", "svg"})
	v.SetDefault("upload.default_expire_hours", 0)
	v.SetDefault("upload.max_expire_days", 0)
	v.SetDefault("processing.compress_enabled", false)
	v.SetDefault("processing.compress_quality", 85)
	v.SetDefault("processing.strip_exif", true)
	v.SetDefault("processing.thumbnail.enabled", true)
	v.SetDefault("processing.thumbnail.max_width", 300)
	v.SetDefault("processing.thumbnail.max_height", 300)
	v.SetDefault("processing.watermark.enabled", false)
	v.SetDefault("processing.watermark.type", "text")
	v.SetDefault("processing.watermark.text", "doyo-img")
	v.SetDefault("processing.watermark.font_path", "")
	v.SetDefault("processing.watermark.font_size", 24.0)
	v.SetDefault("processing.watermark.opacity", 0.3)
	v.SetDefault("processing.watermark.position", "bottom-right")
	v.SetDefault("processing.watermark.image_path", "")
	v.SetDefault("processing.watermark.padding", 20)
	v.SetDefault("id.length", 8)
	v.SetDefault("id.alphabet", "0123456789abcdefghijklmnopqrstuvwxyz")
	v.SetDefault("cors.allowed_origins", []string{"*"})
	v.SetDefault("rate_limit.enabled", true)
	v.SetDefault("rate_limit.requests_per_minute", 30)
	v.SetDefault("rate_limit.burst", 10)
	v.SetDefault("cleanup.enabled", false)
	v.SetDefault("cleanup.retention_days", 30)
	v.SetDefault("log.level", "info")

	// 配置文件路径
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	// 环境变量覆盖，前缀 DOYO_，例如 DOYO_SERVER_PORT=8080
	v.SetEnvPrefix("DOYO")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
		// 配置文件不存在时使用默认值，不报错
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	// 本地存储模式下，确保数据目录存在
	if cfg.Storage.Type == "local" {
		if err := os.MkdirAll(cfg.Storage.Local.DataDir, 0755); err != nil {
			return nil, fmt.Errorf("创建数据目录失败: %w", err)
		}
	}

	return &cfg, nil
}

// Address 返回服务监听地址，格式为 host:port
func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
