# 后端文档分类规范

> 更新日期: 2026-05-22  
> 适用范围: `Qingyu_backend/docs/**`

本文档定义后端文档目录的分类、状态和落盘规则，目标是避免继续把“当前规范、历史设计、阶段报告、迁移说明”混写在同一层。

## 1. 状态模型

所有文档目录只允许落入以下 4 类状态之一：

- `current-owner`: 当前主入口，可以继续新增和维护。
- `current-bounded`: 当前仍可保留，但只承载某一类专题内容，不能扩张成第二个主入口。
- `legacy-readonly`: 历史资料区，只允许补迁移说明、状态提示或必要索引，不再新增正文。
- `archive-only`: 统一归档区，只保留历史镜像、阶段资料和恢复线索。

## 2. 根目录分类

### 2.1 当前主入口

| 目录 | 状态 | 允许内容 | 不允许内容 |
|------|------|----------|------------|
| `architecture/` | `current-owner` | 当前架构边界、模块地图、运行链路 | 阶段性完成报告、历史镜像 |
| `standards/` | `current-owner` | 现行规范、分层规则、文档规范 | 实施记录、历史设计稿 |
| `guides/` | `current-owner` | 上手指南、操作手册、协作路径 | 旧镜像双写、阶段总结 |
| `api/` | `current-owner` | API 文档、Swagger 说明、接口入口 | 历史镜像目录 |
| `ops/` | `current-owner` | 运维入口、环境、运行与排障规范 | 与 `deployment/` 重复的第二套总入口 |
| `testing/` | `current-owner` | 测试入口、测试专题、测试报告索引 | 早期 `test/` 镜像双写 |

### 2.2 当前受限专题区

| 目录 | 状态 | 用途边界 |
|------|------|----------|
| `analysis/` | `current-bounded` | 分析、审计、治理评估，不写现行规范 |
| `review/` | `current-bounded` | 审查结论、风险结论，不写实施标准 |
| `issues/` | `current-bounded` | 问题台账、缺陷记录、归档问题 |
| `database/` | `current-bounded` | 数据库专题说明、种子与结构操作 |
| `deployment/` | `current-bounded` | 部署专题文档；总入口仍以 `ops/` 为准 |
| `implementation/` | `current-bounded` | 已落地实现说明、阶段实现记录 |
| `migration/` | `current-bounded` | 迁移说明、兼容清理说明 |
| `tools/` | `current-bounded` | 工具使用说明、脚本专题入口 |
| `middleware/` | `current-bounded` | 临时专题目录；若继续扩张，后续需并入 `guides/`、`architecture/` 或 `standards/` |

### 2.3 历史只读区

| 目录 | 状态 | 当前规则 |
|------|------|----------|
| `design/` | `legacy-readonly` | 历史设计资料区，只补迁移说明与历史索引 |
| `engineering/` | `legacy-readonly` | 历史需求/SRS 路径，只保留迁移说明 |
| `usage/` | `legacy-readonly` | 历史使用说明路径，只保留迁移说明 |
| `test/` | `legacy-readonly` | 旧测试目录路径，只保留迁移说明 |
| `plans/` | `legacy-readonly` | 本地长期计划已迁往父仓；本目录不再接收新计划 |

### 2.4 归档与系统目录

| 目录 | 状态 | 当前规则 |
|------|------|----------|
| `archive/` | `archive-only` | 历史镜像、阶段资料、归档桶 |
| `.obsidian/` | `archive-only` | 本地工具目录，不作为协作内容 owner |

## 3. 新文档落盘规则

1. 如果文档在回答“当前怎么做”，优先落到 `architecture/`、`standards/`、`guides/`、`api/`、`ops/`、`testing/`。
2. 如果文档在回答“这次发现了什么问题”，优先落到 `analysis/`、`review/` 或 `issues/`。
3. 如果文档只是阶段实现记录，不要写到 `standards/` 或 `architecture/`，应落到 `implementation/`。
4. 如果文档只解释兼容、迁移或弃用，不要混入主入口，应落到 `migration/`。
5. 若文档已不再是当前入口，但仍需保留追溯价值，应先备份，再进入 `archive/`，原路径保留迁移说明 README。

## 4. README 规则

1. 每个 `current-owner` 目录必须有 `README.md` 作为唯一入口。
2. 每个仍在使用的 `current-bounded` 专题目录也应有 `README.md`，至少说明用途边界、当前 owner 和去向。
3. 每个 `legacy-readonly` 目录必须明确写出“历史区/已归档/不再作为主入口”；若目录仍保留正文，优先补标准 `README.md` 而不是继续依赖 `README_*.md`。
4. 每个 `archive-only` 批次必须记录归档日期、归档对象和备份位置；工具目录如 `.obsidian/` 也应有最小说明，避免被误用为协作 owner。
5. 新增并行目录前，必须先判断能否归入现有 owner；不能时应在父仓 plan 或 ADR 中说明原因。

## 5. 目录治理优先级

当前后端 docs 的下一轮治理优先级如下：

1. `deployment/` 与 `ops/` 的边界收紧
2. `middleware/` 是否并入更稳定 owner 的专项判断
3. 其余仍保留正文的 `design/*` 历史专题树持续归档
4. 根层仍存在但只保留迁移说明的目录，持续改成“轻 README + 归档索引”模式
