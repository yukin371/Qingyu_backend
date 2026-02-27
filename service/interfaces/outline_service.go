package interfaces

import (
	"Qingyu_backend/models/dto"
	"Qingyu_backend/models/writer"
	"context"
)

// OutlineService 大纲服务接口
type OutlineService interface {
	// 基础CRUD
	Create(ctx context.Context, projectID, userID string, req *CreateOutlineRequest) (*writer.OutlineNode, error)
	GetByID(ctx context.Context, outlineID, projectID string) (*writer.OutlineNode, error)
	List(ctx context.Context, projectID string) ([]*writer.OutlineNode, error)
	Update(ctx context.Context, outlineID, projectID string, req *UpdateOutlineRequest) (*writer.OutlineNode, error)
	Delete(ctx context.Context, outlineID, projectID string) error

	// 树形结构操作
	GetTree(ctx context.Context, projectID string) ([]*OutlineTreeNode, error)
	GetChildren(ctx context.Context, projectID, parentID string) ([]*writer.OutlineNode, error)
}

// CreateOutlineRequest 创建大纲请求
// Deprecated: 使用 dto.CreateOutlineRequest 替代
type CreateOutlineRequest = dto.CreateOutlineRequest

// UpdateOutlineRequest 更新大纲请求
// Deprecated: 使用 dto.UpdateOutlineRequest 替代
type UpdateOutlineRequest = dto.UpdateOutlineRequest

// OutlineTreeNode 大纲树节点
// Deprecated: 保留大纲模块特定类型
type OutlineTreeNode struct {
	*writer.OutlineNode
	Children []*OutlineTreeNode `json:"children,omitempty"`
}
