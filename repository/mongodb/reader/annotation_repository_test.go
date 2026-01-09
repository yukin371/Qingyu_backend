package reading

import (
	"Qingyu_backend/models/reader"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"

	"Qingyu_backend/global"
	readerRepo "Qingyu_backend/repository/mongodb/reader"
	"Qingyu_backend/test/testutil"
)

var annotationRepo *readerRepo.MongoAnnotationRepository

func setupAnnotationTest(t *testing.T) {
	testutil.SetupTestDB(t)
	annotationRepo = readerRepo.NewMongoAnnotationRepository(global.DB)

	// 清理测试数据
	ctx := context.Background()
	_ = global.DB.Collection("annotations").Drop(ctx)
}

func createTestAnnotation(userID, bookID, chapterID string, annotationType string) *reader.Annotation {
	ann := &reader.Annotation{
		UserID:    userID,
		BookID:    bookID,
		ChapterID: chapterID,
		Range:     "0-100",
		Text:      "Test annotation text",
		Note:      "Test note",
		Type:      annotationType,
	}
	return ann
}

// 计数器确保ID唯一性
var annIDCounter int64

// createAndInsertAnnotation 创建标注并插入数据库
func createAndInsertAnnotation(ctx context.Context, userID, bookID, chapterID string, annotationType string) (*reader.Annotation, error) {
	ann := createTestAnnotation(userID, bookID, chapterID, annotationType)

	// 生成唯一ID（使用纳秒时间+计数器）
	if ann.ID == "" {
		annIDCounter++
		ann.ID = fmt.Sprintf("ann_%d_%d", time.Now().UnixNano(), annIDCounter)
	}
	ann.CreatedAt = time.Now()
	ann.UpdatedAt = time.Now()

	// 直接插入
	_, err := global.DB.Collection("annotations").InsertOne(ctx, ann)
	if err != nil {
		return nil, err
	}

	return ann, nil
}

// ============ 基础CRUD操作测试 ============

func TestAnnotationRepository_Create(t *testing.T) {
	setupAnnotationTest(t)
	ctx := context.Background()

	annotation := createTestAnnotation("user1", "book1", "chapter1", string(reader.AnnotationTypeNote))

	err := annotationRepo.Create(ctx, annotation)
	require.NoError(t, err)
	assert.NotEmpty(t, annotation.ID)
	assert.NotZero(t, annotation.CreatedAt)
	assert.NotZero(t, annotation.UpdatedAt)
}

func TestAnnotationRepository_GetByID(t *testing.T) {
	setupAnnotationTest(t)
	ctx := context.Background()

	// 创建标注
	annotation := createTestAnnotation("user1", "book1", "chapter1", string(reader.AnnotationTypeNote))
	err := annotationRepo.Create(ctx, annotation)
	require.NoError(t, err)

	// 查询标注
	result, err := annotationRepo.GetByID(ctx, annotation.ID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, annotation.UserID, result.UserID)
	assert.Equal(t, annotation.BookID, result.BookID)
	assert.Equal(t, annotation.ChapterID, result.ChapterID)
}

func TestAnnotationRepository_GetByID_NotFound(t *testing.T) {
	setupAnnotationTest(t)
	ctx := context.Background()

	result, err := annotationRepo.GetByID(ctx, "nonexistent")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestAnnotationRepository_Update(t *testing.T) {
	setupAnnotationTest(t)
	ctx := context.Background()

	// 创建标注
	annotation := createTestAnnotation("user1", "book1", "chapter1", string(reader.AnnotationTypeNote))
	err := annotationRepo.Create(ctx, annotation)
	require.NoError(t, err)

	// 更新标注
	updates := map[string]interface{}{
		"note": "Updated note",
		"text": "Updated text",
	}
	err = annotationRepo.Update(ctx, annotation.ID, updates)
	require.NoError(t, err)

	// 验证更新
	result, err := annotationRepo.GetByID(ctx, annotation.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated note", result.Note)
	assert.Equal(t, "Updated text", result.Text)
}

func TestAnnotationRepository_Delete(t *testing.T) {
	setupAnnotationTest(t)
	ctx := context.Background()

	// 创建标注
	annotation := createTestAnnotation("user1", "book1", "chapter1", string(reader.AnnotationTypeNote))
	err := annotationRepo.Create(ctx, annotation)
	require.NoError(t, err)

	// 删除标注
	err = annotationRepo.Delete(ctx, annotation.ID)
	require.NoError(t, err)

	// 验证已删除
	result, err := annotationRepo.GetByID(ctx, annotation.ID)
	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============ 查询操作测试 ============

func TestAnnotationRepository_GetByUserAndBook(t *testing.T) {
	setupAnnotationTest(t)
	ctx := context.Background()

	// 创建多个标注
	ann1 := createTestAnnotation("user1", "book1", "chapter1", string(reader.AnnotationTypeNote))
	ann2 := createTestAnnotation("user1", "book1", "chapter2", string(reader.AnnotationTypeBookmark))
	ann3 := createTestAnnotation("user1", "book2", "chapter1", string(reader.AnnotationTypeNote))
	ann4 := createTestAnnotation("user2", "book1", "chapter1", string(reader.AnnotationTypeNote))

	err := annotationRepo.Create(ctx, ann1)
	require.NoError(t, err)
	err = annotationRepo.Create(ctx, ann2)
	require.NoError(t, err)
	err = annotationRepo.Create(ctx, ann3)
	require.NoError(t, err)
	err = annotationRepo.Create(ctx, ann4)
	require.NoError(t, err)

	// 查询user1在book1的标注
	results, err := annotationRepo.GetByUserAndBook(ctx, "user1", "book1")
	require.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestAnnotationRepository_GetByUserAndChapter(t *testing.T) {
	setupAnnotationTest(t)
	ctx := context.Background()

	// 创建标注
	ann1 := createTestAnnotation("user1", "book1", "chapter1", string(reader.AnnotationTypeNote))
	ann2 := createTestAnnotation("user1", "book1", "chapter1", string(reader.AnnotationTypeBookmark))
	ann3 := createTestAnnotation("user1", "book1", "chapter2", string(reader.AnnotationTypeNote))

	err := annotationRepo.Create(ctx, ann1)
	require.NoError(t, err)
	err = annotationRepo.Create(ctx, ann2)
	require.NoError(t, err)
	err = annotationRepo.Create(ctx, ann3)
	require.NoError(t, err)

	// 查询user1在book1/chapter1的标注
	results, err := annotationRepo.GetByUserAndChapter(ctx, "user1", "book1", "chapter1")
	require.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestAnnotationRepository_GetByType(t *testing.T) {
	setupAnnotationTest(t)
	ctx := context.Background()

	// 创建不同类型的标注（使用直接插入以支持int类型）
	_, err := createAndInsertAnnotation(ctx, "user1", "book1", "chapter1", string(reader.AnnotationTypeNote)) // note
	require.NoError(t, err)
	_, err = createAndInsertAnnotation(ctx, "user1", "book1", "chapter1", string(reader.AnnotationTypeBookmark)) // bookmark
	require.NoError(t, err)
	_, err = createAndInsertAnnotation(ctx, "user1", "book1", "chapter2", string(reader.AnnotationTypeNote)) // note
	require.NoError(t, err)

	// 查询type=note的标注
	results, err := annotationRepo.GetByType(ctx, "user1", "book1", string(reader.AnnotationTypeNote))
	require.NoError(t, err)
	assert.Len(t, results, 2)
}

// ============ 笔记操作测试 ============

func TestAnnotationRepository_GetNotes(t *testing.T) {
	setupAnnotationTest(t)
	ctx := context.Background()

	// 创建标注（使用直接插入）
	_, err := createAndInsertAnnotation(ctx, "user1", "book1", "chapter1", string(reader.AnnotationTypeNote)) // note
	require.NoError(t, err)
	_, err = createAndInsertAnnotation(ctx, "user1", "book1", "chapter1", string(reader.AnnotationTypeBookmark)) // bookmark
	require.NoError(t, err)

	// 查询笔记
	results, err := annotationRepo.GetNotes(ctx, "user1", "book1")
	require.NoError(t, err)
	assert.Len(t, results, 1)
}

func TestAnnotationRepository_GetNotesByChapter(t *testing.T) {
	setupAnnotationTest(t)
	ctx := context.Background()

	// 创建标注
	ann1, err := createAndInsertAnnotation(ctx, "user1", "book1", "chapter1", string(reader.AnnotationTypeNote))
	require.NoError(t, err)
	_, err = createAndInsertAnnotation(ctx, "user1", "book1", "chapter2", string(reader.AnnotationTypeNote))
	require.NoError(t, err)

	// 查询chapter1的笔记
	results, err := annotationRepo.GetNotesByChapter(ctx, "user1", "book1", "chapter1")
	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, ann1.ChapterID, results[0].ChapterID)
}

func TestAnnotationRepository_SearchNotes(t *testing.T) {
	setupAnnotationTest(t)
	ctx := context.Background()

	// 创建标注（使用bson.M以便设置int type和自定义note）
	annIDCounter++
	doc1 := bson.M{
		"_id":        fmt.Sprintf("ann_%d_%d", time.Now().UnixNano(), annIDCounter),
		"user_id":    "user1",
		"book_id":    "book1",
		"chapter_id": "chapter1",
		"text":       "Test text",
		"note":       "Important content",
		"type":       string(reader.AnnotationTypeNote), // 使用字符串类型
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}
	_, err := global.DB.Collection("annotations").InsertOne(ctx, doc1)
	require.NoError(t, err)

	annIDCounter++
	doc2 := bson.M{
		"_id":        fmt.Sprintf("ann_%d_%d", time.Now().UnixNano(), annIDCounter),
		"user_id":    "user1",
		"book_id":    "book1",
		"chapter_id": "chapter2",
		"text":       "Test text",
		"note":       "Regular note",
		"type":       string(reader.AnnotationTypeNote), // 使用字符串类型
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}
	_, err = global.DB.Collection("annotations").InsertOne(ctx, doc2)
	require.NoError(t, err)

	// 搜索包含"Important"的笔记
	results, err := annotationRepo.SearchNotes(ctx, "user1", "Important")
	require.NoError(t, err)
	require.Len(t, results, 1, "应该找到1条包含Important的笔记")
	if len(results) > 0 {
		assert.Contains(t, results[0].Note, "Important")
	}
}

// ============ 书签操作测试 ============

func TestAnnotationRepository_GetBookmarks(t *testing.T) {
	setupAnnotationTest(t)
	ctx := context.Background()

	// 创建标注
	_, err := createAndInsertAnnotation(ctx, "user1", "book1", "chapter1", string(reader.AnnotationTypeNote)) // note
	require.NoError(t, err)
	_, err = createAndInsertAnnotation(ctx, "user1", "book1", "chapter1", string(reader.AnnotationTypeBookmark)) // bookmark
	require.NoError(t, err)

	// 查询书签
	results, err := annotationRepo.GetBookmarks(ctx, "user1", "book1")
	require.NoError(t, err)
	assert.Len(t, results, 1)
}

func TestAnnotationRepository_GetLatestBookmark(t *testing.T) {
	setupAnnotationTest(t)
	ctx := context.Background()

	// 创建多个书签
	_, err := createAndInsertAnnotation(ctx, "user1", "book1", "chapter1", string(reader.AnnotationTypeBookmark))
	require.NoError(t, err)
	time.Sleep(1 * time.Millisecond)
	_, err = createAndInsertAnnotation(ctx, "user1", "book1", "chapter2", string(reader.AnnotationTypeBookmark))
	require.NoError(t, err)

	// 查询最新书签
	result, err := annotationRepo.GetLatestBookmark(ctx, "user1", "book1")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestAnnotationRepository_GetLatestBookmark_NotFound(t *testing.T) {
	setupAnnotationTest(t)
	ctx := context.Background()

	result, err := annotationRepo.GetLatestBookmark(ctx, "user1", "book1")
	require.NoError(t, err)
	assert.Nil(t, result) // 不存在时返回nil
}

// ============ 高亮操作测试 ============

func TestAnnotationRepository_GetHighlights(t *testing.T) {
	setupAnnotationTest(t)
	ctx := context.Background()

	// 创建标注
	_, err := createAndInsertAnnotation(ctx, "user1", "book1", "chapter1", string(reader.AnnotationTypeHighlight)) // highlight
	require.NoError(t, err)
	_, err = createAndInsertAnnotation(ctx, "user1", "book1", "chapter1", string(reader.AnnotationTypeNote)) // note
	require.NoError(t, err)

	// 查询高亮
	results, err := annotationRepo.GetHighlights(ctx, "user1", "book1")
	require.NoError(t, err)
	assert.Len(t, results, 1)
}

func TestAnnotationRepository_GetHighlightsByChapter(t *testing.T) {
	setupAnnotationTest(t)
	ctx := context.Background()

	// 创建标注
	ann1, err := createAndInsertAnnotation(ctx, "user1", "book1", "chapter1", string(reader.AnnotationTypeHighlight))
	require.NoError(t, err)
	_, err = createAndInsertAnnotation(ctx, "user1", "book1", "chapter2", string(reader.AnnotationTypeHighlight))
	require.NoError(t, err)

	// 查询chapter1的高亮
	results, err := annotationRepo.GetHighlightsByChapter(ctx, "user1", "book1", "chapter1")
	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, ann1.ChapterID, results[0].ChapterID)
}

// ============ 统计操作测试 ============

func TestAnnotationRepository_CountByUser(t *testing.T) {
	setupAnnotationTest(t)
	ctx := context.Background()

	// 创建标注
	ann1 := createTestAnnotation("user1", "book1", "chapter1", string(reader.AnnotationTypeNote))
	ann2 := createTestAnnotation("user1", "book2", "chapter1", string(reader.AnnotationTypeNote))
	ann3 := createTestAnnotation("user2", "book1", "chapter1", string(reader.AnnotationTypeNote))

	err := annotationRepo.Create(ctx, ann1)
	require.NoError(t, err)
	err = annotationRepo.Create(ctx, ann2)
	require.NoError(t, err)
	err = annotationRepo.Create(ctx, ann3)
	require.NoError(t, err)

	// 统计user1的标注数
	count, err := annotationRepo.CountByUser(ctx, "user1")
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

func TestAnnotationRepository_CountByBook(t *testing.T) {
	setupAnnotationTest(t)
	ctx := context.Background()

	// 创建标注
	ann1 := createTestAnnotation("user1", "book1", "chapter1", string(reader.AnnotationTypeNote))
	ann2 := createTestAnnotation("user1", "book1", "chapter2", string(reader.AnnotationTypeNote))
	ann3 := createTestAnnotation("user1", "book2", "chapter1", string(reader.AnnotationTypeNote))

	err := annotationRepo.Create(ctx, ann1)
	require.NoError(t, err)
	err = annotationRepo.Create(ctx, ann2)
	require.NoError(t, err)
	err = annotationRepo.Create(ctx, ann3)
	require.NoError(t, err)

	// 统计user1在book1的标注数
	count, err := annotationRepo.CountByBook(ctx, "user1", "book1")
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

func TestAnnotationRepository_CountByType(t *testing.T) {
	setupAnnotationTest(t)
	ctx := context.Background()

	// 创建不同类型的标注
	_, err := createAndInsertAnnotation(ctx, "user1", "book1", "chapter1", string(reader.AnnotationTypeNote)) // note
	require.NoError(t, err)
	_, err = createAndInsertAnnotation(ctx, "user1", "book1", "chapter1", string(reader.AnnotationTypeBookmark)) // bookmark
	require.NoError(t, err)
	_, err = createAndInsertAnnotation(ctx, "user1", "book2", "chapter1", string(reader.AnnotationTypeNote)) // note
	require.NoError(t, err)

	// 统计user1的note数量
	count, err := annotationRepo.CountByType(ctx, "user1", string(reader.AnnotationTypeNote))
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

// ============ 批量操作测试 ============

func TestAnnotationRepository_BatchCreate(t *testing.T) {
	setupAnnotationTest(t)
	ctx := context.Background()

	// 准备批量数据
	annotations := []*reader.Annotation{
		createTestAnnotation("user1", "book1", "chapter1", string(reader.AnnotationTypeNote)),
		createTestAnnotation("user1", "book1", "chapter2", string(reader.AnnotationTypeBookmark)),
		createTestAnnotation("user1", "book1", "chapter3", string(reader.AnnotationTypeHighlight)),
	}

	// 批量创建
	err := annotationRepo.BatchCreate(ctx, annotations)
	require.NoError(t, err)

	// 验证创建成功
	count, err := annotationRepo.CountByUser(ctx, "user1")
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

func TestAnnotationRepository_BatchCreate_Empty(t *testing.T) {
	setupAnnotationTest(t)
	ctx := context.Background()

	// 空数组应该正常返回
	err := annotationRepo.BatchCreate(ctx, []*reader.Annotation{})
	require.NoError(t, err)
}

func TestAnnotationRepository_BatchDelete(t *testing.T) {
	setupAnnotationTest(t)
	ctx := context.Background()

	// 创建标注
	ann1 := createTestAnnotation("user1", "book1", "chapter1", string(reader.AnnotationTypeNote))
	ann2 := createTestAnnotation("user1", "book1", "chapter2", string(reader.AnnotationTypeBookmark))
	ann3 := createTestAnnotation("user1", "book1", "chapter3", string(reader.AnnotationTypeHighlight))

	err := annotationRepo.Create(ctx, ann1)
	require.NoError(t, err)
	err = annotationRepo.Create(ctx, ann2)
	require.NoError(t, err)
	err = annotationRepo.Create(ctx, ann3)
	require.NoError(t, err)

	// 批量删除
	ids := []string{ann1.ID, ann2.ID}
	err = annotationRepo.BatchDelete(ctx, ids)
	require.NoError(t, err)

	// 验证删除结果
	count, err := annotationRepo.CountByUser(ctx, "user1")
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestAnnotationRepository_DeleteByBook(t *testing.T) {
	setupAnnotationTest(t)
	ctx := context.Background()

	// 创建标注
	ann1 := createTestAnnotation("user1", "book1", "chapter1", string(reader.AnnotationTypeNote))
	ann2 := createTestAnnotation("user1", "book1", "chapter2", string(reader.AnnotationTypeBookmark))
	ann3 := createTestAnnotation("user1", "book2", "chapter1", string(reader.AnnotationTypeNote))

	err := annotationRepo.Create(ctx, ann1)
	require.NoError(t, err)
	err = annotationRepo.Create(ctx, ann2)
	require.NoError(t, err)
	err = annotationRepo.Create(ctx, ann3)
	require.NoError(t, err)

	// 删除book1的所有标注
	err = annotationRepo.DeleteByBook(ctx, "user1", "book1")
	require.NoError(t, err)

	// 验证删除结果
	count, err := annotationRepo.CountByBook(ctx, "user1", "book1")
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)

	count, err = annotationRepo.CountByBook(ctx, "user1", "book2")
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestAnnotationRepository_DeleteByChapter(t *testing.T) {
	setupAnnotationTest(t)
	ctx := context.Background()

	// 创建标注
	ann1 := createTestAnnotation("user1", "book1", "chapter1", string(reader.AnnotationTypeNote))
	ann2 := createTestAnnotation("user1", "book1", "chapter2", string(reader.AnnotationTypeBookmark))

	err := annotationRepo.Create(ctx, ann1)
	require.NoError(t, err)
	err = annotationRepo.Create(ctx, ann2)
	require.NoError(t, err)

	// 删除chapter1的标注
	err = annotationRepo.DeleteByChapter(ctx, "user1", "book1", "chapter1")
	require.NoError(t, err)

	// 验证删除结果
	results, err := annotationRepo.GetByUserAndChapter(ctx, "user1", "book1", "chapter1")
	require.NoError(t, err)
	assert.Len(t, results, 0)

	results, err = annotationRepo.GetByUserAndChapter(ctx, "user1", "book1", "chapter2")
	require.NoError(t, err)
	assert.Len(t, results, 1)
}

// ============ 数据同步测试 ============

func TestAnnotationRepository_SyncAnnotations(t *testing.T) {
	setupAnnotationTest(t)
	ctx := context.Background()

	// 准备同步数据
	annotations := []*reader.Annotation{
		{
			ID:        "sync_ann_1",
			UserID:    "user1",
			BookID:    "book1",
			ChapterID: "chapter1",
			Type:      "1",
			Text:      "Synced annotation",
		},
	}

	// 同步
	err := annotationRepo.SyncAnnotations(ctx, "user1", annotations)
	require.NoError(t, err)

	// 验证
	result, err := annotationRepo.GetByID(ctx, "sync_ann_1")
	require.NoError(t, err)
	assert.Equal(t, "Synced annotation", result.Text)
}

func TestAnnotationRepository_GetRecentAnnotations(t *testing.T) {
	setupAnnotationTest(t)
	ctx := context.Background()

	// 创建标注
	for i := 0; i < 5; i++ {
		ann := createTestAnnotation("user1", "book1", "chapter1", string(reader.AnnotationTypeNote))
		err := annotationRepo.Create(ctx, ann)
		require.NoError(t, err)
	}

	// 获取最近3条
	results, err := annotationRepo.GetRecentAnnotations(ctx, "user1", 3)
	require.NoError(t, err)
	assert.Len(t, results, 3)
}

// ============ 健康检查测试 ============

func TestAnnotationRepository_Health(t *testing.T) {
	setupAnnotationTest(t)
	ctx := context.Background()

	err := annotationRepo.Health(ctx)
	assert.NoError(t, err)
}
