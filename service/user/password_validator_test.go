package user

import (
	"testing"
)

// TestValidatePasswordStrength_Strong 测试强密码验证
func TestValidatePasswordStrength_Strong(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "包含大小写字母数字特殊字符的强密码",
			password: "MyStr0ng!Pass",
			want:     true,
		},
		{
			name:     "16位包含所有字符类型的密码",
			password: "Complex!Pass9247",
			want:     true,
		},
		{
			name:     "无特殊字符但足够长的密码",
			password: "MyStr0ngPassword",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, msg := validator.ValidateStrength(tt.password)
			if got != tt.want {
				t.Errorf("ValidateStrength() = %v, want %v, msg = %v", got, tt.want, msg)
			}
			if got && msg != "" {
				t.Errorf("ValidateStrength() returned success but got message: %v", msg)
			}
		})
	}
}

// TestValidatePasswordStrength_Medium 测试中等密码验证
func TestValidatePasswordStrength_Medium(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name     string
		password string
		want     bool
		level    string
	}{
		{
			name:     "8位基本密码",
			password: "Pass9247",
			want:     true,
			level:    "一般",
		},
		{
			name:     "12位包含大小写和数字",
			password: "MyPass924753",
			want:     true,
			level:    "中等",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, msg := validator.ValidateStrength(tt.password)
			if got != tt.want {
				t.Errorf("ValidateStrength() = %v, want %v, msg = %v", got, tt.want, msg)
			}
			if got {
				level := validator.GetStrengthLevel(tt.password)
				if level != tt.level {
					t.Errorf("GetStrengthLevel() = %v, want %v", level, tt.level)
				}
			}
		})
	}
}

// TestValidatePasswordStrength_Weak 测试弱密码验证
func TestValidatePasswordStrength_Weak(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name     string
		password string
		want     bool
		level    string
	}{
		{
			name:     "纯小写字母",
			password: "password",
			want:     false,
			level:    "弱",
		},
		{
			name:     "纯数字",
			password: "12345678",
			want:     false,
			level:    "弱",
		},
		{
			name:     "常见弱密码",
			password: "password123",
			want:     false,
			level:    "弱",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, msg := validator.ValidateStrength(tt.password)
			if got != tt.want {
				t.Errorf("ValidateStrength() = %v, want %v, msg = %v", got, tt.want, msg)
			}
			level := validator.GetStrengthLevel(tt.password)
			if level != tt.level {
				t.Errorf("GetStrengthLevel() = %v, want %v", level, tt.level)
			}
		})
	}
}

// TestValidatePasswordComplexity_Valid 测试满足复杂度要求
func TestValidatePasswordComplexity_Valid(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "满足所有复杂度要求",
			password: "ValidPass924",
			want:     true,
		},
		{
			name:     "包含特殊字符",
			password: "Valid!Pass924",
			want:     true,
		},
		{
			name:     "长密码",
			password: "VeryLongPassword9247",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, msg := validator.ValidateStrength(tt.password)
			if got != tt.want {
				t.Errorf("ValidateStrength() = %v, want %v, msg = %v", got, tt.want, msg)
			}
		})
	}
}

// TestValidatePasswordComplexity_TooShort 测试密码太短
func TestValidatePasswordComplexity_TooShort(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name     string
		password string
		wantErr  string
	}{
		{
			name:     "7位密码",
			password: "Pass123",
			wantErr:  "密码长度不能少于8位",
		},
		{
			name:     "1位密码",
			password: "a",
			wantErr:  "密码长度不能少于8位",
		},
		{
			name:     "空密码",
			password: "",
			wantErr:  "密码长度不能少于8位",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, msg := validator.ValidateStrength(tt.password)
			if got {
				t.Errorf("ValidateStrength() should fail for short password, got = %v", got)
			}
			if msg != tt.wantErr {
				t.Errorf("ValidateStrength() error message = %v, want %v", msg, tt.wantErr)
			}
		})
	}
}

// TestValidatePasswordComplexity_NoLowercase 测试无小写字母
func TestValidatePasswordComplexity_NoLowercase(t *testing.T) {
	validator := NewPasswordValidator()

	password := "NOPASS123"
	got, msg := validator.ValidateStrength(password)

	if got {
		t.Errorf("ValidateStrength() should fail without lowercase, got = %v", got)
	}
	if msg != "密码必须包含小写字母" {
		t.Errorf("ValidateStrength() error message = %v, want %v", msg, "密码必须包含小写字母")
	}
}

// TestValidatePasswordComplexity_NoUppercase 测试无大写字母
func TestValidatePasswordComplexity_NoUppercase(t *testing.T) {
	validator := NewPasswordValidator()

	password := "nopass123"
	got, msg := validator.ValidateStrength(password)

	if got {
		t.Errorf("ValidateStrength() should fail without uppercase, got = %v", got)
	}
	if msg != "密码必须包含大写字母" {
		t.Errorf("ValidateStrength() error message = %v, want %v", msg, "密码必须包含大写字母")
	}
}

// TestValidatePasswordComplexity_NoNumber 测试无数字
func TestValidatePasswordComplexity_NoNumber(t *testing.T) {
	validator := NewPasswordValidator()

	password := "NoNumberPass"
	got, msg := validator.ValidateStrength(password)

	if got {
		t.Errorf("ValidateStrength() should fail without number, got = %v", got)
	}
	if msg != "密码必须包含数字" {
		t.Errorf("ValidateStrength() error message = %v, want %v", msg, "密码必须包含数字")
	}
}

// TestIsCommonPassword_Yes 测试常见弱密码检测
func TestIsCommonPassword_Yes(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "password小写",
			password: "password",
			want:     true,
		},
		{
			name:     "password混合大小写",
			password: "Password",
			want:     true,
		},
		{
			name:     "12345678",
			password: "12345678",
			want:     true,
		},
		{
			name:     "admin123",
			password: "admin123",
			want:     true,
		},
		{
			name:     "qwerty",
			password: "qwerty",
			want:     true,
		},
		{
			name:     "root",
			password: "root",
			want:     true,
		},
		{
			name:     "abc123",
			password: "abc123",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validator.IsCommonPassword(tt.password)
			if got != tt.want {
				t.Errorf("IsCommonPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsCommonPassword_No 测试正常密码检测
func TestIsCommonPassword_No(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "复杂密码",
			password: "MyComplex!Pass123",
			want:     false,
		},
		{
			name:     "随机字符串",
			password: "Xk9#mP2$vL5@nQ8",
			want:     false,
		},
		{
			name:     "包含随机单词",
			password: "PurpleElephant123",
			want:     false,
		},
		{
			name:     "包含特殊字符",
			password: "Secure@Pass456",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validator.IsCommonPassword(tt.password)
			if got != tt.want {
				t.Errorf("IsCommonPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetStrengthScore 测试密码强度评分
func TestGetStrengthScore(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name     string
		password string
		wantMin  int
		wantMax  int
	}{
		{
			name:     "强密码高分",
			password: "Complex!Pass9247",
			wantMin:  80,
			wantMax:  100,
		},
		{
			name:     "一般密码中分",
			password: "Pass9247",
			wantMin:  50,
			wantMax:  59,
		},
		{
			name:     "弱密码低分",
			password: "password",
			wantMin:  0,
			wantMax:  39,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validator.GetStrengthScore(tt.password)
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("GetStrengthScore() = %v, want between %v and %v", got, tt.wantMin, tt.wantMax)
			}
		})
	}
}

// TestHasSequentialChars 测试连续字符检测
func TestHasSequentialChars(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "包含连续数字123",
			password: "Pass123word",
			want:     true,
		},
		{
			name:     "包含连续字母abc",
			password: "Passabc123",
			want:     true,
		},
		{
			name:     "包含连续大写字母ABC",
			password: "PassABC123",
			want:     true,
		},
		{
			name:     "不包含连续字符",
			password: "MyStr0ng!Pass",
			want:     false,
		},
		{
			name:     "包含数字但不连续",
			password: "Pass159word",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasSequentialChars(tt.password)
			if got != tt.want {
				t.Errorf("hasSequentialChars() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetStrengthLevel_Levels 测试强度等级分类
func TestGetStrengthLevel_Levels(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name     string
		password string
		want     string
	}{
		{
			name:     "强密码",
			password: "VeryComplex!Pass9247",
			want:     "强",
		},
		{
			name:     "中等密码",
			password: "MyPass9247Secure",
			want:     "中等",
		},
		{
			name:     "一般密码",
			password: "Mypassword",
			want:     "一般",
		},
		{
			name:     "弱密码",
			password: "12345678",
			want:     "弱",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validator.GetStrengthLevel(tt.password)
			if got != tt.want {
				t.Errorf("GetStrengthLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestValidateStrength_SequentialChars 测试连续字符验证
func TestValidateStrength_SequentialChars(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "包含连续数字",
			password: "Pass123word",
			want:     false,
		},
		{
			name:     "包含连续字母",
			password: "Passabc123",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := validator.ValidateStrength(tt.password)
			if got != tt.want {
				t.Errorf("ValidateStrength() = %v, want %v", got, tt.want)
			}
		})
	}
}


