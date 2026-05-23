# Ops Docs Hub

本目录是 `Qingyu_backend` 当前运维文档的 canonical owner，也是部署相关文档的总入口。

## Current Rule

- 部署、监控、CI/CD、性能优化等总规则，统一写在 `Qingyu_backend/docs/ops/`。
- `Qingyu_backend/docs/deployment/` 只保留部署专题，不再作为第二个运维总入口。
- `Qingyu_backend/docs/guides/ops/` 仅保留历史镜像与迁移提示，不再作为主入口维护。

## Recommended Read Path

1. [README_当前指南.md](./README_当前指南.md)
2. [README_历史记录.md](./README_历史记录.md)
3. [../deployment/README.md](../deployment/README.md)

## Current Layers

### 当前指南

- [README_当前指南.md](./README_当前指南.md): 当前部署、监控、CI/CD、性能文档索引
- [部署指南.md](./部署指南.md): 部署与环境准备主说明
- [监控体系使用指南.md](./监控体系使用指南.md): 监控、指标与排障入口
- [性能优化指南.md](./性能优化指南.md): 性能诊断与优化参考
- [CI_CD配置指南.md](./CI_CD配置指南.md): CI/CD 主配置说明
- [CI_CD问题解决方案.md](./CI_CD问题解决方案.md): CI/CD 常见问题排障
- [快速参考-CI_CD命令.md](./快速参考-CI_CD命令.md): CI/CD 命令速查

### 部署专题

- [../deployment/README.md](../deployment/README.md): 本地启动、Docker 测试环境、OAuth 注册等部署专题入口

### 历史记录

- [README_历史记录.md](./README_历史记录.md): 历史运维与 CI/CD 记录索引
- [CICD配置说明.md](./CICD配置说明.md): 历史 CI/CD 说明，使用时需结合现行配置核验
- [CI_DOCKER_检查报告.md](./CI_DOCKER_检查报告.md): 阶段性检查报告
- [CI_CD_improvements_2025-10-20.md](./CI_CD_improvements_2025-10-20.md): 历史改进记录
- [MVP部署指南_2025-10-23.md](./MVP部署指南_2025-10-23.md): 历史阶段部署说明

## Maintenance Note

1. 新增运维正文前，先确认是否属于 `ops/` 而不是 `guides/ops/`。
2. 若文档只是部署专题而非运维总规则，优先落到 `../deployment/`。
3. 若引用历史 CI/CD 文档，请优先补“适用范围”说明，避免把阶段报告当作当前规范。
4. 目录内优先通过 `README_当前指南.md` 和 `README_历史记录.md` 分流，减少直接面对平铺文件列表的成本。
