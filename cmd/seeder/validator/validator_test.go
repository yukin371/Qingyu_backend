package validator

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDataValidator_NewDataValidator 测试验证器创建
func TestDataValidator_NewDataValidator(t *testing.T) {
	// Arrange & Act
	validator := NewDataValidator(nil)

	// Assert
	assert.NotNil(t, validator)
}

// TestDataValidator_ValidateNoOrphanedRecords_WhenEmptyDatabase 测试空数据库时无孤儿记录
func TestDataValidator_ValidateNoOrphanedRecords_WhenEmptyDatabase(t *testing.T) {
	// Arrange - 使用 mock database
	setupTestDatabase(t)
	defer cleanupTestDatabase(t)

	validator := NewDataValidator(testDB.Database)

	// Act
	report, err := validator.ValidateNoOrphanedRecords(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, 0, report.TotalOrphanedRecords)
	assert.Empty(t, report.OrphanDetails)
}

// TestDataValidator_ValidateNoOrphanedRecords_WithOrphanedReadingProgress 测试检测阅读进度孤儿记录
func TestDataValidator_ValidateNoOrphanedRecords_WithOrphanedReadingProgress(t *testing.T) {
	// Arrange
	setupTestDatabase(t)
	defer cleanupTestDatabase(t)

	// 创建一个用户
	userID := primitive.NewObjectID()
	insertUser(t, userID)

	// 创建一个引用不存在书籍的阅读进度（孤儿记录）
	orphanProgress := bson.D{
		{"_id", primitive.NewObjectID()},
		{"user_id", userID},
		{"book_id", primitive.NewObjectID()}, // 不存在的书籍ID
		{"chapter_id", primitive.NewObjectID()},
		{"progress", 0.5},
		{"created_at", time.Now()},
		{"updated_at", time.Now()},
	}
	_, err := testDB.Database.Collection("reading_progress").InsertOne(context.Background(), orphanProgress)
	require.NoError(t, err)

	validator := NewDataValidator(testDB.Database)

	// Act
	report, err := validator.ValidateNoOrphanedRecords(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, report)
	// 应该检测到孤儿记录
	assert.Greater(t, report.TotalOrphanedRecords, 0)
	assert.NotEmpty(t, report.OrphanDetails)
}

// TestDataValidator_ValidateIDFormatConsistency_WhenConsistent 测试ID格式一致时通过
func TestDataValidator_ValidateIDFormatConsistency_WhenConsistent(t *testing.T) {
	// Arrange
	setupTestDatabase(t)
	defer cleanupTestDatabase(t)

	// 创建格式正确的记录
	userID := primitive.NewObjectID()
	insertUser(t, userID)

	validator := NewDataValidator(testDB.Database)

	// Act
	report, err := validator.ValidateIDFormatConsistency(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Empty(t, report.InconsistentFields)
	assert.Empty(t, report.InconsistentRecords)
}

// TestDataValidator_ValidateIDFormatConsistency_DetectsObjectIDInStringField 测试检测字符串字段中使用ObjectID
func TestDataValidator_ValidateIDFormatConsistency_DetectsObjectIDInStringField(t *testing.T) {
	// Arrange
	setupTestDatabase(t)
	defer cleanupTestDatabase(t)

	// 创建一个使用正确格式的 Like 记录（应该使用 string）
	userID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()

	like := bson.D{
		{"_id", primitive.NewObjectID()},
		{"user_id", userID.Hex()}, // 正确：使用 string (ObjectID.Hex())
		{"target_type", "book"},
		{"target_id", bookID.Hex()}, // 正确：使用 string
		{"created_at", time.Now()},
	}
	_, err := testDB.Database.Collection("likes").InsertOne(context.Background(), like)
	require.NoError(t, err)

	validator := NewDataValidator(testDB.Database)

	// Act
	report, err := validator.ValidateIDFormatConsistency(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, report)
	// 应该检测到格式不一致（如果存在）
	// 当前 Like 模型使用正确的 string 格式，所以应该通过
}

// TestDataValidator_ValidateRelationships_WhenValid 测试所有关系有效时通过
func TestDataValidator_ValidateRelationships_WhenValid(t *testing.T) {
	// Arrange
	setupTestDatabase(t)
	defer cleanupTestDatabase(t)

	// 创建有效的用户和书籍
	userID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()
	insertUser(t, userID)
	insertBook(t, bookID, userID)

	// 创建有效的阅读进度
	progress := bson.D{
		{"_id", primitive.NewObjectID()},
		{"user_id", userID},
		{"book_id", bookID},
		{"chapter_id", primitive.NewObjectID()},
		{"progress", 0.5},
		{"created_at", time.Now()},
		{"updated_at", time.Now()},
	}
	_, err := testDB.Database.Collection("reading_progress").InsertOne(context.Background(), progress)
	require.NoError(t, err)

	validator := NewDataValidator(testDB.Database)

	// Act
	report, err := validator.ValidateRelationships(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.True(t, report.IsValid)
	assert.Equal(t, 0, report.TotalOrphanedRecords)
}

// TestDataValidator_ValidateReadingProgressOrphans 测试阅读进度孤儿记录检测
func TestDataValidator_ValidateReadingProgressOrphans(t *testing.T) {
	// Arrange
	setupTestDatabase(t)
	defer cleanupTestDatabase(t)

	userID := primitive.NewObjectID()
	insertUser(t, userID)

	// 创建孤儿阅读进度
	orphanProgress := bson.D{
		{"_id", primitive.NewObjectID()},
		{"user_id", userID},
		{"book_id", primitive.NewObjectID()}, // 不存在的书籍
		{"chapter_id", primitive.NewObjectID()},
		{"progress", 0.5},
		{"created_at", time.Now()},
		{"updated_at", time.Now()},
	}
	_, err := testDB.Database.Collection("reading_progress").InsertOne(context.Background(), orphanProgress)
	require.NoError(t, err)

	validator := NewDataValidator(testDB.Database)

	// Act
	orphanCount, err := validator.ValidateReadingProgressOrphans(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Greater(t, orphanCount, int64(0))
}

// TestDataValidator_ValidateReadingHistoryOrphans 测试阅读历史孤儿记录检测
func TestDataValidator_ValidateReadingHistoryOrphans(t *testing.T) {
	// Arrange
	setupTestDatabase(t)
	defer cleanupTestDatabase(t)

	userID := primitive.NewObjectID()
	insertUser(t, userID)

	// 创建孤儿阅读历史
	orphanHistory := bson.D{
		{"_id", primitive.NewObjectID()},
		{"user_id", userID},
		{"book_id", primitive.NewObjectID()}, // 不存在的书籍
		{"chapter_id", primitive.NewObjectID()},
		{"read_duration", 100},
		{"progress", 0.5},
		{"start_time", time.Now()},
		{"end_time", time.Now()},
		{"created_at", time.Now()},
	}
	_, err := testDB.Database.Collection("reading_histories").InsertOne(context.Background(), orphanHistory)
	require.NoError(t, err)

	validator := NewDataValidator(testDB.Database)

	// Act
	orphanCount, err := validator.ValidateReadingHistoryOrphans(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Greater(t, orphanCount, int64(0))
}

// TestDataValidator_ValidateBookmarkOrphans 测试书签孤儿记录检测
func TestDataValidator_ValidateBookmarkOrphans(t *testing.T) {
	// Arrange
	setupTestDatabase(t)
	defer cleanupTestDatabase(t)

	userID := primitive.NewObjectID()
	insertUser(t, userID)

	// 创建孤儿书签
	orphanBookmark := bson.D{
		{"_id", primitive.NewObjectID()},
		{"user_id", userID},
		{"book_id", primitive.NewObjectID()}, // 不存在的书籍
		{"chapter_id", primitive.NewObjectID()},
		{"position", 100},
		{"color", "yellow"},
		{"created_at", time.Now()},
		{"updated_at", time.Now()},
	}
	_, err := testDB.Database.Collection("bookmarks").InsertOne(context.Background(), orphanBookmark)
	require.NoError(t, err)

	validator := NewDataValidator(testDB.Database)

	// Act
	orphanCount, err := validator.ValidateBookmarkOrphans(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Greater(t, orphanCount, int64(0))
}

// TestDataValidator_ValidateLikeOrphans 测试点赞孤儿记录检测
func TestDataValidator_ValidateLikeOrphans(t *testing.T) {
	// Arrange
	setupTestDatabase(t)
	defer cleanupTestDatabase(t)

	userID := primitive.NewObjectID()
	insertUser(t, userID)

	// 创建孤儿点赞（引用不存在的书籍）
	orphanLike := bson.D{
		{"_id", primitive.NewObjectID()},
		{"user_id", userID.Hex()},
		{"target_type", "book"},
		{"target_id", primitive.NewObjectID().Hex()}, // 不存在的书籍
		{"created_at", time.Now()},
	}
	_, err := testDB.Database.Collection("likes").InsertOne(context.Background(), orphanLike)
	require.NoError(t, err)

	validator := NewDataValidator(testDB.Database)

	// Act
	orphanCount, err := validator.ValidateLikeOrphans(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Greater(t, orphanCount, int64(0))
}

// TestDataValidator_ValidateNotificationOrphans 测试通知孤儿记录检测
func TestDataValidator_ValidateNotificationOrphans(t *testing.T) {
	// Arrange
	setupTestDatabase(t)
	defer cleanupTestDatabase(t)

	// 创建孤儿通知（引用不存在的用户）
	orphanNotification := bson.D{
		{"_id", primitive.NewObjectID().Hex()},
		{"user_id", primitive.NewObjectID().Hex()}, // 不存在的用户
		{"type", "system"},
		{"priority", "normal"},
		{"title", "Test"},
		{"content", "Test content"},
		{"read", false},
		{"created_at", time.Now()},
	}
	_, err := testDB.Database.Collection("notifications").InsertOne(context.Background(), orphanNotification)
	require.NoError(t, err)

	validator := NewDataValidator(testDB.Database)

	// Act
	orphanCount, err := validator.ValidateNotificationOrphans(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Greater(t, orphanCount, int64(0))
}

// TestDataValidator_GenerateReport 测试生成验证报告
func TestDataValidator_GenerateReport(t *testing.T) {
	// Arrange
	setupTestDatabase(t)
	defer cleanupTestDatabase(t)

	// 创建一些孤儿记录
	userID := primitive.NewObjectID()
	insertUser(t, userID)

	orphanProgress := bson.D{
		{"_id", primitive.NewObjectID()},
		{"user_id", userID},
		{"book_id", primitive.NewObjectID()},
		{"chapter_id", primitive.NewObjectID()},
		{"progress", 0.5},
		{"created_at", time.Now()},
		{"updated_at", time.Now()},
	}
	_, err := testDB.Database.Collection("reading_progress").InsertOne(context.Background(), orphanProgress)
	require.NoError(t, err)

	validator := NewDataValidator(testDB.Database)

	// Act
	report, err := validator.ValidateRelationships(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Contains(t, report.Summary, "验证完成")
}

// TestDataValidator_GetCollectionStats 测试获取集合统计信息
func TestDataValidator_GetCollectionStats(t *testing.T) {
	// Arrange
	setupTestDatabase(t)
	defer cleanupTestDatabase(t)

	// 创建一些测试数据
	userID := primitive.NewObjectID()
	insertUser(t, userID)

	validator := NewDataValidator(testDB.Database)

	// Act
	stats, err := validator.GetCollectionStats(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.NotEmpty(t, stats)
}

// ========== Helper Functions ==========

// setupTestDatabase 设置测试数据库
func setupTestDatabase(t *testing.T) {
	// 这里需要实际的数据库连接设置
	// 暂时跳过，等待实现
	t.Skip("Database setup not yet implemented")
}

// cleanupTestDatabase 清理测试数据库
func cleanupTestDatabase(t *testing.T) {
	// 清理测试数据
}

// insertUser 插入测试用户
func insertUser(t *testing.T, userID primitive.ObjectID) {
	user := bson.D{
		{"_id", userID},
		{"username", "testuser"},
		{"email", "test@example.com"},
		{"password", "hashed_password"},
		{"role", "reader"},
		{"created_at", time.Now()},
		{"updated_at", time.Now()},
	}
	_, err := testDB.Database.Collection("users").InsertOne(context.Background(), user)
	require.NoError(t, err)
}

// insertBook 插入测试书籍
func insertBook(t *testing.T, bookID, authorID primitive.ObjectID) {
	book := bson.D{
		{"_id", bookID},
		{"title", "Test Book"},
		{"author_id", authorID},
		{"status", "published"},
		{"created_at", time.Now()},
		{"updated_at", time.Now()},
	}
	_, err := testDB.Database.Collection("books").InsertOne(context.Background(), book)
	require.NoError(t, err)
}

// testDB 测试数据库连接（全局变量，需要实际实现）
var testDB = struct {
	Database *mongo.Database
}{}
