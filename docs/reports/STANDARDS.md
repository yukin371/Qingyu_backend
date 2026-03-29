# Qingyu 后端开发规范

> **重要提示**: 本仓库的规范文档已迁移至父仓库统一管理。

**当前状态**: 已迁移至父仓库
**迁移日期**: 2026-01-27
**新规范位置**: `../../docs/architecture/BACKEND_DEVELOPMENT_STANDARD.md`

---

## 📚 最新规范文档

### 统一后端开发规范 v2.0

**位置**: [父仓库 docs/architecture/BACKEND_DEVELOPMENT_STANDARD.md](../../docs/architecture/BACKEND_DEVELOPMENT_STANDARD.md)

**版本**: v2.0
**更新日期**: 2026-01-27
**基于**: 2026-01-26架构全面审查报告

**包含内容**:
1. ✅ 总览与原则
2. ✅ 架构设计规范（四层架构、依赖注入、接口分离）
3. ✅ API设计规范（URL、HTTP方法、响应格式、错误码）
4. ✅ 数据模型规范（ID类型统一、基础模型、共享类型）
5. ✅ 数据库规范（索引、缓存）
6. ✅ Service层规范（事务管理、事件驱动）
7. ✅ 中间件规范（目录结构、优先级系统、CORS修复）
8. ✅ 错误处理规范（统一错误码体系）
9. ✅ 测试规范（单元测试、集成测试）
10. ✅ 部署与运维规范（环境配置、健康检查）

**快速访问**:
- 📘 [在线阅读](../../docs/architecture/BACKEND_DEVELOPMENT_STANDARD.md)
- 📥 [下载PDF](../../docs/architecture/BACKEND_DEVELOPMENT_STANDARD.pdf) (如可用)

---

## 🗂️ 旧规范归档

本仓库的旧规范文档已归档至 `docs/standards/archive/` 目录，以下为归档清单：

### 已归档文档

| 文档名称 | 原位置 | 归档位置 | 归档原因 |
|---------|--------|---------|---------|
| API状态码规范 | `docs/architecture/api-status-code-standard.md` | `docs/standards/archive/api-status-code-standard.md` | 已整合到v2.0第3章 |
| ID类型统一标准 | `docs/architecture/id-type-unification-standard.md` | `docs/standards/archive/id-type-unification-standard.md` | 已整合到v2.0第4章 |
| RESTful API规范 | `docs/guides/standards/restful-api-design-standard.md` | （父仓库） | 已整合到v2.0第3章 |

### 归档说明

1. **为什么归档？**
   - 解决规范分散、重复、冲突的问题
   - 基于审查报告（P0/P1问题）进行更新
   - 统一全项目规范管理

2. **旧文档是否还有用？**
   - ✅ 有参考价值：保留了设计思路和示例代码
   - ⚠️ 不再更新：后续只更新父仓库的统一规范
   - 📌 历史参考：可作为理解设计意图的补充材料

3. **如何使用旧文档？**
   - 如需查看历史讨论，可参考归档文档
   - 开发时请以父仓库的v2.0规范为准
   - 如发现冲突，以v2.0为准

---

## 🔗 快速参考

### 核心规范要点

#### 1. API响应码（P0修复）
```go
// ✅ 正确
shared.Success(c, http.StatusOK, "操作成功", data)
// 返回: {"code": 0, "message": "success", ...}

// ❌ 错误（已修复）
// 返回: {"code": 200, ...}
```

#### 2. URL前缀（P0修复）
```
✅ /api/v1/users
✅ /api/v1/system/health
❌ /system/health（已修复）
```

#### 3. PATCH方法（P0修复）
```
PATCH /api/v1/users/{id}
Body: {"nickname": "新昵称"}  // 只更新提供的字段
```

#### 4. ID类型（P0修复）
```go
// Model层: primitive.ObjectID
type Book struct {
    ID primitive.ObjectID `bson:"_id"`
}

// Service层: string
type BookService interface {
    GetBook(ctx context.Context, bookID string) (*dto.BookDTO, error)
}

// API层: string
type BookDTO struct {
    ID string `json:"id"`
}
```

#### 5. 中间件目录（P0修复）
```
internal/api/middleware/
├── request_id.go
├── recovery.go
├── cors.go          # ⚠️ 必须在路由前
└── auth.go
```

---

## 📋 迁移检查清单

如果您正在使用旧的规范文档，请检查以下迁移事项：

### 代码层面
- [ ] 确认所有API响应使用 `code: 0` 表示成功
- [ ] 确认所有URL使用 `/api/v1` 前缀
- [ ] 为需要部分更新的端点添加PATCH方法
- [ ] 确认Model层使用 `primitive.ObjectID`
- [ ] 确认Service层接口使用 `string`
- [ ] 将中间件移至 `internal/api/middleware/`
- [ ] 确认CORS中间件在路由前注册
- [ ] 实现事务管理机制

### 文档层面
- [ ] 更新项目README，引用父仓库规范
- [ ] 更新人脸文档，指向新规范位置
- [ ] 将旧规范移至归档目录
- [ ] 通知团队成员规范已迁移

### 测试层面
- [ ] 更新测试用例以符合新规范
- [ ] 添加规范一致性测试
- [ ] 验证所有P0问题已修复

---

## 🤝 贡献指南

### 如何更新规范？

1. **小修改**：直接提交PR到父仓库
   - 修改 `docs/architecture/BACKEND_DEVELOPMENT_STANDARD.md`
   - 说明修改原因和影响范围

2. **大修改**：先讨论再实施
   - 在团队讨论群提出建议
   - 达成共识后创建Issue
   - 实施前更新版本号和变更日志

3. **冲突处理**：
   - 如发现规范内容冲突，以父仓库v2.0为准
   - 在团队群讨论解决
   - 更新后同步到后端仓库引用文档

---

## 📞 联系方式

如有疑问或建议，请联系：
- 技术负责人：[待填写]
- 架构组：[待填写]
- 或者直接在团队群讨论

---

**最后更新**: 2026-01-27
**维护者**: Qingyu 开发团队
**归档状态**: ✅ 已迁移至父仓库
