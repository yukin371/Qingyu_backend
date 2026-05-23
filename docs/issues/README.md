# 后端问题台账入口

> 最后整理: 2026-05-22

## 目录角色

- `Qingyu_backend/docs/issues/` 保存仍有参考价值的问题台账、问题拆解和整改背景。
- 这里不是新的跨仓库治理方案 owner；长期方案请回到父仓 `docs/plans/submodules/backend/`。
- 历史审计/诊断报告已经统一挂到父仓 `docs/backend/archive/reports/`，本目录不再复制一套 `reports/` 子树。

## 当前活跃台账

- [003-test-infrastructure-improvements.md](./003-test-infrastructure-improvements.md)
- [004-code-quality-improvements.md](./004-code-quality-improvements.md)
- [005-api-standardization-issues.md](./005-api-standardization-issues.md)
- [009-test-coverage-issues.md](./009-test-coverage-issues.md)
- [010-repository-layer-business-logic-permeation.md](./010-repository-layer-business-logic-permeation.md)
- [011-frontend-backend-data-type-inconsistency.md](./011-frontend-backend-data-type-inconsistency.md)
- [012-service-layer-id-conversion-refactor.md](./012-service-layer-id-conversion-refactor.md)

## 历史归档

- [archived/001-unify-id-type-in-models.md](./archived/001-unify-id-type-in-models.md)
- [archived/002-create-method-id-not-set-bug.md](./archived/002-create-method-id-not-set-bug.md)
- [archived/006-database-index-issues.md](./archived/006-database-index-issues.md)
- [archived/007-transaction-management.md](./archived/007-transaction-management.md)
- [archived/008-middleware-architecture-issues.md](./archived/008-middleware-architecture-issues.md)
- [archived/012-auth-401-and-permission-issues.md](./archived/012-auth-401-and-permission-issues.md)
- [archived/013-test-user-seed-id-not-set.md](./archived/013-test-user-seed-id-not-set.md)

## Topic READMEs

- [archived/README.md](./archived/README.md)
- [todo/README.md](./todo/README.md)

## 配套入口

- [ISSUE_RELATIONSHIPS.md](./ISSUE_RELATIONSHIPS.md)
- [../implementation/README.md](../implementation/README.md)
- [../testing/README.md](../testing/README.md)
- [../standards/README.md](../standards/README.md)
- [../../../docs/plans/submodules/backend/README.md](../../../docs/plans/submodules/backend/README.md)

## 使用说明

- 如果只是回看问题背景，优先从本页进入。
- 如果需要查原始诊断报告，直接使用父仓 `docs/backend/archive/reports/` 中的历史报告，不要在后端子模块里补回旧 `reports/` 目录。
- 如果某个问题已经稳定收敛到长期规则，应继续沉淀到父仓 `plans/`、`standards/` 或 `ADR`，而不是在本目录继续扩写。
