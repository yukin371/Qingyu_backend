package admin

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/bookstore"
	adminrep "Qingyu_backend/repository/interfaces/admin"
)

// CategoryAdminService 分类管理服务接口
type CategoryAdminService interface {
	CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*bookstore.Category, error)
	UpdateCategory(ctx context.Context, id string, req *UpdateCategoryRequest) (*bookstore.Category, error)
	DeleteCategory(ctx context.Context, id string) error
	GetCategoryTree(ctx context.Context) ([]*CategoryTreeNode, error)
	GetCategories(ctx context.Context, filter *CategoryFilter) ([]*bookstore.Category, error)
	GetCategoryByID(ctx context.Context, id string) (*bookstore.Category, error)
	MoveCategory(ctx context.Context, id string, req *MoveCategoryRequest) error
	SortCategory(ctx context.Context, id string, sortOrder int) error
}

// CategoryAdminServiceImpl 分类管理服务实现
type CategoryAdminServiceImpl struct {
	categoryRepo adminrep.CategoryAdminRepository
}

// NewCategoryAdminService 创建分类管理服务
func NewCategoryAdminService(categoryRepo adminrep.CategoryAdminRepository) CategoryAdminService {
	return &CategoryAdminServiceImpl{
		categoryRepo: categoryRepo,
	}
}

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Icon        string  `json:"icon"`
	ParentID    *string `json:"parent_id"`
	SortOrder   int     `json:"sort_order"`
}

// UpdateCategoryRequest 更新分类请求
type UpdateCategoryRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Icon        *string `json:"icon"`
	SortOrder   *int    `json:"sort_order"`
}

// CategoryFilter 分类筛选条件
type CategoryFilter struct {
	ParentID *string
	IsActive *bool
	Level    *int
}

// CategoryTreeNode 分类树节点
type CategoryTreeNode struct {
	*bookstore.Category
	Children []*CategoryTreeNode `json:"children,omitempty"`
}

// MoveCategoryRequest 移动分类请求
type MoveCategoryRequest struct {
	ParentID *string `json:"parent_id"`
	TargetID string  `json:"target_id"`
}

// CreateCategory 创建分类
func (s *CategoryAdminServiceImpl) CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*bookstore.Category, error) {
	exists, err := s.categoryRepo.NameExistsAtLevel(ctx, req.ParentID, req.Name, "")
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("CATEGORY_NAME_DUPLICATE: 分类名称已存在")
	}

	level := 0
	if req.ParentID != nil {
		parent, err := s.categoryRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return nil, errors.New("CATEGORY_INVALID_PARENT: 父分类不存在")
		}
		level = parent.Level + 1

		if level > 2 {
			return nil, errors.New("CATEGORY_LEVEL_EXCEEDED: 分类层级超过限制（最多3级）")
		}
	}

	now := time.Now()
	category := &bookstore.Category{
		ID:          primitive.NewObjectID().Hex(),
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		ParentID:    req.ParentID,
		Level:       level,
		SortOrder:   req.SortOrder,
		BookCount:   0,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// UpdateCategory 更新分类
func (s *CategoryAdminServiceImpl) UpdateCategory(ctx context.Context, id string, req *UpdateCategoryRequest) (*bookstore.Category, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("CATEGORY_NOT_FOUND: 分类不存在")
	}

	if req.Name != nil && *req.Name != category.Name {
		exists, err := s.categoryRepo.NameExistsAtLevel(ctx, category.ParentID, *req.Name, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("CATEGORY_NAME_DUPLICATE: 分类名称已存在")
		}
		category.Name = *req.Name
	}

	if req.Description != nil {
		category.Description = *req.Description
	}
	if req.Icon != nil {
		category.Icon = *req.Icon
	}
	if req.SortOrder != nil {
		category.SortOrder = *req.SortOrder
	}
	category.UpdatedAt = time.Now()

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory 删除分类
func (s *CategoryAdminServiceImpl) DeleteCategory(ctx context.Context, id string) error {
	hasChildren, err := s.categoryRepo.HasChildren(ctx, id)
	if err != nil {
		return err
	}
	if hasChildren {
		return errors.New("CATEGORY_HAS_CHILDREN: 无法删除：该分类下存在子分类")
	}

	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("CATEGORY_NOT_FOUND: 分类不存在")
	}
	if category.BookCount > 0 {
		return errors.New("CATEGORY_HAS_BOOKS: 无法删除：该分类下存在关联作品")
	}

	return s.categoryRepo.Delete(ctx, id)
}

// GetCategoryTree 获取分类树
func (s *CategoryAdminServiceImpl) GetCategoryTree(ctx context.Context) ([]*CategoryTreeNode, error) {
	allCategories, err := s.categoryRepo.List(ctx, nil)
	if err != nil {
		return nil, err
	}

	categoryMap := make(map[string]*CategoryTreeNode)
	for _, cat := range allCategories {
		categoryMap[cat.ID] = &CategoryTreeNode{Category: cat}
	}

	var rootNodes []*CategoryTreeNode
	for _, cat := range allCategories {
		node := categoryMap[cat.ID]
		if cat.ParentID == nil {
			rootNodes = append(rootNodes, node)
		} else {
			if parent, exists := categoryMap[*cat.ParentID]; exists {
				parent.Children = append(parent.Children, node)
			}
		}
	}

	// 对根节点和每个节点的子节点按 sort_order 排序
	s.sortCategoryNodes(rootNodes)

	return rootNodes, nil
}

// sortCategoryNodes 递归排序分类节点
func (s *CategoryAdminServiceImpl) sortCategoryNodes(nodes []*CategoryTreeNode) {
	if len(nodes) == 0 {
		return
	}

	// 按sort_order升序排序
	for i := 0; i < len(nodes)-1; i++ {
		for j := i + 1; j < len(nodes); j++ {
			if nodes[i].SortOrder > nodes[j].SortOrder {
				nodes[i], nodes[j] = nodes[j], nodes[i]
			}
		}
	}

	// 递归排序子节点
	for _, node := range nodes {
		s.sortCategoryNodes(node.Children)
	}
}

// GetCategories 获取分类列表
func (s *CategoryAdminServiceImpl) GetCategories(ctx context.Context, filter *CategoryFilter) ([]*bookstore.Category, error) {
	query := map[string]interface{}{}
	if filter != nil {
		if filter.ParentID != nil {
			query["parent_id"] = *filter.ParentID
		}
		if filter.IsActive != nil {
			query["is_active"] = *filter.IsActive
		}
		if filter.Level != nil {
			query["level"] = *filter.Level
		}
	}

	return s.categoryRepo.List(ctx, query)
}

// GetCategoryByID 获取分类详情
func (s *CategoryAdminServiceImpl) GetCategoryByID(ctx context.Context, id string) (*bookstore.Category, error) {
	return s.categoryRepo.GetByID(ctx, id)
}

// MoveCategory 移动分类
func (s *CategoryAdminServiceImpl) MoveCategory(ctx context.Context, id string, req *MoveCategoryRequest) error {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("CATEGORY_NOT_FOUND: 分类不存在")
	}

	if req.ParentID != nil {
		if err := s.detectCircularReference(ctx, id, *req.ParentID); err != nil {
			return err
		}
	}

	level := 0
	if req.ParentID != nil {
		parent, err := s.categoryRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return errors.New("CATEGORY_INVALID_PARENT: 父分类不存在")
		}
		level = parent.Level + 1

		if level > 2 {
			return errors.New("CATEGORY_LEVEL_EXCEEDED: 分类层级超过限制（最多3级）")
		}
	}

	category.ParentID = req.ParentID
	category.Level = level
	category.UpdatedAt = time.Now()

	return s.categoryRepo.Update(ctx, category)
}

// detectCircularReference 检测循环引用
func (s *CategoryAdminServiceImpl) detectCircularReference(ctx context.Context, categoryID, newParentID string) error {
	if categoryID == newParentID {
		return errors.New("CATEGORY_CIRCULAR_REF: 不能将分类移动到自己下面")
	}

	current := newParentID
	visited := make(map[string]bool)
	visited[current] = true

	for current != "" {
		parent, err := s.categoryRepo.GetByID(ctx, current)
		if err != nil {
			break
		}

		if parent.ID == categoryID {
			return errors.New("CATEGORY_CIRCULAR_REF: 检测到循环引用")
		}

		if parent.ParentID == nil {
			break
		}

		current = *parent.ParentID
		if visited[current] {
			break
		}
		visited[current] = true
	}

	return nil
}

// SortCategory 调整分类排序
func (s *CategoryAdminServiceImpl) SortCategory(ctx context.Context, id string, sortOrder int) error {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("CATEGORY_NOT_FOUND: 分类不存在")
	}

	category.SortOrder = sortOrder
	category.UpdatedAt = time.Now()

	return s.categoryRepo.Update(ctx, category)
}
