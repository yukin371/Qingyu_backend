package storage

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultListPage     = 1
	defaultListPageSize = 20
	maxListPageSize     = 100
	defaultDownloadTTL  = 15 * time.Minute
)

// TODO: 完善文件上传功能（分片上传、断点续传）
// TODO: 完善文件下载功能（断点续传、流式下载）
// TODO: 添加图片处理功能（压缩、裁剪、水印）
// TODO: 集成云存储服务（阿里云OSS、腾讯云COS、AWS S3）
// TODO: 实现文件版本管理

// StorageServiceImpl 存储服务实现
type StorageServiceImpl struct {
	backend     StorageBackend
	fileRepo    FileRepository
	initialized bool // 初始化标志
}

// StorageBackend 存储后端接口
type StorageBackend interface {
	// 文件操作
	Save(ctx context.Context, path string, reader io.Reader) error
	Load(ctx context.Context, path string) (io.ReadCloser, error)
	Delete(ctx context.Context, path string) error
	Exists(ctx context.Context, path string) (bool, error)
	// 生成访问URL
	GetURL(ctx context.Context, path string, expiresIn time.Duration) (string, error)
}

// FileRepository 文件元数据Repository
type FileRepository interface {
	Create(ctx context.Context, file *FileInfo) error
	Get(ctx context.Context, fileID string) (*FileInfo, error)
	Update(ctx context.Context, fileID string, updates map[string]interface{}) error
	Delete(ctx context.Context, fileID string) error
	List(ctx context.Context, userID, category string, page, pageSize int) ([]*FileInfo, error)
	// 访问控制
	GrantAccess(ctx context.Context, fileID, userID string) error
	RevokeAccess(ctx context.Context, fileID, userID string) error
	CheckAccess(ctx context.Context, fileID, userID string) (bool, error)
}

// NewStorageService 创建存储服务
func NewStorageService(backend StorageBackend, fileRepo FileRepository) StorageService {
	return &StorageServiceImpl{
		backend:     backend,
		fileRepo:    fileRepo,
		initialized: true, // 简单实现直接标记为已初始化
	}
}

// ============ 文件操作 ============

// Upload 上传文件
func (s *StorageServiceImpl) Upload(ctx context.Context, req *UploadRequest) (*FileInfo, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, fmt.Errorf("upload request is required")
	}
	if req.File == nil {
		return nil, fmt.Errorf("file is required")
	}
	if strings.TrimSpace(req.Filename) == "" {
		return nil, fmt.Errorf("filename is required")
	}
	if strings.TrimSpace(req.UserID) == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	fileData, err := io.ReadAll(req.File)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	// 1. 生成文件ID和存储路径
	fileID := generateFileID()
	ext := filepath.Ext(req.Filename)
	storedFilename := fileID + ext

	// 按分类和日期组织目录
	datePath := time.Now().Format("2006/01/02")
	category := normalizeStorageCategory(req.Category)
	storagePath := filepath.Join(category, datePath, storedFilename)

	// 2. 计算MD5（用于去重）
	md5Hash := calculateBytesMD5(fileData)
	size := req.Size
	if size <= 0 {
		size = int64(len(fileData))
	}

	// 3. 保存文件到存储后端
	if err := s.backend.Save(ctx, storagePath, bytes.NewReader(fileData)); err != nil {
		return nil, fmt.Errorf("保存文件失败: %w", err)
	}

	// 4. 创建文件元数据
	now := time.Now()
	fileInfo := &FileInfo{
		ID:           fileID,
		Filename:     storedFilename,
		OriginalName: req.Filename,
		ContentType:  req.ContentType,
		Size:         size,
		Path:         storagePath,
		UserID:       req.UserID,
		IsPublic:     req.IsPublic,
		Category:     category,
		MD5:          md5Hash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// 5. 保存元数据到数据库
	if err := s.fileRepo.Create(ctx, fileInfo); err != nil {
		// 回滚：删除已保存的文件
		s.backend.Delete(ctx, storagePath)
		return nil, fmt.Errorf("保存文件元数据失败: %w", err)
	}

	return fileInfo, nil
}

// Download 下载文件
func (s *StorageServiceImpl) Download(ctx context.Context, fileID string) (io.ReadCloser, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if strings.TrimSpace(fileID) == "" {
		return nil, fmt.Errorf("fileID is required")
	}

	// 1. 获取文件元数据
	fileInfo, err := s.fileRepo.Get(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("文件不存在: %w", err)
	}

	// 2. 从存储后端加载文件
	reader, err := s.backend.Load(ctx, fileInfo.Path)
	if err != nil {
		return nil, fmt.Errorf("加载文件失败: %w", err)
	}

	// 3. 更新最后访问时间（可选）
	go func() {
		s.fileRepo.Update(context.Background(), fileID, map[string]interface{}{
			"updated_at": time.Now(),
		})
	}()

	return reader, nil
}

// Delete 删除文件
func (s *StorageServiceImpl) Delete(ctx context.Context, fileID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if strings.TrimSpace(fileID) == "" {
		return fmt.Errorf("fileID is required")
	}

	// 1. 获取文件元数据
	fileInfo, err := s.fileRepo.Get(ctx, fileID)
	if err != nil {
		return fmt.Errorf("文件不存在: %w", err)
	}

	// 2. 删除存储后端的文件
	backendErr := s.backend.Delete(ctx, fileInfo.Path)

	// 3. 删除元数据
	if err := s.fileRepo.Delete(ctx, fileID); err != nil {
		if backendErr != nil {
			return fmt.Errorf("删除文件元数据失败: %w", errors.Join(err, backendErr))
		}
		return fmt.Errorf("删除文件元数据失败: %w", err)
	}
	if backendErr != nil {
		return fmt.Errorf("删除文件失败: %w", backendErr)
	}

	return nil
}

// GetFileInfo 获取文件信息
func (s *StorageServiceImpl) GetFileInfo(ctx context.Context, fileID string) (*FileInfo, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if strings.TrimSpace(fileID) == "" {
		return nil, fmt.Errorf("fileID is required")
	}

	fileInfo, err := s.fileRepo.Get(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("文件不存在: %w", err)
	}
	return fileInfo, nil
}

// ============ 权限控制 ============

// GrantAccess 授予访问权限
func (s *StorageServiceImpl) GrantAccess(ctx context.Context, fileID, userID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if strings.TrimSpace(fileID) == "" || strings.TrimSpace(userID) == "" {
		return fmt.Errorf("fileID and userID are required")
	}
	return s.fileRepo.GrantAccess(ctx, fileID, userID)
}

// RevokeAccess 撤销访问权限
func (s *StorageServiceImpl) RevokeAccess(ctx context.Context, fileID, userID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if strings.TrimSpace(fileID) == "" || strings.TrimSpace(userID) == "" {
		return fmt.Errorf("fileID and userID are required")
	}
	return s.fileRepo.RevokeAccess(ctx, fileID, userID)
}

// CheckAccess 检查访问权限
func (s *StorageServiceImpl) CheckAccess(ctx context.Context, fileID, userID string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	if strings.TrimSpace(fileID) == "" || strings.TrimSpace(userID) == "" {
		return false, fmt.Errorf("fileID and userID are required")
	}

	// 1. 获取文件信息
	fileInfo, err := s.fileRepo.Get(ctx, fileID)
	if err != nil {
		return false, err
	}

	// 2. 如果是公开文件，允许访问
	if fileInfo.IsPublic {
		return true, nil
	}

	// 3. 如果是文件所有者，允许访问
	if fileInfo.UserID == userID {
		return true, nil
	}

	// 4. 检查显式授权
	return s.fileRepo.CheckAccess(ctx, fileID, userID)
}

// ============ 文件管理 ============

// ListFiles 查询文件列表
func (s *StorageServiceImpl) ListFiles(ctx context.Context, req *ListFilesRequest) ([]*FileInfo, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, fmt.Errorf("list request is required")
	}

	// 设置默认分页
	page := req.Page
	if page <= 0 {
		page = defaultListPage
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = defaultListPageSize
	}
	if pageSize > maxListPageSize {
		pageSize = maxListPageSize
	}

	return s.fileRepo.List(ctx, req.UserID, normalizeStorageCategory(req.Category), page, pageSize)
}

// GetDownloadURL 生成下载链接
func (s *StorageServiceImpl) GetDownloadURL(ctx context.Context, fileID string, expiresIn time.Duration) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if strings.TrimSpace(fileID) == "" {
		return "", fmt.Errorf("fileID is required")
	}
	if expiresIn <= 0 {
		expiresIn = defaultDownloadTTL
	}

	// 1. 获取文件信息
	fileInfo, err := s.fileRepo.Get(ctx, fileID)
	if err != nil {
		return "", fmt.Errorf("文件不存在: %w", err)
	}

	// 2. 生成访问URL
	url, err := s.backend.GetURL(ctx, fileInfo.Path, expiresIn)
	if err != nil {
		return "", fmt.Errorf("生成下载链接失败: %w", err)
	}

	return url, nil
}

// ============ 辅助方法 ============

// getFileExtension 获取文件扩展名
func getFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	return strings.ToLower(ext)
}

// isImageFile 判断是否为图片文件
func isImageFile(contentType string) bool {
	return strings.HasPrefix(contentType, "image/")
}

// isVideoFile 判断是否为视频文件
func isVideoFile(contentType string) bool {
	return strings.HasPrefix(contentType, "video/")
}

// isAudioFile 判断是否为音频文件
func isAudioFile(contentType string) bool {
	return strings.HasPrefix(contentType, "audio/")
}

// generateFileID 生成文件ID
func generateFileID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func normalizeStorageCategory(category string) string {
	if strings.TrimSpace(category) == "" {
		return defaultCategory
	}
	return category
}

// ============ BaseService 接口实现 ============

// Initialize 初始化服务
func (s *StorageServiceImpl) Initialize(ctx context.Context) error {
	if s.initialized {
		return nil
	}

	// 验证依赖项
	if s.backend == nil {
		return fmt.Errorf("backend is nil")
	}
	if s.fileRepo == nil {
		return fmt.Errorf("fileRepo is nil")
	}

	// 初始化完成
	s.initialized = true
	return nil
}

// Health 健康检查
func (s *StorageServiceImpl) Health(ctx context.Context) error {
	if !s.initialized {
		return fmt.Errorf("service not initialized")
	}

	// 检查存储后端是否可用
	exists, err := s.backend.Exists(ctx, "health_check")
	if err != nil {
		return fmt.Errorf("存储后端不可用: %w", err)
	}
	_ = exists
	return nil
}

// Close 关闭服务，清理资源
func (s *StorageServiceImpl) Close(ctx context.Context) error {
	// 存储服务暂无需要清理的资源
	s.initialized = false
	return nil
}

// GetServiceName 获取服务名称
func (s *StorageServiceImpl) GetServiceName() string {
	return "StorageService"
}

// GetVersion 获取服务版本
func (s *StorageServiceImpl) GetVersion() string {
	return "v1.0.0"
}
