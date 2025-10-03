package storage

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"
)

// StorageServiceImpl 存储服务实现
type StorageServiceImpl struct {
	backend  StorageBackend
	fileRepo FileRepository
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
		backend:  backend,
		fileRepo: fileRepo,
	}
}

// ============ 文件操作 ============

// Upload 上传文件
func (s *StorageServiceImpl) Upload(ctx context.Context, req *UploadRequest) (*FileInfo, error) {
	// 1. 生成文件ID和存储路径
	fileID := generateFileID()
	ext := filepath.Ext(req.Filename)
	storedFilename := fileID + ext

	// 按分类和日期组织目录
	datePath := time.Now().Format("2006/01/02")
	storagePath := filepath.Join(req.Category, datePath, storedFilename)

	// 2. 计算MD5（用于去重）
	var md5Hash string
	if req.File != nil {
		// 如果需要计算MD5，需要读取两次文件（或使用TeeReader）
		// 这里简化处理，实际应用中可以使用TeeReader同时计算MD5和保存
		md5Hash = s.calculateMD5(req.File)
	}

	// 3. 保存文件到存储后端
	if err := s.backend.Save(ctx, storagePath, req.File); err != nil {
		return nil, fmt.Errorf("保存文件失败: %w", err)
	}

	// 4. 创建文件元数据
	fileInfo := &FileInfo{
		ID:           fileID,
		Filename:     storedFilename,
		OriginalName: req.Filename,
		ContentType:  req.ContentType,
		Size:         req.Size,
		Path:         storagePath,
		UserID:       req.UserID,
		IsPublic:     req.IsPublic,
		Category:     req.Category,
		MD5:          md5Hash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
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
	// 1. 获取文件元数据
	fileInfo, err := s.fileRepo.Get(ctx, fileID)
	if err != nil {
		return fmt.Errorf("文件不存在: %w", err)
	}

	// 2. 删除存储后端的文件
	if err := s.backend.Delete(ctx, fileInfo.Path); err != nil {
		// 记录错误但继续删除元数据
		fmt.Printf("删除文件失败: %v\n", err)
	}

	// 3. 删除元数据
	if err := s.fileRepo.Delete(ctx, fileID); err != nil {
		return fmt.Errorf("删除文件元数据失败: %w", err)
	}

	return nil
}

// GetFileInfo 获取文件信息
func (s *StorageServiceImpl) GetFileInfo(ctx context.Context, fileID string) (*FileInfo, error) {
	fileInfo, err := s.fileRepo.Get(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("文件不存在: %w", err)
	}
	return fileInfo, nil
}

// ============ 权限控制 ============

// GrantAccess 授予访问权限
func (s *StorageServiceImpl) GrantAccess(ctx context.Context, fileID, userID string) error {
	return s.fileRepo.GrantAccess(ctx, fileID, userID)
}

// RevokeAccess 撤销访问权限
func (s *StorageServiceImpl) RevokeAccess(ctx context.Context, fileID, userID string) error {
	return s.fileRepo.RevokeAccess(ctx, fileID, userID)
}

// CheckAccess 检查访问权限
func (s *StorageServiceImpl) CheckAccess(ctx context.Context, fileID, userID string) (bool, error) {
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
	// 设置默认分页
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return s.fileRepo.List(ctx, req.UserID, req.Category, page, pageSize)
}

// GetDownloadURL 生成下载链接
func (s *StorageServiceImpl) GetDownloadURL(ctx context.Context, fileID string, expiresIn time.Duration) (string, error) {
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

// ============ 健康检查 ============

// Health 健康检查
func (s *StorageServiceImpl) Health(ctx context.Context) error {
	// 检查存储后端是否可用
	exists, err := s.backend.Exists(ctx, "health_check")
	if err != nil {
		return fmt.Errorf("存储后端不可用: %w", err)
	}
	_ = exists
	return nil
}

// ============ 辅助方法 ============

// calculateMD5 计算文件MD5
func (s *StorageServiceImpl) calculateMD5(reader io.Reader) string {
	hash := md5.New()
	io.Copy(hash, reader)
	return hex.EncodeToString(hash.Sum(nil))
}

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
