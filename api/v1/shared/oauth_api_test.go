package shared

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	authModel "Qingyu_backend/models/auth"
	authsvc "Qingyu_backend/service/auth"
)

func TestOAuthAPI_GetAuthorizeURL_PassesLinkModeAndUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewOAuthAPI(&oauthServiceStub{
		getAuthURLFunc: func(ctx context.Context, provider authModel.OAuthProvider, redirectURI, state string, linkMode bool, userID ...string) (string, error) {
			assert.Equal(t, authModel.OAuthProviderGitHub, provider)
			assert.Equal(t, "https://app.example.com/oauth/callback", redirectURI)
			assert.Equal(t, "state-123", state)
			assert.True(t, linkMode)
			require.Len(t, userID, 1)
			assert.Equal(t, "user-123", userID[0])
			return "https://github.com/login/oauth/authorize?state=state-123", nil
		},
	}, &authServiceStub{}, zap.NewNop())

	body := bytes.NewBufferString(`{"redirect_uri":"https://app.example.com/oauth/callback","state":"state-123"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shared/oauth/github/authorize", body)
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Params = gin.Params{{Key: "provider", Value: "github"}}
	c.Set("user_id", "user-123")
	c.Request = req

	api.GetAuthorizeURL(c)

	require.Equal(t, http.StatusOK, recorder.Code)
	resp := decodeSharedAPIResponse(t, recorder)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "https://github.com/login/oauth/authorize?state=state-123", data["authorize_url"])
	assert.Equal(t, "github", data["provider"])
}

func TestOAuthAPI_HandleCallback_LinkModeUsesBindFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var oauthLoginCalled bool
	api := NewOAuthAPI(&oauthServiceStub{
		exchangeCodeFunc: func(ctx context.Context, provider authModel.OAuthProvider, code, state string) (*oauth2.Token, *authModel.OAuthSession, error) {
			assert.Equal(t, "code-123", code)
			assert.Equal(t, "state-123", state)
			return &oauth2.Token{AccessToken: "oauth-access"}, &authModel.OAuthSession{
				LinkMode: true,
				UserID:   "user-777",
			}, nil
		},
		getUserInfoFunc: func(ctx context.Context, provider authModel.OAuthProvider, token *oauth2.Token) (*authModel.UserIdentity, error) {
			assert.Equal(t, "oauth-access", token.AccessToken)
			return &authModel.UserIdentity{
				Provider:   provider,
				ProviderID: "provider-user-1",
				Email:      "link@example.com",
				Username:   "linker",
			}, nil
		},
		linkAccountFunc: func(ctx context.Context, userID string, provider authModel.OAuthProvider, token *oauth2.Token, identity *authModel.UserIdentity) (*authModel.OAuthAccount, error) {
			assert.Equal(t, "user-777", userID)
			assert.Equal(t, authModel.OAuthProviderGoogle, provider)
			assert.Equal(t, "provider-user-1", identity.ProviderID)
			return &authModel.OAuthAccount{
				ID:             primitive.NewObjectID(),
				UserID:         userID,
				Provider:       provider,
				ProviderUserID: identity.ProviderID,
			}, nil
		},
	}, &authServiceStub{
		oauthLoginFunc: func(ctx context.Context, req *authsvc.OAuthLoginRequest) (*authsvc.LoginResponse, error) {
			oauthLoginCalled = true
			return nil, nil
		},
	}, zap.NewNop())

	body := bytes.NewBufferString(`{"code":"code-123","state":"state-123"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shared/oauth/google/callback", body)
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Params = gin.Params{{Key: "provider", Value: "google"}}
	c.Request = req

	api.HandleCallback(c)

	require.Equal(t, http.StatusOK, recorder.Code)
	resp := decodeSharedAPIResponse(t, recorder)
	assert.Equal(t, "账号绑定成功", resp["message"])
	assert.False(t, oauthLoginCalled)
}

func TestOAuthAPI_HandleCallback_LoginModeCallsAuthService(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewOAuthAPI(&oauthServiceStub{
		exchangeCodeFunc: func(ctx context.Context, provider authModel.OAuthProvider, code, state string) (*oauth2.Token, *authModel.OAuthSession, error) {
			return &oauth2.Token{AccessToken: "oauth-access"}, &authModel.OAuthSession{}, nil
		},
		getUserInfoFunc: func(ctx context.Context, provider authModel.OAuthProvider, token *oauth2.Token) (*authModel.UserIdentity, error) {
			return &authModel.UserIdentity{
				Provider:   provider,
				ProviderID: "provider-user-2",
				Email:      "oauth@example.com",
				Name:       "OAuth User",
				Avatar:     "https://example.com/avatar.png",
				Username:   "oauth-user",
			}, nil
		},
	}, &authServiceStub{
		oauthLoginFunc: func(ctx context.Context, req *authsvc.OAuthLoginRequest) (*authsvc.LoginResponse, error) {
			assert.Equal(t, authModel.OAuthProviderQQ, req.Provider)
			assert.Equal(t, "provider-user-2", req.ProviderID)
			assert.Equal(t, "oauth@example.com", req.Email)
			assert.Equal(t, "OAuth User", req.Name)
			assert.Equal(t, "https://example.com/avatar.png", req.Avatar)
			assert.Equal(t, "oauth-user", req.Username)
			return &authsvc.LoginResponse{
				Token: "app-token",
				User: &authsvc.UserInfo{
					ID:       "user-1",
					Username: "oauth-user",
					Email:    "oauth@example.com",
					Roles:    []string{"reader"},
				},
			}, nil
		},
	}, zap.NewNop())

	body := bytes.NewBufferString(`{"code":"code-456","state":"state-456"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shared/oauth/qq/callback", body)
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Params = gin.Params{{Key: "provider", Value: "qq"}}
	c.Request = req

	api.HandleCallback(c)

	require.Equal(t, http.StatusOK, recorder.Code)
	resp := decodeSharedAPIResponse(t, recorder)
	assert.Equal(t, "OAuth登录成功", resp["message"])
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "app-token", data["token"])
}

func TestOAuthAPI_HandleCallback_ExchangeCodeFailureStopsFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var getUserInfoCalled bool
	var oauthLoginCalled bool
	api := NewOAuthAPI(&oauthServiceStub{
		exchangeCodeFunc: func(ctx context.Context, provider authModel.OAuthProvider, code, state string) (*oauth2.Token, *authModel.OAuthSession, error) {
			assert.Equal(t, authModel.OAuthProviderGoogle, provider)
			assert.Equal(t, "bad-code", code)
			assert.Equal(t, "state-err", state)
			return nil, nil, errors.New("invalid oauth code")
		},
		getUserInfoFunc: func(ctx context.Context, provider authModel.OAuthProvider, token *oauth2.Token) (*authModel.UserIdentity, error) {
			getUserInfoCalled = true
			return nil, nil
		},
	}, &authServiceStub{
		oauthLoginFunc: func(ctx context.Context, req *authsvc.OAuthLoginRequest) (*authsvc.LoginResponse, error) {
			oauthLoginCalled = true
			return nil, nil
		},
	}, zap.NewNop())

	body := bytes.NewBufferString(`{"code":"bad-code","state":"state-err"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shared/oauth/google/callback", body)
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Params = gin.Params{{Key: "provider", Value: "google"}}
	c.Request = req

	api.HandleCallback(c)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	resp := decodeSharedAPIResponse(t, recorder)
	assert.Equal(t, "授权码交换失败: invalid oauth code", resp["message"])
	assert.False(t, getUserInfoCalled)
	assert.False(t, oauthLoginCalled)
}

func TestOAuthAPI_HandleCallback_GetUserInfoFailureStopsLaterFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var linkAccountCalled bool
	var oauthLoginCalled bool
	api := NewOAuthAPI(&oauthServiceStub{
		exchangeCodeFunc: func(ctx context.Context, provider authModel.OAuthProvider, code, state string) (*oauth2.Token, *authModel.OAuthSession, error) {
			assert.Equal(t, authModel.OAuthProviderGitHub, provider)
			assert.Equal(t, "code-userinfo", code)
			assert.Equal(t, "state-userinfo", state)
			return &oauth2.Token{AccessToken: "oauth-access"}, &authModel.OAuthSession{}, nil
		},
		getUserInfoFunc: func(ctx context.Context, provider authModel.OAuthProvider, token *oauth2.Token) (*authModel.UserIdentity, error) {
			assert.Equal(t, "oauth-access", token.AccessToken)
			return nil, errors.New("provider userinfo unavailable")
		},
		linkAccountFunc: func(ctx context.Context, userID string, provider authModel.OAuthProvider, token *oauth2.Token, identity *authModel.UserIdentity) (*authModel.OAuthAccount, error) {
			linkAccountCalled = true
			return nil, nil
		},
	}, &authServiceStub{
		oauthLoginFunc: func(ctx context.Context, req *authsvc.OAuthLoginRequest) (*authsvc.LoginResponse, error) {
			oauthLoginCalled = true
			return nil, nil
		},
	}, zap.NewNop())

	body := bytes.NewBufferString(`{"code":"code-userinfo","state":"state-userinfo"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shared/oauth/github/callback", body)
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Params = gin.Params{{Key: "provider", Value: "github"}}
	c.Request = req

	api.HandleCallback(c)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	resp := decodeSharedAPIResponse(t, recorder)
	assert.Equal(t, "获取用户信息失败: provider userinfo unavailable", resp["message"])
	assert.False(t, linkAccountCalled)
	assert.False(t, oauthLoginCalled)
}

func TestOAuthAPI_HandleCallback_LinkModeLinkAccountFailureReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var oauthLoginCalled bool
	api := NewOAuthAPI(&oauthServiceStub{
		exchangeCodeFunc: func(ctx context.Context, provider authModel.OAuthProvider, code, state string) (*oauth2.Token, *authModel.OAuthSession, error) {
			return &oauth2.Token{AccessToken: "oauth-access"}, &authModel.OAuthSession{
				LinkMode: true,
				UserID:   "user-888",
			}, nil
		},
		getUserInfoFunc: func(ctx context.Context, provider authModel.OAuthProvider, token *oauth2.Token) (*authModel.UserIdentity, error) {
			return &authModel.UserIdentity{
				Provider:   provider,
				ProviderID: "provider-user-link-fail",
				Email:      "link@example.com",
				Username:   "linker",
			}, nil
		},
		linkAccountFunc: func(ctx context.Context, userID string, provider authModel.OAuthProvider, token *oauth2.Token, identity *authModel.UserIdentity) (*authModel.OAuthAccount, error) {
			assert.Equal(t, "user-888", userID)
			assert.Equal(t, authModel.OAuthProviderGoogle, provider)
			assert.Equal(t, "provider-user-link-fail", identity.ProviderID)
			return nil, errors.New("already linked elsewhere")
		},
	}, &authServiceStub{
		oauthLoginFunc: func(ctx context.Context, req *authsvc.OAuthLoginRequest) (*authsvc.LoginResponse, error) {
			oauthLoginCalled = true
			return nil, nil
		},
	}, zap.NewNop())

	body := bytes.NewBufferString(`{"code":"code-link","state":"state-link"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shared/oauth/google/callback", body)
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Params = gin.Params{{Key: "provider", Value: "google"}}
	c.Request = req

	api.HandleCallback(c)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	resp := decodeSharedAPIResponse(t, recorder)
	assert.Equal(t, "绑定账号失败: already linked elsewhere", resp["message"])
	assert.False(t, oauthLoginCalled)
}

func TestOAuthAPI_HandleCallback_LoginModeAuthServiceFailureReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewOAuthAPI(&oauthServiceStub{
		exchangeCodeFunc: func(ctx context.Context, provider authModel.OAuthProvider, code, state string) (*oauth2.Token, *authModel.OAuthSession, error) {
			return &oauth2.Token{AccessToken: "oauth-access"}, &authModel.OAuthSession{}, nil
		},
		getUserInfoFunc: func(ctx context.Context, provider authModel.OAuthProvider, token *oauth2.Token) (*authModel.UserIdentity, error) {
			return &authModel.UserIdentity{
				Provider:   provider,
				ProviderID: "provider-user-3",
				Email:      "oauth@example.com",
				Name:       "OAuth User",
				Username:   "oauth-user",
			}, nil
		},
	}, &authServiceStub{
		oauthLoginFunc: func(ctx context.Context, req *authsvc.OAuthLoginRequest) (*authsvc.LoginResponse, error) {
			assert.Equal(t, authModel.OAuthProviderGitHub, req.Provider)
			assert.Equal(t, "provider-user-3", req.ProviderID)
			return nil, errors.New("user disabled")
		},
	}, zap.NewNop())

	body := bytes.NewBufferString(`{"code":"code-789","state":"state-789"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shared/oauth/github/callback", body)
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Params = gin.Params{{Key: "provider", Value: "github"}}
	c.Request = req

	api.HandleCallback(c)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	resp := decodeSharedAPIResponse(t, recorder)
	assert.Equal(t, "OAuth登录失败: user disabled", resp["message"])
}
