package reading

import (
	"context"

	"Qingyu_backend/models/reader"
)

// CollectionRepository 收藏仓储接口
type CollectionRepository interface {
	// ========== 收藏管理 ==========

	// Create 创建收藏
	Create(ctx context.Context, collection *reader.Collection) error

	// GetByID 根据ID获取收藏
	GetByID(ctx context.Context, id string) (*reader.Collection, error)

	// GetByUserAndBook 根据用户ID和书籍ID获取收藏
	GetByUserAndBook(ctx context.Context, userID, bookID string) (*reader.Collection, error)

	// GetCollectionsByUser 获取用户的收藏列表
	GetCollectionsByUser(ctx context.Context, userID string, folderID string, page, size int) ([]*reader.Collection, int64, error)

	// GetCollectionsByTag 根据标签获取收藏
	GetCollectionsByTag(ctx context.Context, userID string, tag string, page, size int) ([]*reader.Collection, int64, error)

	// Update 更新收藏
	Update(ctx context.Context, id string, updates map[string]interface{}) error

	// Delete 删除收藏
	Delete(ctx context.Context, id string) error

	// ========== 收藏夹管理 ==========

	// CreateFolder 创建收藏夹
	CreateFolder(ctx context.Context, folder *reader.CollectionFolder) error

	// GetFolderByID 根据ID获取收藏夹
	GetFolderByID(ctx context.Context, id string) (*reader.CollectionFolder, error)

	// GetFoldersByUser 获取用户的收藏夹列表
	GetFoldersByUser(ctx context.Context, userID string) ([]*reader.CollectionFolder, error)

	// UpdateFolder 更新收藏夹
	UpdateFolder(ctx context.Context, id string, updates map[string]interface{}) error

	// DeleteFolder 删除收藏夹
	DeleteFolder(ctx context.Context, id string) error

	// IncrementFolderBookCount 增加收藏夹书籍数量
	IncrementFolderBookCount(ctx context.Context, folderID string) error

	// DecrementFolderBookCount 减少收藏夹书籍数量
	DecrementFolderBookCount(ctx context.Context, folderID string) error

	// ========== 公开收藏 ==========

	// GetPublicCollections 获取公开收藏列表
	GetPublicCollections(ctx context.Context, page, size int) ([]*reader.Collection, int64, error)

	// GetPublicFolders 获取公开收藏夹列表
	GetPublicFolders(ctx context.Context, page, size int) ([]*reader.CollectionFolder, int64, error)

	// ========== 统计 ==========

	// CountUserCollections 统计用户收藏数
	CountUserCollections(ctx context.Context, userID string) (int64, error)

	// Health 健康检查
	Health(ctx context.Context) error
}
