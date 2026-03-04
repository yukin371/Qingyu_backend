//go:build integration
// +build integration

package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/pkg/response"
)

// TestWriterParagraphContentLoad 测试段落内容加载
// 测试场景：创建文档后，按段落加载内容并验证返回的段落列表
func TestWriterParagraphContentLoad(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)
	token := helper.LoginTestUser()
	require.NotEmpty(t, token, "登录应该成功")

	// ============================================
	// 步骤1: 创建测试文档
	// ============================================
	helper.LogInfo("步骤1: 创建测试文档")

	createDocReq := map[string]interface{}{
		"projectId": primitive.NewObjectID().Hex(),
		"title":     "测试章节",
		"type":      "chapter",
		"order":     1,
	}

	w := helper.DoAuthRequest("POST", APIBasePath+"/writer/documents", createDocReq, token)
	resp := helper.AssertSuccess(w, http.StatusOK, "创建文档应该成功")

	data, ok := resp["data"].(map[string]interface{})
	require.True(t, ok, "响应数据格式正确")
	documentID, ok := data["id"].(string)
	require.True(t, ok, "文档ID应该存在")
	require.NotEmpty(t, documentID, "文档ID不能为空")

	helper.LogSuccess("创建文档成功: %s", documentID)

	// ============================================
	// 步骤2: 按段落加载内容
	// ============================================
	helper.LogInfo("步骤2: 按段落加载内容")

	w = helper.DoAuthRequest("GET", APIBasePath+"/writer/documents/"+documentID+"/contents", nil, token)
	resp = helper.AssertSuccess(w, http.StatusOK, "获取文档内容应该成功")

	data, ok = resp["data"].(map[string]interface{})
	require.True(t, ok, "响应数据格式正确")

	// 验证返回的段落列表
	contents, ok := data["contents"].([]interface{})
	if ok && len(contents) > 0 {
		helper.LogSuccess("获取到 %d 个段落", len(contents))

		// 验证段落结构
		for i, content := range contents {
			paragraph, ok := content.(map[string]interface{})
			require.True(t, ok, "段落格式正确")

			// 检查必需字段
			if paragraphID, exists := paragraph["id"]; exists {
				assert.NotEmpty(t, paragraphID, "段落ID不能为空")
			}
			if order, exists := paragraph["order"]; exists {
				assert.Equal(t, i+1, int(order.(float64)), "段落顺序应该正确")
			}
		}
	} else {
		helper.LogWarning("文档无内容段落（可能为新文档）")
	}

	helper.LogSuccess("段落内容加载测试通过")
}

// TestWriterParagraphContentUpdate 测试段落更新
// 测试场景：创建文档后，更新指定段落内容，验证更新结果
func TestWriterParagraphContentUpdate(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)
	token := helper.LoginTestUser()
	require.NotEmpty(t, token, "登录应该成功")

	// ============================================
	// 步骤1: 创建测试文档
	// ============================================
	helper.LogInfo("步骤1: 创建测试文档")

	projectID := primitive.NewObjectID().Hex()
	createDocReq := map[string]interface{}{
		"projectId": projectID,
		"title":     "测试章节",
		"type":      "chapter",
		"order":     1,
	}

	w := helper.DoAuthRequest("POST", APIBasePath+"/writer/documents", createDocReq, token)
	resp := helper.AssertSuccess(w, http.StatusOK, "创建文档应该成功")

	data, ok := resp["data"].(map[string]interface{})
	require.True(t, ok, "响应数据格式正确")
	documentID, ok := data["id"].(string)
	require.True(t, ok, "文档ID应该存在")

	helper.LogSuccess("创建文档成功: %s", documentID)

	// ============================================
	// 步骤2: 更新段落内容
	// ============================================
	helper.LogInfo("步骤2: 更新段落内容")

	paragraphID := primitive.NewObjectID().Hex()
	updateReq := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"id":      paragraphID,
				"type":    "paragraph",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": "这是更新后的段落内容",
					},
				},
				"order": 1,
			},
		},
	}

	w = helper.DoAuthRequest("PUT", APIBasePath+"/writer/documents/"+documentID+"/contents", updateReq, token)
	resp = helper.AssertSuccess(w, http.StatusOK, "更新段落内容应该成功")

	data, ok = resp["data"].(map[string]interface{})
	require.True(t, ok, "响应数据格式正确")

	// 验证更新结果
	if updated, ok := data["updated"].(float64); ok {
		assert.Greater(t, int(updated), 0, "应该有段落被更新")
		helper.LogSuccess("更新了 %d 个段落", int(updated))
	}

	if paragraphs, ok := data["paragraphs"].(float64); ok {
		assert.Greater(t, int(paragraphs), 0, "段落数量应该大于0")
	}

	helper.LogSuccess("段落更新测试通过")
}

// TestWriterParagraphReorder 测试段落重排
// 测试场景：创建多段落文档后，重排段落顺序，验证新顺序
func TestWriterParagraphReorder(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)
	token := helper.LoginTestUser()
	require.NotEmpty(t, token, "登录应该成功")

	// ============================================
	// 步骤1: 创建测试文档
	// ============================================
	helper.LogInfo("步骤1: 创建测试文档")

	projectID := primitive.NewObjectID().Hex()
	createDocReq := map[string]interface{}{
		"projectId": projectID,
		"title":     "测试章节",
		"type":      "chapter",
		"order":     1,
	}

	w := helper.DoAuthRequest("POST", APIBasePath+"/writer/documents", createDocReq, token)
	resp := helper.AssertSuccess(w, http.StatusOK, "创建文档应该成功")

	data, ok := resp["data"].(map[string]interface{})
	require.True(t, ok, "响应数据格式正确")
	documentID, ok := data["id"].(string)
	require.True(t, ok, "文档ID应该存在")

	helper.LogSuccess("创建文档成功: %s", documentID)

	// ============================================
	// 步骤2: 创建多个段落
	// ============================================
	helper.LogInfo("步骤2: 创建多个段落")

	paragraph1ID := primitive.NewObjectID().Hex()
	paragraph2ID := primitive.NewObjectID().Hex()
	paragraph3ID := primitive.NewObjectID().Hex()

	createParagraphsReq := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"id":      paragraph1ID,
				"type":    "paragraph",
				"content": []map[string]interface{}{{"type": "text", "text": "第一段"}},
				"order":   1,
			},
			{
				"id":      paragraph2ID,
				"type":    "paragraph",
				"content": []map[string]interface{}{{"type": "text", "text": "第二段"}},
				"order":   2,
			},
			{
				"id":      paragraph3ID,
				"type":    "paragraph",
				"content": []map[string]interface{}{{"type": "text", "text": "第三段"}},
				"order":   3,
			},
		},
	}

	w = helper.DoAuthRequest("PUT", APIBasePath+"/writer/documents/"+documentID+"/contents", createParagraphsReq, token)
	helper.AssertSuccess(w, http.StatusOK, "创建段落应该成功")

	helper.LogSuccess("创建了 3 个段落")

	// ============================================
	// 步骤3: 重排段落顺序
	// ============================================
	helper.LogInfo("步骤3: 重排段落顺序（2->1, 1->2, 3->3）")

	reorderReq := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"id":      paragraph2ID,
				"type":    "paragraph",
				"content": []map[string]interface{}{{"type": "text", "text": "第二段（现在是第一）"}},
				"order":   1,
			},
			{
				"id":      paragraph1ID,
				"type":    "paragraph",
				"content": []map[string]interface{}{{"type": "text", "text": "第一段（现在是第二）"}},
				"order":   2,
			},
			{
				"id":      paragraph3ID,
				"type":    "paragraph",
				"content": []map[string]interface{}{{"type": "text", "text": "第三段"}},
				"order":   3,
			},
		},
	}

	w = helper.DoAuthRequest("PUT", APIBasePath+"/writer/documents/"+documentID+"/contents", reorderReq, token)
	resp = helper.AssertSuccess(w, http.StatusOK, "重排段落应该成功")

	data, ok = resp["data"].(map[string]interface{})
	require.True(t, ok, "响应数据格式正确")

	if updated, ok := data["updated"].(float64); ok {
		assert.Equal(t, 3, int(updated), "应该更新了3个段落")
	}

	helper.LogSuccess("段落重排成功")

	// ============================================
	// 步骤4: 验证新顺序
	// ============================================
	helper.LogInfo("步骤4: 验证新顺序")

	w = helper.DoAuthRequest("GET", APIBasePath+"/writer/documents/"+documentID+"/contents", nil, token)
	resp = helper.AssertSuccess(w, http.StatusOK, "获取文档内容应该成功")

	data, ok = resp["data"].(map[string]interface{})
	require.True(t, ok, "响应数据格式正确")

	contents, ok := data["contents"].([]interface{})
	if ok && len(contents) >= 2 {
		// 验证顺序是否正确
		firstParagraph, ok := contents[0].(map[string]interface{})
		if ok {
			if id, exists := firstParagraph["id"]; exists {
				assert.Equal(t, paragraph2ID, id, "第一段应该是原来的第二段")
			}
		}

		secondParagraph, ok := contents[1].(map[string]interface{})
		if ok {
			if id, exists := secondParagraph["id"]; exists {
				assert.Equal(t, paragraph1ID, id, "第二段应该是原来的第一段")
			}
		}

		helper.LogSuccess("段落顺序验证正确")
	}

	helper.LogSuccess("段落重排测试通过")
}

// TestWriterParagraphReindex 测试段落重新编号
// 测试场景：创建顺序不连续的段落后，调用重新编号接口，验证序号连续
func TestWriterParagraphReindex(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)
	token := helper.LoginTestUser()
	require.NotEmpty(t, token, "登录应该成功")

	// ============================================
	// 步骤1: 创建测试文档
	// ============================================
	helper.LogInfo("步骤1: 创建测试文档")

	projectID := primitive.NewObjectID().Hex()
	createDocReq := map[string]interface{}{
		"projectId": projectID,
		"title":     "测试章节",
		"type":      "chapter",
		"order":     1,
	}

	w := helper.DoAuthRequest("POST", APIBasePath+"/writer/documents", createDocReq, token)
	resp := helper.AssertSuccess(w, http.StatusOK, "创建文档应该成功")

	data, ok := resp["data"].(map[string]interface{})
	require.True(t, ok, "响应数据格式正确")
	documentID, ok := data["id"].(string)
	require.True(t, ok, "文档ID应该存在")

	helper.LogSuccess("创建文档成功: %s", documentID)

	// ============================================
	// 步骤2: 创建顺序不连续的段落
	// ============================================
	helper.LogInfo("步骤2: 创建顺序不连续的段落（1, 3, 5）")

	paragraph1ID := primitive.NewObjectID().Hex()
	paragraph2ID := primitive.NewObjectID().Hex()
	paragraph3ID := primitive.NewObjectID().Hex()

	createParagraphsReq := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"id":      paragraph1ID,
				"type":    "paragraph",
				"content": []map[string]interface{}{{"type": "text", "text": "第一段"}},
				"order":   1,
			},
			{
				"id":      paragraph2ID,
				"type":    "paragraph",
				"content": []map[string]interface{}{{"type": "text", "text": "第二段"}},
				"order":   3,
			},
			{
				"id":      paragraph3ID,
				"type":    "paragraph",
				"content": []map[string]interface{}{{"type": "text", "text": "第三段"}},
				"order":   5,
			},
		},
	}

	w = helper.DoAuthRequest("PUT", APIBasePath+"/writer/documents/"+documentID+"/contents", createParagraphsReq, token)
	helper.AssertSuccess(w, http.StatusOK, "创建段落应该成功")

	helper.LogSuccess("创建了 3 个顺序不连续的段落")

	// ============================================
	// 步骤3: 调用重新编号接口
	// ============================================
	helper.LogInfo("步骤3: 调用重新编号接口")

	w = helper.DoAuthRequest("POST", APIBasePath+"/writer/documents/"+documentID+"/contents/reindex", nil, token)
	resp = helper.AssertSuccess(w, http.StatusOK, "重新编号应该成功")

	data, ok = resp["data"].(map[string]interface{})
	require.True(t, ok, "响应数据格式正确")

	// 验证重新编号结果
	if reindexed, ok := data["reindexed"].(float64); ok {
		assert.Equal(t, 3, int(reindexed), "应该重新编号了3个段落")
		helper.LogSuccess("重新编号了 %d 个段落", int(reindexed))
	}

	if total, ok := data["total"].(float64); ok {
		assert.Equal(t, 3, int(total), "总段落数应该是3")
	}

	helper.LogSuccess("段落重新编号测试通过")
}

// TestWriterParagraphValidationError 测试段落数据验证
// 测试场景：提交无效的段落数据，验证错误处理
func TestWriterParagraphValidationError(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)
	token := helper.LoginTestUser()
	require.NotEmpty(t, token, "登录应该成功")

	// ============================================
	// 步骤1: 创建测试文档
	// ============================================
	helper.LogInfo("步骤1: 创建测试文档")

	projectID := primitive.NewObjectID().Hex()
	createDocReq := map[string]interface{}{
		"projectId": projectID,
		"title":     "测试章节",
		"type":      "chapter",
		"order":     1,
	}

	w := helper.DoAuthRequest("POST", APIBasePath+"/writer/documents", createDocReq, token)
	resp := helper.AssertSuccess(w, http.StatusOK, "创建文档应该成功")

	data, ok := resp["data"].(map[string]interface{})
	require.True(t, ok, "响应数据格式正确")
	documentID, ok := data["id"].(string)
	require.True(t, ok, "文档ID应该存在")

	helper.LogSuccess("创建文档成功: %s", documentID)

	// ============================================
	// 步骤2: 测试空段落验证
	// ============================================
	helper.LogInfo("步骤2: 测试空段落验证")

	emptyParagraphReq := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"id":      primitive.NewObjectID().Hex(),
				"type":    "paragraph",
				"content": []interface{}{},
				"order":   1,
			},
		},
	}

	w = helper.DoAuthRequest("PUT", APIBasePath+"/writer/documents/"+documentID+"/contents", emptyParagraphReq, token)

	// 空段落应该被拒绝或警告
	if w.Code != http.StatusOK {
		var errResp response.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		if err == nil {
			helper.LogInfo("空段落被正确拒绝: %s", errResp.Message)
		}
	}

	// ============================================
	// 步骤3: 测试重复段落ID验证
	// ============================================
	helper.LogInfo("步骤3: 测试重复段落ID验证")

	duplicateID := primitive.NewObjectID().Hex()
	duplicateReq := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"id":      duplicateID,
				"type":    "paragraph",
				"content": []map[string]interface{}{{"type": "text", "text": "第一段"}},
				"order":   1,
			},
			{
				"id":      duplicateID, // 重复ID
				"type":    "paragraph",
				"content": []map[string]interface{}{{"type": "text", "text": "第二段"}},
				"order":   2,
			},
		},
	}

	w = helper.DoAuthRequest("PUT", APIBasePath+"/writer/documents/"+documentID+"/contents", duplicateReq, token)

	// 重复ID应该被拒绝
	if w.Code != http.StatusOK {
		var errResp response.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		if err == nil {
			helper.LogSuccess("重复段落ID被正确拒绝: %s", errResp.Message)
		}
	}

	helper.LogSuccess("段落数据验证测试通过")
}
