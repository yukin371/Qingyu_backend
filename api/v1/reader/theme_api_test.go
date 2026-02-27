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

func setupThemeTestRouter(userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		c.Next()
	})

	api := NewThemeAPI()

	r.GET("/api/v1/reader/themes", api.GetThemes)
	r.GET("/api/v1/reader/themes/:name", api.GetThemeByName)
	r.POST("/api/v1/reader/themes", api.CreateCustomTheme)
	r.PUT("/api/v1/reader/themes/:id", api.UpdateTheme)
	r.DELETE("/api/v1/reader/themes/:id", api.DeleteTheme)
	r.POST("/api/v1/reader/themes/:name/activate", api.ActivateTheme)

	return r
}

func TestThemeAPI_GetThemes_Success(t *testing.T) {
	// Given
	router := setupThemeTestRouter("")

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/themes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestThemeAPI_GetThemes_BuiltinOnly(t *testing.T) {
	// Given
	router := setupThemeTestRouter("")

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/themes?builtin=true", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestThemeAPI_GetThemes_PublicOnly(t *testing.T) {
	// Given
	router := setupThemeTestRouter("")

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/themes?public=true", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestThemeAPI_GetThemeByName_Success(t *testing.T) {
	// Given
	router := setupThemeTestRouter("")

	// When - "light" is a built-in theme
	req, _ := http.NewRequest("GET", "/api/v1/reader/themes/light", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestThemeAPI_GetThemeByName_NotFound(t *testing.T) {
	// Given
	router := setupThemeTestRouter("")

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/themes/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestThemeAPI_GetThemeByName_EmptyName(t *testing.T) {
	// Given
	router := setupThemeTestRouter("")

	// When
	req, _ := http.NewRequest("GET", "/api/v1/reader/themes/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - Gin returns 301 for trailing slash without path param
	assert.Equal(t, http.StatusMovedPermanently, w.Code)
}

func TestThemeAPI_CreateCustomTheme_Success(t *testing.T) {
	// Given
	userID := primitive.NewObjectID().Hex()
	router := setupThemeTestRouter(userID)

	colors := readerModels.ThemeColors{
		Background:          "#FFFFFF",
		SecondaryBackground: "#F5F5F5",
		TextPrimary:         "#000000",
		TextSecondary:       "#666666",
		TextDisabled:        "#CCCCCC",
		LinkColor:           "#3498db",
		AccentColor:         "#3498db",
		AccentHover:         "#2980b9",
		BorderColor:         "#ECF0F1",
		DividerColor:        "#BDC3C7",
		HighlightColor:      "#f1c40f",
		BookmarkColor:       "#f39c12",
		AnnotationColor:     "#e74c3c",
		ShadowColor:         "#bdc3c7",
	}

	reqBody := readerModels.CreateCustomThemeRequest{
		Name:        "custom-theme",
		DisplayName: "自定义主题",
		Description: "我的自定义主题",
		IsPublic:    false,
		Colors:      colors,
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/themes", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestThemeAPI_CreateCustomTheme_Unauthorized(t *testing.T) {
	// Given
	router := setupThemeTestRouter("") // No userID

	colors := readerModels.ThemeColors{
		Background:  "#FFFFFF",
		TextPrimary: "#000000",
	}

	reqBody := readerModels.CreateCustomThemeRequest{
		Name:   "custom-theme",
		Colors: colors,
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/reader/themes", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestThemeAPI_UpdateTheme_Success(t *testing.T) {
	// Given
	userID := primitive.NewObjectID().Hex()
	themeID := primitive.NewObjectID().Hex()
	router := setupThemeTestRouter(userID)

	reqBody := readerModels.UpdateThemeRequest{
		DisplayName: strPtr("更新后的主题名称"),
		Description: strPtr("更新后的描述"),
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/reader/themes/"+themeID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestThemeAPI_UpdateTheme_Unauthorized(t *testing.T) {
	// Given
	router := setupThemeTestRouter("") // No userID

	reqBody := readerModels.UpdateThemeRequest{
		DisplayName: strPtr("更新后的主题名称"),
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/reader/themes/"+primitive.NewObjectID().Hex(), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestThemeAPI_UpdateTheme_EmptyThemeID(t *testing.T) {
	// Given
	userID := primitive.NewObjectID().Hex()
	router := setupThemeTestRouter(userID)

	reqBody := readerModels.UpdateThemeRequest{
		DisplayName: strPtr("更新后的主题名称"),
	}

	// When
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/api/v1/reader/themes/", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestThemeAPI_DeleteTheme_Success(t *testing.T) {
	// Given
	userID := primitive.NewObjectID().Hex()
	themeID := primitive.NewObjectID().Hex()
	router := setupThemeTestRouter(userID)

	// When
	req, _ := http.NewRequest("DELETE", "/api/v1/reader/themes/"+themeID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestThemeAPI_DeleteTheme_Unauthorized(t *testing.T) {
	// Given
	router := setupThemeTestRouter("") // No userID

	// When
	req, _ := http.NewRequest("DELETE", "/api/v1/reader/themes/"+primitive.NewObjectID().Hex(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestThemeAPI_ActivateTheme_Success(t *testing.T) {
	// Given
	userID := primitive.NewObjectID().Hex()
	router := setupThemeTestRouter(userID)

	// When - "light" is a built-in theme
	req, _ := http.NewRequest("POST", "/api/v1/reader/themes/light/activate", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestThemeAPI_ActivateTheme_NotFound(t *testing.T) {
	// Given
	userID := primitive.NewObjectID().Hex()
	router := setupThemeTestRouter(userID)

	// When
	req, _ := http.NewRequest("POST", "/api/v1/reader/themes/nonexistent/activate", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestThemeAPI_ActivateTheme_Unauthorized(t *testing.T) {
	// Given
	router := setupThemeTestRouter("") // No userID

	// When
	req, _ := http.NewRequest("POST", "/api/v1/reader/themes/light/activate", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestThemeAPI_ActivateTheme_EmptyName(t *testing.T) {
	// Given
	userID := primitive.NewObjectID().Hex()
	router := setupThemeTestRouter(userID)

	// When
	req, _ := http.NewRequest("POST", "/api/v1/reader/themes//activate", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then - Empty name parameter returns 400 bad request
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
