package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 生成密码的bcrypt哈希
	password := "123456"

	// 生成哈希，cost factor为10（标准值）
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("生成哈希失败: %v\n", err)
		return
	}

	fmt.Printf("密码: %s\n", password)
	fmt.Printf("Bcrypt哈希: %s\n", string(hash))
	fmt.Printf("哈希长度: %d\n", len(hash))

	// 验证哈希
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		fmt.Printf("验证失败: %v\n", err)
	} else {
		fmt.Printf("✓ 哈希验证成功\n")
	}
}
