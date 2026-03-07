# 完整的发布-审核-订阅系统架构设计

**设计日期**: 2026-03-06
**设计者**: Kore
**优先级**: 🔴 P0
**需求**: Project发布→审核→订阅→版本管理的完整流程

---

## 系统流程图

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           完整业务流程                                           │
└─────────────────────────────────────────────────────────────────────────────────┘

Author                    Publish System              Review System         Reader
   │                           │                           │                   │
   │  1. 创建Publication         │                           │                   │
   │  (单章/全书)                │                           │                   │
   │──────────────────────────>│                           │                   │
   │                           │                           │                   │
   │                           │  2. 创建Review             │                   │
   │                           │─────────────────────────>│                   │
   │                           │                           │                   │
   │  3. 查看审核状态            │                           │  4. 审核            │
   │  <─────────────────────────│──────────────────────────│<──────────────────│
   │                           │                           │                   │
   │                           │  5. 审核通过，创建Book      │                   │
   │                           │  (或BookVersion)          │                   │
   │                           │                           │                   │
   │                           │  6. 书籍上架               │                   │
   │                           │──────────────────────────────────────────────>│
   │                           │                           │                   │
   │                           │                           │  7. 订阅           │
   │                           │                           │<──────────────────│
   │                           │                           │                   │
   │  8. 继续写作，更新Publication│                           │                   │
   │  (提交新版本)               │                           │                   │
   │──────────────────────────>│                           │                   │
   │                           │                           │                   │
   │                           │  9. 新版本审核              │                   │
   │                           │─────────────────────────>│                   │
   │                           │                           │                   │
   │                           │  10. 创建BookVersion(新)   │                   │
   │                           │                           │                   │
   │                           │  11. 通知已订阅用户         │                   │
   │                           │──────────────────────────────────────────────>│
   │                           │                           │                   │
   │                           │                           │  12. 查看版本差异    │
   │                           │                           │<──────────────────│
```

---

## 核心数据模型

### 1. Publication（发布单）- 作者发布内容

发布单是作者提交发布请求的核心实体，记录发布类型、范围和内容。

```go
// models/publishing/publication.go

package publishing

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "Qingyu_backend/models/shared"
)

// PublicationStatus 发布单状态
type PublicationStatus string

const (
    PublicationStatusDraft      PublicationStatus = "draft"       // 草稿（未提交）
    PublicationStatusSubmitted  PublicationStatus = "submitted"   // 已提交审核
    PublicationStatusReviewing  PublicationStatus = "reviewing"   // 审核中
    PublicationStatusApproved   PublicationStatus = "approved"    // 审核通过
    PublicationStatusRejected   PublicationStatus = "rejected"    // 审核拒绝
    PublicationStatusPublished  PublicationStatus = "published"   // 已发布
    PublicationStatusCancelled  PublicationStatus = "cancelled"   // 已取消
)

// PublicationType 发布类型
type PublicationType string

const (
    PublicationTypeFullBook    PublicationType = "full_book"    // 全书发布
    PublicationTypeSingleChapter PublicationType = "single_chapter" // 单章发布
    PublicationTypeBatchChapter PublicationType = "batch_chapter"  // 批量章节
    PublicationTypeUpdate      PublicationType = "update"        // 内容更新
)

// Publication 发布单
type Publication struct {
    shared.IdentifiedEntity `bson:",inline"`
    shared.BaseEntity       `bson:",inline"`

    // 基本信息
    AuthorID       string           `bson:"author_id" json:"authorId"`
    ProjectID      string           `bson:"project_id" json:"projectId"`           // ← 关联的写作项目
    BookID         *string          `bson:"book_id,omitempty" json:"bookId,omitempty"` // ← 关联的书籍（更新时）

    // 发布类型和范围
    Type           PublicationType  `bson:"type" json:"type"`                         // 发布类型
    Status         PublicationStatus `bson:"status" json:"status"`                   // 发布单状态

    // 章节范围（全书发布时为空）
    ChapterIDs     []string         `bson:"chapter_ids,omitempty" json:"chapterIds,omitempty"` // 包含的章节ID

    // 快照内容（发布时的内容快照）
    Snapshot       *PublicationSnapshot `bson:"snapshot,omitempty" json:"snapshot,omitempty"`

    // 审核信息
    ReviewID       *string          `bson:"review_id,omitempty" json:"reviewId,omitempty"`     // 关联的审核单
    SubmitTime     *time.Time       `bson:"submit_time,omitempty" json:"submitTime,omitempty"` // 提交时间
    ApprovedTime   *time.Time       `bson:"approved_time,omitempty" json:"approvedTime,omitempty"` // 审核通过时间

    // 发布信息
    PublishedTime  *time.Time       `bson:"published_time,omitempty" json:"publishedTime,omitempty"` // 发布时间
    BookVersionID  *string          `bson:"book_version_id,omitempty" json:"bookVersionId,omitempty"` // 产生的版本ID

    // 拒绝原因
    RejectionReason *string         `bson:"rejection_reason,omitempty" json:"rejectionReason,omitempty"` // 拒绝原因

    // 元数据
    Title          string           `bson:"title" json:"title"`                           // 发布单标题
    Description    string           `bson:"description,omitempty" json:"description,omitempty"` // 描述
    Metadata       map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"` // 扩展元数据
}

// PublicationSnapshot 发布快照
type PublicationSnapshot struct {
    // Book基本信息
    BookTitle      string           `bson:"book_title" json:"bookTitle"`
    BookCover      string           `bson:"book_cover" json:"bookCover"`
    BookSummary    string           `bson:"book_summary" json:"bookSummary"`
    BookTags       []string         `bson:"book_tags" json:"bookTags"`
    CategoryIDs    []primitive.ObjectID `bson:"category_ids" json:"categoryIds"`

    // 定价信息
    Price          float64          `bson:"price" json:"price"`
    IsFree         bool             `bson:"is_free" json:"isFree"`

    // 章节信息快照
    Chapters       []ChapterSnapshot `bson:"chapters" json:"chapters"`

    // 内容哈希（用于变更检测）
    ContentHash    string           `bson:"content_hash" json:"contentHash"`

    // 快照时间
    SnapshotTime   time.Time        `bson:"snapshot_time" json:"snapshotTime"`
}

// ChapterSnapshot 章节快照
type ChapterSnapshot struct {
    ChapterID      string           `bson:"chapter_id" json:"chapterId"`
    ProjectChapterID string         `bson:"project_chapter_id" json:"projectChapterId"`
    ChapterNumber  int              `bson:"chapter_number" json:"chapterNumber"`
    Title          string           `bson:"title" json:"title"`
    Content        string           `bson:"content" json:"content"`
    WordCount      int              `bson:"word_count" json:"wordCount"`
    Price          float64          `bson:"price" json:"price"`
    IsFree         bool             `bson:"is_free" json:"isFree"`
    ContentHash    string           `bson:"content_hash" json:"contentHash"`
}
```

### 2. Review（审核单）- 管理员审核

```go
// models/publishing/review.go

package publishing

// ReviewStatus 审核状态
type ReviewStatus string

const (
    ReviewStatusPending   ReviewStatus = "pending"    // 待审核
    ReviewStatusReviewing ReviewStatus = "reviewing"  // 审核中
    ReviewStatusApproved  ReviewStatus = "approved"   // 已通过
    ReviewStatusRejected  ReviewStatus = "rejected"   // 已拒绝
    ReviewStatusCancelled ReviewStatus = "cancelled"  // 已取消
)

// ReviewPriority 审核优先级
type ReviewPriority string

const (
    ReviewPriorityLow    ReviewPriority = "low"       // 低优先级
    ReviewPriorityNormal ReviewPriority = "normal"    // 普通优先级
    ReviewPriorityHigh   ReviewPriority = "high"      // 高优先级
    ReviewPriorityUrgent ReviewPriority = "urgent"    // 紧急
)

// Review 审核单
type Review struct {
    shared.IdentifiedEntity `bson:",inline"`
    shared.BaseEntity       `bson:",inline"`

    // 关联信息
    PublicationID  string           `bson:"publication_id" json:"publicationId"`   // ← 关联发布单
    BookID         *string          `bson:"book_id,omitempty" json:"bookId,omitempty"` // 关联书籍（审核更新时）

    // 审核状态
    Status         ReviewStatus     `bson:"status" json:"status"`
    Priority       ReviewPriority   `bson:"priority" json:"priority"`

    // 审核人信息
    ReviewerID     *string          `bson:"reviewer_id,omitempty" json:"reviewerId,omitempty"` // 审核人ID
    ReviewerName   *string          `bson:"reviewer_name,omitempty" json:"reviewerName,omitempty"` // 审核人名称
    AssignedTime   *time.Time       `bson:"assigned_time,omitempty" json:"assignedTime,omitempty"` // 分配时间

    // 审核时间
    StartTime      *time.Time       `bson:"start_time,omitempty" json:"startTime,omitempty"` // 开始审核时间
    CompletedTime  *time.Time       `bson:"completed_time,omitempty" json:"completedTime,omitempty"` // 完成时间

    // 审核结果
    Result         *ReviewResult    `bson:"result,omitempty" json:"result,omitempty"` // 审核结果

    // 审核历史（记录流转过程）
    History        []ReviewHistory  `bson:"history,omitempty" json:"history,omitempty"` // 审核历史

    // 元数据
    Type           PublicationType  `bson:"type" json:"type"` // 发布类型
    Title          string           `bson:"title" json:"title"` // 审核标题
    Tags           []string         `bson:"tags,omitempty" json:"tags,omitempty"` // 标签
}

// ReviewResult 审核结果
type ReviewResult struct {
    Decision       ReviewDecision   `bson:"decision" json:"decision"` // 审核决定
    Reason         string           `bson:"reason,omitempty" json:"reason,omitempty"` // 原因说明
    Comments       []string         `bson:"comments,omitempty" json:"comments,omitempty"` // 详细意见
    Issues         []ReviewIssue    `bson:"issues,omitempty" json:"issues,omitempty"` // 问题列表
    ApprovedVersion *string         `bson:"approved_version,omitempty" json:"approvedVersion,omitempty"` // 批准的版本
}

// ReviewDecision 审核决定
type ReviewDecision string

const (
    ReviewDecisionApprove ReviewDecision = "approve" // 通过
    ReviewDecisionReject  ReviewDecision = "reject"  // 拒绝
    ReviewDecisionRequest ReviewDecision = "request" // 需要修改
)

// ReviewIssue 审核问题
type ReviewIssue struct {
    Type           string           `bson:"type" json:"type"` // 问题类型：content, formatting, legal, other
    Severity       string           `bson:"severity" json:"severity"` // 严重程度：error, warning, info
    Location       string           `bson:"location,omitempty" json:"location,omitempty"` // 位置（章节ID或段落）
    Description    string           `bson:"description" json:"description"` // 问题描述
    Suggestion     string           `bson:"suggestion,omitempty" json:"suggestion,omitempty"` // 修改建议
}

// ReviewHistory 审核历史
type ReviewHistory struct {
    Action         string           `bson:"action" json:"action"` // 动作：submitted, assigned, approved, rejected
    ActorID        string           `bson:"actor_id" json:"actorId"` // 操作人ID
    ActorName      string           `bson:"actor_name" json:"actorName"` // 操作人名称
    ActorRole      string           `bson:"actor_role" json:"actorRole"` // 操作人角色
    Timestamp      time.Time        `bson:"timestamp" json:"timestamp"` // 操作时间
    Comment        string           `bson:"comment,omitempty" json:"comment,omitempty"` // 备注
}
```

### 3. BookVersion（书籍版本）- 版本管理

```go
// models/bookstore/book_version.go

package bookstore

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "Qingyu_backend/models/shared"
)

// BookVersion 书籍版本
type BookVersion struct {
    shared.IdentifiedEntity `bson:",inline"`
    shared.BaseEntity       `bson:",inline"`

    // 关联信息
    BookID         primitive.ObjectID `bson:"book_id" json:"bookId"`                 // ← 关联书籍
    PublicationID  string             `bson:"publication_id" json:"publicationId"` // ← 来源发布单

    // 版本信息
    VersionNumber  int                `bson:"version_number" json:"versionNumber"` // 版本号
    VersionName    string             `bson:"version_name" json:"versionName"`     // 版本名称
    VersionType    VersionType        `bson:"version_type" json:"versionType"`     // 版本类型

    // 快照内容
    Content        *VersionContent    `bson:"content" json:"content"`               // 版本内容快照

    // 变更统计
    ChangeSummary  *ChangeSummary     `bson:"change_summary,omitempty" json:"changeSummary,omitempty"` // 变更摘要

    // 发布信息
    PublishedAt    time.Time          `bson:"published_at" json:"publishedAt"`     // 发布时间
    PublishedBy    string             `bson:"published_by" json:"publishedBy"`     // 发布人

    // 版本状态
    Status         VersionStatus      `bson:"status" json:"status"`                 // 版本状态
    IsLatest       bool               `bson:"is_latest" json:"isLatest"`           // 是否最新版本

    // 统计信息
    ViewCount      int64              `bson:"view_count" json:"viewCount"`         // 查看次数
    SubscriptionCount int64           `bson:"subscription_count" json:"subscriptionCount"` // 订阅该版本的用户数

    // 元数据
    Description    string             `bson:"description,omitempty" json:"description,omitempty"` // 版本说明
    Changelog      string             `bson:"changelog,omitempty" json:"changelog,omitempty"` // 更新日志
    Tags           []string           `bson:"tags,omitempty" json:"tags,omitempty"` // 标签
}

// VersionType 版本类型
type VersionType string

const (
    VersionTypeInitial    VersionType = "initial"     // 初始版本
    VersionTypeMajor      VersionType = "major"       // 主要版本（大改）
    VersionTypeMinor      VersionType = "minor"       // 次要版本（新增章节）
    VersionTypePatch      VersionType = "patch"       // 补丁版本（修正错别字）
    VersionTypeEmergency  VersionType = "emergency"   // 紧急版本
)

// VersionStatus 版本状态
type VersionStatus string

const (
    VersionStatusActive     VersionStatus = "active"       // 活跃版本
    VersionStatusArchived   VersionStatus = "archived"     // 已归档
    VersionStatusDeprecated VersionStatus = "deprecated"   // 已废弃
    VersionStatusDraft      VersionStatus = "draft"        // 草稿
)

// VersionContent 版本内容快照
type VersionContent struct {
    // Book信息
    Title          string             `bson:"title" json:"title"`
    Cover          string             `bson:"cover" json:"cover"`
    Summary        string             `bson:"summary" json:"summary"`
    Tags           []string           `bson:"tags" json:"tags"`
    CategoryIDs    []primitive.ObjectID `bson:"category_ids" json:"categoryIds"`

    // 章节内容快照
    Chapters       []VersionChapter   `bson:"chapters" json:"chapters"`

    // 统计信息
    TotalWords     int64              `bson:"total_words" json:"totalWords"`
    TotalChapters  int                `bson:"total_chapters" json:"totalChapters"`

    // 内容哈希
    ContentHash    string             `bson:"content_hash" json:"contentHash"`

    // 快照时间
    SnapshotTime   time.Time          `bson:"snapshot_time" json:"snapshotTime"`
}

// VersionChapter 章节版本
type VersionChapter struct {
    ChapterID          string         `bson:"chapter_id" json:"chapterId"`
    ProjectChapterID   string         `bson:"project_chapter_id" json:"projectChapterId"`
    ChapterNumber      int            `bson:"chapter_number" json:"chapterNumber"`
    Title              string         `bson:"title" json:"title"`
    Content            string         `bson:"content" json:"content"`
    WordCount          int            `bson:"word_count" json:"wordCount"`
    Price              float64        `bson:"price" json:"price"`
    IsFree             bool           `bson:"is_free" json:"isFree"`
    ContentHash        string         `bson:"content_hash" json:"contentHash"`

    // 变更标记
    Action             ChapterAction  `bson:"action" json:"action"` // 章节动作
    PreviousVersionID  *string        `bson:"previous_version_id,omitempty" json:"previousVersionId,omitempty"` // 上一版本
}

// ChapterAction 章节动作
type ChapterAction string

const (
    ChapterActionAdded    ChapterAction = "added"     // 新增章节
    ChapterActionModified ChapterAction = "modified"  // 修改章节
    ChapterActionDeleted  ChapterAction = "deleted"   // 删除章节
    ChapterActionUnchanged ChapterAction = "unchanged" // 无变化
)

// ChangeSummary 变更摘要
type ChangeSummary struct {
    AddedChapters     int              `bson:"added_chapters" json:"addedChapters"`       // 新增章节数
    ModifiedChapters  int              `bson:"modified_chapters" json:"modifiedChapters"`  // 修改章节数
    DeletedChapters   int              `bson:"deleted_chapters" json:"deletedChapters"`    // 删除章节数
    TotalWordsDiff    int64            `bson:"total_words_diff" json:"totalWordsDiff"`     // 字数变化
    ChangedChapterIDs []string         `bson:"changed_chapter_ids" json:"changedChapterIds"` // 变更的章节ID列表
}
```

### 4. Subscription（订阅）- 读者订阅

```go
// models/reader/subscription.go

package reader

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "Qingyu_backend/models/shared"
)

// SubscriptionStatus 订阅状态
type SubscriptionStatus string

const (
    SubscriptionStatusActive    SubscriptionStatus = "active"      // 活跃
    SubscriptionStatusPaused    SubscriptionStatus = "paused"      // 暂停
    SubscriptionStatusCancelled SubscriptionStatus = "cancelled"  // 已取消
    SubscriptionStatusExpired   SubscriptionStatus = "expired"     // 已过期
)

// SubscriptionType 订阅类型
type SubscriptionType string

const (
    SubscriptionTypeFullBook SubscriptionType = "full_book" // 全书订阅
    SubscriptionTypePerChapter SubscriptionType = "per_chapter" // 按章订阅
    SubscriptionTypeVIP SubscriptionType = "vip" // VIP订阅
)

// Subscription 订阅
type Subscription struct {
    shared.IdentifiedEntity `bson:",inline"`
    shared.BaseEntity       `bson:",inline"`

    // 订阅关系
    UserID         string             `bson:"user_id" json:"userId"`           // ← 用户ID
    BookID         primitive.ObjectID `bson:"book_id" json:"bookId"`          // ← 书籍ID
    ProjectID      *string            `bson:"project_id,omitempty" json:"projectId,omitempty"` // 项目ID（可选）

    // 订阅信息
    Type           SubscriptionType  `bson:"type" json:"type"`                 // 订阅类型
    Status         SubscriptionStatus `bson:"status" json:"status"`           // 订阅状态

    // 版本跟踪
    CurrentVersionID *string         `bson:"current_version_id,omitempty" json:"currentVersionId,omitempty"` // 当前阅读版本
    SubscribedVersionID *string      `bson:"subscribed_version_id,omitempty" json:"subscribedVersionId,omitempty"` // 订阅时的版本

    // 订阅范围
    ChapterIDs     []primitive.ObjectID `bson:"chapter_ids,omitempty" json:"chapterIds,omitempty"` // 已订阅章节

    // 时间信息
    SubscribedAt   time.Time          `bson:"subscribed_at" json:"subscribedAt"`   // 订阅时间
    ExpiresAt      *time.Time         `bson:"expires_at,omitempty" json:"expiresAt,omitempty"` // 过期时间

    // 费用信息
    TotalAmount    float64            `bson:"total_amount" json:"totalAmount"`     // 总金额
    PaidAmount     float64            `bson:"paid_amount" json:"paidAmount"`       // 已付金额

    // 自动更新设置
    AutoUpdate     bool               `bson:"auto_update" json:"autoUpdate"`       // 是否自动更新到新版本
    NotifyUpdate   bool               `bson:"notify_update" json:"notifyUpdate"`   // 是否通知更新

    // 元数据
    Metadata       map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
}

// ChapterSubscription 章节订阅记录
type ChapterSubscription struct {
    shared.IdentifiedEntity `bson:",inline"`

    SubscriptionID primitive.ObjectID `bson:"subscription_id" json:"subscriptionId"` // ← 关联合约
    UserID          string             `bson:"user_id" json:"userId"`               // 用户ID
    BookID          primitive.ObjectID `bson:"book_id" json:"bookId"`              // 书籍ID
    ChapterID       primitive.ObjectID `bson:"chapter_id" json:"chapterId"`        // 章节ID
    VersionID       string             `bson:"version_id" json:"versionId"`         // 购买的版本

    PurchasedAt     time.Time          `bson:"purchased_at" json:"purchasedAt"`     // 购买时间
    Price           float64            `bson:"price" json:"price"`                  // 价格
    IsCurrent       bool               `bson:"is_current" json:"isCurrent"`         // 是否当前版本
}
```

### 5. Book（书籍）- 关联字段更新

```go
// models/bookstore/book.go

type Book struct {
    shared.IdentifiedEntity `bson:",inline"`
    shared.BaseEntity       `bson:",inline"`

    Title         string `bson:"title" json:"title"`
    Author        string `bson:"author" json:"author"`
    AuthorID      string `bson:"author_id,omitempty" json:"authorId,omitempty"`
    Introduction  string `bson:"introduction" json:"introduction"`
    Cover         string `bson:"cover" json:"cover"`
    CategoryIDs   []primitive.ObjectID `bson:"category_ids" json:"categoryIds"`
    Tags          []string `bson:"tags" json:"tags"`
    Status        BookStatus `bson:"status" json:"status"`

    // ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
    //  新增：发布和版本关联字段
    // ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
    ProjectID         *string          `bson:"project_id,omitempty" json:"projectId,omitempty"` // ← 关联Project
    CurrentVersionID  *string          `bson:"current_version_id,omitempty" json:"currentVersionId,omitempty"` // ← 当前版本ID
    LatestVersionID   *string          `bson:"latest_version_id,omitempty" json:"latestVersionId,omitempty"` // ← 最新版本ID
    PublicationID     *string          `bson:"publication_id,omitempty" json:"publicationId,omitempty"` // ← 来源发布单
    TotalVersions     int              `bson:"total_versions" json:"totalVersions"` // 版本总数
    // ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

    // 阅读相关
    Rating        types.Rating `bson:"rating" json:"rating"`
    ViewCount     int64        `bson:"view_count" json:"viewCount"`
    WordCount     int64        `bson:"word_count" json:"wordCount"`
    ChapterCount  int          `bson:"chapter_count" json:"chapterCount"`
    Price         float64      `bson:"price" json:"price"`
    IsFree        bool         `bson:"is_free" json:"isFree"`

    // 统计
    SubscriptionCount int64      `bson:"subscription_count" json:"subscriptionCount"` // 订阅数
    VersionViewCount  int64      `bson:"version_view_count" json:"versionViewCount"` // 版本查看次数

    PublishedAt   *time.Time   `bson:"published_at,omitempty" json:"publishedAt,omitempty"`
    LastUpdateAt  *time.Time   `bson:"last_update_at,omitempty" json:"lastUpdateAt,omitempty"`
}
```

---

## 模型关系图

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              模型关系图                                          │
└─────────────────────────────────────────────────────────────────────────────────┘

┌──────────────┐         ┌──────────────┐         ┌──────────────┐
│   Project    │         │ Publication  │         │    Review    │
│  (写作项目)   │────────>│  (发布单)     │────────>│   (审核单)   │
│              │ 创建    │              │ 提交    │              │
│ - ID         │         │ - ID         │         │ - ID         │
│ - AuthorID   │         │ - ProjectID  │         │ - PublicationID
│ - Chapters[] │         │ - BookID     │         │ - ReviewerID │
│              │         │ - Snapshot   │         │ - Status     │
└──────────────┘         │ - Status     │         │ - Result     │
                          └──────┬───────┘         └──────────────┘
                                 │ 审核通过
                                 ▼
                          ┌──────────────┐         ┌──────────────┐
                          │     Book     │────────>│ BookVersion  │
                          │   (书籍)      │ 创建    │  (书籍版本)   │
                          ├──────────────┤         ├──────────────┤
                          │ - ID         │         │ - ID         │
                          │ - ProjectID  │◀───────│ - BookID     │
                          │ - CurrentVer │         │ - VersionNo  │
                          │ - LatestVer  │         │ - Content    │
                          └──────┬───────┘         │ - ChangeSum  │
                                 │                 └──────────────┘
                                 │ 订阅
                                 ▼
                          ┌──────────────┐
                          │ Subscription │
                          │   (订阅)      │
                          ├──────────────┤
                          │ - UserID     │
                          │ - BookID     │
                          │ - CurrentVer │
                          │ - AutoUpdate │
                          └──────────────┘
```

---

## 数据库索引

```go
// migration/mongodb/publication_indexes.go

func CreatePublicationIndexes(ctx context.Context, db *mongo.Database) error {
    collection := db.Collection("publications")

    indexes := []mongo.IndexModel{
        // Project ID索引（查询项目的发布单）
        {Keys: bson.D{{Key: "project_id", Value: 1}}, Options: options.Index().SetName("idx_project_id")},

        // Book ID索引（查询书籍的发布单）
        {Keys: bson.D{{Key: "book_id", Value: 1}}, Options: options.Index().SetName("idx_book_id")},

        // 作者状态索引（查询作者的待审核发布单）
        {
            Keys: bson.D{
                {Key: "author_id", Value: 1},
                {Key: "status", Value: 1},
            },
            Options: options.Index().SetName("idx_author_status"),
        },

        // 提交时间索引（按时间排序）
        {Keys: bson.D{{Key: "submit_time", Value: -1}}, Options: options.Index().SetName("idx_submit_time")},
    }

    _, err := collection.Indexes().CreateMany(ctx, indexes)
    return err
}

// migration/mongodb/review_indexes.go

func CreateReviewIndexes(ctx context.Context, db *mongo.Database) error {
    collection := db.Collection("reviews")

    indexes := []mongo.IndexModel{
        // Publication ID索引
        {Keys: bson.D{{Key: "publication_id", Value: 1}}, Options: options.Index().SetName("idx_publication_id")},

        // 审核人状态索引（分配给审核人的任务）
        {
            Keys: bson.D{
                {Key: "reviewer_id", Value: 1},
                {Key: "status", Value: 1},
            },
            Options: options.Index().SetName("idx_reviewer_status"),
        },

        // 优先级状态索引（按优先级排序）
        {
            Keys: bson.D{
                {Key: "priority", Value: -1},
                {Key: "status", Value: 1},
            },
            Options: options.Index().SetName("idx_priority_status"),
        },
    }

    _, err := collection.Indexes().CreateMany(ctx, indexes)
    return err
}

// migration/mongodb/book_version_indexes.go

func CreateBookVersionIndexes(ctx context.Context, db *mongo.Database) error {
    collection := db.Collection("book_versions")

    indexes := []mongo.IndexModel{
        // Book ID + 版本号索引（查询书籍的所有版本）
        {
            Keys: bson.D{
                {Key: "book_id", Value: 1},
                {Key: "version_number", Value: -1},
            },
            Options: options.Index().SetName("idx_book_version"),
        },

        // 最新版本索引
        {Keys: bson.D{{Key: "is_latest", Value: 1}}}, Options: options.Index().SetName("idx_is_latest")},

        // 发布时间索引
        {Keys: bson.D{{Key: "published_at", Value: -1}}}, Options: options.Index().SetName("idx_published_at")},
    }

    _, err := collection.Indexes().CreateMany(ctx, indexes)
    return err
}

// migration/mongodb/subscription_indexes.go

func CreateSubscriptionIndexes(ctx context.Context, db *mongo.Database) error {
    collection := db.Collection("subscriptions")

    indexes := []mongo.IndexModel{
        // 用户书籍索引（唯一）
        {
            Keys: bson.D{
                {Key: "user_id", Value: 1},
                {Key: "book_id", Value: 1},
            },
            Options: options.Index().SetUnique(true).SetName("idx_user_book_unique"),
        },

        // 用户状态索引
        {
            Keys: bson.D{
                {Key: "user_id", Value: 1},
                {Key: "status", Value: 1},
            },
            Options: options.Index().SetName("idx_user_status"),
        },

        // 过期时间索引
        {Keys: bson.D{{Key: "expires_at", Value: 1}}}, Options: options.Index().SetName("idx_expires_at")},
    }

    _, err := collection.Indexes().CreateMany(ctx, indexes)
    return err
}
```

---

## API 设计

### 发布相关 API

```go
// api/v1/writer/publish_api.go

// CreatePublication 创建发布单
// @Summary 创建发布单
// @Description 创建新的发布单（单章/全书/更新）
// @Tags Publishing
// @Router /writer/publications [post]
func (api *PublishAPI) CreatePublication(c *gin.Context) {
    var req CreatePublicationRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // 获取Project并创建快照
    publication, err := api.publishService.CreatePublication(
        c.Request.Context(),
        api.GetUserID(c),
        &req,
    )
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(201, publication)
}

// SubmitPublication 提交审核
// @Summary 提交发布单审核
// @Description 将草稿状态的发布单提交审核
// @Tags Publishing
// @Router /writer/publications/{id}/submit [post]
func (api *PublishAPI) SubmitPublication(c *gin.Context) {
    id := c.Param("id")
    userID := api.GetUserID(c)

    publication, err := api.publishService.SubmitPublication(
        c.Request.Context(),
        id,
        userID,
    )
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, publication)
}

// ListPublications 查询发布单列表
// @Summary 查询发布单列表
// @Description 查询当前用户的发布单列表
// @Tags Publishing
// @Router /writer/publications [get]
func (api *PublishAPI) ListPublications(c *gin.Context) {
    userID := api.GetUserID(c)
    status := c.Query("status")

    publications, err := api.publishService.ListPublications(
        c.Request.Context(),
        userID,
        status,
    )
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{
        "data": publications,
        "total": len(publications),
    })
}

// GetPublication 获取发布单详情
// @Summary 获取发布单详情
// @Tags Publishing
// @Router /writer/publications/{id} [get]
func (api *PublishAPI) GetPublication(c *gin.Context) {
    id := c.Param("id")
    userID := api.GetUserID(c)

    publication, err := api.publishService.GetPublication(c.Request.Context(), id, userID)
    if err != nil {
        c.JSON(404, gin.H{"error": "发布单不存在"})
        return
    }

    c.JSON(200, publication)
}

type CreatePublicationRequest struct {
    ProjectID    string           `json:"projectId" validate:"required"`
    Type         PublicationType  `json:"type" validate:"required"`
    ChapterIDs   []string         `json:"chapterIds,omitempty"` // 单章发布时指定
    Title        string           `json:"title" validate:"required"`
    Description  string           `json:"description,omitempty"`

    // Book信息（全书发布时使用）
    BookTitle    string           `json:"bookTitle,omitempty"`
    BookCover    string           `json:"bookCover,omitempty"`
    BookSummary  string           `json:"bookSummary,omitempty"`
    CategoryIDs  []string         `json:"categoryIds,omitempty"`
    Tags         []string         `json:"tags,omitempty"`
    Price        float64          `json:"price"`
    IsFree       bool             `json:"isFree"`
}
```

### 审核相关 API

```go
// api/v1/admin/review_api.go

// ListPendingReviews 获取待审核列表
// @Summary 获取待审核列表
// @Description 获取当前待审核的发布单列表
// @Tags Review
// @Router /admin/reviews/pending [get]
func (api *ReviewAPI) ListPendingReviews(c *gin.Context) {
    priority := c.Query("priority")
    limit := c.DefaultQuery("limit", "20")

    reviews, err := api.reviewService.ListPendingReviews(
        c.Request.Context(),
        priority,
        limit,
    )
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{
        "data": reviews,
        "total": len(reviews),
    })
}

// AssignReview 分配审核任务
// @Summary 分配审核任务
// @Description 将审核任务分配给审核人
// @Tags Review
// @Router /admin/reviews/{id}/assign [post]
func (api *ReviewAPI) AssignReview(c *gin.Context) {
    id := c.Param("id")
    reviewerID := api.GetUserID(c)

    review, err := api.reviewService.AssignReview(
        c.Request.Context(),
        id,
        reviewerID,
    )
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, review)
}

// ApproveReview 审核通过
// @Summary 审核通过
// @Description 审核通过发布单，创建Book或BookVersion
// @Tags Review
// @Router /admin/reviews/{id}/approve [post]
func (api *ReviewAPI) ApproveReview(c *gin.Context) {
    id := c.Param("id")
    reviewerID := api.GetUserID(c)

    var req ApproveRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    result, err := api.reviewService.ApproveReview(
        c.Request.Context(),
        id,
        reviewerID,
        &req,
    )
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, result)
}

// RejectReview 审核拒绝
// @Summary 审核拒绝
// @Description 拒绝发布单，记录拒绝原因
// @Tags Review
// @Router /admin/reviews/{id}/reject [post]
func (api *ReviewAPI) RejectReview(c *gin.Context) {
    id := c.Param("id")
    reviewerID := api.GetUserID(c)

    var req RejectRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    err := api.reviewService.RejectReview(
        c.Request.Context(),
        id,
        reviewerID,
        &req,
    )
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"message": "已拒绝"})
}

type ApproveRequest struct {
    Comments      string   `json:"comments,omitempty"`
    ApprovedVersion string  `json:"approvedVersion,omitempty"` // 批准的版本
}

type RejectRequest struct {
    Reason        string   `json:"reason" validate:"required"`
    Comments      string   `json:"comments,omitempty"`
    Issues        []ReviewIssue `json:"issues,omitempty"`
}
```

### 版本相关 API

```go
// api/v1/reader/version_api.go

// ListBookVersions 查询书籍版本列表
// @Summary 查询书籍版本列表
// @Description 查询指定书籍的所有版本
// @Tags BookVersion
// @Router /reader/books/{bookId}/versions [get]
func (api *VersionAPI) ListBookVersions(c *gin.Context) {
    bookID := c.Param("bookId")
    userID := api.GetUserID(c)

    versions, err := api.versionService.ListBookVersions(
        c.Request.Context(),
        bookID,
        userID,
    )
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{
        "data": versions,
        "total": len(versions),
    })
}

// GetVersionDiff 查询版本差异
// @Summary 查询版本差异
// @Description 对比两个版本之间的差异
// @Tags BookVersion
// @Router /reader/books/{bookId}/versions/diff [get]
func (api *VersionAPI) GetVersionDiff(c *gin.Context) {
    bookID := c.Param("bookId")
    fromVersion := c.Query("from")
    toVersion := c.Query("to")
    userID := api.GetUserID(c)

    diff, err := api.versionService.GetVersionDiff(
        c.Request.Context(),
        bookID,
        fromVersion,
        toVersion,
        userID,
    )
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, diff)
}

// SwitchVersion 切换阅读版本
// @Summary 切换阅读版本
// @Description 切换到指定版本进行阅读
// @Tags BookVersion
// @Router /reader/books/{bookId}/versions/{versionId}/switch [post]
func (api *VersionAPI) SwitchVersion(c *gin.Context) {
    bookID := c.Param("bookId")
    versionID := c.Param("versionId")
    userID := api.GetUserID(c)

    err := api.versionService.SwitchVersion(
        c.Request.Context(),
        bookID,
        versionID,
        userID,
    )
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"message": "已切换版本"})
}

// UpdateVersionPreference 更新版本偏好设置
// @Summary 更新版本偏好
// @Description 设置是否自动更新到新版本
// @Tags BookVersion
// @Router /reader/books/{bookId}/version-preference [put]
func (api *VersionAPI) UpdateVersionPreference(c *gin.Context) {
    bookID := c.Param("bookId")
    userID := api.GetUserID(c)

    var req UpdatePreferenceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    err := api.versionService.UpdatePreference(
        c.Request.Context(),
        bookID,
        userID,
        &req,
    )
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"message": "设置已更新"})
}

type UpdatePreferenceRequest struct {
    AutoUpdate   bool `json:"autoUpdate"`   // 是否自动更新
    NotifyUpdate bool `json:"notifyUpdate"` // 是否通知更新
}
```

### 订阅相关 API

```go
// api/v1/reader/subscription_api.go

// SubscribeBook 订阅书籍
// @Summary 订阅书籍
// @Description 订阅指定书籍
// @Tags Subscription
// @Router /reader/books/{bookId}/subscribe [post]
func (api *SubscriptionAPI) SubscribeBook(c *gin.Context) {
    bookID := c.Param("bookId")
    userID := api.GetUserID(c)

    var req SubscribeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    subscription, err := api.subscriptionService.SubscribeBook(
        c.Request.Context(),
        userID,
        bookID,
        &req,
    )
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(201, subscription)
}

// ListSubscriptions 查询订阅列表
// @Summary 查询订阅列表
// @Description 查询当前用户的订阅列表
// @Tags Subscription
// @Router /reader/subscriptions [get]
func (api *SubscriptionAPI) ListSubscriptions(c *gin.Context) {
    userID := api.GetUserID(c)
    status := c.Query("status")

    subscriptions, err := api.subscriptionService.ListSubscriptions(
        c.Request.Context(),
        userID,
        status,
    )
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{
        "data": subscriptions,
        "total": len(subscriptions),
    })
}

// GetSubscription 获取订阅详情
// @Summary 获取订阅详情
// @Tags Subscription
// @Router /reader/subscriptions/{id} [get]
func (api *SubscriptionAPI) GetSubscription(c *gin.Context) {
    id := c.Param("id")
    userID := api.GetUserID(c)

    subscription, err := api.subscriptionService.GetSubscription(
        c.Request.Context(),
        id,
        userID,
    )
    if err != nil {
        c.JSON(404, gin.H{"error": "订阅不存在"})
        return
    }

    c.JSON(200, subscription)
}

type SubscribeRequest struct {
    Type         SubscriptionType `json:"type" validate:"required"`
    AutoUpdate   bool              `json:"autoUpdate"`
    NotifyUpdate bool              `json:"notifyUpdate"`
    ChapterIDs   []string          `json:"chapterIds,omitempty"` // 按章订阅时指定
}
```

---

## 服务层设计

### PublishService（发布服务）

```go
// service/publishing/publish_service.go

type PublishService struct {
    publicationRepo interfaces.PublicationRepository
    projectRepo     interfaces.ProjectRepository
    bookRepo        interfaces.BookRepository
    versionRepo     interfaces.BookVersionRepository
    chapterRepo     interfaces.ChapterRepository
}

// CreatePublication 创建发布单
func (s *PublishService) CreatePublication(
    ctx context.Context,
    authorID string,
    req *CreatePublicationRequest,
) (*Publication, error) {
    // 1. 验证Project所有权
    project, err := s.projectRepo.GetByID(ctx, req.ProjectID)
    if err != nil {
        return nil, fmt.Errorf("项目不存在: %w", err)
    }
    if project.AuthorID != authorID {
        return nil, fmt.Errorf("无权限操作此项目")
    }

    // 2. 创建内容快照
    snapshot, err := s.createSnapshot(ctx, project, req)
    if err != nil {
        return nil, err
    }

    // 3. 创建发布单
    publication := &Publication{
        AuthorID:    authorID,
        ProjectID:   req.ProjectID,
        Type:        req.Type,
        Status:      PublicationStatusDraft,
        ChapterIDs:  req.ChapterIDs,
        Snapshot:    snapshot,
        Title:       req.Title,
        Description: req.Description,
    }

    if err := s.publicationRepo.Create(ctx, publication); err != nil {
        return nil, err
    }

    return publication, nil
}

// SubmitPublication 提交审核
func (s *PublishService) SubmitPublication(
    ctx context.Context,
    publicationID string,
    authorID string,
) (*Publication, error) {
    publication, err := s.publicationRepo.GetByID(ctx, publicationID)
    if err != nil {
        return nil, err
    }

    if publication.AuthorID != authorID {
        return nil, fmt.Errorf("无权限操作此发布单")
    }

    if publication.Status != PublicationStatusDraft {
        return nil, fmt.Errorf("只有草稿状态可以提交审核")
    }

    // 更新状态
    publication.Status = PublicationStatusSubmitted
    now := time.Now()
    publication.SubmitTime = &now

    if err := s.publicationRepo.Update(ctx, publication); err != nil {
        return nil, err
    }

    // 创建审核单
    review := &Review{
        PublicationID: publicationID,
        Status:       ReviewStatusPending,
        Priority:     ReviewPriorityNormal,
        Type:         publication.Type,
        Title:        publication.Title,
        History: []ReviewHistory{
            {
                Action:    "submitted",
                ActorID:   authorID,
                ActorRole: "author",
                Timestamp: time.Now(),
                Comment:   "提交审核",
            },
        },
    }

    if err := s.reviewRepo.Create(ctx, review); err != nil {
        return nil, err
    }

    // 关联审核单
    publication.ReviewID = &review.ID.Hex()
    s.publicationRepo.Update(ctx, publication)

    return publication, nil
}
```

### ReviewService（审核服务）

```go
// service/publishing/review_service.go

type ReviewService struct {
    reviewRepo      interfaces.ReviewRepository
    publicationRepo interfaces.PublicationRepository
    bookRepo        interfaces.BookRepository
    versionRepo     interfaces.BookVersionRepository
    notificationSvc interfaces.NotificationService
}

// ApproveReview 审核通过
func (s *ReviewService) ApproveReview(
    ctx context.Context,
    reviewID string,
    reviewerID string,
    req *ApproveRequest,
) (*ApprovalResult, error) {
    review, err := s.reviewRepo.GetByID(ctx, reviewID)
    if err != nil {
        return nil, err
    }

    if review.Status != ReviewStatusReviewing {
        return nil, fmt.Errorf("审核单状态不正确")
    }

    // 获取发布单
    publication, err := s.publicationRepo.GetByID(ctx, review.PublicationID)
    if err != nil {
        return nil, err
    }

    var result *ApprovalResult

    // 根据发布类型执行不同的逻辑
    switch publication.Type {
    case PublicationTypeFullBook:
        result, err = s.publishFullBook(ctx, publication, req)
    case PublicationTypeSingleChapter, PublicationTypeBatchChapter:
        result, err = s.publishChapterUpdate(ctx, publication, req)
    case PublicationTypeUpdate:
        result, err = s.publishBookUpdate(ctx, publication, req)
    }

    if err != nil {
        return nil, err
    }

    // 更新审核单状态
    review.Status = ReviewStatusApproved
    review.Result = &ReviewResult{
        Decision: ReviewDecisionApprove,
        Reason:   req.Comments,
    }
    now := time.Now()
    review.CompletedTime = &now
    s.reviewRepo.Update(ctx, review)

    // 更新发布单状态
    publication.Status = PublicationStatusApproved
    publication.ApprovedTime = &now
    publication.BookVersionID = &result.VersionID
    s.publicationRepo.Update(ctx, publication)

    // 通知作者
    s.notificationSvc.NotifyPublicationApproved(ctx, publication.AuthorID, reviewID)

    return result, nil
}

// publishFullBook 发布全书（创建新Book）
func (s *ReviewService) publishFullBook(
    ctx context.Context,
    publication *Publication,
    req *ApproveRequest,
) (*ApprovalResult, error) {
    // 创建Book
    book := &Book{
        Title:        publication.Snapshot.BookTitle,
        AuthorID:     publication.AuthorID,
        Introduction: publication.Snapshot.BookSummary,
        Cover:        publication.Snapshot.BookCover,
        CategoryIDs:  publication.Snapshot.CategoryIDs,
        Tags:         publication.Snapshot.BookTags,
        Status:       BookStatusOngoing,
        Price:        publication.Snapshot.Price,
        IsFree:       publication.Snapshot.IsFree,
        ProjectID:    &publication.ProjectID,
        PublicationID: &publication.ID.Hex(),
    }

    if err := s.bookRepo.Create(ctx, book); err != nil {
        return nil, err
    }

    // 创建BookVersion（初始版本）
    version := &BookVersion{
        BookID:         book.ID,
        PublicationID:  publication.ID.Hex(),
        VersionNumber:  1,
        VersionName:    "初始版本",
        VersionType:    VersionTypeInitial,
        Content:        s.createVersionContent(publication.Snapshot),
        PublishedAt:    time.Now(),
        PublishedBy:    publication.AuthorID,
        Status:         VersionStatusActive,
        IsLatest:       true,
    }

    if err := s.versionRepo.Create(ctx, version); err != nil {
        return nil, err
    }

    // 更新Book关联
    book.CurrentVersionID = &version.ID.Hex()
    book.LatestVersionID = &version.ID.Hex()
    book.TotalVersions = 1
    s.bookRepo.Update(ctx, book)

    // 发布章节
    for _, chapterSnap := range publication.Snapshot.Chapters {
        chapter := &Chapter{
            BookID:            book.ID,
            ProjectID:         &publication.ProjectID,
            ProjectChapterID: &chapterSnap.ProjectChapterID,
            ChapterNumber:     chapterSnap.ChapterNumber,
            Title:             chapterSnap.Title,
            Content:           chapterSnap.Content,
            WordCount:         chapterSnap.WordCount,
            Price:             chapterSnap.Price,
            IsFree:            chapterSnap.IsFree,
            ContentHash:       chapterSnap.ContentHash,
        }

        s.chapterRepo.Create(ctx, chapter)
    }

    return &ApprovalResult{
        BookID:     book.ID.Hex(),
        VersionID:  version.ID.Hex(),
        VersionNo:  1,
        Action:     "created",
    }, nil
}

// publishBookUpdate 发布更新（创建新版本）
func (s *ReviewService) publishBookUpdate(
    ctx context.Context,
    publication *Publication,
    req *ApproveRequest,
) (*ApprovalResult, error) {
    if publication.BookID == nil {
        return nil, fmt.Errorf("更新发布必须关联现有书籍")
    }

    book, err := s.bookRepo.GetByID(ctx, *publication.BookID)
    if err != nil {
        return nil, err
    }

    // 获取当前最新版本号
    latestVersionNo := book.TotalVersions + 1

    // 创建新版本
    version := &BookVersion{
        BookID:         book.ID,
        PublicationID:  publication.ID.Hex(),
        VersionNumber:  latestVersionNo,
        VersionName:    fmt.Sprintf("v%d", latestVersionNo),
        VersionType:    VersionTypeMinor,
        Content:        s.createVersionContent(publication.Snapshot),
        ChangeSummary:  s.calculateChanges(ctx, book.ID, publication.Snapshot),
        PublishedAt:    time.Now(),
        PublishedBy:    publication.AuthorID,
        Status:         VersionStatusActive,
        IsLatest:       true,
    }

    if err := s.versionRepo.Create(ctx, version); err != nil {
        return nil, err
    }

    // 更新Book
    book.LatestVersionID = &version.ID.Hex()
    book.TotalVersions = latestVersionNo
    s.bookRepo.Update(ctx, book)

    // 标记旧版本为非最新
    s.versionRepo.MarkOldVersions(ctx, book.ID.Hex(), version.ID.Hex())

    // 通知已订阅用户
    s.notifySubscribers(ctx, book.ID.Hex(), version.ID.Hex())

    return &ApprovalResult{
        BookID:     book.ID.Hex(),
        VersionID:  version.ID.Hex(),
        VersionNo:  latestVersionNo,
        Action:     "updated",
    }, nil
}

type ApprovalResult struct {
    BookID     string `json:"bookId"`
    VersionID  string `json:"versionId"`
    VersionNo  int    `json:"versionNo"`
    Action     string `json:"action"` // created, updated
}
```

### VersionService（版本服务）

```go
// service/bookstore/version_service.go

type VersionService struct {
    versionRepo      interfaces.BookVersionRepository
    bookRepo         interfaces.BookRepository
    subscriptionRepo interfaces.SubscriptionRepository
    diffService      *DiffService
}

// GetVersionDiff 获取版本差异
func (s *VersionService) GetVersionDiff(
    ctx context.Context,
    bookID string,
    fromVersion string,
    toVersion string,
    userID string,
) (*VersionDiff, error) {
    // 验证用户权限
    subscription, err := s.subscriptionRepo.GetByUserAndBook(ctx, userID, bookID)
    if err != nil {
        return nil, fmt.Errorf("未订阅该书籍")
    }

    // 获取两个版本
    from, err := s.versionRepo.GetByID(ctx, fromVersion)
    if err != nil {
        return nil, err
    }

    to, err := s.versionRepo.GetByID(ctx, toVersion)
    if err != nil {
        return nil, err
    }

    // 计算差异
    diff := s.diffService.CompareVersions(from, to)

    return diff, nil
}

// SwitchVersion 切换阅读版本
func (s *VersionService) SwitchVersion(
    ctx context.Context,
    bookID string,
    versionID string,
    userID string,
) error {
    subscription, err := s.subscriptionRepo.GetByUserAndBook(ctx, userID, bookID)
    if err != nil {
        return fmt.Errorf("未订阅该书籍")
    }

    version, err := s.versionRepo.GetByID(ctx, versionID)
    if err != nil {
        return err
    }

    if version.BookID.Hex() != bookID {
        return fmt.Errorf("版本不属于该书")
    }

    // 更新订阅的当前版本
    subscription.CurrentVersionID = &versionID
    return s.subscriptionRepo.Update(ctx, subscription)
}

type VersionDiff struct {
    FromVersion    int              `json:"fromVersion"`
    ToVersion      int              `json:"toVersion"`
    Summary        *ChangeSummary  `json:"summary"`
    ChapterChanges []ChapterDiff   `json:"chapterChanges"`
    DiffStats      DiffStats       `json:"diffStats"`
}

type ChapterDiff struct {
    ChapterNumber  int             `json:"chapterNumber"`
    Title          string          `json:"title"`
    Action         ChapterAction   `json:"action"`
    Diff           *TextDiff       `json:"diff,omitempty"`
}

type DiffStats struct {
    AddedWords     int             `json:"addedWords"`
    DeletedWords   int             `json:"deletedWords"`
    ChangedWords   int             `json:"changedWords"`
}
```

---

## 差异对比服务

```go
// service/bookstore/diff_service.go

type DiffService struct{}

// CompareVersions 对比两个版本
func (s *DiffService) CompareVersions(
    from *BookVersion,
    to *BookVersion,
) *VersionDiff {
    diff := &VersionDiff{
        FromVersion: from.VersionNumber,
        ToVersion:   to.VersionNumber,
        Summary:     to.ChangeSummary,
    }

    // 对比章节
    diff.ChapterChanges = s.compareChapters(from.Content.Chapters, to.Content.Chapters)

    // 计算统计
    diff.DiffStats = s.calculateStats(diff.ChapterChanges)

    return diff
}

// compareChapters 对比章节列表
func (s *DiffService) compareChapters(
    fromChapters []VersionChapter,
    toChapters []VersionChapter,
) []ChapterDiff {
    var changes []ChapterDiff

    // 创建章节映射
    fromMap := make(map[string]VersionChapter)
    toMap := make(map[string]VersionChapter)

    for _, ch := range fromChapters {
        fromMap[ch.ProjectChapterID] = ch
    }

    for _, ch := range toChapters {
        toMap[ch.ProjectChapterID] = ch
    }

    // 检查新增和修改
    for projectChapterID, toChapter := range toMap {
        fromChapter, exists := fromMap[projectChapterID]

        if !exists {
            // 新增章节
            changes = append(changes, ChapterDiff{
                ChapterNumber: toChapter.ChapterNumber,
                Title:         toChapter.Title,
                Action:        ChapterActionAdded,
            })
        } else if fromChapter.ContentHash != toChapter.ContentHash {
            // 内容修改
            changes = append(changes, ChapterDiff{
                ChapterNumber: toChapter.ChapterNumber,
                Title:         toChapter.Title,
                Action:        ChapterActionModified,
                Diff:          s.computeTextDiff(fromChapter.Content, toChapter.Content),
            })
        }
    }

    // 检查删除
    for projectChapterID, fromChapter := range fromMap {
        if _, exists := toMap[projectChapterID]; !exists {
            changes = append(changes, ChapterDiff{
                ChapterNumber: fromChapter.ChapterNumber,
                Title:         fromChapter.Title,
                Action:        ChapterActionDeleted,
            })
        }
    }

    return changes
}

// computeTextDiff 计算文本差异
func (s *DiffService) computeTextDiff(from, to string) *TextDiff {
    // 使用 diff 算法
    // 这里简化实现，实际可以使用 github.com/sergi/go-diff/diffmatchpatch
    linesFrom := strings.Split(from, "\n")
    linesTo := strings.Split(to, "\n")

    return &TextDiff{
        Additions:    s.findAddedLines(linesFrom, linesTo),
        Deletions:    s.findDeletedLines(linesFrom, linesTo),
        Modified:     s.findModifiedLines(linesFrom, linesTo),
        TotalChanges: len(linesTo) - len(linesFrom),
    }
}

type TextDiff struct {
    Additions    []string `json:"additions"`
    Deletions    []string `json:"deletions"`
    Modified     []string `json:"modified"`
    TotalChanges int      `json:"totalChanges"`
}
```

---

## 实施计划

### Phase 1: 数据模型（2天）

- [ ] 创建Publication模型和Repository
- [ ] 创建Review模型和Repository
- [ ] 创建BookVersion模型和Repository
- [ ] 更新Book模型添加关联字段
- [ ] 创建Subscription模型和Repository
- [ ] 创建数据库迁移脚本

### Phase 2: 核心服务（3天）

- [ ] 实现PublishService（创建、提交、发布）
- [ ] 实现ReviewService（审核、通过、拒绝）
- [ ] 实现VersionService（版本管理、差异对比）
- [ ] 实现DiffService（文本差异计算）

### Phase 3: API层（2天）

- [ ] 发布相关API
- [ ] 审核相关API
- [ ] 版本相关API
- [ ] 订阅相关API

### Phase 4: 前端（3天）

- [ ] 发布对话框UI
- [ ] 审核管理后台
- [ ] 版本对比界面
- [ ] 订阅管理界面

### Phase 5: 测试（2天）

- [ ] 单元测试
- [ ] 集成测试
- [ ] E2E测试

---

## 相关文档

- [Project-Book分离架构](./2026-03-06-project-book-separation-architecture.md)
- [BookStatus枚举统一](./2026-03-05-book-status-unification-design.md)
- [事务管理器设计](./2026-03-05-transaction-manager-design.md)

---

**设计完成时间**: 2026-03-06
**预计实施时间**: 12天
**建议执行者**: 后端团队 + 前端团队协同
