# 青羽写作平台 - 通知系统API快速参考

## 基础路径
```
http://localhost:8080/api/v1
```

## 认证
所有API都需要在请求头中携带JWT Token:
```
Authorization: Bearer <token>
```

---

## 站内通知API

### 获取通知列表
```http
GET /notifications?type=system&read=false&limit=20&offset=0
```

### 获取通知详情
```http
GET /notifications/{id}
```

### 标记为已读
```http
PUT /notifications/{id}/read
```

### 批量标记为已读
```http
PUT /notifications/mark-read
Content-Type: application/json

{
  "ids": ["id1", "id2", "id3"]
}
```

### 全部标记为已读
```http
PUT /notifications/read-all
```

### 删除通知
```http
DELETE /notifications/{id}
```

### 批量删除通知
```http
DELETE /notifications/batch-delete
Content-Type: application/json

{
  "ids": ["id1", "id2", "id3"]
}
```

### 删除所有通知
```http
DELETE /notifications/delete-all
```

### 获取未读数量
```http
GET /notifications/unread-count
```

### 获取通知统计
```http
GET /notifications/stats
```

---

## 通知偏好设置API

### 获取偏好设置
```http
GET /notifications/preferences
```

### 更新偏好设置
```http
PUT /notifications/preferences
Content-Type: application/json

{
  "enableSystem": true,
  "enableSocial": true,
  "enableContent": true,
  "enableReward": true,
  "enableMessage": true,
  "enableUpdate": true,
  "enableMembership": true,
  "pushNotification": true,
  "emailNotification": {
    "enabled": false,
    "types": ["system", "reward"],
    "frequency": "immediate"
  },
  "smsNotification": {
    "enabled": false,
    "types": ["system"]
  },
  "quietHoursStart": "22:00",
  "quietHoursEnd": "08:00"
}
```

### 重置偏好设置
```http
POST /notifications/preferences/reset
```

---

## 邮件通知设置API

### 获取邮件通知设置
```http
GET /user-management/email-notifications
```

### 更新邮件通知设置
```http
PUT /user-management/email-notifications
Content-Type: application/json

{
  "enabled": true,
  "types": ["system", "content", "reward", "membership"],
  "frequency": "immediate"
}
```

**频率选项**:
- `immediate` - 立即发送
- `hourly` - 每小时汇总
- `daily` - 每天汇总

---

## 短信通知设置API

### 获取短信通知设置
```http
GET /user-management/sms-notifications
```

### 更新短信通知设置
```http
PUT /user-management/sms-notifications
Content-Type: application/json

{
  "enabled": true,
  "types": ["system", "reward", "membership"]
}
```

---

## 推送设备管理API

### 注册推送设备
```http
POST /notifications/push/register
Content-Type: application/json

{
  "deviceType": "ios",
  "deviceToken": "device_token_here",
  "deviceId": "unique_device_id"
}
```

**设备类型**:
- `ios` - iOS设备
- `android` - Android设备
- `web` - Web浏览器

### 取消注册推送设备
```http
DELETE /notifications/push/unregister/{deviceId}
```

### 获取推送设备列表
```http
GET /notifications/push/devices
```

---

## 通知类型

| 类型 | 代码 | 描述 |
|-----|------|------|
| 系统通知 | `system` | 平台公告、活动通知 |
| 社交通知 | `social` | 关注、点赞、评论 |
| 内容通知 | `content` | 作品审核、上架、下架 |
| 打赏通知 | `reward` | 收到打赏 |
| 私信通知 | `message` | 收到私信 |
| 更新通知 | `update` | 关注作品更新 |
| 会员通知 | `membership` | 会员到期、续费提醒 |

## 通知优先级

| 优先级 | 代码 |
|-------|------|
| 低 | `low` |
| 普通 | `normal` |
| 高 | `high` |
| 紧急 | `urgent` |

---

## 响应格式

### 成功响应
```json
{
  "code": 200,
  "message": "success",
  "data": { ... }
}
```

### 错误响应
```json
{
  "code": 400,
  "message": "INVALID_REQUEST",
  "data": "请求参数错误"
}
```

---

## 查询参数说明

### 通知列表筛选
- `type` - 通知类型（可选）
- `read` - 是否已读 true/false（可选）
- `priority` - 优先级 low/normal/high/urgent（可选）
- `keyword` - 关键词搜索（可选）
- `limit` - 每页数量，默认20，最大100（可选）
- `offset` - 偏移量，默认0（可选）
- `sortBy` - 排序字段 created_at/priority/read_at（可选）
- `sortDesc` - 是否降序 true/false，默认true（可选）

---

## 使用示例

### JavaScript/TypeScript
```typescript
// 获取通知列表
const response = await fetch('/api/v1/notifications?read=false&limit=20', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});
const data = await response.json();

// 标记为已读
await fetch('/api/v1/notifications/id123/read', {
  method: 'PUT',
  headers: {
    'Authorization': `Bearer ${token}`
  }
});

// 更新偏好设置
await fetch('/api/v1/notifications/preferences', {
  method: 'PUT',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    enableSystem: true,
    pushNotification: true
  })
});
```

### cURL
```bash
# 获取通知列表
curl -X GET "http://localhost:8080/api/v1/notifications?read=false&limit=20" \
  -H "Authorization: Bearer <token>"

# 标记为已读
curl -X PUT "http://localhost:8080/api/v1/notifications/id123/read" \
  -H "Authorization: Bearer <token>"

# 批量标记为已读
curl -X PUT "http://localhost:8080/api/v1/notifications/mark-read" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"ids": ["id1", "id2", "id3"]}'
```

---

## 常见错误码

| 错误码 | 描述 |
|-------|------|
| 400 | 请求参数错误 |
| 401 | 未授权访问 |
| 403 | 无权限操作 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

**文档版本**: 1.0
**最后更新**: 2026-01-03
