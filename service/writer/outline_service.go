package writer

import (
	"context"
	"time"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/pkg/errors"
	writerRepo "Qingyu_backend/repository/interfaces/writer"
	"Qingyu_backend/service/base"
	serviceInterfaces "Qingyu_backend/service/interfaces"
)

// OutlineService 大纲服务实现
type OutlineService struct {
	outlineRepo writerRepo.OutlineRepository
	eventBus    base.EventBus
}

// NewOutlineService 创建OutlineService实例
func NewOutlineService(
	outlineRepo writerRepo.OutlineRepository,
	eventBus base.EventBus,
) serviceInterfaces.OutlineService {
	return &OutlineService{
		outlineRepo: outlineRepo,
		eventBus:    eventBus,
	}
}

// Create 创建大纲节点
func (s *OutlineService) Create(
	ctx context.Context,
	projectID, userID string,
	req *serviceInterfaces.CreateOutlineRequest,
) (*writer.OutlineNode, error) {
	// 如果指定了父节点，验证父节点是否存在且属于同一项目
	if req.ParentID != "" {
		parent, err := s.outlineRepo.FindByID(ctx, req.ParentID)
		if err != nil {
			return nil, errors.NewServiceError("OutlineService", errors.ServiceErrorNotFound, "parent outline not found", "", err)
		}
		if parent.ProjectID != projectID {
			return nil, errors.NewServiceError("OutlineService", errors.ServiceErrorForbidden, "parent outline does not belong to this project", "", nil)
		}
	}

	// 如果没有指定order，设置为同级末尾
	order := req.Order
	if order == 0 && req.ParentID != "" {
		// 获取同级最大order
		count, _ := s.outlineRepo.CountByParentID(ctx, projectID, req.ParentID)
		order = int(count)
	}

	// 构建大纲节点
	outline := &writer.OutlineNode{}
	outline.ProjectID = projectID
	outline.Title = req.Title
	outline.ParentID = req.ParentID
	outline.Summary = req.Summary
	outline.Type = req.Type
	outline.Tension = req.Tension
	outline.ChapterID = req.ChapterID
	outline.Characters = req.Characters
	outline.Items = req.Items
	outline.Order = order

	// 验证数据
	if err := outline.Validate(); err != nil {
		return nil, errors.NewServiceError("OutlineService", errors.ServiceErrorValidation, "invalid outline data", err.Error(), nil)
	}

	// 保存到数据库
	if err := s.outlineRepo.Create(ctx, outline); err != nil {
		return nil, errors.NewServiceError("OutlineService", errors.ServiceErrorInternal, "create outline failed", "", err)
	}

	// 发布事件
	if s.eventBus != nil {
		event := &base.BaseEvent{
			EventType: "outline.created",
			EventData: map[string]interface{}{
				"outline_id": outline.ID,
				"project_id": projectID,
				"user_id":    userID,
				"title":      outline.Title,
			},
			Timestamp: time.Now(),
			Source:    "OutlineService",
		}
		s.eventBus.PublishAsync(ctx, event)
	}

	return outline, nil
}

// GetByID 根据ID获取大纲节点
func (s *OutlineService) GetByID(
	ctx context.Context,
	outlineID, projectID string,
) (*writer.OutlineNode, error) {
	outline, err := s.outlineRepo.FindByID(ctx, outlineID)
	if err != nil {
		return nil, err
	}

	// 验证项目权限
	if outline.ProjectID != projectID {
		return nil, errors.NewServiceError("OutlineService", errors.ServiceErrorForbidden, "no permission to access this outline", "", nil)
	}

	return outline, nil
}

// List 获取项目下的所有大纲节点
func (s *OutlineService) List(
	ctx context.Context,
	projectID string,
) ([]*writer.OutlineNode, error) {
	return s.outlineRepo.FindByProjectID(ctx, projectID)
}

// Update 更新大纲节点
func (s *OutlineService) Update(
	ctx context.Context,
	outlineID, projectID string,
	req *serviceInterfaces.UpdateOutlineRequest,
) (*writer.OutlineNode, error) {
	// 获取现有大纲
	outline, err := s.GetByID(ctx, outlineID, projectID)
	if err != nil {
		return nil, err
	}

	// 如果更新父节点，验证新父节点是否存在且属于同一项目
	if req.ParentID != nil && *req.ParentID != outline.ParentID {
		if *req.ParentID != "" {
			parent, err := s.outlineRepo.FindByID(ctx, *req.ParentID)
			if err != nil {
				return nil, errors.NewServiceError("OutlineService", errors.ServiceErrorNotFound, "parent outline not found", "", err)
			}
			if parent.ProjectID != projectID {
				return nil, errors.NewServiceError("OutlineService", errors.ServiceErrorForbidden, "parent outline does not belong to this project", "", nil)
			}
			// 防止将节点设置为自己的后代
			if parent.ID == outline.ID {
				return nil, errors.NewServiceError("OutlineService", errors.ServiceErrorValidation, "cannot set parent to self", "", nil)
			}
		}
	}

	// 更新字段（仅更新非nil字段）
	if req.Title != nil {
		outline.Title = *req.Title
	}
	if req.ParentID != nil {
		outline.ParentID = *req.ParentID
	}
	if req.Summary != nil {
		outline.Summary = *req.Summary
	}
	if req.Type != nil {
		outline.Type = *req.Type
	}
	if req.Tension != nil {
		outline.Tension = *req.Tension
	}
	if req.ChapterID != nil {
		outline.ChapterID = *req.ChapterID
	}
	if req.Characters != nil {
		outline.Characters = *req.Characters
	}
	if req.Items != nil {
		outline.Items = *req.Items
	}
	if req.Order != nil {
		outline.Order = *req.Order
	}

	// 验证数据
	if err := outline.Validate(); err != nil {
		return nil, errors.NewServiceError("OutlineService", errors.ServiceErrorValidation, "invalid outline data", err.Error(), nil)
	}

	// 保存更新
	if err := s.outlineRepo.Update(ctx, outline); err != nil {
		return nil, err
	}

	// 发布事件
	if s.eventBus != nil {
		event := &base.BaseEvent{
			EventType: "outline.updated",
			EventData: map[string]interface{}{
				"outline_id": outline.ID,
				"project_id": projectID,
			},
			Timestamp: time.Now(),
			Source:    "OutlineService",
		}
		s.eventBus.PublishAsync(ctx, event)
	}

	return outline, nil
}

// Delete 删除大纲节点
func (s *OutlineService) Delete(
	ctx context.Context,
	outlineID, projectID string,
) error {
	// 验证权限
	outline, err := s.GetByID(ctx, outlineID, projectID)
	if err != nil {
		return err
	}

	// 检查是否有子节点
	childCount, err := s.outlineRepo.CountByParentID(ctx, projectID, outlineID)
	if err != nil {
		return errors.NewServiceError("OutlineService", errors.ServiceErrorInternal, "check children failed", "", err)
	}
	if childCount > 0 {
		return errors.NewServiceError("OutlineService", errors.ServiceErrorBusiness, "cannot delete outline with children", "please delete or move children first", nil)
	}

	// 删除大纲
	if err := s.outlineRepo.Delete(ctx, outlineID); err != nil {
		return err
	}

	// 发布事件
	if s.eventBus != nil {
		event := &base.BaseEvent{
			EventType: "outline.deleted",
			EventData: map[string]interface{}{
				"outline_id": outlineID,
				"project_id": projectID,
			},
			Timestamp: time.Now(),
			Source:    "OutlineService",
		}
		s.eventBus.PublishAsync(ctx, event)
	}

	_ = outline // 避免未使用变量警告

	return nil
}

// GetTree 获取大纲树
func (s *OutlineService) GetTree(
	ctx context.Context,
	projectID string,
) ([]*serviceInterfaces.OutlineTreeNode, error) {
	// 获取所有根节点
	roots, err := s.outlineRepo.FindRoots(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// 构建树形结构
	tree := make([]*serviceInterfaces.OutlineTreeNode, 0, len(roots))
	for _, root := range roots {
		treeNode := s.buildTree(ctx, root, projectID)
		tree = append(tree, treeNode)
	}

	return tree, nil
}

// buildTree 递归构建树形结构
func (s *OutlineService) buildTree(
	ctx context.Context,
	node *writer.OutlineNode,
	projectID string,
) *serviceInterfaces.OutlineTreeNode {
	treeNode := &serviceInterfaces.OutlineTreeNode{
		OutlineNode: node,
	}

	// 获取子节点
	children, _ := s.outlineRepo.FindByParentID(ctx, projectID, node.ID.Hex())
	if len(children) > 0 {
		treeNode.Children = make([]*serviceInterfaces.OutlineTreeNode, 0, len(children))
		for _, child := range children {
			childTreeNode := s.buildTree(ctx, child, projectID)
			treeNode.Children = append(treeNode.Children, childTreeNode)
		}
	}

	return treeNode
}

// GetChildren 获取子节点列表
func (s *OutlineService) GetChildren(
	ctx context.Context,
	projectID, parentID string,
) ([]*writer.OutlineNode, error) {
	// 如果parentID为空，返回根节点
	if parentID == "" {
		return s.outlineRepo.FindRoots(ctx, projectID)
	}

	// 验证父节点权限
	parent, err := s.outlineRepo.FindByID(ctx, parentID)
	if err != nil {
		return nil, err
	}
	if parent.ProjectID != projectID {
		return nil, errors.NewServiceError("OutlineService", errors.ServiceErrorForbidden, "no permission to access this outline", "", nil)
	}

	// 获取子节点
	return s.outlineRepo.FindByParentID(ctx, projectID, parentID)
}
