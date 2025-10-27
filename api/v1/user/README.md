# User API 模块结构说明

## 📁 文件结构

```
api/v1/user/
├── user_api.go      # 用户API处理器
├── user_dto.go      # 用户DTO定义
└── README.md        # 本文件
```

## 🎯 模块职责

**职责**: 用户认证和个人信息管理

**核心功能**:
- ✅ 用户注册
- ✅ 用户登录
- ✅ 获取个人信息
- ✅ 更新个人信息
- ✅ 修改密码

**注意**: 用户管理的管理员功能（查询所有用户、封禁用户等）已迁移到 `admin` 模块。

---

## 📋 API端点列表

### 公开端点（无需认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/register | 用户注册 |
| POST | /api/v1/login | 用户登录 |

### 认证端点（需要JWT Token）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/users/profile | 获取当前用户信息 |
| PUT | /api/v1/users/profile | 更新当前用户信息 |
| PUT | /api/v1/users/password | 修改密码 |

---

## 🔄 API调用流程

### 公开端点流程
```
客户端请求 
  → Router 
  → API Handler 
  → Service层 
  → Repository层 
  → 数据库
```

### 认证端点流程
```
客户端请求 
  → Router 
  → JWTAuth中间件（验证Token）
  → API Handler 
  → Service层 
  → Repository层 
  → 数据库
```

---

## 📊 请求/响应示例

### 1. 用户注册

**请求**:
```http
POST /api/v1/register
Content-Type: application/json

{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}
```

**响应**:
```json
{
  "code": 201,
  "message": "注册成功",
  "data": {
    "user_id": "670f1a2b3c4d5e6f7a8b9c0d",
    "username": "testuser",
    "email": "test@example.com",
    "role": "user",
    "status": "active",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**错误响应**:
```json
{
  "code": 400,
  "message": "注册失败",
  "error": "用户名已存在"
}
```

---

### 2. 用户登录

**请求**:
```http
POST /api/v1/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}
```

**响应**:
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "user_id": "670f1a2b3c4d5e6f7a8b9c0d",
    "username": "testuser",
    "email": "test@example.com",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**错误响应**:
```json
{
  "code": 401,
  "message": "用户名或密码错误"
}
```

---

### 3. 获取个人信息

**请求**:
```http
GET /api/v1/users/profile
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "user_id": "670f1a2b3c4d5e6f7a8b9c0d",
    "username": "testuser",
    "email": "test@example.com",
    "phone": "+8613800138000",
    "role": "user",
    "status": "active",
    "avatar": "https://example.com/avatar.jpg",
    "nickname": "测试用户",
    "bio": "这是我的个人简介",
    "email_verified": true,
    "phone_verified": false,
    "last_login_at": "2025-10-24T10:30:00Z",
    "last_login_ip": "192.168.1.100",
    "created_at": "2025-10-20T08:00:00Z",
    "updated_at": "2025-10-24T10:30:00Z"
  }
}
```

---

### 4. 更新个人信息

**请求**:
```http
PUT /api/v1/users/profile
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "nickname": "新昵称",
  "bio": "这是更新后的个人简介",
  "avatar": "https://example.com/new-avatar.jpg"
}
```

**响应**:
```json
{
  "code": 200,
  "message": "更新成功",
  "data": null
}
```

**说明**:
- 所有字段都是可选的，只更新提供的字段
- `nickname`: 昵称，最大50字符
- `bio`: 个人简介，最大500字符
- `avatar`: 头像URL，必须是有效的URL
- `phone`: 手机号，必须符合E.164格式

---

### 5. 修改密码

**请求**:
```http
PUT /api/v1/users/password
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "old_password": "old_password123",
  "new_password": "new_password456"
}
```

**响应**:
```json
{
  "code": 200,
  "message": "密码修改成功",
  "data": null
}
```

**错误响应**:
```json
{
  "code": 401,
  "message": "旧密码错误"
}
```

---

## 🛡️ 认证说明

### JWT Token

所有认证端点都需要在请求头中携带JWT Token：

```http
Authorization: Bearer <token>
```

### Token获取

通过注册或登录接口获取Token：
- 注册成功后会返回Token
- 登录成功后会返回Token

### Token过期

Token有效期由服务器配置决定，过期后需要重新登录获取新Token。

### 错误码

| 错误码 | 说明 |
|--------|------|
| 401 | Token无效或已过期 |
| 403 | 无权限访问 |

---

## 🔧 设计原则

### 1. 安全性优先

- 密码使用bcrypt加密存储
- 敏感操作需要验证旧密码
- Token采用JWT标准

### 2. 职责分离

- 用户模块只处理当前用户自己的信息
- 管理员功能独立到admin模块
- 清晰的权限边界

### 3. RESTful风格

- 使用标准HTTP方法
- 资源路径语义化
- 状态码使用规范

### 4. 统一响应格式

所有接口使用统一的响应格式：
```json
{
  "code": 200,
  "message": "操作描述",
  "data": {...}
}
```

### 5. 参数验证

使用binding标签进行参数验证：
- `required`: 必填
- `min/max`: 长度限制
- `email`: 邮箱格式
- `url`: URL格式

---

## 📝 开发规范

### 1. 命名规范

- **结构体**: UserAPI
- **构造函数**: NewUserAPI
- **方法名**: Register, Login, GetProfile等
- **请求DTO**: `<操作>Request`
- **响应DTO**: `<操作>Response`

### 2. 错误处理

```go
if err != nil {
    if serviceErr, ok := err.(*serviceInterfaces.ServiceError); ok {
        switch serviceErr.Type {
        case serviceInterfaces.ErrorTypeNotFound:
            shared.NotFound(c, "用户不存在")
        case serviceInterfaces.ErrorTypeUnauthorized:
            shared.Unauthorized(c, "认证失败")
        default:
            shared.InternalError(c, "操作失败", err)
        }
        return
    }
    shared.InternalError(c, "操作失败", err)
    return
}
```

### 3. 获取当前用户

从Context中获取当前登录用户ID（由JWT中间件设置）：
```go
userID, exists := c.Get("user_id")
if !exists {
    shared.Unauthorized(c, "未认证")
    return
}
```

### 4. Service层调用

```go
serviceReq := &serviceInterfaces.GetUserRequest{
    ID: userID.(string),
}

resp, err := api.userService.GetUser(c.Request.Context(), serviceReq)
```

---

## 🚀 扩展建议

### 未来可添加的功能

1. **邮箱验证**
   - 发送验证邮件
   - 验证邮箱
   - 重新发送验证邮件

2. **手机验证**
   - 发送验证短信
   - 验证手机号
   - 绑定/解绑手机号

3. **第三方登录**
   - 微信登录
   - QQ登录
   - 微博登录

4. **账号安全**
   - 登录日志查询
   - 在线设备管理
   - 异地登录提醒

5. **密码找回**
   - 通过邮箱找回
   - 通过手机找回
   - 安全问题验证

---

## 🔒 安全注意事项

### 1. 密码安全

- 密码必须至少6位
- 使用bcrypt加密
- 不返回密码哈希

### 2. Token安全

- Token存储在客户端安全位置
- HTTPS传输
- 定期刷新Token

### 3. 输入验证

- 所有用户输入都需要验证
- 防止SQL注入
- 防止XSS攻击

### 4. 频率限制

- 登录失败次数限制
- 注册频率限制
- 密码重置频率限制

---

## 📚 相关文档

- [用户管理API使用指南](../../../doc/api/用户管理API使用指南.md)
- [管理员用户管理API](../admin/README.md)
- [JWT中间件说明](../../../middleware/jwt.go)
- [架构设计规范](../../../doc/architecture/架构设计规范.md)

---

## 🆚 与Admin模块的区别

| 功能 | User模块 | Admin模块 |
|------|----------|-----------|
| 用户注册 | ✅ | ❌ |
| 用户登录 | ✅ | ❌ |
| 查看自己信息 | ✅ | ❌ |
| 修改自己信息 | ✅ | ❌ |
| 修改自己密码 | ✅ | ❌ |
| 查看所有用户 | ❌ | ✅ |
| 修改他人信息 | ❌ | ✅ |
| 封禁用户 | ❌ | ✅ |
| 删除用户 | ❌ | ✅ |

---

**版本**: v1.0  
**创建日期**: 2025-10-24  
**维护者**: User模块开发组

