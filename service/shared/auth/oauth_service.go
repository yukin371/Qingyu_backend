package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"go.uber.org/zap"

	authModel "Qingyu_backend/models/auth"
	"Qingyu_backend/repository/interfaces/auth"
)

// OAuthService OAuth服务
type OAuthService struct {
	logger     *zap.Logger
	repo       auth.OAuthRepository
	config     map[authModel.OAuthProvider]*oauth2.Config
	stateStore map[string]*authModel.OAuthSession // 内存状态存储（生产环境应该用Redis）
}

// NewOAuthService 创建OAuth服务
func NewOAuthService(logger *zap.Logger, repo auth.OAuthRepository, configs map[string]*authModel.OAuthConfig) (*OAuthService, error) {
	service := &OAuthService{
		logger:     logger,
		repo:       repo,
		config:     make(map[authModel.OAuthProvider]*oauth2.Config),
		stateStore: make(map[string]*authModel.OAuthSession),
	}

	// 初始化各平台的OAuth配置
	for provider, cfg := range configs {
		if !cfg.Enabled {
			continue
		}

		switch provider {
		case "google":
			service.config[authModel.OAuthProviderGoogle] = &oauth2.Config{
				ClientID:     cfg.ClientID,
				ClientSecret: cfg.ClientSecret,
				RedirectURL:  cfg.RedirectURI,
				Scopes:       strings.Split(cfg.Scopes, " "),
				Endpoint:     google.Endpoint,
			}

		case "github":
			service.config[authModel.OAuthProviderGitHub] = &oauth2.Config{
				ClientID:     cfg.ClientID,
				ClientSecret: cfg.ClientSecret,
				RedirectURL:  cfg.RedirectURI,
				Scopes:       strings.Split(cfg.Scopes, " "),
				Endpoint: oauth2.Endpoint{
					AuthURL:  "https://github.com/login/oauth/authorize",
					TokenURL: "https://github.com/login/oauth/access_token",
				},
			}

		case "qq":
			// QQ OAuth2.0配置（自定义端点）
			service.config[authModel.OAuthProviderQQ] = &oauth2.Config{
				ClientID:     cfg.ClientID,
				ClientSecret: cfg.ClientSecret,
				RedirectURL:  cfg.RedirectURI,
				Scopes:       strings.Split(cfg.Scopes, " "),
				Endpoint: oauth2.Endpoint{
					AuthURL:  cfg.AuthURL,
					TokenURL: cfg.TokenURL,
				},
			}
		}
	}

	logger.Info("OAuth service initialized",
		zap.Int("providers", len(service.config)),
	)

	return service, nil
}

// GetAuthURL 获取OAuth授权URL
func (s *OAuthService) GetAuthURL(ctx context.Context, provider authModel.OAuthProvider, redirectURI, state string, linkMode bool, userID ...string) (string, error) {
	config, exists := s.config[provider]
	if !exists {
		return "", fmt.Errorf("unsupported OAuth provider: %s", provider)
	}

	// 创建OAuth会话
	oauthState := s.generateState()
	session := &authModel.OAuthSession{
		ID:          generateSessionID(),
		State:       oauthState,
		Provider:    provider,
		RedirectURI: redirectURI,
		ExpiresAt:   time.Now().Add(10 * time.Minute), // 10分钟有效期
		CreatedAt:   time.Now(),
	}

	if linkMode && len(userID) > 0 {
		session.LinkMode = true
		session.UserID = userID[0]
	}

	// 保存会话
	s.stateStore[session.ID] = session
	s.logger.Debug("OAuth session created",
		zap.String("session_id", session.ID),
		zap.String("provider", string(provider)),
		zap.Bool("link_mode", session.LinkMode),
	)

	// 生成授权URL
	authURL := config.AuthCodeURL(oauthState)

	return authURL, nil
}

// ExchangeCode 交换授权码获取Token
func (s *OAuthService) ExchangeCode(ctx context.Context, provider authModel.OAuthProvider, code, state string) (*oauth2.Token, *authModel.OAuthSession, error) {
	config, exists := s.config[provider]
	if !exists {
		return nil, nil, fmt.Errorf("unsupported OAuth provider: %s", provider)
	}

	// 查找并验证会话
	var session *authModel.OAuthSession
	for _, s := range s.stateStore {
		if s.State == state {
			session = s
			break
		}
	}

	if session == nil {
		return nil, nil, fmt.Errorf("invalid OAuth state: %s", state)
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, nil, fmt.Errorf("OAuth session expired")
	}

	// 交换授权码
	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	return token, session, nil
}

// GetUserInfo 获取用户信息
func (s *OAuthService) GetUserInfo(ctx context.Context, provider authModel.OAuthProvider, token *oauth2.Token) (*authModel.UserIdentity, error) {
	switch provider {
	case authModel.OAuthProviderGoogle:
		return s.getGoogleUserInfo(ctx, token)
	case authModel.OAuthProviderGitHub:
		return s.getGitHubUserInfo(ctx, token)
	case authModel.OAuthProviderQQ:
		return s.getQQUserInfo(ctx, token)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// getGoogleUserInfo 获取Google用户信息
func (s *OAuthService) getGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*authModel.UserIdentity, error) {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var googleResp struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
		Locale        string `json:"locale"`
	}

	if err := json.Unmarshal(body, &googleResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &authModel.UserIdentity{
		Provider:     authModel.OAuthProviderGoogle,
		ProviderID:   googleResp.ID,
		Email:        googleResp.Email,
		EmailVerified: googleResp.VerifiedEmail,
		Name:         googleResp.Name,
		Avatar:       googleResp.Picture,
		Username:     strings.Split(googleResp.Email, "@")[0], // 使用邮箱前缀作为用户名
		Locale:       googleResp.Locale,
	}, nil
}

// getGitHubUserInfo 获取GitHub用户信息
func (s *OAuthService) getGitHubUserInfo(ctx context.Context, token *oauth2.Token) (*authModel.UserIdentity, error) {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))

	// 获取用户信息
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var githubResp struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
		Bio       string `json:"bio"`
		Location  string `json:"location"`
		Blog      string `json:"blog"`
	}

	if err := json.Unmarshal(body, &githubResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// 获取用户邮箱（可能为空）
	var email string
	var emailVerified bool

	// 获取用户邮箱（需要额外权限）
	emailResp, err := client.Get("https://api.github.com/user/emails")
	if err == nil && emailResp.StatusCode == http.StatusOK {
		defer emailResp.Body.Close()
		emailBody, _ := io.ReadAll(emailResp.Body)

		var emails []struct {
			Email   string `json:"email"`
			Primary bool   `json:"primary"`
			Verified bool `json:"verified"`
		}
		json.Unmarshal(emailBody, &emails)

		for _, e := range emails {
			if e.Primary {
				email = e.Email
				emailVerified = e.Verified
				break
			}
		}
	}

	return &authModel.UserIdentity{
		Provider:       authModel.OAuthProviderGitHub,
		ProviderID:     fmt.Sprintf("%d", githubResp.ID),
		Email:          email,
		EmailVerified:  emailVerified,
		Name:           githubResp.Name,
		Username:       githubResp.Login,
		Avatar:         githubResp.AvatarURL,
	}, nil
}

// getQQUserInfo 获取QQ用户信息
func (s *OAuthService) getQQUserInfo(ctx context.Context, token *oauth2.Token) (*authModel.UserIdentity, error) {
	// QQ OAuth2.0需要使用OpenID
	// 这里简化处理，实际需要调用QQ的OpenID接口
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))

	// 获取OpenID
	resp, err := client.Get("https://graph.qq.com/oauth2.0/me?access_token=" + url.QueryEscape(token.AccessToken) + "&fmt=json")
	if err != nil {
		return nil, fmt.Errorf("failed to get QQ OpenID: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// QQ返回的是JSONP格式，需要解析
	// 示例: callback({"client_id":"...","openid":"..."})
	var qqOpenIDResp struct {
		ClientID string `json:"client_id"`
		OpenID   string `json:"openid"`
	}

	// 移除JSONP包装
	strBody := string(body)
	start := strings.Index(strBody, "{")
	end := strings.LastIndex(strBody, "}")
	if start == -1 || end == -1 {
		return nil, fmt.Errorf("invalid OpenID response format")
	}
	jsonStr := strBody[start : end+1]

	if err := json.Unmarshal([]byte(jsonStr), &qqOpenIDResp); err != nil {
		return nil, fmt.Errorf("failed to parse OpenID response: %w", err)
	}

	// 获取用户信息
	userInfoResp, err := client.Get("https://graph.qq.com/user/get_user_info?access_token=" +
		url.QueryEscape(token.AccessToken) + "&oauth_consumer_key=" + s.config[authModel.OAuthProviderQQ].ClientID +
		"&openid=" + qqOpenIDResp.OpenID + "&format=json")
	if err != nil {
		return nil, fmt.Errorf("failed to get QQ user info: %w", err)
	}
	defer userInfoResp.Body.Close()

	if userInfoResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", userInfoResp.StatusCode)
	}

	userInfoBody, _ := io.ReadAll(userInfoResp.Body)

	var qqUserInfoResp struct {
		Ret        int    `json:"ret"`
		Msg        string `json:"msg"`
		Nickname   string `json:"nickname"`
		Gender     string `json:"gender"`
		FigureURLQQ string `json:"figureurl_qq_1"` // 中等尺寸头像
	}

	if err := json.Unmarshal(userInfoBody, &qqUserInfoResp); err != nil {
		return nil, fmt.Errorf("failed to parse user info response: %w", err)
	}

	if qqUserInfoResp.Ret != 0 {
		return nil, fmt.Errorf("QQ API error: %s", qqUserInfoResp.Msg)
	}

	return &authModel.UserIdentity{
		Provider:       authModel.OAuthProviderQQ,
		ProviderID:     qqOpenIDResp.OpenID,
		Name:           qqUserInfoResp.Nickname,
		Avatar:         qqUserInfoResp.FigureURLQQ,
		Username:       fmt.Sprintf("qq_%s", qqOpenIDResp.OpenID[:8]),
		EmailVerified:  false, // QQ不提供邮箱验证信息
	}, nil
}

// CleanupSession 清理过期的OAuth会话
func (s *OAuthService) CleanupSession(ctx context.Context) {
	now := time.Now()
	for id, session := range s.stateStore {
		if now.After(session.ExpiresAt) {
			delete(s.stateStore, id)
		}
	}
}

// 辅助函数

func (s *OAuthService) generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func generateSessionID() string {
	return fmt.Sprintf("oauth_%d", time.Now().UnixNano())
}

// getTokenExtraString 从token的Extra中安全地获取字符串值
func getTokenExtraString(token *oauth2.Token, key string) string {
	val := token.Extra(key)
	if val == nil {
		return ""
	}
	if str, ok := val.(string); ok {
		return str
	}
	return ""
}

// ==================== OAuth账号管理 ====================

// LinkAccount 绑定OAuth账号到用户
func (s *OAuthService) LinkAccount(ctx context.Context, userID string, provider authModel.OAuthProvider, token *oauth2.Token, identity *authModel.UserIdentity) (*authModel.OAuthAccount, error) {
	// 检查是否已经绑定过该账号
	existing, err := s.repo.FindByProviderAndProviderID(ctx, provider, identity.ProviderID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing account: %w", err)
	}

	if existing != nil {
		if existing.UserID == userID {
			// 已经绑定到当前用户
			return existing, nil
		}
		return nil, fmt.Errorf("this %s account is already linked to another user", provider)
	}

	// 检查该用户是否已有主账号
	hasPrimary, _ := s.repo.GetPrimaryAccount(ctx, userID)

	// 创建OAuth账号记录
	account := &authModel.OAuthAccount{
		UserID:          userID,
		Provider:        provider,
		ProviderUserID:  identity.ProviderID,
		Email:           identity.Email,
		Username:        identity.Username,
		Avatar:          identity.Avatar,
		AccessToken:     token.AccessToken,
		RefreshToken:    token.RefreshToken,
		TokenExpiresAt:  time.Now().Add(time.Duration(token.ExpiresIn) * time.Second),
		Scope:           getTokenExtraString(token, "scope"),
		IsPrimary:       hasPrimary == nil, // 如果没有主账号，则设为主账号
		LastLoginAt:     time.Now(),
		Metadata:        make(map[string]interface{}),
	}

	if err := s.repo.Create(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to create OAuth account: %w", err)
	}

	s.logger.Info("OAuth account linked",
		zap.String("user_id", userID),
		zap.String("provider", string(provider)),
		zap.String("provider_user_id", identity.ProviderID),
	)

	return account, nil
}

// UnlinkAccount 解绑OAuth账号
func (s *OAuthService) UnlinkAccount(ctx context.Context, userID, accountID string) error {
	account, err := s.repo.FindByID(ctx, accountID)
	if err != nil {
		return fmt.Errorf("failed to find OAuth account: %w", err)
	}

	if account == nil {
		return fmt.Errorf("OAuth account not found")
	}

	if account.UserID != userID {
		return fmt.Errorf("OAuth account does not belong to user")
	}

	// 检查是否是主账号且用户只有一个账�?
	count, _ := s.repo.CountByUserID(ctx, userID)
	if account.IsPrimary && count <= 1 {
		return fmt.Errorf("cannot unlink primary account when it's the only account")
	}

	if err := s.repo.Delete(ctx, accountID); err != nil {
		return fmt.Errorf("failed to delete OAuth account: %w", err)
	}

	s.logger.Info("OAuth account unlinked",
		zap.String("user_id", userID),
		zap.String("account_id", accountID),
	)

	return nil
}

// GetLinkedAccounts 获取用户绑定的所有OAuth账号
func (s *OAuthService) GetLinkedAccounts(ctx context.Context, userID string) ([]*authModel.OAuthAccount, error) {
	accounts, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get linked accounts: %w", err)
	}

	return accounts, nil
}

// SetPrimaryAccount 设置主账�?
func (s *OAuthService) SetPrimaryAccount(ctx context.Context, userID, accountID string) error {
	if err := s.repo.SetPrimaryAccount(ctx, userID, accountID); err != nil {
		return fmt.Errorf("failed to set primary account: %w", err)
	}

	s.logger.Info("Primary account set",
		zap.String("user_id", userID),
		zap.String("account_id", accountID),
	)

	return nil
}

// RefreshToken 刷新OAuth令牌
func (s *OAuthService) RefreshToken(ctx context.Context, accountID string) (*oauth2.Token, error) {
	account, err := s.repo.FindByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to find OAuth account: %w", err)
	}

	if account == nil {
		return nil, fmt.Errorf("OAuth account not found")
	}

	config, exists := s.config[account.Provider]
	if !exists {
		return nil, fmt.Errorf("OAuth provider not configured: %s", account.Provider)
	}

	// 使用refresh_token获取新的token
	tokenSource := config.TokenSource(ctx, &oauth2.Token{
		AccessToken:  account.AccessToken,
		RefreshToken: account.RefreshToken,
		Expiry:       account.TokenExpiresAt,
	})

	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	// 更新数据库中的token
	var expiresAt primitive.DateTime
	if !newToken.Expiry.IsZero() {
		expiresAt = primitive.NewDateTimeFromTime(newToken.Expiry)
	}

	if err := s.repo.UpdateTokens(ctx, accountID, newToken.AccessToken, newToken.RefreshToken, expiresAt); err != nil {
		s.logger.Warn("Failed to update refreshed token", zap.Error(err))
	}

	s.logger.Info("OAuth token refreshed",
		zap.String("account_id", accountID),
		zap.String("provider", string(account.Provider)),
	)

	return newToken, nil
}

// CleanupExpiredSessions 清理过期的OAuth会话
func (s *OAuthService) CleanupExpiredSessions(ctx context.Context) error {
	// 清理内存中的会话
	now := time.Now()
	for id, session := range s.stateStore {
		if now.After(session.ExpiresAt) {
			delete(s.stateStore, id)
		}
	}

	// 清理数据库中的会�?
	count, err := s.repo.CleanupExpiredSessions(ctx)
	if err != nil {
		s.logger.Warn("Failed to cleanup expired sessions from database", zap.Error(err))
	} else {
		s.logger.Debug("Cleaned up expired sessions", zap.Int64("count", count))
	}

	return nil
}
