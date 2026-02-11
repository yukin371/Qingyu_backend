package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"Qingyu_backend/service/auth"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket升级器配置
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 生产环境应该检查Origin
	},
	// 使用子协议传递token，避免在URL中暴露敏感信息
	Subprotocols: []string{"Bearer-Token"},
}

// WSHub WebSocket连接中心
type WSHub struct {
	clients    map[string]*WSClient
	register   chan *WSClient
	unregister chan *WSClient
	broadcast  chan *BroadcastMessage
	jwtService auth.JWTService
	mu         sync.RWMutex // 保护clients map的并发访问
}

// WSClient WebSocket客户端
type WSClient struct {
	ID     string
	UserID string
	Conn   *websocket.Conn
	Send   chan []byte
	Hub    *WSHub
}

// BroadcastMessage 广播消息
type BroadcastMessage struct {
	UserID  string
	Message interface{}
}

// NotificationMessage WebSocket通知消息格式
type NotificationMessage struct {
	Type      string      `json:"type"` // notification, ping, error
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

// NewWSHub 创建WebSocket Hub
func NewWSHub(jwtService auth.JWTService) *WSHub {
	hub := &WSHub{
		clients:    make(map[string]*WSClient),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
		broadcast:  make(chan *BroadcastMessage, 256),
		jwtService: jwtService,
	}
	go hub.Run()
	return hub
}

// HandleWebSocket WebSocket连接处理器
// @Summary WebSocket通知端点
// @Description 实时接收用户通知，通过WebSocket子协议传递token
// @Tags Notifications
// @Param token header string true "JWT认证token (通过Sec-WebSocket-Protocol子协议传递)"
// @Router /ws/notifications [get]
func (h *WSHub) HandleWebSocket(c *gin.Context) {
	// 从子协议中获取token（避免在URL中暴露）
	token := ""
	for _, subprotocol := range c.Request.Header["Sec-WebSocket-Protocol"] {
		if subprotocol != "Bearer-Token" && subprotocol != "" {
			token = subprotocol
			break
		}
	}

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少token参数，请通过Sec-WebSocket-Protocol传递"})
		return
	}

	// 验证token并获取userID
	userID, err := h.validateToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
		return
	}

	// 升级HTTP连接为WebSocket，选择子协议
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}

	// 创建客户端
	client := &WSClient{
		ID:     generateClientID(),
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		Hub:    h,
	}

	// 注册客户端
	h.register <- client

	// 启动读写协程
	go client.writePump()
	go client.readPump()
}

// validateToken 验证JWT token并返回userID
func (h *WSHub) validateToken(ctx context.Context, token string) (string, error) {
	if h.jwtService == nil {
		return "", fmt.Errorf("JWT服务未初始化")
	}

	// 验证token
	claims, err := h.jwtService.ValidateToken(ctx, token)
	if err != nil {
		return "", fmt.Errorf("token验证失败: %w", err)
	}

	// 返回userID
	return claims.UserID, nil
}

// generateClientID 生成客户端ID
func generateClientID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// readPump 从WebSocket连接读取消息
func (c *WSClient) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket读取错误: %v", err)
			}
			break
		}
	}
}

// writePump 向WebSocket连接写入消息
func (c *WSClient) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 排队消息
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Run 启动WebSocket Hub
func (h *WSHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.ID] = client
			h.mu.Unlock()
			log.Printf("客户端已注册: %s (用户: %s)", client.ID, client.UserID)
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("客户端已注销: %s (用户: %s)", client.ID, client.UserID)
		case message := <-h.broadcast:
			h.broadcastToUser(message.UserID, message.Message)
		}
	}
}

// broadcastToUser 向特定用户广播消息
func (h *WSHub) broadcastToUser(userID string, message interface{}) {
	// 第一阶段：在读锁下收集需要发送消息的客户端和需要删除的客户端ID
	var clientsToSend []*WSClient
	var clientIDsToDelete []string

	h.mu.RLock()
	for _, client := range h.clients {
		if client.UserID == userID {
			// 序列化消息
			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("消息序列化失败: %v", err)
				continue
			}

			// 尝试发送消息
			select {
			case client.Send <- data:
				clientsToSend = append(clientsToSend, client)
			default:
				// 发送通道已满，标记需要删除
				clientIDsToDelete = append(clientIDsToDelete, client.ID)
			}
		}
	}
	h.mu.RUnlock()

	// 第二阶段：如果有需要删除的客户端，使用写锁删除
	if len(clientIDsToDelete) > 0 {
		h.mu.Lock()
		for _, clientID := range clientIDsToDelete {
			if client, ok := h.clients[clientID]; ok {
				delete(h.clients, clientID)
				close(client.Send)
			}
		}
		h.mu.Unlock()
	}
}

// BroadcastNotification 广播通知给指定用户
func (h *WSHub) BroadcastNotification(userID string, notification interface{}) {
	msg := NotificationMessage{
		Type:      "notification",
		Data:      notification,
		Timestamp: time.Now().Unix(),
	}

	h.broadcast <- &BroadcastMessage{
		UserID:  userID,
		Message: msg,
	}
}

// GetConnectedCount 获取连接的客户端数量
func (h *WSHub) GetConnectedCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// GetUserConnections 获取指定用户的连接数
func (h *WSHub) GetUserConnections(userID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	count := 0
	for _, client := range h.clients {
		if client.UserID == userID {
			count++
		}
	}
	return count
}
