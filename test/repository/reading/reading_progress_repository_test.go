package reading

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"Qingyu_backend/global"
	"Qingyu_backend/models/reading/reader"
	"Qingyu_backend/repository/mongodb/reading"
	"Qingyu_backend/test/testutil"
)

var repo *reading.MongoReadingProgressRepository

func setupTest(t *testing.T) {
	// 设置测试数据库
	testutil.SetupTestDB(t)

	// 创建Repository实例
	repo = reading.NewMongoReadingProgressRepository(global.DB)

	// 清理当前测试的数据
	ctx := context.Background()
	_ = global.DB.Collection("reading_progress").Drop(ctx)
}

func createTestProgress(userID, bookID string) *reader.ReadingProgress {
	return &reader.ReadingProgress{
		// 不设置ID，让Repository自动生成
		UserID:      userID,
		BookID:      bookID,
		ChapterID:   "chapter_1",
		Progress:    0.5,
		ReadingTime: 3600, // 1小时
		LastReadAt:  time.Now(),
	}
}

// ============ 基础CRUD操作测试 ============

func TestReadingProgressRepository_Create(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	progress := createTestProgress("user123", "book456")

	err := repo.Create(ctx, progress)
	require.NoError(t, err)
	assert.NotEmpty(t, progress.ID)
	assert.NotZero(t, progress.CreatedAt)
	assert.NotZero(t, progress.UpdatedAt)
}

func TestReadingProgressRepository_GetByID(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 先创建
	progress := createTestProgress("user123", "book456")
	err := repo.Create(ctx, progress)
	require.NoError(t, err)

	// 再查询
	result, err := repo.GetByID(ctx, progress.ID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, progress.UserID, result.UserID)
	assert.Equal(t, progress.BookID, result.BookID)
	assert.Equal(t, progress.ChapterID, result.ChapterID)
	assert.Equal(t, progress.Progress, result.Progress)
}

func TestReadingProgressRepository_GetByID_NotFound(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	result, err := repo.GetByID(ctx, "nonexistent")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestReadingProgressRepository_Update(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建进度
	progress := createTestProgress("user123", "book456")
	err := repo.Create(ctx, progress)
	require.NoError(t, err)

	// 更新进度
	updates := map[string]interface{}{
		"progress":   0.75,
		"chapter_id": "chapter_5",
	}
	err = repo.Update(ctx, progress.ID, updates)
	require.NoError(t, err)

	// 验证更新
	result, err := repo.GetByID(ctx, progress.ID)
	require.NoError(t, err)
	assert.Equal(t, 0.75, result.Progress)
	assert.Equal(t, "chapter_5", result.ChapterID)
}

func TestReadingProgressRepository_Delete(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建进度
	progress := createTestProgress("user123", "book456")
	err := repo.Create(ctx, progress)
	require.NoError(t, err)

	// 删除进度
	err = repo.Delete(ctx, progress.ID)
	require.NoError(t, err)

	// 验证已删除
	result, err := repo.GetByID(ctx, progress.ID)
	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============ 查询操作测试 ============

func TestReadingProgressRepository_GetByUserAndBook(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建进度
	progress := createTestProgress("user123", "book456")
	err := repo.Create(ctx, progress)
	require.NoError(t, err)

	// 查询用户对特定书籍的进度
	result, err := repo.GetByUserAndBook(ctx, "user123", "book456")
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "user123", result.UserID)
	assert.Equal(t, "book456", result.BookID)
}

func TestReadingProgressRepository_GetByUserAndBook_NotFound(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 查询不存在的记录
	result, err := repo.GetByUserAndBook(ctx, "user123", "book999")
	require.NoError(t, err)
	assert.Nil(t, result) // 没有记录返回nil而不是错误
}

func TestReadingProgressRepository_GetByUser(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建多个进度
	progress1 := createTestProgress("user123", "book1")
	progress2 := createTestProgress("user123", "book2")
	progress3 := createTestProgress("user456", "book3")

	err := repo.Create(ctx, progress1)
	require.NoError(t, err)
	time.Sleep(10 * time.Millisecond) // 确保时间不同

	err = repo.Create(ctx, progress2)
	require.NoError(t, err)

	err = repo.Create(ctx, progress3)
	require.NoError(t, err)

	// 查询user123的所有进度
	results, err := repo.GetByUser(ctx, "user123")
	require.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestReadingProgressRepository_GetRecentReadingByUser(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建多个进度
	for i := 0; i < 5; i++ {
		progress := createTestProgress("user123", "book"+string(rune('1'+i)))
		err := repo.Create(ctx, progress)
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond)
	}

	// 获取最近3条记录
	results, err := repo.GetRecentReadingByUser(ctx, "user123", 3)
	require.NoError(t, err)
	assert.Len(t, results, 3)
}

// ============ 进度保存和更新测试 ============

func TestReadingProgressRepository_SaveProgress(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 首次保存（upsert创建）
	err := repo.SaveProgress(ctx, "user123", "book456", "chapter_1", 0.3)
	require.NoError(t, err)

	// 验证创建成功
	progress, err := repo.GetByUserAndBook(ctx, "user123", "book456")
	require.NoError(t, err)
	assert.NotNil(t, progress)
	assert.Equal(t, "chapter_1", progress.ChapterID)
	assert.Equal(t, 0.3, progress.Progress)

	// 再次保存（upsert更新）
	err = repo.SaveProgress(ctx, "user123", "book456", "chapter_2", 0.6)
	require.NoError(t, err)

	// 验证更新成功
	progress, err = repo.GetByUserAndBook(ctx, "user123", "book456")
	require.NoError(t, err)
	assert.Equal(t, "chapter_2", progress.ChapterID)
	assert.Equal(t, 0.6, progress.Progress)
}

func TestReadingProgressRepository_UpdateReadingTime(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建进度
	progress := createTestProgress("user123", "book456")
	progress.ReadingTime = 1000
	err := repo.Create(ctx, progress)
	require.NoError(t, err)

	// 增加阅读时长
	err = repo.UpdateReadingTime(ctx, "user123", "book456", 500)
	require.NoError(t, err)

	// 验证时长增加
	result, err := repo.GetByUserAndBook(ctx, "user123", "book456")
	require.NoError(t, err)
	assert.Equal(t, int64(1500), result.ReadingTime)
}

func TestReadingProgressRepository_UpdateReadingTime_CreateIfNotExists(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 对不存在的记录更新时长（会创建新记录）
	err := repo.UpdateReadingTime(ctx, "user123", "book456", 300)
	require.NoError(t, err)

	// 验证创建成功
	progress, err := repo.GetByUserAndBook(ctx, "user123", "book456")
	require.NoError(t, err)
	assert.NotNil(t, progress)
	assert.Equal(t, int64(300), progress.ReadingTime)
}

func TestReadingProgressRepository_UpdateLastReadAt(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建进度
	progress := createTestProgress("user123", "book456")
	oldTime := time.Now().Add(-1 * time.Hour)
	progress.LastReadAt = oldTime
	err := repo.Create(ctx, progress)
	require.NoError(t, err)

	// 等待一下确保时间不同
	time.Sleep(10 * time.Millisecond)

	// 更新最后阅读时间
	err = repo.UpdateLastReadAt(ctx, "user123", "book456")
	require.NoError(t, err)

	// 验证时间更新
	result, err := repo.GetByUserAndBook(ctx, "user123", "book456")
	require.NoError(t, err)
	assert.True(t, result.LastReadAt.After(oldTime))
}

// ============ 批量操作测试 ============

func TestReadingProgressRepository_BatchUpdateProgress(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 准备多个进度记录（BatchUpdateProgress需要有ID）
	progresses := []*reader.ReadingProgress{
		{
			ID:          "batch_prog_1",
			UserID:      "user123",
			BookID:      "book1",
			ChapterID:   "chapter_1",
			Progress:    0.5,
			ReadingTime: 100,
			LastReadAt:  time.Now(),
		},
		{
			ID:          "batch_prog_2",
			UserID:      "user123",
			BookID:      "book2",
			ChapterID:   "chapter_2",
			Progress:    0.8,
			ReadingTime: 200,
			LastReadAt:  time.Now(),
		},
	}

	// 批量更新
	err := repo.BatchUpdateProgress(ctx, progresses)
	require.NoError(t, err)

	// 验证
	result1, err := repo.GetByUserAndBook(ctx, "user123", "book1")
	require.NoError(t, err)
	assert.Equal(t, 0.5, result1.Progress)

	result2, err := repo.GetByUserAndBook(ctx, "user123", "book2")
	require.NoError(t, err)
	assert.Equal(t, 0.8, result2.Progress)
}

func TestReadingProgressRepository_BatchUpdateProgress_Empty(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 空数组应该正常返回
	err := repo.BatchUpdateProgress(ctx, []*reader.ReadingProgress{})
	require.NoError(t, err)
}

// ============ 统计查询测试 ============

func TestReadingProgressRepository_GetTotalReadingTime(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建多个进度记录
	progress1 := createTestProgress("user123", "book1")
	progress1.ReadingTime = 1000
	progress2 := createTestProgress("user123", "book2")
	progress2.ReadingTime = 2000
	progress3 := createTestProgress("user456", "book3")
	progress3.ReadingTime = 3000

	err := repo.Create(ctx, progress1)
	require.NoError(t, err)
	err = repo.Create(ctx, progress2)
	require.NoError(t, err)
	err = repo.Create(ctx, progress3)
	require.NoError(t, err)

	// 查询user123的总阅读时长
	total, err := repo.GetTotalReadingTime(ctx, "user123")
	require.NoError(t, err)
	assert.Equal(t, int64(3000), total)
}

func TestReadingProgressRepository_GetTotalReadingTime_NoData(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 查询不存在的用户
	total, err := repo.GetTotalReadingTime(ctx, "nonexistent")
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
}

func TestReadingProgressRepository_GetReadingTimeByBook(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建进度
	progress := createTestProgress("user123", "book456")
	progress.ReadingTime = 5000
	err := repo.Create(ctx, progress)
	require.NoError(t, err)

	// 查询特定书籍的阅读时长
	time, err := repo.GetReadingTimeByBook(ctx, "user123", "book456")
	require.NoError(t, err)
	assert.Equal(t, int64(5000), time)
}

func TestReadingProgressRepository_GetReadingTimeByPeriod(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建进度记录
	now := time.Now()

	progress1 := createTestProgress("user123", "book1")
	progress1.ReadingTime = 1000
	err := repo.Create(ctx, progress1)
	require.NoError(t, err)
	// Create后更新LastReadAt
	err = repo.Update(ctx, progress1.ID, map[string]interface{}{
		"last_read_at": now.Add(-2 * time.Hour),
	})
	require.NoError(t, err)

	progress2 := createTestProgress("user123", "book2")
	progress2.ReadingTime = 2000
	err = repo.Create(ctx, progress2)
	require.NoError(t, err)
	// Create后更新LastReadAt
	err = repo.Update(ctx, progress2.ID, map[string]interface{}{
		"last_read_at": now.Add(-30 * time.Minute),
	})
	require.NoError(t, err)

	progress3 := createTestProgress("user123", "book3")
	progress3.ReadingTime = 3000
	err = repo.Create(ctx, progress3)
	require.NoError(t, err)
	// Create后更新LastReadAt（超出时间范围）
	err = repo.Update(ctx, progress3.ID, map[string]interface{}{
		"last_read_at": now.Add(-5 * time.Hour),
	})
	require.NoError(t, err)

	// 查询最近3小时的阅读时长
	startTime := now.Add(-3 * time.Hour)
	endTime := now
	total, err := repo.GetReadingTimeByPeriod(ctx, "user123", startTime, endTime)
	require.NoError(t, err)
	assert.Equal(t, int64(3000), total) // progress1 + progress2
}

func TestReadingProgressRepository_CountReadingBooks(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建多个进度
	progress1 := createTestProgress("user123", "book1")
	progress2 := createTestProgress("user123", "book2")
	progress3 := createTestProgress("user123", "book3")

	err := repo.Create(ctx, progress1)
	require.NoError(t, err)
	err = repo.Create(ctx, progress2)
	require.NoError(t, err)
	err = repo.Create(ctx, progress3)
	require.NoError(t, err)

	// 统计书籍数量
	count, err := repo.CountReadingBooks(ctx, "user123")
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

// ============ 阅读记录测试 ============

func TestReadingProgressRepository_GetReadingHistory(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建多个进度记录
	for i := 0; i < 10; i++ {
		progress := createTestProgress("user123", "book"+string(rune('1'+i)))
		err := repo.Create(ctx, progress)
		require.NoError(t, err)
		time.Sleep(5 * time.Millisecond)
	}

	// 分页查询阅读历史
	results, err := repo.GetReadingHistory(ctx, "user123", 5, 2)
	require.NoError(t, err)
	assert.Len(t, results, 5)
}

func TestReadingProgressRepository_GetUnfinishedBooks(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建进度记录
	progress1 := createTestProgress("user123", "book1")
	progress1.Progress = 0.5 // 未读完

	progress2 := createTestProgress("user123", "book2")
	progress2.Progress = 1.0 // 已读完

	progress3 := createTestProgress("user123", "book3")
	progress3.Progress = 0.8 // 未读完

	err := repo.Create(ctx, progress1)
	require.NoError(t, err)
	err = repo.Create(ctx, progress2)
	require.NoError(t, err)
	err = repo.Create(ctx, progress3)
	require.NoError(t, err)

	// 查询未读完的书
	results, err := repo.GetUnfinishedBooks(ctx, "user123")
	require.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestReadingProgressRepository_GetFinishedBooks(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建进度记录
	progress1 := createTestProgress("user123", "book1")
	progress1.Progress = 0.5 // 未读完

	progress2 := createTestProgress("user123", "book2")
	progress2.Progress = 1.0 // 已读完

	progress3 := createTestProgress("user123", "book3")
	progress3.Progress = 1.0 // 已读完

	err := repo.Create(ctx, progress1)
	require.NoError(t, err)
	err = repo.Create(ctx, progress2)
	require.NoError(t, err)
	err = repo.Create(ctx, progress3)
	require.NoError(t, err)

	// 查询已读完的书
	results, err := repo.GetFinishedBooks(ctx, "user123")
	require.NoError(t, err)
	assert.Len(t, results, 2)
}

// ============ 数据同步测试 ============

func TestReadingProgressRepository_SyncProgress(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 准备同步数据（需要有ID）
	progresses := []*reader.ReadingProgress{
		{
			ID:          "sync_prog_1",
			UserID:      "user123",
			BookID:      "book1",
			ChapterID:   "chapter_1",
			Progress:    0.5,
			ReadingTime: 100,
			LastReadAt:  time.Now(),
		},
	}

	// 同步
	err := repo.SyncProgress(ctx, "user123", progresses)
	require.NoError(t, err)

	// 验证
	result, err := repo.GetByUserAndBook(ctx, "user123", "book1")
	require.NoError(t, err)
	assert.Equal(t, 0.5, result.Progress)
}

func TestReadingProgressRepository_GetProgressesByUser(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	oldTime := time.Now().Add(-2 * time.Hour)
	newTime := time.Now()

	// 创建旧的进度记录
	oldProgress := createTestProgress("user123", "book1")
	err := repo.Create(ctx, oldProgress)
	require.NoError(t, err)
	// 使用MongoDB直接操作更新UpdatedAt，避免Update方法自动更新时间
	_, err = global.DB.Collection("reading_progress").UpdateOne(
		ctx,
		map[string]interface{}{"_id": oldProgress.ID},
		map[string]interface{}{"$set": map[string]interface{}{"updated_at": oldTime}},
	)
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	// 创建新的进度记录
	newProgress := createTestProgress("user123", "book2")
	err = repo.Create(ctx, newProgress)
	require.NoError(t, err)
	// 确保新记录的updated_at是最近的
	_, err = global.DB.Collection("reading_progress").UpdateOne(
		ctx,
		map[string]interface{}{"_id": newProgress.ID},
		map[string]interface{}{"$set": map[string]interface{}{"updated_at": newTime}},
	)
	require.NoError(t, err)

	// 查询最近1小时更新的进度
	updatedAfter := time.Now().Add(-1 * time.Hour)
	results, err := repo.GetProgressesByUser(ctx, "user123", updatedAfter)
	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "book2", results[0].BookID)
}

// ============ 清理操作测试 ============

func TestReadingProgressRepository_DeleteOldProgress(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建旧的进度记录
	oldProgress := createTestProgress("user123", "book1")
	err := repo.Create(ctx, oldProgress)
	require.NoError(t, err)
	// Create后更新LastReadAt为1年前
	err = repo.Update(ctx, oldProgress.ID, map[string]interface{}{
		"last_read_at": time.Now().Add(-365 * 24 * time.Hour),
	})
	require.NoError(t, err)

	// 创建新的进度记录
	newProgress := createTestProgress("user123", "book2")
	err = repo.Create(ctx, newProgress)
	require.NoError(t, err)

	// 删除6个月前的记录
	beforeTime := time.Now().Add(-180 * 24 * time.Hour)
	err = repo.DeleteOldProgress(ctx, beforeTime)
	require.NoError(t, err)

	// 验证旧记录被删除
	result, err := repo.GetByUserAndBook(ctx, "user123", "book1")
	assert.NoError(t, err)
	assert.Nil(t, result)

	// 验证新记录仍然存在
	result, err = repo.GetByUserAndBook(ctx, "user123", "book2")
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestReadingProgressRepository_DeleteByBook(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	// 创建多个用户对同一本书的进度
	progress1 := createTestProgress("user1", "book_to_delete")
	progress2 := createTestProgress("user2", "book_to_delete")
	progress3 := createTestProgress("user1", "book_to_keep")

	err := repo.Create(ctx, progress1)
	require.NoError(t, err)
	err = repo.Create(ctx, progress2)
	require.NoError(t, err)
	err = repo.Create(ctx, progress3)
	require.NoError(t, err)

	// 删除某本书的所有进度
	err = repo.DeleteByBook(ctx, "book_to_delete")
	require.NoError(t, err)

	// 验证该书的进度都被删除
	result1, err := repo.GetByUserAndBook(ctx, "user1", "book_to_delete")
	assert.NoError(t, err)
	assert.Nil(t, result1)

	result2, err := repo.GetByUserAndBook(ctx, "user2", "book_to_delete")
	assert.NoError(t, err)
	assert.Nil(t, result2)

	// 验证其他书的进度仍然存在
	result3, err := repo.GetByUserAndBook(ctx, "user1", "book_to_keep")
	require.NoError(t, err)
	assert.NotNil(t, result3)
}

// ============ 健康检查测试 ============

func TestReadingProgressRepository_Health(t *testing.T) {
	setupTest(t)
	ctx := context.Background()

	err := repo.Health(ctx)
	assert.NoError(t, err)
}
