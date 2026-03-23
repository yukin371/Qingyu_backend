# Social Repository (MongoDB)

社交模块的 MongoDB Repository 实现，提供评论、点赞、关注、收藏、私信、书评、书单等数据的持久化操作。

## 模块职责

负责社交相关数据的 MongoDB 存储和查询，包括数据 CRUD、索引管理、事务支持、批量操作等。

## 架构图

```mermaid
graph TB
    subgraph "Service Layer"
        CS[CommentService]
        LS[LikeService]
        FS[FollowService]
        ColS[CollectionService]
        MS[MessageService]
        RS[ReviewService]
        BLS[BookListService]
        URS[UserRelationService]
    end

    subgraph "Repository Interface"
        CR_I[CommentRepository<br/>接口]
        LR_I[LikeRepository<br/>接口]
        FR_I[FollowRepository<br/>接口]
        ColR_I[CollectionRepository<br/>接口]
        MR_I[MessageRepository<br/>接口]
        RR_I[ReviewRepository<br/>接口]
        BLR_I[BookListRepository<br/>接口]
        URR_I[UserRelationRepository<br/>接口]
    end

    subgraph "MongoDB Implementation"
        CR[MongoCommentRepository]
        LR[MongoLikeRepository]
        FR[MongoFollowRepository]
        ColR[MongoCollectionRepository]
        MR[MongoMessageRepository]
        RR[MongoReviewRepository]
        BLR[MongoBookListRepository]
        URR[MongoUserRelationRepository]
    end

    subgraph "Base"
        Base[BaseMongoRepository<br/>通用CRUD/ID解析]
    end

    CS --> CR_I
    LS --> LR_I
    FS --> FR_I
    ColS --> ColR_I
    MS --> MR_I
    RS --> RR_I
    BLS --> BLR_I
    URS --> URR_I

    CR_I --> CR
    LR_I --> LR
    FR_I --> FR
    ColR_I --> ColR
    MR_I --> MR
    RR_I --> RR
    BLR_I --> BLR
    URR_I --> URR

    CR --> Base
    LR --> Base
    FR --> Base
    ColR --> Base
    MR --> Base
    RR --> Base
    BLR --> Base
    URR --> Base
```

## Repository 列表

### 1. MongoCommentRepository

评论仓储，管理评论数据的持久化。

**数据集合**: `comments`

**索引设计**:
| 索引字段 | 用途 |
|----------|------|
| `{book_id: 1, created_at: -1}` | 书籍评论列表查询 |
| `{user_id: 1, created_at: -1}` | 用户评论列表查询 |
| `{parent_id: 1}` | 回复列表查询 |
| `{status: 1, created_at: -1}` | 待审核评论查询 |
| `{book_id: 1, like_count: -1}` | 热门评论排序 |
| `{chapter_id: 1, created_at: -1}` | 章节评论查询 |

**核心方法**:
| 方法 | 职责 |
|------|------|
| `Create` | 创建评论 |
| `GetByID` | 根据ID获取评论 |
| `Update` | 更新评论 |
| `Delete` | 软删除评论（标记为deleted状态） |
| `Exists` | 检查评论是否存在 |
| `GetCommentsByBookID` | 获取书籍评论列表 |
| `GetCommentsByBookIDSorted` | 获取排序的书籍评论 |
| `GetCommentsByUserID` | 获取用户评论列表 |
| `GetRepliesByCommentID` | 获取评论的回复列表 |
| `GetCommentsByChapterID` | 获取章节评论列表 |
| `IncrementLikeCount` | 增加点赞数 |
| `DecrementLikeCount` | 减少点赞数 |
| `IncrementReplyCount` | 增加回复数 |
| `DecrementReplyCount` | 减少回复数 |
| `GetBookRatingStats` | 获取书籍评分统计（聚合） |
| `GetCommentCount` | 获取评论总数 |
| `GetCommentsByIDs` | 批量获取评论 |
| `DeleteCommentsByBookID` | 删除书籍所有评论 |
| `RunInTransaction` | 事务执行 |

### 2. MongoLikeRepository

点赞仓储，管理点赞数据。

**数据集合**: `likes`

**索引设计**:
| 索引字段 | 用途 |
|----------|------|
| `{user_id: 1, target_type: 1, target_id: 1}` (唯一) | 防止重复点赞 |
| `{user_id: 1, created_at: -1}` | 用户点赞列表查询 |
| `{target_type: 1, target_id: 1}` | 目标点赞数统计 |

**核心方法**:
| 方法 | 职责 |
|------|------|
| `AddLike` | 添加点赞（唯一索引防重） |
| `RemoveLike` | 取消点赞 |
| `IsLiked` | 检查是否已点赞 |
| `GetByID` | 根据ID获取点赞记录 |
| `GetUserLikes` | 获取用户点赞列表 |
| `GetLikeCount` | 获取目标点赞数 |
| `GetLikesCountBatch` | 批量获取点赞数（聚合） |
| `GetUserLikeStatusBatch` | 批量检查点赞状态 |
| `CountUserLikes` | 统计用户点赞总数 |
| `RunInTransaction` | 事务执行 |

### 3. MongoFollowRepository

关注仓储，管理用户关注关系。

**数据集合**: `follows`, `author_follows`

**核心方法**:
| 方法 | 职责 |
|------|------|
| `CreateFollow` | 创建关注关系 |
| `DeleteFollow` | 删除关注关系 |
| `IsFollowing` | 检查是否关注 |
| `GetFollowers` | 获取粉丝列表 |
| `GetFollowing` | 获取关注列表 |
| `UpdateMutualStatus` | 更新互关状态 |
| `UpdateFollowStats` | 更新关注统计 |
| `GetFollowStats` | 获取关注统计 |
| `CreateAuthorFollow` | 关注作者 |
| `DeleteAuthorFollow` | 取消关注作者 |
| `GetAuthorFollow` | 获取作者关注记录 |
| `GetUserFollowingAuthors` | 获取用户关注的作者 |
| `RunInTransaction` | 事务执行 |

### 4. MongoCollectionRepository

收藏仓储，管理用户收藏和收藏夹。

**数据集合**: `collections`, `collection_folders`

**核心方法**:
| 方法 | 职责 |
|------|------|
| `Create` | 创建收藏 |
| `GetByID` | 获取收藏详情 |
| `GetByUserAndBook` | 获取用户对书籍的收藏 |
| `GetCollectionsByUser` | 获取用户收藏列表 |
| `GetCollectionsByTag` | 按标签获取收藏 |
| `Update` | 更新收藏 |
| `Delete` | 删除收藏 |
| `CreateFolder` | 创建收藏夹 |
| `GetFolderByID` | 获取收藏夹 |
| `GetFoldersByUser` | 获取用户收藏夹列表 |
| `UpdateFolder` | 更新收藏夹 |
| `DeleteFolder` | 删除收藏夹 |
| `IncrementFolderBookCount` | 增加收藏夹书籍数 |
| `DecrementFolderBookCount` | 减少收藏夹书籍数 |
| `GetPublicCollections` | 获取公开收藏 |
| `GetByShareID` | 根据分享ID获取收藏 |
| `CountUserCollections` | 统计用户收藏数 |
| `RunInTransaction` | 事务执行 |

### 5. MongoMessageRepository

私信仓储，管理会话和消息。

**数据集合**: `conversations`, `messages`, `mentions`

**核心方法**:
| 方法 | 职责 |
|------|------|
| `CreateConversation` | 创建会话 |
| `GetConversationByID` | 获取会话 |
| `GetConversationByParticipants` | 根据参与者获取会话 |
| `GetUserConversations` | 获取用户会话列表 |
| `UpdateLastMessage` | 更新最后一条消息 |
| `IncrementUnreadCount` | 增加未读数 |
| `ClearUnreadCount` | 清空未读数 |
| `CreateMessage` | 创建消息 |
| `GetMessageByID` | 获取消息 |
| `GetMessagesByConversation` | 获取会话消息列表 |
| `MarkMessageAsRead` | 标记消息已读 |
| `DeleteMessage` | 删除消息 |
| `CreateMention` | 创建@提醒 |
| `GetMentionByID` | 获取@提醒 |
| `GetUserMentions` | 获取用户@提醒列表 |
| `MarkMentionAsRead` | 标记@提醒已读 |
| `RunInTransaction` | 事务执行 |

### 6. MongoReviewRepository

书评仓储，管理书评数据。

**数据集合**: `reviews`, `review_likes`

**核心方法**:
| 方法 | 职责 |
|------|------|
| `CreateReview` | 创建书评 |
| `GetReviewByID` | 获取书评 |
| `GetReviewsByBook` | 获取书籍书评列表 |
| `UpdateReview` | 更新书评 |
| `DeleteReview` | 删除书评 |
| `CreateReviewLike` | 创建书评点赞 |
| `DeleteReviewLike` | 删除书评点赞 |
| `IsReviewLiked` | 检查书评是否已点赞 |
| `IncrementReviewLikeCount` | 增加书评点赞数 |
| `DecrementReviewLikeCount` | 减少书评点赞数 |
| `RunInTransaction` | 事务执行 |

### 7. MongoBookListRepository

书单仓储，管理用户书单。

**数据集合**: `booklists`, `booklist_likes`

**核心方法**:
| 方法 | 职责 |
|------|------|
| `CreateBookList` | 创建书单 |
| `GetBookListByID` | 获取书单 |
| `GetPublicBookLists` | 获取公开书单列表 |
| `UpdateBookList` | 更新书单 |
| `DeleteBookList` | 删除书单 |
| `CreateBookListLike` | 创建书单点赞 |
| `IsBookListLiked` | 检查书单是否已点赞 |
| `IncrementBookListLikeCount` | 增加书单点赞数 |
| `IncrementViewCount` | 增加浏览数 |
| `ForkBookList` | 复制书单 |
| `IncrementForkCount` | 增加复制数 |
| `GetBooksInList` | 获取书单中的书籍 |
| `RunInTransaction` | 事务执行 |

### 8. MongoUserRelationRepository

用户关系仓储，管理用户间的关系状态。

**数据集合**: `user_relations`

**核心方法**:
| 方法 | 职责 |
|------|------|
| `Create` | 创建关系记录 |
| `GetByID` | 获取关系记录 |
| `GetRelation` | 获取两个用户的关系 |
| `Update` | 更新关系 |
| `GetFollowers` | 获取粉丝列表 |
| `GetFollowing` | 获取关注列表 |
| `IsFollowing` | 检查是否关注 |
| `CountFollowers` | 统计粉丝数 |
| `CountFollowing` | 统计关注数 |

## 通用特性

### 1. BaseMongoRepository 基类

所有 Repository 都继承自 `BaseMongoRepository`，提供通用功能：

- `ParseID(id string)` - ID 解析（支持 ObjectID）
- `GetCollection()` - 获取集合
- `GetDB()` - 获取数据库实例
- 通用 CRUD 方法

### 2. 事务支持

关键 Repository 实现了 `RunInTransaction` 方法：

```go
func (r *MongoCommentRepository) RunInTransaction(ctx context.Context, fn func(context.Context) error) error {
    session, err := r.GetDB().Client().StartSession()
    if err != nil {
        return fmt.Errorf("failed to start session: %w", err)
    }
    defer session.EndSession(ctx)

    _, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
        return nil, fn(sessCtx)
    })
    return err
}
```

### 3. 软删除

评论使用软删除策略，标记 `state` 字段为 `deleted`：

```go
func (r *MongoCommentRepository) Delete(ctx context.Context, id string) error {
    result, err := r.GetCollection().UpdateOne(
        ctx,
        bson.M{"_id": objectID},
        bson.M{
            "$set": bson.M{
                "state":      social.CommentStateDeleted,
                "updated_at": time.Now(),
            },
        },
    )
    // ...
}
```

### 4. 批量操作

支持批量查询和聚合操作：

```go
// 批量获取点赞数
func (r *MongoLikeRepository) GetLikesCountBatch(ctx context.Context, targetType string, targetIDs []string) (map[string]int64, error) {
    pipeline := mongo.Pipeline{
        {{Key: "$match", Value: bson.M{
            "target_type": targetType,
            "target_id":   bson.M{"$in": targetIDs},
        }}},
        {{Key: "$group", Value: bson.M{
            "_id":   "$target_id",
            "count": bson.M{"$sum": 1},
        }}},
    }
    // ...
}
```

### 5. 参数安全验证

对查询参数进行安全验证：

```go
func sanitizeSocialCommentQueryToken(field, value string) (string, error) {
    value = strings.TrimSpace(value)
    if value == "" {
        return "", fmt.Errorf("%s不能为空", field)
    }
    objectID, err := primitive.ObjectIDFromHex(value)
    if err != nil {
        return value, nil // 非ObjectID格式也接受
    }
    return objectID.Hex(), nil
}
```

## 索引初始化

Repository 在创建时会自动初始化必要的索引：

```go
func NewMongoCommentRepository(db *mongo.Database) *MongoCommentRepository {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    indexes := []mongo.IndexModel{
        {Keys: bson.D{{Key: "book_id", Value: 1}, {Key: "created_at", Value: -1}}},
        // ...更多索引
    }

    collection := db.Collection("comments")
    _, err := collection.Indexes().CreateMany(ctx, indexes)
    if err != nil {
        fmt.Printf("Warning: Failed to create indexes: %v\n", err)
    }
    // ...
}
```

## 文件列表

| 文件 | 职责 |
|------|------|
| `comment_repository_mongo.go` | 评论仓储实现 |
| `like_repository_mongo.go` | 点赞仓储实现 |
| `follow_repository_mongo.go` | 关注仓储实现 |
| `collection_repository_mongo.go` | 收藏仓储实现 |
| `message_repository_mongo.go` | 私信仓储实现 |
| `review_repository_mongo.go` | 书评仓储实现 |
| `booklist_repository_mongo.go` | 书单仓储实现 |
| `user_relation_repository_impl.go` | 用户关系仓储实现 |
