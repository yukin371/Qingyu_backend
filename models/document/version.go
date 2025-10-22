package document

import (
	"fmt"
	"time"
)

// Version 文档版本
// 用于存储文档的历史版本（新架构）
type Version struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	DocumentID  string    `bson:"document_id" json:"documentId" validate:"required"`
	VersionNum  int       `bson:"version_num" json:"versionNum"`                 // 版本号
	Content     string    `bson:"content" json:"content"`                        // 版本内容
	GridFSID    string    `bson:"gridfs_id,omitempty" json:"gridfsId,omitempty"` // 大文件GridFS ID
	ContentType string    `bson:"content_type" json:"contentType"`               // 内容类型
	WordCount   int       `bson:"word_count" json:"wordCount"`                   // 字数
	Comment     string    `bson:"comment,omitempty" json:"comment,omitempty"`    // 版本说明
	CreatedBy   string    `bson:"created_by" json:"createdBy"`                   // 创建人
	CreatedAt   time.Time `bson:"created_at" json:"createdAt"`                   // 创建时间
	IsAutoSave  bool      `bson:"is_auto_save" json:"isAutoSave"`                // 是否自动保存
}

// GetVersionLabel 获取版本标签
func (v *Version) GetVersionLabel() string {
	if v.IsAutoSave {
		return fmt.Sprintf("v%d (自动保存)", v.VersionNum)
	}
	return fmt.Sprintf("v%d", v.VersionNum)
}

// Validate 验证版本数据
func (v *Version) Validate() error {
	if v.DocumentID == "" {
		return fmt.Errorf("文档ID不能为空")
	}
	if v.VersionNum <= 0 {
		return fmt.Errorf("版本号必须大于0")
	}
	return nil
}

// TouchForCreate 设置创建时的默认值
func (v *Version) TouchForCreate() {
	if v.CreatedAt.IsZero() {
		v.CreatedAt = time.Now()
	}
}

// ================================================
// 以下是旧的版本控制结构（保留用于兼容）
// ================================================

// Commit 表示一次批量提交，可包含多个文件的修改
type Commit struct {
	ID        string                 `bson:"_id,omitempty" json:"id"`
	ProjectID string                 `bson:"project_id" json:"projectId"`
	AuthorID  string                 `bson:"author_id" json:"authorId"`
	Message   string                 `bson:"message,omitempty" json:"message,omitempty"`
	FileCount int                    `bson:"file_count" json:"fileCount"`                  // 本次提交涉及的文件数量
	Metadata  map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"` // 可选元数据
	CreatedAt time.Time              `bson:"created_at" json:"createdAt"`
}

// FileRevision 文件修订记录
type FileRevision struct {
	ID         string                 `bson:"_id,omitempty" json:"id"`
	ProjectID  string                 `bson:"project_id,omitempty" json:"projectId,omitempty"`
	NodeID     string                 `bson:"node_id" json:"nodeId"`
	CommitID   string                 ``
	Version    int                    `bson:"version" json:"version"`
	AuthorID   string                 `bson:"author_id" json:"authorId"` // 可为AI或用户
	Message    string                 `bson:"message,omitempty" json:"message,omitempty"`
	Snapshot   string                 `bson:"snapshot,omitempty" json:"snapshot,omitempty"` // 全文快照（可能为空，指向外部存储）
	ParentVers int                    `bson:"parent_version" json:"parentVersion"`
	Metadata   map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`      // 可选元数据
	StorageRef string                 `bson:"storage_ref,omitempty" json:"storageRef,omitempty"` // 外部存储引用（S3 path 等）
	Compressed bool                   `bson:"compressed,omitempty" json:"compressed,omitempty"`  // Snapshot 是否压缩
	CreatedAt  time.Time              `bson:"created_at" json:"createdAt"`
}

type Status string

const (
	StatusPending  Status = "pending"
	StatusApproved Status = "approved"
	StatusRejected Status = "rejected"
	StatusApplied  Status = "applied"
)

// FilePatch 表示一次候选变更（补丁）
type FilePatch struct {
	ID          string                 `bson:"_id,omitempty" json:"id"`
	ProjectID   string                 `bson:"project_id" json:"projectId"`
	NodeID      string                 `bson:"node_id" json:"nodeId"`
	BaseVersion int                    `bson:"base_version" json:"baseVersion"`
	DiffFormat  string                 `bson:"diff_format" json:"diffFormat"`              // unified|json-patch|full
	DiffPayload string                 `bson:"diff_payload" json:"diffPayload"`            // 原始diff内容
	CreatedBy   string                 `bson:"created_by" json:"createdBy"`                // ai|user:ID
	Status      Status                 `bson:"status" json:"status"`                       // pending|approved|rejected|applied
	Preview     string                 `bson:"preview,omitempty" json:"preview,omitempty"` // 可选：预览合并结果
	Metadata    map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
	CreatedAt   time.Time              `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time              `bson:"updated_at" json:"updatedAt"`
}

// ConflictInfo 冲突信息结构
type ConflictInfo struct {
	HasConflict          bool           `json:"hasConflict"`
	CurrentVersion       int            `json:"currentVersion"`
	ExpectedVersion      int            `json:"expectedVersion"`
	ConflictingRevisions []FileRevision `json:"conflictingRevisions,omitempty"`
	LastModified         time.Time      `json:"lastModified,omitempty"`
}

// BatchConflictResult 批量冲突检测结果
type BatchConflictResult struct {
	ProjectID    string                   `json:"projectId"`
	HasConflicts bool                     `json:"hasConflicts"`
	Conflicts    map[string]*ConflictInfo `json:"conflicts"`
}

// SnapshotStorage 快照存储策略配置
type SnapshotStorage struct {
	Strategy     string `json:"strategy"`     // inline, external, hybrid
	Threshold    int    `json:"threshold"`    // 内容大小阈值（字节）
	ExternalPath string `json:"externalPath"` // 外部存储路径
}

// GetSnapshotStrategy 根据内容大小决定存储策略
func GetSnapshotStrategy(contentSize int) string {
	const INLINE_THRESHOLD = 64 * 1024 // 64KB
	if contentSize <= INLINE_THRESHOLD {
		return "inline"
	}
	return "external"
}

// ConflictResolution 冲突解决策略
type ConflictResolution struct {
	Strategy      string `json:"strategy"`      // auto, manual, force
	ResolvedBy    string `json:"resolvedBy"`    // 解决者ID
	Resolution    string `json:"resolution"`    // 解决方案描述
	MergedContent string `json:"mergedContent"` // 合并后的内容
}

// BatchConflictResolution 批量冲突解决请求
type BatchConflictResolution struct {
	ProjectID   string                         `json:"projectId" binding:"required"`
	AuthorID    string                         `json:"authorId" binding:"required"`
	Message     string                         `json:"message"`
	Resolutions map[string]*ConflictResolution `json:"resolutions" binding:"required"` // nodeID -> resolution
}

// CommitFile 提交文件结构
type CommitFile struct {
	NodeID          string `json:"nodeId"`
	Content         string `json:"content"`
	ExpectedVersion int    `json:"expectedVersion"`
}
