package mocks

import (
	"Qingyu_backend/models/bookstore"
	"Qingyu_backend/models/reader"
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockWalletService is a mock implementation of the wallet service
type MockWalletService struct {
	mock.Mock
}

func (m *MockWalletService) GetBalance(ctx context.Context, userID string) (float64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockWalletService) Consume(ctx context.Context, userID string, amount float64, description string) (string, error) {
	args := m.Called(ctx, userID, amount, description)
	return args.String(0), args.Error(1)
}

func (m *MockWalletService) Refund(ctx context.Context, userID string, amount float64, description string) (string, error) {
	args := m.Called(ctx, userID, amount, description)
	return args.String(0), args.Error(1)
}

func (m *MockWalletService) Recharge(ctx context.Context, userID string, amount float64, orderID, paymentMethod string) (string, error) {
	args := m.Called(ctx, userID, amount, orderID, paymentMethod)
	return args.String(0), args.Error(1)
}

// MockBookStoreRepository is a mock implementation of the book store repository
type MockBookStoreRepository struct {
	mock.Mock
}

func (m *MockBookStoreRepository) GetByID(ctx context.Context, id primitive.ObjectID) (interface{}, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}

func (m *MockBookStoreRepository) Create(ctx context.Context, book interface{}) error {
	args := m.Called(ctx, book)
	return args.Error(0)
}

func (m *MockBookStoreRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockBookStoreRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBookStoreRepository) List(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]interface{}, error) {
	args := m.Called(ctx, filter, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]interface{}), args.Error(1)
}

// MockChapterRepository is a mock implementation of the chapter repository
type MockChapterRepository struct {
	mock.Mock
}

func (m *MockChapterRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstore.Chapter, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetByBookID(ctx context.Context, bookID primitive.ObjectID, limit, offset int) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetFreeChapters(ctx context.Context, bookID primitive.ObjectID, limit, offset int) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetPaidChapters(ctx context.Context, bookID primitive.ObjectID, limit, offset int) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) CountByBookID(ctx context.Context, bookID primitive.ObjectID) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterRepository) GetTotalWordCount(ctx context.Context, bookID primitive.ObjectID) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

// MockChapterPurchaseRepository is a mock implementation of the chapter purchase repository
type MockChapterPurchaseRepository struct {
	mock.Mock
}

func (m *MockChapterPurchaseRepository) GetByUserAndChapter(ctx context.Context, userID, chapterID primitive.ObjectID) (*bookstore.ChapterPurchase, error) {
	args := m.Called(ctx, userID, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.ChapterPurchase), args.Error(1)
}

func (m *MockChapterPurchaseRepository) GetBookPurchaseByUserAndBook(ctx context.Context, userID, bookID primitive.ObjectID) (*bookstore.BookPurchase, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.BookPurchase), args.Error(1)
}

func (m *MockChapterPurchaseRepository) Create(ctx context.Context, purchase *bookstore.ChapterPurchase) error {
	args := m.Called(ctx, purchase)
	return args.Error(0)
}

func (m *MockChapterPurchaseRepository) CreateBatch(ctx context.Context, batch *bookstore.ChapterPurchaseBatch) error {
	args := m.Called(ctx, batch)
	return args.Error(0)
}

func (m *MockChapterPurchaseRepository) CreateBookPurchase(ctx context.Context, purchase *bookstore.BookPurchase) error {
	args := m.Called(ctx, purchase)
	return args.Error(0)
}

func (m *MockChapterPurchaseRepository) GetByUser(ctx context.Context, userID primitive.ObjectID, page, pageSize int) ([]*bookstore.ChapterPurchase, int64, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*bookstore.ChapterPurchase), args.Get(1).(int64), args.Error(2)
}

func (m *MockChapterPurchaseRepository) GetByUserAndBook(ctx context.Context, userID, bookID primitive.ObjectID, page, pageSize int) ([]*bookstore.ChapterPurchase, int64, error) {
	args := m.Called(ctx, userID, bookID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*bookstore.ChapterPurchase), args.Get(1).(int64), args.Error(2)
}

func (m *MockChapterPurchaseRepository) CheckUserPurchasedChapter(ctx context.Context, userID, chapterID primitive.ObjectID) (bool, error) {
	args := m.Called(ctx, userID, chapterID)
	return args.Bool(0), args.Error(1)
}

func (m *MockChapterPurchaseRepository) CheckUserPurchasedBook(ctx context.Context, userID, bookID primitive.ObjectID) (bool, error) {
	args := m.Called(ctx, userID, bookID)
	return args.Bool(0), args.Error(1)
}

func (m *MockChapterPurchaseRepository) GetPurchasedChapterIDs(ctx context.Context, userID, bookID primitive.ObjectID) ([]primitive.ObjectID, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]primitive.ObjectID), args.Error(1)
}

func (m *MockChapterPurchaseRepository) GetTotalSpentByUser(ctx context.Context, userID primitive.ObjectID) (float64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockChapterPurchaseRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

func (m *MockChapterPurchaseRepository) GetBatchesByUser(ctx context.Context, userID primitive.ObjectID, page, pageSize int) ([]*bookstore.ChapterPurchaseBatch, int64, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*bookstore.ChapterPurchaseBatch), args.Get(1).(int64), args.Error(2)
}

func (m *MockChapterPurchaseRepository) GetBookPurchasesByUser(ctx context.Context, userID primitive.ObjectID, page, pageSize int) ([]*bookstore.BookPurchase, int64, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*bookstore.BookPurchase), args.Get(1).(int64), args.Error(2)
}

// MockThemeService is a mock implementation of the theme service
type MockThemeService struct {
	mock.Mock
}

func (m *MockThemeService) GetBuiltInThemes() []*reader.ReaderTheme {
	args := m.Called()
	if args.Get(0) == nil {
		return []*reader.ReaderTheme{}
	}
	return args.Get(0).([]*reader.ReaderTheme)
}

func (m *MockThemeService) GetThemeByName(name string) (*reader.ReaderTheme, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.ReaderTheme), args.Error(1)
}

func (m *MockThemeService) CreateCustomTheme(userID string, theme *reader.ReaderTheme) error {
	args := m.Called(userID, theme)
	return args.Error(0)
}

func (m *MockThemeService) UpdateTheme(userID, themeID string, updates map[string]interface{}) error {
	args := m.Called(userID, themeID, updates)
	return args.Error(0)
}

func (m *MockThemeService) DeleteTheme(userID, themeID string) error {
	args := m.Called(userID, themeID)
	return args.Error(0)
}

func (m *MockThemeService) ActivateTheme(userID, themeName string) error {
	args := m.Called(userID, themeName)
	return args.Error(0)
}

// MockCommentService is a mock implementation of the comment service
type MockCommentService struct {
	mock.Mock
}

func (m *MockCommentService) GetChapterComments(ctx context.Context, chapterID string, filter *reader.ChapterCommentFilter) (*reader.ChapterCommentListResponse, error) {
	args := m.Called(ctx, chapterID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.ChapterCommentListResponse), args.Error(1)
}

func (m *MockCommentService) CreateComment(ctx context.Context, comment *reader.ChapterComment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

func (m *MockCommentService) GetCommentByID(ctx context.Context, commentID string) (*reader.ChapterComment, error) {
	args := m.Called(ctx, commentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.ChapterComment), args.Error(1)
}

func (m *MockCommentService) UpdateComment(ctx context.Context, commentID string, updates map[string]interface{}) error {
	args := m.Called(ctx, commentID, updates)
	return args.Error(0)
}

func (m *MockCommentService) DeleteComment(ctx context.Context, commentID string) error {
	args := m.Called(ctx, commentID)
	return args.Error(0)
}

func (m *MockCommentService) LikeComment(ctx context.Context, commentID, userID string) error {
	args := m.Called(ctx, commentID, userID)
	return args.Error(0)
}

func (m *MockCommentService) UnlikeComment(ctx context.Context, commentID, userID string) error {
	args := m.Called(ctx, commentID, userID)
	return args.Error(0)
}

func (m *MockCommentService) GetParagraphComments(ctx context.Context, chapterID string, paragraphIndex int) (*reader.ParagraphCommentResponse, error) {
	args := m.Called(ctx, chapterID, paragraphIndex)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.ParagraphCommentResponse), args.Error(1)
}

// MockCacheService is a mock implementation of the cache service
type MockCacheService struct {
	mock.Mock
}

func (m *MockCacheService) Get(ctx context.Context, key string) (interface{}, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}

func (m *MockCacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockCacheService) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateChapterCache(ctx context.Context, chapterID string) {
	m.Called(ctx, chapterID)
}

func (m *MockCacheService) InvalidateBookChaptersCache(ctx context.Context, bookID string) {
	m.Called(ctx, bookID)
}

func (m *MockCacheService) InvalidateBookDetailCache(ctx context.Context, bookID string) {
	m.Called(ctx, bookID)
}

func (m *MockCacheService) InvalidateUserCache(ctx context.Context, userID string) {
	m.Called(ctx, userID)
}

// MockUserService is a mock implementation of the user service
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUserByID(ctx context.Context, userID string) (interface{}, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}

func (m *MockUserService) IsVIPUser(ctx context.Context, userID string) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserService) UpdateUserReadingSettings(ctx context.Context, userID string, settings map[string]interface{}) error {
	args := m.Called(ctx, userID, settings)
	return args.Error(0)
}

// MockNotificationService is a mock implementation of the notification service
type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) SendNotification(ctx context.Context, userID string, notification interface{}) error {
	args := m.Called(ctx, userID, notification)
	return args.Error(0)
}

func (m *MockNotificationService) SendBatchNotifications(ctx context.Context, userIDs []string, notification interface{}) error {
	args := m.Called(ctx, userIDs, notification)
	return args.Error(0)
}

// MockAnalyticsService is a mock implementation of the analytics service
type MockAnalyticsService struct {
	mock.Mock
}

func (m *MockAnalyticsService) TrackEvent(ctx context.Context, event string, properties map[string]interface{}) error {
	args := m.Called(ctx, event, properties)
	return args.Error(0)
}

func (m *MockAnalyticsService) TrackPurchase(ctx context.Context, userID string, purchaseDetails map[string]interface{}) error {
	args := m.Called(ctx, userID, purchaseDetails)
	return args.Error(0)
}

func (m *MockAnalyticsService) TrackReadingProgress(ctx context.Context, userID, chapterID string, progress float64) error {
	args := m.Called(ctx, userID, chapterID, progress)
	return args.Error(0)
}

// Test Helper Functions

// CreateTestChapter creates a test chapter for testing purposes
func CreateTestChapter(id, bookID primitive.ObjectID, num int, isFree bool, price float64) *bookstore.Chapter {
	return &bookstore.Chapter{
		ID:          id,
		BookID:      bookID,
		Title:       "Test Chapter " + string(rune(num)),
		ChapterNum:  num,
		WordCount:   2000,
		IsFree:      isFree,
		Price:       price,
		PublishTime: time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// CreateTestBook creates a test book map for testing purposes
func CreateTestBook(id primitive.ObjectID, title string) map[string]interface{} {
	return map[string]interface{}{
		"_id":       id,
		"title":     title,
		"cover_url": "http://example.com/cover.jpg",
		"author":    "Test Author",
		"status":    "published",
	}
}

// CreateTestPurchase creates a test chapter purchase for testing purposes
func CreateTestPurchase(id, userID, chapterID, bookID primitive.ObjectID, price float64) *bookstore.ChapterPurchase {
	return &bookstore.ChapterPurchase{
		ID:           id,
		UserID:       userID,
		ChapterID:    chapterID,
		BookID:       bookID,
		Price:        price,
		PurchaseTime: time.Now(),
		CreatedAt:    time.Now(),
	}
}

// CreateTestComment creates a test chapter comment for testing purposes
func CreateTestComment(id, chapterID, bookID, userID string, content string, rating int) *reader.ChapterComment {
	return &reader.ChapterComment{
		ID:        id,
		ChapterID: chapterID,
		BookID:    bookID,
		UserID:    userID,
		Content:   content,
		Rating:    rating,
		LikeCount: 0,
		IsVisible: true,
		IsDeleted: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
