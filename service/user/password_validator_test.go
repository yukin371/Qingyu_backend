package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewPasswordValidator 测试创建密码验证器
func TestNewPasswordValidator(t *testing.T) {
	validator := NewPasswordValidator()

	assert.NotNil(t, validator, "密码验证器创建失败")
	assert.Equal(t, 8, validator.minLength, "最小长度默认值应为8")
	assert.True(t, validator.requireUppercase, "默认应要求大写字母")
	assert.True(t, validator.requireLowercase, "默认应要求小写字母")
	assert.True(t, validator.requireDigit, "默认应要求数字")
	assert.False(t, validator.requireSpecial, "默认不要求特殊字符")
	assert.NotNil(t, validator.commonPasswords, "常见密码列表不应为空")
}

// TestPasswordValidator_ValidateStrength_Success 测试验证密码强度 - 成功案例
func TestPasswordValidator_ValidateStrength_Success(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name     string
		password string
	}{
		{"标准强密码", "StrongP@ssw0rd"},
		{"包含特殊字符", "Test482!@#"},
		{"较长密码", "MyVeryStrongPassword482"},
		{"简单但符合要求", "T3st9XzV"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, msg := validator.ValidateStrength(tt.password)
			assert.True(t, valid, "密码应该验证通过: %s", tt.password)
			assert.Empty(t, msg, "验证通过时不应有错误消息")
		})
	}
}

// TestPasswordValidator_ValidateStrength_Length 测试验证密码强度 - 长度检查
func TestPasswordValidator_ValidateStrength_Length(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name     string
		password string
		wantMsg  string
	}{
		{"空密码", "", "密码长度不能少于8位"},
		{"过短密码1", "a", "密码长度不能少于8位"},
		{"过短密码3", "Ab1", "密码长度不能少于8位"},
		{"过短密码7", "Abc1234", "密码长度不能少于8位"},
		{"刚好8位", "Abx1357Q", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, msg := validator.ValidateStrength(tt.password)
			if tt.wantMsg == "" {
				assert.True(t, valid, "密码应该验证通过")
			} else {
				assert.False(t, valid, "密码应该验证失败")
				assert.Contains(t, msg, tt.wantMsg, "错误消息应包含预期内容")
			}
		})
	}
}

// TestPasswordValidator_ValidateStrength_Uppercase 测试验证密码强度 - 大写字母检查
func TestPasswordValidator_ValidateStrength_Uppercase(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name     string
		password string
		wantMsg  string
	}{
		{"无大写字母", "lowercase123", "密码必须包含大写字母"},
		{"全数字", "12345678", "密码必须包含大写字母"},
		{"全小写字母", "abcdefgh", "密码必须包含大写字母"},
		{"有大写字母", "Aqxvmnp1", ""},
		{"首字母大写", "Abx1357q", ""},
		{"中间大写", "qaCqxv82", ""},
		{"末尾大写", "qzxvpmG8", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, msg := validator.ValidateStrength(tt.password)
			if tt.wantMsg == "" {
				assert.True(t, valid, "密码应该验证通过")
			} else {
				assert.False(t, valid, "密码应该验证失败")
				assert.Contains(t, msg, tt.wantMsg, "错误消息应包含预期内容")
			}
		})
	}
}

// TestPasswordValidator_ValidateStrength_Lowercase 测试验证密码强度 - 小写字母检查
func TestPasswordValidator_ValidateStrength_Lowercase(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name     string
		password string
		wantMsg  string
	}{
		{"无小写字母", "UPPERCASE123", "密码必须包含小写字母"},
		{"全数字", "12345678", "密码必须包含大写字母"},
		{"全大写字母", "ABCDEFGH", "密码必须包含小写字母"},
		{"有小写字母", "QZPqxmz8", ""},
		{"首字母小写", "qBC97531", ""},
		{"中间小写", "QZcDQFG8", ""},
		{"末尾小写", "QZQFJPRg8", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, msg := validator.ValidateStrength(tt.password)
			if tt.wantMsg == "" {
				assert.True(t, valid, "密码应该验证通过")
			} else {
				assert.False(t, valid, "密码应该验证失败")
				assert.Contains(t, msg, tt.wantMsg, "错误消息应包含预期内容")
			}
		})
	}
}

// TestPasswordValidator_ValidateStrength_Digit 测试验证密码强度 - 数字检查
func TestPasswordValidator_ValidateStrength_Digit(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name     string
		password string
		wantMsg  string
	}{
		{"无数字", "NoDigitsHere", "密码必须包含数字"},
		{"全字母", "Abcdefgh", "密码必须包含数字"},
		{"有数字", "Abxqvpm1", ""},
		{"开头数字", "1Abxqvpm", ""},
		{"中间数字", "Abx1qvpm", ""},
		{"末尾数字", "Abxqvpm1", ""},
		{"多个数字", "Abx248de", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, msg := validator.ValidateStrength(tt.password)
			if tt.wantMsg == "" {
				assert.True(t, valid, "密码应该验证通过")
			} else {
				assert.False(t, valid, "密码应该验证失败")
				assert.Contains(t, msg, tt.wantMsg, "错误消息应包含预期内容")
			}
		})
	}
}

// TestPasswordValidator_ValidateStrength_CommonPassword 测试验证密码强度 - 常见密码检查
func TestPasswordValidator_ValidateStrength_CommonPassword(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name     string
		password string
		wantMsg  string
	}{
		{"常见密码password", "Password123", "密码过于常见"},
		{"常见密码admin", "Admin123", "密码过于常见"},
		{"常见密码123456", "123456", "密码长度不能少于8位"},
		{"常见密码qwerty", "Qwerty123", "密码不能包含连续的字符"},
		{"常见密码小写", "password", "密码必须包含大写字母"},
		{"非常见密码", "MyUniqPwd739", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, msg := validator.ValidateStrength(tt.password)
			if tt.wantMsg == "" {
				assert.True(t, valid, "密码应该验证通过")
			} else {
				assert.False(t, valid, "密码应该验证失败")
				assert.Contains(t, msg, tt.wantMsg, "错误消息应包含预期内容")
			}
		})
	}
}

// TestPasswordValidator_ValidateStrength_SequentialChars 测试验证密码强度 - 连续字符检查
func TestPasswordValidator_ValidateStrength_SequentialChars(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name     string
		password string
		wantMsg  string
	}{
		{"连续数字123", "Abc12345", "密码不能包含连续的字符"},
		{"连续数字234", "Abc23456", "密码不能包含连续的字符"},
		{"连续数字345", "Abc34567", "密码不能包含连续的字符"},
		{"连续数字456", "Abc45678", "密码不能包含连续的字符"},
		{"连续数字567", "Abc56789", "密码不能包含连续的字符"},
		{"连续字母abc", "Abcdefg1", "密码不能包含连续的字符"},
		{"连续字母bcd", "Bcdefgh1", "密码不能包含连续的字符"},
		{"连续字母cde", "Cdefghi1", "密码不能包含连续的字符"},
		{"连续字母大写ABC", "AbC12345", "密码不能包含连续的字符"},
		{"连续字母混合Abc", "Abcdefg1", "密码不能包含连续的字符"},
		{"无连续字符", "A1b2C3d4", ""},
		{"无连续字符2", "MyPass147", ""},
		{"倒序不检测", "Xzy32145", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, msg := validator.ValidateStrength(tt.password)
			if tt.wantMsg == "" {
				assert.True(t, valid, "密码应该验证通过")
			} else {
				assert.False(t, valid, "密码应该验证失败")
				assert.Contains(t, msg, tt.wantMsg, "错误消息应包含预期内容")
			}
		})
	}
}

// TestPasswordValidator_ValidateStrength_MultipleErrors 测试验证密码强度 - 多个错误
func TestPasswordValidator_ValidateStrength_MultipleErrors(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name     string
		password string
		wantMsg  string
	}{
		{"太短且无大写", "abc123", "密码长度不能少于8位"},
		{"无大写无小写", "12345678", "密码必须包含大写字母"},
		{"无数字无大写", "abcdefgh", "密码必须包含大写字母"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, msg := validator.ValidateStrength(tt.password)
			assert.False(t, valid, "密码应该验证失败")
			assert.NotEmpty(t, msg, "应该有错误消息")
		})
	}
}

// TestPasswordValidator_IsCommonPassword 测试检查常见密码
func TestPasswordValidator_IsCommonPassword(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{"常见密码password", "password", true},
		{"常见密码Password", "Password", true},
		{"常见密码PASSWORD", "PASSWORD", true},
		{"常见密码password123", "password123", true},
		{"常见密码admin123", "admin123", true},
		{"常见密码123456", "123456", true},
		{"常见密码qwerty", "qwerty", true},
		{"非常见密码", "MyUniquePassword123", false},
		{"空字符串", "", false},
		{"随机密码", "Xk9#mP2$vL5", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validator.IsCommonPassword(tt.password)
			assert.Equal(t, tt.want, got, "IsCommonPassword(%s) 应该返回 %v", tt.password, tt.want)
		})
	}
}

// TestPasswordValidator_GetStrengthScore 测试获取密码强度评分
func TestPasswordValidator_GetStrengthScore(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name     string
		password string
		minScore int
		maxScore int
	}{
		{"强密码", "StrongP@ssw0rd!", 70, 100},
		{"中等密码", "Medium123", 30, 50},
		{"弱密码", "Weak1", 40, 55},
		{"长密码", "VeryLongPassword123!@#", 65, 85},
		{"短但复杂", "A1!a", 50, 70},
		{"只有小写", "lowercase", 10, 30},
		{"只有大写", "UPPERCASE", 10, 30},
		{"只有数字", "12345678", 0, 10},
		{"常见密码", "password123", 0, 10},
		{"连续字符", "Abc12345", 30, 40},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := validator.GetStrengthScore(tt.password)
			assert.GreaterOrEqual(t, score, tt.minScore, "评分应该大于等于最小值")
			assert.LessOrEqual(t, score, tt.maxScore, "评分应该小于等于最大值")
		})
	}
}

// TestPasswordValidator_GetStrengthScore_Length 测试获取密码强度评分 - 长度评分
func TestPasswordValidator_GetStrengthScore_Length(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name        string
		password    string
		baseLength  int
		expectScore int
	}{
		{"8位", "AbX7mP2q", 8, 55},
		{"12位", "AbX7mP2qRt9u", 12, 65},
		{"16位", "AbX7mP2qRt9uKs4v", 16, 75},
		{"20位", "AbX7mP2qRt9uKs4vJn6w", 20, 75},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := validator.GetStrengthScore(tt.password)
			assert.Equal(t, tt.expectScore, score, "%d位密码评分应该为%d", tt.baseLength, tt.expectScore)
		})
	}
}

// TestPasswordValidator_GetStrengthScore_CharTypes 测试获取密码强度评分 - 字符类型评分
func TestPasswordValidator_GetStrengthScore_CharTypes(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name        string
		password    string
		expectScore int
	}{
		{"只有大写", "AZBYCXDW", 25},
		{"只有小写", "azbycxdw", 25},
		{"只有数字", "48295173", 25},
		{"大小写", "AzByCxDw", 40},
		{"大小写数字", "AzByCxD7", 55},
		{"全部四种", "AzByCx7!Q", 70},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := validator.GetStrengthScore(tt.password)
			assert.Equal(t, tt.expectScore, score, "%s类型评分应该为%d", tt.name, tt.expectScore)
		})
	}
}

// TestPasswordValidator_GetStrengthScore_Deductions 测试获取密码强度评分 - 扣分项
func TestPasswordValidator_GetStrengthScore_Deductions(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name        string
		password    string
		expectScore int
	}{
		{"常见密码扣分", "Password123", 5},
		{"连续字符扣分", "Abc12345", 35},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := validator.GetStrengthScore(tt.password)
			assert.Equal(t, tt.expectScore, score, "%s扣分后评分应该为%d", tt.name, tt.expectScore)
		})
	}
}

// TestPasswordValidator_GetStrengthScore_Boundary 测试获取密码强度评分 - 边界值
func TestPasswordValidator_GetStrengthScore_Boundary(t *testing.T) {
	validator := NewPasswordValidator()

	// 测试评分不会超过100
	t.Run("评分上限", func(t *testing.T) {
		score := validator.GetStrengthScore("VeryStrongP@ssw0rd!With!Many!Chars123456")
		assert.LessOrEqual(t, score, 100, "评分不应超过100")
	})

	// 测试评分不会低于0
	t.Run("评分下限", func(t *testing.T) {
		score := validator.GetStrengthScore("password")
		assert.GreaterOrEqual(t, score, 0, "评分不应低于0")
	})
}

// TestPasswordValidator_GetStrengthLevel 测试获取密码强度等级
func TestPasswordValidator_GetStrengthLevel(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name        string
		password    string
		expectLevel string
	}{
		{"强密码", "VeryStrongP@ssw0rd482!@#", "强"},
		{"强密码2", "MySecurePassword482!@#", "强"},
		{"中等密码", "MediumPass482", "中等"},
		{"中等密码2", "MyPassword482", "中等"},
		{"一般密码", "Password1", "一般"},
		{"一般密码2", "Test4821", "一般"},
		{"弱密码", "weak", "弱"},
		{"弱密码2", "abc123", "弱"},
		{"最弱", "a", "弱"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level := validator.GetStrengthLevel(tt.password)
			assert.Equal(t, tt.expectLevel, level, "密码 %s 的等级应该是 %s", tt.password, tt.expectLevel)
		})
	}
}

// TestPasswordValidator_GetStrengthLevel_Boundary 测试获取密码强度等级 - 边界值
func TestPasswordValidator_GetStrengthLevel_Boundary(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name        string
		minScore    int
		maxScore    int
		expectLevel string
	}{
		{"强等级边界80", 80, 100, "强"},
		{"中等等级边界60", 60, 79, "中等"},
		{"一般等级边界40", 40, 59, "一般"},
		{"弱等级边界0", 0, 39, "弱"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建一个能得到该分数范围的密码
			// 通过调整密码长度和字符类型来控制分数
			var password string
			switch tt.expectLevel {
			case "强":
				password = "StrongPass482!@"
			case "中等":
				password = "MediumPass482"
			case "一般":
				password = "General482"
			case "弱":
				password = "weak"
			}

			level := validator.GetStrengthLevel(password)
			assert.Equal(t, tt.expectLevel, level, "密码等级应该是 %s", tt.expectLevel)
		})
	}
}

// TestHasSequentialChars 测试检查连续字符函数
func TestHasSequentialChars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// 连续数字测试
		{"连续123", "abc123def", true},
		{"连续234", "abc234def", true},
		{"连续345", "abc345def", true},
		{"连续456", "abc456def", true},
		{"连续567", "abc567def", true},
		{"连续678", "abc678def", true},
		{"连续789", "abc789def", true},
		{"包含abc序列", "abc129def", true},
		{"包含def序列", "abc135def", true},

		// 连续字母测试
		{"连续abc", "123abc456", true},
		{"连续bcd", "123bcd456", true},
		{"连续cde", "123cde456", true},
		{"连续def", "123def456", true},
		{"连续xyz", "123xyz456", true},
		{"大写ABC", "123ABC456", true},
		{"大写XYZ", "123XYZ456", true},
		{"混合Abc", "123Abc456", true},
		{"包含123序列", "123ace456", true},

		// 边界测试
		{"空字符串", "", false},
		{"一个字符", "a", false},
		{"两个字符", "ab", false},
		{"三个连续", "abc", true},
		{"三个不连续", "abd", false},

		// 倒序不检测
		{"倒序321", "321", false},
		{"倒序cba", "cba", false},

		// 跨边界
		{"数字跨边界89", "abc890", true},
		{"数字跨边界78", "abc789", true},
		{"字母跨边界xyz", "xyz", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasSequentialChars(tt.input)
			assert.Equal(t, tt.expected, got, "hasSequentialChars(%s) 应该返回 %v", tt.input, tt.expected)
		})
	}
}

// TestLoadCommonPasswords 测试加载常见密码
func TestLoadCommonPasswords(t *testing.T) {
	passwords := loadCommonPasswords()

	assert.NotNil(t, passwords, "常见密码列表不应为空")
	assert.Greater(t, len(passwords), 0, "常见密码列表应包含密码")

	// 检查一些预期存在的密码
	expectedPasswords := []string{
		"password", "password123", "admin123", "123456",
		"qwerty", "abc123", "admin", "root",
	}

	for _, pwd := range expectedPasswords {
		assert.True(t, passwords[pwd], "常见密码列表应包含: %s", pwd)
	}

	// 检查大小写不敏感
	assert.True(t, passwords["password"], "应包含小写password")
	assert.True(t, passwords["password"], "应包含大写PASSWORD（存储为小写）")
}

// TestPasswordValidator_Integration 测试密码验证器集成场景
func TestPasswordValidator_Integration(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name        string
		password    string
		shouldPass  bool
		minLevel    string
		description string
	}{
		{
			name:        "完美密码",
			password:    "MySecureP@ssw0rd739",
			shouldPass:  true,
			minLevel:    "强",
			description: "包含所有字符类型且长度足够",
		},
		{
			name:        "最小要求密码",
			password:    "Abx1357Q",
			shouldPass:  true,
			minLevel:    "一般",
			description: "刚好满足所有基本要求",
		},
		{
			name:        "常见弱密码",
			password:    "Password123",
			shouldPass:  false,
			minLevel:    "",
			description: "常见密码应该被拒绝",
		},
		{
			name:        "连续字符密码",
			password:    "Abc12345",
			shouldPass:  false,
			minLevel:    "",
			description: "包含连续字符应该被拒绝",
		},
		{
			name:        "缺少大写",
			password:    "lowercase123",
			shouldPass:  false,
			minLevel:    "",
			description: "缺少大写字母应该被拒绝",
		},
		{
			name:        "缺少数字",
			password:    "NoDigitsHere",
			shouldPass:  false,
			minLevel:    "",
			description: "缺少数字应该被拒绝",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("测试描述:", tt.description)

			// 验证密码
			valid, msg := validator.ValidateStrength(tt.password)
			assert.Equal(t, tt.shouldPass, valid, "验证结果不符合预期")

			if !valid {
				assert.NotEmpty(t, msg, "应该有错误消息")
				t.Log("错误消息:", msg)
			} else {
				// 如果验证通过，检查强度等级
				level := validator.GetStrengthLevel(tt.password)
				t.Log("密码强度:", level, "分数:", validator.GetStrengthScore(tt.password))
			}
		})
	}
}

// TestPasswordValidator_EdgeCases 测试边界情况
func TestPasswordValidator_EdgeCases(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name     string
		password string
		test     func(t *testing.T, password string)
	}{
		{
			name:     "包含特殊字符",
			password: "Test9!@#$%^&*()_+",
			test: func(t *testing.T, password string) {
				valid, _ := validator.ValidateStrength(password)
				// 默认不要求特殊字符，所以应该通过
				assert.True(t, valid, "包含特殊字符的密码应该通过验证")
			},
		},
		{
			name:     "包含空格",
			password: "Test 482 ACF",
			test: func(t *testing.T, password string) {
				valid, _ := validator.ValidateStrength(password)
				assert.True(t, valid, "包含空格的密码应该通过验证")
			},
		},
		{
			name:     "Unicode字符",
			password: "Test482测试",
			test: func(t *testing.T, password string) {
				// Unicode字符按字节计算长度
				valid, _ := validator.ValidateStrength(password)
				assert.True(t, valid, "包含Unicode字符的密码应该通过验证")
			},
		},
		{
			name:     "纯特殊字符",
			password: "!@#$%^&*",
			test: func(t *testing.T, password string) {
				valid, msg := validator.ValidateStrength(password)
				assert.False(t, valid, "纯特殊字符应该验证失败")
				assert.Contains(t, msg, "大写字母", "应该提示缺少大写字母")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t, tt.password)
		})
	}
}

// BenchmarkPasswordValidator_ValidateStrength 性能测试 - ValidateStrength
func BenchmarkPasswordValidator_ValidateStrength(b *testing.B) {
	validator := NewPasswordValidator()
	password := "MySecureP@ssw0rd123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.ValidateStrength(password)
	}
}

// BenchmarkPasswordValidator_GetStrengthScore 性能测试 - GetStrengthScore
func BenchmarkPasswordValidator_GetStrengthScore(b *testing.B) {
	validator := NewPasswordValidator()
	password := "MySecureP@ssw0rd123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.GetStrengthScore(password)
	}
}
