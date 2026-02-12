package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// setupTestDB 创建测试数据库连接
func setupTestDB(t *testing.T) *mongo.Client {
	t.Helper()

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
		t.Skip("跳过测试: 需要MongoDB测试环境 (设置MONGO_URI环境变量)")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	require.NoError(t, err)

	// 验证连接
	require.NoError(t, client.Ping(ctx, nil))

	t.Cleanup(func() {
		if err := client.Disconnect(context.Background()); err != nil {
			t.Errorf("关闭数据库连接失败: %v", err)
		}
	})

	return client
}

// getIndexNames 获取集合的所有索引名称
func getIndexNames(t *testing.T, client *mongo.Client, dbName, collectionName string) []string {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(dbName).Collection(collectionName)
	cursor, err := collection.Indexes().List(ctx)
	require.NoError(t, err)

	var indexes []bson.M
	require.NoError(t, cursor.All(ctx, &indexes))

	names := make([]string, 0, len(indexes))
	for _, idx := range indexes {
		if name, ok := idx["name"].(string); ok {
			names = append(names, name)
		}
	}

	return names
}

// assertIndexExists 验证索引是否存在
func assertIndexExists(t *testing.T, client *mongo.Client, dbName, collectionName, expectedIndexName string) {
	t.Helper()

	indexes := getIndexNames(t, client, dbName, collectionName)

	found := false
	for _, name := range indexes {
		if name == expectedIndexName {
			found = true
			break
		}
	}

	require.True(t, found, "索引 %s 在集合 %s.%s 中不存在，当前索引: %v",
		expectedIndexName, dbName, collectionName, indexes)
}

// TestVerifyIndexes_Users 验证Users集合索引
func TestVerifyIndexes_Users(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	client := setupTestDB(t)
	dbName := "qingyu_test"
	collectionName := "users"

	t.Run("验证status_created_at索引", func(t *testing.T) {
		assertIndexExists(t, client, dbName, collectionName, "status_1_created_at_-1")
	})

	t.Run("验证roles索引", func(t *testing.T) {
		assertIndexExists(t, client, dbName, collectionName, "roles_1")
	})

	t.Run("验证last_login_at索引", func(t *testing.T) {
		assertIndexExists(t, client, dbName, collectionName, "last_login_at_-1")
	})
}

// TestVerifyIndexes_BooksP0 验证Books集合P0索引
func TestVerifyIndexes_BooksP0(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	client := setupTestDB(t)
	dbName := "qingyu_test"
	collectionName := "books"

	t.Run("验证status_created_at索引", func(t *testing.T) {
		assertIndexExists(t, client, dbName, collectionName, "status_1_created_at_-1")
	})

	t.Run("验证status_rating索引", func(t *testing.T) {
		assertIndexExists(t, client, dbName, collectionName, "status_1_rating_-1")
	})

	t.Run("验证author_id_status_created_at复合索引", func(t *testing.T) {
		assertIndexExists(t, client, dbName, collectionName, "author_id_1_status_1_created_at_-1")
	})

	t.Run("验证category_ids_rating复合索引", func(t *testing.T) {
		assertIndexExists(t, client, dbName, collectionName, "category_ids_1_rating_-1")
	})

	t.Run("验证is_completed_status复合索引", func(t *testing.T) {
		assertIndexExists(t, client, dbName, collectionName, "is_completed_1_status_1")
	})
}

// TestVerifyIndexes_Chapters 验证Chapters集合索引
func TestVerifyIndexes_Chapters(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	client := setupTestDB(t)
	dbName := "qingyu_test"
	collectionName := "chapters"

	t.Run("验证book_id_chapter_num复合索引", func(t *testing.T) {
		assertIndexExists(t, client, dbName, collectionName, "book_id_1_chapter_num_1")
	})

	t.Run("验证book_id_status_chapter_num复合索引", func(t *testing.T) {
		assertIndexExists(t, client, dbName, collectionName, "book_id_1_status_1_chapter_num_1")
	})
}

// TestVerifyIndexes_ReadingProgress 验证ReadingProgress集合索引
func TestVerifyIndexes_ReadingProgress(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	client := setupTestDB(t)
	dbName := "qingyu_test"
	collectionName := "reading_progress"

	t.Run("验证user_id_book_id复合索引", func(t *testing.T) {
		assertIndexExists(t, client, dbName, collectionName, "user_id_1_book_id_1")
	})

	t.Run("验证user_id_last_read_at索引", func(t *testing.T) {
		assertIndexExists(t, client, dbName, collectionName, "user_id_1_last_read_at_-1")
	})

	t.Run("验证book_id索引", func(t *testing.T) {
		assertIndexExists(t, client, dbName, collectionName, "book_id_1")
	})
}

// TestListAllIndexes 列出所有集合的索引（调试用）
func TestListAllIndexes(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	client := setupTestDB(t)
	dbName := "qingyu_test"

	collections := []string{"users", "books", "chapters", "reading_progress"}

	for _, collName := range collections {
		t.Run(collName, func(t *testing.T) {
			indexes := getIndexNames(t, client, dbName, collName)
			t.Logf("集合 %s 的索引: %v", collName, indexes)
		})
	}
}
