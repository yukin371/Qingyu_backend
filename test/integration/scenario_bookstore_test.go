package integration

import (
	"fmt"
	"testing"

	"Qingyu_backend/global"

	"go.mongodb.org/mongo-driver/bson"
)

// API路径常量
const (
	BookstoreHomepagePath        = "/api/v1/bookstore/homepage"
	BookstoreCategoriesTreePath  = "/api/v1/bookstore/categories/tree"
	BookstoreRecommendedPath     = "/api/v1/bookstore/books/recommended"
	BookstoreFeaturedPath        = "/api/v1/bookstore/books/featured"
	BookstoreBannersPath         = "/api/v1/bookstore/banners"
	BookstoreRankingRealtimePath = "/api/v1/bookstore/rankings/realtime"
	BookstoreRankingWeeklyPath   = "/api/v1/bookstore/rankings/weekly"
	BookstoreRankingMonthlyPath  = "/api/v1/bookstore/rankings/monthly"
	BookstoreRankingNewbiePath   = "/api/v1/bookstore/rankings/newbie"
)

// 书城测试 - 完整书城浏览流程
func TestBookstoreScenario(t *testing.T) {
	// 设置测试环境
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// 创建TestHelper
	helper := NewTestHelper(t, router)

	t.Run("1.书城首页_获取首页数据", func(t *testing.T) {
		w := helper.DoRequest("GET", BookstoreHomepagePath, nil, "")
		response := helper.AssertSuccess(w, 200, "获取首页数据失败")

		helper.LogSuccess("首页数据获取成功")
		t.Logf("  响应: %v", response["message"])
	})

	t.Run("2.书城首页_获取分类树", func(t *testing.T) {
		w := helper.DoRequest("GET", BookstoreCategoriesTreePath, nil, "")

		if w.Code != 200 {
			t.Logf("⚠ 获取分类树失败，HTTP状态: %d", w.Code)
			t.Logf("  响应体: %s", w.Body.String())
		}

		response := helper.AssertSuccess(w, 200, "获取分类树失败")

		if response["data"] == nil {
			t.Logf("⚠ 分类树数据为空（可能还没有导入分类数据）")
		} else {
			helper.LogSuccess("分类树获取成功")
		}
	})

	t.Run("3.书城首页_获取推荐书籍", func(t *testing.T) {
		url := fmt.Sprintf("%s?page=1&size=10", BookstoreRecommendedPath)
		w := helper.DoRequest("GET", url, nil, "")

		if w.Code != 200 {
			t.Logf("⚠ 获取推荐书籍失败，HTTP状态: %d", w.Code)
			t.Logf("  响应体: %s", w.Body.String())
		}

		response := helper.AssertSuccess(w, 200, "获取推荐书籍失败")

		if data, ok := response["data"].([]interface{}); ok && len(data) > 0 {
			helper.LogSuccess(fmt.Sprintf("推荐书籍获取成功，数量: %d", len(data)))
		} else {
			t.Logf("○ 推荐书籍为空（可能还没有导入数据）")
		}
	})

	t.Run("4.书城首页_获取精选书籍", func(t *testing.T) {
		url := fmt.Sprintf("%s?page=1&size=10", BookstoreFeaturedPath)
		w := helper.DoRequest("GET", url, nil, "")

		if w.Code != 200 {
			t.Logf("⚠ 获取精选书籍失败，HTTP状态: %d", w.Code)
			t.Logf("  响应体: %s", w.Body.String())
		}

		helper.AssertSuccess(w, 200, "获取精选书籍失败")
		helper.LogSuccess("精选书籍获取成功")
	})

	t.Run("5.书城首页_获取活动Banner", func(t *testing.T) {
		url := fmt.Sprintf("%s?limit=5", BookstoreBannersPath)
		w := helper.DoRequest("GET", url, nil, "")
		helper.AssertSuccess(w, 200, "获取Banner失败")
		helper.LogSuccess("Banner列表获取成功")
	})

	// 榜单测试
	t.Run("6.榜单_获取实时榜", func(t *testing.T) {
		url := fmt.Sprintf("%s?limit=20", BookstoreRankingRealtimePath)
		w := helper.DoRequest("GET", url, nil, "")
		response := helper.AssertSuccess(w, 200, "获取实时榜失败")

		if data, ok := response["data"].([]interface{}); ok && len(data) > 0 {
			helper.LogSuccess(fmt.Sprintf("实时榜获取成功，书籍数量: %d", len(data)))
		} else {
			t.Logf("○ 实时榜为空（需要先更新榜单）")
		}
	})

	t.Run("7.榜单_获取周榜", func(t *testing.T) {
		url := fmt.Sprintf("%s?limit=20", BookstoreRankingWeeklyPath)
		w := helper.DoRequest("GET", url, nil, "")
		helper.AssertSuccess(w, 200, "获取周榜失败")
		helper.LogSuccess("周榜获取成功")
	})

	t.Run("8.榜单_获取月榜", func(t *testing.T) {
		url := fmt.Sprintf("%s?limit=20", BookstoreRankingMonthlyPath)
		w := helper.DoRequest("GET", url, nil, "")
		helper.AssertSuccess(w, 200, "获取月榜失败")
		helper.LogSuccess("月榜获取成功")
	})

	t.Run("9.榜单_获取新人榜", func(t *testing.T) {
		url := fmt.Sprintf("%s?limit=20", BookstoreRankingNewbiePath)
		w := helper.DoRequest("GET", url, nil, "")
		helper.AssertSuccess(w, 200, "获取新人榜失败")
		helper.LogSuccess("新人榜获取成功")
	})

	// 验证数据库中有数据
	t.Run("10.验证_数据库书籍统计", func(t *testing.T) {
		if global.DB == nil {
			t.Skip("global.DB 未注入（迁移到 ServiceContainer），跳过数据库直连统计")
		}

		bookCount, err := global.DB.Collection("books").CountDocuments(helper.ctx, bson.M{})
		if err != nil {
			t.Fatalf("统计书籍失败: %v", err)
		}

		chapterCount, err := global.DB.Collection("chapters").CountDocuments(helper.ctx, bson.M{})
		if err != nil {
			t.Fatalf("统计章节失败: %v", err)
		}

		helper.LogSuccess("数据库统计:")
		t.Logf("  - 书籍总数: %d", bookCount)
		t.Logf("  - 章节总数: %d", chapterCount)

		if bookCount == 0 {
			t.Logf("  ⚠ 提示：数据库中没有书籍，请先运行导入脚本")
		}
	})

	t.Logf("\n=== 书城流程测试完成 ===")
	t.Logf("测试场景: 首页数据 → 分类 → 推荐/精选 → 榜单")
}

// 测试分类筛选功能
func TestBookstoreCategoryFilter(t *testing.T) {
	// 设置测试环境
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()
	if global.DB == nil {
		t.Skip("global.DB 未注入（迁移到 ServiceContainer），跳过数据库直连分类筛选测试")
	}

	// 创建TestHelper
	helper := NewTestHelper(t, nil)

	t.Run("分类筛选_按分类浏览书籍", func(t *testing.T) {
		// 首先获取一个分类
		categories := []string{"玄幻", "都市", "科幻", "历史"}

		for _, category := range categories {
			// 直接从数据库查询该分类的书籍
			cursor, err := global.DB.Collection("books").Find(helper.ctx, bson.M{
				"categories": category,
			})

			if err != nil {
				t.Logf("查询分类 %s 失败: %v", category, err)
				continue
			}

			var books []bson.M
			err = cursor.All(helper.ctx, &books)
			cursor.Close(helper.ctx)

			if err == nil && len(books) > 0 {
				helper.LogSuccess(fmt.Sprintf("分类 [%s] 书籍数量: %d", category, len(books)))
			}
		}
	})
}
