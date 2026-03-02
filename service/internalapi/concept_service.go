package internalapi

import (
	"context"
	"errors"

	"Qingyu_backend/models/writer"
	writerRepo "Qingyu_backend/repository/interfaces/writer"
)

// ConceptService 设定百科服务
// 处理Concept的CRUD操作，支持AI写作助手功能中的设定百科管理喵~
type ConceptService struct {
	repo writerRepo.ConceptRepository
}

// NewConceptService 创建ConceptService实例
func NewConceptService(repo writerRepo.ConceptRepository) *ConceptService {
	return &ConceptService{
		repo: repo,
	}
}

// CreateConceptRequest 创建概念请求
type CreateConceptRequest struct {
	UserID    string   `json:"user_id" binding:"required"`
	ProjectID string   `json:"project_id" binding:"required"`
	Name      string   `json:"name" binding:"required"`
	Category  string   `json:"category"`
	Content   string   `json:"content"`
	Tags      []string `json:"tags"`
}

// UpdateConceptRequest 更新概念请求
type UpdateConceptRequest struct {
	UserID    string   `json:"user_id" binding:"required"`
	ProjectID string   `json:"project_id" binding:"required"`
	Name      string   `json:"name"`
	Content   string   `json:"content"`
	Tags      []string `json:"tags"`
}

// Create 创建设定
//
// 创建新的设定概念，自动设置创建和更新时间戳喵~
//
// 参数：
//   - ctx: 上下文
//   - req: 创建请求，包含用户ID、项目ID、名称、分类、内容和标签
//
// 返回：
//   - *writer.Concept: 创建的设定对象
//   - error: 创建错误
func (s *ConceptService) Create(ctx context.Context, req *CreateConceptRequest) (*writer.Concept, error) {
	concept := &writer.Concept{}
	concept.ProjectID = req.ProjectID
	concept.Name = req.Name
	concept.Category = req.Category
	concept.Content = req.Content
	concept.Tags = req.Tags

	// 设置默认值和计算字段
	concept.BeforeCreate()

	if err := s.repo.Create(ctx, concept); err != nil {
		return nil, err
	}
	return concept, nil
}

// GetConcept 获取设定
//
// 根据ID获取单个设定对象喵~
//
// 参数：
//   - ctx: 上下文
//   - userID: 用户ID（预留权限检查）
//   - projectID: 项目ID（预留权限检查）
//   - conceptID: 设定ID
//
// 返回：
//   - *writer.Concept: 设定对象
//   - error: 查询错误
func (s *ConceptService) GetConcept(ctx context.Context, userID, projectID, conceptID string) (*writer.Concept, error) {
	concept, err := s.repo.GetByID(ctx, conceptID)
	if err != nil {
		return nil, err
	}
	if !sameProjectID(concept.ProjectID, projectID) {
		return nil, errors.New("concept not found")
	}
	return concept, nil
}

// Update 更新设定
//
// 更新现有设定的内容，仅更新提供的字段喵~
//
// 参数：
//   - ctx: 上下文
//   - conceptID: 设定ID
//   - req: 更新请求，包含名称、内容和标签
//
// 返回：
//   - *writer.Concept: 更新后的设定对象
//   - error: 更新错误
func (s *ConceptService) Update(ctx context.Context, conceptID string, req *UpdateConceptRequest) (*writer.Concept, error) {
	concept, err := s.repo.GetByID(ctx, conceptID)
	if err != nil {
		return nil, errors.New("concept not found")
	}

	if req.Name != "" {
		concept.Name = req.Name
	}
	if req.Content != "" {
		concept.Content = req.Content
	}
	if req.Tags != nil {
		concept.Tags = req.Tags
	}

	// 更新时间戳
	if err := s.repo.Update(ctx, concept); err != nil {
		return nil, err
	}
	return concept, nil
}

// Delete 删除设定
//
// 根据ID删除设定喵~
//
// 参数：
//   - ctx: 上下文
//   - userID: 用户ID（预留权限检查）
//   - projectID: 项目ID（预留权限检查）
//   - conceptID: 设定ID
//
// 返回：
//   - error: 删除错误
func (s *ConceptService) Delete(ctx context.Context, userID, projectID, conceptID string) error {
	concept, err := s.repo.GetByID(ctx, conceptID)
	if err != nil {
		return err
	}
	if !sameProjectID(concept.ProjectID, projectID) {
		return errors.New("concept not found")
	}
	return s.repo.Delete(ctx, conceptID)
}

// Search 搜索设定
//
// 支持按分类和关键词搜索设定喵~
//
// 参数：
//   - ctx: 上下文
//   - userID: 用户ID（预留权限检查）
//   - projectID: 项目ID
//   - category: 分类（可选）
//   - keyword: 关键词（可选）
//   - limit: 返回数量限制
//
// 返回：
//   - []*writer.Concept: 设定列表
//   - int: 总数
//   - error: 查询错误
func (s *ConceptService) Search(ctx context.Context, userID, projectID, category, keyword string, limit int) ([]*writer.Concept, int, error) {
	concepts, err := s.repo.Search(ctx, projectID, category, keyword)
	if err != nil {
		return nil, 0, err
	}

	// 应用limit限制
	if limit > 0 && len(concepts) > limit {
		concepts = concepts[:limit]
	}

	return concepts, len(concepts), nil
}

// BatchGet 批量获取设定
//
// 根据ID列表批量获取设定对象喵~
//
// 参数：
//   - ctx: 上下文
//   - userID: 用户ID（预留权限检查）
//   - projectID: 项目ID（预留权限检查）
//   - conceptIDs: 设定ID列表
//
// 返回：
//   - []*writer.Concept: 设定列表
//   - error: 查询错误
func (s *ConceptService) BatchGet(ctx context.Context, userID, projectID string, conceptIDs []string) ([]*writer.Concept, error) {
	concepts, err := s.repo.BatchGetByIDs(ctx, conceptIDs)
	if err != nil {
		return nil, err
	}
	filtered := make([]*writer.Concept, 0, len(concepts))
	for _, concept := range concepts {
		if sameProjectID(concept.ProjectID, projectID) {
			filtered = append(filtered, concept)
		}
	}
	return filtered, nil
}

// ListByProject 获取项目的设定列表
//
// 获取指定项目的所有设定喵~
//
// 参数：
//   - ctx: 上下文
//   - userID: 用户ID（预留权限检查）
//   - projectID: 项目ID
//
// 返回：
//   - []*writer.Concept: 设定列表
//   - int: 总数
//   - error: 查询错误
func (s *ConceptService) ListByProject(ctx context.Context, userID, projectID string) ([]*writer.Concept, int, error) {
	concepts, err := s.repo.ListByProject(ctx, projectID)
	return concepts, len(concepts), err
}
