# VIP权限服务设计

> **版本**: v1.0  
> **创建日期**: 2025-10-21  
> **状态**: ✅ 已实现，补充设计文档

---

## 1. 设计概述

### 1.1 业务价值

- VIP章节解锁权限验证
- 会员等级管理
- 会员权益控制

### 1.2 实现情况

**已实现**：
- VIP权限中间件：`middleware/vip_permission.go`
- VIP权限验证逻辑

---

## 2. VIP等级设计

### 2.1 会员等级

```go
type VIPLevel string

const (
    VIPNone    VIPLevel = "none"       // 非会员
    VIPMonthly VIPLevel = "monthly"    // 月度会员
    VIPYearly  VIPLevel = "yearly"     // 年度会员
    VIPLifetime VIPLevel = "lifetime"  // 终身会员
)

type VIPInfo struct {
    Level     VIPLevel
    ExpiresAt *time.Time  // 过期时间（lifetime为nil）
    IsActive  bool
}
```

### 2.2 会员权益

```go
var VIPPrivileges = map[VIPLevel][]string{
    VIPMonthly: {
        "chapter:unlock",      // 解锁VIP章节
        "reading:ad_free",     // 无广告阅读
        "book:download",       // 离线下载
    },
    VIPYearly: {
        "chapter:unlock",
        "reading:ad_free",
        "book:download",
        "ai:priority",         // AI功能优先
        "storage:extra",       // 额外存储空间
    },
    VIPLifetime: {
        "chapter:unlock",
        "reading:ad_free",
        "book:download",
        "ai:priority",
        "storage:extra",
        "badge:lifetime",      // 终身会员徽章
    },
}
```

---

## 3. 权限验证设计

### 3.1 VIP中间件

```go
// middleware/vip_permission.go
func VIPPermission() gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, _ := c.Get("userId")
        
        // 检查VIP状态
        vipInfo := getVIPInfo(userID.(string))
        
        if !vipInfo.IsActive {
            response.Error(c, http.StatusForbidden, "需要VIP权限", "")
            c.Abort()
            return
        }
        
        // 检查是否过期
        if vipInfo.ExpiresAt != nil && time.Now().After(*vipInfo.ExpiresAt) {
            response.Error(c, http.StatusForbidden, "VIP已过期", "")
            c.Abort()
            return
        }
        
        c.Set("vipLevel", vipInfo.Level)
        c.Next()
    }
}
```

### 3.2 章节解锁验证

```go
func CheckChapterAccess(userID string, chapterID string) (bool, error) {
    // 1. 获取章节信息
    chapter := getChapter(chapterID)
    
    // 2. 如果是免费章节，直接允许
    if !chapter.IsVIP {
        return true, nil
    }
    
    // 3. 检查VIP状态
    vipInfo := getVIPInfo(userID)
    if !vipInfo.IsActive {
        return false, errors.New("需要VIP权限")
    }
    
    // 4. 检查权益是否包含chapter:unlock
    if hasPrivilege(vipInfo.Level, "chapter:unlock") {
        return true, nil
    }
    
    return false, errors.New("权限不足")
}
```

---

## 4. VIPService设计

### 4.1 核心方法

```go
type VIPService interface {
    // GetVIPInfo 获取VIP信息
    GetVIPInfo(ctx context.Context, userID string) (*VIPInfo, error)
    
    // ActivateVIP 激活VIP
    ActivateVIP(ctx context.Context, userID string, level VIPLevel, duration time.Duration) error
    
    // RenewVIP 续费VIP
    RenewVIP(ctx context.Context, userID string, duration time.Duration) error
    
    // CancelVIP 取消VIP
    CancelVIP(ctx context.Context, userID string) error
    
    // CheckPrivilege 检查权益
    CheckPrivilege(ctx context.Context, userID string, privilege string) (bool, error)
}
```

### 4.2 VIP激活逻辑

```go
func (s *VIPService) ActivateVIP(ctx context.Context, userID string, level VIPLevel, duration time.Duration) error {
    expiresAt := time.Now().Add(duration)
    
    vipInfo := &VIPInfo{
        UserID:    userID,
        Level:     level,
        ExpiresAt: &expiresAt,
        IsActive:  true,
        ActivatedAt: time.Now(),
    }
    
    // 保存VIP信息
    s.vipRepo.Create(ctx, vipInfo)
    
    // 发布事件
    s.eventBus.Publish("vip.activated", VIPActivatedEvent{
        UserID: userID,
        Level:  level,
    })
    
    return nil
}
```

---

## 5. 与v2.1架构的关系

```
Reading Module
  └─ VIP Permission System
      ├─ VIPInfo (会员信息)
      ├─ VIPService (会员服务)
      └─ VIP Middleware (权限中间件)
```

**集成点**：
- 章节阅读前验证VIP权限
- 书城展示时区分VIP/非VIP章节
- 与钱包系统集成（购买VIP）

---

## 6. 实现参考

**代码文件**：
- `middleware/vip_permission.go` - VIP中间件
- `service/shared/vip_service.go`（待创建） - VIP服务
- `models/shared/vip_info.go`（待创建） - VIP数据模型

---

**文档状态**: ✅ 已完成  
**优先级**: P0

