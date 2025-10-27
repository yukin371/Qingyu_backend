package storage

import (
	"context"
	"io"
	"time"
	
	storageModel "Qingyu_backend/models/shared/storage"
)

// FileInfo 文件信息类型别名（使用models中的定义）
type FileInfo = storageModel.FileInfo

// StorageService 文件存储服务接口（对外暴露）
type StorageService interface {
	// 文件操作
	Upload(ctx context.Context, req *UploadRequest) (*FileInfo, error)
	Download(ctx context.Context, fileID string) (io.ReadCloser, error)
	Delete(ctx context.Context, fileID string) error
	GetFileInfo(ctx context.Context, fileID string) (*FileInfo, error)

	// 权限控制
	GrantAccess(ctx context.Context, fileID, userID string) error
	RevokeAccess(ctx context.Context, fileID, userID string) error
	CheckAccess(ctx context.Context, fileID, userID string) (bool, error)

	// 文件管理
	ListFiles(ctx context.Context, req *ListFilesRequest) ([]*FileInfo, error)
	GetDownloadURL(ctx context.Context, fileID string, expiresIn time.Duration) (string, error)

	// Health 健康检查
	Health(ctx context.Context) error
}

// ============ 请求结构 ============

// UploadRequest 上传请求
type UploadRequest struct {
	File        io.Reader `json:"-"`
	Filename    string    `json:"filename" binding:"required"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
	UserID      string    `json:"user_id" binding:"required"`
	IsPublic    bool      `json:"is_public"`
	Category    string    `json:"category"` // avatar, book, cover, attachment
}

// ListFilesRequest 查询文件列表请求
type ListFilesRequest struct {
	UserID   string `json:"user_id"`
	Category string `json:"category"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

// ============ 数据结构 ============
// FileInfo 已在上面通过type alias定义，使用models/shared/storage.FileInfo
