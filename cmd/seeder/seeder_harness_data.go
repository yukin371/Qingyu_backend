package main

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/cmd/seeder/story"
	"Qingyu_backend/cmd/seeder/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HarnessDataSeeder Harness测试数据填充器
type HarnessDataSeeder struct {
	db *utils.Database
}

// NewHarnessDataSeeder 创建Harness数据填充器
func NewHarnessDataSeeder(db *utils.Database) *HarnessDataSeeder {
	return &HarnessDataSeeder{db: db}
}

// CreateHarnessData 创建Harness测试数据
func (s *HarnessDataSeeder) CreateHarnessData() error {

	// 1. 检查项目是否已存在
	if s.projectExists(story.EmploymentGuideProject.Title) {
		fmt.Println("✓ Harness项目「就业指南：我在异世界当劝退专员」已存在，跳过创建")
		return nil
	}

	fmt.Println("🚀 开始创建Harness数据...")

	// 2. 获取/创建 Harness 专用用户
	harnessUserID, err := s.getOrCreateHarnessUser()
	if err != nil {
		return fmt.Errorf("获取 Harness 账号失败: %w", err)
	}
	fmt.Printf("✓ Harness 账号: %s\n", harnessUserID)

	// 3. 创建项目
	projectID, err := s.createProject(harnessUserID)
	if err != nil {
		return fmt.Errorf("创建项目失败: %w", err)
	}
	fmt.Printf("✓ 创建项目: %s\n", projectID)

	// 4. 创建角色
	characterIDs, err := s.createCharacters(projectID)
	if err != nil {
		return fmt.Errorf("创建角色失败: %w", err)
	}
	fmt.Printf("✓ 创建角色: %d个\n", len(characterIDs))

	// 5. 验证数据完整性
	if err := s.validateHarnessData(projectID); err != nil {
		return fmt.Errorf("数据验证失败: %w", err)
	}

	fmt.Println("\n✅ Harness数据创建完成！")
	fmt.Println("📖 项目名称: 就业指南：我在异世界当劝退专员")
	fmt.Println("🔐 登录账号: harness_writer / qingyu2024")

	return nil
}

// Clean 清空Harness数据
func (s *HarnessDataSeeder) Clean() error {
	ctx := context.Background()

	// 获取项目ID
	var project struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	err := s.db.Collection("projects").FindOne(ctx, bson.M{"title": story.EmploymentGuideProject.Title}).Decode(&project)
	if err != nil {
		fmt.Println("✓ Harness项目不存在，无需清理")
		return nil
	}

	// 按依赖顺序删除
	collections := []string{
		"character_relations",
		"characters",
		"projects",
	}

	for _, colName := range collections {
		filter := bson.M{"project_id": project.ID}
		if colName == "projects" {
			filter = bson.M{"_id": project.ID}
		}
		result, err := s.db.Collection(colName).DeleteMany(ctx, filter)
		if err != nil {
			return fmt.Errorf("清空 %s 失败: %w", colName, err)
		}
		if result.DeletedCount > 0 {
			fmt.Printf("✓ 清空 %s: %d条\n", colName, result.DeletedCount)
		}
	}

	fmt.Println("✓ Harness数据清理完成")
	return nil
}

// projectExists 检查项目是否存在
func (s *HarnessDataSeeder) projectExists(title string) bool {
	ctx := context.Background()
	count, _ := s.db.Collection("projects").CountDocuments(ctx, bson.M{"title": title})
	return count > 0
}

// getOrCreateHarnessUser 获取或创建Harness专用用户
// 独立于admin，避免被 seeder users --clean 误删
func (s *HarnessDataSeeder) getOrCreateHarnessUser() (string, error) {
	ctx := context.Background()
	collection := s.db.Collection("users")

	// 查找专用的 harness 用户
	harnessUsername := "harness_writer"
	var user struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	err := collection.FindOne(ctx, bson.M{"username": harnessUsername}).Decode(&user)
	if err == nil {
		fmt.Printf("✓ 使用已存在Harness用户: %s\n", harnessUsername)
		return user.ID.Hex(), nil
	}

	// 不存在则创建专用用户
	return s.createHarnessUser()
}

// createHarnessUser 创建Harness专用用户
func (s *HarnessDataSeeder) createHarnessUser() (string, error) {
	ctx := context.Background()
	collection := s.db.Collection("users")

	now := time.Now()
	id := primitive.NewObjectID()

	// 使用 bcrypt 哈希密码
	hashedPassword, err := utils.HashPassword("qingyu2024")
	if err != nil {
		return "", fmt.Errorf("密码哈希失败: %w", err)
	}

	user := bson.M{
		"_id":        id,
		"username":   "harness_writer",
		"email":      "harness@qingyu.com",
		"password":   hashedPassword,
		"roles":      []string{"author"},
		"status":     "active",
		"nickname":   "Harness作家",
		"avatar":     "/images/avatars/default.png",
		"bio":        "Harness测试专用账号",
		"created_at": now,
		"updated_at": now,
	}

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	fmt.Println("✓ 创建 Harness 账号: harness_writer")
	return id.Hex(), nil
}

// createProject 创建项目
func (s *HarnessDataSeeder) createProject(authorID string) (string, error) {
	ctx := context.Background()
	collection := s.db.Collection("projects")

	now := time.Now()
	id := primitive.NewObjectID()

	project := bson.M{
		"_id":          id,
		"author_id":    func() primitive.ObjectID { oid, _ := primitive.ObjectIDFromHex(authorID); return oid }(),
		"title":        story.EmploymentGuideProject.Title,
		"summary":      story.EmploymentGuideProject.Summary,
		"cover_url":    story.EmploymentGuideProject.CoverURL,
		"category":     story.EmploymentGuideProject.Category,
		"writing_type": story.EmploymentGuideProject.WritingType,
		"status":       story.EmploymentGuideProject.Status,
		"visibility":   story.EmploymentGuideProject.Visibility,
		"tags":         story.EmploymentGuideProject.Tags,
		"statistics": bson.M{
			"total_words":    0,
			"chapter_count":  0,
			"document_count": 0,
			"last_update_at": now,
		},
		"settings": bson.M{
			"auto_backup":     true,
			"backup_interval": 24,
		},
		"collaborators": []interface{}{},
		"created_at":    now,
		"updated_at":    now,
	}

	_, err := collection.InsertOne(ctx, project)
	if err != nil {
		return "", err
	}

	return id.Hex(), nil
}

// createCharacters 创建角色
func (s *HarnessDataSeeder) createCharacters(projectID string) (map[string]string, error) {
	ctx := context.Background()
	collection := s.db.Collection("characters")

	now := time.Now()
	projectOID, _ := primitive.ObjectIDFromHex(projectID)
	characterIDs := make(map[string]string)

	for _, char := range story.NewCharacters {
		charID := primitive.NewObjectID()
		character := bson.M{
			"_id":                charID,
			"project_id":         projectOID,
			"name":               char.Name,
			"alias":              char.Alias,
			"summary":            char.Summary,
			"traits":             char.Traits,
			"background":         char.Background,
			"avatar_url":         char.AvatarURL,
			"short_description":  char.Summary,
			"personality_prompt": char.PersonalityPrompt,
			"speech_pattern":     char.SpeechPattern,
			"current_state":      "",
			"created_at":         now,
			"updated_at":         now,
		}
		_, err := collection.InsertOne(ctx, character)
		if err != nil {
			return nil, err
		}
		characterIDs[char.ID] = charID.Hex()
	}

	return characterIDs, nil
}

// validateHarnessData 验证数据完整性
func (s *HarnessDataSeeder) validateHarnessData(projectID string) error {
	ctx := context.Background()

	// 检查各集合数量
	collections := []struct {
		name     string
		expected int64
	}{
		{"characters", 3},
	}

	fmt.Println("\n📊 数据完整性验证:")
	for _, col := range collections {
		count, err := s.db.Collection(col.name).CountDocuments(ctx, bson.M{"project_id": projectID})
		if err != nil {
			return err
		}
		status := "✓"
		if count != col.expected {
			status = "⚠"
		}
		fmt.Printf("  %s %s: %d/%d\n", status, col.name, count, col.expected)
	}

	return nil
}
