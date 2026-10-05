package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type documentSceneStateStub struct {
	id      string
	updates map[string]interface{}
}

func (s *documentSceneStateStub) Update(_ context.Context, id string, updates map[string]interface{}) error {
	s.id = id
	s.updates = updates
	return nil
}

func TestStoryWriteAPI_UpdateSceneStateUsesMinimalPort(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &documentSceneStateStub{}
	api := NewStoryWriteApi(nil, nil, stub)

	engine := gin.New()
	engine.PUT("/documents/:id/scene-state", api.UpdateSceneState)
	req := httptest.NewRequest(http.MethodPut, "/documents/doc-1/scene-state", bytes.NewBufferString(`{"sceneGoal":"目标","activeConflict":"冲突"}`))
	res := httptest.NewRecorder()
	engine.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, "doc-1", stub.id)
	require.Equal(t, "目标", stub.updates["scene_goal"])
	require.Equal(t, "冲突", stub.updates["active_conflict"])
	var payload map[string]interface{}
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &payload))
	require.Equal(t, float64(0), payload["code"])
}
