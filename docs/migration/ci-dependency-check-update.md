# CI依赖检查规则更新总结

## 更新概述

本次更新为依赖检查工具添加了对auth模块迁移的支持，包括废弃路径检测和迁移指南。

## 更新时间

- **日期**: 2026-02-09
- **版本**: v1.1.0
- **状态**: 已完成并测试通过

## 主要变更

### 1. 新增功能

#### 废弃路径检测
- 检测已废弃的`Qingyu_backend/service/shared/auth`导入
- 区分生产代码（warning）和测试代码（deprecated/info）
- 新路径`Qingyu_backend/service/auth`不受影响

#### 新的严重级别
- `error`: 阻塞CI的错误
- `warning`: 警告但不阻塞CI
- `deprecated`: 信息提示（主要用于测试文件）

### 2. 代码变更

#### main.go
```go
// 新增废弃路径映射
var deprecatedImports = map[string]string{
    `Qingyu_backend/service/shared/auth`: `auth模块已迁移到service/auth，请使用新的导入路径`,
}

// 更新Violation结构
type Violation struct {
    File     string
    Line     int
    Import   string
    Rule     string
    Severity string // "error", "warning", or "deprecated"
}

// 新增检测逻辑
func checkImport(...) {
    // 规则0: 检查废弃路径
    if deprecationMsg, isDeprecated := deprecatedImports[importPath]; isDeprecated {
        // 返回warning或deprecated
    }
    // ... 其他规则
}
```

#### main_test.go
```go
// 新增测试用例
func TestDeprecatedImports(t *testing.T) {
    // 测试废弃路径检测
    // 测试测试文件的deprecated级别
    // 测试新路径不被标记
}

// 更新现有测试
func TestImportPatterns(t *testing.T) {
    // 添加auth路径相关测试
}
```

### 3. 文档更新

#### README.md
- 添加废弃路径警告章节
- 更新输出示例
- 添加迁移指南说明
- 更新扩展规则章节

#### 新增文档
- `docs/migration/auth-module-migration.md`: 完整的迁移指南
- `docs/migration/ci-dependency-check-update.md`: 本文档

## 验证结果

### 测试通过
```
=== RUN   TestImportPatterns
--- PASS: TestImportPatterns (0.00s)
=== RUN   TestDeprecatedImports
--- PASS: TestDeprecatedImports (0.00s)
PASS
ok      Qingyu_backend/scripts/check-dependencies    3.983s
```

### 实际运行验证
```bash
$ go run ./scripts/check-dependencies
🔍 检查代码依赖关系...
✅ 未发现依赖违规
```

### CI集成验证
- ✅ architecture-ci.yml已集成
- ✅ 检查命令正确
- ✅ 失败时不阻塞其他任务
- ✅ deprecated不影响CI通过

## 检查规则总结

### 当前检查规则

| 规则 | 目标 | 严重级别 | 说明 |
|------|------|----------|------|
| 废弃auth路径 | `Qingyu_backend/service/shared/auth` | warning (生产)<br>deprecated (测试) | 建议迁移到新路径 |
| 业务服务导入shared | `service/shared/*` | error | 应该使用Port接口 |

### 允许列表

以下模块可以导入`service/shared/*`：
- `service/container`
- `service/interfaces/shared`
- `router/shared`
- `api/v1/auth`
- `api/v1/shared`
- `realtime/websocket`
- `middleware`

### 禁止列表

以下模块不应该直接导入`service/shared/*`：
- `service/user`
- `service/writer`
- `service/reader`
- `service/ai`

## 使用指南

### 开发者

1. **日常开发**：依赖检查会在CI中自动运行
2. **本地检查**：`go run ./scripts/check-dependencies`
3. **迁移指导**：参考`docs/migration/auth-module-migration.md`

### CI/CD

- **自动运行**：每次push和PR都会触发检查
- **失败处理**：error和warning会失败，deprecated不影响
- **报告查看**：可在GitHub Actions中查看详细报告

## 后续计划

### 短期（已完成）
- ✅ 添加auth模块废弃路径检测
- ✅ 更新文档和测试
- ✅ 验证CI集成

### 中期
- 🔄 监控废弃路径使用情况
- 🔄 收集迁移反馈
- 🔄 计划其他shared子模块的迁移

### 长期
- 📋 移除兼容层（当所有迁移完成）
- 📋 扩展检查规则到其他模块
- 📋 添加循环依赖检测

## 兼容性

- **向后兼容**: ✅ 完全兼容
- **CI不中断**: ✅ deprecated不影响通过
- **兼容层保留**: ✅ 旧代码继续工作

## 相关链接

- [依赖检查工具](../../scripts/check-dependencies/README.md)
- [Auth模块迁移指南](./auth-module-migration.md)
- [架构设计文档](../architecture/dependency-rules.md)
- [CI配置](../../.github/workflows/architecture-ci.yml)

## 变更日志

### v1.1.0 (2026-02-09)
- ✅ 添加废弃路径检测功能
- ✅ 新增deprecated严重级别
- ✅ 更新测试用例
- ✅ 完善文档

### v1.0.0 (之前)
- ✅ 基础依赖检查功能
- ✅ shared模块导入检测
- ✅ Port接口模式验证

---

**维护者**: 架构团队
**审核**: 代码审查委员会
**状态**: 已发布
