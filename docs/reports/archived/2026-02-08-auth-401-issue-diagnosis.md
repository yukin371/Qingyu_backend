# 401认证错误诊断报告

## 问题描述

用户使用 `testadmin001` 账号成功登录后，访问 `/api/v1/writer/projects` 返回 401 Unauthorized 错误。

## 已确认的正确配置

### 1. 前端API路径修复
已修复以下文件的API路径配置：
- `project.ts`: `/projects` → `/writer/projects`
- `character.ts`: `/projects` → `/writer/projects`
- `document.ts`: `/projects` → `/writer/projects`
- `location.ts`: `/projects` → `/writer/projects`
- `timeline.ts`: `/projects` → `/writer/projects`

### 2. Token存储和发送机制
- **存储**: `authStore.ts` 使用 `storage.set('token', value)` → 实际存储到 `qingyu_token`
- **读取**: `http.service.ts` 使用 `localStorage.getItem('qingyu_token')`
- **发送**: 在请求头中添加 `Authorization: Bearer ${token}`

### 3. JWT配置一致性
- `config.yaml`: `secret: "qingyu_secret_key"`
- `config.go` 默认值: `"qingyu_secret_key"`
- `jwt.go` 默认值: `"qingyu-secret-key-change-in-production"`（仅当配置加载失败时使用）

## 可能的问题原因

### 1. 角色权限配置问题（最可能）

**问题描述**：
- `testadmin001` 拥有 `admin` 角色
- `testauthor001` 拥有 `author` 角色
- `permissions.yaml` 中 `admin` 有 `*:*` 权限
- 但 `author` 没有明确的 `project:read/write` 权限定义

**权限配置片段**（`permissions.yaml:56-70`）：
```yaml
# 作者权限
author:
  - "book:read"
  - "book:create"
  - "book:update"
  - "book:delete"
  - "chapter:read"
  - "chapter:create"
  - "chapter:update"
  - "chapter:delete"
  - "ai:generate"
  - "ai:chat"
  - "document:read"
  - "document:create"
```

注意：`author` 角色没有 `project:*` 相关权限！

### 2. 权限中间件未正确启用

**发现**：
- `internal/router/setup.go` 中的 `SetupAuthMiddleware` 返回一个空的中间件（第54-61行）
- `writer.go` 使用 `auth.JWTAuth()` 作为认证中间件，但只做JWT验证，不做权限检查

### 3. JWT验证可能的问题

虽然配置看起来一致，但可能存在：
- 环境变量覆盖了配置
- 配置文件加载失败，使用了默认值但与登录时的secret不一致

## 解决方案

### 方案1：使用 author 角色账号测试（推荐）

使用 `testauthor001` 账号登录测试：
```
用户名: testauthor001
密码: password
```

**注意**：这可能仍然会有问题，因为 author 角色没有明确的 project 权限。

### 方案2：更新 permissions.yaml

在 `permissions.yaml` 中为 `author` 角色添加 project 权限：
```yaml
author:
  - "book:read"
  - "book:create"
  - "book:update"
  - "book:delete"
  - "chapter:read"
  - "chapter:create"
  - "chapter:update"
  - "chapter:delete"
  - "ai:generate"
  - "ai:chat"
  - "document:read"
  - "document:create"
  - "project:read"      # 新增
  - "project:create"    # 新增
  - "project:update"    # 新增
```

### 方案3：检查后端日志

查看后端日志中的实际错误信息，确认：
1. JWT 验证是否真的通过了
2. 401 错误是在哪个中间件产生的
3. 实际的角色和权限检查结果

### 方案4：使用调试工具

已创建 `src/utils/debug-auth.html` 调试工具，可以：
1. 查看 LocalStorage 中的认证信息
2. 解析 JWT Token 内容
3. 测试 API 请求

**使用方法**：
```bash
# 在前端开发模式下
cp src/utils/debug-auth.html public/debug.html
# 然后访问 http://localhost:5173/debug.html
```

## 推荐的排查步骤

1. **使用调试工具确认Token信息**
   - 访问 `debug.html`
   - 检查 Token 是否存在
   - 检查 Token 中的 roles 字段是否正确

2. **使用 author 角色账号测试**
   - 登录 testauthor001
   - 尝试访问 writer 端点
   - 对比两种角色的行为差异

3. **检查后端日志**
   - 查看是否有 JWT 验证失败的日志
   - 查看是否有权限检查失败的日志

4. **更新权限配置**
   - 如果确认是权限配置问题，更新 `permissions.yaml`
   - 重启后端服务

## 文件修改记录

### 已修改的文件
1. `Qingyu_fronted/src/modules/writer/api/project.ts`
2. `Qingyu_fronted/src/modules/writer/api/character.ts`
3. `Qingyu_fronted/src/modules/writer/api/document.ts`
4. `Qingyu_fronted/src/modules/writer/api/location.ts`
5. `Qingyu_fronted/src/modules/writer/api/timeline.ts`

### 新创建的文件
1. `Qingyu_fronted/src/utils/debug-auth.html` - 认证调试工具

## 待确认项

1. ✅ 前端API路径已修复
2. ❓ 后端权限中间件是否正确启用
3. ❓ JWT secret 是否在所有地方都使用一致的值
4. ❓ testadmin001 的角色是否在数据库中正确设置
5. ❓ 是否需要使用 testauthor001 而不是 testadmin001

## 下一步行动

请主人按照以下顺序操作喵：

1. 使用 testauthor001 账号登录测试
2. 如果仍然失败，访问 debug.html 查看详细信息
3. 检查后端日志获取更多错误信息
4. 根据结果决定是否需要修改权限配置

---

**报告生成时间**: 2026-02-08
**问题状态**: 待解决
**优先级**: 高
