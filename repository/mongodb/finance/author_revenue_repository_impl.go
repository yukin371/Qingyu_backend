package finance

import (
	financeModel "Qingyu_backend/models/finance"
	"context"
	"fmt"
	"time"

	financeInterfaces "Qingyu_backend/repository/interfaces/finance"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AuthorRevenueRepositoryImpl 作者收入Repository实现
type AuthorRevenueRepositoryImpl struct {
	db                  *mongo.Database
	earningCollection   *mongo.Collection
	withdrawalCollection *mongo.Collection
	settlementCollection *mongo.Collection
	detailCollection    *mongo.Collection
	statisticsCollection *mongo.Collection
	taxInfoCollection   *mongo.Collection
}

// NewAuthorRevenueRepository 创建作者收入Repository
func NewAuthorRevenueRepository(db *mongo.Database) financeInterfaces.AuthorRevenueRepository {
	return &AuthorRevenueRepositoryImpl{
		db:                  db,
		earningCollection:   db.Collection("author_earnings"),
		withdrawalCollection: db.Collection("withdrawal_requests"),
		settlementCollection: db.Collection("settlements"),
		detailCollection:    db.Collection("revenue_details"),
		statisticsCollection: db.Collection("revenue_statistics"),
		taxInfoCollection:   db.Collection("tax_info"),
	}
}

// ============ 收入记录管理 ============

// CreateEarning 创建收入记录
func (r *AuthorRevenueRepositoryImpl) CreateEarning(ctx context.Context, earning *financeModel.AuthorEarning) error {
	now := time.Now()
	earning.CreatedAt = now
	earning.UpdatedAt = now

	result, err := r.earningCollection.InsertOne(ctx, earning)
	if err != nil {
		return fmt.Errorf("创建收入记录失败: %w", err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		earning.ID = oid
	}

	return nil
}

// GetEarning 获取收入记录
func (r *AuthorRevenueRepositoryImpl) GetEarning(ctx context.Context, earningID primitive.ObjectID) (*financeModel.AuthorEarning, error) {
	var earning financeModel.AuthorEarning
	err := r.earningCollection.FindOne(ctx, bson.M{"_id": earningID}).Decode(&earning)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("收入记录不存在: %s", earningID.Hex())
		}
		return nil, fmt.Errorf("查询收入记录失败: %w", err)
	}

	return &earning, nil
}

// ListEarnings 列出收入记录
func (r *AuthorRevenueRepositoryImpl) ListEarnings(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]*financeModel.AuthorEarning, int64, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize))

	cursor, err := r.earningCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("查询收入记录列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var earnings []*financeModel.AuthorEarning
	if err := cursor.All(ctx, &earnings); err != nil {
		return nil, 0, fmt.Errorf("解析收入记录列表失败: %w", err)
	}

	// 获取总数
	total, err := r.earningCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("统计收入记录失败: %w", err)
	}

	return earnings, total, nil
}

// UpdateEarning 更新收入记录
func (r *AuthorRevenueRepositoryImpl) UpdateEarning(ctx context.Context, earningID primitive.ObjectID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	result, err := r.earningCollection.UpdateOne(ctx, bson.M{"_id": earningID}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("更新收入记录失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("收入记录不存在: %s", earningID.Hex())
	}

	return nil
}

// BatchUpdateEarnings 批量更新收入记录
func (r *AuthorRevenueRepositoryImpl) BatchUpdateEarnings(ctx context.Context, earningIDs []primitive.ObjectID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	_, err := r.earningCollection.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": earningIDs}}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("批量更新收入记录失败: %w", err)
	}

	return nil
}

// GetEarningsByAuthor 获取作者收入记录
func (r *AuthorRevenueRepositoryImpl) GetEarningsByAuthor(ctx context.Context, authorID string, page, pageSize int) ([]*financeModel.AuthorEarning, int64, error) {
	filter := bson.M{"author_id": authorID}
	return r.ListEarnings(ctx, filter, page, pageSize)
}

// GetEarningsByBook 获取书籍收入记录
func (r *AuthorRevenueRepositoryImpl) GetEarningsByBook(ctx context.Context, bookID primitive.ObjectID, page, pageSize int) ([]*financeModel.AuthorEarning, int64, error) {
	filter := bson.M{"book_id": bookID}
	return r.ListEarnings(ctx, filter, page, pageSize)
}

// ============ 提现申请管理 ============

// CreateWithdrawalRequest 创建提现申请
func (r *AuthorRevenueRepositoryImpl) CreateWithdrawalRequest(ctx context.Context, request *financeModel.WithdrawalRequest) error {
	now := time.Now()
	request.CreatedAt = now
	request.UpdatedAt = now

	result, err := r.withdrawalCollection.InsertOne(ctx, request)
	if err != nil {
		return fmt.Errorf("创建提现申请失败: %w", err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		request.ID = oid
	}

	return nil
}

// GetWithdrawalRequest 获取提现申请
func (r *AuthorRevenueRepositoryImpl) GetWithdrawalRequest(ctx context.Context, requestID primitive.ObjectID) (*financeModel.WithdrawalRequest, error) {
	var request financeModel.WithdrawalRequest
	err := r.withdrawalCollection.FindOne(ctx, bson.M{"_id": requestID}).Decode(&request)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("提现申请不存在: %s", requestID.Hex())
		}
		return nil, fmt.Errorf("查询提现申请失败: %w", err)
	}

	return &request, nil
}

// ListWithdrawalRequests 列出提现申请
func (r *AuthorRevenueRepositoryImpl) ListWithdrawalRequests(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]*financeModel.WithdrawalRequest, int64, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize))

	cursor, err := r.withdrawalCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("查询提现申请列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var requests []*financeModel.WithdrawalRequest
	if err := cursor.All(ctx, &requests); err != nil {
		return nil, 0, fmt.Errorf("解析提现申请列表失败: %w", err)
	}

	// 获取总数
	total, err := r.withdrawalCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("统计提现申请失败: %w", err)
	}

	return requests, total, nil
}

// UpdateWithdrawalRequest 更新提现申请
func (r *AuthorRevenueRepositoryImpl) UpdateWithdrawalRequest(ctx context.Context, requestID primitive.ObjectID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	result, err := r.withdrawalCollection.UpdateOne(ctx, bson.M{"_id": requestID}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("更新提现申请失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("提现申请不存在: %s", requestID.Hex())
	}

	return nil
}

// GetUserWithdrawalRequests 获取用户提现申请
func (r *AuthorRevenueRepositoryImpl) GetUserWithdrawalRequests(ctx context.Context, userID string, page, pageSize int) ([]*financeModel.WithdrawalRequest, int64, error) {
	filter := bson.M{"user_id": userID}
	return r.ListWithdrawalRequests(ctx, filter, page, pageSize)
}

// ============ 结算管理 ============

// CreateSettlement 创建结算记录
func (r *AuthorRevenueRepositoryImpl) CreateSettlement(ctx context.Context, settlement *financeModel.Settlement) error {
	now := time.Now()
	settlement.CreatedAt = now
	settlement.UpdatedAt = now

	result, err := r.settlementCollection.InsertOne(ctx, settlement)
	if err != nil {
		return fmt.Errorf("创建结算记录失败: %w", err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		settlement.ID = oid
	}

	return nil
}

// GetSettlement 获取结算记录
func (r *AuthorRevenueRepositoryImpl) GetSettlement(ctx context.Context, settlementID primitive.ObjectID) (*financeModel.Settlement, error) {
	var settlement financeModel.Settlement
	err := r.settlementCollection.FindOne(ctx, bson.M{"_id": settlementID}).Decode(&settlement)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("结算记录不存在: %s", settlementID.Hex())
		}
		return nil, fmt.Errorf("查询结算记录失败: %w", err)
	}

	return &settlement, nil
}

// ListSettlements 列出结算记录
func (r *AuthorRevenueRepositoryImpl) ListSettlements(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]*financeModel.Settlement, int64, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize))

	cursor, err := r.settlementCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("查询结算记录列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var settlements []*financeModel.Settlement
	if err := cursor.All(ctx, &settlements); err != nil {
		return nil, 0, fmt.Errorf("解析结算记录列表失败: %w", err)
	}

	// 获取总数
	total, err := r.settlementCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("统计结算记录失败: %w", err)
	}

	return settlements, total, nil
}

// UpdateSettlement 更新结算记录
func (r *AuthorRevenueRepositoryImpl) UpdateSettlement(ctx context.Context, settlementID primitive.ObjectID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	result, err := r.settlementCollection.UpdateOne(ctx, bson.M{"_id": settlementID}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("更新结算记录失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("结算记录不存在: %s", settlementID.Hex())
	}

	return nil
}

// GetAuthorSettlements 获取作者结算记录
func (r *AuthorRevenueRepositoryImpl) GetAuthorSettlements(ctx context.Context, authorID string, page, pageSize int) ([]*financeModel.Settlement, int64, error) {
	filter := bson.M{"author_id": authorID}
	return r.ListSettlements(ctx, filter, page, pageSize)
}

// GetPendingSettlements 获取待结算记录
func (r *AuthorRevenueRepositoryImpl) GetPendingSettlements(ctx context.Context) ([]*financeModel.Settlement, error) {
	filter := bson.M{"status": financeModel.SettlementStatusPending}

	cursor, err := r.settlementCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("查询待结算记录失败: %w", err)
	}
	defer cursor.Close(ctx)

	var settlements []*financeModel.Settlement
	if err := cursor.All(ctx, &settlements); err != nil {
		return nil, fmt.Errorf("解析待结算记录失败: %w", err)
	}

	return settlements, nil
}

// ============ 收入统计 ============

// GetRevenueStatistics 获取收入统计
func (r *AuthorRevenueRepositoryImpl) GetRevenueStatistics(ctx context.Context, authorID string, period string, limit int) ([]*financeModel.RevenueStatistics, error) {
	filter := bson.M{"author_id": authorID}
	if period != "" {
		filter["period"] = period
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "period_start", Value: -1}})
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}

	cursor, err := r.statisticsCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询收入统计失败: %w", err)
	}
	defer cursor.Close(ctx)

	var statistics []*financeModel.RevenueStatistics
	if err := cursor.All(ctx, &statistics); err != nil {
		return nil, fmt.Errorf("解析收入统计失败: %w", err)
	}

	return statistics, nil
}

// CreateRevenueStatistics 创建收入统计
func (r *AuthorRevenueRepositoryImpl) CreateRevenueStatistics(ctx context.Context, statistics *financeModel.RevenueStatistics) error {
	now := time.Now()
	statistics.CreatedAt = now
	statistics.UpdatedAt = now

	_, err := r.statisticsCollection.InsertOne(ctx, statistics)
	if err != nil {
		return fmt.Errorf("创建收入统计失败: %w", err)
	}

	return nil
}

// UpdateRevenueStatistics 更新收入统计
func (r *AuthorRevenueRepositoryImpl) UpdateRevenueStatistics(ctx context.Context, statisticsID primitive.ObjectID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	result, err := r.statisticsCollection.UpdateOne(ctx, bson.M{"_id": statisticsID}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("更新收入统计失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("收入统计不存在: %s", statisticsID.Hex())
	}

	return nil
}

// ============ 收入明细 ============

// GetRevenueDetails 获取收入明细
func (r *AuthorRevenueRepositoryImpl) GetRevenueDetails(ctx context.Context, authorID string, page, pageSize int) ([]*financeModel.RevenueDetail, int64, error) {
	filter := bson.M{"author_id": authorID}

	opts := options.Find().
		SetSort(bson.D{{Key: "last_earning_at", Value: -1}}).
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize))

	cursor, err := r.detailCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("查询收入明细失败: %w", err)
	}
	defer cursor.Close(ctx)

	var details []*financeModel.RevenueDetail
	if err := cursor.All(ctx, &details); err != nil {
		return nil, 0, fmt.Errorf("解析收入明细失败: %w", err)
	}

	// 获取总数
	total, err := r.detailCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("统计收入明细失败: %w", err)
	}

	return details, total, nil
}

// GetRevenueDetailByBook 获取书籍收入明细
func (r *AuthorRevenueRepositoryImpl) GetRevenueDetailByBook(ctx context.Context, authorID string, bookID primitive.ObjectID) (*financeModel.RevenueDetail, error) {
	var detail financeModel.RevenueDetail
	err := r.detailCollection.FindOne(ctx, bson.M{"author_id": authorID, "book_id": bookID}).Decode(&detail)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("收入明细不存在: author=%s, book=%s", authorID, bookID.Hex())
		}
		return nil, fmt.Errorf("查询收入明细失败: %w", err)
	}

	return &detail, nil
}

// CreateRevenueDetail 创建收入明细
func (r *AuthorRevenueRepositoryImpl) CreateRevenueDetail(ctx context.Context, detail *financeModel.RevenueDetail) error {
	now := time.Now()
	detail.CreatedAt = now
	detail.UpdatedAt = now

	result, err := r.detailCollection.InsertOne(ctx, detail)
	if err != nil {
		return fmt.Errorf("创建收入明细失败: %w", err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		detail.ID = oid
	}

	return nil
}

// UpdateRevenueDetail 更新收入明细
func (r *AuthorRevenueRepositoryImpl) UpdateRevenueDetail(ctx context.Context, detailID primitive.ObjectID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	result, err := r.detailCollection.UpdateOne(ctx, bson.M{"_id": detailID}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("更新收入明细失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("收入明细不存在: %s", detailID.Hex())
	}

	return nil
}

// ============ 税务信息 ============

// CreateTaxInfo 创建税务信息
func (r *AuthorRevenueRepositoryImpl) CreateTaxInfo(ctx context.Context, taxInfo *financeModel.TaxInfo) error {
	now := time.Now()
	taxInfo.CreatedAt = now
	taxInfo.UpdatedAt = now

	result, err := r.taxInfoCollection.InsertOne(ctx, taxInfo)
	if err != nil {
		return fmt.Errorf("创建税务信息失败: %w", err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		taxInfo.ID = oid
	}

	return nil
}

// GetTaxInfo 获取税务信息
func (r *AuthorRevenueRepositoryImpl) GetTaxInfo(ctx context.Context, userID string) (*financeModel.TaxInfo, error) {
	var taxInfo financeModel.TaxInfo
	err := r.taxInfoCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&taxInfo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("税务信息不存在: user %s", userID)
		}
		return nil, fmt.Errorf("查询税务信息失败: %w", err)
	}

	return &taxInfo, nil
}

// UpdateTaxInfo 更新税务信息
func (r *AuthorRevenueRepositoryImpl) UpdateTaxInfo(ctx context.Context, userID string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	result, err := r.taxInfoCollection.UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("更新税务信息失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("税务信息不存在: user %s", userID)
	}

	return nil
}
