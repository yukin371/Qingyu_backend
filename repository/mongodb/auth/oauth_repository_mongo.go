package auth

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	authModel "Qingyu_backend/models/auth"
	authrepo "Qingyu_backend/repository/interfaces/auth"
)

const (
	OAuthAccountCollection  = "oauth_accounts"
	OAuthSessionCollection  = "oauth_sessions"
)

// MongoOAuthRepository MongoDB OAuth仓储实现
type MongoOAuthRepository struct {
	db *mongo.Database
}

// NewMongoOAuthRepository 创建MongoDB OAuth仓储
func NewMongoOAuthRepository(db *mongo.Database) authrepo.OAuthRepository {
	return &MongoOAuthRepository{db: db}
}

// ==================== OAuth账号管理 ====================

// FindByProviderAndProviderID 根据提供商和提供商用户ID查找OAuth账号
func (r *MongoOAuthRepository) FindByProviderAndProviderID(ctx context.Context, provider authModel.OAuthProvider, providerUserID string) (*authModel.OAuthAccount, error) {
	var account authModel.OAuthAccount
	filter := bson.M{
		"provider":        provider,
		"provider_user_id": providerUserID,
	}

	err := r.db.Collection(OAuthAccountCollection).FindOne(ctx, filter).Decode(&account)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 未找到不返回错误
		}
		return nil, err
	}

	return &account, nil
}

// FindByUserID 根据用户ID查找所有OAuth账号
func (r *MongoOAuthRepository) FindByUserID(ctx context.Context, userID string) ([]*authModel.OAuthAccount, error) {
	filter := bson.M{"user_id": userID}
	cursor, err := r.db.Collection(OAuthAccountCollection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var accounts []*authModel.OAuthAccount
	if err = cursor.All(ctx, &accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}

// FindByID 根据ID查找OAuth账号
func (r *MongoOAuthRepository) FindByID(ctx context.Context, id string) (*authModel.OAuthAccount, error) {
	var account authModel.OAuthAccount
	filter := bson.M{"_id": id}

	err := r.db.Collection(OAuthAccountCollection).FindOne(ctx, filter).Decode(&account)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &account, nil
}

// Create 创建OAuth账号
func (r *MongoOAuthRepository) Create(ctx context.Context, account *authModel.OAuthAccount) error {
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()

	if account.ID == "" {
		account.ID = primitive.NewObjectID().Hex()
	}

	if account.LastLoginAt.IsZero() {
		account.LastLoginAt = time.Now()
	}

	_, err := r.db.Collection(OAuthAccountCollection).InsertOne(ctx, account)
	return err
}

// Update 更新OAuth账号
func (r *MongoOAuthRepository) Update(ctx context.Context, account *authModel.OAuthAccount) error {
	account.UpdatedAt = time.Now()

	filter := bson.M{"_id": account.ID}
	update := bson.M{"$set": account}

	result, err := r.db.Collection(OAuthAccountCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("oauth account not found")
	}

	return nil
}

// Delete 删除OAuth账号
func (r *MongoOAuthRepository) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	result, err := r.db.Collection(OAuthAccountCollection).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("oauth account not found")
	}

	return nil
}

// UpdateLastLogin 更新最后登录时间
func (r *MongoOAuthRepository) UpdateLastLogin(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"last_login_at": time.Now(),
			"updated_at":    time.Now(),
		},
	}

	_, err := r.db.Collection(OAuthAccountCollection).UpdateOne(ctx, filter, update)
	return err
}

// UpdateTokens 更新访问令牌
func (r *MongoOAuthRepository) UpdateTokens(ctx context.Context, id string, accessToken, refreshToken string, expiresAt primitive.DateTime) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"access_token":      accessToken,
			"refresh_token":     refreshToken,
			"token_expires_at":  expiresAt,
			"updated_at":        time.Now(),
		},
	}

	_, err := r.db.Collection(OAuthAccountCollection).UpdateOne(ctx, filter, update)
	return err
}

// SetPrimaryAccount 设置主账号
func (r *MongoOAuthRepository) SetPrimaryAccount(ctx context.Context, userID string, accountID string) error {
	// 取消该用户所有账号的主账号状态
	filter := bson.M{"user_id": userID}
	unsetUpdate := bson.M{"$set": bson.M{"is_primary": false, "updated_at": time.Now()}}
	_, err := r.db.Collection(OAuthAccountCollection).UpdateMany(ctx, filter, unsetUpdate)
	if err != nil {
		return err
	}

	// 设置指定账号为主账号
	filter = bson.M{"_id": accountID, "user_id": userID}
	setUpdate := bson.M{"$set": bson.M{"is_primary": true, "updated_at": time.Now()}}
	result, err := r.db.Collection(OAuthAccountCollection).UpdateOne(ctx, filter, setUpdate)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("oauth account not found or does not belong to user")
	}

	return nil
}

// GetPrimaryAccount 获取用户的主账号
func (r *MongoOAuthRepository) GetPrimaryAccount(ctx context.Context, userID string) (*authModel.OAuthAccount, error) {
	var account authModel.OAuthAccount
	filter := bson.M{
		"user_id":    userID,
		"is_primary": true,
	}

	err := r.db.Collection(OAuthAccountCollection).FindOne(ctx, filter).Decode(&account)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &account, nil
}

// CountByUserID 统计用户的OAuth账号数量
func (r *MongoOAuthRepository) CountByUserID(ctx context.Context, userID string) (int64, error) {
	filter := bson.M{"user_id": userID}
	count, err := r.db.Collection(OAuthAccountCollection).CountDocuments(ctx, filter)
	return count, err
}

// ==================== OAuth会话管理 ====================

// CreateSession 创建OAuth会话
func (r *MongoOAuthRepository) CreateSession(ctx context.Context, session *authModel.OAuthSession) error {
	session.CreatedAt = time.Now()

	if session.ID == "" {
		session.ID = primitive.NewObjectID().Hex()
	}

	_, err := r.db.Collection(OAuthSessionCollection).InsertOne(ctx, session)
	return err
}

// FindSessionByID 根据ID查找OAuth会话
func (r *MongoOAuthRepository) FindSessionByID(ctx context.Context, id string) (*authModel.OAuthSession, error) {
	var session authModel.OAuthSession
	filter := bson.M{"_id": id}

	err := r.db.Collection(OAuthSessionCollection).FindOne(ctx, filter).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &session, nil
}

// FindSessionByState 根据state查找OAuth会话
func (r *MongoOAuthRepository) FindSessionByState(ctx context.Context, state string) (*authModel.OAuthSession, error) {
	var session authModel.OAuthSession
	filter := bson.M{"state": state}

	err := r.db.Collection(OAuthSessionCollection).FindOne(ctx, filter).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &session, nil
}

// DeleteSession 删除OAuth会话
func (r *MongoOAuthRepository) DeleteSession(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	_, err := r.db.Collection(OAuthSessionCollection).DeleteOne(ctx, filter)
	return err
}

// CleanupExpiredSessions 清理过期的OAuth会话
func (r *MongoOAuthRepository) CleanupExpiredSessions(ctx context.Context) (int64, error) {
	filter := bson.M{"expires_at": bson.M{"$lt": time.Now()}}
	result, err := r.db.Collection(OAuthSessionCollection).DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

// EnsureIndexes 创建索引（在初始化时调用）
func (r *MongoOAuthRepository) EnsureIndexes(ctx context.Context) error {
	// OAuth账号索引
	oauthAccountIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "provider", Value: 1},
				{Key: "provider_user_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "is_primary", Value: 1},
			},
		},
	}

	_, err := r.db.Collection(OAuthAccountCollection).Indexes().CreateMany(ctx, oauthAccountIndexes)
	if err != nil {
		return err
	}

	// OAuth会话索引
	oauthSessionIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "state", Value: 1}},
			Options: options.Index().
				SetUnique(true).
				SetExpireAfterSeconds(600), // 10分钟自动过期
		},
		{
			Keys: bson.D{{Key: "expires_at", Value: 1}},
		},
	}

	_, err = r.db.Collection(OAuthSessionCollection).Indexes().CreateMany(ctx, oauthSessionIndexes)
	return err
}
