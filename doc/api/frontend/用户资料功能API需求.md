# 用户资料功能API需求文档

**文档版本**: v1.0  
**创建日期**: 2025-10-29  
**状态**: 待实现  
**优先级**: P2

---

## 📋 目录

1. [概述](#概述)
2. [API列表](#api列表)
3. [详细设计](#详细设计)
4. [数据模型](#数据模型)
5. [实施计划](#实施计划)

---

## 概述

### 功能描述

用户资料功能为用户提供个人信息管理、头像上传、社交关系管理等能力，完善用户profile体系。

### 业务价值

- 丰富用户个人信息展示
- 支持用户社交互动（关注/粉丝）
- 提升平台活跃度和用户粘性
- 为推荐算法提供数据支持

---

## API列表

### 2.1 头像管理

| 序号 | API路径 | 方法 | 说明 | 优先级 |
|-----|---------|------|------|--------|
| 1 | `/api/v1/user/profile/avatar/upload` | POST | 上传头像 | P0 |
| 2 | `/api/v1/user/profile/avatar` | PUT | 更换头像（使用已上传的图片） | P1 |
| 3 | `/api/v1/user/profile/avatar` | DELETE | 删除头像（恢复默认） | P2 |

### 2.2 社交关系

| 序号 | API路径 | 方法 | 说明 | 优先级 |
|-----|---------|------|------|--------|
| 4 | `/api/v1/user/follow/:user_id` | POST | 关注用户 | P0 |
| 5 | `/api/v1/user/unfollow/:user_id` | DELETE | 取消关注 | P0 |
| 6 | `/api/v1/user/following` | GET | 获取关注列表 | P0 |
| 7 | `/api/v1/user/followers` | GET | 获取粉丝列表 | P0 |
| 8 | `/api/v1/user/:user_id/following` | GET | 获取指定用户的关注列表 | P1 |
| 9 | `/api/v1/user/:user_id/followers` | GET | 获取指定用户的粉丝列表 | P1 |
| 10 | `/api/v1/user/follow/status/:user_id` | GET | 查询关注状态 | P1 |
| 11 | `/api/v1/user/mutual-follow` | GET | 获取互关列表 | P2 |

---

## 详细设计

### 3.1 上传头像

#### 基本信息
- **路径**: `/api/v1/user/profile/avatar/upload`
- **方法**: `POST`
- **认证**: 必须
- **Content-Type**: `multipart/form-data`

#### 请求参数

```
Form Data:
- file: 图片文件（必填）
```

**文件要求**:
- 格式：JPG, PNG, GIF, WEBP
- 大小：最大5MB
- 尺寸：建议正方形，最小100x100px

#### 响应示例

```json
{
  "code": 200,
  "message": "头像上传成功",
  "data": {
    "avatar_id": "xxx",
    "url": "https://cdn.qingyu.com/avatars/xxx.jpg",
    "thumbnail_url": "https://cdn.qingyu.com/avatars/xxx_thumb.jpg",
    "width": 800,
    "height": 800,
    "size": 125678,
    "uploaded_at": "2025-10-29T10:30:00Z"
  }
}
```

#### 业务规则

1. **图片处理**:
   - 自动裁剪为正方形
   - 生成多个尺寸：原图、大图(400x400)、中图(200x200)、小图(100x100)
   - 压缩优化
   - 添加水印（可选）

2. **存储策略**:
   - 上传到OSS（阿里云/腾讯云）
   - CDN加速
   - 旧头像保留7天后删除

3. **安全检查**:
   - 图片内容审核（鉴黄、暴恐、违禁）
   - 病毒扫描
   - 文件类型验证

4. **限制**:
   - 每日最多上传10次
   - 单个文件最大5MB

#### 技术实现

```go
// 伪代码
func UploadAvatar(file multipart.File) (*Avatar, error) {
    // 1. 验证文件
    if err := validateImageFile(file); err != nil {
        return nil, err
    }
    
    // 2. 图片处理
    img, err := processImage(file)
    if err != nil {
        return nil, err
    }
    
    // 3. 内容审核
    if err := contentModeration(img); err != nil {
        return nil, err
    }
    
    // 4. 上传到OSS
    url, err := uploadToOSS(img)
    if err != nil {
        return nil, err
    }
    
    // 5. 生成缩略图
    thumbnails, err := generateThumbnails(img)
    if err != nil {
        return nil, err
    }
    
    // 6. 保存记录
    avatar := &Avatar{
        URL: url,
        Thumbnails: thumbnails,
        // ...
    }
    return avatar, nil
}
```

---

### 3.2 更换头像

#### 基本信息
- **路径**: `/api/v1/user/profile/avatar`
- **方法**: `PUT`
- **认证**: 必须

#### 请求参数

```json
{
  "avatar_id": "xxx"  // 已上传的头像ID
}
```

#### 响应示例

```json
{
  "code": 200,
  "message": "头像更换成功",
  "data": {
    "avatar_url": "https://cdn.qingyu.com/avatars/xxx.jpg"
  }
}
```

#### 业务规则

- 只能使用自己上传的头像
- 更新用户资料中的头像字段
- 通知相关服务更新缓存

---

### 3.3 关注用户

#### 基本信息
- **路径**: `/api/v1/user/follow/:user_id`
- **方法**: `POST`
- **认证**: 必须

#### Path参数

- `user_id`: 要关注的用户ID

#### 响应示例

```json
{
  "code": 200,
  "message": "关注成功",
  "data": {
    "user_id": "target_user_id",
    "username": "targetuser",
    "avatar": "https://cdn.qingyu.com/avatars/xxx.jpg",
    "is_mutual": false,  // 是否互关
    "followed_at": "2025-10-29T10:30:00Z"
  }
}
```

#### 业务规则

1. **关注限制**:
   - 不能关注自己
   - 不能重复关注
   - 每日最多关注200人
   - 总关注数上限5000人

2. **关注后操作**:
   - 创建关注关系记录
   - 增加被关注者粉丝数
   - 增加关注者关注数
   - 发送关注通知（如果对方开启）
   - 触发推荐算法更新

3. **互关检测**:
   - 检查对方是否也关注了自己
   - 标记互关状态

#### 数据结构

```go
type Follow struct {
    ID          string    `bson:"_id"`
    FollowerID  string    `bson:"follower_id"`  // 关注者
    FollowingID string    `bson:"following_id"` // 被关注者
    IsMutual    bool      `bson:"is_mutual"`    // 是否互关
    CreatedAt   time.Time `bson:"created_at"`
}
```

---

### 3.4 取消关注

#### 基本信息
- **路径**: `/api/v1/user/unfollow/:user_id`
- **方法**: `DELETE`
- **认证**: 必须

#### Path参数

- `user_id`: 要取消关注的用户ID

#### 响应示例

```json
{
  "code": 200,
  "message": "已取消关注",
  "data": {
    "user_id": "target_user_id"
  }
}
```

#### 业务规则

1. **取消关注操作**:
   - 删除关注关系记录
   - 减少被关注者粉丝数
   - 减少关注者关注数
   - 更新对方的互关状态（如果是互关）

2. **静默操作**:
   - 不发送通知
   - 不在对方动态中显示

---

### 3.5 获取关注列表

#### 基本信息
- **路径**: `/api/v1/user/following`
- **方法**: `GET`
- **认证**: 必须

#### 请求参数

```
Query参数:
- page: 页码（默认1）
- page_size: 每页数量（默认20，最大100）
- keyword: 搜索关键词（可选）
- sort: 排序方式（latest, earliest）
```

#### 响应示例

```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "users": [
      {
        "user_id": "xxx",
        "username": "authorname",
        "nickname": "作者昵称",
        "avatar": "https://cdn.qingyu.com/avatars/xxx.jpg",
        "bio": "个人简介",
        "role": "writer",  // writer, reader
        "is_mutual": true,  // 是否互关
        "followed_at": "2025-10-29T10:30:00Z",
        "stats": {
          "followers_count": 1000,
          "following_count": 500,
          "works_count": 10
        }
      }
    ],
    "total": 150,
    "page": 1,
    "page_size": 20
  }
}
```

#### 业务规则

- 按关注时间倒序排列
- 支持关键词搜索（用户名、昵称）
- 显示互关状态
- 显示基本统计数据

---

### 3.6 获取粉丝列表

#### 基本信息
- **路径**: `/api/v1/user/followers`
- **方法**: `GET`
- **认证**: 必须

#### 请求参数

同关注列表

#### 响应示例

同关注列表，额外包含：
- `is_following_back`: 是否回关

#### 业务规则

- 按关注时间倒序
- 标识是否已回关
- 新增粉丝标记（24小时内）

---

### 3.7 查询关注状态

#### 基本信息
- **路径**: `/api/v1/user/follow/status/:user_id`
- **方法**: `GET`
- **认证**: 必须

#### 响应示例

```json
{
  "code": 200,
  "message": "查询成功",
  "data": {
    "user_id": "xxx",
    "is_following": true,  // 我是否关注了对方
    "is_follower": false,  // 对方是否关注了我
    "is_mutual": false     // 是否互关
  }
}
```

---

### 3.8 获取互关列表

#### 基本信息
- **路径**: `/api/v1/user/mutual-follow`
- **方法**: `GET`
- **认证**: 必须

#### 请求参数

```
Query参数:
- page: 页码
- page_size: 每页数量
```

#### 响应示例

```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "users": [
      {
        "user_id": "xxx",
        "username": "mutualuser",
        "avatar": "https://cdn.qingyu.com/avatars/xxx.jpg",
        "mutual_since": "2025-10-29T10:30:00Z"  // 互关时间
      }
    ],
    "total": 50
  }
}
```

---

## 数据模型

### 4.1 关注关系

```go
type Follow struct {
    ID          string    `bson:"_id"`
    FollowerID  string    `bson:"follower_id"`   // 关注者ID
    FollowingID string    `bson:"following_id"`  // 被关注者ID
    IsMutual    bool      `bson:"is_mutual"`     // 是否互关
    CreatedAt   time.Time `bson:"created_at"`    // 关注时间
}

// 索引
// 1. follower_id + following_id (unique)
// 2. follower_id + created_at
// 3. following_id + created_at
// 4. follower_id + is_mutual
```

### 4.2 用户统计

```go
type UserStats struct {
    UserID          string    `bson:"user_id"`
    FollowersCount  int64     `bson:"followers_count"`   // 粉丝数
    FollowingCount  int64     `bson:"following_count"`   // 关注数
    MutualCount     int64     `bson:"mutual_count"`      // 互关数
    UpdatedAt       time.Time `bson:"updated_at"`
}
```

### 4.3 头像记录

```go
type Avatar struct {
    ID           string    `bson:"_id"`
    UserID       string    `bson:"user_id"`
    URL          string    `bson:"url"`           // 原图URL
    ThumbnailURL string    `bson:"thumbnail_url"` // 缩略图URL
    Width        int       `bson:"width"`
    Height       int       `bson:"height"`
    Size         int64     `bson:"size"`          // 文件大小（字节）
    Format       string    `bson:"format"`        // 图片格式
    Status       string    `bson:"status"`        // active, deleted
    UploadedAt   time.Time `bson:"uploaded_at"`
}
```

---

## 实施计划

### 5.1 Phase 1 - 头像上传 (P0)

**预计工时**: 2-3天

**任务列表**:
1. 图片上传和处理 - 1天
2. OSS集成 - 0.5天
3. 内容审核集成 - 0.5天
4. 测试和优化 - 1天

**技术栈**:
- OSS: 阿里云OSS
- 图片处理: github.com/disintegration/imaging
- 内容审核: 阿里云内容安全

### 5.2 Phase 2 - 社交关系 (P0)

**预计工时**: 3-4天

**任务列表**:
1. 关注/取消关注 - 1天
2. 关注列表查询 - 1天
3. 粉丝列表查询 - 0.5天
4. 关注状态查询 - 0.5天
5. 互关列表 - 0.5天
6. 性能优化和缓存 - 1天

**优化方案**:
- Redis缓存关注关系
- 关注列表分页优化
- 计数器缓存

### 5.3 Phase 3 - 扩展功能 (P2)

**预计工时**: 2天

**任务列表**:
1. 头像历史管理 - 0.5天
2. 批量关注操作 - 0.5天
3. 关注推荐 - 0.5天
4. 监控和告警 - 0.5天

---

## 附录

### A. OSS配置

```yaml
oss:
  endpoint: oss-cn-beijing.aliyuncs.com
  bucket: qingyu-avatars
  access_key_id: ${OSS_ACCESS_KEY_ID}
  access_key_secret: ${OSS_ACCESS_KEY_SECRET}
  cdn_domain: cdn.qingyu.com
  upload_dir: avatars/
```

### B. 图片规格

| 规格 | 尺寸 | 用途 |
|-----|------|------|
| 原图 | 原始尺寸 | 详情页展示 |
| 大图 | 400x400 | 个人主页 |
| 中图 | 200x200 | 评论区头像 |
| 小图 | 100x100 | 列表缩略图 |
| 微小 | 50x50 | 通知头像 |

### C. 内容审核策略

- **鉴黄**: 阻止级别 >= 90分
- **暴恐**: 阻止级别 >= 80分
- **违禁**: 阻止级别 >= 70分
- **人工复审**: 60-70分区间

---

**文档维护者**: 青羽后端架构团队  
**最后更新**: 2025-10-29

