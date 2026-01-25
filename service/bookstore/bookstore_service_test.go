package bookstore

import (
	"context"
	"errors"
	"testing"
	"time"

	bookstoreModel "Qingyu_backend/models/bookstore"
	"Qingyu_backend/repository/interfaces/infrastructure"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// =========================
// 辅助函数
// =========================

// newTestBook 创建一个测试用的 Book，自动生成 ID
func newTestBook(title, author string, status bookstoreModel.BookStatus) *bookstoreModel.Book {
	book := &bookstoreModel.Book{
		Title:  title,
		Author: author,
		Status: status,
	}
	book.ID = primitive.NewObjectID()
	return book
}

// newTestCategory 创建一个测试用的 Category，自动生成 ID
func newTestCategory(name string) *bookstoreModel.Category {
	category := &bookstoreModel.Category{Name: name}
	category.ID = primitive.NewObjectID().Hex()
	return category
}

// =========================
// 简化的 Mock Repository 实现
// =========================

// MockBookRepositoryForService Mock书籍仓储 - 仅包含service中使用的方法
type MockBookRepositoryForService struct {
	mock.Mock
}

// 基础CRUD方法
func (m *MockBookRepositoryForService) Create(ctx context.Context, book *bookstoreModel.Book) error {
	args := m.Called(ctx, book)
	return args.Error(0)
}

func (m *MockBookRepositoryForService) GetByID(ctx context.Context, id string) (*bookstoreModel.Book, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreModel.Book), args.Error(1)
}

func (m *MockBookRepositoryForService) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockBookRepositoryForService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBookRepositoryForService) List(ctx context.Context, filter infrastructure.Filter) ([]*bookstoreModel.Book, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Book), args.Error(1)
}

func (m *MockBookRepositoryForService) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookRepositoryForService) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockBookRepositoryForService) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// 列表查询方法
func (m *MockBookRepositoryForService) GetByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*bookstoreModel.Book, error) {
	args := m.Called(ctx, categoryID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Book), args.Error(1)
}

func (m *MockBookRepositoryForService) GetByAuthor(ctx context.Context, author string, limit, offset int) ([]*bookstoreModel.Book, error) {
	args := m.Called(ctx, author, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Book), args.Error(1)
}

func (m *MockBookRepositoryForService) GetByAuthorID(ctx context.Context, authorID string, limit, offset int) ([]*bookstoreModel.Book, error) {
	args := m.Called(ctx, authorID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Book), args.Error(1)
}

func (m *MockBookRepositoryForService) GetByStatus(ctx context.Context, status bookstoreModel.BookStatus, limit, offset int) ([]*bookstoreModel.Book, error) {
	args := m.Called(ctx, status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Book), args.Error(1)
}

func (m *MockBookRepositoryForService) GetRecommended(ctx context.Context, limit, offset int) ([]*bookstoreModel.Book, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Book), args.Error(1)
}

func (m *MockBookRepositoryForService) GetFeatured(ctx context.Context, limit, offset int) ([]*bookstoreModel.Book, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Book), args.Error(1)
}

func (m *MockBookRepositoryForService) GetHotBooks(ctx context.Context, limit, offset int) ([]*bookstoreModel.Book, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Book), args.Error(1)
}

func (m *MockBookRepositoryForService) GetNewReleases(ctx context.Context, limit, offset int) ([]*bookstoreModel.Book, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Book), args.Error(1)
}

func (m *MockBookRepositoryForService) GetFreeBooks(ctx context.Context, limit, offset int) ([]*bookstoreModel.Book, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Book), args.Error(1)
}

func (m *MockBookRepositoryForService) GetByPriceRange(ctx context.Context, minPrice, maxPrice float64, limit, offset int) ([]*bookstoreModel.Book, error) {
	args := m.Called(ctx, minPrice, maxPrice, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Book), args.Error(1)
}

// 搜索方法
func (m *MockBookRepositoryForService) Search(ctx context.Context, keyword string, limit, offset int) ([]*bookstoreModel.Book, error) {
	args := m.Called(ctx, keyword, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Book), args.Error(1)
}

func (m *MockBookRepositoryForService) SearchWithFilter(ctx context.Context, filter *bookstoreModel.BookFilter) ([]*bookstoreModel.Book, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Book), args.Error(1)
}

// 统计方法
func (m *MockBookRepositoryForService) CountByCategory(ctx context.Context, categoryID string) (int64, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookRepositoryForService) CountByAuthor(ctx context.Context, author string) (int64, error) {
	args := m.Called(ctx, author)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookRepositoryForService) CountByStatus(ctx context.Context, status bookstoreModel.BookStatus) (int64, error) {
	args := m.Called(ctx, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookRepositoryForService) CountByFilter(ctx context.Context, filter *bookstoreModel.BookFilter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

// 批量操作
func (m *MockBookRepositoryForService) BatchUpdateStatus(ctx context.Context, bookIDs []string, status bookstoreModel.BookStatus) error {
	args := m.Called(ctx, bookIDs, status)
	return args.Error(0)
}

func (m *MockBookRepositoryForService) BatchUpdateCategory(ctx context.Context, bookIDs []string, categoryIDs []string) error {
	args := m.Called(ctx, bookIDs, categoryIDs)
	return args.Error(0)
}

func (m *MockBookRepositoryForService) BatchUpdateRecommended(ctx context.Context, bookIDs []string, isRecommended bool) error {
	args := m.Called(ctx, bookIDs, isRecommended)
	return args.Error(0)
}

func (m *MockBookRepositoryForService) BatchUpdateFeatured(ctx context.Context, bookIDs []string, isFeatured bool) error {
	args := m.Called(ctx, bookIDs, isFeatured)
	return args.Error(0)
}

// 统计和计数操作
func (m *MockBookRepositoryForService) GetStats(ctx context.Context) (*bookstoreModel.BookStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreModel.BookStats), args.Error(1)
}

func (m *MockBookRepositoryForService) IncrementViewCount(ctx context.Context, bookID string) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

// 事务支持
func (m *MockBookRepositoryForService) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

// =========================
// MockCategoryRepositoryForService
// =========================

type MockCategoryRepositoryForService struct {
	mock.Mock
}

func (m *MockCategoryRepositoryForService) Create(ctx context.Context, category *bookstoreModel.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepositoryForService) GetByID(ctx context.Context, id string) (*bookstoreModel.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreModel.Category), args.Error(1)
}

func (m *MockCategoryRepositoryForService) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockCategoryRepositoryForService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCategoryRepositoryForService) List(ctx context.Context, filter infrastructure.Filter) ([]*bookstoreModel.Category, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Category), args.Error(1)
}

func (m *MockCategoryRepositoryForService) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCategoryRepositoryForService) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockCategoryRepositoryForService) GetByName(ctx context.Context, name string) (*bookstoreModel.Category, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreModel.Category), args.Error(1)
}

func (m *MockCategoryRepositoryForService) GetByParent(ctx context.Context, parentID string, limit, offset int) ([]*bookstoreModel.Category, error) {
	args := m.Called(ctx, parentID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Category), args.Error(1)
}

func (m *MockCategoryRepositoryForService) GetByLevel(ctx context.Context, level int, limit, offset int) ([]*bookstoreModel.Category, error) {
	args := m.Called(ctx, level, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Category), args.Error(1)
}

func (m *MockCategoryRepositoryForService) GetRootCategories(ctx context.Context) ([]*bookstoreModel.Category, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Category), args.Error(1)
}

func (m *MockCategoryRepositoryForService) GetCategoryTree(ctx context.Context) ([]*bookstoreModel.CategoryTree, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.CategoryTree), args.Error(1)
}

func (m *MockCategoryRepositoryForService) GetByParentID(ctx context.Context, parentID string) ([]*bookstoreModel.Category, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Category), args.Error(1)
}

func (m *MockCategoryRepositoryForService) GetBySlug(ctx context.Context, slug string) (*bookstoreModel.Category, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreModel.Category), args.Error(1)
}

func (m *MockCategoryRepositoryForService) BatchUpdateStatus(ctx context.Context, categoryIDs []string, isActive bool) error {
	args := m.Called(ctx, categoryIDs, isActive)
	return args.Error(0)
}

func (m *MockCategoryRepositoryForService) CountByParent(ctx context.Context, parentID string) (int64, error) {
	args := m.Called(ctx, parentID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCategoryRepositoryForService) UpdateBookCount(ctx context.Context, categoryID string, count int64) error {
	args := m.Called(ctx, categoryID, count)
	return args.Error(0)
}

func (m *MockCategoryRepositoryForService) GetChildren(ctx context.Context, parentID string) ([]*bookstoreModel.Category, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Category), args.Error(1)
}

func (m *MockCategoryRepositoryForService) GetAncestors(ctx context.Context, categoryID string) ([]*bookstoreModel.Category, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Category), args.Error(1)
}

func (m *MockCategoryRepositoryForService) GetDescendants(ctx context.Context, categoryID string) ([]*bookstoreModel.Category, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Category), args.Error(1)
}

func (m *MockCategoryRepositoryForService) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

// =========================
// MockBannerRepositoryForService
// =========================

type MockBannerRepositoryForService struct {
	mock.Mock
}

func (m *MockBannerRepositoryForService) Create(ctx context.Context, banner *bookstoreModel.Banner) error {
	args := m.Called(ctx, banner)
	return args.Error(0)
}

func (m *MockBannerRepositoryForService) GetByID(ctx context.Context, id string) (*bookstoreModel.Banner, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreModel.Banner), args.Error(1)
}

func (m *MockBannerRepositoryForService) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockBannerRepositoryForService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBannerRepositoryForService) List(ctx context.Context, filter infrastructure.Filter) ([]*bookstoreModel.Banner, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Banner), args.Error(1)
}

func (m *MockBannerRepositoryForService) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBannerRepositoryForService) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockBannerRepositoryForService) GetActive(ctx context.Context, limit, offset int) ([]*bookstoreModel.Banner, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Banner), args.Error(1)
}

func (m *MockBannerRepositoryForService) IncrementClickCount(ctx context.Context, bannerID string) error {
	args := m.Called(ctx, bannerID)
	return args.Error(0)
}

func (m *MockBannerRepositoryForService) GetClickStats(ctx context.Context, bannerID string) (int64, error) {
	args := m.Called(ctx, bannerID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBannerRepositoryForService) BatchUpdateStatus(ctx context.Context, bannerIDs []string, isActive bool) error {
	args := m.Called(ctx, bannerIDs, isActive)
	return args.Error(0)
}

func (m *MockBannerRepositoryForService) GetByTargetType(ctx context.Context, targetType string, limit, offset int) ([]*bookstoreModel.Banner, error) {
	args := m.Called(ctx, targetType, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Banner), args.Error(1)
}

func (m *MockBannerRepositoryForService) GetByTimeRange(ctx context.Context, startTime, endTime *time.Time, limit, offset int) ([]*bookstoreModel.Banner, error) {
	args := m.Called(ctx, startTime, endTime, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.Banner), args.Error(1)
}

func (m *MockBannerRepositoryForService) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

// =========================
// MockRankingRepositoryForService
// =========================

type MockRankingRepositoryForService struct {
	mock.Mock
}

func (m *MockRankingRepositoryForService) Create(ctx context.Context, item *bookstoreModel.RankingItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockRankingRepositoryForService) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstoreModel.RankingItem, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreModel.RankingItem), args.Error(1)
}

func (m *MockRankingRepositoryForService) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockRankingRepositoryForService) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRankingRepositoryForService) List(ctx context.Context, filter infrastructure.Filter) ([]*bookstoreModel.RankingItem, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.RankingItem), args.Error(1)
}

func (m *MockRankingRepositoryForService) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRankingRepositoryForService) Exists(ctx context.Context, id primitive.ObjectID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockRankingRepositoryForService) GetByType(ctx context.Context, rankingType bookstoreModel.RankingType, period string, limit, offset int) ([]*bookstoreModel.RankingItem, error) {
	args := m.Called(ctx, rankingType, period, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.RankingItem), args.Error(1)
}

func (m *MockRankingRepositoryForService) GetByTypeWithBooks(ctx context.Context, rankingType bookstoreModel.RankingType, period string, limit, offset int) ([]*bookstoreModel.RankingItem, error) {
	args := m.Called(ctx, rankingType, period, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.RankingItem), args.Error(1)
}

func (m *MockRankingRepositoryForService) GetByBookID(ctx context.Context, bookID primitive.ObjectID, rankingType bookstoreModel.RankingType, period string) (*bookstoreModel.RankingItem, error) {
	args := m.Called(ctx, bookID, rankingType, period)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreModel.RankingItem), args.Error(1)
}

func (m *MockRankingRepositoryForService) GetByPeriod(ctx context.Context, period string, limit, offset int) ([]*bookstoreModel.RankingItem, error) {
	args := m.Called(ctx, period, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.RankingItem), args.Error(1)
}

func (m *MockRankingRepositoryForService) GetRankingStats(ctx context.Context, rankingType bookstoreModel.RankingType, period string) (*bookstoreModel.RankingStats, error) {
	args := m.Called(ctx, rankingType, period)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreModel.RankingStats), args.Error(1)
}

func (m *MockRankingRepositoryForService) CountByType(ctx context.Context, rankingType bookstoreModel.RankingType, period string) (int64, error) {
	args := m.Called(ctx, rankingType, period)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRankingRepositoryForService) GetTopBooks(ctx context.Context, rankingType bookstoreModel.RankingType, period string, limit int) ([]*bookstoreModel.RankingItem, error) {
	args := m.Called(ctx, rankingType, period, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.RankingItem), args.Error(1)
}

func (m *MockRankingRepositoryForService) UpsertRankingItem(ctx context.Context, item *bookstoreModel.RankingItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockRankingRepositoryForService) BatchUpsertRankingItems(ctx context.Context, items []*bookstoreModel.RankingItem) error {
	args := m.Called(ctx, items)
	return args.Error(0)
}

func (m *MockRankingRepositoryForService) UpdateRankings(ctx context.Context, rankingType bookstoreModel.RankingType, period string, items []*bookstoreModel.RankingItem) error {
	args := m.Called(ctx, rankingType, period, items)
	return args.Error(0)
}

func (m *MockRankingRepositoryForService) DeleteByPeriod(ctx context.Context, period string) error {
	args := m.Called(ctx, period)
	return args.Error(0)
}

func (m *MockRankingRepositoryForService) DeleteByType(ctx context.Context, rankingType bookstoreModel.RankingType) error {
	args := m.Called(ctx, rankingType)
	return args.Error(0)
}

func (m *MockRankingRepositoryForService) DeleteExpiredRankings(ctx context.Context, beforeDate time.Time) error {
	args := m.Called(ctx, beforeDate)
	return args.Error(0)
}

func (m *MockRankingRepositoryForService) CalculateRealtimeRanking(ctx context.Context, period string) ([]*bookstoreModel.RankingItem, error) {
	args := m.Called(ctx, period)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.RankingItem), args.Error(1)
}

func (m *MockRankingRepositoryForService) CalculateWeeklyRanking(ctx context.Context, period string) ([]*bookstoreModel.RankingItem, error) {
	args := m.Called(ctx, period)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.RankingItem), args.Error(1)
}

func (m *MockRankingRepositoryForService) CalculateMonthlyRanking(ctx context.Context, period string) ([]*bookstoreModel.RankingItem, error) {
	args := m.Called(ctx, period)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.RankingItem), args.Error(1)
}

func (m *MockRankingRepositoryForService) CalculateNewbieRanking(ctx context.Context, period string) ([]*bookstoreModel.RankingItem, error) {
	args := m.Called(ctx, period)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.RankingItem), args.Error(1)
}

func (m *MockRankingRepositoryForService) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

// =========================
// 测试辅助函数
// =========================

// setupBookstoreServiceForTest 创建测试用的service实例
func setupBookstoreServiceForTest() (*BookstoreServiceImpl, *MockBookRepositoryForService, *MockCategoryRepositoryForService, *MockBannerRepositoryForService, *MockRankingRepositoryForService) {
	mockBookRepo := new(MockBookRepositoryForService)
	mockCategoryRepo := new(MockCategoryRepositoryForService)
	mockBannerRepo := new(MockBannerRepositoryForService)
	mockRankingRepo := new(MockRankingRepositoryForService)

	service := &BookstoreServiceImpl{
		bookRepo:     mockBookRepo,
		categoryRepo: mockCategoryRepo,
		bannerRepo:   mockBannerRepo,
		rankingRepo:  mockRankingRepo,
	}

	return service, mockBookRepo, mockCategoryRepo, mockBannerRepo, mockRankingRepo
}

// =========================
// 书籍相关方法测试
// =========================

// TestBookstoreService_GetAllBooks 测试获取所有书籍
func TestBookstoreService_GetAllBooks(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockBookRepositoryForService)
		wantErr        bool
		errContains    string
		validateResult func(*testing.T, []*bookstoreModel.Book, int64, error)
	}{
		{
			name: "成功获取书籍列表",
			setupMock: func(m *MockBookRepositoryForService) {
				books := []*bookstoreModel.Book{
					func() *bookstoreModel.Book {
						b := &bookstoreModel.Book{Title: "书籍1", Status: bookstoreModel.BookStatusOngoing}
						b.ID = primitive.NewObjectID()
						return b
					}(),
					func() *bookstoreModel.Book {
						b := &bookstoreModel.Book{Title: "书籍2", Status: bookstoreModel.BookStatusOngoing}
						b.ID = primitive.NewObjectID()
						return b
					}(),
				}
				m.On("GetHotBooks", mock.Anything, 10, 0).Return(books, nil)
				m.On("CountByFilter", mock.Anything, mock.Anything).Return(int64(2), nil)
			},
			wantErr: false,
			validateResult: func(t *testing.T, books []*bookstoreModel.Book, total int64, err error) {
				require.NoError(t, err)
				assert.Len(t, books, 2)
				assert.Equal(t, int64(2), total)
			},
		},
		{
			name: "仓储错误",
			setupMock: func(m *MockBookRepositoryForService) {
				m.On("GetHotBooks", mock.Anything, 10, 0).Return(nil, errors.New("数据库错误"))
			},
			wantErr:     true,
			errContains: "failed to get all books",
			validateResult: func(t *testing.T, books []*bookstoreModel.Book, total int64, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "failed to get all books")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service, mockBookRepo, _, _, _ := setupBookstoreServiceForTest()
			ctx := context.Background()
			tt.setupMock(mockBookRepo)

			// Act
			books, total, err := service.GetAllBooks(ctx, 1, 10)

			// Assert
			tt.validateResult(t, books, total, err)
			mockBookRepo.AssertExpectations(t)
		})
	}
}

// TestBookstoreService_GetBookByID 测试根据ID获取书籍
func TestBookstoreService_GetBookByID(t *testing.T) {
	tests := []struct {
		name        string
		bookID      string
		setupMock   func(*MockBookRepositoryForService)
		wantErr     bool
		errContains string
	}{
		{
			name:   "成功获取书籍",
			bookID: primitive.NewObjectID().Hex(),
			setupMock: func(m *MockBookRepositoryForService) {
				book := newTestBook("测试书籍", "", bookstoreModel.BookStatusOngoing)
				m.On("GetByID", mock.Anything, mock.Anything).Return(book, nil)
			},
			wantErr: false,
		},
		{
			name:        "无效的书籍ID",
			bookID:      "invalid-id",
			setupMock:   func(m *MockBookRepositoryForService) {},
			wantErr:     true,
			errContains: "invalid book ID",
		},
		{
			name:   "书籍不存在",
			bookID: primitive.NewObjectID().Hex(),
			setupMock: func(m *MockBookRepositoryForService) {
				m.On("GetByID", mock.Anything, mock.Anything).Return(nil, nil)
			},
			wantErr:     true,
			errContains: "book not found",
		},
		{
			name:   "书籍未发布",
			bookID: primitive.NewObjectID().Hex(),
			setupMock: func(m *MockBookRepositoryForService) {
				book := newTestBook("草稿书籍", "", "draft")
				m.On("GetByID", mock.Anything, mock.Anything).Return(book, nil)
			},
			wantErr:     true,
			errContains: "book not available",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service, mockBookRepo, _, _, _ := setupBookstoreServiceForTest()
			ctx := context.Background()
			tt.setupMock(mockBookRepo)

			// Act
			book, err := service.GetBookByID(ctx, tt.bookID)

			// Assert
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, book)
			}
			mockBookRepo.AssertExpectations(t)
		})
	}
}

// TestBookstoreService_GetBooksByCategory 测试根据分类获取书籍
func TestBookstoreService_GetBooksByCategory(t *testing.T) {
	// Arrange
	service, mockBookRepo, _, _, _ := setupBookstoreServiceForTest()
	ctx := context.Background()

	categoryID := primitive.NewObjectID()
	books := []*bookstoreModel.Book{
		newTestBook("玄幻小说1", "", bookstoreModel.BookStatusOngoing),
		newTestBook("玄幻小说2", "", bookstoreModel.BookStatusOngoing),
	}

	mockBookRepo.On("GetByCategory", ctx, categoryID, 20, 0).Return(books, nil)
	mockBookRepo.On("CountByCategory", ctx, categoryID).Return(int64(2), nil)

	// Act
	result, total, err := service.GetBooksByCategory(ctx, categoryID.Hex(), 1, 20)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	mockBookRepo.AssertExpectations(t)
}

// TestBookstoreService_GetBooksByAuthorID 测试根据作者ID获取书籍
func TestBookstoreService_GetBooksByAuthorID(t *testing.T) {
	// Arrange
	service, mockBookRepo, _, _, _ := setupBookstoreServiceForTest()
	ctx := context.Background()

	authorID := primitive.NewObjectID()
	books := []*bookstoreModel.Book{
		newTestBook("作品1", "", bookstoreModel.BookStatusOngoing),
	}

	mockBookRepo.On("GetByAuthorID", ctx, authorID.Hex(), 20, 0).Return(books, nil)
	mockBookRepo.On("CountByAuthor", ctx, mock.Anything).Return(int64(1), nil)

	// Act
	result, total, err := service.GetBooksByAuthorID(ctx, authorID.Hex(), 1, 20)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
	mockBookRepo.AssertExpectations(t)
}

// TestBookstoreService_GetRecommendedBooks 测试获取推荐书籍
func TestBookstoreService_GetRecommendedBooks(t *testing.T) {
	// Arrange
	service, mockBookRepo, _, _, _ := setupBookstoreServiceForTest()
	ctx := context.Background()

	books := []*bookstoreModel.Book{
		func() *bookstoreModel.Book { b := &bookstoreModel.Book{Title: "$1"}; b.ID = primitive.NewObjectID(); return b }(),
	}

	mockBookRepo.On("GetRecommended", ctx, 10, 0).Return(books, nil)
	mockBookRepo.On("Count", ctx, nil).Return(int64(1), nil)

	// Act
	result, total, err := service.GetRecommendedBooks(ctx, 1, 10)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
	mockBookRepo.AssertExpectations(t)
}

// TestBookstoreService_GetFeaturedBooks 测试获取精选书籍
func TestBookstoreService_GetFeaturedBooks(t *testing.T) {
	// Arrange
	service, mockBookRepo, _, _, _ := setupBookstoreServiceForTest()
	ctx := context.Background()

	books := []*bookstoreModel.Book{
		func() *bookstoreModel.Book { b := &bookstoreModel.Book{Title: "$1"}; b.ID = primitive.NewObjectID(); return b }(),
	}

	mockBookRepo.On("GetFeatured", ctx, 10, 0).Return(books, nil)
	mockBookRepo.On("CountByFilter", ctx, mock.Anything).Return(int64(1), nil)

	// Act
	result, total, err := service.GetFeaturedBooks(ctx, 1, 10)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
	mockBookRepo.AssertExpectations(t)
}

// TestBookstoreService_GetHotBooks 测试获取热门书籍
func TestBookstoreService_GetHotBooks(t *testing.T) {
	// Arrange
	service, mockBookRepo, _, _, _ := setupBookstoreServiceForTest()
	ctx := context.Background()

	books := []*bookstoreModel.Book{
		func() *bookstoreModel.Book { b := &bookstoreModel.Book{Title: "$1"}; b.ID = primitive.NewObjectID(); return b }(),
	}

	mockBookRepo.On("GetHotBooks", ctx, 10, 0).Return(books, nil)
	mockBookRepo.On("CountByFilter", ctx, mock.Anything).Return(int64(1), nil)

	// Act
	result, total, err := service.GetHotBooks(ctx, 1, 10)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
	mockBookRepo.AssertExpectations(t)
}

// TestBookstoreService_GetNewReleases 测试获取新书
func TestBookstoreService_GetNewReleases(t *testing.T) {
	// Arrange
	service, mockBookRepo, _, _, _ := setupBookstoreServiceForTest()
	ctx := context.Background()

	books := []*bookstoreModel.Book{
		func() *bookstoreModel.Book { b := &bookstoreModel.Book{Title: "$1"}; b.ID = primitive.NewObjectID(); return b }(),
	}

	mockBookRepo.On("GetNewReleases", ctx, 10, 0).Return(books, nil)
	mockBookRepo.On("CountByFilter", ctx, mock.Anything).Return(int64(1), nil)

	// Act
	result, total, err := service.GetNewReleases(ctx, 1, 10)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
	mockBookRepo.AssertExpectations(t)
}

// TestBookstoreService_GetFreeBooks 测试获取免费书籍
func TestBookstoreService_GetFreeBooks(t *testing.T) {
	// Arrange
	service, mockBookRepo, _, _, _ := setupBookstoreServiceForTest()
	ctx := context.Background()

	books := []*bookstoreModel.Book{
		func() *bookstoreModel.Book { b := &bookstoreModel.Book{Title: "$1"}; b.ID = primitive.NewObjectID(); return b }(),
	}

	mockBookRepo.On("GetFreeBooks", ctx, 10, 0).Return(books, nil)
	mockBookRepo.On("CountByFilter", ctx, mock.Anything).Return(int64(1), nil)

	// Act
	result, total, err := service.GetFreeBooks(ctx, 1, 10)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
	mockBookRepo.AssertExpectations(t)
}

// TestBookstoreService_SearchBooks 测试搜索书籍
func TestBookstoreService_SearchBooks(t *testing.T) {
	tests := []struct {
		name        string
		keyword     string
		setupMock   func(*MockBookRepositoryForService)
		wantErr     bool
		errContains string
	}{
		{
			name:    "成功搜索书籍",
			keyword: "玄幻",
			setupMock: func(m *MockBookRepositoryForService) {
				books := []*bookstoreModel.Book{
					func() *bookstoreModel.Book { b := &bookstoreModel.Book{Title: "$1"}; b.ID = primitive.NewObjectID(); return b }(),
				}
				m.On("Search", mock.Anything, "玄幻", 20, 0).Return(books, nil)
				m.On("CountByFilter", mock.Anything, mock.Anything).Return(int64(1), nil)
			},
			wantErr: false,
		},
		{
			name:        "空关键词",
			keyword:     "",
			setupMock:   func(m *MockBookRepositoryForService) {},
			wantErr:     true,
			errContains: "keyword is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service, mockBookRepo, _, _, _ := setupBookstoreServiceForTest()
			ctx := context.Background()
			tt.setupMock(mockBookRepo)

			// Act
			books, total, err := service.SearchBooks(ctx, tt.keyword, 1, 20)

			// Assert
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, books)
				assert.GreaterOrEqual(t, total, int64(0))
			}
			mockBookRepo.AssertExpectations(t)
		})
	}
}

// TestBookstoreService_SearchBooksWithFilter 测试高级搜索
func TestBookstoreService_SearchBooksWithFilter(t *testing.T) {
	// Arrange
	service, mockBookRepo, _, _, _ := setupBookstoreServiceForTest()
	ctx := context.Background()

	filter := &bookstoreModel.BookFilter{
		Keyword: stringPtr("测试"),
	}
	books := []*bookstoreModel.Book{
		func() *bookstoreModel.Book { b := &bookstoreModel.Book{Title: "$1"}; b.ID = primitive.NewObjectID(); return b }(),
	}

	mockBookRepo.On("SearchWithFilter", ctx, filter).Return(books, nil)

	// Act
	result, total, err := service.SearchBooksWithFilter(ctx, filter)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
	mockBookRepo.AssertExpectations(t)
}

// TestBookstoreService_GetBookStats 测试获取书籍统计
func TestBookstoreService_GetBookStats(t *testing.T) {
	// Arrange
	service, mockBookRepo, _, _, _ := setupBookstoreServiceForTest()
	ctx := context.Background()

	stats := &bookstoreModel.BookStats{
		TotalBooks: 1000,
	}

	mockBookRepo.On("GetStats", ctx).Return(stats, nil)

	// Act
	result, err := service.GetBookStats(ctx)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(1000), result.TotalBooks)
	mockBookRepo.AssertExpectations(t)
}

// TestBookstoreService_IncrementBookView 测试增加书籍浏览量
func TestBookstoreService_IncrementBookView(t *testing.T) {
	tests := []struct {
		name        string
		bookID      string
		setupMock   func(*MockBookRepositoryForService)
		wantErr     bool
		errContains string
	}{
		{
			name:   "成功增加浏览量",
			bookID: primitive.NewObjectID().Hex(),
			setupMock: func(m *MockBookRepositoryForService) {
				book := newTestBook("测试书籍", "", bookstoreModel.BookStatusOngoing)
				m.On("GetByID", mock.Anything, mock.Anything).Return(book, nil)
				m.On("IncrementViewCount", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "无效的书籍ID",
			bookID:      "invalid-id",
			setupMock:   func(m *MockBookRepositoryForService) {},
			wantErr:     true,
			errContains: "invalid book ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service, mockBookRepo, _, _, _ := setupBookstoreServiceForTest()
			ctx := context.Background()
			tt.setupMock(mockBookRepo)

			// Act
			err := service.IncrementBookView(ctx, tt.bookID)

			// Assert
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				require.NoError(t, err)
			}
			mockBookRepo.AssertExpectations(t)
		})
	}
}

// =========================
// 分类相关方法测试
// =========================

// TestBookstoreService_GetCategoryTree 测试获取分类树
func TestBookstoreService_GetCategoryTree(t *testing.T) {
	// Arrange
	service, _, mockCategoryRepo, _, _ := setupBookstoreServiceForTest()
	ctx := context.Background()

	tree := []*bookstoreModel.CategoryTree{
		{
			Category: bookstoreModel.Category{
				ID:   primitive.NewObjectID().Hex(),
				Name: "玄幻",
			},
			Children: []*bookstoreModel.CategoryTree{},
		},
	}

	mockCategoryRepo.On("GetCategoryTree", ctx).Return(tree, nil)

	// Act
	result, err := service.GetCategoryTree(ctx)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	mockCategoryRepo.AssertExpectations(t)
}

// TestBookstoreService_GetCategoryByID 测试根据ID获取分类
func TestBookstoreService_GetCategoryByID(t *testing.T) {
	tests := []struct {
		name        string
		categoryID  string
		setupMock   func(*MockCategoryRepositoryForService)
		wantErr     bool
		errContains string
	}{
		{
			name:       "成功获取分类",
			categoryID: primitive.NewObjectID().Hex(),
			setupMock: func(m *MockCategoryRepositoryForService) {
				categoryID := primitive.NewObjectID().Hex()
				category := &bookstoreModel.Category{
					ID:       categoryID,
					Name:     "玄幻",
					IsActive: true,
				}
				m.On("GetByID", mock.Anything, mock.Anything).Return(category, nil)
			},
			wantErr: false,
		},
		{
			name:       "分类未激活",
			categoryID: primitive.NewObjectID().Hex(),
			setupMock: func(m *MockCategoryRepositoryForService) {
				categoryID := primitive.NewObjectID().Hex()
				category := &bookstoreModel.Category{
					ID:       categoryID,
					Name:     "玄幻",
					IsActive: false,
				}
				m.On("GetByID", mock.Anything, mock.Anything).Return(category, nil)
			},
			wantErr:     true,
			errContains: "category not available",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service, _, mockCategoryRepo, _, _ := setupBookstoreServiceForTest()
			ctx := context.Background()
			tt.setupMock(mockCategoryRepo)

			// Act
			category, err := service.GetCategoryByID(ctx, tt.categoryID)

			// Assert
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, category)
			}
			mockCategoryRepo.AssertExpectations(t)
		})
	}
}

// TestBookstoreService_GetRootCategories 测试获取根分类
func TestBookstoreService_GetRootCategories(t *testing.T) {
	// Arrange
	service, _, mockCategoryRepo, _, _ := setupBookstoreServiceForTest()
	ctx := context.Background()

	categories := []*bookstoreModel.Category{
		newTestCategory("玄幻"),
		newTestCategory("言情"),
	}

	mockCategoryRepo.On("GetRootCategories", ctx).Return(categories, nil)

	// Act
	result, err := service.GetRootCategories(ctx)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 2)
	mockCategoryRepo.AssertExpectations(t)
}

// =========================
// Banner相关方法测试
// =========================

// TestBookstoreService_GetActiveBanners 测试获取活跃Banner
func TestBookstoreService_GetActiveBanners(t *testing.T) {
	// Arrange
	service, _, _, mockBannerRepo, _ := setupBookstoreServiceForTest()
	ctx := context.Background()

	banners := []*bookstoreModel.Banner{
		func() *bookstoreModel.Banner { b := &bookstoreModel.Banner{Title: "Banner1"}; b.ID = primitive.NewObjectID(); return b }(),
		func() *bookstoreModel.Banner { b := &bookstoreModel.Banner{Title: "Banner2"}; b.ID = primitive.NewObjectID(); return b }(),
	}

	mockBannerRepo.On("GetActive", ctx, 5, 0).Return(banners, nil)

	// Act
	result, err := service.GetActiveBanners(ctx, 5)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 2)
	mockBannerRepo.AssertExpectations(t)
}

// TestBookstoreService_IncrementBannerClick 测试增加Banner点击数
func TestBookstoreService_IncrementBannerClick(t *testing.T) {
	tests := []struct {
		name        string
		bannerID    string
		setupMock   func(*MockBannerRepositoryForService)
		wantErr     bool
		errContains string
	}{
		{
			name:     "成功增加点击数",
			bannerID: primitive.NewObjectID().Hex(),
			setupMock: func(m *MockBannerRepositoryForService) {
				bannerID := primitive.NewObjectID()
				banner := &bookstoreModel.Banner{
					ID:       bannerID,
					Title:    "测试Banner",
					IsActive: true,
				}
				m.On("GetByID", mock.Anything, mock.Anything).Return(banner, nil)
				m.On("IncrementClickCount", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "无效的Banner ID",
			bannerID:    "invalid-id",
			setupMock:   func(m *MockBannerRepositoryForService) {},
			wantErr:     true,
			errContains: "invalid banner ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service, _, _, mockBannerRepo, _ := setupBookstoreServiceForTest()
			ctx := context.Background()
			tt.setupMock(mockBannerRepo)

			// Act
			err := service.IncrementBannerClick(ctx, tt.bannerID)

			// Assert
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				require.NoError(t, err)
			}
			mockBannerRepo.AssertExpectations(t)
		})
	}
}

// =========================
// 榜单相关方法测试
// =========================

// TestBookstoreService_GetRealtimeRanking 测试获取实时榜单
func TestBookstoreService_GetRealtimeRanking(t *testing.T) {
	// Arrange
	service, _, _, _, mockRankingRepo := setupBookstoreServiceForTest()
	ctx := context.Background()

	rankings := []*bookstoreModel.RankingItem{
		{BookID: primitive.NewObjectID(), Rank: 1},
		{BookID: primitive.NewObjectID(), Rank: 2},
	}

	mockRankingRepo.On("GetByTypeWithBooks", ctx, mock.Anything, mock.Anything, 10, 0).Return(rankings, nil)

	// Act
	result, err := service.GetRealtimeRanking(ctx, 10)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 2)
	mockRankingRepo.AssertExpectations(t)
}

// TestBookstoreService_GetWeeklyRanking 测试获取周榜
func TestBookstoreService_GetWeeklyRanking(t *testing.T) {
	// Arrange
	service, _, _, _, mockRankingRepo := setupBookstoreServiceForTest()
	ctx := context.Background()

	rankings := []*bookstoreModel.RankingItem{
		{BookID: primitive.NewObjectID(), Rank: 1},
	}

	mockRankingRepo.On("GetByTypeWithBooks", ctx, mock.Anything, mock.Anything, 10, 0).Return(rankings, nil)

	// Act
	result, err := service.GetWeeklyRanking(ctx, "", 10)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	mockRankingRepo.AssertExpectations(t)
}

// TestBookstoreService_GetMonthlyRanking 测试获取月榜
func TestBookstoreService_GetMonthlyRanking(t *testing.T) {
	// Arrange
	service, _, _, _, mockRankingRepo := setupBookstoreServiceForTest()
	ctx := context.Background()

	rankings := []*bookstoreModel.RankingItem{
		{BookID: primitive.NewObjectID(), Rank: 1},
	}

	mockRankingRepo.On("GetByTypeWithBooks", ctx, mock.Anything, mock.Anything, 10, 0).Return(rankings, nil)

	// Act
	result, err := service.GetMonthlyRanking(ctx, "", 10)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	mockRankingRepo.AssertExpectations(t)
}

// TestBookstoreService_GetNewbieRanking 测试获取新人榜
func TestBookstoreService_GetNewbieRanking(t *testing.T) {
	// Arrange
	service, _, _, _, mockRankingRepo := setupBookstoreServiceForTest()
	ctx := context.Background()

	rankings := []*bookstoreModel.RankingItem{
		{BookID: primitive.NewObjectID(), Rank: 1},
	}

	mockRankingRepo.On("GetByTypeWithBooks", ctx, mock.Anything, mock.Anything, 10, 0).Return(rankings, nil)

	// Act
	result, err := service.GetNewbieRanking(ctx, "", 10)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	mockRankingRepo.AssertExpectations(t)
}

// TestBookstoreService_GetRankingByType 测试根据类型获取榜单
func TestBookstoreService_GetRankingByType(t *testing.T) {
	// Arrange
	service, _, _, _, mockRankingRepo := setupBookstoreServiceForTest()
	ctx := context.Background()

	rankings := []*bookstoreModel.RankingItem{
		{BookID: primitive.NewObjectID(), Rank: 1},
	}

	mockRankingRepo.On("GetByTypeWithBooks", ctx, mock.Anything, mock.Anything, 10, 0).Return(rankings, nil)

	// Act
	result, err := service.GetRankingByType(ctx, bookstoreModel.RankingTypeWeekly, "", 10)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	mockRankingRepo.AssertExpectations(t)
}

// TestBookstoreService_UpdateRankings 测试更新榜单
func TestBookstoreService_UpdateRankings(t *testing.T) {
	// Arrange
	service, _, _, _, mockRankingRepo := setupBookstoreServiceForTest()
	ctx := context.Background()

	items := []*bookstoreModel.RankingItem{
		{BookID: primitive.NewObjectID(), Rank: 1},
	}

	mockRankingRepo.On("CalculateWeeklyRanking", ctx, mock.Anything).Return(items, nil)
	mockRankingRepo.On("UpdateRankings", ctx, mock.Anything, mock.Anything, items).Return(nil)

	// Act
	err := service.UpdateRankings(ctx, bookstoreModel.RankingTypeWeekly, "")

	// Assert
	require.NoError(t, err)
	mockRankingRepo.AssertExpectations(t)
}

// TestBookstoreService_UpdateRankings_UnsupportedType 测试更新不支持的榜单类型
func TestBookstoreService_UpdateRankings_UnsupportedType(t *testing.T) {
	// Arrange
	service, _, _, _, _ := setupBookstoreServiceForTest()
	ctx := context.Background()

	// Act
	err := service.UpdateRankings(ctx, "unsupported_type", "")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported ranking type")
}

// =========================
// 首页数据聚合测试
// =========================

// TestBookstoreService_GetHomepageData 测试获取首页数据
func TestBookstoreService_GetHomepageData(t *testing.T) {
	// Arrange
	service, mockBookRepo, mockCategoryRepo, mockBannerRepo, mockRankingRepo := setupBookstoreServiceForTest()
	ctx := context.Background()

	// Mock Banner
	banners := []*bookstoreModel.Banner{
		func() *bookstoreModel.Banner { b := &bookstoreModel.Banner{Title: "推荐书"}; b.ID = primitive.NewObjectID(); return b }(),
	}
	mockBannerRepo.On("GetActive", ctx, 5, 0).Return(banners, nil)

	// Mock Recommended Books
	books := []*bookstoreModel.Book{
		func() *bookstoreModel.Book { b := &bookstoreModel.Book{Title: "$1"}; b.ID = primitive.NewObjectID(); return b }(),
	}
	mockBookRepo.On("GetRecommended", ctx, 10, 0).Return(books, nil)
	mockBookRepo.On("Count", ctx, nil).Return(int64(1), nil)

	// Mock Featured Books
	mockBookRepo.On("GetFeatured", ctx, 10, 0).Return(books, nil)
	mockBookRepo.On("CountByFilter", ctx, mock.Anything).Return(int64(1), nil)

	// Mock Categories
	categories := []*bookstoreModel.Category{
		newTestCategory("玄幻"),
	}
	mockCategoryRepo.On("GetRootCategories", ctx).Return(categories, nil)

	// Mock Stats
	stats := &bookstoreModel.BookStats{TotalBooks: 1000}
	mockBookRepo.On("GetStats", ctx).Return(stats, nil)

	// Mock Rankings
	rankings := []*bookstoreModel.RankingItem{
		{BookID: primitive.NewObjectID(), Rank: 1},
	}
	mockRankingRepo.On("GetByTypeWithBooks", ctx, mock.Anything, mock.Anything, 10, 0).Return(rankings, nil).Times(3)

	// Act
	data, err := service.GetHomepageData(ctx)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, data)
	assert.Len(t, data.Banners, 1)
	assert.Len(t, data.RecommendedBooks, 1)
	assert.Len(t, data.FeaturedBooks, 1)
	assert.Len(t, data.Categories, 1)
	assert.Equal(t, int64(1000), data.Stats.TotalBooks)
	assert.Len(t, data.Rankings, 3)

	// Verify all mocks were called
	mockBannerRepo.AssertExpectations(t)
	mockBookRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)
	mockRankingRepo.AssertExpectations(t)
}

// =========================
// 辅助函数
// =========================

// stringPtr 返回字符串指针
func stringPtr(s string) *string {
	return &s
}
