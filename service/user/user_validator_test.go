package user

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	usersModel "Qingyu_backend/models/users"

	"github.com/stretchr/testify/assert"
)

// =========================
// validateUsername 测试
// =========================

func TestValidateUsername(t *testing.T) {
	validator := &UserValidator{}

	tests := []struct {
		name      string
		username  string
		wantErr   bool
		errCode   string
		errField  string
	}{
		{"有效用户名", "validuser", false, "", ""},
		{"有效用户名带数字", "user123", false, "", ""},
		{"有效用户名带下划线", "user_name", false, "", ""},
		{"有效用户名大写", "ValidUser", false, "", ""},
		{"空用户名", "", true, "REQUIRED", "username"},
		{"用户名太短", "ab", true, "MIN_LENGTH", "username"},
		{"用户名刚好3字符", "abc", false, "", ""},
		{"用户名刚好30字符", "this_is_exactly_thirty_char_ok", false, "", ""},
		{"用户名31字符(无特殊字符)", "this_is_thirty_one_characters_x", true, "MAX_LENGTH", "username"},
		{"用户名包含特殊字符", "user@name", true, "INVALID_FORMAT", "username"},
		{"用户名包含空格", "user name", true, "INVALID_FORMAT", "username"},
		{"用户名以数字开头", "123user", true, "INVALID_FORMAT", "username"},
		{"保留用户名admin", "admin", true, "RESERVED_NAME", "username"},
		{"保留用户名root", "root", true, "RESERVED_NAME", "username"},
		{"保留用户名system", "system", true, "RESERVED_NAME", "username"},
		{"保留用户名api", "api", true, "RESERVED_NAME", "username"},
		{"保留用户名www", "www", true, "RESERVED_NAME", "username"},
		{"保留用户名mail", "mail", true, "RESERVED_NAME", "username"},
		{"保留用户名ftp", "ftp", true, "RESERVED_NAME", "username"},
		{"保留用户名大写", "Admin", true, "RESERVED_NAME", "username"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateUsername(tt.username)
			if tt.wantErr {
				assert.Error(t, err)
				if err != nil {
					assert.Equal(t, tt.errCode, err.Code)
					assert.Equal(t, tt.errField, err.Field)
				}
			} else {
				// validateUsername 返回 *ValidationError 类型，需要使用 Nil 而不是 NoError
				assert.Nil(t, err)
			}
		})
	}
}

// =========================
// validateEmail 测试
// =========================

func TestValidateEmail(t *testing.T) {
	validator := &UserValidator{}

	tests := []struct {
		name     string
		email    string
		wantErr  bool
		errCode  string
		errField string
	}{
		{"有效邮箱标准格式", "test@example.com", false, "", ""},
		{"有效邮箱带子域名", "user@mail.example.com", false, "", ""},
		{"有效邮箱带点", "user.name@example.com", false, "", ""},
		{"有效邮箱带加号", "user+tag@example.com", false, "", ""},
		{"有效邮箱带下划线", "user_name@example.com", false, "", ""},
		{"有效邮箱带连字符", "user-name@example.com", false, "", ""},
		{"有效邮箱带数字", "user123@example.com", false, "", ""},
		{"有效邮箱短域名", "test@xy.com", false, "", ""},
		{"有效邮箱长TLD", "test@example.museum", false, "", ""},
		{"有效邮箱短用户名", "a@b.co", false, "", ""},
		{"空邮箱", "", true, "REQUIRED", "email"},
		{"邮箱太长", "this_is_a_very_long_email_address_that_exceeds_one_hundred_characters_and_should_not_be_allowed_by_the_validator@example.com", true, "MAX_LENGTH", "email"},
		{"缺少@符号", "userexample.com", true, "INVALID_FORMAT", "email"},
		{"缺少域名", "user@", true, "INVALID_FORMAT", "email"},
		{"缺少用户名", "@example.com", true, "INVALID_FORMAT", "email"},
		{"缺少顶级域名", "user@example", true, "INVALID_FORMAT", "email"},
		{"TLD只有一个字符", "user@example.c", true, "INVALID_FORMAT", "email"},
		{"多个@符号", "user@name@example.com", true, "INVALID_FORMAT", "email"},
		{"无效字符括号", "user(name)@example.com", true, "INVALID_FORMAT", "email"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateEmail(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
				if err != nil {
					assert.Equal(t, tt.errCode, err.Code)
					assert.Equal(t, tt.errField, err.Field)
				}
			} else {
				// validateEmail 返回 *ValidationError 类型，需要使用 Nil 而不是 NoError
				assert.Nil(t, err)
			}
		})
	}
}

// =========================
// validatePassword 测试
// =========================

func TestValidatePassword(t *testing.T) {
	validator := &UserValidator{}

	tests := []struct {
		name     string
		password string
		wantErr  bool
		errCode  string
		errField string
	}{
		{"有效密码标准格式", "ValidPass123", false, "", ""},
		{"有效密码带特殊字符", "Valid@Pass123", false, "", ""},
		{"有效密码较大", "ThisIsAVeryLongPassword12345", false, "", ""},
		{"有效密码刚好8字符", "Valid123", false, "", ""},
		{"空密码", "", true, "REQUIRED", "password"},
		{"密码太短7字符", "Pass123", true, "MIN_LENGTH", "password"},
		{"密码太短1字符", "a", true, "MIN_LENGTH", "password"},
		{"密码只有小写字母", "onlylowercase", true, "WEAK_PASSWORD", "password"},
		{"密码只有大写字母", "ONLYUPPERCASE", true, "WEAK_PASSWORD", "password"},
		{"密码只有数字", "12345678", true, "WEAK_PASSWORD", "password"},
		{"密码只有字母", "OnlyLetters", true, "WEAK_PASSWORD", "password"},
		{"密码只有特殊字符", "!@#$%^&*()", true, "WEAK_PASSWORD", "password"},
		{"常见弱密码12345678", "12345678", true, "WEAK_PASSWORD", "password"},
		{"常见弱密码password", "password", true, "WEAK_PASSWORD", "password"},
		{"常见弱密码qwerty123", "qwerty123", true, "WEAK_PASSWORD", "password"},
		{"常见弱密码abc123456", "abc123456", true, "WEAK_PASSWORD", "password"},
		{"常见弱密码大写", "PASSWORD123", true, "WEAK_PASSWORD", "password"},
		{"常见弱密码小写", "password123", true, "WEAK_PASSWORD", "password"},
		{"包含弱密码但不在列表中", "password123!@#", false, "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validatePassword(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
				if err != nil {
					assert.Equal(t, tt.errCode, err.Code)
					assert.Equal(t, tt.errField, err.Field)
				}
			} else {
				// validatePassword 返回 *ValidationError 类型，需要使用 Nil 而不是 NoError
				assert.Nil(t, err)
			}
		})
	}
}

// =========================
// validateBasicFields 测试
// =========================

func TestValidateBasicFields(t *testing.T) {
	validator := &UserValidator{}

	tests := []struct {
		name     string
		user     *usersModel.User
		wantErr  bool
		errCount int
	}{
		{
			"有效基础字段",
			&usersModel.User{Username: "validuser", Email: "valid@example.com", Password: "ValidPass123"},
			false, 0,
		},
		{
			"所有字段都无效",
			&usersModel.User{Username: "", Email: "", Password: ""},
			true, 2,
		},
		{
			"只有用户名无效",
			&usersModel.User{Username: "ab", Email: "valid@example.com", Password: "ValidPass123"},
			true, 1,
		},
		{
			"只有邮箱无效",
			&usersModel.User{Username: "validuser", Email: "invalid-email", Password: "ValidPass123"},
			true, 1,
		},
		{
			"只有密码无效",
			&usersModel.User{Username: "validuser", Email: "valid@example.com", Password: "weak"},
			true, 1,
		},
		{
			"密码为空时不验证密码",
			&usersModel.User{Username: "validuser", Email: "valid@example.com", Password: ""},
			false, 0,
		},
		{
			"所有字段都有错误",
			&usersModel.User{Username: "a", Email: "invalid", Password: "weak"},
			true, 3,
		},
		{
			"用户名是保留字",
			&usersModel.User{Username: "admin", Email: "admin@example.com", Password: "ValidPass123"},
			true, 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := validator.validateBasicFields(tt.user)
			if tt.wantErr {
				assert.Greater(t, len(errs), 0)
				if tt.errCount > 0 {
					assert.Equal(t, tt.errCount, len(errs))
				}
			} else {
				assert.Equal(t, 0, len(errs))
			}
		})
	}
}

// =========================
// validateUpdateFields 测试
// =========================

func TestValidateUpdateFields(t *testing.T) {
	validator := &UserValidator{}

	tests := []struct {
		name     string
		updates  *usersModel.User
		wantErr  bool
		errCount int
	}{
		{
			"所有字段都为空",
			&usersModel.User{Username: "", Email: "", Password: ""},
			false, 0,
		},
		{
			"只更新用户名-有效",
			&usersModel.User{Username: "validuser", Email: "", Password: ""},
			false, 0,
		},
		{
			"只更新用户名-无效",
			&usersModel.User{Username: "ab", Email: "", Password: ""},
			true, 1,
		},
		{
			"只更新邮箱-有效",
			&usersModel.User{Username: "", Email: "valid@example.com", Password: ""},
			false, 0,
		},
		{
			"只更新邮箱-无效",
			&usersModel.User{Username: "", Email: "invalid", Password: ""},
			true, 1,
		},
		{
			"只更新密码-有效",
			&usersModel.User{Username: "", Email: "", Password: "ValidPass123"},
			false, 0,
		},
		{
			"只更新密码-无效",
			&usersModel.User{Username: "", Email: "", Password: "weak"},
			true, 1,
		},
		{
			"更新所有字段都有效",
			&usersModel.User{Username: "validuser", Email: "valid@example.com", Password: "ValidPass123"},
			false, 0,
		},
		{
			"更新所有字段都无效",
			&usersModel.User{Username: "ab", Email: "invalid", Password: "weak"},
			true, 3,
		},
		{
			"部分字段有效部分无效",
			&usersModel.User{Username: "validuser", Email: "invalid", Password: "ValidPass123"},
			true, 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := validator.validateUpdateFields(tt.updates)
			if tt.wantErr {
				assert.Greater(t, len(errs), 0)
				if tt.errCount > 0 {
					assert.Equal(t, tt.errCount, len(errs))
				}
			} else {
				assert.Equal(t, 0, len(errs))
			}
		})
	}
}

// =========================
// validateBusinessRules 测试
// =========================

func TestValidateBusinessRules(t *testing.T) {
	validator := &UserValidator{}

	tests := []struct {
		name     string
		user     *usersModel.User
		wantErr  bool
		errCount int
	}{
		{
			"有效的业务规则",
			&usersModel.User{Username: "testuser", Email: "different@example.com"},
			false, 0,
		},
		{
			"用户名与邮箱前缀相同",
			&usersModel.User{Username: "testuser", Email: "testuser@example.com"},
			true, 1,
		},
		{
			"用户名与邮箱前缀相同(忽略大小写)",
			&usersModel.User{Username: "TestUser", Email: "testuser@example.com"},
			true, 1,
		},
		{
			"创建时间是未来时间",
			func() *usersModel.User {
				u := &usersModel.User{Username: "testuser", Email: "different@example.com"}
				u.CreatedAt = time.Now().Add(24 * time.Hour)
				return u
			}(),
			true, 1,
		},
		{
			"创建时间为零值",
			func() *usersModel.User {
				u := &usersModel.User{Username: "testuser", Email: "different@example.com"}
				u.CreatedAt = time.Time{}
				return u
			}(),
			false, 0,
		},
		{
			"创建时间是过去时间",
			func() *usersModel.User {
				u := &usersModel.User{Username: "testuser", Email: "different@example.com"}
				u.CreatedAt = time.Now().Add(-24 * time.Hour)
				return u
			}(),
			false, 0,
		},
		{
			"多个业务规则违规",
			func() *usersModel.User {
				u := &usersModel.User{Username: "testuser", Email: "testuser@example.com"}
				u.CreatedAt = time.Now().Add(24 * time.Hour)
				return u
			}(),
			true, 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := validator.validateBusinessRules(tt.user)
			if tt.wantErr {
				assert.Greater(t, len(errs), 0)
				if tt.errCount > 0 {
					assert.Equal(t, tt.errCount, len(errs))
				}
			} else {
				assert.Equal(t, 0, len(errs))
			}
		})
	}
}

// =========================
// ValidateUserStatus 测试
// =========================

func TestValidateUserStatus(t *testing.T) {
	validator := &UserValidator{}

	tests := []struct {
		name     string
		status   string
		wantErr  bool
		errCode  string
		errField string
	}{
		{"有效状态active", "active", false, "", ""},
		{"有效状态inactive", "inactive", false, "", ""},
		{"有效状态suspended", "suspended", false, "", ""},
		{"有效状态deleted", "deleted", false, "", ""},
		{"无效状态pending", "pending", true, "INVALID_STATUS", "status"},
		{"无效状态banned", "banned", true, "INVALID_STATUS", "status"},
		{"空状态", "", true, "INVALID_STATUS", "status"},
		{"大小写敏感", "Active", true, "INVALID_STATUS", "status"},
		{"包含空格", " active ", true, "INVALID_STATUS", "status"},
		{"包含特殊字符", "active-", true, "INVALID_STATUS", "status"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateUserStatus(tt.status)
			if tt.wantErr {
				assert.Error(t, err)
				if err != nil {
					assert.Equal(t, tt.errCode, err.Code)
					assert.Equal(t, tt.errField, err.Field)
				}
			} else {
				// ValidateUserStatus 返回 *ValidationError 类型，需要使用 Nil 而不是 NoError
				assert.Nil(t, err)
			}
		})
	}
}

// =========================
// ValidateUserID 测试
// =========================

func TestValidateUserID(t *testing.T) {
	validator := &UserValidator{}
	validObjectID := primitive.NewObjectID().Hex()

	tests := []struct {
		name     string
		userID   string
		wantErr  bool
		errCode  string
		errField string
	}{
		{"有效用户ID", validObjectID, false, "", ""},
		{"另一个有效用户ID", primitive.NewObjectID().Hex(), false, "", ""},
		{"空用户ID", "", true, "REQUIRED", "user_id"},
		{"无效格式的用户ID", "invalid-id", true, "INVALID_FORMAT", "user_id"},
		{"太短的用户ID", "abc123", true, "INVALID_FORMAT", "user_id"},
		{"包含非法字符", "507f1f77bcf86cd799439011-extra", true, "INVALID_FORMAT", "user_id"},
		{"只包含数字", "123456789012345678901234", true, "INVALID_FORMAT", "user_id"},
		{"只包含字母", "abcdefabcdefabcdefabcdef", true, "INVALID_FORMAT", "user_id"},
		{"长度不正确", "507f1f77bcf86cd79943901", true, "INVALID_FORMAT", "user_id"},
		{"包含大写字母", "507F1F77BCF86CD799439011", true, "INVALID_FORMAT", "user_id"},
		{"包含空格", "507f1f77bcf86cd799439011 ", true, "INVALID_FORMAT", "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateUserID(tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
				if err != nil {
					assert.Equal(t, tt.errCode, err.Code)
					assert.Equal(t, tt.errField, err.Field)
				}
			} else {
				// ValidateUserID 返回 *ValidationError 类型，需要使用 Nil 而不是 NoError
				assert.Nil(t, err)
			}
		})
	}
}

// =========================
// ValidationError.Error() 测试
// =========================

func TestValidationError_Error(t *testing.T) {
	tests := []struct {
		name   string
		err    ValidationError
		expect string
	}{
		{
			"标准错误格式",
			ValidationError{Field: "username", Message: "用户名不能为空", Code: "REQUIRED"},
			"验证失败 [username]: 用户名不能为空",
		},
		{
			"邮箱错误",
			ValidationError{Field: "email", Message: "邮箱格式不正确", Code: "INVALID_FORMAT"},
			"验证失败 [email]: 邮箱格式不正确",
		},
		{
			"密码错误",
			ValidationError{Field: "password", Message: "密码过于简单", Code: "WEAK_PASSWORD"},
			"验证失败 [password]: 密码过于简单",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, tt.err.Error())
		})
	}
}

// =========================
// ValidationErrors.Error() 测试
// =========================

func TestValidationErrors_Error(t *testing.T) {
	tests := []struct {
		name   string
		errs   ValidationErrors
		expect string
	}{
		{"空错误列表", ValidationErrors{}, ""},
		{"nil错误列表", nil, ""},
		{
			"单个错误",
			ValidationErrors{{Field: "username", Message: "用户名不能为空", Code: "REQUIRED"}},
			"验证失败 [username]: 用户名不能为空",
		},
		{
			"两个错误",
			ValidationErrors{
				{Field: "username", Message: "用户名不能为空", Code: "REQUIRED"},
				{Field: "email", Message: "邮箱格式不正确", Code: "INVALID_FORMAT"},
			},
			"验证失败 [username]: 用户名不能为空; 验证失败 [email]: 邮箱格式不正确",
		},
		{
			"三个错误",
			ValidationErrors{
				{Field: "username", Message: "用户名不能为空", Code: "REQUIRED"},
				{Field: "email", Message: "邮箱格式不正确", Code: "INVALID_FORMAT"},
				{Field: "password", Message: "密码过于简单", Code: "WEAK_PASSWORD"},
			},
			"验证失败 [username]: 用户名不能为空; 验证失败 [email]: 邮箱格式不正确; 验证失败 [password]: 密码过于简单",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, tt.errs.Error())
		})
	}
}

// =========================
// 集成测试标记
// =========================

func TestUserValidator_ValidateCreate_Integration(t *testing.T) {
	t.Skip("需要数据库连接，在集成测试中实现")
}

func TestUserValidator_ValidateUpdate_Integration(t *testing.T) {
	t.Skip("需要数据库连接，在集成测试中实现")
}

func TestUserValidator_CheckUsernameUnique_Integration(t *testing.T) {
	t.Skip("需要数据库连接，在集成测试中实现")
}

func TestUserValidator_CheckEmailUnique_Integration(t *testing.T) {
	t.Skip("需要数据库连接，在集成测试中实现")
}

func TestUserValidator_CheckUsernameUniqueExcluding_Integration(t *testing.T) {
	t.Skip("需要数据库连接，在集成测试中实现")
}

func TestUserValidator_CheckEmailUniqueExcluding_Integration(t *testing.T) {
	t.Skip("需要数据库连接，在集成测试中实现")
}
