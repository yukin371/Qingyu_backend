# 运维当前指南索引

> 最后整理: 2026-05-22
> 当前状态: `current-owner`

这一组文档回答的是“现在怎么部署、怎么监控、怎么排障、怎么处理 CI/CD”。如果你要执行当前运维动作，先从这里进。

## 推荐阅读顺序

1. [README.md](./README.md) - 运维总入口与边界说明
2. [部署指南.md](./部署指南.md) - 环境准备、部署流程与联调基线
3. [监控体系使用指南.md](./监控体系使用指南.md) - 监控指标、排障与可观测性
4. [性能优化指南.md](./性能优化指南.md) - 性能诊断与调优
5. [CI_CD配置指南.md](./CI_CD配置指南.md) - CI/CD 主配置说明
6. [CI_CD问题解决方案.md](./CI_CD问题解决方案.md) - 常见 CI/CD 问题排障
7. [快速参考-CI_CD命令.md](./快速参考-CI_CD命令.md) - 命令速查
8. [../deployment/README.md](../deployment/README.md) - 部署专题与 Docker 测试环境入口

## Current Workflows

| 工作流 | 推荐入口 | 适用场景 |
|------|------|------|
| 部署与联调 | [部署指南.md](./部署指南.md) / [../deployment/README.md](../deployment/README.md) | 环境准备、启动服务、Docker 测试环境、联调基线 |
| 监控与排障 | [监控体系使用指南.md](./监控体系使用指南.md) | 指标观察、日志定位、问题排查 |
| 性能优化 | [性能优化指南.md](./性能优化指南.md) | 性能瓶颈分析、调优动作、压测后的诊断 |
| CI/CD 配置 | [CI_CD配置指南.md](./CI_CD配置指南.md) / [快速参考-CI_CD命令.md](./快速参考-CI_CD命令.md) | 理解流水线配置、查常用命令 |
| CI/CD 排障 | [CI_CD问题解决方案.md](./CI_CD问题解决方案.md) | Runner、Docker、构建失败等常见问题处理 |

## Boundary

- 本页：当前运维动作的导航页。
- [../deployment/README.md](../deployment/README.md): 部署专题入口，尤其是 Docker 测试环境和外部配置。
- `Qingyu_backend/docker/README*.md`: Docker 配置和 compose 文件的真实 owner。
- [README_历史记录.md](./README_历史记录.md): 历史 CI/CD / MVP / 检查报告归档入口。

## Routing Rules

1. 当前仍在使用的部署、监控、CI/CD 与性能文档，优先留在这里。
2. 如果某篇文档只剩历史背景或阶段记录价值，应转到 [README_历史记录.md](./README_历史记录.md)。
3. 如果文档已经成为具体部署专题而不是运维总规则，应优先挂到 [../deployment/README.md](../deployment/README.md)。
