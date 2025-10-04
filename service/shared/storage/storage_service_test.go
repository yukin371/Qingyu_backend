package storage

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"
)

// ============ Mock StorageBackend ============

type MockStorageBackend struct {
	saveFunc   func(ctx context.Context, path string, reader io.Reader) error
	loadFunc   func(ctx context.Context, path string) (io.ReadCloser, error)
	deleteFunc func(ctx context.Context, path string) error
	existsFunc func(ctx context.Context, path string) (bool, error)
	getURLFunc func(ctx context.Context, path string, expiresIn time.Duration) (string, error)
}

func (m *MockStorageBackend) Save(ctx context.Context, path string, reader io.Reader) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, path, reader)
	}
	return nil
}

func (m *MockStorageBackend) Load(ctx context.Context, path string) (io.ReadCloser, error) {
	if m.loadFunc != nil {
		return m.loadFunc(ctx, path)
	}
	return io.NopCloser(strings.NewReader("mock content")), nil
}

func (m *MockStorageBackend) Delete(ctx context.Context, path string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, path)
	}
	return nil
}

func (m *MockStorageBackend) Exists(ctx context.Context, path string) (bool, error) {
	if m.existsFunc != nil {
		return m.existsFunc(ctx, path)
	}
	return true, nil
}

func (m *MockStorageBackend) GetURL(ctx context.Context, path string, expiresIn time.Duration) (string, error) {
	if m.getURLFunc != nil {
		return m.getURLFunc(ctx, path, expiresIn)
	}
	return "http://example.com/" + path, nil
}

// ============ Mock FileRepository ============

type MockFileRepository struct {
	createFunc       func(ctx context.Context, file *FileInfo) error
	getFunc          func(ctx context.Context, fileID string) (*FileInfo, error)
	updateFunc       func(ctx context.Context, fileID string, updates map[string]interface{}) error
	deleteFunc       func(ctx context.Context, fileID string) error
	listFunc         func(ctx context.Context, userID, category string, page, pageSize int) ([]*FileInfo, error)
	grantAccessFunc  func(ctx context.Context, fileID, userID string) error
	revokeAccessFunc func(ctx context.Context, fileID, userID string) error
	checkAccessFunc  func(ctx context.Context, fileID, userID string) (bool, error)
}

func (m *MockFileRepository) Create(ctx context.Context, file *FileInfo) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, file)
	}
	return nil
}

func (m *MockFileRepository) Get(ctx context.Context, fileID string) (*FileInfo, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, fileID)
	}
	return &FileInfo{
		ID:           fileID,
		Filename:     "test.txt",
		OriginalName: "test.txt",
		Path:         "test/test.txt",
		UserID:       "user123",
		IsPublic:     false,
	}, nil
}

func (m *MockFileRepository) Update(ctx context.Context, fileID string, updates map[string]interface{}) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, fileID, updates)
	}
	return nil
}

func (m *MockFileRepository) Delete(ctx context.Context, fileID string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, fileID)
	}
	return nil
}

func (m *MockFileRepository) List(ctx context.Context, userID, category string, page, pageSize int) ([]*FileInfo, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, userID, category, page, pageSize)
	}
	return []*FileInfo{}, nil
}

func (m *MockFileRepository) GrantAccess(ctx context.Context, fileID, userID string) error {
	if m.grantAccessFunc != nil {
		return m.grantAccessFunc(ctx, fileID, userID)
	}
	return nil
}

func (m *MockFileRepository) RevokeAccess(ctx context.Context, fileID, userID string) error {
	if m.revokeAccessFunc != nil {
		return m.revokeAccessFunc(ctx, fileID, userID)
	}
	return nil
}

func (m *MockFileRepository) CheckAccess(ctx context.Context, fileID, userID string) (bool, error) {
	if m.checkAccessFunc != nil {
		return m.checkAccessFunc(ctx, fileID, userID)
	}
	return false, nil
}

// ============ 测试用例 ============

func TestUpload(t *testing.T) {
	mockBackend := &MockStorageBackend{}
	mockRepo := &MockFileRepository{}
	service := NewStorageService(mockBackend, mockRepo)

	content := []byte("test file content")
	req := &UploadRequest{
		File:        bytes.NewReader(content),
		Filename:    "test.txt",
		ContentType: "text/plain",
		Size:        int64(len(content)),
		UserID:      "user123",
		IsPublic:    false,
		Category:    "document",
	}

	fileInfo, err := service.Upload(context.Background(), req)
	if err != nil {
		t.Errorf("上传文件失败: %v", err)
	}

	if fileInfo == nil {
		t.Error("文件信息为空")
	}

	if fileInfo.OriginalName != "test.txt" {
		t.Errorf("原始文件名错误: %s", fileInfo.OriginalName)
	}

	if fileInfo.UserID != "user123" {
		t.Errorf("用户ID错误: %s", fileInfo.UserID)
	}
}

func TestUpload_BackendFailure(t *testing.T) {
	mockBackend := &MockStorageBackend{
		saveFunc: func(ctx context.Context, path string, reader io.Reader) error {
			return errors.New("存储失败")
		},
	}
	mockRepo := &MockFileRepository{}
	service := NewStorageService(mockBackend, mockRepo)

	req := &UploadRequest{
		File:        bytes.NewReader([]byte("test")),
		Filename:    "test.txt",
		ContentType: "text/plain",
		Size:        4,
		UserID:      "user123",
		Category:    "document",
	}

	_, err := service.Upload(context.Background(), req)
	if err == nil {
		t.Error("期望上传失败，但成功了")
	}
}

func TestDownload(t *testing.T) {
	mockBackend := &MockStorageBackend{
		loadFunc: func(ctx context.Context, path string) (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader("file content")), nil
		},
	}
	mockRepo := &MockFileRepository{}
	service := NewStorageService(mockBackend, mockRepo)

	reader, err := service.Download(context.Background(), "file123")
	if err != nil {
		t.Errorf("下载文件失败: %v", err)
	}
	defer reader.Close()

	content, _ := io.ReadAll(reader)
	if string(content) != "file content" {
		t.Errorf("文件内容错误: %s", string(content))
	}
}

func TestDownload_FileNotFound(t *testing.T) {
	mockBackend := &MockStorageBackend{}
	mockRepo := &MockFileRepository{
		getFunc: func(ctx context.Context, fileID string) (*FileInfo, error) {
			return nil, errors.New("文件不存在")
		},
	}
	service := NewStorageService(mockBackend, mockRepo)

	_, err := service.Download(context.Background(), "nonexistent")
	if err == nil {
		t.Error("期望文件不存在错误，但成功了")
	}
}

func TestDelete(t *testing.T) {
	mockBackend := &MockStorageBackend{}
	mockRepo := &MockFileRepository{}
	service := NewStorageService(mockBackend, mockRepo)

	err := service.Delete(context.Background(), "file123")
	if err != nil {
		t.Errorf("删除文件失败: %v", err)
	}
}

func TestDelete_FileNotFound(t *testing.T) {
	mockBackend := &MockStorageBackend{}
	mockRepo := &MockFileRepository{
		getFunc: func(ctx context.Context, fileID string) (*FileInfo, error) {
			return nil, errors.New("文件不存在")
		},
	}
	service := NewStorageService(mockBackend, mockRepo)

	err := service.Delete(context.Background(), "nonexistent")
	if err == nil {
		t.Error("期望删除失败，但成功了")
	}
}

func TestGetFileInfo(t *testing.T) {
	mockBackend := &MockStorageBackend{}
	mockRepo := &MockFileRepository{
		getFunc: func(ctx context.Context, fileID string) (*FileInfo, error) {
			return &FileInfo{
				ID:       fileID,
				Filename: "test.txt",
				UserID:   "user123",
			}, nil
		},
	}
	service := NewStorageService(mockBackend, mockRepo)

	fileInfo, err := service.GetFileInfo(context.Background(), "file123")
	if err != nil {
		t.Errorf("获取文件信息失败: %v", err)
	}

	if fileInfo.Filename != "test.txt" {
		t.Errorf("文件名错误: %s", fileInfo.Filename)
	}
}

func TestCheckAccess_PublicFile(t *testing.T) {
	mockBackend := &MockStorageBackend{}
	mockRepo := &MockFileRepository{
		getFunc: func(ctx context.Context, fileID string) (*FileInfo, error) {
			return &FileInfo{
				ID:       fileID,
				IsPublic: true,
			}, nil
		},
	}
	service := NewStorageService(mockBackend, mockRepo)

	hasAccess, err := service.CheckAccess(context.Background(), "file123", "user456")
	if err != nil {
		t.Errorf("检查权限失败: %v", err)
	}

	if !hasAccess {
		t.Error("公开文件应该允许访问")
	}
}

func TestCheckAccess_Owner(t *testing.T) {
	mockBackend := &MockStorageBackend{}
	mockRepo := &MockFileRepository{
		getFunc: func(ctx context.Context, fileID string) (*FileInfo, error) {
			return &FileInfo{
				ID:       fileID,
				UserID:   "user123",
				IsPublic: false,
			}, nil
		},
	}
	service := NewStorageService(mockBackend, mockRepo)

	hasAccess, err := service.CheckAccess(context.Background(), "file123", "user123")
	if err != nil {
		t.Errorf("检查权限失败: %v", err)
	}

	if !hasAccess {
		t.Error("文件所有者应该有访问权限")
	}
}

func TestCheckAccess_NotAuthorized(t *testing.T) {
	mockBackend := &MockStorageBackend{}
	mockRepo := &MockFileRepository{
		getFunc: func(ctx context.Context, fileID string) (*FileInfo, error) {
			return &FileInfo{
				ID:       fileID,
				UserID:   "user123",
				IsPublic: false,
			}, nil
		},
		checkAccessFunc: func(ctx context.Context, fileID, userID string) (bool, error) {
			return false, nil
		},
	}
	service := NewStorageService(mockBackend, mockRepo)

	hasAccess, err := service.CheckAccess(context.Background(), "file123", "user456")
	if err != nil {
		t.Errorf("检查权限失败: %v", err)
	}

	if hasAccess {
		t.Error("未授权用户不应该有访问权限")
	}
}

func TestGrantAccess(t *testing.T) {
	mockBackend := &MockStorageBackend{}
	mockRepo := &MockFileRepository{}
	service := NewStorageService(mockBackend, mockRepo)

	err := service.GrantAccess(context.Background(), "file123", "user456")
	if err != nil {
		t.Errorf("授予权限失败: %v", err)
	}
}

func TestRevokeAccess(t *testing.T) {
	mockBackend := &MockStorageBackend{}
	mockRepo := &MockFileRepository{}
	service := NewStorageService(mockBackend, mockRepo)

	err := service.RevokeAccess(context.Background(), "file123", "user456")
	if err != nil {
		t.Errorf("撤销权限失败: %v", err)
	}
}

func TestListFiles(t *testing.T) {
	mockBackend := &MockStorageBackend{}
	mockRepo := &MockFileRepository{
		listFunc: func(ctx context.Context, userID, category string, page, pageSize int) ([]*FileInfo, error) {
			return []*FileInfo{
				{ID: "file1", Filename: "test1.txt"},
				{ID: "file2", Filename: "test2.txt"},
			}, nil
		},
	}
	service := NewStorageService(mockBackend, mockRepo)

	req := &ListFilesRequest{
		UserID:   "user123",
		Category: "document",
		Page:     1,
		PageSize: 20,
	}

	files, err := service.ListFiles(context.Background(), req)
	if err != nil {
		t.Errorf("列出文件失败: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("文件数量错误: %d", len(files))
	}
}

func TestListFiles_DefaultPagination(t *testing.T) {
	mockBackend := &MockStorageBackend{}
	mockRepo := &MockFileRepository{
		listFunc: func(ctx context.Context, userID, category string, page, pageSize int) ([]*FileInfo, error) {
			if page != 1 || pageSize != 20 {
				t.Errorf("默认分页参数错误: page=%d, pageSize=%d", page, pageSize)
			}
			return []*FileInfo{}, nil
		},
	}
	service := NewStorageService(mockBackend, mockRepo)

	req := &ListFilesRequest{
		UserID: "user123",
	}

	_, err := service.ListFiles(context.Background(), req)
	if err != nil {
		t.Errorf("列出文件失败: %v", err)
	}
}

func TestGetDownloadURL(t *testing.T) {
	mockBackend := &MockStorageBackend{
		getURLFunc: func(ctx context.Context, path string, expiresIn time.Duration) (string, error) {
			return "http://example.com/download/" + path, nil
		},
	}
	mockRepo := &MockFileRepository{}
	service := NewStorageService(mockBackend, mockRepo)

	url, err := service.GetDownloadURL(context.Background(), "file123", 1*time.Hour)
	if err != nil {
		t.Errorf("生成下载链接失败: %v", err)
	}

	if !strings.Contains(url, "download") {
		t.Errorf("下载链接格式错误: %s", url)
	}
}

func TestHealth(t *testing.T) {
	mockBackend := &MockStorageBackend{}
	mockRepo := &MockFileRepository{}
	service := NewStorageService(mockBackend, mockRepo)

	err := service.Health(context.Background())
	if err != nil {
		t.Errorf("健康检查失败: %v", err)
	}
}

func TestHealth_BackendFailure(t *testing.T) {
	mockBackend := &MockStorageBackend{
		existsFunc: func(ctx context.Context, path string) (bool, error) {
			return false, errors.New("存储后端不可用")
		},
	}
	mockRepo := &MockFileRepository{}
	service := NewStorageService(mockBackend, mockRepo)

	err := service.Health(context.Background())
	if err == nil {
		t.Error("期望健康检查失败，但成功了")
	}
}
