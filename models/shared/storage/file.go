package storage

import "time"

// FileInfo 文件信息
type FileInfo struct {
	ID           string    `json:"id" bson:"_id,omitempty"`
	Filename     string    `json:"filename" bson:"filename"`                     // 存储文件名（UUID）
	OriginalName string    `json:"original_name" bson:"original_name"`           // 原始文件名
	ContentType  string    `json:"content_type" bson:"content_type"`             // MIME类型
	Size         int64     `json:"size" bson:"size"`                             // 文件大小（字节）
	Path         string    `json:"path" bson:"path"`                             // 存储路径
	UserID       string    `json:"user_id" bson:"user_id"`                       // 上传者ID
	IsPublic     bool      `json:"is_public" bson:"is_public"`                   // 是否公开
	Category     string    `json:"category" bson:"category"`                     // 文件分类
	MD5          string    `json:"md5,omitempty" bson:"md5,omitempty"`           // 文件MD5（用于去重）
	Width        int       `json:"width,omitempty" bson:"width,omitempty"`       // 图片宽度
	Height       int       `json:"height,omitempty" bson:"height,omitempty"`     // 图片高度
	Duration     int64     `json:"duration,omitempty" bson:"duration,omitempty"` // 视频/音频时长
	Downloads    int64     `json:"downloads" bson:"downloads"`                   // 下载次数
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}

// FileAccess 文件访问权限（可选，也可以通过IsPublic字段简化）
type FileAccess struct {
	FileID     string    `json:"file_id" bson:"file_id"`
	UserID     string    `json:"user_id" bson:"user_id"`
	Permission string    `json:"permission" bson:"permission"` // read, write, delete
	GrantedAt  time.Time `json:"granted_at" bson:"granted_at"`
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
