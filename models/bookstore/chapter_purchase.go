package bookstore

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChapterPurchase 章节购买记录模型
type ChapterPurchase struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	ChapterID     primitive.ObjectID `bson:"chapter_id" json:"chapter_id"`
	BookID        primitive.ObjectID `bson:"book_id" json:"book_id"`
	Price         int64              `bson:"price" json:"price"` // 价格 (分)
	PurchaseTime  time.Time          `bson:"purchase_time" json:"purchase_time"`
	TransactionID string             `bson:"transaction_id,omitempty" json:"transaction_id,omitempty"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`

	// 冗余字段（用于快速查询）
	ChapterTitle string `bson:"chapter_title,omitempty" json:"chapterTitle,omitempty"`
	ChapterNum   int    `bson:"chapter_num,omitempty" json:"chapterNum,omitempty"`
	BookTitle    string `bson:"book_title,omitempty" json:"bookTitle,omitempty"`
	BookCover    string `bson:"book_cover,omitempty" json:"bookCover,omitempty"`
}

// ChapterPurchaseBatch 章节批量购买记录
type ChapterPurchaseBatch struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID   `bson:"user_id" json:"user_id"`
	BookID        primitive.ObjectID   `bson:"book_id" json:"book_id"`
	ChapterIDs    []primitive.ObjectID `bson:"chapter_ids" json:"chapter_ids"`
	TotalPrice    int64                `bson:"total_price" json:"total_price"` // 总价 (分)
	ChaptersCount int                  `bson:"chapters_count" json:"chapters_count"`
	PurchaseTime  time.Time            `bson:"purchase_time" json:"purchase_time"`
	TransactionID string               `bson:"transaction_id,omitempty" json:"transaction_id,omitempty"`
	CreatedAt     time.Time            `bson:"created_at" json:"created_at"`

	// 冗余字段
	BookTitle string `bson:"book_title,omitempty" json:"bookTitle,omitempty"`
	BookCover string `bson:"book_cover,omitempty" json:"bookCover,omitempty"`
}

// TablePurchase 全书购买记录
type BookPurchase struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id"`
	BookID        primitive.ObjectID `bson:"book_id" json:"book_id"`
	TotalPrice    int64              `bson:"total_price" json:"total_price"`       // 总价 (分)
	OriginalPrice int64              `bson:"original_price" json:"original_price"` // 原价 (分)
	Discount      float64            `bson:"discount" json:"discount"`             // 折扣（0-1，浮点数）
	PurchaseTime  time.Time          `bson:"purchase_time" json:"purchase_time"`
	TransactionID string             `bson:"transaction_id,omitempty" json:"transaction_id,omitempty"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`

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
	ChapterID    primitive.ObjectID `json:"chapter_id"`
	Title        string             `json:"title"`
	ChapterNum   int                `json:"chapter_num"`
	WordCount    int                `json:"word_count"`
	IsFree       bool               `json:"is_free"`
	Price        int64              `json:"price"` // 价格 (分)
	IsPurchased  bool               `json:"is_purchased"`
	IsVIP        bool               `json:"is_vip"`
	CanAccess    bool               `json:"can_access"`
	AccessReason string             `json:"access_reason,omitempty"` // free, purchased, vip
	PurchaseTime *time.Time         `json:"purchase_time,omitempty"`
}

// ChapterCatalogItem 章节目录项
type ChapterCatalogItem struct {
	ChapterID   primitive.ObjectID `json:"chapter_id"`
	Title       string             `json:"title"`
	ChapterNum  int                `json:"chapter_num"`
	WordCount   int                `json:"word_count"`
	IsFree      bool               `json:"is_free"`
	Price       int64              `json:"price"` // 价格 (分)
	PublishTime time.Time          `json:"publish_time"`
	IsPublished bool               `json:"is_published"`
	IsPurchased bool               `json:"is_purchased,omitempty"` // 仅在认证用户的请求中返回
	IsVIP       bool               `json:"is_vip,omitempty"`       // VIP专属章节
}

// ChapterCatalog 章节目录
type ChapterCatalog struct {
	BookID         primitive.ObjectID   `json:"book_id"`
	BookTitle      string               `json:"book_title"`
	TotalChapters  int                  `json:"total_chapters"`
	FreeChapters   int                  `json:"free_chapters"`
	PaidChapters   int                  `json:"paid_chapters"`
	VIPChapters    int                  `json:"vip_chapters"`
	TotalWordCount int64                `json:"total_word_count"`
	Chapters       []ChapterCatalogItem `json:"chapters"`
	TrialCount     int                  `json:"trial_count"` // 可试读章节数量
}
