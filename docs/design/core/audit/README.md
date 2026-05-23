# 审核设计历史入口

> 最后整理: 2026-05-22

本目录属于核心能力中的历史审核设计区，当前审核相关 owner 已迁回 `docs/architecture/`、`docs/api/` 与 `docs/implementation/`。

## Page Role

- legacy-hub
- current-owner: `design/core/audit/`
- current-bounded: 历史审核设计入口，只负责审核旧文档导航

## Recommended Read Path

1. 先读 `../../../architecture/README.md`。
2. 再读 `../../../api/README.md`。
3. 需要回看历史审核设计时，再使用本目录。

## Quick Section Map

- Directory Role
- Historical Documents
- Boundary

## Quick Takeaways

- 这是历史审核设计入口，不是当前审核 owner。

## Skip Guide

- 只看当前审核：跳过本目录。

## Directory Role

- 当前目录入口就是本页 `README.md`。
- 这里保留的是内容审核规则引擎和审核数据模型的历史设计稿，不再作为当前审核系统入口。

## Historical Documents

- [内容审核规则引擎设计.md](./内容审核规则引擎设计.md) - 审核规则引擎历史设计稿
- [审核系统数据模型设计.md](./审核系统数据模型设计.md) - 审核系统数据模型历史设计稿

## Boundary

- 当前架构 owner：`../../../architecture/README.md`
- 当前 API 入口 owner：`../../../api/README.md`
- 当前实施记录 owner：`../../../implementation/README.md`
