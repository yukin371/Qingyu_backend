package finance

import (
	financeModel "Qingyu_backend/models/finance"
	"Qingyu_backend/models/shared/types"
	financeInterface "Qingyu_backend/repository/interfaces/finance"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func sanitizeFinanceUserID(userID string) (string, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return "", fmt.Errorf("user_id格式不合法")
	}
	return objectID.Hex(), nil
}

// WalletRepositoryImpl 钱包Repository实现
type WalletRepositoryImpl struct {
	db                        *mongo.Database
	walletCollection          *mongo.Collection
	transactionCollection     *mongo.Collection
	withdrawRequestCollection *mongo.Collection
}

// NewWalletRepository 创建钱包Repository
func NewWalletRepository(db *mongo.Database) financeInterface.WalletRepository {
	return &WalletRepositoryImpl{
		db:                        db,
		walletCollection:          db.Collection("wallets"),
		transactionCollection:     db.Collection("transactions"),
		withdrawRequestCollection: db.Collection("withdraw_requests"),
	}
}

// ============ 钱包管理 ============

// CreateWallet 创建钱包
func (r *WalletRepositoryImpl) CreateWallet(ctx context.Context, wallet *financeModel.Wallet) error {
	if wallet.ID.IsZero() {
		wallet.ID = primitive.NewObjectID()
	}
	now := time.Now()
	wallet.CreatedAt = now
	wallet.UpdatedAt = now

	_, err := r.walletCollection.InsertOne(ctx, wallet)
	if err != nil {
		return fmt.Errorf("创建钱包失败: %w", err)
	}

	return nil
}

// GetWallet 获取钱包（根据用户ID）
// 注意：这是接口定义的方法，使用userID作为参数
func (r *WalletRepositoryImpl) GetWallet(ctx context.Context, userID string) (*financeModel.Wallet, error) {
	safeUserID, err := sanitizeFinanceUserID(userID)
	if err != nil {
		return nil, err
	}

	var wallet financeModel.Wallet
	err = r.walletCollection.FindOne(ctx, bson.M{"user_id": safeUserID}).Decode(&wallet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("钱包不存在")
		}
		return nil, fmt.Errorf("查询钱包失败: %w", err)
	}

	return &wallet, nil
}

// GetWalletByID 根据钱包ID获取钱包（内部使用）
// 这是一个额外的辅助方法，用于通过钱包ObjectID查询
func (r *WalletRepositoryImpl) GetWalletByID(ctx context.Context, walletID string) (*financeModel.Wallet, error) {
	objectID, err := primitive.ObjectIDFromHex(walletID)
	if err != nil {
		return nil, fmt.Errorf("无效的钱包ID: %w", err)
	}

	var wallet financeModel.Wallet
	err = r.walletCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&wallet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("钱包不存在: %s", walletID)
		}
		return nil, fmt.Errorf("查询钱包失败: %w", err)
	}

	return &wallet, nil
}

// UpdateWallet 更新钱包
func (r *WalletRepositoryImpl) UpdateWallet(ctx context.Context, walletID string, updates map[string]interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(walletID)
	if err != nil {
		return fmt.Errorf("无效的钱包ID: %w", err)
	}

	updates["updated_at"] = time.Now()

	result, err := r.walletCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updates},
	)
	if err != nil {
		return fmt.Errorf("更新钱包失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("钱包不存在: %s", walletID)
	}

	return nil
}

// UpdateBalance 更新余额（原子操作）
// 注意：此方法不进行余额验证，仅用于确定安全的操作（如充值）
// 对于可能减少余额的操作，应使用 UpdateBalanceWithCheck
func (r *WalletRepositoryImpl) UpdateBalance(ctx context.Context, userID string, amount int64) error {
	safeUserID, err := sanitizeFinanceUserID(userID)
	if err != nil {
		return err
	}

	// 使用user_id查询钱包
	result, err := r.walletCollection.UpdateOne(
		ctx,
		bson.M{"user_id": safeUserID},
		bson.M{
			"$inc": bson.M{"balance": amount},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)
	if err != nil {
		return fmt.Errorf("更新余额失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("钱包不存在")
	}

	return nil
}

// UpdateBalanceWithCheck 更新余额并验证（防止负数余额）
// 这是带有余额验证的原子操作，用于扣款等可能减少余额的场景
// 如果扣款后余额会变成负数，则操作会失败
func (r *WalletRepositoryImpl) UpdateBalanceWithCheck(ctx context.Context, userID string, amount int64) error {
	safeUserID, err := sanitizeFinanceUserID(userID)
	if err != nil {
		return err
	}

	// 如果是增加余额（正数），直接使用 UpdateBalance
	if amount >= 0 {
		return r.UpdateBalance(ctx, userID, amount)
	}

	// 对于减少余额的操作，使用条件更新确保余额不会变成负数
	// 条件：当前余额 + amount >= 0
	result, err := r.walletCollection.UpdateOne(
		ctx,
		bson.M{
			"user_id": safeUserID,
			"balance": bson.M{"$gte": -amount}, // 确保扣款后余额 >= 0
		},
		bson.M{
			"$inc": bson.M{"balance": amount},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)
	if err != nil {
		return fmt.Errorf("更新余额失败: %w", err)
	}

	if result.MatchedCount == 0 {
		wallet, err := r.GetWallet(ctx, userID)
		if err != nil {
			return err
		}
		if wallet.Balance < types.Money(-amount) {
			return fmt.Errorf("余额不足: 当前余额 %d, 尝试扣款 %d", wallet.Balance, -amount)
		}
		return fmt.Errorf("钱包不存在")
	}

	return nil
}

// ============ 交易管理 ============

// CreateTransaction 创建交易记录
func (r *WalletRepositoryImpl) CreateTransaction(ctx context.Context, transaction *financeModel.Transaction) error {
	if transaction.ID.IsZero() {
		transaction.ID = primitive.NewObjectID()
	}
	now := time.Now()
	transaction.CreatedAt = now

	_, err := r.transactionCollection.InsertOne(ctx, transaction)
	if err != nil {
		return fmt.Errorf("创建交易记录失败: %w", err)
	}

	return nil
}

// GetTransaction 获取交易记录
func (r *WalletRepositoryImpl) GetTransaction(ctx context.Context, transactionID string) (*financeModel.Transaction, error) {
	objectID, err := primitive.ObjectIDFromHex(transactionID)
	if err != nil {
		return nil, fmt.Errorf("无效的交易ID: %w", err)
	}

	var transaction financeModel.Transaction
	err = r.transactionCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&transaction)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("交易记录不存在: %s", transactionID)
		}
		return nil, fmt.Errorf("查询交易记录失败: %w", err)
	}

	return &transaction, nil
}

// ListTransactions 列出交易记录
func (r *WalletRepositoryImpl) ListTransactions(ctx context.Context, filter *financeInterface.TransactionFilter) ([]*financeModel.Transaction, error) {
	query := bson.M{}

	if filter.UserID != "" {
		query["user_id"] = filter.UserID
	}
	if filter.Type != "" {
		query["type"] = filter.Type
	}
	if filter.Status != "" {
		query["status"] = filter.Status
	}
	if !filter.StartDate.IsZero() {
		query["created_at"] = bson.M{"$gte": filter.StartDate}
	}
	if !filter.EndDate.IsZero() {
		if query["created_at"] != nil {
			query["created_at"].(bson.M)["$lte"] = filter.EndDate
		} else {
			query["created_at"] = bson.M{"$lte": filter.EndDate}
		}
	}

	opts := options.Find().
		SetSort(bson.D{bson.E{Key: "created_at", Value: -1}}).
		SetLimit(filter.Limit)

	if filter.Offset > 0 {
		opts.SetSkip(filter.Offset)
	}

	cursor, err := r.transactionCollection.Find(ctx, query, opts)
	if err != nil {
		return nil, fmt.Errorf("查询交易列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var transactions []*financeModel.Transaction
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, fmt.Errorf("解析交易列表失败: %w", err)
	}

	return transactions, nil
}

// CountTransactions 统计交易数量
func (r *WalletRepositoryImpl) CountTransactions(ctx context.Context, filter *financeInterface.TransactionFilter) (int64, error) {
	query := bson.M{}

	if filter.UserID != "" {
		query["user_id"] = filter.UserID
	}
	if filter.Type != "" {
		query["type"] = filter.Type
	}
	if filter.Status != "" {
		query["status"] = filter.Status
	}
	if !filter.StartDate.IsZero() {
		query["created_at"] = bson.M{"$gte": filter.StartDate}
	}
	if !filter.EndDate.IsZero() {
		if query["created_at"] != nil {
			query["created_at"].(bson.M)["$lte"] = filter.EndDate
		} else {
			query["created_at"] = bson.M{"$lte": filter.EndDate}
		}
	}

	count, err := r.transactionCollection.CountDocuments(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("统计交易数量失败: %w", err)
	}

	return count, nil
}

// ============ 提现管理 ============

// CreateWithdrawRequest 创建提现请求
func (r *WalletRepositoryImpl) CreateWithdrawRequest(ctx context.Context, request *financeModel.WithdrawRequest) error {
	if request.ID.IsZero() {
		request.ID = primitive.NewObjectID()
	}
	now := time.Now()
	request.CreatedAt = now
	request.UpdatedAt = now

	_, err := r.withdrawRequestCollection.InsertOne(ctx, request)
	if err != nil {
		return fmt.Errorf("创建提现请求失败: %w", err)
	}

	return nil
}

// GetWithdrawRequest 获取提现请求
func (r *WalletRepositoryImpl) GetWithdrawRequest(ctx context.Context, requestID string) (*financeModel.WithdrawRequest, error) {
	objectID, err := primitive.ObjectIDFromHex(requestID)
	if err != nil {
		return nil, fmt.Errorf("无效的请求ID: %w", err)
	}

	var request financeModel.WithdrawRequest
	err = r.withdrawRequestCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&request)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("提现请求不存在: %s", requestID)
		}
		return nil, fmt.Errorf("查询提现请求失败: %w", err)
	}

	return &request, nil
}

// UpdateWithdrawRequest 更新提现请求
func (r *WalletRepositoryImpl) UpdateWithdrawRequest(ctx context.Context, requestID string, updates map[string]interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(requestID)
	if err != nil {
		return fmt.Errorf("无效的请求ID: %w", err)
	}

	updates["updated_at"] = time.Now()

	result, err := r.withdrawRequestCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updates},
	)
	if err != nil {
		return fmt.Errorf("更新提现请求失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("提现请求不存在: %s", requestID)
	}

	return nil
}

// ListWithdrawRequests 列出提现请求
func (r *WalletRepositoryImpl) ListWithdrawRequests(ctx context.Context, filter *financeInterface.WithdrawFilter) ([]*financeModel.WithdrawRequest, error) {
	query := bson.M{}

	if filter.UserID != "" {
		query["user_id"] = filter.UserID
	}
	if filter.Status != "" {
		query["status"] = filter.Status
	}
	if !filter.StartDate.IsZero() {
		query["created_at"] = bson.M{"$gte": filter.StartDate}
	}
	if !filter.EndDate.IsZero() {
		if query["created_at"] != nil {
			query["created_at"].(bson.M)["$lte"] = filter.EndDate
		} else {
			query["created_at"] = bson.M{"$lte": filter.EndDate}
		}
	}

	opts := options.Find().
		SetSort(bson.D{bson.E{Key: "created_at", Value: -1}}).
		SetLimit(filter.Limit)

	if filter.Offset > 0 {
		opts.SetSkip(filter.Offset)
	}

	cursor, err := r.withdrawRequestCollection.Find(ctx, query, opts)
	if err != nil {
		return nil, fmt.Errorf("查询提现列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var requests []*financeModel.WithdrawRequest
	if err = cursor.All(ctx, &requests); err != nil {
		return nil, fmt.Errorf("解析提现列表失败: %w", err)
	}

	return requests, nil
}

// CountWithdrawRequests 统计提现请求数量
func (r *WalletRepositoryImpl) CountWithdrawRequests(ctx context.Context, filter *financeInterface.WithdrawFilter) (int64, error) {
	query := bson.M{}

	if filter.UserID != "" {
		query["user_id"] = filter.UserID
	}
	if filter.Status != "" {
		query["status"] = filter.Status
	}
	if !filter.StartDate.IsZero() {
		query["created_at"] = bson.M{"$gte": filter.StartDate}
	}
	if !filter.EndDate.IsZero() {
		if query["created_at"] != nil {
			query["created_at"].(bson.M)["$lte"] = filter.EndDate
		} else {
			query["created_at"] = bson.M{"$lte": filter.EndDate}
		}
	}

	count, err := r.withdrawRequestCollection.CountDocuments(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("统计提现请求数量失败: %w", err)
	}

	return count, nil
}

// ============ 健康检查 ============

// Health 健康检查
func (r *WalletRepositoryImpl) Health(ctx context.Context) error {
	return r.db.Client().Ping(ctx, nil)
}

// isReplicaSet 检查 MongoDB 是否是副本集（支持事务）
func (r *WalletRepositoryImpl) isReplicaSet(ctx context.Context) bool {
	client := r.db.Client()
	adminDB := client.Database("admin")

	var result struct {
		Ok      int    `bson:"ok"`
		SetName string `bson:"setName,omitempty"`
	}

	err := adminDB.RunCommand(ctx, map[string]interface{}{"isMaster": 1}).Decode(&result)
	if err != nil {
		return false
	}

	return result.SetName != ""
}

// RunInTransaction 在事务中执行钱包相关操作
// 如果 MongoDB 不支持事务（单节点），则直接执行操作（降级模式）
func (r *WalletRepositoryImpl) RunInTransaction(ctx context.Context, fn func(context.Context) error) error {
	// 检查是否支持事务
	if !r.isReplicaSet(ctx) {
		// 降级模式：直接执行操作，不使用事务
		return fn(ctx)
	}

	// 副本集模式：使用事务
	session, err := r.db.Client().StartSession()
	if err != nil {
		return fmt.Errorf("启动事务失败: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		if err := fn(sessCtx); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return fmt.Errorf("事务执行失败: %w", err)
	}

	return nil
}
