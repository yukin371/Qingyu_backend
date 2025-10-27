# Admin API 模块结构说明

## 📁 文件结构

```
api/v1/admin/
├── user_admin_api.go      # 用户管理API
├── quota_admin_api.go     # AI配额管理API
├── audit_admin_api.go     # 审核管理API
├── system_admin_api.go    # 系统管理API
├── types.go               # 公共DTO定义
└── README.md              # 本文件
```

## 🎯 模块职责划分

### 1. UserAdminAPI (`user_admin_api.go`)

**职责**: 用户管理功能

**核心功能**:
- ✅ 获取用户列表（分页、筛选）
- ✅ 获取指定用户详情
- ✅ 更新用户信息（角色、状态等）
- ✅ 删除用户
- ✅ 封禁用户
- ✅ 解除封禁

**API端点**:
```
GET    /api/v1/admin/users              # 获取用户列表
GET    /api/v1/admin/users/:id          # 获取用户详情
PUT    /api/v1/admin/users/:id          # 更新用户信息
DELETE /api/v1/admin/users/:id          # 删除用户
POST   /api/v1/admin/users/:id/ban      # 封禁用户
POST   /api/v1/admin/users/:id/unban    # 解除封禁
```

**依赖服务**:
- `serviceInterfaces.UserService` - 用户服务

---

### 2. QuotaAdminAPI (`quota_admin_api.go`)

**职责**: AI配额管理功能

**核心功能**:
- ✅ 获取用户配额详情
- ✅ 更新用户配额
- ✅ 暂停用户配额
- ✅ 激活用户配额

**API端点**:
```
GET  /api/v1/admin/quota/:userId           # 获取用户配额详情
PUT  /api/v1/admin/quota/:userId           # 更新用户配额
POST /api/v1/admin/quota/:userId/suspend   # 暂停用户配额
POST /api/v1/admin/quota/:userId/activate  # 激活用户配额
```

**依赖服务**:
- `ai.QuotaService` - AI配额服务

---

### 3. AuditAdminAPI (`audit_admin_api.go`)

**职责**: 内容审核管理功能

**核心功能**:
- ✅ 获取待审核内容列表
- ✅ 审核内容（通过/拒绝）
- ✅ 审核申诉
- ✅ 获取高风险审核记录
- ✅ 获取审核统计

**API端点**:
```
GET  /api/v1/admin/audit/pending             # 获取待审核内容
GET  /api/v1/admin/audit/high-risk           # 获取高风险审核记录
GET  /api/v1/admin/audit/statistics          # 获取审核统计
POST /api/v1/admin/audit/:id/review          # 审核内容
POST /api/v1/admin/audit/:id/appeal/review   # 审核申诉
```

**依赖服务**:
- `interfaces.ContentAuditService` - 内容审核服务

---

### 4. SystemAdminAPI (`system_admin_api.go`)

**职责**: 系统管理功能

**核心功能**:
- ✅ 审核提现申请
- ✅ 获取用户统计
- ✅ 获取操作日志
- ✅ 获取系统统计
- ✅ 获取/更新系统配置
- ✅ 公告管理

**API端点**:
```
# 提现管理
POST /api/v1/admin/withdraw/review           # 审核提现

# 统计信息
GET  /api/v1/admin/stats                     # 获取系统统计
GET  /api/v1/admin/users/:user_id/statistics # 获取用户统计

# 操作日志
GET  /api/v1/admin/operation-logs            # 获取操作日志

# 系统配置
GET  /api/v1/admin/config                    # 获取系统配置
PUT  /api/v1/admin/config                    # 更新系统配置

# 公告管理
GET  /api/v1/admin/announcements             # 获取公告列表
POST /api/v1/admin/announcements             # 发布公告
```

**依赖服务**:
- `adminService.AdminService` - 管理服务

---

## 🔄 API调用流程

### 标准流程
```
客户端请求 
  → Router 
  → JWTAuth中间件（验证Token）
  → RequireRole中间件（验证管理员权限）
  → Admin API Handler 
  → Service层 
  → Repository层 
  → 数据库
```

### 关键流程说明

1. **认证流程**: 所有admin接口都需要有效的JWT Token
2. **授权流程**: 必须具有管理员角色（role=admin）
3. **审计日志**: 管理员操作应记录到操作日志中
4. **错误处理**: 统一使用shared.Error和shared.Success响应

---

## 🛡️ 中间件配置

### 1. JWT认证中间件
所有Admin接口都需要JWT认证：
```go
adminGroup.Use(middleware.JWTAuth())
```

### 2. 角色权限中间件
需要管理员角色：
```go
adminGroup.Use(middleware.RequireRole("admin"))
```

### 3. 完整中间件链
```go
admin := r.Group("/admin")
admin.Use(middleware.JWTAuth())            // JWT认证
admin.Use(middleware.RequireRole("admin")) // 管理员权限
```

---

## 📊 请求/响应示例

### 获取用户列表
```json
GET /api/v1/admin/users?page=1&page_size=20&role=user
Authorization: Bearer <token>

Response:
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "user_id": "user_123",
      "username": "testuser",
      "email": "test@example.com",
      "role": "user",
      "status": "active",
      "created_at": "2025-10-20T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100
  }
}
```

### 封禁用户
```json
POST /api/v1/admin/users/user_123/ban
Authorization: Bearer <token>
Content-Type: application/json

{
  "reason": "违反社区规则",
  "duration": 7,
  "durationUnit": "days"
}

Response:
{
  "code": 200,
  "message": "封禁成功",
  "data": null
}
```

### 审核内容
```json
POST /api/v1/admin/audit/audit_456/review
Authorization: Bearer <token>
Content-Type: application/json

{
  "action": "approve",
  "review_note": "内容符合规范",
  "penalty_type": ""
}

Response:
{
  "code": 200,
  "message": "审核已通过",
  "data": null
}
```

### 更新系统配置
```json
PUT /api/v1/admin/config
Authorization: Bearer <token>
Content-Type: application/json

{
  "allowRegistration": false,
  "maxUploadSize": 20971520
}

Response:
{
  "code": 200,
  "message": "更新成功",
  "data": null
}
```

---

## 🔧 设计原则

### 1. 单一职责原则
每个API文件只负责一个特定的管理领域，职责清晰、边界明确。

### 2. 依赖注入
通过构造函数注入依赖服务，便于单元测试和依赖管理。

### 3. RESTful风格
- 使用标准HTTP方法（GET/POST/PUT/DELETE）
- 资源路径清晰（/admin/users、/admin/quota）
- 状态码语义明确

### 4. 统一响应格式
使用 `shared.Success`、`shared.Error`、`shared.Paginated` 统一响应格式。

### 5. 权限分离
所有管理员API集中在admin模块，与普通用户API分离。

### 6. 审计日志
管理员的重要操作应记录到操作日志中。

---

## 📝 开发规范

### 1. 命名规范
- **API结构体**：`<功能>AdminAPI`（如 `UserAdminAPI`、`QuotaAdminAPI`）
- **构造函数**：`New<功能>AdminAPI`（如 `NewUserAdminAPI`）
- **方法名**：动词+名词（如 `GetUser`、`UpdateUserQuota`）

### 2. 错误处理
```go
if err != nil {
    shared.Error(c, http.StatusInternalServerError, "操作失败", err.Error())
    return
}
```

### 3. 参数验证
使用 `binding` 标签进行参数验证：
```go
type Request struct {
    Field string `json:"field" binding:"required"`
}
```

### 4. 权限检查
从Context中获取管理员信息：
```go
adminID, exists := c.Get("user_id")
if !exists {
    shared.Error(c, http.StatusUnauthorized, "未授权", "无法获取管理员信息")
    return
}
```

### 5. 日志记录
重要操作应记录日志：
```go
// TODO: 记录管理员操作日志
log.Printf("Admin %s updated user %s", adminID, userID)
```

---

## 🚀 扩展建议

### 未来可添加的管理功能

1. **敏感词管理**
   - 添加/删除敏感词
   - 批量导入敏感词
   - 敏感词分类管理

2. **内容管理**
   - 书籍管理（上架/下架/推荐）
   - 评论管理
   - 章节管理

3. **财务管理**
   - 收入统计
   - 支出统计
   - 结算管理

4. **数据分析**
   - 用户行为分析
   - 内容热度分析
   - 转化率分析

5. **权限管理**
   - 角色管理
   - 权限分配
   - 操作权限细化

---

## 🔒 安全注意事项

### 1. 权限验证
- 所有接口必须经过双重验证：JWT认证 + 管理员角色验证
- 禁止普通用户访问管理员接口

### 2. 敏感操作
- 删除、封禁等敏感操作应记录操作日志
- 重要操作建议添加二次确认

### 3. 数据保护
- 不应返回用户的敏感信息（如密码哈希）
- 日志中不应记录敏感数据

### 4. 操作审计
- 所有管理员操作应记录：操作人、操作时间、操作内容、操作结果
- 操作日志应不可篡改

### 5. 限流保护
- 防止管理员接口被恶意调用
- 建议添加操作频率限制

---

## 📚 相关文档

- [管理员API文档](../../../doc/api/管理员API文档.md)
- [用户管理API](../../../doc/api/用户管理API使用指南.md)
- [审核API文档](../../../doc/api/审核API文档.md)
- [架构设计规范](../../../doc/architecture/架构设计规范.md)

---

## 📋 API端点总览

### 用户管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/admin/users | 获取用户列表 |
| GET | /api/v1/admin/users/:id | 获取用户详情 |
| PUT | /api/v1/admin/users/:id | 更新用户信息 |
| DELETE | /api/v1/admin/users/:id | 删除用户 |
| POST | /api/v1/admin/users/:id/ban | 封禁用户 |
| POST | /api/v1/admin/users/:id/unban | 解除封禁 |

### AI配额管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/admin/quota/:userId | 获取用户配额详情 |
| PUT | /api/v1/admin/quota/:userId | 更新用户配额 |
| POST | /api/v1/admin/quota/:userId/suspend | 暂停用户配额 |
| POST | /api/v1/admin/quota/:userId/activate | 激活用户配额 |

### 审核管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/admin/audit/pending | 获取待审核内容 |
| GET | /api/v1/admin/audit/high-risk | 获取高风险审核记录 |
| GET | /api/v1/admin/audit/statistics | 获取审核统计 |
| POST | /api/v1/admin/audit/:id/review | 审核内容 |
| POST | /api/v1/admin/audit/:id/appeal/review | 审核申诉 |

### 系统管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/admin/stats | 获取系统统计 |
| GET | /api/v1/admin/users/:user_id/statistics | 获取用户统计 |
| GET | /api/v1/admin/operation-logs | 获取操作日志 |
| GET | /api/v1/admin/config | 获取系统配置 |
| PUT | /api/v1/admin/config | 更新系统配置 |
| POST | /api/v1/admin/withdraw/review | 审核提现 |
| GET | /api/v1/admin/announcements | 获取公告列表 |
| POST | /api/v1/admin/announcements | 发布公告 |

---

**版本**: v1.0  
**创建日期**: 2025-10-24  
**维护者**: Admin模块开发组

