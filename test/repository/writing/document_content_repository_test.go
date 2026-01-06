package writing_test

import (
	"Qingyu_backend/models/writer"
	writingInterface "Qingyu_backend/repository/interfaces/writer"
	writerRepo "Qingyu_backend/repository/mongodb/writer"
	"Qingyu_backend/test/testutil"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 测试辅助函数
func setupDocumentContentRepo(t *testing.T) (writingInterface.DocumentContentRepository, context.Context, func()) {
	db, cleanup := testutil.SetupTestDB(t)
	repo := writerRepo.NewMongoDocumentContentRepository(db)
	ctx := context.Background()
	return repo, ctx, cleanup
}

func createTestDocumentContent(documentID, content string) *writer.DocumentContent {
	return &writer.DocumentContent{
		DocumentID:  documentID,
		Content:     content,
		ContentType: "markdown",
		WordCount:   len([]rune(content)),
		CharCount:   len(content),
		Version:     1,
	}
}

// 1. 测试创建文档内容
func TestDocumentContentRepository_Create(t *testing.T) {
	repo, ctx, cleanup := setupDocumentContentRepo(t)
	defer cleanup()

	docContent := createTestDocumentContent("doc123", "这是测试内容")

	err := repo.Create(ctx, docContent)
	require.NoError(t, err)
	assert.NotEmpty(t, docContent.ID)
	assert.NotZero(t, docContent.CreatedAt)
	assert.Equal(t, "markdown", docContent.ContentType)
}

// 2. 测试根据ID获取
func TestDocumentContentRepository_GetByID(t *testing.T) {
	repo, ctx, cleanup := setupDocumentContentRepo(t)
	defer cleanup()

	// 创建
	docContent := createTestDocumentContent("doc123", "测试内容")
	err := repo.Create(ctx, docContent)
	require.NoError(t, err)

	// 获取
	retrieved, err := repo.GetByID(ctx, docContent.ID)
	require.NoError(t, err)
	require.NotNil(t, retrieved)
	assert.Equal(t, docContent.ID, retrieved.ID)
	assert.Equal(t, docContent.Content, retrieved.Content)
}

// 3. 测试根据DocumentID获取
func TestDocumentContentRepository_GetByDocumentID(t *testing.T) {
	repo, ctx, cleanup := setupDocumentContentRepo(t)
	defer cleanup()

	documentID := primitive.NewObjectID().Hex()
	docContent := createTestDocumentContent(documentID, "通过DocumentID查询")
	err := repo.Create(ctx, docContent)
	require.NoError(t, err)

	// 通过DocumentID获取
	retrieved, err := repo.GetByDocumentID(ctx, documentID)
	require.NoError(t, err)
	require.NotNil(t, retrieved)
	assert.Equal(t, documentID, retrieved.DocumentID)
	assert.Equal(t, "通过DocumentID查询", retrieved.Content)
}

// 4. 测试更新
func TestDocumentContentRepository_Update(t *testing.T) {
	repo, ctx, cleanup := setupDocumentContentRepo(t)
	defer cleanup()

	// 创建
	docContent := createTestDocumentContent("doc123", "原始内容")
	err := repo.Create(ctx, docContent)
	require.NoError(t, err)

	// 更新
	updates := map[string]interface{}{
		"content":    "更新后的内容",
		"word_count": 7,
	}
	err = repo.Update(ctx, docContent.ID, updates)
	require.NoError(t, err)

	// 验证更新
	retrieved, err := repo.GetByID(ctx, docContent.ID)
	require.NoError(t, err)
	assert.Equal(t, "更新后的内容", retrieved.Content)
	assert.Equal(t, 7, retrieved.WordCount)
}

// 5. 测试删除
func TestDocumentContentRepository_Delete(t *testing.T) {
	repo, ctx, cleanup := setupDocumentContentRepo(t)
	defer cleanup()

	// 创建
	docContent := createTestDocumentContent("doc123", "待删除内容")
	err := repo.Create(ctx, docContent)
	require.NoError(t, err)

	// 删除
	err = repo.Delete(ctx, docContent.ID)
	require.NoError(t, err)

	// 验证已删除
	retrieved, err := repo.GetByID(ctx, docContent.ID)
	require.NoError(t, err)
	assert.Nil(t, retrieved)
}

// 6. 测试带版本号的更新（乐观锁）
func TestDocumentContentRepository_UpdateWithVersion(t *testing.T) {
	t.Skip("UpdateWithVersion 需要MongoDB实现支持")

	repo, ctx, cleanup := setupDocumentContentRepo(t)
	defer cleanup()

	// 创建
	docContent := createTestDocumentContent("doc123", "版本控制内容")
	err := repo.Create(ctx, docContent)
	require.NoError(t, err)

	// 使用正确的版本号更新
	err = repo.UpdateWithVersion(ctx, docContent.DocumentID, "更新内容v2", 1)
	assert.NoError(t, err)

	// 使用错误的版本号更新（应该失败）
	err = repo.UpdateWithVersion(ctx, docContent.DocumentID, "更新内容v3", 1)
	assert.Error(t, err)
}

// 7. 测试获取内容统计
func TestDocumentContentRepository_GetContentStats(t *testing.T) {
	t.Skip("GetContentStats 需要MongoDB实现支持")

	repo, ctx, cleanup := setupDocumentContentRepo(t)
	defer cleanup()

	// 创建
	content := "这是一段测试文本，包含中文和English。"
	docContent := createTestDocumentContent("doc123", content)
	err := repo.Create(ctx, docContent)
	require.NoError(t, err)

	// 获取统计
	wordCount, charCount, err := repo.GetContentStats(ctx, docContent.DocumentID)
	require.NoError(t, err)
	assert.Greater(t, wordCount, 0)
	assert.Greater(t, charCount, 0)
}

// 8. 测试健康检查
func TestDocumentContentRepository_Health(t *testing.T) {
	repo, ctx, cleanup := setupDocumentContentRepo(t)
	defer cleanup()

	err := repo.Health(ctx)
	assert.NoError(t, err)
}

// 9. 测试批量更新
func TestDocumentContentRepository_BatchUpdateContent(t *testing.T) {
	t.Skip("BatchUpdateContent 需要MongoDB实现支持")

	repo, ctx, cleanup := setupDocumentContentRepo(t)
	defer cleanup()

	// 创建多个文档内容
	doc1 := createTestDocumentContent("doc1", "内容1")
	doc2 := createTestDocumentContent("doc2", "内容2")
	err := repo.Create(ctx, doc1)
	require.NoError(t, err)
	err = repo.Create(ctx, doc2)
	require.NoError(t, err)

	// 批量更新
	updates := map[string]string{
		doc1.DocumentID: "批量更新内容1",
		doc2.DocumentID: "批量更新内容2",
	}
	err = repo.BatchUpdateContent(ctx, updates)
	assert.NoError(t, err)
}

// 10. 测试Exists
func TestDocumentContentRepository_Exists(t *testing.T) {
	repo, ctx, cleanup := setupDocumentContentRepo(t)
	defer cleanup()

	// 创建
	docContent := createTestDocumentContent("doc123", "存在检查")
	err := repo.Create(ctx, docContent)
	require.NoError(t, err)

	// 检查存在
	exists, err := repo.Exists(ctx, docContent.ID)
	require.NoError(t, err)
	assert.True(t, exists)

	// 检查不存在
	exists, err = repo.Exists(ctx, primitive.NewObjectID().Hex())
	require.NoError(t, err)
	assert.False(t, exists)
}
