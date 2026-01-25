package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	serviceInterfaces "Qingyu_backend/service/interfaces/base"
	userServiceInterface "Qingyu_backend/service/interfaces/user"
	"Qingyu_backend/models/dto"
)

// MockUserService 模拟用户服务
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Initialize(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockUserService) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockUserService) Close(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockUserService) GetServiceName() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockUserService) GetVersion() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockUserService) CreateUser(ctx context.Context, req *userServiceInterface.CreateUserRequest) (*userServiceInterface.CreateUserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userServiceInterface.CreateUserResponse), args.Error(1)
}

func (m *MockUserService) GetUser(ctx context.Context, req *userServiceInterface.GetUserRequest) (*userServiceInterface.GetUserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userServiceInterface.GetUserResponse), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, req *userServiceInterface.UpdateUserRequest) (*userServiceInterface.UpdateUserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userServiceInterface.UpdateUserResponse), args.Error(1)
}

func (m *MockUserService) DeleteUser(ctx context.Context, req *userServiceInterface.DeleteUserRequest) (*userServiceInterface.DeleteUserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userServiceInterface.DeleteUserResponse), args.Error(1)
}

func (m *MockUserService) ListUsers(ctx context.Context, req *userServiceInterface.ListUsersRequest) (*userServiceInterface.ListUsersResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userServiceInterface.ListUsersResponse), args.Error(1)
}

func (m *MockUserService) RegisterUser(ctx context.Context, req *userServiceInterface.RegisterUserRequest) (*userServiceInterface.RegisterUserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userServiceInterface.RegisterUserResponse), args.Error(1)
}

func (m *MockUserService) LoginUser(ctx context.Context, req *userServiceInterface.LoginUserRequest) (*userServiceInterface.LoginUserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userServiceInterface.LoginUserResponse), args.Error(1)
}

func (m *MockUserService) LogoutUser(ctx context.Context, req *userServiceInterface.LogoutUserRequest) (*userServiceInterface.LogoutUserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userServiceInterface.LogoutUserResponse), args.Error(1)
}

func (m *MockUserService) ValidateToken(ctx context.Context, req *userServiceInterface.ValidateTokenRequest) (*userServiceInterface.ValidateTokenResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userServiceInterface.ValidateTokenResponse), args.Error(1)
}

func (m *MockUserService) UpdateLastLogin(ctx context.Context, req *userServiceInterface.UpdateLastLoginRequest) (*userServiceInterface.UpdateLastLoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userServiceInterface.UpdateLastLoginResponse), args.Error(1)
}

func (m *MockUserService) UpdatePassword(ctx context.Context, req *userServiceInterface.UpdatePasswordRequest) (*userServiceInterface.UpdatePasswordResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userServiceInterface.UpdatePasswordResponse), args.Error(1)
}

func (m *MockUserService) ResetPassword(ctx context.Context, req *userServiceInterface.ResetPasswordRequest) (*userServiceInterface.ResetPasswordResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userServiceInterface.ResetPasswordResponse), args.Error(1)
}

func (m *MockUserService) SendEmailVerification(ctx context.Context, req *userServiceInterface.SendEmailVerificationRequest) (*userServiceInterface.SendEmailVerificationResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userServiceInterface.SendEmailVerificationResponse), args.Error(1)
}

func (m *MockUserService) VerifyEmail(ctx context.Context, req *userServiceInterface.VerifyEmailRequest) (*userServiceInterface.VerifyEmailResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userServiceInterface.VerifyEmailResponse), args.Error(1)
}

func (m *MockUserService) UnbindEmail(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserService) UnbindPhone(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserService) DeleteDevice(ctx context.Context, userID string, deviceID string) error {
	args := m.Called(ctx, userID, deviceID)
	return args.Error(0)
}

func (m *MockUserService) VerifyPassword(ctx context.Context, userID string, password string) error {
	args := m.Called(ctx, userID, password)
	return args.Error(0)
}

func (m *MockUserService) EmailExists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserService) RequestPasswordReset(ctx context.Context, req *userServiceInterface.RequestPasswordResetRequest) (*userServiceInterface.RequestPasswordResetResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userServiceInterface.RequestPasswordResetResponse), args.Error(1)
}

func (m *MockUserService) ConfirmPasswordReset(ctx context.Context, req *userServiceInterface.ConfirmPasswordResetRequest) (*userServiceInterface.ConfirmPasswordResetResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userServiceInterface.ConfirmPasswordResetResponse), args.Error(1)
}

func (m *MockUserService) AssignRole(ctx context.Context, req *userServiceInterface.AssignRoleRequest) (*userServiceInterface.AssignRoleResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userServiceInterface.AssignRoleResponse), args.Error(1)
}

func (m *MockUserService) RemoveRole(ctx context.Context, req *userServiceInterface.RemoveRoleRequest) (*userServiceInterface.RemoveRoleResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userServiceInterface.RemoveRoleResponse), args.Error(1)
}

func (m *MockUserService) GetUserRoles(ctx context.Context, req *userServiceInterface.GetUserRolesRequest) (*userServiceInterface.GetUserRolesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userServiceInterface.GetUserRolesResponse), args.Error(1)
}

func (m *MockUserService) GetUserPermissions(ctx context.Context, req *userServiceInterface.GetUserPermissionsRequest) (*userServiceInterface.GetUserPermissionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userServiceInterface.GetUserPermissionsResponse), args.Error(1)
}

// setupTestRouter 设置测试路由
func setupTestRouter(userService userServiceInterface.UserService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	handler := NewProfileHandler(userService)

	// 添加认证中间件（简化版本，直接设置user_id）
	router.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user-id")
		c.Next()
	})

	router.GET("/profile", handler.GetProfile)
	router.PUT("/profile", handler.UpdateProfile)
	router.PUT("/password", handler.UpdatePassword)
	router.POST("/avatar", handler.UploadAvatar)

	return router
}

func TestProfileHandler_GetProfile(t *testing.T) {
	mockService := new(MockUserService)
	router := setupTestRouter(mockService)

	// 测试用例1：正常获取用户信息
	t.Run("Success", func(t *testing.T) {
		expectedUser := &dto.UserDTO{
			ID:        "test-user-id",
			Username:  "testuser",
			Email:     "test@example.com",
			Nickname:  "Test User",
			Bio:       "Test bio",
			Gender:    "male",
			Location:  "Beijing",
			Website:   "https://example.com",
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		}

		mockService.On("GetUser", mock.Anything, mock.AnythingOfType("*user.GetUserRequest")).Return(
			&userServiceInterface.GetUserResponse{User: expectedUser}, nil,
		)

		req, _ := http.NewRequest("GET", "/profile", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	// 测试用例2：用户不存在
	t.Run("UserNotFound", func(t *testing.T) {
		mockService.On("GetUser", mock.Anything, mock.AnythingOfType("*user.GetUserRequest")).Return(
			nil, serviceInterfaces.NewServiceError("UserService", serviceInterfaces.ErrorTypeNotFound, "用户不存在", nil),
		)

		req, _ := http.NewRequest("GET", "/profile", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestProfileHandler_UpdateProfile(t *testing.T) {
	mockService := new(MockUserService)
	router := setupTestRouter(mockService)

	// 测试用例1：正常更新
	t.Run("Success", func(t *testing.T) {
		updatedUser := &dto.UserDTO{
			ID:        "test-user-id",
			Username:  "testuser",
			Nickname:  "Updated Nickname",
			Bio:       "Updated bio",
			Gender:    "female",
			UpdatedAt: time.Now().Format(time.RFC3339),
		}

		mockService.On("UpdateUser", mock.Anything, mock.AnythingOfType("*user.UpdateUserRequest")).Return(
			&userServiceInterface.UpdateUserResponse{User: updatedUser}, nil,
		)

		body := map[string]interface{}{
			"nickname": "Updated Nickname",
			"bio":      "Updated bio",
			"gender":   "female",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("PUT", "/profile", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	// 测试用例2：生日格式错误
	t.Run("InvalidBirthdayFormat", func(t *testing.T) {
		body := map[string]interface{}{
			"nickname": "Test",
			"birthday": "invalid-date",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("PUT", "/profile", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestProfileHandler_UpdatePassword(t *testing.T) {
	mockService := new(MockUserService)
	router := setupTestRouter(mockService)

	// 测试用例1：正常修改密码
	t.Run("Success", func(t *testing.T) {
		mockService.On("UpdatePassword", mock.Anything, mock.AnythingOfType("*user.UpdatePasswordRequest")).Return(
			&userServiceInterface.UpdatePasswordResponse{Updated: true}, nil,
		)

		body := map[string]string{
			"old_password": "OldPass123",
			"new_password": "NewPass456",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("PUT", "/password", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	// 测试用例2：旧密码错误
	t.Run("OldPasswordMismatch", func(t *testing.T) {
		mockService.On("UpdatePassword", mock.Anything, mock.AnythingOfType("*user.UpdatePasswordRequest")).Return(
			nil, serviceInterfaces.NewServiceError("UserService", serviceInterfaces.ErrorTypeUnauthorized, "旧密码错误", nil),
		)

		body := map[string]string{
			"old_password": "WrongOldPass",
			"new_password": "NewPass456",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("PUT", "/password", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestProfileHandler_UploadAvatar(t *testing.T) {
	mockService := new(MockUserService)
	router := setupTestRouter(mockService)

	// 测试用例1：没有设置存储服务
	t.Run("StorageServiceNotSet", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		part, _ := writer.CreateFormFile("avatar", "test.jpg")
		part.Write([]byte("fake image data"))

		writer.Close()

		req, _ := http.NewRequest("POST", "/avatar", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

// 创建Multipart表单文件的辅助函数
func createMultipartFormData(fieldName, filename string, content []byte) (*bytes.Buffer, string) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, _ := writer.CreateFormFile(fieldName, filename)
	part.Write(content)

	writer.Close()
	return body, writer.FormDataContentType()
}

// 示例：在设置完storage service后测试头像上传
func TestProfileHandler_UploadAvatar_WithStorage(t *testing.T) {
	// 这个测试需要mock StorageService
	// 由于当前实现中StorageService是可选依赖，这个测试暂时跳过
	t.Skip("需要实现StorageService的mock")
}

// 测试RFC3339日期格式解析
func TestParseBirthdayRFC3339(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		expectOK  bool
	}{
		{"Valid RFC3339", "1990-01-01T00:00:00Z", true},
		{"Valid RFC3339 with offset", "1990-01-01T00:00:00+08:00", true},
		{"Invalid format", "1990-01-01", false},
		{"Invalid format2", "01/01/1990", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := time.Parse(time.RFC3339, tc.input)
			if tc.expectOK {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// 辅助函数：创建测试用的UserDTO
func createTestUserDTO(id string) *dto.UserDTO {
	now := time.Now().Format(time.RFC3339)
	return &dto.UserDTO{
		ID:        id,
		Username:  fmt.Sprintf("user_%s", id),
		Email:     fmt.Sprintf("user_%s@example.com", id),
		Nickname:  fmt.Sprintf("Nickname_%s", id),
		Bio:       fmt.Sprintf("Bio for user %s", id),
		Gender:    "male",
		Location:  "Beijing",
		Website:   "https://example.com",
		CreatedAt: now,
		UpdatedAt: now,
	}
}
