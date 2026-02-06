// Package utils 提供密码加密功能
package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword 生成密码哈希
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// DefaultPasswordHash 返回默认密码（password）的哈希
// 这是 "password" 的 bcrypt 哈希值
func DefaultPasswordHash() string {
	hash, _ := HashPassword("password")
	return hash
}
