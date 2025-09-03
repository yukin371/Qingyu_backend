package system

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User 用户模型
type User struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Username  string    `bson:"username" json:"username"`
	Password  string    `bson:"password" json:"-"` // 密码不会在JSON中返回
	Email     string    `bson:"email" json:"email"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

// SetPassword 设置用户密码（使用bcrypt加密）
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// ValidatePassword 验证密码是否正确
func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
