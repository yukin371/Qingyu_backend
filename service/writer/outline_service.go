package writer

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/pkg/errors"
	writerRepo "Qingyu_backend/repository/interfaces/writer"
	"Qingyu_backend/service/base"
	serviceInterfaces "Qingyu_backend/service/interfaces"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OutlineService 大纲服务实现
type OutlineService struct {
	outlineRepo writerRepo.OutlineRepository
	eventBus    base.EventBus
	syncService *OutlineDocumentSyncService // 大纲-文档双向同步服务
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

// SetSyncService 注入同步服务（避免循环依赖）
func (s *OutlineService) SetSyncService(sync *OutlineDocumentSyncService) {
	s.syncService = sync
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
		if parent.ProjectID.Hex() != projectID {
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
	projectOID, _ := primitive.ObjectIDFromHex(projectID)
	outline.ProjectID = projectOID
	outline.Title = req.Title
	outline.ParentID = req.ParentID
	outline.Summary = req.Summary
	outline.Type = req.Type
	outline.Tension = req.Tension
	outline.DocumentID = req.DocumentID
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

	// 双向同步：如果指定了 document_id，更新文档的 outline_node_id 引用
	if outline.DocumentID != "" && s.syncService != nil {
		// 获取 DocumentRepository 来更新文档
		// 注意：这里需要访问 syncService 的 documentRepo
		// 由于避免循环依赖，我们通过 syncService 来处理
		go s.syncService.syncOutlineToDocument(ctx, outline.DocumentID, outline.ID.Hex())
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

	// 双向同步：自动创建对应文档
	if s.syncService != nil {
		s.syncService.SyncFromOutlineCreation(ctx, projectID, outline)
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
	if outline.ProjectID.Hex() != projectID {
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
			if parent.ProjectID.Hex() != projectID {
				return nil, errors.NewServiceError("OutlineService", errors.ServiceErrorForbidden, "parent outline does not belong to this project", "", nil)
			}
			// 防止将节点设置为自己的后代
			if parent.ID == outline.ID {
				return nil, errors.NewServiceError("OutlineService", errors.ServiceErrorValidation, "cannot set parent to self", "", nil)
			}
		}
	}

	// 记录原始层级用于同步
	oldLevel := s.calculateOutlineLevel(ctx, outline)

	// 更新字段（仅更新非nil字段）
	if req.Title != nil {
		outline.Title = *req.Title
		// 同步标题到关联文档
		if s.syncService != nil {
			go s.syncService.SyncTitleToDocument(ctx, outlineID, *req.Title)
		}
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
	if req.DocumentID != nil {
		outline.DocumentID = *req.DocumentID
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

	// 双向同步：检测层级变化并同步到文档
	if req.ParentID != nil && s.syncService != nil {
		newLevel := s.calculateOutlineLevel(ctx, outline)
		if newLevel != oldLevel {
			// 层级发生变化，需要同步文档类型变化
			go s.syncService.SyncLevelChangeToDocument(ctx, outlineID, newLevel)
		}
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

	// 双向同步：处理关联文档
	if s.syncService != nil {
		s.syncService.HandleOutlineDeletion(ctx, outline)
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
	// 自动填充 tags 字段（向后兼容旧数据）
	// 如果节点有 documentId 但没有 tags，自动添加 chapter-binding 标签
	if node.DocumentID != "" && len(node.Tags) == 0 {
		node.Tags = []string{fmt.Sprintf("chapter-binding:%s", node.DocumentID)}
	}

	treeNode := &serviceInterfaces.OutlineTreeNode{
		OutlineNode: node,
	}

	// 获取子节点
	children, err := s.outlineRepo.FindByParentID(ctx, projectID, node.ID.Hex())
	if err != nil {
		fmt.Printf("[buildTree] 获取子节点失败: %v, nodeID: %s, title: %s\n", err, node.ID.Hex(), node.Title)
	}
	if len(children) > 0 {
		fmt.Printf("[buildTree] 节点 %s 有 %d 个子节点\n", node.Title, len(children))
		treeNode.Children = make([]*serviceInterfaces.OutlineTreeNode, 0, len(children))
		for _, child := range children {
			childTreeNode := s.buildTree(ctx, child, projectID)
			treeNode.Children = append(treeNode.Children, childTreeNode)
		}
	} else {
		fmt.Printf("[buildTree] 节点 %s 没有子节点 (children=%d, err=%v)\n", node.Title, len(children), err)
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
	if parent.ProjectID.Hex() != projectID {
		return nil, errors.NewServiceError("OutlineService", errors.ServiceErrorForbidden, "no permission to access this outline", "", nil)
	}

	// 获取子节点
	return s.outlineRepo.FindByParentID(ctx, projectID, parentID)
}

// calculateOutlineLevel 计算大纲节点的层级深度
// Root (ParentID=="") = 0, Level 1 = 1, Level 2 = 2, etc.
func (s *OutlineService) calculateOutlineLevel(ctx context.Context, node *writer.OutlineNode) int {
	if node.ParentID == "" {
		return 0
	}
	parent, err := s.outlineRepo.FindByID(ctx, node.ParentID)
	if err != nil || parent == nil {
		return 0
	}
	return s.calculateOutlineLevel(ctx, parent) + 1
}
