package data

import (
	"testing"

	"Qingyu_backend/models/bookstore"
	"Qingyu_backend/models/social"
	"Qingyu_backend/models/users"
)

// ScenarioBuilder 场景数据构建器
type ScenarioBuilder struct {
	factory *TestDataFactory
	t       *testing.T
}

// NewScenarioBuilder 创建场景构建器
func NewScenarioBuilder(t *testing.T) *ScenarioBuilder {
	return &ScenarioBuilder{
		factory: NewTestDataFactory(t),
		t:       t,
	}
}

// ReaderWithProgress 构建有阅读进度的读者场景
type ReaderWithProgress struct {
	User           *users.User
	Books          []*bookstore.Book
	CurrentBook    *bookstore.Book
	CurrentChapter *bookstore.Chapter
	Progress       float64
}

// BuildReaderWithProgress 创建有阅读进度的读者场景
func (sb *ScenarioBuilder) BuildReaderWithProgress() *ReaderWithProgress {
	// 创建读者
	user := sb.factory.CreateUser(UserOptions{
		Username: "e2e_reader_with_progress",
		VIPLevel: 0,
	})

	// 创建作者和书籍
	author := sb.factory.CreateUser(UserOptions{
		Username: "e2e_author_for_reader",
		Roles:    []string{"reader", "author"},
	})

	book := sb.factory.CreateBook(BookOptions{
		Title:        "读者正在阅读的书",
		AuthorID:     author.ID,
		Price:        0,
		IsFree:       true,
		ChapterCount: 5,
	})

	// 创建章节
	chapters := make([]*bookstore.Chapter, 5)
	for i := 0; i < 5; i++ {
		chapters[i] = sb.factory.CreateChapter(book.ID.Hex(), i+1, i == 0)
	}

	return &ReaderWithProgress{
		User:           user,
		Books:          []*bookstore.Book{book},
		CurrentBook:    book,
		CurrentChapter: chapters[0],
		Progress:       0.3,
	}
}

// SocialInteraction 构建社交互动场景
type SocialInteraction struct {
	Users       []*users.User
	Book        *bookstore.Book
	Comments    []*social.Comment
	Collections []*social.Collection
	Likes       []interface{}
}

// BuildSocialInteraction 创建社交互动场景
func (sb *ScenarioBuilder) BuildSocialInteraction(userCount int) *SocialInteraction {
	// 创建作者
	author := sb.factory.CreateUser(UserOptions{
		Username: "e2e_social_author",
		Roles:    []string{"reader", "author"},
	})

	// 创建书籍
	book := sb.factory.CreateBook(BookOptions{
		Title:        "热门互动书籍",
		AuthorID:     author.ID,
		Price:        0,
		IsFree:       true,
		ChapterCount: 3,
	})

	// 创建互动用户
	users := make([]*users.User, userCount)
	for i := 0; i < userCount; i++ {
		users[i] = sb.factory.CreateUser(UserOptions{
			Username: "",
			VIPLevel: i % 2, // 混合VIP等级
		})
	}

	return &SocialInteraction{
		Users:       users,
		Book:        book,
		Comments:    []*social.Comment{},
		Collections: []*social.Collection{},
		Likes:       []interface{}{},
	}
}

// PaidContent 构建付费内容场景
type PaidContent struct {
	Author      *users.User
	FreeUser    *users.User
	VIPUser     *users.User
	PaidBook    *bookstore.Book
	FreeChapters []*bookstore.Chapter
	PaidChapters []*bookstore.Chapter
}

// BuildPaidContent 创建付费内容场景
func (sb *ScenarioBuilder) BuildPaidContent() *PaidContent {
	// 创建作者
	author := sb.factory.CreateUser(UserOptions{
		Username: "e2e_paid_content_author",
		Roles:    []string{"reader", "author"},
	})

	// 创建免费用户
	freeUser := sb.factory.CreateUser(UserOptions{
		Username: "e2e_free_reader",
		VIPLevel: 0,
		Balance:  0,
	})

	// 创建VIP用户
	vipUser := sb.factory.CreateUser(UserOptions{
		Username: "e2e_vip_reader",
		VIPLevel: 1,
		Balance:  0,
	})

	// 创建付费书籍
	paidBook := sb.factory.CreateBook(BookOptions{
		Title:        "付费书籍",
		AuthorID:     author.ID,
		Price:        100,
		IsFree:       false,
		ChapterCount: 10,
	})

	// 创建章节（前3章免费，后7章付费）
	freeChapters := make([]*bookstore.Chapter, 3)
	paidChapters := make([]*bookstore.Chapter, 7)

	for i := 0; i < 3; i++ {
		freeChapters[i] = sb.factory.CreateChapter(paidBook.ID.Hex(), i+1, true)
	}
	for i := 3; i < 10; i++ {
		paidChapters[i-3] = sb.factory.CreateChapter(paidBook.ID.Hex(), i+1, false)
	}

	return &PaidContent{
		Author:       author,
		FreeUser:     freeUser,
		VIPUser:      vipUser,
		PaidBook:     paidBook,
		FreeChapters: freeChapters,
		PaidChapters: paidChapters,
	}
}
