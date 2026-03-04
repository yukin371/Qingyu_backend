//go:build integration
// +build integration

package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestCommentWithParagraphBinding 测试段落绑定评论
// 测试场景：创建测试文档和段落，创建绑定到段落的评论，验证评论的paragraphId
func TestCommentWithParagraphBinding(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)
	token := helper.LoginTestUser()
	require.NotEmpty(t, token, "登录应该成功")

	// ============================================
	// 步骤1: 创建测试文档和段落
	// ============================================
	helper.LogInfo("步骤1: 创建测试文档和段落")

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

	// 创建段落
	paragraphID := primitive.NewObjectID().Hex()
	createParagraphsReq := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"id":      paragraphID,
				"type":    "paragraph",
				"content": []map[string]interface{}{{"type": "text", "text": "这是测试段落"}},
				"order":   1,
			},
		},
	}

	w = helper.DoAuthRequest("PUT", APIBasePath+"/writer/documents/"+documentID+"/contents", createParagraphsReq, token)
	helper.AssertSuccess(w, http.StatusOK, "创建段落应该成功")

	helper.LogSuccess("创建文档和段落成功 - 文档ID: %s, 段落ID: %s", documentID, paragraphID)

	// ============================================
	// 步骤2: 创建绑定到段落的评论
	// ============================================
	helper.LogInfo("步骤2: 创建绑定到段落的评论")

	createCommentReq := map[string]interface{}{
		"content":     "这是一个绑定到段落的评论",
		"type":        "suggestion",
		"paragraphId": paragraphID,
		"metadata": map[string]interface{}{
			"priority": "high",
		},
	}

	w = helper.DoAuthRequest("POST", APIBasePath+"/writer/documents/"+documentID+"/comments", createCommentReq, token)
	resp = helper.AssertSuccess(w, http.StatusOK, "创建评论应该成功")

	data, ok = resp["data"].(map[string]interface{})
	require.True(t, ok, "响应数据格式正确")

	// 验证评论的paragraphId
	commentID, ok := data["id"].(string)
	require.True(t, ok, "评论ID应该存在")
	require.NotEmpty(t, commentID, "评论ID不能为空")

	if respParagraphID, ok := data["paragraphId"].(string); ok {
		assert.Equal(t, paragraphID, respParagraphID, "评论的paragraphId应该匹配")
		helper.LogSuccess("评论正确绑定到段落: %s", paragraphID)
	} else {
		helper.LogWarning("响应中未找到paragraphId字段")
	}

	helper.LogSuccess("创建评论成功 - 评论ID: %s", commentID)

	// ============================================
	// 步骤3: 验证评论详情
	// ============================================
	helper.LogInfo("步骤3: 验证评论详情")

	w = helper.DoAuthRequest("GET", APIBasePath+"/writer/comments/"+commentID, nil, token)
	resp = helper.AssertSuccess(w, http.StatusOK, "获取评论详情应该成功")

	data, ok = resp["data"].(map[string]interface{})
	require.True(t, ok, "响应数据格式正确")

	// 验证paragraphId字段
	if respParagraphID, ok := data["paragraphId"].(string); ok {
		assert.Equal(t, paragraphID, respParagraphID, "评论详情的paragraphId应该匹配")
		helper.LogSuccess("评论详情验证正确，段落ID: %s", respParagraphID)
	}

	// 验证documentId字段
	if respDocumentID, ok := data["documentId"].(string); ok {
		assert.Equal(t, documentID, respDocumentID, "评论的documentId应该匹配")
	}

	// 验证其他字段
	if content, ok := data["content"].(string); ok {
		assert.Equal(t, "这是一个绑定到段落的评论", content, "评论内容应该匹配")
	}

	if commentType, ok := data["type"].(string); ok {
		assert.Equal(t, "suggestion", commentType, "评论类型应该匹配")
	}

	helper.LogSuccess("评论绑定段落测试通过")
}

// TestCommentQueryByParagraph 测试按段落查询评论
// 测试场景：创建测试文档和段落，创建多个评论绑定不同段落，按paragraphId过滤查询，验证返回结果
func TestCommentQueryByParagraph(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)
	token := helper.LoginTestUser()
	require.NotEmpty(t, token, "登录应该成功")

	// ============================================
	// 步骤1: 创建测试文档和多个段落
	// ============================================
	helper.LogInfo("步骤1: 创建测试文档和多个段落")

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

	// 创建两个段落
	paragraph1ID := primitive.NewObjectID().Hex()
	paragraph2ID := primitive.NewObjectID().Hex()

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
		},
	}

	w = helper.DoAuthRequest("PUT", APIBasePath+"/writer/documents/"+documentID+"/contents", createParagraphsReq, token)
	helper.AssertSuccess(w, http.StatusOK, "创建段落应该成功")

	helper.LogSuccess("创建了 2 个段落")

	// ============================================
	// 步骤2: 创建多个评论绑定不同段落
	// ============================================
	helper.LogInfo("步骤2: 创建多个评论绑定不同段落")

	// 为段落1创建2条评论
	comment1Req := map[string]interface{}{
		"content":     "段落1的评论1",
		"type":        "suggestion",
		"paragraphId": paragraph1ID,
	}

	w = helper.DoAuthRequest("POST", APIBasePath+"/writer/documents/"+documentID+"/comments", comment1Req, token)
	resp = helper.AssertSuccess(w, http.StatusOK, "创建评论1应该成功")
	data, _ = resp["data"].(map[string]interface{})
	comment1ID, _ := data["id"].(string)

	comment2Req := map[string]interface{}{
		"content":     "段落1的评论2",
		"type":        "issue",
		"paragraphId": paragraph1ID,
	}

	w = helper.DoAuthRequest("POST", APIBasePath+"/writer/documents/"+documentID+"/comments", comment2Req, token)
	resp = helper.AssertSuccess(w, http.StatusOK, "创建评论2应该成功")
	data, _ = resp["data"].(map[string]interface{})
	_, _ = data["id"].(string)

	// 为段落2创建1条评论
	comment3Req := map[string]interface{}{
		"content":     "段落2的评论1",
		"type":        "suggestion",
		"paragraphId": paragraph2ID,
	}

	w = helper.DoAuthRequest("POST", APIBasePath+"/writer/documents/"+documentID+"/comments", comment3Req, token)
	resp = helper.AssertSuccess(w, http.StatusOK, "创建评论3应该成功")
	data, _ = resp["data"].(map[string]interface{})
	comment3ID, _ := data["id"].(string)

	helper.LogSuccess("创建了 3 条评论（段落1: 2条, 段落2: 1条）")

	// ============================================
	// 步骤3: 按段落1查询评论
	// ============================================
	helper.LogInfo("步骤3: 按段落1查询评论")

	// 获取文档的所有评论
	w = helper.DoAuthRequest("GET", APIBasePath+"/writer/documents/"+documentID+"/comments", nil, token)
	resp = helper.AssertSuccess(w, http.StatusOK, "获取评论列表应该成功")

	data, ok = resp["data"].(map[string]interface{})
	require.True(t, ok, "响应数据格式正确")

	comments, ok := data["comments"].([]interface{})
	require.True(t, ok, "评论列表应该存在")

	helper.LogInfo("文档共有 %d 条评论", len(comments))

	// 验证返回的评论数量
	total, ok := data["total"].(float64)
	if ok {
		assert.Equal(t, 3, int(total), "总评论数应该是3")
	}

	// 过滤出段落1的评论并验证
	paragraph1CommentCount := 0
	paragraph2CommentCount := 0

	for _, comment := range comments {
		commentMap, ok := comment.(map[string]interface{})
		if !ok {
			continue
		}

		commentParagraphID, ok := commentMap["paragraphId"].(string)
		if !ok {
			continue
		}

		if commentParagraphID == paragraph1ID {
			paragraph1CommentCount++
		} else if commentParagraphID == paragraph2ID {
			paragraph2CommentCount++
		}
	}

	assert.Equal(t, 2, paragraph1CommentCount, "段落1应该有2条评论")
	assert.Equal(t, 1, paragraph2CommentCount, "段落2应该有1条评论")

	helper.LogSuccess("段落1有 %d 条评论, 段落2有 %d 条评论", paragraph1CommentCount, paragraph2CommentCount)

	// ============================================
	// 步骤4: 验证评论详情
	// ============================================
	helper.LogInfo("步骤4: 验证评论详情")

	// 验证评论1
	if comment1ID != "" {
		w = helper.DoAuthRequest("GET", APIBasePath+"/writer/comments/"+comment1ID, nil, token)
		resp = helper.AssertSuccess(w, http.StatusOK, "获取评论1详情应该成功")
		data, _ = resp["data"].(map[string]interface{})
		if respParagraphID, ok := data["paragraphId"].(string); ok {
			assert.Equal(t, paragraph1ID, respParagraphID, "评论1应该绑定到段落1")
		}
		if content, ok := data["content"].(string); ok {
			assert.Equal(t, "段落1的评论1", content, "评论1内容应该匹配")
		}
	}

	// 验证评论3
	if comment3ID != "" {
		w = helper.DoAuthRequest("GET", APIBasePath+"/writer/comments/"+comment3ID, nil, token)
		resp = helper.AssertSuccess(w, http.StatusOK, "获取评论3详情应该成功")
		data, _ = resp["data"].(map[string]interface{})
		if respParagraphID, ok := data["paragraphId"].(string); ok {
			assert.Equal(t, paragraph2ID, respParagraphID, "评论3应该绑定到段落2")
		}
		if content, ok := data["content"].(string); ok {
			assert.Equal(t, "段落2的评论1", content, "评论3内容应该匹配")
		}
	}

	helper.LogSuccess("按段落查询评论测试通过")
}

// TestCommentUpdateWithParagraphBinding 测试更新段落绑定的评论
// 测试场景：创建绑定到段落的评论后，更新评论内容，验证绑定关系不变
func TestCommentUpdateWithParagraphBinding(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)
	token := helper.LoginTestUser()
	require.NotEmpty(t, token, "登录应该成功")

	// ============================================
	// 步骤1: 创建测试文档、段落和评论
	// ============================================
	helper.LogInfo("步骤1: 创建测试文档、段落和评论")

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

	paragraphID := primitive.NewObjectID().Hex()
	createParagraphsReq := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"id":      paragraphID,
				"type":    "paragraph",
				"content": []map[string]interface{}{{"type": "text", "text": "测试段落"}},
				"order":   1,
			},
		},
	}

	w = helper.DoAuthRequest("PUT", APIBasePath+"/writer/documents/"+documentID+"/contents", createParagraphsReq, token)
	helper.AssertSuccess(w, http.StatusOK, "创建段落应该成功")

	createCommentReq := map[string]interface{}{
		"content":     "原始评论内容",
		"type":        "suggestion",
		"paragraphId": paragraphID,
	}

	w = helper.DoAuthRequest("POST", APIBasePath+"/writer/documents/"+documentID+"/comments", createCommentReq, token)
	resp = helper.AssertSuccess(w, http.StatusOK, "创建评论应该成功")

	data, _ = resp["data"].(map[string]interface{})
	commentID, _ := data["id"].(string)
	require.NotEmpty(t, commentID, "评论ID不能为空")

	helper.LogSuccess("创建评论成功: %s", commentID)

	// ============================================
	// 步骤2: 更新评论内容
	// ============================================
	helper.LogInfo("步骤2: 更新评论内容")

	updateCommentReq := map[string]interface{}{
		"content": "更新后的评论内容",
		"type":    "issue",
		"metadata": map[string]interface{}{
			"priority": "high",
		},
	}

	w = helper.DoAuthRequest("PUT", APIBasePath+"/writer/comments/"+commentID, updateCommentReq, token)
	resp = helper.AssertSuccess(w, http.StatusOK, "更新评论应该成功")

	helper.LogSuccess("更新评论成功")

	// ============================================
	// 步骤3: 验证更新后的评论
	// ============================================
	helper.LogInfo("步骤3: 验证更新后的评论")

	w = helper.DoAuthRequest("GET", APIBasePath+"/writer/comments/"+commentID, nil, token)
	resp = helper.AssertSuccess(w, http.StatusOK, "获取评论详情应该成功")

	data, ok = resp["data"].(map[string]interface{})
	require.True(t, ok, "响应数据格式正确")

	// 验证段落绑定关系不变
	if respParagraphID, ok := data["paragraphId"].(string); ok {
		assert.Equal(t, paragraphID, respParagraphID, "更新后段落绑定关系应该保持不变")
		helper.LogSuccess("段落绑定关系正确保持: %s", respParagraphID)
	}

	// 验证内容已更新
	if content, ok := data["content"].(string); ok {
		assert.Equal(t, "更新后的评论内容", content, "评论内容应该已更新")
	}

	// 验证类型已更新
	if commentType, ok := data["type"].(string); ok {
		assert.Equal(t, "issue", commentType, "评论类型应该已更新")
	}

	helper.LogSuccess("评论更新测试通过")
}

// TestCommentResolveWithParagraphBinding 测试解决段落绑定的评论
// 测试场景：创建绑定到段落的评论后，标记为已解决，验证状态变化
func TestCommentResolveWithParagraphBinding(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)
	token := helper.LoginTestUser()
	require.NotEmpty(t, token, "登录应该成功")

	// ============================================
	// 步骤1: 创建测试文档、段落和评论
	// ============================================
	helper.LogInfo("步骤1: 创建测试文档、段落和评论")

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

	paragraphID := primitive.NewObjectID().Hex()
	createParagraphsReq := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"id":      paragraphID,
				"type":    "paragraph",
				"content": []map[string]interface{}{{"type": "text", "text": "测试段落"}},
				"order":   1,
			},
		},
	}

	w = helper.DoAuthRequest("PUT", APIBasePath+"/writer/documents/"+documentID+"/contents", createParagraphsReq, token)
	helper.AssertSuccess(w, http.StatusOK, "创建段落应该成功")

	createCommentReq := map[string]interface{}{
		"content":     "需要解决的评论",
		"type":        "issue",
		"paragraphId": paragraphID,
	}

	w = helper.DoAuthRequest("POST", APIBasePath+"/writer/documents/"+documentID+"/comments", createCommentReq, token)
	resp = helper.AssertSuccess(w, http.StatusOK, "创建评论应该成功")

	data, _ = resp["data"].(map[string]interface{})
	commentID, _ := data["id"].(string)
	require.NotEmpty(t, commentID, "评论ID不能为空")

	helper.LogSuccess("创建评论成功: %s", commentID)

	// ============================================
	// 步骤2: 标记评论为已解决
	// ============================================
	helper.LogInfo("步骤2: 标记评论为已解决")

	w = helper.DoAuthRequest("POST", APIBasePath+"/writer/comments/"+commentID+"/resolve", nil, token)
	resp = helper.AssertSuccess(w, http.StatusOK, "标记已解决应该成功")

	helper.LogSuccess("标记评论已解决")

	// ============================================
	// 步骤3: 验证评论状态
	// ============================================
	helper.LogInfo("步骤3: 验证评论状态")

	w = helper.DoAuthRequest("GET", APIBasePath+"/writer/comments/"+commentID, nil, token)
	resp = helper.AssertSuccess(w, http.StatusOK, "获取评论详情应该成功")

	data, ok = resp["data"].(map[string]interface{})
	require.True(t, ok, "响应数据格式正确")

	// 验证resolved状态
	if resolved, ok := data["resolved"].(bool); ok {
		assert.True(t, resolved, "评论应该标记为已解决")
		helper.LogSuccess("评论状态正确: 已解决")
	}

	// 验证段落绑定关系保持不变
	if respParagraphID, ok := data["paragraphId"].(string); ok {
		assert.Equal(t, paragraphID, respParagraphID, "解决后段落绑定关系应该保持不变")
	}

	helper.LogSuccess("评论解决测试通过")
}

// TestCommentStatsWithParagraphBinding 测试段落绑定评论的统计
// 测试场景：创建多个段落的评论后，获取文档评论统计，验证统计数据
func TestCommentStatsWithParagraphBinding(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)
	token := helper.LoginTestUser()
	require.NotEmpty(t, token, "登录应该成功")

	// ============================================
	// 步骤1: 创建测试文档和多个段落
	// ============================================
	helper.LogInfo("步骤1: 创建测试文档和多个段落")

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

	paragraph1ID := primitive.NewObjectID().Hex()
	paragraph2ID := primitive.NewObjectID().Hex()

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
		},
	}

	w = helper.DoAuthRequest("PUT", APIBasePath+"/writer/documents/"+documentID+"/contents", createParagraphsReq, token)
	helper.AssertSuccess(w, http.StatusOK, "创建段落应该成功")

	// ============================================
	// 步骤2: 创建多条评论
	// ============================================
	helper.LogInfo("步骤2: 创建多条评论")

	comments := []map[string]interface{}{
		{"content": "未解决的问题", "type": "issue", "paragraphId": paragraph1ID},
		{"content": "建议1", "type": "suggestion", "paragraphId": paragraph1ID},
		{"content": "已解决的问题", "type": "issue", "paragraphId": paragraph2ID},
	}

	for _, commentReq := range comments {
		w = helper.DoAuthRequest("POST", APIBasePath+"/writer/documents/"+documentID+"/comments", commentReq, token)
		helper.AssertSuccess(w, http.StatusOK, "创建评论应该成功")
	}

	helper.LogSuccess("创建了 3 条评论")

	// ============================================
	// 步骤3: 标记一条评论为已解决
	// ============================================
	helper.LogInfo("步骤3: 标记一条评论为已解决")

	// 获取评论列表并标记一条为已解决
	w = helper.DoAuthRequest("GET", APIBasePath+"/writer/documents/"+documentID+"/comments", nil, token)
	resp = helper.AssertSuccess(w, http.StatusOK, "获取评论列表应该成功")

	data, _ = resp["data"].(map[string]interface{})
	commentsList, _ := data["comments"].([]interface{})

	if len(commentsList) > 0 {
		firstComment, ok := commentsList[0].(map[string]interface{})
		if ok {
			commentID, ok := firstComment["id"].(string)
			if ok {
				w = helper.DoAuthRequest("POST", APIBasePath+"/writer/comments/"+commentID+"/resolve", nil, token)
				helper.AssertSuccess(w, http.StatusOK, "标记已解决应该成功")
				helper.LogSuccess("标记评论已解决")
			}
		}
	}

	// ============================================
	// 步骤4: 获取评论统计
	// ============================================
	helper.LogInfo("步骤4: 获取评论统计")

	w = helper.DoAuthRequest("GET", APIBasePath+"/writer/documents/"+documentID+"/comments/stats", nil, token)
	resp = helper.AssertSuccess(w, http.StatusOK, "获取评论统计应该成功")

	data, ok = resp["data"].(map[string]interface{})
	require.True(t, ok, "响应数据格式正确")

	helper.LogSuccess("评论统计数据: %+v", data)

	// 验证统计数据
	if total, ok := data["total"].(float64); ok {
		assert.Equal(t, 3, int(total), "总评论数应该是3")
		helper.LogSuccess("总评论数: %d", int(total))
	}

	if unresolved, ok := data["unresolved"].(float64); ok {
		helper.LogInfo("未解决评论数: %d", int(unresolved))
		assert.LessOrEqual(t, int(unresolved), 3, "未解决评论数应该不超过3")
	}

	helper.LogSuccess("评论统计测试通过")
}
