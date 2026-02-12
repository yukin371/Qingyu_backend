package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	usersModel "Qingyu_backend/models/users"
)

// ==================== Test Helpers ====================

func createTestObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

// ==================== TransactionManager Tests ====================

func TestNewTransactionManager(t *testing.T) {
	tm := &TransactionManager{}
	assert.NotNil(t, tm)
}

// ==================== UserRegistrationTransaction Tests ====================

func TestUserRegistrationTransaction_GetDescription(t *testing.T) {
	userID := createTestObjectID()
	user := &usersModel.User{}
	user.ID = userID
	user.Username = "testuser"

	urt := &UserRegistrationTransaction{
		User: user,
	}

	description := urt.GetDescription()

	assert.Contains(t, description, "testuser")
	assert.Contains(t, description, "用户注册事务")
}

func TestUserRegistrationTransaction_Execute_NoID(t *testing.T) {
	user := &usersModel.User{}
	user.Username = "testuser"
	user.Email = "test@example.com"

	urt := &UserRegistrationTransaction{
		User: user,
	}

	assert.True(t, urt.User.ID.IsZero())
}

func TestUserRegistrationTransaction_Execute_WithID(t *testing.T) {
	userID := createTestObjectID()
	user := &usersModel.User{}
	user.ID = userID
	user.Username = "testuser"

	urt := &UserRegistrationTransaction{
		User: user,
	}

	assert.False(t, urt.User.ID.IsZero())
	assert.Equal(t, userID, urt.User.ID)
}

func TestUserRegistrationTransaction_Execute_WithUserRole(t *testing.T) {
	userID := createTestObjectID()
	roleID := createTestObjectID()
	user := &usersModel.User{}
	user.ID = userID
	user.Username = "testuser"

	urt := &UserRegistrationTransaction{
		User: user,
		UserRole: &UserRole{
			UserID: userID,
			RoleID: roleID,
		},
	}

	assert.NotNil(t, urt.UserRole)
	assert.Equal(t, userID, urt.UserRole.UserID)
	assert.Equal(t, roleID, urt.UserRole.RoleID)
}

func TestUserRegistrationTransaction_Execute_WithUserConfig(t *testing.T) {
	userID := createTestObjectID()
	user := &usersModel.User{}
	user.ID = userID
	user.Username = "testuser"

	urt := &UserRegistrationTransaction{
		User: user,
		UserConfig: &UserConfig{
			UserID:   userID,
			Theme:    "dark",
			Language: "en-US",
		},
	}

	assert.NotNil(t, urt.UserConfig)
	assert.Equal(t, userID, urt.UserConfig.UserID)
	assert.Equal(t, "dark", urt.UserConfig.Theme)
	assert.Equal(t, "en-US", urt.UserConfig.Language)
}

func TestUserRegistrationTransaction_Execute_DefaultValues(t *testing.T) {
	userID := createTestObjectID()
	user := &usersModel.User{}
	user.ID = userID
	user.Username = "testuser"

	urt := &UserRegistrationTransaction{
		User: user,
		UserConfig: &UserConfig{
			UserID: userID,
		},
	}

	assert.NotNil(t, urt.UserConfig)
	assert.Equal(t, userID, urt.UserConfig.UserID)
	assert.Empty(t, urt.UserConfig.Theme)
	assert.Empty(t, urt.UserConfig.Language)
}

func TestUserRegistrationTransaction_Execute_NilUser(t *testing.T) {
	urt := &UserRegistrationTransaction{
		User: nil,
	}

	assert.Nil(t, urt.User)
}

func TestUserRegistrationTransaction_Execute_NilUserRole(t *testing.T) {
	userID := createTestObjectID()
	user := &usersModel.User{}
	user.ID = userID

	urt := &UserRegistrationTransaction{
		User:     user,
		UserRole: nil,
	}

	assert.Nil(t, urt.UserRole)
}

func TestUserRegistrationTransaction_Execute_NilUserConfig(t *testing.T) {
	userID := createTestObjectID()
	user := &usersModel.User{}
	user.ID = userID

	urt := &UserRegistrationTransaction{
		User:       user,
		UserConfig: nil,
	}

	assert.Nil(t, urt.UserConfig)
}

// ==================== UserDeletionTransaction Tests ====================

func TestUserDeletionTransaction_GetDescription_SoftDelete(t *testing.T) {
	userID := createTestObjectID()
	udt := &UserDeletionTransaction{
		UserID:     userID,
		SoftDelete: true,
	}

	description := udt.GetDescription()

	assert.Contains(t, description, userID.Hex())
	assert.Contains(t, description, "软删除")
}

func TestUserDeletionTransaction_GetDescription_HardDelete(t *testing.T) {
	userID := createTestObjectID()
	udt := &UserDeletionTransaction{
		UserID:     userID,
		SoftDelete: false,
	}

	description := udt.GetDescription()

	assert.Contains(t, description, userID.Hex())
	assert.Contains(t, description, "硬删除")
}

func TestUserDeletionTransaction_Execute_SoftDelete(t *testing.T) {
	userID := createTestObjectID()
	udt := &UserDeletionTransaction{
		UserID:     userID,
		SoftDelete: true,
	}

	assert.True(t, udt.SoftDelete)
}

func TestUserDeletionTransaction_Execute_HardDelete(t *testing.T) {
	userID := createTestObjectID()
	udt := &UserDeletionTransaction{
		UserID:     userID,
		SoftDelete: false,
	}

	assert.False(t, udt.SoftDelete)
}

// ==================== CascadeManager Tests ====================

func TestNewCascadeManager(t *testing.T) {
	cm := NewCascadeManager(nil, nil)

	assert.NotNil(t, cm)
	assert.Nil(t, cm.db)
	assert.Nil(t, cm.tm)
}

// ==================== UserUpdateTransaction Tests ====================

func TestUserUpdateTransaction_GetDescription(t *testing.T) {
	userID := createTestObjectID()
	uut := &UserUpdateTransaction{
		UserID:  userID,
		Updates: bson.M{"username": "newusername"},
	}

	description := uut.GetDescription()

	assert.Contains(t, description, userID.Hex())
	assert.Contains(t, description, "用户更新事务")
}

func TestUserUpdateTransaction_Execute_WithUpdates(t *testing.T) {
	userID := createTestObjectID()
	updates := bson.M{
		"username": "newusername",
		"email":    "newemail@example.com",
	}

	uut := &UserUpdateTransaction{
		UserID:  userID,
		Updates: updates,
	}

	assert.Equal(t, userID, uut.UserID)
	assert.NotNil(t, uut.Updates)
	assert.Equal(t, "newusername", uut.Updates["username"])
}

func TestUserUpdateTransaction_Execute_AddsTimestamp(t *testing.T) {
	userID := createTestObjectID()
	uut := &UserUpdateTransaction{
		UserID:  userID,
		Updates: bson.M{"username": "newusername"},
	}

	_, exists := uut.Updates["updated_at"]
	assert.False(t, exists)
}

func TestUserUpdateTransaction_Execute_EmptyUpdates(t *testing.T) {
	userID := createTestObjectID()
	uut := &UserUpdateTransaction{
		UserID:  userID,
		Updates: bson.M{},
	}

	assert.Empty(t, uut.Updates)
}

func TestUserUpdateTransaction_Execute_NilUpdates(t *testing.T) {
	userID := createTestObjectID()
	uut := &UserUpdateTransaction{
		UserID:  userID,
		Updates: nil,
	}

	assert.Nil(t, uut.Updates)
}

// ==================== SagaManager Tests ====================

func TestNewSagaManager(t *testing.T) {
	sm := NewSagaManager(nil)

	assert.NotNil(t, sm)
	assert.Nil(t, sm.db)
}

func TestSagaManager_ExecuteSaga_Success(t *testing.T) {
	sm := &SagaManager{}
	ctx := context.Background()

	executedOrder := []string{}
	steps := []SagaStep{
		{
			Name: "step1",
			Execute: func(ctx context.Context) error {
				executedOrder = append(executedOrder, "step1")
				return nil
			},
			Compensate: func(ctx context.Context) error {
				executedOrder = append(executedOrder, "compensate1")
				return nil
			},
		},
		{
			Name: "step2",
			Execute: func(ctx context.Context) error {
				executedOrder = append(executedOrder, "step2")
				return nil
			},
			Compensate: func(ctx context.Context) error {
				executedOrder = append(executedOrder, "compensate2")
				return nil
			},
		},
	}

	err := sm.ExecuteSaga(ctx, steps)

	assert.NoError(t, err)
	assert.Equal(t, []string{"step1", "step2"}, executedOrder)
}

func TestSagaManager_ExecuteSaga_StepError(t *testing.T) {
	sm := &SagaManager{}
	ctx := context.Background()

	executedOrder := []string{}
	steps := []SagaStep{
		{
			Name: "step1",
			Execute: func(ctx context.Context) error {
				executedOrder = append(executedOrder, "step1")
				return nil
			},
			Compensate: func(ctx context.Context) error {
				executedOrder = append(executedOrder, "compensate1")
				return nil
			},
		},
		{
			Name: "step2",
			Execute: func(ctx context.Context) error {
				executedOrder = append(executedOrder, "step2")
				return errors.New("step2 failed")
			},
			Compensate: func(ctx context.Context) error {
				executedOrder = append(executedOrder, "compensate2")
				return nil
			},
		},
	}

	err := sm.ExecuteSaga(ctx, steps)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "step2 failed")
	assert.Equal(t, []string{"step1", "step2", "compensate1"}, executedOrder)
}

func TestSagaManager_ExecuteSaga_CompensateError(t *testing.T) {
	sm := &SagaManager{}
	ctx := context.Background()

	executedOrder := []string{}
	steps := []SagaStep{
		{
			Name: "step1",
			Execute: func(ctx context.Context) error {
				executedOrder = append(executedOrder, "step1")
				return nil
			},
			Compensate: func(ctx context.Context) error {
				executedOrder = append(executedOrder, "compensate1")
				return errors.New("compensate1 failed")
			},
		},
		{
			Name: "step2",
			Execute: func(ctx context.Context) error {
				executedOrder = append(executedOrder, "step2")
				return errors.New("step2 failed")
			},
			Compensate: func(ctx context.Context) error {
				executedOrder = append(executedOrder, "compensate2")
				return nil
			},
		},
	}

	err := sm.ExecuteSaga(ctx, steps)

	assert.Error(t, err)
	assert.Equal(t, []string{"step1", "step2", "compensate1"}, executedOrder)
}

func TestSagaManager_ExecuteSaga_EmptySteps(t *testing.T) {
	sm := &SagaManager{}
	ctx := context.Background()

	steps := []SagaStep{}

	err := sm.ExecuteSaga(ctx, steps)

	assert.NoError(t, err)
}

func TestSagaManager_ExecuteSaga_SingleStep(t *testing.T) {
	sm := &SagaManager{}
	ctx := context.Background()

	executed := false
	steps := []SagaStep{
		{
			Name: "single",
			Execute: func(ctx context.Context) error {
				executed = true
				return nil
			},
			Compensate: func(ctx context.Context) error {
				return nil
			},
		},
	}

	err := sm.ExecuteSaga(ctx, steps)

	assert.NoError(t, err)
	assert.True(t, executed)
}

// ==================== ReferenceIntegrityManager Tests ====================

func TestNewReferenceIntegrityManager(t *testing.T) {
	rim := NewReferenceIntegrityManager(nil)

	assert.NotNil(t, rim)
	assert.Nil(t, rim.db)
}

func TestReferenceIntegrityManager_ValidateReferences_UnknownCollection(t *testing.T) {
	rim := &ReferenceIntegrityManager{}
	ctx := context.Background()

	doc := map[string]interface{}{"key": "value"}
	err := rim.ValidateReferences(ctx, doc, "unknown_collection")

	assert.NoError(t, err)
}

func TestReferenceIntegrityManager_validateProjectReferences_InvalidFormat(t *testing.T) {
	rim := &ReferenceIntegrityManager{}
	ctx := context.Background()

	doc := "not a map"
	err := rim.validateProjectReferences(ctx, doc)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "无效的项目文档格式")
}

func TestReferenceIntegrityManager_validateProjectReferences_MissingCreatorID(t *testing.T) {
	rim := &ReferenceIntegrityManager{}
	ctx := context.Background()

	doc := map[string]interface{}{
		"name": "test project",
	}
	err := rim.validateProjectReferences(ctx, doc)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "项目创建者ID格式错误")
}

func TestReferenceIntegrityManager_validateProjectReferences_InvalidCreatorID(t *testing.T) {
	rim := &ReferenceIntegrityManager{}
	ctx := context.Background()

	doc := map[string]interface{}{
		"creator_id": "invalid_id",
	}
	err := rim.validateProjectReferences(ctx, doc)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "项目创建者ID格式错误")
}

func TestReferenceIntegrityManager_validateProjectReferences_NilCreatorID(t *testing.T) {
	rim := &ReferenceIntegrityManager{}
	ctx := context.Background()

	doc := map[string]interface{}{
		"creator_id": nil,
	}
	err := rim.validateProjectReferences(ctx, doc)

	assert.Error(t, err)
}

func TestReferenceIntegrityManager_validateUserRoleReferences_InvalidFormat(t *testing.T) {
	rim := &ReferenceIntegrityManager{}
	ctx := context.Background()

	doc := "not a map"
	err := rim.validateUserRoleReferences(ctx, doc)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "无效的用户角色文档格式")
}

func TestReferenceIntegrityManager_validateUserRoleReferences_MissingUserID(t *testing.T) {
	rim := &ReferenceIntegrityManager{}
	ctx := context.Background()

	roleID := createTestObjectID()
	doc := map[string]interface{}{
		"role_id": roleID,
	}
	err := rim.validateUserRoleReferences(ctx, doc)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "用户ID格式错误")
}

func TestReferenceIntegrityManager_validateUserRoleReferences_MissingRoleID(t *testing.T) {
	rim := &ReferenceIntegrityManager{}
	ctx := context.Background()

	userID := createTestObjectID()
	doc := map[string]interface{}{
		"user_id": userID,
	}
	err := rim.validateUserRoleReferences(ctx, doc)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "角色ID格式错误")
}

func TestReferenceIntegrityManager_validateUserRoleReferences_InvalidUserID(t *testing.T) {
	rim := &ReferenceIntegrityManager{}
	ctx := context.Background()

	roleID := createTestObjectID()
	doc := map[string]interface{}{
		"user_id": "invalid",
		"role_id": roleID,
	}
	err := rim.validateUserRoleReferences(ctx, doc)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "用户ID格式错误")
}

func TestReferenceIntegrityManager_validateUserRoleReferences_InvalidRoleID(t *testing.T) {
	rim := &ReferenceIntegrityManager{}
	ctx := context.Background()

	userID := createTestObjectID()
	doc := map[string]interface{}{
		"user_id": userID,
		"role_id": "invalid",
	}
	err := rim.validateUserRoleReferences(ctx, doc)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "角色ID格式错误")
}

func TestReferenceIntegrityManager_validateUserRoleReferences_NilUserID(t *testing.T) {
	rim := &ReferenceIntegrityManager{}
	ctx := context.Background()

	roleID := createTestObjectID()
	doc := map[string]interface{}{
		"user_id": nil,
		"role_id": roleID,
	}
	err := rim.validateUserRoleReferences(ctx, doc)

	assert.Error(t, err)
}

func TestReferenceIntegrityManager_validateUserRoleReferences_NilRoleID(t *testing.T) {
	rim := &ReferenceIntegrityManager{}
	ctx := context.Background()

	userID := createTestObjectID()
	doc := map[string]interface{}{
		"user_id": userID,
		"role_id": nil,
	}
	err := rim.validateUserRoleReferences(ctx, doc)

	assert.Error(t, err)
}

// ==================== Type Tests ====================

func TestUserRole_Struct(t *testing.T) {
	userID := createTestObjectID()
	roleID := createTestObjectID()
	now := time.Now()

	ur := UserRole{}
	ur.ID = createTestObjectID()
	ur.UserID = userID
	ur.RoleID = roleID
	ur.CreatedAt = now
	ur.UpdatedAt = now

	assert.False(t, ur.ID.IsZero())
	assert.Equal(t, userID, ur.UserID)
	assert.Equal(t, roleID, ur.RoleID)
	assert.False(t, ur.CreatedAt.IsZero())
	assert.False(t, ur.UpdatedAt.IsZero())
}

func TestUserConfig_Struct(t *testing.T) {
	userID := createTestObjectID()
	now := time.Now()

	uc := UserConfig{}
	uc.ID = createTestObjectID()
	uc.UserID = userID
	uc.Theme = "dark"
	uc.Language = "zh-CN"
	uc.Settings = map[string]interface{}{"key": "value"}
	uc.CreatedAt = now
	uc.UpdatedAt = now

	assert.False(t, uc.ID.IsZero())
	assert.Equal(t, userID, uc.UserID)
	assert.Equal(t, "dark", uc.Theme)
	assert.Equal(t, "zh-CN", uc.Language)
	assert.NotNil(t, uc.Settings)
	assert.Equal(t, "value", uc.Settings["key"])
	assert.False(t, uc.CreatedAt.IsZero())
	assert.False(t, uc.UpdatedAt.IsZero())
}

func TestUserConfig_DefaultSettings(t *testing.T) {
	uc := UserConfig{
		Settings: nil,
	}

	assert.Nil(t, uc.Settings)
}

// ==================== Table-Driven Tests ====================

func TestSagaManager_ExecuteSaga_TableDriven(t *testing.T) {
	sm := &SagaManager{}
	ctx := context.Background()

	tests := []struct {
		name        string
		steps       []SagaStep
		expectError bool
		errorMsg    string
	}{
		{
			name: "成功执行单个步骤",
			steps: []SagaStep{
				{
					Name:       "step1",
					Execute:    func(ctx context.Context) error { return nil },
					Compensate: func(ctx context.Context) error { return nil },
				},
			},
			expectError: false,
		},
		{
			name: "成功执行多个步骤",
			steps: []SagaStep{
				{
					Name:       "step1",
					Execute:    func(ctx context.Context) error { return nil },
					Compensate: func(ctx context.Context) error { return nil },
				},
				{
					Name:       "step2",
					Execute:    func(ctx context.Context) error { return nil },
					Compensate: func(ctx context.Context) error { return nil },
				},
			},
			expectError: false,
		},
		{
			name: "第一个步骤失败",
			steps: []SagaStep{
				{
					Name:       "step1",
					Execute:    func(ctx context.Context) error { return errors.New("step1 error") },
					Compensate: func(ctx context.Context) error { return nil },
				},
			},
			expectError: true,
			errorMsg:    "step1 error",
		},
		{
			name: "中间步骤失败并补偿",
			steps: []SagaStep{
				{
					Name:       "step1",
					Execute:    func(ctx context.Context) error { return nil },
					Compensate: func(ctx context.Context) error { return nil },
				},
				{
					Name:       "step2",
					Execute:    func(ctx context.Context) error { return errors.New("step2 error") },
					Compensate: func(ctx context.Context) error { return nil },
				},
			},
			expectError: true,
			errorMsg:    "step2 error",
		},
		{
			name:        "空步骤列表",
			steps:       []SagaStep{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sm.ExecuteSaga(ctx, tt.steps)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ==================== Integration Tests ====================

func TestIntegration_TransactionManager(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}
	t.Skip("需要配置MongoDB集成测试环境")
}

func TestIntegration_UserRegistrationFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}
	t.Skip("需要配置MongoDB集成测试环境")
}
