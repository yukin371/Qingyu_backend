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
	return []*bookstoreModel.Chapter{}, 0, nil
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
	)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/bookstore/books/507f1f77bcf86cd799439011/chapters/published", nil)
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)
}
