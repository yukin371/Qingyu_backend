# 配置管理 API 文档

## 概述

配置管理API允许管理员通过Web界面管理系统配置，支持实时更新、配置验证和备份恢复功能。

**基础路径**: `/api/v1/admin/config`  
**认证方式**: JWT Token + 管理员权限  
**版本**: v1.0

---

## 权限要求

所有配置管理API都需要：
1. ✅ JWT Token认证（Header: `Authorization: Bearer <token>`）
2. ✅ 管理员角色（role: admin）

---

## API 列表

### 1. 获取所有配置

获取系统所有可配置项（分组显示）

**请求**

```http
GET /api/v1/admin/config
Authorization: Bearer <admin_token>
```

**响应**

```json
{
  "code": 200,
  "message": "获取配置成功",
  "data": {
    "groups": [
      {
        "name": "server",
        "description": "服务器配置",
        "items": [
          {
            "key": "server.port",
            "value": "8080",
            "type": "string",
            "description": "服务器端口",
            "editable": true,
            "sensitive": false
          },
          {
            "key": "server.mode",
            "value": "debug",
            "type": "string",
            "description": "运行模式 (debug/release)",
            "editable": true,
            "sensitive": false
          }
        ]
      },
      {
        "name": "database",
        "description": "数据库配置",
        "items": [
          {
            "key": "database.uri",
            "value": "mong****:27017",
            "type": "string",
            "description": "MongoDB连接URI",
            "editable": true,
            "sensitive": true
          },
          {
            "key": "database.name",
            "value": "Qingyu_backend",
            "type": "string",
            "description": "数据库名称",
            "editable": true,
            "sensitive": false
          },
          {
            "key": "database.max_pool_size",
            "value": 100,
            "type": "number",
            "description": "最大连接池大小",
            "editable": true,
            "sensitive": false
          }
        ]
      },
      {
        "name": "jwt",
        "description": "JWT配置",
        "items": [
          {
            "key": "jwt.secret",
            "value": "qing****_key",
            "type": "string",
            "description": "JWT密钥",
            "editable": true,
            "sensitive": true
          },
          {
            "key": "jwt.expiration_hours",
            "value": 24,
            "type": "number",
            "description": "Token过期时间（小时）",
            "editable": true,
            "sensitive": false
          }
        ]
      },
      {
        "name": "redis",
        "description": "Redis配置",
        "items": [
          {
            "key": "redis.host",
            "value": "localhost",
            "type": "string",
            "description": "Redis主机地址",
            "editable": true,
            "sensitive": false
          },
          {
            "key": "redis.port",
            "value": 6379,
            "type": "number",
            "description": "Redis端口",
            "editable": true,
            "sensitive": false
          },
          {
            "key": "redis.password",
            "value": "",
            "type": "string",
            "description": "Redis密码",
            "editable": true,
            "sensitive": true
          },
          {
            "key": "redis.db",
            "value": 0,
            "type": "number",
            "description": "Redis数据库编号",
            "editable": true,
            "sensitive": false
          }
        ]
      },
      {
        "name": "ai",
        "description": "AI服务配置",
        "items": [
          {
            "key": "ai.api_key",
            "value": "AIza****ceE",
            "type": "string",
            "description": "AI API密钥",
            "editable": true,
            "sensitive": true
          },
          {
            "key": "ai.base_url",
            "value": "https://generativelanguage.googleapis.com/v1beta",
            "type": "string",
            "description": "AI API基础URL",
            "editable": true,
            "sensitive": false
          },
          {
            "key": "ai.max_tokens",
            "value": 2000,
            "type": "number",
            "description": "最大Token数",
            "editable": true,
            "sensitive": false
          }
        ]
      }
    ]
  },
  "timestamp": "2025-10-25T12:00:00Z"
}
```

**字段说明**

| 字段 | 类型 | 说明 |
|-----|------|------|
| `key` | string | 配置键（用于更新时引用） |
| `value` | any | 配置值（敏感信息会脱敏显示） |
| `type` | string | 值类型：string/number/boolean/object |
| `description` | string | 配置说明 |
| `editable` | boolean | 是否可编辑 |
| `sensitive` | boolean | 是否敏感信息（如密码、密钥） |

---

### 2. 获取单个配置

根据配置键获取单个配置项

**请求**

```http
GET /api/v1/admin/config/:key
Authorization: Bearer <admin_token>
```

**路径参数**

| 参数 | 类型 | 必填 | 说明 |
|-----|------|------|------|
| key | string | ✅ | 配置键（如 server.port） |

**示例**

```http
GET /api/v1/admin/config/server.port
```

**响应**

```json
{
  "code": 200,
  "message": "获取配置成功",
  "data": {
    "key": "server.port",
    "value": "8080",
    "type": "string",
    "description": "服务器端口",
    "editable": true,
    "sensitive": false
  },
  "timestamp": "2025-10-25T12:00:00Z"
}
```

---

### 3. 更新配置

更新单个配置项

**请求**

```http
PUT /api/v1/admin/config
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "key": "server.port",
  "value": "9090"
}
```

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|-----|------|------|------|
| key | string | ✅ | 配置键 |
| value | any | ✅ | 新的配置值 |

**响应**

```json
{
  "code": 200,
  "message": "配置更新成功",
  "data": null,
  "timestamp": "2025-10-25T12:00:00Z"
}
```

**注意事项**

1. 配置更新后会**自动备份**原配置文件
2. 更新成功后会**自动重新加载**配置
3. 如果重新加载失败，会**自动恢复备份**
4. 只有 `editable: true` 的配置项可以修改

---

### 4. 批量更新配置

批量更新多个配置项

**请求**

```http
PUT /api/v1/admin/config/batch
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "updates": [
    {
      "key": "server.port",
      "value": "9090"
    },
    {
      "key": "jwt.expiration_hours",
      "value": 48
    },
    {
      "key": "ai.max_tokens",
      "value": 4000
    }
  ]
}
```

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|-----|------|------|------|
| updates | array | ✅ | 更新列表，至少包含1项 |
| updates[].key | string | ✅ | 配置键 |
| updates[].value | any | ✅ | 新的配置值 |

**响应**

```json
{
  "code": 200,
  "message": "配置批量更新成功",
  "data": null,
  "timestamp": "2025-10-25T12:00:00Z"
}
```

**优势**

- ⚡ 一次性更新多个配置，减少重启次数
- 🔒 原子性操作，要么全部成功，要么全部失败
- 📦 只创建一个备份文件

---

### 5. 验证配置

验证YAML配置文件格式是否正确

**请求**

```http
POST /api/v1/admin/config/validate
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "yaml_content": "server:\n  port: 8080\n  mode: debug\n\njwt:\n  secret: my_secret\n  expiration_hours: 24"
}
```

**请求参数**

| 字段 | 类型 | 必填 | 说明 |
|-----|------|------|------|
| yaml_content | string | ✅ | YAML格式的配置内容 |

**响应 - 验证成功**

```json
{
  "code": 200,
  "message": "配置验证通过",
  "data": null,
  "timestamp": "2025-10-25T12:00:00Z"
}
```

**响应 - 验证失败**

```json
{
  "code": 400,
  "message": "配置验证失败",
  "error": "YAML格式错误: yaml: unmarshal errors:\n  line 3: cannot unmarshal !!str `invalid` into int",
  "timestamp": "2025-10-25T12:00:00Z"
}
```

**使用场景**

在实际应用配置前进行验证，避免错误配置导致系统异常。

---

### 6. 获取配置备份列表

获取所有可用的配置备份文件

**请求**

```http
GET /api/v1/admin/config/backups
Authorization: Bearer <admin_token>
```

**响应**

```json
{
  "code": 200,
  "message": "获取备份列表成功",
  "data": {
    "backups": [
      "config.yaml.backup"
    ]
  },
  "timestamp": "2025-10-25T12:00:00Z"
}
```

---

### 7. 恢复配置备份

将配置恢复到最近的备份

**请求**

```http
POST /api/v1/admin/config/restore
Authorization: Bearer <admin_token>
```

**响应**

```json
{
  "code": 200,
  "message": "配置恢复成功",
  "data": null,
  "timestamp": "2025-10-25T12:00:00Z"
}
```

**注意事项**

1. 恢复前会备份当前配置
2. 恢复成功后会自动重新加载配置
3. 如果没有备份文件，会返回错误

---

## 前端集成示例

### React/Vue 示例

```javascript
// 配置管理服务
class ConfigService {
  constructor(baseURL, token) {
    this.baseURL = baseURL;
    this.token = token;
  }

  // 获取所有配置
  async getAllConfigs() {
    const response = await fetch(`${this.baseURL}/api/v1/admin/config`, {
      headers: {
        'Authorization': `Bearer ${this.token}`
      }
    });
    return response.json();
  }

  // 更新配置
  async updateConfig(key, value) {
    const response = await fetch(`${this.baseURL}/api/v1/admin/config`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${this.token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ key, value })
    });
    return response.json();
  }

  // 批量更新
  async batchUpdateConfigs(updates) {
    const response = await fetch(`${this.baseURL}/api/v1/admin/config/batch`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${this.token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ updates })
    });
    return response.json();
  }

  // 恢复备份
  async restoreBackup() {
    const response = await fetch(`${this.baseURL}/api/v1/admin/config/restore`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${this.token}`
      }
    });
    return response.json();
  }
}

// 使用示例
const configService = new ConfigService('http://localhost:8080', adminToken);

// 加载配置
const { data } = await configService.getAllConfigs();
console.log('配置组:', data.groups);

// 更新端口
await configService.updateConfig('server.port', '9090');

// 批量更新
await configService.batchUpdateConfigs([
  { key: 'server.port', value: '9090' },
  { key: 'jwt.expiration_hours', value: 48 }
]);
```

### 前端UI示例

**配置管理页面结构**

```jsx
<div class="config-manager">
  {/* 配置组列表 */}
  <div class="config-groups">
    {groups.map(group => (
      <div class="config-group" key={group.name}>
        <h3>{group.description}</h3>
        
        {/* 配置项列表 */}
        {group.items.map(item => (
          <div class="config-item" key={item.key}>
            <label>{item.description}</label>
            
            {/* 根据类型显示不同输入框 */}
            {item.type === 'string' && (
              <input 
                type={item.sensitive ? 'password' : 'text'}
                value={item.value}
                disabled={!item.editable}
                onChange={e => handleChange(item.key, e.target.value)}
              />
            )}
            
            {item.type === 'number' && (
              <input 
                type="number"
                value={item.value}
                disabled={!item.editable}
                onChange={e => handleChange(item.key, Number(e.target.value))}
              />
            )}
            
            {item.type === 'boolean' && (
              <input 
                type="checkbox"
                checked={item.value}
                disabled={!item.editable}
                onChange={e => handleChange(item.key, e.target.checked)}
              />
            )}
            
            <span class="config-key">{item.key}</span>
          </div>
        ))}
      </div>
    ))}
  </div>
  
  {/* 操作按钮 */}
  <div class="actions">
    <button onClick={handleSave}>保存配置</button>
    <button onClick={handleBatchSave}>批量保存</button>
    <button onClick={handleRestore}>恢复备份</button>
    <button onClick={handleValidate}>验证配置</button>
  </div>
</div>
```

---

## 安全建议

### 1. 权限控制

✅ **已实现**
- JWT Token认证
- 管理员角色验证
- 路由级别权限控制

### 2. 敏感信息保护

✅ **已实现**
- 敏感字段自动脱敏显示
- 密码、密钥等只显示部分字符
- 传输过程使用HTTPS

### 3. 操作审计

🔄 **建议实现**
```go
// 记录配置变更日志
logger.Info("配置更新",
    zap.String("admin_id", adminID),
    zap.String("key", key),
    zap.Any("old_value", oldValue),
    zap.Any("new_value", newValue),
    zap.String("ip", clientIP),
)
```

### 4. 备份策略

✅ **已实现**
- 每次更新前自动备份
- 支持手动恢复备份
- 备份文件与配置文件同目录

### 5. 配置验证

✅ **已实现**
- 更新前验证YAML格式
- 验证配置值的合法性
- 失败自动回滚

---

## 错误码说明

| 错误码 | 说明 | 解决方案 |
|-------|------|---------|
| 400 | 参数错误 | 检查请求参数格式 |
| 401 | 未认证 | 提供有效的JWT Token |
| 403 | 权限不足 | 确保用户具有管理员角色 |
| 404 | 配置项不存在 | 检查配置键是否正确 |
| 500 | 服务器错误 | 查看服务器日志 |

---

## 测试工具

### Postman集合

```json
{
  "info": {
    "name": "配置管理API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "获取所有配置",
      "request": {
        "method": "GET",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{admin_token}}"
          }
        ],
        "url": "{{base_url}}/api/v1/admin/config"
      }
    },
    {
      "name": "更新配置",
      "request": {
        "method": "PUT",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{admin_token}}"
          },
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"key\": \"server.port\",\n  \"value\": \"9090\"\n}"
        },
        "url": "{{base_url}}/api/v1/admin/config"
      }
    }
  ]
}
```

### cURL示例

```bash
# 获取所有配置
curl -X GET \
  http://localhost:8080/api/v1/admin/config \
  -H 'Authorization: Bearer <admin_token>'

# 更新配置
curl -X PUT \
  http://localhost:8080/api/v1/admin/config \
  -H 'Authorization: Bearer <admin_token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "key": "server.port",
    "value": "9090"
  }'

# 批量更新
curl -X PUT \
  http://localhost:8080/api/v1/admin/config/batch \
  -H 'Authorization: Bearer <admin_token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "updates": [
      {"key": "server.port", "value": "9090"},
      {"key": "jwt.expiration_hours", "value": 48}
    ]
  }'

# 恢复备份
curl -X POST \
  http://localhost:8080/api/v1/admin/config/restore \
  -H 'Authorization: Bearer <admin_token>'
```

---

## 更新日志

| 版本 | 日期 | 变更内容 |
|-----|------|---------|
| 1.0 | 2025-10-25 | 初始版本，支持配置的增删改查和备份恢复 |

---

**维护者**：青羽后端团队  
**最后更新**：2025-10-25

