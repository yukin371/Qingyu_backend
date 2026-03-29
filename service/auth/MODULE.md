# Auth Service

> 最后更新：2026-03-29

## 职责

认证授权层，管理用户注册/登录、OAuth 第三方登录、JWT Token 签发与刷新、RBAC 角色权限、Session 管理。不管理用户资料编辑（由 User 模块负责）。

## 数据流

```
API Handler → AuthServiceImpl → UserRepository → MongoDB
                ↓
         JWT Token（签发/验证）
                ↓
         Session Store（Redis）
```

## 约定 & 陷阱

- **OAuth 用户名生成**：第三方登录首次注册时 `generateUsernameFromProvider` 自动生成用户名，可能需要用户后续修改
- **Session 与 JWT 双轨**：系统同时维护 JWT Token 和 Redis Session，Token 刷新时必须同步更新 Session
- **权限检查链**：`CheckPermission` → `GetUserPermissions` → 角色权限合并，权限变更需要等待 Session 过期或主动刷新
- **Token 刷新窗口**：RefreshToken 有独立的过期时间，与 Access Token 过期时间不同
