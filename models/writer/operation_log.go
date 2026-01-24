package writer

import (
	"Qingyu_backend/models/writer/base"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OperationLog 操作日志（用于Undo/Redo）
type OperationLog struct {
	base.IdentifiedEntity `bson:",inline"`
	base.Timestamps       `bson:",inline"`

	// 关联信息
	ProjectID primitive.ObjectID  `bson:"project_id" json:"projectId" validate:"required"`
	UserID    primitive.ObjectID  `bson:"user_id" json:"userId" validate:"required"`
	BatchOpID *primitive.ObjectID `bson:"batch_op_id,omitempty" json:"batchOpId,omitempty"`
	ChainID   string              `bson:"chain_id" json:"chainId"` // 链ID，用于批量操作

	// 命令信息
	CommandType    DocumentCommandType   `bson:"command_type" json:"commandType" validate:"required"`
	TargetIDs      []string              `bson:"target_ids" json:"targetIds" validate:"required,min=1"`
	CommandPayload map[string]interface{} `bson:"command_payload,omitempty" json:"commandPayload,omitempty"`

	// 状态
	Status      OperationLogStatus `bson:"status" json:"status"`
	IsCommitted bool               `bson:"is_committed" json:"isCommitted"` // 是否已提交到数据库

	// 撤销信息
	InverseCommand map[string]interface{} `bson:"inverse_command,omitempty" json:"inverseCommand,omitempty"`
	UndoneAt       *time.Time             `bson:"undone_at,omitempty" json:"undoneAt,omitempty"`
	RedoneAt       *time.Time             `bson:"redone_at,omitempty" json:"redoneAt,omitempty"`
}

// DocumentCommandType 文档命令类型
type DocumentCommandType string

const (
	CommandCreate  DocumentCommandType = "create"
	CommandUpdate  DocumentCommandType = "update"
	CommandMove    DocumentCommandType = "move"
	CommandCopy    DocumentCommandType = "copy"
	CommandDelete  DocumentCommandType = "delete"
	CommandRestore DocumentCommandType = "restore"
)

// OperationLogStatus 操作日志状态
type OperationLogStatus string

const (
	OpLogStatusExecuted OperationLogStatus = "executed"
	OpLogStatusUndone   OperationLogStatus = "undone"
	OpLogStatusRedone   OperationLogStatus = "redone"
)

// TouchForCreate 创建时设置默认值
func (o *OperationLog) TouchForCreate() {
	o.IdentifiedEntity.GenerateID()
	o.Timestamps.TouchForCreate()
	if o.Status == "" {
		o.Status = OpLogStatusExecuted
	}
	if o.ChainID == "" {
		o.ChainID = o.ID.Hex()
	}
}

// IsUndoable 判断是否可撤销
func (o *OperationLog) IsUndoable() bool {
	return o.Status == OpLogStatusExecuted && o.IsCommitted
}

// IsRedoable 判断是否可重做
func (o *OperationLog) IsRedoable() bool {
	return o.Status == OpLogStatusUndone && o.InverseCommand != nil
}
