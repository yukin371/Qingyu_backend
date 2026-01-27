package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	authModel "Qingyu_backend/models/auth"
)

// Config 配置
type Config struct {
	MongoURI string
	DBName   string
}

// TestData 测试数据
type TestData struct {
	Roles []authModel.Role
	Users []User
}

// User 用户
type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username       string             `bson:"username" json:"username"`
	Email          string             `bson:"email" json:"email"`
	Password       string             `bson:"password" json:"password"`
	HashedPassword string             `bson:"hashed_password" json:"-"`
	Roles          []string           `bson:"roles" json:"roles"`
	IsVip          bool               `bson:"is_vip" json:"is_vip"`
	IsAdmin        bool               `bson:"is_admin" json:"is_admin"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

var (
	dbName  = flag.String("db", "qingyu_permission_test", "测试数据库名称")
	mongoURI = flag.String("uri", "mongodb://localhost:27017", "MongoDB连接URI")
	verbose = flag.Bool("v", false, "详细输出")
)

func main() {
	flag.Parse()

	log.Println("========================================")
	log.Println(" Qingyu Backend - 权限测试数据填充工具")
	log.Println("========================================")
	log.Println("")

	// 1. 连接MongoDB
	log.Printf("[1/5] 连接MongoDB: %s", *mongoURI)
	client, err := connectMongoDB(*mongoURI)
	if err != nil {
		log.Fatalf("连接MongoDB失败: %v", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("断开连接失败: %v", err)
		}
	}()
	log.Println("✓ MongoDB连接成功")

	db := client.Database(*dbName)

	// 2. 清理旧数据
	log.Println("[2/5] 清理旧数据...")
	if err := cleanOldData(db); err != nil {
		log.Printf("清理旧数据失败: %v", err)
	} else {
		log.Println("✓ 旧数据清理完成")
	}

	// 3. 创建角色
	log.Println("[3/5] 创建角色...")
	roles := createTestRoles()
	if err := insertRoles(db, roles); err != nil {
		log.Fatalf("插入角色失败: %v", err)
	}
	log.Printf("✓ 创建角色 %d 个", len(roles))

	if *verbose {
		for _, role := range roles {
			log.Printf("  - %s: %d 个权限", role.Name, len(role.Permissions))
		}
	}

	// 4. 创建用户
	log.Println("[4/5] 创建用户...")
	users := createTestUsers()
	if err := insertUsers(db, users); err != nil {
		log.Fatalf("插入用户失败: %v", err)
	}
	log.Printf("✓ 创建用户 %d 个", len(users))

	if *verbose {
		for _, user := range users {
			log.Printf("  - %s (%s): 角色 %v", user.Username, user.Email, user.Roles)
		}
	}

	// 5. 验证数据
	log.Println("[5/5] 验证数据...")
	if err := verifyData(db); err != nil {
		log.Fatalf("数据验证失败: %v", err)
	}
	log.Println("✓ 数据验证完成")

	// 打印摘要
	printSummary(users)

	log.Println("")
	log.Println("========================================")
	log.Println(" 测试数据填充完成！")
	log.Println("========================================")
}

// connectMongoDB 连接MongoDB
func connectMongoDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// 测试连接
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}

// cleanOldData 清理旧数据
func cleanOldData(db *mongo.Database) error {
	collections := []string{"roles", "users"}

	for _, colName := range collections {
		col := db.Collection(colName)
		count, err := col.CountDocuments(context.Background(), bson.M{})
		if err != nil {
			return err
		}

		if count > 0 {
			_, err := col.DeleteMany(context.Background(), bson.M{})
			if err != nil {
				return err
			}
			log.Printf("  清理集合 %s: 删除 %d 条记录", colName, count)
		}
	}

	return nil
}

// createTestRoles 创建测试角色
func createTestRoles() []authModel.Role {
	now := time.Now()

	return []authModel.Role{
		{
			ID:          primitive.NewObjectID().Hex(),
			Name:        "admin",
			Description: "系统管理员，拥有所有权限",
			Permissions: []string{
				"*:*", // 完全通配符权限
			},
			IsSystem:  true,
			IsDefault: false,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:          primitive.NewObjectID().Hex(),
			Name:        "author",
			Description: "作者，可以创作和管理自己的作品",
			Permissions: []string{
				"book:read",
				"book:create",
				"book:update",
				"book:delete",
				"chapter:read",
				"chapter:create",
				"chapter:update",
				"chapter:delete",
				"ai:generate",
				"ai:chat",
				"document:read",
				"document:create",
				"document:update",
				"comment:read",
				"comment:create",
				"comment:update",
			},
			IsSystem:  true,
			IsDefault: false,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:          primitive.NewObjectID().Hex(),
			Name:        "reader",
			Description: "读者，可以阅读内容",
			Permissions: []string{
				"book:read",
				"chapter:read",
				"document:read",
				"comment:read",
				"comment:create",
			},
			IsSystem:  true,
			IsDefault: true, // 新用户默认为读者
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:          primitive.NewObjectID().Hex(),
			Name:        "editor",
			Description: "编辑，可以审核和管理内容",
			Permissions: []string{
				"book:read",
				"book:update",
				"book:review",
				"chapter:read",
				"chapter:update",
				"chapter:review",
				"comment:read",
				"comment:update",
				"comment:delete",
				"document:read",
				"document:update",
			},
			IsSystem:  true,
			IsDefault: false,
			CreatedAt: now,
			UpdatedAt: now,
		},
		// 测试特殊角色
		{
			ID:          primitive.NewObjectID().Hex(),
			Name:        "limited_user",
			Description: "受限用户，只有基本读取权限",
			Permissions: []string{
				"book:read",
			},
			IsSystem:  false,
			IsDefault: false,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

// createTestUsers 创建测试用户
func createTestUsers() []User {
	now := time.Now()
	adminPassword := hashPassword("Admin@123")
	authorPassword := hashPassword("Author@123")
	readerPassword := hashPassword("Reader@123")
	editorPassword := hashPassword("Editor@123")
	limitedPassword := hashPassword("Limited@123")

	return []User{
		{
			ID:             primitive.NewObjectID(),
			Username:       "admin@test.com",
			Email:          "admin@test.com",
			Password:       "Admin@123",
			HashedPassword: adminPassword,
			Roles:          []string{"admin"},
			IsVip:          true,
			IsAdmin:        true,
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		{
			ID:             primitive.NewObjectID(),
			Username:       "author@test.com",
			Email:          "author@test.com",
			Password:       "Author@123",
			HashedPassword: authorPassword,
			Roles:          []string{"author"},
			IsVip:          true,
			IsAdmin:        false,
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		{
			ID:             primitive.NewObjectID(),
			Username:       "reader@test.com",
			Email:          "reader@test.com",
			Password:       "Reader@123",
			HashedPassword: readerPassword,
			Roles:          []string{"reader"},
			IsVip:          false,
			IsAdmin:        false,
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		{
			ID:             primitive.NewObjectID(),
			Username:       "editor@test.com",
			Email:          "editor@test.com",
			Password:       "Editor@123",
			HashedPassword: editorPassword,
			Roles:          []string{"editor"},
			IsVip:          false,
			IsAdmin:        false,
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		{
			ID:             primitive.NewObjectID(),
			Username:       "limited@test.com",
			Email:          "limited@test.com",
			Password:       "Limited@123",
			HashedPassword: limitedPassword,
			Roles:          []string{"limited_user"},
			IsVip:          false,
			IsAdmin:        false,
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		// 多角色用户测试
		{
			ID:             primitive.NewObjectID(),
			Username:       "author_reader@test.com",
			Email:          "author_reader@test.com",
			Password:       "MultiRole@123",
			HashedPassword: hashPassword("MultiRole@123"),
			Roles:          []string{"author", "reader"},
			IsVip:          true,
			IsAdmin:        false,
			CreatedAt:      now,
			UpdatedAt:      now,
		},
	}
}

// insertRoles 插入角色
func insertRoles(db *mongo.Database, roles []authModel.Role) error {
	collection := db.Collection("roles")

	// 批量插入
	docs := make([]interface{}, len(roles))
	for i, role := range roles {
		docs[i] = role
	}

	_, err := collection.InsertMany(context.Background(), docs)
	return err
}

// insertUsers 插入用户
func insertUsers(db *mongo.Database, users []User) error {
	collection := db.Collection("users")

	// 批量插入
	docs := make([]interface{}, len(users))
	for i, user := range users {
		docs[i] = user
	}

	_, err := collection.InsertMany(context.Background(), docs)
	return err
}

// verifyData 验证数据
func verifyData(db *mongo.Database) error {
	// 验证角色数量
	roleCount, err := db.Collection("roles").CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return fmt.Errorf("统计角色数量失败: %w", err)
	}

	if roleCount == 0 {
		return fmt.Errorf("角色数量为0")
	}

	// 验证用户数量
	userCount, err := db.Collection("users").CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return fmt.Errorf("统计用户数量失败: %w", err)
	}

	if userCount == 0 {
		return fmt.Errorf("用户数量为0")
	}

	// 验证admin角色是否有通配符权限
	var adminRole authModel.Role
	err = db.Collection("roles").FindOne(context.Background(), bson.M{"name": "admin"}).Decode(&adminRole)
	if err != nil {
		return fmt.Errorf("查找admin角色失败: %w", err)
	}

	hasWildcard := false
	for _, perm := range adminRole.Permissions {
		if perm == "*:*" || perm == "*" {
			hasWildcard = true
			break
		}
	}

	if !hasWildcard {
		return fmt.Errorf("admin角色缺少通配符权限")
	}

	return nil
}

// hashPassword 哈希密码
func hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("哈希密码失败: %v", err)
	}
	return string(hashedPassword)
}

// printSummary 打印摘要
func printSummary(users []User) {
	fmt.Println("")
	fmt.Println("========================================")
	fmt.Println(" 测试数据摘要")
	fmt.Println("========================================")
	fmt.Println("")
	fmt.Println("角色列表:")
	fmt.Println("  1. admin   - 系统管理员 (所有权限)")
	fmt.Println("  2. author  - 作者 (作品管理)")
	fmt.Println("  3. reader  - 读者 (只读)")
	fmt.Println("  4. editor  - 编辑 (审核管理)")
	fmt.Println("  5. limited_user - 受限用户 (基本读取)")
	fmt.Println("")
	fmt.Println("测试账号:")
	fmt.Println("")

	for _, user := range users {
		roles := ""
		for i, role := range user.Roles {
			if i > 0 {
				roles += ", "
			}
			roles += role
		}

		fmt.Printf("  %s\n", user.Username)
		fmt.Printf("    密码: %s\n", user.Password)
		fmt.Printf("    角色: %s\n", roles)
		fmt.Println("")
	}

	fmt.Println("========================================")
	fmt.Println("")
}
