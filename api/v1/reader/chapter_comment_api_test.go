package reader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	readerModels "Qingyu_backend/models/reader"
)

func setupChapterCommentTestRouter(userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("userId", userID)
		}
		c.Next()
	})

	api := NewChapterCommentAPI()

	v1 := r.Group("/api/v1/reader")
	{
		v1.GET("/chapters/:chapterId/comments", api.GetChapterComments)
		v1.POST("/chapters/:chapterId/comments", api.CreateChapterComment)
		v1.GET("/comments/:commentId", api.GetChapterComment)
		v1.PUT("/comments/:commentId", api.UpdateChapterComment)
		v1.DELETE("/comments/:commentId", api.DeleteChapterComment)
		v1.POST("/comments/:commentId/like", api.LikeChapterComment)
		v1.DELETE("/comments/:commentId/like", api.UnlikeChapterComment)
		v1.GET("/chapters/:chapterId/paragraphs/:paragraphIndex/comments", api.GetParagraphComments)
		v1.POST("/chapters/:chapterId/paragraph-comments", api.CreateParagraphComment)
		v1.GET("/chapters/:chapterId/paragraph-comments", api.GetChapterParagraphComments)
	}

	return r
}

func TestChapterCommentAPI_GetChapterComments_Success(t *testing.T) {
	// Given
	router := setupChapterCommentTestRouter("")
	chapterID := primitive.NewObjectID().Hex()

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/chapters/"+chapterID+"/comments", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestChapterCommentAPI_GetChapterComments_WithPagination(t *testing.T) {
	// Given
	router := setupChapterCommentTestRouter("")
	chapterID := primitive.NewObjectID().Hex()

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/chapters/"+chapterID+"/comments?page=1&pageSize=20&sortBy=created_at&sortOrder=desc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestChapterCommentAPI_GetChapterComments_InvalidChapterID(t *testing.T) {
	// Given
	router := setupChapterCommentTestRouter("")

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/chapters/invalid-id/comments", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChapterCommentAPI_CreateChapterComment_Success(t *testing.T) {
	// Given
	userID := primitive.NewObjectID().Hex()
	router := setupChapterCommentTestRouter(userID)
	chapterID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()

	rating := 5
	reqBody := readerModels.CreateChapterCommentRequest{
		ChapterID: chapterID,
		BookID:    bookID,
		Content:   "这是一条很棒的评论",
		Rating:    rating,
		ParentID:  nil,
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/chapters/"+chapterID+"/comments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestChapterCommentAPI_CreateChapterComment_Unauthorized(t *testing.T) {
	// Given
	router := setupChapterCommentTestRouter("") // No userID
	chapterID := primitive.NewObjectID().Hex()

	reqBody := readerModels.CreateChapterCommentRequest{
		Content: "这是一条评论",
		Rating:  5,
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/chapters/"+chapterID+"/comments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestChapterCommentAPI_CreateChapterComment_InvalidRating(t *testing.T) {
	// Given
	userID := primitive.NewObjectID().Hex()
	router := setupChapterCommentTestRouter(userID)
	chapterID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()

	reqBody := readerModels.CreateChapterCommentRequest{
		ChapterID: chapterID,
		BookID:    bookID,
		Content:   "这是一条评论",
		Rating:    6, // Invalid rating > 5
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/chapters/"+chapterID+"/comments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChapterCommentAPI_CreateChapterComment_WithParent(t *testing.T) {
	// Given
	userID := primitive.NewObjectID().Hex()
	router := setupChapterCommentTestRouter(userID)
	chapterID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	parentID := primitive.NewObjectID().Hex()

	reqBody := readerModels.CreateChapterCommentRequest{
		ChapterID: chapterID,
		BookID:    bookID,
		Content:   "这是回复评论",
		Rating:    0,
		ParentID:  &parentID,
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/chapters/"+chapterID+"/comments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestChapterCommentAPI_GetChapterComment_NotFound(t *testing.T) {
	// Given
	router := setupChapterCommentTestRouter("")
	commentID := primitive.NewObjectID().Hex()

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/comments/"+commentID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestChapterCommentAPI_GetChapterComment_InvalidID(t *testing.T) {
	// Given
	router := setupChapterCommentTestRouter("")

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/comments/invalid-id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChapterCommentAPI_UpdateChapterComment_Success(t *testing.T) {
	// Given
	userID := primitive.NewObjectID().Hex()
	router := setupChapterCommentTestRouter(userID)
	commentID := primitive.NewObjectID().Hex()

	newContent := "更新后的评论内容"
	reqBody := readerModels.UpdateChapterCommentRequest{
		Content: &newContent,
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/reader/comments/"+commentID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestChapterCommentAPI_UpdateChapterComment_Unauthorized(t *testing.T) {
	// Given
	router := setupChapterCommentTestRouter("") // No userID
	commentID := primitive.NewObjectID().Hex()

	reqBody := readerModels.UpdateChapterCommentRequest{}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/reader/comments/"+commentID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestChapterCommentAPI_UpdateChapterComment_InvalidRating(t *testing.T) {
	// Given
	userID := primitive.NewObjectID().Hex()
	router := setupChapterCommentTestRouter(userID)
	commentID := primitive.NewObjectID().Hex()

	invalidRating := 10
	reqBody := readerModels.UpdateChapterCommentRequest{
		Rating: &invalidRating,
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/reader/comments/"+commentID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChapterCommentAPI_DeleteChapterComment_Success(t *testing.T) {
	// Given
	userID := primitive.NewObjectID().Hex()
	router := setupChapterCommentTestRouter(userID)
	commentID := primitive.NewObjectID().Hex()

	// When
	req, _ := http.NewRequest("DELETE", "/api/v1/reader/comments/"+commentID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestChapterCommentAPI_DeleteChapterComment_Unauthorized(t *testing.T) {
	// Given
	router := setupChapterCommentTestRouter("") // No userID

	// When
	req, _ := http.NewRequest("DELETE", "/api/v1/reader/comments/"+primitive.NewObjectID().Hex(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestChapterCommentAPI_LikeChapterComment_Success(t *testing.T) {
	// Given
	userID := primitive.NewObjectID().Hex()
	router := setupChapterCommentTestRouter(userID)
	commentID := primitive.NewObjectID().Hex()

	// When
	req, _ := http.NewRequest("POST", "/api/v1/reader/comments/"+commentID+"/like", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestChapterCommentAPI_LikeChapterComment_Unauthorized(t *testing.T) {
	// Given
	router := setupChapterCommentTestRouter("") // No userID

	// When
	req, _ := http.NewRequest("POST", "/api/v1/reader/comments/"+primitive.NewObjectID().Hex()+"/like", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestChapterCommentAPI_UnlikeChapterComment_Success(t *testing.T) {
	// Given
	userID := primitive.NewObjectID().Hex()
	router := setupChapterCommentTestRouter(userID)
	commentID := primitive.NewObjectID().Hex()

	// When
	req, _ := http.NewRequest("DELETE", "/api/v1/reader/comments/"+commentID+"/like", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestChapterCommentAPI_UnlikeChapterComment_Unauthorized(t *testing.T) {
	// Given
	router := setupChapterCommentTestRouter("") // No userID

	// When
	req, _ := http.NewRequest("DELETE", "/api/v1/reader/comments/"+primitive.NewObjectID().Hex()+"/like", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestChapterCommentAPI_GetParagraphComments_Success(t *testing.T) {
	// Given
	router := setupChapterCommentTestRouter("")
	chapterID := primitive.NewObjectID().Hex()
	paragraphIndex := 5

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/chapters/"+chapterID+"/paragraphs/"+fmt.Sprint(paragraphIndex)+"/comments", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestChapterCommentAPI_GetParagraphComments_InvalidIndex(t *testing.T) {
	// Given
	router := setupChapterCommentTestRouter("")
	chapterID := primitive.NewObjectID().Hex()

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/chapters/"+chapterID+"/paragraphs/invalid/comments", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChapterCommentAPI_CreateParagraphComment_Success(t *testing.T) {
	// Given
	userID := primitive.NewObjectID().Hex()
	router := setupChapterCommentTestRouter(userID)
	chapterID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	paragraphIndex := 10
	charStart := 0
	charEnd := 50

	reqBody := readerModels.CreateChapterCommentRequest{
		ChapterID:     chapterID,
		BookID:        bookID,
		Content:       "这段写得很好",
		ParagraphIndex: &paragraphIndex,
		CharStart:     &charStart,
		CharEnd:       &charEnd,
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/chapters/"+chapterID+"/paragraph-comments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestChapterCommentAPI_CreateParagraphComment_MissingParagraphIndex(t *testing.T) {
	// Given
	userID := primitive.NewObjectID().Hex()
	router := setupChapterCommentTestRouter(userID)
	chapterID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()

	reqBody := readerModels.CreateChapterCommentRequest{
		ChapterID: chapterID,
		BookID:    bookID,
		Content:   "段落评论",
		ParagraphIndex: nil, // Missing
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/chapters/"+chapterID+"/paragraph-comments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChapterCommentAPI_GetChapterParagraphComments_Success(t *testing.T) {
	// Given
	router := setupChapterCommentTestRouter("")
	chapterID := primitive.NewObjectID().Hex()

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/chapters/"+chapterID+"/paragraph-comments", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}
