package project

import (
	"context"

	"Qingyu_backend/pkg/errors"
	model "Qingyu_backend/models/document"
	documentRepo "Qingyu_backend/repository/interfaces/writing"
	"Qingyu_backend/service/base"
)

// ProjectService 项目服务
type ProjectService struct {
	projectRepo  documentRepo.ProjectRepository
	indexManager documentRepo.ProjectIndexManager
	serviceName  string
	version      string
}

// NewProjectService 创建项目服务
func NewProjectService(projectRepo documentRepo.ProjectRepository, indexManager documentRepo.ProjectIndexManager) *ProjectService {
	return &ProjectService{
		projectRepo:  projectRepo,
		indexManager: indexManager,
		serviceName:  "ProjectService",
		version:      "1.0.0",
	}
}

// Initialize 初始化服务
func (s *ProjectService) Initialize(ctx context.Context) error {
	// 确保数据库索引
	if err := s.indexManager.EnsureIndexes(ctx); err != nil {
		return errors.ProjectFactory.InternalError("初始化索引失败", err).
			WithOperation("Initialize")
	}

	// 检查Repository健康状态
	if err := s.Health(ctx); err != nil {
		return errors.ProjectFactory.InternalError("项目Repository健康检查失败", err).
			WithOperation("Initialize")
	}

	return nil
}

// Health 健康检查
func (s *ProjectService) Health(ctx context.Context) error {
	return s.Health(ctx)
}

// Close 关闭服务
func (s *ProjectService) Close(ctx context.Context) error {
	return nil
}

// GetServiceName 获取服务名称
func (s *ProjectService) GetServiceName() string {
	return s.serviceName
}

// GetVersion 获取服务版本
func (s *ProjectService) GetVersion() string {
	return s.version
}

// CreateProject 创建项目
func (s *ProjectService) CreateProject(ctx context.Context, p *model.Project) (*model.Project, error) {
	if p.Name == "" {
		return nil, errors.ProjectFactory.ValidationError("项目名称不能为空", nil).
			WithOperation("CreateProject").
			WithMetadata(map[string]interface{}{
				"owner_id": p.OwnerID,
			})
	}

	if p.OwnerID == "" {
		return nil, errors.ProjectFactory.ValidationError("项目所有者不能为空", nil).
			WithOperation("CreateProject").
			WithMetadata(map[string]interface{}{
				"name": p.Name,
			})
	}

	err := s.projectRepo.Create(ctx, p)
	if err != nil {
		return nil, errors.ProjectFactory.InternalError("创建项目失败", err).
			WithOperation("CreateProject").
			WithMetadata(map[string]interface{}{
				"name": p.Name,
				"owner_id": p.OwnerID,
			})
	}

	return p, nil
}

// GetProjectByID 根据项目ID获取项目详情
func (s *ProjectService) GetProjectByID(ctx context.Context, projectID string) (*model.Project, error) {
	if projectID == "" {
		return nil, errors.ProjectFactory.ValidationError("项目ID不能为空", nil).
			WithOperation("GetProjectByID")
	}

	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, errors.ProjectFactory.InternalError("获取项目失败", err).
			WithOperation("GetProjectByID").
			WithMetadata(map[string]interface{}{
				"project_id": projectID,
			})
	}

	if project == nil {
		return nil, errors.ProjectFactory.NotFoundError("项目不存在", nil).
			WithOperation("GetProjectByID").
			WithMetadata(map[string]interface{}{
				"project_id": projectID,
			})
	}

	return project, nil
}

// GetProjectList 获取项目列表
func (s *ProjectService) GetProjectList(ctx context.Context, ownerID string, status string, limit, offset int64) ([]*model.Project, error) {
	if limit <= 0 {
		limit = 10 // 默认限制
	}

	if limit > 100 {
		limit = 100 // 最大限制
	}

	var projects []*model.Project
	var err error

	if ownerID != "" {
		projects, err = s.projectRepo.GetByOwnerAndStatus(ctx, ownerID, status, limit, offset)
	} else {
		// 如果没有指定所有者，使用通用查询（管理员功能）
		filter := make(map[string]interface{})
		if status != "" {
			filter["status"] = status
		}
		// 这里需要扩展Repository接口支持通用过滤
		return nil, errors.ProjectFactory.BusinessError("需要指定项目所有者", nil).
			WithOperation("GetProjectList").
			WithMetadata(map[string]interface{}{
				"status": status,
				"limit": limit,
				"offset": offset,
			})
	}

	if err != nil {
		return nil, errors.ProjectFactory.InternalError("获取项目列表失败", err).
			WithOperation("GetProjectList").
			WithMetadata(map[string]interface{}{
				"owner_id": ownerID,
				"status": status,
				"limit": limit,
				"offset": offset,
			})
	}

	return projects, nil
}

// UpdateProjectByID 更新项目
func (s *ProjectService) UpdateProjectByID(ctx context.Context, projectID, ownerID string, upd *model.Project) error {
	if projectID == "" || ownerID == "" {
		return errors.ProjectFactory.ValidationError("项目ID和所有者ID不能为空", nil).
			WithOperation("UpdateProjectByID").
			WithMetadata(map[string]interface{}{
				"project_id": projectID,
				"owner_id": ownerID,
			})
	}

	if upd == nil {
		return errors.ProjectFactory.ValidationError("更新信息不能为空", nil).
			WithOperation("UpdateProjectByID").
			WithMetadata(map[string]interface{}{
				"project_id": projectID,
				"owner_id": ownerID,
			})
	}

	// 构建更新字段
	updates := make(map[string]interface{})
	if upd.Name != "" {
		updates["name"] = upd.Name
	}
	if upd.Description != "" {
		updates["description"] = upd.Description
	}
	if upd.Status != "" {
		updates["status"] = upd.Status
	}

	if len(updates) == 0 {
		return errors.ProjectFactory.ValidationError("没有需要更新的字段", nil).
			WithOperation("UpdateProjectByID").
			WithMetadata(map[string]interface{}{
				"project_id": projectID,
				"owner_id": ownerID,
			})
	}

	err := s.projectRepo.UpdateByOwner(ctx, projectID, ownerID, updates)
	if err != nil {
		return errors.ProjectFactory.InternalError("更新项目失败", err).
			WithOperation("UpdateProjectByID").
			WithMetadata(map[string]interface{}{
				"project_id": projectID,
				"owner_id": ownerID,
				"updates": updates,
			})
	}

	return nil
}

// DeleteProjectByID 软删除项目
func (s *ProjectService) DeleteProjectByID(ctx context.Context, projectID, ownerID string) error {
	if projectID == "" || ownerID == "" {
		return errors.ProjectFactory.ValidationError("项目ID和所有者ID不能为空", nil).
			WithOperation("DeleteProjectByID").
			WithMetadata(map[string]interface{}{
				"project_id": projectID,
				"owner_id": ownerID,
			})
	}

	err := s.projectRepo.SoftDelete(ctx, projectID, ownerID)
	if err != nil {
		return errors.ProjectFactory.InternalError("删除项目失败", err).
			WithOperation("DeleteProjectByID").
			WithMetadata(map[string]interface{}{
				"project_id": projectID,
				"owner_id": ownerID,
			})
	}

	return nil
}

// RestoreProjectByID 恢复项目
func (s *ProjectService) RestoreProjectByID(ctx context.Context, projectID, ownerID string) error {
	if projectID == "" || ownerID == "" {
		return errors.ProjectFactory.ValidationError("项目ID和所有者ID不能为空", nil).
			WithOperation("RestoreProjectByID").
			WithMetadata(map[string]interface{}{
				"project_id": projectID,
				"owner_id": ownerID,
			})
	}

	err := s.projectRepo.Restore(ctx, projectID, ownerID)
	if err != nil {
		return errors.ProjectFactory.InternalError("恢复项目失败", err).
			WithOperation("RestoreProjectByID").
			WithMetadata(map[string]interface{}{
				"project_id": projectID,
				"owner_id": ownerID,
			})
	}

	return nil
}

// DeleteHard 硬删除项目（管理员功能）
func (s *ProjectService) DeleteHard(ctx context.Context, projectID string) error {
	if projectID == "" {
		return errors.ProjectFactory.ValidationError("项目ID不能为空", nil).
			WithOperation("DeleteHard").
			WithMetadata(map[string]interface{}{
				"project_id": projectID,
			})
	}

	err := s.projectRepo.HardDelete(ctx, projectID)
	if err != nil {
		return errors.ProjectFactory.InternalError("硬删除项目失败", err).
			WithOperation("DeleteHard").
			WithMetadata(map[string]interface{}{
				"project_id": projectID,
			})
	}

	return nil
}

// IsOwner 判断用户是否为项目所有者
func (s *ProjectService) IsOwner(ctx context.Context, projectID, userID string) (bool, error) {
	if projectID == "" || userID == "" {
		return false, errors.ProjectFactory.ValidationError("项目ID和用户ID不能为空", nil).
			WithOperation("IsOwner").
			WithMetadata(map[string]interface{}{
				"project_id": projectID,
				"user_id": userID,
			})
	}

	isOwner, err := s.projectRepo.IsOwner(ctx, projectID, userID)
	if err != nil {
		return false, errors.ProjectFactory.InternalError("检查所有者权限失败", err).
			WithOperation("IsOwner").
			WithMetadata(map[string]interface{}{
				"project_id": projectID,
				"user_id": userID,
			})
	}

	return isOwner, nil
}

// CreateWithRootNode 创建项目并初始化根节点（事务操作）
func (s *ProjectService) CreateWithRootNode(ctx context.Context, p *model.Project, rootNode *model.Node) error {
	if p == nil || rootNode == nil {
		return errors.ProjectFactory.ValidationError("项目和根节点信息不能为空", nil).
			WithOperation("CreateWithRootNode").
			WithMetadata(map[string]interface{}{
				"project_nil": p == nil,
				"root_node_nil": rootNode == nil,
			})
	}

	err := s.projectRepo.CreateWithTransaction(ctx, p, func(txCtx context.Context) error {
		// 在事务中创建根节点
		rootNode.ProjectID = p.ID
		rootNode.TouchForCreate()

		// 这里需要注入NodeRepository来创建根节点
		// 为了保持示例简单，暂时返回nil
		// 实际实现中应该通过依赖注入获取NodeRepository
		return nil
	})

	if err != nil {
		return errors.ProjectFactory.InternalError("创建项目和根节点失败", err).
			WithOperation("CreateWithRootNode").
			WithMetadata(map[string]interface{}{
				"project_id": p.ID,
				"project_name": p.Name,
				"owner_id": p.OwnerID,
			})
	}

	return nil
}
