package ai

import (
	"Qingyu_backend/models/writer"
	"context"
	"fmt"
	"strings"

	"Qingyu_backend/models/ai"
	writerRepo "Qingyu_backend/repository/interfaces/writer"
	documentService "Qingyu_backend/service/writer/project"
)

// ContextService AI上下文服务
type ContextService struct {
	documentService *documentService.DocumentService
	projectService  *documentService.ProjectService
	nodeService     *documentService.NodeService
	versionService  *documentService.VersionService

	// documentContentRepo: 临时架构债务
	// TODO(架构重构): 当前使用 nil 是因为 ai_service.go 中采用了旧的直接实例化方式
	// 而非依赖注入。待整体架构迁移到 Repository Factory 模式后统一解决。
	// 相关讨论: doc/architecture/架构设计规范.md - 依赖注入原则
	documentContentRepo writerRepo.DocumentContentRepository
}

// NewContextService 创建AI上下文服务
func NewContextService(
	documentService *documentService.DocumentService,
	projectService *documentService.ProjectService,
	nodeService *documentService.NodeService,
	versionService *documentService.VersionService,
	documentContentRepo writerRepo.DocumentContentRepository,
) *ContextService {
	return &ContextService{
		documentService:     documentService,
		projectService:      projectService,
		nodeService:         nodeService,
		versionService:      versionService,
		documentContentRepo: documentContentRepo,
	}
}

// BuildContext 构建AI上下文
func (s *ContextService) BuildContext(ctx context.Context, projectID string, chapterID string) (*ai.AIContext, error) {
	// 获取项目信息（直接从repository，跳过权限检查）
	// AI上下文构建不应该受权限限制，因为已经通过配额中间件验证
	project, err := s.projectService.GetByIDWithoutAuth(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("获取项目信息失败: %w", err)
	}
	if project == nil {
		return nil, fmt.Errorf("项目不存在")
	}

	// 构建章节信息（可选）
	var chapterInfo *ai.ChapterInfo
	if chapterID != "" {
		var err error
		chapterInfo, err = s.buildChapterInfo(ctx, projectID, chapterID)
		if err != nil {
			return nil, fmt.Errorf("构建章节信息失败: %w", err)
		}
	}

	aiContext := &ai.AIContext{
		ProjectID:        projectID,
		CurrentChapter:   chapterInfo,
		ActiveCharacters: []*ai.CharacterInfo{},  // TODO: 实现角色信息获取
		CurrentLocations: []*ai.LocationInfo{},   // TODO: 实现地点信息获取
		RelevantEvents:   []*ai.TimelineEvent{},  // TODO: 实现时间线事件获取
		PreviousChapters: []*ai.ChapterSummary{}, // TODO: 实现前序章节摘要
		NextChapters:     []*ai.ChapterOutline{}, // TODO: 实现下一章节大纲
		WorldSettings:    nil,                    // TODO: 实现世界观设定获取
		PlotThreads:      []*ai.PlotThread{},     // TODO: 实现情节线索获取
		TokenCount:       0,                      // TODO: 计算token数量
	}

	return aiContext, nil
}

// buildChapterInfo 构建章节信息
func (s *ContextService) buildChapterInfo(ctx context.Context, projectID string, chapterID string) (*ai.ChapterInfo, error) {
	// 获取章节文档元数据
	doc, err := s.documentService.GetByID(chapterID)
	if err != nil {
		return nil, fmt.Errorf("获取章节文档失败: %w", err)
	}
	if doc == nil {
		return nil, fmt.Errorf("章节不存在")
	}

	// 获取章节内容
	// 注意: documentContentRepo 可能为 nil（旧架构遗留），需要防御性检查
	var docContent *writer.DocumentContent
	if s.documentContentRepo != nil {
		var err error
		docContent, err = s.documentContentRepo.GetByDocumentID(ctx, chapterID)
		if err != nil {
			return nil, fmt.Errorf("获取章节内容失败: %w", err)
		}
	}

	// 提取摘要（从内容中提取前200字符）
	summary := ""
	if docContent != nil && docContent.Content != "" {
		content := strings.TrimSpace(docContent.Content)
		if len(content) > 200 {
			summary = content[:200] + "..."
		} else {
			summary = content
		}
	}

	// 获取实际内容
	content := ""
	if docContent != nil {
		content = docContent.Content
	}

	chapterInfo := &ai.ChapterInfo{
		ID:           chapterID,
		Title:        doc.Title,
		Summary:      summary,
		Content:      content,
		CharacterIDs: doc.CharacterIDs,                     // 使用Document中的CharacterIDs字段
		LocationIDs:  doc.LocationIDs,                      // 使用Document中的LocationIDs字段
		TimelineIDs:  doc.TimelineIDs,                      // 使用Document中的TimelineIDs字段
		PlotThreads:  doc.PlotThreads,                      // 使用Document中的PlotThreads字段
		KeyPoints:    doc.KeyPoints,                        // 使用Document中的KeyPoints字段
		WritingHints: strings.Join(doc.WritingHints, "\n"), // 将字符串数组转换为单个字符串
	}

	return chapterInfo, nil
}

// buildPreviousChaptersSummary 构建前面章节的摘要
func (s *ContextService) buildPreviousChaptersSummary(ctx context.Context, projectID, currentChapterID string) (string, error) {
	// 暂时返回空字符串，需要实现章节顺序和摘要逻辑
	// TODO: 实现获取项目中当前章节之前的所有章节，并生成摘要
	return "", nil
}

// generateChapterSummary 生成章节摘要
// 注意: 此方法依赖 documentContentRepo，如果为 nil 则降级使用 KeyPoints
// 建议调用方直接使用 buildChapterInfo，它已经包含了摘要生成逻辑
func (s *ContextService) generateChapterSummary(ctx context.Context, doc *writer.Document) string {
	// 注意: documentContentRepo 可能为 nil（旧架构遗留），需要防御性检查
	if s.documentContentRepo == nil {
		// 降级方案：从 KeyPoints 生成摘要
		if len(doc.KeyPoints) > 0 {
			return strings.Join(doc.KeyPoints, "; ")
		}
		return ""
	}

	// 通过 DocumentContentRepository 获取内容
	docContent, err := s.documentContentRepo.GetByDocumentID(ctx, doc.ID.Hex())
	if err != nil || docContent == nil {
		// 降级方案：从 KeyPoints 生成摘要
		if len(doc.KeyPoints) > 0 {
			return strings.Join(doc.KeyPoints, "; ")
		}
		return ""
	}

	// 从内容中提取前200字符作为摘要
	content := strings.TrimSpace(docContent.Content)
	if len(content) > 200 {
		return content[:200] + "..."
	}
	return content
}

// BuildContextWithOptions 根据选项构建AI上下文
func (s *ContextService) BuildContextWithOptions(ctx context.Context, projectID string, chapterID string, options *ai.ContextOptions) (*ai.AIContext, error) {
	// 构建基础上下文
	aiContext, err := s.BuildContext(ctx, projectID, chapterID)
	if err != nil {
		return nil, err
	}

	// 根据选项调整上下文内容
	if options != nil {
		// TODO: 根据options调整上下文内容
		// 例如：限制token数量、包含历史章节、包含大纲等
	}

	return aiContext, nil
}

// UpdateContextWithFeedback 根据反馈更新上下文
func (s *ContextService) UpdateContextWithFeedback(ctx context.Context, aiContext *ai.AIContext, feedback string) error {
	// TODO: 根据反馈更新上下文
	// 例如：调整角色状态、更新情节线索等
	return nil
}
