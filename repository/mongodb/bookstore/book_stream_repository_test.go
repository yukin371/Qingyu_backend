package mongodb

import (
	"context"
	"testing"

	"Qingyu_backend/models/bookstore"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestStreamSearch 测试流式搜索基础功能
func TestStreamSearch(t *testing.T) {
	t.Skip("需要真实MongoDB连接，跳过单元测试")

	// 创建测试repository
	repo := createTestBookStreamRepository()
	ctx := context.Background()

	t.Run("流式搜索 - 空过滤条件", func(t *testing.T) {
		filter := &bookstore.BookFilter{
			Limit: 10,
		}

		cursor, err := repo.StreamSearch(ctx, filter)
		if err != nil {
			t.Fatalf("StreamSearch failed: %v", err)
		}
		if cursor == nil {
			t.Fatal("cursor should not be nil")
		}
		defer cursor.Close(ctx)

		// 尝试读取一条记录
		hasData := cursor.Next(ctx)
		if !hasData {
			t.Log("No data in collection (expected if empty)")
		}
	})

	t.Run("流式搜索 - 带关键词", func(t *testing.T) {
		keyword := "测试"
		filter := &bookstore.BookFilter{
			Keyword: &keyword,
			Limit:   10,
		}

		cursor, err := repo.StreamSearch(ctx, filter)
		if err != nil {
			t.Fatalf("StreamSearch failed: %v", err)
		}
		if cursor == nil {
			t.Fatal("cursor should not be nil")
		}
		defer cursor.Close(ctx)
	})

	t.Run("流式搜索 - 带排序", func(t *testing.T) {
		filter := &bookstore.BookFilter{
			SortBy:    "created_at",
			SortOrder: "desc",
			Limit:     10,
		}

		cursor, err := repo.StreamSearch(ctx, filter)
		if err != nil {
			t.Fatalf("StreamSearch failed: %v", err)
		}
		if cursor == nil {
			t.Fatal("cursor should not be nil")
		}
		defer cursor.Close(ctx)
	})
}

// TestStreamByCursor 测试基于游标的流式读取
func TestStreamByCursor(t *testing.T) {
	t.Skip("需要真实MongoDB连接，跳过单元测试")
	repo := createTestBookStreamRepository()
	cm := NewCursorManager()
	ctx := context.Background()

	t.Run("使用Timestamp游标继续读取", func(t *testing.T) {
		// 首先获取第一批数据
		filter := &bookstore.BookFilter{
			SortBy:    "created_at",
			SortOrder: "desc",
			Limit:     5,
		}

		cursor, err := repo.StreamSearch(ctx, filter)
		if err != nil {
			t.Fatalf("StreamSearch failed: %v", err)
		}
		defer cursor.Close(ctx)

		// 读取数据并记录最后一个游标
		var lastBook *bookstore.Book
		for cursor.Next(ctx) {
			var book bookstore.Book
			if err := cursor.Decode(&book); err != nil {
				t.Fatalf("Decode failed: %v", err)
			}
			lastBook = &book
		}

		if lastBook == nil {
			t.Skip("No books to test cursor pagination")
		}

		// 生成游标
		nextCursor, err := cm.GenerateNextCursor(lastBook, bookstore.CursorTypeTimestamp, "created_at")
		if err != nil {
			t.Fatalf("GenerateNextCursor failed: %v", err)
		}

		// 使用游标继续读取
		nextFilter := &bookstore.BookFilter{
			SortBy:    "created_at",
			SortOrder: "desc",
			Limit:     5,
			Cursor:    &nextCursor,
		}

		nextCursorResult, err := repo.StreamByCursor(ctx, nextFilter)
		if err != nil {
			t.Fatalf("StreamByCursor failed: %v", err)
		}
		if nextCursorResult == nil {
			t.Fatal("cursor should not be nil")
		}
		defer nextCursorResult.Close(ctx)
	})

	t.Run("使用ID游标继续读取", func(t *testing.T) {
		filter := &bookstore.BookFilter{
			SortBy:    "_id",
			SortOrder: "desc",
			Limit:     5,
		}

		cursor, err := repo.StreamSearch(ctx, filter)
		if err != nil {
			t.Fatalf("StreamSearch failed: %v", err)
		}
		defer cursor.Close(ctx)

		var lastBook *bookstore.Book
		for cursor.Next(ctx) {
			var book bookstore.Book
			if err := cursor.Decode(&book); err != nil {
				t.Fatalf("Decode failed: %v", err)
			}
			lastBook = &book
		}

		if lastBook == nil {
			t.Skip("No books to test ID cursor pagination")
		}

		// 生成ID游标
		nextCursor, err := cm.GenerateNextCursor(lastBook, bookstore.CursorTypeID, "_id")
		if err != nil {
			t.Fatalf("GenerateNextCursor failed: %v", err)
		}

		nextFilter := &bookstore.BookFilter{
			SortBy:    "_id",
			SortOrder: "desc",
			Limit:     5,
			Cursor:    &nextCursor,
		}

		nextCursorResult, err := repo.StreamByCursor(ctx, nextFilter)
		if err != nil {
			t.Fatalf("StreamByCursor failed: %v", err)
		}
		defer nextCursorResult.Close(ctx)
	})

	t.Run("无效游标", func(t *testing.T) {
		invalidCursor := "invalid-cursor"
		filter := &bookstore.BookFilter{
			Limit:  10,
			Cursor: &invalidCursor,
		}

		_, err := repo.StreamByCursor(ctx, filter)
		if err == nil {
			t.Error("expected error for invalid cursor")
		}
	})
}

// TestStreamBatch 测试批量流式读取
func TestStreamBatch(t *testing.T) {
	t.Skip("需要真实MongoDB连接，跳过单元测试")
	repo := createTestBookStreamRepository()
	ctx := context.Background()

	t.Run("批量读取 - 每批10条", func(t *testing.T) {
		batchSize := 10
		filter := &bookstore.BookFilter{
			Limit: batchSize,
		}

		cursor, err := repo.StreamSearch(ctx, filter)
		if err != nil {
			t.Fatalf("StreamSearch failed: %v", err)
		}
		defer cursor.Close(ctx)

		count := 0
		for cursor.Next(ctx) {
			count++
		}

		if count > batchSize {
			t.Errorf("Expected at most %d books, got %d", batchSize, count)
		}
	})

	t.Run("批量读取 - 每批20条", func(t *testing.T) {
		batchSize := 20
		filter := &bookstore.BookFilter{
			Limit: batchSize,
		}

		cursor, err := repo.StreamSearch(ctx, filter)
		if err != nil {
			t.Fatalf("StreamSearch failed: %v", err)
		}
		defer cursor.Close(ctx)

		count := 0
		for cursor.Next(ctx) {
			count++
		}

		if count > batchSize {
			t.Errorf("Expected at most %d books, got %d", batchSize, count)
		}
	})
}

// TestStreamWithFilter 测试带过滤条件的流式搜索
func TestStreamWithFilter(t *testing.T) {
	t.Skip("需要真实MongoDB连接，跳过单元测试")
	repo := createTestBookStreamRepository()
	ctx := context.Background()

	t.Run("按分类流式搜索", func(t *testing.T) {
		categoryID := primitive.NewObjectID().Hex()
		filter := &bookstore.BookFilter{
			CategoryID: &categoryID,
			Limit:      10,
		}

		cursor, err := repo.StreamSearch(ctx, filter)
		if err != nil {
			t.Fatalf("StreamSearch failed: %v", err)
		}
		defer cursor.Close(ctx)
	})

	t.Run("按作者流式搜索", func(t *testing.T) {
		author := "测试作者"
		filter := &bookstore.BookFilter{
			Author: &author,
			Limit:  10,
		}

		cursor, err := repo.StreamSearch(ctx, filter)
		if err != nil {
			t.Fatalf("StreamSearch failed: %v", err)
		}
		defer cursor.Close(ctx)
	})

	t.Run("按状态流式搜索", func(t *testing.T) {
		status := bookstore.BookStatusOngoing
		filter := &bookstore.BookFilter{
			Status: &status,
			Limit:  10,
		}

		cursor, err := repo.StreamSearch(ctx, filter)
		if err != nil {
			t.Fatalf("StreamSearch failed: %v", err)
		}
		defer cursor.Close(ctx)
	})

	t.Run("按标签流式搜索", func(t *testing.T) {
		filter := &bookstore.BookFilter{
			Tags:  []string{"玄幻", "武侠"},
			Limit: 10,
		}

		cursor, err := repo.StreamSearch(ctx, filter)
		if err != nil {
			t.Fatalf("StreamSearch failed: %v", err)
		}
		defer cursor.Close(ctx)
	})
}

// TestStreamErrorHandling 测试错误处理
func TestStreamErrorHandling(t *testing.T) {
	t.Skip("需要真实MongoDB连接，跳过单元测试")
	repo := createTestBookStreamRepository()
	ctx := context.Background()

	t.Run("取消的context", func(t *testing.T) {
		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel() // 立即取消

		filter := &bookstore.BookFilter{
			Limit: 10,
		}

		_, err := repo.StreamSearch(cancelledCtx, filter)
		if err == nil {
			t.Error("expected error for cancelled context")
		}
	})

	t.Run("Limit为0", func(t *testing.T) {
		filter := &bookstore.BookFilter{
			Limit: 0,
		}

		cursor, err := repo.StreamSearch(ctx, filter)
		if err != nil {
			t.Fatalf("StreamSearch failed: %v", err)
		}
		if cursor == nil {
			t.Fatal("cursor should not be nil")
		}
		defer cursor.Close(ctx)
	})
}

// TestStreamPerformance 测试性能相关
func TestStreamPerformance(t *testing.T) {
	t.Skip("需要真实MongoDB连接，跳过单元测试")
	repo := createTestBookStreamRepository()
	ctx := context.Background()

	t.Run("cursor应该及时关闭", func(t *testing.T) {
		filter := &bookstore.BookFilter{
			Limit: 5,
		}

		cursor, err := repo.StreamSearch(ctx, filter)
		if err != nil {
			t.Fatalf("StreamSearch failed: %v", err)
		}

		// 测试cursor可以正确关闭
		if err := cursor.Close(ctx); err != nil {
			t.Errorf("Close failed: %v", err)
		}
	})

	t.Run("流式读取不应占用大量内存", func(t *testing.T) {
		// 这是一个概念性测试，实际内存分析需要更复杂的工具
		filter := &bookstore.BookFilter{
			Limit: 100,
		}

		cursor, err := repo.StreamSearch(ctx, filter)
		if err != nil {
			t.Fatalf("StreamSearch failed: %v", err)
		}
		defer cursor.Close(ctx)

		// 逐条读取，而不是一次性加载所有数据
		count := 0
		for cursor.Next(ctx) {
			count++
			// 在实际应用中，这里处理每条记录后应该可以释放内存
		}

		t.Logf("Streamed %d books", count)
	})
}

// createTestBookStreamRepository 创建测试用的Repository
// 注意：这是简化版，实际集成测试需要真实的MongoDB连接
func createTestBookStreamRepository() *BookStreamRepository {
	// 创建一个mock repository，用于测试编译
	// 实际测试中需要注入真实的MongoDB client
	return &BookStreamRepository{
		baseRepo:  &MongoBookRepository{}, // 空的base repo用于测试编译
		cursorMgr: NewCursorManager(),
	}
}
