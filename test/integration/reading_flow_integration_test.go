package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/config"
	"Qingyu_backend/core"
	"Qingyu_backend/global"
	bookModel "Qingyu_backend/models/bookstore"
)

// TestCompleteReadingFlow 完整阅读流程集成测试
// 测试流程: 书城 -> 榜单 -> 书籍详情 -> 章节列表 -> 阅读章节 -> 添加书签/笔记
func TestCompleteReadingFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试（使用 -short 标志）")
	}

	// 1. 初始化配置和数据库
	_, err := config.LoadConfig("../..")
	require.NoError(t, err, "加载配置失败")

	err = core.InitDB()
	require.NoError(t, err, "初始化数据库失败")

	ctx := context.Background()

	// 清理测试数据
	defer cleanupReadingFlowTestData(t, ctx)

	// 2. 准备测试数据
	testUserID := "test_reader_" + time.Now().Format("20060102150405")
	testBook := createTestBook(t, ctx)
	testChapters := createTestChapters(t, ctx, testBook.ID.Hex())
	createTestRanking(t, ctx, testBook.ID)

	// ========== 阶段1：书城首页 - 获取榜单 ==========
	t.Run("书城首页_获取热门榜单", func(t *testing.T) {
		// 获取榜单数据
		cursor, err := global.DB.Collection("ranking_items").Find(ctx, bson.M{
			"type": "realtime",
		})
		require.NoError(t, err)

		var rankings []bookModel.RankingItem
		err = cursor.All(ctx, &rankings)
		require.NoError(t, err, "应该能获取到热门榜单")
		assert.NotEmpty(t, rankings, "榜单应该包含书籍")
		assert.Equal(t, testBook.ID, rankings[0].BookID, "榜单第一本应该是测试书籍")

		t.Logf("✓ 成功获取热门榜单，包含 %d 本书籍", len(rankings))
	})

	// ========== 阶段2：从榜单点击进入书籍详情 ==========
	t.Run("从榜单进入书籍详情", func(t *testing.T) {
		// 模拟从榜单获取书籍ID
		var rankingItem bookModel.RankingItem
		err := global.DB.Collection("ranking_items").FindOne(ctx, bson.M{
			"type": "realtime",
		}).Decode(&rankingItem)
		require.NoError(t, err)

		bookID := rankingItem.BookID

		// 获取书籍详情
		var book bookModel.Book
		err = global.DB.Collection("books").FindOne(ctx, bson.M{
			"_id": bookID,
		}).Decode(&book)

		require.NoError(t, err, "应该能获取到书籍详情")
		assert.Equal(t, testBook.Title, book.Title, "书籍标题应该匹配")
		assert.Equal(t, testBook.Author, book.Author, "作者应该匹配")
		assert.Greater(t, book.WordCount, int64(0), "字数应该大于0")

		t.Logf("✓ 成功获取书籍详情: %s (作者: %s, 字数: %d)", book.Title, book.Author, book.WordCount)
	})

	// ========== 阶段3：查看书籍章节列表 ==========
	t.Run("查看书籍章节列表", func(t *testing.T) {
		// 获取章节列表
		cursor, err := global.DB.Collection("chapters").Find(ctx, bson.M{
			"book_id": testBook.ID,
		})
		require.NoError(t, err)

		var chapters []bookModel.Chapter
		err = cursor.All(ctx, &chapters)
		require.NoError(t, err)

		assert.Len(t, chapters, 3, "应该有3个章节")
		assert.Equal(t, "第一章", chapters[0].Title, "第一章标题应该匹配")

		t.Logf("✓ 成功获取章节列表，共 %d 章", len(chapters))
		for i, ch := range chapters {
			t.Logf("  章节%d: %s (字数: %d, 免费: %v)", i+1, ch.Title, ch.WordCount, ch.IsFree)
		}
	})

	// ========== 阶段4：阅读第一章 ==========
	t.Run("阅读第一章节内容", func(t *testing.T) {
		// 获取第一章
		chapter := testChapters[0]

		// 验证章节信息（Content 字段已移至 ChapterContent）
		assert.NotEmpty(t, chapter.Title, "章节标题不应为空")
		assert.Equal(t, "第一章", chapter.Title, "章节标题应该匹配")
		assert.True(t, chapter.IsFree, "第一章应该是免费章节")
		assert.Greater(t, chapter.WordCount, 0, "章节字数应该大于0")

		// 记录阅读进度（模拟）
		progress := bson.M{
			"user_id":    testUserID,
			"book_id":    testBook.ID.Hex(),
			"chapter_id": chapter.ID,
			"position":   50, // 读到50%
			"updated_at": time.Now(),
		}

		_, err := global.DB.Collection("reading_progress").InsertOne(ctx, progress)
		require.NoError(t, err, "记录阅读进度应该成功")

		t.Logf("✓ 成功阅读第一章: %s (字数: %d)", chapter.Title, chapter.WordCount)
		t.Logf("✓ 阅读进度已保存: 50%%")
	})

	// ========== 阶段5：添加书签和笔记 ==========
	t.Run("添加书签和阅读笔记", func(t *testing.T) {
		chapter := testChapters[0]

		// 添加书签
		bookmark := bson.M{
			"_id":        "bookmark_" + time.Now().Format("20060102150405"),
			"user_id":    testUserID,
			"book_id":    testBook.ID.Hex(),
			"chapter_id": chapter.ID,
			"type":       "bookmark",
			"range":      "100-150",
			"text":       "重要情节",
			"note":       "",
			"created_at": time.Now(),
			"updated_at": time.Now(),
		}

		_, err := global.DB.Collection("annotations").InsertOne(ctx, bookmark)
		require.NoError(t, err, "添加书签应该成功")

		// 添加笔记
		note := bson.M{
			"_id":        "note_" + time.Now().Format("20060102150405"),
			"user_id":    testUserID,
			"book_id":    testBook.ID.Hex(),
			"chapter_id": chapter.ID,
			"type":       "note",
			"range":      "200-250",
			"text":       "精彩描写",
			"note":       "这段描写很生动",
			"created_at": time.Now(),
			"updated_at": time.Now(),
		}

		_, err = global.DB.Collection("annotations").InsertOne(ctx, note)
		require.NoError(t, err, "添加笔记应该成功")

		t.Logf("✓ 成功添加书签")
		t.Logf("✓ 成功添加阅读笔记")
	})

	// ========== 阶段6：继续阅读第二章（付费章节） ==========
	t.Run("检查付费章节访问控制", func(t *testing.T) {
		chapter := testChapters[1] // 第二章是付费章节

		assert.False(t, chapter.IsFree, "第二章应该是付费章节")
		assert.Greater(t, chapter.Price, 0.0, "付费章节应该有价格")

		// 模拟用户访问付费章节
		// 这里仅验证章节的付费属性，实际访问控制在Service层
		t.Logf("✓ 确认第二章为付费章节: %s (价格: %.2f)", chapter.Title, chapter.Price)
		t.Logf("  需要VIP权限或单独购买才能阅读")
	})

	// ========== 阶段7：查看阅读历史 ==========
	t.Run("查看阅读历史", func(t *testing.T) {
		// 获取用户的阅读进度
		var progress bson.M
		err := global.DB.Collection("reading_progress").FindOne(ctx, bson.M{
			"user_id": testUserID,
			"book_id": testBook.ID.Hex(),
		}).Decode(&progress)

		require.NoError(t, err, "应该能获取到阅读进度")
		assert.Equal(t, testBook.ID.Hex(), progress["book_id"], "书籍ID应该匹配")
		assert.Equal(t, testChapters[0].ID, progress["chapter_id"], "章节ID应该匹配")
		assert.Equal(t, int32(50), progress["position"], "阅读位置应该是50%")

		t.Logf("✓ 成功获取阅读历史")
		t.Logf("  当前阅读: %s - %s", testBook.Title, testChapters[0].Title)
		t.Logf("  阅读进度: %d%%", progress["position"])
	})

	// ========== 阶段8：查看我的书签和笔记 ==========
	t.Run("查看书签和笔记列表", func(t *testing.T) {
		// 获取书签
		cursor, err := global.DB.Collection("annotations").Find(ctx, bson.M{
			"user_id": testUserID,
			"type":    "bookmark",
		})
		require.NoError(t, err)

		var bookmarks []bson.M
		err = cursor.All(ctx, &bookmarks)
		require.NoError(t, err)
		assert.Len(t, bookmarks, 1, "应该有1个书签")

		// 获取笔记
		cursor, err = global.DB.Collection("annotations").Find(ctx, bson.M{
			"user_id": testUserID,
			"type":    "note",
		})
		require.NoError(t, err)

		var notes []bson.M
		err = cursor.All(ctx, &notes)
		require.NoError(t, err)
		assert.Len(t, notes, 1, "应该有1条笔记")

		t.Logf("✓ 成功获取书签列表: %d 个", len(bookmarks))
		t.Logf("✓ 成功获取笔记列表: %d 条", len(notes))
	})

	t.Logf("\n=== 完整阅读流程测试通过 ===")
	t.Logf("流程: 书城榜单 -> 书籍详情 -> 章节列表 -> 阅读章节 -> 添加笔记 -> 查看历史")
}

// createTestBook 创建测试书籍
func createTestBook(t *testing.T, ctx context.Context) *bookModel.Book {
	now := time.Now()
	book := &bookModel.Book{
		Title:        "测试小说-完整流程测试",
		Author:       "测试作者",
		Cover:        "http://example.com/cover.jpg",
		Introduction: "这是一本用于完整流程测试的小说",
		Categories:   []string{"玄幻"},
		Tags:         []string{"测试", "玄幻"},
		Status:       bookModel.BookStatusOngoing,
		WordCount:    30000,
		ChapterCount: 3,
		LastUpdateAt: &now,
	}

	result, err := global.DB.Collection("books").InsertOne(ctx, book)
	require.NoError(t, err)

	book.ID = result.InsertedID.(primitive.ObjectID)
	return book
}

// createTestChapters 创建测试章节
func createTestChapters(t *testing.T, ctx context.Context, bookIDHex string) []bookModel.Chapter {
	now := time.Now()
	chapters := []bookModel.Chapter{
		{
			BookID:      bookIDHex,
			Title:       "第一章",
			ChapterNum:  1,
			WordCount:   1000,
			IsFree:      true,
			Price:       0,
			PublishTime: now,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			BookID:      bookIDHex,
			Title:       "第二章",
			ChapterNum:  2,
			WordCount:   2000,
			IsFree:      false,
			Price:       10,
			PublishTime: now,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			BookID:      bookIDHex,
			Title:       "第三章",
			ChapterNum:  3,
			WordCount:   1500,
			IsFree:      false,
			Price:       10,
			PublishTime: now,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	for i := range chapters {
		result, err := global.DB.Collection("chapters").InsertOne(ctx, &chapters[i])
		require.NoError(t, err)
		// Chapter.ID 是 string 类型，从 InsertedID 获取
		insertedID := result.InsertedID.(primitive.ObjectID)
		chapters[i].ID = insertedID.Hex()
	}

	return chapters
}

// createTestRanking 创建测试榜单
func createTestRanking(t *testing.T, ctx context.Context, bookID primitive.ObjectID) {
	now := time.Now()
	period := bookModel.GetPeriodString(bookModel.RankingTypeRealtime, now)

	ranking := bookModel.RankingItem{
		BookID:    bookID,
		Type:      bookModel.RankingTypeRealtime,
		Rank:      1,
		Score:     95.5,
		ViewCount: 1000,
		LikeCount: 100,
		Period:    period,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err := global.DB.Collection("ranking_items").InsertOne(ctx, ranking)
	require.NoError(t, err)
}

// generateContent 生成指定长度的内容
func generateContent(length int) string {
	content := ""
	for i := 0; i < length/10; i++ {
		content += "测试内容。"
	}
	return content
}

// cleanupReadingFlowTestData 清理测试数据
func cleanupReadingFlowTestData(t *testing.T, ctx context.Context) {
	// 删除测试书籍
	global.DB.Collection("books").DeleteMany(ctx, bson.M{
		"title": bson.M{"$regex": "测试小说-完整流程测试"},
	})

	// 删除测试章节
	global.DB.Collection("chapters").DeleteMany(ctx, bson.M{
		"title":      bson.M{"$regex": "^第.*章$"},
		"created_at": bson.M{"$gte": time.Now().Add(-1 * time.Hour)},
	})

	// 删除测试榜单
	global.DB.Collection("ranking_items").DeleteMany(ctx, bson.M{
		"type":       "realtime",
		"updated_at": bson.M{"$gte": time.Now().Add(-1 * time.Hour)},
	})

	// 删除测试阅读进度
	global.DB.Collection("reading_progress").DeleteMany(ctx, bson.M{
		"user_id": bson.M{"$regex": "^test_reader_"},
	})

	// 删除测试标注
	global.DB.Collection("annotations").DeleteMany(ctx, bson.M{
		"user_id": bson.M{"$regex": "^test_reader_"},
	})
}

// TestBookstoreToReading 书城浏览到阅读的流程测试
func TestBookstoreToReading(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试（使用 -short 标志）")
	}

	_, err := config.LoadConfig("../..")
	require.NoError(t, err)

	err = core.InitDB()
	require.NoError(t, err)

	ctx := context.Background()
	defer cleanupReadingFlowTestData(t, ctx)

	// 创建测试数据
	testBook := createTestBook(t, ctx)
	createTestChapters(t, ctx, testBook.ID.Hex())

	t.Run("书城分类浏览", func(t *testing.T) {
		// 按分类查询书籍（categories是数组）
		cursor, err := global.DB.Collection("books").Find(ctx, bson.M{
			"categories": bson.M{"$in": []string{"玄幻"}},
		})
		require.NoError(t, err)

		var books []bookModel.Book
		err = cursor.All(ctx, &books)
		require.NoError(t, err)

		assert.Greater(t, len(books), 0, "应该能找到玄幻类书籍")
		t.Logf("✓ 找到 %d 本玄幻类书籍", len(books))
	})

	t.Run("书籍搜索", func(t *testing.T) {
		// 按标题搜索
		cursor, err := global.DB.Collection("books").Find(ctx, bson.M{
			"title": bson.M{"$regex": "测试小说", "$options": "i"},
		})
		require.NoError(t, err)

		var books []bookModel.Book
		err = cursor.All(ctx, &books)
		require.NoError(t, err)

		assert.Len(t, books, 1, "应该找到1本书")
		assert.Equal(t, testBook.Title, books[0].Title)

		t.Logf("✓ 搜索成功: %s", books[0].Title)
	})

	t.Run("书籍收藏", func(t *testing.T) {
		testUserID := "test_user_collect"

		// 添加到收藏
		collection := bson.M{
			"user_id":    testUserID,
			"book_id":    testBook.ID.Hex(),
			"created_at": time.Now(),
		}

		_, err := global.DB.Collection("user_collections").InsertOne(ctx, collection)
		require.NoError(t, err)

		// 验证收藏
		var result bson.M
		err = global.DB.Collection("user_collections").FindOne(ctx, bson.M{
			"user_id": testUserID,
			"book_id": testBook.ID.Hex(),
		}).Decode(&result)

		require.NoError(t, err)
		t.Logf("✓ 收藏成功")

		// 清理
		global.DB.Collection("user_collections").DeleteOne(ctx, bson.M{
			"user_id": testUserID,
		})
	})
}
