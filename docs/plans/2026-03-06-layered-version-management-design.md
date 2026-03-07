# 分层版本管理架构设计

**设计日期**: 2026-03-06
**设计者**: Kore
**优先级**: 🔴 P0
**问题来源**: 版本快照存储优化需求

---

## 问题分析

### 当前完整快照方案的问题

```
┌─────────────────────────────────────────────────────────────────────┐
│ 场景：10万字书籍，作者修改1个错别字                                  │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  BookVersion v1:                                                    │
│  - 存储大小: ~5MB (完整内容)                                         │
│                                                                     │
│  BookVersion v2:                                                    │
│  - 存储大小: ~5MB (完整内容)                                        │
│  - 实际差异: 几字节                                                 │
│                                                                     │
│  存储浪费: 99.99%                                                    │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

**问题**：
1. ❌ 存储空间浪费巨大
2. ❌ 版本越多，浪费越严重
3. ❌ 传输成本高（读者切换版本需要下载大量重复数据）

---

## 解决方案：分层版本管理

### 核心思想

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        存储策略分层                                     │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  书籍版本层 (BookVersion)          - 只存储元数据和变更摘要              │
│         ↓                                                          只存索引
│  章节版本层 (ChapterVersion)        - 存储章节内容，采用增量存储        │
│         ↓                                                          按需加载
│  内容块层 (ContentBlock) - 可选 - 超大章节的块级增量存储                 │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 存储策略

| 变更规模 | 存储策略 | 说明 |
|---------|---------|------|
| 小修改（<10%） | Diff存储 | 只存储差异 |
| 中等修改（10-50%） | Diff存储 | 存储差异 |
| 大修改（>50%） | Full存储 | 存储完整内容 |
| 新增章节 | Full存储 | 必须完整存储 |

---

## 数据模型设计

### 1. BookVersion（书籍版本）- 元数据层

```go
// models/bookstore/book_version.go

package bookstore

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "Qingyu_backend/models/shared"
)

// BookVersion 书籍版本（元数据层）
type BookVersion struct {
    shared.IdentifiedEntity `bson:",inline"`
    shared.BaseEntity       `bson:",inline"`

    // 关联信息
    BookID         primitive.ObjectID `bson:"book_id" json:"bookId"`
    PublicationID  string             `bson:"publication_id" json:"publicationId"`

    // 版本标识
    VersionNumber  int                `bson:"version_number" json:"versionNumber"`
    VersionName    string             `bson:"version_name" json:"versionName"` // v1.0, v2.0
    VersionType    VersionType        `bson:"version_type" json:"versionType"`

    // ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
    //  变更摘要（不存储实际内容）
    // ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
    ChangeSummary  *ChangeSummary     `bson:"change_summary" json:"changeSummary"`
    ChangedChapters []ChangedChapter `bson:"changed_chapters" json:"changedChapters"` // 变更的章节列表
    BaseVersionID  *string            `bson:"base_version_id,omitempty" json:"baseVersionId,omitempty"` // 基础版本

    // 发布信息
    PublishedAt    time.Time          `bson:"published_at" json:"publishedAt"`
    PublishedBy    string             `bson:"published_by" json:"publishedBy"`

    // 版本状态
    Status         VersionStatus      `bson:"status" json:"status"`
    IsLatest       bool               `bson:"is_latest" json:"isLatest"`

    // 统计信息（不依赖内容）
    TotalChapters  int                `bson:"total_chapters" json:"totalChapters"`
    TotalWords     int64              `bson:"total_words" json:"totalWords"`
    ViewCount      int64              `bson:"view_count" json:"viewCount"`

    // 元数据
    Description    string             `bson:"description,omitempty" json:"description,omitempty"`
    Changelog      string             `bson:"changelog,omitempty" json:"changelog,omitempty"` // 更新日志
    Tags           []string           `bson:"tags,omitempty" json:"tags,omitempty"`
}

// ChangedChapter 章节变更记录
type ChangedChapter struct {
    ChapterID      string         `bson:"chapter_id" json:"chapterId"`
    ChapterNumber  int            `bson:"chapter_number" json:"chapterNumber"`
    Title          string         `bson:"title" json:"title"`
    Action         ChapterAction  `bson:"action" json:"action"` // added, modified, deleted

    // 版本关联
    CurrentVersionID string       `bson:"current_version_id" json:"currentVersionId"` // 当前章节版本ID
    PreviousVersionID *string     `bson:"previous_version_id,omitempty" json:"previousVersionId,omitempty"` // 上一版本ID

    // 变更统计
    WordCountDiff  int            `bson:"word_count_diff,omitempty" json:"wordCountDiff,omitempty"` // 字数变化
    ChangeRatio    float64        `bson:"change_ratio,omitempty" json:"changeRatio,omitempty"` // 变更比例 0-1
}

// ChangeSummary 变更摘要
type ChangeSummary struct {
    AddedChapters     int              `bson:"added_chapters" json:"addedChapters"`
    ModifiedChapters  int              `bson:"modified_chapters" json:"modifiedChapters"`
    DeletedChapters   int              `bson:"deleted_chapters" json:"deletedChapters"`
    TotalWordsDiff    int64            `bson:"total_words_diff" json:"totalWordsDiff"`
    StorageSaved      int64            `bson:"storage_saved" json:"storageSaved"` // 节省的存储（字节）
    CompressionRatio  float64          `bson:"compression_ratio" json:"compressionRatio"` // 压缩比
}
```

### 2. ChapterVersion（章节版本）- 内容层

```go
// models/bookstore/chapter_version.go

package bookstore

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "Qingyu_backend/models/shared"
)

// StorageType 存储类型
type StorageType string

const (
    StorageTypeFull StorageType = "full" // 完整存储
    StorageTypeDiff StorageType = "diff" // 差异存储
)

// ChapterVersion 章节版本（内容层）
type ChapterVersion struct {
    shared.IdentifiedEntity `bson:",inline"`
    shared.BaseEntity       `bson:",inline"`

    // 关联信息
    ChapterID      primitive.ObjectID `bson:"chapter_id" json:"chapterId"` // 章节ID
    BookID         primitive.ObjectID `bson:"book_id" json:"bookId"`       // 书籍ID
    BookVersionID  string             `bson:"book_version_id" json:"bookVersionId"` // 所属书籍版本

    // 版本信息
    Version        int                `bson:"version" json:"version"` // 章节版本号
    ChapterNumber  int                `bson:"chapter_number" json:"chapterNumber"`
    Title          string             `bson:"title" json:"title"`

    // ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
    //  内容存储（根据StorageType决定）
    // ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
    StorageType    StorageType        `bson:"storage_type" json:"storageType"`

    // 完整内容（StorageType=Full）
    FullContent    *string            `bson:"full_content,omitempty" json:"fullContent,omitempty"`

    // 差异内容（StorageType=Diff）
    DiffContent    *ChapterDiff       `bson:"diff_content,omitempty" json:"diffContent,omitempty"`
    BaseVersionID  *string            `bson:"base_version_id,omitempty" json:"baseVersionId,omitempty"` // 基础版本ID

    // 统计信息（两种存储类型都有）
    WordCount      int                `bson:"word_count" json:"wordCount"`
    ContentHash    string             `bson:"content_hash" json:"contentHash"` // 当前内容哈希
    BaseContentHash *string           `bson:"base_content_hash,omitempty" json:"baseContentHash,omitempty"` // 基础内容哈希

    // 元数据
    Price          float64            `bson:"price,omitempty" json:"price,omitempty"`
    IsFree         bool               `bson:"is_free,omitempty" json:"isFree,omitempty"`
    PublishedAt    time.Time          `bson:"published_at" json:"publishedAt"`
}

// ChapterDiff 章节差异
type ChapterDiff struct {
    // 文本差异（行级）
    LineDiff       *TextLineDiff      `bson:"line_diff,omitempty" json:"lineDiff,omitempty"`

    // 或者字符级差异（用于小修改）
    CharDiff       *TextCharDiff      `bson:"char_diff,omitempty" json:"charDiff,omitempty"`

    // 差异统计
    AddedChars     int                `bson:"added_chars" json:"addedChars"`
    DeletedChars   int                `bson:"deleted_chars" json:"deletedChars"`
    ChangedRatio   float64            `bson:"changed_ratio" json:"changedRatio"` // 0-1

    // 差异位置（用于快速定位）
    DiffPositions  []DiffPosition     `bson:"diff_positions,omitempty" json:"diffPositions,omitempty"`
}

// TextLineDiff 行级差异（使用Myers算法）
type TextLineDiff struct {
    OldLines       []Line             `bson:"old_lines,omitempty" json:"oldLines,omitempty"` // 被删除的行
    NewLines       []Line             `bson:"new_lines,omitempty" json:"new_lines,omitempty"` // 新增的行
    Hunks          []DiffHunk         `bson:"hunks" json:"hunks"` // 差异块
}

// TextCharDiff 字符级差异（用于精确跟踪）
type TextCharDiff struct {
    OldText        string             `bson:"old_text,omitempty" json:"oldText,omitempty"` // 被删除的文本
    NewText        string             `bson:"new_text,omitempty" json:"newText,omitempty"` // 新增的文本
    Operations     []DiffOp           `bson:"operations" json:"operations"` // 操作序列
}

// Line 文本行
type Line struct {
    Number         int                `bson:"number" json:"number"`
    Content        string             `bson:"content" json:"content"`
}

// DiffHunk 差异块
type DiffHunk struct {
    OldStart       int                `bson:"old_start" json:"oldStart"` // 旧文件起始行
    OldCount       int                `bson:"old_count" json:"oldCount"` // 旧文件行数
    NewStart       int                `bson:"new_start" json:"newStart"` // 新文件起始行
    NewCount       int                `bson:"new_count" json:"newCount"` // 新文件行数
    Lines          []Line             `bson:"lines" json:"lines"` // 差异行
}

// DiffOp 差异操作
type DiffOp struct {
    Type           DiffOpType         `bson:"type" json:"type"` // equal, insert, delete
    Position       int                `bson:"position" json:"position"`
    Text           string             `bson:"text,omitempty" json:"text,omitempty"`
}

type DiffOpType string

const (
    DiffOpEqual   DiffOpType = "equal"
    DiffOpInsert  DiffOpType = "insert"
    DiffOpDelete  DiffOpType = "delete"
)

// DiffPosition 差异位置
type DiffPosition struct {
    Type           string             `bson:"type" json:"type"` // paragraph, sentence, word
    Index          int                `bson:"index" json:"index"`
    Offset         int                `bson:"offset" json:"offset"`
    Length         int                `bson:"length" json:"length"`
}
```

### 3. ContentBlock（内容块）- 可选的大文件分块层

```go
// models/bookstore/content_block.go

package bookstore

// ContentBlock 内容块（可选，用于超大章节）
type ContentBlock struct {
    shared.IdentifiedEntity `bson:",inline"`

    ChapterVersionID string             `bson:"chapter_version_id" json:"chapterVersionId"`
    BlockIndex      int                `bson:"block_index" json:"blockIndex"` // 块索引
    BlockType       BlockType          `bson:"block_type" json:"blockType"` // full, diff

    // 块内容
    Content         string             `bson:"content" json:"content"`
    ContentLength   int                `bson:"content_length" json:"contentLength"`

    // 压缩信息
    Compressed      bool               `bson:"compressed" json:"compressed"`
    OriginalSize    int                `bson:"original_size" json:"originalSize"`
    CompressedSize  int                `bson:"compressed_size" json:"compressedSize"`

    // 哈希校验
    ContentHash     string             `bson:"content_hash" json:"contentHash"`
}

type BlockType string

const (
    BlockTypeFull     BlockType = "full" // 完整块
    BlockTypeDiff     BlockType = "diff" // 差异块
)
```

---

## 版本创建策略

### 存储类型选择算法

```go
// service/bookstore/version_strategy.go

type VersionStrategy struct {
    fullThreshold      float64 // >50%变化使用Full
    diffThreshold      float64 // <=50%变化使用Diff
    smallDiffThreshold int     // <100字符使用CharDiff
}

// DecideStorageType 决定存储类型
func (s *VersionStrategy) DecideStorageType(
    oldContent *string,
    newContent string,
) (StorageType, *ChangeStats) {
    // 首次发布：必须Full
    if oldContent == nil || *oldContent == "" {
        return StorageTypeFull, &ChangeStats{
            ChangeRatio: 1.0,
        }
    }

    // 计算差异
    stats := s.calculateStats(*oldContent, newContent)

    // 决策
    if stats.ChangeRatio > s.fullThreshold {
        return StorageTypeFull, stats
    }

    return StorageTypeDiff, stats
}

// ChangeStats 变更统计
type ChangeStats struct {
    AddedChars      int     `json:"addedChars"`
    DeletedChars    int     `json:"deletedChars"`
    ChangedChars    int     `json:"changedChars"`
    TotalChars      int     `json:"totalChars"`
    ChangeRatio     float64 `json:"changeRatio"` // 0-1
    LinesChanged    int     `json:"linesChanged"`
}

// calculateStats 计算变更统计
func (s *VersionStrategy) calculateStats(old, new string) *ChangeStats {
    // 使用 Myers Diff 算法
    diff := s.computeMyersDiff(old, new)

    stats := &ChangeStats{
        TotalChars:   len(new),
    }

    for _, hunk := range diff.Hunks {
        stats.DeletedChars += len(hunk.OldText)
        stats.AddedChars += len(hunk.NewText)
        stats.LinesChanged++
    }

    // 计算变更比例
    if stats.TotalChars > 0 {
        stats.ChangedChars = stats.AddedChars + stats.DeletedChars
        stats.ChangeRatio = float64(stats.ChangedChars) / float64(stats.TotalChars*2)
    }

    return stats
}
```

---

## 版本恢复服务

### 完整内容重建

```go
// service/bookstore/version_reconstruction.go

type VersionReconstructionService struct {
    chapterVersionRepo interfaces.ChapterVersionRepository
    contentBlockRepo   interfaces.ContentBlockRepository
    cache              interfaces.CacheService
}

// GetChapterContent 获取章节完整内容
func (s *VersionReconstructionService) GetChapterContent(
    ctx context.Context,
    chapterVersionID string,
) (string, error) {
    // 1. 尝试从缓存获取
    if cached, found := s.cache.Get(ctx, "chapter:"+chapterVersionID); found {
        return cached.(string), nil
    }

    // 2. 获取章节版本
    cv, err := s.chapterVersionRepo.GetByID(ctx, chapterVersionID)
    if err != nil {
        return "", err
    }

    var content string

    // 3. 根据存储类型获取内容
    switch cv.StorageType {
    case StorageTypeFull:
        content = *cv.FullContent

    case StorageTypeDiff:
        // 递归重建完整内容
        content, err = s.rebuildFromDiff(ctx, cv)
        if err != nil {
            return "", err
        }
    }

    // 4. 缓存完整内容
    s.cache.Set(ctx, "chapter:"+chapterVersionID, content, 1*time.Hour)

    return content, nil
}

// rebuildFromDiff 从差异重建完整内容
func (s *VersionReconstructionService) rebuildFromDiff(
    ctx context.Context,
    cv *ChapterVersion,
) (string, error) {
    // 获取基础版本
    if cv.BaseVersionID == nil {
        return "", fmt.Errorf("差异存储缺少基础版本ID")
    }

    baseCV, err := s.chapterVersionRepo.GetByID(ctx, *cv.BaseVersionID)
    if err != nil {
        return "", fmt.Errorf("基础版本不存在: %w", err)
    }

    // 递归获取基础版本内容
    baseContent, err := s.GetChapterContent(ctx, baseCV.ID.Hex())
    if err != nil {
        return "", err
    }

    // 应用差异
    newContent, err := s.applyDiff(baseContent, cv.DiffContent)
    if err != nil {
        return "", err
    }

    return newContent, nil
}

// applyDiff 应用差异
func (s *VersionReconstructionService) applyDiff(
    baseContent string,
    diff *ChapterDiff,
) (string, error) {
    if diff.LineDiff != nil {
        return s.applyLineDiff(baseContent, diff.LineDiff)
    }

    if diff.CharDiff != nil {
        return s.applyCharDiff(baseContent, diff.CharDiff)
    }

    return "", fmt.Errorf("无效的差异内容")
}

// applyLineDiff 应用行级差异
func (s *VersionReconstructionService) applyLineDiff(
    baseContent string,
    diff *TextLineDiff,
) (string, error) {
    lines := strings.Split(baseContent, "\n")
    result := make([]string, 0, len(lines))

    oldLineIdx := 0
    for _, hunk := range diff.Hunks {
        // 添加差异块之前的行
        for oldLineIdx < hunk.OldStart-1 {
            result = append(result, lines[oldLineIdx])
            oldLineIdx++
        }

        // 跳过被删除的行
        oldLineIdx += hunk.OldCount

        // 添加差异块中的行
        for _, line := range hunk.Lines {
            result = append(result, line.Content)
        }
    }

    // 添加剩余行
    for oldLineIdx < len(lines) {
        result = append(result, lines[oldLineIdx])
        oldLineIdx++
    }

    return strings.Join(result, "\n"), nil
}

// applyCharDiff 应用字符级差异
func (s *VersionReconstructionService) applyCharDiff(
    baseContent string,
    diff *TextCharDiff,
) (string, error) {
    result := baseContent
    offset := 0

    for _, op := range diff.Operations {
        switch op.Type {
        case DiffOpDelete:
            // 删除字符
            result = result[:offset] + result[offset+op.Length:]
        case DiffOpInsert:
            // 插入字符
            result = result[:offset] + op.Text + result[offset:]
            offset += len(op.Text)
        case DiffOpEqual:
            offset += op.Position
        }
    }

    return result, nil
}
```

---

## 版本对比服务

### 差异计算

```go
// service/bookstore/diff_service.go

import (
    "github.com/sergi/go-diff/diffmatchpatch"
)

type DiffService struct {
    dmp *diffmatchpatch.DiffMatchPatch
}

// ComputeChapterDiff 计算章节差异
func (s *DiffService) ComputeChapterDiff(
    oldContent string,
    newContent string,
    strategy DiffStrategy,
) (*ChapterDiff, error) {
    switch strategy {
    case DiffStrategyLine:
        return s.computeLineDiff(oldContent, newContent), nil

    case DiffStrategyChar:
        return s.computeCharDiff(oldContent, newContent), nil

    case DiffStrategyWord:
        return s.computeWordDiff(oldContent, newContent), nil

    default:
        return s.computeAutoDiff(oldContent, newContent), nil
    }
}

// DiffStrategy 差异策略
type DiffStrategy string

const (
    DiffStrategyLine  DiffStrategy = "line"  // 行级
    DiffStrategyChar  DiffStrategy = "char"  // 字符级
    DiffStrategyWord  DiffStrategy = "word"  // 单词级
    DiffStrategyAuto  DiffStrategy = "auto"  // 自动选择
)

// computeLineDiff 计算行级差异
func (s *DiffService) computeLineDiff(old, new string) *ChapterDiff {
    diffs := s.dmp.DiffMain(old, new, false)

    lineDiff := &TextLineDiff{
        Hunks: make([]DiffHunk, 0),
    }

    oldLineNum := 1
    newLineNum := 1

    for _, diff := range diffs {
        lines := strings.Split(diff.Text, "\n")

        for _, line := range lines {
            switch diff.Type {
            case diffmatchpatch.DiffDelete:
                lineDiff.OldLines = append(lineDiff.OldLines, Line{
                    Number:  oldLineNum,
                    Content: line,
                })
                oldLineNum++

            case diffmatchpatch.DiffInsert:
                lineDiff.NewLines = append(lineDiff.NewLines, Line{
                    Number:  newLineNum,
                    Content: line,
                })
                newLineNum++

            case diffmatchpatch.DiffEqual:
                // 相同行不需要记录在差异中
                oldLineNum++
                newLineNum++
            }
        }
    }

    return &ChapterDiff{
        LineDiff:    lineDiff,
        AddedChars:  s.countChars(lineDiff.NewLines),
        DeletedChars: s.countChars(lineDiff.OldLines),
    }
}

// computeCharDiff 计算字符级差异
func (s *DiffService) computeCharDiff(old, new string) *ChapterDiff {
    diffs := s.dmp.DiffMain(old, new, true)

    charDiff := &TextCharDiff{
        Operations: make([]DiffOp, 0),
    }

    position := 0
    for _, diff := range diffs {
        for _, r := range diff.Text {
            op := DiffOp{
                Position: position,
            }

            switch diff.Type {
            case diffmatchpatch.DiffDelete:
                op.Type = DiffOpDelete
                op.Text = string(r)
                position++

            case diffmatchpatch.DiffInsert:
                op.Type = DiffOpInsert
                op.Text = string(r)
                position++

            case diffmatchpatch.DiffEqual:
                op.Type = DiffOpEqual
                op.Text = string(r)
                position++
            }

            charDiff.Operations = append(charDiff.Operations, op)
        }
    }

    return &ChapterDiff{
        CharDiff: charDiff,
    }
}

// computeAutoDiff 自动选择差异策略
func (s *DiffService) computeAutoDiff(old, new string) *ChapterDiff {
    // 小修改（<10%）：字符级
    // 中等修改（10-50%）：行级
    // 大修改（>50%）：返回完整内容

    changeRatio := float64(len(new)-len(old)) / float64(len(old))
    if changeRatio < 0 {
        changeRatio = -changeRatio
    }

    if changeRatio < 0.1 {
        return s.computeCharDiff(old, new)
    }

    return s.computeLineDiff(old, new)
}

// countChars 统计行中的字符数
func (s *DiffService) countChars(lines []Line) int {
    count := 0
    for _, line := range lines {
        count += len(line.Content)
    }
    return count
}
```

---

## API 设计

### 获取章节内容

```go
// api/v1/reader/chapter_api.go

// GetChapterContent 获取章节内容
// @Summary 获取章节内容
// @Description 获取指定版本的章节完整内容（自动处理Diff恢复）
// @Tags Reader
// @Router /reader/books/{bookId}/chapters/{chapterId}/versions/{versionId} [get]
func (api *ChapterAPI) GetChapterContent(c *gin.Context) {
    bookID := c.Param("bookId")
    chapterID := c.Param("chapterId")
    versionID := c.Param("versionId")
    userID := api.GetUserID(c)

    // 验证订阅
    subscription, err := api.subscriptionService.GetByUserAndBook(
        c.Request.Context(),
        userID,
        bookID,
    )
    if err != nil {
        c.JSON(403, gin.H{"error": "未订阅该书籍"})
        return
    }

    // 获取章节版本
    chapterVersion, err := api.chapterVersionRepo.GetByVersionID(
        c.Request.Context(),
        chapterID,
        versionID,
    )
    if err != nil {
        c.JSON(404, gin.H{"error": "章节版本不存在"})
        return
    }

    // 检查版本权限
    if !api.canAccessVersion(subscription, versionID) {
        c.JSON(403, gin.H{"error": "无权访问该版本"})
        return
    }

    // 获取完整内容（自动处理Diff恢复）
    content, err := api.reconstructionService.GetChapterContent(
        c.Request.Context(),
        chapterVersion.ID.Hex(),
    )
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{
        "chapterId":   chapterID,
        "versionId":   versionID,
        "title":       chapterVersion.Title,
        "content":     content,
        "wordCount":   chapterVersion.WordCount,
        "storageType": chapterVersion.StorageType,
    })
}

// GetChapterDiff 获取章节差异
// @Summary 获取章节差异
// @Description 对比两个版本之间的差异
// @Tags Reader
// @Router /reader/books/{bookId}/chapters/{chapterId}/diff [get]
func (api *ChapterAPI) GetChapterDiff(c *gin.Context) {
    chapterID := c.Param("chapterId")
    fromVersion := c.Query("from")
    toVersion := c.Query("to")
    userID := api.GetUserID(c)

    // 验证订阅...

    // 获取两个版本
    fromCV, _ := api.chapterVersionRepo.GetByVersionID(c.Request.Context(), chapterID, fromVersion)
    toCV, _ := api.chapterVersionRepo.GetByVersionID(c.Request.Context(), chapterID, toVersion)

    // 获取完整内容
    fromContent, _ := api.reconstructionService.GetChapterContent(
        c.Request.Context(),
        fromCV.ID.Hex(),
    )
    toContent, _ := api.reconstructionService.GetChapterContent(
        c.Request.Context(),
        toCV.ID.Hex(),
    )

    // 计算差异
    diff, _ := api.diffService.ComputeChapterDiff(
        fromContent,
        toContent,
        DiffStrategyAuto,
    )

    c.JSON(200, gin.H{
        "fromVersion": fromVersion,
        "toVersion":   toVersion,
        "diff":        diff,
    })
}
```

---

## 存储优化效果

### 对比示例

```
┌─────────────────────────────────────────────────────────────────────────┐
│  场景：100章书籍，每章5000字，修改1章（改动100字）                        │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  【完整快照方案】                                                        │
│  BookVersion v1: 100章 x 5000字 = ~500KB                              │
│  BookVersion v2: 100章 x 5000字 = ~500KB                              │
│  ─────────────────────────────────────────                             │
│  总存储: 1MB                                                            │
│  实际差异: ~100字                                                       │
│  浪费率: 99.99%                                                         │
│                                                                         │
│  【分层版本方案】                                                        │
│  BookVersion v1: 元数据 ~5KB                                           │
│    ├─ ChapterVersion ch_001_v1 (Full): ~5KB                          │
│    ├─ ChapterVersion ch_002_v1 (Full): ~5KB                          │
│    └─ ... 100个章节                                                    │
│  ─────────────────────────────────────────                             │
│  BookVersion v2: 元数据 ~5KB + 变更摘要 ~1KB                         │
│    ├─ ChapterVersion ch_001_v2 (Diff): ~200 bytes  ← 只有差异        │
│    └─ 其他章节指向 v1 版本                                             │
│  ─────────────────────────────────────────                             │
│  总存储: 500KB (v1) + 6KB (v2) = 506KB                               │
│  节省: 494KB (49.4%)                                                   │
│                                                                         │
│  【修改50章后】                                                          │
│  BookVersion v3: 元数据 ~5KB + 50个Diff引用 ~10KB                    │
│  总存储: 506KB + 15KB = 521KB                                         │
│  完整快照方案: 1MB x 3 = 3MB                                            │
│  节省: 2.48MB (82.7%)                                                  │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 实施建议

### Phase 1: 核心功能（3天）

1. **创建ChapterVersion模型**
2. **实现Diff算法**（使用go-diff库）
3. **实现版本重建服务**
4. **实现存储策略决策**

### Phase 2: 优化功能（2天）

1. **添加缓存层**（重建的内容缓存）
2. **实现ContentBlock分块**（超大章节）
3. **性能测试和优化**

### Phase 3: API和前端（2天）

1. **章节内容API**
2. **版本对比API**
3. **前端差异展示**

---

## 相关文档

- [发布-审核-订阅系统](./2026-03-06-publication-review-subsystem-design.md)
- [Project-Book分离架构](./2026-03-06-project-book-separation-architecture.md)

---

**设计完成时间**: 2026-03-06
**预计实施时间**: 7天
**建议执行者**: 后端团队
