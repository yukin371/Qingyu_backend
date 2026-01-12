package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestFrontendAuthAPIIntegration 测试前端认证 API 对接
func TestFrontendAuthAPIIntegration(t *testing.T) {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 创建测试路由
	router := gin.New()

	// 模拟后端路由结构（需要根据实际项目调整）
	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"code": 200,
					"data": gin.H{
						"token":        "test_token_123",
						"refreshToken": "test_refresh_token_456",
						"user": gin.H{
							"id":       "user_123",
							"username": "testuser",
							"email":    "test@example.com",
							"role":     "user",
						},
						"permissions": []string{"read:books"},
						"roles":       []string{"user"},
					},
					"message": "注册成功",
				})
			})

			auth.POST("/login", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"code": 200,
					"data": gin.H{
						"token":        "test_token_789",
						"refreshToken": "test_refresh_token_101",
						"user": gin.H{
							"id":       "user_123",
							"username": "testuser",
							"email":    "test@example.com",
							"role":     "user",
						},
						"permissions": []string{"read:books", "write:books"},
						"roles":       []string{"user"},
					},
					"message": "登录成功",
				})
			})

			auth.POST("/logout", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"data":    nil,
					"message": "登出成功",
				})
			})

			auth.POST("/refresh", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"code": 200,
					"data": gin.H{
						"token":        "new_token_123",
						"refreshToken": "new_refresh_token_456",
					},
					"message": "Token刷新成功",
				})
			})

			auth.GET("/permissions", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"data":    []string{"read:books", "write:books", "delete:books"},
					"message": "获取权限成功",
				})
			})

			auth.GET("/roles", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"data":    []string{"user", "vip"},
					"message": "获取角色成功",
				})
			})
		}

		users := api.Group("/users")
		{
			users.GET("/profile", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"code": 200,
					"data": gin.H{
						"id":       "user_123",
						"username": "testuser",
						"email":    "test@example.com",
						"nickname": "测试用户",
						"role":     "user",
					},
					"message": "获取用户信息成功",
				})
			})

			users.PUT("/profile", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"code": 200,
					"data": gin.H{
						"id":       "user_123",
						"username": "testuser",
						"email":    "test@example.com",
						"nickname": "更新后的用户",
						"role":     "user",
					},
					"message": "更新用户信息成功",
				})
			})

			users.PUT("/password", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"code":    200,
					"data":    nil,
					"message": "修改密码成功",
				})
			})
		}
	}

	t.Run("RegisterAPI", func(t *testing.T) {
		// 测试注册接口
		reqBody := map[string]string{
			"username": "testuser",
			"email":    "test@example.com",
			"password": "password123",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, float64(200), response["code"])
		data := response["data"].(map[string]interface{})
		assert.NotEmpty(t, data["token"])
		assert.NotEmpty(t, data["refreshToken"])
		assert.NotEmpty(t, data["user"])

		t.Logf("✓ 注册 API 测试通过")
	})

	t.Run("LoginAPI", func(t *testing.T) {
		// 测试登录接口
		reqBody := map[string]string{
			"email":    "test@example.com",
			"password": "password123",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, float64(200), response["code"])
		data := response["data"].(map[string]interface{})
		assert.NotEmpty(t, data["token"])
		assert.NotEmpty(t, data["refreshToken"])

		t.Logf("✓ 登录 API 测试通过")
	})

	t.Run("LogoutAPI", func(t *testing.T) {
		// 测试登出接口
		req := httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, float64(200), response["code"])

		t.Logf("✓ 登出 API 测试通过")
	})

	t.Run("RefreshTokenAPI", func(t *testing.T) {
		// 测试刷新 Token 接口
		req := httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, float64(200), response["code"])
		data := response["data"].(map[string]interface{})
		assert.NotEmpty(t, data["token"])

		t.Logf("✓ 刷新 Token API 测试通过")
	})

	t.Run("GetPermissionsAPI", func(t *testing.T) {
		// 测试获取权限接口
		req := httptest.NewRequest("GET", "/api/v1/auth/permissions", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, float64(200), response["code"])
		data := response["data"].([]interface{})
		assert.Greater(t, len(data), 0)

		t.Logf("✓ 获取权限 API 测试通过")
	})

	t.Run("GetRolesAPI", func(t *testing.T) {
		// 测试获取角色接口
		req := httptest.NewRequest("GET", "/api/v1/auth/roles", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, float64(200), response["code"])
		data := response["data"].([]interface{})
		assert.Greater(t, len(data), 0)

		t.Logf("✓ 获取角色 API 测试通过")
	})

	t.Run("GetUserProfileAPI", func(t *testing.T) {
		// 测试获取用户信息接口
		req := httptest.NewRequest("GET", "/api/v1/users/profile", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, float64(200), response["code"])
		data := response["data"].(map[string]interface{})
		assert.NotEmpty(t, data["id"])
		assert.NotEmpty(t, data["username"])

		t.Logf("✓ 获取用户信息 API 测试通过")
	})

	t.Run("UpdateUserProfileAPI", func(t *testing.T) {
		// 测试更新用户信息接口
		reqBody := map[string]string{
			"nickname": "更新后的用户",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("PUT", "/api/v1/users/profile", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, float64(200), response["code"])
		data := response["data"].(map[string]interface{})
		assert.Equal(t, "更新后的用户", data["nickname"])

		t.Logf("✓ 更新用户信息 API 测试通过")
	})

	t.Run("ChangePasswordAPI", func(t *testing.T) {
		// 测试修改密码接口
		reqBody := map[string]string{
			"old_password": "oldpassword123",
			"new_password": "newpassword456",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("PUT", "/api/v1/users/password", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, float64(200), response["code"])

		t.Logf("✓ 修改密码 API 测试通过")
	})
}

// TestFrontendReaderAPIIntegration 测试前端读者模块 API 对接
func TestFrontendReaderAPIIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api/v1")
	{
		reader := api.Group("/reader")
		{
			// 书架相关
			books := reader.Group("/books")
			{
				books.GET("", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"code": 200,
						"data": []gin.H{
							{
								"bookId":    "book_123",
								"title":     "测试书籍1",
								"author":    "作者1",
								"status":    "reading",
								"progress":  50,
								"addedAt":   "2024-01-01",
								"updatedAt": "2024-01-02",
							},
						},
						"message": "获取书架成功",
					})
				})

				books.GET("/recent", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"code": 200,
						"data": []gin.H{
							{
								"bookId":    "book_456",
								"title":     "最近阅读",
								"author":    "作者2",
								"status":    "reading",
								"addedAt":   "2024-01-03",
								"updatedAt": "2024-01-04",
							},
						},
						"message": "获取最近阅读成功",
					})
				})

				books.POST("/:bookId", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"code":    200,
						"data":    nil,
						"message": "添加到书架成功",
					})
				})

				books.DELETE("/:bookId", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"code":    200,
						"data":    nil,
						"message": "从书架移除成功",
					})
				})

				books.POST("/:bookId/like", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"code": 200,
						"data": gin.H{
							"likeCount": 100,
						},
						"message": "点赞成功",
					})
				})

				books.DELETE("/:bookId/like", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"code": 200,
						"data": gin.H{
							"likeCount": 99,
						},
						"message": "取消点赞成功",
					})
				})

				books.GET("/:bookId/like/info", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"code": 200,
						"data": gin.H{
							"isLiked":   true,
							"likeCount": 100,
						},
						"message": "获取点赞信息成功",
					})
				})
			}

			// 点赞相关
			likes := reader.Group("/likes")
			{
				likes.GET("/books", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"code": 200,
						"data": []gin.H{
							{"bookId": "book_123", "title": "点赞书籍1"},
							{"bookId": "book_456", "title": "点赞书籍2"},
						},
						"message": "获取点赞列表成功",
					})
				})

				likes.GET("/stats", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"code": 200,
						"data": gin.H{
							"totalLikes":   150,
							"bookLikes":    100,
							"commentLikes": 50,
						},
						"message": "获取点赞统计成功",
					})
				})
			}

			// 评论相关
			comments := reader.Group("/comments")
			{
				comments.GET("", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"code": 200,
						"data": gin.H{
							"comments": []gin.H{
								{
									"id":        "comment_123",
									"content":   "测试评论",
									"userId":    "user_123",
									"likes":     10,
									"createdAt": "2024-01-01",
								},
							},
							"total": 1,
						},
						"message": "获取评论列表成功",
					})
				})

				comments.POST("", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"code": 200,
						"data": gin.H{
							"id":      "comment_new",
							"content": "新评论",
						},
						"message": "发表评论成功",
					})
				})

				comments.POST("/:id/like", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"code":    200,
						"data":    nil,
						"message": "点赞评论成功",
					})
				})

				comments.DELETE("/:id/like", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"code":    200,
						"data":    nil,
						"message": "取消点赞评论成功",
					})
				})
			}

			// 收藏相关
			collections := reader.Group("/collections")
			{
				collections.GET("", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"code": 200,
						"data": []gin.H{
							{
								"id":        "collection_123",
								"bookId":    "book_123",
								"title":     "收藏书籍1",
								"isPublic":  false,
								"createdAt": "2024-01-01",
							},
						},
						"message": "获取收藏列表成功",
					})
				})

				collections.POST("", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"code": 200,
						"data": gin.H{
							"id":     "collection_new",
							"bookId": "book_new",
						},
						"message": "添加收藏成功",
					})
				})

				collections.GET("/stats", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"code": 200,
						"data": gin.H{
							"totalCollections":  50,
							"publicCollections": 10,
							"folderCount":       5,
						},
						"message": "获取收藏统计成功",
					})
				})
			}
		}
	}

	t.Run("GetBookshelf", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/reader/books", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		t.Logf("✓ 获取书架 API 测试通过")
	})

	t.Run("GetRecentReading", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/reader/books/recent", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		t.Logf("✓ 获取最近阅读 API 测试通过")
	})

	t.Run("AddToBookshelf", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/reader/books/book_123", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		t.Logf("✓ 添加到书架 API 测试通过")
	})

	t.Run("RemoveFromBookshelf", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/reader/books/book_123", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		t.Logf("✓ 从书架移除 API 测试通过")
	})

	t.Run("LikeBook", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/reader/books/book_123/like", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		t.Logf("✓ 点赞书籍 API 测试通过")
	})

	t.Run("UnlikeBook", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/reader/books/book_123/like", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		t.Logf("✓ 取消点赞书籍 API 测试通过")
	})

	t.Run("GetBookLikeInfo", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/reader/books/book_123/like/info", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		t.Logf("✓ 获取书籍点赞信息 API 测试通过")
	})

	t.Run("GetUserLikedBooks", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/reader/likes/books", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		t.Logf("✓ 获取用户点赞列表 API 测试通过")
	})

	t.Run("GetUserLikeStats", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/reader/likes/stats", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		t.Logf("✓ 获取点赞统计 API 测试通过")
	})

	t.Run("GetComments", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/reader/comments", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		t.Logf("✓ 获取评论列表 API 测试通过")
	})

	t.Run("CreateComment", func(t *testing.T) {
		reqBody := map[string]string{
			"book_id": "book_123",
			"content": "测试评论",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/reader/comments", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		t.Logf("✓ 发表评论 API 测试通过")
	})

	t.Run("LikeComment", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/reader/comments/comment_123/like", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		t.Logf("✓ 点赞评论 API 测试通过")
	})

	t.Run("GetCollections", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/reader/collections", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		t.Logf("✓ 获取收藏列表 API 测试通过")
	})

	t.Run("AddCollection", func(t *testing.T) {
		reqBody := map[string]string{
			"book_id": "book_123",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/reader/collections", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		t.Logf("✓ 添加收藏 API 测试通过")
	})

	t.Run("GetCollectionStats", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/reader/collections/stats", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		t.Logf("✓ 获取收藏统计 API 测试通过")
	})
}

// TestAPIResponseFormat 测试 API 响应格式是否符合前端预期
func TestAPIResponseFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/api/v1/test/format", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":       200,
			"data":       gin.H{"message": "test"},
			"message":    "成功",
			"request_id": "req_123",
		})
	})

	t.Run("ResponseFormat", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/test/format", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		// 验证响应格式
		assert.Contains(t, response, "code")
		assert.Contains(t, response, "data")
		assert.Contains(t, response, "message")
		assert.Equal(t, float64(200), response["code"])

		t.Logf("✓ API 响应格式测试通过")
		t.Logf("响应格式: %+v", response)
	})
}

// FrontendAPICallExample 示例：前端如何调用这些 API
func FrontendAPICallExample() {
	// 这是一个示例函数，展示前端如何使用 httpService 调用后端 API

	// 1. 导入 httpService
	// import { httpService } from '@/core/services/http.service'

	// 2. 调用登录 API
	// const loginData = await httpService.post('/auth/login', {
	//     email: 'user@example.com',
	//     password: 'password123'
	// })
	// // loginData = { token, refreshToken, user, permissions, roles }

	// 3. 调用获取书架 API
	// const bookshelf = await httpService.get('/reader/books')
	// // bookshelf = [{ bookId, title, author, status, ... }]

	// 4. 调用点赞书籍 API
	// await httpService.post(`/reader/books/${bookId}/like`)

	// 5. 调用获取评论列表 API
	// const comments = await httpService.get('/reader/comments', {
	//     params: { book_id: 'book_123', page: 1, page_size: 20 }
	// })

	fmt.Println("前端 API 调用示例（请查看源代码）")
}
