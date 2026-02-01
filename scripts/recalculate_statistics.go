// 统计数据重新计算工具
// 用于修复 books 和 users 集合中不准确的统计字段
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database MongoDB 数据库包装器
type Database struct {
	*mongo.Database
}

// ConnectDB 连接 MongoDB 数据库
func ConnectDB(uri, database string) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("创建 MongoDB 客户端失败: %w", err)
	}

	// 测试连接
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("连接 MongoDB 失败: %w", err)
	}

	db := client.Database(database)
	return &Database{db}, nil
}

type Book struct {
	ID           interface{} `bson:"_id"`
	Title        string      `bson:"title"`
	LikesCount   interface{} `bson:"likes_count"`
	CommentsCount interface{} `bson:"comments_count"`
}

type UpdateExample struct {
	Title  string
	Before interface{}
	After  int64
}

func main() {
	// 连接数据库
	uri := "mongodb://localhost:27017"
	databaseName := "qingyu"

	db, err := ConnectDB(uri, databaseName)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Client().Disconnect(context.Background())

	fmt.Println("=== 统计数据重新计算工具 ===")
	fmt.Printf("连接到数据库: %s\n", databaseName)
	fmt.Printf("时间: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	// 1. 重新计算书籍点赞数
	fmt.Println("1. 重新计算书籍点赞数...")
	booksLikesUpdated, booksLikesExamples := recalculateBooksLikes(db)
	fmt.Printf("   更新数量: %d\n", booksLikesUpdated)
	if len(booksLikesExamples) > 0 {
		fmt.Println("   更新示例:")
		for _, ex := range booksLikesExamples {
			fmt.Printf("     \"%s\": %v -> %d\n", ex.Title, ex.Before, ex.After)
		}
	}
	fmt.Println()

	// 2. 重新计算书籍评论数
	fmt.Println("2. 重新计算书籍评论数...")
	booksCommentsUpdated, booksCommentsExamples := recalculateBooksComments(db)
	fmt.Printf("   更新数量: %d\n", booksCommentsUpdated)
	if len(booksCommentsExamples) > 0 {
		fmt.Println("   更新示例:")
		for _, ex := range booksCommentsExamples {
			fmt.Printf("     \"%s\": %v -> %d\n", ex.Title, ex.Before, ex.After)
		}
	}
	fmt.Println()

	// 3. 重新计算用户粉丝数
	fmt.Println("3. 重新计算用户粉丝数...")
	usersFollowersUpdated, usersFollowersExamples := recalculateUsersFollowers(db)
	fmt.Printf("   更新数量: %d\n", usersFollowersUpdated)
	if len(usersFollowersExamples) > 0 {
		fmt.Println("   更新示例:")
		for _, ex := range usersFollowersExamples {
			fmt.Printf("     用户ID %s: %v -> %d\n", ex.Title, ex.Before, ex.After)
		}
	}
	fmt.Println()

	// 4. 重新计算用户关注数
	fmt.Println("4. 重新计算用户关注数...")
	usersFollowingUpdated, usersFollowingExamples := recalculateUsersFollowing(db)
	fmt.Printf("   更新数量: %d\n", usersFollowingUpdated)
	if len(usersFollowingExamples) > 0 {
		fmt.Println("   更新示例:")
		for _, ex := range usersFollowingExamples {
			fmt.Printf("     用户ID %s: %v -> %d\n", ex.Title, ex.Before, ex.After)
		}
	}
	fmt.Println()

	// 5. 验证结果
	fmt.Println("5. 验证结果...")
	validationResults := validateStatistics(db)
	fmt.Printf("   书籍点赞数不一致: %d\n", validationResults["books_likes_inconsistent"])
	fmt.Printf("   书籍评论数不一致: %d\n", validationResults["books_comments_inconsistent"])
	fmt.Printf("   用户粉丝数不一致: %d\n", validationResults["users_followers_inconsistent"])
	fmt.Printf("   用户关注数不一致: %d\n", validationResults["users_following_inconsistent"])
	fmt.Println()

	// 生成报告
	fmt.Println("=== 生成报告 ===")
	reportPath := "docs/reports/2026-02-01-statistics-recalc-report.md"
	generateReport(reportPath, booksLikesUpdated, booksCommentsUpdated, usersFollowersUpdated, usersFollowingUpdated,
		booksLikesExamples, booksCommentsExamples, usersFollowersExamples, usersFollowingExamples, validationResults)
	fmt.Printf("报告已保存到: %s\n", reportPath)

	// 检查是否所有数据都一致
	allConsistent := validationResults["books_likes_inconsistent"] == 0 &&
		validationResults["books_comments_inconsistent"] == 0 &&
		validationResults["users_followers_inconsistent"] == 0 &&
		validationResults["users_following_inconsistent"] == 0

	if allConsistent {
		fmt.Println("\n✅ 所有统计字段已修复完成！")
		os.Exit(0)
	} else {
		fmt.Println("\n⚠️ 部分统计字段仍存在不一致，请检查报告了解详情")
		os.Exit(1)
	}
}

// recalculateBooksLikes 重新计算书籍点赞数
func recalculateBooksLikes(db *Database) (int, []UpdateExample) {
	ctx := context.Background()
	collection := db.Collection("books")
	likesCollection := db.Collection("likes")

	// 获取所有书籍
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("获取书籍列表失败: %v", err)
		return 0, nil
	}
	defer cursor.Close(ctx)

	updatedCount := 0
	examples := []UpdateExample{}

	for cursor.Next(ctx) {
		var book Book
		if err := cursor.Decode(&book); err != nil {
			continue
		}

		// 计算实际点赞数
		actualCount, _ := likesCollection.CountDocuments(ctx, bson.M{"target_id": book.ID})

		// 检查是否需要更新
		needsUpdate := false
		switch v := book.LikesCount.(type) {
		case int:
			needsUpdate = int64(v) != actualCount
		case int32:
			needsUpdate = int64(v) != actualCount
		case int64:
			needsUpdate = v != actualCount
		case float64:
			needsUpdate = int64(v) != actualCount
		case string:
			// 如果是字符串类型，可能需要解析
			needsUpdate = true
		default:
			needsUpdate = true
		}

		if needsUpdate {
			_, err := collection.UpdateOne(ctx, bson.M{"_id": book.ID}, bson.M{"$set": bson.M{"likes_count": actualCount}})
			if err == nil {
				updatedCount++
				if len(examples) < 5 {
					examples = append(examples, UpdateExample{
						Title:  book.Title,
						Before: book.LikesCount,
						After:  actualCount,
					})
				}
			}
		}
	}

	return updatedCount, examples
}

// recalculateBooksComments 重新计算书籍评论数
func recalculateBooksComments(db *Database) (int, []UpdateExample) {
	ctx := context.Background()
	collection := db.Collection("books")
	commentsCollection := db.Collection("comments")

	// 获取所有书籍
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("获取书籍列表失败: %v", err)
		return 0, nil
	}
	defer cursor.Close(ctx)

	updatedCount := 0
	examples := []UpdateExample{}

	for cursor.Next(ctx) {
		var book Book
		if err := cursor.Decode(&book); err != nil {
			continue
		}

		// 计算实际评论数
		actualCount, _ := commentsCollection.CountDocuments(ctx, bson.M{"target_id": book.ID})

		// 检查是否需要更新
		needsUpdate := false
		switch v := book.CommentsCount.(type) {
		case int:
			needsUpdate = int64(v) != actualCount
		case int32:
			needsUpdate = int64(v) != actualCount
		case int64:
			needsUpdate = v != actualCount
		case float64:
			needsUpdate = int64(v) != actualCount
		case string:
			needsUpdate = true
		default:
			needsUpdate = true
		}

		if needsUpdate {
			_, err := collection.UpdateOne(ctx, bson.M{"_id": book.ID}, bson.M{"$set": bson.M{"comments_count": actualCount}})
			if err == nil {
				updatedCount++
				if len(examples) < 5 {
					examples = append(examples, UpdateExample{
						Title:  book.Title,
						Before: book.CommentsCount,
						After:  actualCount,
					})
				}
			}
		}
	}

	return updatedCount, examples
}

// recalculateUsersFollowers 重新计算用户粉丝数
func recalculateUsersFollowers(db *Database) (int, []UpdateExample) {
	ctx := context.Background()
	collection := db.Collection("users")
	relationsCollection := db.Collection("user_relations")

	// 获取所有用户
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("获取用户列表失败: %v", err)
		return 0, nil
	}
	defer cursor.Close(ctx)

	updatedCount := 0
	examples := []UpdateExample{}

	for cursor.Next(ctx) {
		var user struct {
			ID            interface{} `bson:"_id"`
			Username      string      `bson:"username"`
			FollowersCount interface{} `bson:"followers_count"`
		}
		if err := cursor.Decode(&user); err != nil {
			continue
		}

		// 计算实际粉丝数
		actualCount, _ := relationsCollection.CountDocuments(ctx, bson.M{"following_id": user.ID})

		// 检查是否需要更新
		needsUpdate := false
		switch v := user.FollowersCount.(type) {
		case int:
			needsUpdate = int64(v) != actualCount
		case int32:
			needsUpdate = int64(v) != actualCount
		case int64:
			needsUpdate = v != actualCount
		case float64:
			needsUpdate = int64(v) != actualCount
		case string:
			needsUpdate = true
		default:
			needsUpdate = true
		}

		if needsUpdate {
			_, err := collection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"followers_count": actualCount}})
			if err == nil {
				updatedCount++
				if len(examples) < 5 {
					examples = append(examples, UpdateExample{
						Title:  fmt.Sprintf("%v", user.ID),
						Before: user.FollowersCount,
						After:  actualCount,
					})
				}
			}
		}
	}

	return updatedCount, examples
}

// recalculateUsersFollowing 重新计算用户关注数
func recalculateUsersFollowing(db *Database) (int, []UpdateExample) {
	ctx := context.Background()
	collection := db.Collection("users")
	relationsCollection := db.Collection("user_relations")

	// 获取所有用户
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("获取用户列表失败: %v", err)
		return 0, nil
	}
	defer cursor.Close(ctx)

	updatedCount := 0
	examples := []UpdateExample{}

	for cursor.Next(ctx) {
		var user struct {
			ID            interface{} `bson:"_id"`
			Username      string      `bson:"username"`
			FollowingCount interface{} `bson:"following_count"`
		}
		if err := cursor.Decode(&user); err != nil {
			continue
		}

		// 计算实际关注数
		actualCount, _ := relationsCollection.CountDocuments(ctx, bson.M{"follower_id": user.ID})

		// 检查是否需要更新
		needsUpdate := false
		switch v := user.FollowingCount.(type) {
		case int:
			needsUpdate = int64(v) != actualCount
		case int32:
			needsUpdate = int64(v) != actualCount
		case int64:
			needsUpdate = v != actualCount
		case float64:
			needsUpdate = int64(v) != actualCount
		case string:
			needsUpdate = true
		default:
			needsUpdate = true
		}

		if needsUpdate {
			_, err := collection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"following_count": actualCount}})
			if err == nil {
				updatedCount++
				if len(examples) < 5 {
					examples = append(examples, UpdateExample{
						Title:  fmt.Sprintf("%v", user.ID),
						Before: user.FollowingCount,
						After:  actualCount,
					})
				}
			}
		}
	}

	return updatedCount, examples
}

// validateStatistics 验证统计字段
func validateStatistics(db *Database) map[string]int64 {
	ctx := context.Background()
	results := make(map[string]int64)

	// 验证书籍点赞数
	booksCollection := db.Collection("books")
	likesCollection := db.Collection("likes")

	booksCursor, _ := booksCollection.Find(ctx, bson.M{})
	defer booksCursor.Close(ctx)

	booksLikesInconsistent := int64(0)
	for booksCursor.Next(ctx) {
		var book struct {
			ID         interface{} `bson:"_id"`
			LikesCount interface{} `bson:"likes_count"`
		}
		booksCursor.Decode(&book)

		actualCount, _ := likesCollection.CountDocuments(ctx, bson.M{"target_id": book.ID})

		storedCount := int64(0)
		switch v := book.LikesCount.(type) {
		case int:
			storedCount = int64(v)
		case int32:
			storedCount = int64(v)
		case int64:
			storedCount = v
		case float64:
			storedCount = int64(v)
		}

		if storedCount != actualCount {
			booksLikesInconsistent++
		}
	}
	results["books_likes_inconsistent"] = booksLikesInconsistent

	// 验证书籍评论数
	commentsCollection := db.Collection("comments")

	booksCursor2, _ := booksCollection.Find(ctx, bson.M{})
	defer booksCursor2.Close(ctx)

	booksCommentsInconsistent := int64(0)
	for booksCursor2.Next(ctx) {
		var book struct {
			ID            interface{} `bson:"_id"`
			CommentsCount interface{} `bson:"comments_count"`
		}
		booksCursor2.Decode(&book)

		actualCount, _ := commentsCollection.CountDocuments(ctx, bson.M{"target_id": book.ID})

		storedCount := int64(0)
		switch v := book.CommentsCount.(type) {
		case int:
			storedCount = int64(v)
		case int32:
			storedCount = int64(v)
		case int64:
			storedCount = v
		case float64:
			storedCount = int64(v)
		}

		if storedCount != actualCount {
			booksCommentsInconsistent++
		}
	}
	results["books_comments_inconsistent"] = booksCommentsInconsistent

	// 验证用户粉丝数
	usersCollection := db.Collection("users")
	relationsCollection := db.Collection("user_relations")

	usersCursor, _ := usersCollection.Find(ctx, bson.M{})
	defer usersCursor.Close(ctx)

	usersFollowersInconsistent := int64(0)
	for usersCursor.Next(ctx) {
		var user struct {
			ID             interface{} `bson:"_id"`
			FollowersCount interface{} `bson:"followers_count"`
		}
		usersCursor.Decode(&user)

		actualCount, _ := relationsCollection.CountDocuments(ctx, bson.M{"following_id": user.ID})

		storedCount := int64(0)
		switch v := user.FollowersCount.(type) {
		case int:
			storedCount = int64(v)
		case int32:
			storedCount = int64(v)
		case int64:
			storedCount = v
		case float64:
			storedCount = int64(v)
		}

		if storedCount != actualCount {
			usersFollowersInconsistent++
		}
	}
	results["users_followers_inconsistent"] = usersFollowersInconsistent

	// 验证用户关注数
	usersCursor2, _ := usersCollection.Find(ctx, bson.M{})
	defer usersCursor2.Close(ctx)

	usersFollowingInconsistent := int64(0)
	for usersCursor2.Next(ctx) {
		var user struct {
			ID             interface{} `bson:"_id"`
			FollowingCount interface{} `bson:"following_count"`
		}
		usersCursor2.Decode(&user)

		actualCount, _ := relationsCollection.CountDocuments(ctx, bson.M{"follower_id": user.ID})

		storedCount := int64(0)
		switch v := user.FollowingCount.(type) {
		case int:
			storedCount = int64(v)
		case int32:
			storedCount = int64(v)
		case int64:
			storedCount = v
		case float64:
			storedCount = int64(v)
		}

		if storedCount != actualCount {
			usersFollowingInconsistent++
		}
	}
	results["users_following_inconsistent"] = usersFollowingInconsistent

	return results
}

// generateReport 生成报告
func generateReport(path string, booksLikesUpdated, booksCommentsUpdated, usersFollowersUpdated, usersFollowingUpdated int,
	booksLikesExamples, booksCommentsExamples, usersFollowersExamples, usersFollowingExamples []UpdateExample,
	validationResults map[string]int64) {

	content := fmt.Sprintf(`# 统计数据重新计算报告

**生成时间**: %s
**执行人**: 数据库修复专家女仆

## 执行摘要

本次任务重新计算了所有不准确的统计字段，确保数据一致性。

## 更新详情

### 1. 书籍点赞数 (books.likes_count)

- **更新数量**: %d
- **状态**: %s

更新示例:
`, time.Now().Format("2006-01-02 15:04:05"), booksLikesUpdated, getStatus(booksLikesUpdated > 0))

	if len(booksLikesExamples) > 0 {
		for _, ex := range booksLikesExamples {
			content += fmt.Sprintf("- `%s`: %v → %d\n", ex.Title, ex.Before, ex.After)
		}
	} else {
		content += "无更新（数据已正确）\n"
	}

	content += fmt.Sprintf(`
### 2. 书籍评论数 (books.comments_count)

- **更新数量**: %d
- **状态**: %s

更新示例:
`, booksCommentsUpdated, getStatus(booksCommentsUpdated > 0))

	if len(booksCommentsExamples) > 0 {
		for _, ex := range booksCommentsExamples {
			content += fmt.Sprintf("- `%s`: %v → %d\n", ex.Title, ex.Before, ex.After)
		}
	} else {
		content += "无更新（数据已正确）\n"
	}

	content += fmt.Sprintf(`
### 3. 用户粉丝数 (users.followers_count)

- **更新数量**: %d
- **状态**: %s

更新示例:
`, usersFollowersUpdated, getStatus(usersFollowersUpdated > 0))

	if len(usersFollowersExamples) > 0 {
		for _, ex := range usersFollowersExamples {
			content += fmt.Sprintf("- 用户ID %s: %v → %d\n", ex.Title, ex.Before, ex.After)
		}
	} else {
		content += "无更新（数据已正确）\n"
	}

	content += fmt.Sprintf(`
### 4. 用户关注数 (users.following_count)

- **更新数量**: %d
- **状态**: %s

更新示例:
`, usersFollowingUpdated, getStatus(usersFollowingUpdated > 0))

	if len(usersFollowingExamples) > 0 {
		for _, ex := range usersFollowingExamples {
			content += fmt.Sprintf("- 用户ID %s: %v → %d\n", ex.Title, ex.Before, ex.After)
		}
	} else {
		content += "无更新（数据已正确）\n"
	}

	// 验证结果
	content += `
## 验证结果

执行后的数据一致性验证:

`
	allConsistent := validationResults["books_likes_inconsistent"] == 0 &&
		validationResults["books_comments_inconsistent"] == 0 &&
		validationResults["users_followers_inconsistent"] == 0 &&
		validationResults["users_following_inconsistent"] == 0

	if validationResults["books_likes_inconsistent"] == 0 {
		content += "- ✅ 书籍点赞数: 全部一致\n"
	} else {
		content += fmt.Sprintf("- ❌ 书籍点赞数: %d 条不一致\n", validationResults["books_likes_inconsistent"])
	}

	if validationResults["books_comments_inconsistent"] == 0 {
		content += "- ✅ 书籍评论数: 全部一致\n"
	} else {
		content += fmt.Sprintf("- ❌ 书籍评论数: %d 条不一致\n", validationResults["books_comments_inconsistent"])
	}

	if validationResults["users_followers_inconsistent"] == 0 {
		content += "- ✅ 用户粉丝数: 全部一致\n"
	} else {
		content += fmt.Sprintf("- ❌ 用户粉丝数: %d 条不一致\n", validationResults["users_followers_inconsistent"])
	}

	if validationResults["users_following_inconsistent"] == 0 {
		content += "- ✅ 用户关注数: 全部一致\n"
	} else {
		content += fmt.Sprintf("- ❌ 用户关注数: %d 条不一致\n", validationResults["users_following_inconsistent"])
	}

	// 总结
	content += `
## 总结

`

	if allConsistent {
		content += `✅ **任务完成！所有统计字段已修复，数据完全一致。**

本次重新计算成功修复了所有不准确的统计字段，确保了:
- 书籍的点赞数与实际的 likes 记录一致
- 书籍的评论数与实际的 comments 记录一致
- 用户的粉丝数与实际的 user_relations 记录一致
- 用户的关注数与实际的 user_relations 记录一致
`
	} else {
		content += `⚠️ **部分数据仍存在问题**

尽管进行了重新计算，仍有部分数据存在不一致。可能的原因:
- 数据库中存在孤立记录（如 likes 表中的 target_id 在 books 中不存在）
- 数据类型转换问题
- 并发写入导致的竞态条件

建议:
1. 检查是否存在孤立记录
2. 检查应用程序的统计字段更新逻辑
3. 考虑添加数据库约束或触发器
`
	}

	content += `
## 执行信息

- **数据库**: qingyu
- **连接地址**: mongodb://localhost:27017
- **执行时间**: ` + time.Now().Format("2006-01-02 15:04:05") + `
- **报告生成**: 2026-02-01-statistics-recalc-report.md
`

	// 确保目录存在
	os.MkdirAll("docs/reports", 0755)

	// 写入文件
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		log.Printf("生成报告失败: %v", err)
	}
}

func getStatus(updated bool) string {
	if updated {
		return "已更新"
	}
	return "无需更新"
}
