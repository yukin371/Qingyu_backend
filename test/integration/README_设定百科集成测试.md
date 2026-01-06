# 设定百科系统集成测试指南

**最后更新**: 2025-10-28  
**测试覆盖**: Character、Location、Timeline 完整测试

---

## 📋 测试文件清单

### 1. 单元测试（按功能模块）

| 测试文件 | 测试内容 | 状态 |
|---------|---------|------|
| `scenario_character_test.go` | 角色管理完整流程 | ✅ |
| `scenario_location_test.go` | 地点管理完整流程 | ✅ |
| `scenario_timeline_test.go` | 时间线管理完整流程 | ✅ |

### 2. 端到端测试

| 测试文件 | 测试内容 | 状态 |
|---------|---------|------|
| `writer_encyclopedia_e2e_test.go` | 设定百科完整创作流程 | ✅ |
| `scenario_writing_test.go` | 基础写作流程 | ✅ |
| `scenario_reading_test.go` | 阅读流程 | ✅ |

---

## 🎯 测试覆盖范围

### Character（角色）系统测试

**测试文件**: `scenario_character_test.go`

**测试场景**:
1. ✅ 创建测试项目
2. ✅ 创建第一个角色（李逍遥）
3. ✅ 创建第二个角色（赵灵儿）
4. ✅ 获取角色列表
5. ✅ 获取角色详情
6. ✅ 更新角色信息
7. ✅ 创建角色关系（恋人）
8. ✅ 获取角色关系图
9. ✅ 删除角色关系
10. ✅ 删除角色

**API端点覆盖**:
- ✅ POST   /api/v1/projects/:projectId/characters
- ✅ GET    /api/v1/projects/:projectId/characters
- ✅ GET    /api/v1/characters/:characterId
- ✅ PUT    /api/v1/characters/:characterId
- ✅ DELETE /api/v1/characters/:characterId
- ✅ POST   /api/v1/characters/relations
- ✅ GET    /api/v1/projects/:projectId/characters/graph
- ✅ DELETE /api/v1/characters/relations/:relationId

---

### Location（地点）系统测试

**测试文件**: `scenario_location_test.go`

**测试场景**:
1. ✅ 创建测试项目
2. ✅ 创建顶层地点（修真大陆）
3. ✅ 创建子地点（东部仙域）
4. ✅ 创建三级地点（天剑宗）
5. ✅ 获取地点列表
6. ✅ 获取地点层级树
7. ✅ 获取地点详情
8. ✅ 更新地点信息
9. ✅ 创建地点关系
10. ✅ 删除地点关系
11. ✅ 删除地点

**API端点覆盖**:
- ✅ POST   /api/v1/projects/:projectId/locations
- ✅ GET    /api/v1/projects/:projectId/locations
- ✅ GET    /api/v1/projects/:projectId/locations/tree
- ✅ GET    /api/v1/locations/:locationId
- ✅ PUT    /api/v1/locations/:locationId
- ✅ DELETE /api/v1/locations/:locationId
- ✅ POST   /api/v1/locations/relations
- ✅ DELETE /api/v1/locations/relations/:relationId

**特色功能测试**:
- ✅ 层级树构建（3级：大陆→区域→城市）
- ✅ 父子关系验证

---

### Timeline（时间线）系统测试

**测试文件**: `scenario_timeline_test.go`

**测试场景**:
1. ✅ 创建测试项目
2. ✅ 创建时间线（主线剧情）
3. ✅ 获取时间线列表
4. ✅ 获取时间线详情
5. ✅ 创建第一个事件（主角出生）
6. ✅ 创建第二个事件（拜师学艺）
7. ✅ 获取事件列表
8. ✅ 获取事件详情
9. ✅ 更新事件信息
10. ✅ 获取可视化数据
11. ✅ 删除事件
12. ✅ 删除时间线

**API端点覆盖**:
- ✅ POST   /api/v1/projects/:projectId/timelines
- ✅ GET    /api/v1/projects/:projectId/timelines
- ✅ GET    /api/v1/timelines/:timelineId
- ✅ DELETE /api/v1/timelines/:timelineId
- ✅ POST   /api/v1/timelines/:timelineId/events
- ✅ GET    /api/v1/timelines/:timelineId/events
- ✅ GET    /api/v1/timeline-events/:eventId
- ✅ PUT    /api/v1/timeline-events/:eventId
- ✅ DELETE /api/v1/timeline-events/:eventId
- ✅ GET    /api/v1/timelines/:timelineId/visualization

**特色功能测试**:
- ✅ 事件类型验证（character, milestone）
- ✅ 时间线可视化数据生成
- ✅ 事件关联（角色、地点）

---

### 端到端测试

**测试文件**: `writer_encyclopedia_e2e_test.go`

**测试流程**:
```
项目创建
   ↓
角色创建（云无极）
   ↓
地点创建（天剑宗）
   ↓
时间线创建（主角成长线）
   ↓
事件创建（拜入天剑宗）
   ↓
章节创建（第一章 拜师）
   ↓
完整性验证
```

**验证项**:
- ✅ 所有设定数据可正常创建
- ✅ 角色、地点、时间线数据关联
- ✅ 可视化数据生成正常
- ✅ 各模块数据一致性

---

## 🚀 运行测试

### 运行所有设定百科测试

```bash
# 进入测试目录
cd test/integration

# 运行所有测试
go test -v -run "Character|Location|Timeline|Encyclopedia"

# 或单独运行
go test -v -run TestCharacterScenario
go test -v -run TestLocationScenario
go test -v -run TestTimelineScenario
go test -v -run TestWriterEncyclopediaE2E
```

### 运行指定测试

```bash
# 只测试角色系统
go test -v -run TestCharacterScenario

# 只测试地点系统
go test -v -run TestLocationScenario

# 只测试时间线系统
go test -v -run TestTimelineScenario

# 只测试端到端
go test -v -run TestWriterEncyclopediaE2E
```

### 跳过长测试

```bash
# 使用 -short 跳过集成测试
go test -v -short
```

---

## 📊 测试统计

### 测试覆盖率

| 系统 | API端点数 | 测试场景数 | 覆盖率 |
|------|----------|----------|--------|
| Character | 8 | 10 | 100% |
| Location | 8 | 11 | 100% |
| Timeline | 10 | 12 | 100% |
| **总计** | **26** | **33** | **100%** |

### 测试数据

**测试用例总数**: 33个  
**API端点覆盖**: 26个  
**完整流程测试**: 4个  
**平均测试时长**: ~10秒/文件

---

## 🎯 测试最佳实践

### 1. 使用 TestHelper

```go
helper := NewTestHelper(t, router)

// 登录
token := helper.LoginTestUser()

// 发送请求
w := helper.DoAuthRequest("POST", "/api/v1/projects/:id/characters", data, token)

// 验证响应
data := helper.AssertSuccess(w, 201, "创建应该成功")

// 记录日志
helper.LogSuccess("创建成功 - ID: %s", id)
```

### 2. 测试数据清理

每个测试应该：
- ✅ 创建独立的测试数据
- ✅ 使用唯一标识符（时间戳）
- ✅ 测试结束后清理（通过defer cleanup()）

### 3. 错误处理

```go
if projectID == "" {
    t.Skip("无法创建项目，跳过后续测试")
}

// 继续测试...
```

### 4. 逐步验证

```go
t.Run("步骤1：创建", func(t *testing.T) { ... })
t.Run("步骤2：验证", func(t *testing.T) { ... })
t.Run("步骤3：删除", func(t *testing.T) { ... })
```

---

## 🐛 常见问题

### Q1: 测试失败：404 Not Found

**原因**: 路由未注册或ServiceContainer未初始化

**解决**:
1. 检查 `router/enter.go` 是否注册了writer路由
2. 检查 ServiceContainer 是否注册了相关服务
3. 检查 main.go 的服务初始化

### Q2: 测试失败：无法登录

**原因**: 测试数据库中没有测试用户

**解决**:
```bash
# 运行测试数据准备脚本
go run cmd/prepare_test_data/main.go
```

### Q3: 测试超时

**原因**: 数据库连接或服务初始化慢

**解决**:
```go
// 增加超时时间
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
```

---

## 📝 测试报告示例

```
=== RUN   TestCharacterScenario
=== RUN   TestCharacterScenario/1.创建测试项目
✓ 项目创建成功 - ID: 6543210abc
=== RUN   TestCharacterScenario/2.创建第一个角色
✓ 角色创建成功 - ID: char_001, 名字: 李逍遥
=== RUN   TestCharacterScenario/3.创建第二个角色
✓ 角色创建成功 - ID: char_002
...
--- PASS: TestCharacterScenario (2.34s)
    --- PASS: TestCharacterScenario/1.创建测试项目 (0.12s)
    --- PASS: TestCharacterScenario/2.创建第一个角色 (0.23s)
    ...
PASS
ok      Qingyu_backend/test/integration    2.456s
```

---

## ✅ 验收标准

### 功能完整性
- ✅ 所有API端点都有测试覆盖
- ✅ 创建、读取、更新、删除操作全部测试
- ✅ 关系管理功能测试完整
- ✅ 可视化数据生成测试通过

### 测试质量
- ✅ 测试用例独立运行
- ✅ 测试数据隔离
- ✅ 错误场景处理
- ✅ 清晰的测试日志

### 覆盖率
- ✅ API端点覆盖率: 100%
- ✅ 核心业务流程覆盖率: 100%
- ✅ 边界条件测试: 完成

---

## 🔗 相关文档

- [测试Helper使用指南](./README_TestHelper使用指南.md)
- [集成测试说明](./README_集成测试说明.md)
- [测试运行指南](../README_测试运行指南.md)
- [阶段3前置准备实施报告](../../doc/implementation/00进度指导/阶段3前置准备_实施报告.md)

---

**最后更新**: 2025-10-28  
**维护者**: Qingyu Test Team  
**测试状态**: ✅ 全部通过

