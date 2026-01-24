# OAuth应用注册指南

**版本**: v1.0  
**创建日期**: 2026-01-23  
**状态**: 已完成

## 概述

本文档介绍如何在 Google 和 GitHub 平台注册 OAuth 应用，以便为青羽写作平台启用第三方登录功能。

## 前提条件

- 本地开发环境已配置好
- 项目代码已克隆到本地
- MongoDB 和 Redis 服务已启动（如需要）

## 开发环境配置

### 回调地址说明

**本地开发环境的回调地址格式**：

```
http://localhost:你的前端端口/oauth/callback
```

例如：如果你的前端运行在 `localhost:3000`，回调地址应为：
```
http://localhost:3000/oauth/callback
```

**重要说明**：
- 本地开发使用 `http://` 协议（不需要 HTTPS）
- 域名使用 `localhost`
- 端口号根据你的前端配置而定

---

## Google OAuth 应用注册

### 步骤 1：访问 Google Cloud Console

1. 访问 [Google Cloud Console](https://console.cloud.google.com/)
2. 登录你的 Google 账号
3. 创建新项目或选择现有项目

### 步骤 2：启用 OAuth 2.0

1. 在左侧菜单中，选择 **API 和服务** > **凭据**
2. 点击顶部的 **+ 创建凭据** 按钮
3. 选择 **OAuth 客户端 ID**

### 步骤 3：配置 OAuth 客户端

1. **应用类型**：选择 **Web 应用**
2. **名称**：输入应用名称，例如 `Qingyu (Dev)`
3. **已授权的重定向 URI**：添加以下地址：
   ```
   http://localhost:3000/oauth/callback
   ```
   > 注意：根据前端实际端口调整端口号

4. 点击 **创建**

### 步骤 4：获取凭据

创建成功后，你会看到：
- **客户端 ID**：复制到 `GOOGLE_CLIENT_ID` 环境变量
- **客户端密钥**：点击复制，保存到 `GOOGLE_CLIENT_SECRET` 环境变量

### 步骤 5：配置环境变量

在 `Qingyu_backend/.env` 文件中添加：

```bash
GOOGLE_CLIENT_ID=你的客户端ID
GOOGLE_CLIENT_SECRET=你的客户端密钥
```

### 示例

```bash
# .env 文件
GOOGLE_CLIENT_ID=123456789-abcdefg.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=GOCSPX-xxxxxxxxxxxxx
```

---

## GitHub OAuth 应用注册

### 步骤 1：访问 GitHub 设置

1. 访问 [GitHub](https://github.com/) 并登录
2. 点击右上角头像 > **Settings**
3. 在左侧菜单最下方，点击 **Developer settings**

### 步骤 2：创建 OAuth App

1. 点击 **OAuth apps** > **New OAuth App**
2. 填写应用信息：
   - **Application name**: `Qingyu (Dev)`
   - **Homepage URL**: `http://localhost:5173`
   - **Application description**: `Qingyu 本地开发环境`
   - **Authorization callback URL**: `http://localhost:5173/oauth/callback`

3. 点击 **Register application**

### 步骤 3：获取凭据

注册成功后，你会看到：
- **Client ID**：显示在页面顶部
- **Client Secret**：点击 **Generate a new client secret** 按钮生成

### 步骤 4：配置环境变量

在 `Qingyu_backend/.env` 文件中添加：

```bash
GITHUB_CLIENT_ID=你的GitHub客户端ID
GITHUB_CLIENT_SECRET=你的GitHub客户端密钥
```

### 示例

```bash
# .env 文件
GITHUB_CLIENT_ID=Iv1xxxxxxxxxxxxxxxx
GITHUB_CLIENT_SECRET=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

---

## 配置验证

### 验证步骤

1. **重启后端服务**
   ```bash
   cd Qingyu_backend
   go run cmd/server/main.go
   ```

2. **查看启动日志**
   
   成功配置后，你应该看到类似的日志：
   ```
   ✓ OAuthService初始化完成 (启用的提供商: [google github])
   ```

3. **测试 API**
   
   使用 curl 测试获取授权 URL：
   ```bash
   curl -X POST http://localhost:8080/api/v1/shared/oauth/google/authorize \
     -H "Content-Type: application/json" \
     -d '{"redirect_uri":"http://localhost:3000/oauth/callback","state":"test123"}'
   ```

   成功响应示例：
   ```json
   {
     "code": 200,
     "message": "获取授权URL成功",
     "data": {
       "authorize_url": "https://accounts.google.com/o/oauth2/v2/auth?...",
       "provider": "google"
     }
   }
   ```

---

## 常见问题

### Q1: Google OAuth 提示 "redirect_uri_mismatch"

**原因**：回调地址不匹配

**解决方法**：
1. 检查 Google Cloud Console 中的回调地址配置
2. 确保与请求中的 `redirect_uri` 完全一致（包括协议、域名、端口、路径）
3. 注意：`http://localhost:3000/oauth/callback` 和 `http://localhost:3000/oauth/callback/` 是不同的

### Q2: GitHub OAuth 提示 "Redirect URI mismatch"

**原因**：回调地址不匹配

**解决方法**：
1. 检查 GitHub OAuth App 的回调地址配置
2. 确保与请求中的 `redirect_uri` 完全一致
3. 每次修改回调地址后需要重新注册 OAuth App

### Q3: 后端启动时提示 "OAuth配置为空"

**原因**：环境变量未正确加载

**解决方法**：
1. 确保 `.env` 文件在项目根目录
2. 检查环境变量名称是否正确（区分大小写）
3. 确保客户端 ID 和密钥已正确填写
4. 尝试重启服务

### Q4: 本地开发是否需要 HTTPS？

**回答**：不需要。Google 和 GitHub 都支持 `http://localhost` 的回调地址，但不支持其他域名的 HTTP 回调。

生产环境部署时必须使用 HTTPS。

### Q5: 如何同时启用多个 OAuth 提供商？

**回答**：只需在 `.env` 文件中配置多个提供商的凭据即可：

```bash
# 启用 Google
GOOGLE_CLIENT_ID=xxx
GOOGLE_CLIENT_SECRET=xxx

# 启用 GitHub
GITHUB_CLIENT_ID=xxx
GITHUB_CLIENT_SECRET=xxx

# 启用 QQ
QQ_CLIENT_ID=xxx
QQ_CLIENT_SECRET=xxx
```

所有已配置的提供商都会自动启用。

---

## 安全注意事项

### 开发环境

1. **不要将凭据提交到代码仓库**
   - `.env` 文件已在 `.gitignore` 中
   - 使用 `.env.example` 提供配置模板

2. **定期更换密钥**
   - 如果密钥意外泄露，立即在平台上重新生成

3. **限制应用权限范围**
   - Google: `openid email profile`
   - GitHub: `read:user user:email`

### 生产环境

1. **必须使用 HTTPS**
2. **使用环境变量或密钥管理服务**
3. **设置严格的回调地址白名单**
4. **定期审查授权的应用**

---

## API 端点参考

### 获取授权 URL

```
POST /api/v1/shared/oauth/{provider}/authorize
```

**请求体**：
```json
{
  "redirect_uri": "http://localhost:3000/oauth/callback",
  "state": "随机字符串（用于防止CSRF攻击）"
}
```

**支持的平台**：`google`、`github`、`qq`

### 处理回调

```
POST /api/v1/shared/oauth/{provider}/callback
```

**请求体**：
```json
{
  "code": "授权码",
  "state": "与请求时相同的state"
}
```

### 获取绑定账号列表（需认证）

```
GET /api/v1/shared/oauth/accounts
```

---

## 相关文档

- [OAuth 2.0 设计文档](../design/auth/第三方登录OAuth设计文档.md)
- [OAuth功能待完善清单](../todo/OAuth功能待完善清单.md)
- [OAuth集成完成报告](../design/auth/OAuth集成完成报告.md)

---

## 附录：快速配置检查清单

- [ ] 已注册 Google OAuth 应用
- [ ] 已注册 GitHub OAuth 应用（可选）
- [ ] 已将凭据添加到 `.env` 文件
- [ ] 后端服务启动时显示 "OAuthService初始化完成"
- [ ] 前端回调地址与OAuth应用配置一致
- [ ] 已测试授权 URL 获取 API

---

**文档维护者**: AI Assistant  
**最后更新**: 2026-01-23
