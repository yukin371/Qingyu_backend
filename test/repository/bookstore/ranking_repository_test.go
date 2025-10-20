package bookstore

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/global"
	"Qingyu_backend/models/reading/bookstore"
	mongoBookstore "Qingyu_backend/repository/mongodb/bookstore"
	"Qingyu_backend/test/testutil"
)

func setupRankingTest(t *testing.T) context.Context {
	testutil.SetupTestDB(t)
	ctx := context.Background()

	// 清空测试数据
	_ = global.DB.Collection("rankings").Drop(ctx)
	_ = global.DB.Collection("books").Drop(ctx)

	return ctx
}

// 辅助函数：创建测试书籍
func createTestBook(title, author string) *bookstore.Book {
	return &bookstore.Book{
		ID:            primitive.NewObjectID(),
		Title:         title,
		Author:        author,
		Status:        bookstore.BookStatusPublished,
		WordCount:     50000,
		ChapterCount:  50,
		IsRecommended: true,
		IsFeatured:    false,
		IsHot:         false,
	}
}

// 辅助函数：创建测试榜单项
func createTestRankingItem(bookID primitive.ObjectID, rankType bookstore.RankingType, rank int, period string) *bookstore.RankingItem {
	return &bookstore.RankingItem{
		BookID:    bookID,
		Type:      rankType,
		Rank:      rank,
		Score:     100.0 - float64(rank),
		ViewCount: int64(1000 - rank*10),
		LikeCount: int64(100 - rank),
		Period:    period,
	}
}

// ==================== 基础CRUD测试 ====================

func TestRankingRepository_Create(t *testing.T) {
	ctx := setupRankingTest(t)
	repo := mongoBookstore.NewMongoRankingRepository(global.MongoClient, global.DB.Name())

	t.Run("成功创建榜单项", func(t *testing.T) {
		bookID := primitive.NewObjectID()
		item := createTestRankingItem(bookID, bookstore.RankingTypeRealtime, 1, "2024-01-01")

		err := repo.Create(ctx, item)
		assert.NoError(t, err)
		assert.NotEmpty(t, item.ID)
		assert.False(t, item.CreatedAt.IsZero())
		assert.False(t, item.UpdatedAt.IsZero())
	})

	t.Run("空榜单项", func(t *testing.T) {
		err := repo.Create(ctx, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be nil")
	})
}

func TestRankingRepository_GetByID(t *testing.T) {
	ctx := setupRankingTest(t)
	repo := mongoBookstore.NewMongoRankingRepository(global.MongoClient, global.DB.Name())

	t.Run("成功获取榜单项", func(t *testing.T) {
		bookID := primitive.NewObjectID()
		item := createTestRankingItem(bookID, bookstore.RankingTypeRealtime, 1, "2024-01-01")
		err := repo.Create(ctx, item)
		require.NoError(t, err)

		found, err := repo.GetByID(ctx, item.ID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, item.BookID, found.BookID)
		assert.Equal(t, item.Type, found.Type)
		assert.Equal(t, item.Rank, found.Rank)
	})

	t.Run("榜单项不存在", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		found, err := repo.GetByID(ctx, nonExistentID)
		assert.NoError(t, err)
		assert.Nil(t, found)
	})
}

func TestRankingRepository_Update(t *testing.T) {
	ctx := setupRankingTest(t)
	repo := mongoBookstore.NewMongoRankingRepository(global.MongoClient, global.DB.Name())

	t.Run("成功更新榜单项", func(t *testing.T) {
		bookID := primitive.NewObjectID()
		item := createTestRankingItem(bookID, bookstore.RankingTypeRealtime, 1, "2024-01-01")
		err := repo.Create(ctx, item)
		require.NoError(t, err)

		updates := map[string]interface{}{
			"rank":  2,
			"score": 95.0,
		}
		err = repo.Update(ctx, item.ID, updates)
		assert.NoError(t, err)

		// 验证更新
		found, err := repo.GetByID(ctx, item.ID)
		assert.NoError(t, err)
		assert.Equal(t, 2, found.Rank)
		assert.Equal(t, 95.0, found.Score)
	})

	t.Run("更新不存在的榜单项", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		updates := map[string]interface{}{"rank": 2}
		err := repo.Update(ctx, nonExistentID, updates)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestRankingRepository_Delete(t *testing.T) {
	ctx := setupRankingTest(t)
	repo := mongoBookstore.NewMongoRankingRepository(global.MongoClient, global.DB.Name())

	t.Run("成功删除榜单项", func(t *testing.T) {
		bookID := primitive.NewObjectID()
		item := createTestRankingItem(bookID, bookstore.RankingTypeRealtime, 1, "2024-01-01")
		err := repo.Create(ctx, item)
		require.NoError(t, err)

		err = repo.Delete(ctx, item.ID)
		assert.NoError(t, err)

		// 验证已删除
		found, err := repo.GetByID(ctx, item.ID)
		assert.NoError(t, err)
		assert.Nil(t, found)
	})

	t.Run("删除不存在的榜单项", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		err := repo.Delete(ctx, nonExistentID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

// ==================== 榜单特定查询测试 ====================

func TestRankingRepository_GetByType(t *testing.T) {
	ctx := setupRankingTest(t)
	repo := mongoBookstore.NewMongoRankingRepository(global.MongoClient, global.DB.Name())

	// 准备测试数据
	period := "2024-01-01"
	for i := 1; i <= 5; i++ {
		item := createTestRankingItem(primitive.NewObjectID(), bookstore.RankingTypeRealtime, i, period)
		err := repo.Create(ctx, item)
		require.NoError(t, err)
	}

	// 创建其他类型的榜单
	for i := 1; i <= 3; i++ {
		item := createTestRankingItem(primitive.NewObjectID(), bookstore.RankingTypeWeekly, i, "2024-W01")
		err := repo.Create(ctx, item)
		require.NoError(t, err)
	}

	t.Run("获取实时榜", func(t *testing.T) {
		items, err := repo.GetByType(ctx, bookstore.RankingTypeRealtime, period, 10, 0)
		assert.NoError(t, err)
		assert.Len(t, items, 5)
		// 验证排序（按rank升序）
		for i, item := range items {
			assert.Equal(t, i+1, item.Rank)
		}
	})

	t.Run("分页查询", func(t *testing.T) {
		// 第1页（前3条）
		items, err := repo.GetByType(ctx, bookstore.RankingTypeRealtime, period, 3, 0)
		assert.NoError(t, err)
		assert.Len(t, items, 3)
		assert.Equal(t, 1, items[0].Rank)

		// 第2页（后2条）
		items, err = repo.GetByType(ctx, bookstore.RankingTypeRealtime, period, 3, 3)
		assert.NoError(t, err)
		assert.Len(t, items, 2)
		assert.Equal(t, 4, items[0].Rank)
	})

	t.Run("查询不存在的周期", func(t *testing.T) {
		items, err := repo.GetByType(ctx, bookstore.RankingTypeRealtime, "2099-12-31", 10, 0)
		assert.NoError(t, err)
		assert.Empty(t, items)
	})
}

func TestRankingRepository_GetByTypeWithBooks(t *testing.T) {
	ctx := setupRankingTest(t)
	repo := mongoBookstore.NewMongoRankingRepository(global.MongoClient, global.DB.Name())
	bookCollection := global.DB.Collection("books")

	// 准备测试书籍
	book1 := createTestBook("Test Book 1", "Author 1")
	book2 := createTestBook("Test Book 2", "Author 2")
	_, err := bookCollection.InsertOne(ctx, book1)
	require.NoError(t, err)
	_, err = bookCollection.InsertOne(ctx, book2)
	require.NoError(t, err)

	// 准备榜单数据
	period := "2024-01-01"
	item1 := createTestRankingItem(book1.ID, bookstore.RankingTypeRealtime, 1, period)
	item2 := createTestRankingItem(book2.ID, bookstore.RankingTypeRealtime, 2, period)
	err = repo.Create(ctx, item1)
	require.NoError(t, err)
	err = repo.Create(ctx, item2)
	require.NoError(t, err)

	t.Run("成功获取榜单（包含书籍信息）", func(t *testing.T) {
		items, err := repo.GetByTypeWithBooks(ctx, bookstore.RankingTypeRealtime, period, 10, 0)
		assert.NoError(t, err)
		assert.Len(t, items, 2)

		// 验证书籍信息已填充
		assert.NotNil(t, items[0].Book)
		assert.Equal(t, "Test Book 1", items[0].Book.Title)
		assert.NotNil(t, items[1].Book)
		assert.Equal(t, "Test Book 2", items[1].Book.Title)
	})
}

func TestRankingRepository_GetByBookID(t *testing.T) {
	ctx := setupRankingTest(t)
	repo := mongoBookstore.NewMongoRankingRepository(global.MongoClient, global.DB.Name())

	bookID := primitive.NewObjectID()
	period := "2024-01-01"
	item := createTestRankingItem(bookID, bookstore.RankingTypeRealtime, 1, period)
	err := repo.Create(ctx, item)
	require.NoError(t, err)

	t.Run("成功获取", func(t *testing.T) {
		found, err := repo.GetByBookID(ctx, bookID, bookstore.RankingTypeRealtime, period)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, bookID, found.BookID)
		assert.Equal(t, 1, found.Rank)
	})

	t.Run("书籍不在榜单", func(t *testing.T) {
		nonExistentBookID := primitive.NewObjectID()
		found, err := repo.GetByBookID(ctx, nonExistentBookID, bookstore.RankingTypeRealtime, period)
		assert.NoError(t, err)
		assert.Nil(t, found)
	})
}

// ==================== 榜单统计测试 ====================

func TestRankingRepository_GetRankingStats(t *testing.T) {
	ctx := setupRankingTest(t)
	repo := mongoBookstore.NewMongoRankingRepository(global.MongoClient, global.DB.Name())

	// 准备测试数据
	period := "2024-01-01"
	totalViews := int64(0)
	totalLikes := int64(0)
	for i := 1; i <= 5; i++ {
		item := createTestRankingItem(primitive.NewObjectID(), bookstore.RankingTypeRealtime, i, period)
		totalViews += item.ViewCount
		totalLikes += item.LikeCount
		err := repo.Create(ctx, item)
		require.NoError(t, err)
	}

	t.Run("成功获取统计", func(t *testing.T) {
		stats, err := repo.GetRankingStats(ctx, bookstore.RankingTypeRealtime, period)
		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, int64(5), stats.TotalBooks)
		assert.Equal(t, totalViews, stats.TotalViews)
		assert.Equal(t, totalLikes, stats.TotalLikes)
		assert.Greater(t, stats.AverageScore, 0.0)
	})

	t.Run("空榜单统计", func(t *testing.T) {
		stats, err := repo.GetRankingStats(ctx, bookstore.RankingTypeRealtime, "2099-12-31")
		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, int64(0), stats.TotalBooks)
	})
}

func TestRankingRepository_CountByType(t *testing.T) {
	ctx := setupRankingTest(t)
	repo := mongoBookstore.NewMongoRankingRepository(global.MongoClient, global.DB.Name())

	period := "2024-01-01"
	for i := 1; i <= 5; i++ {
		item := createTestRankingItem(primitive.NewObjectID(), bookstore.RankingTypeRealtime, i, period)
		err := repo.Create(ctx, item)
		require.NoError(t, err)
	}

	t.Run("统计榜单数量", func(t *testing.T) {
		count, err := repo.CountByType(ctx, bookstore.RankingTypeRealtime, period)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)
	})
}

// ==================== 榜单更新测试 ====================

func TestRankingRepository_UpsertRankingItem(t *testing.T) {
	ctx := setupRankingTest(t)
	repo := mongoBookstore.NewMongoRankingRepository(global.MongoClient, global.DB.Name())

	bookID := primitive.NewObjectID()
	period := "2024-01-01"

	t.Run("插入新榜单项", func(t *testing.T) {
		item := createTestRankingItem(bookID, bookstore.RankingTypeRealtime, 1, period)
		err := repo.UpsertRankingItem(ctx, item)
		assert.NoError(t, err)

		// 验证插入
		count, err := repo.CountByType(ctx, bookstore.RankingTypeRealtime, period)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})

	t.Run("更新已存在的榜单项", func(t *testing.T) {
		item := createTestRankingItem(bookID, bookstore.RankingTypeRealtime, 2, period)
		item.Score = 88.0
		err := repo.UpsertRankingItem(ctx, item)
		assert.NoError(t, err)

		// 验证数量不变（是更新而非插入）
		count, err := repo.CountByType(ctx, bookstore.RankingTypeRealtime, period)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)

		// 验证更新内容
		found, err := repo.GetByBookID(ctx, bookID, bookstore.RankingTypeRealtime, period)
		assert.NoError(t, err)
		assert.Equal(t, 2, found.Rank)
		assert.Equal(t, 88.0, found.Score)
	})
}

func TestRankingRepository_BatchUpsertRankingItems(t *testing.T) {
	ctx := setupRankingTest(t)
	repo := mongoBookstore.NewMongoRankingRepository(global.MongoClient, global.DB.Name())

	period := "2024-01-01"
	items := make([]*bookstore.RankingItem, 5)
	for i := 0; i < 5; i++ {
		items[i] = createTestRankingItem(primitive.NewObjectID(), bookstore.RankingTypeRealtime, i+1, period)
	}

	t.Run("批量插入", func(t *testing.T) {
		err := repo.BatchUpsertRankingItems(ctx, items)
		assert.NoError(t, err)

		count, err := repo.CountByType(ctx, bookstore.RankingTypeRealtime, period)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)
	})

	t.Run("批量更新", func(t *testing.T) {
		// 修改数据后再次upsert
		for _, item := range items {
			item.Score += 10.0
		}

		err := repo.BatchUpsertRankingItems(ctx, items)
		assert.NoError(t, err)

		// 数量应该不变
		count, err := repo.CountByType(ctx, bookstore.RankingTypeRealtime, period)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)
	})

	t.Run("空列表", func(t *testing.T) {
		err := repo.BatchUpsertRankingItems(ctx, []*bookstore.RankingItem{})
		assert.NoError(t, err)
	})
}

func TestRankingRepository_UpdateRankings(t *testing.T) {
	ctx := setupRankingTest(t)
	repo := mongoBookstore.NewMongoRankingRepository(global.MongoClient, global.DB.Name())

	period := "2024-01-01"

	// 创建旧榜单
	for i := 1; i <= 3; i++ {
		item := createTestRankingItem(primitive.NewObjectID(), bookstore.RankingTypeRealtime, i, period)
		err := repo.Create(ctx, item)
		require.NoError(t, err)
	}

	t.Run("替换榜单", func(t *testing.T) {
		t.Skip("此测试需要MongoDB副本集支持事务功能")

		// 创建新榜单数据
		newItems := make([]*bookstore.RankingItem, 5)
		for i := 0; i < 5; i++ {
			newItems[i] = createTestRankingItem(primitive.NewObjectID(), bookstore.RankingTypeRealtime, i+1, period)
		}

		// 使用事务替换榜单
		err := repo.UpdateRankings(ctx, bookstore.RankingTypeRealtime, period, newItems)
		assert.NoError(t, err)

		// 验证新榜单
		items, err := repo.GetByType(ctx, bookstore.RankingTypeRealtime, period, 10, 0)
		assert.NoError(t, err)
		assert.Len(t, items, 5)
	})
}

// ==================== 榜单维护测试 ====================

func TestRankingRepository_DeleteByPeriod(t *testing.T) {
	ctx := setupRankingTest(t)
	repo := mongoBookstore.NewMongoRankingRepository(global.MongoClient, global.DB.Name())

	period1 := "2024-01-01"
	period2 := "2024-01-02"

	// 创建两个周期的数据
	for i := 1; i <= 3; i++ {
		item1 := createTestRankingItem(primitive.NewObjectID(), bookstore.RankingTypeRealtime, i, period1)
		err := repo.Create(ctx, item1)
		require.NoError(t, err)

		item2 := createTestRankingItem(primitive.NewObjectID(), bookstore.RankingTypeRealtime, i, period2)
		err = repo.Create(ctx, item2)
		require.NoError(t, err)
	}

	t.Run("删除指定周期", func(t *testing.T) {
		err := repo.DeleteByPeriod(ctx, period1)
		assert.NoError(t, err)

		// 验证period1已删除
		count1, err := repo.CountByType(ctx, bookstore.RankingTypeRealtime, period1)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count1)

		// 验证period2还在
		count2, err := repo.CountByType(ctx, bookstore.RankingTypeRealtime, period2)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), count2)
	})
}

func TestRankingRepository_DeleteExpiredRankings(t *testing.T) {
	ctx := setupRankingTest(t)
	repo := mongoBookstore.NewMongoRankingRepository(global.MongoClient, global.DB.Name())

	// 创建过期和未过期的数据
	oldPeriod := "2023-01-01"
	newPeriod := "2024-01-01"

	for i := 1; i <= 2; i++ {
		oldItem := createTestRankingItem(primitive.NewObjectID(), bookstore.RankingTypeRealtime, i, oldPeriod)
		err := repo.Create(ctx, oldItem)
		require.NoError(t, err)

		newItem := createTestRankingItem(primitive.NewObjectID(), bookstore.RankingTypeRealtime, i, newPeriod)
		err = repo.Create(ctx, newItem)
		require.NoError(t, err)
	}

	t.Run("删除过期榜单", func(t *testing.T) {
		// 删除2024-01-01之前的数据
		cutoffDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		err := repo.DeleteExpiredRankings(ctx, cutoffDate)
		assert.NoError(t, err)

		// 验证旧数据已删除
		count1, err := repo.CountByType(ctx, bookstore.RankingTypeRealtime, oldPeriod)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count1)

		// 验证新数据还在
		count2, err := repo.CountByType(ctx, bookstore.RankingTypeRealtime, newPeriod)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count2)
	})
}

// ==================== 榜单计算测试 ====================

func TestRankingRepository_CalculateRealtimeRanking(t *testing.T) {
	ctx := setupRankingTest(t)
	repo := mongoBookstore.NewMongoRankingRepository(global.MongoClient, global.DB.Name())
	bookCollection := global.DB.Collection("books")

	// 准备测试书籍
	books := []*bookstore.Book{
		createTestBook("Book 1", "Author 1"),
		createTestBook("Book 2", "Author 2"),
		createTestBook("Book 3", "Author 3"),
	}
	// 设置不同的推荐权重
	books[0].IsRecommended = true
	books[0].IsFeatured = true
	books[0].IsHot = true
	books[1].IsRecommended = true
	books[1].IsFeatured = false
	books[2].IsRecommended = false

	for _, book := range books {
		_, err := bookCollection.InsertOne(ctx, book)
		require.NoError(t, err)
	}

	t.Run("计算实时榜", func(t *testing.T) {
		t.Skip("此测试需要BookStats表支持，Book模型中没有view_count和like_count字段")

		period := "2024-01-01"
		items, err := repo.CalculateRealtimeRanking(ctx, period)
		assert.NoError(t, err)
		assert.NotEmpty(t, items)

		// 验证排序（按热度分数降序）
		if len(items) >= 2 {
			assert.Greater(t, items[0].Score, items[1].Score)
			assert.Equal(t, 1, items[0].Rank)
			assert.Equal(t, 2, items[1].Rank)
		}
	})
}

func TestRankingRepository_CalculateNewbieRanking(t *testing.T) {
	ctx := setupRankingTest(t)
	repo := mongoBookstore.NewMongoRankingRepository(global.MongoClient, global.DB.Name())
	bookCollection := global.DB.Collection("books")

	// 准备新书和旧书
	now := time.Now()
	newBook := createTestBook("New Book", "New Author")
	newBook.CreatedAt = now.AddDate(0, -1, 0) // 1个月前

	oldBook := createTestBook("Old Book", "Old Author")
	oldBook.CreatedAt = now.AddDate(0, -6, 0) // 6个月前

	_, err := bookCollection.InsertOne(ctx, newBook)
	require.NoError(t, err)
	_, err = bookCollection.InsertOne(ctx, oldBook)
	require.NoError(t, err)

	t.Run("计算新人榜（只包含新书）", func(t *testing.T) {
		period := now.Format("2006-01")
		items, err := repo.CalculateNewbieRanking(ctx, period)
		assert.NoError(t, err)

		// 验证只有新书在榜单中
		for _, item := range items {
			assert.NotEqual(t, oldBook.ID, item.BookID)
		}
	})
}

// ==================== 边界测试 ====================

func TestRankingRepository_EdgeCases(t *testing.T) {
	ctx := setupRankingTest(t)
	repo := mongoBookstore.NewMongoRankingRepository(global.MongoClient, global.DB.Name())

	t.Run("创建榜单项后立即查询", func(t *testing.T) {
		bookID := primitive.NewObjectID()
		period := "2024-01-01"
		item := createTestRankingItem(bookID, bookstore.RankingTypeRealtime, 1, period)

		err := repo.Create(ctx, item)
		assert.NoError(t, err)

		// 立即查询应该能够查到
		found, err := repo.GetByBookID(ctx, bookID, bookstore.RankingTypeRealtime, period)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, item.ID, found.ID)
	})

	t.Run("空周期字符串", func(t *testing.T) {
		items, err := repo.GetByType(ctx, bookstore.RankingTypeRealtime, "", 10, 0)
		assert.NoError(t, err)
		assert.Empty(t, items)
	})
}
