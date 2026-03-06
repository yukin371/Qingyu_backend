# Issue #002: Repository Create 方法未回设 ID

**优先级**: 高 (P0)
**类型**: 业务逻辑 Bug
**状态**: ✅ 主流已修复（已审查）
**创建日期**: 2026-03-05
**相关报告**: [Writer DTO 重构总结报告](../reports/2026-03-05-dto-refactoring-summary.md#22-bug-outlinerepositorycreate-未回设-id)
**审查日期**: 2026-03-05
**审查报告**: [P0问题审查报告](../reports/2026-03-05-p0-issue-audit-report.md)

---

## 审查结果

**状态**: ✅ 主流已修复

### 审查发现

1. ✅ **主流 Repository 已正确实现 ID 回设**:
   - Writer域（Project, Document, Outline）
   - Reader域
   - Bookstore域
   - Social域

2. ✅ **两种正确模式**:
   - 使用 `result.InsertedID` 回设 ID
   - 预先生成 `primitive.NewObjectID()` 并赋值

3. ⚠️ **剩余问题主要集中在尚未完成 ObjectID 迁移的模型/仓储**:
   - auth 域本轮已补 `PermissionTemplate/Role/Permission/OAuthAccount/OAuthSession` 的 Create 回设
   - messaging 域本轮已补 `Message/MessageTemplate/NotificationDelivery` 的 Create 回设
   - writer 域本轮已补 `Version/Commit/FileRevision/FilePatch/Timeline/TimelineEvent` 的 Create 回设；其他仍保留 string 主键的模型在迁移时仍需逐个复核
   - writer 发布链路本轮已补 `PublicationRecord` 仓储 Create 回设
   - finance 域本轮已补 `Wallet/Transaction/WithdrawRequest` 的 Create 回设
   - bookstore 域本轮已补 `Chapter/Category` 的 Create 回设
   - ai 域活跃仓储 `ChatSession/ChatMessage/UserQuota/QuotaTransaction` 已具备 Create 回设；其余核心 metadata 模型已改为 `BeforeCreate + ObjectID`
   - reader 域本轮已补 `ReadingSettings/ReaderTheme` 的 Create 回设与 ObjectID 边界
   - notification 域本轮已补 `Notification/NotificationPreference/PushDevice/NotificationTemplate` 的 Create 回设与 ObjectID 边界
   - admin 域本轮已补 `AdminLog/AuditRecord` 的 Create 回设与 ObjectID 查询边界
   - recommendation 域本轮已补 `Behavior/ItemFeature/UserProfile/UserBehaviorRecord` 的 Create/Upsert 回设与 ObjectID 边界

### 已修复的 Repository

```go
// ✅ 正确实现示例 (writer/project_repository_mongo.go)
func (r *ProjectRepositoryMongo) Create(ctx context.Context, project *writer.Project) error {
    project.TouchForCreate()

    result, err := r.GetCollection().InsertOne(ctx, project)
    if err != nil {
        return errors.NewRepositoryError(...)
    }

    // ✅ 正确回设 ID
    if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
        project.ID = oid
    }

    return nil
}
```

---

## 问题描述

Repository 的 Create 方法在插入数据后未将 MongoDB 生成的 ID 回设到传入的对象中，导致后续操作（如 Update）因 ID 为空而失败。

### 已修复
- ✅ `repository/mongodb/writer/outline_repository_mongo.go`

### 待检查
- ⚠️ 其他所有 Repository 的 Create 方法

## 问题表现

```go
// 修复前的代码：
func (r *OutlineRepositoryMongo) Create(ctx context.Context, outline *writer.OutlineNode) error {
    outline.TouchForCreate()
    _, err := r.GetCollection().InsertOne(ctx, outline)
    if err != nil {
        return errors.NewRepositoryError(...)
    }
    // ❌ 没有将 ID 设置回 outline
    return nil
}

// 导致的问题：
outline := &writer.OutlineNode{...}
repo.Create(ctx, outline)
// outline.ID 此时仍然是零值！

repo.Update(ctx, outline)  // 💥 失败：找不到记录（ID 为空）
```

## 根本原因

`InsertOne` 返回的 `result.InsertedID` 是新插入文档的 ID，但没有将其设置回传入的对象。Go 语言中参数是值传递，除非直接修改对象的字段，否则调用方的对象不会被修改。

## 修复方案

### 正确的实现

```go
func (r *SomeRepository) Create(ctx context.Context, model *Model) error {
    model.TouchForCreate()

    result, err := r.GetCollection().InsertOne(ctx, model)
    if err != nil {
        return errors.NewRepositoryError(errors.RepositoryErrorInternal, "create failed", err)
    }

    // ✅ 将生成的 ID 回设到模型
    if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
        model.ID = oid
    }

    return nil
}
```

## 需要检查的 Repository 列表

### Writer 模块
- [x] `outline_repository_mongo.go` - 已修复
- [x] `project_repository_mongo.go`
- [x] `document_repository_mongo.go`
- [ ] `batch_operation_repository_mongo.go` - 需检查
- [x] `timeline_repository_mongo.go`
- [x] `version_service.go` 相关 `file_revisions/file_patches/commits` 写入回设

### Social 模块
- [x] `booklist_repository_mongo.go`
- [x] `comment_repository_mongo.go`
- [x] `review_repository_mongo.go`
- [x] `like_repository_mongo.go`
- [x] `follow_repository_mongo.go`
- [x] `message_repository_mongo.go`

### Reader 模块
- [x] `reading_progress_repository_mongo.go`
- [x] `reading_settings_repository_mongo.go`
- [x] `reader_theme_repository_mongo.go`
- [ ] `collection_repository_mongo.go`
- [ ] `comment_repository_mongo.go`

### Auth 模块
- [x] `role_repository_mongo.go`
- [x] `permission_repository_mongo.go`
- [ ] `user_repository_mongo.go`

### Bookstore 模块
- [x] `chapter_repository_mongo.go`
- [x] `category_repository_mongo.go`

### 其他模块
- [x] `finance/wallet_repository_mongo.go` - `wallets/transactions/withdraw_requests`
- [x] `notification/notification_repository_impl.go`
- [x] `notification/preference_repository_impl.go`
- [x] `notification/push_device_repository_impl.go`
- [x] `notification/template_repository_impl.go`
- [ ] `message_repository_mongo.go`
- [ ] 所有其他 Repository

## 检查脚本

可以使用以下脚本快速扫描问题：

```bash
# 搜索所有 InsertOne 但未设置 ID 的情况
grep -r "InsertOne" --include="*_mongo.go" repository/ | \
  grep -v "model.ID = " | \
  grep -v "InsertedID"
```

## 测试验证

每个修复后应添加测试验证：

```go
func TestRepository_Create_SetsID(t *testing.T) {
    repo, ctx, cleanup := setupRepo(t)
    defer cleanup()

    model := &Model{Name: "test"}
    err := repo.Create(ctx, model)
    require.NoError(t, err)

    // 验证 ID 已被设置
    assert.NotEqual(t, primitive.NilObjectID, model.ID)

    // 验证可以使用该 ID 查询
    found, err := repo.GetByID(ctx, model.ID.Hex())
    require.NoError(t, err)
    assert.Equal(t, model.Name, found.Name)
}
```

## 影响评估

| 严重性 | 影响范围 | 说明 |
|--------|----------|------|
| 高 | 数据完整性 | 创建后无法立即更新 |
| 中 | API 兼容性 | 前端可能获取不到 ID |
| 低 | 性能 | 需要额外查询获取 ID |

## 相关代码示例

### 正确示例（参考）

```go
// repository/mongodb/social/booklist_repository_mongo.go
func (r *MongoBookListRepository) CreateBookList(ctx context.Context, bookList *social.BookList) error {
    bookList.CreatedAt = time.Now()
    bookList.UpdatedAt = time.Now()

    result, err := r.GetCollection().InsertOne(ctx, bookList)
    if err != nil {
        return fmt.Errorf("创建书单失败: %w", err)
    }

    // ✅ 正确回设 ID
    if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
        bookList.ID = oid
    }

    return nil
}
```

## 修复优先级

1. **立即修复**: Writer 模块剩余 Repository
2. **高优先级**: Social、Auth 模块
3. **中优先级**: Reader、Notification 模块
4. **低优先级**: 其他辅助模块

## 验证清单

每个 Repository 修复后：
- [ ] 单元测试通过
- [ ] 集成测试通过
- [ ] API 测试验证返回 ID
- [ ] 代码审查通过

## 相关Issue

### 依赖Issue（必须先处理）
- [#001: 统一模型层 ID 字段类型](./001-unify-id-type-in-models.md) - ⚠️ 需要先统一ID类型为ObjectID，Create方法才能正确回设ID

### 相关Issue（联合处理）
- [#010: Repository 层业务逻辑渗透](./010-repository-layer-business-logic-permeation.md) - Create方法问题属于Repository层职责范围
- [#013: 测试用户种子数据ID未设置问题](./013-test-user-seed-id-not-set.md) - 种子数据Create也需要正确回设ID

### 依赖关系
- ID类型统一（#001）是Create方法修复（#002）的前提
- Create方法修复（#002）影响种子数据创建（#013）

## 参考链接

- [MongoDB InsertOne 文档](https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Collection.InsertOne)

