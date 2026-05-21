package user

import (
	"Qingyu_backend/models/shared"
	"Qingyu_backend/models/users"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserStatsUserRepository struct {
	mock.Mock
}

func (m *MockUserStatsUserRepository) GetByID(ctx context.Context, id string) (*users.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.User), args.Error(1)
}

type MockUserStatsProjectRepository struct {
	mock.Mock
}

func (m *MockUserStatsProjectRepository) CountByOwner(ctx context.Context, ownerID string) (int64, error) {
	args := m.Called(ctx, ownerID)
	return args.Get(0).(int64), args.Error(1)
}

func testStatsUser(roles []string, vipLevel int) *users.User {
	now := time.Now()
	return &users.User{
		BaseEntity: shared.BaseEntity{
			CreatedAt: now.Add(-48 * time.Hour),
			UpdatedAt: now,
		},
		Roles:    roles,
		VIPLevel: vipLevel,
	}
}

func TestUserStatsService_GetUserStats_EmptyUserID(t *testing.T) {
	service := NewUserStatsService(&MockUserStatsUserRepository{}, &MockUserStatsProjectRepository{})

	stats, err := service.GetUserStats(context.Background(), "")

	assert.Error(t, err)
	assert.Nil(t, stats)
	assert.Contains(t, err.Error(), "用户ID不能为空")
}

func TestUserStatsService_GetUserStats_Success(t *testing.T) {
	ctx := context.Background()
	userID := "user-123"
	mockUserRepo := new(MockUserStatsUserRepository)
	mockProjectRepo := new(MockUserStatsProjectRepository)
	service := NewUserStatsService(mockUserRepo, mockProjectRepo)

	mockUserRepo.On("GetByID", ctx, userID).Return(testStatsUser([]string{"author"}, 2), nil).Once()
	mockProjectRepo.On("CountByOwner", ctx, userID).Return(int64(5), nil).Once()

	stats, err := service.GetUserStats(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, userID, stats.UserID)
	assert.Equal(t, int64(5), stats.TotalProjects)
	assert.Equal(t, "作者 (VIP Level 2)", stats.MemberLevel)
	mockUserRepo.AssertExpectations(t)
	mockProjectRepo.AssertExpectations(t)
}

func TestUserStatsService_GetUserStats_ProjectCountErrorFallback(t *testing.T) {
	ctx := context.Background()
	userID := "user-456"
	mockUserRepo := new(MockUserStatsUserRepository)
	mockProjectRepo := new(MockUserStatsProjectRepository)
	service := NewUserStatsService(mockUserRepo, mockProjectRepo)

	mockUserRepo.On("GetByID", ctx, userID).Return(testStatsUser([]string{"reader"}, 0), nil).Once()
	mockProjectRepo.On("CountByOwner", ctx, userID).Return(int64(0), errors.New("db down")).Once()

	stats, err := service.GetUserStats(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, int64(0), stats.TotalProjects)
	mockUserRepo.AssertExpectations(t)
	mockProjectRepo.AssertExpectations(t)
}
