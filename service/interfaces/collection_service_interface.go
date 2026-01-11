package interfaces

import (
	"context"
	readerModel "Qingyu_backend/models/reader"
)

// CollectionService 收藏服务接口
type CollectionService interface {
	// =========================
	// 收藏管理
	// =========================
	AddToCollection(ctx context.Context, userID, bookID, folderID, note string, tags []string, isPublic bool) (*readerModel.Collection, error)
	RemoveFromCollection(ctx context.Context, userID, collectionID string) error
	UpdateCollection(ctx context.Context, userID, collectionID string, updates map[string]interface{}) error
	GetUserCollections(ctx context.Context, userID, folderID string, page, size int) ([]*readerModel.Collection, int64, error)
	GetCollectionsByTag(ctx context.Context, userID, tag string, page, size int) ([]*readerModel.Collection, int64, error)
	IsCollected(ctx context.Context, userID, bookID string) (bool, error)

	// =========================
	// 收藏夹管理
	// =========================
	CreateFolder(ctx context.Context, userID, name, description string, isPublic bool) (*readerModel.CollectionFolder, error)
	GetUserFolders(ctx context.Context, userID string) ([]*readerModel.CollectionFolder, error)
	UpdateFolder(ctx context.Context, userID, folderID string, updates map[string]interface{}) error
	DeleteFolder(ctx context.Context, userID, folderID string) error

	// =========================
	// 收藏分享
	// =========================
	ShareCollection(ctx context.Context, userID, collectionID string) error
	UnshareCollection(ctx context.Context, userID, collectionID string) error
	GetPublicCollections(ctx context.Context, page, size int) ([]*readerModel.Collection, int64, error)

	// =========================
	// 统计
	// =========================
	GetUserCollectionStats(ctx context.Context, userID string) (map[string]interface{}, error)
}
