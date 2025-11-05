package shared

import (
	"context"
	"testing"

	"Qingyu_backend/service/shared/admin"
	"Qingyu_backend/service/shared/messaging"
	"Qingyu_backend/service/shared/recommendation"
	"Qingyu_backend/service/shared/storage"
)

// TestAdminServiceBaseService 测试AdminService的BaseService接口实现
func TestAdminServiceBaseService(t *testing.T) {
	// 创建服务实例（使用nil依赖，仅测试BaseService方法）
	svc := &admin.AdminServiceImpl{}

	ctx := context.Background()

	// 测试GetServiceName
	if name := svc.GetServiceName(); name != "AdminService" {
		t.Errorf("GetServiceName() = %v, want %v", name, "AdminService")
	}

	// 测试GetVersion
	if version := svc.GetVersion(); version != "v1.0.0" {
		t.Errorf("GetVersion() = %v, want %v", version, "v1.0.0")
	}

	// 测试未初始化时的Health检查应失败
	if err := svc.Health(ctx); err == nil {
		t.Error("Health() should return error when not initialized")
	}

	// 测试Close（即使未初始化也应该成功）
	if err := svc.Close(ctx); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

// TestStorageServiceBaseService 测试StorageService的BaseService接口实现
func TestStorageServiceBaseService(t *testing.T) {
	svc := &storage.StorageServiceImpl{}

	ctx := context.Background()

	// 测试GetServiceName
	if name := svc.GetServiceName(); name != "StorageService" {
		t.Errorf("GetServiceName() = %v, want %v", name, "StorageService")
	}

	// 测试GetVersion
	if version := svc.GetVersion(); version != "v1.0.0" {
		t.Errorf("GetVersion() = %v, want %v", version, "v1.0.0")
	}

	// 测试未初始化时的Health检查应失败
	if err := svc.Health(ctx); err == nil {
		t.Error("Health() should return error when not initialized")
	}

	// 测试Close
	if err := svc.Close(ctx); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

// TestMessagingServiceBaseService 测试MessagingService的BaseService接口实现
func TestMessagingServiceBaseService(t *testing.T) {
	svc := &messaging.MessagingServiceImpl{}

	ctx := context.Background()

	// 测试GetServiceName
	if name := svc.GetServiceName(); name != "MessagingService" {
		t.Errorf("GetServiceName() = %v, want %v", name, "MessagingService")
	}

	// 测试GetVersion
	if version := svc.GetVersion(); version != "v1.0.0" {
		t.Errorf("GetVersion() = %v, want %v", version, "v1.0.0")
	}

	// 测试未初始化时的Health检查应失败
	if err := svc.Health(ctx); err == nil {
		t.Error("Health() should return error when not initialized")
	}

	// 测试Close
	if err := svc.Close(ctx); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

// TestRecommendationServiceBaseService 测试RecommendationService的BaseService接口实现
func TestRecommendationServiceBaseService(t *testing.T) {
	svc := &recommendation.RecommendationServiceImpl{}

	ctx := context.Background()

	// 测试GetServiceName
	if name := svc.GetServiceName(); name != "RecommendationService" {
		t.Errorf("GetServiceName() = %v, want %v", name, "RecommendationService")
	}

	// 测试GetVersion
	if version := svc.GetVersion(); version != "v1.0.0" {
		t.Errorf("GetVersion() = %v, want %v", version, "v1.0.0")
	}

	// 测试未初始化时的Health检查应失败
	if err := svc.Health(ctx); err == nil {
		t.Error("Health() should return error when not initialized")
	}

	// 测试Close
	if err := svc.Close(ctx); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

// TestAllServicesImplementBaseService 测试所有服务都实现了BaseService接口
func TestAllServicesImplementBaseService(t *testing.T) {
	services := []interface{}{
		&admin.AdminServiceImpl{},
		&storage.StorageServiceImpl{},
		&messaging.MessagingServiceImpl{},
		&recommendation.RecommendationServiceImpl{},
	}

	for _, svc := range services {
		// 检查是否实现了必要的方法
		if _, ok := svc.(interface{ GetServiceName() string }); !ok {
			t.Errorf("%T does not implement GetServiceName()", svc)
		}
		if _, ok := svc.(interface{ GetVersion() string }); !ok {
			t.Errorf("%T does not implement GetVersion()", svc)
		}
		if _, ok := svc.(interface {
			Initialize(context.Context) error
		}); !ok {
			t.Errorf("%T does not implement Initialize()", svc)
		}
		if _, ok := svc.(interface {
			Health(context.Context) error
		}); !ok {
			t.Errorf("%T does not implement Health()", svc)
		}
		if _, ok := svc.(interface {
			Close(context.Context) error
		}); !ok {
			t.Errorf("%T does not implement Close()", svc)
		}
	}
}
