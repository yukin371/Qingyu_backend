package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"Qingyu_backend/config"
	"Qingyu_backend/core"
	"Qingyu_backend/global"
)

// 书城测试 - 完整书城浏览流程
func TestBookstoreScenario(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// 初始化
	_, err := config.LoadConfig("../..")
	require.NoError(t, err, "加载配置失败")

	err = core.InitDB()
	require.NoError(t, err, "初始化数据库失败")

	ctx := context.Background()
	baseURL := "http://localhost:8080"

	t.Run("1.书城首页_获取首页数据", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/api/v1/bookstore/homepage", baseURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "应该成功获取首页数据")

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		t.Logf("✓ 首页数据获取成功")
		t.Logf("  响应: %v", result["message"])
	})

	t.Run("2.书城首页_获取分类树", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/api/v1/bookstore/categories/tree", baseURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			t.Logf("⚠ 获取分类树失败，HTTP状态: %d", resp.StatusCode)
			t.Logf("  响应体: %s", string(body))
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode, "应该成功获取分类树")

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Equal(t, float64(200), result["code"])
		data := result["data"]
		assert.NotNil(t, data, "分类树数据不应为空")

		t.Logf("✓ 分类树获取成功")
	})

	t.Run("3.书城首页_获取推荐书籍", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/api/v1/bookstore/books/recommended?page=1&size=10", baseURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			t.Logf("⚠ 获取推荐书籍失败，HTTP状态: %d", resp.StatusCode)
			t.Logf("  响应体: %s", string(body))
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode, "应该成功获取推荐书籍")

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Equal(t, float64(200), result["code"])
		data := result["data"]
		if data != nil {
			t.Logf("✓ 推荐书籍获取成功，数量: %d", len(data.([]interface{})))
		} else {
			t.Logf("○ 推荐书籍为空（可能还没有导入数据）")
		}
	})

	t.Run("4.书城首页_获取精选书籍", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/api/v1/bookstore/books/featured?page=1&size=10", baseURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			t.Logf("⚠ 获取精选书籍失败，HTTP状态: %d", resp.StatusCode)
			t.Logf("  响应体: %s", string(body))
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode, "应该成功获取精选书籍")

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Equal(t, float64(200), result["code"])
		t.Logf("✓ 精选书籍获取成功")
	})

	t.Run("5.书城首页_获取活动Banner", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/api/v1/bookstore/banners?limit=5", baseURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "应该成功获取Banner")

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Equal(t, float64(200), result["code"])
		t.Logf("✓ Banner列表获取成功")
	})

	// 榜单测试
	t.Run("6.榜单_获取实时榜", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/api/v1/bookstore/rankings/realtime?limit=20", baseURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "应该成功获取实时榜")

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Equal(t, float64(200), result["code"])
		data := result["data"]
		if data != nil && len(data.([]interface{})) > 0 {
			t.Logf("✓ 实时榜获取成功，书籍数量: %d", len(data.([]interface{})))
		} else {
			t.Logf("○ 实时榜为空（需要先更新榜单）")
		}
	})

	t.Run("7.榜单_获取周榜", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/api/v1/bookstore/rankings/weekly?limit=20", baseURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "应该成功获取周榜")

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Equal(t, float64(200), result["code"])
		t.Logf("✓ 周榜获取成功")
	})

	t.Run("8.榜单_获取月榜", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/api/v1/bookstore/rankings/monthly?limit=20", baseURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "应该成功获取月榜")

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Equal(t, float64(200), result["code"])
		t.Logf("✓ 月榜获取成功")
	})

	t.Run("9.榜单_获取新人榜", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/api/v1/bookstore/rankings/newbie?limit=20", baseURL))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "应该成功获取新人榜")

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Equal(t, float64(200), result["code"])
		t.Logf("✓ 新人榜获取成功")
	})

	// 验证数据库中有数据
	t.Run("10.验证_数据库书籍统计", func(t *testing.T) {
		bookCount, err := global.DB.Collection("books").CountDocuments(ctx, map[string]interface{}{})
		require.NoError(t, err)

		chapterCount, err := global.DB.Collection("chapters").CountDocuments(ctx, map[string]interface{}{})
		require.NoError(t, err)

		t.Logf("✓ 数据库统计:")
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
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	_, err := config.LoadConfig("../..")
	require.NoError(t, err)

	err = core.InitDB()
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("分类筛选_按分类浏览书籍", func(t *testing.T) {
		// 首先获取一个分类
		categories := []string{"玄幻", "都市", "科幻", "历史"}

		for _, category := range categories {
			// 直接从数据库查询该分类的书籍
			cursor, err := global.DB.Collection("books").Find(ctx, map[string]interface{}{
				"categories": category,
			})

			if err != nil {
				t.Logf("查询分类 %s 失败: %v", category, err)
				continue
			}

			var books []map[string]interface{}
			err = cursor.All(ctx, &books)
			cursor.Close(ctx)

			if err == nil && len(books) > 0 {
				t.Logf("✓ 分类 [%s] 书籍数量: %d", category, len(books))
			}
		}
	})
}

// 等待服务器启动的辅助函数
func waitForServer(baseURL string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(baseURL + "/api/v1/system/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("服务器在 %v 内未启动", timeout)
}
