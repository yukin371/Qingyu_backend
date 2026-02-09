package shared

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	sharedAuth "Qingyu_backend/service/shared/auth"
	sharedStorage "Qingyu_backend/service/shared/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============ Mock 实现 ============

// MockStorageService Mock存储服务
type MockStorageService struct {
	mock.Mock
}

func (m *MockStorageService) Upload(ctx context.Context, req *sharedStorage.UploadRequest) (*sharedStorage.FileInfo, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sharedStorage.FileInfo), args.Error(1)
}

func (m *MockStorageService) Download(ctx context.Context, fileID string) (io.ReadCloser, error) {
	args := m.Called(ctx, fileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockStorageService) Delete(ctx context.Context, fileID string) error {
	args := m.Called(ctx, fileID)
	return args.Error(0)
}

func (m *MockStorageService) GetFileInfo(ctx context.Context, fileID string) (*sharedStorage.FileInfo, error) {
	args := m.Called(ctx, fileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sharedStorage.FileInfo), args.Error(1)
}

func (m *MockStorageService) GrantAccess(ctx context.Context, fileID, userID string) error {
	args := m.Called(ctx, fileID, userID)
	return args.Error(0)
}

func (m *MockStorageService) RevokeAccess(ctx context.Context, fileID, userID string) error {
	args := m.Called(ctx, fileID, userID)
	return args.Error(0)
}

func (m *MockStorageService) CheckAccess(ctx context.Context, fileID, userID string) (bool, error) {
	args := m.Called(ctx, fileID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockStorageService) ListFiles(ctx context.Context, req *sharedStorage.ListFilesRequest) ([]*sharedStorage.FileInfo, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*sharedStorage.FileInfo), args.Error(1)
}

func (m *MockStorageService) GetDownloadURL(ctx context.Context, fileID string, expiresIn time.Duration) (string, error) {
	args := m.Called(ctx, fileID, expiresIn)
	return args.String(0), args.Error(1)
}

func (m *MockStorageService) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockAuthService Mock认证服务
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, req *sharedAuth.RegisterRequest) (*sharedAuth.RegisterResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sharedAuth.RegisterResponse), args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, req *sharedAuth.LoginRequest) (*sharedAuth.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sharedAuth.LoginResponse), args.Error(1)
}

func (m *MockAuthService) Logout(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockAuthService) ValidateToken(ctx context.Context, token string) (*sharedAuth.TokenClaims, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sharedAuth.TokenClaims), args.Error(1)
}

func (m *MockAuthService) CheckPermission(ctx context.Context, userID, permission string) (bool, error) {
	args := m.Called(ctx, userID, permission)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthService) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockAuthService) HasRole(ctx context.Context, userID, role string) (bool, error) {
	args := m.Called(ctx, userID, role)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthService) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockAuthService) CreateRole(ctx context.Context, req *sharedAuth.CreateRoleRequest) (*sharedAuth.Role, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sharedAuth.Role), args.Error(1)
}

func (m *MockAuthService) UpdateRole(ctx context.Context, roleID string, req *sharedAuth.UpdateRoleRequest) error {
	args := m.Called(ctx, roleID, req)
	return args.Error(0)
}

func (m *MockAuthService) DeleteRole(ctx context.Context, roleID string) error {
	args := m.Called(ctx, roleID)
	return args.Error(0)
}

func (m *MockAuthService) AssignRole(ctx context.Context, userID, roleID string) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

func (m *MockAuthService) RemoveRole(ctx context.Context, userID, roleID string) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

func (m *MockAuthService) CreateSession(ctx context.Context, userID string) (*sharedAuth.Session, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sharedAuth.Session), args.Error(1)
}

func (m *MockAuthService) GetSession(ctx context.Context, sessionID string) (*sharedAuth.Session, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sharedAuth.Session), args.Error(1)
}

func (m *MockAuthService) DestroySession(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockAuthService) RefreshSession(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockAuthService) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockAuthService) OAuthLogin(ctx context.Context, req *sharedAuth.OAuthLoginRequest) (*sharedAuth.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sharedAuth.LoginResponse), args.Error(1)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, token string) (string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.Error(1)
}

// ============ 适配器测试 ============

// TestStorageAdapter 测试存储适配器
func TestStorageAdapter(t *testing.T) {
	t.Run("应该适配上传操作", func(t *testing.T) {
		mockService := new(MockStorageService)
		mockService.On("Upload", mock.Anything, mock.Anything).Return(&sharedStorage.FileInfo{
			ID: "file-123",
		}, nil)

		adapter := NewStorageAdapter(mockService)
		fileID, err := adapter.Upload(context.Background(), "test.txt", []byte("test data"))

		assert.NoError(t, err)
		assert.Equal(t, "file-123", fileID)
		mockService.AssertExpectations(t)
	})

	t.Run("应该适配下载操作", func(t *testing.T) {
		mockService := new(MockStorageService)
		testData := []byte("test data")
		mockService.On("Download", mock.Anything, "file-123").Return(io.NopCloser(&byteReader{data: testData}), nil)

		adapter := NewStorageAdapter(mockService)
		data, err := adapter.Download(context.Background(), "file-123")

		assert.NoError(t, err)
		assert.Equal(t, testData, data)
		mockService.AssertExpectations(t)
	})

	t.Run("应该适配删除操作", func(t *testing.T) {
		mockService := new(MockStorageService)
		mockService.On("Delete", mock.Anything, "file-123").Return(nil)

		adapter := NewStorageAdapter(mockService)
		err := adapter.Delete(context.Background(), "file-123")

		assert.NoError(t, err)
		mockService.AssertExpectations(t)
	})

	t.Run("应该传递上传错误", func(t *testing.T) {
		mockService := new(MockStorageService)
		expectedErr := errors.New("upload failed")
		mockService.On("Upload", mock.Anything, mock.Anything).Return(&sharedStorage.FileInfo{}, expectedErr)

		adapter := NewStorageAdapter(mockService)
		_, err := adapter.Upload(context.Background(), "test.txt", []byte("test data"))

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

// TestAuthAdapter 测试认证适配器
func TestAuthAdapter(t *testing.T) {
	t.Run("应该适配令牌验证操作", func(t *testing.T) {
		mockService := new(MockAuthService)
		mockService.On("ValidateToken", mock.Anything, "valid-token").Return(&sharedAuth.TokenClaims{
			UserID: "user-123",
		}, nil)

		adapter := NewAuthAdapter(mockService)
		userID, err := adapter.ValidateToken(context.Background(), "valid-token")

		assert.NoError(t, err)
		assert.Equal(t, "user-123", userID)
		mockService.AssertExpectations(t)
	})

	t.Run("应该适配权限检查操作", func(t *testing.T) {
		mockService := new(MockAuthService)
		mockService.On("CheckPermission", mock.Anything, "user-123", "read:books").Return(true, nil)

		adapter := NewAuthAdapter(mockService)
		hasPermission, err := adapter.CheckPermission(context.Background(), "user-123", "read:books")

		assert.NoError(t, err)
		assert.True(t, hasPermission)
		mockService.AssertExpectations(t)
	})

	t.Run("应该传递令牌验证错误", func(t *testing.T) {
		mockService := new(MockAuthService)
		expectedErr := errors.New("invalid token")
		mockService.On("ValidateToken", mock.Anything, "invalid-token").Return(&sharedAuth.TokenClaims{}, expectedErr)

		adapter := NewAuthAdapter(mockService)
		_, err := adapter.ValidateToken(context.Background(), "invalid-token")

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

// TestCacheAdapter 测试缓存适配器
func TestCacheAdapter(t *testing.T) {
	t.Run("应该使用nil安全的方式处理空服务", func(t *testing.T) {
		adapter := NewCacheAdapter(nil)

		value, err := adapter.Get(context.Background(), "test-key")
		assert.NoError(t, err)
		assert.Empty(t, value)

		err = adapter.Set(context.Background(), "test-key", "test-value", time.Minute)
		assert.NoError(t, err)

		err = adapter.Delete(context.Background(), "test-key")
		assert.NoError(t, err)
	})
}

// TestPortInterfaceImplementation 测试端口接口实现
func TestPortInterfaceImplementation(t *testing.T) {
	t.Run("StorageAdapter应该实现StoragePort", func(t *testing.T) {
		var _ StoragePort = (*StorageAdapter)(nil)
	})

	t.Run("CacheAdapter应该实现CachePort", func(t *testing.T) {
		var _ CachePort = (*CacheAdapter)(nil)
	})

	t.Run("AuthAdapter应该实现AuthPort", func(t *testing.T) {
		var _ AuthPort = (*AuthAdapter)(nil)
	})
}
