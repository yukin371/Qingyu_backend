package config

import (
	"sync"

	"go.uber.org/zap"
)

var (
	// reloadHandlers 存储配置重载处理器
	reloadHandlers = make(map[string]func())
	// reloadMutex 保护reloadHandlers的并发访问
	reloadMutex sync.RWMutex
)

// RegisterReloadHandler 注册配置重载处理器
func RegisterReloadHandler(name string, handler func()) {
	reloadMutex.Lock()
	defer reloadMutex.Unlock()
	reloadHandlers[name] = handler
}

// UnregisterReloadHandler 注销配置重载处理器
func UnregisterReloadHandler(name string) {
	reloadMutex.Lock()
	defer reloadMutex.Unlock()
	delete(reloadHandlers, name)
}

// handleConfigReload 处理配置重载
func handleConfigReload() {
	reloadMutex.RLock()
	defer reloadMutex.RUnlock()

	for name, handler := range reloadHandlers {
		func() {
			defer func() {
				if r := recover(); r != nil {
					zap.L().Error("Error executing reload handler",
						zap.String("handler", name),
						zap.Any("panic", r))
				}
			}()
			handler()
		}()
	}
}

// EnableHotReload 启用配置热重载
func EnableHotReload() {
	WatchConfig(func() {
		zap.L().Info("Configuration changed, reloading")
		handleConfigReload()
	})
}
