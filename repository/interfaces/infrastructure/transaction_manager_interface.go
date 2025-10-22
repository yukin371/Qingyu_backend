package infrastructure

import (
	"context"
	"fmt"
	"log"
	"time"

	"Qingyu_backend/models/ai"
	"Qingyu_backend/models/document"
	usersModel "Qingyu_backend/models/users"

	"go.mongodb.org/mongo-driver/mongo"
)

// TransactionManager 事务管理器接口
type TransactionManager interface {
	// 执行事务
	ExecuteTransaction(ctx context.Context, fn TransactionFunc) error

	// 执行Saga事务
	ExecuteSaga(ctx context.Context, saga *Saga) error

	// 获取事务上下文
	GetTransactionContext(ctx context.Context) (TransactionContext, error)

	// 健康检查
	Health(ctx context.Context) error
}

// TransactionFunc 事务函数类型
type TransactionFunc func(txCtx TransactionContext) error

// TransactionContext 事务上下文接口
type TransactionContext interface {
	context.Context

	// 获取MongoDB会话
	GetSession() mongo.Session

	// 获取Repository工厂
	GetRepositoryFactory() TransactionRepositoryFactory

	// 事务状态
	IsInTransaction() bool
	GetTransactionID() string

	// 回滚标记
	MarkForRollback(reason string)
	IsMarkedForRollback() bool
	GetRollbackReason() string
}

// TransactionRepositoryFactory 事务Repository工厂接口
type TransactionRepositoryFactory interface {
	// 创建用户Repository
	CreateUserRepository(ctx TransactionContext) TransactionUserRepository

	// 创建AI Repository
	CreateAIRepository(ctx TransactionContext) TransactionAIRepository

	// 创建文档Repository
	CreateDocumentRepository(ctx TransactionContext) TransactionDocumentRepository

	// 创建角色Repository
	CreateRoleRepository(ctx TransactionContext) TransactionRoleRepository
}

// TransactionUserRepository 用户事务接口（简化版，用于事务上下文）
type TransactionUserRepository interface {
	CRUDRepository[*usersModel.User, string]
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	GetByEmail(ctx context.Context, email string) (*usersModel.User, error)
}

// TransactionAIRepository AI事务Repository接口
type TransactionAIRepository interface {
	CRUDRepository[*ai.AIModel, string]
	GetByType(ctx context.Context, modelType string) ([]*ai.AIModel, error)
}

// TransactionDocumentRepository 文档事务Repository接口
type TransactionDocumentRepository interface {
	CRUDRepository[*document.Document, string]
	GetByProjectID(ctx context.Context, projectID string) ([]*document.Document, error)
}

// TransactionRoleRepository 角色事务Repository接口
type TransactionRoleRepository interface {
	CRUDRepository[*usersModel.Role, string]
	GetDefaultRole(ctx context.Context) (*usersModel.Role, error)
	AssignRole(ctx context.Context, userID, roleID string) error
}

// MongoTransactionManager MongoDB事务管理器实现
type MongoTransactionManager struct {
	client                       *mongo.Client
	transactionRepositoryFactory TransactionRepositoryFactory
}

// NewMongoTransactionManager 创建MongoDB事务管理器
func NewMongoTransactionManager(client *mongo.Client, factory TransactionRepositoryFactory) TransactionManager {
	return &MongoTransactionManager{
		client:                       client,
		transactionRepositoryFactory: factory,
	}
}

// ExecuteTransaction 执行事务
func (tm *MongoTransactionManager) ExecuteTransaction(ctx context.Context, fn TransactionFunc) error {
	session, err := tm.client.StartSession()
	if err != nil {
		return fmt.Errorf("启动事务会话失败: %w", err)
	}
	defer session.EndSession(ctx)

	// 创建事务上下文
	txCtx := &mongoTransactionContext{
		Context:                      ctx,
		session:                      session,
		transactionRepositoryFactory: tm.transactionRepositoryFactory,
		transactionID:                generateTransactionID(),
		inTransaction:                true,
	}

	// 执行事务
	_, err = session.WithTransaction(ctx, func(sc mongo.SessionContext) (interface{}, error) {
		// 更新事务上下文
		txCtx.Context = sc

		// 执行业务逻辑
		if err := fn(txCtx); err != nil {
			return nil, err
		}

		// 检查是否标记为回滚
		if txCtx.IsMarkedForRollback() {
			return nil, fmt.Errorf("事务被标记为回滚: %s", txCtx.GetRollbackReason())
		}

		return nil, nil
	})

	if err != nil {
		log.Printf("事务执行失败 [%s]: %v", txCtx.GetTransactionID(), err)
		return fmt.Errorf("事务执行失败: %w", err)
	}

	log.Printf("事务执行成功 [%s]", txCtx.GetTransactionID())
	return nil
}

// ExecuteSaga 执行Saga事务
func (tm *MongoTransactionManager) ExecuteSaga(ctx context.Context, saga *Saga) error {
	sagaCtx := &SagaContext{
		Context: ctx,
		SagaID:  generateSagaID(),
		Steps:   make([]*SagaStepResult, 0, len(saga.Steps)),
	}

	log.Printf("开始执行Saga [%s]: %s", sagaCtx.SagaID, saga.Name)

	// 执行所有步骤
	for i, step := range saga.Steps {
		stepResult := &SagaStepResult{
			StepName:  step.Name,
			StepIndex: i,
			StartTime: time.Now(),
		}

		log.Printf("执行Saga步骤 [%s-%d]: %s", sagaCtx.SagaID, i, step.Name)

		// 执行步骤
		err := tm.ExecuteTransaction(ctx, step.Execute)
		stepResult.EndTime = time.Now()
		stepResult.Duration = stepResult.EndTime.Sub(stepResult.StartTime)

		if err != nil {
			stepResult.Error = err
			sagaCtx.Steps = append(sagaCtx.Steps, stepResult)

			log.Printf("Saga步骤失败 [%s-%d]: %v", sagaCtx.SagaID, i, err)

			// 执行补偿操作
			if compensateErr := tm.compensateSaga(ctx, sagaCtx, saga); compensateErr != nil {
				log.Printf("Saga补偿失败 [%s]: %v", sagaCtx.SagaID, compensateErr)
				return fmt.Errorf("Saga执行失败且补偿失败: 原错误=%w, 补偿错误=%v", err, compensateErr)
			}

			return fmt.Errorf("Saga执行失败: %w", err)
		}

		stepResult.Success = true
		sagaCtx.Steps = append(sagaCtx.Steps, stepResult)

		log.Printf("Saga步骤成功 [%s-%d]: %s (耗时: %v)", sagaCtx.SagaID, i, step.Name, stepResult.Duration)
	}

	log.Printf("Saga执行成功 [%s]: %s", sagaCtx.SagaID, saga.Name)
	return nil
}

// compensateSaga 执行Saga补偿
func (tm *MongoTransactionManager) compensateSaga(ctx context.Context, sagaCtx *SagaContext, saga *Saga) error {
	log.Printf("开始Saga补偿 [%s]", sagaCtx.SagaID)

	// 逆序执行补偿操作
	for i := len(sagaCtx.Steps) - 1; i >= 0; i-- {
		stepResult := sagaCtx.Steps[i]
		if !stepResult.Success {
			continue // 跳过失败的步骤
		}

		step := saga.Steps[stepResult.StepIndex]
		if step.Compensate == nil {
			continue // 跳过没有补偿操作的步骤
		}

		log.Printf("执行Saga补偿步骤 [%s-%d]: %s", sagaCtx.SagaID, i, step.Name)

		if err := tm.ExecuteTransaction(ctx, step.Compensate); err != nil {
			log.Printf("Saga补偿步骤失败 [%s-%d]: %v", sagaCtx.SagaID, i, err)
			return fmt.Errorf("补偿步骤 %s 失败: %w", step.Name, err)
		}

		log.Printf("Saga补偿步骤成功 [%s-%d]: %s", sagaCtx.SagaID, i, step.Name)
	}

	log.Printf("Saga补偿完成 [%s]", sagaCtx.SagaID)
	return nil
}

// GetTransactionContext 获取事务上下文
func (tm *MongoTransactionManager) GetTransactionContext(ctx context.Context) (TransactionContext, error) {
	// 检查是否已经在事务中
	if txCtx, ok := ctx.(TransactionContext); ok {
		return txCtx, nil
	}

	// 创建新的事务上下文（非事务模式）
	return &mongoTransactionContext{
		Context:                      ctx,
		transactionRepositoryFactory: tm.transactionRepositoryFactory,
		transactionID:                generateTransactionID(),
		inTransaction:                false,
	}, nil
}

// Health 健康检查
func (tm *MongoTransactionManager) Health(ctx context.Context) error {
	return tm.client.Ping(ctx, nil)
}

// mongoTransactionContext MongoDB事务上下文实现
type mongoTransactionContext struct {
	context.Context
	session                      mongo.Session
	transactionRepositoryFactory TransactionRepositoryFactory
	transactionID                string
	inTransaction                bool
	rollbackMarked               bool
	rollbackReason               string
}

// GetSession 获取MongoDB会话
func (ctx *mongoTransactionContext) GetSession() mongo.Session {
	return ctx.session
}

// GetRepositoryFactory 获取Repository工厂
func (ctx *mongoTransactionContext) GetRepositoryFactory() TransactionRepositoryFactory {
	return ctx.transactionRepositoryFactory
}

// IsInTransaction 是否在事务中
func (ctx *mongoTransactionContext) IsInTransaction() bool {
	return ctx.inTransaction
}

// GetTransactionID 获取事务ID
func (ctx *mongoTransactionContext) GetTransactionID() string {
	return ctx.transactionID
}

// MarkForRollback 标记为回滚
func (ctx *mongoTransactionContext) MarkForRollback(reason string) {
	ctx.rollbackMarked = true
	ctx.rollbackReason = reason
}

// IsMarkedForRollback 是否标记为回滚
func (ctx *mongoTransactionContext) IsMarkedForRollback() bool {
	return ctx.rollbackMarked
}

// GetRollbackReason 获取回滚原因
func (ctx *mongoTransactionContext) GetRollbackReason() string {
	return ctx.rollbackReason
}

// Saga 事务定义
type Saga struct {
	Name  string      `json:"name"`
	Steps []*SagaStep `json:"steps"`
}

// SagaStep Saga步骤
type SagaStep struct {
	Name       string          `json:"name"`
	Execute    TransactionFunc `json:"-"`
	Compensate TransactionFunc `json:"-"`
}

// SagaContext Saga上下文
type SagaContext struct {
	context.Context
	SagaID string            `json:"sagaId"`
	Steps  []*SagaStepResult `json:"steps"`
}

// SagaStepResult Saga步骤结果
type SagaStepResult struct {
	StepName  string        `json:"stepName"`
	StepIndex int           `json:"stepIndex"`
	Success   bool          `json:"success"`
	Error     error         `json:"error,omitempty"`
	StartTime time.Time     `json:"startTime"`
	EndTime   time.Time     `json:"endTime"`
	Duration  time.Duration `json:"duration"`
}

// 工具函数
func generateTransactionID() string {
	return fmt.Sprintf("tx_%d", time.Now().UnixNano())
}

func generateSagaID() string {
	return fmt.Sprintf("saga_%d", time.Now().UnixNano())
}

// 预定义的Saga事务

// NewUserRegistrationSaga 创建用户注册Saga
func NewUserRegistrationSaga(userReq *UserRegistrationRequest) *Saga {
	return &Saga{
		Name: "用户注册",
		Steps: []*SagaStep{
			{
				Name: "创建用户",
				Execute: func(txCtx TransactionContext) error {
					userRepo := txCtx.GetRepositoryFactory().CreateUserRepository(txCtx)

					// 检查用户是否已存在
					exists, err := userRepo.ExistsByEmail(txCtx, userReq.Email)
					if err != nil {
						return fmt.Errorf("检查用户邮箱失败: %w", err)
					}
					if exists {
						return fmt.Errorf("邮箱已被注册: %s", userReq.Email)
					}

					// 创建用户
					user := &usersModel.User{
						Username:  userReq.Username,
						Email:     userReq.Email,
						Password:  userReq.HashedPassword,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}

					// Note: 实际实现中需要调用 userRepo.Create 方法
					return userRepo.Create(txCtx, user)
				},
				Compensate: func(txCtx TransactionContext) error {
					userRepo := txCtx.GetRepositoryFactory().CreateUserRepository(txCtx)

					// 查找用户
					user, err := userRepo.GetByEmail(txCtx, userReq.Email)
					if err != nil {
						return fmt.Errorf("获取用户失败: %w", err) // 用户不存在，无需补偿
					}

					// 删除用户
					return userRepo.Delete(txCtx, user.ID)
				},
			},
			{
				Name: "分配默认角色",
				Execute: func(txCtx TransactionContext) error {
					userRepo := txCtx.GetRepositoryFactory().CreateUserRepository(txCtx)
					roleRepo := txCtx.GetRepositoryFactory().CreateRoleRepository(txCtx)

					// 获取用户
					user, err := userRepo.GetByEmail(txCtx, userReq.Email)
					if err != nil {
						return fmt.Errorf("获取用户失败: %w", err)
					}

					// 获取默认角色
					defaultRole, err := roleRepo.GetDefaultRole(txCtx)
					if err != nil {
						return fmt.Errorf("获取默认角色失败: %w", err)
					}

					// 分配角色
					return roleRepo.AssignRole(txCtx, user.ID, defaultRole.ID)
				},
				Compensate: func(txCtx TransactionContext) error {
					// 角色分配的补偿操作
					// 这里可以实现取消角色分配的逻辑
					return nil
				},
			},
		},
	}
}

// UserRegistrationRequest 用户注册请求
type UserRegistrationRequest struct {
	Username       string `json:"username" validate:"required,min=3,max=50"`
	Email          string `json:"email" validate:"required,email"`
	HashedPassword string `json:"-"`
}
