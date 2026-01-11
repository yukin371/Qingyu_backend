package shared

// Likable 支持点赞的实体混入
type Likable struct {
	LikeCount int `bson:"like_count" json:"likeCount"`
}

// AddLike 增加点赞计数
func (l *Likable) AddLike() {
	l.LikeCount++
}

// RemoveLike 减少点赞计数
func (l *Likable) RemoveLike() {
	if l.LikeCount > 0 {
		l.LikeCount--
	}
}

// ThreadedConversation 支持回复/评论的实体混入
type ThreadedConversation struct {
	ParentID   *string `bson:"parent_id,omitempty" json:"parentId,omitempty"`
	RootID     *string `bson:"root_id,omitempty" json:"rootId,omitempty"`
	ReplyCount int     `bson:"reply_count" json:"replyCount"`
}

// AddReply 增加回复计数
func (t *ThreadedConversation) AddReply() {
	t.ReplyCount++
}

// RemoveReply 减少回复计数
func (t *ThreadedConversation) RemoveReply() {
	if t.ReplyCount > 0 {
		t.ReplyCount--
	}
}

// IsReply 判断是否为回复
func (t *ThreadedConversation) IsReply() bool {
	return t.ParentID != nil
}
