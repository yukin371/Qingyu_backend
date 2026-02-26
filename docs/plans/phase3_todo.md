# 通信模块重构 - 待办事项 (TODO)

> 创建日期: 2026-02-27
> 状态: Part 1 已完成

## 已完成 ✅

- [x] 废弃简化版 notification 模块
- [x] 优化 announcements 公开API
- [x] 优化 announcements 管理员API

## 待完成任务

### 1. Messaging 模块整合 (P2 - 中优先级)
- **预计工作量**: 8-12小时
- **依赖条件**: 功能差异分析完成
- **说明**: 旧版有@提醒功能需要迁移到新版

### 2. 路由配置清理
- **预计工作量**: 2小时
- **说明**:
  - 移除废弃模块的路由引用
  - 确认所有路由正确注册
  - 更新路由文档

### 3. 测试验证
- **预计工作量**: 4小时
- **说明**:
  - 运行单元测试
  - 运行集成测试
  - 修复发现的问题

### 4. 文档更新
- **预计工作量**: 2小时
- **说明**:
  - 更新 API 文档
  - 更新架构图
  - 记录废弃变更

## 完成标准

- [x] Part 1: notification/announcements 优化完成
- [ ] Part 2: 路由配置清理完成
- [ ] Part 3: 测试验证通过
- [ ] Part 4: 文档更新完成

## 参考资料

- [API重构总计划](./api_refactor_plan.md)
- [Messaging模块迁移TODO](../memory/phase3-messaging-migration-todo)
