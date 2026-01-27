package auth

import (
	"context"
	"fmt"
	"strings"
)

// Permission 权限定义
//
// 表示对特定资源的操作权限
type Permission struct {
	// Resource 资源类型（如：project, document, user）
	Resource string `json:"resource"`

	// Action 操作类型（如：read, write, delete, admin）
	Action string `json:"action"`

	// ResourceID 资源ID（可选，用于细粒度权限控制）
	// 例如：只允许删除自己创建的文档
	ResourceID string `json:"resource_id,omitempty"`
}

// String 返回权限的字符串表示
func (p Permission) String() string {
	if p.ResourceID != "" {
		return fmt.Sprintf("%s:%s:%s", p.Resource, p.Action, p.ResourceID)
	}
	return fmt.Sprintf("%s:%s", p.Resource, p.Action)
}

// ParsePermission 从字符串解析权限
//
// 支持格式：
// - "resource:action" -> Resource: resource, Action: action
// - "resource:action:id" -> Resource: resource, Action: action, ResourceID: id
// - "*:*" -> 通配符权限
func ParsePermission(permStr string) (Permission, error) {
	parts := strings.Split(permStr, ":")

	if len(parts) < 2 {
		return Permission{}, fmt.Errorf("invalid permission format: %s", permStr)
	}

	perm := Permission{
		Resource: parts[0],
	}

	if len(parts) >= 2 {
		perm.Action = parts[1]
	}

	if len(parts) >= 3 {
		perm.ResourceID = parts[2]
	}

	return perm, nil
}

// Checker 权限检查器接口
//
// 定义了权限检查的核心方法，支持多种实现（RBAC、Casbin等）
type Checker interface {
	// Name 返回检查器名称
	//
	// 例如: "rbac", "casbin"
	Name() string

	// Check 检查权限
	//
	// 参数:
	//   - ctx: 上下文
	//   - subject: 主体（通常是用户ID）
	//   - perm: 权限定义
	//
	// 返回:
	//   - bool: 是否有权限
	//   - error: 错误信息
	Check(ctx context.Context, subject string, perm Permission) (bool, error)

	// BatchCheck 批量检查权限
	//
	// 比单独调用Check更高效，因为可以批量查询
	//
	// 参数:
	//   - ctx: 上下文
	//   - subject: 主体（通常是用户ID）
	//   - perms: 权限列表
	//
	// 返回:
	//   - []bool: 每个权限的检查结果
	//   - error: 错误信息
	BatchCheck(ctx context.Context, subject string, perms []Permission) ([]bool, error)

	// Close 关闭检查器
	//
	// 释放资源（如数据库连接、缓存等）
	Close() error
}

// CheckerConfig 检查器配置
type CheckerConfig struct {
	// Strategy 检查器策略（rbac, casbin）
	Strategy string `json:"strategy" yaml:"strategy"`

	// ConfigPath 配置文件路径（用于加载权限规则）
	ConfigPath string `json:"config_path" yaml:"config_path"`
}

// CheckerFactory 检查器工厂函数
//
// 根据配置创建对应的Checker实例
type CheckerFactory func(config *CheckerConfig) (Checker, error)

// registry 检查器注册表
var registry = make(map[string]CheckerFactory)

// RegisterChecker 注册检查器工厂
//
// 例如:
//
//	auth.RegisterChecker("rbac", NewRBACChecker)
//	auth.RegisterChecker("casbin", NewCasbinChecker)
func RegisterChecker(strategy string, factory CheckerFactory) {
	registry[strategy] = factory
}

// CreateChecker 根据配置创建检查器
func CreateChecker(config *CheckerConfig) (Checker, error) {
	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}

	factory, exists := registry[config.Strategy]
	if !exists {
		return nil, fmt.Errorf("unsupported checker strategy: %s", config.Strategy)
	}

	return factory(config)
}

// NoOpChecker 空检查器（用于测试）
//
// 总是返回true
type NoOpChecker struct{}

// NewNoOpChecker 创建空检查器
func NewNoOpChecker(config *CheckerConfig) (Checker, error) {
	return &NoOpChecker{}, nil
}

// Name 返回检查器名称
func (c *NoOpChecker) Name() string {
	return "noop"
}

// Check 检查权限（总是返回true）
func (c *NoOpChecker) Check(ctx context.Context, subject string, perm Permission) (bool, error) {
	return true, nil
}

// BatchCheck 批量检查权限（总是返回true）
func (c *NoOpChecker) BatchCheck(ctx context.Context, subject string, perms []Permission) ([]bool, error) {
	results := make([]bool, len(perms))
	for i := range results {
		results[i] = true
	}
	return results, nil
}

// Close 关闭检查器
func (c *NoOpChecker) Close() error {
	return nil
}

func init() {
	// 注册空检查器
	RegisterChecker("noop", NewNoOpChecker)
}
