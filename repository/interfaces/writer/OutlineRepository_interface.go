package writer

import (
	"Qingyu_backend/models/writer"
	"context"
)

// OutlineRepository 大纲Repository接口
type OutlineRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, outline *writer.OutlineNode) error
	FindByID(ctx context.Context, outlineID string) (*writer.OutlineNode, error)
	FindByProjectID(ctx context.Context, projectID string) ([]*writer.OutlineNode, error)
	Update(ctx context.Context, outline *writer.OutlineNode) error
	Delete(ctx context.Context, outlineID string) error

	// 树形结构支持
	FindByParentID(ctx context.Context, projectID, parentID string) ([]*writer.OutlineNode, error)
	FindRoots(ctx context.Context, projectID string) ([]*writer.OutlineNode, error)

	// 辅助方法
	ExistsByID(ctx context.Context, outlineID string) (bool, error)
	CountByProjectID(ctx context.Context, projectID string) (int64, error)
	CountByParentID(ctx context.Context, projectID, parentID string) (int64, error)
}
