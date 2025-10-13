package users

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// UserStatus 用户状态枚举
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"   // 活跃
	UserStatusInactive UserStatus = "inactive" // 未激活
	UserStatusBanned   UserStatus = "banned"   // 已封禁
	UserStatusDeleted  UserStatus = "deleted"  // 已删除
)

// User 表示系统中的用户数据模型
// 仅包含与数据本身紧密相关的字段与方法
type User struct {
	ID       string     `bson:"_id,omitempty" json:"id"`
	Username string     `bson:"username" json:"username" validate:"required,min=3,max=50"`
	Email    string     `bson:"email,omitempty" json:"email" validate:"omitempty,email"`
	Phone    string     `bson:"phone,omitempty" json:"phone" validate:"omitempty,e164"`
	Password string     `bson:"password" json:"-" validate:"required,min=6"`
	Role     string     `bson:"role" json:"role" validate:"required,oneof=user author admin"`
	Status   UserStatus `bson:"status" json:"status" validate:"required,oneof=active inactive banned deleted"`
	Avatar   string     `bson:"avatar,omitempty" json:"avatar,omitempty"` // 头像URL
	Nickname string     `bson:"nickname,omitempty" json:"nickname,omitempty" validate:"max=50"`
	Bio      string     `bson:"bio,omitempty" json:"bio,omitempty" validate:"max=500"`

	// 认证相关
	EmailVerified bool      `bson:"email_verified" json:"emailVerified"`
	PhoneVerified bool      `bson:"phone_verified" json:"phoneVerified"`
	LastLoginAt   time.Time `bson:"last_login_at,omitempty" json:"lastLoginAt,omitempty"`
	LastLoginIP   string    `bson:"last_login_ip,omitempty" json:"lastLoginIP,omitempty"`

	// 时间戳
	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}

// SetPassword 对明文密码进行哈希并设置到用户模型中
func (u *User) SetPassword(plainPassword string) error {
	if plainPassword == "" {
		return nil
	}
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedBytes)
	return nil
}

// ValidatePassword 校验明文密码是否与已存储的哈希一致
func (u *User) ValidatePassword(plainPassword string) bool {
	if u.Password == "" || plainPassword == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword)) == nil
}

// TouchForCreate 在创建前设置时间戳
func (u *User) TouchForCreate() {
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
}

// TouchForUpdate 在更新前刷新更新时间戳
func (u *User) TouchForUpdate() {
	u.UpdatedAt = time.Now()
}

// IsActive 检查用户是否为活跃状态
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// IsBanned 检查用户是否被封禁
func (u *User) IsBanned() bool {
	return u.Status == UserStatusBanned
}

// IsDeleted 检查用户是否被删除
func (u *User) IsDeleted() bool {
	return u.Status == UserStatusDeleted
}

// GetDisplayName 获取显示名称（优先昵称，其次用户名）
func (u *User) GetDisplayName() string {
	if u.Nickname != "" {
		return u.Nickname
	}
	return u.Username
}

// IsEmailVerified 检查邮箱是否已验证
func (u *User) IsEmailVerified() bool {
	return u.EmailVerified
}

// IsPhoneVerified 检查手机号是否已验证
func (u *User) IsPhoneVerified() bool {
	return u.PhoneVerified
}

// HasRole 检查用户是否拥有指定角色
func (u *User) HasRole(role string) bool {
	return u.Role == role
}

// IsAdmin 检查用户是否为管理员
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// IsAuthor 检查用户是否为作者
func (u *User) IsAuthor() bool {
	return u.Role == "author"
}

// UpdateLastLogin 更新最后登录时间和IP
func (u *User) UpdateLastLogin(ip string) {
	u.LastLoginAt = time.Now()
	u.LastLoginIP = ip
}
