package writer

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	writerModel "Qingyu_backend/models/writer"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type conceptCatalogStub struct {
	concepts []*writerModel.Concept
}

func (s *conceptCatalogStub) GetByID(context.Context, string) (*writerModel.Concept, error) {
	if len(s.concepts) == 0 {
		return nil, nil
	}
	return s.concepts[0], nil
}

func (s *conceptCatalogStub) ListByProject(context.Context, string) ([]*writerModel.Concept, error) {
	return s.concepts, nil
}

func (s *conceptCatalogStub) Search(context.Context, string, string, string) ([]*writerModel.Concept, error) {
	return s.concepts, nil
}

func TestEncyclopediaAPI_UsesConceptCatalogPort(t *testing.T) {
	gin.SetMode(gin.TestMode)
	concept := &writerModel.Concept{ID: "concept-1", ProjectID: "project-1", Name: "世界观"}
	api := NewEncyclopediaApi(&conceptCatalogStub{concepts: []*writerModel.Concept{concept}})

	engine := gin.New()
	engine.GET("/projects/:id/concepts", api.ListConcepts)

	req := httptest.NewRequest(http.MethodGet, "/projects/project-1/concepts", nil)
	res := httptest.NewRecorder()
	engine.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	var payload struct {
		Code int                    `json:"code"`
		Data []*writerModel.Concept `json:"data"`
	}
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &payload))
	require.Equal(t, 0, payload.Code)
	require.Len(t, payload.Data, 1)
	require.Equal(t, "concept-1", payload.Data[0].ID)
}
