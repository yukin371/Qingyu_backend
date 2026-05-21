package document

import (
	"Qingyu_backend/models/writer"
	"context"
	"strings"
	"time"

	pkgErrors "Qingyu_backend/pkg/errors"
)

// ShortcutService 快捷键服务
type ShortcutService struct {
	serviceName string
}

// NewShortcutService 创建快捷键服务
func NewShortcutService() *ShortcutService {
	return &ShortcutService{
		serviceName: "ShortcutService",
	}
}

// GetUserShortcuts 获取用户快捷键配置
func (s *ShortcutService) GetUserShortcuts(ctx context.Context, userID string) (*writer.ShortcutConfig, error) {
	// 1. 参数验证
	if userID == "" {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "用户ID不能为空", "", nil)
	}

	// 2. TODO: 从数据库查询用户配置
	// config, err := s.shortcutRepo.GetByUserID(ctx, userID)
	// if err != nil {
	// 	return nil, err
	// }

	// 3. 如果用户没有自定义配置，返回默认配置
	config := &writer.ShortcutConfig{
		UserID:    userID,
		Shortcuts: writer.GetDefaultShortcuts(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return config, nil
}

// UpdateUserShortcuts 更新用户快捷键配置
func (s *ShortcutService) UpdateUserShortcuts(ctx context.Context, userID string, shortcuts map[string]writer.Shortcut) error {
	// 1. 参数验证
	if userID == "" {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "用户ID不能为空", "", nil)
	}

	if len(shortcuts) == 0 {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "快捷键配置不能为空", "", nil)
	}

	// 2. 验证快捷键配置
	if err := s.validateShortcuts(shortcuts); err != nil {
		return err
	}

	// 3. TODO: 保存到数据库
	// config := &document.ShortcutConfig{
	// 	UserID:    userID,
	// 	Shortcuts: shortcuts,
	// 	UpdatedAt: time.Now(),
	// }
	//
	// if err := s.shortcutRepo.Upsert(ctx, config); err != nil {
	// 	return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "保存快捷键配置失败", "", err)
	// }

	return nil
}

// ResetUserShortcuts 重置用户快捷键为默认配置
func (s *ShortcutService) ResetUserShortcuts(ctx context.Context, userID string) error {
	// 1. 参数验证
	if userID == "" {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "用户ID不能为空", "", nil)
	}

	// 2. TODO: 删除用户自定义配置
	// if err := s.shortcutRepo.Delete(ctx, userID); err != nil {
	// 	return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "重置快捷键配置失败", "", err)
	// }

	return nil
}

// GetShortcutHelp 获取快捷键帮助（按分类）
func (s *ShortcutService) GetShortcutHelp(ctx context.Context, userID string) ([]writer.ShortcutCategory, error) {
	// 1. 获取用户快捷键配置
	config, err := s.GetUserShortcuts(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 2. 按分类整理
	categories := writer.GetShortcutsByCategory(config.Shortcuts)

	return categories, nil
}

// validateShortcuts 验证快捷键配置
func (s *ShortcutService) validateShortcuts(shortcuts map[string]writer.Shortcut) error {
	// 1. 检查是否有重复的按键组合
	usedKeys := make(map[string]string)
	for action, shortcut := range shortcuts {
		if existingAction, exists := usedKeys[shortcut.Key]; exists {
			return pkgErrors.NewServiceError(
				s.serviceName,
				pkgErrors.ServiceErrorValidation,
				"快捷键冲突: 按键 "+shortcut.Key+" 已被 "+existingAction+" 使用",
				"",
				nil,
			)
		}
		usedKeys[shortcut.Key] = action
	}

	// 2. 验证快捷键格式
	for _, shortcut := range shortcuts {
		if shortcut.Key == "" {
			return pkgErrors.NewServiceError(
				s.serviceName,
				pkgErrors.ServiceErrorValidation,
				"快捷键不能为空",
				"操作 "+shortcut.Action+" 的快捷键为空",
				nil,
			)
		}

		if err := validateShortcutKey(shortcut.Key); err != nil {
			return pkgErrors.NewServiceError(
				s.serviceName,
				pkgErrors.ServiceErrorValidation,
				"快捷键格式无效",
				"操作 "+shortcut.Action+" 的快捷键无效: "+err.Error(),
				nil,
			)
		}
	}

	return nil
}

func validateShortcutKey(key string) error {
	parts := strings.Split(key, "+")
	if len(parts) == 0 {
		return pkgErrors.NewServiceError("ShortcutService", pkgErrors.ServiceErrorValidation, "快捷键不能为空", "", nil)
	}

	seenModifiers := make(map[string]struct{})
	for i, rawPart := range parts {
		part := strings.TrimSpace(rawPart)
		if part == "" {
			return pkgErrors.NewServiceError("ShortcutService", pkgErrors.ServiceErrorValidation, "快捷键存在空按键段", "", nil)
		}

		isLast := i == len(parts)-1
		if !isLast {
			if !isModifierKey(part) {
				return pkgErrors.NewServiceError("ShortcutService", pkgErrors.ServiceErrorValidation, "修饰键无效: "+part, "", nil)
			}
			if _, exists := seenModifiers[part]; exists {
				return pkgErrors.NewServiceError("ShortcutService", pkgErrors.ServiceErrorValidation, "修饰键重复: "+part, "", nil)
			}
			seenModifiers[part] = struct{}{}
			continue
		}

		if isModifierKey(part) && len(parts) > 1 {
			return pkgErrors.NewServiceError("ShortcutService", pkgErrors.ServiceErrorValidation, "缺少主按键", "", nil)
		}
		if !isValidPrimaryKey(part) {
			return pkgErrors.NewServiceError("ShortcutService", pkgErrors.ServiceErrorValidation, "主按键无效: "+part, "", nil)
		}
	}

	return nil
}

func isModifierKey(key string) bool {
	switch key {
	case "Ctrl", "Alt", "Shift", "Meta", "Cmd":
		return true
	default:
		return false
	}
}

func isValidPrimaryKey(key string) bool {
	if len([]rune(key)) == 1 {
		return true
	}

	switch key {
	case "Tab", "Enter", "Esc", "Escape", "Space", "Backspace", "Delete", "Insert", "Home", "End", "PageUp", "PageDown", "Up", "Down", "Left", "Right":
		return true
	}

	if len(key) >= 2 && key[0] == 'F' {
		switch key {
		case "F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9", "F10", "F11", "F12":
			return true
		}
	}

	return false
}
