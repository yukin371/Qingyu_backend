package bookstore

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"Qingyu_backend/config"
	bookstoreModel "Qingyu_backend/models/bookstore"
	bookstoreService "Qingyu_backend/service/bookstore"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type publishedChapterRouteService struct {
	bookstoreService.ChapterService
}

func (s publishedChapterRouteService) GetPublishedChaptersByBookID(
	context.Context,
	string,
	int,
	int,
) ([]*bookstoreModel.Chapter, int64, error) {
	projectChapterID := "writer-document-1"
	return []*bookstoreModel.Chapter{{
		ProjectChapterID: &projectChapterID,
		Title:            "第一章",
	}}, 1, nil
}

func (s publishedChapterRouteService) GetFreeChaptersByBookID(
	context.Context,
	string,
	int,
	int,
) ([]*bookstoreModel.Chapter, int64, error) {
	return []*bookstoreModel.Chapter{{Title: "免费章节"}}, 1, nil
}

func (s publishedChapterRouteService) GetPaidChaptersByBookID(
	context.Context,
	string,
	int,
	int,
) ([]*bookstoreModel.Chapter, int64, error) {
	return []*bookstoreModel.Chapter{{Title: "付费章节"}}, 1, nil
}

func (s publishedChapterRouteService) GetFirstChapter(context.Context, string) (*bookstoreModel.Chapter, error) {
	return &bookstoreModel.Chapter{Title: "第一章"}, nil
}

func (s publishedChapterRouteService) GetLastChapter(context.Context, string) (*bookstoreModel.Chapter, error) {
	return &bookstoreModel.Chapter{Title: "最后一章"}, nil
}

func TestInitBookstoreRouterPublishedChaptersUsesBookIDParam(t *testing.T) {
	gin.SetMode(gin.TestMode)
	previousConfig := config.GlobalConfig
	config.GlobalConfig = &config.Config{JWT: &config.JWTConfig{Secret: "router-test-secret", ExpirationHours: 1}}
	t.Cleanup(func() { config.GlobalConfig = previousConfig })
	r := gin.New()
	api := r.Group("/api/v1")

	InitBookstoreRouter(
		api,
		nil,
		nil,
		nil,
		nil,
		publishedChapterRouteService{},
		nil,
		nil,
		nil,
		nil,
	)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/bookstore/books/507f1f77bcf86cd799439011/chapters/published", nil)
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Contains(t, recorder.Body.String(), `"projectChapterId":"writer-document-1"`)
	require.Contains(t, recorder.Body.String(), `"total":1`)
}

func TestInitBookstoreRouterFreeChaptersUsesBookIDParam(t *testing.T) {
	gin.SetMode(gin.TestMode)
	previousConfig := config.GlobalConfig
	config.GlobalConfig = &config.Config{JWT: &config.JWTConfig{Secret: "router-test-secret", ExpirationHours: 1}}
	t.Cleanup(func() { config.GlobalConfig = previousConfig })
	r := gin.New()
	api := r.Group("/api/v1")

	InitBookstoreRouter(
		api,
		nil,
		nil,
		nil,
		nil,
		publishedChapterRouteService{},
		nil,
		nil,
		nil,
		nil,
	)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/bookstore/books/507f1f77bcf86cd799439011/chapters/free", nil)
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Contains(t, recorder.Body.String(), `"title":"免费章节"`)
	require.Contains(t, recorder.Body.String(), `"total":1`)
}

func TestInitBookstoreRouterAdditionalChapterRoutesUseBookIDParam(t *testing.T) {
	gin.SetMode(gin.TestMode)
	previousConfig := config.GlobalConfig
	config.GlobalConfig = &config.Config{JWT: &config.JWTConfig{Secret: "router-test-secret", ExpirationHours: 1}}
	t.Cleanup(func() { config.GlobalConfig = previousConfig })
	r := gin.New()
	api := r.Group("/api/v1")
	InitBookstoreRouter(api, nil, nil, nil, nil, publishedChapterRouteService{}, nil, nil, nil, nil)

	tests := []struct {
		path  string
		value string
	}{
		{"/api/v1/bookstore/books/507f1f77bcf86cd799439011/chapters/paid", "付费章节"},
		{"/api/v1/bookstore/books/507f1f77bcf86cd799439011/chapters/first", "第一章"},
		{"/api/v1/bookstore/books/507f1f77bcf86cd799439011/chapters/last", "最后一章"},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)

			require.Equal(t, http.StatusOK, recorder.Code)
			require.Contains(t, recorder.Body.String(), `"title":"`+tt.value+`"`)
		})
	}
}

func TestInitBookstoreRouterWithAuthUsesInjectedMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1")
	var ratingService bookstoreService.BookRatingService = (*bookstoreService.BookRatingServiceImpl)(nil)

	InitBookstoreRouterWithAuth(
		api,
		nil,
		nil,
		ratingService,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		func(c *gin.Context) {
			c.AbortWithStatus(http.StatusTeapot)
		},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/bookstore/books/507f1f77bcf86cd799439011/rating", nil)
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusTeapot, recorder.Code)
}
