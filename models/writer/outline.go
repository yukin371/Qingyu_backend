package writer

import "time"

// OutlineNode 大纲节点 (细化的 Node)
type OutlineNode struct {
	ID        string `bson:"_id,omitempty" json:"id"`
	ProjectID string `bson:"project_id" json:"projectId"`
	ParentID  string `bson:"parent_id,omitempty" json:"parentId,omitempty"` // 卷/幕
	Title     string `bson:"title" json:"title"`
	Summary   string `bson:"summary" json:"summary"` // 本章/本节摘要

	// 结构属性
	Type    string `bson:"type" json:"type"`       // 英雄之旅阶段(如: 召唤、深渊)、起承转合
	Tension int    `bson:"tension" json:"tension"` // 紧张度/情绪值 (1-10)，用于生成情绪曲线

	// 关联
	ChapterID  string   `bson:"chapter_id,omitempty" json:"chapterId,omitempty"`  // 对应实际写的章节ID
	Characters []string `bson:"characters,omitempty" json:"characters,omitempty"` // 本节登场人物
	Items      []string `bson:"items,omitempty" json:"items,omitempty"`           // 涉及道具

	Order     int       `bson:"order" json:"order"`
	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}
