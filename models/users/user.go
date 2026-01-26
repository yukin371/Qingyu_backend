package users

import (
	"time"

	"Qingyu_backend/models/auth"
	"Qingyu_backend/models/shared"
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
	shared.IdentifiedEntity `bson:",inline"`
	shared.BaseEntity       `bson:",inline"`

	Username string `bson:"username" json:"username" validate:"required,min=3,max=50"`
	Email    string `bson:"email,omitempty" json:"email" validate:"omitempty,email"`
	Phone    string `bson:"phone,omitempty" json:"phone" validate:"omitempty,e164"`
	Password string `bson:"password" json:"-" validate:"required,min=6"`

	// 角色和权限
	Roles    []string `bson:"roles" json:"roles" validate:"required,dive,oneof=reader author admin"` // 多角色支持
	VIPLevel int      `bson:"vip_level" json:"vipLevel" validate:"min=0,max=5"`                      // VIP等级 (0-5)

	Status   UserStatus `bson:"status" json:"status" validate:"required,oneof=active inactive banned deleted"`
	Avatar   string     `bson:"avatar,omitempty" json:"avatar,omitempty"` // 头像URL
	Nickname string     `bson:"nickname,omitempty" json:"nickname,omitempty" validate:"max=50"`
	Bio      string     `bson:"bio,omitempty" json:"bio,omitempty" validate:"max=500"`

	// 个人资料扩展字段
	Gender   string     `bson:"gender,omitempty" json:"gender,omitempty" validate:"omitempty,oneof=male female other"` // 性别
	Birthday *time.Time `bson:"birthday,omitempty" json:"birthday,omitempty"`                                         // 生日
	Location string     `bson:"location,omitempty" json:"location,omitempty" validate:"max=100"`                      // 位置
	Website  string     `bson:"website,omitempty" json:"website,omitempty" validate:"omitempty,url"`                  // 个人网站

	// 认证相关
	EmailVerified bool      `bson:"email_verified" json:"emailVerified"`
	PhoneVerified bool      `bson:"phone_verified" json:"phoneVerified"`
	LastLoginAt   time.Time `bson:"last_login_at,omitempty" json:"lastLoginAt,omitempty"`
	LastLoginIP   string    `bson:"last_login_ip,omitempty" json:"lastLoginIP,omitempty"`
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

// =========================
// 角色相关方法
// =========================

// GetEffectiveRoles 获取有效角色（含角色继承）
// 角色继承规则：admin → author → reader
func (u *User) GetEffectiveRoles() []string {
	roles := make(map[string]bool)
	for _, role := range u.Roles {
		roles[role] = true
		// 角色继承
		if role == auth.RoleAdmin {
			roles[auth.RoleAuthor] = true
			roles[auth.RoleReader] = true
		} else if role == auth.RoleAuthor {
			roles[auth.RoleReader] = true
		}
	}
	result := make([]string, 0, len(roles))
	for role := range roles {
		result = append(result, role)
	}
	return result
}

// HasRole 检查用户是否拥有指定角色（含继承）
func (u *User) HasRole(role string) bool {
	effectiveRoles := u.GetEffectiveRoles()
	for _, r := range effectiveRoles {
		if r == role {
			return true
		}
	}
	return false
}

// HasAnyRole 检查用户是否拥有指定角色中的任意一个
func (u *User) HasAnyRole(roles ...string) bool {
	effectiveRoles := u.GetEffectiveRoles()
	roleSet := make(map[string]bool)
	for _, r := range effectiveRoles {
		roleSet[r] = true
	}
	for _, role := range roles {
		if roleSet[role] {
			return true
		}
	}
	return false
}

// HasAllRoles 检查用户是否拥有所有指定的角色
func (u *User) HasAllRoles(roles ...string) bool {
	effectiveRoles := u.GetEffectiveRoles()
	roleSet := make(map[string]bool)
	for _, r := range effectiveRoles {
		roleSet[r] = true
	}
	for _, role := range roles {
		if !roleSet[role] {
			return false
		}
	}
	return true
}

// IsAdmin 检查用户是否为管理员
func (u *User) IsAdmin() bool {
	return u.HasRole(auth.RoleAdmin)
}

// IsAuthor 检查用户是否为作者
func (u *User) IsAuthor() bool {
	return u.HasRole(auth.RoleAuthor)
}

// IsReader 检查用户是否为读者
func (u *User) IsReader() bool {
	return u.HasRole(auth.RoleReader)
}

// AddRole 添加角色
func (u *User) AddRole(role string) {
	for _, r := range u.Roles {
		if r == role {
			return // 已存在该角色
		}
	}
	u.Roles = append(u.Roles, role)
}

// RemoveRole 移除角色
func (u *User) RemoveRole(role string) {
	for i, r := range u.Roles {
		if r == role {
			u.Roles = append(u.Roles[:i], u.Roles[i+1:]...)
			return
		}
	}
}

// =========================
// VIP相关方法
// =========================

// GetVIPLevel 获取VIP等级
func (u *User) GetVIPLevel() int {
	if u.VIPLevel < 0 {
		return 0
	}
	if u.VIPLevel > 5 {
		return 5
	}
	return u.VIPLevel
}

// IsVIP 检查用户是否为VIP（等级 > 0）
func (u *User) IsVIP() bool {
	return u.VIPLevel > 0
}

// HasVIPLevel 检查用户VIP等级是否达到指定等级
func (u *User) HasVIPLevel(level int) bool {
	return u.VIPLevel >= level
}

// SetVIPLevel 设置VIP等级
func (u *User) SetVIPLevel(level int) {
	if level < 0 {
		u.VIPLevel = 0
	} else if level > 5 {
		u.VIPLevel = 5
	} else {
		u.VIPLevel = level
	}
}

// UpdateLastLogin 更新最后登录时间和IP
func (u *User) UpdateLastLogin(ip string) {
	u.LastLoginAt = time.Now()
	u.LastLoginIP = ip
}
