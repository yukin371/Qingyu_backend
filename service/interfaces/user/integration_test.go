package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/dto"
	"Qingyu_backend/service/interfaces/user"
)

// TestUserServicePort_Integration_EndToEnd 端到端集成测试
// 验证Port接口、Adapter实现和底层服务的完整协作
func TestUserServicePort_Integration_EndToEnd(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// 注意：这是一个集成测试框架示例
	// 实际运行需要：
	// 1. 设置测试数据库
	// 2. 初始化真实的服务实例
	// 3. 执行测试用例
	// 4. 清理测试数据

	t.Run("Port接口完整性验证", func(t *testing.T) {
		// 验证Port接口包含所有必需的方法
		// 这是一个编译时测试，确保接口定义正确
		type portInterface interface {
			// BaseService方法
			Initialize(ctx context.Context) error
			Health(ctx context.Context) error
			Close(ctx context.Context) error
			GetServiceName() string
			GetVersion() string

			// 用户管理方法
			CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error)
			GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error)
			UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.UpdateUserResponse, error)
			DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*user.DeleteUserResponse, error)
			ListUsers(ctx context.Context, req *user.ListUsersRequest) (*user.ListUsersResponse, error)

			// 用户认证方法
			RegisterUser(ctx context.Context, req *user.RegisterUserRequest) (*user.RegisterUserResponse, error)
			LoginUser(ctx context.Context, req *user.LoginUserRequest) (*user.LoginUserResponse, error)
			LogoutUser(ctx context.Context, req *user.LogoutUserRequest) (*user.LogoutUserResponse, error)
			ValidateToken(ctx context.Context, req *user.ValidateTokenRequest) (*user.ValidateTokenResponse, error)

			// 其他方法...
		}

		// 这个变量确保Port接口符合预期
		var _ portInterface = user.UserService(nil)
	})
}

// TestUserServicePort_DTOValidation DTO结构验证测试
func TestUserServicePort_DTOValidation(t *testing.T) {
	t.Run("CreateUserRequest结构验证", func(t *testing.T) {
		req := &user.CreateUserRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
			Role:     "reader",
		}
		assert.Equal(t, "testuser", req.Username)
		assert.Equal(t, "test@example.com", req.Email)
		assert.Equal(t, "password123", req.Password)
		assert.Equal(t, "reader", req.Role)
	})

	t.Run("GetUserRequest结构验证", func(t *testing.T) {
		userID := primitive.NewObjectID().Hex()
		req := &user.GetUserRequest{
			ID: userID,
		}
		assert.Equal(t, userID, req.ID)
	})

	t.Run("UpdateUserRequest结构验证", func(t *testing.T) {
		userID := primitive.NewObjectID().Hex()
		updates := map[string]interface{}{
			"username": "newusername",
			"email":    "new@example.com",
		}
		req := &user.UpdateUserRequest{
			ID:      userID,
			Updates: updates,
		}
		assert.Equal(t, userID, req.ID)
		assert.Equal(t, "newusername", req.Updates["username"])
		assert.Equal(t, "new@example.com", req.Updates["email"])
	})

	t.Run("DeleteUserRequest结构验证", func(t *testing.T) {
		userID := primitive.NewObjectID().Hex()
		req := &user.DeleteUserRequest{
			ID: userID,
		}
		assert.Equal(t, userID, req.ID)
	})

	t.Run("ListUsersRequest结构验证", func(t *testing.T) {
		req := &user.ListUsersRequest{
			Page:     1,
			PageSize: 20,
			Status:   "active",
		}
		assert.Equal(t, 1, req.Page)
		assert.Equal(t, 20, req.PageSize)
		assert.Equal(t, "active", req.Status)
	})

	t.Run("ListUsersResponse结构验证", func(t *testing.T) {
		resp := &user.ListUsersResponse{
			Total:      100,
			Page:       1,
			PageSize:   20,
			TotalPages: 5,
		}
		assert.Equal(t, int64(100), resp.Total)
		assert.Equal(t, 1, resp.Page)
		assert.Equal(t, 20, resp.PageSize)
		assert.Equal(t, 5, resp.TotalPages)
	})

	t.Run("LoginUserRequest结构验证", func(t *testing.T) {
		req := &user.LoginUserRequest{
			Username: "testuser",
			Password: "password123",
			ClientIP: "127.0.0.1",
		}
		assert.Equal(t, "testuser", req.Username)
		assert.Equal(t, "password123", req.Password)
		assert.Equal(t, "127.0.0.1", req.ClientIP)
	})

	t.Run("UpdatePasswordRequest结构验证", func(t *testing.T) {
		userID := primitive.NewObjectID().Hex()
		req := &user.UpdatePasswordRequest{
			ID:          userID,
			OldPassword: "oldpassword",
			NewPassword: "newpassword",
		}
		assert.Equal(t, userID, req.ID)
		assert.Equal(t, "oldpassword", req.OldPassword)
		assert.Equal(t, "newpassword", req.NewPassword)
	})

	t.Run("AssignRoleRequest结构验证", func(t *testing.T) {
		req := &user.AssignRoleRequest{
			UserID: "user123",
			RoleID: "role123",
		}
		assert.Equal(t, "user123", req.UserID)
		assert.Equal(t, "role123", req.RoleID)
	})

	t.Run("DowngradeRoleRequest结构验证", func(t *testing.T) {
		req := &user.DowngradeRoleRequest{
			UserID:     "user123",
			TargetRole: "reader",
			Confirm:    true,
		}
		assert.Equal(t, "user123", req.UserID)
		assert.Equal(t, "reader", req.TargetRole)
		assert.True(t, req.Confirm)
	})
}

// TestUserServicePort_ResponseValidation 响应DTO验证测试
func TestUserServicePort_ResponseValidation(t *testing.T) {
	t.Run("CreateUserResponse验证", func(t *testing.T) {
		resp := &user.CreateUserResponse{
			User: &dto.UserDTO{
				ID:       primitive.NewObjectID().Hex(),
				Username: "testuser",
				Email:    "test@example.com",
			},
		}
		require.NotNil(t, resp.User)
		assert.NotEmpty(t, resp.User.ID)
		assert.Equal(t, "testuser", resp.User.Username)
		assert.Equal(t, "test@example.com", resp.User.Email)
	})

	t.Run("DeleteUserResponse验证", func(t *testing.T) {
		resp := &user.DeleteUserResponse{
			Deleted:   true,
			DeletedAt: time.Now(),
		}
		assert.True(t, resp.Deleted)
		assert.False(t, resp.DeletedAt.IsZero())
	})

	t.Run("UpdatePasswordResponse验证", func(t *testing.T) {
		resp := &user.UpdatePasswordResponse{
			Updated: true,
		}
		assert.True(t, resp.Updated)
	})

	t.Run("AssignRoleResponse验证", func(t *testing.T) {
		resp := &user.AssignRoleResponse{
			Assigned: true,
		}
		assert.True(t, resp.Assigned)
	})

	t.Run("GetUserRolesResponse验证", func(t *testing.T) {
		resp := &user.GetUserRolesResponse{
			Roles: []string{"reader", "author"},
		}
		assert.Len(t, resp.Roles, 2)
		assert.Contains(t, resp.Roles, "reader")
		assert.Contains(t, resp.Roles, "author")
	})

	t.Run("GetUserPermissionsResponse验证", func(t *testing.T) {
		resp := &user.GetUserPermissionsResponse{
			Permissions: []string{"read:books", "write:books"},
		}
		assert.Len(t, resp.Permissions, 2)
		assert.Contains(t, resp.Permissions, "read:books")
		assert.Contains(t, resp.Permissions, "write:books")
	})
}

// TestUserServicePort_EmailVerification 邮箱验证相关DTO验证
func TestUserServicePort_EmailVerification(t *testing.T) {
	t.Run("SendEmailVerificationRequest验证", func(t *testing.T) {
		req := &user.SendEmailVerificationRequest{
			UserID: primitive.NewObjectID().Hex(),
			Email:  "test@example.com",
		}
		assert.NotEmpty(t, req.UserID)
		assert.Equal(t, "test@example.com", req.Email)
	})

	t.Run("SendEmailVerificationResponse验证", func(t *testing.T) {
		resp := &user.SendEmailVerificationResponse{
			Success:   true,
			Message:   "验证码已发送",
			ExpiresIn: 1800,
		}
		assert.True(t, resp.Success)
		assert.Equal(t, "验证码已发送", resp.Message)
		assert.Equal(t, 1800, resp.ExpiresIn)
	})

	t.Run("VerifyEmailRequest验证", func(t *testing.T) {
		req := &user.VerifyEmailRequest{
			UserID: primitive.NewObjectID().Hex(),
			Code:   "123456",
		}
		assert.NotEmpty(t, req.UserID)
		assert.Equal(t, "123456", req.Code)
	})

	t.Run("VerifyEmailResponse验证", func(t *testing.T) {
		resp := &user.VerifyEmailResponse{
			Success: true,
			Message: "邮箱验证成功",
		}
		assert.True(t, resp.Success)
		assert.Equal(t, "邮箱验证成功", resp.Message)
	})
}

// TestUserServicePort_PasswordReset 密码重置相关DTO验证
func TestUserServicePort_PasswordReset(t *testing.T) {
	t.Run("RequestPasswordResetRequest验证", func(t *testing.T) {
		req := &user.RequestPasswordResetRequest{
			Email: "test@example.com",
		}
		assert.Equal(t, "test@example.com", req.Email)
	})

	t.Run("RequestPasswordResetResponse验证", func(t *testing.T) {
		resp := &user.RequestPasswordResetResponse{
			Success:   true,
			Message:   "重置邮件已发送",
			ExpiresIn: 3600,
		}
		assert.True(t, resp.Success)
		assert.Equal(t, "重置邮件已发送", resp.Message)
		assert.Equal(t, 3600, resp.ExpiresIn)
	})

	t.Run("ConfirmPasswordResetRequest验证", func(t *testing.T) {
		req := &user.ConfirmPasswordResetRequest{
			Email:    "test@example.com",
			Token:    "reset-token",
			Password: "newpassword",
		}
		assert.Equal(t, "test@example.com", req.Email)
		assert.Equal(t, "reset-token", req.Token)
		assert.Equal(t, "newpassword", req.Password)
	})

	t.Run("ConfirmPasswordResetResponse验证", func(t *testing.T) {
		resp := &user.ConfirmPasswordResetResponse{
			Success: true,
			Message: "密码重置成功",
		}
		assert.True(t, resp.Success)
		assert.Equal(t, "密码重置成功", resp.Message)
	})
}
