package reader

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"Qingyu_backend/api/v1/shared"
	readerModels "Qingyu_backend/models/reader"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupCommentTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	commentAPI := NewChapterCommentAPI()

	router.GET("/api/v1/reader/chapters/:chapterId/comments", commentAPI.GetChapterComments)
	router.POST("/api/v1/reader/chapters/:chapterId/comments", setupAuthMiddleware(), commentAPI.CreateChapterComment)
	router.GET("/api/v1/reader/comments/:commentId", commentAPI.GetChapterComment)
	router.PUT("/api/v1/reader/comments/:commentId", setupAuthMiddleware(), commentAPI.UpdateChapterComment)
	router.DELETE("/api/v1/reader/comments/:commentId", setupAuthMiddleware(), commentAPI.DeleteChapterComment)
	router.POST("/api/v1/reader/comments/:commentId/like", setupAuthMiddleware(), commentAPI.LikeChapterComment)
	router.DELETE("/api/v1/reader/comments/:commentId/like", setupAuthMiddleware(), commentAPI.UnlikeChapterComment)
	router.GET("/api/v1/reader/chapters/:chapterId/paragraphs/:paragraphIndex/comments", commentAPI.GetParagraphComments)
	router.POST("/api/v1/reader/chapters/:chapterId/paragraph-comments", setupAuthMiddleware(), commentAPI.CreateParagraphComment)
	router.GET("/api/v1/reader/chapters/:chapterId/paragraph-comments", commentAPI.GetChapterParagraphComments)

	return router
}

// Test: GetChapterComments

func TestChapterCommentAPI_GetChapterComments_Success(t *testing.T) {
	router := setupCommentTestRouter()

	chapterID := primitive.NewObjectID()

	req, _ := http.NewRequest("GET", "/api/v1/reader/chapters/"+chapterID.Hex()+"/comments?page=1&pageSize=20", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response shared.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取成功", response.Message)

	// Data 字段应该包含 ChapterCommentListResponse
	dataBytes, err := json.Marshal(response.Data)
	assert.NoError(t, err)

	var commentList reader.ChapterCommentListResponse
	err = json.Unmarshal(dataBytes, &commentList)
	assert.NoError(t, err)
	assert.Equal(t, 1, commentList.Page)
	assert.Equal(t, 20, commentList.PageSize)
	assert.NotNil(t, commentList.Comments)
}

func TestChapterCommentAPI_GetChapterComments_InvalidChapterID(t *testing.T) {
	router := setupCommentTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/reader/chapters/invalid-id/comments", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response shared.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Message, "章节ID格式无效")
}

func TestChapterCommentAPI_GetChapterComments_EmptyChapterID(t *testing.T) {
	router := setupCommentTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/reader/chapters//comments", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response shared.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Message, "章节ID不能为空")
}

func TestChapterCommentAPI_GetChapterComments_WithSorting(t *testing.T) {
	router := setupCommentTestRouter()

	chapterID := primitive.NewObjectID()

	req, _ := http.NewRequest("GET", "/api/v1/reader/chapters/"+chapterID.Hex()+"/comments?sortBy=like_count&sortOrder=desc", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestChapterCommentAPI_GetChapterComments_InvalidSortField(t *testing.T) {
	router := setupCommentTestRouter()

	chapterID := primitive.NewObjectID()

	req, _ := http.NewRequest("GET", "/api/v1/reader/chapters/"+chapterID.Hex()+"/comments?sortBy=invalid_field", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code) // Should default to created_at
}

func TestChapterCommentAPI_GetChapterComments_WithParentFilter(t *testing.T) {
	router := setupCommentTestRouter()

	chapterID := primitive.NewObjectID()

	req, _ := http.NewRequest("GET", "/api/v1/reader/chapters/"+chapterID.Hex()+"/comments?parentId=", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestChapterCommentAPI_GetChapterComments_PageSizeTooLarge(t *testing.T) {
	router := setupCommentTestRouter()

	chapterID := primitive.NewObjectID()

	req, _ := http.NewRequest("GET", "/api/v1/reader/chapters/"+chapterID.Hex()+"/comments?pageSize=150", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code) // Should default to max 100
}

// Test: CreateChapterComment

func TestChapterCommentAPI_CreateChapterComment_Success(t *testing.T) {
	router := setupCommentTestRouter()

	chapterID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()

	reqBody := reader.CreateChapterCommentRequest{
		ChapterID: chapterID.Hex(),
		BookID:    bookID.Hex(),
		Content:   "This is a test comment",
		Rating:    5,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/chapters/"+chapterID.Hex()+"/comments", strings.NewReader(string(jsonBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response shared.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 201, response.Code)

	data := response.Data.(map[string]interface{})
	commentData := data["comment"].(map[string]interface{})

	commentBytes, _ := json.Marshal(commentData)
	var comment reader.ChapterComment
	err = json.Unmarshal(commentBytes, &comment)
	assert.NoError(t, err)
	assert.Equal(t, "This is a test comment", comment.Content)
	assert.Equal(t, 5, comment.Rating)
	assert.True(t, comment.IsVisible)
	assert.False(t, comment.IsDeleted)
}

func TestChapterCommentAPI_CreateChapterComment_WithReply(t *testing.T) {
	router := setupCommentTestRouter()

	chapterID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()
	parentID := primitive.NewObjectID()
	parentIDStr := parentID.Hex()

	reqBody := reader.CreateChapterCommentRequest{
		ChapterID: chapterID.Hex(),
		BookID:    bookID.Hex(),
		Content:   "This is a reply",
		Rating:    0,
		ParentID:  &parentIDStr,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/chapters/"+chapterID.Hex()+"/comments", strings.NewReader(string(jsonBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response shared.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data := response.Data.(map[string]interface{})
	commentData := data["comment"].(map[string]interface{})

	commentBytes, _ := json.Marshal(commentData)
	var comment reader.ChapterComment
	err = json.Unmarshal(commentBytes, &comment)
	assert.NoError(t, err)
	assert.NotNil(t, comment.ParentID)
}

func TestChapterCommentAPI_CreateChapterComment_InvalidChapterID(t *testing.T) {
	router := setupCommentTestRouter()

	bookID := primitive.NewObjectID()

	reqBody := reader.CreateChapterCommentRequest{
		ChapterID: "invalid-id",
		BookID:    bookID.Hex(),
		Content:   "Test comment",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/chapters/invalid-id/comments", strings.NewReader(string(jsonBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChapterCommentAPI_CreateChapterComment_InvalidRating(t *testing.T) {
	router := setupCommentTestRouter()

	chapterID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()

	reqBody := reader.CreateChapterCommentRequest{
		ChapterID: chapterID.Hex(),
		BookID:    bookID.Hex(),
		Content:   "Test comment",
		Rating:    6, // Invalid rating
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/chapters/"+chapterID.Hex()+"/comments", strings.NewReader(string(jsonBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response shared.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Message, "评分必须在0-5之间")
}

func TestChapterCommentAPI_CreateChapterComment_EmptyContent(t *testing.T) {
	router := setupCommentTestRouter()

	chapterID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()

	reqBody := reader.CreateChapterCommentRequest{
		ChapterID: chapterID.Hex(),
		BookID:    bookID.Hex(),
		Content:   "", // Empty content
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/chapters/"+chapterID.Hex()+"/comments", strings.NewReader(string(jsonBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Note: The API currently doesn't validate content emptiness, returns 201
	// This test documents current behavior
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestChapterCommentAPI_CreateChapterComment_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	commentAPI := NewChapterCommentAPI()
	router.POST("/api/v1/reader/chapters/:chapterId/comments", commentAPI.CreateChapterComment)

	chapterID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()

	reqBody := reader.CreateChapterCommentRequest{
		ChapterID: chapterID.Hex(),
		BookID:    bookID.Hex(),
		Content:   "Test comment",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/chapters/"+chapterID.Hex()+"/comments", strings.NewReader(string(jsonBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestChapterCommentAPI_CreateChapterComment_InvalidParentID(t *testing.T) {
	router := setupCommentTestRouter()

	chapterID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()
	invalidParentID := "invalid-parent-id"

	reqBody := reader.CreateChapterCommentRequest{
		ChapterID: chapterID.Hex(),
		BookID:    bookID.Hex(),
		Content:   "Test reply",
		ParentID:  &invalidParentID,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/chapters/"+chapterID.Hex()+"/comments", strings.NewReader(string(jsonBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response shared.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Message, "父评论ID格式无效")
}

// Test: GetChapterComment

func TestChapterCommentAPI_GetChapterComment_InvalidID(t *testing.T) {
	router := setupCommentTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/reader/comments/invalid-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response shared.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Message, "评论ID格式无效")
}

func TestChapterCommentAPI_GetChapterComment_EmptyID(t *testing.T) {
	router := setupCommentTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/reader/comments/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestChapterCommentAPI_GetChapterComment_NotFound(t *testing.T) {
	router := setupCommentTestRouter()

	commentID := primitive.NewObjectID()

	req, _ := http.NewRequest("GET", "/api/v1/reader/comments/"+commentID.Hex(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response shared.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Message, "评论不存在")
}

// Test: UpdateChapterComment

func TestChapterCommentAPI_UpdateChapterComment_Success(t *testing.T) {
	router := setupCommentTestRouter()

	commentID := primitive.NewObjectID()

	newContent := "Updated comment content"
	reqBody := reader.UpdateChapterCommentRequest{
		Content: &newContent,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/reader/comments/"+commentID.Hex(), strings.NewReader(string(jsonBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response shared.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Code)

	data := response.Data.(map[string]interface{})
	assert.Equal(t, commentID.Hex(), data["commentId"])
}

func TestChapterCommentAPI_UpdateChapterComment_WithRating(t *testing.T) {
	router := setupCommentTestRouter()

	commentID := primitive.NewObjectID()

	newRating := 4
	reqBody := reader.UpdateChapterCommentRequest{
		Rating: &newRating,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/reader/comments/"+commentID.Hex(), strings.NewReader(string(jsonBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestChapterCommentAPI_UpdateChapterComment_InvalidRating(t *testing.T) {
	router := setupCommentTestRouter()

	commentID := primitive.NewObjectID()

	invalidRating := 6
	reqBody := reader.UpdateChapterCommentRequest{
		Rating: &invalidRating,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/reader/comments/"+commentID.Hex(), strings.NewReader(string(jsonBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response shared.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Message, "评分必须在0-5之间")
}

func TestChapterCommentAPI_UpdateChapterComment_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	commentAPI := NewChapterCommentAPI()
	router.PUT("/api/v1/reader/comments/:commentId", commentAPI.UpdateChapterComment)

	commentID := primitive.NewObjectID()

	newContent := "Updated content"
	reqBody := reader.UpdateChapterCommentRequest{
		Content: &newContent,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/reader/comments/"+commentID.Hex(), strings.NewReader(string(jsonBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// Test: DeleteChapterComment

func TestChapterCommentAPI_DeleteChapterComment_Success(t *testing.T) {
	router := setupCommentTestRouter()

	commentID := primitive.NewObjectID()

	req, _ := http.NewRequest("DELETE", "/api/v1/reader/comments/"+commentID.Hex(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response shared.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Code)

	data := response.Data.(map[string]interface{})
	assert.Equal(t, commentID.Hex(), data["commentId"])
	assert.NotNil(t, data["userId"])
}

func TestChapterCommentAPI_DeleteChapterComment_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	commentAPI := NewChapterCommentAPI()
	router.DELETE("/api/v1/reader/comments/:commentId", commentAPI.DeleteChapterComment)

	commentID := primitive.NewObjectID()

	req, _ := http.NewRequest("DELETE", "/api/v1/reader/comments/"+commentID.Hex(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// Test: LikeChapterComment

func TestChapterCommentAPI_LikeChapterComment_Success(t *testing.T) {
	router := setupCommentTestRouter()

	commentID := primitive.NewObjectID()

	req, _ := http.NewRequest("POST", "/api/v1/reader/comments/"+commentID.Hex()+"/like", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response shared.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Contains(t, response.Message, "点赞成功")
}

func TestChapterCommentAPI_LikeChapterComment_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	commentAPI := NewChapterCommentAPI()
	router.POST("/api/v1/reader/comments/:commentId/like", commentAPI.LikeChapterComment)

	commentID := primitive.NewObjectID()

	req, _ := http.NewRequest("POST", "/api/v1/reader/comments/"+commentID.Hex()+"/like", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// Test: UnlikeChapterComment

func TestChapterCommentAPI_UnlikeChapterComment_Success(t *testing.T) {
	router := setupCommentTestRouter()

	commentID := primitive.NewObjectID()

	req, _ := http.NewRequest("DELETE", "/api/v1/reader/comments/"+commentID.Hex()+"/like", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response shared.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Message, "取消点赞成功")
}

func TestChapterCommentAPI_UnlikeChapterComment_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	commentAPI := NewChapterCommentAPI()
	router.DELETE("/api/v1/reader/comments/:commentId/like", commentAPI.UnlikeChapterComment)

	commentID := primitive.NewObjectID()

	req, _ := http.NewRequest("DELETE", "/api/v1/reader/comments/"+commentID.Hex()+"/like", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// Test: Paragraph Comments

func TestChapterCommentAPI_GetParagraphComments_Success(t *testing.T) {
	router := setupCommentTestRouter()

	chapterID := primitive.NewObjectID()
	paragraphIndex := 5

	req, _ := http.NewRequest("GET", "/api/v1/reader/chapters/"+chapterID.Hex()+"/paragraphs/"+strconv.Itoa(paragraphIndex)+"/comments", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response shared.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "获取成功", response.Message)

	// Parse Data field
	dataBytes, _ := json.Marshal(response.Data)
	var paragraphResp reader.ParagraphCommentResponse
	err = json.Unmarshal(dataBytes, &paragraphResp)
	assert.NoError(t, err)
	assert.Equal(t, paragraphIndex, paragraphResp.ParagraphIndex)
}

func TestChapterCommentAPI_GetParagraphComments_InvalidChapterID(t *testing.T) {
	router := setupCommentTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/reader/chapters/invalid-id/paragraphs/5/comments", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Note: The API currently doesn't validate chapterID format for this endpoint
	// It returns 200 with empty results
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestChapterCommentAPI_GetParagraphComments_InvalidParagraphIndex(t *testing.T) {
	router := setupCommentTestRouter()

	chapterID := primitive.NewObjectID()

	req, _ := http.NewRequest("GET", "/api/v1/reader/chapters/"+chapterID.Hex()+"/paragraphs/invalid/comments", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response shared.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Message, "段落索引格式无效")
}

func TestChapterCommentAPI_GetParagraphComments_NegativeParagraphIndex(t *testing.T) {
	router := setupCommentTestRouter()

	chapterID := primitive.NewObjectID()

	req, _ := http.NewRequest("GET", "/api/v1/reader/chapters/"+chapterID.Hex()+"/paragraphs/-1/comments", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChapterCommentAPI_CreateParagraphComment_Success(t *testing.T) {
	router := setupCommentTestRouter()

	chapterID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()
	paragraphIndex := 3
	charStart := 0
	charEnd := 50

	reqBody := reader.CreateChapterCommentRequest{
		ChapterID:      chapterID.Hex(),
		BookID:         bookID.Hex(),
		Content:        "This is a paragraph comment",
		ParagraphIndex: &paragraphIndex,
		CharStart:      &charStart,
		CharEnd:        &charEnd,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/chapters/"+chapterID.Hex()+"/paragraph-comments", strings.NewReader(string(jsonBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response shared.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data := response.Data.(map[string]interface{})
	commentData := data["comment"].(map[string]interface{})

	commentBytes, _ := json.Marshal(commentData)
	var comment reader.ChapterComment
	err = json.Unmarshal(commentBytes, &comment)
	assert.NoError(t, err)
	assert.NotNil(t, comment.ParagraphIndex)
	assert.Equal(t, paragraphIndex, *comment.ParagraphIndex)
}

func TestChapterCommentAPI_CreateParagraphComment_MissingParagraphIndex(t *testing.T) {
	router := setupCommentTestRouter()

	chapterID := primitive.NewObjectID()
	bookID := primitive.NewObjectID()

	reqBody := reader.CreateChapterCommentRequest{
		ChapterID: chapterID.Hex(),
		BookID:    bookID.Hex(),
		Content:   "This should fail - no paragraph index",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/chapters/"+chapterID.Hex()+"/paragraph-comments", strings.NewReader(string(jsonBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response shared.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Message, "段落索引不能为空")
}

func TestChapterCommentAPI_GetChapterParagraphComments_Success(t *testing.T) {
	router := setupCommentTestRouter()

	chapterID := primitive.NewObjectID()

	req, _ := http.NewRequest("GET", "/api/v1/reader/chapters/"+chapterID.Hex()+"/paragraph-comments", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response shared.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data := response.Data.(map[string]interface{})
	assert.Equal(t, chapterID.Hex(), data["chapterId"])
	assert.NotNil(t, data["paragraphStats"])
}

// Test: ChapterComment Model Methods

func TestChapterComment_IsParagraphComment(t *testing.T) {
	comment := &reader.ChapterComment{
		ParagraphIndex: func() *int { i := 5; return &i }(),
	}

	assert.True(t, comment.IsParagraphComment())

	comment2 := &reader.ChapterComment{}
	assert.False(t, comment2.IsParagraphComment())
}

func TestChapterComment_IsTopLevel(t *testing.T) {
	comment := &reader.ChapterComment{
		ParentID: nil,
	}

	assert.True(t, comment.IsTopLevel())

	parentID := "some-parent-id"
	comment2 := &reader.ChapterComment{
		ParentID: &parentID,
	}

	assert.False(t, comment2.IsTopLevel())
}

func TestChapterComment_CanEdit_WithinTime(t *testing.T) {
	comment := &reader.ChapterComment{
		CreatedAt: time.Now().Add(-10 * time.Minute),
		IsDeleted: false,
	}

	assert.True(t, comment.CanEdit())
}

func TestChapterComment_CanEdit_ExceedsTime(t *testing.T) {
	comment := &reader.ChapterComment{
		CreatedAt: time.Now().Add(-31 * time.Minute),
		IsDeleted: false,
	}

	assert.False(t, comment.CanEdit())
}

func TestChapterComment_CanEdit_Deleted(t *testing.T) {
	comment := &reader.ChapterComment{
		CreatedAt: time.Now(),
		IsDeleted: true,
	}

	assert.False(t, comment.CanEdit())
}

// Benchmark Tests

func BenchmarkChapterCommentAPI_GetChapterComments(b *testing.B) {
	router := setupCommentTestRouter()
	chapterID := primitive.NewObjectID()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/api/v1/reader/chapters/"+chapterID.Hex()+"/comments", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
