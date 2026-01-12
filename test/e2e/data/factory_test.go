package data

import (
	"testing"
)

// init 在包加载时初始化测试环境
func init() {
	// 注意：这个init会在所有测试之前运行
	// 但由于Testing.T在init中不可用，我们需要在每个测试中调用SetupTestEnvironment
}

// TestDataFactory_CreateUser 测试创建用户
func TestDataFactory_CreateUser(t *testing.T) {
	SetupTestEnvironment(t)

	factory := NewTestDataFactory(t)

	user := factory.CreateUser(UserOptions{
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

	factory := NewTestDataFactory(t)

	users := factory.CreateUsers(3, UserOptions{
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

	factory := NewTestDataFactory(t)

	// 先创建作者
	author := factory.CreateUser(UserOptions{
		Username: "e2e_test_author",
		Roles:    []string{"reader", "author"},
	})

	book := factory.CreateBook(BookOptions{
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

	factory := NewTestDataFactory(t)

	// 创建作者和书籍
	author := factory.CreateUser(UserOptions{
		Username: "e2e_chapter_author",
		Roles:    []string{"reader", "author"},
	})

	book := factory.CreateBook(BookOptions{
		Title:    "E2E章节测试书籍",
		AuthorID: author.ID,
		Price:    0,
		IsFree:   true,
	})

	chapter := factory.CreateChapter(book.ID.Hex(), 1, true)

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

	factory := NewTestDataFactory(t)

	// 创建用户和书籍
	user := factory.CreateUser(UserOptions{
		Username: "e2e_comment_user",
	})

	author := factory.CreateUser(UserOptions{
		Username: "e2e_comment_author",
		Roles:    []string{"reader", "author"},
	})

	book := factory.CreateBook(BookOptions{
		Title:    "E2E评论测试书籍",
		AuthorID: author.ID,
		Price:    0,
		IsFree:   true,
	})

	comment := factory.CreateComment(CommentOptions{
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

	factory := NewTestDataFactory(t)

	// 创建用户和书籍
	user := factory.CreateUser(UserOptions{
		Username: "e2e_collection_user",
	})

	author := factory.CreateUser(UserOptions{
		Username: "e2e_collection_author",
		Roles:    []string{"reader", "author"},
	})

	book := factory.CreateBook(BookOptions{
		Title:    "E2E收藏测试书籍",
		AuthorID: author.ID,
		Price:    0,
		IsFree:   true,
	})

	collection := factory.CreateCollection(user.ID, book.ID.Hex())

	if collection.ID.IsZero() {
		t.Error("expected collection ID to be set")
	}
	if collection.UserID != user.ID {
		t.Errorf("expected user ID '%s', got '%s'", user.ID, collection.UserID)
	}
}
