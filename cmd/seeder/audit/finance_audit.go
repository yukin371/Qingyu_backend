package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuditResult struct {
	Name  string
	Count int64
	Details []interface{}
}

type FinanceAuditReport struct {
	Timestamp time.Time
	Results   map[string]*AuditResult
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("连接MongoDB失败:", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")

	report := &FinanceAuditReport{
		Timestamp: time.Now(),
		Results:   make(map[string]*AuditResult),
	}

	fmt.Println("=== 开始财务视角数据关联审查 ===")
	fmt.Printf("审查时间: %s\n\n", report.Timestamp.Format("2006-01-02 15:04:05"))

	// 1. 钱包外键验证 - 验证 wallets.user_id → users._id
	fmt.Println("1. 验证 wallets.user_id → users._id")
	orphanedWallets := countOrphanedWallets(ctx, db)
	report.Results["orphaned_wallets"] = &AuditResult{
		Name:  "孤儿钱包 (wallets.user_id 不存在)",
		Count: orphanedWallets,
	}
	fmt.Printf("   结果: 发现 %d 个孤儿钱包\n\n", orphanedWallets)

	// 2. 验证每个用户都有钱包
	fmt.Println("2. 验证每个用户都有钱包")
	usersWithoutWallets := countUsersWithoutWallets(ctx, db)
	report.Results["users_without_wallets"] = &AuditResult{
		Name:  "无钱包用户",
		Count: usersWithoutWallets,
	}
	fmt.Printf("   结果: 发现 %d 个无钱包用户\n\n", usersWithoutWallets)

	// 3. 交易记录外键验证 - 验证 transactions.user_id → users._id
	fmt.Println("3. 验证 transactions.user_id → users._id")
	orphanedTransactions := countOrphanedTransactions(ctx, db)
	report.Results["orphaned_transactions"] = &AuditResult{
		Name:  "孤儿交易记录 (transactions.user_id 不存在)",
		Count: orphanedTransactions,
	}
	fmt.Printf("   结果: 发现 %d 条孤儿交易记录\n\n", orphanedTransactions)

	// 4. 验证交易类型有效性
	fmt.Println("4. 验证交易类型有效性")
	invalidTransactionTypes := checkTransactionTypes(ctx, db)
	report.Results["invalid_transaction_types"] = &AuditResult{
		Name:  "无效交易类型",
		Count: int64(len(invalidTransactionTypes)),
	}
	fmt.Printf("   结果: 发现 %d 种无效交易类型: %v\n\n", len(invalidTransactionTypes), invalidTransactionTypes)

	// 5. 会员订阅外键验证
	fmt.Println("5. 验证 memberships.user_id → users._id")
	orphanedMemberships := countOrphanedMemberships(ctx, db)
	report.Results["orphaned_memberships"] = &AuditResult{
		Name:  "孤儿会员记录 (memberships.user_id 不存在)",
		Count: orphanedMemberships,
	}
	fmt.Printf("   结果: 发现 %d 条孤儿会员记录\n\n", orphanedMemberships)

	// 6. 作者收益外键验证
	fmt.Println("6. 验证 author_revenue.user_id → users._id")
	orphanedRevenue := countOrphanedRevenue(ctx, db)
	report.Results["orphaned_revenue"] = &AuditResult{
		Name:  "孤儿收益记录 (author_revenue.user_id 不存在)",
		Count: orphanedRevenue,
	}
	fmt.Printf("   结果: 发现 %d 条孤儿收益记录\n\n", orphanedRevenue)

	// 7. 验证非作者用户的收益记录
	fmt.Println("7. 验证非作者用户的收益记录")
	revenueForNonAuthors := countRevenueForNonAuthors(ctx, db)
	report.Results["revenue_for_non_authors"] = &AuditResult{
		Name:  "非作者用户的收益记录",
		Count: revenueForNonAuthors,
	}
	fmt.Printf("   结果: 发现 %d 条非作者用户的收益记录\n\n", revenueForNonAuthors)

	// 8. 钱包余额非负验证
	fmt.Println("8. 验证钱包余额非负")
	negativeBalanceWallets := countNegativeBalanceWallets(ctx, db)
	report.Results["negative_balance_wallets"] = &AuditResult{
		Name:  "负余额钱包",
		Count: negativeBalanceWallets,
	}
	fmt.Printf("   结果: 发现 %d 个负余额钱包\n\n", negativeBalanceWallets)

	// 9. 交易金额合理性验证
	fmt.Println("9. 验证交易金额合理性 (amount > 0)")
	invalidAmountTransactions := countInvalidAmountTransactions(ctx, db)
	report.Results["invalid_amount_transactions"] = &AuditResult{
		Name:  "金额不合理的交易 (amount <= 0)",
		Count: invalidAmountTransactions,
	}
	fmt.Printf("   结果: 发现 %d 条金额不合理的交易\n\n", invalidAmountTransactions)

	// 10. 会员状态一致性验证
	fmt.Println("10. 验证会员状态一致性")
	inconsistentMemberships := countInconsistentMemberships(ctx, db)
	report.Results["inconsistent_memberships"] = &AuditResult{
		Name:  "状态不一致的会员记录",
		Count: inconsistentMemberships,
	}
	fmt.Printf("   结果: 发现 %d 条状态不一致的会员记录\n\n", inconsistentMemberships)

	// 打印汇总报告
	fmt.Println("=== 财务视角数据关联审查汇总 ===")
	fmt.Println()

	totalIssues := int64(0)
	for _, result := range report.Results {
		fmt.Printf("%s: %d\n", result.Name, result.Count)
		totalIssues += result.Count
	}
	fmt.Printf("\n总计问题数: %d\n", totalIssues)

	// 数据质量评估
	fmt.Println("\n=== 数据质量评估 ===")
	qualityScore := calculateQualityScore(report)
	fmt.Printf("整体评分: %.1f/100 (%s)\n", qualityScore, getQualityLevel(qualityScore))

	fmt.Println("\n主要问题汇总:")
	printTopIssues(report)

	fmt.Println("\n修复建议:")
	printRecommendations(report)
}

// 辅助函数
func countOrphanedWallets(ctx context.Context, db *mongo.Database) int64 {
	collection := db.Collection("wallets")

	pipeline := mongo.Pipeline{
		bson.D{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "user_id"},
			{"foreignField", "_id"},
			{"as", "user"},
		}}},
		bson.D{{"$match", bson.D{
			{"user", bson.D{{"$size", 0}}},
		}}},
		bson.D{{"$count", "count"}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("查询孤儿钱包失败: %v", err)
		return 0
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return 0
	}

	if len(results) > 0 {
		if count, ok := results[0]["count"].(int32); ok {
			return int64(count)
		}
	}
	return 0
}

func countUsersWithoutWallets(ctx context.Context, db *mongo.Database) int64 {
	pipeline := mongo.Pipeline{
		bson.D{{"$lookup", bson.D{
			{"from", "wallets"},
			{"localField", "_id"},
			{"foreignField", "user_id"},
			{"as", "wallet"},
		}}},
		bson.D{{"$match", bson.D{
			{"wallet", bson.D{{"$size", 0}}},
		}}},
		bson.D{{"$count", "count"}},
	}

	cursor, err := db.Collection("users").Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("查询无钱包用户失败: %v", err)
		return 0
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return 0
	}

	if len(results) > 0 {
		if count, ok := results[0]["count"].(int32); ok {
			return int64(count)
		}
	}
	return 0
}

func countOrphanedTransactions(ctx context.Context, db *mongo.Database) int64 {
	pipeline := mongo.Pipeline{
		bson.D{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "user_id"},
			{"foreignField", "_id"},
			{"as", "user"},
		}}},
		bson.D{{"$match", bson.D{
			{"user", bson.D{{"$size", 0}}},
		}}},
		bson.D{{"$count", "count"}},
	}

	cursor, err := db.Collection("transactions").Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("查询孤儿交易失败: %v", err)
		return 0
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return 0
	}

	if len(results) > 0 {
		if count, ok := results[0]["count"].(int32); ok {
			return int64(count)
		}
	}
	return 0
}

func checkTransactionTypes(ctx context.Context, db *mongo.Database) []string {
	validTypes := map[string]bool{
		"recharge": true,
		"consume":  true,
		"withdraw": true,
		"reward":    true,
		"refund":    true,
	}

	collection := db.Collection("transactions")
	distinctTypes, err := collection.Distinct(ctx, "type", bson.D{})
	if err != nil {
		log.Printf("获取交易类型失败: %v", err)
		return nil
	}

	var invalidTypes []string
	for _, t := range distinctTypes {
		if typeStr, ok := t.(string); ok {
			if !validTypes[typeStr] {
				invalidTypes = append(invalidTypes, typeStr)
			}
		}
	}
	return invalidTypes
}

func countOrphanedMemberships(ctx context.Context, db *mongo.Database) int64 {
	pipeline := mongo.Pipeline{
		bson.D{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "user_id"},
			{"foreignField", "_id"},
			{"as", "user"},
		}}},
		bson.D{{"$match", bson.D{
			{"user", bson.D{{"$size", 0}}},
		}}},
		bson.D{{"$count", "count"}},
	}

	cursor, err := db.Collection("memberships").Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("查询孤儿会员失败: %v", err)
		return 0
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return 0
	}

	if len(results) > 0 {
		if count, ok := results[0]["count"].(int32); ok {
			return int64(count)
		}
	}
	return 0
}

func countOrphanedRevenue(ctx context.Context, db *mongo.Database) int64 {
	pipeline := mongo.Pipeline{
		bson.D{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "user_id"},
			{"foreignField", "_id"},
			{"as", "user"},
		}}},
		bson.D{{"$match", bson.D{
			{"user", bson.D{{"$size", 0}}},
		}}},
		bson.D{{"$count", "count"}},
	}

	cursor, err := db.Collection("author_revenue").Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("查询孤儿收益失败: %v", err)
		return 0
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return 0
	}

	if len(results) > 0 {
		if count, ok := results[0]["count"].(int32); ok {
			return int64(count)
		}
	}
	return 0
}

func countRevenueForNonAuthors(ctx context.Context, db *mongo.Database) int64 {
	pipeline := mongo.Pipeline{
		bson.D{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "user_id"},
			{"foreignField", "_id"},
			{"as", "user"},
		}}},
		bson.D{{"$match", bson.D{
			{"user.role", bson.D{{"$ne", "author"}}},
		}}},
		bson.D{{"$count", "count"}},
	}

	cursor, err := db.Collection("author_revenue").Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("查询非作者收益失败: %v", err)
		return 0
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return 0
	}

	if len(results) > 0 {
		if count, ok := results[0]["count"].(int32); ok {
			return int64(count)
		}
	}
	return 0
}

func countNegativeBalanceWallets(ctx context.Context, db *mongo.Database) int64 {
	count, err := db.Collection("wallets").CountDocuments(ctx, bson.D{
		{"balance", bson.D{{"$lt", 0}}},
	})
	if err != nil {
		log.Printf("查询负余额钱包失败: %v", err)
		return 0
	}
	return count
}

func countInvalidAmountTransactions(ctx context.Context, db *mongo.Database) int64 {
	count, err := db.Collection("transactions").CountDocuments(ctx, bson.D{
		{"amount", bson.D{{"$lte", 0}}},
	})
	if err != nil {
		log.Printf("查询无效金额交易失败: %v", err)
		return 0
	}
	return count
}

func countInconsistentMemberships(ctx context.Context, db *mongo.Database) int64 {
	now := primitive.NewDateTimeFromTime(time.Now())

	pipeline := mongo.Pipeline{
		bson.D{{"$project", bson.D{
			{"_id", 1},
			{"status", 1},
			{"expire_at", 1},
			{"is_expired", bson.D{{"$lt", bson.A{"$expire_at", now}}}},
		}}},
		bson.D{{"$match", bson.D{
			{"$expr", bson.D{
				{"$ne", bson.A{
					"$status",
					bson.D{{"$cond", bson.A{
						bson.D{{"$lt", bson.A{"$expire_at", now}}},
						"expired",
						"active",
					}}},
				}},
			}},
		}}},
		bson.D{{"$count", "count"}},
	}

	cursor, err := db.Collection("memberships").Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("查询状态不一致会员失败: %v", err)
		return 0
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return 0
	}

	if len(results) > 0 {
		if count, ok := results[0]["count"].(int32); ok {
			return int64(count)
		}
	}
	return 0
}

func calculateQualityScore(report *FinanceAuditReport) float64 {
	totalChecks := float64(len(report.Results))
	if totalChecks == 0 {
		return 100
	}

	issueCount := int64(0)
	for _, result := range report.Results {
		issueCount += result.Count
	}

	// 基础分数100，每个问题扣1分，最低0分
	score := 100.0 - float64(issueCount)
	if score < 0 {
		score = 0
	}
	return score
}

func getQualityLevel(score float64) string {
	if score >= 90 {
		return "优秀"
	} else if score >= 75 {
		return "良好"
	} else if score >= 60 {
		return "一般"
	}
	return "差"
}

func printTopIssues(report *FinanceAuditReport) {
	type issue struct {
		name  string
		count int64
	}

	var issues []issue
	for _, result := range report.Results {
		if result.Count > 0 {
			issues = append(issues, issue{
				name:  result.Name,
				count: result.Count,
			})
		}
	}

	if len(issues) == 0 {
		fmt.Println("  未发现问题")
		return
	}

	for i, iss := range issues {
		fmt.Printf("  %d. %s: %d\n", i+1, iss.name, iss.count)
	}
}

func printRecommendations(report *FinanceAuditReport) {
	// 根据问题类型给出具体建议
	hasIssues := false

	if report.Results["orphaned_wallets"] != nil && report.Results["orphaned_wallets"].Count > 0 {
		fmt.Println("  1. 清理孤儿钱包记录，或为这些钱包创建对应用户")
		hasIssues = true
	}

	if report.Results["users_without_wallets"] != nil && report.Results["users_without_wallets"].Count > 0 {
		fmt.Println("  2. 为所有缺失钱包的用户创建钱包")
		hasIssues = true
	}

	if report.Results["orphaned_transactions"] != nil && report.Results["orphaned_transactions"].Count > 0 {
		fmt.Println("  3. 清理孤儿交易记录，或修正 user_id")
		hasIssues = true
	}

	if report.Results["invalid_transaction_types"] != nil && report.Results["invalid_transaction_types"].Count > 0 {
		fmt.Println("  4. 规范化交易类型，只允许: recharge, consume, withdraw, reward, refund")
		hasIssues = true
	}

	if report.Results["orphaned_memberships"] != nil && report.Results["orphaned_memberships"].Count > 0 {
		fmt.Println("  5. 清理孤儿会员记录")
		hasIssues = true
	}

	if report.Results["orphaned_revenue"] != nil && report.Results["orphaned_revenue"].Count > 0 {
		fmt.Println("  6. 清理孤儿收益记录")
		hasIssues = true
	}

	if report.Results["revenue_for_non_authors"] != nil && report.Results["revenue_for_non_authors"].Count > 0 {
		fmt.Println("  7. 移除非作者用户的收益记录，或更正用户角色")
		hasIssues = true
	}

	if report.Results["negative_balance_wallets"] != nil && report.Results["negative_balance_wallets"].Count > 0 {
		fmt.Println("  8. 修正负余额钱包，确保余额不小于0")
		hasIssues = true
	}

	if report.Results["invalid_amount_transactions"] != nil && report.Results["invalid_amount_transactions"].Count > 0 {
		fmt.Println("  9. 修正金额不合理的交易，确保交易金额大于0")
		hasIssues = true
	}

	if report.Results["inconsistent_memberships"] != nil && report.Results["inconsistent_memberships"].Count > 0 {
		fmt.Println("  10. 同步会员状态与过期时间")
		hasIssues = true
	}

	if !hasIssues {
		fmt.Println("  数据质量良好，建议继续使用 Seeder 工具的 ID 类型规范")
	}

	fmt.Println("\n  预防措施:")
	fmt.Println("  - 在 Seeder 工具中统一使用 ObjectID 生成 ID")
	fmt.Println("  - 添加数据生成后的事务性验证")
	fmt.Println("  - 实施数据库约束和触发器")
}
