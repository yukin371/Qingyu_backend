package middleware

import (
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strings"
)

// UploadConfig 文件上传配置
type UploadConfig struct {
	MaxSize      int64    // 最大文件大小（字节）
	AllowedTypes []string // 允许的MIME类型
	AllowedExts  []string // 允许的文件扩展名
}

// DefaultUploadConfig 默认上传配置
func DefaultUploadConfig() *UploadConfig {
	return &UploadConfig{
		MaxSize:      10 * 1024 * 1024, // 10MB
		AllowedTypes: []string{"image/jpeg", "image/png", "image/gif", "application/pdf"},
		AllowedExts:  []string{".jpg", ".jpeg", ".png", ".gif", ".pdf"},
	}
}

// ValidateUpload 验证上传的文件
func ValidateUpload(fileHeader *multipart.FileHeader, config *UploadConfig) error {
	if config == nil {
		config = DefaultUploadConfig()
	}

	// 验证文件大小
	if err := ValidateFileSize(fileHeader.Size, config.MaxSize); err != nil {
		return err
	}

	// 验证文件扩展名
	if !ValidateFileExtension(fileHeader.Filename, config.AllowedExts) {
		return fmt.Errorf("file extension not allowed: %s", filepath.Ext(fileHeader.Filename))
	}

	// 验证文件类型（通过读取文件头）
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// 读取前512字节用于检测文件类型
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && n == 0 {
		return fmt.Errorf("failed to read file header: %w", err)
	}

	detectedType, err := DetectFileType(buffer[:n])
	if err != nil {
		return fmt.Errorf("failed to detect file type: %w", err)
	}

	// 验证检测到的类型是否在允许列表中
	if !isAllowedType(detectedType, config.AllowedTypes) {
		return fmt.Errorf("file type not allowed: %s (detected: %s)", filepath.Ext(fileHeader.Filename), detectedType)
	}

	return nil
}

// ValidateFileSize 验证文件大小
func ValidateFileSize(size int64, maxSize int64) error {
	if size > maxSize {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size %d bytes", size, maxSize)
	}
	return nil
}

// ValidateFileExtension 验证文件扩展名
func ValidateFileExtension(filename string, allowedExts []string) bool {
	if filename == "" || len(allowedExts) == 0 {
		return false
	}

	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedExt := range allowedExts {
		if ext == strings.ToLower(allowedExt) {
			return true
		}
	}
	return false
}

// DetectFileType 通过文件头检测文件类型
func DetectFileType(content []byte) (string, error) {
	if len(content) < 4 {
		return "", errors.New("file too small to detect type")
	}

	// JPEG
	if len(content) >= 2 && content[0] == 0xFF && content[1] == 0xD8 {
		return "image/jpeg", nil
	}

	// PNG
	if len(content) >= 8 &&
		content[0] == 0x89 && content[1] == 0x50 &&
		content[2] == 0x4E && content[3] == 0x47 &&
		content[4] == 0x0D && content[5] == 0x0A &&
		content[6] == 0x1A && content[7] == 0x0A {
		return "image/png", nil
	}

	// GIF
	if len(content) >= 6 &&
		content[0] == 0x47 && content[1] == 0x49 &&
		content[2] == 0x46 && content[3] == 0x38 &&
		(content[4] == 0x37 || content[4] == 0x39) &&
		content[5] == 0x61 {
		return "image/gif", nil
	}

	// PDF
	if len(content) >= 4 &&
		content[0] == 0x25 && content[1] == 0x50 &&
		content[2] == 0x44 && content[3] == 0x46 {
		return "application/pdf", nil
	}

	// WebP
	if len(content) >= 12 &&
		content[0] == 0x52 && content[1] == 0x49 &&
		content[2] == 0x46 && content[3] == 0x46 &&
		content[8] == 0x57 && content[9] == 0x45 &&
		content[10] == 0x42 && content[11] == 0x50 {
		return "image/webp", nil
	}

	// BMP
	if len(content) >= 2 &&
		content[0] == 0x42 && content[1] == 0x4D {
		return "image/bmp", nil
	}

	// TIFF (little-endian)
	if len(content) >= 4 &&
		content[0] == 0x49 && content[1] == 0x49 &&
		content[2] == 0x2A && content[3] == 0x00 {
		return "image/tiff", nil
	}

	// TIFF (big-endian)
	if len(content) >= 4 &&
		content[0] == 0x4D && content[1] == 0x4D &&
		content[2] == 0x00 && content[3] == 0x2A {
		return "image/tiff", nil
	}

	// ZIP (could be docx, xlsx, etc.)
	if len(content) >= 4 &&
		content[0] == 0x50 && content[1] == 0x4B &&
		(content[2] == 0x03 || content[2] == 0x05 ||
			content[2] == 0x07 || content[2] == 0x00) &&
		(content[3] == 0x04 || content[3] == 0x06 ||
			content[3] == 0x08 || content[3] == 0x00) {
		return "application/zip", nil
	}

	return "", errors.New("unknown file type")
}

// isAllowedType 检查文件类型是否在允许列表中
func isAllowedType(fileType string, allowedTypes []string) bool {
	for _, allowedType := range allowedTypes {
		if strings.EqualFold(fileType, allowedType) {
			return true
		}
	}
	return false
}

// GenerateSafeFilename 生成安全的文件名
func GenerateSafeFilename(filename string) string {
	if filename == "" {
		return ""
	}

	// 获取文件名（不含路径）
	filename = filepath.Base(filename)

	// 获取扩展名
	ext := filepath.Ext(filename)
	base := filename[:len(filename)-len(ext)]

	// 移除路径遍历字符
	base = strings.ReplaceAll(base, "..", "")
	base = strings.ReplaceAll(base, "/", "")
	base = strings.ReplaceAll(base, "\\", "")

	// 替换空格为下划线（在删除特殊字符之前）
	base = strings.ReplaceAll(base, " ", "_")

	// 只保留字母、数字、下划线、连字符和点
	safeRegex := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	base = safeRegex.ReplaceAllString(base, "")

	// 移除连续的特殊字符
	specialCharRegex := regexp.MustCompile(`[._-]{2,}`)
	base = specialCharRegex.ReplaceAllString(base, ".")

	// 去除首尾的特殊字符
	base = strings.Trim(base, "._-")

	// 如果base为空，只保留扩展名
	if base == "" {
		return ext
	}

	return base + ext
}

// ValidateImageUpload 验证图片上传（便捷方法）
func ValidateImageUpload(fileHeader *multipart.FileHeader, maxSize int64) error {
	config := &UploadConfig{
		MaxSize: maxSize,
		AllowedTypes: []string{
			"image/jpeg",
			"image/png",
			"image/gif",
			"image/webp",
		},
		AllowedExts: []string{".jpg", ".jpeg", ".png", ".gif", ".webp"},
	}
	return ValidateUpload(fileHeader, config)
}

// ValidateDocumentUpload 验证文档上传（便捷方法）
func ValidateDocumentUpload(fileHeader *multipart.FileHeader, maxSize int64) error {
	config := &UploadConfig{
		MaxSize: maxSize,
		AllowedTypes: []string{
			"application/pdf",
			"application/msword",
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		},
		AllowedExts: []string{".pdf", ".doc", ".docx"},
	}
	return ValidateUpload(fileHeader, config)
}
