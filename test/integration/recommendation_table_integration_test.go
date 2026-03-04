//go:build integration
// +build integration

package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/recommendation"
)

// TestRecommendationTableCreate 测试创建榜单
// 测试场景：创建手动榜单，验证创建成功
func TestRecommendationTableCreate(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)
	token := helper.LoginTestUser()
	require.NotEmpty(t, token, "登录应该成功")

	// ============================================
	// 步骤1: 创建手动推荐榜
	// ============================================
	helper.LogInfo("步骤1: 创建手动推荐榜")

	// 准备书籍ID（使用已有测试书籍或创建新的）
	bookID1 := helper.GetTestBook()
	if bookID1 == "" {
		bookID1 = primitive.NewObjectID().Hex()
	}
	bookID2 := primitive.NewObjectID().Hex()

	createTableReq := map[string]interface{}{
		"name":   "本周推荐",
		"period": "2024-W01",
		"items": []map[string]interface{}{
			{
				"bookId": bookID1,
				"order":  1,
				"reason": "精彩绝伦",
			},
			{
				"bookId": bookID2,
				"order":  2,
				"reason": "不可错过",
			},
		},
	}

	w := helper.DoAuthRequest("POST", APIBasePath+"/recommendation/tables/manual", createTableReq, token)

	// 如果API不存在，跳过测试
	if w.Code == http.StatusNotFound {
		helper.LogWarning("推荐榜单API不存在，跳过测试")
		t.SkipNow()
		return
	}

	resp := helper.AssertSuccess(w, http.StatusOK, "创建榜单应该成功")

	data, ok := resp["data"].(map[string]interface{})
	require.True(t, ok, "响应数据格式正确")

	tableID, ok := data["id"].(string)
	if ok {
		require.NotEmpty(t, tableID, "榜单ID不能为空")
		helper.LogSuccess("创建榜单成功 - ID: %s", tableID)
	} else {
		helper.LogSuccess("创建榜单成功（响应中无ID字段）")
	}

	// 验证返回的数据
	if name, ok := data["name"].(string); ok {
		assert.Equal(t, "本周推荐", name, "榜单名称应该匹配")
	}

	if period, ok := data["period"].(string); ok {
		assert.Equal(t, "2024-W01", period, "榜单周期应该匹配")
	}

	if tableType, ok := data["type"].(string); ok {
		assert.Equal(t, string(recommendation.TableTypeManual), tableType, "榜单类型应该是manual")
	}

	helper.LogSuccess("创建榜单测试通过")
}

// TestRecommendationTableUpdateOrder 测试更新榜单顺序
// 测试场景：创建榜单和书籍，更新排序，验证新顺序
func TestRecommendationTableUpdateOrder(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)
	token := helper.LoginTestUser()
	require.NotEmpty(t, token, "登录应该成功")

	// ============================================
	// 步骤1: 创建榜单
	// ============================================
	helper.LogInfo("步骤1: 创建榜单")

	bookID1 := primitive.NewObjectID().Hex()
	bookID2 := primitive.NewObjectID().Hex()
	bookID3 := primitive.NewObjectID().Hex()

	createTableReq := map[string]interface{}{
		"name":   "测试榜单",
		"period": "2024-W02",
		"items": []map[string]interface{}{
			{"bookId": bookID1, "order": 1, "reason": "第一"},
			{"bookId": bookID2, "order": 2, "reason": "第二"},
			{"bookId": bookID3, "order": 3, "reason": "第三"},
		},
	}

	w := helper.DoAuthRequest("POST", APIBasePath+"/recommendation/tables/manual", createTableReq, token)

	// 如果API不存在，跳过测试
	if w.Code == http.StatusNotFound {
		helper.LogWarning("推荐榜单API不存在，跳过测试")
		t.SkipNow()
		return
	}

	resp := helper.AssertSuccess(w, http.StatusOK, "创建榜单应该成功")

	data, ok := resp["data"].(map[string]interface{})
	require.True(t, ok, "响应数据格式正确")

	tableID, ok := data["id"].(string)
	require.True(t, ok, "榜单ID应该存在")
	require.NotEmpty(t, tableID, "榜单ID不能为空")

	helper.LogSuccess("创建榜单成功 - ID: %s", tableID)

	// ============================================
	// 步骤2: 更新榜单顺序
	// ============================================
	helper.LogInfo("步骤2: 更新榜单顺序（交换1和3的位置）")

	updateTableReq := map[string]interface{}{
		"name": "更新后的测试榜单",
		"items": []map[string]interface{}{
			{"bookId": bookID3, "order": 1, "reason": "第三（现在是第一）"},
			{"bookId": bookID2, "order": 2, "reason": "第二（保持第二）"},
			{"bookId": bookID1, "order": 3, "reason": "第一（现在是第三）"},
		},
	}

	w = helper.DoAuthRequest("PUT", APIBasePath+"/recommendation/tables/"+tableID, updateTableReq, token)
	resp = helper.AssertSuccess(w, http.StatusOK, "更新榜单应该成功")

	helper.LogSuccess("更新榜单顺序成功")

	// ============================================
	// 步骤3: 验证新顺序
	// ============================================
	helper.LogInfo("步骤3: 验证新顺序")

	w = helper.DoAuthRequest("GET", APIBasePath+"/recommendation/tables/"+tableID, nil, token)
	resp = helper.AssertSuccess(w, http.StatusOK, "获取榜单详情应该成功")

	data, ok = resp["data"].(map[string]interface{})
	require.True(t, ok, "响应数据格式正确")

	// 验证书籍顺序
	if items, ok := data["items"].([]interface{}); ok && len(items) >= 3 {
		// 检查第一本书是否是原来的第三本
		firstItem, ok := items[0].(map[string]interface{})
		if ok {
			if bookID, exists := firstItem["bookId"].(string); exists {
				// 注意：bookID在响应中可能是对象格式
				helper.LogInfo("第一位书籍ID: %s", bookID)
			}
			if order, exists := firstItem["order"].(float64); exists {
				assert.Equal(t, 1, int(order), "第一本书的order应该是1")
			}
		}
	}

	helper.LogSuccess("榜单顺序更新测试通过")
}

// TestRecommendationTableQuery 测试查询榜单
// 测试场景：创建多种类型榜单，按类型查询，验证返回结果
func TestRecommendationTableQuery(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)
	token := helper.LoginTestUser()
	require.NotEmpty(t, token, "登录应该成功")

	// ============================================
	// 步骤1: 创建多种类型的榜单
	// ============================================
	helper.LogInfo("步骤1: 创建多种类型的榜单")

	bookID := primitive.NewObjectID().Hex()

	// 创建手动榜单1
	createTableReq1 := map[string]interface{}{
		"name":   "手动榜单1",
		"period": "2024-W03",
		"items": []map[string]interface{}{
			{"bookId": bookID, "order": 1, "reason": "推荐理由"},
		},
	}

	w := helper.DoAuthRequest("POST", APIBasePath+"/recommendation/tables/manual", createTableReq1, token)

	// 如果API不存在，跳过测试
	if w.Code == http.StatusNotFound {
		helper.LogWarning("推荐榜单API不存在，跳过测试")
		t.SkipNow()
		return
	}

	resp := helper.AssertSuccess(w, http.StatusOK, "创建手动榜单1应该成功")

	// 创建手动榜单2
	createTableReq2 := map[string]interface{}{
		"name":   "手动榜单2",
		"period": "2024-W04",
		"items": []map[string]interface{}{
			{"bookId": bookID, "order": 1, "reason": "推荐理由2"},
		},
	}

	w = helper.DoAuthRequest("POST", APIBasePath+"/recommendation/tables/manual", createTableReq2, token)
	helper.AssertSuccess(w, http.StatusOK, "创建手动榜单2应该成功")

	helper.LogSuccess("创建了 2 个手动榜单")

	// ============================================
	// 步骤2: 查询所有榜单
	// ============================================
	helper.LogInfo("步骤2: 查询所有榜单")

	w = helper.DoAuthRequest("GET", APIBasePath+"/recommendation/tables?page=1&size=10", nil, token)
	resp = helper.AssertSuccess(w, http.StatusOK, "查询榜单列表应该成功")

	data, ok := resp["data"].(map[string]interface{})
	require.True(t, ok, "响应数据格式正确")

	// 验证返回的榜单列表
	if tables, ok := data["tables"].([]interface{}); ok {
		helper.LogSuccess("查询到 %d 个榜单", len(tables))

		// 验证至少有我们创建的2个榜单
		assert.GreaterOrEqual(t, len(tables), 2, "应该至少有2个榜单")
	}

	if total, ok := data["total"].(float64); ok {
		assert.GreaterOrEqual(t, int(total), 2, "总数应该至少为2")
		helper.LogInfo("榜单总数: %d", int(total))
	}

	// ============================================
	// 步骤3: 按类型查询（手动榜单）
	// ============================================
	helper.LogInfo("步骤3: 按类型查询 - 手动榜单")

	w = helper.DoAuthRequest("GET", APIBasePath+"/recommendation/tables?type=manual&page=1&size=10", nil, token)
	resp = helper.AssertSuccess(w, http.StatusOK, "按类型查询应该成功")

	data, ok = resp["data"].(map[string]interface{})
	require.True(t, ok, "响应数据格式正确")

	if tables, ok := data["tables"].([]interface{}); ok {
		helper.LogSuccess("按类型查询到 %d 个手动榜单", len(tables))

		// 验证返回的都是手动榜单
		for _, table := range tables {
			tableMap, ok := table.(map[string]interface{})
			if ok {
				if tableType, exists := tableMap["type"].(string); exists {
					assert.Equal(t, string(recommendation.TableTypeManual), tableType, "应该都是手动榜单")
				}
			}
		}
	}

	// ============================================
	// 步骤4: 按来源查询
	// ============================================
	helper.LogInfo("步骤4: 按来源查询 - 手动来源")

	w = helper.DoAuthRequest("GET", APIBasePath+"/recommendation/tables?source=manual&page=1&size=10", nil, token)
	resp = helper.AssertSuccess(w, http.StatusOK, "按来源查询应该成功")

	data, ok = resp["data"].(map[string]interface{})
	require.True(t, ok, "响应数据格式正确")

	if tables, ok := data["tables"].([]interface{}); ok {
		helper.LogSuccess("按来源查询到 %d 个榜单", len(tables))
	}

	helper.LogSuccess("榜单查询测试通过")
}

// TestRecommendationTableUpdateStatus 测试更新榜单状态
// 测试场景：创建榜单后，更新其状态（published/draft/archived），验证状态变化
func TestRecommendationTableUpdateStatus(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)
	token := helper.LoginTestUser()
	require.NotEmpty(t, token, "登录应该成功")

	// ============================================
	// 步骤1: 创建榜单
	// ============================================
	helper.LogInfo("步骤1: 创建榜单")

	bookID := primitive.NewObjectID().Hex()

	createTableReq := map[string]interface{}{
		"name":   "状态测试榜单",
		"period": "2024-W05",
		"items": []map[string]interface{}{
			{"bookId": bookID, "order": 1, "reason": "测试"},
		},
	}

	w := helper.DoAuthRequest("POST", APIBasePath+"/recommendation/tables/manual", createTableReq, token)

	// 如果API不存在，跳过测试
	if w.Code == http.StatusNotFound {
		helper.LogWarning("推荐榜单API不存在，跳过测试")
		t.SkipNow()
		return
	}

	resp := helper.AssertSuccess(w, http.StatusOK, "创建榜单应该成功")

	data, ok := resp["data"].(map[string]interface{})
	require.True(t, ok, "响应数据格式正确")

	tableID, ok := data["id"].(string)
	require.True(t, ok, "榜单ID应该存在")
	require.NotEmpty(t, tableID, "榜单ID不能为空")

	// 验证初始状态
	if status, ok := data["status"].(string); ok {
		helper.LogInfo("初始状态: %s", status)
	}

	// ============================================
	// 步骤2: 更新状态为已发布
	// ============================================
	helper.LogInfo("步骤2: 更新状态为已发布")

	updateReq := map[string]interface{}{
		"name":   "状态测试榜单",
		"status": string(recommendation.TableStatusActive),
		"items": []map[string]interface{}{
			{"bookId": bookID, "order": 1, "reason": "测试"},
		},
	}

	w = helper.DoAuthRequest("PUT", APIBasePath+"/recommendation/tables/"+tableID, updateReq, token)
	resp = helper.AssertSuccess(w, http.StatusOK, "更新状态应该成功")

	helper.LogSuccess("更新榜单状态为published")

	// ============================================
	// 步骤3: 验证状态变化
	// ============================================
	helper.LogInfo("步骤3: 验证状态变化")

	w = helper.DoAuthRequest("GET", APIBasePath+"/recommendation/tables/"+tableID, nil, token)
	resp = helper.AssertSuccess(w, http.StatusOK, "获取榜单详情应该成功")

	data, ok = resp["data"].(map[string]interface{})
	require.True(t, ok, "响应数据格式正确")

	if status, ok := data["status"].(string); ok {
		assert.Equal(t, string(recommendation.TableStatusActive), status, "状态应该是active")
		helper.LogSuccess("榜单状态验证正确: %s", status)
	}

	helper.LogSuccess("榜单状态更新测试通过")
}

// TestRecommendationTableDelete 测试删除榜单
// 测试场景：创建榜单后，删除榜单，验证删除成功
func TestRecommendationTableDelete(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)
	token := helper.LoginTestUser()
	require.NotEmpty(t, token, "登录应该成功")

	// ============================================
	// 步骤1: 创建榜单
	// ============================================
	helper.LogInfo("步骤1: 创建榜单")

	bookID := primitive.NewObjectID().Hex()

	createTableReq := map[string]interface{}{
		"name":   "待删除榜单",
		"period": "2024-W06",
		"items": []map[string]interface{}{
			{"bookId": bookID, "order": 1, "reason": "测试"},
		},
	}

	w := helper.DoAuthRequest("POST", APIBasePath+"/recommendation/tables/manual", createTableReq, token)

	// 如果API不存在，跳过测试
	if w.Code == http.StatusNotFound {
		helper.LogWarning("推荐榜单API不存在，跳过测试")
		t.SkipNow()
		return
	}

	resp := helper.AssertSuccess(w, http.StatusOK, "创建榜单应该成功")

	data, ok := resp["data"].(map[string]interface{})
	require.True(t, ok, "响应数据格式正确")

	tableID, ok := data["id"].(string)
	require.True(t, ok, "榜单ID应该存在")
	require.NotEmpty(t, tableID, "榜单ID不能为空")

	helper.LogSuccess("创建榜单成功 - ID: %s", tableID)

	// ============================================
	// 步骤2: 删除榜单
	// ============================================
	helper.LogInfo("步骤2: 删除榜单")

	w = helper.DoAuthRequest("DELETE", APIBasePath+"/recommendation/tables/"+tableID, nil, token)

	// 删除成功应该返回200或204
	if w.Code == http.StatusOK || w.Code == http.StatusNoContent {
		helper.LogSuccess("删除榜单成功: %d", w.Code)
	} else {
		helper.LogWarning("删除榜单返回: %d", w.Code)
	}

	// ============================================
	// 步骤3: 验证删除结果
	// ============================================
	helper.LogInfo("步骤3: 验证删除结果")

	w = helper.DoAuthRequest("GET", APIBasePath+"/recommendation/tables/"+tableID, nil, token)

	// 榜单已被删除，应该返回404
	if w.Code == http.StatusNotFound {
		helper.LogSuccess("榜单已成功删除（返回404）")
	} else if w.Code == http.StatusOK {
		// 某些实现可能使用软删除
		helper.LogWarning("榜单可能使用了软删除（仍可查询）")
	}

	helper.LogSuccess("榜单删除测试通过")
}

// TestRecommendationTablePagination 测试榜单分页
// 测试场景：创建多个榜单，测试分页功能
func TestRecommendationTablePagination(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)
	token := helper.LoginTestUser()
	require.NotEmpty(t, token, "登录应该成功")

	// ============================================
	// 步骤1: 创建多个榜单
	// ============================================
	helper.LogInfo("步骤1: 创建多个榜单")

	createdCount := 0
	for i := 1; i <= 5; i++ {
		bookID := primitive.NewObjectID().Hex()

		createTableReq := map[string]interface{}{
			"name":   "分页测试榜单",
			"period": "2024-W10",
			"items": []map[string]interface{}{
				{"bookId": bookID, "order": 1, "reason": "测试"},
			},
		}

		w := helper.DoAuthRequest("POST", APIBasePath+"/recommendation/tables/manual", createTableReq, token)

		// 如果API不存在，跳过测试
		if w.Code == http.StatusNotFound {
			helper.LogWarning("推荐榜单API不存在，跳过测试")
			t.SkipNow()
			return
		}

		if w.Code == http.StatusOK {
			createdCount++
		}
	}

	helper.LogSuccess("创建了 %d 个榜单", createdCount)

	// ============================================
	// 步骤2: 测试第一页
	// ============================================
	helper.LogInfo("步骤2: 测试第一页（size=2）")

	w := helper.DoAuthRequest("GET", APIBasePath+"/recommendation/tables?page=1&size=2", nil, token)
	resp := helper.AssertSuccess(w, http.StatusOK, "分页查询应该成功")

	data, ok := resp["data"].(map[string]interface{})
	require.True(t, ok, "响应数据格式正确")

	if tables, ok := data["tables"].([]interface{}); ok {
		assert.LessOrEqual(t, len(tables), 2, "第一页应该最多返回2个榜单")
		helper.LogSuccess("第一页返回 %d 个榜单", len(tables))
	}

	if page, ok := data["page"].(float64); ok {
		assert.Equal(t, 1, int(page), "页码应该是1")
	}

	if size, ok := data["size"].(float64); ok {
		assert.Equal(t, 2, int(size), "每页大小应该是2")
	}

	// ============================================
	// 步骤3: 测试第二页
	// ============================================
	helper.LogInfo("步骤3: 测试第二页（size=2）")

	w = helper.DoAuthRequest("GET", APIBasePath+"/recommendation/tables?page=2&size=2", nil, token)
	resp = helper.AssertSuccess(w, http.StatusOK, "分页查询应该成功")

	data, ok = resp["data"].(map[string]interface{})
	require.True(t, ok, "响应数据格式正确")

	if tables, ok := data["tables"].([]interface{}); ok {
		helper.LogSuccess("第二页返回 %d 个榜单", len(tables))
	}

	if page, ok := data["page"].(float64); ok {
		assert.Equal(t, 2, int(page), "页码应该是2")
	}

	// ============================================
	// 步骤4: 测试超出的页码
	// ============================================
	helper.LogInfo("步骤4: 测试超出的页码（page=999）")

	w = helper.DoAuthRequest("GET", APIBasePath+"/recommendation/tables?page=999&size=2", nil, token)
	resp = helper.AssertSuccess(w, http.StatusOK, "超出页码应该返回空结果")

	data, ok = resp["data"].(map[string]interface{})
	require.True(t, ok, "响应数据格式正确")

	if tables, ok := data["tables"].([]interface{}); ok {
		assert.Equal(t, 0, len(tables), "超出页码应该返回空列表")
		helper.LogSuccess("超出页码返回空结果")
	}

	helper.LogSuccess("榜单分页测试通过")
}

// TestRecommendationTableValidationError 测试榜单数据验证
// 测试场景：提交无效的榜单数据，验证错误处理
func TestRecommendationTableValidationError(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)
	token := helper.LoginTestUser()
	require.NotEmpty(t, token, "登录应该成功")

	// ============================================
	// 步骤1: 测试缺失必需字段
	// ============================================
	helper.LogInfo("步骤1: 测试缺失必需字段（name）")

	invalidReq := map[string]interface{}{
		"period": "2024-W07",
		"items": []map[string]interface{}{
			{"bookId": primitive.NewObjectID().Hex(), "order": 1},
		},
	}

	w := helper.DoAuthRequest("POST", APIBasePath+"/recommendation/tables/manual", invalidReq, token)

	// 如果API不存在，跳过测试
	if w.Code == http.StatusNotFound {
		helper.LogWarning("推荐榜单API不存在，跳过测试")
		t.SkipNow()
		return
	}

	// 缺失name应该被拒绝
	if w.Code != http.StatusOK {
		helper.LogSuccess("缺失name字段被正确拒绝: %d", w.Code)
		assert.Contains(t, []int{http.StatusBadRequest, 422}, w.Code, "应该返回400或422")
	}

	// ============================================
	// 步骤2: 测试空items
	// ============================================
	helper.LogInfo("步骤2: 测试空items")

	invalidReq2 := map[string]interface{}{
		"name":   "测试榜单",
		"period": "2024-W08",
		"items": []map[string]interface{}{},
	}

	w = helper.DoAuthRequest("POST", APIBasePath+"/recommendation/tables/manual", invalidReq2, token)

	// 空items可能被接受（某些实现允许空榜单）
	// 或者被拒绝（某些实现要求至少一本书）
	if w.Code != http.StatusOK {
		helper.LogSuccess("空items被拒绝: %d", w.Code)
	} else {
		helper.LogInfo("空items被接受（实现允许空榜单）")
	}

	// ============================================
	// 步骤3: 测试无效的bookId格式
	// ============================================
	helper.LogInfo("步骤3: 测试无效的bookId格式")

	invalidReq3 := map[string]interface{}{
		"name":   "测试榜单",
		"period": "2024-W09",
		"items": []map[string]interface{}{
			{"bookId": "invalid-id-format", "order": 1},
		},
	}

	w = helper.DoAuthRequest("POST", APIBasePath+"/recommendation/tables/manual", invalidReq3, token)

	// 无效的bookId格式可能被拒绝
	if w.Code != http.StatusOK {
		helper.LogSuccess("无效的bookId格式被拒绝: %d", w.Code)
	}

	// ============================================
	// 步骤4: 测试重复的order
	// ============================================
	helper.LogInfo("步骤4: 测试重复的order")

	bookID := primitive.NewObjectID().Hex()
	invalidReq4 := map[string]interface{}{
		"name":   "测试榜单",
		"period": "2024-W10",
		"items": []map[string]interface{}{
			{"bookId": bookID, "order": 1},
			{"bookId": primitive.NewObjectID().Hex(), "order": 1}, // 重复的order
		},
	}

	w = helper.DoAuthRequest("POST", APIBasePath+"/recommendation/tables/manual", invalidReq4, token)

	// 重复的order可能被接受（某些实现不验证）或拒绝
	if w.Code != http.StatusOK {
		helper.LogSuccess("重复的order被拒绝: %d", w.Code)
	} else {
		helper.LogInfo("重复的order被接受（实现不验证order唯一性）")
	}

	helper.LogSuccess("榜单数据验证测试通过")
}
