package bookstore

import (
	"context"
	"errors"
	"testing"
	"time"

	bookstoreModel "Qingyu_backend/models/bookstore"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ============ Mock 实现 ============

// MockRatingRepository Mock评分仓储
type MockRatingRepository struct {
	ratings    map[string]*bookstoreModel.BookRating
	ratingKeys map[string]string // key: "bookID_userID" -> ratingID
	nextID     int
}

func NewMockRatingRepository() *MockRatingRepository {
	return &MockRatingRepository{
		ratings:    make(map[string]*bookstoreModel.BookRating),
		ratingKeys: make(map[string]string),
		nextID:     1,
	}
}

func (m *MockRatingRepository) Create(ctx context.Context, rating *bookstoreModel.BookRating) error {
	rating.ID = primitive.NewObjectID()
	rating.CreatedAt = time.Now()
	rating.UpdatedAt = time.Now()
	m.ratings[rating.ID.Hex()] = rating

	// 创建复合键
	key := rating.BookID.Hex() + "_" + rating.UserID.Hex()
	m.ratingKeys[key] = rating.ID.Hex()

	return nil
}

func (m *MockRatingRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstoreModel.BookRating, error) {
	return m.ratings[id.Hex()], nil
}

func (m *MockRatingRepository) GetByBookIDAndUserID(ctx context.Context, bookID, userID primitive.ObjectID) (*bookstoreModel.BookRating, error) {
	key := bookID.Hex() + "_" + userID.Hex()
	ratingID, ok := m.ratingKeys[key]
	if !ok {
		return nil, nil
	}
	return m.ratings[ratingID], nil
}

func (m *MockRatingRepository) GetByBookID(ctx context.Context, bookID primitive.ObjectID, page, pageSize int) ([]*bookstoreModel.BookRating, int64, error) {
	result := make([]*bookstoreModel.BookRating, 0)
	for _, rating := range m.ratings {
		if rating.BookID == bookID {
			result = append(result, rating)
		}
	}
	return result, int64(len(result)), nil
}

func (m *MockRatingRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID, page, pageSize int) ([]*bookstoreModel.BookRating, int64, error) {
	result := make([]*bookstoreModel.BookRating, 0)
	for _, rating := range m.ratings {
		if rating.UserID == userID {
			result = append(result, rating)
		}
	}
	return result, int64(len(result)), nil
}

func (m *MockRatingRepository) Update(ctx context.Context, rating *bookstoreModel.BookRating) error {
	rating.UpdatedAt = time.Now()
	m.ratings[rating.ID.Hex()] = rating
	return nil
}

func (m *MockRatingRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	rating := m.ratings[id.Hex()]
	if rating != nil {
		key := rating.BookID.Hex() + "_" + rating.UserID.Hex()
		delete(m.ratingKeys, key)
	}
	delete(m.ratings, id.Hex())
	return nil
}

func (m *MockRatingRepository) GetAverageRating(ctx context.Context, bookID primitive.ObjectID) (float64, error) {
	sum := 0.0
	count := 0
	for _, rating := range m.ratings {
		if rating.BookID == bookID {
			sum += rating.Rating
			count++
		}
	}
	if count == 0 {
		return 0, nil
	}
	return sum / float64(count), nil
}

func (m *MockRatingRepository) GetRatingCount(ctx context.Context, bookID primitive.ObjectID) (int64, error) {
	count := 0
	for _, rating := range m.ratings {
		if rating.BookID == bookID {
			count++
		}
	}
	return int64(count), nil
}

func (m *MockRatingRepository) GetTopRatedBooks(ctx context.Context, limit int) ([]*bookstoreModel.BookRating, error) {
	result := make([]*bookstoreModel.BookRating, 0)
	for _, rating := range m.ratings {
		result = append(result, rating)
		if len(result) >= limit {
			break
		}
	}
	return result, nil
}

// MockCacheService Mock缓存服务
type MockCacheService struct {
	data map[string]interface{}
}

func NewMockCacheService() *MockCacheService {
	return &MockCacheService{
		data: make(map[string]interface{}),
	}
}

func (m *MockCacheService) Get(ctx context.Context, key string) (interface{}, error) {
	return m.data[key], nil
}

func (m *MockCacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	m.data[key] = value
	return nil
}

func (m *MockCacheService) Delete(ctx context.Context, key string) error {
	delete(m.data, key)
	return nil
}

func (m *MockCacheService) Exists(ctx context.Context, key string) (bool, error) {
	_, ok := m.data[key]
	return ok, nil
}

func (m *MockCacheService) Invalidate(ctx context.Context, pattern string) error {
	for key := range m.data {
		if len(key) >= len(pattern) && key[:len(pattern)] == pattern {
			delete(m.data, key)
		}
	}
	return nil
}

// ============ 测试辅助函数 ============

func setupTestRatingService() (*BookRatingServiceImpl, *MockRatingRepository, *MockCacheService) {
	ratingRepo := NewMockRatingRepository()
	cacheService := NewMockCacheService()

	service := NewBookRatingService(ratingRepo, cacheService).(*BookRatingServiceImpl)

	return service, ratingRepo, cacheService
}

func createTestRating(bookID, userID primitive.ObjectID, rating float64, comment string) *bookstoreModel.BookRating {
	return &bookstoreModel.BookRating{
		BookID:     bookID,
		UserID:     userID,
		Rating:     rating,
		Comment:    comment,
		Tags:       []string{"精彩", "好看"},
		IsHelpful:  true,
		HelpfulCount: 0,
	}
}

// ============ 测试用例 ============

// TestCreateRating 测试创建评分
func TestCreateRating(t *testing.T) {
	service, ratingRepo, _ := setupTestRatingService()
	ctx := context.Background()

	bookID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	rating := createTestRating(bookID, userID, 4.5, "这本书非常好！")

	err := service.CreateRating(ctx, rating)
	if err != nil {
		t.Fatalf("创建评分失败: %v", err)
	}

	// 验证评分已创建
	if len(ratingRepo.ratings) != 1 {
		t.Errorf("评分数量错误: 期望1个，实际%d个", len(ratingRepo.ratings))
	}

	t.Logf("创建评分成功: 书籍ID=%s, 用户ID=%s, 评分=%.1f",
		bookID.Hex(), userID.Hex(), rating.Rating)
}

// TestCreateRating_InvalidRating 测试创建无效评分
func TestCreateRating_InvalidRating(t *testing.T) {
	service, _, _ := setupTestRatingService()
	ctx := context.Background()

	bookID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	rating := createTestRating(bookID, userID, 6.0, "评分超出范围") // 无效评分

	err := service.CreateRating(ctx, rating)
	if err == nil {
		t.Fatal("应该拒绝无效评分，但成功了")
	}

	t.Logf("正确拒绝了无效评分: %v", err)
}

// TestCreateRating_Duplicate 测试重复评分
func TestCreateRating_Duplicate(t *testing.T) {
	service, _, _ := setupTestRatingService()
	ctx := context.Background()

	bookID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	// 第一次评分
	rating1 := createTestRating(bookID, userID, 4.0, "第一次评分")
	err := service.CreateRating(ctx, rating1)
	if err != nil {
		t.Fatalf("创建评分失败: %v", err)
	}

	// 第二次评分（应该失败）
	rating2 := createTestRating(bookID, userID, 5.0, "第二次评分")
	err = service.CreateRating(ctx, rating2)
	if err == nil {
		t.Fatal("应该拒绝重复评分，但成功了")
	}

	t.Logf("正确拒绝了重复评分: %v", err)
}

// TestGetRatingByID 测试根据ID获取评分
func TestGetRatingByID(t *testing.T) {
	service, ratingRepo, _ := setupTestRatingService()
	ctx := context.Background()

	bookID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	rating := createTestRating(bookID, userID, 4.5, "这本书非常好！")
	_ = service.CreateRating(ctx, rating)

	// 获取评分ID
	var ratingID string
	for id := range ratingRepo.ratings {
		ratingID = id
		break
	}

	// 查询评分
	found, err := service.GetRatingByID(ctx, primitive.NewObjectID())
	if err != nil {
		t.Fatalf("获取评分失败: %v", err)
	}

	if found == nil {
		t.Fatal("评分不应为空")
	}

	t.Logf("获取评分成功: ID=%s, 评分=%.1f", ratingID, found.Rating)
}

// TestGetRatingsByBookID 测试根据书籍ID获取评分
func TestGetRatingsByBookID(t *testing.T) {
	service, _, _ := setupTestRatingService()
	ctx := context.Background()

	bookID := primitive.NewObjectID()

	// 创建多个评分
	for i := 1; i <= 3; i++ {
		userID := primitive.NewObjectID()
		rating := createTestRating(bookID, userID, float64(i), "评分")
		_ = service.CreateRating(ctx, rating)
	}

	// 查询评分
	ratings, total, err := service.GetRatingsByBookID(ctx, bookID, 1, 10)
	if err != nil {
		t.Fatalf("获取评分失败: %v", err)
	}

	if len(ratings) != 3 {
		t.Errorf("评分数量错误: 期望3个，实际%d个", len(ratings))
	}
	if total != 3 {
		t.Errorf("总数错误: 期望3，实际%d", total)
	}

	t.Logf("根据书籍ID获取评分成功: %d个，总计%d", len(ratings), total)
}

// TestGetRatingsByUserID 测试根据用户ID获取评分
func TestGetRatingsByUserID(t *testing.T) {
	service, _, _ := setupTestRatingService()
	ctx := context.Background()

	userID := primitive.NewObjectID()

	// 创建多个评分
	for i := 1; i <= 2; i++ {
		bookID := primitive.NewObjectID()
		rating := createTestRating(bookID, userID, float64(i), "评分")
		_ = service.CreateRating(ctx, rating)
	}

	// 查询评分
	ratings, total, err := service.GetRatingsByUserID(ctx, userID, 1, 10)
	if err != nil {
		t.Fatalf("获取评分失败: %v", err)
	}

	if len(ratings) != 2 {
		t.Errorf("评分数量错误: 期望2个，实际%d个", len(ratings))
	}
	if total != 2 {
		t.Errorf("总数错误: 期望2，实际%d", total)
	}

	t.Logf("根据用户ID获取评分成功: %d个，总计%d", len(ratings), total)
}

// TestGetRatingByBookIDAndUserID 测试根据书籍ID和用户ID获取评分
func TestGetRatingByBookIDAndUserID(t *testing.T) {
	service, _, _ := setupTestRatingService()
	ctx := context.Background()

	bookID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	rating := createTestRating(bookID, userID, 4.5, "这本书非常好！")
	_ = service.CreateRating(ctx, rating)

	// 查询评分
	found, err := service.GetRatingByBookIDAndUserID(ctx, bookID, userID)
	if err != nil {
		t.Fatalf("获取评分失败: %v", err)
	}

	if found == nil {
		t.Fatal("评分不应为空")
	}
	if found.Rating != 4.5 {
		t.Errorf("评分错误: 期望4.5，实际%.1f", found.Rating)
	}

	t.Logf("根据书籍ID和用户ID获取评分成功: 评分=%.1f", found.Rating)
}

// TestUpdateRating 测试更新评分
func TestUpdateRating(t *testing.T) {
	service, ratingRepo, _ := setupTestRatingService()
	ctx := context.Background()

	bookID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	rating := createTestRating(bookID, userID, 4.0, "原始评分")
	_ = service.CreateRating(ctx, rating)

	// 获取评分ID
	var ratingID primitive.ObjectID
	for id, r := range ratingRepo.ratings {
		if r.BookID == bookID && r.UserID == userID {
			ratingID, _ = primitive.ObjectIDFromHex(id)
			break
		}
	}

	// 更新评分
	rating.Rating = 5.0
	rating.Comment = "更新后的评分"
	err := service.UpdateRating(ctx, rating)
	if err != nil {
		t.Fatalf("更新评分失败: %v", err)
	}

	// 验证更新
	updated, _ := service.GetRatingByID(ctx, ratingID)
	if updated.Rating != 5.0 {
		t.Errorf("评分未更新: 期望5.0，实际%.1f", updated.Rating)
	}

	t.Logf("更新评分成功: 新评分=%.1f", updated.Rating)
}

// TestDeleteRating 测试删除评分
func TestDeleteRating(t *testing.T) {
	service, ratingRepo, _ := setupTestRatingService()
	ctx := context.Background()

	bookID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	rating := createTestRating(bookID, userID, 4.5, "要删除的评分")
	_ = service.CreateRating(ctx, rating)

	// 获取评分ID
	var ratingID primitive.ObjectID
	for id := range ratingRepo.ratings {
		ratingID, _ = primitive.ObjectIDFromHex(id)
		break
	}

	// 删除评分
	err := service.DeleteRating(ctx, ratingID)
	if err != nil {
		t.Fatalf("删除评分失败: %v", err)
	}

	// 验证删除
	if len(ratingRepo.ratings) != 0 {
		t.Errorf("评分未删除: 期望0个，实际%d个", len(ratingRepo.ratings))
	}

	t.Logf("删除评分成功")
}

// TestGetAverageRating 测试获取平均评分
func TestGetAverageRating(t *testing.T) {
	service, _, _ := setupTestRatingService()
	ctx := context.Background()

	bookID := primitive.NewObjectID()

	// 创建多个评分
	ratings := []float64{4.0, 5.0, 3.0}
	for _, r := range ratings {
		userID := primitive.NewObjectID()
		rating := createTestRating(bookID, userID, r, "评分")
		_ = service.CreateRating(ctx, rating)
	}

	// 获取平均分
	avg, err := service.GetAverageRating(ctx, bookID)
	if err != nil {
		t.Fatalf("获取平均分失败: %v", err)
	}

	expected := 4.0
	if avg != expected {
		t.Errorf("平均分错误: 期望%.1f，实际%.1f", expected, avg)
	}

	t.Logf("获取平均分成功: %.1f", avg)
}

// TestGetRatingCount 测试获取评分数量
func TestGetRatingCount(t *testing.T) {
	service, _, _ := setupTestRatingService()
	ctx := context.Background()

	bookID := primitive.NewObjectID()

	// 创建多个评分
	for i := 1; i <= 5; i++ {
		userID := primitive.NewObjectID()
		rating := createTestRating(bookID, userID, float64(i), "评分")
		_ = service.CreateRating(ctx, rating)
	}

	// 获取评分数量
	count, err := service.GetRatingCount(ctx, bookID)
	if err != nil {
		t.Fatalf("获取评分数量失败: %v", err)
	}

	if count != 5 {
		t.Errorf("评分数量错误: 期望5，实际%d", count)
	}

	t.Logf("获取评分数量成功: %d", count)
}

// TestGetRatingDistribution 测试获取评分分布
func TestGetRatingDistribution(t *testing.T) {
	service, _, _ := setupTestRatingService()
	ctx := context.Background()

	bookID := primitive.NewObjectID()

	// 创建不同评分
	ratings := []float64{5.0, 5.0, 4.0, 3.0, 2.0, 1.0}
	for _, r := range ratings {
		userID := primitive.NewObjectID()
		rating := createTestRating(bookID, userID, r, "评分")
		_ = service.CreateRating(ctx, rating)
	}

	// 获取评分分布
	distribution, err := service.GetRatingDistribution(ctx, bookID)
	if err != nil {
		t.Fatalf("获取评分分布失败: %v", err)
	}

	// 验证分布
	if distribution[5] != 2 { // 两个5分
		t.Errorf("5分数量错误: 期望2，实际%d", distribution[5])
	}
	if distribution[4] != 1 { // 一个4分
		t.Errorf("4分数量错误: 期望1，实际%d", distribution[4])
	}

	t.Logf("获取评分分布成功: %+v", distribution)
}

// TestHasUserRated 测试检查用户是否已评分
func TestHasUserRated(t *testing.T) {
	service, _, _ := setupTestRatingService()
	ctx := context.Background()

	bookID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	// 检查未评分
	hasRated, err := service.HasUserRated(ctx, bookID, userID)
	if err != nil {
		t.Fatalf("检查评分失败: %v", err)
	}
	if hasRated {
		t.Error("用户应该未评分")
	}

	// 创建评分
	rating := createTestRating(bookID, userID, 4.5, "评分")
	_ = service.CreateRating(ctx, rating)

	// 检查已评分
	hasRated, err = service.HasUserRated(ctx, bookID, userID)
	if err != nil {
		t.Fatalf("检查评分失败: %v", err)
	}
	if !hasRated {
		t.Error("用户应该已评分")
	}

	t.Logf("检查用户评分状态成功")
}

// TestUpdateUserRating 测试更新用户评分
func TestUpdateUserRating(t *testing.T) {
	service, _, _ := setupTestRatingService()
	ctx := context.Background()

	bookID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	// 创建评分
	rating := createTestRating(bookID, userID, 4.0, "原始评分")
	_ = service.CreateRating(ctx, rating)

	// 更新评分
	err := service.UpdateUserRating(ctx, bookID, userID, 5.0, "更新后的评分", []string{"精彩"})
	if err != nil {
		t.Fatalf("更新用户评分失败: %v", err)
	}

	// 验证更新
	updated, _ := service.GetRatingByBookIDAndUserID(ctx, bookID, userID)
	if updated.Rating != 5.0 {
		t.Errorf("评分未更新: 期望5.0，实际%.1f", updated.Rating)
	}

	t.Logf("更新用户评分成功: 新评分=%.1f", updated.Rating)
}

// TestDeleteUserRating 测试删除用户评分
func TestDeleteUserRating(t *testing.T) {
	service, ratingRepo, _ := setupTestRatingService()
	ctx := context.Background()

	bookID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	// 创建评分
	rating := createTestRating(bookID, userID, 4.5, "要删除的评分")
	_ = service.CreateRating(ctx, rating)

	// 删除评分
	err := service.DeleteUserRating(ctx, bookID, userID)
	if err != nil {
		t.Fatalf("删除用户评分失败: %v", err)
	}

	// 验证删除
	if len(ratingRepo.ratings) != 0 {
		t.Errorf("评分未删除: 期望0个，实际%d个", len(ratingRepo.ratings))
	}

	t.Logf("删除用户评分成功")
}

// TestGetTopRatedBooks 测试获取高分书籍
func TestGetTopRatedBooks(t *testing.T) {
	service, _, _ := setupTestRatingService()
	ctx := context.Background()

	// 创建多个高分评分
	for i := 0; i < 5; i++ {
		bookID := primitive.NewObjectID()
		userID := primitive.NewObjectID()
		rating := createTestRating(bookID, userID, 5.0, "高分评分")
		_ = service.CreateRating(ctx, rating)
	}

	// 获取高分书籍
	topRated, err := service.GetTopRatedBooks(ctx, 10)
	if err != nil {
		t.Fatalf("获取高分书籍失败: %v", err)
	}

	if len(topRated) != 5 {
		t.Errorf("高分书籍数量错误: 期望5本，实际%d本", len(topRated))
	}

	t.Logf("获取高分书籍成功: %d本", len(topRated))
}

// TestGetRatingStats 测试获取评分统计
func TestGetRatingStats(t *testing.T) {
	service, _, _ := setupTestRatingService()
	ctx := context.Background()

	bookID := primitive.NewObjectID()

	// 创建多个评分
	ratings := []float64{5.0, 4.0, 4.0, 3.0, 2.0}
	for _, r := range ratings {
		userID := primitive.NewObjectID()
		rating := createTestRating(bookID, userID, r, "评分")
		_ = service.CreateRating(ctx, rating)
	}

	// 获取统计
	stats, err := service.GetRatingStats(ctx, bookID)
	if err != nil {
		t.Fatalf("获取评分统计失败: %v", err)
	}

	if stats == nil {
		t.Fatal("统计不应为空")
	}

	t.Logf("获取评分统计成功: %+v", stats)
}

// BenchmarkGetAverageRating 性能测试：获取平均评分
func BenchmarkGetAverageRating(b *testing.B) {
	service, _, _ := setupTestRatingService()
	ctx := context.Background()

	bookID := primitive.NewObjectID()

	// 创建多个评分
	for i := 0; i < 100; i++ {
		userID := primitive.NewObjectID()
		rating := createTestRating(bookID, userID, 4.5, "评分")
		_ = service.CreateRating(ctx, rating)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetAverageRating(ctx, bookID)
	}
}
