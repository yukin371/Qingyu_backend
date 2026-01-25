package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// MessagingWSHub 消息WebSocket Hub
type MessagingWSHub struct {
	clients    map[string]*MessagingWSClient // key: userID
	mu         sync.RWMutex
	register   chan *MessagingWSClient
	unregister chan *MessagingWSClient
	broadcast  chan *MessageBroadcast
}

// MessagingWSClient 消息WebSocket客户端
type MessagingWSClient struct {
	ID     string
	UserID string
	Conn   *websocket.Conn
	Send   chan []byte
	Hub    *MessagingWSHub
}

// MessageBroadcast 消息广播
type MessageBroadcast struct {
	ConversationID string
	Message        interface{}
	ExcludeUserID  string // 排除某个用户（发送者）
}

// MessageWSMessage WebSocket消息格式
type MessageWSMessage struct {
	Type      string      `json:"type"` // new_message, read, typing, error
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

var messagingUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // TODO: 在生产环境中应该检查origin
	},
}

// NewMessagingWSHub 创建消息WebSocket Hub
func NewMessagingWSHub() *MessagingWSHub {
	hub := &MessagingWSHub{
		clients:    make(map[string]*MessagingWSClient),
		register:   make(chan *MessagingWSClient),
		unregister: make(chan *MessagingWSClient),
		broadcast:  make(chan *MessageBroadcast, 256),
	}
	go hub.Run()
	return hub
}

// Run 启动Hub
func (h *MessagingWSHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()
			log.Printf("用户 %s 连接到消息WebSocket", client.UserID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("用户 %s 断开消息WebSocket", client.UserID)

		case broadcast := <-h.broadcast:
			// broadcast从channel接收的是指针类型
			h.broadcastToConversation(broadcast)
		}
	}
}

// broadcastToConversation 向会话的所有参与者广播消息
func (h *MessagingWSHub) broadcastToConversation(broadcast *MessageBroadcast) {
	message := MessageWSMessage{
		Type:      "new_message",
		Data:      broadcast.Message,
		Timestamp: time.Now().Unix(),
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("消息序列化失败: %v", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for userID, client := range h.clients {
		// 排除发送者
		if userID == broadcast.ExcludeUserID {
			continue
		}

		// TODO: 检查用户是否在conversation中

		select {
		case client.Send <- data:
		default:
			close(client.Send)
			delete(h.clients, userID)
		}
	}
}

// SendMessage 发送消息到会话
func (h *MessagingWSHub) SendMessage(conversationID string, message interface{}, excludeUserID string) {
	broadcast := &MessageBroadcast{
		ConversationID: conversationID,
		Message:        message,
		ExcludeUserID:  excludeUserID,
	}
	h.broadcast <- broadcast
}

// HandleMessagingWebSocket WebSocket连接处理器
// @Summary WebSocket消息端点
// @Description 实时接收会话消息，需要在URL中传递token参数
// @Tags Social Messages
// @Param token query string true "JWT认证token"
// @Router /ws/messages [get]
func (h *MessagingWSHub) HandleMessagingWebSocket(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少token参数"})
		return
	}

	userID, err := validateMessagingToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
		return
	}

	conn, err := messagingUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}

	client := &MessagingWSClient{
		ID:     generateMessagingClientID(),
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		Hub:    h,
	}

	h.register <- client

	go client.writePump()
	go client.readPump()
}

// readPump 从WebSocket读取消息（处理typing状态等）
func (c *MessagingWSClient) readPump() {
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
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket读取错误: %v", err)
			}
			break
		}

		// 处理客户端消息（如typing状态）
		var wsMsg MessageWSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			continue
		}

		// 处理不同类型的消息
		switch wsMsg.Type {
		case "typing":
			// 广播typing状态给会话其他参与者
			// TODO: 实现typing广播
		}
	}
}

// writePump 向WebSocket写入消息
func (c *MessagingWSClient) writePump() {
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

// validateMessagingToken 验证JWT token（TODO: 实现真正的token验证）
func validateMessagingToken(token string) (string, error) {
	// TODO: 实现真正的JWT token验证
	// 这里应该调用认证服务验证token并返回userID
	return "user123", nil
}

// generateMessagingClientID 生成客户端ID
func generateMessagingClientID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString 生成随机字符串
func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[len(letters)%n]
	}
	return string(b)
}
