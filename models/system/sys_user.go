package system

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User 表示系统中的用户数据模型
// 仅包含与数据本身紧密相关的字段与方法
type User struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Username  string    `bson:"username" json:"username"`
	Email     string    `bson:"email" json:"email"`
	Password  string    `bson:"password" json:"-"`
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
