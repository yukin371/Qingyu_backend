# Deployment Topic Hub

本目录是 `Qingyu_backend/docs` 下的“部署专题区”，不是第二个运维总入口。

## Current Boundary

- `../ops/README.md` 是当前运维与部署文档的 canonical owner。
- 本目录只保留部署专题文档，例如本地启动、Docker 测试环境、OAuth 应用注册这类专门主题。
- `Qingyu_backend/docker/` 下的 `README.md`、`README_TEST.md` 和 compose 文件是 Docker 配置的真实文件 owner；本目录只做索引与使用说明，不重复维护第二套长正文。

## Recommended Read Path

1. [../ops/README.md](../ops/README.md)
2. [服务启动指南.md](./服务启动指南.md)
3. [README_TEST.md](./README_TEST.md)
4. [测试环境快速验证指南.md](./测试环境快速验证指南.md)
5. [oauth-app-registration-guide.md](./oauth-app-registration-guide.md)

## Current Topics

### 启动与联调

- [服务启动指南.md](./服务启动指南.md): 本地启动、联调端口、服务健康检查与基础排障

### Docker 测试环境

- [README_TEST.md](./README_TEST.md): Docker 测试环境入口，优先指向 `Qingyu_backend/docker/README_TEST.md`
- [测试环境快速验证指南.md](./测试环境快速验证指南.md): 快速验证专题

### 外部配置专题

- [oauth-app-registration-guide.md](./oauth-app-registration-guide.md): OAuth 第三方应用注册与环境变量配置专题

## Related Owners

- [../ops/README.md](../ops/README.md): 运维与部署总入口。
- [../testing/README.md](../testing/README.md): 测试入口。
- [../../docker/README.md](../../docker/README.md): Docker 开发/生产配置 owner。
- [../../docker/README_TEST.md](../../docker/README_TEST.md): Docker 测试环境 owner。
- [../../docker/测试环境快速验证指南.md](../../docker/测试环境快速验证指南.md): Docker 测试环境快速验证 owner。

## Rule

1. 新增部署类专题文档，可以写在本目录；但若是运维总规则、监控、CI/CD、性能优化，应优先写到 `../ops/`。
2. 不要把 `docker/README*.md` 的长正文复制到本目录；只保留入口和使用场景说明。
3. 若未来某篇部署文档已经失去专题价值，应先确认备份，再考虑归档到 `../archive/`。
4. 若只是 Docker 测试环境子话题，优先挂到 `README_TEST.md` 或 `测试环境快速验证指南.md` 这条路径，不要再新长第二套测试环境入口。
