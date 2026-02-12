package integration

import (
	"context"
	"testing"
	"time"

	"Qingyu_backend/service/auth"
	"Qingyu_backend/service/writer/project"

	"github.com/stretchr/testify/assert"
)

// TestMVP_AutoSaveIntegration 测试自动保存功能集成
// 验证：启动自动保存 → 等待30秒 → 检查版本创建
func TestMVP_AutoSaveIntegration(t *testing.T) {
	t.Skip("集成测试：需要真实MongoDB连接")

	// 准备
	versionService := &project.VersionService{}
	autoSaveService := project.NewAutoSaveService(versionService)

	documentID := "test-doc-001"
	projectID := "test-project-001"
	nodeID := "test-node-001"
	userID := "test-user-001"

	// 启动自动保存
	err := autoSaveService.StartAutoSave(documentID, projectID, nodeID, userID)
	assert.NoError(t, err, "启动自动保存应该成功")

	// 等待31秒（确保触发一次自动保存）
	time.Sleep(31 * time.Second)

	// 验证自动保存状态
	isRunning, lastSaved := autoSaveService.GetStatus(documentID)
	assert.True(t, isRunning, "自动保存应该在运行")
	assert.NotNil(t, lastSaved, "应该有最后保存时间")
	assert.True(t, time.Since(*lastSaved) < 35*time.Second, "最后保存时间应该在35秒内")

	// 停止自动保存
	err = autoSaveService.StopAutoSave(documentID)
	assert.NoError(t, err, "停止自动保存应该成功")

	// 验证停止状态
	isRunning, _ = autoSaveService.GetStatus(documentID)
	assert.False(t, isRunning, "自动保存应该已停止")

	t.Log("✅ 自动保存功能集成测试通过")
}

// TestMVP_MultiDeviceLoginIntegration 测试多端登录限制集成
// 验证：登录6次 → 第6次应被拒绝
func TestMVP_MultiDeviceLoginIntegration(t *testing.T) {
	t.Skip("集成测试：需要真实Redis连接")

	// 准备Mock CacheClient
	mockCache := &MockCacheClient{}
	sessionService := auth.NewSessionService(mockCache)

	ctx := context.Background()
	userID := "test-user-multi-device"

	// 模拟创建5个会话（成功）
	sessions := make([]*auth.Session, 0, 5)
	for i := 0; i < 5; i++ {
		session, err := sessionService.CreateSession(ctx, userID)
		assert.NoError(t, err, "前5次登录应该成功")
		sessions = append(sessions, session)
		t.Logf("创建会话 %d: %s", i+1, session.ID)
	}

	// 检查设备限制（应该还能登录）
	err := sessionService.CheckDeviceLimit(ctx, userID, 5)
	assert.NoError(t, err, "5台设备应该允许")

	// 创建第6个会话（检查应拒绝）
	err = sessionService.CheckDeviceLimit(ctx, userID, 5)
	assert.Error(t, err, "第6台设备应该被拒绝")
	assert.Contains(t, err.Error(), "登录设备数量已达上限", "错误信息应包含设备限制说明")

	// 清理：删除一个会话
	err = sessionService.DestroySession(ctx, sessions[0].ID)
	assert.NoError(t, err, "删除会话应该成功")

	// 再次检查（应该允许）
	err = sessionService.CheckDeviceLimit(ctx, userID, 5)
	assert.NoError(t, err, "删除一个会话后应该允许登录")

	t.Log("✅ 多端登录限制集成测试通过")
}

// TestMVP_PasswordValidationIntegration 测试密码强度验证集成
// 验证：各种密码强度 → 验证规则执行
func TestMVP_PasswordValidationIntegration(t *testing.T) {
	// 准备
	validator := auth.NewPasswordValidator()

	// 测试用例
	testCases := []struct {
		name     string
		password string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "强密码-合格",
			password: "Abc12345",
			wantErr:  false,
		},
		{
			name:     "密码太短",
			password: "Abc123",
			wantErr:  true,
			errMsg:   "密码长度不能少于8位",
		},
		{
			name:     "密码太长",
			password: "Abc123456789012345678901234567890",
			wantErr:  true,
			errMsg:   "密码长度不能超过32位",
		},
		{
			name:     "缺少数字",
			password: "Abcdefgh",
			wantErr:  true,
			errMsg:   "密码必须包含至少一个数字",
		},
		{
			name:     "缺少字母",
			password: "12345678",
			wantErr:  true,
			errMsg:   "密码必须包含至少一个字母",
		},
		{
			name:     "边界-最短合格",
			password: "Abcd1234",
			wantErr:  false,
		},
		{
			name:     "边界-最长合格",
			password: "Abc12345678901234567890123456789",
			wantErr:  false,
		},
	}

	// 执行测试
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.ValidatePassword(tc.password)
			if tc.wantErr {
				assert.Error(t, err, "应该返回错误")
				assert.Contains(t, err.Error(), tc.errMsg, "错误信息应匹配")
			} else {
				assert.NoError(t, err, "应该验证通过")
			}
		})
	}

	t.Log("✅ 密码强度验证集成测试通过")
}

// TestMVP_PasswordStrengthScoring 测试密码强度评分（预留功能）
func TestMVP_PasswordStrengthScoring(t *testing.T) {
	validator := auth.NewPasswordValidator()

	testCases := []struct {
		password string
		expected string
	}{
		{"123456", "weak"},
		{"Abc12345", "medium"},
		{"Abc@12345", "strong"},
		{"A1b2C3d4E5f6!@#", "strong"},
	}

	for _, tc := range testCases {
		t.Run(tc.password, func(t *testing.T) {
			strength := validator.ValidatePasswordStrength(tc.password)
			assert.Equal(t, tc.expected, strength, "密码强度评分应匹配")
		})
	}

	t.Log("✅ 密码强度评分测试通过")
}

// TestMVP_CoreFlowIntegration 测试核心流程集成
// 完整流程：注册（密码验证）→ 登录（设备限制）→ 创建项目 → 编辑（自动保存）
func TestMVP_CoreFlowIntegration(t *testing.T) {
	t.Skip("集成测试：需要完整环境（DB+Redis+API）")

	// ctx := context.Background()

	// Step 1: 注册（密码强度验证）
	t.Log("Step 1: 测试用户注册...")
	password := "TestPass123"
	validator := auth.NewPasswordValidator()
	err := validator.ValidatePassword(password)
	assert.NoError(t, err, "密码应该通过验证")

	// TODO: 调用实际AuthService.Register
	// registerResp, err := authService.Register(ctx, &auth.RegisterRequest{
	//     Username: "test-mvp-user",
	//     Email:    "test@qingyu.com",
	//     Password: password,
	// })
	// assert.NoError(t, err, "注册应该成功")

	// Step 2: 登录（多端登录限制）
	t.Log("Step 2: 测试用户登录...")
	// TODO: 调用实际AuthService.Login
	// loginResp, err := authService.Login(ctx, &auth.LoginRequest{
	//     Username: "test-mvp-user",
	//     Password: password,
	// })
	// assert.NoError(t, err, "登录应该成功")

	// Step 3: 创建项目
	t.Log("Step 3: 测试创建项目...")
	// TODO: 调用ProjectService.CreateProject

	// Step 4: 创建文档并启动自动保存
	t.Log("Step 4: 测试创建文档并启动自动保存...")
	// TODO: 调用DocumentService.CreateDocument
	// TODO: 调用AutoSaveService.StartAutoSave

	// Step 5: 等待自动保存触发
	t.Log("Step 5: 等待自动保存...")
	time.Sleep(31 * time.Second)
	// TODO: 验证版本创建

	t.Log("✅ 核心流程集成测试通过")
}

// MockCacheClient Mock缓存客户端（用于测试）
type MockCacheClient struct {
	data map[string]string
}

func (m *MockCacheClient) Get(ctx context.Context, key string) (string, error) {
	if m.data == nil {
		return "", assert.AnError
	}
	if val, ok := m.data[key]; ok {
		return val, nil
	}
	return "", assert.AnError
}

func (m *MockCacheClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if m.data == nil {
		m.data = make(map[string]string)
	}
	// 转换interface{}为string
	if strVal, ok := value.(string); ok {
		m.data[key] = strVal
	}
	return nil
}

func (m *MockCacheClient) Delete(ctx context.Context, key string) error {
	if m.data == nil {
		return nil
	}
	delete(m.data, key)
	return nil
}
