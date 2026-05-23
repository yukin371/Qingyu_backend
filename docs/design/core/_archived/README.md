# 核心设计归档入口

> 最后整理: 2026-05-22

本目录属于 `design/core/` 下的冻结归档区，只保留已经搁置、被替代或明确不再作为当前方案推进的历史设计稿。

## Page Role

- legacy-frozen-hub
- current-owner: `design/core/_archived/`
- current-bounded: 核心设计冻结归档入口，只负责 frozen / superseded 文档导航

## Recommended Read Path

1. 先读 `../README.md`。
2. 再读 `../../../architecture/README.md`。
3. 只有需要追溯冻结方案时，再使用本目录。

## Quick Section Map

- Directory Role
- Historical Documents
- Boundary

## Quick Takeaways

- 这是冻结归档入口，不是当前核心设计入口。

## Skip Guide

- 只看当前核心架构：跳过本目录。

## Directory Role

- 当前目录入口就是本页 `README.md`。
- 这里保留的是 `legacy-frozen` 与 `legacy-superseded` 文档，不再作为当前功能设计、当前计划或当前架构入口。

## Historical Documents

- [插件扩展系统(搁置).md](<./插件扩展系统(搁置).md>) - 已搁置的插件扩展系统历史方案
- [协作编辑系统（搁置）.md](./协作编辑系统（搁置）.md) - 已搁置的协作编辑历史方案
- [项目CRUD.md](./项目CRUD.md) - 已被后续项目结构设计替代的早期 CRUD 方案
- [项目_节点_文档_CRUD_Design.md](./项目_节点_文档_CRUD_Design.md) - 已被四层架构取代的三层 CRUD 方案

## Boundary

- 当前架构 owner：`../../../architecture/README.md`
- 当前 API 入口 owner：`../../../api/README.md`
- 当前实施记录 owner：`../../../implementation/README.md`
- 本目录只回答“旧方案当时怎么想”，不回答“现在应该怎么做”。
