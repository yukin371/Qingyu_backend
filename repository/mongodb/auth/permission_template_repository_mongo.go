package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	authModel "Qingyu_backend/models/auth"
	repoAuth "Qingyu_backend/repository/interfaces/auth"
)

// PermissionTemplateRepositoryMongo 权限模板MongoDB仓储实现
type PermissionTemplateRepositoryMongo struct {
	client     *mongo.Client
	database   string
	collection string
	redis      *redis.Client
}

// NewPermissionTemplateRepositoryMongo 创建权限模板仓储
func NewPermissionTemplateRepositoryMongo(client *mongo.Client, database string, redisClient *redis.Client) repoAuth.PermissionTemplateRepository {
	return &PermissionTemplateRepositoryMongo{
		client:     client,
		database:   database,
		collection: "permission_templates",
		redis:      redisClient,
	}
}

// getCollection 获取集合
func (r *PermissionTemplateRepositoryMongo) getCollection() *mongo.Collection {
	return r.client.Database(r.database).Collection(r.collection)
}

// ============ 模板管理实现 ============

// CreateTemplate 创建模板
func (r *PermissionTemplateRepositoryMongo) CreateTemplate(ctx context.Context, template *authModel.PermissionTemplate) error {
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()

	// 生成ID
	if template.ID == "" {
		template.ID = primitive.NewObjectID().Hex()
	}

	// 检查代码唯一性
	existing, _ := r.GetTemplateByCode(ctx, template.Code)
	if existing != nil {
		return authModel.ErrTemplateCodeExists
	}

	_, err := r.getCollection().InsertOne(ctx, template)
	if err != nil {
		return fmt.Errorf("插入模板失败: %w", err)
	}

	// 清除缓存
	r.clearTemplateCache(ctx, template.ID)
	r.clearTemplateListCache(ctx)

	return nil
}

// GetTemplateByID 根据ID获取模板
func (r *PermissionTemplateRepositoryMongo) GetTemplateByID(ctx context.Context, templateID string) (*authModel.PermissionTemplate, error) {
	// 尝试从缓存获取
	if r.redis != nil {
		cacheKey := r.getTemplateCacheKey(templateID)
		cached, err := r.redis.Get(ctx, cacheKey).Result()
		if err == nil && cached != "" {
			var template authModel.PermissionTemplate
			if err := bson.Unmarshal([]byte(cached), &template); err == nil {
				return &template, nil
			}
		}
	}

	// 从数据库获取
	objectID, err := primitive.ObjectIDFromHex(templateID)
	if err != nil {
		// 如果不是有效的ObjectID，尝试按ID字符串查询
		filter := bson.M{"_id": templateID}
		return r.findOneByFilter(ctx, filter)
	}

	filter := bson.M{"_id": objectID}
	template, err := r.findOneByFilter(ctx, filter)
	if err != nil {
		return nil, authModel.ErrTemplateNotFound
	}

	// 存入缓存
	if r.redis != nil {
		cacheKey := r.getTemplateCacheKey(templateID)
		data, _ := bson.Marshal(template)
		_ = r.redis.Set(ctx, cacheKey, data, 10*time.Minute)
	}

	return template, nil
}

// GetTemplateByCode 根据代码获取模板
func (r *PermissionTemplateRepositoryMongo) GetTemplateByCode(ctx context.Context, code string) (*authModel.PermissionTemplate, error) {
	filter := bson.M{"code": code}
	return r.findOneByFilter(ctx, filter)
}

// UpdateTemplate 更新模板
func (r *PermissionTemplateRepositoryMongo) UpdateTemplate(ctx context.Context, templateID string, updates map[string]interface{}) error {
	// 获取模板
	template, err := r.GetTemplateByID(ctx, templateID)
	if err != nil {
		return err
	}

	// 检查是否是系统模板
	if template.IsSystem {
		return authModel.ErrTemplateIsSystem
	}

	// 添加更新时间
	updates["updated_at"] = time.Now()

	// 执行更新
	objectID, err := primitive.ObjectIDFromHex(templateID)
	if err != nil {
		filter := bson.M{"_id": templateID}
		_, err = r.getCollection().UpdateOne(ctx, filter, bson.M{"$set": updates})
	} else {
		filter := bson.M{"_id": objectID}
		_, err = r.getCollection().UpdateOne(ctx, filter, bson.M{"$set": updates})
	}

	if err != nil {
		return fmt.Errorf("更新模板失败: %w", err)
	}

	// 清除缓存
	r.clearTemplateCache(ctx, templateID)
	r.clearTemplateListCache(ctx)

	return nil
}

// DeleteTemplate 删除模板
func (r *PermissionTemplateRepositoryMongo) DeleteTemplate(ctx context.Context, templateID string) error {
	// 获取模板
	template, err := r.GetTemplateByID(ctx, templateID)
	if err != nil {
		return err
	}

	// 检查是否是系统模板
	if template.IsSystem {
		return authModel.ErrTemplateIsSystem
	}

	// 执行删除
	objectID, err := primitive.ObjectIDFromHex(templateID)
	var result *mongo.DeleteResult
	if err != nil {
		filter := bson.M{"_id": templateID}
		result, err = r.getCollection().DeleteOne(ctx, filter)
	} else {
		filter := bson.M{"_id": objectID}
		result, err = r.getCollection().DeleteOne(ctx, filter)
	}

	if err != nil {
		return fmt.Errorf("删除模板失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return authModel.ErrTemplateNotFound
	}

	// 清除缓存
	r.clearTemplateCache(ctx, templateID)
	r.clearTemplateListCache(ctx)

	return nil
}

// ListTemplates 列出所有模板
func (r *PermissionTemplateRepositoryMongo) ListTemplates(ctx context.Context) ([]*authModel.PermissionTemplate, error) {
	// 尝试从缓存获取
	if r.redis != nil {
		cacheKey := r.getTemplateListCacheKey("")
		cached, err := r.redis.Get(ctx, cacheKey).Result()
		if err == nil && cached != "" {
			var templates []*authModel.PermissionTemplate
			if err := bson.Unmarshal([]byte(cached), &templates); err == nil {
				return templates, nil
			}
		}
	}

	// 从数据库获取
	filter := bson.M{}
	opts := options.Find().SetSort(bson.M{"category": 1, "created_at": -1})

	cursor, err := r.getCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询模板列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var templates []*authModel.PermissionTemplate
	if err := cursor.All(ctx, &templates); err != nil {
		return nil, fmt.Errorf("解析模板列表失败: %w", err)
	}

	// 存入缓存
	if r.redis != nil && len(templates) > 0 {
		cacheKey := r.getTemplateListCacheKey("")
		data, _ := bson.Marshal(templates)
		_ = r.redis.Set(ctx, cacheKey, data, 5*time.Minute)
	}

	return templates, nil
}

// ListTemplatesByCategory 根据分类列出模板
func (r *PermissionTemplateRepositoryMongo) ListTemplatesByCategory(ctx context.Context, category string) ([]*authModel.PermissionTemplate, error) {
	// 尝试从缓存获取
	if r.redis != nil {
		cacheKey := r.getTemplateListCacheKey(category)
		cached, err := r.redis.Get(ctx, cacheKey).Result()
		if err == nil && cached != "" {
			var templates []*authModel.PermissionTemplate
			if err := bson.Unmarshal([]byte(cached), &templates); err == nil {
				return templates, nil
			}
		}
	}

	// 从数据库获取
	filter := bson.M{"category": category}
	opts := options.Find().SetSort(bson.M{"created_at": -1})

	cursor, err := r.getCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询模板列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var templates []*authModel.PermissionTemplate
	if err := cursor.All(ctx, &templates); err != nil {
		return nil, fmt.Errorf("解析模板列表失败: %w", err)
	}

	// 存入缓存
	if r.redis != nil && len(templates) > 0 {
		cacheKey := r.getTemplateListCacheKey(category)
		data, _ := bson.Marshal(templates)
		_ = r.redis.Set(ctx, cacheKey, data, 5*time.Minute)
	}

	return templates, nil
}

// ============ 模板应用实现 ============

// ApplyTemplateToRole 将模板应用到角色
func (r *PermissionTemplateRepositoryMongo) ApplyTemplateToRole(ctx context.Context, templateID, roleID string) error {
	// 1. 获取模板
	template, err := r.GetTemplateByID(ctx, templateID)
	if err != nil {
		return fmt.Errorf("获取模板失败: %w", err)
	}

	// 2. 更新角色权限
	// 注意：这里需要访问角色集合，所以需要使用数据库连接
	roleCollection := r.client.Database(r.database).Collection("roles")

	objectID, err := primitive.ObjectIDFromHex(roleID)
	var filter primitive.M
	if err != nil {
		filter = bson.M{"_id": roleID}
	} else {
		filter = bson.M{"_id": objectID}
	}

	// 替换角色的权限列表
	update := bson.M{
		"$set": bson.M{
			"permissions": template.Permissions,
			"updated_at":  time.Now(),
		},
	}

	result, err := roleCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新角色权限失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return errors.New("角色不存在")
	}

	return nil
}

// GetSystemTemplates 获取所有系统模板
func (r *PermissionTemplateRepositoryMongo) GetSystemTemplates(ctx context.Context) ([]*authModel.PermissionTemplate, error) {
	filter := bson.M{"is_system": true}
	cursor, err := r.getCollection().Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("查询系统模板失败: %w", err)
	}
	defer cursor.Close(ctx)

	var templates []*authModel.PermissionTemplate
	if err := cursor.All(ctx, &templates); err != nil {
		return nil, fmt.Errorf("解析系统模板失败: %w", err)
	}

	return templates, nil
}

// InitializeSystemTemplates 初始化系统预设模板
func (r *PermissionTemplateRepositoryMongo) InitializeSystemTemplates(ctx context.Context) error {
	// 定义系统预设模板
	systemTemplates := []*authModel.PermissionTemplate{
		{
			Name:        "读者模板",
			Code:        authModel.TemplateReader,
			Description: "读者角色默认权限模板",
			Permissions: []string{
				"book.read",
				"document.read",
				"comment.read",
			},
			IsSystem: true,
			Category: authModel.CategoryReader,
		},
		{
			Name:        "作者模板",
			Code:        authModel.TemplateAuthor,
			Description: "作者角色默认权限模板",
			Permissions: []string{
				"book.read",
				"book.write",
				"book.review",
				"document.read",
				"document.write",
				"comment.read",
				"comment.write",
				"wallet.read",
			},
			IsSystem: true,
			Category: authModel.CategoryAuthor,
		},
		{
			Name:        "管理员模板",
			Code:        authModel.TemplateAdmin,
			Description: "管理员角色默认权限模板",
			Permissions: []string{
				"user.read",
				"user.write",
				"user.delete",
				"book.read",
				"book.write",
				"book.delete",
				"book.review",
				"document.read",
				"document.write",
				"document.delete",
				"document.publish",
				"comment.read",
				"comment.write",
				"comment.delete",
				"wallet.read",
				"wallet.recharge",
				"wallet.withdraw",
				"wallet.review",
				"admin.access",
				"admin.review",
				"admin.manage",
			},
			IsSystem: true,
			Category: authModel.CategoryAdmin,
		},
	}

	// 创建或更新系统模板
	for _, template := range systemTemplates {
		existing, err := r.GetTemplateByCode(ctx, template.Code)
		if err != nil {
			// 不存在，创建新模板
			if errors.Is(err, authModel.ErrTemplateNotFound) {
				template.CreatedAt = time.Now()
				template.UpdatedAt = time.Now()
				template.ID = primitive.NewObjectID().Hex()

				_, err := r.getCollection().InsertOne(ctx, template)
				if err != nil {
					return fmt.Errorf("创建系统模板 %s 失败: %w", template.Code, err)
				}
			} else {
				return fmt.Errorf("检查系统模板 %s 失败: %w", template.Code, err)
			}
		} else {
			// 已存在，更新（保持系统模板的权限最新）
			objectID, err := primitive.ObjectIDFromHex(existing.ID)
			if err != nil {
				continue
			}

			filter := bson.M{"_id": objectID}
			update := bson.M{
				"$set": bson.M{
					"name":        template.Name,
					"description": template.Description,
					"permissions": template.Permissions,
					"category":    template.Category,
					"updated_at":  time.Now(),
				},
			}

			_, err = r.getCollection().UpdateOne(ctx, filter, update)
			if err != nil {
				return fmt.Errorf("更新系统模板 %s 失败: %w", template.Code, err)
			}
		}
	}

	// 清除缓存
	r.clearTemplateListCache(ctx)

	return nil
}

// Health 健康检查
func (r *PermissionTemplateRepositoryMongo) Health(ctx context.Context) error {
	return r.client.Ping(ctx, nil)
}

// ============ 辅助函数 ============

// findOneByFilter 根据条件查询单个模板
func (r *PermissionTemplateRepositoryMongo) findOneByFilter(ctx context.Context, filter bson.M) (*authModel.PermissionTemplate, error) {
	var template authModel.PermissionTemplate
	err := r.getCollection().FindOne(ctx, filter).Decode(&template)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, authModel.ErrTemplateNotFound
		}
		return nil, fmt.Errorf("查询模板失败: %w", err)
	}
	return &template, nil
}

// getTemplateCacheKey 获取模板缓存Key
func (r *PermissionTemplateRepositoryMongo) getTemplateCacheKey(templateID string) string {
	return fmt.Sprintf("permission_template:%s", templateID)
}

// getTemplateListCacheKey 获取模板列表缓存Key
func (r *PermissionTemplateRepositoryMongo) getTemplateListCacheKey(category string) string {
	if category != "" {
		return fmt.Sprintf("permission_templates:category:%s", category)
	}
	return "permission_templates:all"
}

// clearTemplateCache 清除模板缓存
func (r *PermissionTemplateRepositoryMongo) clearTemplateCache(ctx context.Context, templateID string) {
	if r.redis != nil {
		cacheKey := r.getTemplateCacheKey(templateID)
		_ = r.redis.Del(ctx, cacheKey)
	}
}

// clearTemplateListCache 清除模板列表缓存
func (r *PermissionTemplateRepositoryMongo) clearTemplateListCache(ctx context.Context) {
	if r.redis != nil {
		_ = r.redis.Del(ctx, "permission_templates:all")
		_ = r.redis.Del(ctx, "permission_templates:category:"+authModel.CategoryReader)
		_ = r.redis.Del(ctx, "permission_templates:category:"+authModel.CategoryAuthor)
		_ = r.redis.Del(ctx, "permission_templates:category:"+authModel.CategoryAdmin)
		_ = r.redis.Del(ctx, "permission_templates:category:"+authModel.CategoryCustom)
	}
}
