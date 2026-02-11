package shared

import (
	"context"
	"time"
)

// StoragePort 文件存储端口（最小领域契约）
// 这是Port/Adapter模式中的"端口"，定义文件存储的抽象接口
type StoragePort interface {
	// Upload 上传文件
	// 返回：文件ID和可能的错误
	Upload(ctx context.Context, filename string, data []byte) (string, error)

	// Download 下载文件
	// 返回：文件内容和可能的错误
	Download(ctx context.Context, fileID string) ([]byte, error)

	// Delete 删除文件
	Delete(ctx context.Context, fileID string) error
}

// CachePort 缓存端口（最小领域契约）
// 这是Port/Adapter模式中的"端口"，定义缓存的抽象接口
type CachePort interface {
	// Get 获取缓存值
	// 返回：缓存值和可能的错误
	Get(ctx context.Context, key string) (string, error)

	// Set 设置缓存值
	// ttl: 过期时间，0表示永不过期
	Set(ctx context.Context, key string, value string, ttl time.Duration) error

	// Delete 删除缓存值
	Delete(ctx context.Context, key string) error
}

// AuthPort 认证端口（最小领域契约）
// 这是Port/Adapter模式中的"端口"，定义认证的抽象接口
type AuthPort interface {
	// ValidateToken 验证令牌
	// 返回：用户ID和可能的错误
	ValidateToken(ctx context.Context, token string) (string, error)

	// CheckPermission 检查权限
	// 返回：是否有权限和可能的错误
	CheckPermission(ctx context.Context, userID, permission string) (bool, error)
}
