package document

import "time"

// FileRevision 表示文件的一个历史版本
type FileRevision struct {
	ID         string    `bson:"_id,omitempty" json:"id"`
	ProjectID  string    `bson:"project_id" json:"projectId"`
	NodeID     string    `bson:"node_id" json:"nodeId"`
	Version    int       `bson:"version" json:"version"`
	AuthorID   string    `bson:"author_id" json:"authorId"` // 可为AI或用户
	Message    string    `bson:"message,omitempty" json:"message,omitempty"`
	Snapshot   string    `bson:"snapshot" json:"snapshot"` // 可选：全文快照
	ParentVers int       `bson:"parent_version" json:"parentVersion"`
	CreatedAt  time.Time `bson:"created_at" json:"createdAt"`
}

// FilePatch 表示一次候选变更（补丁）
type FilePatch struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	ProjectID   string    `bson:"project_id" json:"projectId"`
	NodeID      string    `bson:"node_id" json:"nodeId"`
	BaseVersion int       `bson:"base_version" json:"baseVersion"`
	DiffFormat  string    `bson:"diff_format" json:"diffFormat"`              // unified|json-patch
	DiffPayload string    `bson:"diff_payload" json:"diffPayload"`            // 原始diff内容
	CreatedBy   string    `bson:"created_by" json:"createdBy"`                // ai|user:ID
	Status      string    `bson:"status" json:"status"`                       // pending|approved|rejected|applied
	Preview     string    `bson:"preview,omitempty" json:"preview,omitempty"` // 可选：预览合并结果
	CreatedAt   time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updatedAt"`
}
