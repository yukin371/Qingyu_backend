# 通信设计文档

> 最后整理: 2026-05-22  
> 当前状态: `summary-draft`

本文档是 `design/communication/` 的早期汇总稿，适合快速回看当时如何把“实时通信 + 消息推送”放在同一专题下组织；它不是当前目录的标准入口，也不等同于当前通信能力的 owner 文档。

## Page Role

- 这里负责：历史通信专题总览、早期分组方式、当时的设计目标与技术草图。
- 不负责：当前目录标准入口、当前通信模块实现边界、现行 API / implementation / architecture owner。

## Recommended Read Path

1. [README.md](./README.md)
2. [../../architecture/README.md](../../architecture/README.md)
3. [../../implementation/README.md](../../implementation/README.md)
4. [../../api/README.md](../../api/README.md)

## Boundary

- 如果你要找“communication 目录现在怎么读”，优先看 [README.md](./README.md)。
- 如果你要找“当前通信能力的实现与分层 owner”，优先看 [../../architecture/README.md](../../architecture/README.md) 和 [../../implementation/README.md](../../implementation/README.md)。
- 如果你要找“当前对外接口或路由入口”，优先看 [../../api/README.md](../../api/README.md)。

## 📁 文档目录

### 实时通信

- [WebSocket实时通信设计](./WebSocket实时通信设计.md) - WebSocket协议、连接管理、消息格式

### 消息推送

- [消息推送系统设计](./消息推送系统设计.md) - 推送机制、消息队列、通知管理

## 🎯 设计目标

### 实时性

- 低延迟的消息传输
- 实时状态同步
- 即时通知推送

### 可靠性

- 消息可达保证
- 断线重连机制
- 消息持久化存储

### 可扩展性

- 支持大规模并发连接
- 水平扩展能力
- 负载均衡策略

## 📊 技术架构

### WebSocket服务

```
Client ←──WebSocket──→ Gateway ←──→ Message Queue ←──→ Service
                          ↓
                       Redis Pub/Sub
```

### 消息推送流程

```
Event Source → Message Queue → Push Service → User Device
                    ↓
              Notification DB
```

## 🔧 核心功能

### WebSocket通信

- 连接建立与认证
- 心跳检测与保活
- 消息发送与接收
- 错误处理与重连

### 消息推送

- 系统通知推送
- 用户消息推送
- 实时提醒推送
- 推送统计分析

## 🔗 相关文档

- [共享服务设计归档入口](../shared/README.md) - 消息队列等共享服务旧设计已迁入归档区
- [平台设计归档入口](../platform/README.md) - 通知系统等平台层旧设计已迁入归档区

## 📝 更新日志

- 2025-09-30: 更新通信设计文档目录，整理现有设计文档
