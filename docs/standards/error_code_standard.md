# 错误码标准

**版本**: v1.0  
**创建日期**: 2026-02-26  
**状态**: ✅ 正式实施  

---

## 一、错误码设计原则

### 1.1 设计目标

1. **可读性**: 错误码应该容易理解
2. **扩展性**: 预留足够空间给未来需求
3. **一致性**: 全局统一的错误码体系
4. **可维护性**: 清晰的文档和管理

### 1.2 命名规则

```
错误码格式: ABCD (4位数字)

A - 错误类别 (1-5)
  1: 通用客户端错误
  2: 用户相关错误
  3: 业务逻辑错误
  4: 限流配额错误
  5: 服务器内部错误

B - 子类别 (0-9)
  按功能模块或错误类型划分

C-D - 具体错误编号 (00-99)
  该类别下的具体错误
```

### 1.3 命名约定

Go代码中的常量命名使用大驼峰 + 描述性名称：

```go
const (
    UserNotFound      ErrorCode = 2001
    TokenExpired      ErrorCode = 2008
    BookNotFound      ErrorCode = 3001
    InsufficientQuota ErrorCode = 3010
)
```

---

## 二、错误码分类详解

### 2.1 通用客户端错误 (1000-1099)

适用于所有模块的通用错误。

#### 2.1.1 参数验证 (1000-1019)

| 错误码 | 常量名 | HTTP状态码 | 说明 |
|--------|--------|-----------|------|
| 1001 | InvalidParams | 400 | 请求参数无效 |
| 1008 | MissingParam | 400 | 缺少必填参数 |
| 1009 | InvalidFormat | 400 | 参数格式无效 |
| 1010 | InvalidLength | 400 | 参数长度无效 |
| 1011 | InvalidType | 400 | 参数类型无效 |
| 1012 | OutOfRange | 400 | 参数超出范围 |
| 1013 | DuplicateField | 409 | 字段重复 |
| 1014 | UnknownField | 400 | 未知字段 |
| 1015 | ValidationFailed | 400 | 验证失败 |
| 1016 | InvalidValue | 400 | 值无效 |
| 1017 | InvalidOperation | 400 | 无效操作 |

#### 2.1.2 认证授权 (1020-1039)

| 错误码 | 常量名 | HTTP状态码 | 说明 |
|--------|--------|-----------|------|
| 1002 | Unauthorized | 401 | 未授权访问 |
| 1003 | Forbidden | 403 | 禁止访问 |

#### 2.1.3 资源相关 (1040-1059)

| 错误码 | 常量名 | HTTP状态码 | 说明 |
|--------|--------|-----------|------|
| 1004 | NotFound | 404 | 资源不存在 |
| 1005 | AlreadyExists | 409 | 资源已存在 |
| 1006 | Conflict | 409 | 冲突 |
| 1018 | ResourceGone | 410 | 资源已删除 |

#### 2.1.4 请求相关 (1060-1079)

| 错误码 | 常量名 | HTTP状态码 | 说明 |
|--------|--------|-----------|------|
| 1019 | MethodNotAllowed | 405 | 方法不允许 |
| 1020 | RequestTimeout | 408 | 请求超时 |

---

### 2.2 用户相关错误 (2000-2999)

用户认证、授权、账户管理相关错误。

#### 2.2.1 用户认证 (2000-2019)

| 错误码 | 常量名 | HTTP状态码 | 说明 |
|--------|--------|-----------|------|
| 2001 | UserNotFound | 404 | 用户不存在 |
| 2002 | InvalidCredentials | 401 | 无效凭证 |
| 2003 | UsernameAlreadyUsed | 409 | 用户名已被使用 |
| 2004 | EmailAlreadyUsed | 409 | 邮箱已被使用 |
| 2005 | EmailSendFailed | 500 | 邮件发送失败 |
| 2006 | InvalidCode | 400 | 验证码无效 |
| 2007 | CodeExpired | 400 | 验证码过期 |
| 2008 | TokenExpired | 401 | Token过期 |
| 2009 | TokenInvalid | 401 | Token无效 |
| 2010 | TokenFormatError | 400 | Token格式错误 |
| 2011 | TokenMissing | 401 | Token缺失 |
| 2012 | RefreshTokenExpired | 401 | Refresh Token过期 |
| 2013 | RefreshTokenInvalid | 401 | Refresh Token无效 |
| 2014 | PasswordTooWeak | 400 | 密码强度不足 |
| 2015 | AccountLocked | 403 | 账户已锁定 |
| 2016 | AccountDisabled | 403 | 账户已禁用 |
| 2017 | TokenRevoked | 401 | Token已被撤销 |
| 2018 | SessionExpired | 401 | 会话过期 |
| 2019 | TooManyAttempts | 429 | 尝试次数过多 |
| 2020 | AccountNotVerified | 403 | 账户未验证 |

#### 2.2.2 邮箱和手机 (2020-2039)

| 错误码 | 常量名 | HTTP状态码 | 说明 |
|--------|--------|-----------|------|
| 2021 | PhoneAlreadyUsed | 409 | 手机号已使用 |
| 2022 | InvalidPhoneFormat | 400 | 手机号格式无效 |
| 2023 | EmailNotVerified | 403 | 邮箱未验证 |
| 2024 | PhoneNotVerified | 403 | 手机号未验证 |
| 2025 | SmsSendFailed | 500 | 短信发送失败 |
| 2026 | EmailSendTooFrequent | 429 | 邮件发送过于频繁 |
| 2027 | SmsSendTooFrequent | 429 | 短信发送过于频繁 |

#### 2.2.3 评分相关 (2500-2599)

| 错误码 | 常量名 | HTTP状态码 | 说明 |
|--------|--------|-----------|------|
| 2501 | RatingNotFound | 404 | 评分不存在 |
| 2502 | RatingInvalid | 400 | 评分值无效 |
| 2503 | RatingAlreadyExists | 409 | 用户已评分 |
| 2504 | RatingUnauthorized | 403 | 无权操作此评分 |
| 2505 | RatingTargetNotFound | 404 | 评分目标不存在 |

---

### 2.3 业务逻辑错误 (3000-3999)

业务逻辑相关的错误，按模块划分。

#### 2.3.1 书籍相关 (3000-3039)

| 错误码 | 常量名 | HTTP状态码 | 说明 |
|--------|--------|-----------|------|
| 3001 | BookNotFound | 404 | 书籍不存在 |
| 3004 | BookAlreadyExists | 409 | 书籍已存在 |
| 3005 | InvalidBookStatus | 400 | 书籍状态无效 |
| 3006 | BookDeleted | 410 | 书籍已删除 |

#### 2.3.2 章节相关 (3040-3069)

| 错误码 | 常量名 | HTTP状态码 | 说明 |
|--------|--------|-----------|------|
| 3002 | ChapterNotFound | 404 | 章节不存在 |
| 3007 | ChapterAlreadyExists | 409 | 章节已存在 |
| 3008 | InvalidChapterStatus | 400 | 章节状态无效 |
| 3009 | ChapterDeleted | 410 | 章节已删除 |

#### 2.3.3 财务相关 (3070-3099)

| 错误码 | 常量名 | HTTP状态码 | 说明 |
|--------|--------|-----------|------|
| 3003 | InsufficientBalance | 400 | 余额不足 |
| 3010 | InsufficientQuota | 400 | 配额不足 |
| 3011 | WalletFrozen | 403 | 钱包已冻结 |
| 3017 | TransactionFailed | 500 | 交易失败 |

#### 2.3.4 内容相关 (3100-3129)

| 错误码 | 常量名 | HTTP状态码 | 说明 |
|--------|--------|-----------|------|
| 3012 | ContentNotPublished | 403 | 内容未发布 |
| 3013 | ChapterLocked | 403 | 章节已锁定 |
| 3018 | ContentLocked | 403 | 内容已锁定 |
| 3019 | ContentDeleted | 410 | 内容已删除 |

#### 2.3.5 内容审核 (3130-3159)

| 错误码 | 常量名 | HTTP状态码 | 说明 |
|--------|--------|-----------|------|
| 3014 | ContentPendingReview | 202 | 内容待审核 |
| 3015 | ContentRejected | 403 | 内容被拒绝 |
| 3016 | ContentViolation | 403 | 内容违规 |

#### 2.3.6 角色相关 (3160-3179)

| 错误码 | 常量名 | HTTP状态码 | 说明 |
|--------|--------|-----------|------|
| 3020 | CharacterNotFound | 404 | 角色不存在 |
| 3021 | InvalidCharacterData | 400 | 角色数据无效 |

#### 2.3.7 评论相关 (3180-3199)

| 错误码 | 常量名 | HTTP状态码 | 说明 |
|--------|--------|-----------|------|
| 3022 | ReviewNotFound | 404 | 评论不存在 |

#### 2.3.8 收藏和关注 (3200-3219)

| 错误码 | 常量名 | HTTP状态码 | 说明 |
|--------|--------|-----------|------|
| 3023 | CollectionNotFound | 404 | 收藏不存在 |
| 3024 | AlreadyCollected | 409 | 已收藏 |
| 3025 | AlreadyFollowed | 409 | 已关注 |

#### 2.3.9 Writer模块 (3300-3399)

| 错误码 | 常量名 | HTTP状态码 | 说明 |
|--------|--------|-----------|------|
| 3301 | ProjectNotFound | 404 | 项目不存在 |
| 3302 | DocumentNotFound | 404 | 文档不存在 |
| 3303 | CommentNotFound | 404 | 批注不存在 |
| 3304 | InvalidProjectInput | 400 | 项目输入无效 |
| 3305 | VersionConflict | 409 | 版本冲突 |
| 3306 | PublishFailed | 500 | 发布失败 |

#### 2.3.10 Reader模块 (3400-3499)

| 错误码 | 常量名 | HTTP状态码 | 说明 |
|--------|--------|-----------|------|
| 3401 | ProgressNotFound | 404 | 阅读进度不存在 |
| 3402 | AnnotationNotFound | 404 | 标注不存在 |
| 3403 | InvalidProgress | 400 | 无效的阅读进度 |
| 3404 | BookAccessDenied | 403 | 书籍访问被拒绝 |

---

### 2.4 限流配额错误 (4000-4099)

频率限制和配额相关错误。

| 错误码 | 常量名 | HTTP状态码 | 说明 |
|--------|--------|-----------|------|
| 4000 | RateLimitExceeded | 429 | 频率限制超出 |
| 4001 | DailyLimitExceeded | 429 | 每日限制超出 |
| 4002 | HourlyLimitExceeded | 429 | 每小时限制超出 |
| 4003 | MinuteLimitExceeded | 429 | 每分钟限制超出 |
| 4004 | UploadLimitExceeded | 429 | 上传限制超出 |
| 4005 | StorageLimitExceeded | 507 | 存储限制超出 |
| 4006 | ApiQuotaExceeded | 429 | API配额超出 |
| 4007 | ConcurrentLimitExceeded | 429 | 并发限制超出 |
| 4008 | RateLimitLogin | 429 | 登录频率限制 |
| 4009 | RateLimitEmailSend | 429 | 邮件发送频率限制 |
| 4010 | RateLimitSmsSend | 429 | 短信发送频率限制 |

---

### 2.5 服务器内部错误 (5000-5099)

系统、数据库、外部服务相关错误。

#### 2.5.1 系统错误 (5000-5019)

| 错误码 | 常量名 | HTTP状态码 | 说明 |
|--------|--------|-----------|------|
| 5000 | InternalError | 500 | 内部错误 |
| 5001 | DatabaseError | 500 | 数据库错误 |
| 5002 | ServiceUnavailable | 503 | 服务不可用 |
| 5003 | RedisError | 500 | Redis错误 |

#### 2.5.2 外部服务 (5020-5039)

| 错误码 | 常量名 | HTTP状态码 | 说明 |
|--------|--------|-----------|------|
| 5004 | ExternalAPIError | 502 | 外部API错误 |
| 5005 | CacheError | 500 | 缓存错误 |
| 5006 | QueueError | 500 | 队列错误 |
| 5007 | StorageError | 500 | 存储错误 |
| 5008 | NetworkError | 500 | 网络错误 |
| 5009 | ConfigurationError | 500 | 配置错误 |
| 5010 | DependencyError | 502 | 依赖错误 |
| 5011 | TimeoutError | 504 | 超时错误 |
| 5012 | OverloadedError | 503 | 过载错误 |
| 5013 | MaintenanceError | 503 | 维护中 |

#### 2.5.3 数据库详细错误 (5040-5059)

| 错误码 | 常量名 | HTTP状态码 | 说明 |
|--------|--------|-----------|------|
| 5014 | DatabaseConnectionFailed | 500 | 数据库连接失败 |
| 5015 | DatabaseQueryTimeout | 500 | 数据库查询超时 |
| 5016 | DatabaseTransactionFailed | 500 | 数据库事务失败 |

---

## 三、错误码使用规范

### 3.1 新增错误码流程

1. **评估需求**: 确认是否真的需要新的错误码
2. **选择类别**: 根据错误类型选择合适的类别
3. **分配编号**: 在对应范围内选择未使用的编号
4. **更新文档**: 更新本文档和 codes.go
5. **添加测试**: 添加对应的错误创建和转换测试

### 3.2 错误码分配规则

1. **预留空间**: 每个子类别预留一定空间
2. **模块化**: 同一模块的错误码集中管理
3. **可扩展**: 预留足够空间给未来需求

### 3.3 已废弃错误码

| 旧错误码 | 新错误码 | 废弃原因 |
|---------|---------|----------|
| 4291 | 4002 | 小时级限制超出，统一为HourlyLimitExceeded |

---

## 四、错误码与HTTP状态码映射

### 4.1 映射原则

| HTTP状态码 | 使用场景 | 错误码范围 |
|-----------|----------|-----------|
| 200 OK | 成功 | 0 |
| 202 Accepted | 异步处理中 | 3014 |
| 400 Bad Request | 客户端请求错误 | 1000-1099, 部分3000 |
| 401 Unauthorized | 未认证 | 1002, 2008-2018 |
| 403 Forbidden | 禁止访问 | 1003, 部分2000, 部分3000 |
| 404 Not Found | 资源不存在 | 1004, 2001, 3001-3002, 等 |
| 405 Method Not Allowed | 方法不允许 | 1019 |
| 408 Request Timeout | 请求超时 | 1020 |
| 409 Conflict | 资源冲突 | 1005-1006, 2003-2004 |
| 410 Gone | 资源已删除 | 1018, 3006, 3009 |
| 429 Too Many Requests | 请求过多 | 4000-4010 |
| 500 Internal Server Error | 服务器错误 | 5000-5016 |
| 502 Bad Gateway | 外部服务错误 | 5004, 5010 |
| 503 Service Unavailable | 服务不可用 | 5002, 5012-5013 |
| 504 Gateway Timeout | 超时 | 5011 |
| 507 Insufficient Storage | 存储不足 | 4005 |

### 4.2 特殊映射

有些错误可能根据上下文返回不同的HTTP状态码：

```go
// 示例：内容审核状态
switch content.Status {
case "pending":
    return http.StatusAccepted  // 202, 内容待审核
case "rejected":
    return http.StatusForbidden  // 403, 内容被拒绝
case "published":
    return http.StatusOK        // 200, 成功
}
```

---

## 五、错误码管理

### 5.1 维护责任人

- **技术负责人**: 审核和批准新错误码
- **架构师**: 确保错误码体系一致性
- **开发团队**: 实施和维护

### 5.2 版本控制

错误码定义随项目版本管理：
- 破坏性变更需要大版本号升级
- 新增错误码可以使用小版本号
- 废弃错误码应该至少保留一个版本

### 5.3 文档更新

每次修改错误码时：
1. 更新本文档
2. 更新 `pkg/errors/codes.go`
3. 更新相关测试
4. 通知相关开发人员

---

## 六、最佳实践

### 6.1 选择合适的错误码

```go
// ✅ 推荐: 使用通用错误码
return factory.NotFoundError("用户", userID)

// ❌ 避免: 创建过于具体的错误码
const ErrUserNotFoundByIdInUserService = 29991
```

### 6.2 错误消息

```go
// ✅ 推荐: 清晰明确的错误消息
"用户名长度必须在3-20个字符之间"

// ❌ 避免: 模糊的错误消息
"输入无效"
```

### 6.3 错误详情

```go
// ✅ 推荐: 提供有用的错误详情
err := factory.ValidationError(
    "1001",
    "用户名格式无效",
    "用户名 'ab' 过短，至少需要3个字符",
)

// ❌ 避免: 不提供错误详情或暴露敏感信息
err := factory.ValidationError(
    "1001",
    "错误",
    "数据库查询失败: SELECT * FROM users WHERE username='ab' AND password='...'",  // 不要暴露SQL
)
```

---

## 七、相关文档

- [错误处理指南](./error_handling_guide.md)
- [pkg/errors包文档](../../pkg/errors/)
- [API设计规范](./api/API设计规范.md)
