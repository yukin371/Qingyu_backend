package bookstore

import (
	"Qingyu_backend/models/bookstore"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	BookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
	"Qingyu_backend/service/finance/wallet"
)

// ChapterPurchaseService 章节购买服务接口
type ChapterPurchaseService interface {
	// 章节目录和权限
	GetChapterCatalog(ctx context.Context, userID, bookID string) (*bookstore.ChapterCatalog, error)
	GetTrialChapters(ctx context.Context, bookID string, trialCount int) ([]*bookstore.Chapter, error)
	GetVIPChapters(ctx context.Context, bookID string) ([]*bookstore.Chapter, error)

	// 购买章节
	PurchaseChapter(ctx context.Context, userID, chapterID string) (*bookstore.ChapterPurchase, error)
	PurchaseChapters(ctx context.Context, userID string, chapterIDs []string) (*bookstore.ChapterPurchaseBatch, error)
	PurchaseBook(ctx context.Context, userID, bookID string) (*bookstore.BookPurchase, error)

	// 查询购买记录
	GetChapterPurchases(ctx context.Context, userID string, page, pageSize int) ([]*bookstore.ChapterPurchase, int64, error)
	GetBookPurchases(ctx context.Context, userID, bookID string, page, pageSize int) ([]*bookstore.ChapterPurchase, int64, error)
	GetAllPurchases(ctx context.Context, userID string, page, pageSize int) (map[string]interface{}, error)

	// 权限检查
	CheckChapterAccess(ctx context.Context, userID, chapterID string) (*bookstore.ChapterAccessInfo, error)
	GetPurchasedChapterIDs(ctx context.Context, userID, bookID string) ([]string, error)

	// 价格查询
	GetChapterPrice(ctx context.Context, chapterID string) (float64, error)
	CalculateBookPrice(ctx context.Context, bookID string) (int64, int64, error) // originalPrice, discountedPrice (分)

	// VIP检查
	IsVIPUser(ctx context.Context, userID string) (bool, error)
}

// ChapterPurchaseServiceImpl 章节购买服务实现
type ChapterPurchaseServiceImpl struct {
	chapterRepo   BookstoreRepo.ChapterRepository
	purchaseRepo  BookstoreRepo.ChapterPurchaseRepository
	bookRepo      BookstoreRepo.BookRepository
	walletService wallet.WalletService
	cacheService  CacheService
}

// NewChapterPurchaseService 创建章节购买服务实例
func NewChapterPurchaseService(
	chapterRepo BookstoreRepo.ChapterRepository,
	purchaseRepo BookstoreRepo.ChapterPurchaseRepository,
	bookRepo BookstoreRepo.BookRepository,
	walletService wallet.WalletService,
	cacheService CacheService,
) ChapterPurchaseService {
	return &ChapterPurchaseServiceImpl{
		chapterRepo:   chapterRepo,
		purchaseRepo:  purchaseRepo,
		bookRepo:      bookRepo,
		walletService: walletService,
		cacheService:  cacheService,
	}
}

// GetChapterCatalog 获取章节目录
func (s *ChapterPurchaseServiceImpl) GetChapterCatalog(ctx context.Context, userID, bookID string) (*bookstore.ChapterCatalog, error) {
	if bookID == "" {
		return nil, errors.New("book ID cannot be empty")
	}

	// 获取书籍信息
	book, err := s.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get book: %w", err)
	}
	if book == nil {
		return nil, errors.New("book not found")
	}

	// 获取所有章节
	chapters, err := s.chapterRepo.GetByBookID(ctx, bookID, 10000, 0) // 获取所有章节
	if err != nil {
		return nil, fmt.Errorf("failed to get chapters: %w", err)
	}

	// 获取用户已购买的章节ID（如果提供了用户ID）
	var purchasedChapterIDs []string
	if userID != "" {
		purchasedChapterIDs, err = s.purchaseRepo.GetPurchasedChapterIDs(ctx, userID, bookID)
		if err != nil {
			// 不影响主流程，继续处理
			purchasedChapterIDs = []string{}
		}
	}

	// 构建章节目录项
	catalogItems := make([]bookstore.ChapterCatalogItem, 0, len(chapters))
	freeCount := 0
	paidCount := 0
	vipCount := 0

	purchasedIDSet := make(map[string]bool)
	for _, id := range purchasedChapterIDs {
		purchasedIDSet[id] = true
	}

	for _, chapter := range chapters {
		chapterOID, _ := primitive.ObjectIDFromHex(chapter.ID)
		item := bookstore.ChapterCatalogItem{
			ChapterID:   chapterOID,
			Title:       chapter.Title,
			ChapterNum:  chapter.ChapterNum,
			WordCount:   chapter.WordCount,
			IsFree:      chapter.IsFree,
			Price:       chapter.Price,
			PublishTime: chapter.PublishTime,
			IsPublished: chapter.IsPublished(),
		}

		// 标记是否已购买
		if userID != "" {
			item.IsPurchased = purchasedIDSet[chapter.ID]
		}

		// 统计
		if chapter.IsFree {
			freeCount++
		} else {
			paidCount++
		}

		catalogItems = append(catalogItems, item)
	}

	// 获取统计信息
	totalChapters, _ := s.chapterRepo.CountByBookID(ctx, bookID)
	totalWordCount, _ := s.chapterRepo.GetTotalWordCount(ctx, bookID)

	// 构建目录
	bookOID, _ := primitive.ObjectIDFromHex(bookID)
	catalog := &bookstore.ChapterCatalog{
		BookID:         bookOID,
		BookTitle:      book.Title,
		TotalChapters:  int(totalChapters),
		FreeChapters:   freeCount,
		PaidChapters:   paidCount,
		VIPChapters:    vipCount,
		TotalWordCount: totalWordCount,
		Chapters:       catalogItems,
		TrialCount:     10, // 默认试读前10章
	}

	return catalog, nil
}

// GetTrialChapters 获取试读章节
func (s *ChapterPurchaseServiceImpl) GetTrialChapters(ctx context.Context, bookID string, trialCount int) ([]*bookstore.Chapter, error) {
	if bookID == "" {
		return nil, errors.New("book ID cannot be empty")
	}

	if trialCount <= 0 {
		trialCount = 10 // 默认试读前10章
	}

	// 获取免费章节
	chapters, err := s.chapterRepo.GetFreeChapters(ctx, bookID, trialCount, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get trial chapters: %w", err)
	}

	// 如果免费章节数量不足，补充付费章节（前N章）
	if len(chapters) < trialCount {
		needed := trialCount - len(chapters)
		allChapters, _ := s.chapterRepo.GetByBookID(ctx, bookID, needed, len(chapters))
		for _, ch := range allChapters {
			if !ch.IsFree {
				chapters = append(chapters, ch)
			}
		}
	}

	return chapters, nil
}

// GetVIPChapters 获取VIP章节
func (s *ChapterPurchaseServiceImpl) GetVIPChapters(ctx context.Context, bookID string) ([]*bookstore.Chapter, error) {
	if bookID == "" {
		return nil, errors.New("book ID cannot be empty")
	}

	// VIP章节定义为价格高于平均价格的章节
	// 或者可以通过书籍模型的VIP标签来确定
	chapters, err := s.chapterRepo.GetPaidChapters(ctx, bookID, 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get VIP chapters: %w", err)
	}

	// 简化处理：所有付费章节都是VIP章节
	// 实际可以根据业务规则筛选
	return chapters, nil
}

// PurchaseChapter 购买单个章节
func (s *ChapterPurchaseServiceImpl) PurchaseChapter(ctx context.Context, userID, chapterID string) (*bookstore.ChapterPurchase, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}
	if chapterID == "" {
		return nil, errors.New("chapter ID cannot be empty")
	}

	// 检查章节是否已购买
	existingPurchase, err := s.purchaseRepo.GetByUserAndChapter(ctx, userID, chapterID)
	if err == nil && existingPurchase != nil {
		return nil, errors.New("chapter already purchased")
	}

	// 获取章节信息
	chapter, err := s.chapterRepo.GetByID(ctx, chapterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chapter: %w", err)
	}
	if chapter == nil {
		return nil, errors.New("chapter not found")
	}

	// 检查是否为免费章节
	if chapter.IsFree {
		return nil, errors.New("cannot purchase free chapter")
	}

	// 获取书籍信息
	book, err := s.bookRepo.GetByID(ctx, chapter.BookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get book: %w", err)
	}
	if book == nil {
		return nil, errors.New("book not found")
	}

	// 检查用户余额
	balance, err := s.walletService.GetBalance(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user balance: %w", err)
	}
	if balance < chapter.Price {
		return nil, errors.New("insufficient balance")
	}

	// 使用事务处理购买
	var purchase *bookstore.ChapterPurchase
	err = s.purchaseRepo.Transaction(ctx, func(txCtx context.Context) error {
		// 扣除用户余额
		_, err := s.walletService.Consume(ctx, userID, chapter.Price, fmt.Sprintf("购买章节: %s", chapter.Title))
		if err != nil {
			return fmt.Errorf("failed to deduct balance: %w", err)
		}

		// 创建购买记录 - 需要将 string 转换为 primitive.ObjectID
		userOID, _ := primitive.ObjectIDFromHex(userID)
		chapterOID, _ := primitive.ObjectIDFromHex(chapterID)
		bookOID, _ := primitive.ObjectIDFromHex(chapter.BookID)
		purchase = &bookstore.ChapterPurchase{
			UserID:       userOID,
			ChapterID:    chapterOID,
			BookID:       bookOID,
			Price:        chapter.Price,
			PurchaseTime: time.Now(),
			ChapterTitle: chapter.Title,
			ChapterNum:   chapter.ChapterNum,
			BookTitle:    book.Title,
			BookCover:    book.Cover,
		}
		purchase.BeforeCreate()

		if err := s.purchaseRepo.Create(txCtx, purchase); err != nil {
			return fmt.Errorf("failed to create purchase record: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 清除缓存
	if s.cacheService != nil {
		s.cacheService.InvalidateChapterCache(ctx, chapterID)
		s.cacheService.InvalidateBookChaptersCache(ctx, chapter.BookID)
	}

	return purchase, nil
}

// PurchaseChapters 批量购买章节
func (s *ChapterPurchaseServiceImpl) PurchaseChapters(ctx context.Context, userID string, chapterIDs []string) (*bookstore.ChapterPurchaseBatch, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}
	if len(chapterIDs) == 0 {
		return nil, errors.New("chapter IDs cannot be empty")
	}

	// 获取所有章节信息
	chapters := make([]*bookstore.Chapter, 0, len(chapterIDs))
	totalPrice := int64(0)
	bookID := ""

	for _, chapterID := range chapterIDs {
		chapter, err := s.chapterRepo.GetByID(ctx, chapterID)
		if err != nil {
			return nil, fmt.Errorf("failed to get chapter %s: %w", chapterID, err)
		}
		if chapter == nil {
			return nil, fmt.Errorf("chapter %s not found", chapterID)
		}

		// 检查是否已购买
		existingPurchase, _ := s.purchaseRepo.GetByUserAndChapter(ctx, userID, chapterID)
		if existingPurchase != nil {
			continue // 跳过已购买的章节
		}

		if chapter.IsFree {
			continue // 跳过免费章节
		}

		chapters = append(chapters, chapter)
		totalPrice += chapter.Price

		if bookID == "" {
			bookID = chapter.BookID
		}
	}

	if len(chapters) == 0 {
		return nil, errors.New("no chapters to purchase")
	}

	// 检查用户余额
	balance, err := s.walletService.GetBalance(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user balance: %w", err)
	}
	if balance < totalPrice {
		return nil, errors.New("insufficient balance")
	}

	// 获取书籍信息
	book, err := s.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	// 使用事务处理购买
	var batch *bookstore.ChapterPurchaseBatch
	purchasedChapterIDs := make([]string, 0)

	err = s.purchaseRepo.Transaction(ctx, func(txCtx context.Context) error {
		// 扣除用户余额
		_, err := s.walletService.Consume(ctx, userID, totalPrice, fmt.Sprintf("批量购买章节: %d章", len(chapters)))
		if err != nil {
			return fmt.Errorf("failed to deduct balance: %w", err)
		}

		// 创建批量购买记录 - 需要转换类型
		userOID, _ := primitive.ObjectIDFromHex(userID)
		bookOID, _ := primitive.ObjectIDFromHex(bookID)
		chapterOIDs := make([]primitive.ObjectID, len(chapterIDs))
		for i, id := range chapterIDs {
			chapterOIDs[i], _ = primitive.ObjectIDFromHex(id)
		}

		batch = &bookstore.ChapterPurchaseBatch{
			UserID:        userOID,
			BookID:        bookOID,
			ChapterIDs:    chapterOIDs,
			TotalPrice:    totalPrice,
			ChaptersCount: len(chapters),
			BookTitle:     book.Title,
			BookCover:     book.Cover,
		}
		batch.BeforeCreate()

		if err := s.purchaseRepo.CreateBatch(txCtx, batch); err != nil {
			return fmt.Errorf("failed to create batch purchase record: %w", err)
		}

		// 为每个章节创建单独的购买记录
		for _, chapter := range chapters {
			chapterOID, _ := primitive.ObjectIDFromHex(chapter.ID)
			bookOID, _ := primitive.ObjectIDFromHex(chapter.BookID)
			purchase := &bookstore.ChapterPurchase{
				UserID:       userOID,
				ChapterID:    chapterOID,
				BookID:       bookOID,
				Price:        chapter.Price,
				PurchaseTime: time.Now(),
				ChapterTitle: chapter.Title,
				ChapterNum:   chapter.ChapterNum,
				BookTitle:    book.Title,
				BookCover:    book.Cover,
			}
			purchase.BeforeCreate()

			if err := s.purchaseRepo.Create(txCtx, purchase); err != nil {
				return fmt.Errorf("failed to create purchase record: %w", err)
			}
			purchasedChapterIDs = append(purchasedChapterIDs, chapter.ID)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 清除缓存
	if s.cacheService != nil {
		s.cacheService.InvalidateBookChaptersCache(ctx, bookID)
	}

	return batch, nil
}

// PurchaseBook 购买全书
func (s *ChapterPurchaseServiceImpl) PurchaseBook(ctx context.Context, userID, bookID string) (*bookstore.BookPurchase, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}
	if bookID == "" {
		return nil, errors.New("book ID cannot be empty")
	}

	// 检查是否已购买全书
	existingPurchase, err := s.purchaseRepo.GetBookPurchaseByUserAndBook(ctx, userID, bookID)
	if err == nil && existingPurchase != nil {
		return nil, errors.New("book already purchased")
	}

	// 获取书籍信息
	book, err := s.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get book: %w", err)
	}
	if book == nil {
		return nil, errors.New("book not found")
	}

	// 计算全书价格
	originalPrice, discountedPrice, err := s.CalculateBookPrice(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate book price: %w", err)
	}

	// 检查用户余额
	balance, err := s.walletService.GetBalance(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user balance: %w", err)
	}
	if balance < discountedPrice {
		return nil, errors.New("insufficient balance")
	}

	// 获取所有付费章节
	chapters, err := s.chapterRepo.GetPaidChapters(ctx, bookID, 10000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get paid chapters: %w", err)
	}

	// 使用事务处理购买
	var purchase *bookstore.BookPurchase
	chapterIDs := make([]string, 0, len(chapters))

	err = s.purchaseRepo.Transaction(ctx, func(txCtx context.Context) error {
		// 扣除用户余额
		_, err := s.walletService.Consume(ctx, userID, discountedPrice, fmt.Sprintf("购买全书: %s", book.Title))
		if err != nil {
			return fmt.Errorf("failed to deduct balance: %w", err)
		}

		// 创建全书购买记录 - 需要转换类型
		userOID, _ := primitive.ObjectIDFromHex(userID)
		bookOID, _ := primitive.ObjectIDFromHex(bookID)
		purchase = &bookstore.BookPurchase{
			UserID:        userOID,
			BookID:        bookOID,
			TotalPrice:    discountedPrice,
			OriginalPrice: originalPrice,
			Discount:      1 - float64(discountedPrice)/float64(originalPrice),
			BookTitle:     book.Title,
			BookCover:     book.Cover,
			ChapterCount:  len(chapters),
		}
		purchase.BeforeCreate()

		if err := s.purchaseRepo.CreateBookPurchase(txCtx, purchase); err != nil {
			return fmt.Errorf("failed to create book purchase record: %w", err)
		}

		// 为每个付费章节创建购买记录
		for _, chapter := range chapters {
			chapterOID, _ := primitive.ObjectIDFromHex(chapter.ID)
			chapterPurchase := &bookstore.ChapterPurchase{
				UserID:       userOID,
				ChapterID:    chapterOID,
				BookID:       bookOID,
				Price:        0, // 全书购买后，章节单价为0
				PurchaseTime: time.Now(),
				ChapterTitle: chapter.Title,
				ChapterNum:   chapter.ChapterNum,
				BookTitle:    book.Title,
				BookCover:    book.Cover,
			}
			chapterPurchase.BeforeCreate()

			if err := s.purchaseRepo.Create(txCtx, chapterPurchase); err != nil {
				return fmt.Errorf("failed to create chapter purchase record: %w", err)
			}
			chapterIDs = append(chapterIDs, chapter.ID)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 清除缓存
	if s.cacheService != nil {
		s.cacheService.InvalidateBookDetailCache(ctx, bookID)
		s.cacheService.InvalidateBookChaptersCache(ctx, bookID)
	}

	return purchase, nil
}

// GetChapterPurchases 获取章节购买记录
func (s *ChapterPurchaseServiceImpl) GetChapterPurchases(ctx context.Context, userID string, page, pageSize int) ([]*bookstore.ChapterPurchase, int64, error) {
	if userID == "" {
		return nil, 0, errors.New("user ID cannot be empty")
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	purchases, total, err := s.purchaseRepo.GetByUser(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get chapter purchases: %w", err)
	}

	return purchases, total, nil
}

// GetBookPurchases 获取某本书的购买记录
func (s *ChapterPurchaseServiceImpl) GetBookPurchases(ctx context.Context, userID, bookID string, page, pageSize int) ([]*bookstore.ChapterPurchase, int64, error) {
	if userID == "" {
		return nil, 0, errors.New("user ID cannot be empty")
	}
	if bookID == "" {
		return nil, 0, errors.New("book ID cannot be empty")
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	purchases, total, err := s.purchaseRepo.GetByUserAndBook(ctx, userID, bookID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get book purchases: %w", err)
	}

	return purchases, total, nil
}

// GetAllPurchases 获取所有购买记录（包括单章、批量、全书）
func (s *ChapterPurchaseServiceImpl) GetAllPurchases(ctx context.Context, userID string, page, pageSize int) (map[string]interface{}, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 获取单章购买记录
	chapterPurchases, chapterTotal, err := s.purchaseRepo.GetByUser(ctx, userID, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get chapter purchases: %w", err)
	}

	// 获取批量购买记录
	batchPurchases, batchTotal, err := s.purchaseRepo.GetBatchesByUser(ctx, userID, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get batch purchases: %w", err)
	}

	// 获取全书购买记录
	bookPurchases, bookTotal, err := s.purchaseRepo.GetBookPurchasesByUser(ctx, userID, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get book purchases: %w", err)
	}

	result := map[string]interface{}{
		"chapter_purchases": chapterPurchases,
		"chapter_total":     chapterTotal,
		"batch_purchases":   batchPurchases,
		"batch_total":       batchTotal,
		"book_purchases":    bookPurchases,
		"book_total":        bookTotal,
		"total_spent":       0,
	}

	// 计算总消费
	totalSpent, _ := s.purchaseRepo.GetTotalSpentByUser(ctx, userID)
	result["total_spent"] = totalSpent

	return result, nil
}

// CheckChapterAccess 检查章节访问权限
func (s *ChapterPurchaseServiceImpl) CheckChapterAccess(ctx context.Context, userID, chapterID string) (*bookstore.ChapterAccessInfo, error) {
	if chapterID == "" {
		return nil, errors.New("chapter ID cannot be empty")
	}

	// 获取章节信息
	chapter, err := s.chapterRepo.GetByID(ctx, chapterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chapter: %w", err)
	}
	if chapter == nil {
		return nil, errors.New("chapter not found")
	}

	chapterOID, _ := primitive.ObjectIDFromHex(chapter.ID)
	accessInfo := &bookstore.ChapterAccessInfo{
		ChapterID:   chapterOID,
		Title:       chapter.Title,
		ChapterNum:  chapter.ChapterNum,
		WordCount:   chapter.WordCount,
		IsFree:      chapter.IsFree,
		Price:       chapter.Price,
		IsPurchased: false,
		IsVIP:       false,
		CanAccess:   false,
	}

	// 检查是否为免费章节
	if chapter.IsFree {
		accessInfo.CanAccess = true
		accessInfo.AccessReason = "free"
		return accessInfo, nil
	}

	// 检查用户是否已购买
	if userID != "" {
		purchased, err := s.purchaseRepo.CheckUserPurchasedChapter(ctx, userID, chapterID)
		if err == nil && purchased {
			accessInfo.IsPurchased = true
			accessInfo.CanAccess = true
			accessInfo.AccessReason = "purchased"

			// 获取购买时间
			purchaseRecord, _ := s.purchaseRepo.GetByUserAndChapter(ctx, userID, chapterID)
			if purchaseRecord != nil {
				accessInfo.PurchaseTime = &purchaseRecord.PurchaseTime
			}
			return accessInfo, nil
		}

		// 检查是否已购买全书
		bookPurchased, err := s.purchaseRepo.CheckUserPurchasedBook(ctx, userID, chapter.BookID)
		if err == nil && bookPurchased {
			accessInfo.IsPurchased = true
			accessInfo.CanAccess = true
			accessInfo.AccessReason = "purchased_book"
			return accessInfo, nil
		}
	}

	// 检查是否为VIP用户（待实现）
	// isVIP, _ := s.IsVIPUser(ctx, userID)
	// if isVIP {
	// 	accessInfo.IsVIP = true
	// 	accessInfo.CanAccess = true
	// 	accessInfo.AccessReason = "vip"
	// 	return accessInfo, nil
	// }

	return accessInfo, nil
}

// GetPurchasedChapterIDs 获取已购买的章节ID列表
func (s *ChapterPurchaseServiceImpl) GetPurchasedChapterIDs(ctx context.Context, userID, bookID string) ([]string, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}

	chapterIDs, err := s.purchaseRepo.GetPurchasedChapterIDs(ctx, userID, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get purchased chapter IDs: %w", err)
	}

	return chapterIDs, nil
}

// GetChapterPrice 获取章节价格
func (s *ChapterPurchaseServiceImpl) GetChapterPrice(ctx context.Context, chapterID string) (float64, error) {
	if chapterID == "" {
		return 0, errors.New("chapter ID cannot be empty")
	}

	chapter, err := s.chapterRepo.GetByID(ctx, chapterID)
	if err != nil {
		return 0, fmt.Errorf("failed to get chapter: %w", err)
	}
	if chapter == nil {
		return 0, errors.New("chapter not found")
	}

	// 转换为元 (除以100)
	return float64(chapter.Price) / 100.0, nil
}

// CalculateBookPrice 计算全书价格
func (s *ChapterPurchaseServiceImpl) CalculateBookPrice(ctx context.Context, bookID string) (int64, int64, error) {
	if bookID == "" {
		return 0, 0, errors.New("book ID cannot be empty")
	}

	// 获取所有付费章节
	chapters, err := s.chapterRepo.GetPaidChapters(ctx, bookID, 10000, 0)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get paid chapters: %w", err)
	}

	// 计算原价 (分)
	originalPrice := int64(0)
	for _, chapter := range chapters {
		originalPrice += chapter.Price
	}

	// 计算折扣价（全书购买通常有折扣，例如8折）
	discount := 0.8
	discountedPrice := int64(float64(originalPrice) * discount)

	return originalPrice, discountedPrice, nil
}

// IsVIPUser 检查是否为VIP用户
func (s *ChapterPurchaseServiceImpl) IsVIPUser(ctx context.Context, userID string) (bool, error) {
	// TODO: 实现VIP用户检查逻辑
	// 这里需要与用户系统集成，检查用户的VIP状态和有效期
	return false, nil
}
