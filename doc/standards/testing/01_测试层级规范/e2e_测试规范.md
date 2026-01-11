# E2E测试规范

## 概述

E2E（End-to-End）测试模拟完整的用户业务流程，从前端到后端到数据库的完整链路。本规范定义了E2E测试的详细要求和最佳实践。

## 核心原则

### ✅ 模拟真实用户场景

```go
// ✅ 正确示例：完整业务流程
func TestE2E_BookPurchaseFlow(t *testing.T) {
    // Setup - 完整测试环境
    env, cleanup := e2e.SetupTestEnvironment(t)
    defer cleanup()

    // 场景：读者购买书籍并阅读
    // 1. 作者发布书籍
    author := env.CreateUser(&users.User{
        Username: "author",
        Email:    "author@example.com",
        Password: "password123",
    })
    authorToken := env.Login(author.Email, "password123")

    book := env.CreateBook(&bookstore.Book{
        Title:       "测试书籍",
        Description: "精彩绝伦的小说",
        AuthorID:    author.ID,
        Price:       100,
    })

    // 2. 读者注册并登录
    reader := env.CreateUser(&users.User{
        Username: "reader",
        Email:    "reader@example.com",
        Password: "password123",
    })
    readerToken := env.Login(reader.Email, "password123")

    // 3. 读者充值
    env.Recharge(reader.ID, 500)

    // 4. 读者购买书籍
    purchase := env.PurchaseBook(reader.ID, book.ID, readerToken)
    require.NotNil(t, purchase)

    // 5. 验证读者余额扣减
    readerAfter := env.GetUser(reader.ID)
    assert.Equal(t, 400, readerAfter.Balance)

    // 6. 验证作者收入增加
    authorAfter := env.GetUser(author.ID)
    assert.Equal(t, 100, authorAfter.Balance)

    // 7. 读者开始阅读
    progress := env.StartReading(reader.ID, book.ID, "chapter1")
    assert.NotNil(t, progress)

    // 8. 读者添加书签
    bookmark := env.AddBookmark(reader.ID, book.ID, "chapter1", 100)
    assert.NotNil(t, bookmark)

    // 9. 验证完整数据链路
    env.VerifyPurchaseRecord(reader.ID, book.ID)
    env.VerifyReadingProgress(reader.ID, book.ID)
}
```

### ❌ 严格禁止

```go
// ❌ 错误：只测试单个API接口
func TestE2E_BookPurchase_SingleAPI(t *testing.T) {
    // 问题：这不是E2E测试，这只是API测试
    router, cleanup := integration.SetupTestEnvironment(t)
    defer cleanup()

    helper := integration.NewTestHelper(t, router)

    // 只测试了一个购买接口
    w := helper.DoAuthRequest("POST", "/api/v1/bookstore/purchase", reqBody, token)
    assert.Equal(t, 200, w.Code)

    // 没有验证：
    // - 用户余额是否真的扣减
    // - 作者是否真的收到钱
    // - 购买记录是否正确
    // - 后续是否能阅读
}
```

## 测试组织结构

### 文件位置

```
test/e2e/{scenario}_test.go
```

示例：
```
test/e2e/book_purchase_flow_test.go
test/e2e/reader_reading_flow_test.go
test/e2e/author_publishing_flow_test.go
test/e2e/payment_flow_test.go
```

### 包命名

```go
package e2e_test

import (
    "testing"
    "time"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "Qingyu_backend/test/e2e/framework"
)
```

## E2E测试框架

### 框架结构

```go
// test/e2e/framework/environment.go
package e2e

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/require"
    "go.mongodb.org/mongo-driver/bson/primitive"

    "Qingyu_backend/container"
    "Qingyu_backend/models/bookstore"
    "Qingyu_backend/models/users"
    "Qingyu_backend/repository/mongodb/bookstore"
    "Qingyu_backend/repository/mongodb/reader"
    "Qingyu_backend/repository/mongodb/user"
)

type TestEnvironment struct {
    T         *testing.T
    Container *container.ServiceContainer
    Router    *gin.Engine
    Client    *http.Client
}

func SetupTestEnvironment(t *testing.T) (*TestEnvironment, func()) {
    t.Helper()

    // 初始化完整环境
    cfg, err := config.LoadConfig("config/config_test.yaml")
    require.NoError(t, err)

    c := container.NewServiceContainer(cfg)
    err = c.Initialize(context.Background())
    require.NoError(t, err)

    // 设置HTTP路由
    gin.SetMode(gin.TestMode)
    router := gin.New()
    routerGroup := router.Group("/api/v1")
    router.RegisterRoutes(routerGroup, c)

    env := &TestEnvironment{
        T:         t,
        Container: c,
        Router:    router,
        Client:    &http.Client{Timeout: 10 * time.Second},
    }

    cleanup := func() {
        // 清理所有集合
        db := c.GetMongoDB()
        collections := []string{
            "users", "books", "projects", "documents",
            "roles", "permissions", "transactions",
            "purchases", "reading_progress", "bookmarks",
            "chapters", "annotations", "announcements",
        }
        ctx := context.Background()
        for _, coll := range collections {
            _ = db.Collection(coll).Drop(ctx)
        }
        _ = c.Close(ctx)
    }

    return env, cleanup
}

// CreateUser 创建测试用户
func (env *TestEnvironment) CreateUser(user *users.User) *users.User {
    env.T.Helper()

    if user == nil {
        user = &users.User{
            Username: "user_" + primitive.NewObjectID().Hex(),
            Email:    primitive.NewObjectID().Hex() + "@example.com",
            Password: "password123",
        }
    }

    userRepo := userRepo.NewUserRepository(env.Container.GetMongoDB())
    err := userRepo.Create(context.Background(), user)
    require.NoError(env.T, err)

    return user
}

// Login 用户登录并返回token
func (env *TestEnvironment) Login(email, password string) string {
    env.T.Helper()

    // 调用登录API
    reqBody := map[string]interface{}{
        "email":    email,
        "password": password,
    }
    w := env.DoRequest("POST", "/api/v1/auth/login", reqBody, nil)

    require.Equal(env.T, 200, w.Code, "登录失败")

    var resp map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &resp)

    token := resp["data"].(map[string]interface{})["token"].(string)
    return token
}

// CreateBook 创建测试书籍
func (env *TestEnvironment) CreateBook(book *bookstore.Book) *bookstore.Book {
    env.T.Helper()

    bookRepo := bookstoreRepo.NewBookRepository(env.Container.GetMongoDB())
    err := bookRepo.Create(context.Background(), book)
    require.NoError(env.T, err)

    return book
}

// Recharge 用户充值
func (env *TestEnvironment) Recharge(userID string, amount int) {
    env.T.Helper()

    // 创建充值交易
    tx := &Transaction{
        UserID:          userID,
        Type:            "recharge",
        Amount:          amount,
        Status:          "completed",
        TransactionID:   primitive.NewObjectID().Hex(),
    }

    txRepo := repository.NewTransactionRepository(env.Container.GetMongoDB())
    err := txRepo.Create(context.Background(), tx)
    require.NoError(env.T, err)

    // 更新用户余额
    userRepo := userRepo.NewUserRepository(env.Container.GetMongoDB())
    user, err := userRepo.GetByID(context.Background(), userID)
    require.NoError(env.T, err)

    user.Balance += amount
    err = userRepo.Update(context.Background(), user)
    require.NoError(env.T, err)
}

// PurchaseBook 购买书籍
func (env *TestEnvironment) PurchaseBook(userID, bookID, token string) *Purchase {
    env.T.Helper()

    reqBody := map[string]interface{}{
        "book_id": bookID,
    }
    w := env.DoAuthRequest("POST", "/api/v1/bookstore/purchase", reqBody, token)

    require.Equal(env.T, 200, w.Code, "购买失败")

    var resp map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &resp)

    purchase := &Purchase{}
    data := resp["data"].(map[string]interface{})
    purchase.ID = data["id"].(string)
    purchase.BookID = data["book_id"].(string)
    purchase.UserID = data["user_id"].(string)

    return purchase
}

// GetUser 获取用户信息
func (env *TestEnvironment) GetUser(userID string) *users.User {
    env.T.Helper()

    userRepo := userRepo.NewUserRepository(env.Container.GetMongoDB())
    user, err := userRepo.GetByID(context.Background(), userID)
    require.NoError(env.T, err)

    return user
}

// StartReading 开始阅读
func (env *TestEnvironment) StartReading(userID, bookID, chapterID string) *ReadingProgress {
    env.T.Helper()

    reqBody := map[string]interface{}{
        "book_id":    bookID,
        "chapter_id": chapterID,
    }
    w := env.DoAuthRequest("POST", "/api/v1/reader/progress", reqBody, token)

    require.Equal(env.T, 200, w.Code, "开始阅读失败")

    var resp map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &resp)

    progress := &ReadingProgress{}
    data := resp["data"].(map[string]interface{})
    progress.ID = data["id"].(string)

    return progress
}

// AddBookmark 添加书签
func (env *TestEnvironment) AddBookmark(userID, bookID, chapterID string, position int) *Bookmark {
    env.T.Helper()

    reqBody := map[string]interface{}{
        "book_id":    bookID,
        "chapter_id": chapterID,
        "position":   position,
    }
    w := env.DoAuthRequest("POST", "/api/v1/reader/bookmarks", reqBody, token)

    require.Equal(env.T, 200, w.Code, "添加书签失败")

    var resp map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &resp)

    bookmark := &Bookmark{}
    data := resp["data"].(map[string]interface{})
    bookmark.ID = data["id"].(string)

    return bookmark
}

// VerifyPurchaseRecord 验证购买记录
func (env *TestEnvironment) VerifyPurchaseRecord(userID, bookID string) {
    env.T.Helper()

    purchaseRepo := readerRepo.NewPurchaseRepository(env.Container.GetMongoDB())
    purchase, err := purchaseRepo.GetByUserAndBook(context.Background(), userID, bookID)
    require.NoError(env.T, err, "购买记录不存在")
    assert.NotNil(env.T, purchase)
}

// VerifyReadingProgress 验证阅读进度
func (env *TestEnvironment) VerifyReadingProgress(userID, bookID string) {
    env.T.Helper()

    progressRepo := readerRepo.NewReadingProgressRepository(env.Container.GetMongoDB())
    progress, err := progressRepo.GetByUserAndBook(context.Background(), userID, bookID)
    require.NoError(env.T, err, "阅读进度不存在")
    assert.NotNil(env.T, progress)
}

// DoRequest 执行HTTP请求
func (env *TestEnvironment) DoRequest(method, path string, body interface{}, token string) *httptest.ResponseRecorder {
    env.T.Helper()

    var bodyReader *bytes.Reader
    if body != nil {
        bodyBytes, _ := json.Marshal(body)
        bodyReader = bytes.NewReader(bodyBytes)
    } else {
        bodyReader = bytes.NewReader([]byte{})
    }

    req := httptest.NewRequest(method, path, bodyReader)
    req.Header.Set("Content-Type", "application/json")

    if token != "" {
        req.Header.Set("Authorization", "Bearer "+token)
    }

    w := httptest.NewRecorder()
    env.Router.ServeHTTP(w, req)

    return w
}
```

## 测试用例设计

### 1. 书籍发布流程

```go
func TestE2E_AuthorPublishBookFlow(t *testing.T) {
    // Setup
    env, cleanup := e2e.SetupTestEnvironment(t)
    defer cleanup()

    // 1. 作者注册并登录
    author := env.CreateUser(&users.User{
        Username: "author",
        Email:    "author@example.com",
        Password: "password123",
    })
    authorToken := env.Login(author.Email, "password123")

    // 2. 创建书籍
    book := env.CreateBook(&bookstore.Book{
        Title:       "我的第一本书",
        Description: "这是一本关于写作的书",
        Category:    "小说",
        AuthorID:    author.ID,
        Price:       0, // 免费
    })

    // 3. 添加章节
    chapter1 := env.CreateChapter(&bookstore.Chapter{
        BookID:    book.ID,
        Title:     "第一章",
        Content:   "这是第一章的内容",
        ChapterNo: 1,
        IsFree:    true,
    })

    chapter2 := env.CreateChapter(&bookstore.Chapter{
        BookID:    book.ID,
        Title:     "第二章",
        Content:   "这是第二章的内容",
        ChapterNo: 2,
        IsFree:    false,
        Price:     10,
    })

    // 4. 发布书籍
    env.PublishBook(book.ID, authorToken)

    // 5. 验证书籍状态
    publishedBook := env.GetBook(book.ID)
    assert.Equal(t, "published", publishedBook.Status)

    // 6. 验证章节关联
    chapters := env.GetBookChapters(book.ID)
    assert.Len(t, chapters, 2)

    // 7. 验证作者作品列表
    authorBooks := env.GetAuthorBooks(author.ID)
    assert.GreaterOrEqual(t, len(authorBooks), 1)

    // 8. 读者可以发现并查看书籍
    reader := env.CreateUser(&users.User{
        Username: "reader",
        Email:    "reader@example.com",
        Password: "password123",
    })

    foundBook := env.SearchBook(book.Title)
    assert.NotNil(t, foundBook)
    assert.Equal(t, book.ID, foundBook.ID)

    // 9. 读者可以查看免费章节
    freeChapter := env.GetChapter(chapter1.ID, nil)
    assert.NotNil(t, freeChapter)
    assert.NotEmpty(t, freeChapter.Content)

    // 10. 读者不能直接查看付费章节
    paidChapter := env.GetChapter(chapter2.ID, nil)
    assert.Nil(t, paidChapter)
}
```

### 2. 购买和阅读流程

```go
func TestE2E_PurchaseAndReadingFlow(t *testing.T) {
    // Setup
    env, cleanup := e2e.SetupTestEnvironment(t)
    defer cleanup()

    // 1. 准备书籍
    author := env.CreateUser(&users.User{
        Username: "author",
        Email:    "author@example.com",
        Password: "password123",
    })
    authorToken := env.Login(author.Email, "password123")

    book := env.CreateBook(&bookstore.Book{
        Title:       "付费书籍",
        Description: "需要购买才能阅读",
        AuthorID:    author.ID,
        Price:       100,
    })

    chapter := env.CreateChapter(&bookstore.Chapter{
        BookID:    book.ID,
        Title:     "付费章节",
        Content:   "精彩的付费内容",
        ChapterNo: 1,
        Price:     100,
    })

    env.PublishBook(book.ID, authorToken)

    // 2. 读者注册并充值
    reader := env.CreateUser(&users.User{
        Username: "reader",
        Email:    "reader@example.com",
        Password: "password123",
    })
    readerToken := env.Login(reader.Email, "password123")

    env.Recharge(reader.ID, 500)

    // 3. 读者尝试购买
    purchase := env.PurchaseBook(reader.ID, book.ID, readerToken)
    assert.NotNil(t, purchase)

    // 4. 验证余额变化
    readerAfterPurchase := env.GetUser(reader.ID)
    assert.Equal(t, 400, readerAfterPurchase.Balance)

    // 5. 验证作者收入
    authorAfterPurchase := env.GetUser(author.ID)
    assert.Equal(t, 100, authorAfterPurchase.Balance)

    // 6. 读者可以查看付费章节
    paidChapter := env.GetChapter(chapter.ID, readerToken)
    assert.NotNil(t, paidChapter)
    assert.NotEmpty(t, paidChapter.Content)

    // 7. 读者开始阅读
    progress := env.StartReading(reader.ID, book.ID, chapter.ID)
    assert.NotNil(t, progress)

    // 8. 读者添加书签
    bookmark := env.AddBookmark(reader.ID, book.ID, chapter.ID, 100)
    assert.NotNil(t, bookmark)

    // 9. 验证阅读记录
    progressRecord := env.GetReadingProgress(reader.ID, book.ID)
    assert.NotNil(t, progressRecord)
    assert.Equal(t, chapter.ID, progressRecord.CurrentChapterID)

    // 10. 验证书签记录
    bookmarks := env.GetBookmarks(reader.ID, book.ID)
    assert.GreaterOrEqual(t, len(bookmarks), 1)

    // 11. 读者查看阅读历史
    history := env.GetReadingHistory(reader.ID)
    assert.GreaterOrEqual(t, len(history), 1)
}
```

### 3. 会员订阅流程

```go
func TestE2E_MembershipSubscriptionFlow(t *testing.T) {
    // Setup
    env, cleanup := e2e.SetupTestEnvironment(t)
    defer cleanup()

    // 1. 管理员创建会员套餐
    admin := env.CreateAdminUser()
    adminToken := env.Login(admin.Email, "password123")

    vipPlan := env.CreateMembershipPlan(&MembershipPlan{
        Name:        "VIP会员",
        Description: "畅享所有书籍",
        Price:       3000, // 30元
        Duration:    30,   // 30天
        Benefits: []string{
            "免费阅读所有书籍",
            "专属标识",
            "优先体验新书",
        },
    })

    // 2. 用户注册并充值
    user := env.CreateUser(&users.User{
        Username: "reader",
        Email:    "reader@example.com",
        Password: "password123",
    })
    userToken := env.Login(user.Email, "password123")

    env.Recharge(user.ID, 5000)

    // 3. 用户订阅会员
    subscription := env.SubscribeMembership(user.ID, vipPlan.ID, userToken)
    assert.NotNil(t, subscription)

    // 4. 验证会员状态
    userAfter := env.GetUser(user.ID)
    assert.True(t, userAfter.IsVIP)
    assert.NotNil(t, userAfter.VIPExpireAt)
    assert.True(t, userAfter.VIPExpireAt.After(time.Now()))

    // 5. 验证余额扣减
    assert.Equal(t, 2000, userAfter.Balance)

    // 6. VIP用户免费阅读付费书籍
    book := env.CreateBook(&bookstore.Book{
        Title:    "VIP专属书籍",
        AuthorID: admin.ID,
        Price:    100,
    })

    chapter := env.CreateChapter(&bookstore.Chapter{
        BookID:  book.ID,
        Title:   "VIP章节",
        Content: "VIP专属内容",
        Price:   100,
    })

    // VIP用户不需要购买就能阅读
    vipChapter := env.GetChapter(chapter.ID, userToken)
    assert.NotNil(t, vipChapter)
    assert.NotEmpty(t, vipChapter.Content)

    // 7. 验证会员权益记录
    benefits := env.GetMembershipBenefits(user.ID)
    assert.GreaterOrEqual(t, len(benefits), 1)

    // 8. 测试会员到期
    env.SetVIPExpired(user.ID)

    expiredUser := env.GetUser(user.ID)
    assert.False(t, expiredUser.IsVIP)

    // 会员到期后需要购买才能阅读
    expiredChapter := env.GetChapter(chapter.ID, userToken)
    assert.Nil(t, expiredChapter)
}
```

### 4. 作者收入提现流程

```go
func TestE2E_AuthorWithdrawalFlow(t *testing.T) {
    // Setup
    env, cleanup := e2e.SetupTestEnvironment(t)
    defer cleanup()

    // 1. 作者发布书籍并获得收入
    author := env.CreateUser(&users.User{
        Username: "author",
        Email:    "author@example.com",
        Password: "password123",
    })

    // 模拟书籍销售
    env.AddAuthorIncome(author.ID, 10000) // 100元收入

    authorWithIncome := env.GetUser(author.ID)
    assert.Equal(t, 10000, authorWithIncome.Balance)

    // 2. 作者申请提现
    authorToken := env.Login(author.Email, "password123")
    withdrawal := env.RequestWithdrawal(author.ID, 5000, authorToken)

    assert.NotNil(t, withdrawal)
    assert.Equal(t, "pending", withdrawal.Status)

    // 3. 验证余额冻结
    authorAfterRequest := env.GetUser(author.ID)
    assert.Equal(t, 5000, authorAfterRequest.Balance)

    // 4. 管理员审核通过
    admin := env.CreateAdminUser()
    adminToken := env.Login(admin.Email, "password123")

    env.ApproveWithdrawal(withdrawal.ID, adminToken)

    // 5. 验证提现完成
    approvedWithdrawal := env.GetWithdrawal(withdrawal.ID)
    assert.Equal(t, "completed", approvedWithdrawal.Status)

    // 6. 验证最终余额
    finalAuthor := env.GetUser(author.ID)
    assert.Equal(t, 5000, finalAuthor.Balance)

    // 7. 验证交易记录
    transactions := env.GetWithdrawalTransactions(author.ID)
    assert.GreaterOrEqual(t, len(transactions), 1)

    // 8. 测试提现失败场景
    env.AddAuthorIncome(author.ID, 3000) // 再增加30元

    // 余额不足无法提现
    failedWithdrawal, err := env.RequestWithdrawalWithError(author.ID, 10000, authorToken)
    assert.Error(t, err)
    assert.Nil(t, failedWithdrawal)
}
```

### 5. 互动功能流程

```go
func TestE2E_BookInteractionFlow(t *testing.T) {
    // Setup
    env, cleanup := e2e.SetupTestEnvironment(t)
    defer cleanup()

    // 1. 准备作者和书籍
    author := env.CreateUser(&users.User{
        Username: "author",
        Email:    "author@example.com",
        Password: "password123",
    })
    authorToken := env.Login(author.Email, "password123")

    book := env.CreateBook(&bookstore.Book{
        Title:    "互动书籍",
        AuthorID: author.ID,
        Price:    0,
    })

    chapter := env.CreateChapter(&bookstore.Chapter{
        BookID:  book.ID,
        Title:   "第一章",
        Content: "精彩内容",
    })

    // 2. 读者评论
    reader := env.CreateUser(&users.User{
        Username: "reader",
        Email:    "reader@example.com",
        Password: "password123",
    })
    readerToken := env.Login(reader.Email, "password123")

    comment := env.AddComment(reader.ID, book.ID, chapter.ID, readerToken, "写得太好了！")
    assert.NotNil(t, comment)

    // 3. 验证评论显示
    comments := env.GetBookComments(book.ID)
    assert.GreaterOrEqual(t, len(comments), 1)

    // 4. 作者回复评论
    reply := env.ReplyComment(author.ID, comment.ID, authorToken, "谢谢支持！")
    assert.NotNil(t, reply)

    // 5. 读者点赞章节
    env.LikeChapter(reader.ID, chapter.ID, readerToken)

    // 6. 验证点赞记录
    likes := env.GetChapterLikes(chapter.ID)
    assert.GreaterOrEqual(t, len(likes), 1)

    // 7. 读者收藏书籍
    collection := env.CollectBook(reader.ID, book.ID, readerToken)
    assert.NotNil(t, collection)

    // 8. 验证收藏列表
    collections := env.GetReaderCollections(reader.ID)
    assert.GreaterOrEqual(t, len(collections), 1)

    // 9. 读者打赏作者
    env.Recharge(reader.ID, 1000)
    reward := env.RewardAuthor(reader.ID, author.ID, book.ID, readerToken, 500)
    assert.NotNil(t, reward)

    // 10. 验证打赏记录
    rewards := env.GetBookRewards(book.ID)
    assert.GreaterOrEqual(t, len(rewards), 1)

    // 11. 验证作者收到打赏
    authorAfterReward := env.GetUser(author.ID)
    assert.GreaterOrEqual(t, authorAfterReward.Balance, 500)

    // 12. 测试取消收藏
    env.UncollectBook(reader.ID, book.ID, readerToken)

    collectionsAfter := env.GetReaderCollections(reader.ID)
    assert.Equal(t, len(collections)-1, len(collectionsAfter))

    // 13. 测试删除评论
    env.DeleteComment(reader.ID, comment.ID, readerToken)

    commentsAfter := env.GetBookComments(book.ID)
    // 剩下作者的回复
    assert.Len(t, commentsAfter, 1)
}
```

## 数据验证

### 完整数据链路验证

```go
func (env *TestEnvironment) VerifyCompletePurchaseFlow(userID, bookID string) {
    env.T.Helper()

    // 1. 验证购买记录
    purchaseRepo := readerRepo.NewPurchaseRepository(env.Container.GetMongoDB())
    purchase, err := purchaseRepo.GetByUserAndBook(context.Background(), userID, bookID)
    require.NoError(env.T, err)
    assert.NotNil(env.T, purchase)
    assert.Equal(env.T, "completed", purchase.Status)

    // 2. 验证交易记录
    txRepo := repository.NewTransactionRepository(env.Container.GetMongoDB())
    transactions, _ := txRepo.GetByUserID(context.Background(), userID)
    found := false
    for _, tx := range transactions {
        if tx.Type == "purchase" && tx.RelatedID == purchase.ID {
            found = true
            assert.Equal(env.T, "completed", tx.Status)
            break
        }
    }
    assert.True(env.T, found, "交易记录不存在")

    // 3. 验证阅读权限
    accessRepo := readerRepo.NewReadingAccessRepository(env.Container.GetMongoDB())
    hasAccess, _ := accessRepo.CheckAccess(context.Background(), userID, bookID)
    assert.True(env.T, hasAccess, "没有阅读权限")

    // 4. 验证用户余额
    userRepo := userRepo.NewUserRepository(env.Container.GetMongoDB())
    user, _ := userRepo.GetByID(context.Background(), userID)
    assert.NotNil(env.T, user)
}
```

## 性能监控

```go
func TestE2E_Performance_Monitoring(t *testing.T) {
    // Setup
    env, cleanup := e2e.SetupTestEnvironment(t)
    defer cleanup()

    // 准备测试数据
    author := env.CreateUser(&users.User{
        Username: "author",
        Email:    "author@example.com",
        Password: "password123",
    })

    for i := 0; i < 100; i++ {
        env.CreateBook(&bookstore.Book{
            Title:    fmt.Sprintf("书籍 %d", i),
            AuthorID: author.ID,
            Price:    0,
        })
    }

    // 测试列表查询性能
    start := time.Now()
    w := env.DoRequest("GET", "/api/v1/books?page=1&page_size=50", nil, nil)
    duration := time.Since(start)

    // 断言性能
    require.Equal(t, 200, w.Code)
    assert.Less(t, duration, 500*time.Millisecond, "列表查询应该在500ms内完成")

    var resp map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &resp)

    data := resp["data"].(map[string]interface{})
    items := data["items"].([]interface{})
    assert.Len(t, items, 50)
}
```

## 常见问题

### Q1: E2E测试和API集成测试的区别？

**A**:
- **API集成测试**：测试单个API端点或少数几个相关端点
- **E2E测试**：测试完整的业务流程，跨越多个模块和API

示例对比：
```go
// API集成测试 - 只测试购买接口
func TestAPI_PurchaseBook(t *testing.T) {
    router, cleanup := integration.SetupTestEnvironment(t)
    defer cleanup()

    helper := integration.NewTestHelper(t, router)

    // 只测试购买接口是否正常
    w := helper.DoAuthRequest("POST", "/api/v1/bookstore/purchase", reqBody, token)
    assert.Equal(t, 200, w.Code)
}

// E2E测试 - 测试完整购买流程
func TestE2E_PurchaseBook(t *testing.T) {
    env, cleanup := e2e.SetupTestEnvironment(t)
    defer cleanup()

    // 完整流程：注册 -> 充值 -> 购买 -> 验证余额 -> 验证权限 -> 阅读
}
```

### Q2: E2E测试应该覆盖多少流程？

**A**: 优先级原则：
1. **核心业务流程**：100%覆盖（购买、阅读、发布）
2. **常用功能**：80%覆盖（评论、收藏、点赞）
3. **边缘功能**：按需覆盖（举报、反馈）

不是所有功能都需要E2E测试，但关键流程必须有。

### Q3: 如何处理E2E测试中的异步操作？

**A**: 使用轮询和事件等待：
```go
func TestE2E_AsyncOperation(t *testing.T) {
    env, cleanup := e2e.SetupTestEnvironment(t)
    defer cleanup()

    // 触发异步操作
    jobID := env.StartAsyncJob()

    // 等待完成
    var job *Job
    assert.Eventually(t, func() bool {
        job = env.GetJob(jobID)
        return job != nil && job.Status == "completed"
    }, 30*time.Second, 1*time.Second, "异步操作未在30秒内完成")

    assert.Equal(t, "success", job.Result)
}
```

## 参考文档

- [API层测试规范](./api_层测试规范.md)
- [集成测试详细规范](../02_测试类型规范/集成测试详细规范.md)
- [测试辅助工具使用指南](../03_测试工具指南/测试辅助工具使用指南.md)
