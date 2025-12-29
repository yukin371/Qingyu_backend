package writer

import "time"

// Organization 组织/势力
type Organization struct {
	ID             string `bson:"_id,omitempty" json:"id"`
	ProjectID      string `bson:"project_id" json:"projectId"`
	Name           string `bson:"name" json:"name"`
	Type           string `bson:"type" json:"type"` // 宗门、国家、公会、公司、家族
	Description    string `bson:"description" json:"description"`
	LeaderID       string `bson:"leader_id,omitempty" json:"leaderId,omitempty"`              // 领袖角色ID
	BaseLocationID string `bson:"base_location_id,omitempty" json:"baseLocationId,omitempty"` // 总部地点ID

	// 组织架构
	ParentID string   `bson:"parent_id,omitempty" json:"parentId,omitempty"` // 上级组织
	Members  []string `bson:"members,omitempty" json:"members,omitempty"`    // 核心成员ID列表(冗余，主要通过Character.OrgID查询)

	// 设定细节
	Motto     string `bson:"motto,omitempty" json:"motto,omitempty"`         // 信条/口号
	Resources string `bson:"resources,omitempty" json:"resources,omitempty"` // 拥有的资源/特产

	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}

// OrgRelation 组织外交关系
// 不同于 RoleRelation，这是宏观层面的
type OrgRelation struct {
	ID        string `bson:"_id,omitempty" json:"id"`
	FromOrgID string `bson:"from_org_id" json:"fromOrgId"`
	ToOrgID   string `bson:"to_org_id" json:"toOrgId"`
	Relation  string `bson:"relation" json:"relation"` // 敌对、盟友、从属、中立
	Notes     string `bson:"notes" json:"notes"`
}
