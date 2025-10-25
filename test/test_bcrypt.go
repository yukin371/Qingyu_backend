package test

import (
	"fmt"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestBcrypt(t *testing.T) {
	fmt.Println("==============================================")
	fmt.Println("Bcrypt 密码兼容性测试")
	fmt.Println("==============================================")

	// 测试案例：从数据库中复制的实际密码哈希
	testCases := []struct {
		name     string
		hash     string
		password string
	}{
		{
			name:     "管理员 (admin)",
			hash:     "$2b$12$gKA1RJooGO50tDDHZpmYCuJDiZtELDcA1X1aZfU0WIYhnxZrJXT1e",
			password: "Admin@123456",
		},
		{
			name:     "测试用 - Go 生成的哈希",
			hash:     "", // 将由 Go 生成
			password: "Admin@123456",
		},
	}

	// 测试1：验证 Node.js bcrypt 生成的哈希 ($2b$)
	fmt.Println("测试 1: 验证 Node.js bcrypt 哈希 ($2b$)")
	fmt.Println("-----------------------------------------------")
	hash := testCases[0].hash
	password := testCases[0].password

	fmt.Printf("密码: %s\n", password)
	fmt.Printf("哈希: %s\n", hash)

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Printf("❌ 密码验证失败: %v\n", err)
		fmt.Println("\n⚠️  Go 的 bcrypt 无法验证 Node.js 生成的 $2b$ 哈希！")
		fmt.Println("这就是登录失败的原因。")
	} else {
		fmt.Printf("✅ 密码验证成功\n")
	}

	// 测试2：生成 Go 的 bcrypt 哈希并验证
	fmt.Println("\n测试 2: Go bcrypt 生成和验证")
	fmt.Println("-----------------------------------------------")

	goHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("❌ 生成哈希失败: %v\n", err)
		return
	}

	fmt.Printf("密码: %s\n", password)
	fmt.Printf("Go 生成的哈希: %s\n", string(goHash))

	err = bcrypt.CompareHashAndPassword(goHash, []byte(password))
	if err != nil {
		fmt.Printf("❌ 验证失败: %v\n", err)
	} else {
		fmt.Printf("✅ 验证成功\n")
	}

	// 测试3：尝试将 $2b$ 替换为 $2a$ 后验证
	fmt.Println("\n测试 3: 尝试 $2b$ → $2a$ 转换")
	fmt.Println("-----------------------------------------------")

	// 将 $2b$ 替换为 $2a$
	modifiedHash := hash
	if len(hash) > 3 && hash[0:3] == "$2b" {
		modifiedHash = "$2a" + hash[3:]
		fmt.Printf("原始哈希: %s\n", hash)
		fmt.Printf("修改哈希: %s\n", modifiedHash)

		err = bcrypt.CompareHashAndPassword([]byte(modifiedHash), []byte(password))
		if err != nil {
			fmt.Printf("❌ 验证失败: %v\n", err)
		} else {
			fmt.Printf("✅ 验证成功 - 可以通过替换前缀解决！\n")
		}
	}

	// 结论
	fmt.Println("\n==============================================")
	fmt.Println("结论")
	fmt.Println("==============================================")
	fmt.Println("\n如果测试1失败:")
	fmt.Println("  问题: Go 的 bcrypt 不兼容 Node.js 的 $2b$ 格式")
	fmt.Println("\n解决方案:")
	fmt.Println("  方案1: 重新用 Go 导入测试用户（推荐）")
	fmt.Println("         go run scripts/testing/import_test_users.go")
	fmt.Println("\n  方案2: 在 ValidatePassword 中转换 $2b$ → $2a$")
	fmt.Println("         修改 models/users/user.go")
	fmt.Println("\n  方案3: 使用兼容的 bcrypt 库")
	fmt.Println()
}
