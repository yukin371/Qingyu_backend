package repository

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/models/system"
)

// UserFilter 用户查询过滤器
type UserFilter struct {
	Username string    `json:"username,omitempty"`
	Email    string    `json:"email,omitempty"`
	Status   string    `json:"status,omitempty"`
	FromDate time.Time `json:"from_date,omitempty"`
	ToDate   time.Time `json:"to_date,omitempty"`
	Limit    int64     `json:"limit,omitempty"`
	Offset   int64     `json:"offset,omitempty"`
}

// UserRepository 用户仓储接口
// 定义了用户数据访问的标准接口，支持多种数据库实现
type UserRepository interface {
	// Create 创建用户
	Create(ctx context.Context, user *system.User) error

	// GetByID 根据ID获取用户
	GetByID(ctx context.Context, id string) (*system.User, error)

	// GetByUsername 根据用户名获取用户
	GetByUsername(ctx context.Context, username string) (*system.User, error)

	// GetByEmail 根据邮箱获取用户
	GetByEmail(ctx context.Context, email string) (*system.User, error)

	// Update 更新用户信息
	Update(ctx context.Context, id string, updates map[string]interface{}) error

	// Delete 删除用户（软删除）
	Delete(ctx context.Context, id string) error

	// HardDelete 硬删除用户
	HardDelete(ctx context.Context, id string) error

	// List 获取用户列表
	List(ctx context.Context, filter UserFilter) ([]*system.User, error)

	// Count 统计用户数量
	Count(ctx context.Context, filter UserFilter) (int64, error)

	// Exists 检查用户是否存在
	Exists(ctx context.Context, id string) (bool, error)

	// ExistsByUsername 检查用户名是否存在
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	// ExistsByEmail 检查邮箱是否存在
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// UpdateLastLogin 更新最后登录时间
	UpdateLastLogin(ctx context.Context, id string) error

	// UpdatePassword 更新密码
	UpdatePassword(ctx context.Context, id string, hashedPassword string) error

	// GetActiveUsers 获取活跃用户
	GetActiveUsers(ctx context.Context, limit int64) ([]*system.User, error)

	// BatchUpdate 批量更新用户
	BatchUpdate(ctx context.Context, ids []string, updates map[string]interface{}) error

	// Transaction 执行事务操作
	Transaction(ctx context.Context, fn func(ctx context.Context, repo UserRepository) error) error
}

// ProjectRepository 项目仓储接口
type ProjectRepository interface {
	// Create 创建项目
	Create(ctx context.Context, project interface{}) error

	// GetByID 根据ID获取项目
	GetByID(ctx context.Context, id string) (interface{}, error)

	// GetByCreatorID 根据创建者ID获取项目列表
	GetByCreatorID(ctx context.Context, creatorID string) ([]interface{}, error)

	// Update 更新项目
	Update(ctx context.Context, id string, updates map[string]interface{}) error

	// Delete 删除项目
	Delete(ctx context.Context, id string) error

	// List 获取项目列表
	List(ctx context.Context, filter interface{}) ([]interface{}, error)

	// Count 统计项目数量
	Count(ctx context.Context, filter interface{}) (int64, error)

	// ArchiveByCreatorID 根据创建者ID归档项目
	ArchiveByCreatorID(ctx context.Context, creatorID string) error
}

// RoleRepository 角色仓储接口
type RoleRepository interface {
	// Create 创建角色
	Create(ctx context.Context, role interface{}) error

	// GetByID 根据ID获取角色
	GetByID(ctx context.Context, id string) (interface{}, error)

	// GetByName 根据名称获取角色
	GetByName(ctx context.Context, name string) (interface{}, error)

	// Update 更新角色
	Update(ctx context.Context, id string, updates map[string]interface{}) error

	// Delete 删除角色
	Delete(ctx context.Context, id string) error

	// List 获取角色列表
	List(ctx context.Context) ([]interface{}, error)

	// GetDefaultRole 获取默认角色
	GetDefaultRole(ctx context.Context) (interface{}, error)

	// GetUserRoles 获取用户角色
	GetUserRoles(ctx context.Context, userID string) ([]interface{}, error)

	// AssignRole 分配角色给用户
	AssignRole(ctx context.Context, userID, roleID string) error

	// RemoveRole 移除用户角色
	RemoveRole(ctx context.Context, userID, roleID string) error

	// GetUserPermissions 获取用户权限
	GetUserPermissions(ctx context.Context, userID string) ([]string, error)
}

// RepositoryFactory 仓储工厂接口
// 用于创建不同类型的仓储实例，支持多种数据库实现
type RepositoryFactory interface {
	// CreateUserRepository 创建用户仓储
	CreateUserRepository() UserRepository

	// CreateProjectRepository 创建项目仓储
	CreateProjectRepository() ProjectRepository

	// CreateRoleRepository 创建角色仓储
	CreateRoleRepository() RoleRepository

	// Close 关闭数据库连接
	Close() error

	// Health 检查数据库健康状态
	Health(ctx context.Context) error

	// GetDatabaseType 获取数据库类型
	GetDatabaseType() string
}

// DatabaseConfig 数据库配置接口
type DatabaseConfig interface {
	// GetConnectionString 获取连接字符串
	GetConnectionString() string

	// GetDatabaseName 获取数据库名称
	GetDatabaseName() string

	// GetMaxConnections 获取最大连接数
	GetMaxConnections() int

	// GetTimeout 获取超时时间
	GetTimeout() time.Duration

	// Validate 验证配置
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
		return fmt.Errorf("MongoDB URI不能为空")
	}
	if c.Database == "" {
		return fmt.Errorf("数据库名称不能为空")
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

// PostgreSQLConfig PostgreSQL配置（预留）
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
		return fmt.Errorf("PostgreSQL主机地址不能为空")
	}
	if c.Port == 0 {
		c.Port = 5432
	}
	if c.Database == "" {
		return fmt.Errorf("数据库名称不能为空")
	}
	if c.Username == "" {
		return fmt.Errorf("用户名不能为空")
	}
	if c.SSLMode == "" {
		c.SSLMode = "disable"
	}
	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = 100
	}
	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = 10
	}
	if c.ConnTimeout == 0 {
		c.ConnTimeout = 10 * time.Second
	}
	return nil
}

// DatabaseType 数据库类型常量
const (
	DatabaseTypeMongoDB    = "mongodb"
	DatabaseTypePostgreSQL = "postgresql"
	DatabaseTypeMySQL      = "mysql"
)

// RepositoryError 仓储错误类型
type RepositoryError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Cause   error  `json:"cause,omitempty"`
}

func (e *RepositoryError) Error() string {
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
	ErrorTypeTransaction = "TRANSACTION"
	ErrorTypeInternal    = "INTERNAL"
)

// NewRepositoryError 创建仓储错误
func NewRepositoryError(errorType, message string, cause error) *RepositoryError {
	return &RepositoryError{
		Type:    errorType,
		Message: message,
		Cause:   cause,
	}
}

// IsNotFoundError 检查是否为未找到错误
func IsNotFoundError(err error) bool {
	if repoErr, ok := err.(*RepositoryError); ok {
		return repoErr.Type == ErrorTypeNotFound
	}
	return false
}

// IsDuplicateError 检查是否为重复错误
func IsDuplicateError(err error) bool {
	if repoErr, ok := err.(*RepositoryError); ok {
		return repoErr.Type == ErrorTypeDuplicate
	}
	return false
}

// IsValidationError 检查是否为验证错误
func IsValidationError(err error) bool {
	if repoErr, ok := err.(*RepositoryError); ok {
		return repoErr.Type == ErrorTypeValidation
	}
	return false
}
