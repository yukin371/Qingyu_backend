package bookstore

import (
	"Qingyu_backend/models/bookstore"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock implementations

type MockChapterRepository struct {
	mock.Mock
}

func (m *MockChapterRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockChapterRepository) Create(ctx context.Context, entity *bookstore.Chapter) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockChapterRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstore.Chapter, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockChapterRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockChapterRepository) List(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, filter, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetByBookID(ctx context.Context, bookID primitive.ObjectID, limit, offset int) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetByBookIDAndChapterNum(ctx context.Context, bookID primitive.ObjectID, chapterNum int) (*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID, chapterNum)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetByTitle(ctx context.Context, title string, limit, offset int) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, title, limit, offset)
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

func (m *MockChapterRepository) GetPublishedChapters(ctx context.Context, bookID primitive.ObjectID, limit, offset int) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetChapterRange(ctx context.Context, bookID primitive.ObjectID, startChapter, endChapter int) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID, startChapter, endChapter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) Search(ctx context.Context, keyword string, limit, offset int) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, keyword, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) SearchByFilter(ctx context.Context, filter interface{}) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) CountByBookID(ctx context.Context, bookID primitive.ObjectID) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterRepository) CountFreeChapters(ctx context.Context, bookID primitive.ObjectID) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterRepository) CountPaidChapters(ctx context.Context, bookID primitive.ObjectID) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterRepository) CountPublishedChapters(ctx context.Context, bookID primitive.ObjectID) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterRepository) GetTotalWordCount(ctx context.Context, bookID primitive.ObjectID) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterRepository) GetPreviousChapter(ctx context.Context, bookID primitive.ObjectID, chapterNum int) (*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID, chapterNum)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetNextChapter(ctx context.Context, bookID primitive.ObjectID, chapterNum int) (*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID, chapterNum)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetFirstChapter(ctx context.Context, bookID primitive.ObjectID) (*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetLastChapter(ctx context.Context, bookID primitive.ObjectID) (*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) BatchUpdatePrice(ctx context.Context, chapterIDs []primitive.ObjectID, price float64) error {
	args := m.Called(ctx, chapterIDs, price)
	return args.Error(0)
}

func (m *MockChapterRepository) BatchDelete(ctx context.Context, chapterIDs []primitive.ObjectID) error {
	args := m.Called(ctx, chapterIDs)
	return args.Error(0)
}

func (m *MockChapterRepository) BatchUpdateFreeStatus(ctx context.Context, chapterIDs []primitive.ObjectID, isFree bool) error {
	args := m.Called(ctx, chapterIDs, isFree)
	return args.Error(0)
}

func (m *MockChapterRepository) BatchUpdatePublishTime(ctx context.Context, chapterIDs []primitive.ObjectID, publishTime time.Time) error {
	args := m.Called(ctx, chapterIDs, publishTime)
	return args.Error(0)
}

func (m *MockChapterRepository) DeleteByBookID(ctx context.Context, bookID primitive.ObjectID) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockChapterRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

type MockChapterPurchaseRepository struct {
	mock.Mock
}

func (m *MockChapterPurchaseRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockChapterPurchaseRepository) Create(ctx context.Context, purchase *bookstore.ChapterPurchase) error {
	args := m.Called(ctx, purchase)
	return args.Error(0)
}

func (m *MockChapterPurchaseRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstore.ChapterPurchase, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.ChapterPurchase), args.Error(1)
}

func (m *MockChapterPurchaseRepository) GetByUserAndChapter(ctx context.Context, userID, chapterID primitive.ObjectID) (*bookstore.ChapterPurchase, error) {
	args := m.Called(ctx, userID, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.ChapterPurchase), args.Error(1)
}

func (m *MockChapterPurchaseRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockChapterPurchaseRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
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

func (m *MockChapterPurchaseRepository) CreateBatch(ctx context.Context, batch *bookstore.ChapterPurchaseBatch) error {
	args := m.Called(ctx, batch)
	return args.Error(0)
}

func (m *MockChapterPurchaseRepository) GetBatchByID(ctx context.Context, id primitive.ObjectID) (*bookstore.ChapterPurchaseBatch, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.ChapterPurchaseBatch), args.Error(1)
}

func (m *MockChapterPurchaseRepository) GetBatchesByUser(ctx context.Context, userID primitive.ObjectID, page, pageSize int) ([]*bookstore.ChapterPurchaseBatch, int64, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*bookstore.ChapterPurchaseBatch), args.Get(1).(int64), args.Error(2)
}

func (m *MockChapterPurchaseRepository) GetBatchesByUserAndBook(ctx context.Context, userID, bookID primitive.ObjectID, page, pageSize int) ([]*bookstore.ChapterPurchaseBatch, int64, error) {
	args := m.Called(ctx, userID, bookID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*bookstore.ChapterPurchaseBatch), args.Get(1).(int64), args.Error(2)
}

func (m *MockChapterPurchaseRepository) CreateBookPurchase(ctx context.Context, purchase *bookstore.BookPurchase) error {
	args := m.Called(ctx, purchase)
	return args.Error(0)
}

func (m *MockChapterPurchaseRepository) GetBookPurchaseByID(ctx context.Context, id primitive.ObjectID) (*bookstore.BookPurchase, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.BookPurchase), args.Error(1)
}

func (m *MockChapterPurchaseRepository) GetBookPurchaseByUserAndBook(ctx context.Context, userID, bookID primitive.ObjectID) (*bookstore.BookPurchase, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.BookPurchase), args.Error(1)
}

func (m *MockChapterPurchaseRepository) GetBookPurchasesByUser(ctx context.Context, userID primitive.ObjectID, page, pageSize int) ([]*bookstore.BookPurchase, int64, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*bookstore.BookPurchase), args.Get(1).(int64), args.Error(2)
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

func (m *MockChapterPurchaseRepository) CountByUser(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterPurchaseRepository) CountByUserAndBook(ctx context.Context, userID, bookID primitive.ObjectID) (int64, error) {
	args := m.Called(ctx, userID, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterPurchaseRepository) GetTotalSpentByUser(ctx context.Context, userID primitive.ObjectID) (float64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockChapterPurchaseRepository) GetTotalSpentByUserAndBook(ctx context.Context, userID, bookID primitive.ObjectID) (float64, error) {
	args := m.Called(ctx, userID, bookID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockChapterPurchaseRepository) GetPurchasesByTimeRange(ctx context.Context, userID primitive.ObjectID, startTime, endTime time.Time) ([]*bookstore.ChapterPurchase, error) {
	args := m.Called(ctx, userID, startTime, endTime)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.ChapterPurchase), args.Error(1)
}

func (m *MockChapterPurchaseRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

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

type MockWalletService struct {
	mock.Mock
}

func (m *MockWalletService) GetBalance(userID string) (float64, error) {
	args := m.Called(userID)
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

type MockCacheService struct {
	mock.Mock
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

// Test Helper Functions

func createTestChapter(id, bookID primitive.ObjectID, num int, isFree bool, price float64) *bookstore.Chapter {
	return &bookstore.Chapter{
		ID:         id,
		BookID:     bookID,
		Title:      "Test Chapter",
		ChapterNum: num,
		WordCount:  2000,
		IsFree:     isFree,
		Price:      price,
		PublishTime: time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func createTestBook(id primitive.ObjectID, title string) map[string]interface{} {
	return map[string]interface{}{
		"_id":       id,
		"title":     title,
		"cover_url": "http://example.com/cover.jpg",
		"author":    "Test Author",
		"status":    "published",
	}
}

// Test: GetChapterCatalog

func TestChapterPurchaseService_GetChapterCatalog_Success(t *testing.T) {
	chapterRepo := new(MockChapterRepository)
	purchaseRepo := new(MockChapterPurchaseRepository)
	bookRepo := new(MockBookStoreRepository)
	walletService := new(MockWalletService)
	cacheService := new(MockCacheService)

	service := NewChapterPurchaseService(chapterRepo, purchaseRepo, bookRepo, walletService, cacheService)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()

	chapters := []*bookstore.Chapter{
		createTestChapter(primitive.NewObjectID(), bookID, 1, true, 0),
		createTestChapter(primitive.NewObjectID(), bookID, 2, false, 0.99),
		createTestChapter(primitive.NewObjectID(), bookID, 3, false, 1.99),
	}

	book := createTestBook(bookID, "Test Book")

	bookRepo.On("GetByID", ctx, bookID).Return(book, nil)
	chapterRepo.On("GetByBookID", ctx, bookID, 10000, 0).Return(chapters, nil)
	purchaseRepo.On("GetPurchasedChapterIDs", ctx, userID, bookID).Return([]primitive.ObjectID{chapters[1].ID}, nil)
	chapterRepo.On("CountByBookID", ctx, bookID).Return(int64(3), nil)
	chapterRepo.On("GetTotalWordCount", ctx, bookID).Return(int64(6000), nil)

	catalog, err := service.GetChapterCatalog(ctx, userID, bookID)

	assert.NoError(t, err)
	assert.NotNil(t, catalog)
	assert.Equal(t, bookID, catalog.BookID)
	assert.Equal(t, "Test Book", catalog.BookTitle)
	assert.Equal(t, 3, catalog.TotalChapters)
	assert.Equal(t, 1, catalog.FreeChapters)
	assert.Equal(t, 2, catalog.PaidChapters)
	assert.Len(t, catalog.Chapters, 3)

	// Check purchased status
	assert.True(t, catalog.Chapters[1].IsPurchased)
	assert.False(t, catalog.Chapters[2].IsPurchased)

	bookRepo.AssertExpectations(t)
	chapterRepo.AssertExpectations(t)
	purchaseRepo.AssertExpectations(t)
}

func TestChapterPurchaseService_GetChapterCatalog_BookNotFound(t *testing.T) {
	chapterRepo := new(MockChapterRepository)
	purchaseRepo := new(MockChapterPurchaseRepository)
	bookRepo := new(MockBookStoreRepository)
	walletService := new(MockWalletService)
	cacheService := new(MockCacheService)

	service := NewChapterPurchaseService(chapterRepo, purchaseRepo, bookRepo, walletService, cacheService)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()

	bookRepo.On("GetByID", ctx, bookID).Return(nil, nil)

	catalog, err := service.GetChapterCatalog(ctx, userID, bookID)

	assert.Error(t, err)
	assert.Nil(t, catalog)
	assert.Contains(t, err.Error(), "book not found")

	bookRepo.AssertExpectations(t)
}

func TestChapterPurchaseService_GetChapterCatalog_EmptyBookID(t *testing.T) {
	chapterRepo := new(MockChapterRepository)
	purchaseRepo := new(MockChapterPurchaseRepository)
	bookRepo := new(MockBookStoreRepository)
	walletService := new(MockWalletService)
	cacheService := new(MockCacheService)

	service := NewChapterPurchaseService(chapterRepo, purchaseRepo, bookRepo, walletService, cacheService)

	ctx := context.Background()
	userID := primitive.NewObjectID()

	catalog, err := service.GetChapterCatalog(ctx, userID, primitive.NilObjectID)

	assert.Error(t, err)
	assert.Nil(t, catalog)
	assert.Contains(t, err.Error(), "book ID cannot be empty")
}

// Test: GetTrialChapters

func TestChapterPurchaseService_GetTrialChapters_Success(t *testing.T) {
	chapterRepo := new(MockChapterRepository)
	purchaseRepo := new(MockChapterPurchaseRepository)
	bookRepo := new(MockBookStoreRepository)
	walletService := new(MockWalletService)
	cacheService := new(MockCacheService)

	service := NewChapterPurchaseService(chapterRepo, purchaseRepo, bookRepo, walletService, cacheService)

	ctx := context.Background()
	bookID := primitive.NewObjectID()

	freeChapters := []*bookstore.Chapter{
		createTestChapter(primitive.NewObjectID(), bookID, 1, true, 0),
		createTestChapter(primitive.NewObjectID(), bookID, 2, true, 0),
	}

	chapterRepo.On("GetFreeChapters", ctx, bookID, 10, 0).Return(freeChapters, nil)

	chapters, err := service.GetTrialChapters(ctx, bookID, 10)

	assert.NoError(t, err)
	assert.NotNil(t, chapters)
	assert.Len(t, chapters, 2)

	chapterRepo.AssertExpectations(t)
}

func TestChapterPurchaseService_GetTrialChapters_DefaultCount(t *testing.T) {
	chapterRepo := new(MockChapterRepository)
	purchaseRepo := new(MockChapterPurchaseRepository)
	bookRepo := new(MockBookStoreRepository)
	walletService := new(MockWalletService)
	cacheService := new(MockCacheService)

	service := NewChapterPurchaseService(chapterRepo, purchaseRepo, bookRepo, walletService, cacheService)

	ctx := context.Background()
	bookID := primitive.NewObjectID()

	freeChapters := []*bookstore.Chapter{
		createTestChapter(primitive.NewObjectID(), bookID, 1, true, 0),
	}

	chapterRepo.On("GetFreeChapters", ctx, bookID, 10, 0).Return(freeChapters, nil)

	chapters, err := service.GetTrialChapters(ctx, bookID, 0)

	assert.NoError(t, err)
	assert.NotNil(t, chapters)
	assert.Len(t, chapters, 1)

	chapterRepo.AssertExpectations(t)
}

// Test: GetVIPChapters

func TestChapterPurchaseService_GetVIPChapters_Success(t *testing.T) {
	chapterRepo := new(MockChapterRepository)
	purchaseRepo := new(MockChapterPurchaseRepository)
	bookRepo := new(MockBookStoreRepository)
	walletService := new(MockWalletService)
	cacheService := new(MockCacheService)

	service := NewChapterPurchaseService(chapterRepo, purchaseRepo, bookRepo, walletService, cacheService)

	ctx := context.Background()
	bookID := primitive.NewObjectID()

	vipChapters := []*bookstore.Chapter{
		createTestChapter(primitive.NewObjectID(), bookID, 5, false, 2.99),
		createTestChapter(primitive.NewObjectID(), bookID, 6, false, 3.99),
	}

	chapterRepo.On("GetPaidChapters", ctx, bookID, 1000, 0).Return(vipChapters, nil)

	chapters, err := service.GetVIPChapters(ctx, bookID)

	assert.NoError(t, err)
	assert.NotNil(t, chapters)
	assert.Len(t, chapters, 2)

	chapterRepo.AssertExpectations(t)
}

// Test: PurchaseChapter

func TestChapterPurchaseService_PurchaseChapter_Success(t *testing.T) {
	chapterRepo := new(MockChapterRepository)
	purchaseRepo := new(MockChapterPurchaseRepository)
	bookRepo := new(MockBookStoreRepository)
	walletService := new(MockWalletService)
	cacheService := new(MockCacheService)

	service := NewChapterPurchaseService(chapterRepo, purchaseRepo, bookRepo, walletService, cacheService)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	chapterID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()

	chapter := createTestChapter(chapterID, bookID, 1, false, 1.99)
	book := createTestBook(bookID, "Test Book")

	purchaseRepo.On("GetByUserAndChapter", ctx, userID, chapterID).Return(nil, errors.New("not found"))
	chapterRepo.On("GetByID", ctx, chapterID).Return(chapter, nil)
	bookRepo.On("GetByID", ctx, bookID).Return(book, nil)
	walletService.On("GetBalance", userID.Hex()).Return(10.0, nil)

	purchaseRepo.On("Transaction", ctx, mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			walletService.On("Consume", ctx, userID.Hex(), 1.99, "购买章节: Test Chapter").Return("txn_123", nil)
			purchaseRepo.On("Create", ctx, mock.AnythingOfType("*bookstore.ChapterPurchase")).Return(nil)
			err := fn(ctx)
			assert.NoError(t, err)
		}).Return(nil)

	cacheService.On("InvalidateChapterCache", ctx, chapterID.Hex())
	cacheService.On("InvalidateBookChaptersCache", ctx, bookID.Hex())

	purchase, err := service.PurchaseChapter(ctx, userID, chapterID)

	assert.NoError(t, err)
	assert.NotNil(t, purchase)
	assert.Equal(t, userID, purchase.UserID)
	assert.Equal(t, chapterID, purchase.ChapterID)
	assert.Equal(t, bookID, purchase.BookID)
	assert.Equal(t, 1.99, purchase.Price)

	purchaseRepo.AssertExpectations(t)
	chapterRepo.AssertExpectations(t)
	bookRepo.AssertExpectations(t)
	walletService.AssertExpectations(t)
	cacheService.AssertExpectations(t)
}

func TestChapterPurchaseService_PurchaseChapter_AlreadyPurchased(t *testing.T) {
	chapterRepo := new(MockChapterRepository)
	purchaseRepo := new(MockChapterPurchaseRepository)
	bookRepo := new(MockBookStoreRepository)
	walletService := new(MockWalletService)
	cacheService := new(MockCacheService)

	service := NewChapterPurchaseService(chapterRepo, purchaseRepo, bookRepo, walletService, cacheService)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	chapterID := primitive.NewObjectID()

	existingPurchase := &bookstore.ChapterPurchase{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		ChapterID: chapterID,
	}

	purchaseRepo.On("GetByUserAndChapter", ctx, userID, chapterID).Return(existingPurchase, nil)

	purchase, err := service.PurchaseChapter(ctx, userID, chapterID)

	assert.Error(t, err)
	assert.Nil(t, purchase)
	assert.Contains(t, err.Error(), "already purchased")

	purchaseRepo.AssertExpectations(t)
}

func TestChapterPurchaseService_PurchaseChapter_InsufficientBalance(t *testing.T) {
	chapterRepo := new(MockChapterRepository)
	purchaseRepo := new(MockChapterPurchaseRepository)
	bookRepo := new(MockBookStoreRepository)
	walletService := new(MockWalletService)
	cacheService := new(MockCacheService)

	service := NewChapterPurchaseService(chapterRepo, purchaseRepo, bookRepo, walletService, cacheService)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	chapterID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()

	chapter := createTestChapter(chapterID, bookID, 1, false, 5.99)
	book := createTestBook(bookID, "Test Book")

	purchaseRepo.On("GetByUserAndChapter", ctx, userID, chapterID).Return(nil, errors.New("not found"))
	chapterRepo.On("GetByID", ctx, chapterID).Return(chapter, nil)
	bookRepo.On("GetByID", ctx, bookID).Return(book, nil)
	walletService.On("GetBalance", userID.Hex()).Return(1.0, nil)

	purchase, err := service.PurchaseChapter(ctx, userID, chapterID)

	assert.Error(t, err)
	assert.Nil(t, purchase)
	assert.Contains(t, err.Error(), "insufficient balance")

	purchaseRepo.AssertExpectations(t)
	chapterRepo.AssertExpectations(t)
	bookRepo.AssertExpectations(t)
	walletService.AssertExpectations(t)
}

func TestChapterPurchaseService_PurchaseChapter_FreeChapter(t *testing.T) {
	chapterRepo := new(MockChapterRepository)
	purchaseRepo := new(MockChapterPurchaseRepository)
	bookRepo := new(MockBookStoreRepository)
	walletService := new(MockWalletService)
	cacheService := new(MockCacheService)

	service := NewChapterPurchaseService(chapterRepo, purchaseRepo, bookRepo, walletService, cacheService)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	chapterID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()

	chapter := createTestChapter(chapterID, bookID, 1, true, 0)

	purchaseRepo.On("GetByUserAndChapter", ctx, userID, chapterID).Return(nil, errors.New("not found"))
	chapterRepo.On("GetByID", ctx, chapterID).Return(chapter, nil)

	purchase, err := service.PurchaseChapter(ctx, userID, chapterID)

	assert.Error(t, err)
	assert.Nil(t, purchase)
	assert.Contains(t, err.Error(), "cannot purchase free chapter")

	purchaseRepo.AssertExpectations(t)
	chapterRepo.AssertExpectations(t)
}

// Test: PurchaseChapters (Batch)

func TestChapterPurchaseService_PurchaseChapters_Success(t *testing.T) {
	chapterRepo := new(MockChapterRepository)
	purchaseRepo := new(MockChapterPurchaseRepository)
	bookRepo := new(MockBookStoreRepository)
	walletService := new(MockWalletService)
	cacheService := new(MockCacheService)

	service := NewChapterPurchaseService(chapterRepo, purchaseRepo, bookRepo, walletService, cacheService)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()

	chapterIDs := []primitive.ObjectID{
		primitive.NewObjectID(),
		primitive.NewObjectID(),
	}

	chapters := []*bookstore.Chapter{
		createTestChapter(chapterIDs[0], bookID, 1, false, 1.99),
		createTestChapter(chapterIDs[1], bookID, 2, false, 2.99),
	}

	book := createTestBook(bookID, "Test Book")

	for _, chapterID := range chapterIDs {
		purchaseRepo.On("GetByUserAndChapter", ctx, userID, chapterID).Return(nil, errors.New("not found"))
	}

	for _, chapter := range chapters {
		chapterRepo.On("GetByID", ctx, chapter.ID).Return(chapter, nil)
	}

	bookRepo.On("GetByID", ctx, bookID).Return(book, nil)
	walletService.On("GetBalance", userID.Hex()).Return(10.0, nil)

	purchaseRepo.On("Transaction", ctx, mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			walletService.On("Consume", ctx, userID.Hex(), 4.98, "批量购买章节: 2章").Return("txn_batch_123", nil)
			purchaseRepo.On("CreateBatch", ctx, mock.AnythingOfType("*bookstore.ChapterPurchaseBatch")).Return(nil)
			purchaseRepo.On("Create", ctx, mock.AnythingOfType("*bookstore.ChapterPurchase")).Return(nil).Times(2)
			err := fn(ctx)
			assert.NoError(t, err)
		}).Return(nil)

	cacheService.On("InvalidateBookChaptersCache", ctx, bookID.Hex())

	batch, err := service.PurchaseChapters(ctx, userID, chapterIDs)

	assert.NoError(t, err)
	assert.NotNil(t, batch)
	assert.Equal(t, userID, batch.UserID)
	assert.Equal(t, bookID, batch.BookID)
	assert.Equal(t, 4.98, batch.TotalPrice)
	assert.Equal(t, 2, batch.ChaptersCount)

	purchaseRepo.AssertExpectations(t)
	chapterRepo.AssertExpectations(t)
	bookRepo.AssertExpectations(t)
	walletService.AssertExpectations(t)
	cacheService.AssertExpectations(t)
}

// Test: PurchaseBook

func TestChapterPurchaseService_PurchaseBook_Success(t *testing.T) {
	chapterRepo := new(MockChapterRepository)
	purchaseRepo := new(MockChapterPurchaseRepository)
	bookRepo := new(MockBookStoreRepository)
	walletService := new(MockWalletService)
	cacheService := new(MockCacheService)

	service := NewChapterPurchaseService(chapterRepo, purchaseRepo, bookRepo, walletService, cacheService)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()

	book := createTestBook(bookID, "Test Book")

	paidChapters := []*bookstore.Chapter{
		createTestChapter(primitive.NewObjectID(), bookID, 3, false, 1.99),
		createTestChapter(primitive.NewObjectID(), bookID, 4, false, 2.99),
		createTestChapter(primitive.NewObjectID(), bookID, 5, false, 3.99),
	}

	purchaseRepo.On("GetBookPurchaseByUserAndBook", ctx, userID, bookID).Return(nil, errors.New("not found"))
	bookRepo.On("GetByID", ctx, bookID).Return(book, nil)
	chapterRepo.On("GetPaidChapters", ctx, bookID, 10000, 0).Return(paidChapters, int64(3), nil)
	walletService.On("GetBalance", userID.Hex()).Return(10.0, nil)

	purchaseRepo.On("Transaction", ctx, mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			walletService.On("Consume", ctx, userID.Hex(), 7.136, "购买全书: Test Book").Return("txn_book_123", nil)
			purchaseRepo.On("CreateBookPurchase", ctx, mock.AnythingOfType("*bookstore.BookPurchase")).Return(nil)
			purchaseRepo.On("Create", ctx, mock.AnythingOfType("*bookstore.ChapterPurchase")).Return(nil).Times(3)
			err := fn(ctx)
			assert.NoError(t, err)
		}).Return(nil)

	cacheService.On("InvalidateBookDetailCache", ctx, bookID.Hex())
	cacheService.On("InvalidateBookChaptersCache", ctx, bookID.Hex())

	purchase, err := service.PurchaseBook(ctx, userID, bookID)

	assert.NoError(t, err)
	assert.NotNil(t, purchase)
	assert.Equal(t, userID, purchase.UserID)
	assert.Equal(t, bookID, purchase.BookID)
	assert.Equal(t, 3, purchase.ChapterCount)
	assert.True(t, purchase.TotalPrice < 8.97) // Should have discount

	purchaseRepo.AssertExpectations(t)
	bookRepo.AssertExpectations(t)
	chapterRepo.AssertExpectations(t)
	walletService.AssertExpectations(t)
	cacheService.AssertExpectations(t)
}

// Test: CheckChapterAccess

func TestChapterPurchaseService_CheckChapterAccess_FreeChapter(t *testing.T) {
	chapterRepo := new(MockChapterRepository)
	purchaseRepo := new(MockChapterPurchaseRepository)
	bookRepo := new(MockBookStoreRepository)
	walletService := new(MockWalletService)
	cacheService := new(MockCacheService)

	service := NewChapterPurchaseService(chapterRepo, purchaseRepo, bookRepo, walletService, cacheService)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	chapterID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()

	chapter := createTestChapter(chapterID, bookID, 1, true, 0)

	chapterRepo.On("GetByID", ctx, chapterID).Return(chapter, nil)

	accessInfo, err := service.CheckChapterAccess(ctx, userID, chapterID)

	assert.NoError(t, err)
	assert.NotNil(t, accessInfo)
	assert.True(t, accessInfo.CanAccess)
	assert.True(t, accessInfo.IsFree)
	assert.Equal(t, "free", accessInfo.AccessReason)

	chapterRepo.AssertExpectations(t)
}

func TestChapterPurchaseService_CheckChapterAccess_PurchasedChapter(t *testing.T) {
	chapterRepo := new(MockChapterRepository)
	purchaseRepo := new(MockChapterPurchaseRepository)
	bookRepo := new(MockBookStoreRepository)
	walletService := new(MockWalletService)
	cacheService := new(MockCacheService)

	service := NewChapterPurchaseService(chapterRepo, purchaseRepo, bookRepo, walletService, cacheService)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	chapterID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()

	chapter := createTestChapter(chapterID, bookID, 1, false, 1.99)
	existingPurchase := &bookstore.ChapterPurchase{
		ID:           primitive.NewObjectID(),
		UserID:       userID,
		ChapterID:    chapterID,
		PurchaseTime: time.Now(),
	}

	chapterRepo.On("GetByID", ctx, chapterID).Return(chapter, nil)
	purchaseRepo.On("CheckUserPurchasedChapter", ctx, userID, chapterID).Return(true, nil)
	purchaseRepo.On("GetByUserAndChapter", ctx, userID, chapterID).Return(existingPurchase, nil)

	accessInfo, err := service.CheckChapterAccess(ctx, userID, chapterID)

	assert.NoError(t, err)
	assert.NotNil(t, accessInfo)
	assert.True(t, accessInfo.CanAccess)
	assert.True(t, accessInfo.IsPurchased)
	assert.Equal(t, "purchased", accessInfo.AccessReason)
	assert.NotNil(t, accessInfo.PurchaseTime)

	chapterRepo.AssertExpectations(t)
	purchaseRepo.AssertExpectations(t)
}

func TestChapterPurchaseService_CheckChapterAccess_NoAccess(t *testing.T) {
	chapterRepo := new(MockChapterRepository)
	purchaseRepo := new(MockChapterPurchaseRepository)
	bookRepo := new(MockBookStoreRepository)
	walletService := new(MockWalletService)
	cacheService := new(MockCacheService)

	service := NewChapterPurchaseService(chapterRepo, purchaseRepo, bookRepo, walletService, cacheService)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	chapterID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()

	chapter := createTestChapter(chapterID, bookID, 1, false, 1.99)

	chapterRepo.On("GetByID", ctx, chapterID).Return(chapter, nil)
	purchaseRepo.On("CheckUserPurchasedChapter", ctx, userID, chapterID).Return(false, nil)
	purchaseRepo.On("CheckUserPurchasedBook", ctx, userID, bookID).Return(false, nil)

	accessInfo, err := service.CheckChapterAccess(ctx, userID, chapterID)

	assert.NoError(t, err)
	assert.NotNil(t, accessInfo)
	assert.False(t, accessInfo.CanAccess)
	assert.False(t, accessInfo.IsPurchased)
	assert.False(t, accessInfo.IsFree)

	chapterRepo.AssertExpectations(t)
	purchaseRepo.AssertExpectations(t)
}

// Test: GetChapterPurchases (Pagination)

func TestChapterPurchaseService_GetChapterPurchases_Success(t *testing.T) {
	chapterRepo := new(MockChapterRepository)
	purchaseRepo := new(MockChapterPurchaseRepository)
	bookRepo := new(MockBookStoreRepository)
	walletService := new(MockWalletService)
	cacheService := new(MockCacheService)

	service := NewChapterPurchaseService(chapterRepo, purchaseRepo, bookRepo, walletService, cacheService)

	ctx := context.Background()
	userID := primitive.NewObjectID()

	purchases := []*bookstore.ChapterPurchase{
		{
			ID:        primitive.NewObjectID(),
			UserID:    userID,
			ChapterID: primitive.NewObjectID(),
			Price:     1.99,
		},
		{
			ID:        primitive.NewObjectID(),
			UserID:    userID,
			ChapterID: primitive.NewObjectID(),
			Price:     2.99,
		},
	}

	purchaseRepo.On("GetByUser", ctx, userID, 1, 20).Return(purchases, int64(2), nil)

	result, total, err := service.GetChapterPurchases(ctx, userID, 1, 20)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(2), total)
	assert.Len(t, result, 2)

	purchaseRepo.AssertExpectations(t)
}

func TestChapterPurchaseService_GetChapterPurchases_EmptyUserID(t *testing.T) {
	chapterRepo := new(MockChapterRepository)
	purchaseRepo := new(MockChapterPurchaseRepository)
	bookRepo := new(MockBookStoreRepository)
	walletService := new(MockWalletService)
	cacheService := new(MockCacheService)

	service := NewChapterPurchaseService(chapterRepo, purchaseRepo, bookRepo, walletService, cacheService)

	ctx := context.Background()

	result, total, err := service.GetChapterPurchases(ctx, primitive.NilObjectID, 1, 20)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, int64(0), total)
	assert.Contains(t, err.Error(), "user ID cannot be empty")
}

// Test: CalculateBookPrice

func TestChapterPurchaseService_CalculateBookPrice_Success(t *testing.T) {
	chapterRepo := new(MockChapterRepository)
	purchaseRepo := new(MockChapterPurchaseRepository)
	bookRepo := new(MockBookStoreRepository)
	walletService := new(MockWalletService)
	cacheService := new(MockCacheService)

	service := NewChapterPurchaseService(chapterRepo, purchaseRepo, bookRepo, walletService, cacheService)

	ctx := context.Background()
	bookID := primitive.NewObjectID()

	paidChapters := []*bookstore.Chapter{
		createTestChapter(primitive.NewObjectID(), bookID, 3, false, 1.99),
		createTestChapter(primitive.NewObjectID(), bookID, 4, false, 2.99),
		createTestChapter(primitive.NewObjectID(), bookID, 5, false, 3.99),
	}

	chapterRepo.On("GetPaidChapters", ctx, bookID, 10000, 0).Return(paidChapters, nil)

	originalPrice, discountedPrice, err := service.CalculateBookPrice(ctx, bookID)

	assert.NoError(t, err)
	assert.Equal(t, 8.97, originalPrice)
	assert.True(t, discountedPrice < originalPrice) // Should have 20% discount
	assert.Equal(t, 7.176, discountedPrice)

	chapterRepo.AssertExpectations(t)
}
