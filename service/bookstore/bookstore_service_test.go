package bookstore

import (
	"context"
	"testing"
	"time"

	bookstoreModel "Qingyu_backend/models/bookstore"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ============ Mock 实现 ============

// MockBookRepository Mock书籍仓储
type MockBookRepository struct {
	books       map[string]*bookstoreModel.Book
	bookDetails map[string]*bookstoreModel.BookDetail
	nextID      int
}

func NewMockBookRepository() *MockBookRepository {
	return &MockBookRepository{
		books:       make(map[string]*bookstoreModel.Book),
		bookDetails: make(map[string]*bookstoreModel.BookDetail),
		nextID:      1,
	}
}

func (m *MockBookRepository) GetHotBooks(ctx context.Context, limit int, offset int) ([]*bookstoreModel.Book, error) {
	result := make([]*bookstoreModel.Book, 0)
	count := 0
	for _, book := range m.books {
		if count >= offset && len(result) < limit {
			result = append(result, book)
		}
		count++
	}
	return result, nil
}

func (m *MockBookRepository) GetByID(ctx context.Context, id string) (*bookstoreModel.Book, error) {
	return m.books[id], nil
}

func (m *MockBookRepository) Create(ctx context.Context, book *bookstoreModel.Book) error {
	book.ID = primitive.NewObjectID()
	m.books[book.ID.Hex()] = book
	return nil
}

func (m *MockBookRepository) Update(ctx context.Context, book *bookstoreModel.Book) error {
	m.books[book.ID.Hex()] = book
	return nil
}

func (m *MockBookRepository) Delete(ctx context.Context, id string) error {
	delete(m.books, id)
	return nil
}

func (m *MockBookRepository) GetByTitle(ctx context.Context, title string) (*bookstoreModel.Book, error) {
	for _, book := range m.books {
		if book.Title == title {
			return book, nil
		}
	}
	return nil, nil
}

func (m *MockBookRepository) GetByAuthorID(ctx context.Context, authorID primitive.ObjectID) ([]*bookstoreModel.Book, error) {
	result := make([]*bookstoreModel.Book, 0)
	for _, book := range m.books {
		if book.AuthorID == authorID {
			result = append(result, book)
		}
	}
	return result, nil
}

func (m *MockBookRepository) GetByCategoryID(ctx context.Context, categoryID primitive.ObjectID) ([]*bookstoreModel.Book, error) {
	return []*bookstoreModel.Book{}, nil
}

func (m *MockBookRepository) GetByTags(ctx context.Context, tags []string) ([]*bookstoreModel.Book, error) {
	return []*bookstoreModel.Book{}, nil
}

func (m *MockBookRepository) Search(ctx context.Context, keyword string) ([]*bookstoreModel.Book, error) {
	result := make([]*bookstoreModel.Book, 0)
	for _, book := range m.books {
		if contains(book.Title, keyword) || contains(book.Description, keyword) {
			result = append(result, book)
		}
	}
	return result, nil
}

func (m *MockBookRepository) IncrementViewCount(ctx context.Context, id string) error {
	if book, ok := m.books[id]; ok {
		book.ViewCount++
	}
	return nil
}

func (m *MockBookRepository) GetRecommended(ctx context.Context, limit int) ([]*bookstoreModel.Book, error) {
	return []*bookstoreModel.Book{}, nil
}

func (m *MockBookRepository) GetFeatured(ctx context.Context, limit int) ([]*bookstoreModel.Book, error) {
	return []*bookstoreModel.Book{}, nil
}

func (m *MockBookRepository) GetNewReleases(ctx context.Context, limit int) ([]*bookstoreModel.Book, error) {
	return []*bookstoreModel.Book{}, nil
}

func (m *MockBookRepository) GetFreeBooks(ctx context.Context, limit int) ([]*bookstoreModel.Book, error) {
	return []*bookstoreModel.Book{}, nil
}

func (m *MockBookRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.books)), nil
}

func (m *MockBookRepository) GetStats(ctx context.Context) (*bookstoreModel.BookStats, error) {
	return &bookstoreModel.BookStats{
		TotalBooks:     int64(len(m.books)),
		PublishedBooks: int64(len(m.books)),
	}, nil
}

// MockCategoryRepository Mock分类仓储
type MockCategoryRepository struct {
	categories map[string]*bookstoreModel.Category
	nextID     int
}

func NewMockCategoryRepository() *MockCategoryRepository {
	repo := &MockCategoryRepository{
		categories: make(map[string]*bookstoreModel.Category),
		nextID:     1,
	}

	// 添加默认分类
	fantasy := &bookstoreModel.Category{
		ID:          primitive.NewObjectID(),
		Name:        "奇幻",
		Slug:        "fantasy",
		Description: "奇幻小说分类",
		ParentID:    nil,
	}
	repo.categories[fantasy.ID.Hex()] = fantasy

	romance := &bookstoreModel.Category{
		ID:          primitive.NewObjectID(),
		Name:        "言情",
		Slug:        "romance",
		Description: "言情小说分类",
		ParentID:    nil,
	}
	repo.categories[romance.ID.Hex()] = romance

	return repo
}

func (m *MockCategoryRepository) GetByID(ctx context.Context, id string) (*bookstoreModel.Category, error) {
	return m.categories[id], nil
}

func (m *MockCategoryRepository) GetBySlug(ctx context.Context, slug string) (*bookstoreModel.Category, error) {
	for _, cat := range m.categories {
		if cat.Slug == slug {
			return cat, nil
		}
	}
	return nil, nil
}

func (m *MockCategoryRepository) GetRootCategories(ctx context.Context) ([]*bookstoreModel.Category, error) {
	result := make([]*bookstoreModel.Category, 0)
	for _, cat := range m.categories {
		if cat.ParentID == nil {
			result = append(result, cat)
		}
	}
	return result, nil
}

func (m *MockCategoryRepository) GetCategoryTree(ctx context.Context) ([]*bookstoreModel.CategoryTree, error) {
	return []*bookstoreModel.CategoryTree{}, nil
}

func (m *MockCategoryRepository) Create(ctx context.Context, category *bookstoreModel.Category) error {
	category.ID = primitive.NewObjectID()
	m.categories[category.ID.Hex()] = category
	return nil
}

func (m *MockCategoryRepository) Update(ctx context.Context, category *bookstoreModel.Category) error {
	m.categories[category.ID.Hex()] = category
	return nil
}

func (m *MockCategoryRepository) Delete(ctx context.Context, id string) error {
	delete(m.categories, id)
	return nil
}

func (m *MockCategoryRepository) GetChildren(ctx context.Context, parentID string) ([]*bookstoreModel.Category, error) {
	return []*bookstoreModel.Category{}, nil
}

// MockBannerRepository Mock横幅仓储
type MockBannerRepository struct {
	banners map[string]*bookstoreModel.Banner
	nextID  int
}

func NewMockBannerRepository() *MockBannerRepository {
	repo := &MockBannerRepository{
		banners: make(map[string]*bookstoreModel.Banner),
		nextID:  1,
	}

	// 添加测试横幅
	banner1 := &bookstoreModel.Banner{
		ID:        primitive.NewObjectID(),
		Title:     "新书推荐",
		ImageURL:  "https://example.com/banner1.jpg",
		LinkURL:   "https://example.com/book/1",
		Position:  "home",
		IsActive:  true,
		SortOrder: 1,
	}
	repo.banners[banner1.ID.Hex()] = banner1

	banner2 := &bookstoreModel.Banner{
		ID:        primitive.NewObjectID(),
		Title:     "热门榜单",
		ImageURL:  "https://example.com/banner2.jpg",
		LinkURL:   "https://example.com/ranking",
		Position:  "home",
		IsActive:  true,
		SortOrder: 2,
	}
	repo.banners[banner2.ID.Hex()] = banner2

	return repo
}

func (m *MockBannerRepository) GetActiveBanners(ctx context.Context, limit int) ([]*bookstoreModel.Banner, error) {
	result := make([]*bookstoreModel.Banner, 0)
	for _, banner := range m.banners {
		if banner.IsActive {
			result = append(result, banner)
			if len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

func (m *MockBannerRepository) IncrementClickCount(ctx context.Context, id string) error {
	return nil
}

func (m *MockBannerRepository) Create(ctx context.Context, banner *bookstoreModel.Banner) error {
	banner.ID = primitive.NewObjectID()
	m.banners[banner.ID.Hex()] = banner
	return nil
}

func (m *MockBannerRepository) Update(ctx context.Context, banner *bookstoreModel.Banner) error {
	m.banners[banner.ID.Hex()] = banner
	return nil
}

func (m *MockBannerRepository) Delete(ctx context.Context, id string) error {
	delete(m.banners, id)
	return nil
}

// MockRankingRepository Mock榜单仓储
type MockRankingRepository struct {
	rankings map[string][]*bookstoreModel.RankingItem
}

func NewMockRankingRepository() *MockRankingRepository {
	repo := &MockRankingRepository{
		rankings: make(map[string][]*bookstoreModel.RankingItem),
	}

	// 添加测试榜单
	dailyRanking := []*bookstoreModel.RankingItem{
		{
			BookID:    primitive.NewObjectID(),
			BookTitle: "玄幻霸主",
			Rank:      1,
			Score:     9800,
		},
		{
			BookID:    primitive.NewObjectID(),
			BookTitle: "都市传说",
			Rank:      2,
			Score:     9500,
		},
	}
	repo.rankings["daily"] = dailyRanking

	weeklyRanking := []*bookstoreModel.RankingItem{
		{
			BookID:    primitive.NewObjectID(),
			BookTitle: "修真之路",
			Rank:      1,
			Score:     8500,
		},
		{
			BookID:    primitive.NewObjectID(),
			BookTitle: "仙侠世界",
			Rank:      2,
			Score:     8200,
		},
	}
	repo.rankings["weekly"] = weeklyRanking

	return repo
}

func (m *MockRankingRepository) GetRealtimeRanking(ctx context.Context, limit int) ([]*bookstoreModel.RankingItem, error) {
	return m.rankings["daily"][:limit], nil
}

func (m *MockRankingRepository) GetWeeklyRanking(ctx context.Context, period string, limit int) ([]*bookstoreModel.RankingItem, error) {
	return m.rankings["weekly"][:limit], nil
}

func (m *MockRankingRepository) GetMonthlyRanking(ctx context.Context, period string, limit int) ([]*bookstoreModel.RankingItem, error) {
	return m.rankings["weekly"][:limit], nil
}

// ============ 测试辅助函数 ============

func setupTestBookstoreService() (*BookstoreServiceImpl, *MockBookRepository, *MockCategoryRepository, *MockBannerRepository, *MockRankingRepository) {
	bookRepo := NewMockBookRepository()
	categoryRepo := NewMockCategoryRepository()
	bannerRepo := NewMockBannerRepository()
	rankingRepo := NewMockRankingRepository()

	// 添加测试书籍
	book1 := &bookstoreModel.Book{
		Title:       "测试书籍1",
		Description: "这是一本测试书籍",
		AuthorID:    primitive.NewObjectID(),
		CategoryID:  primitive.NewObjectID(),
		Price:       0,
		IsFree:      true,
		Status:      bookstoreModel.BookStatusPublished,
		ViewCount:   100,
		LikeCount:   50,
	}
	bookRepo.Create(context.Background(), book1)

	book2 := &bookstoreModel.Book{
		Title:       "测试书籍2",
		Description: "这是另一本测试书籍",
		AuthorID:    primitive.NewObjectID(),
		CategoryID:  primitive.NewObjectID(),
		Price:       0,
		IsFree:      true,
		Status:      bookstoreModel.BookStatusPublished,
		ViewCount:   200,
		LikeCount:   100,
	}
	bookRepo.Create(context.Background(), book2)

	service := NewBookstoreService(bookRepo, categoryRepo, bannerRepo, rankingRepo).(*BookstoreServiceImpl)

	return service, bookRepo, categoryRepo, bannerRepo, rankingRepo
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && (s[:1] == substr || s[len(s)-1:] == substr || containsInner(s, substr))))
}

func containsInner(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ============ 测试用例 ============

// TestGetAllBooks 测试获取所有书籍
func TestGetAllBooks(t *testing.T) {
	service, _, _, _, _ := setupTestBookstoreService()
	ctx := context.Background()

	books, total, err := service.GetAllBooks(ctx, 1, 10)
	if err != nil {
		t.Fatalf("获取书籍列表失败: %v", err)
	}

	if len(books) != 2 {
		t.Errorf("书籍数量错误: 期望2本，实际%d本", len(books))
	}
	if total != 2 {
		t.Errorf("总数错误: 期望2，实际%d", total)
	}

	t.Logf("获取书籍列表成功: %d本，总计%d", len(books), total)
}

// TestGetBookByID 测试根据ID获取书籍
func TestGetBookByID(t *testing.T) {
	service, bookRepo, _, _, _ := setupTestBookstoreService()
	ctx := context.Background()

	// 获取第一本书的ID
	var bookID string
	for id, book := range bookRepo.books {
		bookID = id
		t.Logf("测试书籍ID: %s, 标题: %s", id, book.Title)
		break
	}

	book, err := service.GetBookByID(ctx, bookID)
	if err != nil {
		t.Fatalf("获取书籍失败: %v", err)
	}

	if book == nil {
		t.Fatal("书籍不应为空")
	}

	t.Logf("获取书籍成功: %s", book.Title)
}

// TestGetHotBooks 测试获取热门书籍
func TestGetHotBooks(t *testing.T) {
	service, _, _, _, _ := setupTestBookstoreService()
	ctx := context.Background()

	books, total, err := service.GetHotBooks(ctx, 1, 10)
	if err != nil {
		t.Fatalf("获取热门书籍失败: %v", err)
	}

	if len(books) != 2 {
		t.Errorf("热门书籍数量错误: %d", len(books))
	}
	if total != 2 {
		t.Errorf("总数错误: %d", total)
	}

	t.Logf("获取热门书籍成功: %d本", len(books))
}

// TestGetFreeBooks 测试获取免费书籍
func TestGetFreeBooks(t *testing.T) {
	service, _, _, _, _ := setupTestBookstoreService()
	ctx := context.Background()

	books, total, err := service.GetFreeBooks(ctx, 1, 10)
	if err != nil {
		t.Fatalf("获取免费书籍失败: %v", err)
	}

	t.Logf("获取免费书籍成功: %d本，总计%d", len(books), total)
}

// TestSearchBooks 测试搜索书籍
func TestSearchBooks(t *testing.T) {
	service, _, _, _, _ := setupTestBookstoreService()
	ctx := context.Background()

	// 搜索"测试"
	books, total, err := service.SearchBooks(ctx, "测试", 1, 10)
	if err != nil {
		t.Fatalf("搜索书籍失败: %v", err)
	}

	if len(books) != 2 {
		t.Errorf("搜索结果数量错误: %d", len(books))
	}
	if total != 2 {
		t.Errorf("总数错误: %d", total)
	}

	t.Logf("搜索书籍成功: %d本，总计%d", len(books), total)
}

// TestGetRootCategories 测试获取根分类
func TestGetRootCategories(t *testing.T) {
	service, _, categoryRepo, _, _ := setupTestBookstoreService()
	ctx := context.Background()

	categories, err := service.GetRootCategories(ctx)
	if err != nil {
		t.Fatalf("获取根分类失败: %v", err)
	}

	expectedCount := len(categoryRepo.categories)
	if len(categories) != expectedCount {
		t.Errorf("分类数量错误: 期望%d，实际%d", expectedCount, len(categories))
	}

	t.Logf("获取根分类成功: %d个", len(categories))
}

// TestGetCategoryTree 测试获取分类树
func TestGetCategoryTree(t *testing.T) {
	service, _, _, _, _ := setupTestBookstoreService()
	ctx := context.Background()

	tree, err := service.GetCategoryTree(ctx)
	if err != nil {
		t.Fatalf("获取分类树失败: %v", err)
	}

	t.Logf("获取分类树成功: %d个节点", len(tree))
}

// TestGetActiveBanners 测试获取活跃横幅
func TestGetActiveBanners(t *testing.T) {
	service, _, _, bannerRepo, _ := setupTestBookstoreService()
	ctx := context.Background()

	banners, err := service.GetActiveBanners(ctx, 10)
	if err != nil {
		t.Fatalf("获取横幅失败: %v", err)
	}

	expectedCount := 0
	for _, banner := range bannerRepo.banners {
		if banner.IsActive {
			expectedCount++
		}
	}

	if len(banners) != expectedCount {
		t.Errorf("横幅数量错误: 期望%d，实际%d", expectedCount, len(banners))
	}

	t.Logf("获取活跃横幅成功: %d个", len(banners))
}

// TestIncrementBannerClick 测试增加横幅点击次数
func TestIncrementBannerClick(t *testing.T) {
	service, _, _, bannerRepo, _ := setupTestBookstoreService()
	ctx := context.Background()

	// 获取第一个横幅ID
	var bannerID string
	for id := range bannerRepo.banners {
		bannerID = id
		break
	}

	err := service.IncrementBannerClick(ctx, bannerID)
	if err != nil {
		t.Fatalf("增加横幅点击次数失败: %v", err)
	}

	t.Logf("增加横幅点击次数成功")
}

// TestGetRealtimeRanking 测试获取实时榜单
func TestGetRealtimeRanking(t *testing.T) {
	service, _, _, _, rankingRepo := setupTestBookstoreService()
	ctx := context.Background()

	rankings, err := service.GetRealtimeRanking(ctx, 10)
	if err != nil {
		t.Fatalf("获取实时榜单失败: %v", err)
	}

	expectedCount := len(rankingRepo.rankings["daily"])
	if len(rankings) != expectedCount {
		t.Errorf("榜单数量错误: 期望%d，实际%d", expectedCount, len(rankings))
	}

	t.Logf("获取实时榜单成功: %d个", len(rankings))
}

// TestGetWeeklyRanking 测试获取周榜
func TestGetWeeklyRanking(t *testing.T) {
	service, _, _, _, _ := setupTestBookstoreService()
	ctx := context.Background()

	rankings, err := service.GetWeeklyRanking(ctx, "2024-01", 10)
	if err != nil {
		t.Fatalf("获取周榜失败: %v", err)
	}

	t.Logf("获取周榜成功: %d个", len(rankings))
}

// TestGetBookStats 测试获取书籍统计
func TestGetBookStats(t *testing.T) {
	service, bookRepo, _, _, _ := setupTestBookstoreService()
	ctx := context.Background()

	stats, err := service.GetBookStats(ctx)
	if err != nil {
		t.Fatalf("获取书籍统计失败: %v", err)
	}

	if stats == nil {
		t.Fatal("统计不应为空")
	}

	expectedCount := len(bookRepo.books)
	if stats.TotalBooks != int64(expectedCount) {
		t.Errorf("书籍总数错误: 期望%d，实际%d", expectedCount, stats.TotalBooks)
	}

	t.Logf("获取书籍统计成功: 总数%d", stats.TotalBooks)
}

// TestIncrementBookView 测试增加书籍浏览次数
func TestIncrementBookView(t *testing.T) {
	service, bookRepo, _, _, _ := setupTestBookstoreService()
	ctx := context.Background()

	// 获取第一本书的ID
	var bookID string
	for id := range bookRepo.books {
		bookID = id
		break
	}

	err := service.IncrementBookView(ctx, bookID)
	if err != nil {
		t.Fatalf("增加浏览次数失败: %v", err)
	}

	// 验证浏览次数增加
	book := bookRepo.books[bookID]
	if book.ViewCount != 101 { // 原始值100 + 1
		t.Errorf("浏览次数错误: 期望101，实际%d", book.ViewCount)
	}

	t.Logf("增加书籍浏览次数成功: %d", book.ViewCount)
}

// TestGetHomepageData 测试获取首页数据
func TestGetHomepageData(t *testing.T) {
	service, _, categoryRepo, bannerRepo, _ := setupTestBookstoreService()
	ctx := context.Background()

	data, err := service.GetHomepageData(ctx)
	if err != nil {
		t.Fatalf("获取首页数据失败: %v", err)
	}

	if data == nil {
		t.Fatal("首页数据不应为空")
	}

	// 验证横幅数据
	expectedBanners := 0
	for _, banner := range bannerRepo.banners {
		if banner.IsActive {
			expectedBanners++
		}
	}
	if len(data.Banners) != expectedBanners {
		t.Errorf("横幅数量错误: 期望%d，实际%d", expectedBanners, len(data.Banners))
	}

	// 验证分类数据
	expectedCategories := len(categoryRepo.categories)
	if len(data.Categories) != expectedCategories {
		t.Errorf("分类数量错误: 期望%d，实际%d", expectedCategories, len(data.Categories))
	}

	// 验证统计数据
	if data.Stats == nil {
		t.Error("统计数据不应为空")
	}

	t.Logf("获取首页数据成功: 横幅%d，分类%d，榜单%d",
		len(data.Banners), len(data.Categories), len(data.Rankings))
}

// TestGetRecommendedBooks 测试获取推荐书籍
func TestGetRecommendedBooks(t *testing.T) {
	service, _, _, _, _ := setupTestBookstoreService()
	ctx := context.Background()

	books, total, err := service.GetRecommendedBooks(ctx, 1, 10)
	if err != nil {
		t.Fatalf("获取推荐书籍失败: %v", err)
	}

	t.Logf("获取推荐书籍成功: %d本，总计%d", len(books), total)
}

// TestGetFeaturedBooks 测试获取精选书籍
func TestGetFeaturedBooks(t *testing.T) {
	service, _, _, _, _ := setupTestBookstoreService()
	ctx := context.Background()

	books, total, err := service.GetFeaturedBooks(ctx, 1, 10)
	if err != nil {
		t.Fatalf("获取精选书籍失败: %v", err)
	}

	t.Logf("获取精选书籍成功: %d本，总计%d", len(books), total)
}

// TestGetNewReleases 测试获取新书
func TestGetNewReleases(t *testing.T) {
	service, _, _, _, _ := setupTestBookstoreService()
	ctx := context.Background()

	books, total, err := service.GetNewReleases(ctx, 1, 10)
	if err != nil {
		t.Fatalf("获取新书失败: %v", err)
	}

	t.Logf("获取新书成功: %d本，总计%d", len(books), total)
}

// TestGetBooksByCategory 测试根据分类获取书籍
func TestGetBooksByCategory(t *testing.T) {
	service, _, categoryRepo, _, _ := setupTestBookstoreService()
	ctx := context.Background()

	// 获取第一个分类ID
	var categoryID string
	for id := range categoryRepo.categories {
		categoryID = id
		break
	}

	books, total, err := service.GetBooksByCategory(ctx, categoryID, 1, 10)
	if err != nil {
		t.Fatalf("根据分类获取书籍失败: %v", err)
	}

	t.Logf("根据分类获取书籍成功: %d本，总计%d", len(books), total)
}

// TestGetMonthlyRanking 测试获取月榜
func TestGetMonthlyRanking(t *testing.T) {
	service, _, _, _, _ := setupTestBookstoreService()
	ctx := context.Background()

	rankings, err := service.GetMonthlyRanking(ctx, "2024-01", 10)
	if err != nil {
		t.Fatalf("获取月榜失败: %v", err)
	}

	t.Logf("获取月榜成功: %d个", len(rankings))
}

// TestGetNewbieRanking 测试获取新书榜
func TestGetNewbieRanking(t *testing.T) {
	service, _, _, _, _ := setupTestBookstoreService()
	ctx := context.Background()

	rankings, err := service.GetNewbieRanking(ctx, "2024-01", 10)
	if err != nil {
		t.Fatalf("获取新书榜失败: %v", err)
	}

	t.Logf("获取新书榜成功: %d个", len(rankings))
}

// BenchmarkGetAllBooks 性能测试：获取所有书籍
func BenchmarkGetAllBooks(b *testing.B) {
	service, _, _, _, _ := setupTestBookstoreService()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = service.GetAllBooks(ctx, 1, 10)
	}
}

// BenchmarkSearchBooks 性能测试：搜索书籍
func BenchmarkSearchBooks(b *testing.B) {
	service, _, _, _, _ := setupTestBookstoreService()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = service.SearchBooks(ctx, "测试", 1, 10)
	}
}
