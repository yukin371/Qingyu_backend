package dto

import "time"

// ===========================
// Writer DTO（符合分层架构规范）
// ===========================
//
// 本文件包含 Writer 模块的数据传输对象（DTO）
//
// 命名和标签规范：
// - DTO 结构体使用驼峰命名（PascalCase）
// - JSON 字段标签使用驼峰命名（camelCase）
// - 对应的 MongoDB 模型（位于 models/writer/）使用蛇形命名（snake_case）的 BSON 标签
//
// 用途：
// - 用于 Service 层和 API 层之间的数据传输
// - ID 和时间字段统一使用字符串类型
// - 避免直接暴露 MongoDB 模型到 API 层

// ===========================
// Project DTOs
// ===========================

// CreateProjectRequest 创建项目请求
type CreateProjectRequest struct {
	Title       string   `json:"title" validate:"required,min=1,max=100"`
	Summary     string   `json:"summary,omitempty" validate:"max=500"`
	CoverURL    string   `json:"coverUrl,omitempty" validate:"omitempty,url,max=500"`
	Tags        []string `json:"tags,omitempty" validate:"max=10,dive,min=1,max=50"`
	Category    string   `json:"category,omitempty" validate:"max=50"`
	WritingType string   `json:"writingType,omitempty" validate:"omitempty,oneof=novel article script"`
}

// UpdateProjectRequest 更新项目请求
type UpdateProjectRequest struct {
	Title       *string   `json:"title,omitempty" validate:"omitempty,min=1,max=100"`
	Summary     *string   `json:"summary,omitempty" validate:"omitempty,max=500"`
	CoverURL    *string   `json:"coverUrl,omitempty" validate:"omitempty,url,max=500"`
	Tags        *[]string `json:"tags,omitempty" validate:"omitempty,max=10,dive,min=1,max=50"`
	Category    *string   `json:"category,omitempty" validate:"omitempty,max=50"`
	Status      *string   `json:"status,omitempty" validate:"omitempty,oneof=draft serializing completed suspended archived"`
	WritingType *string   `json:"writingType,omitempty" validate:"omitempty,oneof=novel article script"`
}

// ProjectResponse 项目响应
type ProjectResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	CoverURL  string    `json:"coverUrl"`
	Tags      []string  `json:"tags"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ListProjectsRequest 查询参数用于列出项目
type ListProjectsRequest struct {
	Page     int    `form:"page" validate:"min=1"`
	PageSize int    `form:"page_size" validate:"min=1,max=100"`
	Status   string `form:"status" validate:"omitempty,oneof=draft published archived"`
	Sort     string `form:"sort" validate:"omitempty,oneof=created_at updated_at title"`
	Order    string `form:"order" validate:"omitempty,oneof=asc desc"`
}

// ProjectListResponse 分页项目列表响应
type ProjectListResponse struct {
	Items    []ProjectResponse `json:"items"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"pageSize"`
}

// ===========================
// Document DTOs
// ===========================

// DocumentType represents the type of a document node
type DocumentType string

const (
	DocumentTypeVolume  DocumentType = "volume"  // 卷
	DocumentTypeChapter DocumentType = "chapter" // 章
	DocumentTypeSection DocumentType = "section" // 节
	DocumentTypeScene   DocumentType = "scene"   // 场景
)

// DocumentStatus represents the writing status of a document
type DocumentStatus string

const (
	DocumentStatusPlanned   DocumentStatus = "planned"   // 计划中
	DocumentStatusWriting   DocumentStatus = "writing"   // 写作中
	DocumentStatusCompleted DocumentStatus = "completed" // 已完成
)

// CreateDocumentRequest 创建文档请求（树形结构）
type CreateDocumentRequest struct {
	ProjectID    string       `json:"projectId" validate:"required"`
	ParentID     *string      `json:"parentId,omitempty"`
	Title        string       `json:"title" validate:"required,min=1,max=200"`
	Type         DocumentType `json:"type" validate:"required,oneof=volume chapter section scene"`
	Level        int          `json:"level" validate:"min=0,max=2"`    // 0=volume, 1=chapter, 2=section/scene
	Order        int          `json:"order,omitempty" validate:"min=0"` // 向后兼容字段
	OrderKey     string       `json:"orderKey,omitempty"`              // LexoRank排序键，可选，不传则自动生成
	CharacterIDs []string     `json:"characterIds,omitempty"`
	LocationIDs  []string     `json:"locationIds,omitempty"`
	TimelineIDs  []string     `json:"timelineIds,omitempty"`
	Tags         []string     `json:"tags,omitempty"`
	Notes        string       `json:"notes,omitempty"`
}

// UpdateDocumentRequest 更新文档元数据请求
type UpdateDocumentRequest struct {
	Title        *string         `json:"title,omitempty" validate:"omitempty,min=1,max=200"`
	Status       *DocumentStatus `json:"status,omitempty" validate:"omitempty,oneof=planned writing completed"`
	CharacterIDs *[]string       `json:"characterIds,omitempty" validate:"omitempty,max=50"`
	LocationIDs  *[]string       `json:"locationIds,omitempty" validate:"omitempty,max=50"`
	TimelineIDs  *[]string       `json:"timelineIds,omitempty" validate:"omitempty,max=50"`
	Tags         *[]string       `json:"tags,omitempty" validate:"omitempty,max=20"`
	Notes        *string         `json:"notes,omitempty" validate:"omitempty,max=1000"`
	OrderKey     *string         `json:"orderKey,omitempty"`
}

// DocumentResponse 文档响应（树形结构）
type DocumentResponse struct {
	ID           string         `json:"id"`
	ProjectID    string         `json:"projectId"`
	ParentID     *string        `json:"parentId,omitempty"`
	Title        string         `json:"title"`
	Type         DocumentType   `json:"type"`
	Level        int            `json:"level"`
	Order        int            `json:"order"`    // 向后兼容字段
	OrderKey     string         `json:"orderKey"` // LexoRank排序键
	Status       DocumentStatus `json:"status"`
	WordCount    int            `json:"wordCount"`
	CharacterIDs []string       `json:"characterIds,omitempty"`
	LocationIDs  []string       `json:"locationIds,omitempty"`
	TimelineIDs  []string       `json:"timelineIds,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	Notes        string         `json:"notes,omitempty"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
}

// DocumentTreeResponse 文档树响应（嵌套结构）
type DocumentTreeResponse struct {
	ProjectID string              `json:"projectId"`
	Documents []*DocumentTreeItem `json:"documents"`
}

// DocumentTreeItem 文档树节点
type DocumentTreeItem struct {
	ID        string              `json:"id"`
	ParentID  *string             `json:"parentId,omitempty"`
	Title     string              `json:"title"`
	Type      DocumentType        `json:"type"`
	Level     int                 `json:"level"`
	OrderKey  string              `json:"orderKey"`
	WordCount int                 `json:"wordCount"`
	Children  []*DocumentTreeItem `json:"children,omitempty"`
}

// ListDocumentsRequest 查询文档列表请求
type ListDocumentsRequest struct {
	ProjectID string `form:"project_id" validate:"required"`
	ParentID  string `form:"parent_id,omitempty"`
	Type      string `form:"type,omitempty" validate:"omitempty,oneof=volume chapter section scene"`
	Status    string `form:"status,omitempty" validate:"omitempty,oneof=planned writing completed"`
	Page      int    `form:"page" validate:"min=1"`
	PageSize  int    `form:"page_size" validate:"min=1,max=100"`
}

// DocumentListResponse 文档列表响应
type DocumentListResponse struct {
	Items    []DocumentResponse `json:"items"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"pageSize"`
}

// ReorderDocumentsRequest 重排序文档请求
type ReorderDocumentsRequest struct {
	ProjectID string        `json:"projectId" validate:"required"`
	ParentID  *string       `json:"parentId,omitempty"`
	Items     []ReorderItem `json:"items" validate:"required,min=1,dive"`
}

// ReorderItem 重排序项
type ReorderItem struct {
	DocumentID string  `json:"documentId" validate:"required"`
	ParentID   *string `json:"parentId,omitempty"`
	OrderKey   string  `json:"orderKey" validate:"required"`
}

// ===========================
// Additional Response DTOs
// ===========================

// CreateProjectResponse 创建项目响应（用于兼容旧代码）
type CreateProjectResponse struct {
	ProjectID string    `json:"projectId"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

// CreateDocumentResponse 创建文档响应
type CreateDocumentResponse struct {
	DocumentID string    `json:"documentId"`
	Title      string    `json:"title"`
	Type       string    `json:"type"`
	CreatedAt  time.Time `json:"createdAt"`
}

// MoveDocumentRequest 移动文档请求
type MoveDocumentRequest struct {
	DocumentID  string  `json:"documentId" validate:"required"`
	NewParentID *string `json:"newParentId,omitempty"`
	OrderKey    string  `json:"orderKey,omitempty"`
}
