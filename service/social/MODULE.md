# Social Service

> 最后更新：2026-03-29

## 职责

社交互动层，管理关注/粉丝、点赞、收藏、评论（含敏感词审核）、评分、书评、私信消息。不管理用户认证和内容创作。

## 数据流

```
API Handler → {LikeService/CommentService/FollowService/...} → Repository → MongoDB
                ↓
         EventBus（异步事件：like.book、comment.published、follow.user 等）
```

## 约定 & 陷阱

- **评论敏感词审核**：`checkSensitiveWords` 在发布前自动调用，命中敏感词的评论会被标记为 pending 状态
- **评论树结构**：评论支持嵌套回复，通过 `rootID` 关联线程，`getRootID` 自动处理层级
- **批量点赞**：`BatchLikeBooks` 返回 `BatchLikeBooksResult` 包含成功/失败分项，部分失败不回滚
- **评分统计**：`RatingService` 独立实现，评分变更后需要触发统计更新
- **事务包裹**：关注、收藏等写操作使用 MongoDB 事务，确保计数器与关系记录一致性
