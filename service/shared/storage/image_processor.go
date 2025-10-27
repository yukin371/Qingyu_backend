package storage

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

// ImageProcessor 图片处理服务
type ImageProcessor struct {
	backend StorageBackend
}

// NewImageProcessor 创建图片处理服务
func NewImageProcessor(backend StorageBackend) *ImageProcessor {
	return &ImageProcessor{
		backend: backend,
	}
}

// ImageProcessOptions 图片处理选项
type ImageProcessOptions struct {
	// 缩略图选项
	ThumbnailWidth  int  // 缩略图宽度
	ThumbnailHeight int  // 缩略图高度
	KeepAspectRatio bool // 保持宽高比

	// 压缩选项
	Quality        int  // JPEG质量 (1-100)
	EnableCompress bool // 是否启用压缩

	// 裁剪选项
	CropX      int  // 裁剪起点X
	CropY      int  // 裁剪起点Y
	CropWidth  int  // 裁剪宽度
	CropHeight int  // 裁剪高度
	EnableCrop bool // 是否启用裁剪

	// 水印选项
	WatermarkPath    string  // 水印图片路径
	WatermarkPos     string  // 水印位置: top-left, top-right, bottom-left, bottom-right, center
	WatermarkOpacity float64 // 水印透明度 (0-1)

	// 格式转换
	OutputFormat string // 输出格式: jpeg, png
}

// ProcessImageRequest 处理图片请求
type ProcessImageRequest struct {
	SourcePath string               // 源图片路径
	DestPath   string               // 目标图片路径
	Options    *ImageProcessOptions // 处理选项
}

// ProcessImageResponse 处理图片响应
type ProcessImageResponse struct {
	DestPath string `json:"dest_path"` // 处理后的图片路径
	Width    int    `json:"width"`     // 图片宽度
	Height   int    `json:"height"`    // 图片高度
	Size     int64  `json:"size"`      // 文件大小
}

// ProcessImage 处理图片
func (p *ImageProcessor) ProcessImage(ctx context.Context, req *ProcessImageRequest) (*ProcessImageResponse, error) {
	// 1. 加载源图片
	reader, err := p.backend.Load(ctx, req.SourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load source image: %w", err)
	}
	defer reader.Close()

	// 2. 解码图片
	img, format, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// 3. 应用处理选项
	processedImg := img
	if req.Options != nil {
		processedImg, err = p.applyOptions(img, req.Options)
		if err != nil {
			return nil, fmt.Errorf("failed to process image: %w", err)
		}
	}

	// 4. 确定输出格式
	outputFormat := format
	if req.Options != nil && req.Options.OutputFormat != "" {
		outputFormat = req.Options.OutputFormat
	}

	// 5. 编码并保存
	var buf bytes.Buffer
	err = p.encodeImage(&buf, processedImg, outputFormat, req.Options)
	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	// 6. 保存到存储后端
	err = p.backend.Save(ctx, req.DestPath, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to save processed image: %w", err)
	}

	// 7. 返回结果
	bounds := processedImg.Bounds()
	return &ProcessImageResponse{
		DestPath: req.DestPath,
		Width:    bounds.Dx(),
		Height:   bounds.Dy(),
		Size:     int64(buf.Len()),
	}, nil
}

// GenerateThumbnail 生成缩略图
func (p *ImageProcessor) GenerateThumbnail(ctx context.Context, sourcePath string, width, height int, keepAspectRatio bool) (string, error) {
	// 生成缩略图路径
	thumbnailPath := p.getThumbnailPath(sourcePath, width, height)

	// 处理图片
	_, err := p.ProcessImage(ctx, &ProcessImageRequest{
		SourcePath: sourcePath,
		DestPath:   thumbnailPath,
		Options: &ImageProcessOptions{
			ThumbnailWidth:  width,
			ThumbnailHeight: height,
			KeepAspectRatio: keepAspectRatio,
			EnableCompress:  true,
			Quality:         85,
		},
	})
	if err != nil {
		return "", err
	}

	return thumbnailPath, nil
}

// CompressImage 压缩图片
func (p *ImageProcessor) CompressImage(ctx context.Context, sourcePath string, quality int) (string, error) {
	// 生成压缩后的路径
	compressedPath := p.getCompressedPath(sourcePath)

	// 处理图片
	_, err := p.ProcessImage(ctx, &ProcessImageRequest{
		SourcePath: sourcePath,
		DestPath:   compressedPath,
		Options: &ImageProcessOptions{
			EnableCompress: true,
			Quality:        quality,
		},
	})
	if err != nil {
		return "", err
	}

	return compressedPath, nil
}

// CropImage 裁剪图片
func (p *ImageProcessor) CropImage(ctx context.Context, sourcePath string, x, y, width, height int) (string, error) {
	// 生成裁剪后的路径
	croppedPath := p.getCroppedPath(sourcePath)

	// 处理图片
	_, err := p.ProcessImage(ctx, &ProcessImageRequest{
		SourcePath: sourcePath,
		DestPath:   croppedPath,
		Options: &ImageProcessOptions{
			EnableCrop: true,
			CropX:      x,
			CropY:      y,
			CropWidth:  width,
			CropHeight: height,
		},
	})
	if err != nil {
		return "", err
	}

	return croppedPath, nil
}

// ConvertFormat 转换图片格式
func (p *ImageProcessor) ConvertFormat(ctx context.Context, sourcePath string, outputFormat string) (string, error) {
	// 生成新格式路径
	convertedPath := p.getConvertedPath(sourcePath, outputFormat)

	// 处理图片
	_, err := p.ProcessImage(ctx, &ProcessImageRequest{
		SourcePath: sourcePath,
		DestPath:   convertedPath,
		Options: &ImageProcessOptions{
			OutputFormat: outputFormat,
		},
	})
	if err != nil {
		return "", err
	}

	return convertedPath, nil
}

// ============ 私有方法 ============

// applyOptions 应用处理选项
func (p *ImageProcessor) applyOptions(img image.Image, opts *ImageProcessOptions) (image.Image, error) {
	result := img

	// 1. 裁剪
	if opts.EnableCrop {
		result = imaging.Crop(result, image.Rect(
			opts.CropX,
			opts.CropY,
			opts.CropX+opts.CropWidth,
			opts.CropY+opts.CropHeight,
		))
	}

	// 2. 缩略图/调整大小
	if opts.ThumbnailWidth > 0 || opts.ThumbnailHeight > 0 {
		if opts.KeepAspectRatio {
			// 保持宽高比
			if opts.ThumbnailWidth > 0 && opts.ThumbnailHeight > 0 {
				result = imaging.Fit(result, opts.ThumbnailWidth, opts.ThumbnailHeight, imaging.Lanczos)
			} else if opts.ThumbnailWidth > 0 {
				result = imaging.Resize(result, opts.ThumbnailWidth, 0, imaging.Lanczos)
			} else {
				result = imaging.Resize(result, 0, opts.ThumbnailHeight, imaging.Lanczos)
			}
		} else {
			// 不保持宽高比，强制缩放
			width := opts.ThumbnailWidth
			height := opts.ThumbnailHeight
			if width == 0 {
				width = result.Bounds().Dx()
			}
			if height == 0 {
				height = result.Bounds().Dy()
			}
			result = imaging.Resize(result, width, height, imaging.Lanczos)
		}
	}

	// 3. 水印（暂未实现）
	// if opts.WatermarkPath != "" {
	//     result = p.applyWatermark(result, opts)
	// }

	return result, nil
}

// encodeImage 编码图片
func (p *ImageProcessor) encodeImage(w io.Writer, img image.Image, format string, opts *ImageProcessOptions) error {
	quality := 90 // 默认质量
	if opts != nil && opts.Quality > 0 {
		quality = opts.Quality
	}

	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		return jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
	case "png":
		encoder := png.Encoder{CompressionLevel: png.DefaultCompression}
		return encoder.Encode(w, img)
	default:
		// 默认使用JPEG
		return jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
	}
}

// getThumbnailPath 获取缩略图路径
func (p *ImageProcessor) getThumbnailPath(sourcePath string, width, height int) string {
	ext := filepath.Ext(sourcePath)
	base := strings.TrimSuffix(sourcePath, ext)
	return fmt.Sprintf("%s_thumb_%dx%d%s", base, width, height, ext)
}

// getCompressedPath 获取压缩后路径
func (p *ImageProcessor) getCompressedPath(sourcePath string) string {
	ext := filepath.Ext(sourcePath)
	base := strings.TrimSuffix(sourcePath, ext)
	return fmt.Sprintf("%s_compressed%s", base, ext)
}

// getCroppedPath 获取裁剪后路径
func (p *ImageProcessor) getCroppedPath(sourcePath string) string {
	ext := filepath.Ext(sourcePath)
	base := strings.TrimSuffix(sourcePath, ext)
	return fmt.Sprintf("%s_cropped%s", base, ext)
}

// getConvertedPath 获取转换格式后路径
func (p *ImageProcessor) getConvertedPath(sourcePath string, outputFormat string) string {
	base := strings.TrimSuffix(sourcePath, filepath.Ext(sourcePath))
	return fmt.Sprintf("%s.%s", base, outputFormat)
}

// GetImageInfo 获取图片信息
func (p *ImageProcessor) GetImageInfo(ctx context.Context, sourcePath string) (*ImageInfo, error) {
	reader, err := p.backend.Load(ctx, sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load image: %w", err)
	}
	defer reader.Close()

	config, format, err := image.DecodeConfig(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image config: %w", err)
	}

	return &ImageInfo{
		Width:  config.Width,
		Height: config.Height,
		Format: format,
	}, nil
}

// ImageInfo 图片信息
type ImageInfo struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Format string `json:"format"`
}
