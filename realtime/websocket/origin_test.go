package websocket

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsRequestOriginAllowed(t *testing.T) {
	SetAllowedOrigins(nil)

	t.Run("empty origin", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/ws/test", nil)
		assert.True(t, IsRequestOriginAllowed(req))
	})

	t.Run("default localhost origin", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/ws/test", nil)
		req.Header.Set("Origin", "http://localhost:5173")
		assert.True(t, IsRequestOriginAllowed(req))
	})

	t.Run("custom origin", func(t *testing.T) {
		SetAllowedOrigins([]string{"https://editor.example.com"})

		req, _ := http.NewRequest(http.MethodGet, "/ws/test", nil)
		req.Header.Set("Origin", "https://editor.example.com")
		assert.True(t, IsRequestOriginAllowed(req))
	})

	t.Run("blocked origin", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/ws/test", nil)
		req.Header.Set("Origin", "https://malicious.example.com")
		assert.False(t, IsRequestOriginAllowed(req))
	})
}
