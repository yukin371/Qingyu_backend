package storage

import (
	storageModel "Qingyu_backend/models/storage"
	"context"
	"time"
)

// StorageRepository 文件存储相关Repository
type StorageRepository interface {
	// 文件元数据管理
	CreateFile(ctx context.Context, file *storageModel.FileInfo) error
	GetFile(ctx context.Context, fileID string) (*storageModel.FileInfo, error)
	GetFileByMD5(ctx context.Context, md5Hash string) (*storageModel.FileInfo, error)
	UpdateFile(ctx context.Context, fileID string, updates map[string]interface{}) error
	DeleteFile(ctx context.Context, fileID string) error
	ListFiles(ctx context.Context, filter *FileFilter) ([]*storageModel.FileInfo, int64, error)
	CountFiles(ctx context.Context, filter *FileFilter) (int64, error)

	// 权限管理
	GrantAccess(ctx context.Context, fileID, userID string) error
	RevokeAccess(ctx context.Context, fileID, userID string) error
	CheckAccess(ctx context.Context, fileID, userID string) (bool, error)

	// 分片上传管理
	CreateMultipartUpload(ctx context.Context, upload *storageModel.MultipartUpload) error
	GetMultipartUpload(ctx context.Context, uploadID string) (*storageModel.MultipartUpload, error)
	UpdateMultipartUpload(ctx context.Context, uploadID string, updates map[string]interface{}) error
	CompleteMultipartUpload(ctx context.Context, uploadID string) error
	AbortMultipartUpload(ctx context.Context, uploadID string) error
	ListMultipartUploads(ctx context.Context, userID string, status string) ([]*storageModel.MultipartUpload, error)

	// 统计功能
	IncrementDownloadCount(ctx context.Context, fileID string) error

	// Health 健康检查
	Health(ctx context.Context) error
}

// FileFilter 文件过滤器
type FileFilter struct {
	UserID    string
	Category  string
	FileType  string
	Status    string
	IsPublic  *bool
	Tags      []string
	Keyword   string
	StartDate *time.Time
	EndDate   *time.Time
	MinSize   *int64
	MaxSize   *int64
	Page      int
	PageSize  int
	SortBy    string
	SortOrder string
	Limit     int64
	Offset    int64
}
