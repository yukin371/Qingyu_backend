package user

import (
	usersModel "Qingyu_backend/models/users"
	base "Qingyu_backend/repository/interfaces/infrastructure"
	writingRepoInterface "Qingyu_backend/repository/interfaces/writer"
	"context"
	"fmt"
	"time"
)

// UserFilter 用户查询过滤器
type UserFilter struct {
	ID       string    `json:"id,omitempty"`
	Username string    `json:"username,omitempty"`
	Email    string    `json:"email,omitempty"`
	Status   string    `json:"status,omitempty"`
	FromDate time.Time `json:"from_date,omitempty"`
	ToDate   time.Time `json:"to_date,omitempty"`
	Limit    int64     `json:"limit,omitempty"`
	Offset   int64     `json:"offset,omitempty"`
}

// GetConditions 获取查询条件
func (f UserFilter) GetConditions() map[string]interface{} {
	conditions := make(map[string]interface{})

	if f.ID != "" {
		conditions["id"] = f.ID
	}
	if f.Username != "" {
		conditions["username"] = f.Username
	}
	if f.Email != "" {
		conditions["email"] = f.Email
	}
	if f.Status != "" {
		conditions["status"] = f.Status
	}
	if !f.FromDate.IsZero() {
		conditions["from_date"] = f.FromDate
	}
	if !f.ToDate.IsZero() {
		conditions["to_date"] = f.ToDate
	}

	return conditions
}

// GetSort 获取排序条件
func (f UserFilter) GetSort() map[string]int {
	// 默认按创建时间倒序排列
	return map[string]int{
		"created_at": -1,
	}
}

// GetFields 获取字段选择
func (f UserFilter) GetFields() []string {
	// 返回空切片表示选择所有字段
	return []string{}
}

// Validate 验证过滤器
func (f UserFilter) Validate() error {
	if f.Limit < 0 {
		return NewUserRepositoryError(ErrorTypeValidation, "Limit不能为负数", nil)
	}
	if f.Offset < 0 {
		return NewUserRepositoryError(ErrorTypeValidation, "Offset不能为负数", nil)
	}
	if !f.FromDate.IsZero() && !f.ToDate.IsZero() && f.FromDate.After(f.ToDate) {
		return NewUserRepositoryError(ErrorTypeValidation, "开始时间不能晚于结束时间", nil)
	}
	return nil
}

// UserRepository 用户仓储接口
// 继承BaseRepository的通用CRUD操作，并添加用户特定的业务方法
type UserRepository interface {
	// 继承CRUDRepository接口
	base.CRUDRepository[*usersModel.User, string]
	// 继承 HealthRepository 接口
	base.HealthRepository
	// 继承BaseUserRepository共享基础接口
	BaseUserRepository

	// === user模块特有方法 ===

	// 用户特定的查询方法（user特有）
	GetByPhone(ctx context.Context, phone string) (*usersModel.User, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)

	// 用户状态管理（user特有）
	UpdateLastLogin(ctx context.Context, id string, ip string) error
	UpdatePasswordByEmail(ctx context.Context, email string, hashedPassword string) error
	GetActiveUsers(ctx context.Context, limit int64) ([]*usersModel.User, error)
	GetUsersByRole(ctx context.Context, role string, limit int64) ([]*usersModel.User, error)

	// 验证状态管理（user特有）
	SetPhoneVerified(ctx context.Context, id string, verified bool) error
	UnbindEmail(ctx context.Context, id string) error
	UnbindPhone(ctx context.Context, id string) error

	// 设备管理（user特有）
	DeleteDevice(ctx context.Context, userID string, deviceID string) error
	GetDevices(ctx context.Context, userID string) ([]interface{}, error)

	// 高级查询（user特有）
	FindWithFilter(ctx context.Context, filter *usersModel.UserFilter) ([]*usersModel.User, int64, error)
	SearchUsers(ctx context.Context, keyword string, limit int) ([]*usersModel.User, error)

	// 事务操作（user特有）
	Transaction(ctx context.Context,
		user *usersModel.User,
		fn func(ctx context.Context, repo UserRepository) error) error
}

// 注意：以下方法已通过 BaseUserRepository 提供，不再重复定义：
// - GetByUsername(ctx, string) (*usersModel.User, error)
// - GetByEmail(ctx, string) (*usersModel.User, error)
// - UpdateStatus(ctx, string, UserStatus) error
// - UpdatePassword(ctx, string, string) error
// - SetEmailVerified(ctx, string, bool) error
// - BatchUpdateStatus(ctx, []string, UserStatus) error
// - BatchDelete(ctx, []string) error
// - CountByRole(ctx, string) (int64, error)
// - CountByStatus(ctx, UserStatus) (int64, error)

// RepositoryFactory 仓储工厂接口
type RepositoryFactory interface {
	CreateUserRepository() UserRepository
	CreateProjectRepository() writingRepoInterface.ProjectRepository
	CreateRoleRepository() RoleRepository
	Close() error
	Health(ctx context.Context) error
	GetDatabaseType() string
}

// DatabaseConfig 数据库配置接口
type DatabaseConfig interface {
	GetConnectionString() string
	GetDatabaseName() string
	GetMaxConnections() int
	GetTimeout() time.Duration
	Validate() error
}

// MongoConfig MongoDB配置
type MongoConfig struct {
	URI            string        `yaml:"uri" json:"uri"`
	Database       string        `yaml:"database" json:"database"`
	MaxPoolSize    uint64        `yaml:"max_pool_size" json:"max_pool_size"`
	MinPoolSize    uint64        `yaml:"min_pool_size" json:"min_pool_size"`
	ConnectTimeout time.Duration `yaml:"connect_timeout" json:"connect_timeout"`
	ServerTimeout  time.Duration `yaml:"server_timeout" json:"server_timeout"`
}

func (c *MongoConfig) GetConnectionString() string {
	return c.URI
}

func (c *MongoConfig) GetDatabaseName() string {
	return c.Database
}

func (c *MongoConfig) GetMaxConnections() int {
	return int(c.MaxPoolSize)
}

func (c *MongoConfig) GetTimeout() time.Duration {
	return c.ConnectTimeout
}

func (c *MongoConfig) Validate() error {
	if c.URI == "" {
		return NewUserRepositoryError(ErrorTypeValidation, "MongoDB URI is required", nil)
	}
	if c.Database == "" {
		return NewUserRepositoryError(ErrorTypeValidation, "MongoDB database name is required", nil)
	}
	if c.MaxPoolSize == 0 {
		c.MaxPoolSize = 100
	}
	if c.MinPoolSize == 0 {
		c.MinPoolSize = 5
	}
	if c.ConnectTimeout == 0 {
		c.ConnectTimeout = 10 * time.Second
	}
	if c.ServerTimeout == 0 {
		c.ServerTimeout = 30 * time.Second
	}
	return nil
}

// PostgreSQLConfig PostgreSQL配置
type PostgreSQLConfig struct {
	Host         string        `yaml:"host" json:"host"`
	Port         int           `yaml:"port" json:"port"`
	Database     string        `yaml:"database" json:"database"`
	Username     string        `yaml:"username" json:"username"`
	Password     string        `yaml:"password" json:"password"`
	SSLMode      string        `yaml:"ssl_mode" json:"ssl_mode"`
	MaxOpenConns int           `yaml:"max_open_conns" json:"max_open_conns"`
	MaxIdleConns int           `yaml:"max_idle_conns" json:"max_idle_conns"`
	ConnTimeout  time.Duration `yaml:"conn_timeout" json:"conn_timeout"`
}

func (c *PostgreSQLConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Password, c.Database, c.SSLMode)
}

func (c *PostgreSQLConfig) GetDatabaseName() string {
	return c.Database
}

func (c *PostgreSQLConfig) GetMaxConnections() int {
	return c.MaxOpenConns
}

func (c *PostgreSQLConfig) GetTimeout() time.Duration {
	return c.ConnTimeout
}

func (c *PostgreSQLConfig) Validate() error {
	if c.Host == "" {
		return NewUserRepositoryError(ErrorTypeValidation, "PostgreSQL host is required", nil)
	}
	if c.Port == 0 {
		c.Port = 5432
	}
	if c.Database == "" {
		return NewUserRepositoryError(ErrorTypeValidation, "PostgreSQL database name is required", nil)
	}
	if c.Username == "" {
		return NewUserRepositoryError(ErrorTypeValidation, "PostgreSQL username is required", nil)
	}
	if c.SSLMode == "" {
		c.SSLMode = "disable"
	}
	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = 25
	}
	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = 5
	}
	if c.ConnTimeout == 0 {
		c.ConnTimeout = 10 * time.Second
	}
	return nil
}

// 数据库类型常量
const (
	DatabaseTypeMongoDB    = "mongodb"
	DatabaseTypePostgreSQL = "postgresql"
	DatabaseTypeMySQL      = "mysql"
)

// RepositoryError 仓储错误类型
type UserRepositoryError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Cause   error  `json:"cause,omitempty"`
}

func (e *UserRepositoryError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// 错误类型常量
const (
	ErrorTypeNotFound    = "NOT_FOUND"
	ErrorTypeDuplicate   = "DUPLICATE"
	ErrorTypeValidation  = "VALIDATION"
	ErrorTypeConnection  = "CONNECTION"
	ErrorTypeTimeout     = "TIMEOUT"
	ErrorTypeInternal    = "INTERNAL"
	ErrorTypeTransaction = "TRANSACTION"
)

// NewRepositoryError 创建仓储错误
func NewUserRepositoryError(errorType, message string, cause error) *UserRepositoryError {
	return &UserRepositoryError{
		Type:    errorType,
		Message: message,
		Cause:   cause,
	}
}

// IsNotFoundError 检查是否为未找到错误
func IsNotFoundError(err error) bool {
	if repoErr, ok := err.(*UserRepositoryError); ok {
		return repoErr.Type == ErrorTypeNotFound
	}
	return false
}

// IsDuplicateError 检查是否为重复错误
func IsDuplicateError(err error) bool {
	if repoErr, ok := err.(*UserRepositoryError); ok {
		return repoErr.Type == ErrorTypeDuplicate
	}
	return false
}

// IsValidationError 检查是否为验证错误
func IsValidationError(err error) bool {
	if repoErr, ok := err.(*UserRepositoryError); ok {
		return repoErr.Type == ErrorTypeValidation
	}
	return false
}
