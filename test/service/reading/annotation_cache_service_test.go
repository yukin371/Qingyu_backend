package reading

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"

	"Qingyu_backend/models/reading/reader"
	"Qingyu_backend/service/reading"
)

// setupAnnotationCacheMock 创建AnnotationCacheService和Redis Mock
func setupAnnotationCacheMock() (*reading.AnnotationCacheService, redismock.ClientMock) {
	db, mock := redismock.NewClientMock()
	service := reading.NewAnnotationCacheService(db)
	return service, mock
}

// createTestAnnotation 创建测试用标注
func createTestAnnotation(id, userID, bookID string) *reader.Annotation {
	return &reader.Annotation{
		ID:        id,
		UserID:    userID,
		BookID:    bookID,
		ChapterID: "chapter1",
		Range:     "0-100",
		Text:      "标注文本",
		Note:      "我的笔记",
		Type:      "bookmark",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// createTestAnnotationStats 创建测试用标注统计
func createTestAnnotationStats() *reading.AnnotationStats {
	return &reading.AnnotationStats{
		TotalCount:     100,
		BookmarkCount:  30,
		HighlightCount: 50,
		NoteCount:      20,
		LastUpdated:    time.Now(),
	}
}

// TestAnnotationCacheService_GetAnnotation_Hit 测试获取单个标注缓存命中
func TestAnnotationCacheService_GetAnnotation_Hit(t *testing.T) {
	service, mock := setupAnnotationCacheMock()
	ctx := context.Background()

	annotation := createTestAnnotation("ann123", "user123", "book123")
	annotationJSON, _ := json.Marshal(annotation)

	// 设置Mock期望
	mock.ExpectGet("annotation:ann123").SetVal(string(annotationJSON))

	// 执行测试
	result, err := service.GetAnnotation(ctx, "ann123")

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, annotation.ID, result.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestAnnotationCacheService_GetAnnotation_Miss 测试获取单个标注缓存未命中
func TestAnnotationCacheService_GetAnnotation_Miss(t *testing.T) {
	service, mock := setupAnnotationCacheMock()
	ctx := context.Background()

	// 设置Mock期望：缓存未命中
	mock.ExpectGet("annotation:ann123").RedisNil()

	// 执行测试
	result, err := service.GetAnnotation(ctx, "ann123")

	// 验证结果
	assert.NoError(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestAnnotationCacheService_GetAnnotation_Error 测试Redis错误
func TestAnnotationCacheService_GetAnnotation_Error(t *testing.T) {
	service, mock := setupAnnotationCacheMock()
	ctx := context.Background()

	// 设置Mock期望：Redis返回错误
	mock.ExpectGet("annotation:ann123").SetErr(fmt.Errorf("redis error"))

	// 执行测试
	result, err := service.GetAnnotation(ctx, "ann123")

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestAnnotationCacheService_SetAnnotation 测试设置单个标注缓存
func TestAnnotationCacheService_SetAnnotation(t *testing.T) {
	service, mock := setupAnnotationCacheMock()
	ctx := context.Background()

	annotation := createTestAnnotation("ann123", "user123", "book123")
	annotationJSON, _ := json.Marshal(annotation)

	// 设置Mock期望
	mock.ExpectSet("annotation:ann123", annotationJSON, 30*time.Minute).SetVal("OK")

	// 执行测试
	err := service.SetAnnotation(ctx, annotation)

	// 验证结果
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestAnnotationCacheService_DeleteAnnotation 测试删除标注缓存
func TestAnnotationCacheService_DeleteAnnotation(t *testing.T) {
	service, mock := setupAnnotationCacheMock()
	ctx := context.Background()

	// 设置Mock期望
	mock.ExpectDel("annotation:ann123").SetVal(1)

	// 执行测试
	err := service.DeleteAnnotation(ctx, "ann123")

	// 验证结果
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestAnnotationCacheService_GetChapterAnnotations_Hit 测试获取章节标注列表缓存命中
func TestAnnotationCacheService_GetChapterAnnotations_Hit(t *testing.T) {
	service, mock := setupAnnotationCacheMock()
	ctx := context.Background()

	annotations := []*reader.Annotation{
		createTestAnnotation("ann1", "user123", "book123"),
		createTestAnnotation("ann2", "user123", "book123"),
	}
	annotationsJSON, _ := json.Marshal(annotations)

	// 设置Mock期望
	mock.ExpectGet("annotation:user:user123:book:book123:chapter:chapter1:list").SetVal(string(annotationsJSON))

	// 执行测试
	result, err := service.GetChapterAnnotations(ctx, "user123", "book123", "chapter1")

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result))
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestAnnotationCacheService_GetChapterAnnotations_Miss 测试获取章节标注列表缓存未命中
func TestAnnotationCacheService_GetChapterAnnotations_Miss(t *testing.T) {
	service, mock := setupAnnotationCacheMock()
	ctx := context.Background()

	// 设置Mock期望：缓存未命中
	mock.ExpectGet("annotation:user:user123:book:book123:chapter:chapter1:list").RedisNil()

	// 执行测试
	result, err := service.GetChapterAnnotations(ctx, "user123", "book123", "chapter1")

	// 验证结果
	assert.NoError(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestAnnotationCacheService_SetChapterAnnotations 测试设置章节标注列表缓存
func TestAnnotationCacheService_SetChapterAnnotations(t *testing.T) {
	service, mock := setupAnnotationCacheMock()
	ctx := context.Background()

	annotations := []*reader.Annotation{
		createTestAnnotation("ann1", "user123", "book123"),
	}
	annotationsJSON, _ := json.Marshal(annotations)

	// 设置Mock期望
	mock.ExpectSet("annotation:user:user123:book:book123:chapter:chapter1:list", annotationsJSON, 30*time.Minute).SetVal("OK")

	// 执行测试
	err := service.SetChapterAnnotations(ctx, "user123", "book123", "chapter1", annotations)

	// 验证结果
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestAnnotationCacheService_InvalidateChapterAnnotations 测试使章节标注缓存失效
func TestAnnotationCacheService_InvalidateChapterAnnotations(t *testing.T) {
	service, mock := setupAnnotationCacheMock()
	ctx := context.Background()

	// 设置Mock期望
	mock.ExpectDel("annotation:user:user123:book:book123:chapter:chapter1:list").SetVal(1)

	// 执行测试
	err := service.InvalidateChapterAnnotations(ctx, "user123", "book123", "chapter1")

	// 验证结果
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestAnnotationCacheService_GetBookAnnotations_Hit 测试获取书籍标注列表缓存命中
func TestAnnotationCacheService_GetBookAnnotations_Hit(t *testing.T) {
	service, mock := setupAnnotationCacheMock()
	ctx := context.Background()

	annotations := []*reader.Annotation{
		createTestAnnotation("ann1", "user123", "book123"),
	}
	annotationsJSON, _ := json.Marshal(annotations)

	// 设置Mock期望
	mock.ExpectGet("annotation:user:user123:book:book123:list").SetVal(string(annotationsJSON))

	// 执行测试
	result, err := service.GetBookAnnotations(ctx, "user123", "book123")

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestAnnotationCacheService_SetBookAnnotations 测试设置书籍标注列表缓存
func TestAnnotationCacheService_SetBookAnnotations(t *testing.T) {
	service, mock := setupAnnotationCacheMock()
	ctx := context.Background()

	annotations := []*reader.Annotation{
		createTestAnnotation("ann1", "user123", "book123"),
	}
	annotationsJSON, _ := json.Marshal(annotations)

	// 设置Mock期望
	mock.ExpectSet("annotation:user:user123:book:book123:list", annotationsJSON, 30*time.Minute).SetVal("OK")

	// 执行测试
	err := service.SetBookAnnotations(ctx, "user123", "book123", annotations)

	// 验证结果
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestAnnotationCacheService_InvalidateBookAnnotations 测试使书籍标注缓存失效
func TestAnnotationCacheService_InvalidateBookAnnotations(t *testing.T) {
	service, mock := setupAnnotationCacheMock()
	ctx := context.Background()

	// 设置Mock期望
	mock.ExpectDel("annotation:user:user123:book:book123:list").SetVal(1)

	// 执行测试
	err := service.InvalidateBookAnnotations(ctx, "user123", "book123")

	// 验证结果
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestAnnotationCacheService_InvalidateUserAnnotations 测试使用户所有标注缓存失效
func TestAnnotationCacheService_InvalidateUserAnnotations(t *testing.T) {
	service, mock := setupAnnotationCacheMock()
	ctx := context.Background()

	keys := []string{
		"annotation:user:user123:book:book1:list",
		"annotation:user:user123:book:book2:list",
	}

	// 设置Mock期望：使用SCAN命令
	mock.ExpectScan(0, "annotation:user:user123:*", 0).SetVal(keys, 0)
	// 期望删除每个找到的键
	for _, key := range keys {
		mock.ExpectDel(key).SetVal(1)
	}

	// 执行测试
	err := service.InvalidateUserAnnotations(ctx, "user123")

	// 验证结果
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestAnnotationCacheService_BatchGetAnnotations 测试批量获取标注
func TestAnnotationCacheService_BatchGetAnnotations(t *testing.T) {
	service, mock := setupAnnotationCacheMock()
	ctx := context.Background()

	ann1 := createTestAnnotation("ann1", "user123", "book123")
	ann2 := createTestAnnotation("ann2", "user123", "book123")
	ann1JSON, _ := json.Marshal(ann1)
	ann2JSON, _ := json.Marshal(ann2)

	// 设置Mock期望：MGET批量获取
	mock.ExpectMGet("annotation:ann1", "annotation:ann2").SetVal([]interface{}{string(ann1JSON), string(ann2JSON)})

	// 执行测试
	result, err := service.BatchGetAnnotations(ctx, []string{"ann1", "ann2"})

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result))
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestAnnotationCacheService_BatchGetAnnotations_Empty 测试批量获取空列表
func TestAnnotationCacheService_BatchGetAnnotations_Empty(t *testing.T) {
	service, _ := setupAnnotationCacheMock()
	ctx := context.Background()

	// 执行测试：传入空列表
	result, err := service.BatchGetAnnotations(ctx, []string{})

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result))
}

// TestAnnotationCacheService_BatchGetAnnotations_PartialMiss 测试批量获取部分缺失
func TestAnnotationCacheService_BatchGetAnnotations_PartialMiss(t *testing.T) {
	service, mock := setupAnnotationCacheMock()
	ctx := context.Background()

	ann1 := createTestAnnotation("ann1", "user123", "book123")
	ann1JSON, _ := json.Marshal(ann1)

	// 设置Mock期望：第一个存在，第二个不存在
	mock.ExpectMGet("annotation:ann1", "annotation:ann2").SetVal([]interface{}{string(ann1JSON), nil})

	// 执行测试
	result, err := service.BatchGetAnnotations(ctx, []string{"ann1", "ann2"})

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result)) // 只返回存在的
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestAnnotationCacheService_BatchSetAnnotations 测试批量设置标注缓存
func TestAnnotationCacheService_BatchSetAnnotations(t *testing.T) {
	service, mock := setupAnnotationCacheMock()
	ctx := context.Background()

	annotations := []*reader.Annotation{
		createTestAnnotation("ann1", "user123", "book123"),
		createTestAnnotation("ann2", "user123", "book123"),
	}

	// 设置Mock期望：Pipeline批量设置
	for _, ann := range annotations {
		annJSON, _ := json.Marshal(ann)
		mock.ExpectSet(fmt.Sprintf("annotation:%s", ann.ID), annJSON, 30*time.Minute).SetVal("OK")
	}

	// 执行测试
	err := service.BatchSetAnnotations(ctx, annotations)

	// 验证结果
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestAnnotationCacheService_BatchSetAnnotations_Empty 测试批量设置空列表
func TestAnnotationCacheService_BatchSetAnnotations_Empty(t *testing.T) {
	service, _ := setupAnnotationCacheMock()
	ctx := context.Background()

	// 执行测试：传入空列表
	err := service.BatchSetAnnotations(ctx, []*reader.Annotation{})

	// 验证结果
	assert.NoError(t, err)
}

// TestAnnotationCacheService_GetAnnotationStats_Hit 测试获取标注统计缓存命中
func TestAnnotationCacheService_GetAnnotationStats_Hit(t *testing.T) {
	service, mock := setupAnnotationCacheMock()
	ctx := context.Background()

	stats := createTestAnnotationStats()
	statsJSON, _ := json.Marshal(stats)

	// 设置Mock期望
	mock.ExpectGet("annotation:user:user123:book:book123:stats").SetVal(string(statsJSON))

	// 执行测试
	result, err := service.GetAnnotationStats(ctx, "user123", "book123")

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, stats.TotalCount, result.TotalCount)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestAnnotationCacheService_GetAnnotationStats_Miss 测试获取标注统计缓存未命中
func TestAnnotationCacheService_GetAnnotationStats_Miss(t *testing.T) {
	service, mock := setupAnnotationCacheMock()
	ctx := context.Background()

	// 设置Mock期望：缓存未命中
	mock.ExpectGet("annotation:user:user123:book:book123:stats").RedisNil()

	// 执行测试
	result, err := service.GetAnnotationStats(ctx, "user123", "book123")

	// 验证结果
	assert.NoError(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestAnnotationCacheService_SetAnnotationStats 测试设置标注统计缓存
func TestAnnotationCacheService_SetAnnotationStats(t *testing.T) {
	service, mock := setupAnnotationCacheMock()
	ctx := context.Background()

	stats := createTestAnnotationStats()
	statsJSON, _ := json.Marshal(stats)

	// 设置Mock期望
	mock.ExpectSet("annotation:user:user123:book:book123:stats", statsJSON, 30*time.Minute).SetVal("OK")

	// 执行测试
	err := service.SetAnnotationStats(ctx, "user123", "book123", stats)

	// 验证结果
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestAnnotationCacheService_GetAnnotation_InvalidJSON 测试JSON解析失败
func TestAnnotationCacheService_GetAnnotation_InvalidJSON(t *testing.T) {
	service, mock := setupAnnotationCacheMock()
	ctx := context.Background()

	// 设置Mock期望：返回无效JSON
	mock.ExpectGet("annotation:ann123").SetVal("invalid json")

	// 执行测试
	result, err := service.GetAnnotation(ctx, "ann123")

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "解析注记数据失败")
	assert.NoError(t, mock.ExpectationsWereMet())
}
