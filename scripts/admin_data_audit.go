//go:build ignore

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
	Value interface{}
}

type Report struct {
	Section string
	Results []AuditResult
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/qingyu"))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")

	fmt.Println("========================================")
	fmt.Println("管理员视角数据关联审查报告")
	fmt.Println("========================================")
	fmt.Printf("审查时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println()

	// 外键关系完整性
	fmt.Println("## 外键关系完整性")
	fmt.Println()

	// 1. 用户角色一致性
	fmt.Println("### 1. 用户角色一致性")
	roles := distinctRoles(ctx, db)
	fmt.Printf("- 角色类型: %v\n", roles)
	fmt.Printf("- 预期角色: [reader, author, admin, vip]\n")
	fmt.Println()

	// 2. 作者是否有作品
	authorsWithoutBooks := countAuthorsWithoutBooks(ctx, db)
	fmt.Printf("- 无作品的作者: %d 个\n", authorsWithoutBooks)
	fmt.Println()

	// 3. 分类有效性
	invalidCategories := countInvalidCategories(ctx, db)
	fmt.Printf("- 无效分类的书籍: %d 本\n", invalidCategories)
	fmt.Println()

	// 4. 空分类
	emptyCategories := countEmptyCategories(ctx, db)
	fmt.Printf("- 空分类: %d 个\n", emptyCategories)
	fmt.Println()

	// 5. 榜单中的无效书籍
	invalidBooksInRankings := countInvalidBooksInRankings(ctx, db)
	fmt.Printf("- 榜单中无效书籍: %d 本\n", invalidBooksInRankings)
	fmt.Println()

	// 业务规则一致性
	fmt.Println("## 业务规则一致性")
	fmt.Println()

	// 1. 过期公告
	expiredAnnouncements := countExpiredAnnouncements(ctx, db)
	fmt.Printf("- 过期公告: %d 条\n", expiredAnnouncements)
	fmt.Println()

	// 2. 未发布公告
	unpublishedAnnouncements := countUnpublishedAnnouncements(ctx, db)
	fmt.Printf("- 未发布公告: %d 条\n", unpublishedAnnouncements)
	fmt.Println()

	// 3. 封禁用户异常活跃
	bannedUsersActive := countBannedUsersActive(ctx, db)
	fmt.Printf("- 封禁用户异常活跃: %d 个\n", bannedUsersActive)
	fmt.Println()

	// 4. 过期VIP用户
	expiredVIPUsers := countExpiredVIPUsers(ctx, db)
	fmt.Printf("- 过期VIP用户: %d 个\n", expiredVIPUsers)
	fmt.Println()

	// 5. 用户状态统计
	userStatusStats := getUserStatusStats(ctx, db)
	fmt.Println("### 用户状态统计:")
	if len(userStatusStats) > 0 {
		for status, count := range userStatusStats {
			fmt.Printf("- %s: %d 个\n", status, count)
		}
	} else {
		fmt.Println("- 无法获取用户状态统计")
	}
	fmt.Println()

	// 6. 详细问题分析
	fmt.Println("## 详细问题分析")
	fmt.Println()

	// 分析无效分类的书籍详情
	invalidCategoriesDetail := getInvalidCategoriesDetail(ctx, db)
	fmt.Printf("### 无效分类书籍详情 (%d 本)\n", len(invalidCategoriesDetail))
	for _, book := range invalidCategoriesDetail[:min(10, len(invalidCategoriesDetail))] {
		fmt.Printf("- 书籍ID: %s, 分类: '%s'\n", book.ID.Hex(), book.Category)
	}
	if len(invalidCategoriesDetail) > 10 {
		fmt.Printf("- ... 还有 %d 本书籍\n", len(invalidCategoriesDetail)-10)
	}
	fmt.Println()

	// 分析空分类详情
	emptyCategoriesDetail := getEmptyCategoriesDetail(ctx, db)
	fmt.Printf("### 空分类详情 (%d 个)\n", len(emptyCategoriesDetail))
	for _, cat := range emptyCategoriesDetail {
		fmt.Printf("- 分类ID: %s, 名称: %s\n", cat.ID.Hex(), cat.Name)
	}
	fmt.Println()

	// 分析未发布公告详情
	unpublishedAnnouncementsDetail := getUnpublishedAnnouncementsDetail(ctx, db)
	fmt.Printf("### 未发布公告详情 (%d 条)\n", len(unpublishedAnnouncementsDetail))
	for _, ann := range unpublishedAnnouncementsDetail {
		fmt.Printf("- 公告ID: %s, 状态: %s, 标题: %s\n", ann.ID.Hex(), ann.Status, ann.Title)
	}
	fmt.Println()

	// 数据质量评估
	fmt.Println("## 数据质量评估")
	fmt.Println()

	totalIssues := authorsWithoutBooks + invalidCategories + emptyCategories +
		invalidBooksInRankings + expiredAnnouncements + unpublishedAnnouncements +
		bannedUsersActive + expiredVIPUsers

	fmt.Printf("- 总问题数: %d\n", totalIssues)
	fmt.Println()

	var rating string
	if totalIssues == 0 {
		rating = "优秀"
	} else if totalIssues <= 10 {
		rating = "良好"
	} else if totalIssues <= 50 {
		rating = "一般"
	} else {
		rating = "差"
	}
	fmt.Printf("- 整体评分: %s\n", rating)
	fmt.Println()

	fmt.Println("========================================")
}

func distinctRoles(ctx context.Context, db *mongo.Database) []string {
	collection := db.Collection("users")
	distinct, err := collection.Distinct(ctx, "role", bson.M{})
	if err != nil {
		log.Printf("Error getting distinct roles: %v", err)
		return []string{"error"}
	}
	roles := make([]string, 0, len(distinct))
	for _, v := range distinct {
		if s, ok := v.(string); ok {
			roles = append(roles, s)
		}
	}
	return roles
}

func countAuthorsWithoutBooks(ctx context.Context, db *mongo.Database) int64 {
	usersCollection := db.Collection("users")
	booksCollection := db.Collection("books")

	// 获取所有作者
	cursor, err := usersCollection.Find(ctx, bson.M{"role": "author"})
	if err != nil {
		log.Printf("Error finding authors: %v", err)
		return 0
	}
	defer cursor.Close(ctx)

	type User struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	var authors []User
	if err = cursor.All(ctx, &authors); err != nil {
		log.Printf("Error decoding authors: %v", err)
		return 0
	}

	if len(authors) == 0 {
		return 0
	}

	// 检查每个作者是否有作品
	authorsWithoutBooks := 0
	for _, author := range authors {
		count, err := booksCollection.CountDocuments(ctx, bson.M{"author_id": author.ID})
		if err != nil {
			continue
		}
		if count == 0 {
			authorsWithoutBooks++
		}
	}

	return int64(authorsWithoutBooks)
}

func countInvalidCategories(ctx context.Context, db *mongo.Database) int64 {
	booksCollection := db.Collection("books")
	categoriesCollection := db.Collection("categories")

	// 获取所有有效的分类名称
	cursor, err := categoriesCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error finding categories: %v", err)
		return 0
	}
	defer cursor.Close(ctx)

	type Category struct {
		Name string `bson:"name"`
	}
	var categories []Category
	if err = cursor.All(ctx, &categories); err != nil {
		log.Printf("Error decoding categories: %v", err)
		return 0
	}

	validCategories := make([]string, 0, len(categories))
	for _, cat := range categories {
		validCategories = append(validCategories, cat.Name)
	}

	// 统计无效分类的书籍 - 使用 $not 和 $in 组合
	var invalidCount int64
	if len(validCategories) > 0 {
		// 有分类存在时，查找不属于任何有效分类的书籍
		cursor, err := booksCollection.Find(ctx, bson.M{})
		if err != nil {
			log.Printf("Error finding books: %v", err)
			return 0
		}
		defer cursor.Close(ctx)

		type Book struct {
			ID       primitive.ObjectID `bson:"_id"`
			Category string             `bson:"category"`
		}
		var books []Book
		if err = cursor.All(ctx, &books); err != nil {
			log.Printf("Error decoding books: %v", err)
			return 0
		}

		validCatMap := make(map[string]bool)
		for _, cat := range validCategories {
			validCatMap[cat] = true
		}

		for _, book := range books {
			if book.Category == "" || !validCatMap[book.Category] {
				invalidCount++
			}
		}
	} else {
		// 没有分类存在时，所有有分类字段的书籍都是无效的
		invalidCount, _ = booksCollection.CountDocuments(ctx, bson.M{
			"category": bson.M{"$exists": true, "$ne": ""},
		})
	}

	return invalidCount
}

func countEmptyCategories(ctx context.Context, db *mongo.Database) int64 {
	categoriesCollection := db.Collection("categories")
	booksCollection := db.Collection("books")

	cursor, err := categoriesCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error finding categories: %v", err)
		return 0
	}
	defer cursor.Close(ctx)

	type Category struct {
		ID   primitive.ObjectID `bson:"_id"`
		Name string             `bson:"name"`
	}
	var categories []Category
	if err = cursor.All(ctx, &categories); err != nil {
		log.Printf("Error decoding categories: %v", err)
		return 0
	}

	emptyCategories := 0
	for _, cat := range categories {
		count, err := booksCollection.CountDocuments(ctx, bson.M{"category": cat.Name})
		if err != nil {
			continue
		}
		if count == 0 {
			emptyCategories++
		}
	}

	return int64(emptyCategories)
}

func countInvalidBooksInRankings(ctx context.Context, db *mongo.Database) int64 {
	rankingsCollection := db.Collection("rankings")
	booksCollection := db.Collection("books")

	cursor, err := rankingsCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error finding rankings: %v", err)
		return 0
	}
	defer cursor.Close(ctx)

	type Ranking struct {
		ID    primitive.ObjectID   `bson:"_id"`
		Books []primitive.ObjectID `bson:"books"`
	}
	var rankings []Ranking
	if err = cursor.All(ctx, &rankings); err != nil {
		log.Printf("Error decoding rankings: %v", err)
		return 0
	}

	invalidBooks := 0
	for _, ranking := range rankings {
		for _, bookID := range ranking.Books {
			count, err := booksCollection.CountDocuments(ctx, bson.M{"_id": bookID})
			if err != nil {
				continue
			}
			if count == 0 {
				invalidBooks++
			}
		}
	}

	return int64(invalidBooks)
}

func countExpiredAnnouncements(ctx context.Context, db *mongo.Database) int64 {
	collection := db.Collection("announcements")
	count, err := collection.CountDocuments(ctx, bson.M{
		"end_time": bson.M{"$lt": time.Now()},
	})
	if err != nil {
		log.Printf("Error counting expired announcements: %v", err)
		return 0
	}
	return count
}

func countUnpublishedAnnouncements(ctx context.Context, db *mongo.Database) int64 {
	collection := db.Collection("announcements")
	count, err := collection.CountDocuments(ctx, bson.M{
		"status": bson.M{"$ne": "published"},
	})
	if err != nil {
		log.Printf("Error counting unpublished announcements: %v", err)
		return 0
	}
	return count
}

func countBannedUsersActive(ctx context.Context, db *mongo.Database) int64 {
	collection := db.Collection("users")
	count, err := collection.CountDocuments(ctx, bson.M{
		"status":      "banned",
		"last_login_at": bson.M{"$gt": time.Now().Add(-24 * time.Hour)},
	})
	if err != nil {
		log.Printf("Error counting banned active users: %v", err)
		return 0
	}
	return count
}

func countExpiredVIPUsers(ctx context.Context, db *mongo.Database) int64 {
	collection := db.Collection("users")
	count, err := collection.CountDocuments(ctx, bson.M{
		"role":          "vip",
		"vip_expire_at": bson.M{"$lt": time.Now()},
	})
	if err != nil {
		log.Printf("Error counting expired VIP users: %v", err)
		return 0
	}
	return count
}

func getUserStatusStats(ctx context.Context, db *mongo.Database) map[string]int64 {
	collection := db.Collection("users")

	// 直接查询所有用户，按状态分组统计
	type User struct {
		Status string `bson:"status"`
	}

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error finding users for status stats: %v", err)
		return nil
	}
	defer cursor.Close(ctx)

	var users []User
	if err = cursor.All(ctx, &users); err != nil {
		log.Printf("Error decoding users for status stats: %v", err)
		return nil
	}

	stats := make(map[string]int64)
	for _, user := range users {
		status := user.Status
		if status == "" {
			status = "(empty)"
		}
		stats[status]++
	}
	return stats
}

type InvalidBook struct {
	ID       primitive.ObjectID `bson:"_id"`
	Category string             `bson:"category"`
	Title    string             `bson:"title"`
}

func getInvalidCategoriesDetail(ctx context.Context, db *mongo.Database) []InvalidBook {
	booksCollection := db.Collection("books")
	categoriesCollection := db.Collection("categories")

	// 获取所有有效的分类名称
	cursor, err := categoriesCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error finding categories: %v", err)
		return nil
	}
	defer cursor.Close(ctx)

	type Category struct {
		Name string `bson:"name"`
	}
	var categories []Category
	if err = cursor.All(ctx, &categories); err != nil {
		log.Printf("Error decoding categories: %v", err)
		return nil
	}

	validCatMap := make(map[string]bool)
	for _, cat := range categories {
		validCatMap[cat.Name] = true
	}

	// 查找所有有分类的书籍
	cursor, err = booksCollection.Find(ctx, bson.M{"category": bson.M{"$exists": true, "$ne": ""}})
	if err != nil {
		log.Printf("Error finding books: %v", err)
		return nil
	}
	defer cursor.Close(ctx)

	var books []InvalidBook
	if err = cursor.All(ctx, &books); err != nil {
		log.Printf("Error decoding books: %v", err)
		return nil
	}

	var invalidBooks []InvalidBook
	for _, book := range books {
		if !validCatMap[book.Category] {
			invalidBooks = append(invalidBooks, book)
		}
	}

	return invalidBooks
}

type EmptyCategory struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}

func getEmptyCategoriesDetail(ctx context.Context, db *mongo.Database) []EmptyCategory {
	categoriesCollection := db.Collection("categories")
	booksCollection := db.Collection("books")

	cursor, err := categoriesCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error finding categories: %v", err)
		return nil
	}
	defer cursor.Close(ctx)

	type Category struct {
		ID   primitive.ObjectID `bson:"_id"`
		Name string             `bson:"name"`
	}
	var categories []Category
	if err = cursor.All(ctx, &categories); err != nil {
		log.Printf("Error decoding categories: %v", err)
		return nil
	}

	var emptyCats []Category
	for _, cat := range categories {
		count, err := booksCollection.CountDocuments(ctx, bson.M{"category": cat.Name})
		if err != nil {
			continue
		}
		if count == 0 {
			emptyCats = append(emptyCats, cat)
		}
	}

	result := make([]EmptyCategory, len(emptyCats))
	for i, cat := range emptyCats {
		result[i].ID = cat.ID
		result[i].Name = cat.Name
	}

	return result
}

type UnpublishedAnnouncement struct {
	ID    primitive.ObjectID `bson:"_id"`
	Title string             `bson:"title"`
	Status string             `bson:"status"`
}

func getUnpublishedAnnouncementsDetail(ctx context.Context, db *mongo.Database) []UnpublishedAnnouncement {
	collection := db.Collection("announcements")

	cursor, err := collection.Find(ctx, bson.M{"status": bson.M{"$ne": "published"}})
	if err != nil {
		log.Printf("Error finding unpublished announcements: %v", err)
		return nil
	}
	defer cursor.Close(ctx)

	var announcements []UnpublishedAnnouncement
	if err = cursor.All(ctx, &announcements); err != nil {
		log.Printf("Error decoding announcements: %v", err)
		return nil
	}

	return announcements
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
