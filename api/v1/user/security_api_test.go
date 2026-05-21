package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	userServiceInterface "Qingyu_backend/service/interfaces/user"
)

type apiResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
	Errors  map[string]string      `json:"errors"`
}

func setupSecurityCompatRouter(api *SecurityAPI) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/api/v1/user/email/send-code", api.SendEmailVerification)
	router.POST("/api/v1/user/email/verify", api.VerifyEmail)
	router.POST("/api/v1/user/password/reset-request", api.RequestPasswordReset)
	router.POST("/api/v1/user/password/reset", api.ConfirmPasswordReset)

	return router
}

func performJSONRequest(t *testing.T, router *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()

	payload, err := json.Marshal(body)
	require.NoError(t, err)

	req, err := http.NewRequest(method, path, bytes.NewBuffer(payload))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func decodeAPIResponse(t *testing.T, body []byte) apiResponse {
	t.Helper()

	var resp apiResponse
	require.NoError(t, json.Unmarshal(body, &resp))
	return resp
}

func TestSecurityAPI_CompatRoutes_Success(t *testing.T) {
	mockUserService := new(MockVerificationUserService)
	api := NewSecurityAPI(mockUserService)
	router := setupSecurityCompatRouter(api)

	t.Run("send_email_verification", func(t *testing.T) {
		userID := primitive.NewObjectID().Hex()
		email := "register@example.com"
		expectedResp := &userServiceInterface.SendEmailVerificationResponse{
			Success:   true,
			Message:   "验证码已发送到您的邮箱",
			ExpiresIn: 300,
		}

		mockUserService.On("SendEmailVerification", mock.Anything, mock.MatchedBy(func(req *userServiceInterface.SendEmailVerificationRequest) bool {
			return req.UserID == userID && req.Email == email
		})).Return(expectedResp, nil).Once()

		w := performJSONRequest(t, router, http.MethodPost, "/api/v1/user/email/send-code", map[string]string{
			"user_id": userID,
			"email":   email,
		})

		require.Equal(t, http.StatusOK, w.Code)

		resp := decodeAPIResponse(t, w.Body.Bytes())
		assert.Equal(t, 0, resp.Code)
		assert.Equal(t, "验证码已发送到您的邮箱", resp.Message)
		require.Contains(t, resp.Data, "expires_in")
		assert.Equal(t, float64(300), resp.Data["expires_in"])
		require.Contains(t, resp.Data, "message")
		assert.Equal(t, "验证码已发送到您的邮箱", resp.Data["message"])
		mockUserService.AssertExpectations(t)
	})

	t.Run("request_password_reset", func(t *testing.T) {
		email := "reset@example.com"
		expectedResp := &userServiceInterface.RequestPasswordResetResponse{
			Success:   true,
			Message:   "如果该邮箱已注册，重置链接已发送",
			ExpiresIn: 900,
		}

		mockUserService.On("RequestPasswordReset", mock.Anything, mock.MatchedBy(func(req *userServiceInterface.RequestPasswordResetRequest) bool {
			return req.Email == email
		})).Return(expectedResp, nil).Once()

		w := performJSONRequest(t, router, http.MethodPost, "/api/v1/user/password/reset-request", map[string]string{
			"email": email,
		})

		require.Equal(t, http.StatusOK, w.Code)

		resp := decodeAPIResponse(t, w.Body.Bytes())
		assert.Equal(t, 0, resp.Code)
		assert.Equal(t, "如果该邮箱已注册，重置链接已发送", resp.Message)
		require.Contains(t, resp.Data, "expires_in")
		assert.Equal(t, float64(900), resp.Data["expires_in"])
		mockUserService.AssertExpectations(t)
	})
}

func TestSecurityAPI_CompatRoutes_Negatives(t *testing.T) {
	mockUserService := new(MockVerificationUserService)
	api := NewSecurityAPI(mockUserService)
	router := setupSecurityCompatRouter(api)

	t.Run("send_email_verification_invalid_email", func(t *testing.T) {
		w := performJSONRequest(t, router, http.MethodPost, "/api/v1/user/email/send-code", map[string]string{
			"user_id": primitive.NewObjectID().Hex(),
			"email":   "not-an-email",
		})

		require.Equal(t, http.StatusBadRequest, w.Code)

		resp := decodeAPIResponse(t, w.Body.Bytes())
		assert.Equal(t, 400, resp.Code)
		assert.Equal(t, "请求参数验证失败", resp.Message)
		require.Contains(t, resp.Errors, "email")
	})

	t.Run("verify_email_service_error", func(t *testing.T) {
		userID := primitive.NewObjectID().Hex()
		code := "123456"

		mockUserService.On("VerifyEmail", mock.Anything, mock.MatchedBy(func(req *userServiceInterface.VerifyEmailRequest) bool {
			return req.UserID == userID && req.Code == code
		})).Return(nil, errors.New("验证码无效或已过期")).Once()

		w := performJSONRequest(t, router, http.MethodPost, "/api/v1/user/email/verify", map[string]string{
			"user_id": userID,
			"code":    code,
		})

		require.Equal(t, http.StatusBadRequest, w.Code)

		resp := decodeAPIResponse(t, w.Body.Bytes())
		assert.Equal(t, 1001, resp.Code)
		assert.Equal(t, "验证失败: 验证码无效或已过期", resp.Message)
		mockUserService.AssertExpectations(t)
	})

	t.Run("confirm_password_reset_service_error", func(t *testing.T) {
		email := "reset@example.com"
		token := "invalid-token"
		password := "NewPassword456!"

		mockUserService.On("ConfirmPasswordReset", mock.Anything, mock.MatchedBy(func(req *userServiceInterface.ConfirmPasswordResetRequest) bool {
			return req.Email == email && req.Token == token && req.Password == password
		})).Return(nil, errors.New("Token验证失败")).Once()

		w := performJSONRequest(t, router, http.MethodPost, "/api/v1/user/password/reset", map[string]string{
			"email":    email,
			"token":    token,
			"password": password,
		})

		require.Equal(t, http.StatusBadRequest, w.Code)

		resp := decodeAPIResponse(t, w.Body.Bytes())
		assert.Equal(t, 1001, resp.Code)
		assert.Equal(t, "重置失败: Token验证失败", resp.Message)
		mockUserService.AssertExpectations(t)
	})
}
