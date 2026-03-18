# ServerDemo - 测试用启动器

完整的后端服务启动器，展示完整的服务初始化流程。

## ✨ 测试模式功能

### 1. 跳过认证（skip_auth）

**测试模式**允许你跳过JWT认证，直接访问所有API，无需登录获取Token。这对于API开发和调试非常有用。

#### 启用跳过认证模式

**方式1：修改配置文件**（推荐）

编辑 `configs/config.yaml`：
```yaml
server:
  port: "9090"
  mode: "debug"
  skip_auth: true  # 设置为 true 启用测试模式
```

**方式2：环境变量**

```bash
# Windows PowerShell
$env:SKIP_AUTH="true"; go run cmd/serverdemo/main.go

# Windows CMD
set SKIP_AUTH=true && go run cmd/serverdemo/main.go

# Linux/Mac
SKIP_AUTH=true go run cmd/serverdemo/main.go
```

#### 测试用户信息

测试模式下，所有请求会自动注入以下用户信息：
- `user_id`: `test-user-id`
- `username`: `test-user`
- `roles`: `["reader", "author", "admin"]`

### 2. 固定邮箱验证码（fixed_code）

**固定验证码模式**允许你使用固定的验证码进行注册测试，无需真实发送邮件。

#### 启用固定验证码模式

编辑 `configs/config.yaml`：
```yaml
email:
  enabled: true
  fixed_code: "123456"  # 设置固定验证码
```

#### 使用说明

启用后，注册时可以直接使用 `123456` 作为验证码完成注册，系统会在日志中输出：

```
[EmailCode] 测试模式已启用，使用固定验证码 {"fixed_code": "123456"}
[EmailCode] 测试模式：使用固定验证码 {"email": "test@example.com", "code": "123456"}
```

## API 测试示例

### 正常模式（需要先登录获取 Token）

```bash
# 1. 登录获取 Token
curl -X POST http://localhost:9090/api/v1/shared/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"demo","password":"demo123"}'

# 2. 使用 Token 访问 API
curl http://localhost:9090/api/v1/bookstore/homepage \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### 测试模式（无需 Token）

```bash
# 设置 SKIP_AUTH=true 后，直接访问 API
curl http://localhost:9090/api/v1/bookstore/homepage

# 访问需要认证的 API
curl http://localhost:9090/api/v1/reader/bookmarks

# 访问 AI API
curl http://localhost:9090/api/v1/ai/chat \
  -H "Content-Type: application/json" \
  -d '{"message":"你好"}'
```

### 注册测试（使用固定验证码）

```bash
# 1. 发送注册验证码
curl -X POST http://localhost:9090/api/v1/shared/auth/send-code \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com"}'

# 2. 使用固定验证码完成注册
curl -X POST http://localhost:9090/api/v1/shared/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123456","code":"123456"}'
```

## 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `SKIP_AUTH` | 跳过JWT认证（测试模式） | `false` |
| `CONFIG_FILE` | 配置文件路径 | `.` |
| `GIN_MODE` | Gin运行模式 | `debug` |

## 配置文件

服务启动时会自动加载以下位置的配置文件（按优先级）：
1. `./configs/config.yaml`
2. `./config/config.yaml`
3. `./config.yaml`

## 服务端口

默认端口：`9090`

可通过配置文件 `configs/config.yaml` 中的 `server.port` 修改。

## 常见问题

### Q: 如何测试需要特定角色的API？

A: 测试模式下，用户拥有所有角色（reader、author、admin），可以直接访问任何API。

### Q: 生产环境可以使用测试模式吗？

A: **绝对不可以！** `SKIP_AUTH` 和 `fixed_code` 仅用于本地开发和测试，生产环境必须禁用。

### Q: 测试模式会影响API功能吗？

A: 不会。测试模式只是跳过认证验证或使用固定验证码，所有业务逻辑正常执行。

### Q: 固定验证码是什么？

A: 固定验证码是用于测试的简化功能。启用后，任何邮箱注册都可以使用预设的验证码（如"123456"）完成验证，无需真实发送邮件。

### Q: 可以同时使用跳过认证和固定验证码吗？

A: 可以。两者互不影响，可以同时启用用于快速测试。
