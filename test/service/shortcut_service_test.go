package service

import (
	"Qingyu_backend/models/writer"
	"context"
	"testing"

	documentService "Qingyu_backend/service/writer/document"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortcutService_GetUserShortcuts(t *testing.T) {
	service := documentService.NewShortcutService()
	ctx := context.Background()

	t.Run("成功获取默认快捷键", func(t *testing.T) {
		config, err := service.GetUserShortcuts(ctx, "user123")

		require.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, "user123", config.UserID)
		assert.NotEmpty(t, config.Shortcuts)
		assert.NotZero(t, config.CreatedAt)
		assert.NotZero(t, config.UpdatedAt)
	})

	t.Run("验证默认快捷键数量", func(t *testing.T) {
		config, err := service.GetUserShortcuts(ctx, "user123")

		require.NoError(t, err)
		// 应该有34个默认快捷键
		assert.Equal(t, 34, len(config.Shortcuts))
	})

	t.Run("验证常用快捷键存在", func(t *testing.T) {
		config, err := service.GetUserShortcuts(ctx, "user123")

		require.NoError(t, err)

		// 验证重要的快捷键
		commonShortcuts := []string{
			"save",   // Ctrl+S
			"undo",   // Ctrl+Z
			"redo",   // Ctrl+Y
			"copy",   // Ctrl+C
			"paste",  // Ctrl+V
			"bold",   // Ctrl+B
			"italic", // Ctrl+I
			"find",   // Ctrl+F
		}

		for _, key := range commonShortcuts {
			shortcut, exists := config.Shortcuts[key]
			assert.True(t, exists, "快捷键 %s 应该存在", key)
			assert.NotEmpty(t, shortcut.Key, "快捷键 %s 的按键不应为空", key)
			assert.NotEmpty(t, shortcut.Description, "快捷键 %s 的描述不应为空", key)
			assert.NotEmpty(t, shortcut.Category, "快捷键 %s 的分类不应为空", key)
		}
	})

	t.Run("空用户ID返回错误", func(t *testing.T) {
		config, err := service.GetUserShortcuts(ctx, "")

		assert.Error(t, err)
		assert.Nil(t, config)
		assert.Contains(t, err.Error(), "用户ID不能为空")
	})
}

func TestShortcutService_UpdateUserShortcuts(t *testing.T) {
	service := documentService.NewShortcutService()
	ctx := context.Background()

	t.Run("成功更新快捷键", func(t *testing.T) {
		shortcuts := map[string]writer.Shortcut{
			"save": {
				Action:      "save",
				Key:         "Ctrl+Shift+S",
				Description: "自定义保存",
				Category:    "文件",
				IsCustom:    true,
			},
		}

		err := service.UpdateUserShortcuts(ctx, "user123", shortcuts)

		assert.NoError(t, err)
	})

	t.Run("空用户ID返回错误", func(t *testing.T) {
		shortcuts := map[string]writer.Shortcut{
			"save": {Action: "save", Key: "Ctrl+S"},
		}

		err := service.UpdateUserShortcuts(ctx, "", shortcuts)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "用户ID不能为空")
	})

	t.Run("空快捷键配置返回错误", func(t *testing.T) {
		err := service.UpdateUserShortcuts(ctx, "user123", map[string]writer.Shortcut{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "快捷键配置不能为空")
	})

	t.Run("快捷键冲突检测", func(t *testing.T) {
		shortcuts := map[string]writer.Shortcut{
			"save": {
				Action: "save",
				Key:    "Ctrl+S",
			},
			"custom_save": {
				Action: "custom_save",
				Key:    "Ctrl+S", // 与save冲突
			},
		}

		err := service.UpdateUserShortcuts(ctx, "user123", shortcuts)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "快捷键冲突")
	})

	t.Run("空按键返回错误", func(t *testing.T) {
		shortcuts := map[string]writer.Shortcut{
			"save": {
				Action: "save",
				Key:    "", // 空按键
			},
		}

		err := service.UpdateUserShortcuts(ctx, "user123", shortcuts)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "快捷键不能为空")
	})
}

func TestShortcutService_ResetUserShortcuts(t *testing.T) {
	service := documentService.NewShortcutService()
	ctx := context.Background()

	t.Run("成功重置快捷键", func(t *testing.T) {
		err := service.ResetUserShortcuts(ctx, "user123")

		assert.NoError(t, err)
	})

	t.Run("空用户ID返回错误", func(t *testing.T) {
		err := service.ResetUserShortcuts(ctx, "")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "用户ID不能为空")
	})
}

func TestShortcutService_GetShortcutHelp(t *testing.T) {
	service := documentService.NewShortcutService()
	ctx := context.Background()

	t.Run("成功获取快捷键帮助", func(t *testing.T) {
		categories, err := service.GetShortcutHelp(ctx, "user123")

		require.NoError(t, err)
		assert.NotEmpty(t, categories)
	})

	t.Run("验证分类结构", func(t *testing.T) {
		categories, err := service.GetShortcutHelp(ctx, "user123")

		require.NoError(t, err)

		// 验证每个分类
		expectedCategories := []string{"文件", "编辑", "格式", "段落", "插入", "视图"}
		foundCategories := make(map[string]bool)

		for _, category := range categories {
			foundCategories[category.Name] = true
			assert.NotEmpty(t, category.Shortcuts, "分类 %s 应该有快捷键", category.Name)

			// 验证每个快捷键的结构
			for _, shortcut := range category.Shortcuts {
				assert.NotEmpty(t, shortcut.Action, "快捷键动作不应为空")
				assert.NotEmpty(t, shortcut.Key, "按键不应为空")
				assert.NotEmpty(t, shortcut.Description, "描述不应为空")
				assert.Equal(t, category.Name, shortcut.Category, "分类应该匹配")
			}
		}

		// 验证所有预期分类都存在
		for _, expected := range expectedCategories {
			assert.True(t, foundCategories[expected], "应该包含分类 %s", expected)
		}
	})

	t.Run("空用户ID返回错误", func(t *testing.T) {
		categories, err := service.GetShortcutHelp(ctx, "")

		assert.Error(t, err)
		assert.Nil(t, categories)
	})
}

func TestShortcutService_DefaultShortcuts(t *testing.T) {
	t.Run("验证所有默认快捷键", func(t *testing.T) {
		shortcuts := writer.GetDefaultShortcuts()

		// 按分类统计
		categoryCount := make(map[string]int)
		for _, shortcut := range shortcuts {
			categoryCount[shortcut.Category]++
		}

		// 验证分类数量 (总计34个快捷键)
		assert.Equal(t, 4, categoryCount["文件"], "文件分类应该有4个快捷键")
		assert.Equal(t, 8, categoryCount["编辑"], "编辑分类应该有8个快捷键")
		assert.Equal(t, 7, categoryCount["格式"], "格式分类应该有7个快捷键")
		assert.Equal(t, 5, categoryCount["段落"], "段落分类应该有5个快捷键")
		assert.Equal(t, 4, categoryCount["插入"], "插入分类应该有4个快捷键")
		assert.Equal(t, 6, categoryCount["视图"], "视图分类应该有6个快捷键")
	})

	t.Run("验证默认快捷键不可变", func(t *testing.T) {
		// 获取两次默认配置
		shortcuts1 := writer.GetDefaultShortcuts()
		shortcuts2 := writer.GetDefaultShortcuts()

		// 修改第一个
		shortcuts1["save"] = writer.Shortcut{
			Action: "save",
			Key:    "Modified",
		}

		// 验证第二个没有被修改
		assert.NotEqual(t, "Modified", shortcuts2["save"].Key,
			"默认快捷键应该是不可变的")
	})
}

func TestShortcutService_ValidationLogic(t *testing.T) {
	service := documentService.NewShortcutService()
	ctx := context.Background()

	t.Run("允许多个动作有不同按键", func(t *testing.T) {
		shortcuts := map[string]writer.Shortcut{
			"save": {
				Action: "save",
				Key:    "Ctrl+S",
			},
			"save_all": {
				Action: "save_all",
				Key:    "Ctrl+Shift+S",
			},
		}

		err := service.UpdateUserShortcuts(ctx, "user123", shortcuts)

		assert.NoError(t, err)
	})

	t.Run("检测按键完全相同的冲突", func(t *testing.T) {
		shortcuts := map[string]writer.Shortcut{
			"action1": {
				Action: "action1",
				Key:    "Ctrl+K",
			},
			"action2": {
				Action: "action2",
				Key:    "Ctrl+K", // 完全相同
			},
		}

		err := service.UpdateUserShortcuts(ctx, "user123", shortcuts)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "快捷键冲突")
		assert.Contains(t, err.Error(), "Ctrl+K")
	})

	t.Run("允许单个快捷键", func(t *testing.T) {
		shortcuts := map[string]writer.Shortcut{
			"only_one": {
				Action:      "only_one",
				Key:         "F1",
				Description: "唯一快捷键",
				Category:    "其他",
			},
		}

		err := service.UpdateUserShortcuts(ctx, "user123", shortcuts)

		assert.NoError(t, err)
	})
}

func BenchmarkShortcutService_GetUserShortcuts(b *testing.B) {
	service := documentService.NewShortcutService()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GetUserShortcuts(ctx, "user123")
	}
}

func BenchmarkShortcutService_ValidateShortcuts(b *testing.B) {
	service := documentService.NewShortcutService()
	ctx := context.Background()

	shortcuts := writer.GetDefaultShortcuts()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.UpdateUserShortcuts(ctx, "user123", shortcuts)
	}
}
