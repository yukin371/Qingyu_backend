package ai

import (
	"context"
	"testing"

	"Qingyu_backend/models/ai"
	"Qingyu_backend/models/writer"
	testMock "Qingyu_backend/service/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestBuildContext_WithEmptyContext 测试构建空上下文
func TestBuildContext_WithEmptyContext(t *testing.T) {
	ctx := context.Background()

	// 创建一个最小的ContextService实例
	service := &ContextService{
		documentService:     nil,
		projectService:      nil,
		nodeService:         nil,
		versionService:      nil,
		documentContentRepo: nil,
	}

	// 测试当projectService为nil时的行为
	_, err := service.BuildContext(ctx, "", "")

	// 应该返回错误，因为projectService是nil
	assert.Error(t, err)
}

// TestBuildContext_WithOptions 测试带选项的上下文构建
func TestBuildContext_WithOptions(t *testing.T) {
	ctx := context.Background()

	service := &ContextService{
		documentService:     nil,
		projectService:      nil,
		nodeService:         nil,
		versionService:      nil,
		documentContentRepo: nil,
	}

	options := &ai.ContextOptions{}

	// 测试nil projectService的情况
	_, err := service.BuildContextWithOptions(ctx, "", "", options)

	assert.Error(t, err)
}

// TestUpdateContextWithFeedback 测试更新上下文反馈
func TestUpdateContextWithFeedback(t *testing.T) {
	ctx := context.Background()
	service := &ContextService{}

	aiContext := &ai.AIContext{
		ProjectID: "test-project",
	}

	err := service.UpdateContextWithFeedback(ctx, aiContext, "这个上下文很好")

	assert.NoError(t, err)
}

// TestBuildPreviousChaptersSummary 测试构建前序章节摘要
func TestBuildPreviousChaptersSummary(t *testing.T) {
	service := &ContextService{}

	ctx := context.Background()
	summary, err := service.buildPreviousChaptersSummary(ctx, "project1", "chapter1")

	assert.NoError(t, err)
	assert.Equal(t, "", summary)
}

// TestGenerateChapterSummary_WithContentRepo 测试有ContentRepo时生成摘要
func TestGenerateChapterSummary_WithContentRepo(t *testing.T) {
	mockContentRepo := new(testMock.MockDocumentContentRepository)

	// 创建Document，使用嵌入的ID字段
	doc := &writer.Document{}
	doc.ID = "doc1"
	doc.KeyPoints = []string{"关键点1", "关键点2"}

	docContent := &writer.DocumentContent{
		DocumentID: "doc1",
		Content:    "这是一段很长的内容，用来测试摘要生成功能。这段文字应该超过200个字符，这样可以确保能够正确地截取前200个字符并添加省略号。我们需要更多的内容来达到这个目标，所以这里添加一些额外的文字。",
	}

	mockContentRepo.On("GetByDocumentID", mock.Anything, "doc1").Return(docContent, nil)
	mockContentRepo.On("Count", mock.Anything).Return(int64(1), nil)

	service := &ContextService{
		documentContentRepo: mockContentRepo,
	}

	ctx := context.Background()
	summary := service.generateChapterSummary(ctx, doc)

	assert.NotEmpty(t, summary)
	assert.Contains(t, summary, "...")

	mockContentRepo.AssertExpectations(t)
}

// TestGenerateChapterSummary_WithoutContentRepo 测试没有ContentRepo时使用KeyPoints
func TestGenerateChapterSummary_WithoutContentRepo(t *testing.T) {
	doc := &writer.Document{}
	doc.ID = "doc1"
	doc.KeyPoints = []string{"关键点1", "关键点2", "关键点3"}

	service := &ContextService{
		documentContentRepo: nil,
	}

	ctx := context.Background()
	summary := service.generateChapterSummary(ctx, doc)

	assert.Equal(t, "关键点1; 关键点2; 关键点3", summary)
}

// TestGenerateChapterSummary_ContentRepoError 测试ContentRepo错误时降级到KeyPoints
func TestGenerateChapterSummary_ContentRepoError(t *testing.T) {
	mockContentRepo := new(testMock.MockDocumentContentRepository)

	doc := &writer.Document{}
	doc.ID = "doc1"
	doc.KeyPoints = []string{"关键点1", "关键点2"}

	mockContentRepo.On("GetByDocumentID", mock.Anything, "doc1").Return(nil, assert.AnError)
	mockContentRepo.On("Count", mock.Anything).Return(int64(0), nil)

	service := &ContextService{
		documentContentRepo: mockContentRepo,
	}

	ctx := context.Background()
	summary := service.generateChapterSummary(ctx, doc)

	assert.Equal(t, "关键点1; 关键点2", summary)

	mockContentRepo.AssertExpectations(t)
}

// TestGenerateChapterSummary_ShortContent 测试短内容生成摘要
func TestGenerateChapterSummary_ShortContent(t *testing.T) {
	mockContentRepo := new(testMock.MockDocumentContentRepository)

	doc := &writer.Document{}
	doc.ID = "doc1"
	doc.KeyPoints = []string{"关键点1", "关键点2"}

	docContent := &writer.DocumentContent{
		DocumentID: "doc1",
		Content:    "这是一段短内容",
	}

	mockContentRepo.On("GetByDocumentID", mock.Anything, "doc1").Return(docContent, nil)
	mockContentRepo.On("Count", mock.Anything).Return(int64(1), nil)

	service := &ContextService{
		documentContentRepo: mockContentRepo,
	}

	ctx := context.Background()
	summary := service.generateChapterSummary(ctx, doc)

	assert.Equal(t, "这是一段短内容", summary)

	mockContentRepo.AssertExpectations(t)
}

// TestGenerateChapterSummary_EmptyKeyPoints 测试没有KeyPoints的情况
func TestGenerateChapterSummary_EmptyKeyPoints(t *testing.T) {
	mockContentRepo := new(testMock.MockDocumentContentRepository)

	doc := &writer.Document{}
	doc.ID = "doc1"
	doc.KeyPoints = []string{}

	mockContentRepo.On("GetByDocumentID", mock.Anything, "doc1").Return(nil, assert.AnError)
	mockContentRepo.On("Count", mock.Anything).Return(int64(0), nil)

	service := &ContextService{
		documentContentRepo: mockContentRepo,
	}

	ctx := context.Background()
	summary := service.generateChapterSummary(ctx, doc)

	assert.Equal(t, "", summary)

	mockContentRepo.AssertExpectations(t)
}
