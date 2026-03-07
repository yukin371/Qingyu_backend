# Issue #012: 401认证错误和权限配置问题

**优先级**: 高 (P1)
**类型**: 安全/认证问题
**状态**: ✅ 核心问题已解决（已归档）
**创建日期**: 2026-03-05
**归档日期**: 2026-03-06
**来源报告**: [401认证错误诊断报告](../../reports/archived/2026-02-08-auth-401-issue-diagnosis.md)

---

## 解决结果

本 Issue 对应的核心阻塞项已经关闭：

1. `author` 角色已补齐 `project:*` 与 `document:update/delete` 权限
2. `RBACChecker` 已能从 `configs/permissions.yaml` 真实加载角色、权限和继承关系
3. `internal/router/setup.go` 中遗留的空 `SetupAuthMiddleware` 已接到真实 `auth.JWTAuth()`
4. 已补测试覆盖：
   - 权限配置文件加载
   - `author` 项目/文档权限生效
   - 遗留鉴权入口对匿名请求返回 401

因此，`/api/v1/writer/projects` 这类 401/权限配置阻塞已不再作为活跃问题保留。

---

## 保留说明

原 Issue 中有一部分描述已经被后续代码演进覆盖：

- 实际主路由早已普遍使用 `auth.JWTAuth()` 与路由级角色校验
- 文档中对“权限中间件未启用”的判断只对遗留 `internal/router/setup.go` 成立，现已修复

仍保留为后续优化项的内容：

- 对嵌套路由的统一资源权限映射
- 统一的权限策略文档与治理规则
- 是否将零散路由级角色校验进一步收敛到统一权限中间件

---

## 后续 TODO

- 优化嵌套路由资源解析，避免 `/api/v1/writer/projects` 这类路径在通用权限中间件中被识别为模块名而非资源名
- 形成统一权限治理文档，明确角色校验与权限校验边界
- 如需要更细粒度权限控制，再评估将路由级角色检查进一步统一到权限中间件

---

## 相关 Issue

- 已归档关联: [#008: 中间件架构问题](./008-middleware-architecture-issues.md)
- 持续治理关联: [#005: API 标准化问题](../005-api-standardization-issues.md)
