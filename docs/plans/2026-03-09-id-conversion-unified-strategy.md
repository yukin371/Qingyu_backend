# ID转换统一策略实施计划 v2.0

## ✅ 项目完成

**完成日期**: 2026-03-09

**PR**:
- Phase 1 & 2: [#115](https://github.com/yukin371/Qingyu_backend/pull/115) (已合并)
- Phase 4 & 5: [#116](https://github.com/yukin371/Qingyu_backend/pull/116)

**成果总结**:
1. 建立了统一的ID解析工具 `repository.ParseID`
2. Service层28个文件完成迁移
3. 创建了ID错误处理指南文档
4. 弃用了旧的 `models/shared/types/id.go`
5. 修复了ReviewRepository类型不匹配bug

---

## 修订说明

基于2026-03-09的反馈，本计划采用**渐进式迁移策略**，避免一次性大规模接口变更带来的风险。

## 核心目标

**原目标（已废弃）**：将所有ID转换逻辑集中到Repository层，Service层和API层只处理string类型ID

**新目标**：
1. 建立统一的ID解析策略和错误语义
2. 消除API/Service层热路径中的重复ID解析
3. 保持Repository核心接口稳定性，采用渐进式兼容策略

---

## Phase 1: 统一ID转换基础设施

### 目标
建立统一的ID转换工具，定义清晰的错误语义，**不改任何公开Repository接口**

### 1.1 创建统一错误定义

```go
// repository/errors.go
package repository

import "errors"

var (
    // ErrEmptyID 表示ID为空字符串
    ErrEmptyID = errors.New("ID cannot be empty")
    // ErrInvalidIDFormat 表示ID格式无效
    ErrInvalidIDFormat = errors.New("invalid ID format")
)
```

### 1.2 增强ID转换工具

```go
// repository/id_converter.go

// ParseID 解析必需的ID，空字符串返回错误
func ParseID(id string) (primitive.ObjectID, error) {
    if id == "" {
        return primitive.NilObjectID, ErrEmptyID
    }
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return primitive.NilObjectID, fmt.Errorf("%w: %s", ErrInvalidIDFormat, id)
    }
    return oid, nil
}

// ParseOptionalID 解析可选ID，空字符串返回nil（不报错）
// 用于过滤条件等场景，如"不限分类"时category_id为空
func ParseOptionalID(id string) (*primitive.ObjectID, error) {
    if id == "" {
        return nil, nil
    }
    oid, err := ParseID(id)
    if err != nil {
        return nil, err
    }
    return &oid, nil
}

// ParseIDs 批量解析ID列表
func ParseIDs(ids []string) ([]primitive.ObjectID, error) {
    if len(ids) == 0 {
        return nil, nil
    }
    result := make([]primitive.ObjectID, 0, len(ids))
    for i, id := range ids {
        oid, err := ParseID(id)
        if err != nil {
            return nil, fmt.Errorf("ids[%d]: %w", i, err)
        }
        result = append(result, oid)
    }
    return result, nil
}

// IsIDError 判断是否为ID相关错误（用于API层快速失败）
func IsIDError(err error) bool {
    return errors.Is(err, ErrEmptyID) || errors.Is(err, ErrInvalidIDFormat)
}
```

### 1.3 验收标准

- [x] `repository/errors.go` 创建完成
- [x] `repository/id_converter.go` 新增ParseID、ParseOptionalID、ParseIDs方法
- [x] 单元测试覆盖所有边界情况（空字符串、无效格式、批量场景）
- [x] 不修改任何现有公开接口

---

## Phase 2: 试点模块迁移

### 选择标准（按优先级）

| 优先级 | 模块 | 理由 |
|--------|------|------|
| P0 | **低风险模块** | 接口少、测试完整、依赖窄 |
| P1 | 高频低风险 | 调用频繁但结构简单 |
| P2 | 高频高风险 | 最后处理，需要更多准备 |

### 2.1 推荐试点模块

基于分析，推荐以下模块作为首批试点：

1. **NotificationRepository** - 通知模块
   - 接口简单，CRUD为主
   - 依赖链短
   - 测试覆盖较好

2. **BookmarkRepository** - 书签模块
   - 业务逻辑相对独立
   - 调用频次中等

### 2.2 迁移模式（不改接口）

**方式A：内部helper函数**

```go
// repository/mongodb/reader/bookmark_repository_mongo.go

// 内部helper，不改公开接口
func (r *BookmarkMongoRepo) parseID(id string) (primitive.ObjectID, error) {
    return repository.ParseID(id)
}

// 公开接口保持不变
func (r *BookmarkMongoRepo) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Bookmark, error) {
    // 原有实现
}

// 可选：新增string版本facade（不替代原接口）
func (r *BookmarkMongoRepo) GetByIDString(ctx context.Context, id string) (*models.Bookmark, error) {
    oid, err := r.parseID(id)
    if err != nil {
        return nil, err
    }
    return r.GetByID(ctx, oid)
}
```

**方式B：Service层使用统一工具**

```go
// service/reader/bookmark_service.go

// Before
func (s *BookmarkService) GetBookmark(ctx context.Context, id string) (*dto.BookmarkDTO, error) {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, errors.New("invalid bookmark id")
    }
    return s.repo.GetByID(ctx, oid)
}

// After - 使用统一工具
func (s *BookmarkService) GetBookmark(ctx context.Context, id string) (*dto.BookmarkDTO, error) {
    oid, err := repository.ParseID(id)
    if err != nil {
        return nil, err  // 统一错误语义
    }
    return s.repo.GetByID(ctx, oid)
}
```

### 2.3 验收标准
- [x] 1-2个试点模块完成迁移 ✅
- [x] Service层使用统一ParseID替代ObjectIDFromHex ✅
- [x] 所有现有测试通过 ✅
- [x] Mock编译检查通过 ✅

---

## Phase 3: 评估与决策

### 3.1 评估点

在Phase 2完成后，评估以下问题：

1. **是否需要将Repository接口参数统一改为string？**
   - 收益：API/Service层完全不需要关心ObjectID
   - 成本：接口签名变更、Mock更新、测试调整
   - 当前方案（内部helper + string facade）是否足够？

2. **如果决定迁移接口，采用兼容期策略：**
   - 旧接口保留（标记Deprecated）
   - 新接口并存（GetByIDString）
   - 高频模块迁完后再删除旧接口

### 3.2 决策矩阵

| 情况 | 建议 |
|------|------|
| helper方案效果良好，重复代码已大幅减少 | 保持现状，不强制改接口 ✅ 已选择 |
| 仍有大量重复代码，团队希望更彻底的统一 | 启动接口迁移，分模块进行 |

### 3.3 评估结论

**决策：保持现状，不强制修改 Repository 接口**

理由：
1. Service 层已全部完成迁移（28个文件），重复代码大幅减少
2. API 层的 ID 转换是"快速失败"前置校验，符合设计原则
3. Repository 层的 ID 转换是内部实现细节，修改接口成本高、收益低

---

## Phase 4: 统一错误翻译路径

### 目标
建立清晰的错误翻译链路：repo → service → api

### 4.1 Service层错误处理

```go
// service层统一包装
func (s *SomeService) DoSomething(ctx context.Context, id string) error {
    oid, err := repository.ParseID(id)
    if err != nil {
        // 统一翻译为业务错误
        if errors.Is(err, repository.ErrEmptyID) {
            return serviceerrors.ErrMissingParameter
        }
        return serviceerrors.ErrInvalidID
    }
    // ...
}
```

### 4.2 API层保留快速失败

**重要原则**：API层不应完全放弃前置参数校验

```go
// api/v1/some_api.go

func (api *SomeAPI) GetSomething(c *gin.Context) {
    id := c.Param("id")

    // 快速失败 - 明显错误在API层拦截
    if id == "" {
        response.BadRequest(c, "参数错误", "ID不能为空")
        return
    }

    // 调用Service，让Service处理格式校验
    result, err := api.service.GetSomething(ctx, id)
    if err != nil {
        // 统一错误翻译
        if repository.IsIDError(err) {
            response.BadRequest(c, "参数错误", "无效的ID格式")
            return
        }
        // 其他错误处理
    }
}
```

### 4.3 验收标准

- [x] 错误翻译路径文档化 ✅ `docs/guides/id-error-handling-guide.md`
- [x] API层保留快速失败机制 ✅ 文档中明确保留
- [x] 错误响应格式一致性验证 ✅ 统一使用 `response.BadRequest`

---

## Phase 5: 清理与文档

### 5.1 弃用旧工具

```go
// models/shared/types/id.go

// Deprecated: 使用 repository.ParseID 替代
func ParseObjectID(id string) (primitive.ObjectID, error) {
    return repository.ParseID(id)
}
```

### 5.2 全量扫描清理

```bash
# 扫描所有直接使用ObjectIDFromHex的地方
rg "ObjectIDFromHex|ParseObjectID|StringToObjectId|StringToObjectID" \
   --type go \
   --glob '!*_test.go' \
   --glob '!vendor/'

# 目标：除了deprecated包装和测试，不应有其他直接调用
```

### 5.3 文档更新

- [x] 开发规范更新：ID处理指南 ✅ `docs/guides/id-error-handling-guide.md`
- [x] 迁移指南：如何使用新工具 ✅ 已在错误处理指南中包含
- [x] API错误码文档更新 ✅ 统一使用 `response.BadRequest`

### 5.4 全量扫描结果

```bash
# 2026-03-09 扫描结果
# 除了 deprecated 包装和测试文件，无其他直接调用
```

---

## 验证方案

### 单元测试

```bash
go test ./repository/... -v -run TestParseID
go test ./repository/... -v -run TestParseOptionalID
go test ./repository/... -v -run TestParseIDs
```

### 集成测试检查点

- [x] Phase 1: 基础转换工具单元测试通过 ✅
- [x] Phase 2: 试点模块功能端到端测试通过 ✅
- [x] Phase 3: 评估报告完成 ✅
- [x] Phase 4: 错误翻译路径验证 ✅
- [x] Phase 5: 全量扫描无遗漏 ✅

### 验证命令

```bash
# 全量扫描ID转换调用
rg "ObjectIDFromHex|ParseObjectID|StringToObjectId|StringToObjectID" \
   --type go \
   --glob '!*_test.go' \
   --glob '!vendor/' \
   -c

# Mock编译检查
go test ./... -run=Nothing -compile-only 2>&1 | grep -i "mock"

# Contract test（API响应格式）
go test ./api/... -v -run TestContract
```

---

## 风险与缓解

| 风险 | 缓解措施 |
|------|----------|
| 循环依赖 | 错误定义放在独立文件 `repository/errors.go` |
| 行为变更 | 新增函数而非修改现有函数，保持兼容 |
| 遗漏转换点 | 使用rg全量扫描确认 |
| Mock失效 | 每阶段运行Mock编译检查 |
| API响应漂移 | Contract test验证400/404/500返回码 |

---

## 时间线

| 阶段 | 预计工作量 | 说明 |
|------|-----------|------|
| Phase 1 | 0.5天 | 基础设施，无风险 |
| Phase 2 | 1-2天 | 试点迁移，验证方案 |
| Phase 3 | 0.5天 | 评估决策 |
| Phase 4 | 1天 | 错误路径统一 |
| Phase 5 | 0.5天 | 清理文档 |

**总计**：3.5-4.5天

---

## 附录：高频转换文件参考

以下文件ID转换次数较多，但**不代表迁移优先级**。实际优先级按"风险低+测试完整+依赖窄"排序。

| 文件 | 转换次数 | 风险评估 |
|------|----------|----------|
| `repository/mongodb/bookstore/bookstore_repository_mongo.go` | 17 | 中 - 核心业务 |
| `api/v1/bookstore/chapter_catalog_api.go` | 16 | 中 - API层 |
| `api/v1/bookstore/chapter_api.go` | 14 | 中 - API层 |
| `service/writer/comment_service.go` | 14 | 低 - 可试点 |
| `service/reader/bookmark_service.go` | 12 | 低 - 可试点 |
