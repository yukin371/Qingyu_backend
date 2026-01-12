package notification

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/models/notification"
	repo "Qingyu_backend/repository/interfaces/notification"
)

// NotificationTemplateRepositoryImpl 通知模板仓储实现
type NotificationTemplateRepositoryImpl struct {
	db                 *mongo.Database
	templateCollection *mongo.Collection
}

// NewNotificationTemplateRepository 创建通知模板仓储实例
func NewNotificationTemplateRepository(db *mongo.Database) repo.NotificationTemplateRepository {
	return &NotificationTemplateRepositoryImpl{
		db:                 db,
		templateCollection: db.Collection("notification_templates"),
	}
}

// Create 创建通知模板
func (r *NotificationTemplateRepositoryImpl) Create(ctx context.Context, template *notification.NotificationTemplate) error {
	objectID, err := primitive.ObjectIDFromHex(template.ID)
	if err != nil {
		// 如果ID无效，生成新的
		objectID = primitive.NewObjectID()
		template.ID = objectID.Hex()
	}

	doc := bson.M{
		"_id":        objectID,
		"type":       template.Type,
		"action":     template.Action,
		"title":      template.Title,
		"content":    template.Content,
		"variables":  template.Variables,
		"data":       template.Data,
		"language":   template.Language,
		"is_active":  template.IsActive,
		"created_at": template.CreatedAt,
		"updated_at": template.UpdatedAt,
	}

	_, err = r.templateCollection.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("创建通知模板失败: %w", err)
	}

	return nil
}

// GetByID 根据ID获取通知模板
func (r *NotificationTemplateRepositoryImpl) GetByID(ctx context.Context, id string) (*notification.NotificationTemplate, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("无效的模板ID: %w", err)
	}

	var template notification.NotificationTemplate
	err = r.templateCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&template)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("查询通知模板失败: %w", err)
	}

	return &template, nil
}

// Update 更新通知模板
func (r *NotificationTemplateRepositoryImpl) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的模板ID: %w", err)
	}

	result, err := r.templateCollection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("更新通知模板失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("通知模板不存在")
	}

	return nil
}

// Delete 删除通知模板
func (r *NotificationTemplateRepositoryImpl) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的模板ID: %w", err)
	}

	result, err := r.templateCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("删除通知模板失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("通知模板不存在")
	}

	return nil
}

// Exists 检查通知模板是否存在
func (r *NotificationTemplateRepositoryImpl) Exists(ctx context.Context, id string) (bool, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, fmt.Errorf("无效的模板ID: %w", err)
	}

	count, err := r.templateCollection.CountDocuments(ctx, bson.M{"_id": objectID})
	if err != nil {
		return false, fmt.Errorf("检查通知模板存在性失败: %w", err)
	}

	return count > 0, nil
}

// List 获取通知模板列表
func (r *NotificationTemplateRepositoryImpl) List(ctx context.Context, filter *repo.TemplateFilter) ([]*notification.NotificationTemplate, error) {
	mongoFilter := r.buildFilter(filter)

	// 构建排序
	opts := options.Find()
	if filter.Limit > 0 {
		opts.SetLimit(int64(filter.Limit))
	}
	if filter.Offset > 0 {
		opts.SetSkip(int64(filter.Offset))
	}
	opts.SetSort(bson.M{"created_at": -1})

	cursor, err := r.templateCollection.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询通知模板列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var templates []*notification.NotificationTemplate
	if err = cursor.All(ctx, &templates); err != nil {
		return nil, fmt.Errorf("解析通知模板列表失败: %w", err)
	}

	return templates, nil
}

// GetByTypeAndAction 根据类型和操作获取通知模板列表
func (r *NotificationTemplateRepositoryImpl) GetByTypeAndAction(ctx context.Context, templateType notification.NotificationType, action string) ([]*notification.NotificationTemplate, error) {
	filter := bson.M{
		"type":   templateType,
		"action": action,
	}

	cursor, err := r.templateCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("查询通知模板失败: %w", err)
	}
	defer cursor.Close(ctx)

	var templates []*notification.NotificationTemplate
	if err = cursor.All(ctx, &templates); err != nil {
		return nil, fmt.Errorf("解析通知模板失败: %w", err)
	}

	return templates, nil
}

// GetActiveTemplate 获取活跃的通知模板
func (r *NotificationTemplateRepositoryImpl) GetActiveTemplate(ctx context.Context, templateType notification.NotificationType, action string, language string) (*notification.NotificationTemplate, error) {
	filter := bson.M{
		"type":      templateType,
		"action":    action,
		"language":  language,
		"is_active": true,
	}

	var template notification.NotificationTemplate
	err := r.templateCollection.FindOne(ctx, filter).Decode(&template)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("查询通知模板失败: %w", err)
	}

	return &template, nil
}

// buildFilter 构建MongoDB查询条件
func (r *NotificationTemplateRepositoryImpl) buildFilter(filter *repo.TemplateFilter) bson.M {
	mongoFilter := bson.M{}

	if filter.Type != nil {
		mongoFilter["type"] = *filter.Type
	}

	if filter.Action != nil {
		mongoFilter["action"] = *filter.Action
	}

	if filter.Language != nil {
		mongoFilter["language"] = *filter.Language
	}

	if filter.IsActive != nil {
		mongoFilter["is_active"] = *filter.IsActive
	}

	return mongoFilter
}
