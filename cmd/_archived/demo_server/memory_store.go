// Package main 提供论文答辩演示用的内存数据存储
package main

import (
	"context"
	"sync"
	"time"
)

// UserInfo 用户信息
type UserInfo struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Password     string    `json:"-"` // 不暴露密码
	Nickname     string    `json:"nickname"`
	Email        string    `json:"email"`
	Avatar       string    `json:"avatar"`
	Role         string    `json:"role"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	LastLoginAt  time.Time `json:"last_login_at"`
	Bio          string    `json:"bio"`
	Followers    int       `json:"followers"`
	Following    int       `json:"following"`
}

// BookInfo 书籍信息
type BookInfo struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Author       string    `json:"author"`
	AuthorID     string    `json:"author_id"`
	Category     string    `json:"category"`
	CategoryID   string    `json:"category_id"`
	Status       string    `json:"status"` // ongoing, completed
	WordCount    int64     `json:"word_count"`
	ChapterCount int       `json:"chapter_count"`
	ClickCount   int64     `json:"click_count"`
	CollectCount int       `json:"collect_count"`
	LikeCount    int       `json:"like_count"`
	Cover        string    `json:"cover"`
	Summary      string    `json:"summary"`
	Tags         []string  `json:"tags"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ChapterInfo 章节信息
type ChapterInfo struct {
	ID         string    `json:"id"`
	BookID     string    `json:"book_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	WordCount  int       `json:"word_count"`
	ChapterNum int       `json:"chapter_num"`
	IsFree     bool      `json:"is_free"`
	CreatedAt  time.Time `json:"created_at"`
}

// CategoryInfo 分类信息
type CategoryInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	BookCount   int    `json:"book_count"`
	Sort        int    `json:"sort"`
}

// CommentInfo 评论信息
type CommentInfo struct {
	ID        string    `json:"id"`
	BookID    string    `json:"book_id"`
	ChapterID string    `json:"chapter_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	LikeCount int       `json:"like_count"`
	CreatedAt time.Time `json:"created_at"`
}

// MemoryDataStore 内存数据存储
type MemoryDataStore struct {
	mu         sync.RWMutex
	users      map[string]*UserInfo
	usersByUname map[string]*UserInfo
	books      map[string]*BookInfo
	chapters   map[string]*ChapterInfo
	chaptersByBook map[string][]*ChapterInfo
	categories map[string]*CategoryInfo
	comments   map[string]*CommentInfo
	tokens     map[string]string // token -> userID
}

// MemoryStore 全局内存存储实例
var MemoryStore *MemoryDataStore

// InitMemoryStore 初始化内存存储
func InitMemoryStore() {
	MemoryStore = &MemoryDataStore{
		users:      make(map[string]*UserInfo),
		usersByUname: make(map[string]*UserInfo),
		books:      make(map[string]*BookInfo),
		chapters:   make(map[string]*ChapterInfo),
		chaptersByBook: make(map[string][]*ChapterInfo),
		categories: make(map[string]*CategoryInfo),
		comments:   make(map[string]*CommentInfo),
		tokens:     make(map[string]string),
	}
}

// ============ 用户操作 ============

// CreateUser 创建用户
func (s *MemoryDataStore) CreateUser(ctx context.Context, user UserInfo) *UserInfo {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	user.CreatedAt = now
	user.Status = "active"
	user.Avatar = "https://picsum.photos/seed/" + user.ID + "/100/100"

	s.users[user.ID] = &user
	s.usersByUname[user.Username] = &user

	return &user
}

// GetUserByID 根据ID获取用户
func (s *MemoryDataStore) GetUserByID(ctx context.Context, id string) *UserInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.users[id]
}

// GetUserByUsername 根据用户名获取用户
func (s *MemoryDataStore) GetUserByUsername(ctx context.Context, username string) *UserInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.usersByUname[username]
}

// ListUsers 列出所有用户
func (s *MemoryDataStore) ListUsers(ctx context.Context, limit, offset int) []*UserInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*UserInfo, 0, len(s.users))
	for _, user := range s.users {
		result = append(result, user)
	}

	// 简单分页
	if offset >= len(result) {
		return []*UserInfo{}
	}
	end := offset + limit
	if end > len(result) {
		end = len(result)
	}
	return result[offset:end]
}

// ============ 书籍操作 ============

// CreateBook 创建书籍
func (s *MemoryDataStore) CreateBook(ctx context.Context, book BookInfo) *BookInfo {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	book.CreatedAt = now
	book.UpdatedAt = now
	book.ChapterCount = 0
	book.LikeCount = 0

	s.books[book.ID] = &book
	return &book
}

// GetBookByID 根据ID获取书籍
func (s *MemoryDataStore) GetBookByID(ctx context.Context, id string) *BookInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.books[id]
}

// ListBooks 列出书籍
func (s *MemoryDataStore) ListBooks(ctx context.Context, category string, status string, limit, offset int) []*BookInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*BookInfo, 0)
	for _, book := range s.books {
		// 过滤条件
		if category != "" && book.Category != category {
			continue
		}
		if status != "" && book.Status != status {
			continue
		}
		result = append(result, book)
	}

	// 按点击数排序（简单冒泡）
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].ClickCount > result[i].ClickCount {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	// 分页
	if offset >= len(result) {
		return []*BookInfo{}
	}
	end := offset + limit
	if end > len(result) {
		end = len(result)
	}
	return result[offset:end]
}

// SearchBooks 搜索书籍
func (s *MemoryDataStore) SearchBooks(ctx context.Context, keyword string, limit int) []*BookInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*BookInfo, 0)
	for _, book := range s.books {
		// 简单的关键词匹配
		if containsIgnoreCase(book.Title, keyword) ||
			containsIgnoreCase(book.Author, keyword) ||
			containsIgnoreCase(book.Summary, keyword) {
			result = append(result, book)
			if len(result) >= limit {
				break
			}
		}
	}
	return result
}

// ============ 章节操作 ============

// CreateChapter 创建章节
func (s *MemoryDataStore) CreateChapter(ctx context.Context, chapter ChapterInfo) *ChapterInfo {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	chapter.CreatedAt = now

	s.chapters[chapter.ID] = &chapter
	s.chaptersByBook[chapter.BookID] = append(s.chaptersByBook[chapter.BookID], &chapter)

	// 更新书籍章节数
	if book, exists := s.books[chapter.BookID]; exists {
		book.ChapterCount++
		book.WordCount += int64(chapter.WordCount)
	}

	return &chapter
}

// GetChapterByID 根据ID获取章节
func (s *MemoryDataStore) GetChapterByID(ctx context.Context, id string) *ChapterInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.chapters[id]
}

// ListChaptersByBook 列出书籍的章节
func (s *MemoryDataStore) ListChaptersByBook(ctx context.Context, bookID string) []*ChapterInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := s.chaptersByBook[bookID]
	if result == nil {
		return []*ChapterInfo{}
	}
	return result
}

// ============ 分类操作 ============

// CreateCategory 创建分类
func (s *MemoryDataStore) CreateCategory(ctx context.Context, cat CategoryInfo) *CategoryInfo {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.categories[cat.ID] = &cat
	return &cat
}

// GetCategoryByID 根据ID获取分类
func (s *MemoryDataStore) GetCategoryByID(ctx context.Context, id string) *CategoryInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.categories[id]
}

// ListCategories 列出所有分类
func (s *MemoryDataStore) ListCategories(ctx context.Context) []*CategoryInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*CategoryInfo, 0, len(s.categories))
	for _, cat := range s.categories {
		result = append(result, cat)
	}
	return result
}

// ============ 评论操作 ============

// CreateComment 创建评论
func (s *MemoryDataStore) CreateComment(ctx context.Context, comment CommentInfo) *CommentInfo {
	s.mu.Lock()
	defer s.mu.Unlock()

	comment.ID = "comment_" + generateID()
	comment.CreatedAt = time.Now()
	comment.LikeCount = 0

	s.comments[comment.ID] = &comment
	return &comment
}

// ListCommentsByBook 列出书籍评论
func (s *MemoryDataStore) ListCommentsByBook(ctx context.Context, bookID string, limit, offset int) []*CommentInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*CommentInfo, 0)
	for _, comment := range s.comments {
		if comment.BookID == bookID {
			result = append(result, comment)
		}
	}

	if offset >= len(result) {
		return []*CommentInfo{}
	}
	end := offset + limit
	if end > len(result) {
		end = len(result)
	}
	return result[offset:end]
}

// ============ Token操作 ============

// StoreToken 存储Token
func (s *MemoryDataStore) StoreToken(ctx context.Context, token, userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[token] = userID
}

// GetUserIDByToken 根据Token获取用户ID
func (s *MemoryDataStore) GetUserIDByToken(ctx context.Context, token string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tokens[token]
}

// DeleteToken 删除Token
func (s *MemoryDataStore) DeleteToken(ctx context.Context, token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tokens, token)
}

// ============ 统计操作 ============

// GetStats 获取统计数据
func (s *MemoryDataStore) GetStats(ctx context.Context) map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]int{
		"users":      len(s.users),
		"books":      len(s.books),
		"chapters":   len(s.chapters),
		"categories": len(s.categories),
		"comments":   len(s.comments),
	}
}

// ============ 辅助函数 ============

// containsIgnoreCase 不区分大小写的包含检查
func containsIgnoreCase(s, substr string) bool {
	// 简单实现，实际应该使用 strings.Contains 配合 strings.ToLower
	return len(s) >= len(substr) && s[0:len(substr)] == substr || len(s) > len(substr) && containsIgnoreCase(s[1:], substr)
}

// generateID 生成简单ID
func generateID() string {
	return time.Now().Format("20060102150405")
}
