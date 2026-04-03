package internalapi

import (
	"context"
	"fmt"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/repository/interfaces"
	writerRepo "Qingyu_backend/repository/interfaces/writer"

	"go.uber.org/zap"
)

// ContextAggregator 上下文聚合服务
// 供AI服务内部调用，聚合项目上下文数据（角色、大纲、文档内容、角色关系）
type ContextAggregator struct {
	projectRepo      writerRepo.ProjectRepository
	characterRepo    writerRepo.CharacterRepository
	outlineRepo      writerRepo.OutlineRepository
	documentRepo     writerRepo.DocumentRepository
	documentContent  writerRepo.DocumentContentRepository
	logger           *zap.Logger
}

// NewContextAggregator 创建上下文聚合服务
func NewContextAggregator(factory interfaces.RepositoryFactory, logger *zap.Logger) *ContextAggregator {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	return &ContextAggregator{
		projectRepo:     factory.CreateProjectRepository(),
		characterRepo:   factory.CreateCharacterRepository(),
		outlineRepo:     factory.CreateOutlineRepository(),
		documentRepo:    factory.CreateDocumentRepository(),
		documentContent: factory.CreateDocumentContentRepository(),
		logger:          logger,
	}
}

// ProjectContext 项目上下文汇总
type ProjectContext struct {
	ID          string                 `json:"id"`
	Title       string                 `json:"title"`
	WritingType string                 `json:"writingType,omitempty"`
	Summary     string                 `json:"summary,omitempty"`
	Status      string                 `json:"status"`
	Category    string                 `json:"category,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Statistics  writer.ProjectStats    `json:"statistics"`
	Settings    writer.ProjectSettings `json:"settings"`
}

// CharacterInfo 角色摘要（供AI服务使用）
type CharacterInfo struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Role        string   `json:"role,omitempty"`
	Traits      []string `json:"traits,omitempty"`
	Description string   `json:"description,omitempty"`
	Summary     string   `json:"summary,omitempty"`
}

// OutlineNodeInfo 大纲节点摘要（供AI服务使用）
type OutlineNodeInfo struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	ParentID    string   `json:"parentId,omitempty"`
	Order       int      `json:"order"`
	Level       int      `json:"level"`
	Summary     string   `json:"summary,omitempty"`
	Type        string   `json:"type,omitempty"`
	Tension     int      `json:"tension"`
	DocumentID  string   `json:"documentId,omitempty"`
	Characters  []string `json:"characters,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// DocumentContentInfo 文档内容摘要（供AI服务使用）
type DocumentContentInfo struct {
	Content   string `json:"content"`
	WordCount int    `json:"wordCount"`
}

// RelationInfo 角色关系摘要（供AI服务使用）
type RelationInfo struct {
	ID        string `json:"id"`
	FromID    string `json:"fromId"`
	ToID      string `json:"toId"`
	Type      string `json:"type"`
	Strength  int    `json:"strength"`
	Notes     string `json:"notes,omitempty"`
}

// GetProjectContext 获取项目上下文汇总
func (s *ContextAggregator) GetProjectContext(ctx context.Context, projectID string) (*ProjectContext, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("查找项目失败: %w", err)
	}
	if project == nil {
		return nil, fmt.Errorf("项目不存在: %s", projectID)
	}

	return &ProjectContext{
		ID:          project.ID.Hex(),
		Title:       project.Title,
		WritingType: project.WritingType,
		Summary:     project.Summary,
		Status:      string(project.Status),
		Category:    project.Category,
		Tags:        project.Tags,
		Statistics:  project.Statistics,
		Settings:    project.Settings,
	}, nil
}

// GetCharacters 获取项目角色列表
func (s *ContextAggregator) GetCharacters(ctx context.Context, projectID string) ([]*CharacterInfo, error) {
	characters, err := s.characterRepo.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("查找角色列表失败: %w", err)
	}

	result := make([]*CharacterInfo, 0, len(characters))
	for _, c := range characters {
		info := &CharacterInfo{
			ID:          c.ID.Hex(),
			Name:        c.Name,
			Traits:      c.Traits,
			Description: c.Background,
			Summary:     c.Summary,
		}
		// 使用短描述作为 role 字段
		if c.ShortDescription != "" {
			info.Role = c.ShortDescription
		}
		result = append(result, info)
	}

	return result, nil
}

// GetOutline 获取项目大纲树
func (s *ContextAggregator) GetOutline(ctx context.Context, projectID string) ([]*OutlineNodeInfo, error) {
	nodes, err := s.outlineRepo.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("查找大纲节点失败: %w", err)
	}

	// 构建 parentID -> level 映射用于计算层级
	parentLevel := map[string]int{}
	result := make([]*OutlineNodeInfo, 0, len(nodes))

	for _, node := range nodes {
		level := 1
		if node.ParentID != "" {
			if parentLevel, ok := parentLevel[node.ParentID]; ok {
				level = parentLevel + 1
			} else {
				// 回溯查找层级
				level = s.calculateNodeLevel(nodes, node.ParentID, 1)
			}
		}
		parentLevel[node.ID.Hex()] = level

		info := &OutlineNodeInfo{
			ID:         node.ID.Hex(),
			Title:      node.Title,
			ParentID:   node.ParentID,
			Order:      node.Order,
			Level:      level,
			Summary:    node.Summary,
			Type:       node.Type,
			Tension:    node.Tension,
			DocumentID: node.DocumentID,
			Characters: node.Characters,
			Tags:       node.Tags,
		}
		result = append(result, info)
	}

	return result, nil
}

// calculateNodeLevel 计算节点层级
func (s *ContextAggregator) calculateNodeLevel(nodes []*writer.OutlineNode, parentID string, currentLevel int) int {
	for _, n := range nodes {
		if n.ID.Hex() == parentID {
			if n.ParentID == "" {
				return currentLevel + 1
			}
			return s.calculateNodeLevel(nodes, n.ParentID, currentLevel+1)
		}
	}
	return currentLevel + 1
}

// GetDocumentContent 获取文档内容
func (s *ContextAggregator) GetDocumentContent(ctx context.Context, documentID string) (*DocumentContentInfo, error) {
	content, err := s.documentContent.GetByDocumentID(ctx, documentID)
	if err != nil {
		return nil, fmt.Errorf("查找文档内容失败: %w", err)
	}
	if content == nil {
		return &DocumentContentInfo{
			Content:   "",
			WordCount: 0,
		}, nil
	}

	return &DocumentContentInfo{
		Content:   content.Content,
		WordCount: content.GetDisplayWordCount(),
	}, nil
}

// GetCharacterRelations 获取角色关系列表
func (s *ContextAggregator) GetCharacterRelations(ctx context.Context, projectID string) ([]*RelationInfo, error) {
	// 传入 nil characterID 表示获取该项目下所有关系
	relations, err := s.characterRepo.FindRelations(ctx, projectID, nil)
	if err != nil {
		return nil, fmt.Errorf("查找角色关系失败: %w", err)
	}

	result := make([]*RelationInfo, 0, len(relations))
	for _, r := range relations {
		info := &RelationInfo{
			ID:       r.ID.Hex(),
			FromID:   r.FromID,
			ToID:     r.ToID,
			Type:     string(r.Type),
			Strength: r.Strength,
			Notes:    r.Notes,
		}
		result = append(result, info)
	}

	return result, nil
}
