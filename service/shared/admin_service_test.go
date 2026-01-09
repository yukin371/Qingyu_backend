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

// MockAdminUser 模拟管理员用户模型
type MockAdminUser struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username         string             `bson:"username" json:"username"`
	Email            string             `bson:"email" json:"email"`
	Phone            string             `bson:"phone" json:"phone"`
	RealName         string             `bson:"real_name" json:"realName"`
	Avatar           string             `bson:"avatar" json:"avatar"`
	Department       string             `bson:"department" json:"department"`
	Position         string             `bson:"position" json:"position"`
	Status           string             `bson:"status" json:"status"`
	Roles            []string           `bson:"roles" json:"roles"`
	Permissions      []string           `bson:"permissions" json:"permissions"`
	LastLoginTime    time.Time          `bson:"last_login_time" json:"lastLoginTime"`
	LastLoginIP      string             `bson:"last_login_ip" json:"lastLoginIP"`
	LoginCount       int64              `bson:"login_count" json:"loginCount"`
	PasswordHash     string             `bson:"password_hash" json:"-"`
	Salt             string             `bson:"salt" json:"-"`
	TwoFactorSecret  string             `bson:"two_factor_secret" json:"-"`
	TwoFactorEnabled bool               `bson:"two_factor_enabled" json:"twoFactorEnabled"`
	CreatedBy        string             `bson:"created_by" json:"createdBy"`
	CreatedAt        time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updatedAt"`
}

// MockOperationLog 模拟操作日志模型
type MockOperationLog struct {
	ID           primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	AdminUserID  primitive.ObjectID     `bson:"admin_user_id" json:"adminUserId"`
	Username     string                 `bson:"username" json:"username"`
	Module       string                 `bson:"module" json:"module"`
	Action       string                 `bson:"action" json:"action"`
	Resource     string                 `bson:"resource" json:"resource"`
	ResourceID   string                 `bson:"resource_id" json:"resourceId"`
	Method       string                 `bson:"method" json:"method"`
	URL          string                 `bson:"url" json:"url"`
	IPAddress    string                 `bson:"ip_address" json:"ipAddress"`
	UserAgent    string                 `bson:"user_agent" json:"userAgent"`
	RequestData  map[string]interface{} `bson:"request_data" json:"requestData"`
	ResponseData map[string]interface{} `bson:"response_data" json:"responseData"`
	Status       string                 `bson:"status" json:"status"`
	ErrorMessage string                 `bson:"error_message" json:"errorMessage"`
	Duration     int64                  `bson:"duration" json:"duration"`
	CreatedAt    time.Time              `bson:"created_at" json:"createdAt"`
}

// MockSystemConfig 模拟系统配置模型
type MockSystemConfig struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Key          string             `bson:"key" json:"key"`
	Value        string             `bson:"value" json:"value"`
	Type         string             `bson:"type" json:"type"`
	Category     string             `bson:"category" json:"category"`
	Name         string             `bson:"name" json:"name"`
	Description  string             `bson:"description" json:"description"`
	IsPublic     bool               `bson:"is_public" json:"isPublic"`
	IsEditable   bool               `bson:"is_editable" json:"isEditable"`
	DefaultValue string             `bson:"default_value" json:"defaultValue"`
	Validation   string             `bson:"validation" json:"validation"`
	UpdatedBy    string             `bson:"updated_by" json:"updatedBy"`
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updatedAt"`
}

// MockDataStats 模拟数据统计模型
type MockDataStats struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type       string             `bson:"type" json:"type"`
	Date       string             `bson:"date" json:"date"`
	Metrics    map[string]int64   `bson:"metrics" json:"metrics"`
	Dimensions map[string]string  `bson:"dimensions" json:"dimensions"`
	CreatedAt  time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updatedAt"`
}

// MockAdminUserRepository 模拟管理员用户仓储
type MockAdminUserRepository struct {
	mock.Mock
}

func (m *MockAdminUserRepository) Create(ctx context.Context, user *MockAdminUser) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockAdminUserRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*MockAdminUser, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockAdminUser), args.Error(1)
}

func (m *MockAdminUserRepository) GetByUsername(ctx context.Context, username string) (*MockAdminUser, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockAdminUser), args.Error(1)
}

func (m *MockAdminUserRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockAdminUserRepository) List(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]*MockAdminUser, error) {
	args := m.Called(ctx, filter, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*MockAdminUser), args.Error(1)
}

func (m *MockAdminUserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockOperationLogRepository 模拟操作日志仓储
type MockOperationLogRepository struct {
	mock.Mock
}

func (m *MockOperationLogRepository) Create(ctx context.Context, log *MockOperationLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockOperationLogRepository) List(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]*MockOperationLog, error) {
	args := m.Called(ctx, filter, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*MockOperationLog), args.Error(1)
}

func (m *MockOperationLogRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID, limit, offset int) ([]*MockOperationLog, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*MockOperationLog), args.Error(1)
}

// MockSystemConfigRepository 模拟系统配置仓储
type MockSystemConfigRepository struct {
	mock.Mock
}

func (m *MockSystemConfigRepository) Create(ctx context.Context, config *MockSystemConfig) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *MockSystemConfigRepository) GetByKey(ctx context.Context, key string) (*MockSystemConfig, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockSystemConfig), args.Error(1)
}

func (m *MockSystemConfigRepository) Update(ctx context.Context, key string, value string, updatedBy string) error {
	args := m.Called(ctx, key, value, updatedBy)
	return args.Error(0)
}

func (m *MockSystemConfigRepository) List(ctx context.Context, category string) ([]*MockSystemConfig, error) {
	args := m.Called(ctx, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*MockSystemConfig), args.Error(1)
}

// MockDataStatsRepository 模拟数据统计仓储
type MockDataStatsRepository struct {
	mock.Mock
}

func (m *MockDataStatsRepository) Create(ctx context.Context, stats *MockDataStats) error {
	args := m.Called(ctx, stats)
	return args.Error(0)
}

func (m *MockDataStatsRepository) GetByTypeAndDate(ctx context.Context, statsType, date string) (*MockDataStats, error) {
	args := m.Called(ctx, statsType, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockDataStats), args.Error(1)
}

func (m *MockDataStatsRepository) GetByTypeAndDateRange(ctx context.Context, statsType, startDate, endDate string) ([]*MockDataStats, error) {
	args := m.Called(ctx, statsType, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*MockDataStats), args.Error(1)
}

// MockAdminService 模拟管理后台服务
type MockAdminService struct {
	adminUserRepo    *MockAdminUserRepository
	operationLogRepo *MockOperationLogRepository
	systemConfigRepo *MockSystemConfigRepository
	dataStatsRepo    *MockDataStatsRepository
}

func NewMockAdminService(
	adminUserRepo *MockAdminUserRepository,
	operationLogRepo *MockOperationLogRepository,
	systemConfigRepo *MockSystemConfigRepository,
	dataStatsRepo *MockDataStatsRepository,
) *MockAdminService {
	return &MockAdminService{
		adminUserRepo:    adminUserRepo,
		operationLogRepo: operationLogRepo,
		systemConfigRepo: systemConfigRepo,
		dataStatsRepo:    dataStatsRepo,
	}
}

// CreateAdminUser 创建管理员用户
func (s *MockAdminService) CreateAdminUser(ctx context.Context, username, email, password, realName string, roles []string, createdBy string) (*MockAdminUser, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}

	if email == "" {
		return nil, errors.New("email is required")
	}

	// 检查用户名是否已存在
	existingUser, _ := s.adminUserRepo.GetByUsername(ctx, username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// 创建管理员用户
	adminUser := &MockAdminUser{
		ID:               primitive.NewObjectID(),
		Username:         username,
		Email:            email,
		RealName:         realName,
		Status:           "active",
		Roles:            roles,
		Permissions:      []string{},
		LoginCount:       0,
		PasswordHash:     hashPassword(password),
		Salt:             generateSalt(),
		TwoFactorEnabled: false,
		CreatedBy:        createdBy,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := s.adminUserRepo.Create(ctx, adminUser)
	if err != nil {
		return nil, err
	}

	return adminUser, nil
}

// GetAdminUser 获取管理员用户
func (s *MockAdminService) GetAdminUser(ctx context.Context, userID primitive.ObjectID) (*MockAdminUser, error) {
	return s.adminUserRepo.GetByID(ctx, userID)
}

// UpdateAdminUser 更新管理员用户
func (s *MockAdminService) UpdateAdminUser(ctx context.Context, userID primitive.ObjectID, updates map[string]interface{}) error {
	// 检查用户是否存在
	_, err := s.adminUserRepo.GetByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	updates["updated_at"] = time.Now()
	return s.adminUserRepo.Update(ctx, userID, updates)
}

// ListAdminUsers 获取管理员用户列表
func (s *MockAdminService) ListAdminUsers(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]*MockAdminUser, error) {
	return s.adminUserRepo.List(ctx, filter, limit, offset)
}

// DeleteAdminUser 删除管理员用户
func (s *MockAdminService) DeleteAdminUser(ctx context.Context, userID primitive.ObjectID) error {
	// 检查用户是否存在
	_, err := s.adminUserRepo.GetByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	return s.adminUserRepo.Delete(ctx, userID)
}

// LogOperation 记录操作日志
func (s *MockAdminService) LogOperation(ctx context.Context, adminUserID primitive.ObjectID, username, module, action, resource, resourceID, method, url, ipAddress, userAgent string, requestData, responseData map[string]interface{}, status string, errorMessage string, duration int64) error {
	log := &MockOperationLog{
		ID:           primitive.NewObjectID(),
		AdminUserID:  adminUserID,
		Username:     username,
		Module:       module,
		Action:       action,
		Resource:     resource,
		ResourceID:   resourceID,
		Method:       method,
		URL:          url,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		RequestData:  requestData,
		ResponseData: responseData,
		Status:       status,
		ErrorMessage: errorMessage,
		Duration:     duration,
		CreatedAt:    time.Now(),
	}

	return s.operationLogRepo.Create(ctx, log)
}

// GetOperationLogs 获取操作日志
func (s *MockAdminService) GetOperationLogs(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]*MockOperationLog, error) {
	return s.operationLogRepo.List(ctx, filter, limit, offset)
}

// GetUserOperationLogs 获取用户操作日志
func (s *MockAdminService) GetUserOperationLogs(ctx context.Context, userID primitive.ObjectID, limit, offset int) ([]*MockOperationLog, error) {
	return s.operationLogRepo.GetByUserID(ctx, userID, limit, offset)
}

// GetSystemConfig 获取系统配置
func (s *MockAdminService) GetSystemConfig(ctx context.Context, key string) (*MockSystemConfig, error) {
	return s.systemConfigRepo.GetByKey(ctx, key)
}

// UpdateSystemConfig 更新系统配置
func (s *MockAdminService) UpdateSystemConfig(ctx context.Context, key, value, updatedBy string) error {
	// 检查配置是否存在
	config, err := s.systemConfigRepo.GetByKey(ctx, key)
	if err != nil {
		return errors.New("config not found")
	}

	if !config.IsEditable {
		return errors.New("config is not editable")
	}

	return s.systemConfigRepo.Update(ctx, key, value, updatedBy)
}

// ListSystemConfigs 获取系统配置列表
func (s *MockAdminService) ListSystemConfigs(ctx context.Context, category string) ([]*MockSystemConfig, error) {
	return s.systemConfigRepo.List(ctx, category)
}

// GetDashboardStats 获取仪表盘统计数据
func (s *MockAdminService) GetDashboardStats(ctx context.Context, dateRange string) (map[string]interface{}, error) {
	// 计算日期范围
	endDate := time.Now().Format("2006-01-02")
	startDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02")

	if dateRange == "30d" {
		startDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}

	// 获取用户统计
	userStats, _ := s.dataStatsRepo.GetByTypeAndDateRange(ctx, "user", startDate, endDate)

	// 获取内容统计
	contentStats, _ := s.dataStatsRepo.GetByTypeAndDateRange(ctx, "content", startDate, endDate)

	// 获取订单统计
	orderStats, _ := s.dataStatsRepo.GetByTypeAndDateRange(ctx, "order", startDate, endDate)

	// 构建仪表盘数据
	dashboard := map[string]interface{}{
		"overview": map[string]interface{}{
			"totalUsers":    calculateTotal(userStats, "total_users"),
			"activeUsers":   calculateTotal(userStats, "active_users"),
			"totalContents": calculateTotal(contentStats, "total_contents"),
			"totalOrders":   calculateTotal(orderStats, "total_orders"),
			"totalRevenue":  calculateTotal(orderStats, "total_revenue"),
		},
		"trends": map[string]interface{}{
			"userGrowth": buildTrendData(userStats, "new_users"),
			"orderTrend": buildTrendData(orderStats, "new_orders"),
		},
		"distribution": map[string]interface{}{
			"usersByType": map[string]interface{}{
				"normal": 10200,
				"vip":    2360,
			},
			"contentsByCategory": map[string]interface{}{
				"小说": 25600,
				"漫画": 12000,
				"音频": 8000,
			},
		},
	}

	return dashboard, nil
}

// CreateDataStats 创建数据统计
func (s *MockAdminService) CreateDataStats(ctx context.Context, statsType, date string, metrics map[string]int64, dimensions map[string]string) error {
	stats := &MockDataStats{
		ID:         primitive.NewObjectID(),
		Type:       statsType,
		Date:       date,
		Metrics:    metrics,
		Dimensions: dimensions,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return s.dataStatsRepo.Create(ctx, stats)
}

// 辅助函数
func calculateTotal(stats []*MockDataStats, metric string) int64 {
	var total int64
	for _, stat := range stats {
		if value, exists := stat.Metrics[metric]; exists {
			total += value
		}
	}
	return total
}

func buildTrendData(stats []*MockDataStats, metric string) []map[string]interface{} {
	var trends []map[string]interface{}
	for _, stat := range stats {
		if value, exists := stat.Metrics[metric]; exists {
			trends = append(trends, map[string]interface{}{
				"date":  stat.Date,
				"value": value,
			})
		}
	}
	return trends
}

// 测试用例

func TestAdminService_CreateAdminUser_Success(t *testing.T) {
	adminUserRepo := new(MockAdminUserRepository)
	operationLogRepo := new(MockOperationLogRepository)
	systemConfigRepo := new(MockSystemConfigRepository)
	dataStatsRepo := new(MockDataStatsRepository)
	service := NewMockAdminService(adminUserRepo, operationLogRepo, systemConfigRepo, dataStatsRepo)

	ctx := context.Background()
	username := "admin"
	email := "admin@example.com"
	password := "password123"
	realName := "管理员"
	roles := []string{"super_admin"}
	createdBy := "system"

	// Mock 设置
	adminUserRepo.On("GetByUsername", ctx, username).Return(nil, errors.New("not found"))
	adminUserRepo.On("Create", ctx, mock.AnythingOfType("*shared.MockAdminUser")).Return(nil)

	// 执行测试
	adminUser, err := service.CreateAdminUser(ctx, username, email, password, realName, roles, createdBy)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, adminUser)
	assert.Equal(t, username, adminUser.Username)
	assert.Equal(t, email, adminUser.Email)
	assert.Equal(t, realName, adminUser.RealName)
	assert.Equal(t, roles, adminUser.Roles)
	assert.Equal(t, "active", adminUser.Status)
	assert.Equal(t, createdBy, adminUser.CreatedBy)

	adminUserRepo.AssertExpectations(t)
}

func TestAdminService_CreateAdminUser_UsernameExists(t *testing.T) {
	adminUserRepo := new(MockAdminUserRepository)
	operationLogRepo := new(MockOperationLogRepository)
	systemConfigRepo := new(MockSystemConfigRepository)
	dataStatsRepo := new(MockDataStatsRepository)
	service := NewMockAdminService(adminUserRepo, operationLogRepo, systemConfigRepo, dataStatsRepo)

	ctx := context.Background()
	username := "existing_admin"
	email := "admin@example.com"
	password := "password123"
	realName := "管理员"
	roles := []string{"admin"}
	createdBy := "system"

	existingUser := &MockAdminUser{
		ID:       primitive.NewObjectID(),
		Username: username,
		Email:    "existing@example.com",
	}

	// Mock 设置
	adminUserRepo.On("GetByUsername", ctx, username).Return(existingUser, nil)

	// 执行测试
	adminUser, err := service.CreateAdminUser(ctx, username, email, password, realName, roles, createdBy)

	// 断言
	assert.Error(t, err)
	assert.Nil(t, adminUser)
	assert.Equal(t, "username already exists", err.Error())

	adminUserRepo.AssertExpectations(t)
}

func TestAdminService_UpdateAdminUser_Success(t *testing.T) {
	adminUserRepo := new(MockAdminUserRepository)
	operationLogRepo := new(MockOperationLogRepository)
	systemConfigRepo := new(MockSystemConfigRepository)
	dataStatsRepo := new(MockDataStatsRepository)
	service := NewMockAdminService(adminUserRepo, operationLogRepo, systemConfigRepo, dataStatsRepo)

	ctx := context.Background()
	userID := primitive.NewObjectID()
	updates := map[string]interface{}{
		"real_name": "新管理员",
		"status":    "inactive",
	}

	existingUser := &MockAdminUser{
		ID:       userID,
		Username: "admin",
		Email:    "admin@example.com",
		Status:   "active",
	}

	// 使用 mock.MatchedBy 来匹配包含 time.Time 的 map
	expectedUpdates := mock.MatchedBy(func(updates map[string]interface{}) bool {
		if updates["real_name"] != "新管理员" {
			return false
		}
		if updates["status"] != "inactive" {
			return false
		}
		// 检查是否有 updated_at 字段且为 time.Time 类型
		if _, ok := updates["updated_at"]; !ok {
			return false
		}
		return true
	})

	// Mock 设置
	adminUserRepo.On("GetByID", ctx, userID).Return(existingUser, nil)
	adminUserRepo.On("Update", ctx, userID, expectedUpdates).Return(nil)

	// 执行测试
	err := service.UpdateAdminUser(ctx, userID, updates)

	// 断言
	assert.NoError(t, err)

	adminUserRepo.AssertExpectations(t)
}

func TestAdminService_LogOperation_Success(t *testing.T) {
	adminUserRepo := new(MockAdminUserRepository)
	operationLogRepo := new(MockOperationLogRepository)
	systemConfigRepo := new(MockSystemConfigRepository)
	dataStatsRepo := new(MockDataStatsRepository)
	service := NewMockAdminService(adminUserRepo, operationLogRepo, systemConfigRepo, dataStatsRepo)

	ctx := context.Background()
	adminUserID := primitive.NewObjectID()
	username := "admin"
	module := "user"
	action := "create"
	resource := "admin_user"
	resourceID := "user_123"
	method := "POST"
	url := "/admin/api/v1/users"
	ipAddress := "192.168.1.1"
	userAgent := "Mozilla/5.0"
	requestData := map[string]interface{}{"username": "newuser"}
	responseData := map[string]interface{}{"id": "user_123"}
	status := "success"
	errorMessage := ""
	duration := int64(150)

	// Mock 设置
	operationLogRepo.On("Create", ctx, mock.AnythingOfType("*shared.MockOperationLog")).Return(nil)

	// 执行测试
	err := service.LogOperation(ctx, adminUserID, username, module, action, resource, resourceID, method, url, ipAddress, userAgent, requestData, responseData, status, errorMessage, duration)

	// 断言
	assert.NoError(t, err)

	operationLogRepo.AssertExpectations(t)
}

func TestAdminService_UpdateSystemConfig_Success(t *testing.T) {
	adminUserRepo := new(MockAdminUserRepository)
	operationLogRepo := new(MockOperationLogRepository)
	systemConfigRepo := new(MockSystemConfigRepository)
	dataStatsRepo := new(MockDataStatsRepository)
	service := NewMockAdminService(adminUserRepo, operationLogRepo, systemConfigRepo, dataStatsRepo)

	ctx := context.Background()
	key := "site.title"
	value := "青羽平台"
	updatedBy := "admin"

	existingConfig := &MockSystemConfig{
		ID:         primitive.NewObjectID(),
		Key:        key,
		Value:      "旧标题",
		IsEditable: true,
	}

	// Mock 设置
	systemConfigRepo.On("GetByKey", ctx, key).Return(existingConfig, nil)
	systemConfigRepo.On("Update", ctx, key, value, updatedBy).Return(nil)

	// 执行测试
	err := service.UpdateSystemConfig(ctx, key, value, updatedBy)

	// 断言
	assert.NoError(t, err)

	systemConfigRepo.AssertExpectations(t)
}

func TestAdminService_UpdateSystemConfig_NotEditable(t *testing.T) {
	adminUserRepo := new(MockAdminUserRepository)
	operationLogRepo := new(MockOperationLogRepository)
	systemConfigRepo := new(MockSystemConfigRepository)
	dataStatsRepo := new(MockDataStatsRepository)
	service := NewMockAdminService(adminUserRepo, operationLogRepo, systemConfigRepo, dataStatsRepo)

	ctx := context.Background()
	key := "system.version"
	value := "2.0.0"
	updatedBy := "admin"

	existingConfig := &MockSystemConfig{
		ID:         primitive.NewObjectID(),
		Key:        key,
		Value:      "1.0.0",
		IsEditable: false,
	}

	// Mock 设置
	systemConfigRepo.On("GetByKey", ctx, key).Return(existingConfig, nil)

	// 执行测试
	err := service.UpdateSystemConfig(ctx, key, value, updatedBy)

	// 断言
	assert.Error(t, err)
	assert.Equal(t, "config is not editable", err.Error())

	systemConfigRepo.AssertExpectations(t)
}

func TestAdminService_GetDashboardStats_Success(t *testing.T) {
	adminUserRepo := new(MockAdminUserRepository)
	operationLogRepo := new(MockOperationLogRepository)
	systemConfigRepo := new(MockSystemConfigRepository)
	dataStatsRepo := new(MockDataStatsRepository)
	service := NewMockAdminService(adminUserRepo, operationLogRepo, systemConfigRepo, dataStatsRepo)

	ctx := context.Background()
	dateRange := "7d"

	userStats := []*MockDataStats{
		{
			Type: "user",
			Date: "2024-01-15",
			Metrics: map[string]int64{
				"total_users":  12560,
				"active_users": 8900,
				"new_users":    120,
			},
		},
	}

	contentStats := []*MockDataStats{
		{
			Type: "content",
			Date: "2024-01-15",
			Metrics: map[string]int64{
				"total_contents": 45600,
			},
		},
	}

	orderStats := []*MockDataStats{
		{
			Type: "order",
			Date: "2024-01-15",
			Metrics: map[string]int64{
				"total_orders":  23400,
				"total_revenue": 1234567,
				"new_orders":    45,
			},
		},
	}

	// Mock 设置
	dataStatsRepo.On("GetByTypeAndDateRange", ctx, "user", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(userStats, nil)
	dataStatsRepo.On("GetByTypeAndDateRange", ctx, "content", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(contentStats, nil)
	dataStatsRepo.On("GetByTypeAndDateRange", ctx, "order", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(orderStats, nil)

	// 执行测试
	dashboard, err := service.GetDashboardStats(ctx, dateRange)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, dashboard)

	overview := dashboard["overview"].(map[string]interface{})
	assert.Equal(t, int64(12560), overview["totalUsers"])
	assert.Equal(t, int64(8900), overview["activeUsers"])
	assert.Equal(t, int64(45600), overview["totalContents"])
	assert.Equal(t, int64(23400), overview["totalOrders"])
	assert.Equal(t, int64(1234567), overview["totalRevenue"])

	trends := dashboard["trends"].(map[string]interface{})
	assert.NotNil(t, trends["userGrowth"])
	assert.NotNil(t, trends["orderTrend"])

	dataStatsRepo.AssertExpectations(t)
}

func TestAdminService_CreateDataStats_Success(t *testing.T) {
	adminUserRepo := new(MockAdminUserRepository)
	operationLogRepo := new(MockOperationLogRepository)
	systemConfigRepo := new(MockSystemConfigRepository)
	dataStatsRepo := new(MockDataStatsRepository)
	service := NewMockAdminService(adminUserRepo, operationLogRepo, systemConfigRepo, dataStatsRepo)

	ctx := context.Background()
	statsType := "user"
	date := "2024-01-15"
	metrics := map[string]int64{
		"total_users":  12560,
		"active_users": 8900,
		"new_users":    120,
	}
	dimensions := map[string]string{
		"platform": "web",
		"region":   "cn",
	}

	// Mock 设置
	dataStatsRepo.On("Create", ctx, mock.AnythingOfType("*shared.MockDataStats")).Return(nil)

	// 执行测试
	err := service.CreateDataStats(ctx, statsType, date, metrics, dimensions)

	// 断言
	assert.NoError(t, err)

	dataStatsRepo.AssertExpectations(t)
}
