package document

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"Qingyu_backend/models/writer"
	documentSvc "Qingyu_backend/service/document"
)

// TestShortcutService_GetUserShortcuts 测试获取用户快捷键配置
func TestShortcutService_GetUserShortcuts(t *testing.T) {
	service := documentSvc.NewShortcutService()
	ctx := context.Background()

	t.Run("ValidUserID", func(t *testing.T) {
		userID := "user123"
		config, err := service.GetUserShortcuts(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, userID, config.UserID)
		assert.NotNil(t, config.Shortcuts)
		assert.Greater(t, len(config.Shortcuts), 0)
		t.Logf("✓ 获取快捷键成功: 快捷键数量=%d", len(config.Shortcuts))
	})

	t.Run("EmptyUserID", func(t *testing.T) {
		config, err := service.GetUserShortcuts(ctx, "")

		assert.Error(t, err)
		assert.Nil(t, config)
		t.Logf("✓ 空UserID验证通过")
	})

	t.Run("DefaultShortcutsReturned", func(t *testing.T) {
		userID := "newuser"
		config, err := service.GetUserShortcuts(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, config)
		// 应该返回默认快捷键
		assert.Greater(t, len(config.Shortcuts), 0)
		t.Logf("✓ 返回默认快捷键")
	})

	t.Run("DifferentUsers", func(t *testing.T) {
		// 不同用户应该获得相同的默认配置
		config1, err1 := service.GetUserShortcuts(ctx, "user1")
		config2, err2 := service.GetUserShortcuts(ctx, "user2")

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, len(config1.Shortcuts), len(config2.Shortcuts))
		t.Logf("✓ 不同用户获得相同的默认配置")
	})
}

// TestShortcutService_UpdateUserShortcuts 测试更新用户快捷键
func TestShortcutService_UpdateUserShortcuts(t *testing.T) {
	service := documentSvc.NewShortcutService()
	ctx := context.Background()

	t.Run("ValidUpdate", func(t *testing.T) {
		userID := "user123"
		shortcuts := map[string]writer.Shortcut{
			"save": {
				Action:      "save",
				Key:         "Ctrl+S",
				Description: "保存文档",
				Category:    "editing",
			},
		}

		err := service.UpdateUserShortcuts(ctx, userID, shortcuts)

		assert.NoError(t, err)
		t.Logf("✓ 快捷键更新成功")
	})

	t.Run("EmptyUserID", func(t *testing.T) {
		shortcuts := map[string]writer.Shortcut{
			"save": {
				Action:      "save",
				Key:         "Ctrl+S",
				Description: "保存文档",
				Category:    "editing",
			},
		}

		err := service.UpdateUserShortcuts(ctx, "", shortcuts)

		assert.Error(t, err)
		t.Logf("✓ 空UserID验证通过")
	})

	t.Run("EmptyShortcuts", func(t *testing.T) {
		userID := "user123"
		shortcuts := make(map[string]writer.Shortcut)

		err := service.UpdateUserShortcuts(ctx, userID, shortcuts)

		assert.Error(t, err)
		t.Logf("✓ 空快捷键配置验证通过")
	})

	t.Run("MultipleShortcuts", func(t *testing.T) {
		userID := "user456"
		shortcuts := map[string]writer.Shortcut{
			"save": {
				Action:      "save",
				Key:         "Ctrl+S",
				Description: "保存文档",
				Category:    "editing",
			},
			"undo": {
				Action:      "undo",
				Key:         "Ctrl+Z",
				Description: "撤销操作",
				Category:    "editing",
			},
			"redo": {
				Action:      "redo",
				Key:         "Ctrl+Y",
				Description: "重做操作",
				Category:    "editing",
			},
		}

		err := service.UpdateUserShortcuts(ctx, userID, shortcuts)

		assert.NoError(t, err)
		t.Logf("✓ 多个快捷键更新成功")
	})

	t.Run("InvalidShortcuts", func(t *testing.T) {
		userID := "user789"
		// 创建包含无效数据的快捷键
		shortcuts := map[string]writer.Shortcut{
			"": { // 空的快捷键名
				Action:      "",
				Key:         "",
				Description: "",
				Category:    "",
			},
		}

		err := service.UpdateUserShortcuts(ctx, userID, shortcuts)

		assert.Error(t, err)
		t.Logf("✓ 无效快捷键验证通过")
	})
}

// TestShortcutService_ResetUserShortcuts 测试重置用户快捷键
func TestShortcutService_ResetUserShortcuts(t *testing.T) {
	service := documentSvc.NewShortcutService()
	ctx := context.Background()

	t.Run("ValidReset", func(t *testing.T) {
		userID := "user123"
		err := service.ResetUserShortcuts(ctx, userID)

		assert.NoError(t, err)
		t.Logf("✓ 快捷键重置成功")
	})

	t.Run("EmptyUserID", func(t *testing.T) {
		err := service.ResetUserShortcuts(ctx, "")

		assert.Error(t, err)
		t.Logf("✓ 空UserID验证通过")
	})

	t.Run("MultipleResets", func(t *testing.T) {
		userID := "user456"
		// 多次重置应该都成功
		err1 := service.ResetUserShortcuts(ctx, userID)
		err2 := service.ResetUserShortcuts(ctx, userID)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		t.Logf("✓ 多次重置都成功")
	})
}

// TestShortcutService_GetShortcutHelp 测试获取快捷键帮助
func TestShortcutService_GetShortcutHelp(t *testing.T) {
	service := documentSvc.NewShortcutService()
	ctx := context.Background()

	t.Run("ValidUserHelp", func(t *testing.T) {
		userID := "user123"
		categories, err := service.GetShortcutHelp(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, categories)
		// 应该返回按分类组织的快捷键
		t.Logf("✓ 获取快捷键帮助成功: 分类数=%d", len(categories))
	})

	t.Run("EmptyUserID", func(t *testing.T) {
		categories, err := service.GetShortcutHelp(ctx, "")

		assert.Error(t, err)
		assert.Nil(t, categories)
		t.Logf("✓ 空UserID验证通过")
	})

	t.Run("CategoriesOrganized", func(t *testing.T) {
		userID := "user789"
		categories, err := service.GetShortcutHelp(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, categories)
		// 验证分类都是有效的
		for _, category := range categories {
			assert.NotEmpty(t, category.Name)
			assert.NotNil(t, category.Shortcuts)
			assert.Greater(t, len(category.Shortcuts), 0)
		}
		t.Logf("✓ 快捷键分类组织正确")
	})

	t.Run("ConsistentHelp", func(t *testing.T) {
		// 相同用户应该获得相同的帮助信息
		categories1, err1 := service.GetShortcutHelp(ctx, "user1")
		categories2, err2 := service.GetShortcutHelp(ctx, "user1")

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, len(categories1), len(categories2))
		t.Logf("✓ 帮助信息一致")
	})
}

// TestShortcutService_ConcurrentOperations 并发操作测试
func TestShortcutService_ConcurrentOperations(t *testing.T) {
	service := documentSvc.NewShortcutService()
	ctx := context.Background()

	t.Run("ConcurrentReads", func(t *testing.T) {
		done := make(chan error, 10)

		// 10个并发读取
		for i := 0; i < 10; i++ {
			go func(userID string) {
				_, err := service.GetUserShortcuts(ctx, userID)
				done <- err
			}(string(rune('0' + i)))
		}

		for i := 0; i < 10; i++ {
			err := <-done
			assert.NoError(t, err)
		}
		t.Logf("✓ 10个并发读取成功")
	})

	t.Run("ConcurrentUpdates", func(t *testing.T) {
		done := make(chan error, 5)

		for i := 0; i < 5; i++ {
			go func(userID string) {
				shortcuts := map[string]writer.Shortcut{
					"save": {
						Action:      "save",
						Key:         "Ctrl+S",
						Description: "保存文档",
						Category:    "editing",
					},
				}
				err := service.UpdateUserShortcuts(ctx, userID, shortcuts)
				done <- err
			}(string(rune('A' + rune(i))))
		}

		for i := 0; i < 5; i++ {
			err := <-done
			assert.NoError(t, err)
		}
		t.Logf("✓ 5个并发更新成功")
	})
}

// TestShortcutService_EdgeCases 边界条件测试
func TestShortcutService_EdgeCases(t *testing.T) {
	service := documentSvc.NewShortcutService()
	ctx := context.Background()

	t.Run("VeryLongUserID", func(t *testing.T) {
		userID := ""
		for i := 0; i < 1000; i++ {
			userID += "a"
		}
		config, err := service.GetUserShortcuts(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, config)
		t.Logf("✓ 很长的UserID处理成功")
	})

	t.Run("SpecialCharactersUserID", func(t *testing.T) {
		userID := "user@#$%^&*()"
		config, err := service.GetUserShortcuts(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, config)
		t.Logf("✓ 特殊字符UserID处理成功")
	})

	t.Run("UnicodeUserID", func(t *testing.T) {
		userID := "用户123🎉"
		config, err := service.GetUserShortcuts(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, config)
		t.Logf("✓ Unicode UserID处理成功")
	})

	t.Run("SpaceUserID", func(t *testing.T) {
		userID := "   " // 只有空格
		config, err := service.GetUserShortcuts(ctx, userID)

		// 取决于实现：可能通过也可能失败（取决于trim逻辑）
		if err != nil {
			t.Logf("✓ 空格UserID被拒绝（符合预期）")
		} else {
			assert.NotNil(t, config)
			t.Logf("✓ 空格UserID被接受")
		}
	})
}
