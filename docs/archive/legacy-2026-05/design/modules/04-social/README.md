# 04 - 社交互动模块

> **模块编号**: 04
> **模块名称**: Social Interaction
> **负责功能**: 评论、点赞、关注、收藏等社交互动
> **完成度**: 🟡 40%

## 📋 目录结构

```
社交互动模块/
├── api/v1/
│   └── social/                   # 社交API
│       ├── comment_api.go       # 评论系统
│       ├── like_api.go          # 点赞系统
│       ├── follow_api.go        # 关注系统
│       └── collection_api.go    # 收藏系统
├── service/social/               # 社交服务层
│   ├── comment_service.go       # 评论服务
│   ├── like_service.go          # 点赞服务
│   ├── follow_service.go        # 关注服务
│   └── collection_service.go    # 收藏服务
├── repository/interfaces/social/ # 仓储接口
├── repository/mongodb/social/    # MongoDB仓储实现
└── models/social/                # 数据模型
    ├── comment.go                # 评论
    ├── like.go                   # 点赞
    ├── follow.go                 # 关注
    └── collection.go             # 收藏
```

## 🎯 核心功能

### 1. 评论系统

- **发布评论**: 对书籍、章节、书评进行评论
- **回复评论**: 支持评论嵌套回复
- **评论点赞**: 对评论进行点赞
- **评论举报**: 违规评论举报
- **评论审核**: 敏感词过滤、人工审核

### 2. 点赞系统

- **内容点赞**: 对书籍、章节点赞
- **评论点赞**: 对评论点赞
- **点赞状态**: 查询点赞状态
- **取消点赞**: 取消已点赞内容

### 3. 关注系统

- **关注作者**: 关注感兴趣的作者
- **粉丝管理**: 查看粉丝列表
- **关注列表**: 查看已关注作者
- **取消关注**: 取消关注作者

### 4. 收藏系统

- **收藏内容**: 收藏书籍、章节
- **收藏夹管理**: 创建和管理收藏夹
- **收藏分类**: 按标签分类收藏
- **收藏分享**: 分享收藏内容

## 📊 数据模型

### Comment (评论)

```go
type Comment struct {
    ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    UserID          primitive.ObjectID   `bson:"user_id" json:"userId"`
    TargetType      TargetType          `bson:"target_type" json:"targetType"`
    TargetID        primitive.ObjectID   `bson:"target_id" json:"targetId"`
    ParentID        *primitive.ObjectID  `bson:"parent_id,omitempty" json:"parentId,omitempty"`

    // 评论内容
    Content         string               `bson:"content" json:"content"`
    RichContent     string               `bson:"rich_content" json:"richContent"`

    // 统计信息
    LikeCount       int                  `bson:"like_count" json:"likeCount"`
    ReplyCount      int                  `bson:"reply_count" json:"replyCount"`

    // 状态
    Status          CommentStatus        `bson:"status" json:"status"`

    // 时间戳
    CreatedAt       time.Time            `bson:"created_at" json:"createdAt"`
    UpdatedAt       time.Time            `bson:"updated_at" json:"updatedAt"`
    DeletedAt       *time.Time           `bson:"deleted_at,omitempty" json:"deletedAt,omitempty"`
}

type TargetType string
const (
    TargetTypeBook      TargetType = "book"
    TargetTypeChapter   TargetType = "chapter"
    TargetTypeReview    TargetType = "review"
)

type CommentStatus string
const (
    CommentStatusNormal   CommentStatus = "normal"
    CommentStatusHidden   CommentStatus = "hidden"
    CommentStatusDeleted  CommentStatus = "deleted"
)
```

### Like (点赞)

```go
type Like struct {
    ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    UserID          primitive.ObjectID   `bson:"user_id" json:"userId"`
    TargetType      TargetType          `bson:"target_type" json:"targetType"`
    TargetID        primitive.ObjectID   `bson:"target_id" json:"targetId"`
    CreatedAt       time.Time            `bson:"created_at" json:"createdAt"`
}
```

### Follow (关注)

```go
type Follow struct {
    ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    FollowerID      primitive.ObjectID   `bson:"follower_id" json:"followerId"`
    FollowingID     primitive.ObjectID   `bson:"following_id" json:"followingId"`
    CreatedAt       time.Time            `bson:"created_at" json:"createdAt"`
}
```

### Collection (收藏)

```go
type Collection struct {
    ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    UserID          primitive.ObjectID   `bson:"user_id" json:"userId"`
    TargetType      TargetType          `bson:"target_type" json:"targetType"`
    TargetID        primitive.ObjectID   `bson:"target_id" json:"targetId"`
    FolderID        *primitive.ObjectID  `bson:"folder_id,omitempty" json:"folderId,omitempty"`
    Tags            []string             `bson:"tags" json:"tags"`
    Note            string               `bson:"note" json:"note"`
    CreatedAt       time.Time            `bson:"created_at" json:"createdAt"`
}
```

## 🌐 API端点

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | /api/v1/social/comments | 发布评论 | 是 |
| GET | /api/v1/social/comments/:targetType/:targetId | 获取评论列表 | 否 |
| PUT | /api/v1/social/comments/:id | 更新评论 | 是 |
| DELETE | /api/v1/social/comments/:id | 删除评论 | 是 |
| POST | /api/v1/social/comments/:id/like | 点赞评论 | 是 |
| POST | /api/v1/social/likes | 点赞内容 | 是 |
| DELETE | /api/v1/social/likes/:targetType/:targetId | 取消点赞 | 是 |
| GET | /api/v1/social/likes/status | 查询点赞状态 | 是 |
| POST | /api/v1/social/follow/:userId | 关注用户 | 是 |
| DELETE | /api/v1/social/follow/:userId | 取消关注 | 是 |
| GET | /api/v1/social/following/:userId | 获取关注列表 | 否 |
| GET | /api/v1/social/followers/:userId | 获取粉丝列表 | 否 |
| POST | /api/v1/social/collections | 收藏内容 | 是 |
| DELETE | /api/v1/social/collections/:targetType/:targetId | 取消收藏 | 是 |
| GET | /api/v1/social/collections | 获取收藏列表 | 是 |

## 🔧 依赖关系

### 依赖的模块
- **01 - 认证授权**: 用户身份验证
- **02 - 写作创作**: 获取作品信息
- **03 - 阅读器**: 获取阅读内容

### 被依赖的模块
- **06 - 书城**: 展示评论、点赞数
- **07 - 管理**: 内容审核

## 📈 扩展点

1. **话题系统**
   - 创建话题
   - 话题讨论
   - 热门话题

2. **动态分享**
   - 发布动态
   - 动态互动
   - 动态推荐

3. **私信功能**
   - 一对一私聊
   - 群组聊天
   - 消息提醒

---

**文档维护**: 青羽后端架构团队
**最后更新**: 2025-01-06
**对应实现**: `../../Qingyu_backend/api/v1/social/`
