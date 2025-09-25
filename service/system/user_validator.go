package system

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"Qingyu_backend/models/system"
)

// UserValidator 用户数据验证器
type UserValidator struct {
	db *mongo.Database
}

// NewUserValidator 创建用户验证器实例
func NewUserValidator(db *mongo.Database) *UserValidator {
	return &UserValidator{
		db: db,
	}
}

// ValidationError 验证错误类型
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("验证失败 [%s]: %s", e.Field, e.Message)
}

// ValidationErrors 多个验证错误
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}
	var messages []string
	for _, err := range e {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// ValidateCreate 验证用户创建数据
func (v *UserValidator) ValidateCreate(ctx context.Context, user *system.User) error {
	var errors ValidationErrors

	// 1. 基础字段验证
	if fieldErrors := v.validateBasicFields(user); len(fieldErrors) > 0 {
		errors = append(errors, fieldErrors...)
	}

	// 2. 唯一性验证
	if uniqueErrors := v.validateUniqueness(ctx, user); len(uniqueErrors) > 0 {
		errors = append(errors, uniqueErrors...)
	}

	// 3. 业务规则验证
	if businessErrors := v.validateBusinessRules(user); len(businessErrors) > 0 {
		errors = append(errors, businessErrors...)
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

// ValidateUpdate 验证用户更新数据
func (v *UserValidator) ValidateUpdate(ctx context.Context, userID string, updates *system.User) error {
	var errors ValidationErrors

	// 1. 基础字段验证（排除不可更新字段）
	if fieldErrors := v.validateUpdateFields(updates); len(fieldErrors) > 0 {
		errors = append(errors, fieldErrors...)
	}

	// 2. 唯一性验证（排除当前用户）
	if uniqueErrors := v.validateUniquenessForUpdate(ctx, userID, updates); len(uniqueErrors) > 0 {
		errors = append(errors, uniqueErrors...)
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

// validateBasicFields 验证基础字段
func (v *UserValidator) validateBasicFields(user *system.User) ValidationErrors {
	var errors ValidationErrors

	// 验证用户名
	if err := v.validateUsername(user.Username); err != nil {
		errors = append(errors, *err)
	}

	// 验证邮箱
	if err := v.validateEmail(user.Email); err != nil {
		errors = append(errors, *err)
	}

	// 验证密码（仅在创建时）
	if user.Password != "" {
		if err := v.validatePassword(user.Password); err != nil {
			errors = append(errors, *err)
		}
	}

	return errors
}

// validateUpdateFields 验证更新字段
func (v *UserValidator) validateUpdateFields(updates *system.User) ValidationErrors {
	var errors ValidationErrors

	// 验证用户名（如果提供）
	if updates.Username != "" {
		if err := v.validateUsername(updates.Username); err != nil {
			errors = append(errors, *err)
		}
	}

	// 验证邮箱（如果提供）
	if updates.Email != "" {
		if err := v.validateEmail(updates.Email); err != nil {
			errors = append(errors, *err)
		}
	}

	// 验证密码（如果提供）
	if updates.Password != "" {
		if err := v.validatePassword(updates.Password); err != nil {
			errors = append(errors, *err)
		}
	}

	return errors
}

// validateUsername 验证用户名
func (v *UserValidator) validateUsername(username string) *ValidationError {
	if username == "" {
		return &ValidationError{
			Field:   "username",
			Message: "用户名不能为空",
			Code:    "REQUIRED",
		}
	}

	if len(username) < 3 {
		return &ValidationError{
			Field:   "username",
			Message: "用户名长度不能少于3个字符",
			Code:    "MIN_LENGTH",
		}
	}

	if len(username) > 30 {
		return &ValidationError{
			Field:   "username",
			Message: "用户名长度不能超过30个字符",
			Code:    "MAX_LENGTH",
		}
	}

	// 用户名只能包含字母、数字和下划线
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	if !matched {
		return &ValidationError{
			Field:   "username",
			Message: "用户名只能包含字母、数字和下划线",
			Code:    "INVALID_FORMAT",
		}
	}

	// 不能以数字开头
	if username[0] >= '0' && username[0] <= '9' {
		return &ValidationError{
			Field:   "username",
			Message: "用户名不能以数字开头",
			Code:    "INVALID_FORMAT",
		}
	}

	// 检查保留用户名
	reservedNames := []string{"admin", "root", "system", "api", "www", "mail", "ftp"}
	for _, reserved := range reservedNames {
		if strings.ToLower(username) == reserved {
			return &ValidationError{
				Field:   "username",
				Message: "该用户名为系统保留，不能使用",
				Code:    "RESERVED_NAME",
			}
		}
	}

	return nil
}

// validateEmail 验证邮箱
func (v *UserValidator) validateEmail(email string) *ValidationError {
	if email == "" {
		return &ValidationError{
			Field:   "email",
			Message: "邮箱不能为空",
			Code:    "REQUIRED",
		}
	}

	if len(email) > 100 {
		return &ValidationError{
			Field:   "email",
			Message: "邮箱长度不能超过100个字符",
			Code:    "MAX_LENGTH",
		}
	}

	// 邮箱格式验证
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, email)
	if !matched {
		return &ValidationError{
			Field:   "email",
			Message: "邮箱格式不正确",
			Code:    "INVALID_FORMAT",
		}
	}

	return nil
}

// validatePassword 验证密码
func (v *UserValidator) validatePassword(password string) *ValidationError {
	if password == "" {
		return &ValidationError{
			Field:   "password",
			Message: "密码不能为空",
			Code:    "REQUIRED",
		}
	}

	if len(password) < 8 {
		return &ValidationError{
			Field:   "password",
			Message: "密码长度不能少于8个字符",
			Code:    "MIN_LENGTH",
		}
	}

	if len(password) > 128 {
		return &ValidationError{
			Field:   "password",
			Message: "密码长度不能超过128个字符",
			Code:    "MAX_LENGTH",
		}
	}

	// 密码强度验证：至少包含字母和数字
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasLetter || !hasNumber {
		return &ValidationError{
			Field:   "password",
			Message: "密码必须包含至少一个字母和一个数字",
			Code:    "WEAK_PASSWORD",
		}
	}

	// 检查常见弱密码
	weakPasswords := []string{"12345678", "password", "qwerty123", "abc123456"}
	for _, weak := range weakPasswords {
		if strings.ToLower(password) == weak {
			return &ValidationError{
				Field:   "password",
				Message: "密码过于简单，请使用更复杂的密码",
				Code:    "WEAK_PASSWORD",
			}
		}
	}

	return nil
}

// validateUniqueness 验证唯一性约束
func (v *UserValidator) validateUniqueness(ctx context.Context, user *system.User) ValidationErrors {
	var errors ValidationErrors

	// 检查用户名唯一性
	if err := v.checkUsernameUnique(ctx, user.Username); err != nil {
		errors = append(errors, *err)
	}

	// 检查邮箱唯一性
	if err := v.checkEmailUnique(ctx, user.Email); err != nil {
		errors = append(errors, *err)
	}

	return errors
}

// validateUniquenessForUpdate 验证更新时的唯一性约束
func (v *UserValidator) validateUniquenessForUpdate(ctx context.Context, userID string, updates *system.User) ValidationErrors {
	var errors ValidationErrors

	// 检查用户名唯一性（排除当前用户）
	if updates.Username != "" {
		if err := v.checkUsernameUniqueExcluding(ctx, updates.Username, userID); err != nil {
			errors = append(errors, *err)
		}
	}

	// 检查邮箱唯一性（排除当前用户）
	if updates.Email != "" {
		if err := v.checkEmailUniqueExcluding(ctx, updates.Email, userID); err != nil {
			errors = append(errors, *err)
		}
	}

	return errors
}

// checkUsernameUnique 检查用户名唯一性
func (v *UserValidator) checkUsernameUnique(ctx context.Context, username string) *ValidationError {
	count, err := v.db.Collection("users").CountDocuments(ctx, bson.M{
		"username": username,
	})
	if err != nil {
		return &ValidationError{
			Field:   "username",
			Message: "验证用户名唯一性时发生错误",
			Code:    "VALIDATION_ERROR",
		}
	}

	if count > 0 {
		return &ValidationError{
			Field:   "username",
			Message: "用户名已存在",
			Code:    "DUPLICATE",
		}
	}

	return nil
}

// checkEmailUnique 检查邮箱唯一性
func (v *UserValidator) checkEmailUnique(ctx context.Context, email string) *ValidationError {
	count, err := v.db.Collection("users").CountDocuments(ctx, bson.M{
		"email": email,
	})
	if err != nil {
		return &ValidationError{
			Field:   "email",
			Message: "验证邮箱唯一性时发生错误",
			Code:    "VALIDATION_ERROR",
		}
	}

	if count > 0 {
		return &ValidationError{
			Field:   "email",
			Message: "邮箱已被使用",
			Code:    "DUPLICATE",
		}
	}

	return nil
}

// checkUsernameUniqueExcluding 检查用户名唯一性（排除指定用户）
func (v *UserValidator) checkUsernameUniqueExcluding(ctx context.Context, username, excludeUserID string) *ValidationError {
	objectID, err := primitive.ObjectIDFromHex(excludeUserID)
	if err != nil {
		return &ValidationError{
			Field:   "username",
			Message: "用户ID格式错误",
			Code:    "INVALID_ID",
		}
	}

	count, err := v.db.Collection("users").CountDocuments(ctx, bson.M{
		"username": username,
		"_id":      bson.M{"$ne": objectID},
	})
	if err != nil {
		return &ValidationError{
			Field:   "username",
			Message: "验证用户名唯一性时发生错误",
			Code:    "VALIDATION_ERROR",
		}
	}

	if count > 0 {
		return &ValidationError{
			Field:   "username",
			Message: "用户名已存在",
			Code:    "DUPLICATE",
		}
	}

	return nil
}

// checkEmailUniqueExcluding 检查邮箱唯一性（排除指定用户）
func (v *UserValidator) checkEmailUniqueExcluding(ctx context.Context, email, excludeUserID string) *ValidationError {
	objectID, err := primitive.ObjectIDFromHex(excludeUserID)
	if err != nil {
		return &ValidationError{
			Field:   "email",
			Message: "用户ID格式错误",
			Code:    "INVALID_ID",
		}
	}

	count, err := v.db.Collection("users").CountDocuments(ctx, bson.M{
		"email": email,
		"_id":   bson.M{"$ne": objectID},
	})
	if err != nil {
		return &ValidationError{
			Field:   "email",
			Message: "验证邮箱唯一性时发生错误",
			Code:    "VALIDATION_ERROR",
		}
	}

	if count > 0 {
		return &ValidationError{
			Field:   "email",
			Message: "邮箱已被使用",
			Code:    "DUPLICATE",
		}
	}

	return nil
}

// validateBusinessRules 验证业务规则
func (v *UserValidator) validateBusinessRules(user *system.User) ValidationErrors {
	var errors ValidationErrors

	// 检查用户名和邮箱不能相同
	if user.Username != "" && user.Email != "" {
		if strings.ToLower(user.Username) == strings.ToLower(strings.Split(user.Email, "@")[0]) {
			errors = append(errors, ValidationError{
				Field:   "username",
				Message: "用户名不能与邮箱前缀相同",
				Code:    "BUSINESS_RULE_VIOLATION",
			})
		}
	}

	// 检查创建时间不能是未来时间
	if !user.CreatedAt.IsZero() && user.CreatedAt.After(time.Now()) {
		errors = append(errors, ValidationError{
			Field:   "created_at",
			Message: "创建时间不能是未来时间",
			Code:    "INVALID_TIME",
		})
	}

	return errors
}

// ValidateUserStatus 验证用户状态
func (v *UserValidator) ValidateUserStatus(status string) *ValidationError {
	validStatuses := []string{"active", "inactive", "suspended", "deleted"}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return nil
		}
	}

	return &ValidationError{
		Field:   "status",
		Message: "无效的用户状态",
		Code:    "INVALID_STATUS",
	}
}

// ValidateUserID 验证用户ID格式
func (v *UserValidator) ValidateUserID(userID string) *ValidationError {
	if userID == "" {
		return &ValidationError{
			Field:   "user_id",
			Message: "用户ID不能为空",
			Code:    "REQUIRED",
		}
	}

	if _, err := primitive.ObjectIDFromHex(userID); err != nil {
		return &ValidationError{
			Field:   "user_id",
			Message: "用户ID格式错误",
			Code:    "INVALID_FORMAT",
		}
	}

	return nil
}
