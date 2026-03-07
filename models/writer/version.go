package writer

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Version 文档版本
// 用于存储文档的历史版本（新架构）
type Version struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DocumentID  string             `bson:"document_id" json:"documentId" validate:"required"`
	VersionNum  int                `bson:"version_num" json:"versionNum"`
	Content     string             `bson:"content" json:"content"`
	GridFSID    string             `bson:"gridfs_id,omitempty" json:"gridfsId,omitempty"`
	ContentType string             `bson:"content_type" json:"contentType"`
	WordCount   int                `bson:"word_count" json:"wordCount"`
	Comment     string             `bson:"comment,omitempty" json:"comment,omitempty"`
	CreatedBy   string             `bson:"created_by" json:"createdBy"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	IsAutoSave  bool               `bson:"is_auto_save" json:"isAutoSave"`
}

func (v *Version) GetVersionLabel() string {
	if v.IsAutoSave {
		return fmt.Sprintf("v%d (自动保存)", v.VersionNum)
	}
	return fmt.Sprintf("v%d", v.VersionNum)
}

func (v *Version) Validate() error {
	if v.DocumentID == "" {
		return fmt.Errorf("文档ID不能为空")
	}
	if v.VersionNum <= 0 {
		return fmt.Errorf("版本号必须大于0")
	}
	return nil
}

func (v *Version) TouchForCreate() {
	if v.CreatedAt.IsZero() {
		v.CreatedAt = time.Now()
	}
}

// ================================================
// 以下是旧的版本控制结构（保留用于兼容）
// ================================================

type Commit struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	ProjectID string                 `bson:"project_id" json:"projectId"`
	AuthorID  string                 `bson:"author_id" json:"authorId"`
	Message   string                 `bson:"message,omitempty" json:"message,omitempty"`
	FileCount int                    `bson:"file_count" json:"fileCount"`
	Metadata  map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
	CreatedAt time.Time              `bson:"created_at" json:"createdAt"`
}

type FileRevision struct {
	ID         primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	ProjectID  string                 `bson:"project_id,omitempty" json:"projectId,omitempty"`
	NodeID     string                 `bson:"node_id" json:"nodeId"`
	CommitID   string                 `bson:"commit_id,omitempty" json:"commitId,omitempty"`
	Version    int                    `bson:"version" json:"version"`
	AuthorID   string                 `bson:"author_id" json:"authorId"`
	Message    string                 `bson:"message,omitempty" json:"message,omitempty"`
	Snapshot   string                 `bson:"snapshot,omitempty" json:"snapshot,omitempty"`
	ParentVers int                    `bson:"parent_version" json:"parentVersion"`
	Metadata   map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
	StorageRef string                 `bson:"storage_ref,omitempty" json:"storageRef,omitempty"`
	Compressed bool                   `bson:"compressed,omitempty" json:"compressed,omitempty"`
	CreatedAt  time.Time              `bson:"created_at" json:"createdAt"`
}

type Status string

const (
	StatusPending  Status = "pending"
	StatusApproved Status = "approved"
	StatusRejected Status = "rejected"
	StatusApplied  Status = "applied"
)

type FilePatch struct {
	ID          primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	ProjectID   string                 `bson:"project_id" json:"projectId"`
	NodeID      string                 `bson:"node_id" json:"nodeId"`
	BaseVersion int                    `bson:"base_version" json:"baseVersion"`
	DiffFormat  string                 `bson:"diff_format" json:"diffFormat"`
	DiffPayload string                 `bson:"diff_payload" json:"diffPayload"`
	CreatedBy   string                 `bson:"created_by" json:"createdBy"`
	Status      Status                 `bson:"status" json:"status"`
	Preview     string                 `bson:"preview,omitempty" json:"preview,omitempty"`
	Metadata    map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
	CreatedAt   time.Time              `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time              `bson:"updated_at" json:"updatedAt"`
}

type ConflictInfo struct {
	HasConflict          bool           `json:"hasConflict"`
	CurrentVersion       int            `json:"currentVersion"`
	ExpectedVersion      int            `json:"expectedVersion"`
	ConflictingRevisions []FileRevision `json:"conflictingRevisions,omitempty"`
	LastModified         time.Time      `json:"lastModified,omitempty"`
}

type BatchConflictResult struct {
	ProjectID    string                   `json:"projectId"`
	HasConflicts bool                     `json:"hasConflicts"`
	Conflicts    map[string]*ConflictInfo `json:"conflicts"`
}

type SnapshotStorage struct {
	Strategy     string `json:"strategy"`
	Threshold    int    `json:"threshold"`
	ExternalPath string `json:"externalPath"`
}

func GetSnapshotStrategy(contentSize int) string {
	const InlineThreshold = 64 * 1024
	if contentSize <= InlineThreshold {
		return "inline"
	}
	return "external"
}

type ConflictResolution struct {
	Strategy      string `json:"strategy"`
	ResolvedBy    string `json:"resolvedBy"`
	Resolution    string `json:"resolution"`
	MergedContent string `json:"mergedContent"`
}

type BatchConflictResolution struct {
	ProjectID   string                         `json:"projectId" binding:"required"`
	AuthorID    string                         `json:"authorId" binding:"required"`
	Message     string                         `json:"message"`
	Resolutions map[string]*ConflictResolution `json:"resolutions" binding:"required"`
}

type CommitFile struct {
	NodeID          string `json:"nodeId"`
	Content         string `json:"content"`
	ExpectedVersion int    `json:"expectedVersion"`
}
