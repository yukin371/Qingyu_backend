package social_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	socialAPI "Qingyu_backend/api/v1/social"
	"Qingyu_backend/models/social"
	"Qingyu_backend/service/interfaces"
)

// MockCommentService 模拟评论服务接口
type MockCommentService struct {
	mock.Mock
}

// PublishComment 模拟发表评论
func (m *MockCommentService) PublishComment(ctx context.Context, userID, bookID, chapterID, content string, rating int) (*social.Comment, error) {
	args := m.Called(ctx, userID, bookID, chapterID, content, rating)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*social.Comment), args.Error(1)
}

// ReplyComment 模拟回复评论
func (m *MockCommentService) ReplyComment(ctx context.Context, userID, parentCommentID, content string) (*social.Comment, error) {
	args := m.Called(ctx, userID, parentCommentID, content)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*social.Comment), args.Error(1)
}

// GetCommentList 模拟获取评论列表
func (m *MockCommentService) GetCommentList(ctx context.Context, bookID string, sortBy string, page, size int) ([]*social.Comment, int64, error) {
	args := m.Called(ctx, bookID, sortBy, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.Comment), args.Get(1).(int64), args.Error(2)
}

// GetCommentDetail 模拟获取评论详情
func (m *MockCommentService) GetCommentDetail(ctx context.Context, commentID string) (*social.Comment, error) {
	args := m.Called(ctx, commentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*social.Comment), args.Error(1)
}

// UpdateComment 模拟更新评论
func (m *MockCommentService) UpdateComment(ctx context.Context, userID, commentID, content string) error {
	args := m.Called(ctx, userID, commentID, content)
	return args.Error(0)
}

// DeleteComment 模拟删除评论
func (m *MockCommentService) DeleteComment(ctx context.Context, userID, commentID string) error {
	args := m.Called(ctx, userID, commentID)
	return args.Error(0)
}

// LikeComment 模拟点赞评论
func (m *MockCommentService) LikeComment(ctx context.Context, userID, commentID string) error {
	args := m.Called(ctx, userID, commentID)
	return args.Error(0)
}

// UnlikeComment 模拟取消点赞评论
func (m *MockCommentService) UnlikeComment(ctx context.Context, userID, commentID string) error {
	args := m.Called(ctx, userID, commentID)
	return args.Error(0)
}

// AutoReviewComment 模拟自动审核评论
func (m *MockCommentService) AutoReviewComment(ctx context.Context, comment *social.Comment) (social.CommentState, string, error) {
	args := m.Called(ctx, comment)
	return args.Get(0).(social.CommentState), args.String(1), args.Error(2)
}

// GetBookCommentStats 模拟获取书籍评论统计
func (m *MockCommentService) GetBookCommentStats(ctx context.Context, bookID string) (map[string]interface{}, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// GetUserComments 模拟获取用户评论列表
func (m *MockCommentService) GetUserComments(ctx context.Context, userID string, page, size int) ([]*social.Comment, int64, error) {
	args := m.Called(ctx, userID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.Comment), args.Get(1).(int64), args.Error(2)
}

// GetCommentThread 模拟获取评论线程
func (m *MockCommentService) GetCommentThread(ctx context.Context, commentID string) (*social.CommentThread, error) {
	args := m.Called(ctx, commentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*social.CommentThread), args.Error(1)
}

// GetTopComments 模拟获取热门评论
func (m *MockCommentService) GetTopComments(ctx context.Context, bookID string, limit int) ([]*social.Comment, error) {
	args := m.Called(ctx, bookID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*social.Comment), args.Error(1)
}

// GetCommentReplies 模拟获取评论回复
func (m *MockCommentService) GetCommentReplies(ctx context.Context, commentID string, page, size int) ([]*social.Comment, int64, error) {
	args := m.Called(ctx, commentID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.Comment), args.Get(1).(int64), args.Error(2)
}

// setupCommentTestRouter 设置测试路由
func setupCommentTestRouter(commentService interfaces.CommentService, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 添加middleware来设置userId（用于需要认证的端点）
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("userId", userID)
		}
		c.Next()
	})

	api := socialAPI.NewCommentAPI(commentService)

	v1 := r.Group("/api/v1/reader/comments")
	{
		v1.POST("", api.CreateComment)
		v1.GET("", api.GetCommentList)
		v1.GET("/:id", api.GetCommentDetail)
		v1.PUT("/:id", api.UpdateComment)
		v1.DELETE("/:id", api.DeleteComment)
		v1.POST("/:id/reply", api.ReplyComment)
		v1.POST("/:id/like", api.LikeComment)
		v1.DELETE("/:id/like", api.UnlikeComment)
		v1.GET("/:id/thread", api.GetCommentThread)
		v1.GET("/top", api.GetTopComments)
		v1.GET("/:id/replies", api.GetCommentReplies)
	}

	return r
}

// TestCommentAPI_CreateComment_Success 测试成功发表评论
func TestCommentAPI_CreateComment_Success(t *testing.T) {
	// Given
	mockService := new(MockCommentService)
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	router := setupCommentTestRouter(mockService, userID)

	reqBody := map[string]interface{}{
		"book_id":    bookID,
		"content":    "这是一本非常好的书，推荐大家阅读！",
		"rating":     5,
		"chapter_id": "",
	}

	expectedComment := &social.Comment{}
	expectedComment.ID = primitive.NewObjectID()
	expectedComment.AuthorID = userID
	expectedComment.TargetID = bookID
	expectedComment.Content = "这是一本非常好的书，推荐大家阅读！"
	expectedComment.Rating = 5
	expectedComment.State = social.CommentStateNormal

	mockService.On("PublishComment", mock.Anything, userID, bookID, "", "这是一本非常好的书，推荐大家阅读！", 5).Return(expectedComment, nil)

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/comments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusCreated), response["code"])
	assert.NotNil(t, response["data"])

	mockService.AssertExpectations(t)
}

// TestCommentAPI_CreateComment_MissingBookID 测试缺少book_id
func TestCommentAPI_CreateComment_MissingBookID(t *testing.T) {
	// Given
	mockService := new(MockCommentService)
	userID := primitive.NewObjectID().Hex()
	router := setupCommentTestRouter(mockService, userID)

	reqBody := map[string]interface{}{
		"content": "这是一本非常好的书！",
		"rating":  5,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/comments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusBadRequest), response["code"])
	assert.Contains(t, response["message"], "参数错误")
}

// TestCommentAPI_CreateComment_ContentTooShort 测试评论内容过短
func TestCommentAPI_CreateComment_ContentTooShort(t *testing.T) {
	// Given
	mockService := new(MockCommentService)
	userID := primitive.NewObjectID().Hex()
	router := setupCommentTestRouter(mockService, userID)

	reqBody := map[string]interface{}{
		"book_id": primitive.NewObjectID().Hex(),
		"content": "太短了",
		"rating":  5,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/comments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusBadRequest), response["code"])
}

// TestCommentAPI_CreateComment_Unauthorized 测试未授权访问
func TestCommentAPI_CreateComment_Unauthorized(t *testing.T) {
	// Given
	mockService := new(MockCommentService)
	router := setupCommentTestRouter(mockService, "") // 不设置userID

	reqBody := map[string]interface{}{
		"book_id": primitive.NewObjectID().Hex(),
		"content": "这是一本非常好的书，推荐大家阅读！",
		"rating":  5,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/comments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusUnauthorized), response["code"])
	assert.Contains(t, response["message"], "未授权")
}

// TestCommentAPI_GetCommentList_Success 测试成功获取评论列表
func TestCommentAPI_GetCommentList_Success(t *testing.T) {
	// Given
	mockService := new(MockCommentService)
	router := setupCommentTestRouter(mockService, "")

	bookID := primitive.NewObjectID().Hex()
	comment1 := &social.Comment{}
	comment1.ID = primitive.NewObjectID()
	comment1.Content = "非常好的书！"
	comment1.Rating = 5

	comment2 := &social.Comment{}
	comment2.ID = primitive.NewObjectID()
	comment2.Content = "推荐阅读"
	comment2.Rating = 4

	expectedComments := []*social.Comment{comment1, comment2}

	mockService.On("GetCommentList", mock.Anything, bookID, "latest", 1, 20).Return(expectedComments, int64(2), nil)

	req, _ := http.NewRequest("GET", "/api/v1/reader/comments?book_id="+bookID+"&sortBy=latest&page=1&size=20", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])

	data := response["data"].(map[string]interface{})
	assert.NotNil(t, data["comments"])
	assert.Equal(t, float64(2), data["total"])

	mockService.AssertExpectations(t)
}

// TestCommentAPI_GetCommentList_MissingBookID 测试缺少book_id参数
func TestCommentAPI_GetCommentList_MissingBookID(t *testing.T) {
	// Given
	mockService := new(MockCommentService)
	router := setupCommentTestRouter(mockService, "")

	req, _ := http.NewRequest("GET", "/api/v1/reader/comments?sortBy=latest", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusBadRequest), response["code"])
	assert.Equal(t, "参数错误", response["message"])
	// errorDetail字段可能不存在，如果存在则检查
	if errorDetail, ok := response["error"]; ok {
		assert.Equal(t, "书籍ID不能为空", errorDetail)
	}
}

// TestCommentAPI_DeleteComment_Success 测试成功删除评论
func TestCommentAPI_DeleteComment_Success(t *testing.T) {
	// Given
	mockService := new(MockCommentService)
	userID := primitive.NewObjectID().Hex()
	commentID := primitive.NewObjectID().Hex()
	router := setupCommentTestRouter(mockService, userID)

	mockService.On("DeleteComment", mock.Anything, userID, commentID).Return(nil)

	req, _ := http.NewRequest("DELETE", "/api/v1/reader/comments/"+commentID, nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.Contains(t, response["message"], "删除成功")

	mockService.AssertExpectations(t)
}

// TestCommentAPI_LikeComment_Success 测试成功点赞评论
func TestCommentAPI_LikeComment_Success(t *testing.T) {
	// Given
	mockService := new(MockCommentService)
	userID := primitive.NewObjectID().Hex()
	commentID := primitive.NewObjectID().Hex()
	router := setupCommentTestRouter(mockService, userID)

	mockService.On("LikeComment", mock.Anything, userID, commentID).Return(nil)

	req, _ := http.NewRequest("POST", "/api/v1/reader/comments/"+commentID+"/like", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.Contains(t, response["message"], "点赞成功")

	mockService.AssertExpectations(t)
}

// TestCommentAPI_UpdateComment_Success 测试成功更新评论
func TestCommentAPI_UpdateComment_Success(t *testing.T) {
	// Given
	mockService := new(MockCommentService)
	userID := primitive.NewObjectID().Hex()
	commentID := primitive.NewObjectID().Hex()
	router := setupCommentTestRouter(mockService, userID)

	reqBody := map[string]interface{}{
		"content": "更新后的评论内容，这是一本非常好的书！",
	}

	mockService.On("UpdateComment", mock.Anything, userID, commentID, "更新后的评论内容，这是一本非常好的书！").Return(nil)

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/reader/comments/"+commentID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.Contains(t, response["message"], "更新成功")

	mockService.AssertExpectations(t)
}

// TestCommentAPI_ReplyComment_Success 测试成功回复评论
func TestCommentAPI_ReplyComment_Success(t *testing.T) {
	// Given
	mockService := new(MockCommentService)
	userID := primitive.NewObjectID().Hex()
	parentCommentID := primitive.NewObjectID().Hex()
	router := setupCommentTestRouter(mockService, userID)

	reqBody := map[string]interface{}{
		"content": "我同意你的观点，这本书确实很棒！",
	}

	expectedReply := &social.Comment{}
	expectedReply.ID = primitive.NewObjectID()
	expectedReply.AuthorID = userID
	expectedReply.Content = "我同意你的观点，这本书确实很棒！"
	expectedReply.State = social.CommentStateNormal

	mockService.On("ReplyComment", mock.Anything, userID, parentCommentID, "我同意你的观点，这本书确实很棒！").Return(expectedReply, nil)

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/comments/"+parentCommentID+"/reply", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusCreated), response["code"])
	assert.NotNil(t, response["data"])

	mockService.AssertExpectations(t)
}
