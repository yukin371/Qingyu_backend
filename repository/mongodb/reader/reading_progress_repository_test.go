package reader_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	readerModel "Qingyu_backend/models/reader"
	readerRepo "Qingyu_backend/repository/mongodb/reader"
	"Qingyu_backend/test/testutil"
)

// TestReadingProgressRepository_Create 测试创建阅读进度
func TestReadingProgressRepository_Create(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := readerRepo.NewMongoReadingProgressRepository(db)
	ctx := context.Background()

	progress := &readerModel.ReadingProgress{
		UserID:      "user123",
		BookID:      "book123",
		ChapterID:   "chapter1",
		Progress:    0.5,
		ReadingTime: 3600,
		LastReadAt:  time.Now(),
	}

	// Act
	err := repo.Create(ctx, progress)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, progress.ID)
	assert.NotZero(t, progress.CreatedAt)
	assert.NotZero(t, progress.UpdatedAt)
}

// TestReadingProgressRepository_GetByID 测试根据ID获取阅读进度
func TestReadingProgressRepository_GetByID(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := readerRepo.NewMongoReadingProgressRepository(db)
	ctx := context.Background()

	testProgress := &readerModel.ReadingProgress{
		UserID:      "user123",
		BookID:      "book123",
		ChapterID:   "chapter1",
		Progress:    0.5,
		ReadingTime: 3600,
		LastReadAt:  time.Now(),
	}
	err := repo.Create(ctx, testProgress)
	require.NoError(t, err)

	// Act
	found, err := repo.GetByID(ctx, testProgress.ID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, testProgress.UserID, found.UserID)
	assert.Equal(t, testProgress.BookID, found.BookID)
	assert.Equal(t, testProgress.ChapterID, found.ChapterID)
}

// TestReadingProgressRepository_GetByID_NotFound 测试获取不存在的阅读进度
func TestReadingProgressRepository_GetByID_NotFound(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := readerRepo.NewMongoReadingProgressRepository(db)
	ctx := context.Background()

	// Act
	found, err := repo.GetByID(ctx, "nonexistent_id")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, found)
}

// TestReadingProgressRepository_GetByUserAndBook 测试根据用户和书籍获取阅读进度
func TestReadingProgressRepository_GetByUserAndBook(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := readerRepo.NewMongoReadingProgressRepository(db)
	ctx := context.Background()

	testProgress := &readerModel.ReadingProgress{
		UserID:      "user123",
		BookID:      "book123",
		ChapterID:   "chapter1",
		Progress:    0.5,
		ReadingTime: 3600,
		LastReadAt:  time.Now(),
	}
	err := repo.Create(ctx, testProgress)
	require.NoError(t, err)

	// Act
	found, err := repo.GetByUserAndBook(ctx, "user123", "book123")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, "user123", found.UserID)
	assert.Equal(t, "book123", found.BookID)
}

// TestReadingProgressRepository_Update 测试更新阅读进度
func TestReadingProgressRepository_Update(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := readerRepo.NewMongoReadingProgressRepository(db)
	ctx := context.Background()

	testProgress := &readerModel.ReadingProgress{
		UserID:      "user123",
		BookID:      "book123",
		ChapterID:   "chapter1",
		Progress:    0.5,
		ReadingTime: 3600,
		LastReadAt:  time.Now(),
	}
	err := repo.Create(ctx, testProgress)
	require.NoError(t, err)

	// Act - 更新进度
	updates := map[string]interface{}{
		"chapter_id": "chapter2",
		"progress":   0.8,
		"last_read_at": time.Now(),
	}
	err = repo.Update(ctx, testProgress.ID, updates)

	// Assert
	require.NoError(t, err)

	// 验证更新
	found, err := repo.GetByID(ctx, testProgress.ID)
	require.NoError(t, err)
	assert.Equal(t, "chapter2", found.ChapterID)
	assert.Equal(t, 0.8, found.Progress)
}

// TestReadingProgressRepository_Delete 测试删除阅读进度
func TestReadingProgressRepository_Delete(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := readerRepo.NewMongoReadingProgressRepository(db)
	ctx := context.Background()

	testProgress := &readerModel.ReadingProgress{
		UserID:      "user123",
		BookID:      "book123",
		ChapterID:   "chapter1",
		Progress:    0.5,
		ReadingTime: 3600,
		LastReadAt:  time.Now(),
	}
	err := repo.Create(ctx, testProgress)
	require.NoError(t, err)

	// Act - 删除进度
	err = repo.Delete(ctx, testProgress.ID)

	// Assert
	require.NoError(t, err)

	// 验证已删除
	found, err := repo.GetByID(ctx, testProgress.ID)
	assert.Error(t, err)
	assert.Nil(t, found)
}

// TestReadingProgressRepository_SaveProgress 测试保存阅读进度
func TestReadingProgressRepository_SaveProgress(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := readerRepo.NewMongoReadingProgressRepository(db)
	ctx := context.Background()

	// Act - 保存进度
	err := repo.SaveProgress(ctx, "user123", "book123", "chapter1", 0.5)

	// Assert
	require.NoError(t, err)

	// 验证进度已保存
	found, err := repo.GetByUserAndBook(ctx, "user123", "book123")
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, "chapter1", found.ChapterID)
	assert.Equal(t, 0.5, found.Progress)
}

// TestReadingProgressRepository_UpdateReadingTime 测试更新阅读时间
func TestReadingProgressRepository_UpdateReadingTime(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := readerRepo.NewMongoReadingProgressRepository(db)
	ctx := context.Background()

	testProgress := &readerModel.ReadingProgress{
		UserID:      "user123",
		BookID:      "book123",
		ChapterID:   "chapter1",
		Progress:    0.5,
		ReadingTime: 3600,
		LastReadAt:  time.Now(),
	}
	err := repo.Create(ctx, testProgress)
	require.NoError(t, err)

	// Act - 更新阅读时间
	err = repo.UpdateReadingTime(ctx, "user123", "book123", 1800)

	// Assert
	require.NoError(t, err)

	// 验证阅读时间已更新
	found, err := repo.GetByUserAndBook(ctx, "user123", "book123")
	require.NoError(t, err)
	assert.Equal(t, int64(5400), found.ReadingTime) // 3600 + 1800
}

// TestReadingProgressRepository_UpdateLastReadAt 测试更新最后阅读时间
func TestReadingProgressRepository_UpdateLastReadAt(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := readerRepo.NewMongoReadingProgressRepository(db)
	ctx := context.Background()

	testProgress := &readerModel.ReadingProgress{
		UserID:      "user123",
		BookID:      "book123",
		ChapterID:   "chapter1",
		Progress:    0.5,
		ReadingTime: 3600,
		LastReadAt:  time.Now().Add(-1 * time.Hour),
	}
	err := repo.Create(ctx, testProgress)
	require.NoError(t, err)

	// Act - 更新最后阅读时间
	err = repo.UpdateLastReadAt(ctx, "user123", "book123")

	// Assert
	require.NoError(t, err)

	// 验证最后阅读时间已更新
	found, err := repo.GetByUserAndBook(ctx, "user123", "book123")
	require.NoError(t, err)
	assert.True(t, found.LastReadAt.After(testProgress.LastReadAt))
}

// TestReadingProgressRepository_GetByUser 测试获取用户的所有阅读进度
func TestReadingProgressRepository_GetByUser(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := readerRepo.NewMongoReadingProgressRepository(db)
	ctx := context.Background()

	// 创建多个阅读进度 - 使用不同的BookID避免ID冲突
	userID := "user_getby_test_" + primitive.NewObjectID().Hex()
	progresses := []*readerModel.ReadingProgress{
		{ID: primitive.NewObjectID().Hex(), UserID: userID, BookID: primitive.NewObjectID().Hex(), ChapterID: "chapter1", Progress: 0.3, LastReadAt: time.Now()},
		{ID: primitive.NewObjectID().Hex(), UserID: userID, BookID: primitive.NewObjectID().Hex(), ChapterID: "chapter2", Progress: 0.5, LastReadAt: time.Now()},
		{ID: primitive.NewObjectID().Hex(), UserID: userID, BookID: primitive.NewObjectID().Hex(), ChapterID: "chapter3", Progress: 0.8, LastReadAt: time.Now()},
	}

	for _, progress := range progresses {
		err := repo.Create(ctx, progress)
		require.NoError(t, err)
	}

	// Act
	result, err := repo.GetByUser(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.GreaterOrEqual(t, len(result), 3)
}

// TestReadingProgressRepository_GetRecentReadingByUser 测试获取用户最近的阅读记录
func TestReadingProgressRepository_GetRecentReadingByUser(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := readerRepo.NewMongoReadingProgressRepository(db)
	ctx := context.Background()

	// 创建阅读进度 - 使用唯一ID避免冲突
	userID := "user_recent_test_" + primitive.NewObjectID().Hex()
	for i := 0; i < 5; i++ {
		progress := &readerModel.ReadingProgress{
			ID:          primitive.NewObjectID().Hex(),
			UserID:      userID,
			BookID:      primitive.NewObjectID().Hex(),
			ChapterID:   "chapter1",
			Progress:    float64(i) * 0.2,
			LastReadAt:  time.Now(),
		}
		err := repo.Create(ctx, progress)
		require.NoError(t, err)
	}

	// Act
	result, err := repo.GetRecentReadingByUser(ctx, userID, 3)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.GreaterOrEqual(t, len(result), 3)
}

// TestReadingProgressRepository_CountReadingBooks 测试统计用户阅读的书籍数量
func TestReadingProgressRepository_CountReadingBooks(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := readerRepo.NewMongoReadingProgressRepository(db)
	ctx := context.Background()

	// 创建阅读进度 - 使用唯一ID避免冲突
	userID := "user_count_test_" + primitive.NewObjectID().Hex()
	for i := 0; i < 5; i++ {
		progress := &readerModel.ReadingProgress{
			ID:          primitive.NewObjectID().Hex(),
			UserID:      userID,
			BookID:      primitive.NewObjectID().Hex(),
			ChapterID:   "chapter1",
			Progress:    0.5,
			LastReadAt:  time.Now(),
		}
		err := repo.Create(ctx, progress)
		require.NoError(t, err)
	}

	// Act
	count, err := repo.CountReadingBooks(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(5))
}

// TestReadingProgressRepository_Health 测试健康检查
func TestReadingProgressRepository_Health(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := readerRepo.NewMongoReadingProgressRepository(db)
	ctx := context.Background()

	// Act
	err := repo.Health(ctx)

	// Assert
	assert.NoError(t, err)
}

// TestReadingProgressRepository_DeleteByBook 测试删除书籍的所有阅读进度
func TestReadingProgressRepository_DeleteByBook(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := readerRepo.NewMongoReadingProgressRepository(db)
	ctx := context.Background()

	// 创建同一本书的多个章节进度 - 使用唯一ID避免冲突
	bookID := fmt.Sprintf("book_delete_test_%d", time.Now().UnixNano())
	for i := 0; i < 3; i++ {
		progress := &readerModel.ReadingProgress{
			ID:          primitive.NewObjectID().Hex(),
			UserID:      "user123",
			BookID:      bookID,
			ChapterID:   fmt.Sprintf("chapter%d", i+1),
			Progress:    float64(i+1) * 0.3,
			LastReadAt:  time.Now(),
		}
		err := repo.Create(ctx, progress)
		require.NoError(t, err)
	}

	// Act - 删除书籍的所有进度
	err := repo.DeleteByBook(ctx, bookID)

	// Assert
	require.NoError(t, err)

	// 验证已删除
	allProgress, _ := repo.GetByUser(ctx, "user123")
	for _, progress := range allProgress {
		assert.NotEqual(t, bookID, progress.BookID)
	}
}
