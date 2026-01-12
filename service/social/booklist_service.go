package social

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/models/social"
	socialRepo "Qingyu_backend/repository/interfaces/social"
	"Qingyu_backend/service/base"
)

// BookListService 书单服务
type BookListService struct {
	bookListRepo socialRepo.BookListRepository
	eventBus     base.EventBus
	serviceName  string
	version      string
}

// NewBookListService 创建书单服务实例
func NewBookListService(
	bookListRepo socialRepo.BookListRepository,
	eventBus base.EventBus,
) *BookListService {
	return &BookListService{
		bookListRepo: bookListRepo,
		eventBus:     eventBus,
		serviceName:  "BookListService",
		version:      "1.0.0",
	}
}

// BaseService 接口实现
func (s *BookListService) Initialize(ctx context.Context) error { return nil }
func (s *BookListService) Health(ctx context.Context) error {
	if err := s.bookListRepo.Health(ctx); err != nil {
		return fmt.Errorf("书单Repository健康检查失败: %w", err)
	}
	return nil
}
func (s *BookListService) Close(ctx context.Context) error { return nil }
func (s *BookListService) GetServiceName() string          { return s.serviceName }
func (s *BookListService) GetVersion() string              { return s.version }

// CreateBookList 创建书单
func (s *BookListService) CreateBookList(ctx context.Context, userID, userName, userAvatar, title, description, cover, category string, tags []string, isPublic bool) (*social.BookList, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}
	if title == "" {
		return nil, fmt.Errorf("标题不能为空")
	}
	if len(title) > 100 {
		return nil, fmt.Errorf("标题最多100字")
	}
	if len(description) > 500 {
		return nil, fmt.Errorf("描述最多500字")
	}
	if len(tags) > 10 {
		return nil, fmt.Errorf("最多10个标签")
	}

	bookList := &social.BookList{
		UserID:      userID,
		UserName:    userName,
		UserAvatar:  userAvatar,
		Title:       title,
		Description: description,
		Cover:       cover,
		Books:       []social.BookListItem{},
		BookCount:   0,
		LikeCount:   0,
		ForkCount:   0,
		ViewCount:   0,
		IsPublic:    isPublic,
		Tags:        tags,
		Category:    category,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.bookListRepo.CreateBookList(ctx, bookList); err != nil {
		return nil, fmt.Errorf("创建书单失败: %w", err)
	}

	s.publishBookListEvent(ctx, "booklist.created", userID, bookList.ID.Hex())
	return bookList, nil
}

// GetBookLists 获取书单列表
func (s *BookListService) GetBookLists(ctx context.Context, page, size int) ([]*social.BookList, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	bookLists, total, err := s.bookListRepo.GetPublicBookLists(ctx, page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("获取书单列表失败: %w", err)
	}

	return bookLists, total, nil
}

// GetBookListByID 获取书单详情
func (s *BookListService) GetBookListByID(ctx context.Context, bookListID string) (*social.BookList, error) {
	if bookListID == "" {
		return nil, fmt.Errorf("书单ID不能为空")
	}

	bookList, err := s.bookListRepo.GetBookListByID(ctx, bookListID)
	if err != nil {
		return nil, fmt.Errorf("获取书单详情失败: %w", err)
	}

	// 增加浏览次数
	if err := s.bookListRepo.IncrementViewCount(ctx, bookListID); err != nil {
		fmt.Printf("Warning: Failed to increment view count: %v\n", err)
	}

	return bookList, nil
}

// UpdateBookList 更新书单
func (s *BookListService) UpdateBookList(ctx context.Context, userID, bookListID string, updates map[string]interface{}) error {
	if userID == "" || bookListID == "" {
		return fmt.Errorf("用户ID和书单ID不能为空")
	}

	// 获取书单
	bookList, err := s.bookListRepo.GetBookListByID(ctx, bookListID)
	if err != nil {
		return fmt.Errorf("获取书单失败: %w", err)
	}

	// 权限检查
	if bookList.UserID != userID {
		return fmt.Errorf("无权更新该书单")
	}

	// 更新
	if err := s.bookListRepo.UpdateBookList(ctx, bookListID, updates); err != nil {
		return fmt.Errorf("更新书单失败: %w", err)
	}

	return nil
}

// DeleteBookList 删除书单
func (s *BookListService) DeleteBookList(ctx context.Context, userID, bookListID string) error {
	if userID == "" || bookListID == "" {
		return fmt.Errorf("用户ID和书单ID不能为空")
	}

	// 获取书单
	bookList, err := s.bookListRepo.GetBookListByID(ctx, bookListID)
	if err != nil {
		return fmt.Errorf("获取书单失败: %w", err)
	}

	// 权限检查
	if bookList.UserID != userID {
		return fmt.Errorf("无权删除该书单")
	}

	// 删除
	if err := s.bookListRepo.DeleteBookList(ctx, bookListID); err != nil {
		return fmt.Errorf("删除书单失败: %w", err)
	}

	s.publishBookListEvent(ctx, "booklist.deleted", userID, bookListID)
	return nil
}

// LikeBookList 点赞书单
func (s *BookListService) LikeBookList(ctx context.Context, userID, bookListID string) error {
	if userID == "" || bookListID == "" {
		return fmt.Errorf("用户ID和书单ID不能为空")
	}

	// 检查是否已点赞
	isLiked, err := s.bookListRepo.IsBookListLiked(ctx, bookListID, userID)
	if err != nil {
		return fmt.Errorf("检查点赞状态失败: %w", err)
	}
	if isLiked {
		return fmt.Errorf("已经点赞过该书单")
	}

	// 创建点赞
	bookListLike := &social.BookListLike{
		BookListID: bookListID,
		UserID:     userID,
		CreatedAt:  time.Now(),
	}

	if err := s.bookListRepo.CreateBookListLike(ctx, bookListLike); err != nil {
		return fmt.Errorf("点赞失败: %w", err)
	}

	// 增加点赞数
	if err := s.bookListRepo.IncrementBookListLikeCount(ctx, bookListID); err != nil {
		fmt.Printf("Warning: Failed to increment like count: %v\n", err)
	}

	return nil
}

// ForkBookList 复制书单
func (s *BookListService) ForkBookList(ctx context.Context, userID, bookListID string) (*social.BookList, error) {
	if userID == "" || bookListID == "" {
		return nil, fmt.Errorf("用户ID和书单ID不能为空")
	}

	// 验证原书单存在
	_, err := s.bookListRepo.GetBookListByID(ctx, bookListID)
	if err != nil {
		return nil, fmt.Errorf("获取原书单失败: %w", err)
	}

	// 复制书单
	forkedList, err := s.bookListRepo.ForkBookList(ctx, bookListID, userID)
	if err != nil {
		return nil, fmt.Errorf("复制书单失败: %w", err)
	}

	// 增加被复制次数
	if err := s.bookListRepo.IncrementForkCount(ctx, bookListID); err != nil {
		fmt.Printf("Warning: Failed to increment fork count: %v\n", err)
	}

	s.publishBookListEvent(ctx, "booklist.forked", userID, forkedList.ID.Hex())
	return forkedList, nil
}

// GetBooksInList 获取书单中的书籍
func (s *BookListService) GetBooksInList(ctx context.Context, bookListID string) ([]*social.BookListItem, error) {
	if bookListID == "" {
		return nil, fmt.Errorf("书单ID不能为空")
	}

	books, err := s.bookListRepo.GetBooksInList(ctx, bookListID)
	if err != nil {
		return nil, fmt.Errorf("获取书单书籍失败: %w", err)
	}

	return books, nil
}

func (s *BookListService) publishBookListEvent(ctx context.Context, eventType, userID, bookListID string) {
	if s.eventBus == nil {
		return
	}

	event := &base.BaseEvent{
		EventType: eventType,
		EventData: map[string]interface{}{
			"user_id":     userID,
			"booklist_id": bookListID,
		},
		Timestamp: time.Now(),
		Source:    s.serviceName,
	}

	s.eventBus.PublishAsync(ctx, event)
}
