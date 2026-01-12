package data

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"Qingyu_backend/global"
	"Qingyu_backend/models/bookstore"
	"Qingyu_backend/models/social"
	"Qingyu_backend/models/users"
	bookRepo "Qingyu_backend/repository/mongodb/bookstore"
	socialRepo "Qingyu_backend/repository/mongodb/social"
	userRepo "Qingyu_backend/repository/mongodb/user"
)

// TestDataFactory 测试数据工厂
type TestDataFactory struct {
	t *testing.T
}

// NewTestDataFactory 创建测试数据工厂
func NewTestDataFactory(t *testing.T) *TestDataFactory {
	return &TestDataFactory{t: t}
}

// UserOptions 用户创建选项
type UserOptions struct {
	Username string
	Email    string
	VIPLevel int
	Balance  float64
	Roles    []string
}

// BookOptions 书籍创建选项
type BookOptions struct {
	Title        string
	AuthorID     string
	Price        float64
	IsFree       bool
	Categories   []string
	WordCount    int
	ChapterCount int
}

// CommentOptions 评论创建选项
type CommentOptions struct {
	AuthorID   string
	TargetID   string
	TargetType string
	Content    string
}

// CreateUser 创建测试用户
func (f *TestDataFactory) CreateUser(opts UserOptions) *users.User {
	userID := primitive.NewObjectID()
	testPassword := "Test1234"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
	require.NoError(f.t, err, "密码哈希失败")

	// 默认值处理
	username := opts.Username
	if username == "" {
		username = fmt.Sprintf("e2e_user_%s", userID.Hex()[:8])
	}

	email := opts.Email
	if email == "" {
		email = fmt.Sprintf("e2e_%s@example.com", userID.Hex()[:8])
	}

	roles := opts.Roles
	if len(roles) == 0 {
		roles = []string{"reader"}
	}

	user := &users.User{
		ID:        userID.Hex(),
		Username:  username,
		Email:     email,
		Password:  string(hashedPassword),
		VIPLevel:  opts.VIPLevel,
		Status:    users.UserStatusActive,
		Roles:     roles,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 清理可能存在的同名用户
	userRepository := userRepo.NewMongoUserRepository(global.DB)
	existingUser, _ := userRepository.GetByUsername(context.Background(), user.Username)
	if existingUser != nil && existingUser.ID != user.ID {
		_ = userRepository.Delete(context.Background(), existingUser.ID)
	}

	err = userRepository.Create(context.Background(), user)
	require.NoError(f.t, err, "创建用户失败")

	return user
}

// CreateUsers 批量创建用户
func (f *TestDataFactory) CreateUsers(count int, baseOptions UserOptions) []*users.User {
	createdUsers := make([]*users.User, count)
	for i := 0; i < count; i++ {
		opts := baseOptions
		if baseOptions.Username != "" {
			opts.Username = fmt.Sprintf("%s_%d", baseOptions.Username, i)
		}
		if baseOptions.Email == "" {
			opts.Email = fmt.Sprintf("e2e_batch_%d@example.com", i+rand.Intn(10000))
		}
		createdUsers[i] = f.CreateUser(opts)
	}
	return createdUsers
}

// CreateBook 创建测试书籍
func (f *TestDataFactory) CreateBook(opts BookOptions) *bookstore.Book {
	bookID := primitive.NewObjectID()
	authorObjID, err := primitive.ObjectIDFromHex(opts.AuthorID)
	require.NoError(f.t, err, "作者ID格式错误")

	// 默认值处理
	title := opts.Title
	if title == "" {
		title = fmt.Sprintf("e2e_book_%s", bookID.Hex()[:8])
	}

	categories := opts.Categories
	if len(categories) == 0 {
		categories = []string{"小说"}
	}

	wordCount := int64(opts.WordCount)
	if wordCount == 0 {
		wordCount = 10000
	}

	book := &bookstore.Book{
		ID:           bookID,
		Title:        title,
		AuthorID:     authorObjID,
		Introduction: "E2E测试书籍 - 用于验证系统功能",
		Categories:   categories,
		Price:        opts.Price,
		Status:       bookstore.BookStatusPublished,
		WordCount:    wordCount,
		IsFree:       opts.IsFree,
		ChapterCount: opts.ChapterCount,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	bookRepository := bookRepo.NewMongoBookRepository(global.DB.Client(), global.DB.Name())
	err = bookRepository.Create(context.Background(), book)
	require.NoError(f.t, err, "创建书籍失败")

	return book
}

// CreateChapter 创建测试章节
func (f *TestDataFactory) CreateChapter(bookID string, chapterNum int, isFree bool) *bookstore.Chapter {
	chapterID := primitive.NewObjectID()
	bookObjID, err := primitive.ObjectIDFromHex(bookID)
	require.NoError(f.t, err, "书籍ID格式错误")

	chapter := &bookstore.Chapter{
		ID:         chapterID,
		BookID:     bookObjID,
		Title:      fmt.Sprintf("第%d章", chapterNum),
		ChapterNum: chapterNum,
		WordCount:  2000,
		IsFree:     isFree,
		Price:      0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	chapter.BeforeCreate()

	chapterRepo := bookRepo.NewMongoChapterRepository(global.DB.Client(), global.DB.Name())
	err = chapterRepo.Create(context.Background(), chapter)
	require.NoError(f.t, err, "创建章节失败")

	// 创建章节内容
	chapterContentRepo := bookRepo.NewMongoChapterContentRepository(global.DB)
	content := fmt.Sprintf("这是第%d章的测试内容。用于验证阅读功能和数据处理流程。", chapterNum)
	chapterContent := &bookstore.ChapterContent{
		ID:        primitive.NewObjectID(),
		ChapterID: chapterID,
		Content:   content,
		Format:    "markdown",
		Version:   1,
		WordCount: len(content),
		CreatedAt: time.Now(),
	}
	chapterContent.BeforeCreate()

	err = chapterContentRepo.Create(context.Background(), chapterContent)
	require.NoError(f.t, err, "创建章节内容失败")

	return chapter
}

// CreateComment 创建测试评论
func (f *TestDataFactory) CreateComment(opts CommentOptions) *social.Comment {
	commentID := primitive.NewObjectID()

	// 默认值处理
	content := opts.Content
	if content == "" {
		content = "这是一条E2E测试评论，用于验证评论系统的功能。"
	}

	now := time.Now()
	comment := &social.Comment{
		IdentifiedEntity: social.IdentifiedEntity{
			ID: commentID.Hex(),
		},
		Timestamps: social.Timestamps{
			CreatedAt: now,
			UpdatedAt: now,
		},
		AuthorID:   opts.AuthorID,
		TargetID:   opts.TargetID,
		TargetType: social.CommentTargetType(opts.TargetType),
		Content:    content,
		State:      social.CommentStateNormal,
	}

	commentRepo := socialRepo.NewMongoCommentRepository(global.DB)
	err := commentRepo.Create(context.Background(), comment)
	require.NoError(f.t, err, "创建评论失败")

	return comment
}

// CreateCollection 创建测试收藏
func (f *TestDataFactory) CreateCollection(userID, bookID string) *social.Collection {
	collectionID := primitive.NewObjectID()

	collection := &social.Collection{
		ID:        collectionID,
		UserID:    userID,
		BookID:    bookID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	collectionRepo := socialRepo.NewMongoCollectionRepository(global.DB)
	err := collectionRepo.Create(context.Background(), collection)
	require.NoError(f.t, err, "创建收藏失败")

	return collection
}

// Cleanup 清理测试数据
func (f *TestDataFactory) Cleanup(prefix string) {
	ctx := context.Background()
	collections := []string{
		"users", "books", "chapters", "chapter_contents",
		"comments", "collections", "likes", "reading_progress",
	}

	for _, collName := range collections {
		// 删除带前缀的数据
		filter := map[string]interface{}{
			"$or": []map[string]interface{}{
				{"username": map[string]interface{}{"$regex": "^" + prefix}},
				{"email": map[string]interface{}{"$regex": "^" + prefix}},
				{"title": map[string]interface{}{"$regex": "^" + prefix}},
			},
		}
		result, _ := global.DB.Collection(collName).DeleteMany(ctx, filter)
		if result.DeletedCount > 0 {
			f.t.Logf("清理 %s: %d 条记录", collName, result.DeletedCount)
		}
	}
}
