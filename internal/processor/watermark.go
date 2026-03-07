package processor

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"

	"github.com/disintegration/imaging"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"

	"github.com/wangdoyos/doyo-img/internal/config"
)

// ApplyWatermark 在图片上叠加水印（文本或图片）
func ApplyWatermark(img image.Image, cfg *config.WatermarkConfig) (image.Image, error) {
	if !cfg.Enabled {
		return img, nil
	}

	switch cfg.Type {
	case "text":
		return applyTextWatermark(img, cfg)
	case "image":
		return applyImageWatermark(img, cfg)
	default:
		return nil, fmt.Errorf("不支持的水印类型: %s", cfg.Type)
	}
}

// applyTextWatermark 在图片上叠加文本水印
func applyTextWatermark(img image.Image, cfg *config.WatermarkConfig) (image.Image, error) {
	if cfg.Text == "" {
		return img, nil
	}

	bounds := img.Bounds()
	// 创建可写的 RGBA 画布
	canvas := image.NewRGBA(bounds)
	draw.Draw(canvas, bounds, img, bounds.Min, draw.Src)

	// 加载字体
	face, err := loadFontFace(cfg.FontPath, cfg.FontSize)
	if err != nil {
		return nil, fmt.Errorf("加载字体失败: %w", err)
	}

	// 计算文本尺寸
	textWidth := measureTextWidth(face, cfg.Text)
	metrics := face.Metrics()
	textHeight := metrics.Ascent.Ceil() + metrics.Descent.Ceil()

	// 计算水印位置
	x, y := calcPosition(
		bounds.Dx(), bounds.Dy(),
		textWidth, textHeight,
		cfg.Padding, cfg.Position,
	)

	// 应用透明度
	alpha := uint8(math.Round(cfg.Opacity * 255))
	textColor := color.NRGBA{R: 255, G: 255, B: 255, A: alpha}

	// 绘制文本（带半透明黑色阴影增强可读性）
	shadowColor := color.NRGBA{R: 0, G: 0, B: 0, A: alpha / 2}
	drawText(canvas, face, cfg.Text, x+1, y+metrics.Ascent.Ceil()+1, shadowColor)
	drawText(canvas, face, cfg.Text, x, y+metrics.Ascent.Ceil(), textColor)

	return canvas, nil
}

// applyImageWatermark 在图片上叠加图片水印
func applyImageWatermark(img image.Image, cfg *config.WatermarkConfig) (image.Image, error) {
	if cfg.ImagePath == "" {
		return nil, fmt.Errorf("图片水印路径不能为空")
	}

	// 加载水印图片
	wmFile, err := os.Open(cfg.ImagePath)
	if err != nil {
		return nil, fmt.Errorf("打开水印图片失败: %w", err)
	}
	defer wmFile.Close()

	wmImg, err := png.Decode(wmFile)
	if err != nil {
		return nil, fmt.Errorf("解码水印图片失败（仅支持 PNG）: %w", err)
	}

	bounds := img.Bounds()
	wmBounds := wmImg.Bounds()

	// 计算水印位置
	x, y := calcPosition(
		bounds.Dx(), bounds.Dy(),
		wmBounds.Dx(), wmBounds.Dy(),
		cfg.Padding, cfg.Position,
	)

	// 使用 imaging 库叠加水印
	bgImg := imaging.Clone(img)
	result := imaging.Overlay(bgImg, wmImg, image.Pt(x, y), cfg.Opacity)

	return result, nil
}

// loadFontFace 加载字体，支持外部 TTF/OTF 文件，无指定时使用内置基础字体
func loadFontFace(fontPath string, fontSize float64) (font.Face, error) {
	if fontPath == "" {
		// 使用 Go 内置的基础字体（仅支持 ASCII）
		return basicfont.Face7x13, nil
	}

	// 读取外部字体文件
	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		return nil, fmt.Errorf("读取字体文件失败: %w", err)
	}

	f, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, fmt.Errorf("解析字体文件失败: %w", err)
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("创建字体 Face 失败: %w", err)
	}

	return face, nil
}

// measureTextWidth 测量文本渲染宽度
func measureTextWidth(face font.Face, text string) int {
	width := fixed.Int26_6(0)
	for _, r := range text {
		adv, ok := face.GlyphAdvance(r)
		if ok {
			width += adv
		}
	}
	return width.Ceil()
}

// drawText 在 RGBA 画布上绘制文本
func drawText(canvas *image.RGBA, face font.Face, text string, x, y int, col color.Color) {
	d := &font.Drawer{
		Dst:  canvas,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  fixed.P(x, y),
	}
	d.DrawString(text)
}

// calcPosition 根据位置策略计算水印左上角坐标
func calcPosition(bgW, bgH, wmW, wmH, padding int, position string) (x, y int) {
	switch position {
	case "top-left":
		return padding, padding
	case "top-right":
		return bgW - wmW - padding, padding
	case "bottom-left":
		return padding, bgH - wmH - padding
	case "center":
		return (bgW - wmW) / 2, (bgH - wmH) / 2
	default: // bottom-right
		return bgW - wmW - padding, bgH - wmH - padding
	}
}
