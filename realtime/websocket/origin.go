package websocket

import (
	"net/http"
	"strings"
	"sync"

	gorillawebsocket "github.com/gorilla/websocket"
)

var defaultAllowedOrigins = []string{
	"http://localhost:5173",
	"http://localhost:3000",
	"http://localhost:80",
	"http://localhost",
}

var (
	allowedOriginsMu sync.RWMutex
	allowedOrigins   = buildAllowedOrigins(defaultAllowedOrigins)
)

func buildAllowedOrigins(origins []string) map[string]bool {
	result := make(map[string]bool, len(origins))
	for _, origin := range origins {
		if normalized := normalizeOrigin(origin); normalized != "" {
			result[normalized] = true
		}
	}
	return result
}

func normalizeOrigin(origin string) string {
	return strings.TrimSpace(origin)
}

// SetAllowedOrigins 设置允许的 WebSocket origin 列表。
// 传入的列表会追加在本地默认白名单之上。
func SetAllowedOrigins(origins []string) {
	merged := append([]string{}, defaultAllowedOrigins...)
	merged = append(merged, origins...)

	allowedOriginsMu.Lock()
	allowedOrigins = buildAllowedOrigins(merged)
	allowedOriginsMu.Unlock()
}

// IsRequestOriginAllowed 检查 WebSocket 握手请求的 Origin 是否允许。
func IsRequestOriginAllowed(r *http.Request) bool {
	if r == nil {
		return false
	}

	origin := normalizeOrigin(r.Header.Get("Origin"))
	if origin == "" {
		// 非浏览器请求（如 curl、本地探活）不携带 Origin。
		return true
	}

	allowedOriginsMu.RLock()
	allowed := allowedOrigins[origin]
	allowedOriginsMu.RUnlock()
	return allowed
}

// NewUpgrader 创建带统一 Origin 校验的 WebSocket upgrader。
func NewUpgrader(subprotocols []string) gorillawebsocket.Upgrader {
	return gorillawebsocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     IsRequestOriginAllowed,
		Subprotocols:    subprotocols,
	}
}
