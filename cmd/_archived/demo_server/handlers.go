// Package main 提供论文答辩演示用的API处理函数
package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// JWT签名密钥
var jwtSecret = []byte("qingyu_demo_secret_key_for_thesis_defense")

// ============ 认证处理 ============

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      *UserInfo `json:"user"`
}

// handleDemoLogin 处理登录
func handleDemoLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	// 查找用户
	user := MemoryStore.GetUserByUsername(c.Request.Context(), req.Username)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}

	// 验证密码
	if user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}

	// 生成JWT Token
	expiresAt := time.Now().Add(24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      expiresAt.Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成Token失败"})
		return
	}

	// 存储Token
	MemoryStore.StoreToken(c.Request.Context(), tokenString, user.ID)

	// 清除密码字段
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
		"data": LoginResponse{
			Token:     tokenString,
			ExpiresAt: expiresAt,
			User:      user,
		},
	})
}

// handleDemoRegister 处理注册
func handleDemoRegister(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	// 检查用户是否存在
	if MemoryStore.GetUserByUsername(c.Request.Context(), req.Username) != nil {
		c.JSON(http.StatusConflict, gin.H{"code": 409, "message": "用户名已存在"})
		return
	}

	// 创建用户
	user := UserInfo{
		ID:       "user_" + time.Now().Format("20060102150405"),
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Email:    req.Email,
		Role:     "reader",
	}

	createdUser := MemoryStore.CreateUser(c.Request.Context(), user)
	createdUser.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "注册成功",
		"data":    createdUser,
	})
}

// handleDemoLogout 处理登出
func handleDemoLogout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != "" {
		// 移除前缀 "Bearer "
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
		MemoryStore.DeleteToken(c.Request.Context(), token)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登出成功",
	})
}

// ============ 书城处理 ============

// handleDemoGetBooks 获取书籍列表
func handleDemoGetBooks(c *gin.Context) {
	category := c.Query("category")
	status := c.Query("status")
	page := 1
	pageSize := 20

	books := MemoryStore.ListBooks(c.Request.Context(), category, status, pageSize, (page-1)*pageSize)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":     books,
			"total":    len(books),
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// handleDemoGetBook 获取书籍详情
func handleDemoGetBook(c *gin.Context) {
	bookID := c.Param("id")
	book := MemoryStore.GetBookByID(c.Request.Context(), bookID)

	if book == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "书籍不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    book,
	})
}

// handleDemoGetChapters 获取章节列表
func handleDemoGetChapters(c *gin.Context) {
	bookID := c.Param("id")
	chapters := MemoryStore.ListChaptersByBook(c.Request.Context(), bookID)

	// 简化返回，不包含内容
	list := make([]gin.H, len(chapters))
	for i, ch := range chapters {
		list[i] = gin.H{
			"id":          ch.ID,
			"title":       ch.Title,
			"chapter_num": ch.ChapterNum,
			"word_count":  ch.WordCount,
			"is_free":     ch.IsFree,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":  list,
			"total": len(list),
		},
	})
}

// handleDemoGetChapter 获取章节内容
func handleDemoGetChapter(c *gin.Context) {
	chapterID := c.Param("id")
	chapter := MemoryStore.GetChapterByID(c.Request.Context(), chapterID)

	if chapter == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "章节不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    chapter,
	})
}

// handleDemoGetCategories 获取分类列表
func handleDemoGetCategories(c *gin.Context) {
	categories := MemoryStore.ListCategories(c.Request.Context())

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    categories,
	})
}

// handleDemoSearchBooks 搜索书籍
func handleDemoSearchBooks(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请输入搜索关键词"})
		return
	}

	books := MemoryStore.SearchBooks(c.Request.Context(), keyword, 20)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":  books,
			"total": len(books),
			"keyword": keyword,
		},
	})
}

// handleDemoGetRanking 获取排行榜
func handleDemoGetRanking(c *gin.Context) {
	// 简单返回按点击量排序的书籍
	books := MemoryStore.ListBooks(c.Request.Context(), "", "", 10, 0)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":  books,
			"title": "热门推荐",
		},
	})
}

// ============ 用户处理 ============

// handleDemoGetProfile 获取用户资料
func handleDemoGetProfile(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录"})
		return
	}

	user := MemoryStore.GetUserByID(c.Request.Context(), userID)
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    user,
	})
}

// handleDemoUpdateProfile 更新用户资料
func handleDemoUpdateProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功（演示模式）",
	})
}

// ============ 阅读器处理 ============

// handleDemoGetProgress 获取阅读进度
func handleDemoGetProgress(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"chapter_id":   "chap_001",
			"chapter_num":  1,
			"progress":     0.5,
			"last_read_at": time.Now(),
		},
	})
}

// handleDemoSaveProgress 保存阅读进度
func handleDemoSaveProgress(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "保存成功（演示模式）",
	})
}

// ============ 社交处理 ============

// handleDemoGetComments 获取评论列表
func handleDemoGetComments(c *gin.Context) {
	bookID := c.Param("bookId")
	comments := MemoryStore.ListCommentsByBook(c.Request.Context(), bookID, 10, 0)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":  comments,
			"total": len(comments),
		},
	})
}

// handleDemoCreateComment 创建评论
func handleDemoCreateComment(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录"})
		return
	}

	var req struct {
		BookID  string `json:"book_id" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	comment := CommentInfo{
		BookID:  req.BookID,
		UserID:  userID,
		Content: req.Content,
	}

	created := MemoryStore.CreateComment(c.Request.Context(), comment)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "评论成功",
		"data":    created,
	})
}

// ============ AI处理 ============

// handleDemoAIGenerate AI内容生成（模拟）
func handleDemoAIGenerate(c *gin.Context) {
	var req struct {
		Prompt string `json:"prompt" binding:"required"`
		Type   string `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	// 模拟AI生成响应
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"content":   "这是AI生成的内容演示。在实际系统中，这里会调用真实的AI服务来生成内容。\n\n您输入的提示词：" + req.Prompt,
			"model":     "demo-model",
			"tokens":    100,
			"timestamp": time.Now(),
		},
	})
}

// handleDemoAIContinue AI续写（模拟）
func handleDemoAIContinue(c *gin.Context) {
	var req struct {
		Context string `json:"context" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	// 模拟AI续写响应
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"content":   "\n\n（演示模式）这是AI续写的内容。在实际系统中，AI会根据您提供的上下文继续创作故事。\n\n" + req.Context[:min(100, len(req.Context))] + "...",
			"model":     "demo-model",
			"tokens":    150,
			"timestamp": time.Now(),
		},
	})
}

// ============ 管理处理 ============

// handleDemoGetStats 获取统计数据
func handleDemoGetStats(c *gin.Context) {
	stats := MemoryStore.GetStats(c.Request.Context())

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"users":      stats["users"],
			"books":      stats["books"],
			"chapters":   stats["chapters"],
			"categories": stats["categories"],
			"comments":   stats["comments"],
			"uptime":     time.Since(time.Now()).Seconds(),
		},
	})
}

// handleDemoListUsers 列出用户
func handleDemoListUsers(c *gin.Context) {
	users := MemoryStore.ListUsers(c.Request.Context(), 100, 0)

	// 清除密码
	for _, u := range users {
		u.Password = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":  users,
			"total": len(users),
		},
	})
}

// ============ 辅助函数 ============

// getUserIDFromContext 从上下文获取用户ID
func getUserIDFromContext(c *gin.Context) string {
	token := c.GetHeader("Authorization")
	if token == "" {
		return ""
	}

	// 移除前缀 "Bearer "
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// 从存储中验证
	userID := MemoryStore.GetUserIDByToken(c.Request.Context(), token)
	if userID != "" {
		return userID
	}

	// 尝试解析JWT
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !parsedToken.Valid {
		return ""
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return ""
	}

	if userID, ok := claims["user_id"].(string); ok {
		return userID
	}

	return ""
}

// min 返回较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
