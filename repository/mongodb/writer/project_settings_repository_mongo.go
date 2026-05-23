// --- 项目设置仓库MongoDB实现 ---
package writer

import (
	"context"
	"time"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/pkg/errors"
	writerRepo "Qingyu_backend/repository/interfaces/writer"
	"Qingyu_backend/repository/mongodb/base"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ProjectSettingsRepositoryMongo ProjectSettings Repository的MongoDB实现
type ProjectSettingsRepositoryMongo struct {
	*base.BaseMongoRepository // 嵌入基类，继承ID转换和通用CRUD方法
	db                        *mongo.Database
}

// NewProjectSettingsRepository 创建ProjectSettingsRepository实例
func NewProjectSettingsRepository(db *mongo.Database) writerRepo.ProjectSettingsRepository {
	return &ProjectSettingsRepositoryMongo{
		BaseMongoRepository: base.NewBaseMongoRepository(db, "projects"),
		db:                  db,
	}
}

// Create 创建项目设置
// 注意：ProjectSettings 作为 Project 文档的一部分存储，此方法用于创建包含设置的完整项目
func (r *ProjectSettingsRepositoryMongo) Create(ctx context.Context, settings *writer.ProjectSettings) error {
	// ProjectSettings 应该作为 Project 的一部分创建
	// 此方法主要用于初始化项目时设置默认值
	now := time.Now()

	// 创建一个完整的项目文档（假设 Project 结构体存在于 models/writer/project.go 中）
	project := bson.M{
		"_id":         primitive.NewObjectID(),
		"settings":    settings,
		"created_at":  now,
		"updated_at":  now,
	}

	_, err := r.GetCollection().InsertOne(ctx, project)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "create project settings failed", err)
	}

	return nil
}

// FindByProjectID 根据项目ID查找设置
func (r *ProjectSettingsRepositoryMongo) FindByProjectID(ctx context.Context, projectID string) (*writer.ProjectSettings, error) {
	projectObjectID, err := r.ParseID(projectID)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid project ID", err)
	}

	// 只查询 settings 字段
	filter := bson.M{"_id": projectObjectID}
	projection := bson.M{"settings": 1}

	opts := options.FindOne().SetProjection(projection)
	var result struct {
		Settings writer.ProjectSettings `bson:"settings"`
	}

	err = r.GetCollection().FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// 项目不存在时，返回默认设置
			defaultSettings := writer.NewProjectSettings(projectObjectID)
			return &defaultSettings, nil
		}
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find project settings failed", err)
	}

	return &result.Settings, nil
}

// Update 更新项目设置
func (r *ProjectSettingsRepositoryMongo) Update(ctx context.Context, projectID string, settings *writer.ProjectSettings) error {
	projectObjectID, err := r.ParseID(projectID)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid project ID", err)
	}

	filter := bson.M{"_id": projectObjectID}
	update := bson.M{
		"$set": bson.M{
			"settings":             settings,
			"updated_at":           time.Now(),
		},
	}

	result, err := r.GetCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "update project settings failed", err)
	}

	if result.MatchedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "project not found", nil)
	}

	return nil
}

// AddCharacterRole 添加自定义角色类型
func (r *ProjectSettingsRepositoryMongo) AddCharacterRole(ctx context.Context, projectID string, role *writer.CharacterRole) error {
	projectObjectID, err := r.ParseID(projectID)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid project ID", err)
	}

	filter := bson.M{"_id": projectObjectID}
	update := bson.M{
		"$push": bson.M{
			"settings.character_roles": role,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	result, err := r.GetCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "add character role failed", err)
	}

	if result.MatchedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "project not found", nil)
	}

	return nil
}

// UpdateCharacterRole 更新角色类型
func (r *ProjectSettingsRepositoryMongo) UpdateCharacterRole(ctx context.Context, projectID, roleID string, role *writer.CharacterRole) error {
	projectObjectID, err := r.ParseID(projectID)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid project ID", err)
	}

	roleObjectID, err := r.ParseID(roleID)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid role ID", err)
	}

	filter := bson.M{
		"_id":                       projectObjectID,
		"settings.character_roles.id": roleObjectID,
	}

	update := bson.M{
		"$set": bson.M{
			"settings.character_roles.$": role,
			"updated_at":                 time.Now(),
		},
	}

	result, err := r.GetCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "update character role failed", err)
	}

	if result.MatchedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "project or character role not found", nil)
	}

	return nil
}

// DeleteCharacterRole 删除角色类型
func (r *ProjectSettingsRepositoryMongo) DeleteCharacterRole(ctx context.Context, projectID, roleID string) error {
	projectObjectID, err := r.ParseID(projectID)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid project ID", err)
	}

	roleObjectID, err := r.ParseID(roleID)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid role ID", err)
	}

	filter := bson.M{"_id": projectObjectID}
	update := bson.M{
		"$pull": bson.M{
			"settings.character_roles": bson.M{"id": roleObjectID},
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	result, err := r.GetCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "delete character role failed", err)
	}

	if result.MatchedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "project not found", nil)
	}

	return nil
}

// GetDefaultRoles 获取项目的角色类型列表（包含默认和自定义）
func (r *ProjectSettingsRepositoryMongo) GetDefaultRoles(ctx context.Context, projectID string) ([]writer.CharacterRole, error) {
	settings, err := r.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if settings == nil {
		// 返回默认角色类型
		projectObjectID, err := r.ParseID(projectID)
		if err != nil {
			return nil, errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid project ID", err)
		}
		return writer.GetDefaultCharacterRoles(projectObjectID), nil
	}

	return settings.CharacterRoles, nil
}
