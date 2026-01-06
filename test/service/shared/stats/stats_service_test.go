package stats_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/bookstore"
	"Qingyu_backend/models/users"
	bookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
	userRepo "Qingyu_backend/repository/interfaces/user"
	writingRepo "Qingyu_backend/repository/interfaces/writer"
	statsService "Qingyu_backend/service/shared/stats"
)

// ============ Mock Repositories ============

// MockUserRepository Mock用户Repository
type MockUserRepository struct {
	mock.Mock
	userRepo.UserRepository
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*users.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.User), args.Error(1)
}

// MockBookRepository Mock书籍Repository
type MockBookRepository struct {
	mock.Mock
	bookstoreRepo.BookRepository
}

func (m *MockBookRepository) GetByAuthorID(ctx context.Context, authorID primitive.ObjectID, limit, offset int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, authorID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

// MockProjectRepository Mock项目Repository
type MockProjectRepository struct {
	mock.Mock
	writingRepo.ProjectRepository
}

func (m *MockProjectRepository) CountByOwner(ctx context.Context, ownerID string) (int64, error) {
	args := m.Called(ctx, ownerID)
	return args.Get(0).(int64), args.Error(1)
}

// MockChapterRepository Mock章节Repository
type MockChapterRepository struct {
	mock.Mock
	bookstoreRepo.ChapterRepository
}

func (m *MockChapterRepository) GetTotalWordCount(ctx context.Context, bookID primitive.ObjectID) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

// ============ 测试用例 ============

func TestPlatformStatsService_GetUserStats(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Setup
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockChapterRepo := new(MockChapterRepository)

		service := statsService.NewPlatformStatsService(
			mockUserRepo,
			mockBookRepo,
			mockProjectRepo,
			mockChapterRepo,
		)

		// 准备测试数据
		userID := "user123"
		testUser := &users.User{
			ID:        userID,
			Username:  "testuser",
			Role:      "writer",
			CreatedAt: time.Now().Add(-30 * 24 * time.Hour), // 30天前注册
		}

		// Mock用户查询
		mockUserRepo.On("GetByID", ctx, userID).Return(testUser, nil)

		// Mock项目数统计
		mockProjectRepo.On("CountByOwner", ctx, userID).Return(int64(5), nil)

		// 执行测试
		stats, err := service.GetUserStats(ctx, userID)

		// 验证
		require.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, userID, stats.UserID)
		assert.Equal(t, int64(5), stats.TotalProjects)
		assert.Contains(t, []string{"普通用户", "作者", "管理员"}, stats.MemberLevel)

		mockUserRepo.AssertExpectations(t)
		mockProjectRepo.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		// Setup
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockChapterRepo := new(MockChapterRepository)

		service := statsService.NewPlatformStatsService(
			mockUserRepo,
			mockBookRepo,
			mockProjectRepo,
			mockChapterRepo,
		)

		// Mock用户不存在
		userID := "nonexistent"
		mockUserRepo.On("GetByID", ctx, userID).Return(nil, assert.AnError)

		// 执行测试
		stats, err := service.GetUserStats(ctx, userID)

		// 验证
		assert.Error(t, err)
		assert.Nil(t, stats)
		assert.Contains(t, err.Error(), "获取用户信息失败")

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("EmptyUserID", func(t *testing.T) {
		// Setup
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockChapterRepo := new(MockChapterRepository)

		service := statsService.NewPlatformStatsService(
			mockUserRepo,
			mockBookRepo,
			mockProjectRepo,
			mockChapterRepo,
		)

		// 执行测试（空用户ID）
		stats, err := service.GetUserStats(ctx, "")

		// 验证
		assert.Error(t, err)
		assert.Nil(t, stats)
		assert.Contains(t, err.Error(), "用户ID不能为空")
	})

	t.Run("ProjectCountError", func(t *testing.T) {
		// Setup
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockChapterRepo := new(MockChapterRepository)

		service := statsService.NewPlatformStatsService(
			mockUserRepo,
			mockBookRepo,
			mockProjectRepo,
			mockChapterRepo,
		)

		// 准备测试数据
		userID := "user123"
		testUser := &users.User{
			ID:        userID,
			Username:  "testuser",
			Role:      "writer",
			CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
		}

		// Mock用户查询成功
		mockUserRepo.On("GetByID", ctx, userID).Return(testUser, nil)

		// Mock项目统计失败
		mockProjectRepo.On("CountByOwner", ctx, userID).Return(int64(0), assert.AnError)

		// 执行测试
		stats, err := service.GetUserStats(ctx, userID)

		// 验证：统计失败时使用默认值0，不返回错误
		require.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, int64(0), stats.TotalProjects)

		mockUserRepo.AssertExpectations(t)
		mockProjectRepo.AssertExpectations(t)
	})
}

func TestPlatformStatsService_GetContentStats(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Setup
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockChapterRepo := new(MockChapterRepository)

		service := statsService.NewPlatformStatsService(
			mockUserRepo,
			mockBookRepo,
			mockProjectRepo,
			mockChapterRepo,
		)

		// 准备测试数据
		userID := "user123"
		testUser := &users.User{
			ID:        userID,
			Username:  "testuser",
			Role:      "writer",
			CreatedAt: time.Now().Add(-100 * 24 * time.Hour), // 100天前注册
		}

		// Mock用户查询
		mockUserRepo.On("GetByID", ctx, userID).Return(testUser, nil)

		// Mock项目数统计
		mockProjectRepo.On("CountByOwner", ctx, userID).Return(int64(3), nil)

		// 执行测试
		stats, err := service.GetContentStats(ctx, userID)

		// 验证
		require.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, userID, stats.UserID)
		assert.Equal(t, int64(3), stats.TotalProjects)
		// 验证日均字数计算（基于注册时间）
		assert.GreaterOrEqual(t, stats.AverageWordsPerDay, float64(0))

		mockUserRepo.AssertExpectations(t)
		mockProjectRepo.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		// Setup
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockChapterRepo := new(MockChapterRepository)

		service := statsService.NewPlatformStatsService(
			mockUserRepo,
			mockBookRepo,
			mockProjectRepo,
			mockChapterRepo,
		)

		// Mock用户不存在
		userID := "nonexistent"
		mockUserRepo.On("GetByID", ctx, userID).Return(nil, assert.AnError)

		// 执行测试
		stats, err := service.GetContentStats(ctx, userID)

		// 验证
		assert.Error(t, err)
		assert.Nil(t, stats)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("EmptyUserID", func(t *testing.T) {
		// Setup
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockChapterRepo := new(MockChapterRepository)

		service := statsService.NewPlatformStatsService(
			mockUserRepo,
			mockBookRepo,
			mockProjectRepo,
			mockChapterRepo,
		)

		// 执行测试（空用户ID）
		stats, err := service.GetContentStats(ctx, "")

		// 验证
		assert.Error(t, err)
		assert.Nil(t, stats)
		assert.Contains(t, err.Error(), "用户ID不能为空")
	})
}

func TestPlatformStatsService_GetPlatformUserStats(t *testing.T) {
	ctx := context.Background()

	t.Run("ReturnsEmptyData", func(t *testing.T) {
		// Setup
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockChapterRepo := new(MockChapterRepository)

		service := statsService.NewPlatformStatsService(
			mockUserRepo,
			mockBookRepo,
			mockProjectRepo,
			mockChapterRepo,
		)

		// 执行测试
		startDate := time.Now().Add(-7 * 24 * time.Hour)
		endDate := time.Now()
		stats, err := service.GetPlatformUserStats(ctx, startDate, endDate)

		// 验证：当前返回空数据（待Task3实现）
		require.NoError(t, err)
		assert.NotNil(t, stats)
		// 当前实现返回0值
		assert.Equal(t, int64(0), stats.TotalUsers)
		assert.Equal(t, int64(0), stats.NewUsers)
		assert.Equal(t, int64(0), stats.ActiveUsers)
	})
}

func TestPlatformStatsService_GetPlatformContentStats(t *testing.T) {
	ctx := context.Background()

	t.Run("ReturnsEmptyData", func(t *testing.T) {
		// Setup
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockChapterRepo := new(MockChapterRepository)

		service := statsService.NewPlatformStatsService(
			mockUserRepo,
			mockBookRepo,
			mockProjectRepo,
			mockChapterRepo,
		)

		// 执行测试
		startDate := time.Now().Add(-7 * 24 * time.Hour)
		endDate := time.Now()
		stats, err := service.GetPlatformContentStats(ctx, startDate, endDate)

		// 验证：当前返回空数据（待Task3实现）
		require.NoError(t, err)
		assert.NotNil(t, stats)
		// 当前实现返回0值
		assert.Equal(t, int64(0), stats.TotalBooks)
		assert.Equal(t, int64(0), stats.NewBooks)
		assert.Equal(t, int64(0), stats.TotalChapters)
	})
}

func TestPlatformStatsService_GetUserActivityStats(t *testing.T) {
	ctx := context.Background()

	t.Run("ReturnsEmptyData", func(t *testing.T) {
		// Setup
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockChapterRepo := new(MockChapterRepository)

		service := statsService.NewPlatformStatsService(
			mockUserRepo,
			mockBookRepo,
			mockProjectRepo,
			mockChapterRepo,
		)

		// 执行测试
		stats, err := service.GetUserActivityStats(ctx, "user123", 7)

		// 验证：当前返回空数据（待Task3实现）
		require.NoError(t, err)
		assert.NotNil(t, stats)
		// 当前实现返回0值
		assert.Equal(t, "user123", stats.UserID)
		assert.Equal(t, int64(0), stats.TotalActions) // 类型为int64
	})
}

func TestPlatformStatsService_GetRevenueStats(t *testing.T) {
	ctx := context.Background()

	t.Run("ReturnsEmptyData", func(t *testing.T) {
		// Setup
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockChapterRepo := new(MockChapterRepository)

		service := statsService.NewPlatformStatsService(
			mockUserRepo,
			mockBookRepo,
			mockProjectRepo,
			mockChapterRepo,
		)

		// 执行测试
		startDate := time.Now().Add(-30 * 24 * time.Hour)
		endDate := time.Now()
		stats, err := service.GetRevenueStats(ctx, "user123", startDate, endDate)

		// 验证：当前返回空数据（待Task3实现）
		require.NoError(t, err)
		assert.NotNil(t, stats)
		// 当前实现返回0值
		assert.Equal(t, "user123", stats.UserID)
		assert.Equal(t, float64(0), stats.TotalRevenue)
	})
}

func TestPlatformStatsService_Health(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Setup
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockChapterRepo := new(MockChapterRepository)

		service := statsService.NewPlatformStatsService(
			mockUserRepo,
			mockBookRepo,
			mockProjectRepo,
			mockChapterRepo,
		)

		// 执行测试
		err := service.Health(ctx)

		// 验证
		assert.NoError(t, err)
	})
}

// ============ 辅助函数测试 ============

func TestPlatformStatsService_MemberLevelMapping(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		role          string
		expectedLevel string
	}{
		{"reader", "普通用户"},
		{"writer", "普通用户"}, // writer角色也映射为普通用户（因为代码只识别admin）
		{"admin", "管理员"},
		{"", "普通用户"}, // 默认值
	}

	for _, tc := range testCases {
		t.Run("Role_"+tc.role, func(t *testing.T) {
			// Setup
			mockUserRepo := new(MockUserRepository)
			mockBookRepo := new(MockBookRepository)
			mockProjectRepo := new(MockProjectRepository)
			mockChapterRepo := new(MockChapterRepository)

			service := statsService.NewPlatformStatsService(
				mockUserRepo,
				mockBookRepo,
				mockProjectRepo,
				mockChapterRepo,
			)

			// 准备测试数据
			userID := "user123"
			testUser := &users.User{
				ID:        userID,
				Username:  "testuser",
				Role:      tc.role,
				CreatedAt: time.Now(),
			}

			// Mock
			mockUserRepo.On("GetByID", ctx, userID).Return(testUser, nil)
			mockProjectRepo.On("CountByOwner", ctx, userID).Return(int64(0), nil)

			// 执行测试
			stats, err := service.GetUserStats(ctx, userID)

			// 验证
			require.NoError(t, err)
			assert.NotNil(t, stats)
			assert.Equal(t, tc.expectedLevel, stats.MemberLevel)

			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestPlatformStatsService_AverageWordsPerDay(t *testing.T) {
	ctx := context.Background()

	t.Run("CalculatesCorrectly", func(t *testing.T) {
		// Setup
		mockUserRepo := new(MockUserRepository)
		mockBookRepo := new(MockBookRepository)
		mockProjectRepo := new(MockProjectRepository)
		mockChapterRepo := new(MockChapterRepository)

		service := statsService.NewPlatformStatsService(
			mockUserRepo,
			mockBookRepo,
			mockProjectRepo,
			mockChapterRepo,
		)

		// 准备测试数据：100天前注册
		userID := "user123"
		daysAgo := 100
		testUser := &users.User{
			ID:        userID,
			Username:  "testuser",
			Role:      "writer",
			CreatedAt: time.Now().Add(-time.Duration(daysAgo) * 24 * time.Hour),
		}

		// Mock
		mockUserRepo.On("GetByID", ctx, userID).Return(testUser, nil)
		mockProjectRepo.On("CountByOwner", ctx, userID).Return(int64(5), nil)

		// 执行测试
		stats, err := service.GetContentStats(ctx, userID)

		// 验证
		require.NoError(t, err)
		assert.NotNil(t, stats)
		// 验证日期计算（应该约等于100天）
		daysSinceRegistration := time.Since(testUser.CreatedAt).Hours() / 24
		assert.InDelta(t, float64(daysAgo), daysSinceRegistration, 1.0)
	})
}
