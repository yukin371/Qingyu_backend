package writer

import (
	"context"
	"errors"
	"testing"

	"Qingyu_backend/models/users"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockContentStatsUserRepository struct {
	mock.Mock
}

func (m *MockContentStatsUserRepository) GetByID(ctx context.Context, id string) (*users.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.User), args.Error(1)
}

type MockContentStatsProjectRepository struct {
	mock.Mock
}

func (m *MockContentStatsProjectRepository) CountByOwner(ctx context.Context, ownerID string) (int64, error) {
	args := m.Called(ctx, ownerID)
	return args.Get(0).(int64), args.Error(1)
}

func TestContentStatsService_GetContentStats_EmptyUserID(t *testing.T) {
	service := NewContentStatsService(nil, nil)

	stats, err := service.GetContentStats(context.Background(), "")

	assert.Error(t, err)
	assert.Nil(t, stats)
	assert.Contains(t, err.Error(), "用户ID不能为空")
}

func TestContentStatsService_GetContentStats_UserNotFound(t *testing.T) {
	ctx := context.Background()
	userID := "user-404"
	mockUserRepo := new(MockContentStatsUserRepository)
	mockProjectRepo := new(MockContentStatsProjectRepository)
	service := NewContentStatsService(mockUserRepo, mockProjectRepo)

	mockUserRepo.On("GetByID", ctx, userID).Return(nil, errors.New("用户不存在")).Once()

	stats, err := service.GetContentStats(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, stats)
	assert.Contains(t, err.Error(), "获取用户信息失败")
	mockUserRepo.AssertExpectations(t)
}

func TestContentStatsService_GetContentStats_Success(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"
	mockUserRepo := new(MockContentStatsUserRepository)
	mockProjectRepo := new(MockContentStatsProjectRepository)
	service := NewContentStatsService(mockUserRepo, mockProjectRepo)

	mockUserRepo.On("GetByID", ctx, userID).Return(&users.User{}, nil).Once()
	mockProjectRepo.On("CountByOwner", ctx, userID).Return(int64(10), nil).Once()

	stats, err := service.GetContentStats(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, userID, stats.UserID)
	assert.Equal(t, int64(10), stats.TotalProjects)
	assert.Equal(t, int64(0), stats.PublishedBooks)
	assert.Equal(t, int64(0), stats.TotalWords)
	mockUserRepo.AssertExpectations(t)
	mockProjectRepo.AssertExpectations(t)
}

func TestContentStatsService_GetContentStats_SuccessWithoutProjectRepo(t *testing.T) {
	ctx := context.Background()
	userID := "user-789"
	mockUserRepo := new(MockContentStatsUserRepository)
	service := NewContentStatsService(mockUserRepo, nil)

	mockUserRepo.On("GetByID", ctx, userID).Return(&users.User{}, nil).Once()

	stats, err := service.GetContentStats(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, userID, stats.UserID)
	assert.Equal(t, int64(0), stats.TotalProjects)
	assert.Equal(t, int64(0), stats.PublishedBooks)
	assert.Equal(t, int64(0), stats.TotalWords)
	mockUserRepo.AssertExpectations(t)
}

func TestContentStatsService_GetContentStats_ProjectCountErrorFallback(t *testing.T) {
	ctx := context.Background()
	userID := "user-456"
	mockUserRepo := new(MockContentStatsUserRepository)
	mockProjectRepo := new(MockContentStatsProjectRepository)
	service := NewContentStatsService(mockUserRepo, mockProjectRepo)

	mockUserRepo.On("GetByID", ctx, userID).Return(&users.User{}, nil).Once()
	mockProjectRepo.On("CountByOwner", ctx, userID).Return(int64(0), errors.New("db down")).Once()

	stats, err := service.GetContentStats(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, int64(0), stats.TotalProjects)
	mockUserRepo.AssertExpectations(t)
	mockProjectRepo.AssertExpectations(t)
}
