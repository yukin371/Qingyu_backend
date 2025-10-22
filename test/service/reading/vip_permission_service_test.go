package reading

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"Qingyu_backend/service/reading"
)

// setupVIPPermissionMock 创建VIP权限服务和Redis Mock
func setupVIPPermissionMock(prefix string) (reading.VIPPermissionService, redismock.ClientMock) {
	db, mock := redismock.NewClientMock()
	service := reading.NewVIPPermissionService(db, prefix)
	return service, mock
}

// TestVIPPermissionService_CheckVIPAccess_FreeChapter 测试非VIP章节访问
func TestVIPPermissionService_CheckVIPAccess_FreeChapter(t *testing.T) {
	service, mock := setupVIPPermissionMock("qingyu")
	ctx := context.Background()

	// 执行测试：非VIP章节，任何人都可以访问
	canAccess, err := service.CheckVIPAccess(ctx, "user123", "chapter1", false)

	// 验证结果
	assert.NoError(t, err)
	assert.True(t, canAccess)

	// 验证所有Redis期望都已满足
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestVIPPermissionService_CheckVIPAccess_VIPUser 测试VIP用户访问VIP章节
func TestVIPPermissionService_CheckVIPAccess_VIPUser(t *testing.T) {
	service, mock := setupVIPPermissionMock("qingyu")
	ctx := context.Background()

	// 设置Mock期望：用户是VIP
	mock.ExpectExists("qingyu:vip:user:user123:status").SetVal(1)

	// 执行测试
	canAccess, err := service.CheckVIPAccess(ctx, "user123", "chapter1", true)

	// 验证结果
	assert.NoError(t, err)
	assert.True(t, canAccess)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestVIPPermissionService_CheckVIPAccess_PurchasedChapter 测试已购买章节访问
func TestVIPPermissionService_CheckVIPAccess_PurchasedChapter(t *testing.T) {
	service, mock := setupVIPPermissionMock("qingyu")
	ctx := context.Background()

	// 设置Mock期望：用户不是VIP，但已购买章节
	mock.ExpectExists("qingyu:vip:user:user123:status").SetVal(0)
	mock.ExpectSIsMember("qingyu:vip:purchase:user123:chapters", "chapter1").SetVal(true)

	// 执行测试
	canAccess, err := service.CheckVIPAccess(ctx, "user123", "chapter1", true)

	// 验证结果
	assert.NoError(t, err)
	assert.True(t, canAccess)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestVIPPermissionService_CheckVIPAccess_Denied 测试无权限访问
func TestVIPPermissionService_CheckVIPAccess_Denied(t *testing.T) {
	service, mock := setupVIPPermissionMock("qingyu")
	ctx := context.Background()

	// 设置Mock期望：用户不是VIP，也未购买章节
	mock.ExpectExists("qingyu:vip:user:user456:status").SetVal(0)
	mock.ExpectSIsMember("qingyu:vip:purchase:user456:chapters", "chapter1").SetVal(false)

	// 执行测试
	canAccess, err := service.CheckVIPAccess(ctx, "user456", "chapter1", true)

	// 验证结果
	assert.NoError(t, err)
	assert.False(t, canAccess)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestVIPPermissionService_CheckUserVIPStatus_IsVIP 测试VIP用户状态检查
func TestVIPPermissionService_CheckUserVIPStatus_IsVIP(t *testing.T) {
	service, mock := setupVIPPermissionMock("qingyu")
	ctx := context.Background()

	// 设置Mock期望
	mock.ExpectExists("qingyu:vip:user:user123:status").SetVal(1)

	// 执行测试
	isVIP, err := service.CheckUserVIPStatus(ctx, "user123")

	// 验证结果
	assert.NoError(t, err)
	assert.True(t, isVIP)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestVIPPermissionService_CheckUserVIPStatus_NotVIP 测试普通用户状态检查
func TestVIPPermissionService_CheckUserVIPStatus_NotVIP(t *testing.T) {
	service, mock := setupVIPPermissionMock("qingyu")
	ctx := context.Background()

	// 设置Mock期望
	mock.ExpectExists("qingyu:vip:user:user456:status").SetVal(0)

	// 执行测试
	isVIP, err := service.CheckUserVIPStatus(ctx, "user456")

	// 验证结果
	assert.NoError(t, err)
	assert.False(t, isVIP)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestVIPPermissionService_CheckChapterPurchased_Yes 测试已购买章节
func TestVIPPermissionService_CheckChapterPurchased_Yes(t *testing.T) {
	service, mock := setupVIPPermissionMock("qingyu")
	ctx := context.Background()

	// 设置Mock期望
	mock.ExpectSIsMember("qingyu:vip:purchase:user123:chapters", "chapter1").SetVal(true)

	// 执行测试
	purchased, err := service.CheckChapterPurchased(ctx, "user123", "chapter1")

	// 验证结果
	assert.NoError(t, err)
	assert.True(t, purchased)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestVIPPermissionService_CheckChapterPurchased_No 测试未购买章节
func TestVIPPermissionService_CheckChapterPurchased_No(t *testing.T) {
	service, mock := setupVIPPermissionMock("qingyu")
	ctx := context.Background()

	// 设置Mock期望
	mock.ExpectSIsMember("qingyu:vip:purchase:user456:chapters", "chapter1").SetVal(false)

	// 执行测试
	purchased, err := service.CheckChapterPurchased(ctx, "user456", "chapter1")

	// 验证结果
	assert.NoError(t, err)
	assert.False(t, purchased)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestVIPPermissionService_GrantVIPAccess 测试授予VIP权限
func TestVIPPermissionService_GrantVIPAccess(t *testing.T) {
	service, mock := setupVIPPermissionMock("qingyu")
	ctx := context.Background()

	duration := 30 * 24 * time.Hour // 30天

	// 设置Mock期望
	mock.ExpectSet("qingyu:vip:user:user123:status", "vip", duration).SetVal("OK")

	// 执行测试
	err := service.GrantVIPAccess(ctx, "user123", duration)

	// 验证结果
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestVIPPermissionService_GrantChapterAccess 测试授予章节访问权限
func TestVIPPermissionService_GrantChapterAccess(t *testing.T) {
	service, mock := setupVIPPermissionMock("qingyu")
	ctx := context.Background()

	// 设置Mock期望
	mock.ExpectSAdd("qingyu:vip:purchase:user123:chapters", "chapter1").SetVal(1)
	mock.ExpectExpire("qingyu:vip:purchase:user123:chapters", 365*24*time.Hour).SetVal(true)

	// 执行测试
	err := service.GrantChapterAccess(ctx, "user123", "chapter1")

	// 验证结果
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestVIPPermissionService_GetUserPurchasedChapters 测试获取用户购买的章节列表
func TestVIPPermissionService_GetUserPurchasedChapters(t *testing.T) {
	service, mock := setupVIPPermissionMock("qingyu")
	ctx := context.Background()

	expectedChapters := []string{"chapter1", "chapter2", "chapter3"}

	// 设置Mock期望
	mock.ExpectSMembers("qingyu:vip:purchase:user123:chapters").SetVal(expectedChapters)

	// 执行测试
	impl, ok := service.(*reading.VIPPermissionServiceImpl)
	if !ok {
		t.Fatal("service类型转换失败")
	}
	chapters, err := impl.GetUserPurchasedChapters(ctx, "user123")

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, expectedChapters, chapters)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestVIPPermissionService_GetVIPExpireTime 测试获取VIP过期时间
func TestVIPPermissionService_GetVIPExpireTime(t *testing.T) {
	service, mock := setupVIPPermissionMock("qingyu")
	ctx := context.Background()

	expectedTTL := 15 * 24 * time.Hour // 15天

	// 设置Mock期望
	mock.ExpectTTL("qingyu:vip:user:user123:status").SetVal(expectedTTL)

	// 执行测试
	impl, ok := service.(*reading.VIPPermissionServiceImpl)
	if !ok {
		t.Fatal("service类型转换失败")
	}
	ttl, err := impl.GetVIPExpireTime(ctx, "user123")

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, expectedTTL, ttl)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestVIPPermissionService_RevokeVIPAccess 测试撤销VIP权限
func TestVIPPermissionService_RevokeVIPAccess(t *testing.T) {
	service, mock := setupVIPPermissionMock("qingyu")
	ctx := context.Background()

	// 设置Mock期望
	mock.ExpectDel("qingyu:vip:user:user123:status").SetVal(1)

	// 执行测试
	impl, ok := service.(*reading.VIPPermissionServiceImpl)
	if !ok {
		t.Fatal("service类型转换失败")
	}
	err := impl.RevokeVIPAccess(ctx, "user123")

	// 验证结果
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestVIPPermissionService_RevokeChapterAccess 测试撤销章节访问权限
func TestVIPPermissionService_RevokeChapterAccess(t *testing.T) {
	service, mock := setupVIPPermissionMock("qingyu")
	ctx := context.Background()

	// 设置Mock期望
	mock.ExpectSRem("qingyu:vip:purchase:user123:chapters", "chapter1").SetVal(1)

	// 执行测试
	impl, ok := service.(*reading.VIPPermissionServiceImpl)
	if !ok {
		t.Fatal("service类型转换失败")
	}
	err := impl.RevokeChapterAccess(ctx, "user123", "chapter1")

	// 验证结果
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestVIPPermissionService_CustomPrefix 测试自定义前缀
func TestVIPPermissionService_CustomPrefix(t *testing.T) {
	service, mock := setupVIPPermissionMock("custom")
	ctx := context.Background()

	// 设置Mock期望：使用自定义前缀
	mock.ExpectExists("custom:vip:user:user123:status").SetVal(1)

	// 执行测试
	isVIP, err := service.CheckUserVIPStatus(ctx, "user123")

	// 验证结果
	assert.NoError(t, err)
	assert.True(t, isVIP)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestVIPPermissionService_RedisError 测试Redis错误处理
func TestVIPPermissionService_RedisError(t *testing.T) {
	service, mock := setupVIPPermissionMock("qingyu")
	ctx := context.Background()

	// 设置Mock期望：Redis返回错误
	mock.ExpectExists("qingyu:vip:user:user789:status").SetErr(fmt.Errorf("redis connection error"))

	// 执行测试
	isVIP, err := service.CheckUserVIPStatus(ctx, "user789")

	// 验证结果
	assert.Error(t, err)
	assert.False(t, isVIP)
	assert.Contains(t, err.Error(), "检查VIP状态失败")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestVIPPermissionService_CheckVIPAccess_RedisError 测试CheckVIPAccess时的Redis错误
func TestVIPPermissionService_CheckVIPAccess_RedisError(t *testing.T) {
	service, mock := setupVIPPermissionMock("qingyu")
	ctx := context.Background()

	// 设置Mock期望：Redis返回错误
	mock.ExpectExists("qingyu:vip:user:user789:status").SetErr(redis.Nil)

	// 执行测试
	canAccess, err := service.CheckVIPAccess(ctx, "user789", "chapter1", true)

	// 验证结果
	assert.Error(t, err)
	assert.False(t, canAccess)
	assert.NoError(t, mock.ExpectationsWereMet())
}
