package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	sharedAPI "Qingyu_backend/api/v1/shared"
	"Qingyu_backend/internal/middleware/builtin"
	authModel "Qingyu_backend/models/auth"
	authservice "Qingyu_backend/service/auth"
)

type MockOAuthService struct {
	mock.Mock
}

func (m *MockOAuthService) GetAuthURL(ctx context.Context, provider authModel.OAuthProvider, redirectURI, state string, linkMode bool, userID ...string) (string, error) {
	args := m.Called(ctx, provider, redirectURI, state, linkMode, userID)
	return args.String(0), args.Error(1)
}

func (m *MockOAuthService) ExchangeCode(ctx context.Context, provider authModel.OAuthProvider, code, state string) (*oauth2.Token, *authModel.OAuthSession, error) {
	args := m.Called(ctx, provider, code, state)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*authModel.OAuthSession), args.Error(2)
	}
	return args.Get(0).(*oauth2.Token), args.Get(1).(*authModel.OAuthSession), args.Error(2)
}

func (m *MockOAuthService) GetUserInfo(ctx context.Context, provider authModel.OAuthProvider, token *oauth2.Token) (*authModel.UserIdentity, error) {
	args := m.Called(ctx, provider, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authModel.UserIdentity), args.Error(1)
}

func (m *MockOAuthService) LinkAccount(ctx context.Context, userID string, provider authModel.OAuthProvider, token *oauth2.Token, identity *authModel.UserIdentity) (*authModel.OAuthAccount, error) {
	args := m.Called(ctx, userID, provider, token, identity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authModel.OAuthAccount), args.Error(1)
}

func (m *MockOAuthService) UnlinkAccount(ctx context.Context, userID, accountID string) error {
	args := m.Called(ctx, userID, accountID)
	return args.Error(0)
}

func (m *MockOAuthService) GetLinkedAccounts(ctx context.Context, userID string) ([]*authModel.OAuthAccount, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*authModel.OAuthAccount), args.Error(1)
}

func (m *MockOAuthService) SetPrimaryAccount(ctx context.Context, userID, accountID string) error {
	args := m.Called(ctx, userID, accountID)
	return args.Error(0)
}

func setupOAuthTestRouter(oauthService *MockOAuthService, authService authservice.AuthService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	logger := zap.NewNop()
	errorHandler := builtin.NewErrorHandlerMiddleware(logger)
	r.Use(errorHandler.Handler())

	api := sharedAPI.NewOAuthAPI(oauthService, authService, logger)

	v1 := r.Group("/api/v1/shared/oauth")
	{
		v1.POST("/:provider/authorize", api.GetAuthorizeURL)
		v1.POST("/:provider/callback", api.HandleCallback)
		v1.GET("/accounts", api.GetLinkedAccounts)
		v1.DELETE("/accounts/:accountID", api.UnlinkAccount)
		v1.PUT("/accounts/:accountID/primary", api.SetPrimaryAccount)
	}

	return r
}

func setupOAuthTestRouterWithAuth(oauthService *MockOAuthService, authService authservice.AuthService, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	logger := zap.NewNop()
	errorHandler := builtin.NewErrorHandlerMiddleware(logger)

	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		c.Next()
	}, errorHandler.Handler())

	api := sharedAPI.NewOAuthAPI(oauthService, authService, logger)

	v1 := r.Group("/api/v1/shared/oauth")
	{
		v1.POST("/:provider/authorize", api.GetAuthorizeURL)
		v1.POST("/:provider/callback", api.HandleCallback)
		v1.GET("/accounts", api.GetLinkedAccounts)
		v1.DELETE("/accounts/:accountID", api.UnlinkAccount)
		v1.PUT("/accounts/:accountID/primary", api.SetPrimaryAccount)
	}

	return r
}

func TestOAuthAPI_GetAuthorizeURL_Success(t *testing.T) {
	mockOAuthService := new(MockOAuthService)
	mockAuthService := new(MockAuthService)
	router := setupOAuthTestRouter(mockOAuthService, mockAuthService)

	reqBody := sharedAPI.OAuthAuthorizeRequest{
		RedirectURI: "http://localhost:3000/callback",
		State:       "random-state-123",
	}

	expectedURL := "https://accounts.google.com/o/oauth2/v2/auth?redirect_uri=http://localhost:3000/callback&state=random-state-123"

	mockOAuthService.On("GetAuthURL", mock.Anything, authModel.OAuthProviderGoogle, reqBody.RedirectURI, reqBody.State, false, mock.AnythingOfType("[]string")).Return(expectedURL, nil)

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/shared/oauth/google/authorize", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, expectedURL, data["authorize_url"])
	assert.Equal(t, "google", data["provider"])

	mockOAuthService.AssertExpectations(t)
}

func TestOAuthAPI_GetAuthorizeURL_MissingRedirectURI(t *testing.T) {
	mockOAuthService := new(MockOAuthService)
	mockAuthService := new(MockAuthService)
	router := setupOAuthTestRouter(mockOAuthService, mockAuthService)

	reqBody := map[string]string{
		"state": "random-state-123",
	}

	mockOAuthService.On("GetAuthURL", mock.Anything, authModel.OAuthProviderGoogle, "", "random-state-123", false, mock.AnythingOfType("[]string")).Return("", assert.AnError)

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/shared/oauth/google/authorize", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockOAuthService.AssertExpectations(t)
}

func TestOAuthAPI_GetAuthorizeURL_LinkMode(t *testing.T) {
	mockOAuthService := new(MockOAuthService)
	mockAuthService := new(MockAuthService)
	userID := "user123"
	router := setupOAuthTestRouterWithAuth(mockOAuthService, mockAuthService, userID)

	reqBody := sharedAPI.OAuthAuthorizeRequest{
		RedirectURI: "http://localhost:3000/callback",
		State:       "random-state-123",
	}

	expectedURL := "https://github.com/login/oauth/authorize?redirect_uri=http://localhost:3000/callback"

	mockOAuthService.On("GetAuthURL", mock.Anything, authModel.OAuthProviderGitHub, reqBody.RedirectURI, reqBody.State, true, []string{userID}).Return(expectedURL, nil)

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/shared/oauth/github/authorize", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockOAuthService.AssertExpectations(t)
}

func TestOAuthAPI_GetAuthorizeURL_ServiceError(t *testing.T) {
	mockOAuthService := new(MockOAuthService)
	mockAuthService := new(MockAuthService)
	router := setupOAuthTestRouter(mockOAuthService, mockAuthService)

	reqBody := sharedAPI.OAuthAuthorizeRequest{
		RedirectURI: "invalid://uri",
		State:       "random-state-123",
	}

	mockOAuthService.On("GetAuthURL", mock.Anything, authModel.OAuthProviderGoogle, reqBody.RedirectURI, reqBody.State, false, mock.AnythingOfType("[]string")).Return("", assert.AnError)

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/shared/oauth/google/authorize", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockOAuthService.AssertExpectations(t)
}

func TestOAuthAPI_GetAuthorizeURL_InvalidProvider(t *testing.T) {
	mockOAuthService := new(MockOAuthService)
	mockAuthService := new(MockAuthService)
	router := setupOAuthTestRouter(mockOAuthService, mockAuthService)

	reqBody := sharedAPI.OAuthAuthorizeRequest{
		RedirectURI: "http://localhost:3000/callback",
		State:       "random-state-123",
	}

	mockOAuthService.On("GetAuthURL", mock.Anything, authModel.OAuthProvider("invalid_provider"), reqBody.RedirectURI, reqBody.State, false, mock.AnythingOfType("[]string")).Return("", assert.AnError)

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/shared/oauth/invalid_provider/authorize", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockOAuthService.AssertExpectations(t)
}

func TestOAuthAPI_HandleCallback_LoginMode_Success(t *testing.T) {
	mockOAuthService := new(MockOAuthService)
	mockAuthService := new(MockAuthService)
	router := setupOAuthTestRouter(mockOAuthService, mockAuthService)

	reqBody := sharedAPI.OAuthCallbackRequest{
		Code:  "auth-code-123",
		State: "random-state-123",
	}

	token := &oauth2.Token{
		AccessToken: "access-token-123",
		TokenType:   "Bearer",
	}

	session := &authModel.OAuthSession{
		ID:        primitive.NewObjectID(),
		LinkMode:  false,
		CreatedAt: time.Now(),
	}

	identity := &authModel.UserIdentity{
		ProviderID: "google-id-123",
		Email:      "user@example.com",
		Name:       "Test User",
		Avatar:     "http://example.com/avatar.jpg",
		Username:   "testuser",
	}

	loginResp := &authservice.LoginResponse{
		User: &authservice.UserInfo{
			ID:       "user123",
			Username: "testuser",
			Email:    "user@example.com",
			Roles:    []string{"reader"},
		},
		Token: "jwt-token-123",
	}

	mockOAuthService.On("ExchangeCode", mock.Anything, authModel.OAuthProviderGoogle, reqBody.Code, reqBody.State).Return(token, session, nil)
	mockOAuthService.On("GetUserInfo", mock.Anything, authModel.OAuthProviderGoogle, token).Return(identity, nil)
	mockAuthService.On("OAuthLogin", mock.Anything, mock.MatchedBy(func(req *authservice.OAuthLoginRequest) bool {
		return req.Provider == authModel.OAuthProviderGoogle && req.ProviderID == "google-id-123"
	})).Return(loginResp, nil)

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/shared/oauth/google/callback", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "OAuth登录成功", response["message"])

	mockOAuthService.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
}

func TestOAuthAPI_HandleCallback_LinkMode_Success(t *testing.T) {
	mockOAuthService := new(MockOAuthService)
	mockAuthService := new(MockAuthService)
	userID := "user123"
	router := setupOAuthTestRouterWithAuth(mockOAuthService, mockAuthService, userID)

	reqBody := sharedAPI.OAuthCallbackRequest{
		Code:  "auth-code-123",
		State: "random-state-123",
	}

	token := &oauth2.Token{
		AccessToken: "access-token-123",
		TokenType:   "Bearer",
	}

	session := &authModel.OAuthSession{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		LinkMode:  true,
		CreatedAt: time.Now(),
	}

	identity := &authModel.UserIdentity{
		ProviderID: "github-id-123",
		Email:      "user@example.com",
		Name:       "Test User",
	}

	account := &authModel.OAuthAccount{
		ID:       primitive.NewObjectID(),
		UserID:   userID,
		Provider: authModel.OAuthProviderGitHub,
	}

	mockOAuthService.On("ExchangeCode", mock.Anything, authModel.OAuthProviderGitHub, reqBody.Code, reqBody.State).Return(token, session, nil)
	mockOAuthService.On("GetUserInfo", mock.Anything, authModel.OAuthProviderGitHub, token).Return(identity, nil)
	mockOAuthService.On("LinkAccount", mock.Anything, userID, authModel.OAuthProviderGitHub, token, identity).Return(account, nil)

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/shared/oauth/github/callback", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "账号绑定成功", response["message"])

	data := response["data"].(map[string]interface{})
	assert.NotNil(t, data["account"])

	mockOAuthService.AssertExpectations(t)
}

func TestOAuthAPI_HandleCallback_MissingCode(t *testing.T) {
	mockOAuthService := new(MockOAuthService)
	mockAuthService := new(MockAuthService)
	router := setupOAuthTestRouter(mockOAuthService, mockAuthService)

	reqBody := map[string]string{
		"state": "random-state-123",
	}

	mockOAuthService.On("ExchangeCode", mock.Anything, authModel.OAuthProviderGoogle, "", "random-state-123").Return(nil, &authModel.OAuthSession{}, assert.AnError)

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/shared/oauth/google/callback", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockOAuthService.AssertExpectations(t)
}

func TestOAuthAPI_HandleCallback_ExchangeCodeError(t *testing.T) {
	mockOAuthService := new(MockOAuthService)
	mockAuthService := new(MockAuthService)
	router := setupOAuthTestRouter(mockOAuthService, mockAuthService)

	reqBody := sharedAPI.OAuthCallbackRequest{
		Code:  "invalid-code",
		State: "random-state-123",
	}

	mockOAuthService.On("ExchangeCode", mock.Anything, authModel.OAuthProviderGoogle, reqBody.Code, reqBody.State).Return(nil, &authModel.OAuthSession{}, assert.AnError)

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/shared/oauth/google/callback", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockOAuthService.AssertExpectations(t)
}

func TestOAuthAPI_HandleCallback_GetUserInfoError(t *testing.T) {
	mockOAuthService := new(MockOAuthService)
	mockAuthService := new(MockAuthService)
	router := setupOAuthTestRouter(mockOAuthService, mockAuthService)

	reqBody := sharedAPI.OAuthCallbackRequest{
		Code:  "auth-code-123",
		State: "random-state-123",
	}

	token := &oauth2.Token{AccessToken: "access-token-123"}
	session := &authModel.OAuthSession{ID: primitive.NewObjectID()}

	mockOAuthService.On("ExchangeCode", mock.Anything, authModel.OAuthProviderGoogle, reqBody.Code, reqBody.State).Return(token, session, nil)
	mockOAuthService.On("GetUserInfo", mock.Anything, authModel.OAuthProviderGoogle, token).Return(nil, assert.AnError)

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/shared/oauth/google/callback", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockOAuthService.AssertExpectations(t)
}

func TestOAuthAPI_HandleCallback_OAuthLoginError(t *testing.T) {
	mockOAuthService := new(MockOAuthService)
	mockAuthService := new(MockAuthService)
	router := setupOAuthTestRouter(mockOAuthService, mockAuthService)

	reqBody := sharedAPI.OAuthCallbackRequest{
		Code:  "auth-code-123",
		State: "random-state-123",
	}

	token := &oauth2.Token{AccessToken: "access-token-123"}
	session := &authModel.OAuthSession{ID: primitive.NewObjectID()}
	identity := &authModel.UserIdentity{
		ProviderID: "google-id-123",
		Email:      "user@example.com",
		Name:       "Test User",
	}

	mockOAuthService.On("ExchangeCode", mock.Anything, authModel.OAuthProviderGoogle, reqBody.Code, reqBody.State).Return(token, session, nil)
	mockOAuthService.On("GetUserInfo", mock.Anything, authModel.OAuthProviderGoogle, token).Return(identity, nil)
	mockAuthService.On("OAuthLogin", mock.Anything, mock.AnythingOfType("*auth.OAuthLoginRequest")).Return(nil, assert.AnError)

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/shared/oauth/google/callback", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockOAuthService.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
}

func TestOAuthAPI_GetLinkedAccounts_Success(t *testing.T) {
	mockOAuthService := new(MockOAuthService)
	mockAuthService := new(MockAuthService)
	userID := "user123"
	router := setupOAuthTestRouterWithAuth(mockOAuthService, mockAuthService, userID)

	expectedAccounts := []*authModel.OAuthAccount{
		{
			ID:       primitive.NewObjectID(),
			UserID:   userID,
			Provider: authModel.OAuthProviderGoogle,
			Email:    "user@example.com",
		},
		{
			ID:       primitive.NewObjectID(),
			UserID:   userID,
			Provider: authModel.OAuthProviderGitHub,
			Username: "githubuser",
		},
	}

	mockOAuthService.On("GetLinkedAccounts", mock.Anything, userID).Return(expectedAccounts, nil)

	req, _ := http.NewRequest("GET", "/api/v1/shared/oauth/accounts", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "获取成功", response["message"])

	mockOAuthService.AssertExpectations(t)
}

func TestOAuthAPI_GetLinkedAccounts_Unauthorized(t *testing.T) {
	mockOAuthService := new(MockOAuthService)
	mockAuthService := new(MockAuthService)
	router := setupOAuthTestRouterWithAuth(mockOAuthService, mockAuthService, "")

	req, _ := http.NewRequest("GET", "/api/v1/shared/oauth/accounts", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestOAuthAPI_GetLinkedAccounts_ServiceError(t *testing.T) {
	mockOAuthService := new(MockOAuthService)
	mockAuthService := new(MockAuthService)
	userID := "user123"
	router := setupOAuthTestRouterWithAuth(mockOAuthService, mockAuthService, userID)

	mockOAuthService.On("GetLinkedAccounts", mock.Anything, userID).Return(nil, assert.AnError)

	req, _ := http.NewRequest("GET", "/api/v1/shared/oauth/accounts", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockOAuthService.AssertExpectations(t)
}

func TestOAuthAPI_UnlinkAccount_Success(t *testing.T) {
	mockOAuthService := new(MockOAuthService)
	mockAuthService := new(MockAuthService)
	userID := "user123"
	router := setupOAuthTestRouterWithAuth(mockOAuthService, mockAuthService, userID)

	accountID := "account-123"
	mockOAuthService.On("UnlinkAccount", mock.Anything, userID, accountID).Return(nil)

	req, _ := http.NewRequest("DELETE", "/api/v1/shared/oauth/accounts/"+accountID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "解绑成功", response["message"])

	mockOAuthService.AssertExpectations(t)
}

func TestOAuthAPI_UnlinkAccount_Unauthorized(t *testing.T) {
	mockOAuthService := new(MockOAuthService)
	mockAuthService := new(MockAuthService)
	router := setupOAuthTestRouterWithAuth(mockOAuthService, mockAuthService, "")

	req, _ := http.NewRequest("DELETE", "/api/v1/shared/oauth/accounts/account-123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestOAuthAPI_UnlinkAccount_EmptyAccountID(t *testing.T) {
	mockOAuthService := new(MockOAuthService)
	mockAuthService := new(MockAuthService)
	userID := "user123"
	router := setupOAuthTestRouterWithAuth(mockOAuthService, mockAuthService, userID)

	req, _ := http.NewRequest("DELETE", "/api/v1/shared/oauth/accounts/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestOAuthAPI_UnlinkAccount_ServiceError(t *testing.T) {
	mockOAuthService := new(MockOAuthService)
	mockAuthService := new(MockAuthService)
	userID := "user123"
	router := setupOAuthTestRouterWithAuth(mockOAuthService, mockAuthService, userID)

	accountID := "account-123"
	mockOAuthService.On("UnlinkAccount", mock.Anything, userID, accountID).Return(assert.AnError)

	req, _ := http.NewRequest("DELETE", "/api/v1/shared/oauth/accounts/"+accountID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockOAuthService.AssertExpectations(t)
}

func TestOAuthAPI_SetPrimaryAccount_Success(t *testing.T) {
	mockOAuthService := new(MockOAuthService)
	mockAuthService := new(MockAuthService)
	userID := "user123"
	router := setupOAuthTestRouterWithAuth(mockOAuthService, mockAuthService, userID)

	accountID := "account-123"
	mockOAuthService.On("SetPrimaryAccount", mock.Anything, userID, accountID).Return(nil)

	req, _ := http.NewRequest("PUT", "/api/v1/shared/oauth/accounts/"+accountID+"/primary", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "设置成功", response["message"])

	mockOAuthService.AssertExpectations(t)
}

func TestOAuthAPI_SetPrimaryAccount_Unauthorized(t *testing.T) {
	mockOAuthService := new(MockOAuthService)
	mockAuthService := new(MockAuthService)
	router := setupOAuthTestRouterWithAuth(mockOAuthService, mockAuthService, "")

	req, _ := http.NewRequest("PUT", "/api/v1/shared/oauth/accounts/account-123/primary", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestOAuthAPI_SetPrimaryAccount_EmptyAccountID(t *testing.T) {
	mockOAuthService := new(MockOAuthService)
	mockAuthService := new(MockAuthService)
	userID := "user123"
	router := setupOAuthTestRouterWithAuth(mockOAuthService, mockAuthService, userID)

	req, _ := http.NewRequest("PUT", "/api/v1/shared/oauth/accounts//primary", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestOAuthAPI_SetPrimaryAccount_ServiceError(t *testing.T) {
	mockOAuthService := new(MockOAuthService)
	mockAuthService := new(MockAuthService)
	userID := "user123"
	router := setupOAuthTestRouterWithAuth(mockOAuthService, mockAuthService, userID)

	accountID := "account-123"
	mockOAuthService.On("SetPrimaryAccount", mock.Anything, userID, accountID).Return(assert.AnError)

	req, _ := http.NewRequest("PUT", "/api/v1/shared/oauth/accounts/"+accountID+"/primary", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockOAuthService.AssertExpectations(t)
}
