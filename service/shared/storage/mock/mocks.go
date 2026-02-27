package mock

import (
	storageInterfaces "Qingyu_backend/repository/interfaces/storage"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	storageModel "Qingyu_backend/models/storage"
)

// MockStorageBackend 模拟存储后端
type MockStorageBackend struct {
	mu          sync.RWMutex
	storage     map[string][]byte
	callLog     []string
	saveError   error
	loadError   error
	deleteError error
	existsError error
	getURLError error
}

// NewMockStorageBackend 创建模拟存储后端
func NewMockStorageBackend() *MockStorageBackend {
	return &MockStorageBackend{
		storage: make(map[string][]byte),
		callLog: make([]string, 0),
	}
}

// Save 保存文件
func (m *MockStorageBackend) Save(ctx context.Context, path string, reader io.Reader) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callLog = append(m.callLog, fmt.Sprintf("Save(%s)", path))

	if m.saveError != nil {
		return m.saveError
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read data: %w", err)
	}

	m.storage[path] = data
	return nil
}

// Load 加载文件
func (m *MockStorageBackend) Load(ctx context.Context, path string) (io.ReadCloser, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callLog = append(m.callLog, fmt.Sprintf("Load(%s)", path))

	if m.loadError != nil {
		return nil, m.loadError
	}

	data, exists := m.storage[path]
	if !exists {
		return nil, fmt.Errorf("file not found: %s", path)
	}

	return io.NopCloser(strings.NewReader(string(data))), nil
}

// Delete 删除文件
func (m *MockStorageBackend) Delete(ctx context.Context, path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callLog = append(m.callLog, fmt.Sprintf("Delete(%s)", path))

	if m.deleteError != nil {
		return m.deleteError
	}

	delete(m.storage, path)
	return nil
}

// Exists 检查文件是否存在
func (m *MockStorageBackend) Exists(ctx context.Context, path string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callLog = append(m.callLog, fmt.Sprintf("Exists(%s)", path))

	if m.existsError != nil {
		return false, m.existsError
	}

	_, exists := m.storage[path]
	return exists, nil
}

// GetURL 生成访问URL
func (m *MockStorageBackend) GetURL(ctx context.Context, path string, expiresIn time.Duration) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callLog = append(m.callLog, fmt.Sprintf("GetURL(%s, %v)", path, expiresIn))

	if m.getURLError != nil {
		return "", m.getURLError
	}

	return fmt.Sprintf("http://mock-storage.example.com/%s?expires=%v", path, expiresIn), nil
}

// ============ 辅助方法 ============

// GetData 获取存储的数据（用于测试验证）
func (m *MockStorageBackend) GetData(path string) ([]byte, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, exists := m.storage[path]
	return data, exists
}

// SetData 直接设置存储的数据（用于测试设置）
func (m *MockStorageBackend) SetData(path string, data []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.storage[path] = data
}

// GetCallLog 获取调用日志
func (m *MockStorageBackend) GetCallLog() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return append([]string{}, m.callLog...)
}

// ClearCallLog 清空调用日志
func (m *MockStorageBackend) ClearCallLog() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callLog = make([]string, 0)
}

// SetSaveError 设置保存错误
func (m *MockStorageBackend) SetSaveError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.saveError = err
}

// SetLoadError 设置加载错误
func (m *MockStorageBackend) SetLoadError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.loadError = err
}

// SetDeleteError 设置删除错误
func (m *MockStorageBackend) SetDeleteError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.deleteError = err
}

// SetExistsError 设置检查存在错误
func (m *MockStorageBackend) SetExistsError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.existsError = err
}

// SetGetURLError 设置URL生成错误
func (m *MockStorageBackend) SetGetURLError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.getURLError = err
}

// Reset 重置所有状态
func (m *MockStorageBackend) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.storage = make(map[string][]byte)
	m.callLog = make([]string, 0)
	m.saveError = nil
	m.loadError = nil
	m.deleteError = nil
	m.existsError = nil
	m.getURLError = nil
}

// ============ Mock StorageRepository ============

// MockStorageRepository 模拟存储仓储
type MockStorageRepository struct {
	mu               sync.RWMutex
	files            map[string]*storageModel.FileInfo
	multipartUploads map[string]*storageModel.MultipartUpload
	fileAccess       map[string][]string // fileID -> userIDs
	callLog          []string

	// 可配置的错误
	createFileError        error
	getFileError           error
	updateFileError        error
	deleteFileError        error
	createMultipartError   error
	getMultipartError      error
	updateMultipartError   error
	completeMultipartError error
	abortMultipartError    error
	listMultipartError     error
	grantAccessError       error
	revokeAccessError      error
	checkAccessError       error
}

// NewMockStorageRepository 创建模拟存储仓储
func NewMockStorageRepository() *MockStorageRepository {
	return &MockStorageRepository{
		files:            make(map[string]*storageModel.FileInfo),
		multipartUploads: make(map[string]*storageModel.MultipartUpload),
		fileAccess:       make(map[string][]string),
		callLog:          make([]string, 0),
	}
}

// ============ 文件元数据管理 ============

// CreateFile 创建文件元数据
func (m *MockStorageRepository) CreateFile(ctx context.Context, file *storageModel.FileInfo) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callLog = append(m.callLog, fmt.Sprintf("CreateFile(%s)", file.ID))

	if m.createFileError != nil {
		return m.createFileError
	}

	file.CreatedAt = time.Now()
	file.UpdatedAt = time.Now()
	m.files[file.ID] = file
	return nil
}

// GetFile 获取文件元数据
func (m *MockStorageRepository) GetFile(ctx context.Context, fileID string) (*storageModel.FileInfo, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callLog = append(m.callLog, fmt.Sprintf("GetFile(%s)", fileID))

	if m.getFileError != nil {
		return nil, m.getFileError
	}

	file, exists := m.files[fileID]
	if !exists {
		return nil, fmt.Errorf("file not found: %s", fileID)
	}

	// 返回副本避免修改原始数据
	fileCopy := *file
	return &fileCopy, nil
}

// GetFileByMD5 根据MD5获取文件
func (m *MockStorageRepository) GetFileByMD5(ctx context.Context, md5Hash string) (*storageModel.FileInfo, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callLog = append(m.callLog, fmt.Sprintf("GetFileByMD5(%s)", md5Hash))

	if m.getFileError != nil {
		return nil, m.getFileError
	}

	for _, file := range m.files {
		if file.MD5 == md5Hash {
			fileCopy := *file
			return &fileCopy, nil
		}
	}

	return nil, fmt.Errorf("file not found with MD5: %s", md5Hash)
}

// UpdateFile 更新文件元数据
func (m *MockStorageRepository) UpdateFile(ctx context.Context, fileID string, updates map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callLog = append(m.callLog, fmt.Sprintf("UpdateFile(%s, %v)", fileID, updates))

	if m.updateFileError != nil {
		return m.updateFileError
	}

	file, exists := m.files[fileID]
	if !exists {
		return fmt.Errorf("file not found: %s", fileID)
	}

	// 应用更新
	if contentType, ok := updates["content_type"].(string); ok {
		file.ContentType = contentType
	}
	if size, ok := updates["size"].(int64); ok {
		file.Size = size
	}
	if isPublic, ok := updates["is_public"].(bool); ok {
		file.IsPublic = isPublic
	}
	if downloads, ok := updates["downloads"].(int64); ok {
		file.Downloads = downloads
	}
	file.UpdatedAt = time.Now()

	return nil
}

// DeleteFile 删除文件元数据
func (m *MockStorageRepository) DeleteFile(ctx context.Context, fileID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callLog = append(m.callLog, fmt.Sprintf("DeleteFile(%s)", fileID))

	if m.deleteFileError != nil {
		return m.deleteFileError
	}

	if _, exists := m.files[fileID]; !exists {
		return fmt.Errorf("file not found: %s", fileID)
	}

	delete(m.files, fileID)
	delete(m.fileAccess, fileID)
	return nil
}

// ListFiles 列出文件（简化实现，忽略 filter）
func (m *MockStorageRepository) ListFiles(ctx context.Context, filter *storageInterfaces.FileFilter) ([]*storageModel.FileInfo, int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callLog = append(m.callLog, fmt.Sprintf("ListFiles(%v)", filter))

	files := make([]*storageModel.FileInfo, 0, len(m.files))
	for _, file := range m.files {
		fileCopy := *file
		files = append(files, &fileCopy)
	}

	return files, int64(len(files)), nil
}

// CountFiles 统计文件数量
func (m *MockStorageRepository) CountFiles(ctx context.Context, filter *storageInterfaces.FileFilter) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callLog = append(m.callLog, fmt.Sprintf("CountFiles(%v)", filter))

	return int64(len(m.files)), nil
}

// ============ 分片上传管理 ============

// CreateMultipartUpload 创建分片上传任务
func (m *MockStorageRepository) CreateMultipartUpload(ctx context.Context, upload *storageModel.MultipartUpload) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callLog = append(m.callLog, fmt.Sprintf("CreateMultipartUpload(%s)", upload.UploadID))

	if m.createMultipartError != nil {
		return m.createMultipartError
	}

	upload.CreatedAt = time.Now()
	upload.UpdatedAt = time.Now()
	m.multipartUploads[upload.UploadID] = upload
	return nil
}

// GetMultipartUpload 获取分片上传任务
func (m *MockStorageRepository) GetMultipartUpload(ctx context.Context, uploadID string) (*storageModel.MultipartUpload, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callLog = append(m.callLog, fmt.Sprintf("GetMultipartUpload(%s)", uploadID))

	if m.getMultipartError != nil {
		return nil, m.getMultipartError
	}

	upload, exists := m.multipartUploads[uploadID]
	if !exists {
		return nil, fmt.Errorf("multipart upload not found: %s", uploadID)
	}

	// 返回副本
	uploadCopy := *upload
	uploadCopy.UploadedChunks = append([]int{}, upload.UploadedChunks...)
	return &uploadCopy, nil
}

// UpdateMultipartUpload 更新分片上传任务
func (m *MockStorageRepository) UpdateMultipartUpload(ctx context.Context, uploadID string, updates map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callLog = append(m.callLog, fmt.Sprintf("UpdateMultipartUpload(%s, %v)", uploadID, updates))

	if m.updateMultipartError != nil {
		return m.updateMultipartError
	}

	upload, exists := m.multipartUploads[uploadID]
	if !exists {
		return fmt.Errorf("multipart upload not found: %s", uploadID)
	}

	// 应用更新
	if uploadedChunks, ok := updates["uploaded_chunks"].([]int); ok {
		upload.UploadedChunks = uploadedChunks
	}
	if status, ok := updates["status"].(string); ok {
		upload.Status = status
	}
	upload.UpdatedAt = time.Now()

	return nil
}

// CompleteMultipartUpload 完成分片上传
func (m *MockStorageRepository) CompleteMultipartUpload(ctx context.Context, uploadID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callLog = append(m.callLog, fmt.Sprintf("CompleteMultipartUpload(%s)", uploadID))

	if m.completeMultipartError != nil {
		return m.completeMultipartError
	}

	upload, exists := m.multipartUploads[uploadID]
	if !exists {
		return fmt.Errorf("multipart upload not found: %s", uploadID)
	}

	upload.Status = "completed"
	now := time.Now()
	upload.UpdatedAt = now
	upload.CompletedAt = &now

	return nil
}

// AbortMultipartUpload 中止分片上传
func (m *MockStorageRepository) AbortMultipartUpload(ctx context.Context, uploadID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callLog = append(m.callLog, fmt.Sprintf("AbortMultipartUpload(%s)", uploadID))

	if m.abortMultipartError != nil {
		return m.abortMultipartError
	}

	upload, exists := m.multipartUploads[uploadID]
	if !exists {
		return fmt.Errorf("multipart upload not found: %s", uploadID)
	}

	upload.Status = "aborted"
	upload.UpdatedAt = time.Now()

	return nil
}

// ListMultipartUploads 列出分片上传任务
func (m *MockStorageRepository) ListMultipartUploads(ctx context.Context, userID string, status string) ([]*storageModel.MultipartUpload, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callLog = append(m.callLog, fmt.Sprintf("ListMultipartUploads(%s, %s)", userID, status))

	if m.listMultipartError != nil {
		return nil, m.listMultipartError
	}

	result := make([]*storageModel.MultipartUpload, 0)
	for _, upload := range m.multipartUploads {
		if upload.UploadedBy == userID {
			if status == "" || upload.Status == status {
				uploadCopy := *upload
				uploadCopy.UploadedChunks = append([]int{}, upload.UploadedChunks...)
				result = append(result, &uploadCopy)
			}
		}
	}

	return result, nil
}

// ============ 权限管理 ============

// GrantAccess 授予访问权限
func (m *MockStorageRepository) GrantAccess(ctx context.Context, fileID, userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callLog = append(m.callLog, fmt.Sprintf("GrantAccess(%s, %s)", fileID, userID))

	if m.grantAccessError != nil {
		return m.grantAccessError
	}

	if _, exists := m.files[fileID]; !exists {
		return fmt.Errorf("file not found: %s", fileID)
	}

	m.fileAccess[fileID] = append(m.fileAccess[fileID], userID)
	return nil
}

// RevokeAccess 撤销访问权限
func (m *MockStorageRepository) RevokeAccess(ctx context.Context, fileID, userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callLog = append(m.callLog, fmt.Sprintf("RevokeAccess(%s, %s)", fileID, userID))

	if m.revokeAccessError != nil {
		return m.revokeAccessError
	}

	accessList, exists := m.fileAccess[fileID]
	if !exists {
		return fmt.Errorf("file not found: %s", fileID)
	}

	for i, uid := range accessList {
		if uid == userID {
			m.fileAccess[fileID] = append(accessList[:i], accessList[i+1:]...)
			return nil
		}
	}

	return nil
}

// CheckAccess 检查访问权限
func (m *MockStorageRepository) CheckAccess(ctx context.Context, fileID, userID string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callLog = append(m.callLog, fmt.Sprintf("CheckAccess(%s, %s)", fileID, userID))

	if m.checkAccessError != nil {
		return false, m.checkAccessError
	}

	accessList, exists := m.fileAccess[fileID]
	if !exists {
		return false, nil
	}

	for _, uid := range accessList {
		if uid == userID {
			return true, nil
		}
	}

	return false, nil
}

// ============ 统计功能 ============

// IncrementDownloadCount 增加下载次数
func (m *MockStorageRepository) IncrementDownloadCount(ctx context.Context, fileID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callLog = append(m.callLog, fmt.Sprintf("IncrementDownloadCount(%s)", fileID))

	file, exists := m.files[fileID]
	if !exists {
		return fmt.Errorf("file not found: %s", fileID)
	}

	file.Downloads++
	return nil
}

// ============ 健康检查 ============

// Health 健康检查
func (m *MockStorageRepository) Health(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callLog = append(m.callLog, "Health()")
	return nil
}

// ============ 辅助方法 ============

// GetCallLog 获取调用日志
func (m *MockStorageRepository) GetCallLog() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return append([]string{}, m.callLog...)
}

// ClearCallLog 清空调用日志
func (m *MockStorageRepository) ClearCallLog() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callLog = make([]string, 0)
}

// ============ 错误设置方法 ============

func (m *MockStorageRepository) SetCreateFileError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.createFileError = err
}

func (m *MockStorageRepository) SetGetFileError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.getFileError = err
}

func (m *MockStorageRepository) SetUpdateFileError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.updateFileError = err
}

func (m *MockStorageRepository) SetDeleteFileError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.deleteFileError = err
}

func (m *MockStorageRepository) SetCreateMultipartError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.createMultipartError = err
}

func (m *MockStorageRepository) SetGetMultipartError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.getMultipartError = err
}

func (m *MockStorageRepository) SetUpdateMultipartError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.updateMultipartError = err
}

func (m *MockStorageRepository) SetCompleteMultipartError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.completeMultipartError = err
}

func (m *MockStorageRepository) SetAbortMultipartError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.abortMultipartError = err
}

func (m *MockStorageRepository) SetListMultipartError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.listMultipartError = err
}

func (m *MockStorageRepository) SetGrantAccessError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.grantAccessError = err
}

func (m *MockStorageRepository) SetRevokeAccessError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.revokeAccessError = err
}

func (m *MockStorageRepository) SetCheckAccessError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.checkAccessError = err
}

// SetMultipartUploadExpiresAt 设置分片上传任务过期时间（用于测试）
func (m *MockStorageRepository) SetMultipartUploadExpiresAt(uploadID string, expiresAt time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	upload, exists := m.multipartUploads[uploadID]
	if !exists {
		return fmt.Errorf("multipart upload not found: %s", uploadID)
	}

	upload.ExpiresAt = expiresAt
	upload.UpdatedAt = time.Now()
	return nil
}

// Reset 重置所有状态
func (m *MockStorageRepository) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.files = make(map[string]*storageModel.FileInfo)
	m.multipartUploads = make(map[string]*storageModel.MultipartUpload)
	m.fileAccess = make(map[string][]string)
	m.callLog = make([]string, 0)

	m.createFileError = nil
	m.getFileError = nil
	m.updateFileError = nil
	m.deleteFileError = nil
	m.createMultipartError = nil
	m.getMultipartError = nil
	m.updateMultipartError = nil
	m.completeMultipartError = nil
	m.abortMultipartError = nil
	m.listMultipartError = nil
	m.grantAccessError = nil
	m.revokeAccessError = nil
	m.checkAccessError = nil
}
