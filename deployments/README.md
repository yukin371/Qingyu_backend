# Qingyu 监控系统部署指南

## 概述

本监控系统使用Prometheus + Grafana + Alertmanager构建，提供完整的监控、可视化和告警功能。

## 组件说明

| 组件 | 端口 | 功能 |
|------|------|------|
| Prometheus | 9090 | 指标采集和存储 |
| Grafana | 3000 | 可视化仪表板 |
| Alertmanager | 9093 | 告警管理和通知 |
| Node Exporter | 9100 | 系统指标采集 |

## 快速开始

### 1. 配置环境变量

```bash
cp .env.monitoring.example .env.monitoring
```

编辑`.env.monitoring`文件，配置必要的环境变量：

```env
# Grafana管理员账号
GRAFANA_ADMIN_USER=admin
GRAFANA_ADMIN_PASSWORD=your-secure-password

# SMTP邮件配置（用于告警通知）
SMTP_PASSWORD=your-smtp-password
```

### 2. 启动监控系统

```bash
# 启动所有服务
docker-compose -f docker-compose.monitoring.yml up -d

# 查看服务状态
docker-compose -f docker-compose.monitoring.yml ps

# 查看日志
docker-compose -f docker-compose.monitoring.yml logs -f
```

### 3. 访问服务

- **Grafana**: http://localhost:3000
  - 默认账号: admin
  - 默认密码: admin（首次登录后需修改）

- **Prometheus**: http://localhost:9090

- **Alertmanager**: http://localhost:9093

## 配置说明

### Prometheus配置

**配置文件**: `prometheus/prometheus.yml`

主要配置项：
- 全局抓取间隔：15秒
- 告警规则文件：`prometheus/alerts.yml`
- 数据保留时间：30天

**告警规则**: `prometheus/alerts.yml`

定义的告警规则包括：
- API性能告警（错误率、响应时间）
- 数据库告警（连接池、查询性能）
- 缓存告警（命中率）
- 系统资源告警（CPU、内存）
- 业务指标告警（支付、用户注册）

### Alertmanager配置

**配置文件**: `alertmanager/alertmanager.yml`

主要配置：
- SMTP邮件服务器
- 告警路由规则
- 告警分组和抑制

### Grafana配置

**仪表板**: `grafana/dashboards/qingyu-dashboard.json`

预置仪表板包含：
- 请求速率
- 响应时间（P50, P95, P99）
- 错误率
- 数据库连接数
- 缓存命中率
- 系统资源使用

## 使用指南

### 查看指标

1. 访问Grafana: http://localhost:3000
2. 登录后选择"Qingyu API Dashboard"
3. 查看各项指标的可视化展示

### 查看告警

1. 访问Alertmanager: http://localhost:9093
2. 查看当前活跃的告警
3. 查看告警历史和静默规则

### 查询Prometheus

1. 访问Prometheus: http://localhost:9090
2. 使用PromQL查询指标
3. 示例查询：
   ```
   # 查看请求速率
   rate(http_requests_total[5m])
   
   # 查看P95响应时间
   histogram_quantile(0.95, http_request_duration_seconds)
   
   # 查看错误率
   rate(http_requests_total{status=~"5.."}[5m])
   ```

## 维护操作

### 备份和恢复

```bash
# 备份Prometheus数据
docker exec qingyu-prometheus tar czf /tmp/prometheus-backup.tar.gz /prometheus
docker cp qingyu-prometheus:/tmp/prometheus-backup.tar.gz ./backups/

# 备份Grafana数据
docker exec qingyu-grafana tar czf /tmp/grafana-backup.tar.gz /var/lib/grafana
docker cp qingyu-grafana:/tmp/grafana-backup.tar.gz ./backups/
```

### 更新配置

```bash
# 更新Prometheus配置后重载
curl -X POST http://localhost:9090/-/reload

# 重启Alertmanager
docker-compose -f docker-compose.monitoring.yml restart alertmanager

# 重启Grafana
docker-compose -f docker-compose.monitoring.yml restart grafana
```

### 查看日志

```bash
# 查看所有服务日志
docker-compose -f docker-compose.monitoring.yml logs

# 查看特定服务日志
docker-compose -f docker-compose.monitoring.yml logs prometheus
docker-compose -f docker-compose.monitoring.yml logs grafana
docker-compose -f docker-compose.monitoring.yml logs alertmanager

# 实时跟踪日志
docker-compose -f docker-compose.monitoring.yml logs -f prometheus
```

## 停止服务

```bash
# 停止所有服务
docker-compose -f docker-compose.monitoring.yml down

# 停止并删除数据卷
docker-compose -f docker-compose.monitoring.yml down -v
```

## 故障排查

### Prometheus无法采集指标

1. 检查后端服务是否暴露`/metrics`端点
2. 检查Prometheus配置中的`targets`是否正确
3. 访问 http://localhost:9090/targets 查看采集状态

### Grafana无法连接Prometheus

1. 检查Prometheus服务是否运行：`docker ps | grep prometheus`
2. 检查Grafana数据源配置
3. 查看Grafana日志：`docker logs qingyu-grafana`

### 告警邮件未发送

1. 检查SMTP配置是否正确
2. 检查`.env.monitoring`中的`SMTP_PASSWORD`
3. 查看Alertmanager日志：`docker logs qingyu-alertmanager`
4. 测试SMTP连接：
   ```bash
   docker exec qingyu-alertmanager amtool config validate
   ```

## 监控指标说明

### HTTP请求指标

| 指标名称 | 类型 | 标签 | 说明 |
|---------|------|------|------|
| `http_requests_total` | Counter | method, path, status | HTTP请求总数 |
| `http_request_duration_seconds` | Histogram | method, path | HTTP请求持续时间 |
| `http_request_size_bytes` | Histogram | method, path | HTTP请求大小 |
| `http_response_size_bytes` | Histogram | method, path | HTTP响应大小 |

### 数据库指标

| 指标名称 | 类型 | 标签 | 说明 |
|---------|------|------|------|
| `db_pool_connections` | Gauge | database | 数据库连接数 |
| `db_query_duration_seconds` | Histogram | database, operation | 数据库查询持续时间 |
| `db_queries_total` | Counter | database, operation, status | 数据库查询总数 |

### 缓存指标

| 指标名称 | 类型 | 标签 | 说明 |
|---------|------|------|------|
| `cache_hits_total` | Counter | cache | 缓存命中总数 |
| `cache_misses_total` | Counter | cache | 缓存未命中总数 |
| `cache_duration_seconds` | Histogram | cache, operation | 缓存操作持续时间 |

## 扩展和定制

### 添加自定义告警规则

编辑`prometheus/alerts.yml`，添加新的告警规则：

```yaml
groups:
  - name: custom_alerts
    interval: 30s
    rules:
      - alert: CustomAlert
        expr: your_expression_here
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "自定义告警"
          description: "告警详情"
```

### 添加Grafana仪表板

1. 在Grafana UI中创建新仪表板
2. 导出仪表板JSON
3. 将JSON文件保存到`grafana/dashboards/`目录
4. 重启Grafana或重新加载配置

## 联系方式

如有问题，请联系：
- 技术支持: [待定]
- 监控告警: [待定]

## 参考文档

- [Prometheus官方文档](https://prometheus.io/docs/)
- [Grafana官方文档](https://grafana.com/docs/)
- [Alertmanager官方文档](https://prometheus.io/docs/alerting/latest/alertmanager/)
