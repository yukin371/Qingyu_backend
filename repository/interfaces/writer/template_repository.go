package writer

import (
	"Qingyu_backend/models/writer"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TemplateRepository 模板仓储接口
type TemplateRepository interface {
	// Create 创建模板
	Create(ctx context.Context, template *writer.Template) error

	// GetByID 根据ID获取模板
	GetByID(ctx context.Context, id primitive.ObjectID) (*writer.Template, error)

	// Update 更新模板
	Update(ctx context.Context, template *writer.Template) error

	// Delete 删除模板（硬删除）
	Delete(ctx context.Context, id primitive.ObjectID) error

	// ListByProject 列出项目的所有模板
	ListByProject(ctx context.Context, projectID primitive.ObjectID, filter *TemplateFilter) ([]*writer.Template, error)

	// ListByWorkspace 列出工作区的所有模板（包括全局模板）
	ListByWorkspace(ctx context.Context, workspaceID string, filter *TemplateFilter) ([]*writer.Template, error)

	// ListGlobal 列出所有全局模板
	ListGlobal(ctx context.Context, filter *TemplateFilter) ([]*writer.Template, error)

	// CountByProject 统计项目的模板数量
	CountByProject(ctx context.Context, projectID primitive.ObjectID) (int64, error)

	// ExistsByName 检查同名模板是否存在
	ExistsByName(ctx context.Context, projectID *primitive.ObjectID, name string, excludeID *primitive.ObjectID) (bool, error)
}

// TemplateFilter 模板查询过滤器
type TemplateFilter struct {
	Type      *writer.TemplateType
	Status    *writer.TemplateStatus
	Category  string
	Keyword   string // 搜索名称或描述
	IsSystem  *bool
	Page      int
	PageSize  int
	SortBy    string // name, created_at, updated_at
	SortOrder int // 1: asc, -1: desc
}
