// Package relationships 处理测试数据的关联关系构建
package relationships

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"Qingyu_backend/cmd/seeder/utils"
)

// RelationshipBuilder 关联关系构建器
type RelationshipBuilder struct {
	db *utils.Database
}

// NewRelationshipBuilder 创建关联关系构建器实例
func NewRelationshipBuilder(db *utils.Database) *RelationshipBuilder {
	return &RelationshipBuilder{db: db}
}

// BuildSubscriptions 构建书籍订阅关系
// 根据书籍评分确定订阅数量，并随机选择用户作为订阅者
func (rb *RelationshipBuilder) BuildSubscriptions() error {
	ctx := context.Background()

	// 获取所有书籍及其评分
	books, err := rb.getBooksWithRating(ctx)
	if err != nil {
		return fmt.Errorf("获取书籍失败: %w", err)
	}

	if len(books) == 0 {
		return fmt.Errorf("没有找到书籍数据")
	}

	// 获取所有用户ID
	userIDs, err := rb.getAllUsers(ctx)
	if err != nil {
		return fmt.Errorf("获取用户失败: %w", err)
	}

	if len(userIDs) == 0 {
		return fmt.Errorf("没有找到用户数据")
	}

	fmt.Printf("开始构建订阅关系: %d 本书, %d 个用户\n", len(books), len(userIDs))

	// 为每本书生成订阅关系
	collection := rb.db.Collection("subscriptions")
	totalSubscriptions := 0

	for _, book := range books {
		// 根据评分确定订阅数
		var targetSubscriptions int
		if book.Rating > 8.5 {
			// 热门书籍: 200-500 个订阅
			targetSubscriptions = rand.Intn(300) + 200
		} else if book.Rating > 6.0 {
			// 普通书籍: 20-200 个订阅
			targetSubscriptions = rand.Intn(180) + 20
		} else {
			// 冷门书籍: 0-20 个订阅
			targetSubscriptions = rand.Intn(20)
		}

		// 如果目标订阅数超过用户总数，则限制为用户总数
		if targetSubscriptions > len(userIDs) {
			targetSubscriptions = len(userIDs)
		}

		// 随机选择订阅用户
		subscribers := rb.selectRandomUsers(userIDs, targetSubscriptions)

		// 创建订阅文档
		subscriptions := make([]interface{}, 0, len(subscribers))
		for _, userID := range subscribers {
			// 订阅时间随机分布在书籍发布时间到现在
			subscribedAt := rb.randomTime(book.PublishedAt, time.Now())

			subscription := bson.M{
				"user_id":       userID,
				"book_id":       book.ID,
				"subscribed_at": subscribedAt,
				"created_at":    time.Now(),
			}
			subscriptions = append(subscriptions, subscription)
		}

		// 批量插入订阅关系
		if len(subscriptions) > 0 {
			_, err := collection.InsertMany(ctx, subscriptions)
			if err != nil {
				return fmt.Errorf("插入订阅关系失败 (书籍ID: %s): %w", book.ID, err)
			}
			totalSubscriptions += len(subscriptions)
		}
	}

	fmt.Printf("订阅关系构建完成: 共创建 %d 条订阅记录\n", totalSubscriptions)
	return nil
}

// BookInfo 书籍信息结构体
type BookInfo struct {
	ID          string    `bson:"_id"`
	Title       string    `bson:"title"`
	Rating      float64   `bson:"rating"`
	PublishedAt time.Time `bson:"published_at"`
}

// getBooksWithRating 获取所有书籍及其评分
func (rb *RelationshipBuilder) getBooksWithRating(ctx context.Context) ([]BookInfo, error) {
	collection := rb.db.Collection("books")

	// 查询所有书籍，只获取需要的字段
	cursor, err := collection.Find(ctx, bson.M{},
		// 使用 projection 只获取需要的字段
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var books []BookInfo
	if err := cursor.All(ctx, &books); err != nil {
		return nil, err
	}

	return books, nil
}

// getAllUsers 获取所有用户ID
func (rb *RelationshipBuilder) getAllUsers(ctx context.Context) ([]string, error) {
	collection := rb.db.Collection("users")

	// 查询所有用户ID
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []struct {
		ID string `bson:"_id"`
	}
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	// 提取所有用户ID
	userIDs := make([]string, len(users))
	for i, user := range users {
		userIDs[i] = user.ID
	}

	return userIDs, nil
}

// selectRandomUsers 从用户列表中随机选择指定数量的用户
func (rb *RelationshipBuilder) selectRandomUsers(userIDs []string, count int) []string {
	if count >= len(userIDs) {
		// 如果请求数量大于等于用户总数，返回所有用户（打乱顺序）
		shuffled := make([]string, len(userIDs))
		copy(shuffled, userIDs)
		rand.Shuffle(len(shuffled), func(i, j int) {
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		})
		return shuffled
	}

	// 使用 Fisher-Yates 洗牌算法选择随机用户
	shuffled := make([]string, len(userIDs))
	copy(shuffled, userIDs)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled[:count]
}

// randomTime 生成在指定时间范围内的随机时间
func (rb *RelationshipBuilder) randomTime(start, end time.Time) time.Time {
	// 计算时间差（秒）
	delta := end.Sub(start).Seconds()

	// 生成随机秒数（0到delta之间）
	randomSeconds := rand.Int63n(int64(delta))

	// 返回随机时间
	return start.Add(time.Duration(randomSeconds) * time.Second)
}
