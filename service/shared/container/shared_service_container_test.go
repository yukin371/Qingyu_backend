package container

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============ 测试用例 ============
// Mock实现已移至 test_mocks.go 文件

// TestNewSharedServiceContainer 测试容器创建
func TestNewSharedServiceContainer(t *testing.T) {
	container := NewSharedServiceContainer()

	assert.NotNil(t, container, "容器不应为nil")
	assert.False(t, container.initialized, "新创建的容器应该是未初始化状态")
	assert.Nil(t, container.authService, "服务应该是nil")
	assert.Nil(t, container.walletService, "服务应该是nil")
	assert.Nil(t, container.recommendationService, "服务应该是nil")
	assert.Nil(t, container.messagingService, "服务应该是nil")
	assert.Nil(t, container.storageService, "服务应该是nil")
	assert.Nil(t, container.adminService, "服务应该是nil")
}

// TestSharedServiceContainer_ServiceRegistration 测试服务注册
func TestSharedServiceContainer_ServiceRegistration(t *testing.T) {
	container := NewSharedServiceContainer()

	// 创建Mock服务
	mockAuth := new(MockAuthService)
	mockWallet := new(MockWalletService)
	mockRecommendation := new(MockRecommendationService)
	mockMessaging := new(MockMessagingService)
	mockStorage := new(MockStorageService)
	mockAdmin := new(MockAdminService)

	// 注册服务
	container.SetAuthService(mockAuth)
	container.SetWalletService(mockWallet)
	container.SetRecommendationService(mockRecommendation)
	container.SetMessagingService(mockMessaging)
	container.SetStorageService(mockStorage)
	container.SetAdminService(mockAdmin)

	// 验证服务已注册
	assert.NotNil(t, container.authService, "Auth服务应该已注册")
	assert.NotNil(t, container.walletService, "Wallet服务应该已注册")
	assert.NotNil(t, container.recommendationService, "Recommendation服务应该已注册")
	assert.NotNil(t, container.messagingService, "Messaging服务应该已注册")
	assert.NotNil(t, container.storageService, "Storage服务应该已注册")
	assert.NotNil(t, container.adminService, "Admin服务应该已注册")
}

// TestSharedServiceContainer_GetServices 测试服务获取
func TestSharedServiceContainer_GetServices(t *testing.T) {
	container := NewSharedServiceContainer()

	// 创建Mock服务
	mockAuth := new(MockAuthService)
	mockWallet := new(MockWalletService)
	mockRecommendation := new(MockRecommendationService)
	mockMessaging := new(MockMessagingService)
	mockStorage := new(MockStorageService)
	mockAdmin := new(MockAdminService)

	// 注册服务
	container.SetAuthService(mockAuth)
	container.SetWalletService(mockWallet)
	container.SetRecommendationService(mockRecommendation)
	container.SetMessagingService(mockMessaging)
	container.SetStorageService(mockStorage)
	container.SetAdminService(mockAdmin)

	// 获取并验证服务
	assert.Equal(t, mockAuth, container.AuthService(), "应该返回正确的Auth服务")
	assert.Equal(t, mockWallet, container.WalletService(), "应该返回正确的Wallet服务")
	assert.Equal(t, mockRecommendation, container.RecommendationService(), "应该返回正确的Recommendation服务")
	assert.Equal(t, mockMessaging, container.MessagingService(), "应该返回正确的Messaging服务")
	assert.Equal(t, mockStorage, container.StorageService(), "应该返回正确的Storage服务")
	assert.Equal(t, mockAdmin, container.AdminService(), "应该返回正确的Admin服务")
}

// TestSharedServiceContainer_Initialize_Success 测试成功初始化
func TestSharedServiceContainer_Initialize_Success(t *testing.T) {
	container := NewSharedServiceContainer()
	ctx := context.Background()

	// 创建Mock服务并设置期望
	mockAuth := new(MockAuthService)
	mockWallet := new(MockWalletService)
	mockRecommendation := new(MockRecommendationService)
	mockMessaging := new(MockMessagingService)
	mockStorage := new(MockStorageService)
	mockAdmin := new(MockAdminService)

	mockAuth.On("Health", ctx).Return(nil)
	mockWallet.On("Health", ctx).Return(nil)
	mockRecommendation.On("Health", ctx).Return(nil)
	mockMessaging.On("Health", ctx).Return(nil)
	mockStorage.On("Health", ctx).Return(nil)
	mockAdmin.On("Health", ctx).Return(nil)

	// 注册所有服务
	container.SetAuthService(mockAuth)
	container.SetWalletService(mockWallet)
	container.SetRecommendationService(mockRecommendation)
	container.SetMessagingService(mockMessaging)
	container.SetStorageService(mockStorage)
	container.SetAdminService(mockAdmin)

	// 执行初始化
	err := container.Initialize(ctx)

	// 验证
	assert.NoError(t, err, "初始化应该成功")
	assert.True(t, container.IsInitialized(), "容器应该已初始化")

	// 验证所有Mock调用
	mockAuth.AssertExpectations(t)
	mockWallet.AssertExpectations(t)
	mockRecommendation.AssertExpectations(t)
	mockMessaging.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
	mockAdmin.AssertExpectations(t)
}

// TestSharedServiceContainer_Initialize_MissingService 测试缺少服务时初始化失败
func TestSharedServiceContainer_Initialize_MissingService(t *testing.T) {
	tests := []struct {
		name          string
		setupServices func(*SharedServiceContainer)
		expectedError string
	}{
		{
			name: "缺少Auth服务",
			setupServices: func(c *SharedServiceContainer) {
				// 不设置Auth服务
				c.SetWalletService(new(MockWalletService))
				c.SetRecommendationService(new(MockRecommendationService))
				c.SetMessagingService(new(MockMessagingService))
				c.SetStorageService(new(MockStorageService))
				c.SetAdminService(new(MockAdminService))
			},
			expectedError: "Auth服务未注册",
		},
		{
			name: "缺少Wallet服务",
			setupServices: func(c *SharedServiceContainer) {
				c.SetAuthService(new(MockAuthService))
				// 不设置Wallet服务
				c.SetRecommendationService(new(MockRecommendationService))
				c.SetMessagingService(new(MockMessagingService))
				c.SetStorageService(new(MockStorageService))
				c.SetAdminService(new(MockAdminService))
			},
			expectedError: "Wallet服务未注册",
		},
		{
			name: "缺少Recommendation服务",
			setupServices: func(c *SharedServiceContainer) {
				c.SetAuthService(new(MockAuthService))
				c.SetWalletService(new(MockWalletService))
				// 不设置Recommendation服务
				c.SetMessagingService(new(MockMessagingService))
				c.SetStorageService(new(MockStorageService))
				c.SetAdminService(new(MockAdminService))
			},
			expectedError: "Recommendation服务未注册",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := NewSharedServiceContainer()
			ctx := context.Background()

			// 设置服务
			tt.setupServices(container)

			// 执行初始化
			err := container.Initialize(ctx)

			// 验证
			assert.Error(t, err, "应该返回错误")
			assert.Contains(t, err.Error(), tt.expectedError, "错误信息应该包含预期内容")
			assert.False(t, container.IsInitialized(), "容器不应该被标记为已初始化")
		})
	}
}

// TestSharedServiceContainer_Initialize_HealthCheckFail 测试健康检查失败
func TestSharedServiceContainer_Initialize_HealthCheckFail(t *testing.T) {
	container := NewSharedServiceContainer()
	ctx := context.Background()

	// 创建Mock服务
	mockAuth := new(MockAuthService)
	mockWallet := new(MockWalletService)
	mockRecommendation := new(MockRecommendationService)
	mockMessaging := new(MockMessagingService)
	mockStorage := new(MockStorageService)
	mockAdmin := new(MockAdminService)

	// 设置健康检查 - Wallet服务失败
	// 注意：由于Map遍历顺序不确定，所有服务都可能被调用，所以都要设置mock
	mockAuth.On("Health", ctx).Return(nil).Maybe()
	mockWallet.On("Health", ctx).Return(errors.New("database connection failed")).Maybe()
	mockRecommendation.On("Health", ctx).Return(nil).Maybe()
	mockMessaging.On("Health", ctx).Return(nil).Maybe()
	mockStorage.On("Health", ctx).Return(nil).Maybe()
	mockAdmin.On("Health", ctx).Return(nil).Maybe()

	// 注册所有服务
	container.SetAuthService(mockAuth)
	container.SetWalletService(mockWallet)
	container.SetRecommendationService(mockRecommendation)
	container.SetMessagingService(mockMessaging)
	container.SetStorageService(mockStorage)
	container.SetAdminService(mockAdmin)

	// 执行初始化
	err := container.Initialize(ctx)

	// 验证
	assert.Error(t, err, "初始化应该失败")
	assert.Contains(t, err.Error(), "健康检查失败", "错误信息应该提示健康检查失败")
	assert.False(t, container.IsInitialized(), "容器不应该被标记为已初始化")
}

// TestSharedServiceContainer_Initialize_Idempotent 测试初始化幂等性
func TestSharedServiceContainer_Initialize_Idempotent(t *testing.T) {
	container := NewSharedServiceContainer()
	ctx := context.Background()

	// 创建Mock服务
	mockAuth := new(MockAuthService)
	mockWallet := new(MockWalletService)
	mockRecommendation := new(MockRecommendationService)
	mockMessaging := new(MockMessagingService)
	mockStorage := new(MockStorageService)
	mockAdmin := new(MockAdminService)

	// Health应该只被调用一次（第一次初始化）
	mockAuth.On("Health", ctx).Return(nil).Once()
	mockWallet.On("Health", ctx).Return(nil).Once()
	mockRecommendation.On("Health", ctx).Return(nil).Once()
	mockMessaging.On("Health", ctx).Return(nil).Once()
	mockStorage.On("Health", ctx).Return(nil).Once()
	mockAdmin.On("Health", ctx).Return(nil).Once()

	// 注册所有服务
	container.SetAuthService(mockAuth)
	container.SetWalletService(mockWallet)
	container.SetRecommendationService(mockRecommendation)
	container.SetMessagingService(mockMessaging)
	container.SetStorageService(mockStorage)
	container.SetAdminService(mockAdmin)

	// 第一次初始化
	err := container.Initialize(ctx)
	assert.NoError(t, err)
	assert.True(t, container.IsInitialized())

	// 第二次初始化（应该立即返回）
	err = container.Initialize(ctx)
	assert.NoError(t, err)
	assert.True(t, container.IsInitialized())

	// 第三次初始化
	err = container.Initialize(ctx)
	assert.NoError(t, err)

	// 验证Health只被调用了一次
	mockAuth.AssertExpectations(t)
	mockWallet.AssertExpectations(t)
	mockRecommendation.AssertExpectations(t)
	mockMessaging.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
	mockAdmin.AssertExpectations(t)
}

// TestSharedServiceContainer_Health 测试健康检查
func TestSharedServiceContainer_Health(t *testing.T) {
	t.Run("所有服务健康", func(t *testing.T) {
		container := NewSharedServiceContainer()
		ctx := context.Background()

		// 创建Mock服务
		mockAuth := new(MockAuthService)
		mockWallet := new(MockWalletService)
		mockRecommendation := new(MockRecommendationService)
		mockMessaging := new(MockMessagingService)
		mockStorage := new(MockStorageService)
		mockAdmin := new(MockAdminService)

		mockAuth.On("Health", ctx).Return(nil)
		mockWallet.On("Health", ctx).Return(nil)
		mockRecommendation.On("Health", ctx).Return(nil)
		mockMessaging.On("Health", ctx).Return(nil)
		mockStorage.On("Health", ctx).Return(nil)
		mockAdmin.On("Health", ctx).Return(nil)

		// 注册服务
		container.SetAuthService(mockAuth)
		container.SetWalletService(mockWallet)
		container.SetRecommendationService(mockRecommendation)
		container.SetMessagingService(mockMessaging)
		container.SetStorageService(mockStorage)
		container.SetAdminService(mockAdmin)

		// 执行健康检查
		err := container.Health(ctx)

		assert.NoError(t, err, "健康检查应该通过")
		mockAuth.AssertExpectations(t)
		mockWallet.AssertExpectations(t)
	})

	t.Run("某个服务不健康", func(t *testing.T) {
		container := NewSharedServiceContainer()
		ctx := context.Background()

		// 创建Mock服务
		mockAuth := new(MockAuthService)
		mockWallet := new(MockWalletService)
		mockRecommendation := new(MockRecommendationService)
		mockMessaging := new(MockMessagingService)
		mockStorage := new(MockStorageService)
		mockAdmin := new(MockAdminService)

		// 由于Map遍历顺序不确定，需要为所有服务设置mock
		mockAuth.On("Health", ctx).Return(nil).Maybe()
		mockWallet.On("Health", ctx).Return(errors.New("connection timeout")).Maybe()
		mockRecommendation.On("Health", ctx).Return(nil).Maybe()
		mockMessaging.On("Health", ctx).Return(nil).Maybe()
		mockStorage.On("Health", ctx).Return(nil).Maybe()
		mockAdmin.On("Health", ctx).Return(nil).Maybe()

		// 注册服务
		container.SetAuthService(mockAuth)
		container.SetWalletService(mockWallet)
		container.SetRecommendationService(mockRecommendation)
		container.SetMessagingService(mockMessaging)
		container.SetStorageService(mockStorage)
		container.SetAdminService(mockAdmin)

		// 执行健康检查
		err := container.Health(ctx)

		assert.Error(t, err, "健康检查应该失败")
		assert.Contains(t, err.Error(), "Wallet", "错误应该指明是哪个服务")
		assert.Contains(t, err.Error(), "connection timeout", "错误应该包含原因")
	})

	t.Run("服务未注册", func(t *testing.T) {
		container := NewSharedServiceContainer()
		ctx := context.Background()

		// 只注册部分服务
		mockAuth := new(MockAuthService)
		mockAuth.On("Health", mock.Anything).Return(nil)
		container.SetAuthService(mockAuth)

		// 执行健康检查
		err := container.Health(ctx)

		assert.Error(t, err, "健康检查应该失败")
		assert.Contains(t, err.Error(), "未初始化", "错误应该提示服务未注册")
	})
}

// TestSharedServiceContainer_GetServiceStatus 测试服务状态查询
func TestSharedServiceContainer_GetServiceStatus(t *testing.T) {
	container := NewSharedServiceContainer()
	ctx := context.Background()

	// 创建Mock服务
	mockAuth := new(MockAuthService)
	mockWallet := new(MockWalletService)
	mockRecommendation := new(MockRecommendationService)

	// 设置不同的健康状态
	mockAuth.On("Health", ctx).Return(nil)                                       // 健康
	mockWallet.On("Health", ctx).Return(errors.New("database connection error")) // 不健康
	mockRecommendation.On("Health", ctx).Return(nil)                             // 健康

	// 注册部分服务（Messaging/Storage/Admin未注册）
	container.SetAuthService(mockAuth)
	container.SetWalletService(mockWallet)
	container.SetRecommendationService(mockRecommendation)

	// 获取服务状态
	status := container.GetServiceStatus(ctx)

	// 验证状态
	assert.Equal(t, "healthy", status["Auth"], "Auth服务应该是健康的")
	assert.Contains(t, status["Wallet"], "unhealthy", "Wallet服务应该是不健康的")
	assert.Contains(t, status["Wallet"], "database connection error", "应该包含错误信息")
	assert.Equal(t, "healthy", status["Recommendation"], "Recommendation服务应该是健康的")
	assert.Equal(t, "not_registered", status["Messaging"], "Messaging服务应该是未注册")
	assert.Equal(t, "not_registered", status["Storage"], "Storage服务应该是未注册")
	assert.Equal(t, "not_registered", status["Admin"], "Admin服务应该是未注册")

	// 验证Mock调用
	mockAuth.AssertExpectations(t)
	mockWallet.AssertExpectations(t)
	mockRecommendation.AssertExpectations(t)
}

// TestSharedServiceContainer_IsInitialized 测试初始化状态检查
func TestSharedServiceContainer_IsInitialized(t *testing.T) {
	container := NewSharedServiceContainer()

	// 初始状态
	assert.False(t, container.IsInitialized(), "新容器应该是未初始化状态")

	// 手动设置initialized标志（仅用于测试）
	container.initialized = true
	assert.True(t, container.IsInitialized(), "设置后应该是已初始化状态")
}
