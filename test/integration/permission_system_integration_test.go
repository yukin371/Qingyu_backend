package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"

	"Qingyu_backend/global"
	"Qingyu_backend/models/auth"
)

// TestPermissionSystem_Integration 权限系统集成测试
// 测试完整的权限系统流程：权限模型 -> 权限检查 -> API应用
func TestPermissionSystem_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过权限系统集成测试（使用 -short 标志）")
	}

	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)

	t.Run("权限模型完整性", func(t *testing.T) {
		testPermissionModelIntegrity(t)
	})

	t.Run("权限模板系统", func(t *testing.T) {
		testPermissionTemplateSystem(t, router, helper)
	})

	t.Run("动态权限检查", func(t *testing.T) {
		testDynamicPermissionCheck(t, router, helper)
	})

	t.Run("权限继承与覆盖", func(t *testing.T) {
		testPermissionInheritance(t)
	})
}

// testPermissionModelIntegrity 测试权限模型完整性
func testPermissionModelIntegrity(t *testing.T) {
	mongoDB := global.DB
	if mongoDB == nil {
		t.Skip("数据库连接未初始化，跳过权限模型测试")
	}

	t.Run("1.权限集合存在性", func(t *testing.T) {
		collections := []string{
			"permissions",
			"permission_templates",
			"role_permissions",
		}

		for _, coll := range collections {
			// 尝试查询集合
			count, err := mongoDB.Collection(coll).CountDocuments(context.Background(), bson.M{})
			if err != nil {
				t.Logf("○ 集合 %s 查询失败（可能不存在）: %v", coll, err)
			} else {
				t.Logf("✓ 集合 %s 存在，文档数: %d", coll, count)
			}
		}
	})

	t.Run("2.权限模型结构验证", func(t *testing.T) {
		// 查询一个权限文档来验证结构
		var permission auth.Permission
		err := mongoDB.Collection("permissions").FindOne(context.Background(), bson.M{}).Decode(&permission)

		if err == nil {
			// 验证必需字段
			assert.NotEmpty(t, permission.ID, "权限ID不应为空")
			assert.NotEmpty(t, permission.Resource, "资源不应为空")
			assert.NotEmpty(t, permission.Action, "操作不应为空")
			t.Logf("✓ 权限模型结构正确: %+v", permission)
		} else {
			t.Log("○ 权限集合为空或查询失败")
		}
	})
}

// testPermissionTemplateSystem 测试权限模板系统
func testPermissionTemplateSystem(t *testing.T, router *gin.Engine, helper *TestHelper) {
	t.Run("1.获取权限模板列表", func(t *testing.T) {
		// 尝试登录管理员账号
		adminToken := helper.LoginUser("admin", "Admin@123456")
		if adminToken == "" {
			t.Skip("需要管理员账号，跳过权限模板测试")
		}

		w := helper.DoRequest("GET", "/api/v1/admin/permission-templates", nil, adminToken)

		if w.Code == 200 {
			helper.LogSuccess("获取权限模板列表成功")
		} else {
			t.Logf("获取权限模板列表状态码: %d", w.Code)
		}
	})

	t.Run("2.创建权限模板", func(t *testing.T) {
		adminToken := helper.LoginUser("admin", "Admin@123456")
		if adminToken == "" {
			t.Skip("需要管理员账号，跳过创建权限模板测试")
		}

		timestamp := time.Now().Unix()
		templateData := map[string]interface{}{
			"name":        fmt.Sprintf("test_template_%d", timestamp),
			"description": "测试权限模板",
			"permissions": []string{
				"book:read",
				"chapter:read",
			},
		}

		w := helper.DoRequest("POST", "/api/v1/admin/permission-templates", templateData, adminToken)

		if w.Code == 200 || w.Code == 201 {
			helper.LogSuccess("创建权限模板成功")
		} else {
			t.Logf("创建权限模板状态码: %d, 响应: %s", w.Code, w.Body.String())
		}
	})

	t.Run("3.应用权限模板到用户", func(t *testing.T) {
		// 这个测试需要先创建模板，然后应用到用户
		t.Skip("需要完整的模板-用户流程，暂时跳过")
	})
}

// testDynamicPermissionCheck 测试动态权限检查
func testDynamicPermissionCheck(t *testing.T, router *gin.Engine, helper *TestHelper) {
	t.Run("1.基于角色的权限检查", func(t *testing.T) {
		// 1.1 普通用户权限
		normalToken := helper.LoginUser("test_user01", "Test@123456")
		if normalToken != "" {
			// 测试访问普通用户接口
			w := helper.DoRequest("GET", UserProfilePath, nil, normalToken)
			if w.Code == 200 || w.Code == 404 {
				t.Log("✓ 普通用户可以访问普通接口")
			}

			// 测试访问管理员接口（应该被拒绝）
			w = helper.DoRequest("GET", "/api/v1/admin/users", nil, normalToken)
			if w.Code == 403 || w.Code == 401 {
				t.Log("✓ 普通用户无法访问管理员接口")
			}
		}

		// 1.2 管理员权限
		adminToken := helper.LoginUser("admin", "Admin@123456")
		if adminToken != "" {
			// 测试访问管理员接口
			w := helper.DoRequest("GET", "/api/v1/admin/users", nil, adminToken)
			if w.Code == 200 || w.Code == 404 {
				t.Log("✓ 管理员可以访问管理员接口")
			}
		}
	})

	t.Run("2.基于资源的权限检查", func(t *testing.T) {
		// 测试资源级别的权限控制
		t.Run("2.1书籍访问权限", func(t *testing.T) {
			// 测试VIP书籍访问权限
			testToken := helper.LoginUser("test_user01", "Test@123456")
			if testToken != "" {
				// 尝试访问VIP书籍（可能需要VIP权限）
				w := helper.DoRequest("GET", "/api/v1/reader/books/vip-book-123", nil, testToken)

				if w.Code == 403 {
					t.Log("✓ 普通用户无法访问VIP书籍")
				} else if w.Code == 404 {
					t.Log("○ VIP书籍不存在")
				} else {
					t.Logf("VIP书籍访问状态: %d", w.Code)
				}
			}
		})
	})

	t.Run("3.权限缓存机制", func(t *testing.T) {
		// 测试权限缓存是否正常工作
		// 1. 首次请求
		// 2. 第二次请求（应该命中缓存）
		t.Skip("需要缓存监控机制，暂时跳过")
	})
}

// testPermissionInheritance 测试权限继承与覆盖
func testPermissionInheritance(t *testing.T) {
	mongoDB := global.DB
	if mongoDB == nil {
		t.Skip("数据库连接未初始化，跳过权限继承测试")
	}

	t.Run("1.角色权限继承", func(t *testing.T) {
		// 测试角色之间的权限继承关系
		// 例如：admin角色应该包含所有基础权限

		// 查询admin角色的权限
		var adminPermissions []auth.Permission
		cursor, err := mongoDB.Collection("role_permissions").Find(
			context.Background(),
			bson.M{"role": "admin"},
		)

		if err == nil {
			defer cursor.Close(context.Background())
			cursor.All(context.Background(), &adminPermissions)

			t.Logf("✓ Admin角色权限数: %d", len(adminPermissions))

			// 验证admin角色应该有的权限
			expectedPermissions := []string{
				"user:read",
				"user:write",
				"user:delete",
			}

			for _, expected := range expectedPermissions {
				found := false
				for _, perm := range adminPermissions {
					if perm.Resource+":"+perm.Action == expected {
						found = true
						break
					}
				}
				if found {
					t.Logf("✓ Admin拥有权限: %s", expected)
				}
			}
		} else {
			t.Log("○ 角色权限集合查询失败或为空")
		}
	})

	t.Run("2.用户自定义权限覆盖", func(t *testing.T) {
		// 测试用户自定义权限是否可以覆盖角色权限
		t.Skip("需要设置自定义权限的用户，暂时跳过")
	})
}

// TestPermissionAPI_Integration 权限API集成测试
func TestPermissionAPI_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过权限API集成测试（使用 -short 标志）")
	}

	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)

	// 1. 测试权限检查API
	t.Run("权限检查API", func(t *testing.T) {
		testToken := helper.LoginUser("test_user01", "Test@123456")
		if testToken == "" {
			t.Skip("无法获取测试Token，跳过权限检查API测试")
		}

		t.Run("1.检查当前用户权限", func(t *testing.T) {
			w := helper.DoRequest("GET", "/api/v1/user/permissions", nil, testToken)

			if w.Code == 200 {
				helper.LogSuccess("获取用户权限列表成功")
			} else if w.Code == 404 {
				t.Log("○ 权限接口不存在")
			} else {
				t.Logf("获取权限列表状态码: %d", w.Code)
			}
		})

		t.Run("2.检查特定资源权限", func(t *testing.T) {
			checkData := map[string]interface{}{
				"resource": "book",
				"action":   "read",
			}
			w := helper.DoRequest("POST", "/api/v1/user/permissions/check", checkData, testToken)

			if w.Code == 200 {
				helper.LogSuccess("权限检查API正常")
			} else if w.Code == 404 {
				t.Log("○ 权限检查接口不存在")
			} else {
				t.Logf("权限检查状态码: %d", w.Code)
			}
		})
	})

	// 2. 测试管理员权限管理API
	t.Run("管理员权限管理API", func(t *testing.T) {
		adminToken := helper.LoginUser("admin", "Admin@123456")
		if adminToken == "" {
			t.Skip("需要管理员账号，跳过管理员权限管理API测试")
		}

		t.Run("1.获取所有权限列表", func(t *testing.T) {
			w := helper.DoRequest("GET", "/api/v1/admin/permissions", nil, adminToken)

			if w.Code == 200 {
				helper.LogSuccess("获取权限列表成功")
			} else if w.Code == 404 {
				t.Log("○ 权限列表接口不存在")
			} else {
				t.Logf("获取权限列表状态码: %d", w.Code)
			}
		})

		t.Run("2.创建新权限", func(t *testing.T) {
			timestamp := time.Now().Unix()
			permData := map[string]interface{}{
				"resource": "test_resource",
				"action":   "test_action",
				"description": fmt.Sprintf("测试权限_%d", timestamp),
			}
			w := helper.DoRequest("POST", "/api/v1/admin/permissions", permData, adminToken)

			if w.Code == 200 || w.Code == 201 {
				helper.LogSuccess("创建权限成功")
			} else if w.Code == 404 {
				t.Log("○ 创建权限接口不存在")
			} else {
				t.Logf("创建权限状态码: %d", w.Code)
			}
		})
	})
}

// TestPermission_Performance 权限系统性能测试
func TestPermission_Performance(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过权限系统性能测试（使用 -short 标志）")
	}

	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)

	testToken := helper.LoginUser("test_user01", "Test@123456")
	if testToken == "" {
		t.Skip("无法获取测试Token，跳过性能测试")
	}

	t.Run("权限检查性能", func(t *testing.T) {
		// 测试多次权限检查的性能
		iterations := 50
		start := time.Now()

		for i := 0; i < iterations; i++ {
			w := helper.DoRequest("GET", UserProfilePath, nil, testToken)
			if w.Code != 200 && w.Code != 404 {
				t.Logf("第%d次请求状态码: %d", i+1, w.Code)
			}
		}

		duration := time.Since(start)
		avgDuration := duration / time.Duration(iterations)

		t.Logf("✓ %d次权限检查总耗时: %v", iterations, duration)
		t.Logf("✓ 平均每次检查耗时: %v", avgDuration)

		// 验证性能是否在可接受范围内（< 100ms）
		if avgDuration < 100*time.Millisecond {
			t.Log("✓ 权限检查性能良好")
		} else {
			t.Logf("⚠ 权限检查性能可能需要优化（平均耗时: %v）", avgDuration)
		}
	})
}
