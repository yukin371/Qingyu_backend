# SharedServiceContainer 废弃迁移指南

**日期**: 2025-10-24  
**状态**: 🚧 进行中  
**版本**: v1.0

## 📋 概述

`SharedServiceContainer` 已被废弃，所有共享服务现在统一由 `service/container/ServiceContainer` 管理。这样可以避免两套管理机制的冲突，统一服务入口，简化架构。

## 🎯 迁移原因

### 问题
1. **状态不一致**: 两个容器可能管理不同的服务实例
2. **生命周期混乱**: 服务可能被初始化多次
3. **服务获取不一致**: 从不同容器获取的服务实例可能不同
4. **指标收集混乱**: ServiceContainer 有指标系统，SharedServiceContainer 没有

### 解决方案
- 废弃 `SharedServiceContainer`
- 统一使用 `ServiceContainer` 管理所有服务（包括共享服务）
- 所有共享服务通过 `ServiceContainer` 的 Getter 方法获取

## 📝 已完成的修改

### 1. 标记废弃
- ✅ `service/shared/container/shared_service_container.go` - 添加 DEPRECATED 注释
- ✅ `service/shared/container/shared_service_factory.go` - 添加 DEPRECATED 注释
- ✅ `service/shared/container/shared_service_container_test.go` - 标记测试为废弃

### 2. 修改路由注册
- ✅ `router/shared/shared_router.go` - 改为接收独立服务参数而非容器
- ✅ `router/enter.go` - 改为从 ServiceContainer 获取共享服务

### 3. 编译验证
- ✅ 项目编译通过，无错误

## 🔄 迁移步骤

### 对于新代码

**之前 (废弃)**:
```go
import "Qingyu_backend/service/shared/container"

// 创建独立的共享服务容器
sharedContainer := container.NewSharedServiceContainer()
sharedContainer.SetAuthService(authService)
sharedContainer.SetWalletService(walletService)

// 获取服务
authSvc := sharedContainer.AuthService()
walletSvc := sharedContainer.WalletService()
```

**现在 (推荐)**:
```go
import "Qingyu_backend/service"

// 使用全局服务容器
serviceContainer := service.GetServiceContainer()

// 获取共享服务
authSvc, err := serviceContainer.GetAuthService()
if err != nil {
    log.Printf("AuthService未配置: %v", err)
}

walletSvc, err := serviceContainer.GetWalletService()
if err != nil {
    log.Printf("WalletService未配置: %v", err)
}
```

### 对于现有代码

#### 示例1: 路由注册

**之前**:
```go
// router/enter.go
sharedContainer := container.NewSharedServiceContainer()
sharedGroup := v1.Group("/shared")
sharedRouter.RegisterRoutes(sharedGroup, sharedContainer)
```

**现在**:
```go
// router/enter.go
serviceContainer := service.GetServiceContainer()

// 获取共享服务
authSvc, _ := serviceContainer.GetAuthService()
walletSvc, _ := serviceContainer.GetWalletService()
storageSvc, _ := serviceContainer.GetStorageService()

// 注册路由
sharedGroup := v1.Group("/shared")
sharedRouter.RegisterRoutes(sharedGroup, authSvc, walletSvc, storageSvc)
```

#### 示例2: API 处理器

**之前**:
```go
// 从 SharedServiceContainer 获取服务
authAPI := api.NewAuthAPI(sharedContainer.AuthService())
```

**现在**:
```go
// 从 ServiceContainer 获取服务
authSvc, err := serviceContainer.GetAuthService()
if err != nil {
    return err
}
authAPI := api.NewAuthAPI(authSvc)
```

## 🗑️ 待删除的文件

在确认所有功能正常后，以下文件将被删除：

### 即将删除
- [ ] `service/shared/container/shared_service_container.go`
- [ ] `service/shared/container/shared_service_factory.go`
- [ ] `service/shared/container/shared_service_container_test.go`
- [ ] `service/shared/container/shared_service_factory_test.go` (如果存在)
- [ ] `service/shared/container/test_mocks.go` (如果存在)

### 删除前检查清单
- [ ] 运行所有测试确保无回归
- [ ] 检查是否有其他代码引用这些文件
- [ ] 确认服务容器功能完整
- [ ] 更新相关文档

## ✅ 验证步骤

### 1. 编译检查
```bash
go build ./...
```
**状态**: ✅ 通过

### 2. 运行测试
```bash
# 运行服务容器测试
go test ./test/service/container/... -v

# 运行路由测试
go test ./test/api/... -v

# 运行集成测试
go test ./test/integration/... -v
```
**状态**: 🔄 待运行

### 3. 启动服务验证
```bash
go run cmd/server/main.go
```
检查点:
- [ ] 服务正常启动
- [ ] 路由正确注册
- [ ] 共享服务可用（如果已配置）
- [ ] 日志显示正常

## 📊 影响范围分析

### 受影响的模块
1. **router 模块** - 已修改路由注册逻辑 ✅
2. **api 模块** - 无直接影响，通过路由获取服务
3. **service 模块** - SharedServiceContainer 标记为废弃 ✅
4. **test 模块** - 测试文件标记为废弃 ✅

### 不受影响的模块
- **models** - 无影响
- **repository** - 无影响
- **middleware** - 无影响
- **pkg** - 无影响

## 🚀 后续任务

### 高优先级
1. [ ] 运行完整测试套件验证功能
2. [ ] 在 `service/container/service_container.go` 的 `SetupDefaultServices()` 中添加共享服务初始化
3. [ ] 验证服务启动和运行

### 中优先级
4. [ ] 创建共享服务初始化示例代码
5. [ ] 更新开发文档
6. [ ] 通知团队成员关于此变更

### 低优先级
7. [ ] 在确认一切正常后，删除废弃文件
8. [ ] 清理相关的导入和引用
9. [ ] 更新 CHANGELOG

## 📚 相关文档

- **架构设计规范**: `doc/architecture/架构设计规范.md`
- **服务容器文档**: `doc/architecture/服务容器集成报告.md`
- **共享服务实现报告**: `doc/architecture/共享服务实现报告_2025-10-24.md`

## ⚠️ 注意事项

1. **兼容性**: 为了平稳过渡，废弃的代码暂时保留，但会输出警告
2. **测试覆盖**: 确保新的 ServiceContainer 测试覆盖原有功能
3. **文档更新**: 所有文档中提到 SharedServiceContainer 的地方需要更新
4. **团队沟通**: 及时通知团队成员这一变更

## 📞 联系方式

如有问题或疑问，请联系：
- **负责人**: 青羽后端架构团队
- **文档**: 见上述"相关文档"部分

---

**最后更新**: 2025-10-24  
**状态**: 废弃标记已完成，等待测试验证和最终删除

