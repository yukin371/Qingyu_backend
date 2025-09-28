# WebSocket实时通信设计

## 1. 需求概述
- **功能描述**：设计基于WebSocket的实时通信系统，支持即时消息、在线状态、实时通知等功能
- **业务价值**：提升用户交互体验，实现实时数据同步，增强应用的互动性
- **用户场景**：即时聊天、实时通知、在线状态显示、协同编辑、实时数据推送
- **功能边界**：连接管理、消息路由、房间管理、用户状态管理、消息持久化

## 2. 架构设计

### 2.1 整体架构
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Client    │    │  Mobile Client  │    │   Admin Panel   │
│   (WebSocket)   │    │   (WebSocket)   │    │   (WebSocket)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
┌─────────────────────────────────────────────────────────────────┐
│                    WebSocket Gateway                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │   Load      │  │ Connection  │  │   Message   │             │
│  │  Balancer   │  │   Manager   │  │   Router    │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
└─────────────────────────────────────────────────────────────────┘
                                 │
┌─────────────────────────────────────────────────────────────────┐
│                  WebSocket Service Cluster                     │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │   Node 1    │  │   Node 2    │  │   Node 3    │             │
│  │ (WebSocket) │  │ (WebSocket) │  │ (WebSocket) │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
└─────────────────────────────────────────────────────────────────┘
                                 │
┌─────────────────────────────────────────────────────────────────┐
│                      Service Layer                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │   Message   │  │    Room     │  │    User     │             │
│  │   Service   │  │   Service   │  │   Service   │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
└─────────────────────────────────────────────────────────────────┘
                                 │
┌─────────────────────────────────────────────────────────────────┐
│                    Repository Layer                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │Message Repo │  │  Room Repo  │  │  User Repo  │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
└─────────────────────────────────────────────────────────────────┘
                                 │
┌─────────────────────────────────────────────────────────────────┐
│                      Data Layer                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │   MongoDB   │  │    Redis    │  │  Message    │             │
│  │ (Messages)  │  │(Connections)│  │   Queue     │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
└─────────────────────────────────────────────────────────────────┘
```

### 2.2 模块划分
- **连接管理模块**：WebSocket连接建立、维护、断开处理
- **消息路由模块**：消息分发、广播、点对点传输
- **房间管理模块**：聊天室创建、用户加入/离开、权限管理
- **用户状态模块**：在线状态、心跳检测、用户信息同步
- **消息持久化模块**：消息存储、历史记录、离线消息

### 2.3 数据流向
```
客户端连接 → 身份验证 → 连接注册 → 消息接收 → 消息路由 → 目标推送 → 状态更新
```

### 2.4 技术选型
- **WebSocket框架**：Gorilla WebSocket (Go)
- **消息队列**：Redis Pub/Sub / Apache Kafka
- **连接存储**：Redis Cluster
- **消息存储**：MongoDB
- **负载均衡**：Nginx / HAProxy
- **服务发现**：Consul / etcd
- **监控工具**：Prometheus + Grafana

## 3. 详细设计

### 3.1 WebSocket连接管理
```go
type ConnectionManager struct {
    connections map[string]*Connection
    register    chan *Connection
    unregister  chan *Connection
    broadcast   chan []byte
    mutex       sync.RWMutex
}

type Connection struct {
    ID       string
    UserID   string
    RoomID   string
    Conn     *websocket.Conn
    Send     chan []byte
    Hub      *ConnectionManager
    LastPing time.Time
}

func (cm *ConnectionManager) Run() {
    for {
        select {
        case conn := <-cm.register:
            cm.registerConnection(conn)
        case conn := <-cm.unregister:
            cm.unregisterConnection(conn)
        case message := <-cm.broadcast:
            cm.broadcastMessage(message)
        }
    }
}

func (c *Connection) ReadPump() {
    defer func() {
        c.Hub.unregister <- c
        c.Conn.Close()
    }()
    
    c.Conn.SetReadLimit(maxMessageSize)
    c.Conn.SetReadDeadline(time.Now().Add(pongWait))
    c.Conn.SetPongHandler(func(string) error {
        c.Conn.SetReadDeadline(time.Now().Add(pongWait))
        c.LastPing = time.Now()
        return nil
    })
    
    for {
        _, message, err := c.Conn.ReadMessage()
        if err != nil {
            break
        }
        
        // 处理接收到的消息
        c.handleMessage(message)
    }
}

func (c *Connection) WritePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.Conn.Close()
    }()
    
    for {
        select {
        case message, ok := <-c.Send:
            if !ok {
                c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            
            c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
                return
            }
            
        case <-ticker.C:
            c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}
```

### 3.2 消息协议设计
```go
type WSMessage struct {
    Type      MessageType            `json:"type"`
    ID        string                 `json:"id"`
    From      string                 `json:"from"`
    To        string                 `json:"to,omitempty"`
    RoomID    string                 `json:"room_id,omitempty"`
    Content   interface{}            `json:"content"`
    Timestamp time.Time              `json:"timestamp"`
    Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type MessageType string
const (
    // 系统消息
    MessageTypeConnect    MessageType = "connect"
    MessageTypeDisconnect MessageType = "disconnect"
    MessageTypePing       MessageType = "ping"
    MessageTypePong       MessageType = "pong"
    
    // 聊天消息
    MessageTypeChat       MessageType = "chat"
    MessageTypePrivate    MessageType = "private"
    MessageTypeGroup      MessageType = "group"
    
    // 状态消息
    MessageTypeUserOnline  MessageType = "user_online"
    MessageTypeUserOffline MessageType = "user_offline"
    MessageTypeTyping      MessageType = "typing"
    
    // 房间消息
    MessageTypeJoinRoom    MessageType = "join_room"
    MessageTypeLeaveRoom   MessageType = "leave_room"
    MessageTypeRoomUpdate  MessageType = "room_update"
    
    // 通知消息
    MessageTypeNotification MessageType = "notification"
    MessageTypeSystem       MessageType = "system"
)

// 消息内容类型
type ChatContent struct {
    Text        string                 `json:"text,omitempty"`
    Image       string                 `json:"image,omitempty"`
    File        string                 `json:"file,omitempty"`
    MessageType string                 `json:"message_type"`
    ReplyTo     string                 `json:"reply_to,omitempty"`
    Mentions    []string               `json:"mentions,omitempty"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type NotificationContent struct {
    Title   string                 `json:"title"`
    Body    string                 `json:"body"`
    Icon    string                 `json:"icon,omitempty"`
    Action  string                 `json:"action,omitempty"`
    Data    map[string]interface{} `json:"data,omitempty"`
}
```

### 3.3 Service层设计
```go
type WebSocketService interface {
    HandleConnection(w http.ResponseWriter, r *http.Request)
    SendMessage(ctx context.Context, message *WSMessage) error
    BroadcastToRoom(ctx context.Context, roomID string, message *WSMessage) error
    SendToUser(ctx context.Context, userID string, message *WSMessage) error
    GetOnlineUsers(ctx context.Context, roomID string) ([]string, error)
    JoinRoom(ctx context.Context, userID, roomID string) error
    LeaveRoom(ctx context.Context, userID, roomID string) error
}

type webSocketServiceImpl struct {
    connManager    *ConnectionManager
    messageRepo    repository.MessageRepository
    roomRepo       repository.RoomRepository
    userRepo       repository.UserRepository
    redisClient    *redis.Client
    messageQueue   MessageQueue
    logger         *zap.Logger
}

func (s *webSocketServiceImpl) HandleConnection(w http.ResponseWriter, r *http.Request) {
    // 升级HTTP连接为WebSocket
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        s.logger.Error("Failed to upgrade connection", zap.Error(err))
        return
    }
    
    // 验证用户身份
    userID, err := s.authenticateUser(r)
    if err != nil {
        conn.Close()
        return
    }
    
    // 创建连接对象
    connection := &Connection{
        ID:       generateConnectionID(),
        UserID:   userID,
        Conn:     conn,
        Send:     make(chan []byte, 256),
        Hub:      s.connManager,
        LastPing: time.Now(),
    }
    
    // 注册连接
    s.connManager.register <- connection
    
    // 启动读写协程
    go connection.WritePump()
    go connection.ReadPump()
    
    // 发送连接成功消息
    welcomeMsg := &WSMessage{
        Type:      MessageTypeConnect,
        ID:        generateMessageID(),
        Content:   map[string]interface{}{"status": "connected"},
        Timestamp: time.Now(),
    }
    s.sendToConnection(connection, welcomeMsg)
}

func (s *webSocketServiceImpl) SendMessage(ctx context.Context, message *WSMessage) error {
    // 消息验证
    if err := s.validateMessage(message); err != nil {
        return err
    }
    
    // 持久化消息
    if s.shouldPersistMessage(message.Type) {
        if err := s.messageRepo.Create(ctx, s.convertToMessageModel(message)); err != nil {
            s.logger.Error("Failed to persist message", zap.Error(err))
        }
    }
    
    // 路由消息
    switch message.Type {
    case MessageTypeChat, MessageTypeGroup:
        return s.BroadcastToRoom(ctx, message.RoomID, message)
    case MessageTypePrivate:
        return s.SendToUser(ctx, message.To, message)
    case MessageTypeNotification:
        return s.SendToUser(ctx, message.To, message)
    default:
        return s.handleSystemMessage(ctx, message)
    }
}

func (s *webSocketServiceImpl) BroadcastToRoom(ctx context.Context, roomID string, message *WSMessage) error {
    // 获取房间内的所有连接
    connections := s.connManager.GetRoomConnections(roomID)
    
    messageBytes, err := json.Marshal(message)
    if err != nil {
        return err
    }
    
    // 并发发送消息
    var wg sync.WaitGroup
    for _, conn := range connections {
        wg.Add(1)
        go func(c *Connection) {
            defer wg.Done()
            select {
            case c.Send <- messageBytes:
            default:
                // 连接阻塞，关闭连接
                close(c.Send)
                s.connManager.unregister <- c
            }
        }(conn)
    }
    wg.Wait()
    
    return nil
}
```

### 3.4 Repository层设计
```go
type MessageRepository interface {
    Create(ctx context.Context, message *model.Message) error
    GetByID(ctx context.Context, id string) (*model.Message, error)
    GetRoomMessages(ctx context.Context, roomID string, pagination *Pagination) (*PagedResult[*model.Message], error)
    GetUserMessages(ctx context.Context, userID string, pagination *Pagination) (*PagedResult[*model.Message], error)
    MarkAsRead(ctx context.Context, messageID, userID string) error
    GetUnreadCount(ctx context.Context, userID string) (int64, error)
}

type RoomRepository interface {
    Create(ctx context.Context, room *model.Room) error
    GetByID(ctx context.Context, id string) (*model.Room, error)
    Update(ctx context.Context, room *model.Room) error
    Delete(ctx context.Context, id string) error
    GetUserRooms(ctx context.Context, userID string) ([]*model.Room, error)
    AddMember(ctx context.Context, roomID, userID string) error
    RemoveMember(ctx context.Context, roomID, userID string) error
    GetMembers(ctx context.Context, roomID string) ([]string, error)
}

type ConnectionRepository interface {
    SaveConnection(ctx context.Context, conn *model.Connection) error
    GetConnection(ctx context.Context, connID string) (*model.Connection, error)
    GetUserConnections(ctx context.Context, userID string) ([]*model.Connection, error)
    GetRoomConnections(ctx context.Context, roomID string) ([]*model.Connection, error)
    RemoveConnection(ctx context.Context, connID string) error
    UpdateLastSeen(ctx context.Context, connID string) error
}
```

### 3.5 Model层设计
```go
type Message struct {
    ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
    Type      string                 `bson:"type" json:"type"`
    From      string                 `bson:"from" json:"from"`
    To        string                 `bson:"to,omitempty" json:"to,omitempty"`
    RoomID    string                 `bson:"room_id,omitempty" json:"room_id,omitempty"`
    Content   map[string]interface{} `bson:"content" json:"content"`
    ReadBy    []ReadStatus           `bson:"read_by" json:"read_by"`
    CreatedAt time.Time              `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time              `bson:"updated_at" json:"updated_at"`
}

type ReadStatus struct {
    UserID string    `bson:"user_id" json:"user_id"`
    ReadAt time.Time `bson:"read_at" json:"read_at"`
}

type Room struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name        string             `bson:"name" json:"name"`
    Type        RoomType           `bson:"type" json:"type"`
    Description string             `bson:"description,omitempty" json:"description,omitempty"`
    Owner       string             `bson:"owner" json:"owner"`
    Members     []string           `bson:"members" json:"members"`
    Settings    RoomSettings       `bson:"settings" json:"settings"`
    CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type RoomSettings struct {
    IsPrivate     bool     `bson:"is_private" json:"is_private"`
    MaxMembers    int      `bson:"max_members" json:"max_members"`
    AllowedTypes  []string `bson:"allowed_types" json:"allowed_types"`
    MuteAll       bool     `bson:"mute_all" json:"mute_all"`
    RequireApproval bool   `bson:"require_approval" json:"require_approval"`
}

type Connection struct {
    ID        string    `json:"id"`
    UserID    string    `json:"user_id"`
    RoomID    string    `json:"room_id,omitempty"`
    ServerID  string    `json:"server_id"`
    Status    string    `json:"status"`
    LastSeen  time.Time `json:"last_seen"`
    CreatedAt time.Time `json:"created_at"`
}

type UserStatus struct {
    UserID    string    `json:"user_id"`
    Status    string    `json:"status"` // online, offline, away, busy
    LastSeen  time.Time `json:"last_seen"`
    Device    string    `json:"device"`
    Location  string    `json:"location,omitempty"`
}

// 枚举类型
type RoomType string
const (
    RoomTypePublic  RoomType = "public"
    RoomTypePrivate RoomType = "private"
    RoomTypeGroup   RoomType = "group"
    RoomTypeDirect  RoomType = "direct"
)
```

## 4. 数据设计

### 4.1 MongoDB数据模型
```javascript
// 消息集合
db.messages.createIndex({ "room_id": 1, "created_at": -1 })
db.messages.createIndex({ "from": 1, "created_at": -1 })
db.messages.createIndex({ "to": 1, "created_at": -1 })
db.messages.createIndex({ "type": 1 })

// 房间集合
db.rooms.createIndex({ "owner": 1 })
db.rooms.createIndex({ "members": 1 })
db.rooms.createIndex({ "type": 1 })

// 用户状态集合
db.user_status.createIndex({ "user_id": 1 }, { unique: true })
db.user_status.createIndex({ "status": 1 })
db.user_status.createIndex({ "last_seen": 1 })
```

### 4.2 Redis缓存设计
```
# 连接管理
connections:{server_id}:{user_id} -> connection_info (Hash)
user_connections:{user_id} -> [connection_ids] (Set)
room_connections:{room_id} -> [connection_ids] (Set)

# 用户状态
user_status:{user_id} -> status_info (Hash)
online_users -> [user_ids] (Set)

# 消息队列
message_queue:{room_id} -> messages (List)
user_message_queue:{user_id} -> messages (List)

# 限流控制
rate_limit:{user_id}:{action} -> count (String with TTL)
```

### 4.3 数据库索引策略
- **消息表**：按房间ID和时间复合索引，按用户ID和时间复合索引
- **房间表**：按所有者索引，按成员索引，按类型索引
- **连接表**：按用户ID索引，按房间ID索引，按服务器ID索引

## 5. 接口设计

### 5.1 WebSocket接口
```javascript
// 连接建立
ws://localhost:8080/ws?token={jwt_token}&room_id={room_id}

// 消息格式
{
  "type": "chat|private|notification|system",
  "id": "message_id",
  "from": "user_id",
  "to": "user_id", // 私聊时使用
  "room_id": "room_id", // 群聊时使用
  "content": {
    "text": "消息内容",
    "message_type": "text|image|file",
    "reply_to": "message_id"
  },
  "timestamp": "2024-01-01T00:00:00Z"
}

// 系统消息
{
  "type": "user_online|user_offline|join_room|leave_room",
  "content": {
    "user_id": "user_id",
    "room_id": "room_id",
    "status": "online|offline"
  }
}
```

### 5.2 HTTP管理接口
```go
// 房间管理
POST   /api/v1/rooms                    // 创建房间
GET    /api/v1/rooms/{id}               // 获取房间信息
PUT    /api/v1/rooms/{id}               // 更新房间信息
DELETE /api/v1/rooms/{id}               // 删除房间
POST   /api/v1/rooms/{id}/members       // 添加成员
DELETE /api/v1/rooms/{id}/members/{uid} // 移除成员

// 消息管理
GET    /api/v1/messages/rooms/{id}      // 获取房间消息历史
GET    /api/v1/messages/users/{id}      // 获取用户消息
POST   /api/v1/messages/{id}/read       // 标记消息已读
GET    /api/v1/messages/unread/count    // 获取未读消息数

// 连接管理
GET    /api/v1/connections/online       // 获取在线用户
GET    /api/v1/connections/rooms/{id}   // 获取房间在线用户
POST   /api/v1/connections/broadcast    // 广播消息
```

### 5.3 内部服务接口
```go
type MessageBroker interface {
    PublishToRoom(roomID string, message *WSMessage) error
    PublishToUser(userID string, message *WSMessage) error
    Subscribe(pattern string, handler MessageHandler) error
}

type ConnectionRegistry interface {
    RegisterConnection(conn *Connection) error
    UnregisterConnection(connID string) error
    GetUserConnections(userID string) ([]*Connection, error)
    GetRoomConnections(roomID string) ([]*Connection, error)
}
```

## 6. 安全设计

### 6.1 身份认证
- **JWT Token验证**：WebSocket连接建立时验证JWT token
- **Token刷新机制**：定期刷新token，防止token过期
- **权限验证**：验证用户对房间的访问权限

### 6.2 数据安全
- **消息加密**：敏感消息内容加密存储
- **传输安全**：使用WSS (WebSocket Secure) 协议
- **数据脱敏**：日志中敏感信息脱敏处理

### 6.3 防护措施
```go
// 限流控制
type RateLimiter struct {
    redis       *redis.Client
    maxRequests int
    window      time.Duration
}

func (r *RateLimiter) Allow(userID, action string) bool {
    key := fmt.Sprintf("rate_limit:%s:%s", userID, action)
    count, err := r.redis.Incr(context.Background(), key).Result()
    if err != nil {
        return false
    }
    
    if count == 1 {
        r.redis.Expire(context.Background(), key, r.window)
    }
    
    return count <= int64(r.maxRequests)
}

// 消息过滤
type MessageFilter struct {
    sensitiveWords []string
    maxLength      int
}

func (f *MessageFilter) Filter(content string) (string, error) {
    if len(content) > f.maxLength {
        return "", errors.New("message too long")
    }
    
    for _, word := range f.sensitiveWords {
        if strings.Contains(content, word) {
            content = strings.ReplaceAll(content, word, "***")
        }
    }
    
    return content, nil
}
```

## 7. 测试设计

### 7.1 单元测试
```go
func TestConnectionManager_RegisterConnection(t *testing.T) {
    cm := NewConnectionManager()
    conn := &Connection{
        ID:     "test-conn-1",
        UserID: "user-1",
        RoomID: "room-1",
    }
    
    cm.register <- conn
    time.Sleep(100 * time.Millisecond)
    
    assert.Contains(t, cm.connections, conn.ID)
}

func TestWebSocketService_SendMessage(t *testing.T) {
    service := setupTestService()
    message := &WSMessage{
        Type:    MessageTypeChat,
        From:    "user-1",
        RoomID:  "room-1",
        Content: map[string]interface{}{"text": "hello"},
    }
    
    err := service.SendMessage(context.Background(), message)
    assert.NoError(t, err)
}
```

### 7.2 集成测试
```go
func TestWebSocketIntegration(t *testing.T) {
    // 启动测试服务器
    server := httptest.NewServer(setupTestHandler())
    defer server.Close()
    
    // 建立WebSocket连接
    wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
    conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
    require.NoError(t, err)
    defer conn.Close()
    
    // 发送消息
    message := WSMessage{
        Type:    MessageTypeChat,
        Content: map[string]interface{}{"text": "test message"},
    }
    err = conn.WriteJSON(message)
    require.NoError(t, err)
    
    // 验证响应
    var response WSMessage
    err = conn.ReadJSON(&response)
    require.NoError(t, err)
    assert.Equal(t, MessageTypeChat, response.Type)
}
```

### 7.3 性能测试
```go
func BenchmarkMessageBroadcast(b *testing.B) {
    service := setupBenchmarkService()
    message := &WSMessage{
        Type:    MessageTypeChat,
        RoomID:  "test-room",
        Content: map[string]interface{}{"text": "benchmark"},
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        service.BroadcastToRoom(context.Background(), "test-room", message)
    }
}

// 并发连接测试
func TestConcurrentConnections(t *testing.T) {
    const numConnections = 1000
    var wg sync.WaitGroup
    
    for i := 0; i < numConnections; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            // 建立连接并发送消息
            testConnection(t, id)
        }(i)
    }
    
    wg.Wait()
}
```

## 8. 部署和运维

### 8.1 部署架构
```yaml
# docker-compose.yml
version: '3.8'
services:
  websocket-gateway:
    image: qingyu/websocket-gateway:latest
    ports:
      - "8080:8080"
    environment:
      - REDIS_URL=redis://redis:6379
      - MONGODB_URL=mongodb://mongo:27017
    depends_on:
      - redis
      - mongodb
    
  websocket-service:
    image: qingyu/websocket-service:latest
    deploy:
      replicas: 3
    environment:
      - REDIS_URL=redis://redis:6379
      - MONGODB_URL=mongodb://mongo:27017
    
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    
  mongodb:
    image: mongo:6
    ports:
      - "27017:27017"
```

### 8.2 监控告警
```go
// Prometheus指标
var (
    connectionsTotal = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "websocket_connections_total",
            Help: "Total number of WebSocket connections",
        },
        []string{"server_id"},
    )
    
    messagesTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "websocket_messages_total",
            Help: "Total number of messages processed",
        },
        []string{"type", "status"},
    )
    
    messageLatency = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "websocket_message_duration_seconds",
            Help: "Message processing duration",
        },
        []string{"type"},
    )
)

// 健康检查
func (s *WebSocketService) HealthCheck() error {
    // 检查Redis连接
    if err := s.redisClient.Ping(context.Background()).Err(); err != nil {
        return fmt.Errorf("redis health check failed: %w", err)
    }
    
    // 检查MongoDB连接
    if err := s.mongoClient.Ping(context.Background(), nil); err != nil {
        return fmt.Errorf("mongodb health check failed: %w", err)
    }
    
    return nil
}
```

### 8.3 日志管理
```go
// 结构化日志
type LogEntry struct {
    Level     string    `json:"level"`
    Timestamp time.Time `json:"timestamp"`
    Message   string    `json:"message"`
    UserID    string    `json:"user_id,omitempty"`
    RoomID    string    `json:"room_id,omitempty"`
    ConnID    string    `json:"connection_id,omitempty"`
    Action    string    `json:"action,omitempty"`
    Error     string    `json:"error,omitempty"`
}

func (s *WebSocketService) logMessage(level, message string, fields ...zap.Field) {
    switch level {
    case "info":
        s.logger.Info(message, fields...)
    case "warn":
        s.logger.Warn(message, fields...)
    case "error":
        s.logger.Error(message, fields...)
    }
}
```

## 9. 风险评估

### 9.1 技术风险
- **连接数限制**：单机连接数受限，需要水平扩展
- **消息丢失**：网络异常可能导致消息丢失，需要重试机制
- **内存泄漏**：连接未正确释放可能导致内存泄漏

### 9.2 业务风险
- **消息延迟**：高并发时消息可能出现延迟
- **数据一致性**：分布式环境下数据一致性问题
- **用户体验**：连接断开重连影响用户体验

### 9.3 运维风险
- **服务雪崩**：某个节点故障可能影响整体服务
- **数据备份**：消息数据需要定期备份
- **安全漏洞**：WebSocket协议可能存在安全风险

## 10. 实施计划

### 10.1 开发阶段 (4周)
- **第1周**：基础架构搭建，连接管理模块开发
- **第2周**：消息路由模块，房间管理模块开发
- **第3周**：用户状态模块，消息持久化模块开发
- **第4周**：安全模块，监控模块开发

### 10.2 测试阶段 (2周)
- **第1周**：单元测试，集成测试
- **第2周**：性能测试，压力测试，安全测试

### 10.3 上线阶段 (1周)
- **灰度发布**：先在小范围用户中测试
- **全量发布**：逐步扩大到所有用户
- **监控观察**：密切关注系统指标和用户反馈

---
*本文档为设计模板，需要根据实际需求进行详细设计*