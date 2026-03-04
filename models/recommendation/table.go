package recommendation

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TableType 推荐榜类型
type TableType string

const (
	TableTypeWeekly      TableType = "weekly"
	TableTypeMonthly     TableType = "monthly"
	TableTypeMonthlyVote TableType = "monthly_vote"
	TableTypeManual      TableType = "manual"
)

// TableSource 榜单来源
type TableSource string

const (
	TableSourceAuto   TableSource = "auto"
	TableSourceManual TableSource = "manual"
)

// TableStatus 榜单状态
type TableStatus string

const (
	TableStatusActive   TableStatus = "active"
	TableStatusArchived TableStatus = "archived"
)

// TableItem 榜单项
type TableItem struct {
	BookID  string  `bson:"book_id" json:"bookId"`
	Rank    int     `bson:"rank" json:"rank"`
	Score   float64 `bson:"score,omitempty" json:"score,omitempty"`
	Reason  string  `bson:"reason,omitempty" json:"reason,omitempty"`
	Manual  bool    `bson:"manual" json:"manual"`
	AddedBy string  `bson:"added_by,omitempty" json:"addedBy,omitempty"`
	AddedAt int64   `bson:"added_at,omitempty" json:"addedAt,omitempty"`
}

// RecommendationTable 推荐榜表
type RecommendationTable struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Name      string                 `bson:"name" json:"name"`
	Type      TableType              `bson:"type" json:"type"`
	Period    string                 `bson:"period" json:"period"` // weekly: 2026-W10, monthly: 2026-03, manual: custom
	Source    TableSource            `bson:"source" json:"source"`
	Status    TableStatus            `bson:"status" json:"status"`
	Items     []TableItem            `bson:"items" json:"items"`
	Metadata  map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
	UpdatedBy string                 `bson:"updated_by,omitempty" json:"updatedBy,omitempty"`
	CreatedAt time.Time              `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time              `bson:"updated_at" json:"updatedAt"`
}

func IsValidTableType(v string) bool {
	switch TableType(v) {
	case TableTypeWeekly, TableTypeMonthly, TableTypeMonthlyVote, TableTypeManual:
		return true
	default:
		return false
	}
}

func IsValidTableSource(v string) bool {
	switch TableSource(v) {
	case TableSourceAuto, TableSourceManual:
		return true
	default:
		return false
	}
}
