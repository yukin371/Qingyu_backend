package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "password"
	hash := "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"
	
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		fmt.Println("✓ 密码 'password' 与哈希匹配")
	} else {
		fmt.Printf("❌ 密码不匹配: %v\n", err)
	}
	
	// 生成新的哈希用于对比
	newHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	fmt.Printf("\n新生成的哈希: %s\n", string(newHash))
}
