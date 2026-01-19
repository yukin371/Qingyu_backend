# version_service_test.go 测试修复报告

**修复日期**: 2025-10-17  
**相关任务**: 架构冗余字段修复  
**修复文件**: `service/project/version_service_test.go`

---

## 修复背景

在完成Document模型的架构重构后（移除Content和Version字段），测试文件仍然使用旧的模型结构，导致测试无法正确运行。

### 架构变更

**修改前**:
```go
type Document struct {
    ID        string
    ProjectID string
    Content   string  // ❌ 已移除
    Version   int     // ❌ 已移除
    // ...
}
```

**修改后**:
```go
// Document - 仅包含元数据
type Document struct {
    ID        string
    ProjectID string
    Title     string
    // ... (无Content和Version)
}

// DocumentContent - 包含实际内容和版本号
type DocumentContent struct {
    ID         string
    DocumentID string
    Content    string  // ✅ 内容在这里
    Version    int     // ✅ 版本号在这里
    // ...
}
```

---

## 修复内容

### 1. 更新测试数据准备

**修改前**：
```go
// 直接插入包含content和version的文档
_, err := fileCol().InsertOne(ctx, bson.M{
    "_id": id,
    "project_id": projectID,
    "node_id": nodeID,
    "content": "old",    // ❌ 错误
    "version": 1,        // ❌ 错误
})
```

**修改后**：
```go
// 分别插入元数据和内容
// 1. 插入文档元数据
docID := primitive.NewObjectID().Hex()
_, err := fileCol().InsertOne(ctx, bson.M{
    "_id": docID,
    "project_id": projectID,
    "node_id": nodeID,
    "title": "Test Document",
    "created_at": time.Now(),
    "updated_at": time.Now(),
})

// 2. 插入文档内容
_, err = contentCol().InsertOne(ctx, bson.M{
    "_id": primitive.NewObjectID().Hex(),
    "document_id": docID,
    "content": "old",    // ✅ 正确
    "version": 1,        // ✅ 正确
    "created_at": time.Now(),
    "updated_at": time.Now(),
})
```

### 2. 更新测试验证逻辑

**修改前**：
```go
// 从Document中读取content和version
var f model.Document
fileCol().FindOne(ctx, bson.M{"_id": id}).Decode(&f)
if f.Content != "expected" {  // ❌ Document没有Content字段
    t.Fatalf("content mismatch")
}
if f.Version != 2 {            // ❌ Document没有Version字段
    t.Fatalf("version mismatch")
}
```

**修改后**：
```go
// 从DocumentContent中读取content和version
var content model.DocumentContent
contentCol().FindOne(ctx, bson.M{"document_id": docID}).Decode(&content)
if content.Content != "expected" {  // ✅ 正确
    t.Fatalf("content mismatch")
}
if content.Version != 2 {           // ✅ 正确
    t.Fatalf("version mismatch")
}
```

### 3. 更新清理函数

**修改前**：
```go
func cleanupCollections(t *testing.T, projectID string) {
    ctx := context.Background()
    global.DB.Collection("novel_files").DeleteMany(ctx, bson.M{"project_id": projectID})
    global.DB.Collection("file_revisions").DeleteMany(ctx, bson.M{"project_id": projectID})
    global.DB.Collection("file_patches").DeleteMany(ctx, bson.M{"project_id": projectID})
}
```

**修改后**：
```go
func cleanupCollections(t *testing.T, projectID string) {
    ctx := context.Background()
    global.DB.Collection("novel_files").DeleteMany(ctx, bson.M{"project_id": projectID})
    global.DB.Collection("document_contents").DeleteMany(ctx, bson.M{}) // ✅ 新增
    global.DB.Collection("file_revisions").DeleteMany(ctx, bson.M{"project_id": projectID})
    global.DB.Collection("file_patches").DeleteMany(ctx, bson.M{"project_id": projectID})
}
```

---

## 修复的测试用例

### 1. TestUpdateContentWithVersion_HappyPath
- ✅ 分离Document和DocumentContent的插入
- ✅ 使用DocumentContent验证更新结果
- ✅ 验证版本号增加逻辑

### 2. TestUpdateContentWithVersion_Conflict
- ✅ 分离数据准备
- ✅ 测试版本冲突检测

### 3. TestRollbackToVersion
- ✅ 分离数据准备
- ✅ 使用DocumentContent验证回滚结果

### 4. TestCreateAndApplyPatch_Full
- ✅ 分离数据准备
- ✅ 使用DocumentContent验证补丁应用结果

---

## 测试数据库集合

修复后的测试涉及以下集合：

| 集合名 | 用途 | 存储内容 |
|-------|------|---------|
| `novel_files` | Document元数据 | 文档基本信息（标题、项目ID等） |
| `document_contents` | DocumentContent | 文档实际内容和版本号 |
| `file_revisions` | 版本历史 | 历史版本快照 |
| `file_patches` | 补丁记录 | 待应用的补丁 |

---

## 验证结果

✅ **编译检查**: `go test -c service/project` - 通过  
✅ **Linter检查**: 无错误  
✅ **测试逻辑**: 符合新架构设计  

---

## 架构一致性检查

修复后的测试代码符合以下架构原则：

1. ✅ **内容与元数据分离**: 测试分别操作两个集合
2. ✅ **版本控制在内容层**: Version字段在DocumentContent中
3. ✅ **Repository模式**: 使用`contentCol()`访问内容集合
4. ✅ **数据完整性**: 测试覆盖了版本冲突、回滚等场景

---

## 后续建议

### 1. 增强测试覆盖
建议添加以下测试场景：
- DocumentContent为空时的处理
- 并发更新冲突
- 大文档GridFS存储（当功能实现后）

### 2. Mock测试
对于单元测试，建议使用Mock Repository而非真实数据库：
```go
mockContentRepo := new(MockDocumentContentRepository)
mockContentRepo.On("GetByDocumentID", ...).Return(...)
```

### 3. 测试数据构建器
可以创建测试数据构建器简化测试代码：
```go
func BuildTestDocumentWithContent(projectID, nodeID string) (docID string, cleanup func()) {
    // 创建Document和DocumentContent
    // 返回ID和清理函数
}
```

---

## 相关文档

- **架构设计**: `doc/architecture/架构设计规范.md`
- **架构审计报告**: `doc/implementation/架构冗余字段审计报告.md`
- **Repository设计**: `doc/architecture/repository层设计规范.md`

---

**修复完成**: ✅  
**修复人**: AI架构修复系统  
**验证状态**: 已通过编译和linter检查

