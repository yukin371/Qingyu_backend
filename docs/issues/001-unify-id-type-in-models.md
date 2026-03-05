# Issue #001: 统一模型层 ID 字段类型

**优先级**: 高 (P0)
**类型**: 技术债务
**状态**: 待处理
**创建日期**: 2026-03-05
**相关报告**: [Writer DTO 重构总结报告](../reports/2026-03-05-dto-refactoring-summary.md#21-id-类型不一致问题高优先级)

---

## 问题描述

项目中存在 ID 字段类型不统一的问题，约 37 个模型使用 `string` 类型，约 176 个模型使用 `primitive.ObjectID` 类型。这种不一致导致：

1. 查询失败：使用 ObjectID 查询 string ID 的记录会找不到
2. 类型转换开销：需要在多处进行 `.Hex()` 和 `ObjectIDFromHex()` 转换
3. 代码混乱：开发者不清楚应该使用哪种类型

## 影响范围

### 受影响的模型（string ID）
- `models/auth/permission_template.go` - PermissionTemplate.ID
- `models/social/*` - 部分模型
- 其他约 30+ 个模型

### 具体问题示例

```go
// 不一致的示例：
type PermissionTemplate struct {
    ID string `bson:"_id,omitempty" json:"id"`  // string 类型
}

type BookList struct {
    ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`  // ObjectID 类型
}

// 导致 Repository 层查询逻辑不一致：
func (r *PermissionTemplateRepository) GetByID(id string) {
    // 需要使用 string 查询
    filter := bson.M{"_id": exactMatchRegex(id)}
}

func (r *BookListRepository) GetByID(id string) {
    // 需要转换为 ObjectID 查询
    objectID, _ := primitive.ObjectIDFromHex(id)
    filter := bson.M{"_id": objectID}
}
```

## 推荐解决方案

### 统一原则

```
┌─────────────┐
│   API层      │ → 统一使用 string
├─────────────┤
│   DTO层      │ → 统一使用 string
├─────────────┤
│  Service层   │ → 使用 DTO string，内部转换
├─────────────┤
│ Repository层 │ → 负责类型转换
├─────────────┤
│  Model层     │ → 统一使用 primitive.ObjectID
└─────────────┘
```

### 转换模式

```go
// Model 层（统一使用 ObjectID）
type User struct {
    ID       primitive.ObjectID `bson:"_id,omitempty"`
    Username string             `bson:"username"`
}

// DTO 层（统一使用 string）
type UserResponse struct {
    ID       string `json:"id"`
    Username string `json:"username"`
}

// Converter 负责转换
func ToUserResponse(user *User) UserResponse {
    return UserResponse{
        ID:       user.ID.Hex(),        // ObjectID → string
        Username: user.Username,
    }
}

func ToUserID(id string) (primitive.ObjectID, error) {
    return primitive.ObjectIDFromHex(id)  // string → ObjectID
}
```

## 实施计划

### Phase 1: 评估和准备
- [ ] 使用工具扫描所有模型，统计 string ID 的使用情况
- [ ] 分析每个模型的依赖关系
- [ ] 制定详细的迁移计划

### Phase 2: 逐模块迁移
按以下顺序迁移模块（从依赖少的开始）：

1. [ ] `models/auth/` - PermissionTemplate, Role, Permission
2. [ ] `models/social/` - BookListLike, Comment, Review
3. [ ] `models/writer/` - 检查是否有 string ID
4. [ ] 其他模块

每个模块迁移步骤：
1. 修改 Model 定义：`ID string` → `ID primitive.ObjectID`
2. 修改 Repository 查询逻辑
3. 更新 Converter 函数
4. 运行测试验证
5. 提交并合并

### Phase 3: 清理和验证
- [ ] 移除过时的 ID 转换代码
- [ ] 更新文档
- [ ] 全量测试验证

## 风险评估

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 数据迁移失败 | 高 | 先在测试环境验证，做好备份 |
| API 兼容性 | 中 | DTO 层保持 string，对外透明 |
| 性能影响 | 低 | ObjectID 查询性能更好 |

## 相关Issue

### 依赖Issue（必须先处理）
- 无

### 相关Issue（联合处理）
- [#002: Repository Create 方法未回设 ID](./002-create-method-id-not-set-bug.md) - ID类型统一后需要确保Create方法正确回设ID
- [#011: 前后端数据类型不一致](./011-frontend-backend-data-type-inconsistency.md) - 包含ID类型转换边界不清晰问题
- [#013: 测试用户种子数据ID未设置问题](./013-test-user-seed-id-not-set.md) - 种子数据需要使用正确的ObjectID类型

### 架构相关
- [#010: Repository 层业务逻辑渗透](./010-repository-layer-business-logic-permeation.md) - Repository层重构时需要确保ID类型正确

## 相关代码

需要修改的文件示例：
- `models/auth/permission_template.go`
- `repository/mongodb/auth/permission_template_repository_mongo.go`
- 其他使用 string ID 的模型和 Repository

## 预期收益

1. **一致性**: 所有模型使用相同的 ID 类型
2. **性能**: ObjectID 索引查询更高效
3. **简化**: Repository 层查询逻辑统一
4. **可维护**: 新开发者更容易理解

## 参考链接

- [MongoDB ObjectId 规范](https://www.mongodb.com/docs/manual/reference/method/ObjectId/)
- [Go MongoDB Driver 类型转换](https://pkg.go.dev/go.mongodb.org/mongo-driver/bson/primitive)
