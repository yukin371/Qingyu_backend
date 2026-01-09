package shared

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockFile 模拟文件模型
type MockFile struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FileName     string             `bson:"file_name" json:"fileName"`
	OriginalName string             `bson:"original_name" json:"originalName"`
	FileSize     int64              `bson:"file_size" json:"fileSize"`
	ContentType  string             `bson:"content_type" json:"contentType"`
	FileHash     string             `bson:"file_hash" json:"fileHash"`
	StoragePath  string             `bson:"storage_path" json:"storagePath"`
	PublicURL    string             `bson:"public_url" json:"publicUrl"`
	Status       string             `bson:"status" json:"status"`
	UploadedBy   primitive.ObjectID `bson:"uploaded_by" json:"uploadedBy"`
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updatedAt"`
}

// MockFileVersion 模拟文件版本模型
type MockFileVersion struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FileID      primitive.ObjectID `bson:"file_id" json:"fileId"`
	Version     int                `bson:"version" json:"version"`
	FileSize    int64              `bson:"file_size" json:"fileSize"`
	FileHash    string             `bson:"file_hash" json:"fileHash"`
	StoragePath string             `bson:"storage_path" json:"storagePath"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
}

// MockFilePermission 模拟文件权限模型
type MockFilePermission struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FileID     primitive.ObjectID `bson:"file_id" json:"fileId"`
	UserID     primitive.ObjectID `bson:"user_id" json:"userId"`
	Permission string             `bson:"permission" json:"permission"`
	CreatedAt  time.Time          `bson:"created_at" json:"createdAt"`
}

// MockUploadRequest 模拟上传请求
type MockUploadRequest struct {
	FileName    string
	ContentType string
	FileSize    int64
	Content     io.Reader
	UserID      primitive.ObjectID
}

// MockFileRepository 模拟文件仓储
type MockFileRepository struct {
	mock.Mock
}

func (m *MockFileRepository) Create(ctx context.Context, file *MockFile) error {
	args := m.Called(ctx, file)
	return args.Error(0)
}

func (m *MockFileRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*MockFile, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockFile), args.Error(1)
}

func (m *MockFileRepository) GetByHash(ctx context.Context, hash string) (*MockFile, error) {
	args := m.Called(ctx, hash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockFile), args.Error(1)
}

func (m *MockFileRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockFileRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFileRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID, limit, offset int) ([]*MockFile, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*MockFile), args.Error(1)
}

// MockVersionRepository 模拟版本仓储
type MockVersionRepository struct {
	mock.Mock
}

func (m *MockVersionRepository) Create(ctx context.Context, version *MockFileVersion) error {
	args := m.Called(ctx, version)
	return args.Error(0)
}

func (m *MockVersionRepository) GetByFileID(ctx context.Context, fileID primitive.ObjectID) ([]*MockFileVersion, error) {
	args := m.Called(ctx, fileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*MockFileVersion), args.Error(1)
}

// MockPermissionRepository 模拟权限仓储
type MockPermissionRepository struct {
	mock.Mock
}

func (m *MockPermissionRepository) Create(ctx context.Context, permission *MockFilePermission) error {
	args := m.Called(ctx, permission)
	return args.Error(0)
}

func (m *MockPermissionRepository) CheckPermission(ctx context.Context, fileID, userID primitive.ObjectID, permission string) (bool, error) {
	args := m.Called(ctx, fileID, userID, permission)
	return args.Bool(0), args.Error(1)
}

// MockStorageProvider 模拟存储提供者
type MockStorageProvider struct {
	mock.Mock
}

func (m *MockStorageProvider) Upload(ctx context.Context, path string, content io.Reader, contentType string) (string, error) {
	args := m.Called(ctx, path, content, contentType)
	return args.String(0), args.Error(1)
}

func (m *MockStorageProvider) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	args := m.Called(ctx, path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockStorageProvider) Delete(ctx context.Context, path string) error {
	args := m.Called(ctx, path)
	return args.Error(0)
}

func (m *MockStorageProvider) GetPublicURL(path string) string {
	args := m.Called(path)
	return args.String(0)
}

// MockStorageService 模拟存储服务
type MockStorageService struct {
	fileRepo        *MockFileRepository
	versionRepo     *MockVersionRepository
	permissionRepo  *MockPermissionRepository
	storageProvider *MockStorageProvider
}

func NewMockStorageService(
	fileRepo *MockFileRepository,
	versionRepo *MockVersionRepository,
	permissionRepo *MockPermissionRepository,
	storageProvider *MockStorageProvider,
) *MockStorageService {
	return &MockStorageService{
		fileRepo:        fileRepo,
		versionRepo:     versionRepo,
		permissionRepo:  permissionRepo,
		storageProvider: storageProvider,
	}
}

// UploadFile 上传文件
func (s *MockStorageService) UploadFile(ctx context.Context, req *MockUploadRequest) (*MockFile, error) {
	if req.FileSize <= 0 {
		return nil, errors.New("invalid file size")
	}

	if req.FileName == "" {
		return nil, errors.New("file name is required")
	}

	// 计算文件哈希
	fileHash := calculateFileHash(req.Content)

	// 检查文件是否已存在
	existingFile, _ := s.fileRepo.GetByHash(ctx, fileHash)
	if existingFile != nil {
		return existingFile, nil
	}

	// 生成存储路径
	storagePath := generateStoragePath(req.FileName, fileHash)

	// 上传到存储提供者
	publicURL, err := s.storageProvider.Upload(ctx, storagePath, req.Content, req.ContentType)
	if err != nil {
		return nil, err
	}

	// 创建文件记录
	file := &MockFile{
		ID:           primitive.NewObjectID(),
		FileName:     generateFileName(req.FileName),
		OriginalName: req.FileName,
		FileSize:     req.FileSize,
		ContentType:  req.ContentType,
		FileHash:     fileHash,
		StoragePath:  storagePath,
		PublicURL:    publicURL,
		Status:       "active",
		UploadedBy:   req.UserID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.fileRepo.Create(ctx, file)
	if err != nil {
		return nil, err
	}

	// 创建版本记录
	version := &MockFileVersion{
		ID:          primitive.NewObjectID(),
		FileID:      file.ID,
		Version:     1,
		FileSize:    req.FileSize,
		FileHash:    fileHash,
		StoragePath: storagePath,
		CreatedAt:   time.Now(),
	}

	err = s.versionRepo.Create(ctx, version)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// GetFile 获取文件信息
func (s *MockStorageService) GetFile(ctx context.Context, fileID primitive.ObjectID, userID primitive.ObjectID) (*MockFile, error) {
	file, err := s.fileRepo.GetByID(ctx, fileID)
	if err != nil {
		return nil, err
	}

	// 检查权限
	hasPermission, err := s.CheckFilePermission(ctx, fileID, userID, "read")
	if err != nil {
		return nil, err
	}

	if !hasPermission && file.UploadedBy != userID {
		return nil, errors.New("permission denied")
	}

	return file, nil
}

// DownloadFile 下载文件
func (s *MockStorageService) DownloadFile(ctx context.Context, fileID primitive.ObjectID, userID primitive.ObjectID) (io.ReadCloser, error) {
	file, err := s.GetFile(ctx, fileID, userID)
	if err != nil {
		return nil, err
	}

	return s.storageProvider.Download(ctx, file.StoragePath)
}

// DeleteFile 删除文件
func (s *MockStorageService) DeleteFile(ctx context.Context, fileID primitive.ObjectID, userID primitive.ObjectID) error {
	file, err := s.fileRepo.GetByID(ctx, fileID)
	if err != nil {
		return err
	}

	// 检查权限
	hasPermission, err := s.CheckFilePermission(ctx, fileID, userID, "delete")
	if err != nil {
		return err
	}

	if !hasPermission && file.UploadedBy != userID {
		return errors.New("permission denied")
	}

	// 从存储提供者删除
	err = s.storageProvider.Delete(ctx, file.StoragePath)
	if err != nil {
		return err
	}

	// 删除文件记录
	return s.fileRepo.Delete(ctx, fileID)
}

// CheckFilePermission 检查文件权限
func (s *MockStorageService) CheckFilePermission(ctx context.Context, fileID, userID primitive.ObjectID, permission string) (bool, error) {
	return s.permissionRepo.CheckPermission(ctx, fileID, userID, permission)
}

// GrantFilePermission 授予文件权限
func (s *MockStorageService) GrantFilePermission(ctx context.Context, fileID, userID primitive.ObjectID, permission string) error {
	filePermission := &MockFilePermission{
		ID:         primitive.NewObjectID(),
		FileID:     fileID,
		UserID:     userID,
		Permission: permission,
		CreatedAt:  time.Now(),
	}

	return s.permissionRepo.Create(ctx, filePermission)
}

// GetUserFiles 获取用户文件列表
func (s *MockStorageService) GetUserFiles(ctx context.Context, userID primitive.ObjectID, limit, offset int) ([]*MockFile, error) {
	return s.fileRepo.GetByUserID(ctx, userID, limit, offset)
}

// GetFileVersions 获取文件版本历史
func (s *MockStorageService) GetFileVersions(ctx context.Context, fileID primitive.ObjectID, userID primitive.ObjectID) ([]*MockFileVersion, error) {
	// 检查权限
	hasPermission, err := s.CheckFilePermission(ctx, fileID, userID, "read")
	if err != nil {
		return nil, err
	}

	if !hasPermission {
		file, err := s.fileRepo.GetByID(ctx, fileID)
		if err != nil {
			return nil, err
		}
		if file.UploadedBy != userID {
			return nil, errors.New("permission denied")
		}
	}

	return s.versionRepo.GetByFileID(ctx, fileID)
}

// 辅助函数
func calculateFileHash(content io.Reader) string {
	// 模拟文件哈希计算
	return "hash_" + time.Now().Format("20060102150405")
}

func generateStoragePath(fileName, hash string) string {
	// 模拟存储路径生成
	return "files/" + hash + "/" + fileName
}

func generateFileName(originalName string) string {
	// 模拟文件名生成
	return time.Now().Format("20060102150405") + "_" + originalName
}

// 测试用例

func TestStorageService_UploadFile_Success(t *testing.T) {
	fileRepo := new(MockFileRepository)
	versionRepo := new(MockVersionRepository)
	permissionRepo := new(MockPermissionRepository)
	storageProvider := new(MockStorageProvider)
	service := NewMockStorageService(fileRepo, versionRepo, permissionRepo, storageProvider)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	content := strings.NewReader("test file content")

	req := &MockUploadRequest{
		FileName:    "test.txt",
		ContentType: "text/plain",
		FileSize:    17,
		Content:     content,
		UserID:      userID,
	}

	// Mock 设置
	fileRepo.On("GetByHash", ctx, mock.AnythingOfType("string")).Return(nil, errors.New("not found"))
	storageProvider.On("Upload", ctx, mock.AnythingOfType("string"), content, "text/plain").Return("https://cdn.example.com/test.txt", nil)
	fileRepo.On("Create", ctx, mock.AnythingOfType("*shared.MockFile")).Return(nil)
	versionRepo.On("Create", ctx, mock.AnythingOfType("*shared.MockFileVersion")).Return(nil)

	// 执行测试
	file, err := service.UploadFile(ctx, req)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, file)
	assert.Equal(t, "test.txt", file.OriginalName)
	assert.Equal(t, int64(17), file.FileSize)
	assert.Equal(t, "text/plain", file.ContentType)
	assert.Equal(t, "active", file.Status)
	assert.Equal(t, userID, file.UploadedBy)

	fileRepo.AssertExpectations(t)
	versionRepo.AssertExpectations(t)
	storageProvider.AssertExpectations(t)
}

func TestStorageService_UploadFile_DuplicateFile(t *testing.T) {
	fileRepo := new(MockFileRepository)
	versionRepo := new(MockVersionRepository)
	permissionRepo := new(MockPermissionRepository)
	storageProvider := new(MockStorageProvider)
	service := NewMockStorageService(fileRepo, versionRepo, permissionRepo, storageProvider)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	content := strings.NewReader("test file content")

	req := &MockUploadRequest{
		FileName:    "test.txt",
		ContentType: "text/plain",
		FileSize:    17,
		Content:     content,
		UserID:      userID,
	}

	existingFile := &MockFile{
		ID:           primitive.NewObjectID(),
		FileName:     "existing_test.txt",
		OriginalName: "test.txt",
		FileSize:     17,
		ContentType:  "text/plain",
		Status:       "active",
	}

	// Mock 设置
	fileRepo.On("GetByHash", ctx, mock.AnythingOfType("string")).Return(existingFile, nil)

	// 执行测试
	file, err := service.UploadFile(ctx, req)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, file)
	assert.Equal(t, existingFile, file)

	fileRepo.AssertExpectations(t)
}

func TestStorageService_UploadFile_InvalidFileSize(t *testing.T) {
	fileRepo := new(MockFileRepository)
	versionRepo := new(MockVersionRepository)
	permissionRepo := new(MockPermissionRepository)
	storageProvider := new(MockStorageProvider)
	service := NewMockStorageService(fileRepo, versionRepo, permissionRepo, storageProvider)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	content := strings.NewReader("")

	req := &MockUploadRequest{
		FileName:    "test.txt",
		ContentType: "text/plain",
		FileSize:    0,
		Content:     content,
		UserID:      userID,
	}

	// 执行测试
	file, err := service.UploadFile(ctx, req)

	// 断言
	assert.Error(t, err)
	assert.Nil(t, file)
	assert.Equal(t, "invalid file size", err.Error())
}

func TestStorageService_GetFile_Success(t *testing.T) {
	fileRepo := new(MockFileRepository)
	versionRepo := new(MockVersionRepository)
	permissionRepo := new(MockPermissionRepository)
	storageProvider := new(MockStorageProvider)
	service := NewMockStorageService(fileRepo, versionRepo, permissionRepo, storageProvider)

	ctx := context.Background()
	fileID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	expectedFile := &MockFile{
		ID:           fileID,
		FileName:     "test.txt",
		OriginalName: "test.txt",
		FileSize:     17,
		ContentType:  "text/plain",
		Status:       "active",
		UploadedBy:   userID,
	}

	// Mock 设置
	fileRepo.On("GetByID", ctx, fileID).Return(expectedFile, nil)
	permissionRepo.On("CheckPermission", ctx, fileID, userID, "read").Return(false, nil)

	// 执行测试
	file, err := service.GetFile(ctx, fileID, userID)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, file)
	assert.Equal(t, expectedFile, file)

	fileRepo.AssertExpectations(t)
	permissionRepo.AssertExpectations(t)
}

func TestStorageService_GetFile_PermissionDenied(t *testing.T) {
	fileRepo := new(MockFileRepository)
	versionRepo := new(MockVersionRepository)
	permissionRepo := new(MockPermissionRepository)
	storageProvider := new(MockStorageProvider)
	service := NewMockStorageService(fileRepo, versionRepo, permissionRepo, storageProvider)

	ctx := context.Background()
	fileID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	ownerID := primitive.NewObjectID()

	expectedFile := &MockFile{
		ID:           fileID,
		FileName:     "test.txt",
		OriginalName: "test.txt",
		FileSize:     17,
		ContentType:  "text/plain",
		Status:       "active",
		UploadedBy:   ownerID,
	}

	// Mock 设置
	fileRepo.On("GetByID", ctx, fileID).Return(expectedFile, nil)
	permissionRepo.On("CheckPermission", ctx, fileID, userID, "read").Return(false, nil)

	// 执行测试
	file, err := service.GetFile(ctx, fileID, userID)

	// 断言
	assert.Error(t, err)
	assert.Nil(t, file)
	assert.Equal(t, "permission denied", err.Error())

	fileRepo.AssertExpectations(t)
	permissionRepo.AssertExpectations(t)
}

func TestStorageService_DeleteFile_Success(t *testing.T) {
	fileRepo := new(MockFileRepository)
	versionRepo := new(MockVersionRepository)
	permissionRepo := new(MockPermissionRepository)
	storageProvider := new(MockStorageProvider)
	service := NewMockStorageService(fileRepo, versionRepo, permissionRepo, storageProvider)

	ctx := context.Background()
	fileID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	expectedFile := &MockFile{
		ID:          fileID,
		FileName:    "test.txt",
		StoragePath: "files/hash123/test.txt",
		UploadedBy:  userID,
	}

	// Mock 设置
	fileRepo.On("GetByID", ctx, fileID).Return(expectedFile, nil)
	permissionRepo.On("CheckPermission", ctx, fileID, userID, "delete").Return(false, nil)
	storageProvider.On("Delete", ctx, "files/hash123/test.txt").Return(nil)
	fileRepo.On("Delete", ctx, fileID).Return(nil)

	// 执行测试
	err := service.DeleteFile(ctx, fileID, userID)

	// 断言
	assert.NoError(t, err)

	fileRepo.AssertExpectations(t)
	permissionRepo.AssertExpectations(t)
	storageProvider.AssertExpectations(t)
}

func TestStorageService_GrantFilePermission_Success(t *testing.T) {
	fileRepo := new(MockFileRepository)
	versionRepo := new(MockVersionRepository)
	permissionRepo := new(MockPermissionRepository)
	storageProvider := new(MockStorageProvider)
	service := NewMockStorageService(fileRepo, versionRepo, permissionRepo, storageProvider)

	ctx := context.Background()
	fileID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	permission := "read"

	// Mock 设置
	permissionRepo.On("Create", ctx, mock.AnythingOfType("*shared.MockFilePermission")).Return(nil)

	// 执行测试
	err := service.GrantFilePermission(ctx, fileID, userID, permission)

	// 断言
	assert.NoError(t, err)

	permissionRepo.AssertExpectations(t)
}

func TestStorageService_GetUserFiles_Success(t *testing.T) {
	fileRepo := new(MockFileRepository)
	versionRepo := new(MockVersionRepository)
	permissionRepo := new(MockPermissionRepository)
	storageProvider := new(MockStorageProvider)
	service := NewMockStorageService(fileRepo, versionRepo, permissionRepo, storageProvider)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	limit := 10
	offset := 0

	expectedFiles := []*MockFile{
		{
			ID:           primitive.NewObjectID(),
			FileName:     "file1.txt",
			OriginalName: "file1.txt",
			UploadedBy:   userID,
		},
		{
			ID:           primitive.NewObjectID(),
			FileName:     "file2.txt",
			OriginalName: "file2.txt",
			UploadedBy:   userID,
		},
	}

	// Mock 设置
	fileRepo.On("GetByUserID", ctx, userID, limit, offset).Return(expectedFiles, nil)

	// 执行测试
	files, err := service.GetUserFiles(ctx, userID, limit, offset)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, files)
	assert.Len(t, files, 2)
	assert.Equal(t, expectedFiles, files)

	fileRepo.AssertExpectations(t)
}
