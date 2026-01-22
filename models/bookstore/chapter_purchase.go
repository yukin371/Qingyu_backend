package bookstore

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChapterPurchase 章节购买记录模型
type ChapterPurchase struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"$1$2"`
	ChapterID     primitive.ObjectID `bson:"chapter_id" json:"$1$2"`
	BookID        primitive.ObjectID `bson:"book_id" json:"$1$2"`
	Price         int64              `bson:"price" json:"price"` // 价格 (分)
	PurchaseTime  time.Time          `bson:"purchase_time" json:"$1$2"`
	TransactionID string             `bson:"transaction_id,omitempty" json:"transaction_id,omitempty"`
	CreatedAt     time.Time          `bson:"created_at" json:"$1$2"`

	// 冗余字段（用于快速查询）
	ChapterTitle string `bson:"chapter_title,omitempty" json:"chapterTitle,omitempty"`
	ChapterNum   int    `bson:"chapter_num,omitempty" json:"chapterNum,omitempty"`
	BookTitle    string `bson:"book_title,omitempty" json:"bookTitle,omitempty"`
	BookCover    string `bson:"book_cover,omitempty" json:"bookCover,omitempty"`
}

// ChapterPurchaseBatch 章节批量购买记录
type ChapterPurchaseBatch struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID   `bson:"user_id" json:"$1$2"`
	BookID        primitive.ObjectID   `bson:"book_id" json:"$1$2"`
	ChapterIDs    []primitive.ObjectID `bson:"chapter_ids" json:"$1$2"`
	TotalPrice    int64                `bson:"total_price" json:"$1$2"` // 总价 (分)
	ChaptersCount int                  `bson:"chapters_count" json:"$1$2"`
	PurchaseTime  time.Time            `bson:"purchase_time" json:"$1$2"`
	TransactionID string               `bson:"transaction_id,omitempty" json:"transaction_id,omitempty"`
	CreatedAt     time.Time            `bson:"created_at" json:"$1$2"`

	// 冗余字段
	BookTitle string `bson:"book_title,omitempty" json:"bookTitle,omitempty"`
	BookCover string `bson:"book_cover,omitempty" json:"bookCover,omitempty"`
}

// TablePurchase 全书购买记录
type BookPurchase struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"$1$2"`
	BookID        primitive.ObjectID `bson:"book_id" json:"$1$2"`
	TotalPrice    int64              `bson:"total_price" json:"$1$2"`       // 总价 (分)
	OriginalPrice int64              `bson:"original_price" json:"$1$2"` // 原价 (分)
	Discount      float64            `bson:"discount" json:"discount"`             // 折扣（0-1，浮点数）
	PurchaseTime  time.Time          `bson:"purchase_time" json:"$1$2"`
	TransactionID string             `bson:"transaction_id,omitempty" json:"transaction_id,omitempty"`
	CreatedAt     time.Time          `bson:"created_at" json:"$1$2"`

	// 冗余字段
	BookTitle    string `bson:"book_title,omitempty" json:"bookTitle,omitempty"`
	BookCover    string `bson:"book_cover,omitempty" json:"bookCover,omitempty"`
	ChapterCount int    `bson:"chapter_count,omitempty" json:"chapterCount,omitempty"`
}

// BeforeCreate 在创建前设置时间戳
func (cp *ChapterPurchase) BeforeCreate() {
	now := time.Now()
	cp.CreatedAt = now
	if cp.PurchaseTime.IsZero() {
		cp.PurchaseTime = now
	}
}

// BeforeCreate 在创建前设置时间戳
func (cpb *ChapterPurchaseBatch) BeforeCreate() {
	now := time.Now()
	cpb.CreatedAt = now
	if cpb.PurchaseTime.IsZero() {
		cpb.PurchaseTime = now
	}
}

// BeforeCreate 在创建前设置时间戳
func (bp *BookPurchase) BeforeCreate() {
	now := time.Now()
	bp.CreatedAt = now
	if bp.PurchaseTime.IsZero() {
		bp.PurchaseTime = now
	}
}

// ChapterAccessInfo 章节访问信息（用于API响应）
type ChapterAccessInfo struct {
	ChapterID    primitive.ObjectID `json:"$1$2"`
	Title        string             `json:"title"`
	ChapterNum   int                `json:"$1$2"`
	WordCount    int                `json:"$1$2"`
	IsFree       bool               `json:"$1$2"`
	Price        int64              `json:"price"` // 价格 (分)
	IsPurchased  bool               `json:"$1$2"`
	IsVIP        bool               `json:"$1$2"`
	CanAccess    bool               `json:"$1$2"`
	AccessReason string             `json:"access_reason,omitempty"` // free, purchased, vip
	PurchaseTime *time.Time         `json:"purchase_time,omitempty"`
}

// ChapterCatalogItem 章节目录项
type ChapterCatalogItem struct {
	ChapterID   primitive.ObjectID `json:"$1$2"`
	Title       string             `json:"title"`
	ChapterNum  int                `json:"$1$2"`
	WordCount   int                `json:"$1$2"`
	IsFree      bool               `json:"$1$2"`
	Price       int64              `json:"price"` // 价格 (分)
	PublishTime time.Time          `json:"$1$2"`
	IsPublished bool               `json:"$1$2"`
	IsPurchased bool               `json:"is_purchased,omitempty"` // 仅在认证用户的请求中返回
	IsVIP       bool               `json:"is_vip,omitempty"`       // VIP专属章节
}

// ChapterCatalog 章节目录
type ChapterCatalog struct {
	BookID         primitive.ObjectID   `json:"$1$2"`
	BookTitle      string               `json:"$1$2"`
	TotalChapters  int                  `json:"$1$2"`
	FreeChapters   int                  `json:"$1$2"`
	PaidChapters   int                  `json:"$1$2"`
	VIPChapters    int                  `json:"$1$2"`
	TotalWordCount int64                `json:"$1$2"`
	Chapters       []ChapterCatalogItem `json:"chapters"`
	TrialCount     int                  `json:"$1$2"` // 可试读章节数量
}
