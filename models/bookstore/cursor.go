package bookstore

// CursorType 游标类型
type CursorType string

const (
	CursorTypeOffset    CursorType = "offset"    // 偏移量游标 (兼容老版本)
	CursorTypeTimestamp CursorType = "timestamp" // 时间戳游标
	CursorTypeID        CursorType = "id"        // ID游标
)

// StreamCursor 流式搜索游标
type StreamCursor struct {
	Type      CursorType `json:"type"`
	Value     string     `json:"value"`     // 游标值 (JSON编码后)
	Timestamp int64      `json:"timestamp"` // 创建时间
	TTL       int64      `json:"ttl"`       // 有效期 (秒)
}

// StreamResponse 流式响应类型
type StreamResponseType string

const (
	StreamResponseTypeMeta     StreamResponseType = "meta"     // 元数据
	StreamResponseTypeData     StreamResponseType = "data"     // 数据
	StreamResponseTypeProgress StreamResponseType = "progress" // 进度
	StreamResponseTypeDone     StreamResponseType = "done"     // 完成
	StreamResponseTypeError    StreamResponseType = "error"    // 错误
)

// StreamMessage 流式消息
type StreamMessage struct {
	Type   StreamResponseType `json:"type"`
	Cursor string             `json:"cursor,omitempty"`
	Total  *int64             `json:"total,omitempty"`
	HasMore bool              `json:"hasMore,omitempty"`

	// Data字段
	Books []*Book `json:"books,omitempty"`

	// Progress字段
	Loaded int `json:"loaded,omitempty"`

	// Error字段
	Error string `json:"error,omitempty"`
}
