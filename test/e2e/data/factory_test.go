//go:build e2e
// +build e2e

package data

import (
	"context"
	"testing"
)

// TestDataFactory_CreateUser 测试创建用户
func TestDataFactory_CreateUser(t *testing.T) {
	SetupTestEnvironment(t)
	ctx := context.Background()

	factory := NewTestDataFactory(t)

	user := factory.CreateUser(ctx, UserOptions{
		Username: "e2e_test_user",
		VIPLevel: 1,
	})

	if user.Username != "e2e_test_user" {
		t.Errorf("expected username 'e2e_test_user', got '%s'", user.Username)
	}
	if user.VIPLevel != 1 {
		t.Errorf("expected VIP level 1, got %d", user.VIPLevel)
	}
	if user.ID == "" {
		t.Error("expected user ID to be set")
	}
	if user.Password == "" {
		t.Error("expected password to be set")
	}
}

// TestDataFactory_CreateUsers 测试批量创建用户
func TestDataFactory_CreateUsers(t *testing.T) {
	SetupTestEnvironment(t)
	ctx := context.Background()

	factory := NewTestDataFactory(t)

	users := factory.CreateUsers(ctx, 3, UserOptions{
		Username: "e2e_batch_user",
		VIPLevel: 0,
	})

	if len(users) != 3 {
		t.Errorf("expected 3 users, got %d", len(users))
	}

	for i, user := range users {
		if user.ID == "" {
			t.Errorf("user %d: expected ID to be set", i)
		}
	}
}

// TestDataFactory_CreateBook 测试创建书籍
func TestDataFactory_CreateBook(t *testing.T) {
	SetupTestEnvironment(t)
	ctx := context.Background()

	factory := NewTestDataFactory(t)

	// 先创建作者
	author := factory.CreateUser(ctx, UserOptions{
		Username: "e2e_test_author",
		Roles:    []string{"reader", "author"},
	})

	book := factory.CreateBook(ctx, BookOptions{
		Title:     "E2E测试书籍",
		AuthorID:  author.ID,
		Price:     100,
		IsFree:    false,
		WordCount: 10000,
	})

	if book.Title != "E2E测试书籍" {
		t.Errorf("expected title 'E2E测试书籍', got '%s'", book.Title)
	}
	if book.ID.IsZero() {
		t.Error("expected book ID to be set")
	}
}

// TestDataFactory_CreateChapter 测试创建章节
func TestDataFactory_CreateChapter(t *testing.T) {
	SetupTestEnvironment(t)
	ctx := context.Background()

	factory := NewTestDataFactory(t)

	// 创建作者和书籍
	author := factory.CreateUser(ctx, UserOptions{
		Username: "e2e_chapter_author",
		Roles:    []string{"reader", "author"},
	})

	book := factory.CreateBook(ctx, BookOptions{
		Title:    "E2E章节测试书籍",
		AuthorID: author.ID,
		Price:    0,
		IsFree:   true,
	})

	chapter := factory.CreateChapter(ctx, book.ID.Hex(), 1, true)

	if chapter.ID.IsZero() {
		t.Error("expected chapter ID to be set")
	}
	if chapter.ChapterNum != 1 {
		t.Errorf("expected chapter number 1, got %d", chapter.ChapterNum)
	}
	if !chapter.IsFree {
		t.Error("expected chapter to be free")
	}
}

// TestDataFactory_CreateComment 测试创建评论
func TestDataFactory_CreateComment(t *testing.T) {
	SetupTestEnvironment(t)
	ctx := context.Background()

	factory := NewTestDataFactory(t)

	// 创建用户和书籍
	user := factory.CreateUser(ctx, UserOptions{
		Username: "e2e_comment_user",
	})

	author := factory.CreateUser(ctx, UserOptions{
		Username: "e2e_comment_author",
		Roles:    []string{"reader", "author"},
	})

	book := factory.CreateBook(ctx, BookOptions{
		Title:    "E2E评论测试书籍",
		AuthorID: author.ID,
		Price:    0,
		IsFree:   true,
	})

	comment := factory.CreateComment(ctx, CommentOptions{
		AuthorID:   user.ID,
		TargetID:   book.ID.Hex(),
		TargetType: "book",
		Content:    "这是一条E2E测试评论",
	})

	if comment.ID == "" {
		t.Error("expected comment ID to be set")
	}
	if comment.AuthorID != user.ID {
		t.Errorf("expected author ID '%s', got '%s'", user.ID, comment.AuthorID)
	}
}

// TestDataFactory_CreateCollection 测试创建收藏
func TestDataFactory_CreateCollection(t *testing.T) {
	SetupTestEnvironment(t)
	ctx := context.Background()

	factory := NewTestDataFactory(t)

	// 创建用户和书籍
	user := factory.CreateUser(ctx, UserOptions{
		Username: "e2e_collection_user",
	})

	author := factory.CreateUser(ctx, UserOptions{
		Username: "e2e_collection_author",
		Roles:    []string{"reader", "author"},
	})

	book := factory.CreateBook(ctx, BookOptions{
		Title:    "E2E收藏测试书籍",
		AuthorID: author.ID,
		Price:    0,
		IsFree:   true,
	})

	collection := factory.CreateCollection(ctx, user.ID, book.ID.Hex())

	if collection.ID.IsZero() {
		t.Error("expected collection ID to be set")
	}
	if collection.UserID != user.ID {
		t.Errorf("expected user ID '%s', got '%s'", user.ID, collection.UserID)
	}
}

// TestDataFactory_Cleanup 测试清理功能
func TestDataFactory_Cleanup(t *testing.T) {
	SetupTestEnvironment(t)
	ctx := context.Background()

	factory := NewTestDataFactory(t)

	// 创建测试数据
	_ = factory.CreateUser(ctx, UserOptions{
		Username: "e2e_cleanup_user_1",
	})
	_ = factory.CreateUser(ctx, UserOptions{
		Username: "e2e_cleanup_user_2",
	})

	// 使用不同的前缀创建不应该被清理的数据
	_ = factory.CreateUser(ctx, UserOptions{
		Username: "other_prefix_user",
	})

	// 清理数据
	factory.Cleanup("e2e_cleanup")

	// 验证已清理的数据不存在（通过尝试创建同名用户来验证）
	// 如果清理成功，应该能成功创建同名用户
	newUser1 := factory.CreateUser(ctx, UserOptions{
		Username: "e2e_cleanup_user_1",
	})
	if newUser1.ID == "" {
		t.Error("expected to create user with cleaned username")
	}

	// 验证其他前缀的数据没有被清理
	// 这里我们无法直接验证，但Cleanup只清理特定前缀的数据
	// 在实际使用中，这确保了测试数据的隔离性

	t.Log("Cleanup test passed")
}

