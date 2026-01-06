package reading

import (
	reader2 "Qingyu_backend/models/reader"
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"

	"Qingyu_backend/service/reader"
)

// setupReaderCacheMock 创建ReaderCacheService和Redis Mock
func setupReaderCacheMock(prefix string) (reading.ReaderCacheService, redismock.ClientMock) {
	db, mock := redismock.NewClientMock()
	service := reading.NewRedisReaderCacheService(db, prefix)
	return service, mock
}

// createTestChapter 创建测试用章节
func createTestChapter(chapterID string) *reader2.Chapter {
	return &reader2.Chapter{
		ID:      chapterID,
		BookID:  "book123",
		Title:   "第一章",
		Content: "章节内容...",
		IsVIP:   false,
	}
}

// createTestProgress 创建测试用阅读进度
func createTestProgress(userID, bookID string) *reader2.ReadingProgress {
	return &reader2.ReadingProgress{
		ID:          "progress123",
		UserID:      userID,
		BookID:      bookID,
		ChapterID:   "chapter1",
		Progress:    0.5,
		ReadingTime: 3600,
		LastReadAt:  time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// createTestReadingSettings 创建测试用阅读设置
func createTestReadingSettings(userID string) *reader2.ReadingSettings {
	return &reader2.ReadingSettings{
		ID:          "settings123",
		UserID:      userID,
		FontFamily:  "Arial",
		FontSize:    16,
		LineHeight:  1.5,
		Theme:       "light",
		Background:  "#FFFFFF",
		PageMode:    1,
		AutoScroll:  false,
		ScrollSpeed: 50,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// TestReaderCacheService_GetChapterContent_Hit 测试章节内容缓存命中
func TestReaderCacheService_GetChapterContent_Hit(t *testing.T) {
	service, mock := setupReaderCacheMock("qingyu")
	ctx := context.Background()

	expectedContent := "这是章节内容..."

	// 设置Mock期望
	mock.ExpectGet("qingyu:reader:chapter_content:chapter1").SetVal(expectedContent)

	// 执行测试
	content, err := service.GetChapterContent(ctx, "chapter1")

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, expectedContent, content)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestReaderCacheService_GetChapterContent_Miss 测试章节内容缓存未命中
func TestReaderCacheService_GetChapterContent_Miss(t *testing.T) {
	service, mock := setupReaderCacheMock("qingyu")
	ctx := context.Background()

	// 设置Mock期望：缓存未命中
	mock.ExpectGet("qingyu:reader:chapter_content:chapter1").RedisNil()

	// 执行测试
	content, err := service.GetChapterContent(ctx, "chapter1")

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, "", content)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestReaderCacheService_GetChapterContent_Error 测试Redis错误
func TestReaderCacheService_GetChapterContent_Error(t *testing.T) {
	service, mock := setupReaderCacheMock("qingyu")
	ctx := context.Background()

	// 设置Mock期望：Redis返回错误
	mock.ExpectGet("qingyu:reader:chapter_content:chapter1").SetErr(fmt.Errorf("redis error"))

	// 执行测试
	content, err := service.GetChapterContent(ctx, "chapter1")

	// 验证结果
	assert.Error(t, err)
	assert.Equal(t, "", content)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestReaderCacheService_SetChapterContent 测试设置章节内容缓存
func TestReaderCacheService_SetChapterContent(t *testing.T) {
	service, mock := setupReaderCacheMock("qingyu")
	ctx := context.Background()

	content := "章节内容..."
	expiration := 1 * time.Hour

	// 设置Mock期望
	mock.ExpectSet("qingyu:reader:chapter_content:chapter1", content, expiration).SetVal("OK")

	// 执行测试
	err := service.SetChapterContent(ctx, "chapter1", content, expiration)

	// 验证结果
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestReaderCacheService_InvalidateChapterContent 测试清除章节内容缓存
func TestReaderCacheService_InvalidateChapterContent(t *testing.T) {
	service, mock := setupReaderCacheMock("qingyu")
	ctx := context.Background()

	// 设置Mock期望
	mock.ExpectDel("qingyu:reader:chapter_content:chapter1").SetVal(1)

	// 执行测试
	err := service.InvalidateChapterContent(ctx, "chapter1")

	// 验证结果
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestReaderCacheService_GetChapter_Hit 测试章节信息缓存命中
func TestReaderCacheService_GetChapter_Hit(t *testing.T) {
	service, mock := setupReaderCacheMock("qingyu")
	ctx := context.Background()

	chapter := createTestChapter("chapter1")
	chapterJSON, _ := json.Marshal(chapter)

	// 设置Mock期望
	mock.ExpectGet("qingyu:reader:chapter:chapter1").SetVal(string(chapterJSON))

	// 执行测试
	result, err := service.GetChapter(ctx, "chapter1")

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, chapter.ID, result.ID)
	assert.Equal(t, chapter.Title, result.Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestReaderCacheService_GetChapter_Miss 测试章节信息缓存未命中
func TestReaderCacheService_GetChapter_Miss(t *testing.T) {
	service, mock := setupReaderCacheMock("qingyu")
	ctx := context.Background()

	// 设置Mock期望：缓存未命中
	mock.ExpectGet("qingyu:reader:chapter:chapter1").RedisNil()

	// 执行测试
	result, err := service.GetChapter(ctx, "chapter1")

	// 验证结果
	assert.NoError(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestReaderCacheService_GetChapter_JSONError 测试JSON解析错误
func TestReaderCacheService_GetChapter_JSONError(t *testing.T) {
	service, mock := setupReaderCacheMock("qingyu")
	ctx := context.Background()

	// 设置Mock期望：返回无效JSON
	mock.ExpectGet("qingyu:reader:chapter:chapter1").SetVal("invalid json")

	// 执行测试
	result, err := service.GetChapter(ctx, "chapter1")

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to unmarshal")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestReaderCacheService_SetChapter 测试设置章节信息缓存
func TestReaderCacheService_SetChapter(t *testing.T) {
	service, mock := setupReaderCacheMock("qingyu")
	ctx := context.Background()

	chapter := createTestChapter("chapter1")
	chapterJSON, _ := json.Marshal(chapter)
	expiration := 2 * time.Hour

	// 设置Mock期望
	mock.ExpectSet("qingyu:reader:chapter:chapter1", chapterJSON, expiration).SetVal("OK")

	// 执行测试
	err := service.SetChapter(ctx, "chapter1", chapter, expiration)

	// 验证结果
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestReaderCacheService_InvalidateChapter 测试清除章节信息缓存
func TestReaderCacheService_InvalidateChapter(t *testing.T) {
	service, mock := setupReaderCacheMock("qingyu")
	ctx := context.Background()

	// 设置Mock期望
	mock.ExpectDel("qingyu:reader:chapter:chapter1").SetVal(1)

	// 执行测试
	err := service.InvalidateChapter(ctx, "chapter1")

	// 验证结果
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestReaderCacheService_GetReadingSettings_Hit 测试阅读设置缓存命中
func TestReaderCacheService_GetReadingSettings_Hit(t *testing.T) {
	service, mock := setupReaderCacheMock("qingyu")
	ctx := context.Background()

	settings := createTestReadingSettings("user123")
	settingsJSON, _ := json.Marshal(settings)

	// 设置Mock期望
	mock.ExpectGet("qingyu:reader:settings:user123").SetVal(string(settingsJSON))

	// 执行测试
	result, err := service.GetReadingSettings(ctx, "user123")

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, settings.UserID, result.UserID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestReaderCacheService_GetReadingSettings_Miss 测试阅读设置缓存未命中
func TestReaderCacheService_GetReadingSettings_Miss(t *testing.T) {
	service, mock := setupReaderCacheMock("qingyu")
	ctx := context.Background()

	// 设置Mock期望：缓存未命中
	mock.ExpectGet("qingyu:reader:settings:user123").RedisNil()

	// 执行测试
	result, err := service.GetReadingSettings(ctx, "user123")

	// 验证结果
	assert.NoError(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestReaderCacheService_SetReadingSettings 测试设置阅读设置缓存
func TestReaderCacheService_SetReadingSettings(t *testing.T) {
	service, mock := setupReaderCacheMock("qingyu")
	ctx := context.Background()

	settings := createTestReadingSettings("user123")
	settingsJSON, _ := json.Marshal(settings)
	expiration := 24 * time.Hour

	// 设置Mock期望
	mock.ExpectSet("qingyu:reader:settings:user123", settingsJSON, expiration).SetVal("OK")

	// 执行测试
	err := service.SetReadingSettings(ctx, "user123", settings, expiration)

	// 验证结果
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestReaderCacheService_InvalidateReadingSettings 测试清除阅读设置缓存
func TestReaderCacheService_InvalidateReadingSettings(t *testing.T) {
	service, mock := setupReaderCacheMock("qingyu")
	ctx := context.Background()

	// 设置Mock期望
	mock.ExpectDel("qingyu:reader:settings:user123").SetVal(1)

	// 执行测试
	err := service.InvalidateReadingSettings(ctx, "user123")

	// 验证结果
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestReaderCacheService_GetReadingProgress_Hit 测试阅读进度缓存命中
func TestReaderCacheService_GetReadingProgress_Hit(t *testing.T) {
	service, mock := setupReaderCacheMock("qingyu")
	ctx := context.Background()

	progress := createTestProgress("user123", "book123")
	progressJSON, _ := json.Marshal(progress)

	// 设置Mock期望
	mock.ExpectGet("qingyu:reader:progress:user123:book123").SetVal(string(progressJSON))

	// 执行测试
	result, err := service.GetReadingProgress(ctx, "user123", "book123")

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, progress.UserID, result.UserID)
	assert.Equal(t, progress.BookID, result.BookID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestReaderCacheService_GetReadingProgress_Miss 测试阅读进度缓存未命中
func TestReaderCacheService_GetReadingProgress_Miss(t *testing.T) {
	service, mock := setupReaderCacheMock("qingyu")
	ctx := context.Background()

	// 设置Mock期望：缓存未命中
	mock.ExpectGet("qingyu:reader:progress:user123:book123").RedisNil()

	// 执行测试
	result, err := service.GetReadingProgress(ctx, "user123", "book123")

	// 验证结果
	assert.NoError(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestReaderCacheService_SetReadingProgress 测试设置阅读进度缓存
func TestReaderCacheService_SetReadingProgress(t *testing.T) {
	service, mock := setupReaderCacheMock("qingyu")
	ctx := context.Background()

	progress := createTestProgress("user123", "book123")
	progressJSON, _ := json.Marshal(progress)
	expiration := 30 * time.Minute

	// 设置Mock期望
	mock.ExpectSet("qingyu:reader:progress:user123:book123", progressJSON, expiration).SetVal("OK")

	// 执行测试
	err := service.SetReadingProgress(ctx, "user123", "book123", progress, expiration)

	// 验证结果
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestReaderCacheService_InvalidateReadingProgress 测试清除阅读进度缓存
func TestReaderCacheService_InvalidateReadingProgress(t *testing.T) {
	service, mock := setupReaderCacheMock("qingyu")
	ctx := context.Background()

	// 设置Mock期望
	mock.ExpectDel("qingyu:reader:progress:user123:book123").SetVal(1)

	// 执行测试
	err := service.InvalidateReadingProgress(ctx, "user123", "book123")

	// 验证结果
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestReaderCacheService_InvalidateBookChapters 测试批量清除书籍章节缓存
func TestReaderCacheService_InvalidateBookChapters(t *testing.T) {
	service, mock := setupReaderCacheMock("qingyu")
	ctx := context.Background()

	keys := []string{
		"qingyu:reader:chapter:book123:chapter1",
		"qingyu:reader:chapter:book123:chapter2",
	}

	// 设置Mock期望
	mock.ExpectKeys("qingyu:reader:chapter*:book123:*").SetVal(keys)
	mock.ExpectDel(keys...).SetVal(int64(len(keys)))

	// 执行测试
	err := service.InvalidateBookChapters(ctx, "book123")

	// 验证结果
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestReaderCacheService_InvalidateBookChapters_NoKeys 测试批量清除时无缓存
func TestReaderCacheService_InvalidateBookChapters_NoKeys(t *testing.T) {
	service, mock := setupReaderCacheMock("qingyu")
	ctx := context.Background()

	// 设置Mock期望：没有找到缓存键
	mock.ExpectKeys("qingyu:reader:chapter*:book123:*").SetVal([]string{})

	// 执行测试
	err := service.InvalidateBookChapters(ctx, "book123")

	// 验证结果
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestReaderCacheService_InvalidateUserData 测试清除用户所有缓存
func TestReaderCacheService_InvalidateUserData(t *testing.T) {
	service, mock := setupReaderCacheMock("qingyu")
	ctx := context.Background()

	progressKeys := []string{
		"qingyu:reader:progress:user123:book1",
		"qingyu:reader:progress:user123:book2",
	}

	// 设置Mock期望
	mock.ExpectDel("qingyu:reader:settings:user123").SetVal(1)
	mock.ExpectKeys("qingyu:reader:progress:user123:*").SetVal(progressKeys)
	mock.ExpectDel(progressKeys...).SetVal(int64(len(progressKeys)))

	// 执行测试
	err := service.InvalidateUserData(ctx, "user123")

	// 验证结果
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestReaderCacheService_DefaultPrefix 测试默认前缀
func TestReaderCacheService_DefaultPrefix(t *testing.T) {
	db, mock := redismock.NewClientMock()
	service := reading.NewRedisReaderCacheService(db, "")
	ctx := context.Background()

	// 设置Mock期望：使用默认前缀"qingyu"
	mock.ExpectGet("qingyu:reader:chapter_content:chapter1").SetVal("content")

	// 执行测试
	content, err := service.GetChapterContent(ctx, "chapter1")

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, "content", content)
	assert.NoError(t, mock.ExpectationsWereMet())
}
