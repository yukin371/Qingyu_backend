package container

import (
	"context"
	"testing"

	"Qingyu_backend/repository/interfaces"
	"Qingyu_backend/service/container"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepositoryFactory Repository工厂的Mock实现
type MockRepositoryFactory struct {
	mock.Mock
}

func (m *MockRepositoryFactory) CreateUserRepository() interfaces.UserRepository {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(interfaces.UserRepository)
}

func (m *MockRepositoryFactory) CreateBookRepository() interfaces.BookRepository {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(interfaces.BookRepository)
}

func (m *MockRepositoryFactory) CreateCategoryRepository() interfaces.CategoryRepository {
	args := m.Called()
	return nil
}

func (m *MockRepositoryFactory) CreateBannerRepository() interfaces.BannerRepository {
	args := m.Called()
	return nil
}

func (m *MockRepositoryFactory) CreateChapterRepository() interfaces.ChapterRepository {
	args := m.Called()
	return nil
}

func (m *MockRepositoryFactory) CreateReadingProgressRepository() interfaces.ReadingProgressRepository {
	args := m.Called()
	return nil
}

func (m *MockRepositoryFactory) CreateAnnotationRepository() interfaces.AnnotationRepository {
	args := m.Called()
	return nil
}

func (m *MockRepositoryFactory) CreateReadingSettingsRepository() interfaces.ReadingSettingsRepository {
	args := m.Called()
	return nil
}

func (m *MockRepositoryFactory) CreateQuotaRepository() interfaces.QuotaRepository {
	args := m.Called()
	return nil
}

func (m *MockRepositoryFactory) CreateRankingRepository() interfaces.RankingRepository {
	args := m.Called()
	return nil
}

func (m *MockRepositoryFactory) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockRepositoryFactory) Close() error {
	args := m.Called()
	return args.Error(0)
}

// TestNewServiceContainer 测试服务容器创建
func TestNewServiceContainer(t *testing.T) {
	mockFactory := new(MockRepositoryFactory)

	serviceContainer := container.NewServiceContainer(mockFactory)

	assert.NotNil(t, serviceContainer, "服务容器不应为nil")
	assert.False(t, serviceContainer.IsInitialized(), "新创建的容器应该是未初始化状态")
}

// TestServiceContainer_RegisterService 测试服务注册
func TestServiceContainer_RegisterService(t *testing.T) {
	mockFactory := new(MockRepositoryFactory)
	serviceContainer := container.NewServiceContainer(mockFactory)

	// 创建一个Mock服务
	mockService := new(MockBaseService)
	mockService.On("GetServiceName").Return("TestService")
	mockService.On("GetVersion").Return("v1.0.0")

	// 注册服务
	err := serviceContainer.RegisterService("TestService", mockService)
	assert.NoError(t, err, "注册服务不应该失败")

	// 获取服务
	retrievedService, err := serviceContainer.GetService("TestService")
	assert.NoError(t, err, "获取服务不应该失败")
	assert.NotNil(t, retrievedService, "获取的服务不应为nil")

	// 尝试注册相同名称的服务应该失败
	err = serviceContainer.RegisterService("TestService", mockService)
	assert.Error(t, err, "注册重复服务应该失败")
}

// TestServiceContainer_GetService_NotFound 测试获取不存在的服务
func TestServiceContainer_GetService_NotFound(t *testing.T) {
	mockFactory := new(MockRepositoryFactory)
	serviceContainer := container.NewServiceContainer(mockFactory)

	// 获取不存在的服务
	service, err := serviceContainer.GetService("NonExistentService")
	assert.Error(t, err, "获取不存在的服务应该返回错误")
	assert.Nil(t, service, "获取不存在的服务应该返回nil")
}

// TestServiceContainer_GetServiceMetrics 测试获取服务指标
func TestServiceContainer_GetServiceMetrics(t *testing.T) {
	mockFactory := new(MockRepositoryFactory)
	serviceContainer := container.NewServiceContainer(mockFactory)

	// 创建一个Mock服务
	mockService := new(MockBaseService)
	mockService.On("GetServiceName").Return("TestService")
	mockService.On("GetVersion").Return("v1.0.0")

	// 注册服务
	err := serviceContainer.RegisterService("TestService", mockService)
	assert.NoError(t, err)

	// 获取服务指标
	metrics, err := serviceContainer.GetServiceMetrics("TestService")
	assert.NoError(t, err, "获取服务指标不应该失败")
	assert.NotNil(t, metrics, "服务指标不应为nil")
	assert.Equal(t, "TestService", metrics.ServiceName, "服务名称应该匹配")
	assert.Equal(t, "v1.0.0", metrics.Version, "服务版本应该匹配")
}

// TestServiceContainer_GetServiceNames 测试获取所有服务名称
func TestServiceContainer_GetServiceNames(t *testing.T) {
	mockFactory := new(MockRepositoryFactory)
	serviceContainer := container.NewServiceContainer(mockFactory)

	// 初始时应该没有服务
	names := serviceContainer.GetServiceNames()
	assert.Empty(t, names, "初始时服务列表应该为空")

	// 注册几个服务
	mockService1 := new(MockBaseService)
	mockService1.On("GetServiceName").Return("Service1")
	mockService1.On("GetVersion").Return("v1.0.0")

	mockService2 := new(MockBaseService)
	mockService2.On("GetServiceName").Return("Service2")
	mockService2.On("GetVersion").Return("v1.0.0")

	serviceContainer.RegisterService("Service1", mockService1)
	serviceContainer.RegisterService("Service2", mockService2)

	// 获取所有服务名称
	names = serviceContainer.GetServiceNames()
	assert.Len(t, names, 2, "应该有2个服务")
	assert.Contains(t, names, "Service1", "应该包含Service1")
	assert.Contains(t, names, "Service2", "应该包含Service2")
}

// MockBaseService BaseService的Mock实现
type MockBaseService struct {
	mock.Mock
}

func (m *MockBaseService) Initialize(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockBaseService) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockBaseService) Close(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockBaseService) GetServiceName() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockBaseService) GetVersion() string {
	args := m.Called()
	return args.String(0)
}
