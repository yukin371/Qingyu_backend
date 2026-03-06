package auth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
	ID             primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	UserID         string                 `bson:"user_id" json:"user_id"`
	Provider       OAuthProvider          `bson:"provider" json:"provider"`
	ProviderUserID string                 `bson:"provider_user_id" json:"provider_user_id"`
	Email          string                 `bson:"email" json:"email"`
	Username       string                 `bson:"username" json:"username,omitempty"`
	Avatar         string                 `bson:"avatar" json:"avatar"`
	AccessToken    string                 `bson:"access_token" json:"-"`
	RefreshToken   string                 `bson:"refresh_token" json:"-"`
	ExpiresAt      time.Time              `bson:"expires_at" json:"expires_at"`
	TokenExpiresAt time.Time              `bson:"token_expires_at" json:"token_expires_at"`
	Scope          string                 `bson:"scope" json:"scope"`
	IsPrimary      bool                   `bson:"is_primary" json:"is_primary"`
	LastLoginAt    time.Time              `bson:"last_login_at" json:"last_login_at"`
	Metadata       map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
	CreatedAt      time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time              `bson:"updated_at" json:"updated_at"`
}

// OAuthSession OAuth会话（用于OAuth流程中的临时状态）
type OAuthSession struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	State       string             `bson:"state" json:"state"`
	Provider    OAuthProvider      `bson:"provider" json:"provider"`
	RedirectURI string             `bson:"redirect_uri" json:"redirect_uri"`
	Scope       string             `bson:"scope" json:"scope"`
	UserID      string             `bson:"user_id,omitempty" json:"user_id"`
	LinkMode    bool               `bson:"link_mode" json:"link_mode"`
	ExpiresAt   time.Time          `bson:"expires_at" json:"expires_at"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

// OAuthConfig OAuth配置
type OAuthConfig struct {
	Provider     OAuthProvider `bson:"provider" json:"provider"`
	ClientID     string        `bson:"client_id" json:"client_id"`
	ClientSecret string        `bson:"client_secret" json:"-"`
	AuthURL      string        `bson:"auth_url" json:"auth_url"`
	TokenURL     string        `bson:"token_url" json:"token_url"`
	UserInfoURL  string        `bson:"user_info_url" json:"user_info_url"`
	RedirectURI  string        `bson:"redirect_uri" json:"redirect_uri"`
	Scopes       string        `bson:"scopes" json:"scopes"`
	Enabled      bool          `bson:"enabled" json:"enabled"`
}

// UserIdentity 用户身份信息（从第三方获取）
type UserIdentity struct {
	Provider      OAuthProvider `json:"provider"`
	ProviderID    string        `json:"provider_id"`
	Email         string        `json:"email"`
	EmailVerified bool          `json:"email_verified"`
	Name          string        `json:"name"`
	Avatar        string        `json:"avatar"`
	Username      string        `json:"username,omitempty"`
	Locale        string        `json:"locale,omitempty"`
	Gender        string        `json:"gender,omitempty"`
	Birthday      string        `json:"birthday,omitempty"`
}
