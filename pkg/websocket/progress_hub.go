package websocket

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ProgressMessage 阅读进度消息
type ProgressMessage struct {
	Type      string    `json:"type"`       // sync, ack, conflict
	UserID    string    `json:"userId"`
	BookID    string    `json:"bookId"`
	ChapterID string    `json:"chapterId"`
	Progress  float64   `json:"progress"`   // 0-1
	Timestamp time.Time `json:"timestamp"`
	DeviceID  string    `json:"deviceId"`   // 设备标识
}

// Client WebSocket客户端
type Client struct {
	UserID   string
	DeviceID string
	Hub      *ProgressHub
	Conn     *websocket.Conn
	Send     chan *ProgressMessage
}

// ProgressHub 阅读进度同步中心
type ProgressHub struct {
	// 用户 -> 设备 -> 客户端连接
	clients map[string]map[string]*Client
	// 读写锁
	mu sync.RWMutex
	// 注册新客户端
	Register chan *Client
	// 注销客户端
	Unregister chan *Client
	// 广播消息
	Broadcast chan *ProgressMessage
}

// NewProgressHub 创建进度同步中心
func NewProgressHub() *ProgressHub {
	return &ProgressHub{
		clients:    make(map[string]map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *ProgressMessage, 256),
	}
}

// Run 运行Hub
func (h *ProgressHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)
		case client := <-h.Unregister:
			h.unregisterClient(client)
		case message := <-h.Broadcast:
			h.broadcastMessage(message)
		}
	}
}

// registerClient 注册新客户端
func (h *ProgressHub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[client.UserID] == nil {
		h.clients[client.UserID] = make(map[string]*Client)
	}

	// 如果同一设备已有连接，关闭旧连接
	if oldClient, exists := h.clients[client.UserID][client.DeviceID]; exists {
		oldClient.Conn.Close()
		close(oldClient.Send)
	}

	h.clients[client.UserID][client.DeviceID] = client
}

// unregisterClient 注销客户端
func (h *ProgressHub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client.UserID]; ok {
		if _, ok := h.clients[client.UserID][client.DeviceID]; ok {
			delete(h.clients[client.UserID], client.DeviceID)
			close(client.Send)

			// 如果用户没有其他设备连接，清理用户条目
			if len(h.clients[client.UserID]) == 0 {
				delete(h.clients, client.UserID)
			}
		}
	}
}

// broadcastMessage 广播消息给用户的所有其他设备
func (h *ProgressHub) broadcastMessage(message *ProgressMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	userClients, ok := h.clients[message.UserID]
	if !ok {
		return
	}

	for deviceID, client := range userClients {
		// 跳过发送消息的设备自己
		if deviceID == message.DeviceID {
			continue
		}

		select {
		case client.Send <- message:
		default:
			// 发送缓冲区满，关闭连接
			h.Unregister <- client
		}
	}
}

// SyncProgress 同步阅读进度
func (h *ProgressHub) SyncProgress(message *ProgressMessage) {
	h.Broadcast <- message
}

// GetConnectedDevices 获取用户已连接的设备列表
func (h *ProgressHub) GetConnectedDevices(userID string) []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	userClients, ok := h.clients[userID]
	if !ok {
		return []string{}
	}

	devices := make([]string, 0, len(userClients))
	for deviceID := range userClients {
		devices = append(devices, deviceID)
	}
	return devices
}

// ReadPump 从WebSocket连接读取消息
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
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
				fmt.Printf("WebSocket error: %v\n", err)
			}
			break
		}

		var progressMsg ProgressMessage
		if err := json.Unmarshal(message, &progressMsg); err != nil {
			fmt.Printf("Invalid message format: %v\n", err)
			continue
		}

		// 验证消息
		if progressMsg.UserID != c.UserID {
			continue
		}

		// 广播进度更新
		c.Hub.SyncProgress(&progressMsg)
	}
}

// WritePump 向WebSocket连接写入消息
func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// Hub关闭连接
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
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
