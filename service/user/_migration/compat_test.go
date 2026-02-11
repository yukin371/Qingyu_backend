package user

import (
	"context"
	"testing"

	useriface "Qingyu_backend/service/interfaces/user"
	"Qingyu_backend/models/dto"
)

// MockPort 实现 Port 接口的 Mock
type MockPort struct {
	createUserCalled    bool
	getUserCalled       bool
	updateUserCalled    bool
	deleteUserCalled    bool
	listUsersCalled     bool
	registerUserCalled  bool
	loginUserCalled     bool
	logoutUserCalled    bool
	validateTokenCalled bool
	updatePasswordCalled bool
	resetPasswordCalled  bool
	verifyPasswordCalled bool
	requestPasswordResetCalled bool
	confirmPasswordResetCalled bool
	sendEmailVerificationCalled bool
	verifyEmailCalled      bool
	unbindEmailCalled      bool
	emailExistsCalled      bool
	assignRoleCalled       bool
	removeRoleCalled       bool
	getUserRolesCalled     bool
	getUserPermissionsCalled bool
	downgradeRoleCalled    bool
	updateLastLoginCalled  bool
	deleteDeviceCalled     bool
	unbindPhoneCalled      bool
	initializeCalled       bool
	healthCalled           bool
	closeCalled            bool
	getServiceNameCalled   bool
	getVersionCalled       bool
}

func (m *MockPort) Initialize(ctx context.Context) error {
	m.initializeCalled = true
	return nil
}

func (m *MockPort) Health(ctx context.Context) error {
	m.healthCalled = true
	return nil
}

func (m *MockPort) Close(ctx context.Context) error {
	m.closeCalled = true
	return nil
}

func (m *MockPort) GetServiceName() string {
	m.getServiceNameCalled = true
	return "MockUserService"
}

func (m *MockPort) GetVersion() string {
	m.getVersionCalled = true
	return "1.0.0"
}

func (m *MockPort) CreateUser(ctx context.Context, req *useriface.CreateUserRequest) (*useriface.CreateUserResponse, error) {
	m.createUserCalled = true
	return &useriface.CreateUserResponse{User: &dto.UserDTO{ID: "test-id"}}, nil
}

func (m *MockPort) GetUser(ctx context.Context, req *useriface.GetUserRequest) (*useriface.GetUserResponse, error) {
	m.getUserCalled = true
	return &useriface.GetUserResponse{User: &dto.UserDTO{ID: "test-id"}}, nil
}

func (m *MockPort) UpdateUser(ctx context.Context, req *useriface.UpdateUserRequest) (*useriface.UpdateUserResponse, error) {
	m.updateUserCalled = true
	return &useriface.UpdateUserResponse{User: &dto.UserDTO{ID: "test-id"}}, nil
}

func (m *MockPort) DeleteUser(ctx context.Context, req *useriface.DeleteUserRequest) (*useriface.DeleteUserResponse, error) {
	m.deleteUserCalled = true
	return &useriface.DeleteUserResponse{Deleted: true}, nil
}

func (m *MockPort) ListUsers(ctx context.Context, req *useriface.ListUsersRequest) (*useriface.ListUsersResponse, error) {
	m.listUsersCalled = true
	return &useriface.ListUsersResponse{Users: []*dto.UserDTO{}, Total: 0}, nil
}

func (m *MockPort) RegisterUser(ctx context.Context, req *useriface.RegisterUserRequest) (*useriface.RegisterUserResponse, error) {
	m.registerUserCalled = true
	return &useriface.RegisterUserResponse{User: &dto.UserDTO{ID: "test-id"}}, nil
}

func (m *MockPort) LoginUser(ctx context.Context, req *useriface.LoginUserRequest) (*useriface.LoginUserResponse, error) {
	m.loginUserCalled = true
	return &useriface.LoginUserResponse{User: &dto.UserDTO{ID: "test-id"}, Token: "test-token"}, nil
}

func (m *MockPort) LogoutUser(ctx context.Context, req *useriface.LogoutUserRequest) (*useriface.LogoutUserResponse, error) {
	m.logoutUserCalled = true
	return &useriface.LogoutUserResponse{Success: true}, nil
}

func (m *MockPort) ValidateToken(ctx context.Context, req *useriface.ValidateTokenRequest) (*useriface.ValidateTokenResponse, error) {
	m.validateTokenCalled = true
	return &useriface.ValidateTokenResponse{Valid: true, User: &dto.UserDTO{ID: "test-id"}}, nil
}

func (m *MockPort) UpdatePassword(ctx context.Context, req *useriface.UpdatePasswordRequest) (*useriface.UpdatePasswordResponse, error) {
	m.updatePasswordCalled = true
	return &useriface.UpdatePasswordResponse{Updated: true}, nil
}

func (m *MockPort) ResetPassword(ctx context.Context, req *useriface.ResetPasswordRequest) (*useriface.ResetPasswordResponse, error) {
	m.resetPasswordCalled = true
	return &useriface.ResetPasswordResponse{Success: true}, nil
}

func (m *MockPort) VerifyPassword(ctx context.Context, userID string, password string) error {
	m.verifyPasswordCalled = true
	return nil
}

func (m *MockPort) RequestPasswordReset(ctx context.Context, req *useriface.RequestPasswordResetRequest) (*useriface.RequestPasswordResetResponse, error) {
	m.requestPasswordResetCalled = true
	return &useriface.RequestPasswordResetResponse{Success: true}, nil
}

func (m *MockPort) ConfirmPasswordReset(ctx context.Context, req *useriface.ConfirmPasswordResetRequest) (*useriface.ConfirmPasswordResetResponse, error) {
	m.confirmPasswordResetCalled = true
	return &useriface.ConfirmPasswordResetResponse{Success: true}, nil
}

func (m *MockPort) SendEmailVerification(ctx context.Context, req *useriface.SendEmailVerificationRequest) (*useriface.SendEmailVerificationResponse, error) {
	m.sendEmailVerificationCalled = true
	return &useriface.SendEmailVerificationResponse{Success: true}, nil
}

func (m *MockPort) VerifyEmail(ctx context.Context, req *useriface.VerifyEmailRequest) (*useriface.VerifyEmailResponse, error) {
	m.verifyEmailCalled = true
	return &useriface.VerifyEmailResponse{Success: true}, nil
}

func (m *MockPort) UnbindEmail(ctx context.Context, userID string) error {
	m.unbindEmailCalled = true
	return nil
}

func (m *MockPort) EmailExists(ctx context.Context, email string) (bool, error) {
	m.emailExistsCalled = true
	return false, nil
}

func (m *MockPort) AssignRole(ctx context.Context, req *useriface.AssignRoleRequest) (*useriface.AssignRoleResponse, error) {
	m.assignRoleCalled = true
	return &useriface.AssignRoleResponse{Assigned: true}, nil
}

func (m *MockPort) RemoveRole(ctx context.Context, req *useriface.RemoveRoleRequest) (*useriface.RemoveRoleResponse, error) {
	m.removeRoleCalled = true
	return &useriface.RemoveRoleResponse{Removed: true}, nil
}

func (m *MockPort) GetUserRoles(ctx context.Context, req *useriface.GetUserRolesRequest) (*useriface.GetUserRolesResponse, error) {
	m.getUserRolesCalled = true
	return &useriface.GetUserRolesResponse{Roles: []string{"reader"}}, nil
}

func (m *MockPort) GetUserPermissions(ctx context.Context, req *useriface.GetUserPermissionsRequest) (*useriface.GetUserPermissionsResponse, error) {
	m.getUserPermissionsCalled = true
	return &useriface.GetUserPermissionsResponse{Permissions: []string{"read"}}, nil
}

func (m *MockPort) DowngradeRole(ctx context.Context, req *useriface.DowngradeRoleRequest) (*useriface.DowngradeRoleResponse, error) {
	m.downgradeRoleCalled = true
	return &useriface.DowngradeRoleResponse{CurrentRoles: []string{"reader"}}, nil
}

func (m *MockPort) UpdateLastLogin(ctx context.Context, req *useriface.UpdateLastLoginRequest) (*useriface.UpdateLastLoginResponse, error) {
	m.updateLastLoginCalled = true
	return &useriface.UpdateLastLoginResponse{Updated: true}, nil
}

func (m *MockPort) DeleteDevice(ctx context.Context, userID string, deviceID string) error {
	m.deleteDeviceCalled = true
	return nil
}

func (m *MockPort) UnbindPhone(ctx context.Context, userID string) error {
	m.unbindPhoneCalled = true
	return nil
}

// TestUserServiceAdapter_ImplementsUserServiceInterface 测试适配器实现了 UserService 接口
func TestUserServiceAdapter_ImplementsUserServiceInterface(t *testing.T) {
	ctx := context.Background()
	mockPort := &MockPort{}

	// 创建适配器
	adapter := NewUserServiceAdapter(mockPort, mockPort, mockPort, mockPort, mockPort, mockPort)

	// 测试 UserManagementPort 方法
	t.Run("UserManagementPort methods", func(t *testing.T) {
		_, err := adapter.CreateUser(ctx, &useriface.CreateUserRequest{})
		if err != nil || !mockPort.createUserCalled {
			t.Errorf("CreateUser failed or not called")
		}

		_, err = adapter.GetUser(ctx, &useriface.GetUserRequest{})
		if err != nil || !mockPort.getUserCalled {
			t.Errorf("GetUser failed or not called")
		}

		_, err = adapter.UpdateUser(ctx, &useriface.UpdateUserRequest{})
		if err != nil || !mockPort.updateUserCalled {
			t.Errorf("UpdateUser failed or not called")
		}

		_, err = adapter.DeleteUser(ctx, &useriface.DeleteUserRequest{})
		if err != nil || !mockPort.deleteUserCalled {
			t.Errorf("DeleteUser failed or not called")
		}

		_, err = adapter.ListUsers(ctx, &useriface.ListUsersRequest{})
		if err != nil || !mockPort.listUsersCalled {
			t.Errorf("ListUsers failed or not called")
		}
	})

	// 测试 UserAuthPort 方法
	t.Run("UserAuthPort methods", func(t *testing.T) {
		_, err := adapter.RegisterUser(ctx, &useriface.RegisterUserRequest{})
		if err != nil || !mockPort.registerUserCalled {
			t.Errorf("RegisterUser failed or not called")
		}

		_, err = adapter.LoginUser(ctx, &useriface.LoginUserRequest{})
		if err != nil || !mockPort.loginUserCalled {
			t.Errorf("LoginUser failed or not called")
		}

		_, err = adapter.LogoutUser(ctx, &useriface.LogoutUserRequest{})
		if err != nil || !mockPort.logoutUserCalled {
			t.Errorf("LogoutUser failed or not called")
		}

		_, err = adapter.ValidateToken(ctx, &useriface.ValidateTokenRequest{})
		if err != nil || !mockPort.validateTokenCalled {
			t.Errorf("ValidateToken failed or not called")
		}
	})

	// 测试 PasswordManagementPort 方法
	t.Run("PasswordManagementPort methods", func(t *testing.T) {
		_, err := adapter.UpdatePassword(ctx, &useriface.UpdatePasswordRequest{})
		if err != nil || !mockPort.updatePasswordCalled {
			t.Errorf("UpdatePassword failed or not called")
		}

		_, err = adapter.ResetPassword(ctx, &useriface.ResetPasswordRequest{})
		if err != nil || !mockPort.resetPasswordCalled {
			t.Errorf("ResetPassword failed or not called")
		}

		err = adapter.VerifyPassword(ctx, "user-id", "password")
		if err != nil || !mockPort.verifyPasswordCalled {
			t.Errorf("VerifyPassword failed or not called")
		}

		_, err = adapter.RequestPasswordReset(ctx, &useriface.RequestPasswordResetRequest{})
		if err != nil || !mockPort.requestPasswordResetCalled {
			t.Errorf("RequestPasswordReset failed or not called")
		}

		_, err = adapter.ConfirmPasswordReset(ctx, &useriface.ConfirmPasswordResetRequest{})
		if err != nil || !mockPort.confirmPasswordResetCalled {
			t.Errorf("ConfirmPasswordReset failed or not called")
		}
	})

	// 测试 EmailManagementPort 方法
	t.Run("EmailManagementPort methods", func(t *testing.T) {
		_, err := adapter.SendEmailVerification(ctx, &useriface.SendEmailVerificationRequest{})
		if err != nil || !mockPort.sendEmailVerificationCalled {
			t.Errorf("SendEmailVerification failed or not called")
		}

		_, err = adapter.VerifyEmail(ctx, &useriface.VerifyEmailRequest{})
		if err != nil || !mockPort.verifyEmailCalled {
			t.Errorf("VerifyEmail failed or not called")
		}

		err = adapter.UnbindEmail(ctx, "user-id")
		if err != nil || !mockPort.unbindEmailCalled {
			t.Errorf("UnbindEmail failed or not called")
		}

		_, err = adapter.EmailExists(ctx, "test@example.com")
		if err != nil || !mockPort.emailExistsCalled {
			t.Errorf("EmailExists failed or not called")
		}
	})

	// 测试 UserPermissionPort 方法
	t.Run("UserPermissionPort methods", func(t *testing.T) {
		_, err := adapter.AssignRole(ctx, &useriface.AssignRoleRequest{})
		if err != nil || !mockPort.assignRoleCalled {
			t.Errorf("AssignRole failed or not called")
		}

		_, err = adapter.RemoveRole(ctx, &useriface.RemoveRoleRequest{})
		if err != nil || !mockPort.removeRoleCalled {
			t.Errorf("RemoveRole failed or not called")
		}

		_, err = adapter.GetUserRoles(ctx, &useriface.GetUserRolesRequest{})
		if err != nil || !mockPort.getUserRolesCalled {
			t.Errorf("GetUserRoles failed or not called")
		}

		_, err = adapter.GetUserPermissions(ctx, &useriface.GetUserPermissionsRequest{})
		if err != nil || !mockPort.getUserPermissionsCalled {
			t.Errorf("GetUserPermissions failed or not called")
		}

		_, err = adapter.DowngradeRole(ctx, &useriface.DowngradeRoleRequest{})
		if err != nil || !mockPort.downgradeRoleCalled {
			t.Errorf("DowngradeRole failed or not called")
		}
	})

	// 测试 UserStatusPort 方法
	t.Run("UserStatusPort methods", func(t *testing.T) {
		_, err := adapter.UpdateLastLogin(ctx, &useriface.UpdateLastLoginRequest{})
		if err != nil || !mockPort.updateLastLoginCalled {
			t.Errorf("UpdateLastLogin failed or not called")
		}

		err = adapter.DeleteDevice(ctx, "user-id", "device-id")
		if err != nil || !mockPort.deleteDeviceCalled {
			t.Errorf("DeleteDevice failed or not called")
		}

		err = adapter.UnbindPhone(ctx, "user-id")
		if err != nil || !mockPort.unbindPhoneCalled {
			t.Errorf("UnbindPhone failed or not called")
		}
	})
}

// TestNewUserServiceAdapter 测试适配器构造函数
func TestNewUserServiceAdapter(t *testing.T) {
	mockPort := &MockPort{}
	adapter := NewUserServiceAdapter(mockPort, mockPort, mockPort, mockPort, mockPort, mockPort)

	if adapter == nil {
		t.Fatal("NewUserServiceAdapter returned nil")
	}

	if adapter.managementPort != mockPort {
		t.Error("managementPort not set correctly")
	}

	if adapter.authPort != mockPort {
		t.Error("authPort not set correctly")
	}

	if adapter.passwordPort != mockPort {
		t.Error("passwordPort not set correctly")
	}

	if adapter.emailPort != mockPort {
		t.Error("emailPort not set correctly")
	}

	if adapter.permissionPort != mockPort {
		t.Error("permissionPort not set correctly")
	}

	if adapter.statusPort != mockPort {
		t.Error("statusPort not set correctly")
	}
}
