package reader_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	readerModel "Qingyu_backend/models/reader"
	readerRepo "Qingyu_backend/repository/mongodb/reader"
	"Qingyu_backend/models/shared/types"
	"Qingyu_backend/test/testutil"
)

// TestReadingProgressRepository_Create 测试创建阅读进度
func TestReadingProgressRepository_Create(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := readerRepo.NewMongoReadingProgressRepository(db)
	ctx := context.Background()

	userID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()
	chapterID := primitive.NewObjectID()

	progress := &readerModel.ReadingProgress{
		UserID:      userID,
		BookID:      bookID,
		ChapterID:   chapterID,
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

	userID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()
	chapterID := primitive.NewObjectID()

	testProgress := &readerModel.ReadingProgress{
		UserID:      userID,
		BookID:      bookID,
		ChapterID:   chapterID,
		Progress:    0.5,
		ReadingTime: 3600,
		LastReadAt:  time.Now(),
	}
	err := repo.Create(ctx, testProgress)
	require.NoError(t, err)

	// Act
	found, err := repo.GetByID(ctx, testProgress.ID.Hex())

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

	userID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()
	chapterID := primitive.NewObjectID()

	testProgress := &readerModel.ReadingProgress{
		UserID:      userID,
		BookID:      bookID,
		ChapterID:   chapterID,
		Progress:    0.5,
		ReadingTime: 3600,
		LastReadAt:  time.Now(),
	}
	err := repo.Create(ctx, testProgress)
	require.NoError(t, err)

	// Act
	found, err := repo.GetByUserAndBook(ctx, userID.Hex(), bookID.Hex())

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, userID, found.UserID)
	assert.Equal(t, bookID, found.BookID)
}

// TestReadingProgressRepository_Update 测试更新阅读进度
func TestReadingProgressRepository_Update(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := readerRepo.NewMongoReadingProgressRepository(db)
	ctx := context.Background()

	userID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()
	chapterID := primitive.NewObjectID()
	newChapterID := primitive.NewObjectID()

	testProgress := &readerModel.ReadingProgress{
		UserID:      userID,
		BookID:      bookID,
		ChapterID:   chapterID,
		Progress:    0.5,
		ReadingTime: 3600,
		LastReadAt:  time.Now(),
	}
	err := repo.Create(ctx, testProgress)
	require.NoError(t, err)

	// Act - 更新进度
	updates := map[string]interface{}{
		"chapter_id":   newChapterID,
		"progress":     types.Progress(0.8),
		"last_read_at": time.Now(),
	}
	err = repo.Update(ctx, testProgress.ID.Hex(), updates)

	// Assert
	require.NoError(t, err)

	// 验证更新
	found, err := repo.GetByID(ctx, testProgress.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, newChapterID, found.ChapterID)
	assert.Equal(t, types.Progress(0.8), found.Progress)
}

// TestReadingProgressRepository_Delete 测试删除阅读进度
func TestReadingProgressRepository_Delete(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := readerRepo.NewMongoReadingProgressRepository(db)
	ctx := context.Background()

	userID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()
	chapterID := primitive.NewObjectID()

	testProgress := &readerModel.ReadingProgress{
		UserID:      userID,
		BookID:      bookID,
		ChapterID:   chapterID,
		Progress:    0.5,
		ReadingTime: 3600,
		LastReadAt:  time.Now(),
	}
	err := repo.Create(ctx, testProgress)
	require.NoError(t, err)

	// Act - 删除进度
	err = repo.Delete(ctx, testProgress.ID.Hex())

	// Assert
	require.NoError(t, err)

	// 验证已删除
	found, err := repo.GetByID(ctx, testProgress.ID.Hex())
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

	userID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()
	chapterID := primitive.NewObjectID()

	// Act - 保存进度
	err := repo.SaveProgress(ctx, userID.Hex(), bookID.Hex(), chapterID.Hex(), 0.5)

	// Assert
	require.NoError(t, err)

	// 验证进度已保存
	found, err := repo.GetByUserAndBook(ctx, userID.Hex(), bookID.Hex())
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, chapterID, found.ChapterID)
	assert.Equal(t, types.Progress(0.5), found.Progress)
}

// TestReadingProgressRepository_UpdateReadingTime 测试更新阅读时间
func TestReadingProgressRepository_UpdateReadingTime(t *testing.T) {
	// Arrange
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := readerRepo.NewMongoReadingProgressRepository(db)
	ctx := context.Background()

	userID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()
	chapterID := primitive.NewObjectID()

	testProgress := &readerModel.ReadingProgress{
		UserID:      userID,
		BookID:      bookID,
		ChapterID:   chapterID,
		Progress:    0.5,
		ReadingTime: 3600,
		LastReadAt:  time.Now(),
	}
	err := repo.Create(ctx, testProgress)
	require.NoError(t, err)

	// Act - 更新阅读时间
	err = repo.UpdateReadingTime(ctx, userID.Hex(), bookID.Hex(), 1800)

	// Assert
	require.NoError(t, err)

	// 验证阅读时间已更新
	found, err := repo.GetByUserAndBook(ctx, userID.Hex(), bookID.Hex())
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

	userID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()
	chapterID := primitive.NewObjectID()

	testProgress := &readerModel.ReadingProgress{
		UserID:      userID,
		BookID:      bookID,
		ChapterID:   chapterID,
		Progress:    0.5,
		ReadingTime: 3600,
		LastReadAt:  time.Now().Add(-1 * time.Hour),
	}
	err := repo.Create(ctx, testProgress)
	require.NoError(t, err)

	// Act - 更新最后阅读时间
	err = repo.UpdateLastReadAt(ctx, userID.Hex(), bookID.Hex())

	// Assert
	require.NoError(t, err)

	// 验证最后阅读时间已更新
	found, err := repo.GetByUserAndBook(ctx, userID.Hex(), bookID.Hex())
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
	userID := primitive.NewObjectID()
	progresses := []*readerModel.ReadingProgress{
		{UserID: userID, BookID: primitive.NewObjectID(), ChapterID: primitive.NewObjectID(), Progress: 0.3, LastReadAt: time.Now()},
		{UserID: userID, BookID: primitive.NewObjectID(), ChapterID: primitive.NewObjectID(), Progress: 0.5, LastReadAt: time.Now()},
		{UserID: userID, BookID: primitive.NewObjectID(), ChapterID: primitive.NewObjectID(), Progress: 0.8, LastReadAt: time.Now()},
	}

	for _, progress := range progresses {
		err := repo.Create(ctx, progress)
		require.NoError(t, err)
	}

	// Act
	result, err := repo.GetByUser(ctx, userID.Hex())

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
	userID := primitive.NewObjectID()
	for i := 0; i < 5; i++ {
		progress := &readerModel.ReadingProgress{
			UserID:     userID,
			BookID:     primitive.NewObjectID(),
			ChapterID:  primitive.NewObjectID(),
			Progress:   types.Progress(float32(i) * 0.2),
			LastReadAt: time.Now(),
		}
		err := repo.Create(ctx, progress)
		require.NoError(t, err)
	}

	// Act
	result, err := repo.GetRecentReadingByUser(ctx, userID.Hex(), 3)

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
	userID := primitive.NewObjectID()
	for i := 0; i < 5; i++ {
		progress := &readerModel.ReadingProgress{
			UserID:     userID,
			BookID:     primitive.NewObjectID(),
			ChapterID:  primitive.NewObjectID(),
			Progress:   0.5,
			LastReadAt: time.Now(),
		}
		err := repo.Create(ctx, progress)
		require.NoError(t, err)
	}

	// Act
	count, err := repo.CountReadingBooks(ctx, userID.Hex())

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
	userID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()
	for i := 0; i < 3; i++ {
		progress := &readerModel.ReadingProgress{
			UserID:     userID,
			BookID:     bookID,
			ChapterID:  primitive.NewObjectID(),
			Progress:   types.Progress(float32(i+1) * 0.3),
			LastReadAt: time.Now(),
		}
		err := repo.Create(ctx, progress)
		require.NoError(t, err)
	}

	// Act - 删除书籍的所有进度
	err := repo.DeleteByBook(ctx, bookID.Hex())

	// Assert
	require.NoError(t, err)

	// 验证已删除
	allProgress, _ := repo.GetByUser(ctx, userID.Hex())
	for _, progress := range allProgress {
		assert.NotEqual(t, bookID, progress.BookID)
	}
}
