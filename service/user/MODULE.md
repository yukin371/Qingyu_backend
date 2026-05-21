---
name: User Service
description: 用户属性侧业务层，管理用户资料、密码操作、验证码、角色降级与统计数据
type: module
---

# User Service

> 最后更新：2026-05-21

## 职责

用户属性侧核心业务层，管理用户资料 CRUD、密码修改/重置、邮箱验证码、角色降级与用户统计数据聚合。不负责认证主链（注册、登录、登出、JWT 签发），该部分已收敛到 `service/auth`。

## 数据流

```
api/v1/user Handler → UserServiceImpl → UserRepository → MongoDB
                      PasswordService ↗
                      VerificationService → EmailVerificationTokenManager (内存单例)
                                            PasswordResetTokenManager (内存单例)
                                            channels.EmailService (邮件通道)
UserStatsServiceImpl → UserStatsUserRepository + UserStatsProjectRepository (跨域聚合查询)
```

- `UserServiceImpl`：用户资料 CRUD、邮箱验证、密码重置完整流程、角色降级、设备/绑定管理
- `PasswordService`：密码修改（需旧密码验证）、验证码重置密码（依赖 VerificationService）
- `VerificationService`：验证码生成/校验/标记已用，邮件发送委托 EmailService
- `UserStatsServiceImpl`：聚合用户统计（项目数、字数等），依赖 writer 域 ProjectRepository

## 约定 & 陷阱

- **认证边界**：公开注册、登录、登出和 JWT 签发已迁移到 `service/auth`，user 模块不再负责 token 生命周期
- **分层违规 TECHDEBT(#2026-03-22)**：`UserServiceImpl` 直接依赖 `AuthRepository`，理想链路应为 UserService → AuthService → AuthRepository，待后续迭代重构
- **角色降级边界**：`DowngradeRole` 涉及权限验证应由 AuthService 统管，当前暂时保留在 UserService，后续应通过 AuthService 鉴权后再执行
- **验证码/Token 存储为内存单例**：`EmailVerificationTokenManager` 和 `PasswordResetTokenManager` 均为内存 map + sync.RWMutex，重启即丢失；生产环境需迁移到 Redis
- **5位错误码体系**：40xxx 客户端 / 50xxx 服务端，`UserError` 支持字段级定位和 `IsRetryable` 判断
- **密码强度策略**：最小 8 位，必须包含大小写字母 + 数字，特殊字符可选；检测常见弱密码和连续字符
- **并发创建用户**：`CreateUser` 通过 MongoDB 唯一索引 + 3 次重试 + `isDuplicateKeyError` 检测处理并发冲突
- **防邮箱枚举**：`ResetPassword` 和 `RequestPasswordReset` 在用户不存在时也返回成功响应，避免泄露注册信息
- **TransactionManager**：支持 MongoDB 事务操作，通过 `MongoTransactionsDisabled()` 兜底非副本集环境
