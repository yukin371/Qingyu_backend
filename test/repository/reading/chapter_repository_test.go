package reading

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"

	"Qingyu_backend/global"
	"Qingyu_backend/models/reading/reader"
	"Qingyu_backend/repository/mongodb/reading"
	"Qingyu_backend/test/testutil"
)

var chapterIDCounter int

func setupChapterTest(t *testing.T) context.Context {
	testutil.SetupTestDB(t)
	ctx := context.Background()

	// 清空测试数据
	_ = global.DB.Collection("chapters").Drop(ctx)

	return ctx
}

// 辅助函数：创建测试章节（不设置ID，让Repository生成）
func createTestChapter(bookID string, chapterNum int) *reader.Chapter {
	return &reader.Chapter{
		BookID:      bookID,
		Title:       fmt.Sprintf("第%d章", chapterNum),
		Content:     fmt.Sprintf("这是第%d章的内容", chapterNum),
		WordCount:   1000 + chapterNum*100,
		ChapterNum:  chapterNum,
		IsVIP:       chapterNum > 5, // 第6章及以后是VIP
		Price:       int64(chapterNum * 10),
		Status:      1, // 正常状态
		PublishTime: time.Now().Add(-24 * time.Hour),
	}
}

// 辅助函数：创建并插入章节（直接插入到MongoDB）
func createAndInsertChapter(ctx context.Context, t *testing.T, bookID string, chapterNum int) string {
	chapterIDCounter++
	chapterID := fmt.Sprintf("chapter_%d_%d", time.Now().UnixNano(), chapterIDCounter)

	chapter := bson.M{
		"_id":          chapterID,
		"book_id":      bookID,
		"title":        fmt.Sprintf("第%d章", chapterNum),
		"content":      fmt.Sprintf("这是第%d章的内容", chapterNum),
		"word_count":   1000 + chapterNum*100,
		"chapter_num":  chapterNum,
		"is_vip":       chapterNum > 5,
		"price":        int64(chapterNum * 10),
		"status":       1,
		"publish_time": time.Now().Add(-24 * time.Hour),
		"created_at":   time.Now(),
		"updated_at":   time.Now(),
	}

	_, err := global.DB.Collection("chapters").InsertOne(ctx, chapter)
	require.NoError(t, err)

	return chapterID
}

// ==================== 基础CRUD测试 ====================

func TestChapterRepository_Create(t *testing.T) {
	ctx := setupChapterTest(t)
	repo := reading.NewMongoChapterRepository(global.DB)

	t.Run("成功创建章节", func(t *testing.T) {
		chapter := createTestChapter("book001", 1)

		err := repo.Create(ctx, chapter)
		assert.NoError(t, err)
		assert.NotEmpty(t, chapter.ID)
		assert.False(t, chapter.CreatedAt.IsZero())
		assert.False(t, chapter.UpdatedAt.IsZero())
	})
}

func TestChapterRepository_GetByID(t *testing.T) {
	ctx := setupChapterTest(t)
	repo := reading.NewMongoChapterRepository(global.DB)

	t.Run("成功获取章节", func(t *testing.T) {
		chapterID := createAndInsertChapter(ctx, t, "book001", 1)

		chapter, err := repo.GetByID(ctx, chapterID)
		assert.NoError(t, err)
		assert.NotNil(t, chapter)
		assert.Equal(t, chapterID, chapter.ID)
		assert.Equal(t, "book001", chapter.BookID)
	})

	t.Run("章节不存在", func(t *testing.T) {
		_, err := repo.GetByID(ctx, "nonexistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "章节不存在")
	})
}

func TestChapterRepository_Update(t *testing.T) {
	ctx := setupChapterTest(t)
	repo := reading.NewMongoChapterRepository(global.DB)

	t.Run("成功更新章节", func(t *testing.T) {
		chapterID := createAndInsertChapter(ctx, t, "book001", 1)

		updates := map[string]interface{}{
			"title":      "更新后的标题",
			"word_count": 2000,
		}

		err := repo.Update(ctx, chapterID, updates)
		assert.NoError(t, err)

		// 验证更新
		chapter, _ := repo.GetByID(ctx, chapterID)
		assert.Equal(t, "更新后的标题", chapter.Title)
		assert.Equal(t, 2000, chapter.WordCount)
	})

	t.Run("更新不存在的章节", func(t *testing.T) {
		updates := map[string]interface{}{"title": "新标题"}
		err := repo.Update(ctx, "nonexistent", updates)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "章节不存在")
	})
}

func TestChapterRepository_Delete(t *testing.T) {
	ctx := setupChapterTest(t)
	repo := reading.NewMongoChapterRepository(global.DB)

	t.Run("成功删除章节", func(t *testing.T) {
		chapterID := createAndInsertChapter(ctx, t, "book001", 1)

		err := repo.Delete(ctx, chapterID)
		assert.NoError(t, err)

		// 验证已删除
		_, err = repo.GetByID(ctx, chapterID)
		assert.Error(t, err)
	})

	t.Run("删除不存在的章节", func(t *testing.T) {
		err := repo.Delete(ctx, "nonexistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "章节不存在")
	})
}

// ==================== 章节查询测试 ====================

func TestChapterRepository_GetByBookID(t *testing.T) {
	ctx := setupChapterTest(t)
	repo := reading.NewMongoChapterRepository(global.DB)

	t.Run("成功获取书籍所有章节", func(t *testing.T) {
		// 创建多个章节
		for i := 1; i <= 5; i++ {
			createAndInsertChapter(ctx, t, "book001", i)
		}

		chapters, err := repo.GetByBookID(ctx, "book001")
		assert.NoError(t, err)
		assert.Len(t, chapters, 5)

		// 验证按chapter_num排序
		for i, chapter := range chapters {
			assert.Equal(t, i+1, chapter.ChapterNum)
		}
	})

	t.Run("只返回正常状态的章节", func(t *testing.T) {
		_ = global.DB.Collection("chapters").Drop(ctx)

		// 创建正常章节
		createAndInsertChapter(ctx, t, "book002", 1)

		// 创建删除状态章节
		chapterIDCounter++
		deletedChapterID := fmt.Sprintf("chapter_%d_%d", time.Now().UnixNano(), chapterIDCounter)
		deletedChapter := bson.M{
			"_id":         deletedChapterID,
			"book_id":     "book002",
			"title":       "删除的章节",
			"chapter_num": 2,
			"status":      2, // 删除状态
			"created_at":  time.Now(),
			"updated_at":  time.Now(),
		}
		_, _ = global.DB.Collection("chapters").InsertOne(ctx, deletedChapter)

		chapters, err := repo.GetByBookID(ctx, "book002")
		assert.NoError(t, err)
		assert.Len(t, chapters, 1) // 只有1个正常状态
	})

	t.Run("书籍没有章节", func(t *testing.T) {
		chapters, err := repo.GetByBookID(ctx, "nonexistent_book")
		assert.NoError(t, err)
		assert.Empty(t, chapters)
	})
}

func TestChapterRepository_GetByBookIDWithPagination(t *testing.T) {
	ctx := setupChapterTest(t)
	repo := reading.NewMongoChapterRepository(global.DB)

	t.Run("分页获取章节", func(t *testing.T) {
		// 创建10个章节
		for i := 1; i <= 10; i++ {
			createAndInsertChapter(ctx, t, "book001", i)
		}

		// 获取第1页（前5章）
		chapters, err := repo.GetByBookIDWithPagination(ctx, "book001", 5, 0)
		assert.NoError(t, err)
		assert.Len(t, chapters, 5)
		assert.Equal(t, 1, chapters[0].ChapterNum)
		assert.Equal(t, 5, chapters[4].ChapterNum)

		// 获取第2页（后5章）
		chapters, err = repo.GetByBookIDWithPagination(ctx, "book001", 5, 5)
		assert.NoError(t, err)
		assert.Len(t, chapters, 5)
		assert.Equal(t, 6, chapters[0].ChapterNum)
		assert.Equal(t, 10, chapters[4].ChapterNum)
	})
}

func TestChapterRepository_GetByChapterNum(t *testing.T) {
	ctx := setupChapterTest(t)
	repo := reading.NewMongoChapterRepository(global.DB)

	t.Run("成功获取指定章节号的章节", func(t *testing.T) {
		createAndInsertChapter(ctx, t, "book001", 3)

		chapter, err := repo.GetByChapterNum(ctx, "book001", 3)
		assert.NoError(t, err)
		assert.NotNil(t, chapter)
		assert.Equal(t, 3, chapter.ChapterNum)
	})

	t.Run("章节不存在", func(t *testing.T) {
		_, err := repo.GetByChapterNum(ctx, "book001", 999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "章节不存在")
	})
}

// ==================== 章节导航测试 ====================

func TestChapterRepository_GetPrevChapter(t *testing.T) {
	ctx := setupChapterTest(t)
	repo := reading.NewMongoChapterRepository(global.DB)

	t.Run("获取上一章", func(t *testing.T) {
		// 创建章节1, 2, 3
		for i := 1; i <= 3; i++ {
			createAndInsertChapter(ctx, t, "book001", i)
		}

		// 获取第3章的上一章
		prevChapter, err := repo.GetPrevChapter(ctx, "book001", 3)
		assert.NoError(t, err)
		assert.NotNil(t, prevChapter)
		assert.Equal(t, 2, prevChapter.ChapterNum)
	})

	t.Run("已经是第一章", func(t *testing.T) {
		prevChapter, err := repo.GetPrevChapter(ctx, "book001", 1)
		assert.NoError(t, err)
		assert.Nil(t, prevChapter) // 没有上一章
	})
}

func TestChapterRepository_GetNextChapter(t *testing.T) {
	ctx := setupChapterTest(t)
	repo := reading.NewMongoChapterRepository(global.DB)

	t.Run("获取下一章", func(t *testing.T) {
		// 创建章节1, 2, 3
		for i := 1; i <= 3; i++ {
			createAndInsertChapter(ctx, t, "book001", i)
		}

		// 获取第1章的下一章
		nextChapter, err := repo.GetNextChapter(ctx, "book001", 1)
		assert.NoError(t, err)
		assert.NotNil(t, nextChapter)
		assert.Equal(t, 2, nextChapter.ChapterNum)
	})

	t.Run("已经是最后一章", func(t *testing.T) {
		nextChapter, err := repo.GetNextChapter(ctx, "book001", 999)
		assert.NoError(t, err)
		assert.Nil(t, nextChapter) // 没有下一章
	})
}

func TestChapterRepository_GetFirstChapter(t *testing.T) {
	ctx := setupChapterTest(t)
	repo := reading.NewMongoChapterRepository(global.DB)

	t.Run("获取第一章", func(t *testing.T) {
		// 创建章节3, 1, 2（无序插入）
		createAndInsertChapter(ctx, t, "book001", 3)
		createAndInsertChapter(ctx, t, "book001", 1)
		createAndInsertChapter(ctx, t, "book001", 2)

		firstChapter, err := repo.GetFirstChapter(ctx, "book001")
		assert.NoError(t, err)
		assert.NotNil(t, firstChapter)
		assert.Equal(t, 1, firstChapter.ChapterNum) // 应该返回章节号最小的
	})

	t.Run("书籍没有章节", func(t *testing.T) {
		_, err := repo.GetFirstChapter(ctx, "nonexistent_book")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "该书没有章节")
	})
}

func TestChapterRepository_GetLastChapter(t *testing.T) {
	ctx := setupChapterTest(t)
	repo := reading.NewMongoChapterRepository(global.DB)

	t.Run("获取最后一章", func(t *testing.T) {
		// 创建章节1, 2, 3
		for i := 1; i <= 3; i++ {
			createAndInsertChapter(ctx, t, "book001", i)
		}

		lastChapter, err := repo.GetLastChapter(ctx, "book001")
		assert.NoError(t, err)
		assert.NotNil(t, lastChapter)
		assert.Equal(t, 3, lastChapter.ChapterNum)
	})

	t.Run("书籍没有章节", func(t *testing.T) {
		_, err := repo.GetLastChapter(ctx, "nonexistent_book")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "该书没有章节")
	})
}

// ==================== 章节状态查询测试 ====================

func TestChapterRepository_GetPublishedChapters(t *testing.T) {
	ctx := setupChapterTest(t)
	repo := reading.NewMongoChapterRepository(global.DB)

	t.Run("获取已发布章节", func(t *testing.T) {
		// 创建已发布章节
		chapterIDCounter++
		publishedID := fmt.Sprintf("chapter_%d_%d", time.Now().UnixNano(), chapterIDCounter)
		publishedChapter := bson.M{
			"_id":          publishedID,
			"book_id":      "book001",
			"title":        "已发布章节",
			"chapter_num":  1,
			"status":       1,
			"publish_time": time.Now().Add(-1 * time.Hour), // 1小时前发布
			"created_at":   time.Now(),
			"updated_at":   time.Now(),
		}
		_, _ = global.DB.Collection("chapters").InsertOne(ctx, publishedChapter)

		// 创建未发布章节（发布时间在未来）
		chapterIDCounter++
		unpublishedID := fmt.Sprintf("chapter_%d_%d", time.Now().UnixNano(), chapterIDCounter)
		unpublishedChapter := bson.M{
			"_id":          unpublishedID,
			"book_id":      "book001",
			"title":        "未发布章节",
			"chapter_num":  2,
			"status":       1,
			"publish_time": time.Now().Add(24 * time.Hour), // 24小时后发布
			"created_at":   time.Now(),
			"updated_at":   time.Now(),
		}
		_, _ = global.DB.Collection("chapters").InsertOne(ctx, unpublishedChapter)

		chapters, err := repo.GetPublishedChapters(ctx, "book001")
		assert.NoError(t, err)
		assert.Len(t, chapters, 1) // 只有已发布的
		assert.Equal(t, 1, chapters[0].ChapterNum)
	})
}

func TestChapterRepository_GetVIPChapters(t *testing.T) {
	ctx := setupChapterTest(t)
	repo := reading.NewMongoChapterRepository(global.DB)

	t.Run("获取VIP章节", func(t *testing.T) {
		// 创建免费章节和VIP章节
		for i := 1; i <= 10; i++ {
			createAndInsertChapter(ctx, t, "book001", i)
		}

		vipChapters, err := repo.GetVIPChapters(ctx, "book001")
		assert.NoError(t, err)
		assert.Len(t, vipChapters, 5) // 第6-10章是VIP（根据createTestChapter逻辑，chapterNum > 5）

		// 验证都是VIP章节
		for _, chapter := range vipChapters {
			assert.True(t, chapter.IsVIP)
			assert.Greater(t, chapter.ChapterNum, 5)
		}
	})
}

func TestChapterRepository_GetFreeChapters(t *testing.T) {
	ctx := setupChapterTest(t)
	repo := reading.NewMongoChapterRepository(global.DB)

	t.Run("获取免费章节", func(t *testing.T) {
		// 创建免费章节和VIP章节
		for i := 1; i <= 10; i++ {
			createAndInsertChapter(ctx, t, "book001", i)
		}

		freeChapters, err := repo.GetFreeChapters(ctx, "book001")
		assert.NoError(t, err)
		assert.Len(t, freeChapters, 5) // 第1-5章是免费

		// 验证都是免费章节
		for _, chapter := range freeChapters {
			assert.False(t, chapter.IsVIP)
			assert.LessOrEqual(t, chapter.ChapterNum, 5)
		}
	})
}

// ==================== 统计查询测试 ====================

func TestChapterRepository_CountByBookID(t *testing.T) {
	ctx := setupChapterTest(t)
	repo := reading.NewMongoChapterRepository(global.DB)

	t.Run("统计书籍章节数", func(t *testing.T) {
		// 创建5个章节
		for i := 1; i <= 5; i++ {
			createAndInsertChapter(ctx, t, "book001", i)
		}

		count, err := repo.CountByBookID(ctx, "book001")
		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)
	})

	t.Run("书籍没有章节", func(t *testing.T) {
		count, err := repo.CountByBookID(ctx, "nonexistent_book")
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})
}

func TestChapterRepository_CountByStatus(t *testing.T) {
	ctx := setupChapterTest(t)
	repo := reading.NewMongoChapterRepository(global.DB)

	t.Run("按状态统计章节数", func(t *testing.T) {
		// 创建正常状态章节
		for i := 1; i <= 3; i++ {
			createAndInsertChapter(ctx, t, "book001", i)
		}

		// 创建删除状态章节
		chapterIDCounter++
		deletedID := fmt.Sprintf("chapter_%d_%d", time.Now().UnixNano(), chapterIDCounter)
		deletedChapter := bson.M{
			"_id":         deletedID,
			"book_id":     "book001",
			"title":       "删除的章节",
			"chapter_num": 4,
			"status":      2, // 删除状态
			"created_at":  time.Now(),
			"updated_at":  time.Now(),
		}
		_, _ = global.DB.Collection("chapters").InsertOne(ctx, deletedChapter)

		// 统计正常状态
		normalCount, err := repo.CountByStatus(ctx, "book001", 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), normalCount)

		// 统计删除状态
		deletedCount, err := repo.CountByStatus(ctx, "book001", 2)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), deletedCount)
	})
}

func TestChapterRepository_CountVIPChapters(t *testing.T) {
	ctx := setupChapterTest(t)
	repo := reading.NewMongoChapterRepository(global.DB)

	t.Run("统计VIP章节数", func(t *testing.T) {
		// 创建10个章节（第6-10章是VIP）
		for i := 1; i <= 10; i++ {
			createAndInsertChapter(ctx, t, "book001", i)
		}

		vipCount, err := repo.CountVIPChapters(ctx, "book001")
		assert.NoError(t, err)
		assert.Equal(t, int64(5), vipCount) // 第6-10章
	})
}

// ==================== 批量操作测试 ====================

func TestChapterRepository_BatchCreate(t *testing.T) {
	ctx := setupChapterTest(t)
	repo := reading.NewMongoChapterRepository(global.DB)

	t.Run("批量创建章节", func(t *testing.T) {
		chapters := []*reader.Chapter{
			createTestChapter("book001", 1),
			createTestChapter("book001", 2),
			createTestChapter("book001", 3),
		}

		err := repo.BatchCreate(ctx, chapters)
		assert.NoError(t, err)

		// 验证所有章节都有ID和时间戳
		for _, chapter := range chapters {
			assert.NotEmpty(t, chapter.ID)
			assert.False(t, chapter.CreatedAt.IsZero())
			assert.False(t, chapter.UpdatedAt.IsZero())
		}

		// 验证数据库中的数量
		count, _ := repo.CountByBookID(ctx, "book001")
		assert.Equal(t, int64(3), count)
	})

	t.Run("空数组不报错", func(t *testing.T) {
		err := repo.BatchCreate(ctx, []*reader.Chapter{})
		assert.NoError(t, err)
	})
}

func TestChapterRepository_BatchUpdateStatus(t *testing.T) {
	ctx := setupChapterTest(t)
	repo := reading.NewMongoChapterRepository(global.DB)

	t.Run("批量更新章节状态", func(t *testing.T) {
		// 创建3个章节
		chapterIDs := []string{}
		for i := 1; i <= 3; i++ {
			id := createAndInsertChapter(ctx, t, "book001", i)
			chapterIDs = append(chapterIDs, id)
		}

		// 批量更新状态为删除
		err := repo.BatchUpdateStatus(ctx, chapterIDs, 2)
		assert.NoError(t, err)

		// 验证状态已更新
		for _, id := range chapterIDs {
			chapter, _ := repo.GetByID(ctx, id)
			assert.Equal(t, 2, chapter.Status)
		}
	})

	t.Run("空数组不报错", func(t *testing.T) {
		err := repo.BatchUpdateStatus(ctx, []string{}, 2)
		assert.NoError(t, err)
	})
}

func TestChapterRepository_BatchDelete(t *testing.T) {
	ctx := setupChapterTest(t)
	repo := reading.NewMongoChapterRepository(global.DB)

	t.Run("批量删除章节", func(t *testing.T) {
		// 创建3个章节
		chapterIDs := []string{}
		for i := 1; i <= 3; i++ {
			id := createAndInsertChapter(ctx, t, "book001", i)
			chapterIDs = append(chapterIDs, id)
		}

		// 批量删除
		err := repo.BatchDelete(ctx, chapterIDs)
		assert.NoError(t, err)

		// 验证已删除
		for _, id := range chapterIDs {
			_, err := repo.GetByID(ctx, id)
			assert.Error(t, err)
		}
	})

	t.Run("空数组不报错", func(t *testing.T) {
		err := repo.BatchDelete(ctx, []string{})
		assert.NoError(t, err)
	})
}

// ==================== VIP权限检查测试 ====================

func TestChapterRepository_CheckVIPAccess(t *testing.T) {
	ctx := setupChapterTest(t)
	repo := reading.NewMongoChapterRepository(global.DB)

	t.Run("VIP章节", func(t *testing.T) {
		chapterID := createAndInsertChapter(ctx, t, "book001", 10) // 第10章是VIP

		isVIP, err := repo.CheckVIPAccess(ctx, chapterID)
		assert.NoError(t, err)
		assert.True(t, isVIP)
	})

	t.Run("免费章节", func(t *testing.T) {
		chapterID := createAndInsertChapter(ctx, t, "book001", 3) // 第3章是免费

		isVIP, err := repo.CheckVIPAccess(ctx, chapterID)
		assert.NoError(t, err)
		assert.False(t, isVIP)
	})

	t.Run("章节不存在", func(t *testing.T) {
		_, err := repo.CheckVIPAccess(ctx, "nonexistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "章节不存在")
	})
}

func TestChapterRepository_GetChapterPrice(t *testing.T) {
	ctx := setupChapterTest(t)
	repo := reading.NewMongoChapterRepository(global.DB)

	t.Run("获取章节价格", func(t *testing.T) {
		chapterID := createAndInsertChapter(ctx, t, "book001", 5)

		price, err := repo.GetChapterPrice(ctx, chapterID)
		assert.NoError(t, err)
		assert.Equal(t, int64(50), price) // 第5章价格是50（5 * 10）
	})

	t.Run("章节不存在", func(t *testing.T) {
		_, err := repo.GetChapterPrice(ctx, "nonexistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "章节不存在")
	})
}

// ==================== 章节内容管理测试 ====================

func TestChapterRepository_GetChapterContent(t *testing.T) {
	ctx := setupChapterTest(t)
	repo := reading.NewMongoChapterRepository(global.DB)

	t.Run("获取章节内容", func(t *testing.T) {
		chapterID := createAndInsertChapter(ctx, t, "book001", 5)

		content, err := repo.GetChapterContent(ctx, chapterID)
		assert.NoError(t, err)
		assert.Contains(t, content, "这是第5章的内容")
	})

	t.Run("章节不存在", func(t *testing.T) {
		_, err := repo.GetChapterContent(ctx, "nonexistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "章节不存在")
	})
}

func TestChapterRepository_UpdateChapterContent(t *testing.T) {
	ctx := setupChapterTest(t)
	repo := reading.NewMongoChapterRepository(global.DB)

	t.Run("更新章节内容", func(t *testing.T) {
		chapterID := createAndInsertChapter(ctx, t, "book001", 1)

		newContent := "这是新的章节内容，包含更多字数。"
		err := repo.UpdateChapterContent(ctx, chapterID, newContent)
		assert.NoError(t, err)

		// 验证内容已更新
		chapter, _ := repo.GetByID(ctx, chapterID)
		assert.Equal(t, newContent, chapter.Content)

		// 验证字数已更新
		expectedWordCount := len([]rune(newContent))
		assert.Equal(t, expectedWordCount, chapter.WordCount)
	})

	t.Run("更新不存在的章节", func(t *testing.T) {
		err := repo.UpdateChapterContent(ctx, "nonexistent", "新内容")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "章节不存在")
	})
}

// ==================== 健康检查测试 ====================

func TestChapterRepository_Health(t *testing.T) {
	ctx := setupChapterTest(t)
	repo := reading.NewMongoChapterRepository(global.DB)

	t.Run("健康检查", func(t *testing.T) {
		err := repo.Health(ctx)
		assert.NoError(t, err)
	})
}
