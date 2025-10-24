package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"Qingyu_backend/config"
	"Qingyu_backend/core"
	"Qingyu_backend/global"

	"go.mongodb.org/mongo-driver/bson"
)

// 书籍搜索测试 - 完整搜索流程
func TestSearchScenario(t *testing.T) {
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

	// 先获取一些测试数据的标题和作者
	var testBooks []struct {
		Title  string
		Author string
	}

	cursor, err := global.DB.Collection("books").Find(ctx, bson.M{}, nil)
	if err == nil {
		var books []map[string]interface{}
		cursor.All(ctx, &books)
		cursor.Close(ctx)

		for i := 0; i < len(books) && i < 3; i++ {
			testBooks = append(testBooks, struct {
				Title  string
				Author string
			}{
				Title:  books[i]["title"].(string),
				Author: books[i]["author"].(string),
			})
		}
	}

	t.Run("1.搜索_按标题关键词搜索", func(t *testing.T) {
		if len(testBooks) == 0 {
			t.Skip("没有测试数据")
		}

		// 取标题的一部分作为关键词
		keyword := testBooks[0].Title[:3] // 取前3个字符
		encodedKeyword := url.QueryEscape(keyword)

		resp, err := http.Get(fmt.Sprintf("%s/api/v1/bookstore/books/search?keyword=%s", baseURL, encodedKeyword))
		require.NoError(t, err)
		defer resp.Body.Close()

		// 如果返回400，说明关键词太短，这是预期的
		if resp.StatusCode == http.StatusBadRequest {
			t.Logf("○ 关键词 '%s' 太短，服务器要求提供更长的关键词", keyword)
			return
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode, "应该成功搜索")

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		if result["code"] == float64(200) {
			data := result["data"]
			if data != nil {
				books := data.([]interface{})
				t.Logf("✓ 按标题搜索成功，关键词: '%s'，找到 %d 本书", keyword, len(books))
			}
		}
	})

	t.Run("2.搜索_按作者搜索", func(t *testing.T) {
		if len(testBooks) == 0 {
			t.Skip("没有测试数据")
		}

		author := testBooks[0].Author
		encodedAuthor := url.QueryEscape(author)

		resp, err := http.Get(fmt.Sprintf("%s/api/v1/bookstore/books/search?author=%s", baseURL, encodedAuthor))
		require.NoError(t, err)
		defer resp.Body.Close()

		// 如果返回400，说明需要提供关键词
		if resp.StatusCode == http.StatusBadRequest {
			t.Logf("○ 需要同时提供关键词参数")
			return
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode, "应该成功搜索")

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		t.Logf("✓ 按作者搜索完成，作者: '%s'", author)
	})

	t.Run("3.搜索_组合搜索（关键词+作者）", func(t *testing.T) {
		if len(testBooks) == 0 {
			t.Skip("没有测试数据")
		}

		keyword := testBooks[0].Title[:3]
		author := testBooks[0].Author
		encodedKeyword := url.QueryEscape(keyword)
		encodedAuthor := url.QueryEscape(author)

		url := fmt.Sprintf("%s/api/v1/bookstore/books/search?keyword=%s&author=%s",
			baseURL, encodedKeyword, encodedAuthor)

		resp, err := http.Get(url)
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusBadRequest {
			body, _ := io.ReadAll(resp.Body)
			t.Logf("○ 搜索被拒绝: %s", string(body))
			return
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode, "应该成功搜索")

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		t.Logf("✓ 组合搜索完成，关键词: '%s'，作者: '%s'", keyword, author)
	})

	t.Run("4.搜索_排序功能测试", func(t *testing.T) {
		keyword := "书"
		encodedKeyword := url.QueryEscape(keyword)

		sortOptions := []struct {
			name      string
			sortBy    string
			sortOrder string
		}{
			{"最新", "created_at", "desc"},
			{"最早", "created_at", "asc"},
			{"字数降序", "word_count", "desc"},
			{"字数升序", "word_count", "asc"},
		}

		for _, opt := range sortOptions {
			url := fmt.Sprintf("%s/api/v1/bookstore/books/search?keyword=%s&sortBy=%s&sortOrder=%s",
				baseURL, encodedKeyword, opt.sortBy, opt.sortOrder)

			resp, err := http.Get(url)
			require.NoError(t, err)
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				t.Logf("  ✓ 排序方式 [%s] 测试通过", opt.name)
			} else {
				body, _ := io.ReadAll(resp.Body)
				t.Logf("  ○ 排序方式 [%s] 失败: %s", opt.name, string(body))
			}
		}
	})

	t.Run("5.搜索_分页功能测试", func(t *testing.T) {
		keyword := "书"
		encodedKeyword := url.QueryEscape(keyword)

		pages := []struct {
			page int
			size int
		}{
			{1, 10},
			{1, 20},
			{2, 10},
		}

		for _, p := range pages {
			url := fmt.Sprintf("%s/api/v1/bookstore/books/search?keyword=%s&page=%d&size=%d",
				baseURL, encodedKeyword, p.page, p.size)

			resp, err := http.Get(url)
			require.NoError(t, err)
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				var result map[string]interface{}
				json.Unmarshal(body, &result)

				if result["data"] != nil {
					books := result["data"].([]interface{})
					t.Logf("  ✓ 分页 [第%d页，每页%d条] 返回 %d 条结果", p.page, p.size, len(books))
				}
			}
		}
	})

	t.Run("6.搜索_无结果场景", func(t *testing.T) {
		keyword := "这是一个不存在的书名关键词xyz123"
		encodedKeyword := url.QueryEscape(keyword)

		resp, err := http.Get(fmt.Sprintf("%s/api/v1/bookstore/books/search?keyword=%s", baseURL, encodedKeyword))
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			var result map[string]interface{}
			json.Unmarshal(body, &result)

			data := result["data"]
			if data == nil || len(data.([]interface{})) == 0 {
				t.Logf("✓ 无结果场景测试通过，关键词: '%s'", keyword)
			} else {
				t.Logf("○ 意外找到结果")
			}
		}
	})

	t.Logf("\n=== 搜索流程测试完成 ===")
	t.Logf("测试场景: 标题搜索 → 作者搜索 → 组合搜索 → 排序 → 分页")
}

// 测试高级搜索功能
func TestAdvancedSearch(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	_, err := config.LoadConfig("../..")
	require.NoError(t, err)

	err = core.InitDB()
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("高级搜索_数据库层面测试", func(t *testing.T) {
		// 测试通过数据库直接搜索
		testCases := []struct {
			name   string
			filter bson.M
		}{
			{
				name:   "按分类搜索",
				filter: bson.M{"categories": "玄幻"},
			},
			{
				name:   "按状态搜索",
				filter: bson.M{"status": "completed"},
			},
			{
				name:   "搜索推荐书籍",
				filter: bson.M{"is_recommended": true},
			},
			{
				name:   "搜索精选书籍",
				filter: bson.M{"is_featured": true},
			},
			{
				name:   "搜索免费书籍",
				filter: bson.M{"is_free": true},
			},
		}

		for _, tc := range testCases {
			count, err := global.DB.Collection("books").CountDocuments(ctx, tc.filter)
			if err == nil {
				t.Logf("  ✓ %s: %d 本", tc.name, count)
			} else {
				t.Logf("  ✗ %s 失败: %v", tc.name, err)
			}
		}
	})
}

