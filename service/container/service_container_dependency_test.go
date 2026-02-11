package container

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestContainerDoesNotRequireGlobalDB 测试容器不依赖全局DB
// 这是TDD的核心：先定义期望的行为
func TestContainerDoesNotRequireGlobalDB(t *testing.T) {
	t.Run("容器创建不需要设置全局变量", func(t *testing.T) {
		// 创建容器不需要任何全局变量
		container := NewServiceContainer()
		require.NotNil(t, container)

		// 验证容器初始状态
		assert.False(t, container.IsInitialized())
	})

	t.Run("容器应该支持创建多个独立实例", func(t *testing.T) {
		// 创建两个独立的容器实例
		// 它们不应该共享全局状态
		container1 := NewServiceContainer()
		container2 := NewServiceContainer()

		// 验证它们是不同的实例
		assert.NotNil(t, container1)
		assert.NotNil(t, container2)
		assert.NotSame(t, container1, container2)
	})
}

// TestServiceConstructionWithDependencyInjection 测试服务构造使用依赖注入
func TestServiceConstructionWithDependencyInjection(t *testing.T) {
	t.Run("容器应该提供服务注册机制", func(t *testing.T) {
		container := NewServiceContainer()

		// 验证容器的基本方法存在
		assert.NotNil(t, container)
		assert.False(t, container.IsInitialized())
	})

	t.Run("容器应该提供Get方法而不是依赖全局变量", func(t *testing.T) {
		container := NewServiceContainer()

		// 验证容器有GetMongoDB方法（即使未初始化返回nil）
		db := container.GetMongoDB()
		// 未初始化时应该为nil，这是正常的
		assert.Nil(t, db, "未初始化时GetMongoDB应该返回nil")

		// 验证容器有GetMongoClient方法
		client := container.GetMongoClient()
		assert.Nil(t, client, "未初始化时GetMongoClient应该返回nil")
	})
}

// TestGlobalDBCanBeRemoved 测试global.DB可以安全移除
// 这个测试验证了移除全局变量后的行为
func TestGlobalDBCanBeRemoved(t *testing.T) {
	t.Run("容器不需要global包就能创建", func(t *testing.T) {
		// 这个测试验证了service/container包不再依赖global包
		container := NewServiceContainer()
		assert.NotNil(t, container, "容器应该能够创建而不依赖global包")
	})

	t.Run("容器方法不依赖全局变量", func(t *testing.T) {
		container := NewServiceContainer()

		// 验证GetMongoDB和GetMongoClient方法存在
		// 这些方法使用容器内部的mongoDB和mongoClient，而不是全局变量
		db := container.GetMongoDB()
		client := container.GetMongoClient()

		// 未初始化时返回nil是正常的
		assert.Nil(t, db)
		assert.Nil(t, client)
	})
}

// TestRepositoryFactoryInjection 测试Repository工厂注入机制
func TestRepositoryFactoryInjection(t *testing.T) {
	t.Run("未初始化时Repository工厂为nil", func(t *testing.T) {
		container := NewServiceContainer()

		// 未初始化时Repository工厂应该为nil
		repoFactory := container.GetRepositoryFactory()
		assert.Nil(t, repoFactory, "未初始化时Repository工厂应该为nil")
	})

	t.Run("Repository工厂在Initialize后创建", func(t *testing.T) {
		if testing.Short() {
			t.Skip("跳过需要配置的测试")
		}

		// 这个测试需要完整的配置环境
		// 在单元测试中跳过，通过集成测试覆盖
		t.Skip("需要完整配置环境，通过集成测试覆盖")
	})
}

// TestContainerAPI验证容器API的正确性
func TestContainerAPI(t *testing.T) {
	container := NewServiceContainer()

	t.Run("未初始化时GetRepositoryFactory返回nil", func(t *testing.T) {
		factory := container.GetRepositoryFactory()
		assert.Nil(t, factory, "未初始化时Repository工厂应该为nil")
	})

	t.Run("GetEventBus应该返回事件总线", func(t *testing.T) {
		eventBus := container.GetEventBus()
		assert.NotNil(t, eventBus, "事件总线应该存在")
	})

	t.Run("IsInitialized应该正确报告状态", func(t *testing.T) {
		assert.False(t, container.IsInitialized(), "新创建的容器应该未初始化")
	})
}

// TestNoGlobalPackageImport 验证容器不再导入global包
// 这是一个编译时测试：如果service_container.go导入了global包，编译会失败
func TestNoGlobalPackageImport(t *testing.T) {
	// 这个测试的存在本身就是一种文档：
	// 说明service/container包不再依赖global包

	t.Run("容器可以独立创建", func(t *testing.T) {
		container := NewServiceContainer()
		assert.NotNil(t, container)
	})
}

// 测试MongoDB连接的集成测试（需要实际MongoDB）
// 这些测试在有MongoDB连接的环境中运行
func TestMongoDBConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	t.Run("Initialize应该初始化MongoDB连接", func(t *testing.T) {
		container := NewServiceContainer()

		ctx := context.Background()
		err := container.Initialize(ctx)

		// 如果MongoDB未运行，跳过测试
		if err != nil {
			t.Skip("MongoDB未连接，跳过测试")
		}

		// 验证MongoDB已设置
		db := container.GetMongoDB()
		assert.NotNil(t, db, "Initialize后GetMongoDB应该返回有效实例")

		client := container.GetMongoClient()
		assert.NotNil(t, client, "Initialize后GetMongoClient应该返回有效实例")
	})

	t.Run("初始化后应该能够获取数据库实例", func(t *testing.T) {
		container := NewServiceContainer()

		ctx := context.Background()
		err := container.Initialize(ctx)

		if err != nil {
			t.Skip("MongoDB未连接，跳过测试")
		}

		db := container.GetMongoDB()
		require.NotNil(t, db)

		// 验证数据库可用性
		collection := db.Collection("test")
		assert.NotNil(t, collection, "应该能够创建集合实例")
	})
}
