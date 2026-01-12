package reader

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"go.mongodb.org/mongo-driver/bson/primitive"

	readerModels "Qingyu_backend/models/reader"
)

func setupFontTestRouter(userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("userId", userID)
		}
		c.Next()
	})

	api := NewFontAPI()

	r.GET("/api/v1/reader/fonts", api.GetFonts)
	r.GET("/api/v1/reader/fonts/:name", api.GetFontByName)
	r.POST("/api/v1/reader/fonts", api.CreateCustomFont)
	r.PUT("/api/v1/reader/fonts/:id", api.UpdateFont)
	r.DELETE("/api/v1/reader/fonts/:id", api.DeleteFont)
	r.POST("/api/v1/reader/settings/font", api.SetFontPreference)

	return r
}

func TestFontAPI_GetFonts_Success(t *testing.T) {
	// Given
	router := setupFontTestRouter("")

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/fonts", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFontAPI_GetFonts_WithCategory(t *testing.T) {
	// Given
	router := setupFontTestRouter("")

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/fonts?category=serif", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFontAPI_GetFonts_BuiltinOnly(t *testing.T) {
	// Given
	router := setupFontTestRouter("")

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/fonts?builtin=true", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFontAPI_GetFontByName_Success(t *testing.T) {
	// Given
	router := setupFontTestRouter("")

	// When - "system-serif" is a built-in font
	req, _ := http.NewRequest("GET", "/api/v1/reader/fonts/system-serif", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFontAPI_GetFontByName_NotFound(t *testing.T) {
	// Given
	router := setupFontTestRouter("")

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/fonts/NonExistentFont", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestFontAPI_GetFontByName_EmptyName(t *testing.T) {
	// Given
	router := setupFontTestRouter("")

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/fonts/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - Gin returns 301 redirect for trailing slash without path param
	assert.Equal(t, http.StatusMovedPermanently, w.Code)
}

func TestFontAPI_CreateCustomFont_Success(t *testing.T) {
	// Given
	userID := primitive.NewObjectID().Hex()
	router := setupFontTestRouter(userID)

	reqBody := readerModels.CreateCustomFontRequest{
		Name:        "CustomFont",
		DisplayName: "自定义字体",
		FontFamily:  "CustomFamily",
		Category:    "custom",
		FontURL:     "https://example.com/font.ttf",
		PreviewText: "字体预览文本",
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/fonts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestFontAPI_CreateCustomFont_Unauthorized(t *testing.T) {
	// Given
	router := setupFontTestRouter("")

	reqBody := readerModels.CreateCustomFontRequest{
		Name: "CustomFont",
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/fonts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestFontAPI_UpdateFont_Success(t *testing.T) {
	// Given
	userID := primitive.NewObjectID().Hex()
	fontID := primitive.NewObjectID().Hex()
	router := setupFontTestRouter(userID)

	reqBody := readerModels.UpdateFontRequest{
		DisplayName: strPtr("更新字体名称"),
		Description: strPtr("更新描述"),
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/reader/fonts/"+fontID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFontAPI_UpdateFont_Unauthorized(t *testing.T) {
	// Given
	router := setupFontTestRouter("")

	reqBody := readerModels.UpdateFontRequest{
		DisplayName: strPtr("更新字体名称"),
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/reader/fonts/"+primitive.NewObjectID().Hex(), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestFontAPI_UpdateFont_EmptyFontID(t *testing.T) {
	// Given
	userID := primitive.NewObjectID().Hex()
	router := setupFontTestRouter(userID)

	reqBody := readerModels.UpdateFontRequest{
		DisplayName: strPtr("更新字体名称"),
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/reader/fonts/", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestFontAPI_DeleteFont_Success(t *testing.T) {
	// Given
	userID := primitive.NewObjectID().Hex()
	fontID := primitive.NewObjectID().Hex()
	router := setupFontTestRouter(userID)

	// When
	req, _ := http.NewRequest("DELETE", "/api/v1/reader/fonts/"+fontID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFontAPI_DeleteFont_Unauthorized(t *testing.T) {
	// Given
	router := setupFontTestRouter("")

	// When
	req, _ := http.NewRequest("DELETE", "/api/v1/reader/fonts/"+primitive.NewObjectID().Hex(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestFontAPI_SetFontPreference_Success(t *testing.T) {
	// Given
	userID := primitive.NewObjectID().Hex()
	router := setupFontTestRouter(userID)

	reqBody := readerModels.FontPreference{
		FontName: "system-serif",
		FontSize: 18,
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/settings/font", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFontAPI_SetFontPreference_FontNotFound(t *testing.T) {
	// Given
	userID := primitive.NewObjectID().Hex()
	router := setupFontTestRouter(userID)

	reqBody := readerModels.FontPreference{
		FontName: "NonExistentFont",
		FontSize: 18,
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/settings/font", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestFontAPI_SetFontPreference_Unauthorized(t *testing.T) {
	// Given
	router := setupFontTestRouter("")

	reqBody := readerModels.FontPreference{
		FontName: "Arial",
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/settings/font", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
