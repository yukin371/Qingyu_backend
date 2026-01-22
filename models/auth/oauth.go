package auth

import (
	"time"
)

// OAuthProvider OAuth提供商类型
type OAuthProvider string

const (
	OAuthProviderGoogle OAuthProvider = "google"
	OAuthProviderGitHub OAuthProvider = "github"
	OAuthProviderQQ     OAuthProvider = "qq"
	OAuthProviderWeChat OAuthProvider = "wechat"
	OAuthProviderWeibo  OAuthProvider = "weibo"
)

// OAuthAccount OAuth绑定账号
type OAuthAccount struct {
	ID             string                 `bson:"_id" json:"id"`
	UserID         string                 `bson:"user_id" json:"$1$2"`
	Provider       OAuthProvider          `bson:"provider" json:"provider"`
	ProviderUserID string                 `bson:"provider_user_id" json:"$1$2"`
	Email          string                 `bson:"email" json:"email"`
	Username       string                 `bson:"username" json:"username,omitempty"`
	Avatar         string                 `bson:"avatar" json:"avatar"`
	AccessToken    string                 `bson:"access_token" json:"-"`  // 加密存储
	RefreshToken   string                 `bson:"refresh_token" json:"-"` // 加密存储
	ExpiresAt      time.Time              `bson:"expires_at" json:"$1$2"`
	TokenExpiresAt time.Time              `bson:"token_expires_at" json:"$1$2"`
	Scope          string                 `bson:"scope" json:"scope"`
	IsPrimary      bool                   `bson:"is_primary" json:"$1$2"` // 是否为主账号
	LastLoginAt    time.Time              `bson:"last_login_at" json:"$1$2"`
	Metadata       map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
	CreatedAt      time.Time              `bson:"created_at" json:"$1$2"`
	UpdatedAt      time.Time              `bson:"updated_at" json:"$1$2"`
}

// OAuthSession OAuth会话（用于OAuth流程中的临时状态）
type OAuthSession struct {
	ID          string        `bson:"_id" json:"id"`
	State       string        `bson:"state" json:"state"` // OAuth状态参数
	Provider    OAuthProvider `bson:"provider" json:"provider"`
	RedirectURI string        `bson:"redirect_uri" json:"$1$2"`
	Scope       string        `bson:"scope" json:"scope"`
	UserID      string        `bson:"user_id,omitempty" json:"$1$2"` // 已登录用户ID
	LinkMode    bool          `bson:"link_mode" json:"$1$2"`       // 是否为绑定模式
	ExpiresAt   time.Time     `bson:"expires_at" json:"$1$2"`
	CreatedAt   time.Time     `bson:"created_at" json:"$1$2"`
}

// OAuthConfig OAuth配置
type OAuthConfig struct {
	Provider     OAuthProvider `bson:"provider" json:"provider"`
	ClientID     string        `bson:"client_id" json:"$1$2"`
	ClientSecret string        `bson:"client_secret" json:"-"` // 不输出到JSON
	AuthURL      string        `bson:"auth_url" json:"$1$2"`
	TokenURL     string        `bson:"token_url" json:"$1$2"`
	UserInfoURL  string        `bson:"user_info_url" json:"$1$2"`
	RedirectURI  string        `bson:"redirect_uri" json:"$1$2"` // OAuth回调地址
	Scopes       string        `bson:"scopes" json:"scopes"`
	Enabled      bool          `bson:"enabled" json:"enabled"`
}

// UserIdentity 用户身份信息（从第三方获取）
type UserIdentity struct {
	Provider      OAuthProvider `json:"provider"`
	ProviderID    string        `json:"$1$2"`
	Email         string        `json:"email"`
	EmailVerified bool          `json:"$1$2"`
	Name          string        `json:"name"`
	Avatar        string        `json:"avatar"`
	Username      string        `json:"username,omitempty"`
	Locale        string        `json:"locale,omitempty"`
	Gender        string        `json:"gender,omitempty"`
	Birthday      string        `json:"birthday,omitempty"`
}
