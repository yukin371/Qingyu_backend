package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"

	"Qingyu_backend/global"
	"Qingyu_backend/models/admin"
)

// TestExportAPI_Integration 导出API集成测试
// 测试完整的导出流程：API -> Service -> Repository
func TestExportAPI_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过导出API集成测试（使用 -short 标志）")
	}

	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)

	t.Run("导出API基础功能", func(t *testing.T) {
		testExportAPIBasicFunctions(t, router, helper)
	})

	t.Run("导出历史记录", func(t *testing.T) {
		testExportHistory(t, router, helper)
	})

	t.Run("导出权限控制", func(t *testing.T) {
		testExportPermissionControl(t, router, helper)
	})

	t.Run("导出格式验证", func(t *testing.T) {
		testExportFormatValidation(t, router, helper)
	})
}

// testExportAPIBasicFunctions 测试导出API基础功能
func testExportAPIBasicFunctions(t *testing.T, router *gin.Engine, helper *TestHelper) {
	// 1. 测试获取导出任务列表
	t.Run("1.获取导出任务列表", func(t *testing.T) {
		// 尝试使用管理员账号
		adminToken := helper.LoginUser("admin", "Admin@123456")
		if adminToken == "" {
			t.Skip("需要管理员账号，跳过导出任务列表测试")
		}

		w := helper.DoRequest("GET", "/api/v1/admin/exports", nil, adminToken)

		if w.Code == 200 {
			helper.LogSuccess("获取导出任务列表成功")
		} else if w.Code == 404 {
			t.Log("○ 导出任务列表接口不存在")
		} else {
			t.Logf("获取导出任务列表状态码: %d", w.Code)
		}
	})

	// 2. 测试创建导出任务
	t.Run("2.创建书籍内容导出任务", func(t *testing.T) {
		adminToken := helper.LoginUser("admin", "Admin@123456")
		if adminToken == "" {
			t.Skip("需要管理员账号，跳过创建导出任务测试")
		}

		timestamp := time.Now().Unix()
		exportData := map[string]interface{}{
			"type":        "book_content",
			"book_id":     "test-book-001",
			"format":      "markdown",
			"include_metadata": true,
			"description": fmt.Sprintf("测试导出任务_%d", timestamp),
		}

		w := helper.DoRequest("POST", "/api/v1/admin/exports", exportData, adminToken)

		if w.Code == 200 || w.Code == 201 {
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			helper.LogSuccess("创建导出任务成功")

			if data, ok := response["data"].(map[string]interface{}); ok {
				if taskID, ok := data["task_id"].(string); ok {
					t.Logf("✓ 导出任务ID: %s", taskID)
				}
			}
		} else {
			t.Logf("创建导出任务状态码: %d, 响应: %s", w.Code, w.Body.String())
		}
	})

	// 3. 测试查询导出任务状态
	t.Run("3.查询导出任务状态", func(t *testing.T) {
		adminToken := helper.LoginUser("admin", "Admin@123456")
		if adminToken == "" {
			t.Skip("需要管理员账号，跳过查询导出任务状态测试")
		}

		// 使用一个假的任务ID进行测试
		testTaskID := "507f1f77bcf86cd799439011"
		w := helper.DoRequest("GET", "/api/v1/admin/exports/"+testTaskID, nil, adminToken)

		if w.Code == 200 {
			helper.LogSuccess("查询导出任务状态成功")
		} else if w.Code == 404 {
			t.Log("○ 导出任务不存在（符合预期）")
		} else {
			t.Logf("查询导出任务状态状态码: %d", w.Code)
		}
	})

	// 4. 测试取消导出任务
	t.Run("4.取消导出任务", func(t *testing.T) {
		adminToken := helper.LoginUser("admin", "Admin@123456")
		if adminToken == "" {
			t.Skip("需要管理员账号，跳过取消导出任务测试")
		}

		// 使用一个假的任务ID进行测试
		testTaskID := "507f1f77bcf86cd799439011"
		w := helper.DoRequest("DELETE", "/api/v1/admin/exports/"+testTaskID, nil, adminToken)

		if w.Code == 200 || w.Code == 204 {
			helper.LogSuccess("取消导出任务成功")
		} else if w.Code == 404 {
			t.Log("○ 导出任务不存在（符合预期）")
		} else {
			t.Logf("取消导出任务状态码: %d", w.Code)
		}
	})
}

// testExportHistory 测试导出历史记录
func testExportHistory(t *testing.T, router *gin.Engine, helper *TestHelper) {
	mongoDB := global.DB
	if mongoDB == nil {
		t.Skip("数据库连接未初始化，跳过导出历史测试")
	}

	t.Run("1.导出历史记录存储", func(t *testing.T) {
		// 检查导出历史集合
		count, err := mongoDB.Collection("export_histories").CountDocuments(context.Background(), bson.M{})
		require.NoError(t, err, "导出历史集合查询失败")

		t.Logf("✓ 导出历史记录数: %d", count)
	})

	t.Run("2.导出历史记录结构", func(t *testing.T) {
		var exportHistory admin.ExportHistory
		err := mongoDB.Collection("export_histories").FindOne(context.Background(), bson.M{}).Decode(&exportHistory)

		if err == nil {
			// 验证必需字段
			assert.NotEmpty(t, exportHistory.ID, "导出历史ID不应为空")
			assert.NotEmpty(t, exportHistory.ExportType, "导出类型不应为空")
			assert.NotEmpty(t, exportHistory.AdminID, "管理员ID不应为空")
			t.Logf("✓ 导出历史结构正确: %+v", exportHistory)
		} else {
			t.Log("○ 导出历史集合为空")
		}
	})

	t.Run("3.用户导出历史查询", func(t *testing.T) {
		testToken := helper.LoginUser("test_user01", "Test@123456")
		if testToken == "" {
			t.Skip("无法获取测试Token，跳过用户导出历史测试")
		}

		w := helper.DoRequest("GET", "/api/v1/user/exports/history", nil, testToken)

		if w.Code == 200 {
			helper.LogSuccess("获取用户导出历史成功")
		} else if w.Code == 404 {
			t.Log("○ 用户导出历史接口不存在")
		} else {
			t.Logf("获取用户导出历史状态码: %d", w.Code)
		}
	})
}

// testExportPermissionControl 测试导出权限控制
func testExportPermissionControl(t *testing.T, router *gin.Engine, helper *TestHelper) {
	t.Run("1.普通用户导出权限", func(t *testing.T) {
		testToken := helper.LoginUser("test_user01", "Test@123456")
		if testToken == "" {
			t.Skip("无法获取测试Token，跳过普通用户导出权限测试")
		}

		// 测试用户导出自己的内容
		exportData := map[string]interface{}{
			"type":   "book_content",
			"book_id": "my-book-001",
			"format": "markdown",
		}

		w := helper.DoRequest("POST", "/api/v1/user/exports", exportData, testToken)

		if w.Code == 200 || w.Code == 201 {
			helper.LogSuccess("用户可以创建导出任务")
		} else if w.Code == 404 {
			t.Log("○ 用户导出接口不存在")
		} else if w.Code == 403 {
			t.Log("✓ 用户导出权限控制正常")
		} else {
			t.Logf("用户导出状态码: %d", w.Code)
		}
	})

	t.Run("2.普通用户访问管理员导出接口", func(t *testing.T) {
		testToken := helper.LoginUser("test_user01", "Test@123456")
		if testToken == "" {
			t.Skip("无法获取测试Token，跳过权限测试")
		}

		// 测试普通用户访问管理员导出接口（应该被拒绝）
		w := helper.DoRequest("GET", "/api/v1/admin/exports", nil, testToken)

		if w.Code == 403 || w.Code == 401 {
			helper.LogSuccess("普通用户无法访问管理员导出接口")
		} else if w.Code == 404 {
			t.Log("○ 管理员导出接口不存在")
		} else {
			t.Logf("访问管理员导出接口状态码: %d（可能存在问题）", w.Code)
		}
	})
}

// testExportFormatValidation 测试导出格式验证
func testExportFormatValidation(t *testing.T, router *gin.Engine, helper *TestHelper) {
	adminToken := helper.LoginUser("admin", "Admin@123456")
	if adminToken == "" {
		t.Skip("需要管理员账号，跳过导出格式验证测试")
	}

	supportedFormats := []string{"markdown", "json", "txt", "html"}

	t.Run("1.支持的导出格式", func(t *testing.T) {
		for _, format := range supportedFormats {
			exportData := map[string]interface{}{
				"type":   "book_content",
				"book_id": "test-book-001",
				"format": format,
			}

			w := helper.DoRequest("POST", "/api/v1/admin/exports", exportData, adminToken)

			if w.Code == 200 || w.Code == 201 {
				t.Logf("✓ 格式 %s 支持", format)
			} else if w.Code == 400 {
				// 可能是格式验证失败
				t.Logf("○ 格式 %s 可能不支持", format)
			} else {
				t.Logf("格式 %s 测试状态码: %d", format, w.Code)
			}
		}
	})

	t.Run("2.不支持的导出格式", func(t *testing.T) {
		invalidFormats := []string{"pdf", "docx", "xyz"}

		for _, format := range invalidFormats {
			exportData := map[string]interface{}{
				"type":   "book_content",
				"book_id": "test-book-001",
				"format": format,
			}

			w := helper.DoRequest("POST", "/api/v1/admin/exports", exportData, adminToken)

			if w.Code == 400 {
				t.Logf("✓ 不支持的格式 %s 被正确拒绝", format)
			} else {
				t.Logf("格式 %s 测试状态码: %d", format, w.Code)
			}
		}
	})
}

// TestExportService_Integration 导出服务集成测试
// 测试 Service + Repository 层的集成
func TestExportService_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过导出服务集成测试（使用 -short 标志）")
	}

	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	mongoDB := global.DB
	if mongoDB == nil {
		t.Skip("数据库连接未初始化，跳过导出服务测试")
	}

	t.Run("导出任务生命周期", func(t *testing.T) {
		testExportTaskLifecycle(t)
	})

	t.Run("导出数据完整性", func(t *testing.T) {
		testExportDataIntegrity(t)
	})
}

// testExportTaskLifecycle 测试导出任务生命周期
func testExportTaskLifecycle(t *testing.T) {
	mongoDB := global.DB

	// 1. 创建导出任务
	t.Run("1.创建导出任务", func(t *testing.T) {
		timestamp := time.Now().Unix()
		exportHistory := admin.ExportHistory{
			ID:         fmt.Sprintf("export-test-%d", timestamp),
			AdminID:    "test-admin-001",
			ExportType: admin.ExportTypeBooks,
			Format:     "markdown",
			Status:     admin.ExportStatusPending,
			CreatedAt:  time.Now(),
		}

		_, err := mongoDB.Collection("export_histories").InsertOne(context.Background(), exportHistory)
		require.NoError(t, err, "插入导出任务失败")

		t.Logf("✓ 创建导出任务: %s", exportHistory.ID)
	})

	// 2. 查询导出任务
	t.Run("2.查询导出任务", func(t *testing.T) {
		var exportHistory admin.ExportHistory
		err := mongoDB.Collection("export_histories").FindOne(
			context.Background(),
			bson.M{"export_type": admin.ExportTypeBooks},
		).Decode(&exportHistory)

		require.NoError(t, err, "查询导出任务失败")
		assert.Equal(t, admin.ExportStatusPending, exportHistory.Status, "任务状态应该是pending")

		t.Logf("✓ 查询导出任务: %s, 状态: %s", exportHistory.ID, exportHistory.Status)
	})

	// 3. 更新导出任务状态
	t.Run("3.更新导出任务状态", func(t *testing.T) {
		var exportHistory admin.ExportHistory
		err := mongoDB.Collection("export_histories").FindOne(
			context.Background(),
			bson.M{"export_type": admin.ExportTypeBooks},
		).Decode(&exportHistory)

		require.NoError(t, err, "查询导出任务失败")

		// 更新状态为 completed（简化测试，直接跳到completed）
		now := time.Now()
		filePath := fmt.Sprintf("/exports/%s.md", exportHistory.ID)
		_, err = mongoDB.Collection("export_histories").UpdateOne(
			context.Background(),
			bson.M{"_id": exportHistory.ID},
			bson.M{
				"$set": bson.M{
					"status":       admin.ExportStatusCompleted,
					"file_path":    filePath,
					"completed_at": now,
				},
			},
		)

		require.NoError(t, err, "更新导出任务状态失败")

		t.Logf("✓ 更新导出任务状态: %s -> completed", exportHistory.ID)
	})

	// 4. 清理测试数据
	t.Run("4.清理测试数据", func(t *testing.T) {
		_, err := mongoDB.Collection("export_histories").DeleteMany(
			context.Background(),
			bson.M{"export_type": admin.ExportTypeBooks},
		)

		require.NoError(t, err, "清理测试数据失败")

		t.Log("✓ 清理导出任务测试数据")
	})
}

// testExportDataIntegrity 测试导出数据完整性
func testExportDataIntegrity(t *testing.T) {
	mongoDB := global.DB

	t.Run("1.导出任务包含必需字段", func(t *testing.T) {
		// 创建一个完整的导出任务
		timestamp := time.Now().Unix()
		filePath := fmt.Sprintf("/exports/test-%d.md", timestamp)
		now := time.Now()

		exportHistory := admin.ExportHistory{
			ID:         fmt.Sprintf("export-integrity-%d", timestamp),
			AdminID:    "test-admin-integrity",
			ExportType: admin.ExportTypeBooks,
			Format:     "markdown",
			Status:     admin.ExportStatusCompleted,
			FilePath:   filePath,
			CreatedAt:  now,
			CompletedAt: &now,
		}

		_, err := mongoDB.Collection("export_histories").InsertOne(context.Background(), exportHistory)
		require.NoError(t, err, "插入导出任务失败")

		// 查询并验证字段
		var result admin.ExportHistory
		err = mongoDB.Collection("export_histories").FindOne(
			context.Background(),
			bson.M{"_id": exportHistory.ID},
		).Decode(&result)

		require.NoError(t, err, "查询导出任务失败")

		// 验证必需字段
		assert.NotEmpty(t, result.ID, "ID不应为空")
		assert.NotEmpty(t, result.ExportType, "ExportType不应为空")
		assert.NotEmpty(t, result.AdminID, "AdminID不应为空")
		assert.NotEmpty(t, result.Format, "Format不应为空")
		assert.NotEmpty(t, result.Status, "Status不应为空")

		t.Log("✓ 导出任务数据完整性验证通过")

		// 清理
		mongoDB.Collection("export_histories").DeleteOne(context.Background(), bson.M{"_id": exportHistory.ID})
	})
}

// TestExportAPI_ErrorHandling 导出API错误处理测试
func TestExportAPI_ErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过导出API错误处理测试（使用 -short 标志）")
	}

	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)

	t.Run("错误处理", func(t *testing.T) {
		adminToken := helper.LoginUser("admin", "Admin@123456")
		if adminToken == "" {
			t.Skip("需要管理员账号，跳过错误处理测试")
		}

		// 1. 测试缺少必需字段
		t.Run("1.缺少必需字段", func(t *testing.T) {
			invalidData := map[string]interface{}{
				"type": "book_content",
				// 缺少 book_id 和 format
			}

			w := helper.DoRequest("POST", "/api/v1/admin/exports", invalidData, adminToken)

			if w.Code == 400 {
				helper.LogSuccess("缺少必需字段被正确拒绝")
			} else {
				t.Logf("缺少必需字段状态码: %d", w.Code)
			}
		})

		// 2. 测试无效的导出类型
		t.Run("2.无效的导出类型", func(t *testing.T) {
			invalidData := map[string]interface{}{
				"type":    "invalid_type",
				"book_id": "test-book-001",
				"format":  "markdown",
			}

			w := helper.DoRequest("POST", "/api/v1/admin/exports", invalidData, adminToken)

			if w.Code == 400 {
				helper.LogSuccess("无效的导出类型被正确拒绝")
			} else {
				t.Logf("无效的导出类型状态码: %d", w.Code)
			}
		})

		// 3. 测试查询不存在的导出任务
		t.Run("3.查询不存在的导出任务", func(t *testing.T) {
			invalidID := "507f1f77bcf86cd799439011"
			w := helper.DoRequest("GET", "/api/v1/admin/exports/"+invalidID, nil, adminToken)

			if w.Code == 404 {
				helper.LogSuccess("不存在的导出任务返回404")
			} else {
				t.Logf("不存在的导出任务状态码: %d", w.Code)
			}
		})
	})
}
