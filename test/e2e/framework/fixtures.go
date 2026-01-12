package e2e

import (
	"context"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"Qingyu_backend/global"
	"Qingyu_backend/models/bookstore"
	"Qingyu_backend/models/users"
	bookstoreRepo "Qingyu_backend/repository/mongodb/bookstore"
	userRepo "Qingyu_backend/repository/mongodb/user"
)

// Fixtures 测试数据夹具
type Fixtures struct {
	env *TestEnvironment
}

// NewFixtures 创建测试数据夹具
func (env *TestEnvironment) Fixtures() *Fixtures {
	return &Fixtures{env: env}
}

// ============ 用户选项 ============

// UserOption 用户选项函数
type UserOption func(*users.User)

// WithUsername 设置用户名
func WithUsername(username string) UserOption {
	return func(u *users.User) { u.Username = username }
}

// WithEmail 设置邮箱
func WithEmail(email string) UserOption {
	return func(u *users.User) { u.Email = email }
}

// WithVIPLevel 设置VIP等级
func WithVIPLevel(level int) UserOption {
	return func(u *users.User) { u.VIPLevel = level }
}

// ============ 书籍选项 ============

// BookOption 书籍选项函数
type BookOption func(*bookstore.Book)

// WithBookTitle 设置书名
func WithBookTitle(title string) BookOption {
	return func(b *bookstore.Book) { b.Title = title }
}

// WithBookPrice 设置价格
func WithBookPrice(price float64) BookOption {
	return func(b *bookstore.Book) { b.Price = price }
}

// WithBookCategory 设置分类
func WithBookCategory(category string) BookOption {
	return func(b *bookstore.Book) { b.Categories = []string{category} }
}

// ============ 章节选项 ============

// ChapterOption 章节选项函数
type ChapterOption func(*bookstore.Chapter)

// WithChapterTitle 设置章节标题
func WithChapterTitle(title string) ChapterOption {
	return func(c *bookstore.Chapter) { c.Title = title }
}

// WithChapterFree 设置是否免费
func WithChapterFree(isFree bool) ChapterOption {
	return func(c *bookstore.Chapter) { c.IsFree = isFree }
}

// WithChapterPrice 设置章节价格
func WithChapterPrice(price float64) ChapterOption {
	return func(c *bookstore.Chapter) { c.Price = price }
}

// ============ 创建测试数据 ============

// CreateUser 创建测试用户（带 e2e_test_ 前缀）
func (f *Fixtures) CreateUser(opts ...UserOption) *users.User {
	userID := primitive.NewObjectID()
	testPassword := "Test1234"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
	require.NoError(f.env.T, err, "密码哈希失败")

	// 使用完整的 ID 确保用户名唯一
	uniqueSuffix := userID.Hex()
	user := &users.User{
		ID:       userID.Hex(), // User.ID 是 string 类型
		Username: "e2e_test_user_" + uniqueSuffix,
		Email:    "e2e_test_" + uniqueSuffix + "@example.com",
		Password: string(hashedPassword),
		VIPLevel: 0,
		Status:   users.UserStatusActive,
		Roles:    []string{"reader"},
	}

	// 应用选项（注意：WithUsername 会覆盖上面的用户名）
	for _, opt := range opts {
		opt(user)
	}

	// 清理可能存在的同名用户（从之前的测试运行遗留）
	userRepository := userRepo.NewMongoUserRepository(global.DB)
	existingUser, _ := userRepository.GetByUsername(context.Background(), user.Username)
	if existingUser != nil && existingUser.ID != user.ID {
		_ = userRepository.Delete(context.Background(), existingUser.ID)
	}

	err = userRepository.Create(context.Background(), user)
	require.NoError(f.env.T, err, "创建用户失败")

	// 验证可以通过用户名找到刚创建的用户
	userByUsername, err := userRepository.GetByUsername(context.Background(), user.Username)
	require.NoError(f.env.T, err, "通过用户名查找用户失败")
	require.Equal(f.env.T, user.ID, userByUsername.ID, "找到的用户ID不匹配")
	require.True(f.env.T, userByUsername.ValidatePassword(testPassword), "用户密码验证失败")

	f.env.LogSuccess("创建用户: %s (%s)", user.Username, user.ID)

	return user
}

// CreateBook 创建测试书籍（带 e2e_test_ 前缀）
func (f *Fixtures) CreateBook(authorID string, opts ...BookOption) *bookstore.Book {
	bookID := primitive.NewObjectID()
	authorObjID, _ := primitive.ObjectIDFromHex(authorID)

	book := &bookstore.Book{
		ID:           bookID,
		Title:        "e2e_test_book_" + bookID.Hex()[:8],
		AuthorID:     authorObjID,
		Introduction: "E2E测试书籍",
		Categories:   []string{"小说"},
		Price:        0,
		Status:       bookstore.BookStatusPublished,
		WordCount:    10000,
		IsFree:       true,
	}

	// 应用选项
	for _, opt := range opts {
		opt(book)
	}

	bookRepo := bookstoreRepo.NewMongoBookRepository(global.DB.Client(), global.DB.Name())
	err := bookRepo.Create(context.Background(), book)
	require.NoError(f.env.T, err, "创建书籍失败")

	f.env.LogSuccess("创建书籍: %s (%s)", book.Title, book.ID.Hex())

	return book
}

// CreateChapter 创建测试章节（带 e2e_test_ 前缀）
func (f *Fixtures) CreateChapter(bookID string, opts ...ChapterOption) *bookstore.Chapter {
	chapterID := primitive.NewObjectID()
	bookObjID, _ := primitive.ObjectIDFromHex(bookID)

	chapter := &bookstore.Chapter{
		ID:         chapterID,
		BookID:     bookObjID,
		Title:      "e2e_test_chapter_" + chapterID.Hex()[:8],
		ChapterNum: 1, // 使用 ChapterNum 而不是 ChapterNo
		WordCount:  50,
		IsFree:     true,
		Price:      0,
		// Chapter 没有 Content 字段，内容需要通过 ChapterContent 单独存储
	}

	// 应用选项
	for _, opt := range opts {
		opt(chapter)
	}

	chapterRepo := bookstoreRepo.NewMongoChapterRepository(global.DB.Client(), global.DB.Name())
	err := chapterRepo.Create(context.Background(), chapter)
	require.NoError(f.env.T, err, "创建章节失败")

	// 创建章节内容
	chapterContentRepo := bookstoreRepo.NewMongoChapterContentRepository(global.DB)
	chapterContent := &bookstore.ChapterContent{
		ID:        primitive.NewObjectID(),
		ChapterID: chapterID,
		Content:   "这是 E2E 测试章节内容。这是一段测试文字，用于验证阅读功能。内容需要有足够的长度来模拟真实的章节。",
		Format:    "markdown",
		Version:   1,
		WordCount: 50,
	}
	err = chapterContentRepo.Create(context.Background(), chapterContent)
	require.NoError(f.env.T, err, "创建章节内容失败")

	f.env.LogSuccess("创建章节: %s (%s)", chapter.Title, chapter.ID.Hex())

	return chapter
}

// CreateAdminUser 创建管理员用户
func (f *Fixtures) CreateAdminUser(opts ...UserOption) *users.User {
	userID := primitive.NewObjectID()
	testPassword := "Test1234"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
	require.NoError(f.env.T, err, "密码哈希失败")

	user := &users.User{
		ID:       userID.Hex(), // User.ID 是 string 类型
		Username: "e2e_test_admin_" + userID.Hex()[:8],
		Email:    "e2e_test_admin_" + userID.Hex()[:8] + "@example.com",
		Password: string(hashedPassword),
		VIPLevel: 0,
		Status:   users.UserStatusActive,
		Roles:    []string{"admin"},
	}

	// 应用选项
	for _, opt := range opts {
		opt(user)
	}

	userRepository := userRepo.NewMongoUserRepository(global.DB)
	err = userRepository.Create(context.Background(), user)
	require.NoError(f.env.T, err, "创建管理员用户失败")

	f.env.LogSuccess("创建管理员: %s (%s)", user.Username, user.ID)

	return user
}

// GetUserByUsername 根据用户名获取用户
func (f *Fixtures) GetUserByUsername(username string) *users.User {
	userRepository := userRepo.NewMongoUserRepository(global.DB)
	user, err := userRepository.GetByUsername(context.Background(), username)
	require.NoError(f.env.T, err, "获取用户失败")
	return user
}

// CreateChapters 批量创建章节
func (f *Fixtures) CreateChapters(bookID string, count int) []*bookstore.Chapter {
	chapters := make([]*bookstore.Chapter, count)
	chapterRepo := bookstoreRepo.NewMongoChapterRepository(global.DB.Client(), global.DB.Name())
	chapterContentRepo := bookstoreRepo.NewMongoChapterContentRepository(global.DB)
	bookObjID, _ := primitive.ObjectIDFromHex(bookID)

	for i := 0; i < count; i++ {
		chapterID := primitive.NewObjectID()
		chapter := &bookstore.Chapter{
			ID:         chapterID,
			BookID:     bookObjID,
			Title:      "e2e_test_chapter_" + chapterID.Hex()[:8],
			ChapterNum: i + 1,
			WordCount:  20,
			IsFree:     i == 0, // 第一章免费
			Price:      0,
		}

		err := chapterRepo.Create(context.Background(), chapter)
		require.NoError(f.env.T, err, "创建章节失败")

		// 创建章节内容
		chapterContent := &bookstore.ChapterContent{
			ID:        primitive.NewObjectID(),
			ChapterID: chapterID,
			Content:   "这是章节内容。",
			Format:    "markdown",
			Version:   1,
			WordCount: 20,
		}
		err = chapterContentRepo.Create(context.Background(), chapterContent)
		require.NoError(f.env.T, err, "创建章节内容失败")

		chapters[i] = chapter
	}

	f.env.LogSuccess("批量创建 %d 个章节", count)

	return chapters
}
