package e2e

import (
	"context"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"Qingyu_backend/global"
	userRepo "Qingyu_backend/repository/mongodb/user"
	"Qingyu_backend/test/e2e/data"
)

// E2EAssertions E2E 专用断言
type E2EAssertions struct {
	env *TestEnvironment
}

// Assert 获取断言辅助器
func (env *TestEnvironment) Assert() *E2EAssertions {
	return &E2EAssertions{env: env}
}

// ============ 用户相关断言 ============

// AssertUserVIPLevel 验证用户VIP等级
func (ea *E2EAssertions) AssertUserVIPLevel(userID string, expectedLevel int) {
	userRepository := userRepo.NewMongoUserRepository(global.DB)
	user, err := userRepository.GetByID(context.Background(), userID)
	require.NoError(ea.env.T, err, "获取用户失败")
	assert.Equal(ea.env.T, expectedLevel, user.VIPLevel, "VIP等级不匹配")

	ea.env.LogSuccess("用户VIP等级验证: %s, Level=%d", userID, user.VIPLevel)
}

// AssertUserIsVIP 验证用户是否为VIP
func (ea *E2EAssertions) AssertUserIsVIP(userID string, expectedIsVIP bool) {
	userRepository := userRepo.NewMongoUserRepository(global.DB)
	user, err := userRepository.GetByID(context.Background(), userID)
	require.NoError(ea.env.T, err, "获取用户失败")
	assert.Equal(ea.env.T, expectedIsVIP, user.IsVIP(), "VIP状态不匹配")

	ea.env.LogSuccess("VIP状态验证: %s, IsVIP=%v", userID, user.IsVIP())
}

// AssertUserExists 验证用户存在
func (ea *E2EAssertions) AssertUserExists(userID string) bool {
	userRepository := userRepo.NewMongoUserRepository(global.DB)
	user, err := userRepository.GetByID(context.Background(), userID)
	require.NoError(ea.env.T, err, "获取用户失败")
	assert.NotNil(ea.env.T, user, "用户不存在")

	return user != nil
}

// ============ HTTP 响应断言 ============

// AssertHTTPSuccess 断言 HTTP 请求成功
func (ea *E2EAssertions) AssertHTTPSuccess(statusCode int, expectedStatus int) {
	assert.Equal(ea.env.T, expectedStatus, statusCode, "HTTP 状态码不匹配")
}

// AssertHTTPError 断言 HTTP 请求失败
func (ea *E2EAssertions) AssertHTTPError(statusCode int, expectedStatus int) {
	assert.Equal(ea.env.T, expectedStatus, statusCode, "HTTP 错误状态码不匹配")
}

// AssertResponseContains 断言响应包含指定字段
func (ea *E2EAssertions) AssertResponseContains(response map[string]interface{}, key string) {
	_, exists := response[key]
	assert.True(ea.env.T, exists, "响应不包含字段: %s", key)
}

// AssertResponseEquals 断言响应字段值等于期望值
func (ea *E2EAssertions) AssertResponseEquals(response map[string]interface{}, key string, expectedValue interface{}) {
	actualValue, exists := response[key]
	require.True(ea.env.T, exists, "响应不包含字段: %s", key)
	assert.Equal(ea.env.T, expectedValue, actualValue, "字段 %s 的值不匹配", key)
}

// ============ 数据库断言 ============

// AssertCollectionCount 验证集合中记录数量
func (ea *E2EAssertions) AssertCollectionCount(collectionName string, expectedCount int64) {
	ctx := context.Background()
	count, err := global.DB.Collection(collectionName).CountDocuments(ctx, map[string]interface{}{})
	require.NoError(ea.env.T, err, "统计集合记录失败")
	assert.Equal(ea.env.T, expectedCount, count, "集合 %s 的记录数不匹配", collectionName)

	ea.env.LogSuccess("集合记录数验证: %s, count=%d", collectionName, count)
}

// AssertDocumentExists 验证文档存在
func (ea *E2EAssertions) AssertDocumentExists(collectionName string, filter map[string]interface{}) bool {
	ctx := context.Background()
	count, err := global.DB.Collection(collectionName).CountDocuments(ctx, filter)
	require.NoError(ea.env.T, err, "查询文档失败")
	exists := count > 0

	assert.True(ea.env.T, exists, "文档不存在: collection=%s, filter=%v", collectionName, filter)

	return exists
}

// ============ 业务流程断言 ============

// AssertReadingProgress 验证阅读进度
func (ea *E2EAssertions) AssertReadingProgress(userID, bookID string) {
	exists := ea.AssertDocumentExists("reading_progress", map[string]interface{}{
		"user_id": userID,
		"book_id": bookID,
	})

	if exists {
		ea.env.LogSuccess("阅读进度验证: user=%s, book=%s", userID, bookID)
	}
}

// AssertCommentExists 验证评论存在
func (ea *E2EAssertions) AssertCommentExists(userID, bookID string) {
	// 评论模型使用 target_id 和 target_type（新结构）
	// 而不是旧的 book_id 字段
	exists := ea.AssertDocumentExists("comments", map[string]interface{}{
		"author_id":   userID,
		"target_id":   bookID,
		"target_type": "book",
	})

	if exists {
		ea.env.LogSuccess("评论记录验证: user=%s, book=%s", userID, bookID)
	}
}

// AssertCollectionExists 验证收藏记录存在
func (ea *E2EAssertions) AssertCollectionExists(userID, bookID string) {
	exists := ea.AssertDocumentExists("collections", map[string]interface{}{
		"user_id": userID,
		"book_id": bookID,
	})

	if exists {
		ea.env.LogSuccess("收藏记录验证: user=%s, book=%s", userID, bookID)
	}
}

// ============ 辅助方法 ============

// AssertSuccess 通用成功断言
func (ea *E2EAssertions) AssertSuccess(condition bool, msg string) {
	assert.True(ea.env.T, condition, msg)
	if condition {
		ea.env.LogSuccess("断言成功: %s", msg)
	}
}

// AssertNoError 通用无错误断言
func (ea *E2EAssertions) AssertNoError(err error, msg string) {
	require.NoError(ea.env.T, err, msg)
}

// ============ 数据一致性验证器包装 ============

// ConsistencyValidatorWrapper 数据一致性验证器包装器
type ConsistencyValidatorWrapper struct {
	env *TestEnvironment
}

// ValidateUserData 验证用户数据一致性
func (cv *ConsistencyValidatorWrapper) ValidateUserData(userID string) []data.ConsistencyIssue {
	validator := data.NewConsistencyValidator(cv.env.T)
	return validator.ValidateUserData(cv.env.T.Context(), userID)
}

// ValidateBookData 验证书籍数据一致性
func (cv *ConsistencyValidatorWrapper) ValidateBookData(bookID string) []data.ConsistencyIssue {
	validator := data.NewConsistencyValidator(cv.env.T)
	return validator.ValidateBookData(cv.env.T.Context(), bookID)
}
