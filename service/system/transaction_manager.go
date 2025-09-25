package system

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"Qingyu_backend/models/system"
)

// TransactionManager 事务管理器
type TransactionManager struct {
	client *mongo.Client
	db     *mongo.Database
}

// NewTransactionManager 创建事务管理器实例
func NewTransactionManager(client *mongo.Client, db *mongo.Database) *TransactionManager {
	return &TransactionManager{
		client: client,
		db:     db,
	}
}

// TransactionOperation 事务操作接口
type TransactionOperation interface {
	Execute(ctx mongo.SessionContext, db *mongo.Database) error
	Rollback(ctx mongo.SessionContext, db *mongo.Database) error
	GetDescription() string
}

// ExecuteTransaction 执行事务操作
func (tm *TransactionManager) ExecuteTransaction(ctx context.Context, operations []TransactionOperation) error {
	session, err := tm.client.StartSession()
	if err != nil {
		return fmt.Errorf("启动事务会话失败: %w", err)
	}
	defer session.EndSession(ctx)

	return mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return fmt.Errorf("启动事务失败: %w", err)
		}

		// 执行所有操作
		for i, op := range operations {
			if err := op.Execute(sc, tm.db); err != nil {
				// 回滚事务
				if abortErr := session.AbortTransaction(sc); abortErr != nil {
					log.Printf("回滚事务失败: %v", abortErr)
				}
				return fmt.Errorf("执行操作 %d (%s) 失败: %w", i+1, op.GetDescription(), err)
			}
		}

		// 提交事务
		if err := session.CommitTransaction(sc); err != nil {
			return fmt.Errorf("提交事务失败: %w", err)
		}

		return nil
	})
}

// UserRegistrationTransaction 用户注册事务操作
type UserRegistrationTransaction struct {
	User       *system.User
	UserRole   *UserRole
	UserConfig *UserConfig
}

// UserRole 用户角色关联
type UserRole struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	RoleID    primitive.ObjectID `bson:"role_id" json:"role_id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// UserConfig 用户配置
type UserConfig struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Theme     string             `bson:"theme" json:"theme"`
	Language  string             `bson:"language" json:"language"`
	Settings  map[string]interface{} `bson:"settings" json:"settings"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

func (urt *UserRegistrationTransaction) Execute(ctx mongo.SessionContext, db *mongo.Database) error {
	// 1. 创建用户
	if urt.User.ID == "" {
		urt.User.ID = primitive.NewObjectID().Hex()
	}
	urt.User.TouchForCreate()

	result, err := db.Collection("users").InsertOne(ctx, urt.User)
	if err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}

	userObjectID := result.InsertedID.(primitive.ObjectID)

	// 2. 分配默认角色
	if urt.UserRole != nil {
		urt.UserRole.UserID = userObjectID
		if urt.UserRole.RoleID.IsZero() {
			// 获取默认用户角色ID
			defaultRoleID, err := getDefaultUserRoleID(ctx, db)
			if err != nil {
				return fmt.Errorf("获取默认角色失败: %w", err)
			}
			urt.UserRole.RoleID = defaultRoleID
		}
		urt.UserRole.CreatedAt = time.Now()
		urt.UserRole.UpdatedAt = time.Now()

		_, err = db.Collection("user_roles").InsertOne(ctx, urt.UserRole)
		if err != nil {
			return fmt.Errorf("分配用户角色失败: %w", err)
		}
	}

	// 3. 初始化用户配置
	if urt.UserConfig != nil {
		urt.UserConfig.UserID = userObjectID
		if urt.UserConfig.Theme == "" {
			urt.UserConfig.Theme = "default"
		}
		if urt.UserConfig.Language == "" {
			urt.UserConfig.Language = "zh-CN"
		}
		if urt.UserConfig.Settings == nil {
			urt.UserConfig.Settings = make(map[string]interface{})
		}
		urt.UserConfig.CreatedAt = time.Now()
		urt.UserConfig.UpdatedAt = time.Now()

		_, err = db.Collection("user_configs").InsertOne(ctx, urt.UserConfig)
		if err != nil {
			return fmt.Errorf("初始化用户配置失败: %w", err)
		}
	}

	return nil
}

func (urt *UserRegistrationTransaction) Rollback(ctx mongo.SessionContext, db *mongo.Database) error {
	// MongoDB事务会自动回滚，这里可以添加额外的清理逻辑
	log.Printf("用户注册事务回滚: %s", urt.User.Username)
	return nil
}

func (urt *UserRegistrationTransaction) GetDescription() string {
	return fmt.Sprintf("用户注册事务: %s", urt.User.Username)
}

// getDefaultUserRoleID 获取默认用户角色ID
func getDefaultUserRoleID(ctx mongo.SessionContext, db *mongo.Database) (primitive.ObjectID, error) {
	var role struct {
		ID primitive.ObjectID `bson:"_id"`
	}

	err := db.Collection("roles").FindOne(ctx, bson.M{
		"name":        "user",
		"is_default":  true,
	}).Decode(&role)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// 如果没有找到默认角色，创建一个
			return createDefaultUserRole(ctx, db)
		}
		return primitive.NilObjectID, err
	}

	return role.ID, nil
}

// createDefaultUserRole 创建默认用户角色
func createDefaultUserRole(ctx mongo.SessionContext, db *mongo.Database) (primitive.ObjectID, error) {
	role := bson.M{
		"name":        "user",
		"description": "默认用户角色",
		"is_default":  true,
		"permissions": []string{"read_own_profile", "update_own_profile"},
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	}

	result, err := db.Collection("roles").InsertOne(ctx, role)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("创建默认角色失败: %w", err)
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

// CascadeManager 级联操作管理器
type CascadeManager struct {
	db *mongo.Database
	tm *TransactionManager
}

// NewCascadeManager 创建级联操作管理器
func NewCascadeManager(db *mongo.Database, tm *TransactionManager) *CascadeManager {
	return &CascadeManager{
		db: db,
		tm: tm,
	}
}

// UserDeletionTransaction 用户删除级联事务
type UserDeletionTransaction struct {
	UserID     primitive.ObjectID
	SoftDelete bool // 是否软删除
}

func (udt *UserDeletionTransaction) Execute(ctx mongo.SessionContext, db *mongo.Database) error {
	if udt.SoftDelete {
		return udt.executeSoftDelete(ctx, db)
	}
	return udt.executeHardDelete(ctx, db)
}

func (udt *UserDeletionTransaction) executeSoftDelete(ctx mongo.SessionContext, db *mongo.Database) error {
	now := time.Now()

	// 1. 软删除用户
	_, err := db.Collection("users").UpdateOne(ctx,
		bson.M{"_id": udt.UserID},
		bson.M{"$set": bson.M{
			"status":     "deleted",
			"deleted_at": now,
			"updated_at": now,
		}},
	)
	if err != nil {
		return fmt.Errorf("软删除用户失败: %w", err)
	}

	// 2. 归档用户项目
	_, err = db.Collection("projects").UpdateMany(ctx,
		bson.M{"creator_id": udt.UserID},
		bson.M{"$set": bson.M{
			"status":          "archived",
			"archived_reason": "creator_deleted",
			"archived_at":     now,
			"updated_at":      now,
		}},
	)
	if err != nil {
		return fmt.Errorf("归档用户项目失败: %w", err)
	}

	// 3. 禁用用户会话
	_, err = db.Collection("user_sessions").UpdateMany(ctx,
		bson.M{"user_id": udt.UserID, "status": "active"},
		bson.M{"$set": bson.M{
			"status":     "revoked",
			"revoked_at": now,
			"updated_at": now,
		}},
	)
	if err != nil {
		return fmt.Errorf("禁用用户会话失败: %w", err)
	}

	return nil
}

func (udt *UserDeletionTransaction) executeHardDelete(ctx mongo.SessionContext, db *mongo.Database) error {
	// 1. 删除用户角色关联
	_, err := db.Collection("user_roles").DeleteMany(ctx,
		bson.M{"user_id": udt.UserID},
	)
	if err != nil {
		return fmt.Errorf("删除用户角色关联失败: %w", err)
	}

	// 2. 删除用户配置
	_, err = db.Collection("user_configs").DeleteMany(ctx,
		bson.M{"user_id": udt.UserID},
	)
	if err != nil {
		return fmt.Errorf("删除用户配置失败: %w", err)
	}

	// 3. 删除用户会话
	_, err = db.Collection("user_sessions").DeleteMany(ctx,
		bson.M{"user_id": udt.UserID},
	)
	if err != nil {
		return fmt.Errorf("删除用户会话失败: %w", err)
	}

	// 4. 处理用户项目（转移或删除）
	_, err = db.Collection("projects").UpdateMany(ctx,
		bson.M{"creator_id": udt.UserID},
		bson.M{"$set": bson.M{
			"creator_id": primitive.NilObjectID, // 设置为空，表示已删除用户
			"status":     "orphaned",
			"updated_at": time.Now(),
		}},
	)
	if err != nil {
		return fmt.Errorf("处理用户项目失败: %w", err)
	}

	// 5. 删除用户
	_, err = db.Collection("users").DeleteOne(ctx,
		bson.M{"_id": udt.UserID},
	)
	if err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}

	return nil
}

func (udt *UserDeletionTransaction) Rollback(ctx mongo.SessionContext, db *mongo.Database) error {
	log.Printf("用户删除事务回滚: %s", udt.UserID.Hex())
	return nil
}

func (udt *UserDeletionTransaction) GetDescription() string {
	deleteType := "软删除"
	if !udt.SoftDelete {
		deleteType = "硬删除"
	}
	return fmt.Sprintf("用户%s事务: %s", deleteType, udt.UserID.Hex())
}

// DeleteUser 删除用户（级联操作）
func (cm *CascadeManager) DeleteUser(ctx context.Context, userID primitive.ObjectID, softDelete bool) error {
	transaction := &UserDeletionTransaction{
		UserID:     userID,
		SoftDelete: softDelete,
	}

	return cm.tm.ExecuteTransaction(ctx, []TransactionOperation{transaction})
}

// UserUpdateTransaction 用户更新事务
type UserUpdateTransaction struct {
	UserID  primitive.ObjectID
	Updates bson.M
}

func (uut *UserUpdateTransaction) Execute(ctx mongo.SessionContext, db *mongo.Database) error {
	// 添加更新时间
	uut.Updates["updated_at"] = time.Now()

	_, err := db.Collection("users").UpdateOne(ctx,
		bson.M{"_id": uut.UserID},
		bson.M{"$set": uut.Updates},
	)
	if err != nil {
		return fmt.Errorf("更新用户失败: %w", err)
	}

	return nil
}

func (uut *UserUpdateTransaction) Rollback(ctx mongo.SessionContext, db *mongo.Database) error {
	log.Printf("用户更新事务回滚: %s", uut.UserID.Hex())
	return nil
}

func (uut *UserUpdateTransaction) GetDescription() string {
	return fmt.Sprintf("用户更新事务: %s", uut.UserID.Hex())
}

// SagaManager Saga模式事务管理器
type SagaManager struct {
	db *mongo.Database
}

// NewSagaManager 创建Saga管理器
func NewSagaManager(db *mongo.Database) *SagaManager {
	return &SagaManager{db: db}
}

// SagaStep Saga步骤
type SagaStep struct {
	Name       string
	Execute    func(ctx context.Context) error
	Compensate func(ctx context.Context) error
}

// ExecuteSaga 执行Saga事务
func (sm *SagaManager) ExecuteSaga(ctx context.Context, steps []SagaStep) error {
	executedSteps := make([]SagaStep, 0)

	for _, step := range steps {
		log.Printf("执行Saga步骤: %s", step.Name)
		
		if err := step.Execute(ctx); err != nil {
			log.Printf("Saga步骤 %s 执行失败: %v", step.Name, err)
			
			// 执行补偿操作
			for i := len(executedSteps) - 1; i >= 0; i-- {
				compensateStep := executedSteps[i]
				log.Printf("执行补偿操作: %s", compensateStep.Name)
				
				if compensateErr := compensateStep.Compensate(ctx); compensateErr != nil {
					log.Printf("补偿操作 %s 失败: %v", compensateStep.Name, compensateErr)
				}
			}
			
			return fmt.Errorf("Saga事务失败，步骤 %s: %w", step.Name, err)
		}
		
		executedSteps = append(executedSteps, step)
	}

	log.Printf("Saga事务成功完成，共执行 %d 个步骤", len(steps))
	return nil
}

// ReferenceIntegrityManager 引用完整性管理器
type ReferenceIntegrityManager struct {
	db *mongo.Database
}

// NewReferenceIntegrityManager 创建引用完整性管理器
func NewReferenceIntegrityManager(db *mongo.Database) *ReferenceIntegrityManager {
	return &ReferenceIntegrityManager{db: db}
}

// ValidateReferences 验证引用完整性
func (rim *ReferenceIntegrityManager) ValidateReferences(ctx context.Context, doc interface{}, collection string) error {
	switch collection {
	case "projects":
		return rim.validateProjectReferences(ctx, doc)
	case "user_roles":
		return rim.validateUserRoleReferences(ctx, doc)
	default:
		return nil
	}
}

// validateProjectReferences 验证项目引用
func (rim *ReferenceIntegrityManager) validateProjectReferences(ctx context.Context, doc interface{}) error {
	project, ok := doc.(map[string]interface{})
	if !ok {
		return fmt.Errorf("无效的项目文档格式")
	}

	creatorID, ok := project["creator_id"].(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("项目创建者ID格式错误")
	}

	// 验证创建者存在且未被删除
	count, err := rim.db.Collection("users").CountDocuments(ctx, bson.M{
		"_id":    creatorID,
		"status": bson.M{"$ne": "deleted"},
	})
	if err != nil {
		return fmt.Errorf("验证项目创建者失败: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("项目创建者不存在或已被删除")
	}

	return nil
}

// validateUserRoleReferences 验证用户角色引用
func (rim *ReferenceIntegrityManager) validateUserRoleReferences(ctx context.Context, doc interface{}) error {
	userRole, ok := doc.(map[string]interface{})
	if !ok {
		return fmt.Errorf("无效的用户角色文档格式")
	}

	userID, ok := userRole["user_id"].(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("用户ID格式错误")
	}

	roleID, ok := userRole["role_id"].(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("角色ID格式错误")
	}

	// 验证用户存在
	userCount, err := rim.db.Collection("users").CountDocuments(ctx, bson.M{
		"_id":    userID,
		"status": bson.M{"$ne": "deleted"},
	})
	if err != nil {
		return fmt.Errorf("验证用户存在性失败: %w", err)
	}

	if userCount == 0 {
		return fmt.Errorf("用户不存在或已被删除")
	}

	// 验证角色存在
	roleCount, err := rim.db.Collection("roles").CountDocuments(ctx, bson.M{
		"_id": roleID,
	})
	if err != nil {
		return fmt.Errorf("验证角色存在性失败: %w", err)
	}

	if roleCount == 0 {
		return fmt.Errorf("角色不存在")
	}

	return nil
}