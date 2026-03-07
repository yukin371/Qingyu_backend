# 基于现有3层架构的版本管理设计方案

**设计日期**: 2026-03-06
**设计者**: Kore
**优先级**: 🔴 P0
**原则**: 最小重构，利用现有架构

---

## 当前3层架构分析

### Writer域（写作）

```
┌─────────────────────────────────────────────────────────────────────────┐
│  Project (项目元数据)                                                      │
│  ├─ ID, Title, AuthorID, Status                                         │
│  └─ 1:N                                                                  │
└─────────────────────────────────────────────────────────────────────────┘
         │ 1:N
         ▼
┌─────────────────────────────────────────────────────────────────────────┐
│  Document (文档元数据)                                                    │
│  ├─ ID, ProjectID, Title, Type(volume/chapter/section/scene)              │
│  ├─ ParentID, Level, Order, Status                                         │
│  └─ StableRef, OrderKey                                                   │
└─────────────────────────────────────────────────────────────────────────┘
         │ 1:N
         ▼
┌─────────────────────────────────────────────────────────────────────────┐
│  DocumentContent (内容层 - 分段存储)                                    │
│  ├─ ID, DocumentID                                                       │
│  ├─ Content (tiptap/markdown/richtext)                                  │
│  ├─ ContentType, ParagraphOrder (段落顺序)                               │
│  ├─ WordCount, Version                                                   │
│  └─ 支持多段落存储：一篇文档可以有多个DocumentContent记录                  │
└─────────────────────────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────┐
│  Paragraph (领域层 - 轻量类型，不单独建表)                                │
│  ├─ 从DocumentContent解析得到                                           │
│  └─ 用于编辑器和评论段落绑定                                               │
└─────────────────────────────────────────────────────────────────────────┘
```

### Bookstore域（阅读）

```
┌─────────────────────────────────────────────────────────────────────────┐
│  Book (书籍元数据)                                                        │
│  ├─ ID, Title, Author, Status                                             │
│  └─ 1:N                                                                  │
└─────────────────────────────────────────────────────────────────────────┘
         │ 1:N
         ▼
┌─────────────────────────────────────────────────────────────────────────┐
│  Chapter (章节元数据)                                                     │
│  ├─ ID, BookID, Title, ChapterNum                                        │
│  ├─ Price, IsFree, WordCount                                             │
│  ├─ ContentURL, ContentHash, ContentVersion                             │
│  └─ 1:N                                                                  │
└─────────────────────────────────────────────────────────────────────────┘
         │ 1:N
         ▼
┌─────────────────────────────────────────────────────────────────────────┐
│  ChapterContent (章节内容)                                              │
│  ├─ ID, ChapterID                                                       │
│  ├─ Content (markdown/html/txt)                                        │
│  ├─ Format, Version                                                      │
│  └─ Hash, WordCount                                                      │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 版本管理需求分析

### 问题：当前架构缺少版本管理

1. ❌ **DocumentContent** 虽然有Version字段，但没有历史记录
2. ❌ **ChapterContent** 同样只有当前版本，无历史追踪
3. ❌ **发布快照** 无法记录发布时的内容状态
4. ❌ **版本对比** 无法查看新旧版本差异

### 解决方案：最小重构

**核心思想**：复用现有3层架构，增加版本关联层

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              版本管理层（新增）                             │
│                                                                               │
│  DocumentVersion (文档版本)     BookVersion (书籍版本)                   │
│  ├─ DocumentID                  ├─ BookID                             │
│  ├─ VersionNumber               ├─ VersionNumber                      │
│  ├─ SnapshotRef (引用)         ├─ SnapshotRef                         │
│  └─ ChangeSummary              └─ ChangeSummary                        │
│         │ 引用                      │ 引用                                    │
│         ▼                           ▼                                        │
│  DocumentContentSnapshot    ChapterContentSnapshot                           │
│  (完整内容快照)               (完整内容快照)                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

**关键设计决策**：
- **复用现有存储**：不改变DocumentContent和ChapterContent的存储方式
- **快照独立存储**：版本快照存储在新的集合中
- **引用而非复制**：使用SnapshotRef引用，按需加载完整内容
- **段落级版本管理**：Paragraph级别的差异对比

---

## 数据模型设计

### 1. DocumentVersion（文档版本）

```go
// models/writer/document_version.go

package writer

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "Qingyu_backend/models/shared"
)

// DocumentVersion 文档版本
type DocumentVersion struct {
    shared.IdentifiedEntity `bson:",inline"`
    shared.BaseEntity       `bson:",inline"`

    // 关联信息
    DocumentID     primitive.ObjectID `bson:"document_id" json:"documentId"`     // ← 关联Document
    PublicationID  *string            `bson:"publication_id,omitempty" json:"publicationId,omitempty"` // ← 来源发布单

    // 版本标识
    VersionNumber  int                `bson:"version_number" json:"versionNumber"` // 版本号
    VersionName    string             `bson:"version_name" json:"versionName"`    // 版本名称
    VersionType    DocumentVersionType `bson:"version_type" json:"versionType"` // 版本类型

    // ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
    //  快照引用（而非完整内容复制）
    // ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
    SnapshotRef    *string            `bson:"snapshot_ref,omitempty" json:"snapshotRef,omitempty"` // 快照引用ID
    SnapshotSize   int64              `bson:"snapshot_size,omitempty" json:"snapshotSize,omitempty"` // 快照大小

    // 变更摘要
    ChangeSummary  *DocumentChangeSummary `bson:"change_summary,omitempty" json:"changeSummary,omitempty"` // 变更摘要
    DiffData        *ParagraphDiffData `bson:"diff_data,omitempty" json:"diffData,omitempty"` // 段落级差异

    // 发布信息
    PublishedAt     time.Time          `bson:"published_at" json:"publishedAt"`
    PublishedBy     string             `bson:"published_by" json:"publishedBy"`

    // 版本状态
    Status         DocumentVersionStatus `bson:"status" json:"status"`
    IsLatest       bool               `bson:"is_latest" json:"isLatest"`       // 是否最新版本

    // 统计信息
    ViewCount      int64              `bson:"view_count" json:"viewCount"`

    // 元数据
    Description     string             `bson:"description,omitempty" json:"description,omitempty"`
    Changelog       string             `bson:"changelog,omitempty" json:"changelog,omitempty"`
    Tags            []string           `bson:"tags,omitempty" json:"tags,omitempty"`
}

// DocumentVersionType 文档版本类型
type DocumentVersionType string

const (
    DocumentVersionTypeInitial   DocumentVersionType = "initial"   // 初始版本
    DocumentVersionTypeMajor     DocumentVersionType = "major"     // 主要版本
    DocumentVersionTypeMinor     DocumentVersionType = "minor"     // 次要版本（新增段落）
    DocumentVersionTypePatch     DocumentVersionType = "patch"     // 补丁版本（修改段落）
    DocumentVersionTypeEmergency DocumentVersionType = "emergency" // 紧急版本
)

// DocumentVersionStatus 文档版本状态
type DocumentVersionStatus string

const (
    DocumentVersionStatusActive     DocumentVersionStatus = "active"     // 活跃版本
    DocumentVersionStatusArchived   DocumentVersionStatus = "archived"   // 已归档
    DocumentVersionStatusDeprecated DocumentVersionStatus = "deprecated" // 已废弃
)

// DocumentChangeSummary 文档变更摘要
type DocumentChangeSummary struct {
    ParagraphsAdded     int              `bson:"paragraphs_added" json:"paragraphsAdded"`     // 新增段落数
    ParagraphsModified  int              `bson:"paragraphs_modified" json:"paragraphsModified"`  // 修改段落数
    ParagraphsDeleted   int              `bson:"paragraphs_deleted" json:"paragraphsDeleted"`   // 删除段落数
    WordsChanged        int64            `bson:"words_changed" json:"wordsChanged"`        // 字数变化
    StorageSaved        int64            `bson:"storage_saved" json:"storageSaved"`        // 节省存储（字节）
}

// ParagraphDiffData 段落级差异数据
type ParagraphDiffData struct {
    ParagraphsChanged   []ParagraphDiff   `bson:"paragraphs_changed" json:"paragraphsChanged"` // 变更的段落
    TotalParagraphs     int               `bson:"total_paragraphs" json:"totalParagraphs"`    // 总段落数
}

// ParagraphDiff 单个段落的差异
type ParagraphDiff struct {
    ParagraphOrder    int               `bson:"paragraph_order" json:"paragraphOrder"` // 段落顺序
    Action            DiffAction       `bson:"action" json:"action"`                   // added, modified, deleted

    // 版本关联
    CurrentVersionID  *string          `bson:"current_version_id,omitempty" json:"currentVersionId,omitempty"` // 当前段落版本ID
    PreviousVersionID *string          `bson:"previous_version_id,omitempty" json:"previousVersionId,omitempty"` // 上一段落版本ID

    // 差异统计
    WordsChanged      int               `bson:"words_changed" json:"wordsChanged"`       // 字数变化
    ChangedRatio      float64           `bson:"changed_ratio" json:"changedRatio"`       // 变更比例 0-1
}

// DiffAction 差异动作
type DiffAction string

const (
    DiffActionAdded    DiffAction = "added"    // 新增
    DiffActionModified DiffAction = "modified" // 修改
    DiffActionDeleted  DiffAction = "deleted"  // 删除
)
```

### 2. DocumentContentSnapshot（文档快照）

```go
// models/writer/document_content_snapshot.go

package writer

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// DocumentContentSnapshot 文档内容快照
// 存储发布时的DocumentContent完整状态，用于版本恢复和对比
type DocumentContentSnapshot struct {
    shared.IdentifiedEntity `bson:",inline"`
    shared.BaseEntity       `bson:",inline"`

    // 关联信息
    DocumentID      primitive.ObjectID `bson:"document_id" json:"documentId"`
    DocumentVersionID string             `bson:"document_version_id" json:"documentVersionId"` // ← 关联DocumentVersion

    // 快照内容（完整复制DocumentContent）
    ContentType     string             `bson:"content_type" json:"contentType"`       // markdown | richtext | tiptap
    Content         string             `bson:"content" json:"content"`                 // 完整内容
    ParagraphOrder  int                `bson:"paragraph_order" json:"paragraphOrder"` // 段落顺序

    // 元数据（从DocumentContent复制）
    WordCount       int                `bson:"word_count" json:"wordCount"`
    CharCount       int                `bson:"char_count" json:"charCount"`
    Version         int                `bson:"version" json:"version"`           // 原始版本号

    // 快照元数据
    SnapshotTime    time.Time          `bson:"snapshot_time" json:"snapshotTime"`    // 快照时间
    SnapshotBy      string             `bson:"snapshot_by" json:"snapshotBy"`      // 快照创建者

    // 内容校验
    ContentHash     string             `bson:"content_hash" json:"contentHash"`     // 内容哈希

    // 压缩（可选）
    Compressed      bool               `bson:"compressed" json:"compressed"`      // 是否压缩
    OriginalSize    int64              `bson:"original_size" json:"originalSize"`    // 原始大小
    CompressedSize  int64              `bson:"compressed_size" json:"compressedSize"` // 压缩后大小
}

// GetDisplaySize 获取显示用的快照大小
func (s *DocumentContentSnapshot) GetDisplaySize() string {
    if s.Compressed {
        return formatBytes(s.CompressedSize)
    }
    return formatBytes(s.OriginalSize)
}

// formatBytes 格式化字节大小
func formatBytes(bytes int64) string {
    const unit = 1024
    if bytes < unit {
        return fmt.Sprintf("%d B", bytes)
    }
    div, exp := int64(unit), 0
    for n := bytes / unit; n >= unit; n /= unit {
        div++
        exp++
    }
    return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(int64(math.Pow(unit, float64(exp)))), "KMGTPE"[exp])
}
```

### 3. BookVersion（书籍版本）

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
    BookID          primitive.ObjectID `bson:"book_id" json:"bookId"`
    PublicationID   string             `bson:"publication_id" json:"publicationId"`

    // 版本标识
    VersionNumber   int                `bson:"version_number" json:"versionNumber"`
    VersionName     string             `bson:"version_name" json:"versionName"`
    VersionType     BookVersionType    `bson:"version_type" json:"versionType"`

    // ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
    //  快照引用（而非完整内容复制）
    // ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
    SnapshotRef     *string            `bson:"snapshot_ref,omitempty" json:"snapshotRef,omitempty"` // 快照引用ID
    SnapshotSize    int64              `bson:"snapshot_size,omitempty" json:"snapshotSize,omitempty"` // 快照大小

    // 变更摘要
    ChangeSummary   *BookChangeSummary `bson:"change_summary,omitempty" json:"changeSummary"`
    DiffData         *ChapterDiffData   `bson:"diff_data,omitempty" json:"diffData,omitempty"`

    // 发布信息
    PublishedAt      time.Time          `bson:"published_at" json:"publishedAt"`
    PublishedBy      string             `bson:"published_by" json:"publishedBy"`

    // 版本状态
    Status          BookVersionStatus  `bson:"status" json:"status"`
    IsLatest        bool               `bson:"is_latest" json:"isLatest"`

    // 统计信息
    ViewCount       int64              `bson:"view_count" json:"viewCount"`

    // 元数据
    Description     string             `bson:"description,omitempty" json:"description,omitempty"`
    Changelog       string             `bson:"changelog,omitempty" json:"changelog,omitempty"`
    Tags            []string           `bson:"tags,omitempty" json:"tags,omitempty"`
}

// BookVersionType 书籍版本类型
type BookVersionType string

const (
    BookVersionTypeInitial   BookVersionType = "initial"   // 初始版本
    BookVersionTypeMajor     BookVersionType = "major"     // 主要版本
    BookVersionTypeMinor     BookVersionType = "minor"     // 次要版本（新增章节）
    BookVersionTypePatch     BookVersionType = "patch"     // 补丁版本（修改章节）
    BookVersionTypeEmergency BookVersionType = "emergency" // 紧急版本
)

// BookVersionStatus 书籍版本状态
type BookVersionStatus string

const (
    BookVersionStatusActive     BookVersionStatus = "active"     // 活跃版本
    BookVersionStatusArchived   BookVersionStatus = "archived"   // 已归档
    BookVersionStatusDeprecated BookVersionStatus = "deprecated" // 已废弃
)

// BookChangeSummary 书籍变更摘要
type BookChangeSummary struct {
    ChaptersAdded     int               `bson:"chapters_added" json:"chaptersAdded"`     // 新增章节数
    ChaptersModified  int               `bson:"chapters_modified" json:"chaptersModified"`  // 修改章节数
    ChaptersDeleted   int               `bson:"chapters_deleted" json:"chaptersDeleted"`   // 删除章节数
    WordsChanged     int64             `bson:"words_changed" json:"wordsChanged"`     // 字数变化
    StorageSaved      int64             `bson:"storage_saved" json:"storageSaved"`      // 节省存储（字节）
}

// ChapterDiffData 章节级差异数据
type ChapterDiffData struct {
    ChaptersChanged   []ChapterDiff     `bson:"chapters_changed" json:"chaptersChanged"` // 变更的章节
    TotalChapters    int               `bson:"total_chapters" json:"totalChapters"`    // 总章节数
}

// ChapterDiff 单个章节的差异
type ChapterDiff struct {
    ChapterID       string            `bson:"chapter_id" json:"chapterId"`
    ChapterNumber   int               `bson:"chapter_number" json:"chapterNumber"`
    Title           string            `bson:"title" json:"title"`
    Action          DiffAction        `bson:"action" json:"action"`                   // added, modified, deleted

    // 版本关联
    CurrentVersionID *string          `bson:"current_version_id,omitempty" json:"currentVersionId,omitempty"` // 当前章节版本ID
    PreviousVersionID *string          `bson:"previous_version_id,omitempty" json:"previousVersionId,omitempty"` // 上一章节版本ID

    // 差异统计
    WordsChanged    int               `bson:"words_changed" json:"wordsChanged"`       // 字数变化
    ChangedRatio    float64           `bson:"changed_ratio" json:"changedRatio"`       // 变更比例 0-1

    // 段落级差异（详细）
    ParagraphsChanged []ParagraphDiff `bson:"paragraphs_changed,omitempty" json:"paragraphsChanged,omitempty"` // 段落级变更
}
```

### 4. ChapterContentSnapshot（章节内容快照）

```go
// models/bookstore/chapter_content_snapshot.go

package bookstore

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// ChapterContentSnapshot 章节内容快照
// 存储发布时的ChapterContent完整状态
type ChapterContentSnapshot struct {
    shared.IdentifiedEntity `bson:",inline"`
    shared.BaseEntity       `bson:",inline"`

    // 关联信息
    ChapterID       primitive.ObjectID `bson:"chapter_id" json:"chapterId"`
    ChapterVersionID string             `bson:"chapter_version_id" json:"chapterVersionId"` // ← 关联BookVersion

    // 快照内容（完整复制ChapterContent）
    Content         string             `bson:"content" json:"content"`                   // 完整内容
    Format          string             `bson:"format" json:"format"`                     // markdown | html | txt
    Version         int                `bson:"version" json:"version"`                    // 原始版本号

    // 元数据（从ChapterContent复制）
    WordCount       int                `bson:"word_count" json:"wordCount"`
    Hash            string             `bson:"hash,omitempty" json:"hash,omitempty"`

    // 快照元数据
    SnapshotTime    time.Time          `bson:"snapshot_time" json:"snapshotTime"`    // 快照时间
    SnapshotBy      string             `bson:"snapshot_by" json:"snapshotBy"`      // 快照创建者

    // 压缩（可选）
    Compressed      bool               `bson:"compressed" json:"compressed"`      // 是否压缩
    OriginalSize    int64              `bson:"original_size" json:"originalSize"`    // 原始大小
    CompressedSize  int64              `bson:"compressed_size" json:"compressedSize"` // 压缩后大小
}
```

---

## 版本创建策略

### 存储类型决策算法

```go
// service/writer/version_strategy.go

type VersionStrategy struct {
    fullThreshold      float64 // >50%变化使用Full快照
    diffThreshold      float64 // <=50%变化使用Diff存储
}

// DecideSnapshotType 决定快照类型
func (s *VersionStrategy) DecideSnapshotType(
    oldContent *DocumentContent,
    newContent *DocumentContent,
) (SnapshotType, *ChangeStats) {
    // 首次发布：必须Full
    if oldContent == nil {
        return SnapshotTypeFull, &ChangeStats{ChangeRatio: 1.0}
    }

    // 计算差异
    stats := s.calculateStats(oldContent, newContent)

    // 决策
    if stats.ChangeRatio > s.fullThreshold {
        return SnapshotTypeFull, stats
    }

    return SnapshotTypeDiff, stats
}

// SnapshotType 快照类型
type SnapshotType string

const (
    SnapshotTypeFull SnapshotType = "full" // 完整快照
    SnapshotTypeDiff SnapshotType = "diff" // 差异快照
)

// ChangeStats 变更统计
type ChangeStats struct {
    ParagraphsAdded     int     `json:"paragraphsAdded"`
    ParagraphsModified  int     `json:"paragraphsModified"`
    ParagraphsDeleted   int     `json:"paragraphsDeleted"`
    WordsChanged        int64   `json:"wordsChanged"`
    TotalParagraphs     int     `json:"totalParagraphs"`
    ChangeRatio         float64 `json:"changeRatio"` // 0-1
}
```

---

## 发布流程设计

### 核心流程：保持现有架构，增加版本层

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                          发布流程（最小重构）                                  │
└─────────────────────────────────────────────────────────────────────────────────┘

Author                    PublishService                ReviewService               Reader
   │                           │                              │                     │
   │  1. 创建Publication          │                              │                     │
   │  (指定Document版本)        │                              │                     │
   │──────────────────────────>│                              │                     │
   │                           │  2. 创建DocumentVersion          │                     │
   │                           │  (DocumentContentSnapshot)      │                     │
   │                           │                              │                     │
   │  3. 提交审核                │                              │                     │
   │  <───────────────────────│─────────────────────────────>│                     │
   │                           │                              │                     │
   │                           │  4. 审核通过                    │                     │
   │                           │  ──────────────────────────────>│                     │
   │                           │                              │                     │
   │                           │  5. 创建BookVersion             │                     │
   │                           │  (ChapterContentSnapshot)       │                     │
   │                           │                              │                     │
   │                           │  6. 书籍上架                    │────────────────────>│
│                           │─────────────────────────────────────────────────────│
│                           │                              │                     │
│  7. 继续写作                │                              │                     │
│   (修改DocumentContent)     │                              │                     │
│  │                          │                              │                     │
│  8. 创建更新Publication       │                              │                     │
│  │──────────────────────────>│                              │                     │
│                           │  9. 审核通过 → 新BookVersion       │                     │
│                           │                              │                     │
│                           │ 10. 通知已订阅用户                │────────────────────>│
│                           │                              │                     │
│                           │                              │                     │
│                           │ 11. 读者查看版本差异              │<────────────────────│
│                           │<─────────────────────────────│                     │
```

---

## 服务层实现

### PublishService（发布服务）

```go
// service/writer/publish_service.go

type PublishService struct {
    publicationRepo     interfaces.PublicationRepository
    documentRepo       interfaces.DocumentRepository
    documentContentRepo interfaces.DocumentContentRepository
    bookRepo           interfaces.BookRepository
    chapterRepo        interfaces.ChapterRepository
    documentVersionRepo interfaces.DocumentVersionRepository
    bookVersionRepo    interfaces.BookVersionRepository
    snapshotService    *SnapshotService
}

// PublishDocument 发布文档为书籍章节
func (s *PublishService) PublishDocument(
    ctx context.Context,
    authorID string,
    req *PublishDocumentRequest,
) (*Book, error) {
    // 1. 获取Document
    document, err := s.documentRepo.GetByID(ctx, req.DocumentID)
    if err != nil {
        return nil, fmt.Errorf("文档不存在: %w", err)
    }

    // 验证权限
    if document.ProjectID.AuthorID != authorID {
        return nil, fmt.Errorf("无权限操作此文档")
    }

    // 2. 获取DocumentContent
    documentContent, err := s.documentContentRepo.GetByDocumentID(ctx, req.DocumentID)
    if err != nil {
        return nil, fmt.Errorf("文档内容不存在: %w", err)
    }

    // 3. 创建DocumentVersion（元数据层版本）
    documentVersion, err := s.createDocumentVersion(
        ctx,
        document,
        documentContent,
        req,
    )
    if err != nil {
        return nil, err
    }

    // 4. 创建快照（内容层版本）
    snapshot, err := s.snapshotService.CreateDocumentSnapshot(
        ctx,
        documentContent,
        documentVersion,
    )
    if err != nil {
        return nil, err
    }

    // 5. 创建或更新Book/Chapter
    book, err := s.publishToBookstore(
        ctx,
        document,
        documentContent,
        snapshot,
        req,
    )
    if err != nil {
        return nil, err
    }

    // 6. 关联版本
    documentVersion.SnapshotRef = &snapshot.ID.Hex()
    s.documentVersionRepo.Update(ctx, documentVersion)

    return book, nil
}

// createDocumentVersion 创建文档版本
func (s *PublishService) createDocumentVersion(
    ctx context.Context,
    document *Document,
    content *DocumentContent,
    req *PublishDocumentRequest,
) (*DocumentVersion, error) {
    // 检查是否有上一版本
    latestVersion, _ := s.documentVersionRepo.GetLatestByDocumentID(
        ctx,
        document.ID.Hex(),
    )

    versionNumber := 1
    versionType := DocumentVersionTypeInitial
    if latestVersion != nil {
        versionNumber = latestVersion.VersionNumber + 1
        versionType = DocumentVersionTypeMinor // 默认为次要版本
    }

    documentVersion := &DocumentVersion{
        DocumentID:     document.ID,
        VersionNumber:  versionNumber,
        VersionName:    fmt.Sprintf("v%d", versionNumber),
        VersionType:    versionType,
        PublishedAt:    time.Now(),
        PublishedBy:    document.ProjectID.AuthorID,
        Status:         DocumentVersionStatusActive,
        IsLatest:       true,
    }

    if err := s.documentVersionRepo.Create(ctx, documentVersion); err != nil {
        return nil, err
    }

    // 标记旧版本为非最新
    if latestVersion != nil {
        latestVersion.IsLatest = false
        s.documentVersionRepo.Update(ctx, latestVersion)
    }

    return documentVersion, nil
}

// publishToBookstore 发布到书店域
func (s *PublishService) publishToBookstore(
    ctx context.Context,
    document *Document,
    content *DocumentContent,
    snapshot *DocumentContentSnapshot,
    req *PublishDocumentRequest,
) (*Book, error) {
    // 检查是否有关联的Book
    book, _ := s.bookRepo.GetByProjectIDAndDocumentID(
        ctx,
        document.ProjectID.Hex(),
        document.ID.Hex(),
    )

    var err error

    if book == nil {
        // 首次发布：创建新Book和Chapter
        book, err = s.createBookAndChapter(
            ctx,
            document,
            content,
            snapshot,
            req,
        )
    } else {
        // 更新现有Book和Chapter
        book, err = s.updateBookAndChapter(
            ctx,
            book,
            document,
            content,
            snapshot,
            req,
        )
    }

    return book, err
}

// createBookAndChapter 创建新Book和Chapter
func (s *PublishService) createBookAndChapter(
    ctx context.Context,
    document *Document,
    content *DocumentContent,
    snapshot *DocumentContentSnapshot,
    req *PublishDocumentRequest,
) (*Book, error) {
    // 创建Book
    book := &Book{
        Title:        document.Title,
        Author:       getAuthorName(ctx, document.ProjectID.AuthorID),
        AuthorID:     document.ProjectID.AuthorID,
        Introduction: generateIntroduction(content),
        CategoryIDs:  req.CategoryIDs,
        Tags:         req.Tags,
        Status:       BookStatusOngoing,
        ProjectID:    &document.ProjectID.Hex(),
        WordCount:    int64(content.WordCount),
        ChapterCount: 1,
    }

    if err := s.bookRepo.Create(ctx, book); err != nil {
        return nil, err
    }

    // 创建Chapter
    chapter := &Chapter{
        BookID:       book.ID,
        Title:        document.Title,
        ChapterNum:   req.ChapterNumber,
        WordCount:    content.WordCount,
        Price:        req.Price,
        IsFree:       req.IsFree,
    }

    if err := s.chapterRepo.Create(ctx, chapter); err != nil {
        return nil, err
    }

    // 创建ChapterContent
    chapterContent := &ChapterContent{
        ChapterID:    chapter.ID,
        Content:      content.Content,
        Format:       content.ContentType,
        WordCount:    content.WordCount,
        Version:      1,
    }

    if err := s.chapterContentRepo.Create(ctx, chapterContent); err != nil {
        return nil, err
    }

    // 更新Chapter的Content引用
    chapter.ContentURL = fmt.Sprintf("/api/v1/reader/chapters/%s/content", chapter.ID.Hex())
    chapter.ContentSize = int64(len(content.Content))
    chapter.ContentHash = snapshot.ContentHash
    chapter.ContentVersion = 1

    s.chapterRepo.Update(ctx, chapter)

    return book, nil
}
```

### SnapshotService（快照服务）

```go
// service/writer/snapshot_service.go

type SnapshotService struct {
    documentSnapshotRepo interfaces.DocumentContentSnapshotRepository
    chapterSnapshotRepo  interfaces.ChapterContentSnapshotRepository
    compressionService    *CompressionService
}

// CreateDocumentSnapshot 创建文档快照
func (s *SnapshotService) CreateDocumentSnapshot(
    ctx context.Context,
    content *DocumentContent,
    version *DocumentVersion,
) (*DocumentContentSnapshot, error) {
    snapshot := &DocumentContentSnapshot{
        DocumentID:       content.DocumentID,
        DocumentVersionID: version.ID.Hex(),
        ContentType:      content.ContentType,
        Content:          content.Content,
        ParagraphOrder:    content.ParagraphOrder,
        WordCount:        content.WordCount,
        CharCount:        content.CharCount,
        Version:          content.Version,
        SnapshotTime:     time.Now(),
        ContentHash:       calculateContentHash(content.Content),
        SnapshotBy:       version.PublishedBy,
    }

    // 可选压缩
    if len(content.Content) > 100*1024 { // >100KB
        compressed, err := s.compressionService.Compress(content.Content)
        if err == nil && len(compressed) < len(content.Content)*9/10 { // 压缩率>10%
            snapshot.Content = string(compressed)
            snapshot.Compressed = true
            snapshot.OriginalSize = int64(len(content.Content))
            snapshot.CompressedSize = int64(len(compressed))
        }
    }

    if err := s.documentSnapshotRepo.Create(ctx, snapshot); err != nil {
        return nil, err
    }

    return snapshot, nil
}

// GetDocumentSnapshot 获取文档快照
func (s *SnapshotService) GetDocumentSnapshot(
    ctx context.Context,
    snapshotID string,
) (*DocumentContentSnapshot, error) {
    snapshot, err := s.documentSnapshotRepo.GetByID(ctx, snapshotID)
    if err != nil {
        return nil, err
    }

    // 如果压缩了，解压
    if snapshot.Compressed {
        decompressed, err := s.compressionService.Decompress(snapshot.Content)
        if err != nil {
            return nil, err
        }

        // 返回解压后的内容
        result := *snapshot
        result.Content = string(decompressed)
        return &result, nil
    }

    return snapshot, nil
}

// CompareDocumentSnapshots 对比两个文档快照
func (s *SnapshotService) CompareDocumentSnapshots(
    ctx context.Context,
    fromSnapshotID string,
    toSnapshotID string,
) (*DocumentDiffResult, error) {
    from, err := s.GetDocumentSnapshot(ctx, fromSnapshotID)
    if err != nil {
        return nil, err
    }

    to, err := s.GetDocumentSnapshot(ctx, toSnapshotID)
    if err != nil {
        return nil, err
    }

    // 解析内容为段落
    fromParagraphs := s.parseParagraphs(from.Content, from.ContentType)
    toParagraphs := s.parseParagraphs(to.Content, to.ContentType)

    // 计算差异
    diff := s.computeParagraphDiff(fromParagraphs, toParagraphs)

    return &DocumentDiffResult{
        FromVersion: from.DocumentVersionID,
        ToVersion:   to.DocumentVersionID,
        DiffData:     diff,
    }, nil
}

// parseParagraphs 解析内容为段落
func (s *SnapshotService) parseParagraphs(
    content string,
    contentType string,
) []Paragraph {
    switch contentType {
    case "tiptap":
        return s.parseTipTapParagraphs(content)
    case "markdown":
        return s.parseMarkdownParagraphs(content)
    case "richtext":
        return s.parseRichTextParagraphs(content)
    default:
        return s.parsePlainTextParagraphs(content)
    }
}

// computeParagraphDiff 计算段落差异
func (s *SnapshotService) computeParagraphDiff(
    from []Paragraph,
    to []Paragraph,
) *ParagraphDiffData {
    diff := &ParagraphDiffData{
        ParagraphsChanged: make([]ParagraphDiff, 0),
        TotalParagraphs:     len(to),
    }

    fromMap := make(map[int]*Paragraph)
    for i := range from {
        fromMap[from[i].Order] = &from[i]
    }

    toMap := make(map[int]*Paragraph)
    for i := range to {
        toMap[to[i].Order] = &to[i]
    }

    // 检查新增和修改
    for order, toPara := range toMap {
        fromPara, exists := fromMap[order]

        if !exists {
            // 新增段落
            diff.ParagraphsChanged = append(diff.ParagraphsChanged, ParagraphDiff{
                ParagraphOrder: order,
                Action:        DiffActionAdded,
                WordsChanged:  len([]rune(toPara.Content)),
            })
        } else if fromPara.ContentHash != toPara.ContentHash {
            // 内容修改
            diff.ParagraphsChanged = append(diff.ParagraphsChanged, ParagraphDiff{
                ParagraphOrder: order,
                Action:        DiffActionModified,
                WordsChanged:  len([]rune(toPara.Content)) - len([]rune(fromPara.Content)),
                ChangedRatio:  float64(len([]rune(toPara.Content))-len([]rune(fromPara.Content))) / float64(len([]rune(fromPara.Content))),
            })
        }
    }

    // 检查删除
    for order, fromPara := range fromMap {
        if _, exists := toMap[order]; !exists {
            diff.ParagraphsChanged = append(diff.ParagraphsChanged, ParagraphDiff{
                ParagraphOrder: order,
                Action:        DiffActionDeleted,
                WordsChanged:  -len([]rune(fromPara.Content)),
            })
        }
    }

    return diff
}

type Paragraph struct {
    Order      int     `json:"order"`
    Content    string  `json:"content"`
    ContentHash string  `json:"contentHash"`
}

type DocumentDiffResult struct {
    FromVersion   string             `json:"fromVersion"`
    ToVersion     string             `json:"toVersion"`
    DiffData       *ParagraphDiffData  `json:"diffData"`
}
```

---

## API 设计

### 版本相关API

```go
// api/v1/writer/version_api.go

// GetDocumentVersions 获取文档版本列表
// @Summary 获取文档版本列表
// @Tags Writer
// @Router /writer/documents/{documentId}/versions [get]
func (api *VersionAPI) GetDocumentVersions(c *gin.Context) {
    documentID := c.Param("documentId")
    userID := api.GetUserID(c)

    versions, err := api.documentVersionService.ListByDocumentID(
        c.Request.Context(),
        documentID,
        userID,
    )
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{
        "data":  versions,
        "total": len(versions),
    })
}

// GetDocumentVersionDiff 获取文档版本差异
// @Summary 获取文档版本差异
// @Tags Writer
// @Router /writer/documents/{documentId}/versions/diff [get]
func (api *VersionAPI) GetDocumentVersionDiff(c *gin.Context) {
    documentID := c.Param("documentId")
    fromVersion := c.Query("from")
    toVersion := c.Query("to")
    userID := api.GetUserID(c)

    diff, err := api.snapshotService.CompareDocumentSnapshots(
        c.Request.Context(),
        fromVersion,
        toVersion,
    )
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, diff)
}

// api/v1/reader/book_version_api.go

// GetBookVersions 获取书籍版本列表
// @Summary 获取书籍版本列表
// @Tags Reader
// @Router /reader/books/{bookId}/versions [get]
func (api *BookVersionAPI) GetBookVersions(c *gin.Context) {
    bookID := c.Param("bookId")
    userID := api.GetUserID(c)

    // 验证订阅...

    versions, err := api.bookVersionService.ListByBookID(
        c.Request.Context(),
        bookID,
        userID,
    )
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{
        "data":  versions,
        "total": len(versions),
    })
}

// GetBookVersionDiff 获取书籍版本差异
// @Summary 获取书籍版本差异
// @Tags Reader
// @Router /reader/books/{bookId}/versions/diff [get]
func (api *BookVersionAPI) GetBookVersionDiff(c *gin.Context) {
    bookID := c.Param("bookId")
    fromVersion := c.Query("from")
    toVersion := c.Query("to")
    userID := api.GetUserID(c)

    diff, err := api.snapshotService.CompareChapterSnapshots(
        c.Request.Context(),
        fromVersion,
        toVersion,
    )
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, diff)
}
```

---

## 实施计划

### Phase 1: 数据模型（2天）

- [ ] 创建 `DocumentVersion` 模型和Repository
- [ ] 创建 `DocumentContentSnapshot` 模型和Repository
- [ ] 创建 `BookVersion` 模型和Repository
- [ ] 创建 `ChapterContentSnapshot` 模型和Repository
- [ ] 数据库索引

### Phase 2: 核心服务（3天）

- [ ] 实现 `SnapshotService`（创建、获取、对比）
- [ ] 实现 `VersionService`（版本管理）
- [ ] 实现 `PublishService`（发布逻辑）
- [ ] 实现 `DiffService`（差异计算）

### Phase 3: API层（2天）

- [ ] Writer版本相关API
- [ ] Reader版本相关API
- [ ] 版本对比API

### Phase 4: 前端（3天）

- [ ] 版本列表展示
- [ ] 版本对比UI（段落级高亮）
- [ ] 版本切换功能

---

## 相关文档

- [发布-审核-订阅系统](./2026-03-06-publication-review-subsystem-design.md)
- [分层版本管理](./2026-03-06-layered-version-management-design.md)
- [Editor V2交接单](E:\Github\Qingyu\_wt_qy_backend_editor_v2\docs\architecture\2026-03-04-writer-editor-v2-handoff.md)

---

**设计完成时间**: 2026-03-06
**预计实施时间**: 10天
**建议执行者**: 后端团队 + 前端团队协同
