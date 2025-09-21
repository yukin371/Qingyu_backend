package ai

import (
	"context"
	"fmt"
	"strings"

	"Qingyu_backend/models/ai"
	"Qingyu_backend/models/document"
	documentService "Qingyu_backend/service/document"
)

// ContextService AI上下文服务
type ContextService struct {
	documentService *documentService.DocumentService
	projectService  *documentService.ProjectService
	nodeService     *documentService.NodeService
	versionService  *documentService.VersionService
}

// NewContextService 创建AI上下文服务
func NewContextService(
	documentService *documentService.DocumentService,
	projectService *documentService.ProjectService,
	nodeService *documentService.NodeService,
	versionService *documentService.VersionService,
) *ContextService {
	return &ContextService{
		documentService: documentService,
		projectService:  projectService,
		nodeService:     nodeService,
		versionService:  versionService,
	}
}

// BuildContext 构建AI上下文
func (s *ContextService) BuildContext(ctx context.Context, projectID string, chapterID string) (*ai.AIContext, error) {
	// 获取项目信息
	project, err := s.projectService.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("获取项目信息失败: %w", err)
	}
	if project == nil {
		return nil, fmt.Errorf("项目不存在")
	}

	// 构建章节信息
	chapterInfo, err := s.buildChapterInfo(ctx, projectID, chapterID)
	if err != nil {
		return nil, fmt.Errorf("构建章节信息失败: %w", err)
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
	// 获取章节文档
	doc, err := s.documentService.GetByID(chapterID)
	if err != nil {
		return nil, fmt.Errorf("获取章节文档失败: %w", err)
	}
	if doc == nil {
		return nil, fmt.Errorf("章节不存在")
	}

	chapterInfo := &ai.ChapterInfo{
		ID:           chapterID,
		Title:        doc.Title,
		Summary:      "", // Document模型没有Summary字段，可以从Content中提取或留空
		Content:      doc.Content,
		CharacterIDs: doc.CharacterIDs, // 使用Document中的CharacterIDs字段
		LocationIDs:  doc.LocationIDs,  // 使用Document中的LocationIDs字段
		TimelineIDs:  doc.TimelineIDs,  // 使用Document中的TimelineIDs字段
		PlotThreads:  doc.PlotThreads,  // 使用Document中的PlotThreads字段
		KeyPoints:    doc.KeyPoints,    // 使用Document中的KeyPoints字段
		WritingHints: doc.WritingHints, // 使用Document中的WritingHints字段
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
func (s *ContextService) generateChapterSummary(doc *document.Document) string {
	// 如果没有关键点，从内容中提取前200字符作为摘要
	content := strings.TrimSpace(doc.Content)
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
