# Service层测试改进执行计划

## 🎯 项目目标

**目标**：在3个月内将Service层测试覆盖率从35%提升到75%，确保所有关键业务逻辑都有完整的单元测试。

**时间框架**：
- **第1周（P0）**：关键业务服务
- **第2周（P1）**：重要功能服务
- **第3周及以后（P2）**：次要服务和工具类

---

## 📊 当前状态快照

```
Service层统计：
├─ 总Service数量：~60个
├─ 已有测试：~21个 (35%)
├─ 需要测试：~39个 (65%)
│  ├─ 关键业务：~10个 (P0)
│  ├─ 重要功能：~15个 (P1)
│  └─ 次要服务：~14个 (P2)
└─ 预计工作量：60-80小时
```

---

## 📋 第1周行动计划（P0 - 关键业务）

**优先级**：🔴 最高 | **工作量**：~20小时 | **目标覆盖率**：85%+

### 1️⃣ UserService 改进

**文件**：`service/user/user_service.go`  
**现状**：✅ 有基础测试，但覆盖不完整

**改进任务**：
- [ ] 扩展 `user_service_test.go` 中的测试
- [ ] 新增 `test/service/user/password_validator_test.go`
- [ ] 新增 `test/service/user/password_reset_token_test.go`

**预期测试用例**（共15+）：
```
✓ Register - 成功注册
✓ Register - 邮箱已存在
✓ Register - 密码不符合要求
✓ Register - 参数验证失败
✓ Login - 成功登录
✓ Login - 密码错误
✓ Login - 用户不存在
✓ GetUserProfile - 成功获取
✓ GetUserProfile - 权限不足
✓ PasswordValidator - 各种密码规则
✓ PasswordResetToken - 生成、验证、过期
✓ 并发场景测试
```

**预计时间**：4小时  
**负责人**：@backend-team

---

### 2️⃣ ProjectService 改进

**文件**：`service/project/project_service.go`  
**现状**：✅ 有基础测试，需要扩展

**改进任务**：
- [ ] 扩展 `test/service/project/project_service_test.go`
- [ ] 新增 `test/service/project/autosave_service_test.go`
- [ ] 新增 `test/service/project/node_service_test.go`

**预期测试用例**（共20+）：
```
✓ CreateProject - 成功创建
✓ CreateProject - 项目名称验证
✓ UpdateProject - 成功更新
✓ UpdateProject - 权限检查
✓ UpdateProject - 版本冲突
✓ DeleteProject - 成功删除
✓ DeleteProject - 权限检查
✓ GetProject - 获取项目详情
✓ AutoSave - 启动自动保存
✓ AutoSave - 停止自动保存
✓ NodeService - 树形结构操作
✓ 并发项目操作
✓ 缓存一致性
```

**预计时间**：5小时  
**负责人**：@backend-team

---

### 3️⃣ DocumentService 测试（全新）

**文件**：`service/document/document_service.go`  
**现状**：❌ 无测试

**改进任务**：
- [ ] 新增 `test/service/document/document_service_test.go`

**预期测试用例**（共15+）：
```
✓ CreateDocument - 成功创建
✓ CreateDocument - 项目不存在
✓ CreateDocument - 文档名称重复
✓ UpdateDocument - 成功更新
✓ UpdateDocument - 权限检查
✓ DeleteDocument - 成功删除
✓ DeleteDocument - 权限检查
✓ GetDocumentContent - 获取内容
✓ SaveDocumentVersion - 版本保存
✓ 自动保存触发
✓ 并发编辑处理
```

**预计时间**：5小时  
**负责人**：@backend-team

---

### 4️⃣ AuthService 改进（认证模块）

**文件**：`service/auth/auth_service.go`  
**现状**：✅ 有基础测试，需要扩展

**改进任务**：
- [ ] 扩展 `service/auth/*_test.go`
- [ ] 新增 `service/auth/session_service_test.go`
- [ ] 加强 PermissionService 和 RoleService 测试

**预期测试用例**（共18+）：
```
✓ Authenticate - 成功认证
✓ Authenticate - 凭证无效
✓ TokenRefresh - 刷新Token
✓ TokenRevoke - 撤销Token
✓ SessionManagement - 会话管理
✓ PermissionCheck - 权限检查
✓ RoleInheritance - 角色继承
✓ RBAC - 基于角色的访问控制
✓ 并发会话处理
✓ 安全边界条件
```

**预计时间**：6小时  
**负责人**：@backend-team

---

### 📍 第1周检查点

**预期成果**：
- ✅ 4个关键Service全覆盖
- ✅ ~65个新测试用例
- ✅ 平均覆盖率 ≥ 85%
- ✅ 所有P0 Service测试通过

**验证命令**：
```bash
go test ./test/service/user/... -v -cover
go test ./test/service/project/... -v -cover
go test ./test/service/document/... -v -cover
go test ./service/auth/... -v -cover
```

---

## 📋 第2周行动计划（P1 - 重要功能）

**优先级**：🟠 高 | **工作量**：~25小时 | **目标覆盖率**：80%+

### 5️⃣ 书城模块 (Bookstore) 改进

**改进任务**：
- [ ] 新增 `test/service/bookstore/book_rating_service_test.go`
- [ ] 新增 `test/service/bookstore/book_statistics_service_test.go`
- [ ] 新增 `test/service/bookstore/chapter_service_test.go`
- [ ] 扩展 bookstore 现有测试

**关键测试**：评分计算、统计准确性、章节管理、缓存处理

**预计时间**：6小时

---

### 6️⃣ 阅读模块补充 (Reading) 改进

**改进任务**：
- [ ] 新增 `test/service/reading/reading_history_service_test.go`
- [ ] 加强缓存服务测试

**关键测试**：历史记录管理、缓存一致性、权限检查

**预计时间**：4小时

---

### 7️⃣ 存储模块 (Storage) 改进

**改进任务**：
- [ ] 扩展 `test/service/shared/storage/storage_service_test.go`
- [ ] 新增 `test/service/shared/storage/multipart_upload_service_test.go`
- [ ] 新增 `test/service/shared/storage/image_processor_test.go`

**关键测试**：上传处理、分片上传、图片处理、错误处理

**预计时间**：7小时

---

### 8️⃣ 钱包/支付模块 (Wallet) 改进

**改进任务**：
- [ ] 新增 `test/service/shared/wallet/transaction_service_test.go`
- [ ] 新增 `test/service/shared/wallet/withdraw_service_test.go`

**关键测试**：交易处理、余额验证、并发支付、撤销处理

**预计时间**：5小时

---

### 9️⃣ 审核模块 (Audit) 改进

**改进任务**：
- [ ] 新增 `test/service/audit/rule_engine_test.go`

**关键测试**：规则引擎、规则匹配、审核流程

**预计时间**：3小时

---

### 📍 第2周检查点

**预期成果**：
- ✅ 9个P1 Service全覆盖
- ✅ ~80个新测试用例
- ✅ 平均覆盖率 ≥ 80%

---

## 📋 第3周及以后（P2 - 次要服务）

**优先级**：🟡 中等 | **工作量**：~20小时 | **目标覆盖率**：70%+

### 1️⃣0️⃣ Writer模块

- [ ] `test/service/writer/character_service_test.go`
- [ ] `test/service/writer/location_service_test.go`
- [ ] `test/service/writer/timeline_service_test.go`

**预计时间**：6小时

---

### 1️⃣1️⃣ 缓存和统计

- [ ] 完善缓存服务测试
- [ ] 完善统计服务测试

**预计时间**：4小时

---

### 1️⃣2️⃣ 其他工具类

- [ ] 搜索服务测试
- [ ] 配置服务测试
- [ ] 通知服务补充

**预计时间**：5小时

---

## 🛠️ 技术规范和要求

### 测试结构要求

```go
func Test[ServiceName]_[MethodName](t *testing.T) {
    t.Run("Success", func(t *testing.T) { })        // 正常流程
    t.Run("InvalidInput", func(t *testing.T) { })   // 参数验证
    t.Run("NotFound", func(t *testing.T) { })       // 资源不存在
    t.Run("Unauthorized", func(t *testing.T) { })   // 权限检查
    t.Run("RepositoryError", func(t *testing.T) { }) // 数据库错误
    t.Run("ConcurrentAccess", func(t *testing.T) { }) // 并发场景
}
```

### 最低要求

| 指标 | 要求 |
|------|------|
| 每个Service的测试函数数 | ≥ 3个 |
| 每个方法的子测试数 | ≥ 5个 |
| 代码覆盖率 | ≥ 80% |
| Mock完整性 | 100%实现接口方法 |
| 并发测试 | ✅ 必须 |
| 权限测试 | ✅ 必须（若适用） |

### 代码审查清单

提交PR前检查：

- [ ] 所有测试都通过 (`go test ./test/service/...`)
- [ ] 覆盖率 ≥ 80% (`go test ./test/service/... -cover`)
- [ ] 没有linting错误 (`golangci-lint run`)
- [ ] Mock实现完整
- [ ] 有错误路径测试
- [ ] 有并发场景测试（如适用）
- [ ] 有权限检查测试（如适用）

---

## 📊 质量指标追踪

### 周进度表

| 周次 | 目标 | 计划Service数 | 预计覆盖率 | 实际进度 |
|------|------|-------|---------|---------|
| W1 | P0关键服务 | 4 | 85%+ | ⬜ 待开始 |
| W2 | P1重要功能 | 5 | 80%+ | ⬜ 待开始 |
| W3+ | P2次要服务 | 6 | 70%+ | ⬜ 待开始 |

### 整体进度

```
目前状态：35% (21/60服务有测试)
第1周目标：55% (33/60)
第2周目标：70% (42/60)
第3周目标：85%+ (51/60)
最终目标：90%+ (54/60)
```

---

## 🚀 快速启动命令

### 初始化P0测试框架

```bash
# 创建P0测试文件
mkdir -p test/service/user
mkdir -p test/service/project/
mkdir -p test/service/document/
mkdir -p service/auth/

# 查看现有测试
ls -la test/service/*/

# 运行所有测试
go test ./test/service/... -v

# 检查覆盖率
go test ./test/service/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### 验证单个模块

```bash
# UserService
go test ./test/service/user/... -v -cover

# ProjectService
go test ./test/service/project/... -v -cover

# DocumentService
go test ./test/service/document/... -v -cover

# AuthService
go test ./service/auth/... -v -cover
```

---

## 📚 相关资源

### 文档
- [Service层测试覆盖分析报告](SERVICE_TEST_COVERAGE_REPORT.md)
- [Service层测试改进快速指南](TESTING_IMPROVEMENT_GUIDE.md)
- [Service层架构规范](../../doc/architecture/架构设计规范.md)
- [软件工程规范](../../doc/engineering/软件工程规范_v2.0.md)

### 代码示例
- [现有测试参考](../../../test/service/)
- [Mock实现范例](../../../service/user/mocks/)

### 工具和技能
- 学习 testify/assert 和 testify/mock
- Go testing 最佳实践
- Mock对象设计模式

---

## 🎯 成功标准

### 第1周（P0）
- ✅ 4个关键Service完整测试
- ✅ 所有P0Service测试通过率 = 100%
- ✅ 平均代码覆盖率 ≥ 85%
- ✅ 0个pending/skip的测试

### 整体目标
- ✅ Service层整体覆盖率 ≥ 75%
- ✅ 所有关键业务逻辑都有测试
- ✅ 测试代码质量 = 源代码质量
- ✅ 建立持续的测试维护机制

---

## 🔄 维护机制

### 每周例会
- **时间**：周五下午
- **内容**：覆盖率进度检查、问题排查
- **输出**：进度报告、下周计划调整

### 持续改进
- **新Service必须有测试**：立即要求
- **PR审核检查**：测试覆盖率 ≥ 80%
- **月度总结**：记录改进数据、经验教训

### 工具和自动化
- 设置GitHub Actions自动运行测试
- 集成覆盖率报告到CI/CD
- 建立代码覆盖率基线

---

## 📞 联系和支持

- **技术问题**：后端架构团队
- **进度汇报**：@project-manager
- **代码审查**：@senior-backend-dev

---

**创建日期**：2025-10-31  
**最后更新**：2025-10-31  
**维护者**：后端架构团队  
**审批者**：待审批

