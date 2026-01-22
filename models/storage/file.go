package storage

import "time"

// FileInfo 文件信息
type FileInfo struct {
	ID           string    `json:"id" bson:"_id,omitempty"`
	Filename     string    `json:"filename" bson:"filename"`                     // 存储文件名（UUID）
	OriginalName string    `json:"$1$2" bson:"original_name"`           // 原始文件名
	ContentType  string    `json:"$1$2" bson:"content_type"`             // MIME类型
	Size         int64     `json:"size" bson:"size"`                             // 文件大小（字节）
	Path         string    `json:"path" bson:"path"`                             // 存储路径
	UserID       string    `json:"$1$2" bson:"user_id"`                       // 上传者ID
	IsPublic     bool      `json:"$1$2" bson:"is_public"`                   // 是否公开
	Category     string    `json:"category" bson:"category"`                     // 文件分类
	MD5          string    `json:"md5,omitempty" bson:"md5,omitempty"`           // 文件MD5（用于去重）
	Width        int       `json:"width,omitempty" bson:"width,omitempty"`       // 图片宽度
	Height       int       `json:"height,omitempty" bson:"height,omitempty"`     // 图片高度
	Duration     int64     `json:"duration,omitempty" bson:"duration,omitempty"` // 视频/音频时长
	Downloads    int64     `json:"downloads" bson:"downloads"`                   // 下载次数
	CreatedAt    time.Time `json:"$1$2" bson:"created_at"`
	UpdatedAt    time.Time `json:"$1$2" bson:"updated_at"`
}

// FileAccess 文件访问权限（可选，也可以通过IsPublic字段简化）
type FileAccess struct {
	FileID     string    `json:"$1$2" bson:"file_id"`
	UserID     string    `json:"$1$2" bson:"user_id"`
	Permission string    `json:"permission" bson:"permission"` // read, write, delete
	GrantedAt  time.Time `json:"$1$2" bson:"granted_at"`
	ExpiresAt  time.Time `json:"expires_at,omitempty" bson:"expires_at,omitempty"`
}

// 文件分类
const (
	CategoryAvatar     = "avatar"     // 头像
	CategoryBook       = "book"       // 书籍文件
	CategoryCover      = "cover"      // 封面图片
	CategoryAttachment = "attachment" // 附件
	CategoryDocument   = "document"   // 文档
	CategoryImage      = "image"      // 图片
	CategoryVideo      = "video"      // 视频
	CategoryAudio      = "audio"      // 音频
)

// 权限类型
const (
	PermissionRead   = "read"
	PermissionWrite  = "write"
	PermissionDelete = "delete"
)

// MultipartUpload 分片上传任务
type MultipartUpload struct {
	ID             string            `json:"id" bson:"_id,omitempty"`
	UploadID       string            `json:"$1$2" bson:"upload_id"`             // 上传任务ID
	FileID         string            `json:"$1$2" bson:"file_id"`                 // 文件ID
	FileName       string            `json:"$1$2" bson:"file_name"`             // 文件名
	FileSize       int64             `json:"$1$2" bson:"file_size"`             // 文件总大小
	ChunkSize      int64             `json:"$1$2" bson:"chunk_size"`           // 分片大小
	TotalChunks    int               `json:"$1$2" bson:"total_chunks"`       // 总分片数
	UploadedChunks []int             `json:"$1$2" bson:"uploaded_chunks"` // 已上传分片索引
	MD5Hash        string            `json:"md5_hash,omitempty" bson:"md5_hash"`     // 完整文件MD5
	StoragePath    string            `json:"$1$2" bson:"storage_path"`       // 存储路径
	UploadedBy     string            `json:"$1$2" bson:"uploaded_by"`         // 上传者
	Status         string            `json:"status" bson:"status"`                   // pending, uploading, completed, failed, aborted
	Metadata       map[string]string `json:"metadata,omitempty" bson:"metadata"`     // 扩展元数据
	ExpiresAt      time.Time         `json:"$1$2" bson:"expires_at"`           // 过期时间
	CreatedAt      time.Time         `json:"$1$2" bson:"created_at"`
	UpdatedAt      time.Time         `json:"$1$2" bson:"updated_at"`
	CompletedAt    *time.Time        `json:"completed_at,omitempty" bson:"completed_at"` // 完成时间
}
