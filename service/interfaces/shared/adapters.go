package shared

import (
	"context"
	"io"
	"time"

	sharedAuth "Qingyu_backend/service/shared/auth"
	sharedStorage "Qingyu_backend/service/shared/storage"
)

// ============ 存储适配器 ============

// StorageAdapter 存储适配器
// 将现有的storage.StorageService适配到StoragePort接口
type StorageAdapter struct {
	service sharedStorage.StorageService
}

// NewStorageAdapter 创建存储适配器
func NewStorageAdapter(service sharedStorage.StorageService) StoragePort {
	return &StorageAdapter{service: service}
}

// Upload 实现StoragePort接口
func (a *StorageAdapter) Upload(ctx context.Context, filename string, data []byte) (string, error) {
	// 将[]byte转换为io.Reader
	reader := &byteReader{data: data}

	req := &sharedStorage.UploadRequest{
		File:     reader,
		Filename: filename,
		Size:     int64(len(data)),
		UserID:   "system", // 默认用户，可根据需要调整
	}

	fileInfo, err := a.service.Upload(ctx, req)
	if err != nil {
		return "", err
	}

	return fileInfo.ID, nil
}

// Download 实现StoragePort接口
func (a *StorageAdapter) Download(ctx context.Context, fileID string) ([]byte, error) {
	reader, err := a.service.Download(ctx, fileID)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}

// Delete 实现StoragePort接口
func (a *StorageAdapter) Delete(ctx context.Context, fileID string) error {
	return a.service.Delete(ctx, fileID)
}

// byteReader 辅助类型，用于将[]byte转换为io.Reader
type byteReader struct {
	data []byte
	pos  int
}

func (r *byteReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

// stringReader 辅助类型，用于将string转换为io.Reader
type stringReader struct {
	data string
	pos  int
}

func (r *stringReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

// ============ 缓存适配器 ============

// CacheAdapter 缓存适配器
// 将现有的cache服务适配到CachePort接口
type CacheAdapter struct {
	service interface {
		Get(ctx context.Context, key string) (string, error)
		Set(ctx context.Context, key string, value string, expiration time.Duration) error
		Delete(ctx context.Context, key string) error
	}
}

// NewCacheAdapter 创建缓存适配器
func NewCacheAdapter(service interface{}) CachePort {
	// 尝试转换为CacheService接口
	if cacheSvc, ok := service.(interface {
		Get(ctx context.Context, key string) (string, error)
		Set(ctx context.Context, key string, value string, expiration time.Duration) error
		Delete(ctx context.Context, key string) error
	}); ok {
		return &CacheAdapter{service: cacheSvc}
	}

	// 如果不是预期的接口，返回一个空实现
	return &CacheAdapter{}
}

// Get 实现CachePort接口
func (a *CacheAdapter) Get(ctx context.Context, key string) (string, error) {
	if a.service == nil {
		return "", nil
	}
	return a.service.Get(ctx, key)
}

// Set 实现CachePort接口
func (a *CacheAdapter) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	if a.service == nil {
		return nil
	}
	return a.service.Set(ctx, key, value, ttl)
}

// Delete 实现CachePort接口
func (a *CacheAdapter) Delete(ctx context.Context, key string) error {
	if a.service == nil {
		return nil
	}
	return a.service.Delete(ctx, key)
}

// ============ 认证适配器 ============

// AuthAdapter 认证适配器
// 将现有的auth.AuthService适配到AuthPort接口
type AuthAdapter struct {
	service sharedAuth.AuthService
}

// NewAuthAdapter 创建认证适配器
func NewAuthAdapter(service sharedAuth.AuthService) AuthPort {
	return &AuthAdapter{service: service}
}

// ValidateToken 实现AuthPort接口
func (a *AuthAdapter) ValidateToken(ctx context.Context, token string) (string, error) {
	claims, err := a.service.ValidateToken(ctx, token)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}

// CheckPermission 实现AuthPort接口
func (a *AuthAdapter) CheckPermission(ctx context.Context, userID, permission string) (bool, error) {
	return a.service.CheckPermission(ctx, userID, permission)
}
