package data

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/global"
	"Qingyu_backend/models/bookstore"
)

// TestConsistencyValidator_NormalUserData 测试正常用户数据一致性验证
func TestConsistencyValidator_NormalUserData(t *testing.T) {
	SetupTestEnvironment(t)

	ctx := context.Background()
	factory := NewTestDataFactory(t)
	validator := NewConsistencyValidator(t)

	// 创建有完整数据的用户
	author := factory.CreateUser(ctx, UserOptions{
		Username: "e2e_validator_author",
		Roles:    []string{"reader", "author"},
	})

	book := factory.CreateBook(ctx, BookOptions{
		Title:        "e2e_validator_book",
		AuthorID:     author.ID,
		IsFree:       true,
		ChapterCount: 3,
	})

	// 创建章节
	for i := 0; i < 3; i++ {
		factory.CreateChapter(ctx, book.ID.Hex(), i+1, true)
	}

	// 创建读者
	reader := factory.CreateUser(ctx, UserOptions{
		Username: "e2e_validator_reader",
	})

	// 创建阅读进度
	readerObjID, _ := primitive.ObjectIDFromHex(reader.ID)
	progressID := primitive.NewObjectID()
	readingProgress := bson.M{
		"_id":        progressID,
		"user_id":    reader.ID,
		"book_id":    book.ID.Hex(),
		"chapter_id": book.ID.Hex(),
		"progress":   0.5,
		"created_at": primitive.NewDateTimeFromTime(book.CreatedAt),
		"updated_at": primitive.NewDateTimeFromTime(book.CreatedAt),
	}
	_, err := global.DB.Collection("reading_progress").InsertOne(ctx, readingProgress)
	require.NoError(t, err)

	// 创建评论
	comment := factory.CreateComment(ctx, CommentOptions{
		AuthorID:   reader.ID,
		TargetID:   book.ID.Hex(),
		TargetType: "book",
		Content:    "测试评论",
	})

	// 创建收藏
	collection := factory.CreateCollection(ctx, reader.ID, book.ID.Hex())

	// 创建购买记录
	chapterObjID := book.ID
	purchase := bson.M{
		"_id":           primitive.NewObjectID(),
		"user_id":       readerObjID,
		"book_id":       book.ID,
		"chapter_id":    chapterObjID,
		"price":         10.0,
		"purchase_time": book.CreatedAt,
		"created_at":    book.CreatedAt,
	}
	_, err = global.DB.Collection("chapter_purchases").InsertOne(ctx, purchase)
	require.NoError(t, err)

	t.Run("验证用户数据一致性", func(t *testing.T) {
		issues := validator.ValidateUserData(ctx, reader.ID)
		validator.AssertNoConsistencyIssues(t, issues)
	})

	// 清理
	t.Cleanup(func() {
		global.DB.Collection("reading_progress").DeleteOne(ctx, bson.M{"_id": progressID})
		global.DB.Collection("comments").DeleteOne(ctx, bson.M{"_id": comment.ID})
		global.DB.Collection("collections").DeleteOne(ctx, bson.M{"_id": collection.ID})
		global.DB.Collection("chapter_purchases").DeleteOne(ctx, bson.M{"user_id": readerObjID})
		global.DB.Collection("chapters").DeleteMany(ctx, bson.M{"book_id": book.ID})
		global.DB.Collection("books").DeleteOne(ctx, bson.M{"_id": book.ID})
		global.DB.Collection("users").DeleteOne(ctx, bson.M{"_id": author.ID})
		global.DB.Collection("users").DeleteOne(ctx, bson.M{"_id": reader.ID})
	})
}

// TestConsistencyValidator_OrphanedRecords 测试孤儿记录检测
func TestConsistencyValidator_OrphanedRecords(t *testing.T) {
	SetupTestEnvironment(t)

	ctx := context.Background()
	factory := NewTestDataFactory(t)
	validator := NewConsistencyValidator(t)

	reader := factory.CreateUser(ctx, UserOptions{
		Username: "e2e_validator_orphan_reader",
	})

	// 不存在的书籍ID
	nonExistentBookID := primitive.NewObjectID().Hex()
	nonExistentChapterID := primitive.NewObjectID().Hex()

	t.Run("阅读进度引用不存在的书籍", func(t *testing.T) {
		now := time.Now()
		progressID := primitive.NewObjectID()
		orphanProgress := bson.M{
			"_id":        progressID,
			"user_id":    reader.ID,
			"book_id":    nonExistentBookID,
			"chapter_id": nonExistentChapterID,
			"progress":   0.5,
			"created_at": primitive.NewDateTimeFromTime(now),
			"updated_at": primitive.NewDateTimeFromTime(now),
		}
		_, err := global.DB.Collection("reading_progress").InsertOne(ctx, orphanProgress)
		require.NoError(t, err)

		issues := validator.ValidateUserData(ctx, reader.ID)

		// 应该检测到孤儿记录
		hasBookWarning := false
		for _, issue := range issues {
			if issue.Type == "orphaned_reading_progress_book" {
				hasBookWarning = true
				assert.Equal(t, "warning", issue.Severity)
				assert.Contains(t, issue.Description, "阅读进度引用的书籍不存在")
			}
		}
		assert.True(t, hasBookWarning, "应该检测到阅读进度引用不存在的书籍")

		// 清理
		global.DB.Collection("reading_progress").DeleteOne(ctx, bson.M{"_id": progressID})
	})

	t.Run("评论引用不存在的目标", func(t *testing.T) {
		now := time.Now()
		commentID := primitive.NewObjectID()
		orphanComment := bson.M{
			"_id":         commentID,
			"author_id":   reader.ID,
			"target_id":   nonExistentBookID,
			"target_type": "book",
			"content":     "测试评论",
			"state":       "normal",
			"created_at":  primitive.NewDateTimeFromTime(now),
			"updated_at":  primitive.NewDateTimeFromTime(now),
		}
		_, err := global.DB.Collection("comments").InsertOne(ctx, orphanComment)
		require.NoError(t, err)

		issues := validator.ValidateUserData(ctx, reader.ID)

		// 应该检测到孤儿记录
		hasTargetWarning := false
		for _, issue := range issues {
			if issue.Type == "orphaned_comment_target" {
				hasTargetWarning = true
				assert.Equal(t, "warning", issue.Severity)
				assert.Contains(t, issue.Description, "评论引用的目标不存在")
			}
		}
		assert.True(t, hasTargetWarning, "应该检测到评论引用不存在的目标")

		// 清理
		global.DB.Collection("comments").DeleteOne(ctx, bson.M{"_id": commentID})
	})

	t.Run("收藏引用不存在的书籍", func(t *testing.T) {
		now := time.Now()
		collectionID := primitive.NewObjectID()
		orphanCollection := bson.M{
			"_id":        collectionID,
			"user_id":    reader.ID,
			"book_id":    nonExistentBookID,
			"created_at": primitive.NewDateTimeFromTime(now),
			"updated_at": primitive.NewDateTimeFromTime(now),
		}
		_, err := global.DB.Collection("collections").InsertOne(ctx, orphanCollection)
		require.NoError(t, err)

		issues := validator.ValidateUserData(ctx, reader.ID)

		// 应该检测到孤儿记录
		hasBookWarning := false
		for _, issue := range issues {
			if issue.Type == "orphaned_collection_book" {
				hasBookWarning = true
				assert.Equal(t, "warning", issue.Severity)
				assert.Contains(t, issue.Description, "收藏引用的书籍不存在")
			}
		}
		assert.True(t, hasBookWarning, "应该检测到收藏引用不存在的书籍")

		// 清理
		global.DB.Collection("collections").DeleteOne(ctx, bson.M{"_id": collectionID})
	})

	// 清理
	global.DB.Collection("users").DeleteOne(ctx, bson.M{"_id": reader.ID})
}

// TestConsistencyValidator_NormalBookData 测试正常书籍数据一致性验证
func TestConsistencyValidator_NormalBookData(t *testing.T) {
	SetupTestEnvironment(t)

	ctx := context.Background()
	factory := NewTestDataFactory(t)
	validator := NewConsistencyValidator(t)

	author := factory.CreateUser(ctx, UserOptions{
		Username: "e2e_validator_book_author",
		Roles:    []string{"reader", "author"},
	})

	book := factory.CreateBook(ctx, BookOptions{
		Title:        "e2e_validator_consistent_book",
		AuthorID:     author.ID,
		ChapterCount: 5,
	})

	// 创建正确数量的章节
	for i := 0; i < 5; i++ {
		chapter := factory.CreateChapter(ctx, book.ID.Hex(), i+1, i < 2)
		assert.NotNil(t, chapter)
	}

	t.Run("验证书籍数据一致性", func(t *testing.T) {
		issues := validator.ValidateBookData(ctx, book.ID.Hex())
		validator.AssertNoConsistencyIssues(t, issues)
	})

	// 清理
	t.Cleanup(func() {
		global.DB.Collection("chapters").DeleteMany(ctx, bson.M{"book_id": book.ID})
		global.DB.Collection("chapter_contents").DeleteMany(ctx, bson.M{"chapter_id": bson.M{"$in": func() []primitive.ObjectID {
			cursor, _ := global.DB.Collection("chapters").Find(ctx, bson.M{"book_id": book.ID}, nil)
			var chapters []bookstore.Chapter
			cursor.All(ctx, &chapters)
			ids := make([]primitive.ObjectID, len(chapters))
			for i, c := range chapters {
				ids[i] = c.ID
			}
			return ids
		}()}})
		global.DB.Collection("books").DeleteOne(ctx, bson.M{"_id": book.ID})
		global.DB.Collection("users").DeleteOne(ctx, bson.M{"_id": author.ID})
	})
}

// TestConsistencyValidator_ChapterCountMismatch 测试章节数量不匹配
func TestConsistencyValidator_ChapterCountMismatch(t *testing.T) {
	SetupTestEnvironment(t)

	ctx := context.Background()
	factory := NewTestDataFactory(t)
	validator := NewConsistencyValidator(t)

	author := factory.CreateUser(ctx, UserOptions{
		Username: "e2e_validator_mismatch_author",
		Roles:    []string{"reader", "author"},
	})

	book := factory.CreateBook(ctx, BookOptions{
		Title:        "e2e_validator_mismatch_book",
		AuthorID:     author.ID,
		ChapterCount: 5, // 书籍记录说有5章
	})

	// 只创建3章，制造不一致
	for i := 0; i < 3; i++ {
		factory.CreateChapter(ctx, book.ID.Hex(), i+1, true)
	}

	t.Run("检测章节数量不匹配", func(t *testing.T) {
		issues := validator.ValidateBookData(ctx, book.ID.Hex())

		// 应该检测到章节数量不匹配
		hasCountError := false
		for _, issue := range issues {
			if issue.Type == "chapter_count_mismatch" {
				hasCountError = true
				assert.Equal(t, "error", issue.Severity)
				assert.Contains(t, issue.Description, "章节数量不一致")
				assert.Equal(t, 3, issue.Details["actual_count"])
				assert.Equal(t, 5, issue.Details["expected_count"])
			}
		}
		assert.True(t, hasCountError, "应该检测到章节数量不匹配")
	})

	// 清理
	t.Cleanup(func() {
		global.DB.Collection("chapters").DeleteMany(ctx, bson.M{"book_id": book.ID})
		global.DB.Collection("books").DeleteOne(ctx, bson.M{"_id": book.ID})
		global.DB.Collection("users").DeleteOne(ctx, bson.M{"_id": author.ID})
	})
}

// TestConsistencyValidator_MissingChapterContent 测试章节内容缺失
func TestConsistencyValidator_MissingChapterContent(t *testing.T) {
	SetupTestEnvironment(t)

	ctx := context.Background()
	factory := NewTestDataFactory(t)
	validator := NewConsistencyValidator(t)

	author := factory.CreateUser(ctx, UserOptions{
		Username: "e2e_validator_content_author",
		Roles:    []string{"reader", "author"},
	})

	book := factory.CreateBook(ctx, BookOptions{
		Title:        "e2e_validator_content_book",
		AuthorID:     author.ID,
		ChapterCount: 2,
	})

	// 创建有内容的章节
	factory.CreateChapter(ctx, book.ID.Hex(), 1, true)

	// 创建没有内容的章节（直接插入章节记录，不创建内容）
	now := time.Now()
	chapterWithoutContentID := primitive.NewObjectID()
	chapterWithoutContent := bson.M{
		"_id":         chapterWithoutContentID,
		"book_id":     book.ID,
		"title":       "第二章",
		"chapter_num": 2,
		"word_count":  2000,
		"is_free":     true,
		"price":       0,
		"created_at":  primitive.NewDateTimeFromTime(now),
		"updated_at":  primitive.NewDateTimeFromTime(now),
	}
	_, err := global.DB.Collection("chapters").InsertOne(ctx, chapterWithoutContent)
	require.NoError(t, err)

	t.Run("检测章节内容缺失", func(t *testing.T) {
		issues := validator.ValidateBookData(ctx, book.ID.Hex())

		// 应该检测到章节内容缺失
		hasContentWarning := false
		for _, issue := range issues {
			if issue.Type == "missing_chapter_content" {
				hasContentWarning = true
				assert.Equal(t, "error", issue.Severity)
				assert.Contains(t, issue.Description, "章节缺少内容")
			}
		}
		assert.True(t, hasContentWarning, "应该检测到章节内容缺失")
	})

	// 清理
	t.Cleanup(func() {
		global.DB.Collection("chapters").DeleteMany(ctx, bson.M{"book_id": book.ID})
		global.DB.Collection("books").DeleteOne(ctx, bson.M{"_id": book.ID})
		global.DB.Collection("users").DeleteOne(ctx, bson.M{"_id": author.ID})
	})
}

// TestConsistencyValidator_UserNotExist 测试用户不存在的情况
func TestConsistencyValidator_UserNotExist(t *testing.T) {
	SetupTestEnvironment(t)

	ctx := context.Background()
	validator := NewConsistencyValidator(t)

	nonExistentUserID := primitive.NewObjectID().Hex()

	t.Run("验证不存在的用户", func(t *testing.T) {
		issues := validator.ValidateUserData(ctx, nonExistentUserID)

		// 应该检测到用户不存在
		hasError := false
		for _, issue := range issues {
			if issue.Type == "user_not_found" {
				hasError = true
				assert.Equal(t, "error", issue.Severity)
				assert.Contains(t, issue.Description, "用户不存在")
			}
		}
		assert.True(t, hasError, "应该检测到用户不存在")
	})
}

// TestConsistencyValidator_BookNotExist 测试书籍不存在的情况
func TestConsistencyValidator_BookNotExist(t *testing.T) {
	SetupTestEnvironment(t)

	ctx := context.Background()
	validator := NewConsistencyValidator(t)

	nonExistentBookID := primitive.NewObjectID().Hex()

	t.Run("验证不存在的书籍", func(t *testing.T) {
		issues := validator.ValidateBookData(ctx, nonExistentBookID)

		// 应该检测到书籍不存在
		hasError := false
		for _, issue := range issues {
			if issue.Type == "book_not_found" {
				hasError = true
				assert.Equal(t, "error", issue.Severity)
				assert.Contains(t, issue.Description, "书籍不存在")
			}
		}
		assert.True(t, hasError, "应该检测到书籍不存在")
	})
}
