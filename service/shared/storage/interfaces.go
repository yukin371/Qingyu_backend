package storage

import (
	"context"
	"io"
	"time"
)

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

// FileInfo 文件信息
type FileInfo struct {
	ID           string    `json:"id" bson:"_id,omitempty"`
	Filename     string    `json:"filename" bson:"filename"`
	OriginalName string    `json:"original_name" bson:"original_name"`
	ContentType  string    `json:"content_type" bson:"content_type"`
	Size         int64     `json:"size" bson:"size"`
	Path         string    `json:"path" bson:"path"`
	URL          string    `json:"url,omitempty" bson:"-"` // 访问URL（不存储）
	UserID       string    `json:"user_id" bson:"user_id"`
	IsPublic     bool      `json:"is_public" bson:"is_public"`
	Category     string    `json:"category" bson:"category"`
	MD5          string    `json:"md5,omitempty" bson:"md5,omitempty"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}
