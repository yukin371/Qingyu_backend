package writer

import "time"

// Paragraph 是文档分段内容的领域类型（轻量值对象）。
// 注意：Paragraph 不是独立持久化实体，仍然存储在 document_content 集合中。
type Paragraph struct {
	ID          string
	DocumentID  string
	Order       int
	Content     string
	ContentType string
	Version     int
	UpdatedAt   time.Time
}

// Normalize 归一化段落字段，避免在 Service 层重复处理默认值逻辑。
func (p *Paragraph) Normalize(defaultOrder int) {
	if p.Order <= 0 {
		p.Order = defaultOrder
	}
	if p.ContentType == "" {
		p.ContentType = "tiptap"
	}
	if p.Version < 0 {
		p.Version = 0
	}
}
