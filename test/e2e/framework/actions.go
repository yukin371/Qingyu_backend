package e2e

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stretchr/testify/require"

	"Qingyu_backend/global"
	"Qingyu_backend/models/users"
	userRepo "Qingyu_backend/repository/mongodb/user"
)

// BusinessActions 业务操作辅助
type BusinessActions struct {
	env *TestEnvironment
}

// Actions 获取业务操作辅助器
func (env *TestEnvironment) Actions() *BusinessActions {
	return &BusinessActions{env: env}
}

// ============ 认证相关操作 ============

// Login 用户登录并返回 token
func (ba *BusinessActions) Login(username, password string) string {
	reqBody := map[string]interface{}{
		"username": username,
		"password": password,
	}
	w := ba.env.DoRequest("POST", "/api/v1/user/auth/login", reqBody, "")

	require.Equal(ba.env.T, 200, w.Code, "登录失败")

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	data, ok := resp["data"].(map[string]interface{})
	require.True(ba.env.T, ok, "响应数据格式错误")

	token, ok := data["token"].(string)
	require.True(ba.env.T, ok, "获取 token 失败")

	ba.env.LogSuccess("用户登录: %s (token: %s...)", username, token[:20])

	return token
}

// ============ 用户相关操作 ============

// GetUser 获取用户信息
func (ba *BusinessActions) GetUser(userID string) *users.User {
	userRepository := userRepo.NewMongoUserRepository(global.DB)
	user, err := userRepository.GetByID(context.Background(), userID)
	require.NoError(ba.env.T, err, "获取用户失败")
	return user
}

// SetUserVIPLevel 设置用户VIP等级（直接操作数据库）
func (ba *BusinessActions) SetUserVIPLevel(userID string, level int) {
	userRepository := userRepo.NewMongoUserRepository(global.DB)
	updates := map[string]interface{}{
		"vip_level": level,
	}
	err := userRepository.Update(context.Background(), userID, updates)
	require.NoError(ba.env.T, err, "更新用户VIP等级失败")

	ba.env.LogSuccess("设置用户VIP等级: %s -> %d", userID, level)
}

// ============ 书籍相关操作 ============

// GetBookstoreHomepage 获取书城首页
func (ba *BusinessActions) GetBookstoreHomepage() map[string]interface{} {
	w := ba.env.DoRequest("GET", "/api/v1/bookstore/homepage", nil, "")
	require.Equal(ba.env.T, 200, w.Code, "获取书城首页失败")

	return ba.env.ParseJSONResponse(w)
}

// GetRankings 获取榜单
func (ba *BusinessActions) GetRankings(rankingType string) map[string]interface{} {
	path := fmt.Sprintf("/api/v1/bookstore/rankings/%s", rankingType)
	w := ba.env.DoRequest("GET", path, nil, "")
	require.Equal(ba.env.T, 200, w.Code, "获取榜单失败")

	return ba.env.ParseJSONResponse(w)
}

// GetBookDetail 获取书籍详情
func (ba *BusinessActions) GetBookDetail(bookID string) map[string]interface{} {
	path := fmt.Sprintf("/api/v1/bookstore/books/%s", bookID)
	w := ba.env.DoRequest("GET", path, nil, "")
	require.Equal(ba.env.T, 200, w.Code, "获取书籍详情失败")

	return ba.env.ParseJSONResponse(w)
}

// SearchBooks 搜索书籍
func (ba *BusinessActions) SearchBooks(keyword string) map[string]interface{} {
	w := ba.env.DoRequest("GET", "/api/v1/bookstore/books/search?keyword="+keyword, nil, "")
	require.Equal(ba.env.T, 200, w.Code, "搜索书籍失败")

	return ba.env.ParseJSONResponse(w)
}

// ============ 阅读相关操作 ============

// GetChapter 获取章节内容
func (ba *BusinessActions) GetChapter(chapterID, token string) map[string]interface{} {
	path := fmt.Sprintf("/api/v1/bookstore/chapters/%s", chapterID)
	w := ba.env.DoRequest("GET", path, nil, token)
	require.Equal(ba.env.T, 200, w.Code, "获取章节内容失败")

	return ba.env.ParseJSONResponse(w)
}

// GetChapterList 获取章节列表
func (ba *BusinessActions) GetChapterList(bookID string, token string) map[string]interface{} {
	// 使用 /api/v1/reader/books/:bookId/chapters 路由（需要认证）
	path := fmt.Sprintf("/api/v1/reader/books/%s/chapters", bookID)
	w := ba.env.DoRequest("GET", path, nil, token)
	require.Equal(ba.env.T, 200, w.Code, "获取章节列表失败")

	return ba.env.ParseJSONResponse(w)
}

// StartReading 开始阅读（保存阅读进度）
func (ba *BusinessActions) StartReading(userID, bookID, chapterID, token string) map[string]interface{} {
	reqBody := map[string]interface{}{
		"bookId":    bookID,
		"chapterId": chapterID,
		"progress":  0.1, // 阅读进度 0-1（使用 0.1 避免零值问题）
	}
	w := ba.env.DoRequest("POST", "/api/v1/reader/progress", reqBody, token)
	require.Equal(ba.env.T, 200, w.Code, "保存阅读进度失败")

	ba.env.LogSuccess("开始阅读: 书籍=%s, 章节=%s", bookID, chapterID)

	return ba.env.ParseJSONResponse(w)
}

// GetReadingProgress 获取阅读进度
func (ba *BusinessActions) GetReadingProgress(userID, bookID string) map[string]interface{} {
	w := ba.env.DoRequest("GET", "/api/v1/reader/progress", nil, "")
	require.Equal(ba.env.T, 200, w.Code, "获取阅读进度失败")

	return ba.env.ParseJSONResponse(w)
}

// GetReadingHistory 获取阅读历史
func (ba *BusinessActions) GetReadingHistory(userID string) map[string]interface{} {
	w := ba.env.DoRequest("GET", "/api/v1/reader/progress/history", nil, "")
	require.Equal(ba.env.T, 200, w.Code, "获取阅读历史失败")

	return ba.env.ParseJSONResponse(w)
}

// ============ 社交互动操作 ============

// AddComment 发表评论
func (ba *BusinessActions) AddComment(token, bookID, chapterID, content string) map[string]interface{} {
	reqBody := map[string]interface{}{
		"book_id":    bookID,
		"chapter_id": chapterID,
		"content":    content,
	}
	w := ba.env.DoRequest("POST", "/api/v1/social/comments", reqBody, token)
	require.True(ba.env.T, w.Code == 200 || w.Code == 201, "发表评论失败")

	ba.env.LogSuccess("发表评论: %s", content)

	return ba.env.ParseJSONResponse(w)
}

// GetBookComments 获取书籍评论
func (ba *BusinessActions) GetBookComments(bookID string) map[string]interface{} {
	w := ba.env.DoRequest("GET", "/api/v1/social/comments?book_id="+bookID, nil, "")
	require.Equal(ba.env.T, 200, w.Code, "获取评论失败")

	return ba.env.ParseJSONResponse(w)
}

// CollectBook 收藏书籍
func (ba *BusinessActions) CollectBook(token, bookID string) map[string]interface{} {
	reqBody := map[string]interface{}{
		"book_id": bookID,
	}
	w := ba.env.DoRequest("POST", "/api/v1/social/collections", reqBody, token)
	require.True(ba.env.T, w.Code == 200 || w.Code == 201, "收藏书籍失败")

	ba.env.LogSuccess("收藏书籍: %s", bookID)

	return ba.env.ParseJSONResponse(w)
}

// GetReaderCollections 获取用户收藏列表
func (ba *BusinessActions) GetReaderCollections(userID string) map[string]interface{} {
	w := ba.env.DoRequest("GET", "/api/v1/social/collections", nil, "")
	require.Equal(ba.env.T, 200, w.Code, "获取收藏列表失败")

	return ba.env.ParseJSONResponse(w)
}

// LikeChapter 点赞书籍
func (ba *BusinessActions) LikeChapter(token, bookID string) map[string]interface{} {
	reqBody := map[string]interface{}{}
	w := ba.env.DoRequest("POST", "/api/v1/social/books/"+bookID+"/like", reqBody, token)
	require.True(ba.env.T, w.Code == 200 || w.Code == 201, "点赞失败")

	ba.env.LogSuccess("点赞书籍: %s", bookID)

	return ba.env.ParseJSONResponse(w)
}

// AddBookmark 添加书签
func (ba *BusinessActions) AddBookmark(token, bookID, chapterID string, position int) map[string]interface{} {
	reqBody := map[string]interface{}{
		"chapter_id": chapterID,
		"position":   position,
	}
	path := fmt.Sprintf("/api/v1/reader/books/%s/bookmarks", bookID)
	w := ba.env.DoRequest("POST", path, reqBody, token)
	require.Equal(ba.env.T, 200, w.Code, "添加书签失败")

	ba.env.LogSuccess("添加书签: 书籍=%s, 位置=%d", bookID, position)

	return ba.env.ParseJSONResponse(w)
}

// ============ 写作相关操作 ============

// CreateProject 创建写作项目
func (ba *BusinessActions) CreateProject(token string, req map[string]interface{}) map[string]interface{} {
	w := ba.env.DoRequest("POST", "/api/v1/writer/projects", req, token)
	require.True(ba.env.T, w.Code == 200 || w.Code == 201, "创建项目失败")

	ba.env.LogSuccess("创建项目: %s", req["title"])

	return ba.env.ParseJSONResponse(w)
}
