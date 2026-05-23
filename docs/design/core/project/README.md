# 项目核心设计历史入口

> 最后整理: 2026-05-22

本目录属于核心能力中的历史项目设计区，当前项目结构、文档管理、版本管理相关 owner 已迁回 `docs/architecture/`、`docs/api/`、`docs/implementation/` 与 `docs/standards/`。

## Page Role

- legacy-hub
- current-owner: `design/core/project/`
- current-bounded: 历史项目核心设计入口，只负责项目结构与版本相关旧文档导航

## Recommended Read Path

1. 先读 `../../../architecture/README.md`。
2. 再读 `../../../implementation/README.md`。
3. 需要回看历史项目设计时，再使用本目录。

## Quick Section Map

- Directory Role
- Historical Documents
- Boundary

## Quick Takeaways

- 这是历史项目核心设计入口，不是当前项目能力 owner。

## Skip Guide

- 只看当前项目结构：跳过本目录。

## Directory Role

- 当前目录入口就是本页 `README.md`。
- 这里保留的是项目路径、四层 CRUD、版本控制和版本服务重构的历史设计稿，不再作为当前项目能力主入口。

## Historical Documents

- [项目路径管理设计.md](./项目路径管理设计.md) - 项目路径与树形节点历史设计稿
- [项目四层架构CRUD设计.md](./项目四层架构CRUD设计.md) - 项目/节点/文档/内容四层结构历史设计稿
- [版本管理服务重构设计.md](./版本管理服务重构设计.md) - 版本服务历史重构方案
- [版本控制.md](./版本控制.md) - 版本控制历史方案说明

## Boundary

- 当前架构 owner：`../../../architecture/README.md`
- 当前 API 入口 owner：`../../../api/README.md`
- 当前实施记录 owner：`../../../implementation/README.md`
- 当前标准 owner：`../../../standards/README.md`
