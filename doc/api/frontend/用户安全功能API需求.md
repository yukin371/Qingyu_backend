# 用户安全功能API需求文档

**文档版本**: v1.0  
**创建日期**: 2025-10-29  
**状态**: 待实现  
**优先级**: P1

---

## 📋 目录

1. [概述](#概述)
2. [API列表](#api列表)
3. [详细设计](#详细设计)
4. [数据模型](#数据模型)
5. [错误码定义](#错误码定义)
6. [实施计划](#实施计划)

---

## 概述

### 功能描述

用户安全功能为用户提供账户安全管理能力，包括手机号/邮箱绑定、密码管理、设备管理和账号注销等核心安全功能。

### 业务价值

- 提升账户安全性，降低账户被盗风险
- 完善用户身份验证体系
- 满足监管要求（实名制、数据安全法等）
- 提升用户信任度和平台可信度

### 技术目标

- 短信/邮件验证码发送频率限制（防刷）
- 验证码有效期管理（5-10分钟）
- 敏感操作二次验证
- 设备指纹识别
- 安全审计日志

---

## API列表

### 2.1 验证码相关

| 序号 | API路径 | 方法 | 说明 | 优先级 |
|-----|---------|------|------|--------|
| 1 | `/api/v1/user/security/send-sms-code` | POST | 发送手机验证码 | P0 |
| 2 | `/api/v1/user/security/send-email-code` | POST | 发送邮箱验证码 | P0 |
| 3 | `/api/v1/user/security/verify-code` | POST | 验证验证码（通用） | P1 |

### 2.2 手机号管理

| 序号 | API路径 | 方法 | 说明 | 优先级 |
|-----|---------|------|------|--------|
| 4 | `/api/v1/user/security/bind-phone` | POST | 绑定手机号 | P0 |
| 5 | `/api/v1/user/security/unbind-phone` | POST | 解绑手机号 | P1 |
| 6 | `/api/v1/user/security/change-phone` | POST | 更换手机号 | P1 |

### 2.3 邮箱管理

| 序号 | API路径 | 方法 | 说明 | 优先级 |
|-----|---------|------|------|--------|
| 7 | `/api/v1/user/security/bind-email` | POST | 绑定邮箱 | P0 |
| 8 | `/api/v1/user/security/unbind-email` | POST | 解绑邮箱 | P1 |
| 9 | `/api/v1/user/security/change-email` | POST | 更换邮箱 | P1 |

### 2.4 密码管理

| 序号 | API路径 | 方法 | 说明 | 优先级 |
|-----|---------|------|------|--------|
| 10 | `/api/v1/user/security/change-password` | POST | 修改密码 | P0 |
| 11 | `/api/v1/user/security/reset-password` | POST | 重置密码 | P0 |
| 12 | `/api/v1/user/security/verify-password` | POST | 验证当前密码 | P1 |

### 2.5 设备管理

| 序号 | API路径 | 方法 | 说明 | 优先级 |
|-----|---------|------|------|--------|
| 13 | `/api/v1/user/security/devices` | GET | 获取登录设备列表 | P1 |
| 14 | `/api/v1/user/security/devices/:id` | DELETE | 移除登录设备 | P1 |
| 15 | `/api/v1/user/security/devices/current` | GET | 获取当前设备信息 | P2 |

### 2.6 账号管理

| 序号 | API路径 | 方法 | 说明 | 优先级 |
|-----|---------|------|------|--------|
| 16 | `/api/v1/user/security/deactivate` | POST | 申请注销账号 | P1 |
| 17 | `/api/v1/user/security/cancel-deactivation` | POST | 取消注销申请 | P2 |
| 18 | `/api/v1/user/security/check-status` | GET | 检查账号状态 | P2 |

---

## 详细设计

### 3.1 发送手机验证码

#### 基本信息
- **路径**: `/api/v1/user/security/send-sms-code`
- **方法**: `POST`
- **认证**: 可选（已登录用户使用当前手机号，未登录用户指定手机号）
- **限流**: 60秒/次，每天10次/手机号

#### 请求参数

```json
{
  "phone": "13800138000",
  "scene": "bind_phone",  // 场景：bind_phone, change_phone, login, reset_password
  "captcha_token": "xxx"  // 图形验证码token（防刷）
}
```

**参数说明**:
- `phone`: 手机号（必填）
- `scene`: 使用场景（必填）
- `captcha_token`: 图形验证码token（可选，高频调用时需要）

#### 响应示例

```json
{
  "code": 200,
  "message": "验证码已发送",
  "data": {
    "expires_in": 300,  // 有效期（秒）
    "can_resend_after": 60  // 可重发时间（秒）
  }
}
```

#### 业务规则

1. **频率限制**:
   - 同一手机号：60秒内只能发送1次
   - 同一IP：1分钟内最多3次，1小时内最多20次
   - 同一手机号：24小时内最多10次

2. **验证码规则**:
   - 6位数字验证码
   - 有效期5分钟
   - 同一手机号同一场景只保留最新的验证码

3. **防刷机制**:
   - 短时间内多次失败需要图形验证码
   - 异常IP自动封禁

4. **通知内容**:
   ```
   【青羽写作】您的验证码是：123456，5分钟内有效。请勿泄露给他人。
   ```

---

### 3.2 发送邮箱验证码

#### 基本信息
- **路径**: `/api/v1/user/security/send-email-code`
- **方法**: `POST`
- **认证**: 可选
- **限流**: 60秒/次，每天10次/邮箱

#### 请求参数

```json
{
  "email": "user@example.com",
  "scene": "bind_email",  // 场景：bind_email, change_email, login, reset_password
  "captcha_token": "xxx"
}
```

#### 响应示例

```json
{
  "code": 200,
  "message": "验证码已发送至邮箱",
  "data": {
    "expires_in": 600,  // 有效期10分钟
    "can_resend_after": 60
  }
}
```

#### 业务规则

- 频率限制同手机验证码
- 验证码：6位数字或字母组合
- 有效期：10分钟
- 邮件主题：【青羽写作】邮箱验证码

---

### 3.3 绑定手机号

#### 基本信息
- **路径**: `/api/v1/user/security/bind-phone`
- **方法**: `POST`
- **认证**: 必须
- **权限**: 用户本人

#### 请求参数

```json
{
  "phone": "13800138000",
  "code": "123456",  // 验证码
  "password": "***"  // 当前密码（二次验证）
}
```

#### 响应示例

```json
{
  "code": 200,
  "message": "手机号绑定成功",
  "data": {
    "phone": "138****8000",  // 脱敏显示
    "verified": true,
    "bound_at": "2025-10-29T10:30:00Z"
  }
}
```

#### 业务规则

1. **验证步骤**:
   - 验证用户登录状态
   - 验证当前密码
   - 验证手机验证码
   - 检查手机号是否已被其他账号绑定

2. **限制条件**:
   - 一个手机号只能绑定一个账号
   - 已绑定手机号的账号不能重复绑定
   - 必须先解绑才能绑定新号码

3. **安全审计**:
   - 记录绑定操作日志
   - 发送绑定成功通知

---

### 3.4 绑定邮箱

#### 基本信息
- **路径**: `/api/v1/user/security/bind-email`
- **方法**: `POST`
- **认证**: 必须
- **权限**: 用户本人

#### 请求参数

```json
{
  "email": "user@example.com",
  "code": "123456",
  "password": "***"
}
```

#### 响应示例

```json
{
  "code": 200,
  "message": "邮箱绑定成功",
  "data": {
    "email": "u***@example.com",
    "verified": true,
    "bound_at": "2025-10-29T10:30:00Z"
  }
}
```

#### 业务规则

- 与手机号绑定类似
- 一个邮箱只能绑定一个账号
- 绑定后发送确认邮件

---

### 3.5 修改密码

#### 基本信息
- **路径**: `/api/v1/user/security/change-password`
- **方法**: `POST`
- **认证**: 必须
- **权限**: 用户本人

#### 请求参数

```json
{
  "old_password": "***",  // 当前密码
  "new_password": "***",  // 新密码
  "confirm_password": "***",  // 确认新密码
  "verification_code": "123456"  // 手机或邮箱验证码（可选，增强安全）
}
```

#### 响应示例

```json
{
  "code": 200,
  "message": "密码修改成功，请重新登录",
  "data": {
    "require_relogin": true
  }
}
```

#### 业务规则

1. **密码强度要求**:
   - 长度8-32位
   - 必须包含字母和数字
   - 建议包含特殊字符
   - 不能与用户名相同
   - 不能是常见弱密码

2. **验证步骤**:
   - 验证当前密码
   - 验证新密码格式
   - 验证两次输入一致性
   - （可选）验证验证码

3. **后续操作**:
   - 修改成功后使所有token失效
   - 要求用户重新登录
   - 发送密码修改通知（邮件/短信）
   - 记录修改日志

---

### 3.6 重置密码

#### 基本信息
- **路径**: `/api/v1/user/security/reset-password`
- **方法**: `POST`
- **认证**: 不需要
- **场景**: 忘记密码

#### 请求参数

```json
{
  "identifier": "13800138000 或 user@example.com",  // 手机号或邮箱
  "code": "123456",  // 验证码
  "new_password": "***",
  "confirm_password": "***"
}
```

#### 响应示例

```json
{
  "code": 200,
  "message": "密码重置成功，请登录",
  "data": {
    "require_login": true
  }
}
```

#### 业务规则

1. **重置流程**:
   - 用户输入手机号或邮箱
   - 发送验证码
   - 验证验证码
   - 设置新密码
   - 完成重置

2. **安全措施**:
   - 验证码有效期5分钟
   - 每次重置使所有token失效
   - 发送重置成功通知
   - 记录重置日志（含IP和设备信息）

---

### 3.7 获取登录设备列表

#### 基本信息
- **路径**: `/api/v1/user/security/devices`
- **方法**: `GET`
- **认证**: 必须
- **权限**: 用户本人

#### 请求参数

```
Query参数:
- page: 页码（默认1）
- page_size: 每页数量（默认20）
- status: 设备状态（active, expired）
```

#### 响应示例

```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "device_id": "xxx",
      "device_name": "iPhone 14 Pro",
      "device_type": "mobile",  // mobile, desktop, tablet
      "os": "iOS 17.0",
      "browser": "Safari 17.0",
      "ip": "123.456.789.0",
      "location": "北京市朝阳区",
      "last_active": "2025-10-29T10:30:00Z",
      "login_at": "2025-10-28T08:00:00Z",
      "is_current": true,
      "status": "active"
    }
  ],
  "total": 5,
  "page": 1,
  "page_size": 20
}
```

#### 业务规则

1. **设备识别**:
   - 基于User-Agent解析设备信息
   - 使用设备指纹技术
   - IP地址定位

2. **状态管理**:
   - active: 活跃设备（30天内有活动）
   - expired: 过期设备（超过30天无活动）

3. **安全提示**:
   - 标识当前登录设备
   - 显示最后活跃时间
   - 异常设备预警

---

### 3.8 移除登录设备

#### 基本信息
- **路径**: `/api/v1/user/security/devices/:device_id`
- **方法**: `DELETE`
- **认证**: 必须
- **权限**: 用户本人

#### 请求参数

```
Path参数:
- device_id: 设备ID

Body参数:
{
  "password": "***"  // 当前密码（二次验证）
}
```

#### 响应示例

```json
{
  "code": 200,
  "message": "设备已移除",
  "data": {
    "device_id": "xxx",
    "removed_at": "2025-10-29T10:30:00Z"
  }
}
```

#### 业务规则

1. **验证步骤**:
   - 验证用户密码
   - 验证设备归属
   - 不能移除当前设备

2. **执行操作**:
   - 使该设备的token失效
   - 记录移除日志
   - 发送设备移除通知

---

### 3.9 申请注销账号

#### 基本信息
- **路径**: `/api/v1/user/security/deactivate`
- **方法**: `POST`
- **认证**: 必须
- **权限**: 用户本人

#### 请求参数

```json
{
  "password": "***",  // 当前密码
  "verification_code": "123456",  // 手机或邮箱验证码
  "reason": "不再使用",  // 注销原因（可选）
  "feedback": "建议内容"  // 反馈意见（可选）
}
```

#### 响应示例

```json
{
  "code": 200,
  "message": "注销申请已提交，将在7天后生效",
  "data": {
    "deactivation_id": "xxx",
    "effective_date": "2025-11-05T10:30:00Z",
    "can_cancel_before": "2025-11-05T10:30:00Z"
  }
}
```

#### 业务规则

1. **注销条件检查**:
   - 账户无未完成订单
   - 钱包余额为零或已提现
   - 无违规处罚期
   - 无未结清债务

2. **冷静期机制**:
   - 提交申请后进入7天冷静期
   - 冷静期内可以取消注销
   - 冷静期结束后自动执行注销

3. **注销流程**:
   - 验证密码和验证码
   - 检查注销条件
   - 创建注销申请
   - 发送确认通知
   - 7天后执行注销

4. **数据处理**:
   - 匿名化用户信息
   - 删除个人隐私数据
   - 保留必要的交易记录（法律要求）
   - 发布内容标记为"已注销用户"

5. **不可恢复提示**:
   - 明确告知用户注销后不可恢复
   - 提示备份重要数据
   - 确认理解后才能提交

---

## 数据模型

### 4.1 验证码记录

```go
type VerificationCode struct {
    ID          string    `bson:"_id"`
    Type        string    `bson:"type"`         // sms, email
    Recipient   string    `bson:"recipient"`    // 手机号或邮箱
    Code        string    `bson:"code"`         // 验证码（加密存储）
    Scene       string    `bson:"scene"`        // 使用场景
    ExpiresAt   time.Time `bson:"expires_at"`   // 过期时间
    UsedAt      *time.Time `bson:"used_at"`     // 使用时间
    IP          string    `bson:"ip"`           // 请求IP
    UserAgent   string    `bson:"user_agent"`   // 用户代理
    Status      string    `bson:"status"`       // pending, used, expired
    CreatedAt   time.Time `bson:"created_at"`
}
```

### 4.2 登录设备

```go
type LoginDevice struct {
    ID          string    `bson:"_id"`
    UserID      string    `bson:"user_id"`
    DeviceID    string    `bson:"device_id"`    // 设备唯一标识
    DeviceName  string    `bson:"device_name"`  // 设备名称
    DeviceType  string    `bson:"device_type"`  // mobile, desktop, tablet
    OS          string    `bson:"os"`           // 操作系统
    Browser     string    `bson:"browser"`      // 浏览器
    IP          string    `bson:"ip"`           // IP地址
    Location    string    `bson:"location"`     // 地理位置
    LastActive  time.Time `bson:"last_active"`  // 最后活跃时间
    LoginAt     time.Time `bson:"login_at"`     // 登录时间
    Status      string    `bson:"status"`       // active, expired, removed
    CreatedAt   time.Time `bson:"created_at"`
    UpdatedAt   time.Time `bson:"updated_at"`
}
```

### 4.3 注销申请

```go
type DeactivationRequest struct {
    ID            string     `bson:"_id"`
    UserID        string     `bson:"user_id"`
    Reason        string     `bson:"reason"`         // 注销原因
    Feedback      string     `bson:"feedback"`       // 用户反馈
    EffectiveDate time.Time  `bson:"effective_date"` // 生效日期
    CancelledAt   *time.Time `bson:"cancelled_at"`   // 取消时间
    ExecutedAt    *time.Time `bson:"executed_at"`    // 执行时间
    Status        string     `bson:"status"`         // pending, cancelled, executed
    CreatedAt     time.Time  `bson:"created_at"`
    UpdatedAt     time.Time  `bson:"updated_at"`
}
```

### 4.4 安全日志

```go
type SecurityLog struct {
    ID          string    `bson:"_id"`
    UserID      string    `bson:"user_id"`
    Action      string    `bson:"action"`       // 操作类型
    Resource    string    `bson:"resource"`     // 操作资源
    Result      string    `bson:"result"`       // success, failure
    IP          string    `bson:"ip"`
    UserAgent   string    `bson:"user_agent"`
    Details     map[string]interface{} `bson:"details"`  // 详细信息
    CreatedAt   time.Time `bson:"created_at"`
}
```

---

## 错误码定义

| 错误码 | 说明 | HTTP状态码 |
|-------|------|-----------|
| 40001 | 验证码发送频率过快 | 429 |
| 40002 | 验证码已过期 | 400 |
| 40003 | 验证码错误 | 400 |
| 40004 | 手机号已被绑定 | 409 |
| 40005 | 邮箱已被绑定 | 409 |
| 40006 | 当前密码错误 | 401 |
| 40007 | 新密码格式不符合要求 | 400 |
| 40008 | 设备不存在或已被移除 | 404 |
| 40009 | 不能移除当前设备 | 400 |
| 40010 | 注销条件不满足 | 400 |
| 40011 | 账号已在注销流程中 | 400 |
| 40012 | 需要图形验证码 | 400 |

---

## 实施计划

### 6.1 Phase 1 - 核心功能 (优先级P0)

**预计工时**: 5-7天

**任务列表**:
1. 验证码发送服务（手机+邮箱） - 2天
2. 手机号绑定功能 - 1天
3. 邮箱绑定功能 - 1天
4. 密码修改功能 - 1天
5. 密码重置功能 - 1天
6. 测试和文档 - 1天

### 6.2 Phase 2 - 设备管理 (优先级P1)

**预计工时**: 3-4天

**任务列表**:
1. 设备识别和记录 - 1.5天
2. 设备列表查询 - 0.5天
3. 设备移除功能 - 1天
4. 测试和优化 - 1天

### 6.3 Phase 3 - 账号注销 (优先级P1)

**预计工时**: 3-4天

**任务列表**:
1. 注销条件检查 - 1天
2. 注销申请流程 - 1天
3. 定时任务执行 - 1天
4. 测试和文档 - 1天

### 6.4 Phase 4 - 扩展功能 (优先级P2)

**预计工时**: 2-3天

**任务列表**:
1. 手机号/邮箱更换 - 1天
2. 解绑功能 - 0.5天
3. 安全审计增强 - 1天
4. 监控和告警 - 0.5天

---

## 附录

### A. 第三方服务

#### 短信服务提供商
- 阿里云短信服务
- 腾讯云短信服务
- 备用：云片网络

#### 邮件服务提供商
- SendGrid
- 阿里云邮件推送
- 备用：SMTP自建

### B. 参考文档

- [RFC 6238 - TOTP算法](https://tools.ietf.org/html/rfc6238)
- [OWASP密码存储指南](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
- [《个人信息安全规范》](https://www.gb688.cn/bzgk/gb/newGbInfo?hcno=4FFAA51D63BA21B9EE40C51DD3CC40BE)

---

**文档维护者**: 青羽后端架构团队  
**最后更新**: 2025-10-29

