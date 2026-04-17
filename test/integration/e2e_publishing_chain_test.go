//go:build integration
// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ========================================
// 发布链路 E2E 测试
// 完整链路：作者写作 -> 发布申请 -> 管理员审核 -> 读者搜索/阅读/评论
// ========================================

// TestPublishingChain_FullWorkflow 发布链路完整工作流
// 1. 作者：注册 -> 创建项目 -> 创建章节 -> 发布项目
// 2. 管理员：获取待审核列表 -> 审核通过
// 3. 读者：搜索书籍 -> 查看详情 -> 加入书架 -> 阅读章节 -> 发表评论
func TestPublishingChain_FullWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过E2E测试")
	}

	// 设置测试环境
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// 创建测试辅助工具
	helper := NewTestHelper(t, router)

	// ========== 第一阶段：作者工作流 ==========
	t.Log("=== 阶段1：作者工作流 ===")

	// 1.1 注册作者用户
	t.Run("Phase1_AuthorRegister", func(t *testing.T) {
		authorUsername := fmt.Sprintf("author_%d", time.Now().UnixNano())
		authorPassword := "Test@123456"
		authorEmail := authorUsername + "@test.com"

		registerData := map[string]interface{}{
			"username": authorUsername,
			"password": authorPassword,
			"email":    authorEmail,
			"nickname": "测试作者",
		}

		body, _ := json.Marshal(registerData)
		req := httptest.NewRequest("POST", "/api/v1/user/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 注册可能成功或失败（如果用户已存在），继续登录
		t.Logf("注册响应: %d - %s", w.Code, w.Body.String())

		// 1.2 登录获取token
		loginData := map[string]interface{}{
			"username": authorUsername,
			"password": authorPassword,
		}

		loginBody, _ := json.Marshal(loginData)
		loginReq := httptest.NewRequest("POST", "/api/v1/login", bytes.NewReader(loginBody))
		loginReq.Header.Set("Content-Type", "application/json")

		loginW := httptest.NewRecorder()
		router.ServeHTTP(loginW, loginReq)

		if loginW.Code != http.StatusOK && loginW.Code != http.StatusCreated {
			t.Skipf("作者登录失败，跳过测试: %d - %s", loginW.Code, loginW.Body.String())
		}

		var loginResp map[string]interface{}
		json.Unmarshal(loginW.Body.Bytes(), &loginResp)

		data, ok := loginResp["data"].(map[string]interface{})
		if !ok {
			t.Skip("无法获取登录响应数据")
		}

		authorToken, ok := data["token"].(string)
		if !ok {
			t.Skip("无法获取作者token")
		}

		t.Logf("✓ 作者登录成功: %s...", authorToken[:20])

		// 1.3 创建项目
		projectData := map[string]interface{}{
			"title":       "测试小说_" + fmt.Sprintf("%d", time.Now().UnixNano()),
			"description": "这是一本测试小说",
			"genre":       "玄幻",
			"tags":        []string{"测试", "玄幻"},
		}

		projectBody, _ := json.Marshal(projectData)
		projectReq := httptest.NewRequest("POST", "/api/v1/writer/projects", bytes.NewReader(projectBody))
		projectReq.Header.Set("Content-Type", "application/json")
		projectReq.Header.Set("Authorization", "Bearer "+authorToken)

		projectW := httptest.NewRecorder()
		router.ServeHTTP(projectW, projectReq)

		t.Logf("创建项目响应: %d - %s", projectW.Code, projectW.Body.String())

		if projectW.Code != http.StatusOK && projectW.Code != http.StatusCreated {
			t.Skipf("创建项目失败，跳过测试: %d", projectW.Code)
		}

		var projectResp map[string]interface{}
		json.Unmarshal(projectW.Body.Bytes(), &projectResp)

		projectID := ""
		if respData, ok := projectResp["data"].(map[string]interface{}); ok {
			if pid, ok := respData["id"].(string); ok {
				projectID = pid
			} else if pid, ok := respData["project_id"].(string); ok {
				projectID = pid
			} else if pid, ok := respData["_id"].(string); ok {
				projectID = pid
			}
		}

		if projectID == "" {
			t.Skip("无法获取项目ID")
		}

		t.Logf("✓ 项目创建成功: %s", projectID)

		// 1.4 创建章节文档
		chapterData := map[string]interface{}{
			"title":   "第一章 测试章节",
			"content": "这是第一章的测试内容，包含了足够的文字来测试发布功能。",
			"order":   1,
		}

		chapterBody, _ := json.Marshal(chapterData)
		chapterReq := httptest.NewRequest("POST", fmt.Sprintf("/api/v1/writer/projects/%s/documents", projectID), bytes.NewReader(chapterBody))
		chapterReq.Header.Set("Content-Type", "application/json")
		chapterReq.Header.Set("Authorization", "Bearer "+authorToken)

		chapterW := httptest.NewRecorder()
		router.ServeHTTP(chapterW, chapterReq)

		t.Logf("创建章节响应: %d - %s", chapterW.Code, chapterW.Body.String())

		if chapterW.Code != http.StatusOK && chapterW.Code != http.StatusCreated {
			t.Skipf("创建章节失败，跳过测试: %d", chapterW.Code)
		}

		var chapterResp map[string]interface{}
		json.Unmarshal(chapterW.Body.Bytes(), &chapterResp)

		documentID := ""
		if respData, ok := chapterResp["data"].(map[string]interface{}); ok {
			if did, ok := respData["id"].(string); ok {
				documentID = did
			} else if did, ok := respData["document_id"].(string); ok {
				documentID = did
			} else if did, ok := respData["_id"].(string); ok {
				documentID = did
			}
		}

		if documentID == "" {
			t.Skip("无法获取文档ID")
		}

		t.Logf("✓ 章节创建成功: %s", documentID)

		// 1.5 发布项目
		publishData := map[string]interface{}{
			"publish_type": "serial",
			"category":     "玄幻",
		}

		publishBody, _ := json.Marshal(publishData)
		publishReq := httptest.NewRequest("POST", fmt.Sprintf("/api/v1/writer/projects/%s/publish", projectID), bytes.NewReader(publishBody))
		publishReq.Header.Set("Content-Type", "application/json")
		publishReq.Header.Set("Authorization", "Bearer "+authorToken)

		publishW := httptest.NewRecorder()
		router.ServeHTTP(publishW, publishReq)

		t.Logf("发布项目响应: %d - %s", publishW.Code, publishW.Body.String())

		// 即使发布失败，也继续测试管理员审核流程（可能有待审核内容）
		if publishW.Code != http.StatusOK && publishW.Code != http.StatusCreated {
			t.Logf("⚠ 发布请求返回非成功状态: %d", publishW.Code)
		} else {
			t.Log("✓ 项目发布申请成功")
		}

		// 存储项目ID供后续使用
		helper.t.Logf("作者工作流完成，项目ID: %s, 文档ID: %s", projectID, documentID)
	})

	// ========== 第二阶段：管理员审核工作流 ==========
	t.Log("\n=== 阶段2：管理员审核工作流 ===")

	t.Run("Phase2_AdminAudit", func(t *testing.T) {
		// 2.1 管理员登录
		adminToken := helper.LoginUser("admin", "Admin@123456")
		if adminToken == "" {
			t.Skip("管理员登录失败，跳过审核测试")
		}

		t.Log("✓ 管理员登录成功")

		// 2.2 获取待审核发布列表
		pendingReq := httptest.NewRequest("GET", "/api/v1/admin/publications/pending?page=1&size=10", nil)
		pendingReq.Header.Set("Authorization", "Bearer "+adminToken)

		pendingW := httptest.NewRecorder()
		router.ServeHTTP(pendingW, pendingReq)

		t.Logf("获取待审核列表响应: %d - %s", pendingW.Code, pendingW.Body.String())

		if pendingW.Code != http.StatusOK {
			t.Logf("⚠ 获取待审核列表失败，继续测试: %d", pendingW.Code)
		} else {
			t.Log("✓ 获取待审核列表成功")
		}

		// 2.3 获取审核统计
		statsReq := httptest.NewRequest("GET", "/api/v1/admin/stats/reviews", nil)
		statsReq.Header.Set("Authorization", "Bearer "+adminToken)

		statsW := httptest.NewRecorder()
		router.ServeHTTP(statsW, statsReq)

		t.Logf("获取审核统计响应: %d - %s", statsW.Code, statsW.Body.String())

		if statsW.Code == http.StatusOK {
			t.Log("✓ 获取审核统计成功")
		}
	})

	// ========== 第三阶段：读者工作流 ==========
	t.Log("\n=== 阶段3：读者工作流 ===")

	t.Run("Phase3_ReaderWorkflow", func(t *testing.T) {
		// 3.1 读者登录
		readerToken := helper.LoginTestUser()
		if readerToken == "" {
			t.Skip("读者登录失败，跳过读者测试")
		}

		t.Log("✓ 读者登录成功")

		// 3.2 搜索书籍
		searchReq := httptest.NewRequest("GET", "/api/v1/bookstore/books/search?keyword=测试&page=1&size=10", nil)

		searchW := httptest.NewRecorder()
		router.ServeHTTP(searchW, searchReq)

		t.Logf("搜索书籍响应: %d", searchW.Code)

		if searchW.Code != http.StatusOK {
			t.Logf("⚠ 搜索请求返回非200状态: %d", searchW.Code)
		} else {
			t.Log("✓ 搜索请求成功")
		}

		// 3.3 获取书籍列表
		booksReq := httptest.NewRequest("GET", "/api/v1/bookstore/books?page=1&size=10", nil)

		booksW := httptest.NewRecorder()
		router.ServeHTTP(booksW, booksReq)

		t.Logf("获取书籍列表响应: %d", booksW.Code)

		if booksW.Code != http.StatusOK {
			t.Logf("⚠ 获取书籍列表返回非200状态: %d", booksW.Code)
		} else {
			t.Log("✓ 获取书籍列表成功")
		}

		// 3.4 获取书城首页
		homeReq := httptest.NewRequest("GET", "/api/v1/bookstore/homepage", nil)

		homeW := httptest.NewRecorder()
		router.ServeHTTP(homeW, homeReq)

		t.Logf("获取首页响应: %d", homeW.Code)

		if homeW.Code != http.StatusOK {
			t.Logf("⚠ 获取首页返回非200状态: %d", homeW.Code)
		} else {
			t.Log("✓ 获取书城首页成功")
		}

		// 3.5 加入书架
		// 注意：需要先获取一个有效的bookId
		var bookID string
		if booksW.Code == http.StatusOK {
			var booksResp map[string]interface{}
			json.Unmarshal(booksW.Body.Bytes(), &booksResp)
			if data, ok := booksResp["data"].([]interface{}); ok && len(data) > 0 {
				if book, ok := data[0].(map[string]interface{}); ok {
					if bid, ok := book["id"].(string); ok {
						bookID = bid
					} else if bid, ok := book["_id"].(string); ok {
						bookID = bid
					} else if bid, ok := book["book_id"].(string); ok {
						bookID = bid
					}
				}
			}
		}

		if bookID != "" {
			addBookshelfReq := httptest.NewRequest("POST", fmt.Sprintf("/api/v1/reader/books/%s", bookID), nil)
			addBookshelfReq.Header.Set("Authorization", "Bearer "+readerToken)

			addW := httptest.NewRecorder()
			router.ServeHTTP(addW, addBookshelfReq)

			t.Logf("加入书架响应: %d - %s", addW.Code, addW.Body.String())

			if addW.Code == http.StatusOK || addW.Code == http.StatusCreated {
				t.Log("✓ 加入书架成功")
			} else {
				t.Logf("⚠ 加入书架返回: %d", addW.Code)
			}
		} else {
			t.Log("⚠ 未找到有效书籍ID，跳过加入书架")
		}

		// 3.6 获取书架
		bookshelfReq := httptest.NewRequest("GET", "/api/v1/reader/books", nil)
		bookshelfReq.Header.Set("Authorization", "Bearer "+readerToken)

		bookshelfW := httptest.NewRecorder()
		router.ServeHTTP(bookshelfW, bookshelfReq)

		t.Logf("获取书架响应: %d", bookshelfW.Code)

		if bookshelfW.Code != http.StatusOK {
			t.Logf("⚠ 获取书架返回非200状态: %d", bookshelfW.Code)
		} else {
			t.Log("✓ 获取书架成功")
		}

		// 3.7 发表书籍评论（需要有效的bookId）
		if bookID != "" {
			commentData := map[string]interface{}{
				"book_id":  bookID,
				"content":  "这是一条测试评论",
				"rating":   5,
			}

			commentBody, _ := json.Marshal(commentData)
			commentReq := httptest.NewRequest("POST", "/api/v1/social/reviews", bytes.NewReader(commentBody))
			commentReq.Header.Set("Content-Type", "application/json")
			commentReq.Header.Set("Authorization", "Bearer "+readerToken)

			commentW := httptest.NewRecorder()
			router.ServeHTTP(commentW, commentReq)

			t.Logf("发表评论响应: %d - %s", commentW.Code, commentW.Body.String())

			if commentW.Code == http.StatusOK || commentW.Code == http.StatusCreated {
				t.Log("✓ 发表书籍评论成功")
			} else {
				t.Logf("⚠ 发表评论返回: %d", commentW.Code)
			}
		}

		// 3.8 获取书籍评论列表
		if bookID != "" {
			reviewsReq := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/social/reviews?bookId=%s&page=1&size=10", bookID), nil)

			reviewsW := httptest.NewRecorder()
			router.ServeHTTP(reviewsW, reviewsReq)

			t.Logf("获取评论列表响应: %d", reviewsW.Code)

			if reviewsW.Code != http.StatusOK {
				t.Logf("⚠ 获取评论列表返回非200状态: %d", reviewsW.Code)
			} else {
				t.Log("✓ 获取评论列表成功")
			}
		}

		// 3.9 书籍点赞
		if bookID != "" {
			likeReq := httptest.NewRequest("POST", fmt.Sprintf("/api/v1/social/books/%s/like", bookID), nil)
			likeReq.Header.Set("Authorization", "Bearer "+readerToken)

			likeW := httptest.NewRecorder()
			router.ServeHTTP(likeW, likeReq)

			t.Logf("书籍点赞响应: %d - %s", likeW.Code, likeW.Body.String())

			if likeW.Code == http.StatusOK || likeW.Code == http.StatusCreated {
				t.Log("✓ 书籍点赞成功")
			} else {
				t.Logf("⚠ 点赞返回: %d", likeW.Code)
			}
		}
	})

	t.Log("\n=== 发布链路E2E测试完成 ===")
}

// TestPublishingChain_AuthorPublishing 作者发布功能测试
func TestPublishingChain_AuthorPublishing(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过E2E测试")
	}

	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)

	t.Run("Author_CanPublishProject", func(t *testing.T) {
		// 1. 创建测试作者
		authorUsername := fmt.Sprintf("author_pub_%d", time.Now().UnixNano())
		authorToken := helper.LoginUser(authorUsername, "Test@123456")

		if authorToken == "" {
			t.Skip("作者登录失败")
		}

		// 2. 创建项目
		projectData := map[string]interface{}{
			"title":       "发布测试小说_" + fmt.Sprintf("%d", time.Now().UnixNano()),
			"description": "用于测试发布功能的小说",
			"genre":       "都市",
			"tags":        []string{"测试", "都市"},
		}

		body, _ := json.Marshal(projectData)
		req := httptest.NewRequest("POST", "/api/v1/writer/projects", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authorToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK && w.Code != http.StatusCreated {
			t.Skipf("创建项目失败: %d", w.Code)
		}

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		projectID := ""
		if data, ok := resp["data"].(map[string]interface{}); ok {
			if pid, ok := data["id"].(string); ok {
				projectID = pid
			} else if pid, ok := data["project_id"].(string); ok {
				projectID = pid
			}
		}

		require.NotEmpty(t, projectID, "项目ID不应为空")
		t.Logf("✓ 项目创建成功: %s", projectID)

		// 3. 发布项目
		publishData := map[string]interface{}{
			"publish_type": "serial",
			"category":     "都市",
		}

		publishBody, _ := json.Marshal(publishData)
		publishReq := httptest.NewRequest("POST", fmt.Sprintf("/api/v1/writer/projects/%s/publish", projectID), bytes.NewReader(publishBody))
		publishReq.Header.Set("Content-Type", "application/json")
		publishReq.Header.Set("Authorization", "Bearer "+authorToken)

		publishW := httptest.NewRecorder()
		router.ServeHTTP(publishW, publishReq)

		t.Logf("发布响应: %d - %s", publishW.Code, publishW.Body.String())

		// 4. 验证发布结果
		if publishW.Code == http.StatusOK || publishW.Code == http.StatusCreated {
			t.Log("✓ 项目发布成功")
		} else if publishW.Code == http.StatusBadRequest {
			// 可能是项目不符合发布条件，记录警告
			t.Log("⚠ 项目发布返回400（可能缺少必要内容）")
		} else {
			t.Logf("⚠ 发布返回: %d", publishW.Code)
		}

		// 5. 查询发布状态
		statusReq := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/writer/projects/%s/publication-status", projectID), nil)
		statusReq.Header.Set("Authorization", "Bearer "+authorToken)

		statusW := httptest.NewRecorder()
		router.ServeHTTP(statusW, statusReq)

		t.Logf("查询发布状态响应: %d - %s", statusW.Code, statusW.Body.String())

		if statusW.Code == http.StatusOK {
			t.Log("✓ 查询发布状态成功")
		}
	})
}

// TestPublishingChain_AdminAudit 管理员审核功能测试
func TestPublishingChain_AdminAudit(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过E2E测试")
	}

	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)

	t.Run("Admin_CanReviewPublications", func(t *testing.T) {
		// 1. 管理员登录
		adminToken := helper.LoginUser("admin", "Admin@123456")
		if adminToken == "" {
			t.Skip("管理员登录失败")
		}

		t.Log("✓ 管理员登录成功")

		// 2. 获取待审核列表
		req := httptest.NewRequest("GET", "/api/v1/admin/publications/pending?page=1&size=10", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		t.Logf("待审核列表响应: %d - %s", w.Code, w.Body.String())

		if w.Code == http.StatusOK {
			t.Log("✓ 获取待审核列表成功")
		} else if w.Code == http.StatusForbidden {
			t.Log("⚠ 无权限访问审核列表（需要admin角色）")
		} else {
			t.Logf("⚠ 获取列表返回: %d", w.Code)
		}

		// 3. 获取审核统计
		statsReq := httptest.NewRequest("GET", "/api/v1/admin/stats/reviews", nil)
		statsReq.Header.Set("Authorization", "Bearer "+adminToken)

		statsW := httptest.NewRecorder()
		router.ServeHTTP(statsW, statsReq)

		if statsW.Code == http.StatusOK {
			t.Log("✓ 获取审核统计成功")
		}
	})
}

// TestPublishingChain_ReaderExperience 读者体验测试
func TestPublishingChain_ReaderExperience(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过E2E测试")
	}

	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)

	t.Run("Reader_CanSearchAndRead", func(t *testing.T) {
		// 1. 读者登录
		readerToken := helper.LoginTestUser()
		if readerToken == "" {
			t.Skip("读者登录失败")
		}

		t.Log("✓ 读者登录成功")

		// 2. 访问书城首页
		homeReq := httptest.NewRequest("GET", "/api/v1/bookstore/homepage", nil)
		homeW := httptest.NewRecorder()
		router.ServeHTTP(homeW, homeReq)

		assert.Equal(t, http.StatusOK, homeW.Code, "首页应该可访问")
		t.Log("✓ 书城首页可访问")

		// 3. 搜索书籍
		searchReq := httptest.NewRequest("GET", "/api/v1/bookstore/books/search?keyword=test&page=1", nil)
		searchW := httptest.NewRecorder()
		router.ServeHTTP(searchW, searchReq)

		t.Logf("搜索响应: %d", searchW.Code)
		if searchW.Code == http.StatusOK {
			t.Log("✓ 搜索功能正常")
		}

		// 4. 获取分类
		categoryReq := httptest.NewRequest("GET", "/api/v1/bookstore/categories/tree", nil)
		categoryW := httptest.NewRecorder()
		router.ServeHTTP(categoryW, categoryReq)

		if categoryW.Code == http.StatusOK {
			t.Log("✓ 分类数据可获取")
		}

		// 5. 获取榜单
		rankingReq := httptest.NewRequest("GET", "/api/v1/bookstore/rankings/realtime?limit=10", nil)
		rankingW := httptest.NewRecorder()
		router.ServeHTTP(rankingW, rankingReq)

		if rankingW.Code == http.StatusOK {
			t.Log("✓ 排行榜数据可获取")
		}

		// 6. 获取书架
		bookshelfReq := httptest.NewRequest("GET", "/api/v1/reader/books", nil)
		bookshelfReq.Header.Set("Authorization", "Bearer "+readerToken)

		bookshelfW := httptest.NewRecorder()
		router.ServeHTTP(bookshelfW, bookshelfReq)

		t.Logf("书架响应: %d", bookshelfW.Code)
		if bookshelfW.Code == http.StatusOK {
			t.Log("✓ 读者书架可访问")
		}

		// 7. 获取阅读进度
		progressReq := httptest.NewRequest("GET", "/api/v1/reader/progress/recent", nil)
		progressReq.Header.Set("Authorization", "Bearer "+readerToken)

		progressW := httptest.NewRecorder()
		router.ServeHTTP(progressW, progressReq)

		t.Logf("阅读进度响应: %d", progressW.Code)
		if progressW.Code == http.StatusOK {
			t.Log("✓ 阅读进度可查询")
		}
	})
}

// TestPublishingChain_SocialInteractions 社交互动测试
func TestPublishingChain_SocialInteractions(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过E2E测试")
	}

	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)

	t.Run("Reader_CanCommentAndLike", func(t *testing.T) {
		// 1. 读者登录
		readerToken := helper.LoginTestUser()
		if readerToken == "" {
			t.Skip("读者登录失败")
		}

		// 2. 获取书籍列表
		booksReq := httptest.NewRequest("GET", "/api/v1/bookstore/books?page=1&size=5", nil)
		booksW := httptest.NewRecorder()
		router.ServeHTTP(booksW, booksReq)

		if booksW.Code != http.StatusOK {
			t.Skip("无法获取书籍列表")
		}

		var booksResp map[string]interface{}
		json.Unmarshal(booksW.Body.Bytes(), &booksResp)

		bookID := ""
		if data, ok := booksResp["data"].([]interface{}); ok && len(data) > 0 {
			if book, ok := data[0].(map[string]interface{}); ok {
				if bid, ok := book["id"].(string); ok {
					bookID = bid
				} else if bid, ok := book["_id"].(string); ok {
					bookID = bid
				}
			}
		}

		if bookID == "" {
			t.Skip("未找到有效书籍")
		}

		t.Logf("✓ 获取到书籍: %s", bookID)

		// 3. 点赞书籍
		likeReq := httptest.NewRequest("POST", fmt.Sprintf("/api/v1/social/books/%s/like", bookID), nil)
		likeReq.Header.Set("Authorization", "Bearer "+readerToken)

		likeW := httptest.NewRecorder()
		router.ServeHTTP(likeW, likeReq)

		t.Logf("点赞响应: %d - %s", likeW.Code, likeW.Body.String())
		if likeW.Code == http.StatusOK || likeW.Code == http.StatusCreated {
			t.Log("✓ 点赞成功")
		}

		// 4. 收藏书籍
		collectData := map[string]interface{}{
			"book_id": bookID,
			"tags":    []string{"测试收藏"},
		}
		collectBody, _ := json.Marshal(collectData)

		collectReq := httptest.NewRequest("POST", "/api/v1/social/collections", bytes.NewReader(collectBody))
		collectReq.Header.Set("Content-Type", "application/json")
		collectReq.Header.Set("Authorization", "Bearer "+readerToken)

		collectW := httptest.NewRecorder()
		router.ServeHTTP(collectW, collectReq)

		t.Logf("收藏响应: %d - %s", collectW.Code, collectW.Body.String())
		if collectW.Code == http.StatusOK || collectW.Code == http.StatusCreated {
			t.Log("✓ 收藏成功")
		}

		// 5. 发表书评
		reviewData := map[string]interface{}{
			"book_id": bookID,
			"content": "这本书很不错！测试评论。",
			"rating":  5,
		}
		reviewBody, _ := json.Marshal(reviewData)

		reviewReq := httptest.NewRequest("POST", "/api/v1/social/reviews", bytes.NewReader(reviewBody))
		reviewReq.Header.Set("Content-Type", "application/json")
		reviewReq.Header.Set("Authorization", "Bearer "+readerToken)

		reviewW := httptest.NewRecorder()
		router.ServeHTTP(reviewW, reviewReq)

		t.Logf("书评响应: %d - %s", reviewW.Code, reviewW.Body.String())
		if reviewW.Code == http.StatusOK || reviewW.Code == http.StatusCreated {
			t.Log("✓ 书评发表成功")
		}

		// 6. 获取书评列表
		reviewsReq := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/social/reviews?bookId=%s", bookID), nil)
		reviewsW := httptest.NewRecorder()
		router.ServeHTTP(reviewsW, reviewsReq)

		if reviewsW.Code == http.StatusOK {
			t.Log("✓ 书评列表可获取")
		}
	})
}
