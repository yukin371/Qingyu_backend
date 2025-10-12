# 共享服务数据模型

> 共享底层服务的数据模型定义

---

## 📋 模块列表

### 1. Auth 模块 ✅

**文件**:
- `auth/role.go` - 角色与权限模型
- `auth/session.go` - 会话与Token黑名单

**核心模型**:
- `Role` - 角色定义
- `Permission` - 权限定义
- `UserRole` - 用户角色关联
- `Session` - 会话信息
- `TokenBlacklist` - Token黑名单

---

### 2. Wallet 模块 ✅

**文件**:
- `wallet/wallet.go` - 钱包、交易、提现模型

**核心模型**:
- `Wallet` - 钱包信息
- `Transaction` - 交易记录
- `WithdrawRequest` - 提现申请

---

### 3. Recommendation 模块 ✅

**文件**:
- `recommendation/recommendation.go` - 推荐与行为模型

**核心模型**:
- `UserBehavior` - 用户行为记录
- `RecommendedItem` - 推荐项（主要存缓存）
- `UserProfile` - 用户画像（可选）

---

### 4. Storage 模块 ✅

**文件**:
- `storage/file.go` - 文件存储模型

**核心模型**:
- `FileInfo` - 文件元数据
- `FileAccess` - 文件访问权限（可选）

---

### 5. Admin 模块 ✅

**文件**:
- `admin/admin.go` - 审核与日志模型

**核心模型**:
- `AuditRecord` - 审核记录
- `AdminLog` - 管理员操作日志

---

## 🗄️ 数据库集合

### MongoDB 集合

```
qingyu_db (数据库)
├── roles                    # 角色定义
├── user_roles              # 用户角色关联（可选，也可存在users.roles）
├── wallets                 # 钱包信息
├── transactions            # 交易记录
├── withdraw_requests       # 提现申请
├── user_behaviors          # 用户行为
├── files                   # 文件元数据
├── audit_records           # 审核记录
└── admin_logs              # 管理员日志
```

### Redis 数据

```
# 会话管理
session:{session_id}        → Session JSON

# Token黑名单
token:blacklist:{token}     → TokenBlacklist JSON

# 推荐缓存
recommend:user:{user_id}    → RecommendedItem[] JSON
recommend:hot:{item_type}   → RecommendedItem[] JSON
```

---

## 📚 使用示例

### 导入模型

```go
import (
    authModel "Qingyu_backend/models/shared/auth"
    walletModel "Qingyu_backend/models/shared/wallet"
    recommendationModel "Qingyu_backend/models/shared/recommendation"
    storageModel "Qingyu_backend/models/shared/storage"
    adminModel "Qingyu_backend/models/shared/admin"
)
```

### 创建角色

```go
role := &authModel.Role{
    Name: authModel.RoleReader,
    Description: "普通读者",
    Permissions: []string{
        authModel.PermBookRead,
        authModel.PermUserRead,
    },
    IsSystem: true,
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
}
```

### 创建钱包

```go
wallet := &walletModel.Wallet{
    UserID: "user123",
    Balance: 0.0,
    Frozen: false,
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
}
```

---

## 🔍 索引建议

### MongoDB 索引

```javascript
// roles 集合
db.roles.createIndex({ "name": 1 }, { unique: true })

// wallets 集合
db.wallets.createIndex({ "user_id": 1 }, { unique: true })

// transactions 集合
db.transactions.createIndex({ "user_id": 1, "transaction_time": -1 })
db.transactions.createIndex({ "order_no": 1 }, { unique: true })
db.transactions.createIndex({ "status": 1, "created_at": -1 })

// withdraw_requests 集合
db.withdraw_requests.createIndex({ "user_id": 1, "created_at": -1 })
db.withdraw_requests.createIndex({ "status": 1, "created_at": -1 })
db.withdraw_requests.createIndex({ "order_no": 1 }, { unique: true })

// user_behaviors 集合
db.user_behaviors.createIndex({ "user_id": 1, "created_at": -1 })
db.user_behaviors.createIndex({ "item_id": 1, "action_type": 1 })
db.user_behaviors.createIndex({ "created_at": -1 })

// files 集合
db.files.createIndex({ "user_id": 1, "created_at": -1 })
db.files.createIndex({ "md5": 1 })
db.files.createIndex({ "category": 1, "created_at": -1 })

// audit_records 集合
db.audit_records.createIndex({ "content_id": 1, "content_type": 1 }, { unique: true })
db.audit_records.createIndex({ "status": 1, "created_at": -1 })

// admin_logs 集合
db.admin_logs.createIndex({ "admin_id": 1, "created_at": -1 })
db.admin_logs.createIndex({ "operation": 1, "created_at": -1 })
```

---

*模型定义完成 ✅*
