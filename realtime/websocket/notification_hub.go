package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

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
}

// WSHub WebSocket连接中心
type WSHub struct {
	clients    map[string]*WSClient
	register   chan *WSClient
	unregister chan *WSClient
	broadcast  chan *BroadcastMessage
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
func NewWSHub() *WSHub {
	hub := &WSHub{
		clients:    make(map[string]*WSClient),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
		broadcast:  make(chan *BroadcastMessage, 256),
	}
	go hub.Run()
	return hub
}

// HandleWebSocket WebSocket连接处理器
// @Summary WebSocket通知端点
// @Description 实时接收用户通知，需要在URL中传递token参数
// @Tags Notifications
// @Param token query string true "JWT认证token"
// @Router /ws/notifications [get]
func (h *WSHub) HandleWebSocket(c *gin.Context) {
	// 从query参数获取token
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少token参数"})
		return
	}

	// 验证token并获取userID
	userID, err := validateToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
		return
	}

	// 升级HTTP连接为WebSocket
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
func validateToken(ctx context.Context, token string) (string, error) {
	// TODO: 实现JWT验证逻辑
	// 这里应该调用JWT服务验证token
	// 返回userID或错误
	// 暂时返回模拟数据用于测试
	return "test_user_id", nil
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
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		if client.UserID == userID {
			// 序列化消息
			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("消息序列化失败: %v", err)
				continue
			}

			// 发送消息
			select {
			case client.Send <- data:
			default:
				// 发送通道已满，关闭连接
				close(client.Send)
				delete(h.clients, client.ID)
			}
		}
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
