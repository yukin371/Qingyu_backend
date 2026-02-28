package finance

import (
	financeModel "Qingyu_backend/models/finance"
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
	now := time.Now()
	wallet.CreatedAt = now
	wallet.UpdatedAt = now

	result, err := r.walletCollection.InsertOne(ctx, wallet)
	if err != nil {
		return fmt.Errorf("创建钱包失败: %w", err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		wallet.ID = oid.Hex()
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

// ============ 交易管理 ============

// CreateTransaction 创建交易记录
func (r *WalletRepositoryImpl) CreateTransaction(ctx context.Context, transaction *financeModel.Transaction) error {
	now := time.Now()
	transaction.CreatedAt = now

	result, err := r.transactionCollection.InsertOne(ctx, transaction)
	if err != nil {
		return fmt.Errorf("创建交易记录失败: %w", err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		transaction.ID = oid.Hex()
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
	now := time.Now()
	request.CreatedAt = now
	request.UpdatedAt = now

	result, err := r.withdrawRequestCollection.InsertOne(ctx, request)
	if err != nil {
		return fmt.Errorf("创建提现请求失败: %w", err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		request.ID = oid.Hex()
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
