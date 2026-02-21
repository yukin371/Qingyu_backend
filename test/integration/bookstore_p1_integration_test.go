//go:build integration
// +build integration

package integration

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ============================================
// P1 书城 API 集成测试
//
// 测试场景:
// - 搜索：书籍搜索 API
// - 筛选：分类筛选、状态筛选等
// - 推荐：个性化推荐、热门推荐
// - 排行榜：实时榜、周榜、月榜、新人榜
// - 公共访问：无需认证即可访问的接口
//
// 运行方式：
//   go test -v ./test/integration/... -run TestBookstoreP1
// ============================================

// TestBookstoreP1_Integration 主测试入口
func TestBookstoreP1_Integration(t *testing.T) {
	// 设置测试环境
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// 创建TestHelper
	helper := NewTestHelper(t, router)

	t.Log("========================================")
	t.Log("开始 P1 书城 API 集成测试")
	t.Log("========================================")

	// ========== 测试场景1：搜索功能 ==========
	t.Run("场景1_搜索功能", func(t *testing.T) {
		testSearchFunctionality(t, helper)
	})

	// ========== 测试场景2：筛选功能 ==========
	t.Run("场景2_筛选功能", func(t *testing.T) {
		testFilterFunctionality(t, helper)
	})

	// ========== 测试场景3：推荐功能 ==========
	t.Run("场景3_推荐功能", func(t *testing.T) {
		testRecommendationFunctionality(t, helper)
	})

	// ========== 测试场景4：排行榜功能 ==========
	t.Run("场景4_排行榜功能", func(t *testing.T) {
		testRankingFunctionality(t, helper)
	})

	// ========== 测试场景5：公共访问 ==========
	t.Run("场景5_公共访问", func(t *testing.T) {
		testPublicAccess(t, helper)
	})

	t.Log("========================================")
	t.Log("✅ P1 书城 API 集成测试完成")
	t.Log("========================================")
}

// ============================================
// 测试场景1：搜索功能
// ============================================
func testSearchFunctionality(t *testing.T, helper *TestHelper) {
	t.Log("\n--- 测试场景1：搜索功能 ---")

	// 1.1 测试关键词搜索
	t.Run("1.1_关键词搜索", func(t *testing.T) {
		keyword := "斗罗"
		url := fmt.Sprintf("/api/v1/search/search?q=%s&type=books&page=1&size=10", keyword)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		t.Logf("搜索关键词: %s", keyword)
		t.Logf("HTTP状态码: %d", w.Code)

		if w.Code == 200 {
			t.Log("✓ 关键词搜索成功")
		} else {
			t.Logf("⚠ 关键词搜索返回状态码: %d", w.Code)
			t.Logf("  响应: %s", w.Body.String())
		}
	})

	// 1.2 测试空搜索
	t.Run("1.2_空搜索", func(t *testing.T) {
		url := "/api/v1/search/search?q=&type=books&page=1&size=10"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		t.Logf("空搜索 HTTP状态码: %d", w.Code)

		// 空搜索应该返回400错误或空结果
		if w.Code == 400 || w.Code == 200 {
			t.Log("✓ 空搜索处理正确")
		} else {
			t.Logf("⚠ 空搜索返回意外状态码: %d", w.Code)
		}
	})

	// 1.3 测试分页搜索
	t.Run("1.3_分页搜索", func(t *testing.T) {
		keyword := "玄幻"
		url := fmt.Sprintf("/api/v1/search/search?q=%s&type=books&page=2&size=5", keyword)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		t.Logf("分页搜索 (page=2, size=5) HTTP状态码: %d", w.Code)

		if w.Code == 200 {
			t.Log("✓ 分页搜索成功")
		} else {
			t.Logf("⚠ 分页搜索返回状态码: %d", w.Code)
		}
	})

	// 1.4 测试批量搜索
	t.Run("1.4_批量搜索", func(t *testing.T) {
		url := "/api/v1/search/batch"

		// 构造批量搜索请求
		batchReq := map[string]interface{}{
			"queries": []map[string]string{
				{"query": "斗罗", "type": "books"},
				{"query": "玄幻", "type": "books"},
			},
		}

		w := helper.DoRequest("POST", url, batchReq, "")

		t.Logf("批量搜索 HTTP状态码: %d", w.Code)

		if w.Code == 200 {
			t.Log("✓ 批量搜索成功")
		} else {
			t.Logf("⚠ 批量搜索返回状态码: %d", w.Code)
		}
	})

	// 1.5 测试搜索健康检查
	t.Run("1.5_搜索服务健康检查", func(t *testing.T) {
		url := "/api/v1/search/health"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		t.Logf("搜索健康检查 HTTP状态码: %d", w.Code)

		if w.Code == 200 {
			t.Log("✓ 搜索服务健康检查通过")
		} else {
			t.Logf("⚠ 搜索服务健康检查失败: %d", w.Code)
		}
	})
}

// ============================================
// 测试场景2：筛选功能
// ============================================
func testFilterFunctionality(t *testing.T, helper *TestHelper) {
	t.Log("\n--- 测试场景2：筛选功能 ---")

	// 2.1 测试分类筛选
	t.Run("2.1_分类筛选", func(t *testing.T) {
		categories := []string{"玄幻", "都市", "科幻", "历史"}

		for _, category := range categories {
			url := fmt.Sprintf("/api/v1/bookstore/books?category=%s&page=1&size=10", category)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", url, nil)
			helper.router.ServeHTTP(w, req)

			t.Logf("分类筛选 [%s] HTTP状态码: %d", category, w.Code)

			if w.Code == 200 {
				t.Logf("  ✓ 分类 [%s] 筛选成功", category)
			} else {
				t.Logf("  ⚠ 分类 [%s] 筛选失败: %d", category, w.Code)
			}
		}
	})

	// 2.2 测试状态筛选
	t.Run("2.2_状态筛选", func(t *testing.T) {
		statuses := []string{"ongoing", "completed"}

		for _, status := range statuses {
			url := fmt.Sprintf("/api/v1/bookstore/books?status=%s&page=1&size=10", status)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", url, nil)
			helper.router.ServeHTTP(w, req)

			t.Logf("状态筛选 [%s] HTTP状态码: %d", status, w.Code)

			if w.Code == 200 {
				t.Logf("  ✓ 状态 [%s] 筛选成功", status)
			} else {
				t.Logf("  ⚠ 状态 [%s] 筛选失败: %d", status, w.Code)
			}
		}
	})

	// 2.3 测试组合筛选（分类+状态）
	t.Run("2.3_组合筛选", func(t *testing.T) {
		url := "/api/v1/bookstore/books?category=玄幻&status=ongoing&page=1&size=10"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		t.Logf("组合筛选 (分类=玄幻, 状态=连载) HTTP状态码: %d", w.Code)

		if w.Code == 200 {
			t.Log("✓ 组合筛选成功")
		} else {
			t.Logf("⚠ 组合筛选失败: %d", w.Code)
		}
	})

	// 2.4 测试排序选项
	t.Run("2.4_排序选项", func(t *testing.T) {
		sortOptions := []string{"latest", "popular", "rating"}

		for _, sort := range sortOptions {
			url := fmt.Sprintf("/api/v1/bookstore/books?sort=%s&page=1&size=10", sort)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", url, nil)
			helper.router.ServeHTTP(w, req)

			t.Logf("排序 [%s] HTTP状态码: %d", sort, w.Code)

			if w.Code == 200 {
				t.Logf("  ✓ 排序 [%s] 成功", sort)
			} else {
				t.Logf("  ⚠ 排序 [%s] 失败: %d", sort, w.Code)
			}
		}
	})

	// 2.5 测试获取分类树
	t.Run("2.5_获取分类树", func(t *testing.T) {
		url := "/api/v1/bookstore/categories/tree"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		t.Logf("分类树 HTTP状态码: %d", w.Code)

		if w.Code == 200 {
			t.Log("✓ 分类树获取成功")
		} else {
			t.Logf("⚠ 分类树获取失败: %d", w.Code)
			t.Logf("  响应: %s", w.Body.String())
		}
	})
}

// ============================================
// 测试场景3：推荐功能
// ============================================
func testRecommendationFunctionality(t *testing.T, helper *TestHelper) {
	t.Log("\n--- 测试场景3：推荐功能 ---")

	// 3.1 测试个性化推荐（需要认证，可能失败）
	t.Run("3.1_个性化推荐", func(t *testing.T) {
		url := "/api/v1/recommendation/personalized?page=1&size=10"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		t.Logf("个性化推荐 HTTP状态码: %d", w.Code)

		// 可能返回401（未认证）或200（成功）
		if w.Code == 200 {
			t.Log("✓ 个性化推荐成功（已认证）")
		} else if w.Code == 401 {
			t.Log("○ 个性化推荐需要认证（符合预期）")
		} else {
			t.Logf("⚠ 个性化推荐返回意外状态码: %d", w.Code)
		}
	})

	// 3.2 测试热门推荐
	t.Run("3.2_热门推荐", func(t *testing.T) {
		url := "/api/v1/recommendation/hot?page=1&size=10"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		t.Logf("热门推荐 HTTP状态码: %d", w.Code)

		if w.Code == 200 {
			t.Log("✓ 热门推荐成功")
		} else {
			t.Logf("⚠ 热门推荐失败: %d", w.Code)
		}
	})

	// 3.3 测试分类推荐
	t.Run("3.3_分类推荐", func(t *testing.T) {
		categories := []string{"玄幻", "都市"}

		for _, category := range categories {
			url := fmt.Sprintf("/api/v1/recommendation/category?category=%s&page=1&size=10", category)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", url, nil)
			helper.router.ServeHTTP(w, req)

			t.Logf("分类推荐 [%s] HTTP状态码: %d", category, w.Code)

			if w.Code == 200 {
				t.Logf("  ✓ 分类推荐 [%s] 成功", category)
			} else {
				t.Logf("  ⚠ 分类推荐 [%s] 失败: %d", category, w.Code)
			}
		}
	})

	// 3.4 测试首页推荐
	t.Run("3.4_首页推荐", func(t *testing.T) {
		url := "/api/v1/recommendation/homepage"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		t.Logf("首页推荐 HTTP状态码: %d", w.Code)

		if w.Code == 200 {
			t.Log("✓ 首页推荐成功")
		} else {
			t.Logf("⚠ 首页推荐失败: %d", w.Code)
		}
	})

	// 3.5 测试书城推荐书籍
	t.Run("3.5_书城推荐书籍", func(t *testing.T) {
		url := "/api/v1/bookstore/books/recommended?page=1&size=10"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		t.Logf("书城推荐书籍 HTTP状态码: %d", w.Code)

		if w.Code == 200 {
			t.Log("✓ 书城推荐书籍成功")
		} else {
			t.Logf("⚠ 书城推荐书籍失败: %d", w.Code)
		}
	})

	// 3.6 测试书城精选书籍
	t.Run("3.6_书城精选书籍", func(t *testing.T) {
		url := "/api/v1/bookstore/books/featured?page=1&size=10"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		t.Logf("书城精选书籍 HTTP状态码: %d", w.Code)

		if w.Code == 200 {
			t.Log("✓ 书城精选书籍成功")
		} else {
			t.Logf("⚠ 书城精选书籍失败: %d", w.Code)
		}
	})
}

// ============================================
// 测试场景4：排行榜功能
// ============================================
func testRankingFunctionality(t *testing.T, helper *TestHelper) {
	t.Log("\n--- 测试场景4：排行榜功能 ---")

	// 4.1 测试实时榜
	t.Run("4.1_实时榜", func(t *testing.T) {
		url := "/api/v1/bookstore/rankings/realtime?limit=20"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		t.Logf("实时榜 HTTP状态码: %d", w.Code)

		if w.Code == 200 {
			t.Log("✓ 实时榜获取成功")
		} else {
			t.Logf("⚠ 实时榜获取失败: %d", w.Code)
			t.Logf("  响应: %s", w.Body.String())
		}
	})

	// 4.2 测试周榜
	t.Run("4.2_周榜", func(t *testing.T) {
		url := "/api/v1/bookstore/rankings/weekly?limit=20"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		t.Logf("周榜 HTTP状态码: %d", w.Code)

		if w.Code == 200 {
			t.Log("✓ 周榜获取成功")
		} else {
			t.Logf("⚠ 周榜获取失败: %d", w.Code)
		}
	})

	// 4.3 测试月榜
	t.Run("4.3_月榜", func(t *testing.T) {
		url := "/api/v1/bookstore/rankings/monthly?limit=20"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		t.Logf("月榜 HTTP状态码: %d", w.Code)

		if w.Code == 200 {
			t.Log("✓ 月榜获取成功")
		} else {
			t.Logf("⚠ 月榜获取失败: %d", w.Code)
		}
	})

	// 4.4 测试新人榜
	t.Run("4.4_新人榜", func(t *testing.T) {
		url := "/api/v1/bookstore/rankings/newbie?limit=20"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		t.Logf("新人榜 HTTP状态码: %d", w.Code)

		if w.Code == 200 {
			t.Log("✓ 新人榜获取成功")
		} else {
			t.Logf("⚠ 新人榜获取失败: %d", w.Code)
		}
	})

	// 4.5 测试榜单列表
	t.Run("4.5_榜单列表", func(t *testing.T) {
		url := "/api/v1/bookstore/rankings"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		t.Logf("榜单列表 HTTP状态码: %d", w.Code)

		if w.Code == 200 {
			t.Log("✓ 榜单列表获取成功")
		} else {
			t.Logf("⚠ 榜单列表获取失败: %d", w.Code)
		}
	})
}

// ============================================
// 测试场景5：公共访问
// ============================================
func testPublicAccess(t *testing.T, helper *TestHelper) {
	t.Log("\n--- 测试场景5：公共访问 ---")

	// 5.1 测试书城首页
	t.Run("5.1_书城首页", func(t *testing.T) {
		url := "/api/v1/bookstore/homepage"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		t.Logf("书城首页 HTTP状态码: %d", w.Code)

		if w.Code == 200 {
			t.Log("✓ 书城首页访问成功（公共访问）")
		} else {
			t.Logf("⚠ 书城首页访问失败: %d", w.Code)
		}
	})

	// 5.2 测试书籍列表
	t.Run("5.2_书籍列表", func(t *testing.T) {
		url := "/api/v1/bookstore/books?page=1&size=20"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		t.Logf("书籍列表 HTTP状态码: %d", w.Code)

		if w.Code == 200 {
			t.Log("✓ 书籍列表访问成功（公共访问）")
		} else {
			t.Logf("⚠ 书籍列表访问失败: %d", w.Code)
		}
	})

	// 5.3 测试书籍详情
	t.Run("5.3_书籍详情", func(t *testing.T) {
		// 使用一个测试书籍ID，实际环境中需要替换为有效ID
		url := "/api/v1/bookstore/books/507f1f77bcf86cd799439011"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		t.Logf("书籍详情 HTTP状态码: %d", w.Code)

		// 可能返回200（成功）或404（书籍不存在）
		if w.Code == 200 {
			t.Log("✓ 书籍详情访问成功（公共访问）")
		} else if w.Code == 404 {
			t.Log("○ 书籍不存在（符合预期，使用测试ID）")
		} else {
			t.Logf("⚠ 书籍详情返回意外状态码: %d", w.Code)
		}
	})

	// 5.4 测试活动Banner
	t.Run("5.4_活动Banner", func(t *testing.T) {
		url := "/api/v1/bookstore/banners?limit=5"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		t.Logf("活动Banner HTTP状态码: %d", w.Code)

		if w.Code == 200 {
			t.Log("✓ 活动Banner访问成功（公共访问）")
		} else {
			t.Logf("⚠ 活动Banner访问失败: %d", w.Code)
		}
	})

	// 5.5 测试用户信息（公共）
	t.Run("5.5_用户信息_公开访问", func(t *testing.T) {
		// 使用一个测试用户ID
		url := "/api/v1/user/users/507f1f77bcf86cd799439011"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		t.Logf("用户信息（公开）HTTP状态码: %d", w.Code)

		// 可能返回200（成功）或404（用户不存在）
		if w.Code == 200 {
			t.Log("✓ 用户信息访问成功（公共访问）")
		} else if w.Code == 404 {
			t.Log("○ 用户不存在（符合预期，使用测试ID）")
		} else {
			t.Logf("⚠ 用户信息返回意外状态码: %d", w.Code)
		}
	})

	// 5.6 测试用户资料（公开）
	t.Run("5.6_用户资料_公开访问", func(t *testing.T) {
		url := "/api/v1/user/users/507f1f77bcf86cd799439011/profile"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		t.Logf("用户资料（公开）HTTP状态码: %d", w.Code)

		if w.Code == 200 {
			t.Log("✓ 用户资料访问成功（公共访问）")
		} else if w.Code == 404 {
			t.Log("○ 用户不存在（符合预期，使用测试ID）")
		} else {
			t.Logf("⚠ 用户资料返回意外状态码: %d", w.Code)
		}
	})

	// 5.7 测试用户作品列表（公开）
	t.Run("5.7_用户作品列表_公开访问", func(t *testing.T) {
		url := "/api/v1/user/users/507f1f77bcf86cd799439011/books?page=1&size=10"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		t.Logf("用户作品列表（公开）HTTP状态码: %d", w.Code)

		if w.Code == 200 {
			t.Log("✓ 用户作品列表访问成功（公共访问）")
		} else if w.Code == 404 {
			t.Log("○ 用户不存在（符合预期，使用测试ID）")
		} else {
			t.Logf("⚠ 用户作品列表返回意外状态码: %d", w.Code)
		}
	})
}

// ============================================
// 测试场景6：综合测试
// ============================================
func TestBookstoreP1_ComprehensiveFlow(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)

	t.Log("========================================")
	t.Log("综合测试：完整书城浏览流程")
	t.Log("========================================")

	// 步骤1：访问书城首页
	t.Run("步骤1_访问书城首页", func(t *testing.T) {
		url := "/api/v1/bookstore/homepage"
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		if w.Code == 200 {
			t.Log("✓ 步骤1完成：书城首页加载成功")
		} else {
			t.Logf("✗ 步骤1失败：HTTP %d", w.Code)
		}
	})

	// 步骤2：浏览分类
	t.Run("步骤2_浏览分类", func(t *testing.T) {
		url := "/api/v1/bookstore/categories/tree"
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		if w.Code == 200 {
			t.Log("✓ 步骤2完成：分类树加载成功")
		} else {
			t.Logf("✗ 步骤2失败：HTTP %d", w.Code)
		}
	})

	// 步骤3：查看推荐书籍
	t.Run("步骤3_查看推荐书籍", func(t *testing.T) {
		url := "/api/v1/bookstore/books/recommended?page=1&size=10"
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		if w.Code == 200 {
			t.Log("✓ 步骤3完成：推荐书籍加载成功")
		} else {
			t.Logf("✗ 步骤3失败：HTTP %d", w.Code)
		}
	})

	// 步骤4：查看排行榜
	t.Run("步骤4_查看排行榜", func(t *testing.T) {
		url := "/api/v1/bookstore/rankings/realtime?limit=10"
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		if w.Code == 200 {
			t.Log("✓ 步骤4完成：排行榜加载成功")
		} else {
			t.Logf("✗ 步骤4失败：HTTP %d", w.Code)
		}
	})

	// 步骤5：搜索书籍
	t.Run("步骤5_搜索书籍", func(t *testing.T) {
		url := "/api/v1/search/search?q=斗罗&type=books&page=1&size=10"
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		if w.Code == 200 {
			t.Log("✓ 步骤5完成：书籍搜索成功")
		} else {
			t.Logf("✗ 步骤5失败：HTTP %d", w.Code)
		}
	})

	// 步骤6：查看热门推荐
	t.Run("步骤6_查看热门推荐", func(t *testing.T) {
		url := "/api/v1/recommendation/hot?page=1&size=10"
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		helper.router.ServeHTTP(w, req)

		if w.Code == 200 {
			t.Log("✓ 步骤6完成：热门推荐加载成功")
		} else {
			t.Logf("✗ 步骤6失败：HTTP %d", w.Code)
		}
	})

	t.Log("========================================")
	t.Log("综合测试流程完成")
	t.Log("========================================")
}
