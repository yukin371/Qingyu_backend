# Saga 设计历史入口

> 最后整理: 2026-05-22

本目录属于核心能力中的历史 Saga 设计区，当前架构事实和实施 owner 已迁回 `docs/architecture/` 与 `docs/implementation/`。

## Page Role

- legacy-hub
- current-owner: `design/core/saga/`
- current-bounded: 历史 Saga 设计入口，只负责事务协调旧文档导航

## Recommended Read Path

1. 先读 `../../../architecture/README.md`。
2. 再读 `../../../implementation/README.md`。
3. 需要回看历史 Saga 设计时，再使用本目录。

## Quick Section Map

- Directory Role
- Historical Documents
- Boundary

## Quick Takeaways

- 这是历史 Saga 设计入口，不是当前事务协调 owner。

## Skip Guide

- 只看当前事务协调：跳过本目录。

## Directory Role

- 当前目录入口就是本页 `README.md`。
- 这里保留的是 Saga 事务模式历史设计稿，不再作为当前事务协调机制入口。

## Historical Documents

- [Saga事务模式设计.md](./Saga事务模式设计.md) - Saga 事务模式历史详细设计稿

## Boundary

- 当前架构 owner：`../../../architecture/README.md`
- 当前实施记录 owner：`../../../implementation/README.md`
- 当前标准 owner：`../../../standards/README.md`
