package shared

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	authModel "Qingyu_backend/models/auth"
	authsvc "Qingyu_backend/service/auth"
	"golang.org/x/oauth2"
)

type authServiceStub struct {
	loginFunc        func(ctx context.Context, req *authsvc.LoginRequest) (*authsvc.LoginResponse, error)
	oauthLoginFunc   func(ctx context.Context, req *authsvc.OAuthLoginRequest) (*authsvc.LoginResponse, error)
	logoutFunc       func(ctx context.Context, token string) error
	refreshTokenFunc func(ctx context.Context, token string) (string, error)
}

func (s *authServiceStub) Register(ctx context.Context, req *authsvc.RegisterRequest) (*authsvc.RegisterResponse, error) {
	return nil, nil
}

func (s *authServiceStub) Login(ctx context.Context, req *authsvc.LoginRequest) (*authsvc.LoginResponse, error) {
	if s.loginFunc != nil {
		return s.loginFunc(ctx, req)
	}
	return nil, nil
}

func (s *authServiceStub) OAuthLogin(ctx context.Context, req *authsvc.OAuthLoginRequest) (*authsvc.LoginResponse, error) {
	if s.oauthLoginFunc != nil {
		return s.oauthLoginFunc(ctx, req)
	}
	return nil, nil
}

func (s *authServiceStub) Logout(ctx context.Context, token string) error {
	if s.logoutFunc != nil {
		return s.logoutFunc(ctx, token)
	}
	return nil
}

func (s *authServiceStub) RefreshToken(ctx context.Context, token string) (string, error) {
	if s.refreshTokenFunc != nil {
		return s.refreshTokenFunc(ctx, token)
	}
	return "", nil
}

func (s *authServiceStub) ValidateToken(ctx context.Context, token string) (*authsvc.TokenClaims, error) {
	return nil, nil
}

func (s *authServiceStub) CheckPermission(ctx context.Context, userID, permission string) (bool, error) {
	return false, nil
}

func (s *authServiceStub) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	return nil, nil
}

func (s *authServiceStub) HasRole(ctx context.Context, userID, role string) (bool, error) {
	return false, nil
}

func (s *authServiceStub) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	return nil, nil
}

func (s *authServiceStub) CreateRole(ctx context.Context, req *authsvc.CreateRoleRequest) (*authsvc.Role, error) {
	return nil, nil
}

func (s *authServiceStub) UpdateRole(ctx context.Context, roleID string, req *authsvc.UpdateRoleRequest) error {
	return nil
}

func (s *authServiceStub) DeleteRole(ctx context.Context, roleID string) error {
	return nil
}

func (s *authServiceStub) AssignRole(ctx context.Context, userID, roleID string) error {
	return nil
}

func (s *authServiceStub) RemoveRole(ctx context.Context, userID, roleID string) error {
	return nil
}

func (s *authServiceStub) CreateSession(ctx context.Context, userID string) (*authsvc.Session, error) {
	return nil, nil
}

func (s *authServiceStub) GetSession(ctx context.Context, sessionID string) (*authsvc.Session, error) {
	return nil, nil
}

func (s *authServiceStub) DestroySession(ctx context.Context, sessionID string) error {
	return nil
}

func (s *authServiceStub) RefreshSession(ctx context.Context, sessionID string) error {
	return nil
}

func (s *authServiceStub) Health(ctx context.Context) error {
	return nil
}

type oauthServiceStub struct {
	getAuthURLFunc    func(ctx context.Context, provider authModel.OAuthProvider, redirectURI, state string, linkMode bool, userID ...string) (string, error)
	exchangeCodeFunc  func(ctx context.Context, provider authModel.OAuthProvider, code, state string) (*oauth2.Token, *authModel.OAuthSession, error)
	getUserInfoFunc   func(ctx context.Context, provider authModel.OAuthProvider, token *oauth2.Token) (*authModel.UserIdentity, error)
	linkAccountFunc   func(ctx context.Context, userID string, provider authModel.OAuthProvider, token *oauth2.Token, identity *authModel.UserIdentity) (*authModel.OAuthAccount, error)
	unlinkAccountFunc func(ctx context.Context, userID, accountID string) error
	getLinkedFunc     func(ctx context.Context, userID string) ([]*authModel.OAuthAccount, error)
	setPrimaryFunc    func(ctx context.Context, userID, accountID string) error
}

func (s *oauthServiceStub) GetAuthURL(ctx context.Context, provider authModel.OAuthProvider, redirectURI, state string, linkMode bool, userID ...string) (string, error) {
	if s.getAuthURLFunc != nil {
		return s.getAuthURLFunc(ctx, provider, redirectURI, state, linkMode, userID...)
	}
	return "", nil
}

func (s *oauthServiceStub) ExchangeCode(ctx context.Context, provider authModel.OAuthProvider, code, state string) (*oauth2.Token, *authModel.OAuthSession, error) {
	if s.exchangeCodeFunc != nil {
		return s.exchangeCodeFunc(ctx, provider, code, state)
	}
	return nil, nil, nil
}

func (s *oauthServiceStub) GetUserInfo(ctx context.Context, provider authModel.OAuthProvider, token *oauth2.Token) (*authModel.UserIdentity, error) {
	if s.getUserInfoFunc != nil {
		return s.getUserInfoFunc(ctx, provider, token)
	}
	return nil, nil
}

func (s *oauthServiceStub) LinkAccount(ctx context.Context, userID string, provider authModel.OAuthProvider, token *oauth2.Token, identity *authModel.UserIdentity) (*authModel.OAuthAccount, error) {
	if s.linkAccountFunc != nil {
		return s.linkAccountFunc(ctx, userID, provider, token, identity)
	}
	return nil, nil
}

func (s *oauthServiceStub) UnlinkAccount(ctx context.Context, userID, accountID string) error {
	if s.unlinkAccountFunc != nil {
		return s.unlinkAccountFunc(ctx, userID, accountID)
	}
	return nil
}

func (s *oauthServiceStub) GetLinkedAccounts(ctx context.Context, userID string) ([]*authModel.OAuthAccount, error) {
	if s.getLinkedFunc != nil {
		return s.getLinkedFunc(ctx, userID)
	}
	return nil, nil
}

func (s *oauthServiceStub) SetPrimaryAccount(ctx context.Context, userID, accountID string) error {
	if s.setPrimaryFunc != nil {
		return s.setPrimaryFunc(ctx, userID, accountID)
	}
	return nil
}

func decodeSharedAPIResponse(t *testing.T, recorder *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
	return resp
}

func TestAuthAPI_Login_ServiceErrorReturnsCompatUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewAuthAPI(&authServiceStub{
		loginFunc: func(ctx context.Context, req *authsvc.LoginRequest) (*authsvc.LoginResponse, error) {
			assert.Equal(t, "tester", req.Username)
			assert.Equal(t, "wrong-password", req.Password)
			return nil, errors.New("invalid credentials")
		},
	})

	body := bytes.NewBufferString(`{"username":"tester","password":"wrong-password"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shared/auth/login", body)
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = req

	api.Login(c)

	require.Equal(t, http.StatusUnauthorized, recorder.Code)
	resp := decodeSharedAPIResponse(t, recorder)
	assert.Equal(t, "登录失败: invalid credentials", resp["message"])
}

func TestAuthAPI_RefreshToken_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewAuthAPI(&authServiceStub{
		refreshTokenFunc: func(ctx context.Context, token string) (string, error) {
			assert.Equal(t, "refresh-token", token)
			return "new-access-token", nil
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/shared/auth/refresh", nil)
	req.Header.Set("Authorization", "Bearer refresh-token")
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = req

	api.RefreshToken(c)

	require.Equal(t, http.StatusOK, recorder.Code)
	resp := decodeSharedAPIResponse(t, recorder)
	assert.Equal(t, "Token刷新成功", resp["message"])
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "new-access-token", data["token"])
}

func TestAuthAPI_RefreshToken_ServiceErrorReturnsUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewAuthAPI(&authServiceStub{
		refreshTokenFunc: func(ctx context.Context, token string) (string, error) {
			assert.Equal(t, "expired-refresh-token", token)
			return "", errors.New("token expired")
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/shared/auth/refresh", nil)
	req.Header.Set("Authorization", "Bearer expired-refresh-token")
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = req

	api.RefreshToken(c)

	require.Equal(t, http.StatusUnauthorized, recorder.Code)
	resp := decodeSharedAPIResponse(t, recorder)
	assert.Equal(t, "Token刷新失败: token expired", resp["message"])
}

func TestAuthAPI_Logout_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewAuthAPI(&authServiceStub{
		logoutFunc: func(ctx context.Context, token string) error {
			assert.Equal(t, "token-logout-123", token)
			return nil
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/shared/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer token-logout-123")
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = req

	api.Logout(c)

	require.Equal(t, http.StatusOK, recorder.Code)
	resp := decodeSharedAPIResponse(t, recorder)
	assert.Equal(t, "登出成功", resp["message"])
	_, hasData := resp["data"]
	assert.False(t, hasData)
}

func TestAuthAPI_Logout_ServiceErrorReturnsInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewAuthAPI(&authServiceStub{
		logoutFunc: func(ctx context.Context, token string) error {
			assert.Equal(t, "token-logout-123", token)
			return errors.New("revoke failed")
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/shared/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer token-logout-123")
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = req

	api.Logout(c)

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	resp := decodeSharedAPIResponse(t, recorder)
	assert.Equal(t, "服务器内部错误", resp["message"])
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "revoke failed", data["error"])
}

func TestAuthAPI_Logout_MissingTokenReturnsUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewAuthAPI(&authServiceStub{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shared/auth/logout", nil)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = req

	api.Logout(c)

	require.Equal(t, http.StatusUnauthorized, recorder.Code)
	resp := decodeSharedAPIResponse(t, recorder)
	assert.Equal(t, "未提供Token", resp["message"])
}

func TestAuthAPI_Logout_LegacyAuthorizationHeaderPassesThrough(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewAuthAPI(&authServiceStub{
		logoutFunc: func(ctx context.Context, token string) error {
			assert.Equal(t, "legacy-token-123", token)
			return nil
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/shared/auth/logout", nil)
	req.Header.Set("Authorization", "legacy-token-123")
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = req

	api.Logout(c)

	require.Equal(t, http.StatusOK, recorder.Code)
	resp := decodeSharedAPIResponse(t, recorder)
	assert.Equal(t, "登出成功", resp["message"])
	_, hasData := resp["data"]
	assert.False(t, hasData)
}
