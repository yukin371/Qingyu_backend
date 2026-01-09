package shared

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockUserBehavior 模拟用户行为模型
type MockUserBehavior struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"userId"`
	ContentID primitive.ObjectID `bson:"content_id" json:"contentId"`
	Action    string             `bson:"action" json:"action"`
	Duration  int64              `bson:"duration" json:"duration"`
	Rating    float64            `bson:"rating" json:"rating"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
}

// MockContentFeature 模拟内容特征模型
type MockContentFeature struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ContentID primitive.ObjectID `bson:"content_id" json:"contentId"`
	Title     string             `bson:"title" json:"title"`
	Category  string             `bson:"category" json:"category"`
	Tags      []string           `bson:"tags" json:"tags"`
	Author    string             `bson:"author" json:"author"`
	WordCount int64              `bson:"word_count" json:"wordCount"`
	Rating    float64            `bson:"rating" json:"rating"`
	ViewCount int64              `bson:"view_count" json:"viewCount"`
	Features  map[string]float64 `bson:"features" json:"features"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}

// MockUserProfile 模拟用户画像模型
type MockUserProfile struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID `bson:"user_id" json:"userId"`
	Preferences     map[string]float64 `bson:"preferences" json:"preferences"`
	Categories      []string           `bson:"categories" json:"categories"`
	Tags            []string           `bson:"tags" json:"tags"`
	ReadingHabits   map[string]float64 `bson:"reading_habits" json:"readingHabits"`
	InteractionData map[string]int64   `bson:"interaction_data" json:"interactionData"`
	CreatedAt       time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updatedAt"`
}

// MockRecommendResult 模拟推荐结果模型
type MockRecommendResult struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID   `bson:"user_id" json:"userId"`
	ContentIDs []primitive.ObjectID `bson:"content_ids" json:"contentIds"`
	Scores     []float64            `bson:"scores" json:"scores"`
	Algorithm  string               `bson:"algorithm" json:"algorithm"`
	Context    map[string]string    `bson:"context" json:"context"`
	ExpiresAt  time.Time            `bson:"expires_at" json:"expiresAt"`
	CreatedAt  time.Time            `bson:"created_at" json:"createdAt"`
}

// MockBehaviorRepository 模拟用户行为仓储
type MockBehaviorRepository struct {
	mock.Mock
}

func (m *MockBehaviorRepository) Create(ctx context.Context, behavior *MockUserBehavior) error {
	args := m.Called(ctx, behavior)
	return args.Error(0)
}

func (m *MockBehaviorRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID, limit int) ([]*MockUserBehavior, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*MockUserBehavior), args.Error(1)
}

func (m *MockBehaviorRepository) GetSimilarUsers(ctx context.Context, userID primitive.ObjectID, limit int) ([]primitive.ObjectID, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]primitive.ObjectID), args.Error(1)
}

// MockContentRepository 模拟内容仓储
type MockContentRepository struct {
	mock.Mock
}

func (m *MockContentRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*MockContentFeature, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockContentFeature), args.Error(1)
}

func (m *MockContentRepository) GetByCategory(ctx context.Context, category string, limit int) ([]*MockContentFeature, error) {
	args := m.Called(ctx, category, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*MockContentFeature), args.Error(1)
}

func (m *MockContentRepository) GetSimilarContent(ctx context.Context, contentID primitive.ObjectID, limit int) ([]*MockContentFeature, error) {
	args := m.Called(ctx, contentID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*MockContentFeature), args.Error(1)
}

func (m *MockContentRepository) GetPopular(ctx context.Context, limit int) ([]*MockContentFeature, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*MockContentFeature), args.Error(1)
}

// MockUserProfileRepository 模拟用户画像仓储
type MockUserProfileRepository struct {
	mock.Mock
}

func (m *MockUserProfileRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID) (*MockUserProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockUserProfile), args.Error(1)
}

func (m *MockUserProfileRepository) UpdatePreferences(ctx context.Context, userID primitive.ObjectID, preferences map[string]float64) error {
	args := m.Called(ctx, userID, preferences)
	return args.Error(0)
}

// MockRecommendRepository 模拟推荐结果仓储
type MockRecommendRepository struct {
	mock.Mock
}

func (m *MockRecommendRepository) Save(ctx context.Context, result *MockRecommendResult) error {
	args := m.Called(ctx, result)
	return args.Error(0)
}

func (m *MockRecommendRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID) (*MockRecommendResult, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockRecommendResult), args.Error(1)
}

// MockRecommendationService 模拟推荐服务
type MockRecommendationService struct {
	behaviorRepo    *MockBehaviorRepository
	contentRepo     *MockContentRepository
	userProfileRepo *MockUserProfileRepository
	recommendRepo   *MockRecommendRepository
}

func NewMockRecommendationService(
	behaviorRepo *MockBehaviorRepository,
	contentRepo *MockContentRepository,
	userProfileRepo *MockUserProfileRepository,
	recommendRepo *MockRecommendRepository,
) *MockRecommendationService {
	return &MockRecommendationService{
		behaviorRepo:    behaviorRepo,
		contentRepo:     contentRepo,
		userProfileRepo: userProfileRepo,
		recommendRepo:   recommendRepo,
	}
}

// GetPersonalizedRecommendations 获取个性化推荐
func (s *MockRecommendationService) GetPersonalizedRecommendations(ctx context.Context, userID primitive.ObjectID, limit int) ([]*MockContentFeature, error) {
	// 获取用户画像
	userProfile, err := s.userProfileRepo.GetByUserID(ctx, userID)
	if err != nil {
		// 如果没有用户画像，返回热门推荐
		return s.GetPopularRecommendations(ctx, limit)
	}

	// 基于用户偏好推荐
	var recommendations []*MockContentFeature
	for _, category := range userProfile.Categories {
		contents, err := s.contentRepo.GetByCategory(ctx, category, limit/len(userProfile.Categories)+1)
		if err == nil {
			recommendations = append(recommendations, contents...)
		}
	}

	// 限制返回数量
	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	return recommendations, nil
}

// GetSimilarRecommendations 获取相似推荐
func (s *MockRecommendationService) GetSimilarRecommendations(ctx context.Context, contentID primitive.ObjectID, limit int) ([]*MockContentFeature, error) {
	return s.contentRepo.GetSimilarContent(ctx, contentID, limit)
}

// GetPopularRecommendations 获取热门推荐
func (s *MockRecommendationService) GetPopularRecommendations(ctx context.Context, limit int) ([]*MockContentFeature, error) {
	return s.contentRepo.GetPopular(ctx, limit)
}

// GetCollaborativeRecommendations 获取协同过滤推荐
func (s *MockRecommendationService) GetCollaborativeRecommendations(ctx context.Context, userID primitive.ObjectID, limit int) ([]*MockContentFeature, error) {
	// 获取相似用户
	similarUsers, err := s.behaviorRepo.GetSimilarUsers(ctx, userID, 10)
	if err != nil {
		return nil, err
	}

	if len(similarUsers) == 0 {
		return s.GetPopularRecommendations(ctx, limit)
	}

	// 获取相似用户的行为
	var allBehaviors []*MockUserBehavior
	for _, similarUserID := range similarUsers {
		behaviors, err := s.behaviorRepo.GetByUserID(ctx, similarUserID, 20)
		if err == nil {
			allBehaviors = append(allBehaviors, behaviors...)
		}
	}

	// 统计内容频次
	contentFreq := make(map[primitive.ObjectID]int)
	for _, behavior := range allBehaviors {
		if behavior.Action == "read" || behavior.Action == "like" {
			contentFreq[behavior.ContentID]++
		}
	}

	// 获取推荐内容
	var recommendations []*MockContentFeature
	count := 0
	for contentID := range contentFreq {
		if count >= limit {
			break
		}
		content, err := s.contentRepo.GetByID(ctx, contentID)
		if err == nil {
			recommendations = append(recommendations, content)
			count++
		}
	}

	return recommendations, nil
}

// RecordUserBehavior 记录用户行为
func (s *MockRecommendationService) RecordUserBehavior(ctx context.Context, userID, contentID primitive.ObjectID, action string, duration int64) error {
	behavior := &MockUserBehavior{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		ContentID: contentID,
		Action:    action,
		Duration:  duration,
		CreatedAt: time.Now(),
	}

	return s.behaviorRepo.Create(ctx, behavior)
}

// UpdateUserProfile 更新用户画像
func (s *MockRecommendationService) UpdateUserProfile(ctx context.Context, userID primitive.ObjectID) error {
	// 获取用户行为
	behaviors, err := s.behaviorRepo.GetByUserID(ctx, userID, 100)
	if err != nil {
		return err
	}

	// 分析用户偏好
	categoryCount := make(map[string]int)
	tagCount := make(map[string]int)

	for _, behavior := range behaviors {
		content, err := s.contentRepo.GetByID(ctx, behavior.ContentID)
		if err == nil {
			categoryCount[content.Category]++
			for _, tag := range content.Tags {
				tagCount[tag]++
			}
		}
	}

	// 构建偏好权重
	preferences := make(map[string]float64)
	for category, count := range categoryCount {
		preferences[category] = float64(count) / float64(len(behaviors))
	}

	// 更新用户画像
	return s.userProfileRepo.UpdatePreferences(ctx, userID, preferences)
}

// 测试用例

func TestRecommendationService_GetPersonalizedRecommendations_WithProfile(t *testing.T) {
	behaviorRepo := new(MockBehaviorRepository)
	contentRepo := new(MockContentRepository)
	userProfileRepo := new(MockUserProfileRepository)
	recommendRepo := new(MockRecommendRepository)
	service := NewMockRecommendationService(behaviorRepo, contentRepo, userProfileRepo, recommendRepo)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	limit := 5

	userProfile := &MockUserProfile{
		ID:         primitive.NewObjectID(),
		UserID:     userID,
		Categories: []string{"小说", "漫画"},
		Tags:       []string{"玄幻", "都市"},
	}

	expectedContents := []*MockContentFeature{
		{
			ID:       primitive.NewObjectID(),
			Title:    "测试小说1",
			Category: "小说",
			Tags:     []string{"玄幻"},
		},
		{
			ID:       primitive.NewObjectID(),
			Title:    "测试漫画1",
			Category: "漫画",
			Tags:     []string{"都市"},
		},
	}

	// Mock 设置
	userProfileRepo.On("GetByUserID", ctx, userID).Return(userProfile, nil)
	contentRepo.On("GetByCategory", ctx, "小说", 3).Return(expectedContents[:1], nil)
	contentRepo.On("GetByCategory", ctx, "漫画", 3).Return(expectedContents[1:], nil)

	// 执行测试
	recommendations, err := service.GetPersonalizedRecommendations(ctx, userID, limit)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, recommendations)
	assert.Len(t, recommendations, 2)
	assert.Equal(t, "测试小说1", recommendations[0].Title)
	assert.Equal(t, "测试漫画1", recommendations[1].Title)

	userProfileRepo.AssertExpectations(t)
	contentRepo.AssertExpectations(t)
}

func TestRecommendationService_GetPersonalizedRecommendations_NoProfile(t *testing.T) {
	behaviorRepo := new(MockBehaviorRepository)
	contentRepo := new(MockContentRepository)
	userProfileRepo := new(MockUserProfileRepository)
	recommendRepo := new(MockRecommendRepository)
	service := NewMockRecommendationService(behaviorRepo, contentRepo, userProfileRepo, recommendRepo)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	limit := 5

	expectedContents := []*MockContentFeature{
		{
			ID:        primitive.NewObjectID(),
			Title:     "热门小说1",
			Category:  "小说",
			ViewCount: 10000,
		},
		{
			ID:        primitive.NewObjectID(),
			Title:     "热门漫画1",
			Category:  "漫画",
			ViewCount: 8000,
		},
	}

	// Mock 设置
	userProfileRepo.On("GetByUserID", ctx, userID).Return(nil, errors.New("profile not found"))
	contentRepo.On("GetPopular", ctx, limit).Return(expectedContents, nil)

	// 执行测试
	recommendations, err := service.GetPersonalizedRecommendations(ctx, userID, limit)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, recommendations)
	assert.Len(t, recommendations, 2)
	assert.Equal(t, "热门小说1", recommendations[0].Title)
	assert.Equal(t, "热门漫画1", recommendations[1].Title)

	userProfileRepo.AssertExpectations(t)
	contentRepo.AssertExpectations(t)
}

func TestRecommendationService_GetSimilarRecommendations_Success(t *testing.T) {
	behaviorRepo := new(MockBehaviorRepository)
	contentRepo := new(MockContentRepository)
	userProfileRepo := new(MockUserProfileRepository)
	recommendRepo := new(MockRecommendRepository)
	service := NewMockRecommendationService(behaviorRepo, contentRepo, userProfileRepo, recommendRepo)

	ctx := context.Background()
	contentID := primitive.NewObjectID()
	limit := 5

	expectedContents := []*MockContentFeature{
		{
			ID:       primitive.NewObjectID(),
			Title:    "相似内容1",
			Category: "小说",
			Tags:     []string{"玄幻"},
		},
		{
			ID:       primitive.NewObjectID(),
			Title:    "相似内容2",
			Category: "小说",
			Tags:     []string{"玄幻"},
		},
	}

	// Mock 设置
	contentRepo.On("GetSimilarContent", ctx, contentID, limit).Return(expectedContents, nil)

	// 执行测试
	recommendations, err := service.GetSimilarRecommendations(ctx, contentID, limit)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, recommendations)
	assert.Len(t, recommendations, 2)
	assert.Equal(t, expectedContents, recommendations)

	contentRepo.AssertExpectations(t)
}

func TestRecommendationService_GetCollaborativeRecommendations_Success(t *testing.T) {
	behaviorRepo := new(MockBehaviorRepository)
	contentRepo := new(MockContentRepository)
	userProfileRepo := new(MockUserProfileRepository)
	recommendRepo := new(MockRecommendRepository)
	service := NewMockRecommendationService(behaviorRepo, contentRepo, userProfileRepo, recommendRepo)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	limit := 5

	similarUsers := []primitive.ObjectID{
		primitive.NewObjectID(),
		primitive.NewObjectID(),
	}

	contentID1 := primitive.NewObjectID()
	contentID2 := primitive.NewObjectID()

	behaviors := []*MockUserBehavior{
		{
			ID:        primitive.NewObjectID(),
			UserID:    similarUsers[0],
			ContentID: contentID1,
			Action:    "read",
		},
		{
			ID:        primitive.NewObjectID(),
			UserID:    similarUsers[1],
			ContentID: contentID2,
			Action:    "like",
		},
	}

	expectedContents := []*MockContentFeature{
		{
			ID:       contentID1,
			Title:    "协同推荐1",
			Category: "小说",
		},
		{
			ID:       contentID2,
			Title:    "协同推荐2",
			Category: "漫画",
		},
	}

	// Mock 设置
	behaviorRepo.On("GetSimilarUsers", ctx, userID, 10).Return(similarUsers, nil)
	behaviorRepo.On("GetByUserID", ctx, similarUsers[0], 20).Return([]*MockUserBehavior{behaviors[0]}, nil)
	behaviorRepo.On("GetByUserID", ctx, similarUsers[1], 20).Return([]*MockUserBehavior{behaviors[1]}, nil)
	contentRepo.On("GetByID", ctx, contentID1).Return(expectedContents[0], nil)
	contentRepo.On("GetByID", ctx, contentID2).Return(expectedContents[1], nil)

	// 执行测试
	recommendations, err := service.GetCollaborativeRecommendations(ctx, userID, limit)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, recommendations)
	assert.Len(t, recommendations, 2)

	behaviorRepo.AssertExpectations(t)
	contentRepo.AssertExpectations(t)
}

func TestRecommendationService_RecordUserBehavior_Success(t *testing.T) {
	behaviorRepo := new(MockBehaviorRepository)
	contentRepo := new(MockContentRepository)
	userProfileRepo := new(MockUserProfileRepository)
	recommendRepo := new(MockRecommendRepository)
	service := NewMockRecommendationService(behaviorRepo, contentRepo, userProfileRepo, recommendRepo)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	contentID := primitive.NewObjectID()
	action := "read"
	duration := int64(300)

	// Mock 设置
	behaviorRepo.On("Create", ctx, mock.AnythingOfType("*shared.MockUserBehavior")).Return(nil)

	// 执行测试
	err := service.RecordUserBehavior(ctx, userID, contentID, action, duration)

	// 断言
	assert.NoError(t, err)

	behaviorRepo.AssertExpectations(t)
}

func TestRecommendationService_UpdateUserProfile_Success(t *testing.T) {
	behaviorRepo := new(MockBehaviorRepository)
	contentRepo := new(MockContentRepository)
	userProfileRepo := new(MockUserProfileRepository)
	recommendRepo := new(MockRecommendRepository)
	service := NewMockRecommendationService(behaviorRepo, contentRepo, userProfileRepo, recommendRepo)

	ctx := context.Background()
	userID := primitive.NewObjectID()

	contentID1 := primitive.NewObjectID()
	contentID2 := primitive.NewObjectID()

	behaviors := []*MockUserBehavior{
		{
			ID:        primitive.NewObjectID(),
			UserID:    userID,
			ContentID: contentID1,
			Action:    "read",
		},
		{
			ID:        primitive.NewObjectID(),
			UserID:    userID,
			ContentID: contentID2,
			Action:    "like",
		},
	}

	contents := []*MockContentFeature{
		{
			ID:       contentID1,
			Category: "小说",
			Tags:     []string{"玄幻"},
		},
		{
			ID:       contentID2,
			Category: "小说",
			Tags:     []string{"都市"},
		},
	}

	expectedPreferences := map[string]float64{
		"小说": 1.0,
	}

	// Mock 设置
	behaviorRepo.On("GetByUserID", ctx, userID, 100).Return(behaviors, nil)
	contentRepo.On("GetByID", ctx, contentID1).Return(contents[0], nil)
	contentRepo.On("GetByID", ctx, contentID2).Return(contents[1], nil)
	userProfileRepo.On("UpdatePreferences", ctx, userID, expectedPreferences).Return(nil)

	// 执行测试
	err := service.UpdateUserProfile(ctx, userID)

	// 断言
	assert.NoError(t, err)

	behaviorRepo.AssertExpectations(t)
	contentRepo.AssertExpectations(t)
	userProfileRepo.AssertExpectations(t)
}
