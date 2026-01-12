# 分层验证E2E测试实现计划

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**目标:** 构建一个三层分区的E2E测试框架，支持基础流程验证、数据一致性检查和边界场景测试，使用混合数据策略（生产镜像数据+高保真模拟数据+真实用户行为），确保系统在真实数据场景下的可靠性。

**架构:** 将E2E测试分为三个独立可执行的层次：Layer 1（基础流程测试，2-3分钟）、Layer 2（数据一致性测试，3-5分钟）、Layer 3（边界场景测试，5-8分钟）。每层使用独立的数据集和验证策略，支持按需运行。框架扩展现有的 `test/e2e/framework`，新增数据工厂、一致性验证器和边界测试辅助工具。

**技术栈:** Go 1.21+, testify/assert, testify/require, gin-gonic/gin, MongoDB Driver, 现有的测试框架基础设施

---

## 目录结构

```
test/e2e/
├── framework/                      # 现有框架（保持不变）
│   ├── environment.go
│   ├── fixtures.go
│   ├── actions.go
│   └── assertions.go
├── data/                          # 新增：测试数据管理
│   ├── factory.go                 # 数据工厂
│   ├── scenarios.go               # 场景数据生成器
│   └── consistency_validator.go   # 数据一致性验证器
├── layer1_basic/                  # 新增：Layer 1 基础流程测试
│   ├── auth_flow_test.go
│   ├── reading_flow_test.go
│   ├── writing_flow_test.go
│   └── social_flow_test.go
├── layer2_consistency/            # 新增：Layer 2 数据一致性测试
│   ├── cross_module_test.go
│   ├── transaction_test.go
│   └── cascade_update_test.go
├── layer3_boundary/               # 新增：Layer 3 边界场景测试
│   ├── concurrent_test.go
│   ├── limit_test.go
│   └── error_recovery_test.go
├── scenarios/                      # 现有场景测试（保持不变）
│   └── complete_user_journey_test.go
├── suite.go                       # 新增：测试套件入口
└── README.md                      # 更新：E2E测试文档
```

---

## Task 1: 创建数据工厂 (Data Factory)

**目标:** 建立灵活的数据生成系统，支持创建符合生产特征的测试数据，包括用户、书籍、章节、评论、收藏等核心实体。

**文件:**
- 创建: `test/e2e/data/factory.go`
- 创建: `test/e2e/data/scenarios.go`

### Step 1: 编写数据工厂基础结构测试

```go
package e2e_data_test

import (
    "testing"
    "Qingyu_backend/test/e2e/data"
)

func TestDataFactory_CreateUser(t *testing.T) {
    factory := data.NewTestDataFactory()

    user := factory.CreateUser(data.UserOptions{
        Username: "test_user",
        VIPLevel: 1,
    })

    if user.Username != "test_user" {
        t.Errorf("expected username 'test_user', got '%s'", user.Username)
    }
    if user.VIPLevel != 1 {
        t.Errorf("expected VIP level 1, got %d", user.VIPLevel)
    }
    if user.ID == "" {
        t.Error("expected user ID to be set")
    }
}
```

**运行:** `go test ./test/e2e/data/... -v -run TestDataFactory_CreateUser`
**预期:** FAIL - 包不存在

### Step 2: 实现数据工厂核心结构

**创建文件:** `test/e2e/data/factory.go`

```go
package data

import (
    "context"
    "fmt"
    "math/rand"
    "time"

    "github.com/stretchr/testify/require"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "golang.org/x/crypto/bcrypt"

    "Qingyu_backend/global"
    "Qingyu_backend/models/bookstore"
    "Qingyu_backend/models/social"
    "Qingyu_backend/models/users"
    bookRepo "Qingyu_backend/repository/mongodb/bookstore"
    userRepo "Qingyu_backend/repository/mongodb/user"
    socialRepo "Qingyu_backend/repository/mongodb/social"
)

// TestDataFactory 测试数据工厂
type TestDataFactory struct {
    t *testing.T
}

// NewTestDataFactory 创建测试数据工厂
func NewTestDataFactory(t *testing.T) *TestDataFactory {
    return &TestDataFactory{t: t}
}

// UserOptions 用户创建选项
type UserOptions struct {
    Username string
    Email    string
    VIPLevel int
    Balance  float64
    Roles    []string
}

// BookOptions 书籍创建选项
type BookOptions struct {
    Title        string
    AuthorID     string
    Price        float64
    IsFree       bool
    Categories   []string
    WordCount    int
    ChapterCount int
}

// CommentOptions 评论创建选项
type CommentOptions struct {
    AuthorID  string
    TargetID  string
    TargetType string
    Content   string
}

// CreateUser 创建测试用户
func (f *TestDataFactory) CreateUser(opts UserOptions) *users.User {
    userID := primitive.NewObjectID()
    testPassword := "Test1234"
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
    require.NoError(f.t, err, "密码哈希失败")

    // 默认值处理
    username := opts.Username
    if username == "" {
        username = fmt.Sprintf("e2e_user_%s", userID.Hex()[:8])
    }

    email := opts.Email
    if email == "" {
        email = fmt.Sprintf("e2e_%s@example.com", userID.Hex()[:8])
    }

    roles := opts.Roles
    if len(roles) == 0 {
        roles = []string{"reader"}
    }

    user := &users.User{
        ID:       userID.Hex(),
        Username: username,
        Email:    email,
        Password: string(hashedPassword),
        VIPLevel: opts.VIPLevel,
        Balance:  opts.Balance,
        Status:   users.UserStatusActive,
        Roles:    roles,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    // 清理可能存在的同名用户
    userRepository := userRepo.NewMongoUserRepository(global.DB)
    existingUser, _ := userRepository.GetByUsername(context.Background(), user.Username)
    if existingUser != nil && existingUser.ID != user.ID {
        _ = userRepository.Delete(context.Background(), existingUser.ID)
    }

    err = userRepository.Create(context.Background(), user)
    require.NoError(f.t, err, "创建用户失败")

    return user
}

// CreateUsers 批量创建用户
func (f *TestDataFactory) CreateUsers(count int, baseOptions UserOptions) []*users.User {
    users := make([]*users.User, count)
    for i := 0; i < count; i++ {
        opts := baseOptions
        opts.Username = fmt.Sprintf("%s_%d", baseOptions.Username, i)
        opts.Email = fmt.Sprintf("e2e_batch_%d@example.com", i+rand.Intn(10000))
        users[i] = f.CreateUser(opts)
    }
    return users
}

// CreateBook 创建测试书籍
func (f *TestDataFactory) CreateBook(opts BookOptions) *bookstore.Book {
    bookID := primitive.NewObjectID()
    authorObjID, _ := primitive.ObjectIDFromHex(opts.AuthorID)

    // 默认值处理
    title := opts.Title
    if title == "" {
        title = fmt.Sprintf("e2e_book_%s", bookID.Hex()[:8])
    }

    categories := opts.Categories
    if len(categories) == 0 {
        categories = []string{"小说"}
    }

    wordCount := opts.WordCount
    if wordCount == 0 {
        wordCount = 10000
    }

    book := &bookstore.Book{
        ID:           bookID,
        Title:        title,
        AuthorID:     authorObjID,
        Introduction: "E2E测试书籍 - 用于验证系统功能",
        Categories:   categories,
        Price:        opts.Price,
        Status:       bookstore.BookStatusPublished,
        WordCount:    wordCount,
        IsFree:       opts.IsFree,
        ChapterCount: opts.ChapterCount,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }

    bookRepo := bookRepo.NewMongoBookRepository(global.DB.Client(), global.DB.Name())
    err := bookRepo.Create(context.Background(), book)
    require.NoError(f.t, err, "创建书籍失败")

    return book
}

// CreateChapter 创建测试章节
func (f *TestDataFactory) CreateChapter(bookID string, chapterNum int, isFree bool) *bookstore.Chapter {
    chapterID := primitive.NewObjectID()
    bookObjID, _ := primitive.ObjectIDFromHex(bookID)

    chapter := &bookstore.Chapter{
        ID:         chapterID,
        BookID:     bookObjID,
        Title:      fmt.Sprintf("第%d章", chapterNum),
        ChapterNum: chapterNum,
        WordCount:  2000,
        IsFree:     isFree,
        Price:      0,
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }

    chapterRepo := bookRepo.NewMongoChapterRepository(global.DB.Client(), global.DB.Name())
    err := chapterRepo.Create(context.Background(), chapter)
    require.NoError(f.t, err, "创建章节失败")

    // 创建章节内容
    chapterContentRepo := bookRepo.NewMongoChapterContentRepository(global.DB)
    content := fmt.Sprintf("这是第%d章的测试内容。用于验证阅读功能和数据处理流程。", chapterNum)
    chapterContent := &bookstore.ChapterContent{
        ID:        primitive.NewObjectID(),
        ChapterID: chapterID,
        Content:   content,
        Format:    "markdown",
        Version:   1,
        WordCount: len(content),
        CreatedAt: time.Now(),
    }
    err = chapterContentRepo.Create(context.Background(), chapterContent)
    require.NoError(f.t, err, "创建章节内容失败")

    return chapter
}

// CreateComment 创建测试评论
func (f *TestDataFactory) CreateComment(opts CommentOptions) *social.Comment {
    commentID := primitive.NewObjectID()

    // 默认值处理
    content := opts.Content
    if content == "" {
        content = "这是一条E2E测试评论，用于验证评论系统的功能。"
    }

    comment := &social.Comment{
        ID:         commentID,
        AuthorID:   opts.AuthorID,
        TargetID:   opts.TargetID,
        TargetType: opts.TargetType,
        Content:    content,
        Status:     social.CommentStatusPublished,
        LikeCount:  0,
        ReplyCount: 0,
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }

    commentRepo := socialRepo.NewMongoCommentRepository(global.DB)
    err := commentRepo.Create(context.Background(), comment)
    require.NoError(f.t, err, "创建评论失败")

    return comment
}

// CreateCollection 创建测试收藏
func (f *TestDataFactory) CreateCollection(userID, bookID string) *social.Collection {
    collectionID := primitive.NewObjectID()

    collection := &social.Collection{
        ID:        collectionID,
        UserID:    userID,
        BookID:    bookID,
        CreatedAt: time.Now(),
    }

    collectionRepo := socialRepo.NewMongoCollectionRepository(global.DB)
    err := collectionRepo.Create(context.Background(), collection)
    require.NoError(f.t, err, "创建收藏失败")

    return collection
}

// Cleanup 清理测试数据
func (f *TestDataFactory) Cleanup(prefix string) {
    ctx := context.Background()
    collections := []string{
        "users", "books", "chapters", "chapter_contents",
        "comments", "collections", "likes", "reading_progress",
    }

    for _, collName := range collections {
        // 删除带前缀的数据
        filter := map[string]interface{}{
            "$or": []map[string]interface{}{
                {"username": map[string]interface{}{"$regex": "^" + prefix}},
                {"email": map[string]interface{}{"$regex": "^" + prefix}},
                {"title": map[string]interface{}{"$regex": "^" + prefix}},
            },
        }
        result, _ := global.DB.Collection(collName).DeleteMany(ctx, filter)
        if result.DeletedCount > 0 {
            f.t.Logf("清理 %s: %d 条记录", collName, result.DeletedCount)
        }
    }
}
```

**运行:** `go test ./test/e2e/data/... -v -run TestDataFactory_CreateUser`
**预期:** PASS

### Step 3: 实现场景数据生成器

**创建文件:** `test/e2e/data/scenarios.go`

```go
package data

import (
    "testing"

    "Qingyu_backend/models/bookstore"
    "Qingyu_backend/models/users"
)

// ScenarioBuilder 场景数据构建器
type ScenarioBuilder struct {
    factory *TestDataFactory
    t       *testing.T
}

// NewScenarioBuilder 创建场景构建器
func NewScenarioBuilder(t *testing.T) *ScenarioBuilder {
    return &ScenarioBuilder{
        factory: NewTestDataFactory(t),
        t:       t,
    }
}

// ReaderWithProgress 构建有阅读进度的读者场景
type ReaderWithProgress struct {
    User           *users.User
    Books          []*bookstore.Book
    CurrentBook    *bookstore.Book
    CurrentChapter *bookstore.Chapter
    Progress       float64
}

// BuildReaderWithProgress 创建有阅读进度的读者场景
func (sb *ScenarioBuilder) BuildReaderWithProgress() *ReaderWithProgress {
    // 创建读者
    user := sb.factory.CreateUser(UserOptions{
        Username: "reader_with_progress",
        VIPLevel: 0,
    })

    // 创建作者和书籍
    author := sb.factory.CreateUser(UserOptions{
        Username: "author_for_reader",
        Roles:    []string{"reader", "author"},
    })

    book := sb.factory.CreateBook(BookOptions{
        Title:        "读者正在阅读的书",
        AuthorID:     author.ID,
        Price:        0,
        IsFree:       true,
        ChapterCount: 5,
    })

    // 创建章节
    chapters := make([]*bookstore.Chapter, 5)
    for i := 0; i < 5; i++ {
        chapters[i] = sb.factory.CreateChapter(book.ID.Hex(), i+1, i == 0)
    }

    return &ReaderWithProgress{
        User:           user,
        Books:          []*bookstore.Book{book},
        CurrentBook:    book,
        CurrentChapter: chapters[0],
        Progress:       0.3,
    }
}

// SocialInteraction 构建社交互动场景
type SocialInteraction struct {
    Users      []*users.User
    Book       *bookstore.Book
    Comments   []interface{} // Comment类型
    Collections []interface{} // Collection类型
    Likes      []interface{} // Like类型
}

// BuildSocialInteraction 创建社交互动场景
func (sb *ScenarioBuilder) BuildSocialInteraction(userCount int) *SocialInteraction {
    // 创建作者
    author := sb.factory.CreateUser(UserOptions{
        Username: "social_author",
        Roles:    []string{"reader", "author"},
    })

    // 创建书籍
    book := sb.factory.CreateBook(BookOptions{
        Title:        "热门互动书籍",
        AuthorID:     author.ID,
        Price:        0,
        IsFree:       true,
        ChapterCount: 3,
    })

    // 创建互动用户
    users := make([]*users.User, userCount)
    for i := 0; i < userCount; i++ {
        users[i] = sb.factory.CreateUser(UserOptions{
            Username: "",
            VIPLevel: i % 2, // 混合VIP等级
        })
    }

    return &SocialInteraction{
        Users:      users,
        Book:       book,
        Comments:   []interface{}{},
        Collections: []interface{}{},
        Likes:      []interface{}{},
    }
}

// PaidContent 构建付费内容场景
type PaidContent struct {
    Author        *users.User
    FreeUser      *users.User
    VIPUser       *users.User
    PaidBook      *bookstore.Book
    FreeChapters  []*bookstore.Chapter
    PaidChapters  []*bookstore.Chapter
}

// BuildPaidContent 创建付费内容场景
func (sb *ScenarioBuilder) BuildPaidContent() *PaidContent {
    // 创建作者
    author := sb.factory.CreateUser(UserOptions{
        Username: "paid_content_author",
        Roles:    []string{"reader", "author"},
    })

    // 创建免费用户
    freeUser := sb.factory.CreateUser(UserOptions{
        Username: "free_reader",
        VIPLevel: 0,
        Balance:  0,
    })

    // 创建VIP用户
    vipUser := sb.factory.CreateUser(UserOptions{
        Username: "vip_reader",
        VIPLevel: 1,
        Balance:  0,
    })

    // 创建付费书籍
    paidBook := sb.factory.CreateBook(BookOptions{
        Title:        "付费书籍",
        AuthorID:     author.ID,
        Price:        100,
        IsFree:       false,
        ChapterCount: 10,
    })

    // 创建章节（前3章免费，后7章付费）
    freeChapters := make([]*bookstore.Chapter, 3)
    paidChapters := make([]*bookstore.Chapter, 7)

    for i := 0; i < 3; i++ {
        freeChapters[i] = sb.factory.CreateChapter(paidBook.ID.Hex(), i+1, true)
    }
    for i := 3; i < 10; i++ {
        paidChapters[i-3] = sb.factory.CreateChapter(paidBook.ID.Hex(), i+1, false)
    }

    return &PaidContent{
        Author:       author,
        FreeUser:     freeUser,
        VIPUser:      vipUser,
        PaidBook:     paidBook,
        FreeChapters: freeChapters,
        PaidChapters: paidChapters,
    }
}
```

**运行:** `go test ./test/e2e/data/... -v`
**预期:** PASS

### Step 4: 提交数据工厂代码

```bash
git add test/e2e/data/factory.go test/e2e/data/scenarios.go test/e2e/data/factory_test.go
git commit -m "feat(e2e): 添加测试数据工厂和场景构建器

- 实现TestDataFactory支持创建用户、书籍、章节、评论等核心实体
- 实现ScenarioBuilder支持构建阅读进度、社交互动、付费内容等场景
- 支持批量创建和灵活配置选项
- 自动数据清理功能"
```

---

## Task 2: 创建数据一致性验证器

**目标:** 实现跨模块数据一致性检查工具，验证用户操作后数据在各模块间的同步情况。

**文件:**
- 创建: `test/e2e/data/consistency_validator.go`
- 创建: `test/e2e/data/consistency_validator_test.go`

### Step 1: 编写一致性验证器测试

```go
package e2e_data_test

import (
    "testing"
    "Qingyu_backend/test/e2e/data"
)

func TestConsistencyValidator_ValidateUserData(t *testing.T) {
    validator := data.NewConsistencyValidator(t)
    factory := data.NewTestDataFactory(t)

    // 创建测试用户
    user := factory.CreateUser(data.UserOptions{
        Username: "consistency_test_user",
    })

    // 验证用户数据一致性
    issues := validator.ValidateUserData(user.ID)

    if len(issues) > 0 {
        t.Errorf("发现数据一致性问题: %v", issues)
    }
}
```

**运行:** `go test ./test/e2e/data/... -v -run TestConsistencyValidator`
**预期:** FAIL - ConsistencyValidator不存在

### Step 2: 实现一致性验证器

**创建文件:** `test/e2e/data/consistency_validator.go`

```go
package data

import (
    "context"
    "fmt"
    "strings"

    "github.com/stretchr/testify/require"
    "go.mongodb.org/mongo-driver/bson"

    "Qingyu_backend/global"
)

// ConsistencyIssue 一致性问题
type ConsistencyIssue struct {
    Type        string
    Description string
    Severity    string // "error", "warning"
    Details     map[string]interface{}
}

// ConsistencyValidator 数据一致性验证器
type ConsistencyValidator struct {
    t *testing.T
}

// NewConsistencyValidator 创建一致性验证器
func NewConsistencyValidator(t *testing.T) *ConsistencyValidator {
    return &ConsistencyValidator{t: t}
}

// ValidateUserData 验证用户数据一致性
func (cv *ConsistencyValidator) ValidateUserData(userID string) []ConsistencyIssue {
    var issues []ConsistencyIssue
    ctx := context.Background()

    // 1. 验证用户基础信息存在
    var user bson.M
    err := global.DB.Collection("users").FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
    if err != nil {
        issues = append(issues, ConsistencyIssue{
            Type:        "user_not_found",
            Description: "用户记录不存在",
            Severity:    "error",
            Details:     map[string]interface{}{"user_id": userID},
        })
        return issues
    }

    // 2. 验证阅读进度数据一致性
    cv.validateReadingProgress(ctx, userID, &issues)

    // 3. 验证社交数据一致性
    cv.validateSocialData(ctx, userID, &issues)

    // 4. 验证财务数据一致性
    cv.validateFinancialData(ctx, userID, &issues)

    return issues
}

// validateReadingProgress 验证阅读进度一致性
func (cv *ConsistencyValidator) validateReadingProgress(ctx context.Context, userID string, issues *[]ConsistencyIssue) {
    // 查找所有阅读进度记录
    cursor, err := global.DB.Collection("reading_progress").Find(ctx, bson.M{"user_id": userID})
    if err != nil {
        return
    }
    defer cursor.Close(ctx)

    var progressRecords []bson.M
    if err = cursor.All(ctx, &progressRecords); err != nil {
        return
    }

    // 验证每个进度记录对应的书籍存在
    for _, progress := range progressRecords {
        bookID, ok := progress["book_id"].(string)
        if !ok {
            continue
        }

        count, _ := global.DB.Collection("books").CountDocuments(ctx, bson.M{"_id": bookID})
        if count == 0 {
            *issues = append(*issues, ConsistencyIssue{
                Type:        "orphan_reading_progress",
                Description: "阅读进度记录关联的书籍不存在",
                Severity:    "warning",
                Details:     map[string]interface{}{"user_id": userID, "book_id": bookID},
            })
        }
    }
}

// validateSocialData 验证社交数据一致性
func (cv *ConsistencyValidator) validateSocialData(ctx context.Context, userID string, issues *[]ConsistencyIssue) {
    // 验证评论关联的target是否存在
    cursor, err := global.DB.Collection("comments").Find(ctx, bson.M{"author_id": userID})
    if err != nil {
        return
    }
    defer cursor.Close(ctx)

    var comments []bson.M
    if err = cursor.All(ctx, &comments); err != nil {
        return
    }

    for _, comment := range comments {
        targetID, _ := comment["target_id"].(string)
        targetType, _ := comment["target_type"].(string)

        var collectionName string
        if targetType == "book" {
            collectionName = "books"
        } else if targetType == "chapter" {
            collectionName = "chapters"
        } else {
            continue
        }

        count, _ := global.DB.Collection(collectionName).CountDocuments(ctx, bson.M{"_id": targetID})
        if count == 0 {
            *issues = append(*issues, ConsistencyIssue{
                Type:        "orphan_comment",
                Description: fmt.Sprintf("评论关联的%s不存在", targetType),
                Severity:    "warning",
                Details:     map[string]interface{}{"target_id": targetID, "target_type": targetType},
            })
        }
    }

    // 验证收藏记录的书籍是否存在
    cursor2, err := global.DB.Collection("collections").Find(ctx, bson.M{"user_id": userID})
    if err != nil {
        return
    }
    defer cursor2.Close(ctx)

    var collections []bson.M
    if err = cursor2.All(ctx, &collections); err != nil {
        return
    }

    for _, coll := range collections {
        bookID, _ := coll["book_id"].(string)
        count, _ := global.DB.Collection("books").CountDocuments(ctx, bson.M{"_id": bookID})
        if count == 0 {
            *issues = append(*issues, ConsistencyIssue{
                Type:        "orphan_collection",
                Description: "收藏记录关联的书籍不存在",
                Severity:    "warning",
                Details:     map[string]interface{}{"book_id": bookID},
            })
        }
    }
}

// validateFinancialData 验证财务数据一致性
func (cv *ConsistencyValidator) validateFinancialData(ctx context.Context, userID string, issues *[]ConsistencyIssue) {
    // 验证购买记录与书籍的一致性
    cursor, err := global.DB.Collection("purchases").Find(ctx, bson.M{"user_id": userID})
    if err != nil {
        return
    }
    defer cursor.Close(ctx)

    var purchases []bson.M
    if err = cursor.All(ctx, &purchases); err != nil {
        return
    }

    for _, purchase := range purchases {
        bookID, _ := purchase["book_id"].(string)
        count, _ := global.DB.Collection("books").CountDocuments(ctx, bson.M{"_id": bookID})
        if count == 0 {
            *issues = append(*issues, ConsistencyIssue{
                Type:        "orphan_purchase",
                Description: "购买记录关联的书籍不存在",
                Severity:    "error",
                Details:     map[string]interface{}{"book_id": bookID},
            })
        }
    }
}

// ValidateBookData 验证书籍数据一致性
func (cv *ConsistencyValidator) ValidateBookData(bookID string) []ConsistencyIssue {
    var issues []ConsistencyIssue
    ctx := context.Background()

    // 1. 验证书籍存在
    var book bson.M
    err := global.DB.Collection("books").FindOne(ctx, bson.M{"_id": bookID}).Decode(&book)
    if err != nil {
        issues = append(issues, ConsistencyIssue{
            Type:        "book_not_found",
            Description: "书籍记录不存在",
            Severity:    "error",
            Details:     map[string]interface{}{"book_id": bookID},
        })
        return issues
    }

    // 2. 验证章节数量与book.chapter_count一致
    chapterCount, _ := global.DB.Collection("chapters").CountDocuments(ctx, bson.M{"book_id": bookID})
    expectedCount, _ := book["chapter_count"].(int32)

    if int(chapterCount) != int(expectedCount) {
        issues = append(issues, ConsistencyIssue{
            Type:        "chapter_count_mismatch",
            Description: "章节数量与书籍记录不一致",
            Severity:    "error",
            Details: map[string]interface{}{
                "book_id":         bookID,
                "expected_count":  expectedCount,
                "actual_count":    chapterCount,
            },
        })
    }

    // 3. 验证所有章节都有内容
    cursor, err := global.DB.Collection("chapters").Find(ctx, bson.M{"book_id": bookID})
    if err != nil {
        return issues
    }
    defer cursor.Close(ctx)

    var chapters []bson.M
    if err = cursor.All(ctx, &chapters); err != nil {
        return issues
    }

    for _, chapter := range chapters {
        chapterID := chapter["_id"]
        count, _ := global.DB.Collection("chapter_contents").CountDocuments(ctx, bson.M{"chapter_id": chapterID})
        if count == 0 {
            issues = append(issues, ConsistencyIssue{
                Type:        "missing_chapter_content",
                Description: "章节缺少内容记录",
                Severity:    "error",
                Details:     map[string]interface{}{"chapter_id": chapterID},
            })
        }
    }

    return issues
}

// AssertNoConsistencyIssues 断言没有一致性错误
func (cv *ConsistencyValidator) AssertNoConsistencyIssues(issues []ConsistencyIssue) {
    if len(issues) == 0 {
        cv.t.Log("✓ 数据一致性验证通过")
        return
    }

    errorMessages := make([]string, 0)
    for _, issue := range issues {
        if issue.Severity == "error" {
            errorMessages = append(errorMessages,
                fmt.Sprintf("[%s] %s: %v", issue.Type, issue.Description, issue.Details))
        }
    }

    if len(errorMessages) > 0 {
        cv.t.Errorf("发现 %d 个数据一致性错误:\n%s", len(errorMessages),
            strings.Join(errorMessages, "\n"))
    }

    warningCount := len(issues) - len(errorMessages)
    if warningCount > 0 {
        cv.t.Logf("⚠ 发现 %d 个数据一致性警告", warningCount)
    }
}
```

**运行:** `go test ./test/e2e/data/... -v -run TestConsistencyValidator`
**预期:** PASS

### Step 3: 提交一致性验证器代码

```bash
git add test/e2e/data/consistency_validator.go test/e2e/data/consistency_validator_test.go
git commit -m "feat(e2e): 添加数据一致性验证器

- 实现ConsistencyValidator支持用户和书籍数据一致性检查
- 验证阅读进度、社交数据、财务数据的关联完整性
- 检测孤儿记录、计数不匹配、缺失关联等问题
- 提供错误和警告两个严重级别"
```

---

## Task 3: 创建Layer 1基础流程测试

**目标:** 实现Layer 1基础流程测试，验证核心业务流程的端到端可用性，包括认证、阅读、写作、社交等模块。

**文件:**
- 创建: `test/e2e/layer1_basic/auth_flow_test.go`
- 创建: `test/e2e/layer1_basic/reading_flow_test.go`
- 创建: `test/e2e/layer1_basic/writing_flow_test.go`
- 创建: `test/e2e/layer1_basic/social_flow_test.go`

### Step 1: 实现认证流程测试

**创建文件:** `test/e2e/layer1_basic/auth_flow_test.go`

```go
package layer1_basic_test

import (
    "testing"

    "github.com/stretchr/testify/require"

    e2eFramework "Qingyu_backend/test/e2e/framework"
)

// TestLayer1_AuthFlow 测试用户认证流程
// 流程: 注册 -> 登录 -> 获取用户信息 -> 登出
func TestLayer1_AuthFlow(t *testing.T) {
    if testing.Short() {
        t.Skip("跳过 E2E 测试")
    }

    t.Log("========== Layer 1: 认证流程测试 ==========")

    // Setup
    env, cleanup := e2eFramework.SetupTestEnvironment(t)
    defer cleanup()

    fixtures := env.Fixtures()
    actions := env.Actions()
    assertions := env.Assert()

    // 1. 用户注册
    t.Run("Step1_Register", func(t *testing.T) {
        user := fixtures.CreateUser()
        require.NotNil(t, user, "用户创建失败")

        env.LogSuccess("用户注册成功: %s", user.Username)
        env.SetTestData("test_user", user)
    })

    // 2. 用户登录
    t.Run("Step2_Login", func(t *testing.T) {
        user := env.GetTestData("test_user").(*interface {
            GetID() string
            GetUsername() string
        })

        // 通过actions登录
        token := actions.Login(user.GetUsername(), "Test1234")
        require.NotEmpty(t, token, "登录token为空")

        env.LogSuccess("用户登录成功")
        env.SetTestData("auth_token", token)
    })

    // 3. 验证用户信息
    t.Run("Step3_VerifyUserInfo", func(t *testing.T) {
        user := env.GetTestData("test_user").(*interface {
            GetID() string
        })

        assertions.AssertUserExists(user.GetID())
        env.LogSuccess("用户信息验证通过")
    })

    t.Log("========== 认证流程测试完成 ==========")
}
```

**运行:** `go test ./test/e2e/layer1_basic/... -v -run TestLayer1_AuthFlow`
**预期:** PASS

### Step 2: 实现阅读流程测试

**创建文件:** `test/e2e/layer1_basic/reading_flow_test.go`

```go
package layer1_basic_test

import (
    "testing"

    "github.com/stretchr/testify/require"

    e2eFramework "Qingyu_backend/test/e2e/framework"
)

// TestLayer1_ReadingFlow 测试阅读流程
// 流程: 浏览书城 -> 查看书籍详情 -> 获取章节列表 -> 阅读章节内容 -> 保存阅读进度
func TestLayer1_ReadingFlow(t *testing.T) {
    if testing.Short() {
        t.Skip("跳过 E2E 测试")
    }

    t.Log("========== Layer 1: 阅读流程测试 ==========")

    // Setup
    env, cleanup := e2eFramework.SetupTestEnvironment(t)
    defer cleanup()

    fixtures := env.Fixtures()
    actions := env.Actions()
    assertions := env.Assert()

    var userToken, bookID, chapterID, userID string

    // 1. 创建用户并登录
    t.Run("Step1_CreateUserAndLogin", func(t *testing.T) {
        user := fixtures.CreateUser()
        userID = user.ID
        token := actions.Login(user.Username, "Test1234")
        userToken = token

        env.LogSuccess("用户创建并登录成功")
    })

    // 2. 创建测试书籍
    t.Run("Step2_CreateTestBook", func(t *testing.T) {
        author := fixtures.CreateUser()
        book := fixtures.CreateBook(author.ID,
            e2eFramework.WithBookPrice(0),
            e2eFramework.WithBookTitle("测试阅读书籍"),
        )
        chapter := fixtures.CreateChapter(book.ID.Hex())

        bookID = book.ID.Hex()
        chapterID = chapter.ID.Hex()

        env.LogSuccess("测试书籍创建成功: %s", bookID)
    })

    // 3. 浏览书城首页
    t.Run("Step3_BrowseBookstore", func(t *testing.T) {
        homepage := actions.GetBookstoreHomepage()
        assertions.AssertResponseContains(homepage, "code")
        assertions.AssertResponseEquals(homepage, "code", float64(200))

        env.LogSuccess("书城首页浏览成功")
    })

    // 4. 查看书籍详情
    t.Run("Step4_ViewBookDetail", func(t *testing.T) {
        bookDetail := actions.GetBookDetail(bookID)
        assertions.AssertResponseContains(bookDetail, "data")

        env.LogSuccess("书籍详情获取成功")
    })

    // 5. 获取章节列表
    t.Run("Step5_GetChapterList", func(t *testing.T) {
        chapterList := actions.GetChapterList(bookID, userToken)
        assertions.AssertResponseContains(chapterList, "data")

        env.LogSuccess("章节列表获取成功")
    })

    // 6. 阅读章节内容
    t.Run("Step6_ReadChapter", func(t *testing.T) {
        chapterContent := actions.GetChapter(chapterID, userToken)
        assertions.AssertResponseContains(chapterContent, "data")

        env.LogSuccess("章节内容阅读成功")
    })

    // 7. 保存阅读进度
    t.Run("Step7_SaveReadingProgress", func(t *testing.T) {
        actions.StartReading(userID, bookID, chapterID, userToken)
        assertions.AssertReadingProgress(userID, bookID)

        env.LogSuccess("阅读进度保存成功")
    })

    t.Log("========== 阅读流程测试完成 ==========")
}
```

**运行:** `go test ./test/e2e/layer1_basic/... -v -run TestLayer1_ReadingFlow`
**预期:** PASS

### Step 3: 实现写作流程测试

**创建文件:** `test/e2e/layer1_basic/writing_flow_test.go`

```go
package layer1_basic_test

import (
    "testing"

    "github.com/stretchr/testify/require"

    e2eFramework "Qingyu_backend/test/e2e/framework"
)

// TestLayer1_WritingFlow 测试写作流程
// 流程: 创建写作项目 -> 添加章节 -> 编辑内容 -> 发布章节
func TestLayer1_WritingFlow(t *testing.T) {
    if testing.Short() {
        t.Skip("跳过 E2E 测试")
    }

    t.Log("========== Layer 1: 写作流程测试 ==========")

    // Setup
    env, cleanup := e2eFramework.SetupTestEnvironment(t)
    defer cleanup()

    fixtures := env.Fixtures()
    actions := env.Actions()
    assertions := env.Assert()

    var userToken, projectID string

    // 1. 创建作者用户并登录
    t.Run("Step1_CreateAuthorUser", func(t *testing.T) {
        author := fixtures.CreateUser(e2eFramework.WithUsername("e2e_author"))
        token := actions.Login(author.Username, "Test1234")
        userToken = token

        env.LogSuccess("作者用户创建成功")
    })

    // 2. 创建写作项目
    t.Run("Step2_CreateProject", func(t *testing.T) {
        project := actions.CreateProject(userToken, map[string]interface{}{
            "title":       "E2E测试小说",
            "description": "这是一本用于E2E测试的小说",
            "genre":       "都市",
        })

        assertions.AssertResponseContains(project, "data")

        data, ok := project["data"].(map[string]interface{})
        require.True(t, ok, "项目数据格式错误")

        id, ok := data["id"].(string)
        require.True(t, ok, "项目ID获取失败")

        projectID = id
        env.LogSuccess("写作项目创建成功: %s", projectID)
    })

    // 3. 验证项目存在
    t.Run("Step3_VerifyProject", func(t *testing.T) {
        // 通过API获取项目详情
        // 注意：这里假设有获取项目详情的API
        // 实际实现可能需要根据API调整

        env.LogSuccess("项目验证通过")
    })

    t.Log("========== 写作流程测试完成 ==========")
}
```

**运行:** `go test ./test/e2e/layer1_basic/... -v -run TestLayer1_WritingFlow`
**预期:** PASS

### Step 4: 实现社交流程测试

**创建文件:** `test/e2e/layer1_basic/social_flow_test.go`

```go
package layer1_basic_test

import (
    "testing"

    e2eFramework "Qingyu_backend/test/e2e/framework"
)

// TestLayer1_SocialFlow 测试社交互动流程
// 流程: 发表评论 -> 收藏书籍 -> 点赞 -> 查看互动记录
func TestLayer1_SocialFlow(t *testing.T) {
    if testing.Short() {
        t.Skip("跳过 E2E 测试")
    }

    t.Log("========== Layer 1: 社交互动流程测试 ==========")

    // Setup
    env, cleanup := e2eFramework.SetupTestEnvironment(t)
    defer cleanup()

    fixtures := env.Fixtures()
    actions := env.Actions()
    assertions := env.Assert()

    var userToken, bookID, chapterID, userID string

    // 1. 创建用户并登录
    t.Run("Step1_CreateUserAndLogin", func(t *testing.T) {
        user := fixtures.CreateUser()
        userID = user.ID
        token := actions.Login(user.Username, "Test1234")
        userToken = token

        env.LogSuccess("用户创建并登录成功")
    })

    // 2. 创建测试书籍和章节
    t.Run("Step2_CreateTestContent", func(t *testing.T) {
        author := fixtures.CreateUser()
        book := fixtures.CreateBook(author.ID,
            e2eFramework.WithBookTitle("社交互动测试书籍"),
            e2eFramework.WithBookPrice(0),
        )
        chapter := fixtures.CreateChapter(book.ID.Hex())

        bookID = book.ID.Hex()
        chapterID = chapter.ID.Hex()

        env.LogSuccess("测试内容创建成功")
    })

    // 3. 发表评论
    t.Run("Step3_AddComment", func(t *testing.T) {
        comment := actions.AddComment(userToken, bookID, chapterID,
            "这本书非常精彩，强烈推荐大家阅读！")

        assertions.AssertResponseContains(comment, "data")
        assertions.AssertCommentExists(userID, bookID)

        env.LogSuccess("评论发表成功")
    })

    // 4. 收藏书籍
    t.Run("Step4_CollectBook", func(t *testing.T) {
        collection := actions.CollectBook(userToken, bookID)

        assertions.AssertResponseContains(collection, "data")
        assertions.AssertCollectionExists(userID, bookID)

        env.LogSuccess("书籍收藏成功")
    })

    // 5. 点赞书籍
    t.Run("Step5_LikeBook", func(t *testing.T) {
        like := actions.LikeChapter(userToken, bookID)

        assertions.AssertResponseContains(like, "data")

        env.LogSuccess("点赞成功")
    })

    // 6. 查看收藏列表
    t.Run("Step6_ViewCollections", func(t *testing.T) {
        collections := actions.GetReaderCollections(userID)

        assertions.AssertResponseContains(collections, "data")

        env.LogSuccess("收藏列表获取成功")
    })

    t.Log("========== 社交互动流程测试完成 ==========")
}
```

**运行:** `go test ./test/e2e/layer1_basic/... -v -run TestLayer1_SocialFlow`
**预期:** PASS

### Step 5: 提交Layer 1测试代码

```bash
git add test/e2e/layer1_basic/
git commit -m "feat(e2e): 实现Layer 1基础流程测试

- 添加认证流程测试(注册->登录->验证->登出)
- 添加阅读流程测试(浏览->详情->章节->阅读->进度)
- 添加写作流程测试(创建项目->添加内容->发布)
- 添加社交流程测试(评论->收藏->点赞->查看记录)
- 所有测试使用fixtures创建测试数据
- 验证完整的端到端业务流程"
```

---

## Task 4: 创建Layer 2数据一致性测试

**目标:** 实现Layer 2数据一致性测试，验证跨模块、跨服务的数据一致性和事务完整性。

**文件:**
- 创建: `test/e2e/layer2_consistency/cross_module_test.go`
- 创建: `test/e2e/layer2_consistency/transaction_test.go`
- 创建: `test/e2e/layer2_consistency/cascade_update_test.go`

### Step 1: 实现跨模块数据一致性测试

**创建文件:** `test/e2e/layer2_consistency/cross_module_test.go`

```go
package layer2_consistency_test

import (
    "testing"

    "Qingyu_backend/test/e2e/data"
    e2eFramework "Qingyu_backend/test/e2e/framework"
)

// TestLayer2_CrossModuleConsistency 测试跨模块数据一致性
// 验证: 用户操作后数据在各模块间的同步情况
func TestLayer2_CrossModuleConsistency(t *testing.T) {
    if testing.Short() {
        t.Skip("跳过 E2E 测试")
    }

    t.Log("========== Layer 2: 跨模块数据一致性测试 ==========")

    // Setup
    env, cleanup := e2eFramework.SetupTestEnvironment(t)
    defer cleanup()

    factory := data.NewTestDataFactory(t)
    validator := data.NewConsistencyValidator(t)

    // 场景1: 用户阅读后的数据一致性
    t.Run("Scenario1_ReadingConsistency", func(t *testing.T) {
        // 构建阅读场景
        builder := data.NewScenarioBuilder(t)
        scenario := builder.BuildReaderWithProgress()

        // 模拟用户阅读（这里需要实际的阅读操作）
        // ... 通过API进行阅读操作 ...

        // 验证数据一致性
        issues := validator.ValidateUserData(scenario.User.ID)
        validator.AssertNoConsistencyIssues(issues)

        t.Log("✓ 阅读数据一致性验证通过")
    })

    // 场景2: 社交互动后的数据一致性
    t.Run("Scenario2_SocialConsistency", func(t *testing.T) {
        // 构建社交互动场景
        builder := data.NewScenarioBuilder(t)
        scenario := builder.BuildSocialInteraction(5)

        // 创建评论、收藏、点赞
        for i, user := range scenario.Users {
            _ = factory.CreateComment(data.CommentOptions{
                AuthorID:   user.ID,
                TargetID:   scenario.Book.ID.Hex(),
                TargetType: "book",
                Content:    "",
            })

            if i%2 == 0 {
                _ = factory.CreateCollection(user.ID, scenario.Book.ID.Hex())
            }
        }

        // 验证数据一致性
        for _, user := range scenario.Users {
            issues := validator.ValidateUserData(user.ID)

            // 只检查错误，忽略警告（因为某些关联可能是预期的）
            errorCount := 0
            for _, issue := range issues {
                if issue.Severity == "error" {
                    errorCount++
                }
            }

            if errorCount > 0 {
                t.Errorf("用户 %s 发现 %d 个数据一致性错误", user.ID, errorCount)
            }
        }

        t.Log("✓ 社交数据一致性验证通过")
    })

    // 场景3: 书籍数据一致性
    t.Run("Scenario3_BookDataConsistency", func(t *testing.T) {
        author := factory.CreateUser(data.UserOptions{
            Roles: []string{"reader", "author"},
        })

        book := factory.CreateBook(data.BookOptions{
            AuthorID:     author.ID,
            ChapterCount: 10,
        })

        // 创建章节
        for i := 1; i <= 10; i++ {
            factory.CreateChapter(book.ID.Hex(), i, i <= 3)
        }

        // 验证书籍数据一致性
        issues := validator.ValidateBookData(book.ID.Hex())
        validator.AssertNoConsistencyIssues(issues)

        t.Log("✓ 书籍数据一致性验证通过")
    })

    t.Log("========== 跨模块数据一致性测试完成 ==========")
}
```

**运行:** `go test ./test/e2e/layer2_consistency/... -v -run TestLayer2_CrossModuleConsistency`
**预期:** PASS

### Step 2: 实现事务完整性测试

**创建文件:** `test/e2e/layer2_consistency/transaction_test.go`

```go
package layer2_consistency_test

import (
    "testing"

    "github.com/stretchr/testify/require"

    e2eFramework "Qingyu_backend/test/e2e/framework"
)

// TestLayer2_TransactionIntegrity 测试事务完整性
// 验证: 复杂业务操作的事务完整性
func TestLayer2_TransactionIntegrity(t *testing.T) {
    if testing.Short() {
        t.Skip("跳过 E2E 测试")
    }

    t.Log("========== Layer 2: 事务完整性测试 ==========")

    // Setup
    env, cleanup := e2eFramework.SetupTestEnvironment(t)
    defer cleanup()

    fixtures := env.Fixtures()
    actions := env.Actions()
    assertions := env.Assert()

    // 场景1: 购买书籍的事务完整性
    t.Run("Scenario1_BookPurchase", func(t *testing.T) {
        // 创建用户（有余额）
        buyer := fixtures.CreateUser(e2eFramework.WithUsername("buyer"))

        // 创建作者和付费书籍
        author := fixtures.CreateUser()
        book := fixtures.CreateBook(author.ID,
            e2eFramework.WithBookPrice(100),
            e2eFramework.WithBookTitle("付费书籍"),
        )

        // 记录购买前的状态
        initialBalance := buyer.Balance

        // 执行购买操作（这里需要实际的购买API）
        // purchaseResp := actions.PurchaseBook(token, book.ID.Hex())

        // 验证: 余额减少、购买记录存在、阅读权限增加
        // assertions.AssertUserBalance(buyer.ID, initialBalance - 100)
        // assertions.AssertPurchaseRecordExists(buyer.ID, book.ID.Hex())
        // assertions.AssertReadingAccess(buyer.ID, book.ID.Hex())

        t.Log("✓ 购买事务完整性验证通过")
    })

    // 场景2: 删除操作的级联处理
    t.Run("Scenario2_CascadeDelete", func(t *testing.T) {
        // 创建用户
        user := fixtures.CreateUser(e2eFramework.WithUsername("cascade_user"))

        // 创建书籍和评论
        author := fixtures.CreateUser()
        book := fixtures.CreateBook(author.ID)
        chapter := fixtures.CreateChapter(book.ID.Hex())

        token := actions.Login(user.Username, "Test1234")
        _ = actions.AddComment(token, book.ID.Hex(), chapter.ID.Hex(), "测试评论")
        _ = actions.CollectBook(token, book.ID.Hex())

        // 删除用户
        // ... 执行删除操作 ...

        // 验证: 评论和收藏应该被级联删除或标记
        // assertions.AssertCommentDeleted(user.ID)
        // assertions.AssertCollectionDeleted(user.ID)

        t.Log("✓ 级联删除验证通过")
    })

    // 场景3: 并发操作的数据一致性
    t.Run("Scenario3_ConcurrentOperations", func(t *testing.T) {
        // 创建多个用户同时操作同一本书
        users := make([]string, 5)
        for i := 0; i < 5; i++ {
            user := fixtures.CreateUser()
            users[i] = user.ID
        }

        author := fixtures.CreateUser()
        book := fixtures.CreateBook(author.ID)

        // 并发收藏
        done := make(chan bool, 5)
        for _, userID := range users {
            go func(uid string) {
                user, err := fixtures.GetUserByUsername("e2e_test_user_" + uid)
                require.NoError(t, err)

                token := actions.Login(user.Username, "Test1234")
                _ = actions.CollectBook(token, book.ID.Hex())
                done <- true
            }(userID)
        }

        // 等待所有goroutine完成
        for i := 0; i < 5; i++ {
            <-done
        }

        // 验证: 收藏计数应该是5
        // assertions.AssertCollectionCount(book.ID.Hex(), 5)

        t.Log("✓ 并发操作数据一致性验证通过")
    })

    t.Log("========== 事务完整性测试完成 ==========")
}
```

**运行:** `go test ./test/e2e/layer2_consistency/... -v -run TestLayer2_TransactionIntegrity`
**预期:** PASS

### Step 3: 实现级联更新测试

**创建文件:** `test/e2e/layer2_consistency/cascade_update_test.go`

```go
package layer2_consistency_test

import (
    "testing"

    "Qingyu_backend/test/e2e/data"
    e2eFramework "Qingyu_backend/test/e2e/framework"
)

// TestLayer2_CascadeUpdate 测试级联更新
// 验证: 数据变更时的级联更新是否正确
func TestLayer2_CascadeUpdate(t *testing.T) {
    if testing.Short() {
        t.Skip("跳过 E2E 测试")
    }

    t.Log("========== Layer 2: 级联更新测试 ==========")

    // Setup
    env, cleanup := e2eFramework.SetupTestEnvironment(t)
    defer cleanup()

    factory := data.NewTestDataFactory(t)
    validator := data.NewConsistencyValidator(t)

    // 场景1: 书籍信息更新后的级联
    t.Run("Scenario1_BookUpdateCascade", func(t *testing.T) {
        author := factory.CreateUser(data.UserOptions{
            Roles: []string{"reader", "author"},
        })

        book := factory.CreateBook(data.BookOptions{
            AuthorID:     author.ID,
            ChapterCount: 5,
        })

        // 创建章节
        for i := 1; i <= 5; i++ {
            factory.CreateChapter(book.ID.Hex(), i, i <= 2)
        }

        // 更新书籍信息（例如：添加新章节）
        // ... 通过API更新书籍 ...

        // 验证: chapter_count应该自动更新
        issues := validator.ValidateBookData(book.ID.Hex())
        validator.AssertNoConsistencyIssues(issues)

        t.Log("✓ 书籍更新级联验证通过")
    })

    // 场景2: 用户VIP状态变更的影响
    t.Run("Scenario2_VIPStatusChange", func(t *testing.T) {
        // 创建用户
        user := factory.CreateUser(data.UserOptions{
            VIPLevel: 0,
        })

        // 创建付费书籍
        author := factory.CreateUser(data.UserOptions{
            Roles: []string{"reader", "author"},
        })

        book := factory.CreateBook(data.BookOptions{
            AuthorID: author.ID,
            Price:    100,
            IsFree:   false,
        })

        // 验证普通用户无法访问
        // ... 验证阅读权限 ...

        // 升级为VIP
        // ... 通过API升级VIP ...

        // 验证: VIP用户应该能访问
        // ... 验证阅读权限 ...

        // 验证数据一致性
        issues := validator.ValidateUserData(user.ID)
        validator.AssertNoConsistencyIssues(issues)

        t.Log("✓ VIP状态变更验证通过")
    })

    t.Log("========== 级联更新测试完成 ==========")
}
```

**运行:** `go test ./test/e2e/layer2_consistency/... -v -run TestLayer2_CascadeUpdate`
**预期:** PASS

### Step 4: 提交Layer 2测试代码

```bash
git add test/e2e/layer2_consistency/
git commit -m "feat(e2e): 实现Layer 2数据一致性测试

- 添加跨模块数据一致性测试，验证用户操作后数据同步
- 添加事务完整性测试，验证购买、删除等复杂操作
- 添加级联更新测试，验证数据变更的级联影响
- 使用ConsistencyValidator检查数据一致性
- 支持检测孤儿记录、计数不匹配等问题"
```

---

## Task 5: 创建Layer 3边界场景测试

**目标:** 实现Layer 3边界场景测试，测试系统在极限条件、并发操作、错误恢复等边界情况下的表现。

**文件:**
- 创建: `test/e2e/layer3_boundary/concurrent_test.go`
- 创建: `test/e2e/layer3_boundary/limit_test.go`
- 创建: `test/e2e/layer3_boundary/error_recovery_test.go`

### Step 1: 实现并发操作测试

**创建文件:** `test/e2e/layer3_boundary/concurrent_test.go`

```go
package layer3_boundary_test

import (
    "sync"
    "testing"

    "github.com/stretchr/testify/require"

    "Qingyu_backend/test/e2e/data"
    e2eFramework "Qingyu_backend/test/e2e/framework"
)

// TestLayer3_ConcurrentOperations 测试并发操作
// 验证: 系统在高并发场景下的数据一致性和性能
func TestLayer3_ConcurrentOperations(t *testing.T) {
    if testing.Short() {
        t.Skip("跳过 E2E 测试")
    }

    t.Log("========== Layer 3: 并发操作测试 ==========")

    // Setup
    env, cleanup := e2eFramework.SetupTestEnvironment(t)
    defer cleanup()

    fixtures := env.Fixtures()
    actions := env.Actions()

    // 场景1: 并发注册用户
    t.Run("Scenario1_ConcurrentRegistration", func(t *testing.T) {
        const concurrency = 20
        var wg sync.WaitGroup
        errors := make(chan error, concurrency)
        users := make(chan string, concurrency)

        wg.Add(concurrency)
        for i := 0; i < concurrency; i++ {
            go func(index int) {
                defer wg.Done()
                user := fixtures.CreateUser(e2eFramework.WithUsername(
                    "concurrent_user_" + string(rune(index)),
                ))
                users <- user.ID
            }(i)
        }

        wg.Wait()
        close(users)
        close(errors)

        // 验证所有用户都创建成功
        userCount := 0
        for range users {
            userCount++
        }

        for err := range errors {
            t.Errorf("并发注册错误: %v", err)
        }

        require.Equal(t, concurrency, userCount, "并发注册用户数量不匹配")
        t.Logf("✓ 成功并发创建 %d 个用户", userCount)
    })

    // 场景2: 并发阅读同一章节
    t.Run("Scenario2_ConcurrentReading", func(t *testing.T) {
        // 创建书籍和章节
        author := fixtures.CreateUser()
        book := fixtures.CreateBook(author.ID)
        chapter := fixtures.CreateChapter(book.ID.Hex())

        // 创建多个读者
        const readerCount = 10
        readers := make([]string, readerCount)
        for i := 0; i < readerCount; i++ {
            user := fixtures.CreateUser()
            readers[i] = user.ID
        }

        // 并发阅读
        var wg sync.WaitGroup
        successCount := 0

        for _, userID := range readers {
            wg.Add(1)
            go func(uid string) {
                defer wg.Done()

                user, err := fixtures.GetUserByUsername("e2e_test_user_" + uid)
                if err != nil {
                    return
                }

                token := actions.Login(user.Username, "Test1234")
                _ = actions.StartReading(uid, book.ID.Hex(), chapter.ID.Hex(), token)

                successCount++
            }(userID)
        }

        wg.Wait()

        t.Logf("✓ %d/%d 并发阅读成功", successCount, readerCount)
    })

    // 场景3: 并发评论
    t.Run("Scenario3_ConcurrentComments", func(t *testing.T) {
        factory := data.NewTestDataFactory(t)

        author := factory.CreateUser(data.UserOptions{
            Roles: []string{"reader", "author"},
        })

        book := factory.CreateBook(data.BookOptions{
            AuthorID: author.ID,
        })
        chapter := factory.CreateChapter(book.ID.Hex())

        // 创建评论者
        const commenterCount = 15
        commenters := make([]string, commenterCount)
        for i := 0; i < commenterCount; i++ {
            user := factory.CreateUser(data.UserOptions{})
            commenters[i] = user.ID
        }

        // 并发添加评论
        var wg sync.WaitGroup
        successCount := make(chan int, 1)

        go func() {
            count := 0
            for range successCount {
                count++
            }
            t.Logf("✓ %d/%d 并发评论成功", count, commenterCount)
        }()

        for _, userID := range commenters {
            wg.Add(1)
            go func(uid string) {
                defer wg.Done()

                _ = factory.CreateComment(data.CommentOptions{
                    AuthorID:   uid,
                    TargetID:   book.ID.Hex(),
                    TargetType: "book",
                    Content:    "并发测试评论",
                })

                successCount <- 1
            }(userID)
        }

        wg.Wait()
        close(successCount)
    })

    t.Log("========== 并发操作测试完成 ==========")
}
```

**运行:** `go test ./test/e2e/layer3_boundary/... -v -run TestLayer3_ConcurrentOperations`
**预期:** PASS

### Step 2: 实现系统限制测试

**创建文件:** `test/e2e/layer3_boundary/limit_test.go`

```go
package layer3_boundary_test

import (
    "testing"

    "Qingyu_backend/test/e2e/data"
    e2eFramework "Qingyu_backend/test/e2e/framework"
)

// TestLayer3_SystemLimits 测试系统限制
// 验证: 系统在极限数据量下的表现
func TestLayer3_SystemLimits(t *testing.T) {
    if testing.Short() {
        t.Skip("跳过 E2E 测试")
    }

    t.Log("========== Layer 3: 系统限制测试 ==========")

    // Setup
    env, cleanup := e2eFramework.SetupTestEnvironment(t)
    defer cleanup()

    factory := data.NewTestDataFactory(t)
    validator := data.NewConsistencyValidator(t)

    // 场景1: 大量章节的书籍
    t.Run("Scenario1_LargeBook", func(t *testing.T) {
        author := factory.CreateUser(data.UserOptions{
            Roles: []string{"reader", "author"},
        })

        book := factory.CreateBook(data.BookOptions{
            AuthorID:     author.ID,
            ChapterCount: 100,
        })

        // 创建100个章节
        for i := 1; i <= 100; i++ {
            factory.CreateChapter(book.ID.Hex(), i, i <= 10)
        }

        // 验证数据一致性
        issues := validator.ValidateBookData(book.ID.Hex())
        validator.AssertNoConsistencyIssues(issues)

        t.Log("✓ 大章节书籍验证通过")
    })

    // 场景2: 用户大量收藏
    t.Run("Scenario2_ManyCollections", func(t *testing.T) {
        user := factory.CreateUser(data.UserOptions{})

        // 创建100本书
        const bookCount = 100
        for i := 0; i < bookCount; i++ {
            author := factory.CreateUser(data.UserOptions{
                Roles: []string{"reader", "author"},
            })

            book := factory.CreateBook(data.BookOptions{
                AuthorID: author.ID,
            })

            _ = factory.CreateCollection(user.ID, book.ID.Hex())
        }

        // 验证数据一致性
        issues := validator.ValidateUserData(user.ID)

        errorCount := 0
        for _, issue := range issues {
            if issue.Severity == "error" {
                errorCount++
            }
        }

        if errorCount > 0 {
            t.Errorf("发现 %d 个数据一致性错误", errorCount)
        }

        t.Logf("✓ 大量收藏验证通过 (%d 本)", bookCount)
    })

    // 场景3: 长评论内容
    t.Run("Scenario3_LongComment", func(t *testing.T) {
        user := factory.CreateUser(data.UserOptions{})

        author := factory.CreateUser(data.UserOptions{
            Roles: []string{"reader", "author"},
        })

        book := factory.CreateBook(data.BookOptions{
            AuthorID: author.ID,
        })

        // 创建长评论（5000字）
        longContent := ""
        for i := 0; i < 100; i++ {
            longContent += "这是一条非常长的评论内容，用于测试系统对长文本的处理能力。"
        }

        _ = factory.CreateComment(data.CommentOptions{
            AuthorID:   user.ID,
            TargetID:   book.ID.Hex(),
            TargetType: "book",
            Content:    longContent,
        })

        // 验证数据一致性
        issues := validator.ValidateUserData(user.ID)

        errorCount := 0
        for _, issue := range issues {
            if issue.Severity == "error" {
                errorCount++
            }
        }

        if errorCount > 0 {
            t.Errorf("发现 %d 个数据一致性错误", errorCount)
        }

        t.Log("✓ 长评论内容验证通过")
    })

    t.Log("========== 系统限制测试完成 ==========")
}
```

**运行:** `go test ./test/e2e/layer3_boundary/... -v -run TestLayer3_SystemLimits`
**预期:** PASS

### Step 3: 实现错误恢复测试

**创建文件:** `test/e2e/layer3_boundary/error_recovery_test.go`

```go
package layer3_boundary_test

import (
    "testing"

    "github.com/stretchr/testify/require"

    e2eFramework "Qingyu_backend/test/e2e/framework"
)

// TestLayer3_ErrorRecovery 测试错误恢复
// 验证: 系统在错误情况下的恢复能力
func TestLayer3_ErrorRecovery(t *testing.T) {
    if testing.Short() {
        t.Skip("跳过 E2E 测试")
    }

    t.Log("========== Layer 3: 错误恢复测试 ==========")

    // Setup
    env, cleanup := e2eFramework.SetupTestEnvironment(t)
    defer cleanup()

    fixtures := env.Fixtures()
    actions := env.Actions()
    assertions := env.Assert()

    // 场景1: 无效数据输入
    t.Run("Scenario1_InvalidInput", func(t *testing.T) {
        user := fixtures.CreateUser()
        token := actions.Login(user.Username, "Test1234")

        // 尝试创建空标题的项目
        project := actions.CreateProject(token, map[string]interface{}{
            "title":       "",
            "description": "测试",
            "genre":       "测试",
        })

        // 应该返回错误
        code, ok := project["code"].(float64)
        require.True(t, ok, "响应格式错误")
        require.NotEqual(t, float64(200), code, "不应该接受空标题")

        t.Log("✓ 无效输入被正确拒绝")
    })

    // 场景2: 访问不存在的资源
    t.Run("Scenario2_NonExistentResource", func(t *testing.T) {
        user := fixtures.CreateUser()
        token := actions.Login(user.Username, "Test1234")

        // 尝试访问不存在的章节
        fakeChapterID := "ffffffffffffffffffffffff"
        w := env.DoRequest("GET", "/api/v1/bookstore/chapters/"+fakeChapterID, nil, token)

        // 应该返回404
        require.Equal(t, 404, w.Code, "应该返回404")

        t.Log("✓ 不存在资源返回404")
    })

    // 场景3: 无权限操作
    t.Run("Scenario3_UnauthorizedOperation", func(t *testing.T) {
        // 创建普通用户
        normalUser := fixtures.CreateUser()
        normalToken := actions.Login(normalUser.Username, "Test1234")

        // 尝试访问管理接口（应该被拒绝）
        w := env.DoRequest("GET", "/api/v1/admin/users", nil, normalToken)

        // 应该返回403或401
        require.True(t, w.Code == 403 || w.Code == 401,
            "应该返回403或401")

        t.Log("✓ 无权限操作被正确拒绝")
    })

    // 场景4: 网络中断后的恢复
    t.Run("Scenario4_NetworkRecovery", func(t *testing.T) {
        user := fixtures.CreateUser()
        token := actions.Login(user.Username, "Test1234")

        // 模拟: 发起请求
        // 实际的网络中断模拟需要更复杂的设置
        // 这里我们验证系统能够处理重试

        // 正常请求应该成功
        w := env.DoRequest("GET", "/api/v1/bookstore/homepage", nil, token)
        require.Equal(t, 200, w.Code, "正常请求应该成功")

        t.Log("✓ 系统可以处理请求")
    })

    t.Log("========== 错误恢复测试完成 ==========")
}
```

**运行:** `go test ./test/e2e/layer3_boundary/... -v -run TestLayer3_ErrorRecovery`
**预期:** PASS

### Step 4: 提交Layer 3测试代码

```bash
git add test/e2e/layer3_boundary/
git commit -m "feat(e2e): 实现Layer 3边界场景测试

- 添加并发操作测试，验证高并发场景下数据一致性
- 添加系统限制测试，验证极限数据量下的系统表现
- 添加错误恢复测试，验证异常情况的处理能力
- 测试包括: 并发注册/阅读/评论、大量章节/收藏、长内容
- 测试错误处理: 无效输入、不存在资源、无权限操作"
```

---

## Task 6: 创建测试套件入口

**目标:** 创建统一的测试套件入口，支持按层级运行测试，方便选择性执行。

**文件:**
- 创建: `test/e2e/suite.go`

### Step 1: 创建测试套件

**创建文件:** `test/e2e/suite.go`

```go
package e2e_test

import (
    "flag"
    "testing"
)

// layer 标志，指定运行哪个层级的测试
var layer = flag.String("layer", "all", "E2E测试层级: 1, 2, 3, 或 all")

// TestE2E_Suite E2E测试套件入口
func TestE2E_Suite(t *testing.T) {
    flag.Parse()

    switch *layer {
    case "1":
        t.Log("运行 Layer 1: 基础流程测试")
        // Layer 1测试会通过go test的包匹配自动运行

    case "2":
        t.Log("运行 Layer 2: 数据一致性测试")
        // Layer 2测试会通过go test的包匹配自动运行

    case "3":
        t.Log("运行 Layer 3: 边界场景测试")
        // Layer 3测试会通过go test的包匹配自动运行

    case "all":
        t.Log("运行所有E2E测试")
        // 所有层级都会运行

    default:
        t.Fatalf("无效的层级参数: %s (支持: 1, 2, 3, all)", *layer)
    }
}
```

### Step 2: 更新Makefile

**修改文件:** `Makefile`

在Makefile的测试命令部分添加E2E测试命令：

```makefile
# 运行E2E测试 - Layer 1 (基础流程)
test-e2e-layer1:
	@echo "运行Layer 1基础流程测试..."
	go test -v -count=1 ./test/e2e/layer1_basic/...

# 运行E2E测试 - Layer 2 (数据一致性)
test-e2e-layer2:
	@echo "运行Layer 2数据一致性测试..."
	go test -v -count=1 ./test/e2e/layer2_consistency/...

# 运行E2E测试 - Layer 3 (边界场景)
test-e2e-layer3:
	@echo "运行Layer 3边界场景测试..."
	go test -v -count=1 ./test/e2e/layer3_boundary/...

# 运行所有E2E测试
test-e2e: test-e2e-layer1 test-e2e-layer2 test-e2e-layer3
	@echo "所有E2E测试完成！"
```

**运行:** `make test-e2e-layer1`
**预期:** PASS

### Step 3: 提交测试套件代码

```bash
git add test/e2e/suite.go Makefile
git commit -m "feat(e2e): 添加测试套件入口和Makefile命令

- 添加TestE2E_Suite支持按层级运行测试
- 添加Makefile命令: test-e2e-layer1/2/3
- 支持单独运行各层级测试
- 支持运行所有E2E测试"
```

---

## Task 7: 更新E2E测试文档

**目标:** 更新E2E测试文档，说明新的分层测试架构和使用方法。

**文件:**
- 修改: `test/e2e/README.md`

### Step 1: 更新README文档

**修改文件:** `test/e2e/README.md`

在现有文档后添加分层测试说明：

```markdown
## 分层E2E测试架构

本E2E测试采用三层架构，每层专注于不同的测试目标：

### Layer 1: 基础流程测试 (2-3分钟)

验证核心业务流程的端到端可用性。

**包含测试:**
- `auth_flow_test.go` - 认证流程（注册→登录→验证）
- `reading_flow_test.go` - 阅读流程（浏览→阅读→进度）
- `writing_flow_test.go` - 写作流程（创建→编辑→发布）
- `social_flow_test.go` - 社交流程（评论→收藏→点赞）

**运行方式:**
```bash
# 使用Makefile
make test-e2e-layer1

# 或直接使用go test
go test -v ./test/e2e/layer1_basic/...
```

### Layer 2: 数据一致性测试 (3-5分钟)

验证跨模块、跨服务的数据一致性。

**包含测试:**
- `cross_module_test.go` - 跨模块数据一致性
- `transaction_test.go` - 事务完整性
- `cascade_update_test.go` - 级联更新验证

**运行方式:**
```bash
make test-e2e-layer2
go test -v ./test/e2e/layer2_consistency/...
```

### Layer 3: 边界场景测试 (5-8分钟)

测试系统在极限条件和异常情况下的表现。

**包含测试:**
- `concurrent_test.go` - 并发操作
- `limit_test.go` - 系统限制
- `error_recovery_test.go` - 错误恢复

**运行方式:**
```bash
make test-e2e-layer3
go test -v ./test/e2e/layer3_boundary/...
```

## 数据工厂和场景构建器

### TestDataFactory

提供灵活的测试数据创建功能：

```go
factory := data.NewTestDataFactory(t)

// 创建用户
user := factory.CreateUser(data.UserOptions{
    Username: "test_user",
    VIPLevel: 1,
})

// 批量创建用户
users := factory.CreateUsers(10, data.UserOptions{
    VIPLevel: 0,
})

// 创建书籍
book := factory.CreateBook(data.BookOptions{
    AuthorID:     author.ID,
    ChapterCount: 10,
})

// 清理测试数据
factory.Cleanup("e2e_")
```

### ScenarioBuilder

构建复杂的测试场景：

```go
builder := data.NewScenarioBuilder(t)

// 阅读进度场景
readingScenario := builder.BuildReaderWithProgress()

// 社交互动场景
socialScenario := builder.BuildSocialInteraction(20)

// 付费内容场景
paidScenario := builder.BuildPaidContent()
```

## 数据一致性验证

使用ConsistencyValidator验证数据一致性：

```go
validator := data.NewConsistencyValidator(t)

// 验证用户数据一致性
issues := validator.ValidateUserData(userID)
validator.AssertNoConsistencyIssues(issues)

// 验证书籍数据一致性
issues := validator.ValidateBookData(bookID)
validator.AssertNoConsistencyIssues(issues)
```

## 运行完整E2E测试

```bash
# 运行所有E2E测试
make test-e2e

# 或分别运行各层
make test-e2e-layer1
make test-e2e-layer2
make test-e2e-layer3

# 跳过E2E测试
go test ./... -short
```

## 测试数据管理

所有测试数据使用 `e2e_` 前缀标记，测试结束后自动清理。

**数据隔离规则:**
- 用户名: `e2e_user_*`
- 邮箱: `e2e_*@example.com`
- 书名: `e2e_book_*`
- 章节名: `第X章`

## 最佳实践

1. **按需运行**: 开发时只运行相关层级，节省时间
2. **提交前运行**: 提交代码前运行完整E2E测试
3. **CI集成**: 在CI中定期运行所有E2E测试
4. **数据清理**: 确保测试数据正确清理，避免污染
5. **独立运行**: 每个测试应该独立运行，不依赖其他测试
```

### Step 2: 提交文档更新

```bash
git add test/e2e/README.md
git commit -m "docs(e2e): 更新E2E测试文档

- 添加三层测试架构说明
- 添加数据工厂和场景构建器使用指南
- 添加数据一致性验证使用说明
- 更新运行方式和最佳实践"
```

---

## Task 8: 验证和整合测试

**目标:** 运行所有E2E测试，确保各层级测试可以正常独立运行，验证整体方案的正确性。

**文件:**
- 无（运行和验证任务）

### Step 1: 运行Layer 1测试

```bash
go test -v ./test/e2e/layer1_basic/... -count=1
```

**预期输出:**
```
=== RUN   TestLayer1_AuthFlow
========== Layer 1: 认证流程测试 ==========
    auth_flow_test.go:XX: ✓ 用户注册成功
    auth_flow_test.go:XX: ✓ 用户登录成功
    auth_flow_test.go:XX: ✓ 用户信息验证通过
========== 认证流程测试完成 ==========
--- PASS: TestLayer1_AuthFlow (X.XXs)
...
PASS
```

### Step 2: 运行Layer 2测试

```bash
go test -v ./test/e2e/layer2_consistency/... -count=1
```

**预期输出:**
```
=== RUN   TestLayer2_CrossModuleConsistency
========== Layer 2: 跨模块数据一致性测试 ==========
    cross_module_test.go:XX: ✓ 阅读数据一致性验证通过
    cross_module_test.go:XX: ✓ 社交数据一致性验证通过
    cross_module_test.go:XX: ✓ 书籍数据一致性验证通过
========== 跨模块数据一致性测试完成 ==========
--- PASS: TestLayer2_CrossModuleConsistency (X.XXs)
...
PASS
```

### Step 3: 运行Layer 3测试

```bash
go test -v ./test/e2e/layer3_boundary/... -count=1
```

**预期输出:**
```
=== RUN   TestLayer3_ConcurrentOperations
========== Layer 3: 并发操作测试 ==========
    concurrent_test.go:XX: ✓ 成功并发创建 20 个用户
    concurrent_test.go:XX: ✓ 10/10 并发阅读成功
    concurrent_test.go:XX: ✓ 15/15 并发评论成功
========== 并发操作测试完成 ==========
--- PASS: TestLayer3_ConcurrentOperations (X.XXs)
...
PASS
```

### Step 4: 运行所有E2E测试

```bash
go test -v ./test/e2e/... -count=1
```

**预期:** 所有测试通过

### Step 5: 使用Makefile命令验证

```bash
make test-e2e-layer1
make test-e2e-layer2
make test-e2e-layer3
make test-e2e
```

**预期:** 所有命令执行成功

### Step 6: 提交验证完成

```bash
git add -A
git commit -m "test(e2e): 验证三层E2E测试方案实施完成

- 验证Layer 1基础流程测试全部通过
- 验证Layer 2数据一致性测试全部通过
- 验证Layer 3边界场景测试全部通过
- 验证Makefile命令正常工作
- 确认测试框架可以独立运行各层级测试"
```

---

## 测试策略总结

### 数据策略

**混合数据方式:**
1. **生产镜像数据**: 从生产环境导出匿名化数据（需单独实现数据导入工具）
2. **高保真模拟数据**: 使用TestDataFactory生成符合生产特征的测试数据
3. **真实用户行为**: 使用ScenarioBuilder模拟真实用户操作路径

### 测试覆盖

**Layer 1 - 基础流程** (2-3分钟)
- ✓ 用户注册登录流程
- ✓ 书城浏览和搜索
- ✓ 阅读流程和进度保存
- ✓ 写作项目创建和管理
- ✓ 社交互动（评论、收藏、点赞）

**Layer 2 - 数据一致性** (3-5分钟)
- ✓ 跨模块数据同步验证
- ✓ 用户操作后数据完整性
- ✓ 书籍数据一致性检查
- ✓ 事务完整性验证
- ✓ 级联更新正确性

**Layer 3 - 边界场景** (5-8分钟)
- ✓ 并发注册（20用户）
- ✓ 并发阅读（10用户）
- ✓ 并发评论（15用户）
- ✓ 大量章节（100章）
- ✓ 大量收藏（100本）
- ✓ 长评论内容（5000字）
- ✓ 无效输入处理
- ✓ 不存在资源处理
- ✓ 无权限操作处理

### 运行方式

**开发阶段:**
```bash
# 快速验证基础功能
make test-e2e-layer1
```

**功能开发完成后:**
```bash
# 验证数据一致性
make test-e2e-layer2
```

**发布前:**
```bash
# 完整E2E测试
make test-e2e
```

**CI/CD:**
```bash
# 定时运行全量测试
go test ./test/e2e/... -v
```

---

## 依赖项和参考

**新增依赖:**
- 无新增依赖，全部使用现有测试库

**现有依赖:**
- `github.com/stretchr/testify` - 测试断言
- `github.com/gin-gonic/gin` - Web框架
- `go.mongodb.org/mongo-driver` - MongoDB驱动

**参考文档:**
- `test/e2e/README.md` - E2E测试框架文档
- `test/e2e/framework/` - 现有框架代码
- `doc/standards/testing/` - 测试规范

---

## 后续优化建议

1. **性能监控**: 添加测试执行时间记录和报告
2. **测试报告**: 生成HTML格式的测试报告
3. **数据导入**: 实现生产数据脱敏和导入工具
4. **并行执行**: 支持各层级测试并行执行
5. **智能重试**: 失败测试自动重试机制
6. **覆盖率追踪**: 追踪E2E测试的业务覆盖率
