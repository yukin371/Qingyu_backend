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

// MembershipRepositoryImpl 会员Repository实现
type MembershipRepositoryImpl struct {
	db                      *mongo.Database
	planCollection          *mongo.Collection
	membershipCollection    *mongo.Collection
	cardCollection          *mongo.Collection
	benefitCollection       *mongo.Collection
	usageCollection         *mongo.Collection
}

// NewMembershipRepository 创建会员Repository
func NewMembershipRepository(db *mongo.Database) financeInterfaces.MembershipRepository {
	return &MembershipRepositoryImpl{
		db:                   db,
		planCollection:       db.Collection("membership_plans"),
		membershipCollection: db.Collection("user_memberships"),
		cardCollection:       db.Collection("membership_cards"),
		benefitCollection:    db.Collection("membership_benefits"),
		usageCollection:      db.Collection("membership_usage"),
	}
}

// ============ 套餐管理 ============

// CreatePlan 创建套餐
func (r *MembershipRepositoryImpl) CreatePlan(ctx context.Context, plan *financeModel.MembershipPlan) error {
	now := time.Now()
	plan.CreatedAt = now
	plan.UpdatedAt = now

	result, err := r.planCollection.InsertOne(ctx, plan)
	if err != nil {
		return fmt.Errorf("创建套餐失败: %w", err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		plan.ID = oid
	}

	return nil
}

// GetPlan 获取套餐
func (r *MembershipRepositoryImpl) GetPlan(ctx context.Context, planID primitive.ObjectID) (*financeModel.MembershipPlan, error) {
	var plan financeModel.MembershipPlan
	err := r.planCollection.FindOne(ctx, bson.M{"_id": planID}).Decode(&plan)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("套餐不存在: %s", planID.Hex())
		}
		return nil, fmt.Errorf("查询套餐失败: %w", err)
	}

	return &plan, nil
}

// GetPlanByType 根据类型获取套餐
func (r *MembershipRepositoryImpl) GetPlanByType(ctx context.Context, planType string) (*financeModel.MembershipPlan, error) {
	var plan financeModel.MembershipPlan
	err := r.planCollection.FindOne(ctx, bson.M{"type": planType}).Decode(&plan)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("套餐不存在: type=%s", planType)
		}
		return nil, fmt.Errorf("查询套餐失败: %w", err)
	}

	return &plan, nil
}

// ListPlans 列出套餐
func (r *MembershipRepositoryImpl) ListPlans(ctx context.Context, enabledOnly bool) ([]*financeModel.MembershipPlan, error) {
	filter := bson.M{}
	if enabledOnly {
		filter["is_enabled"] = true
	}

	opts := options.Find().SetSort(bson.D{{Key: "sort_order", Value: 1}, {Key: "created_at", Value: -1}})

	cursor, err := r.planCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询套餐列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var plans []*financeModel.MembershipPlan
	if err := cursor.All(ctx, &plans); err != nil {
		return nil, fmt.Errorf("解析套餐列表失败: %w", err)
	}

	return plans, nil
}

// UpdatePlan 更新套餐
func (r *MembershipRepositoryImpl) UpdatePlan(ctx context.Context, planID primitive.ObjectID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	result, err := r.planCollection.UpdateOne(ctx, bson.M{"_id": planID}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("更新套餐失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("套餐不存在: %s", planID.Hex())
	}

	return nil
}

// DeletePlan 删除套餐
func (r *MembershipRepositoryImpl) DeletePlan(ctx context.Context, planID primitive.ObjectID) error {
	result, err := r.planCollection.DeleteOne(ctx, bson.M{"_id": planID})
	if err != nil {
		return fmt.Errorf("删除套餐失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("套餐不存在: %s", planID.Hex())
	}

	return nil
}

// ============ 用户会员管理 ============

// CreateMembership 创建会员
func (r *MembershipRepositoryImpl) CreateMembership(ctx context.Context, membership *financeModel.UserMembership) error {
	now := time.Now()
	membership.CreatedAt = now
	membership.UpdatedAt = now

	result, err := r.membershipCollection.InsertOne(ctx, membership)
	if err != nil {
		return fmt.Errorf("创建会员失败: %w", err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		membership.ID = oid
	}

	return nil
}

// GetMembership 获取用户会员信息
func (r *MembershipRepositoryImpl) GetMembership(ctx context.Context, userID string) (*financeModel.UserMembership, error) {
	var membership financeModel.UserMembership
	err := r.membershipCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&membership)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("会员不存在: user %s", userID)
		}
		return nil, fmt.Errorf("查询会员失败: %w", err)
	}

	return &membership, nil
}

// GetMembershipByID 根据ID获取会员
func (r *MembershipRepositoryImpl) GetMembershipByID(ctx context.Context, membershipID primitive.ObjectID) (*financeModel.UserMembership, error) {
	var membership financeModel.UserMembership
	err := r.membershipCollection.FindOne(ctx, bson.M{"_id": membershipID}).Decode(&membership)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("会员不存在: %s", membershipID.Hex())
		}
		return nil, fmt.Errorf("查询会员失败: %w", err)
	}

	return &membership, nil
}

// UpdateMembership 更新会员
func (r *MembershipRepositoryImpl) UpdateMembership(ctx context.Context, membershipID primitive.ObjectID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	result, err := r.membershipCollection.UpdateOne(ctx, bson.M{"_id": membershipID}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("更新会员失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("会员不存在: %s", membershipID.Hex())
	}

	return nil
}

// DeleteMembership 删除会员
func (r *MembershipRepositoryImpl) DeleteMembership(ctx context.Context, membershipID primitive.ObjectID) error {
	result, err := r.membershipCollection.DeleteOne(ctx, bson.M{"_id": membershipID})
	if err != nil {
		return fmt.Errorf("删除会员失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("会员不存在: %s", membershipID.Hex())
	}

	return nil
}

// ListMemberships 列出会员
func (r *MembershipRepositoryImpl) ListMemberships(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]*financeModel.UserMembership, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize))

	cursor, err := r.membershipCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询会员列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var memberships []*financeModel.UserMembership
	if err := cursor.All(ctx, &memberships); err != nil {
		return nil, fmt.Errorf("解析会员列表失败: %w", err)
	}

	return memberships, nil
}

// ============ 会员卡管理 ============

// CreateMembershipCard 创建会员卡
func (r *MembershipRepositoryImpl) CreateMembershipCard(ctx context.Context, card *financeModel.MembershipCard) error {
	now := time.Now()
	card.CreatedAt = now
	card.UpdatedAt = now

	result, err := r.cardCollection.InsertOne(ctx, card)
	if err != nil {
		return fmt.Errorf("创建会员卡失败: %w", err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		card.ID = oid
	}

	return nil
}

// GetMembershipCard 获取会员卡
func (r *MembershipRepositoryImpl) GetMembershipCard(ctx context.Context, cardID primitive.ObjectID) (*financeModel.MembershipCard, error) {
	var card financeModel.MembershipCard
	err := r.cardCollection.FindOne(ctx, bson.M{"_id": cardID}).Decode(&card)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("会员卡不存在: %s", cardID.Hex())
		}
		return nil, fmt.Errorf("查询会员卡失败: %w", err)
	}

	return &card, nil
}

// GetMembershipCardByCode 根据卡密获取会员卡
func (r *MembershipRepositoryImpl) GetMembershipCardByCode(ctx context.Context, code string) (*financeModel.MembershipCard, error) {
	var card financeModel.MembershipCard
	err := r.cardCollection.FindOne(ctx, bson.M{"code": code}).Decode(&card)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("会员卡不存在: %s", code)
		}
		return nil, fmt.Errorf("查询会员卡失败: %w", err)
	}

	return &card, nil
}

// UpdateMembershipCard 更新会员卡
func (r *MembershipRepositoryImpl) UpdateMembershipCard(ctx context.Context, cardID primitive.ObjectID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	result, err := r.cardCollection.UpdateOne(ctx, bson.M{"_id": cardID}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("更新会员卡失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("会员卡不存在: %s", cardID.Hex())
	}

	return nil
}

// ListMembershipCards 列出会员卡
func (r *MembershipRepositoryImpl) ListMembershipCards(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]*financeModel.MembershipCard, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize))

	cursor, err := r.cardCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询会员卡列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var cards []*financeModel.MembershipCard
	if err := cursor.All(ctx, &cards); err != nil {
		return nil, fmt.Errorf("解析会员卡列表失败: %w", err)
	}

	return cards, nil
}

// CountMembershipCards 统计会员卡
func (r *MembershipRepositoryImpl) CountMembershipCards(ctx context.Context, filter map[string]interface{}) (int64, error) {
	count, err := r.cardCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("统计会员卡失败: %w", err)
	}

	return count, nil
}

// BatchCreateMembershipCards 批量创建会员卡
func (r *MembershipRepositoryImpl) BatchCreateMembershipCards(ctx context.Context, cards []*financeModel.MembershipCard) error {
	now := time.Now()

	documents := make([]interface{}, len(cards))
	for i, card := range cards {
		card.CreatedAt = now
		card.UpdatedAt = now
		documents[i] = card
	}

	_, err := r.cardCollection.InsertMany(ctx, documents)
	if err != nil {
		return fmt.Errorf("批量创建会员卡失败: %w", err)
	}

	return nil
}

// ============ 会员权益管理 ============

// CreateBenefit 创建权益
func (r *MembershipRepositoryImpl) CreateBenefit(ctx context.Context, benefit *financeModel.MembershipBenefit) error {
	now := time.Now()
	benefit.CreatedAt = now
	benefit.UpdatedAt = now

	result, err := r.benefitCollection.InsertOne(ctx, benefit)
	if err != nil {
		return fmt.Errorf("创建权益失败: %w", err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		benefit.ID = oid
	}

	return nil
}

// GetBenefit 获取权益
func (r *MembershipRepositoryImpl) GetBenefit(ctx context.Context, benefitID primitive.ObjectID) (*financeModel.MembershipBenefit, error) {
	var benefit financeModel.MembershipBenefit
	err := r.benefitCollection.FindOne(ctx, bson.M{"_id": benefitID}).Decode(&benefit)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("权益不存在: %s", benefitID.Hex())
		}
		return nil, fmt.Errorf("查询权益失败: %w", err)
	}

	return &benefit, nil
}

// GetBenefitByCode 根据代码获取权益
func (r *MembershipRepositoryImpl) GetBenefitByCode(ctx context.Context, code string) (*financeModel.MembershipBenefit, error) {
	var benefit financeModel.MembershipBenefit
	err := r.benefitCollection.FindOne(ctx, bson.M{"code": code}).Decode(&benefit)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("权益不存在: %s", code)
		}
		return nil, fmt.Errorf("查询权益失败: %w", err)
	}

	return &benefit, nil
}

// ListBenefits 列出权益
func (r *MembershipRepositoryImpl) ListBenefits(ctx context.Context, level string, enabledOnly bool) ([]*financeModel.MembershipBenefit, error) {
	filter := bson.M{}
	if level != "" {
		filter["level"] = level
	}
	if enabledOnly {
		filter["is_enabled"] = true
	}

	cursor, err := r.benefitCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("查询权益列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var benefits []*financeModel.MembershipBenefit
	if err := cursor.All(ctx, &benefits); err != nil {
		return nil, fmt.Errorf("解析权益列表失败: %w", err)
	}

	return benefits, nil
}

// UpdateBenefit 更新权益
func (r *MembershipRepositoryImpl) UpdateBenefit(ctx context.Context, benefitID primitive.ObjectID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	result, err := r.benefitCollection.UpdateOne(ctx, bson.M{"_id": benefitID}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("更新权益失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("权益不存在: %s", benefitID.Hex())
	}

	return nil
}

// DeleteBenefit 删除权益
func (r *MembershipRepositoryImpl) DeleteBenefit(ctx context.Context, benefitID primitive.ObjectID) error {
	result, err := r.benefitCollection.DeleteOne(ctx, bson.M{"_id": benefitID})
	if err != nil {
		return fmt.Errorf("删除权益失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("权益不存在: %s", benefitID.Hex())
	}

	return nil
}

// ============ 会员权益使用情况 ============

// CreateUsage 创建使用记录
func (r *MembershipRepositoryImpl) CreateUsage(ctx context.Context, usage *financeModel.MembershipUsage) error {
	now := time.Now()
	usage.CreatedAt = now
	usage.UpdatedAt = now

	result, err := r.usageCollection.InsertOne(ctx, usage)
	if err != nil {
		return fmt.Errorf("创建使用记录失败: %w", err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		usage.ID = oid
	}

	return nil
}

// GetUsage 获取使用记录
func (r *MembershipRepositoryImpl) GetUsage(ctx context.Context, userID string, benefitCode string) (*financeModel.MembershipUsage, error) {
	var usage financeModel.MembershipUsage
	err := r.usageCollection.FindOne(ctx, bson.M{"user_id": userID, "benefit_code": benefitCode}).Decode(&usage)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("使用记录不存在: user=%s, benefit=%s", userID, benefitCode)
		}
		return nil, fmt.Errorf("查询使用记录失败: %w", err)
	}

	return &usage, nil
}

// UpdateUsage 更新使用记录
func (r *MembershipRepositoryImpl) UpdateUsage(ctx context.Context, usageID primitive.ObjectID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	result, err := r.usageCollection.UpdateOne(ctx, bson.M{"_id": usageID}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("更新使用记录失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("使用记录不存在: %s", usageID.Hex())
	}

	return nil
}

// ListUsages 列出使用记录
func (r *MembershipRepositoryImpl) ListUsages(ctx context.Context, userID string) ([]*financeModel.MembershipUsage, error) {
	cursor, err := r.usageCollection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("查询使用记录列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var usages []*financeModel.MembershipUsage
	if err := cursor.All(ctx, &usages); err != nil {
		return nil, fmt.Errorf("解析使用记录列表失败: %w", err)
	}

	return usages, nil
}
