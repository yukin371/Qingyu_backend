package ai

import (
	"context"
	"fmt"
	"sort"
	"time"

	"Qingyu_backend/models/ai"
	documentService "Qingyu_backend/service/project"
)

// NovelContextService 小说上下文管理服务
type NovelContextService struct {
	documentService *documentService.DocumentService
	projectService  *documentService.ProjectService
	nodeService     *documentService.NodeService
	versionService  *documentService.VersionService
	vectorDB        VectorDatabase
	memoryStore     MemoryStore
	retriever       RetrievalService
	summarizer      SummaryService
}

// VectorDatabase 向量数据库接口
type VectorDatabase interface {
	Store(ctx context.Context, id string, embedding []float32, metadata map[string]interface{}) error
	SimilaritySearch(ctx context.Context, query []float32, limit int) ([]*ai.RetrievalResult, error)
	Delete(ctx context.Context, id string) error
}

// MemoryStore 记忆存储接口
type MemoryStore interface {
	StoreMemory(ctx context.Context, memory *ai.ContextMemory) error
	GetMemory(ctx context.Context, projectID, memoryType string) ([]*ai.ContextMemory, error)
	UpdateMemoryAccess(ctx context.Context, memoryID string) error
	DeleteMemory(ctx context.Context, memoryID string) error
}

// RetrievalService 检索服务接口
type RetrievalService interface {
	RetrieveRelevant(ctx context.Context, projectID, query string, limit int) ([]*ai.RetrievalResult, error)
	HybridSearch(ctx context.Context, projectID, query string, options *SearchOptions) ([]*ai.RetrievalResult, error)
}

// SummaryService 摘要服务接口
type SummaryService interface {
	SummarizeContent(ctx context.Context, content string, maxTokens int) (string, error)
	CompressContext(ctx context.Context, contexts []*ai.NovelContext, targetTokens int) ([]*ai.NovelContext, error)
}

// SearchOptions 搜索选项
type SearchOptions struct {
	VectorWeight   float32 `json:"vectorWeight"`
	KeywordWeight  float32 `json:"keywordWeight"`
	MetadataWeight float32 `json:"metadataWeight"`
	MaxResults     int     `json:"maxResults"`
	MinScore       float32 `json:"minScore"`
}

// NewNovelContextService 创建小说上下文管理服务
func NewNovelContextService(
	documentService *documentService.DocumentService,
	projectService *documentService.ProjectService,
	nodeService *documentService.NodeService,
	versionService *documentService.VersionService,
	vectorDB VectorDatabase,
	memoryStore MemoryStore,
	retriever RetrievalService,
	summarizer SummaryService,
) *NovelContextService {
	return &NovelContextService{
		documentService: documentService,
		projectService:  projectService,
		nodeService:     nodeService,
		versionService:  versionService,
		vectorDB:        vectorDB,
		memoryStore:     memoryStore,
		retriever:       retriever,
		summarizer:      summarizer,
	}
}

// BuildContext 构建动态上下文
func (s *NovelContextService) BuildContext(ctx context.Context, req *ai.ContextBuildRequest) (*ai.ContextBuildResponse, error) {
	startTime := time.Now()

	// 1. 获取短期记忆（当前章节和最近章节）
	shortTermMemory, err := s.getShortTermMemory(ctx, req.ProjectID, req.CurrentPosition)
	if err != nil {
		return nil, fmt.Errorf("获取短期记忆失败: %w", err)
	}

	// 2. 检索相关的中长期记忆
	relevantMemory, err := s.retriever.RetrieveRelevant(ctx, req.ProjectID, req.CurrentPosition, req.MaxTokens/2)
	if err != nil {
		return nil, fmt.Errorf("检索相关记忆失败: %w", err)
	}

	// 3. 构建动态上下文
	aiContext, tokenCount, err := s.buildDynamicContext(ctx, shortTermMemory, relevantMemory, req.MaxTokens)
	if err != nil {
		return nil, fmt.Errorf("构建动态上下文失败: %w", err)
	}

	buildTime := time.Since(startTime)

	return &ai.ContextBuildResponse{
		Context:       aiContext,
		TokenCount:    tokenCount,
		Sources:       relevantMemory,
		BuildStrategy: "hybrid_retrieval",
		BuildTime:     buildTime,
	}, nil
}

// getShortTermMemory 获取短期记忆
func (s *NovelContextService) getShortTermMemory(ctx context.Context, projectID, currentPosition string) ([]*ai.ContextMemory, error) {
	memories, err := s.memoryStore.GetMemory(ctx, projectID, "short_term")
	if err != nil {
		return nil, err
	}

	// 按重要性和时间排序
	sort.Slice(memories, func(i, j int) bool {
		if memories[i].Importance != memories[j].Importance {
			return memories[i].Importance > memories[j].Importance
		}
		return memories[i].CreatedAt.After(memories[j].CreatedAt)
	})

	// 限制数量，避免过多
	if len(memories) > 10 {
		memories = memories[:10]
	}

	return memories, nil
}

// buildDynamicContext 构建动态上下文
func (s *NovelContextService) buildDynamicContext(
	ctx context.Context,
	shortTerm []*ai.ContextMemory,
	relevant []*ai.RetrievalResult,
	maxTokens int,
) (*ai.AIContext, int, error) {
	aiContext := &ai.AIContext{
		ActiveCharacters: []*ai.CharacterInfo{},
		CurrentLocations: []*ai.LocationInfo{},
		RelevantEvents:   []*ai.TimelineEvent{},
		PreviousChapters: []*ai.ChapterSummary{},
		PlotThreads:      []*ai.PlotThread{},
	}

	tokenCount := 0

	// 处理短期记忆（100%权重）
	for _, memory := range shortTerm {
		if tokenCount+len(memory.Content)/4 > maxTokens { // 粗略估算token数
			break
		}
		// 根据记忆类型添加到相应的上下文中
		s.addMemoryToContext(aiContext, memory)
		tokenCount += len(memory.Content) / 4
	}

	// 处理检索到的相关内容（按权重排序）
	sort.Slice(relevant, func(i, j int) bool {
		return relevant[i].Score > relevant[j].Score
	})

	for _, result := range relevant {
		if tokenCount+len(result.Context.Content)/4 > maxTokens {
			break
		}
		s.addContextToAIContext(aiContext, result.Context)
		tokenCount += len(result.Context.Content) / 4
	}

	aiContext.TokenCount = tokenCount
	return aiContext, tokenCount, nil
}

// addMemoryToContext 将记忆添加到AI上下文中
func (s *NovelContextService) addMemoryToContext(aiContext *ai.AIContext, memory *ai.ContextMemory) {
	// 根据记忆类型和内容，将其添加到相应的上下文字段中
	// 这里需要根据具体的记忆内容进行解析和分类
	// 简化实现，实际应该有更复杂的解析逻辑
}

// addContextToAIContext 将NovelContext添加到AIContext中
func (s *NovelContextService) addContextToAIContext(aiContext *ai.AIContext, context *ai.NovelContext) {
	switch context.Type {
	case "character":
		// 解析角色信息并添加
	case "plot":
		// 解析情节信息并添加
	case "setting":
		// 解析设定信息并添加
	case "chapter":
		// 解析章节信息并添加
	}
}

// StoreContext 存储上下文信息
func (s *NovelContextService) StoreContext(ctx context.Context, novelContext *ai.NovelContext) error {
	// 设置创建时间
	novelContext.BeforeCreate()

	// 生成嵌入向量（这里需要调用嵌入模型）
	embedding, err := s.generateEmbedding(ctx, novelContext.Content)
	if err != nil {
		return fmt.Errorf("生成嵌入向量失败: %w", err)
	}
	novelContext.Embedding = embedding

	// 存储到向量数据库
	metadata := map[string]interface{}{
		"project_id": novelContext.ProjectID,
		"type":       novelContext.Type,
		"title":      novelContext.Title,
		"importance": novelContext.Importance,
		"created_at": novelContext.CreatedAt,
	}

	err = s.vectorDB.Store(ctx, novelContext.ID, embedding, metadata)
	if err != nil {
		return fmt.Errorf("存储到向量数据库失败: %w", err)
	}

	// 这里还需要存储到MongoDB等持久化存储
	// 实际实现中需要添加数据库操作

	return nil
}

// UpdateContext 更新上下文信息
func (s *NovelContextService) UpdateContext(ctx context.Context, novelContext *ai.NovelContext) error {
	novelContext.BeforeUpdate()

	// 重新生成嵌入向量
	embedding, err := s.generateEmbedding(ctx, novelContext.Content)
	if err != nil {
		return fmt.Errorf("生成嵌入向量失败: %w", err)
	}
	novelContext.Embedding = embedding

	// 更新向量数据库
	metadata := map[string]interface{}{
		"project_id": novelContext.ProjectID,
		"type":       novelContext.Type,
		"title":      novelContext.Title,
		"importance": novelContext.Importance,
		"updated_at": novelContext.UpdatedAt,
	}

	err = s.vectorDB.Store(ctx, novelContext.ID, embedding, metadata)
	if err != nil {
		return fmt.Errorf("更新向量数据库失败: %w", err)
	}

	return nil
}

// DeleteContext 删除上下文信息
func (s *NovelContextService) DeleteContext(ctx context.Context, contextID string) error {
	err := s.vectorDB.Delete(ctx, contextID)
	if err != nil {
		return fmt.Errorf("从向量数据库删除失败: %w", err)
	}

	// 这里还需要从MongoDB等持久化存储中删除
	// 实际实现中需要添加数据库操作

	return nil
}

// generateEmbedding 生成嵌入向量（占位实现）
func (s *NovelContextService) generateEmbedding(ctx context.Context, content string) ([]float32, error) {
	// 这里需要调用实际的嵌入模型API
	// 比如OpenAI的text-embedding-ada-002或者本地的BGE模型
	// 暂时返回空向量作为占位
	return make([]float32, 1536), nil // OpenAI embedding维度为1536
}

// SearchContext 搜索上下文
func (s *NovelContextService) SearchContext(ctx context.Context, projectID, query string, options *SearchOptions) ([]*ai.RetrievalResult, error) {
	if options == nil {
		options = &SearchOptions{
			VectorWeight:   0.6,
			KeywordWeight:  0.3,
			MetadataWeight: 0.1,
			MaxResults:     20,
			MinScore:       0.1,
		}
	}

	return s.retriever.HybridSearch(ctx, projectID, query, options)
}