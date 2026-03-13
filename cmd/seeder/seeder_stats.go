// Package main 提供统计数据填充功能
package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	"Qingyu_backend/cmd/seeder/config"
	"Qingyu_backend/cmd/seeder/utils"
	statsModel "Qingyu_backend/models/stats"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type statsBookSeed struct {
	ID          string
	Title       string
	AuthorID    string
	PublishedAt time.Time
	WordCount   int64
}

type statsChapterSeed struct {
	ID         string
	BookID     string
	Title      string
	ChapterNum int
	WordCount  int
	Price      float64
	IsFree     bool
}

type statsUserSeed struct {
	ID    string
	Roles []string
}

type bookDailyAccumulator struct {
	views        int64
	newReaders   map[string]struct{}
	subscribers  int64
	dailyRevenue float64
}

type chapterBehaviorAccumulator struct {
	viewCount      int64
	uniqueReaders  map[string]struct{}
	totalReadTime  int64
	readTimeEvents int64
	completeCount  int64
	dropOffCount   int64
}

// StatsSeeder 统计数据填充器
type StatsSeeder struct {
	db     *utils.Database
	config *config.Config
	rng    *rand.Rand
}

// NewStatsSeeder 创建统计数据填充器
func NewStatsSeeder(db *utils.Database, cfg *config.Config) *StatsSeeder {
	return &StatsSeeder{
		db:     db,
		config: cfg,
		rng:    rand.New(rand.NewSource(20260311)),
	}
}

// SeedStats 填充统计数据
func (s *StatsSeeder) SeedStats() error {
	ctx := context.Background()

	books, err := s.getBooks(ctx)
	if err != nil {
		return fmt.Errorf("获取书籍列表失败: %w", err)
	}
	if len(books) == 0 {
		fmt.Println("  没有找到书籍，请先运行 bookstore 命令创建书籍")
		return nil
	}

	chaptersByBook, err := s.getChaptersByBook(ctx)
	if err != nil {
		return fmt.Errorf("获取章节列表失败: %w", err)
	}
	if len(chaptersByBook) == 0 {
		fmt.Println("  没有找到章节，请先运行 chapters 命令创建章节")
		return nil
	}

	readers, err := s.getReaders(ctx)
	if err != nil {
		return fmt.Errorf("获取用户列表失败: %w", err)
	}
	if len(readers) == 0 {
		fmt.Println("  没有找到读者用户，请先运行 users 命令创建用户")
		return nil
	}

	behaviorCount, err := s.db.Collection("reader_behaviors").CountDocuments(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("检查 reader_behaviors 集合失败: %w", err)
	}
	if behaviorCount == 0 {
		if err := s.seedReaderBehaviors(ctx, readers, books, chaptersByBook); err != nil {
			return err
		}
	}

	if err := s.resetDerivedCollections(ctx); err != nil {
		return err
	}

	subscriptionEvents, subscriptionCounts, err := s.getSubscriptions(ctx)
	if err != nil {
		return fmt.Errorf("获取订阅数据失败: %w", err)
	}

	behaviors, err := s.getReaderBehaviors(ctx)
	if err != nil {
		return fmt.Errorf("获取读者行为失败: %w", err)
	}

	commentCounts, err := s.getGroupedCounts(ctx, "comments", "target_id", bson.M{"target_type": "book"})
	if err != nil {
		return fmt.Errorf("获取评论统计失败: %w", err)
	}
	likeCounts, err := s.getGroupedCounts(ctx, "likes", "target_id", bson.M{"target_type": "book"})
	if err != nil {
		return fmt.Errorf("获取点赞统计失败: %w", err)
	}
	collectionCounts, err := s.getGroupedCounts(ctx, "collections", "book_id", bson.M{})
	if err != nil {
		return fmt.Errorf("获取收藏统计失败: %w", err)
	}
	bookmarkCounts, err := s.getGroupedCounts(ctx, "bookmarks", "book_id", bson.M{})
	if err != nil {
		return fmt.Errorf("获取书签统计失败: %w", err)
	}

	dailyDocs, chapterDocs, bookDocs, retentionDocs := s.buildStatsDocuments(
		books,
		chaptersByBook,
		behaviors,
		subscriptionEvents,
		subscriptionCounts,
		commentCounts,
		likeCounts,
		collectionCounts,
		bookmarkCounts,
	)

	if len(dailyDocs) > 0 {
		if _, err := s.db.Collection("book_stats_daily").InsertMany(ctx, toDocSlice(dailyDocs)); err != nil {
			return fmt.Errorf("插入每日统计失败: %w", err)
		}
	}
	if len(chapterDocs) > 0 {
		if _, err := s.db.Collection("chapter_stats").InsertMany(ctx, toDocSlice(chapterDocs)); err != nil {
			return fmt.Errorf("插入章节统计失败: %w", err)
		}
	}
	if len(bookDocs) > 0 {
		if _, err := s.db.Collection("book_stats").InsertMany(ctx, toDocSlice(bookDocs)); err != nil {
			return fmt.Errorf("插入作品统计失败: %w", err)
		}
	}
	if len(retentionDocs) > 0 {
		if _, err := s.db.Collection("reader_retentions").InsertMany(ctx, toDocSlice(retentionDocs)); err != nil {
			return fmt.Errorf("插入留存统计失败: %w", err)
		}
	}

	fmt.Printf("  创建了 %d 条每日统计记录\n", len(dailyDocs))
	fmt.Printf("  创建了 %d 条章节统计记录\n", len(chapterDocs))
	fmt.Printf("  创建了 %d 条书籍统计记录\n", len(bookDocs))
	fmt.Printf("  创建了 %d 条留存统计记录\n", len(retentionDocs))
	return nil
}

func (s *StatsSeeder) resetDerivedCollections(ctx context.Context) error {
	collections := []string{"book_stats", "book_stats_daily", "chapter_stats", "reader_retentions"}
	for _, collName := range collections {
		if _, err := s.db.Collection(collName).DeleteMany(ctx, bson.M{}); err != nil {
			return fmt.Errorf("清空 %s 集合失败: %w", collName, err)
		}
	}
	return nil
}

func (s *StatsSeeder) seedReaderBehaviors(
	ctx context.Context,
	readers []statsUserSeed,
	books []statsBookSeed,
	chaptersByBook map[string][]statsChapterSeed,
) error {
	collection := s.db.Collection("reader_behaviors")
	now := time.Now()
	behaviors := make([]interface{}, 0)

	for _, book := range books {
		chapters := chaptersByBook[book.ID]
		if len(chapters) == 0 {
			continue
		}

		readerCount := s.pickReaderCount(len(readers))
		selectedReaders := s.pickUsers(readers, readerCount)

		for _, reader := range selectedReaders {
			sessionCount := 1 + s.rng.Intn(minInt(4, len(chapters)))
			for session := 0; session < sessionCount; session++ {
				chapterIndex := minInt(len(chapters)-1, session+s.rng.Intn(minInt(3, len(chapters))))
				chapter := chapters[chapterIndex]
				readAt := s.randomRecentTime(now, 34)
				progress := 0.35 + s.rng.Float64()*0.65
				readDuration := 180 + s.rng.Intn(1800)
				endPosition := int(float64(maxInt(chapter.WordCount, 1)) * progress)

				behaviors = append(behaviors, statsModel.ReaderBehavior{
					ID:            primitive.NewObjectID().Hex(),
					UserID:        reader.ID,
					BookID:        book.ID,
					ChapterID:     chapter.ID,
					BehaviorType:  statsModel.BehaviorTypeView,
					StartPosition: 0,
					EndPosition:   endPosition,
					Progress:      progress,
					ReadDuration:  readDuration,
					ReadAt:        readAt,
					DeviceType:    s.randomDeviceType(),
					ClientIP:      s.randomIP(),
					Source:        s.randomSource(),
					Referrer:      "/bookstore/books/" + book.ID,
					CreatedAt:     readAt,
				})

				if progress >= 0.82 {
					behaviors = append(behaviors, statsModel.ReaderBehavior{
						ID:            primitive.NewObjectID().Hex(),
						UserID:        reader.ID,
						BookID:        book.ID,
						ChapterID:     chapter.ID,
						BehaviorType:  statsModel.BehaviorTypeComplete,
						StartPosition: 0,
						EndPosition:   chapter.WordCount,
						Progress:      1,
						ReadDuration:  readDuration + 90 + s.rng.Intn(420),
						ReadAt:        readAt.Add(time.Duration(10+s.rng.Intn(180)) * time.Second),
						DeviceType:    s.randomDeviceType(),
						ClientIP:      s.randomIP(),
						Source:        statsModel.SourceBookshelf,
						Referrer:      "/reader/books/" + book.ID,
						CreatedAt:     readAt,
					})
				} else if progress <= 0.55 {
					behaviors = append(behaviors, statsModel.ReaderBehavior{
						ID:            primitive.NewObjectID().Hex(),
						UserID:        reader.ID,
						BookID:        book.ID,
						ChapterID:     chapter.ID,
						BehaviorType:  statsModel.BehaviorTypeDropOff,
						StartPosition: 0,
						EndPosition:   endPosition,
						Progress:      progress,
						ReadDuration:  readDuration,
						ReadAt:        readAt.Add(time.Duration(5+s.rng.Intn(120)) * time.Second),
						DeviceType:    s.randomDeviceType(),
						ClientIP:      s.randomIP(),
						Source:        statsModel.SourceRecommendation,
						Referrer:      "/bookstore/home",
						CreatedAt:     readAt,
					})
				}
			}
		}

		s.appendRetentionAnchors(&behaviors, selectedReaders, book, chapters, now)
	}

	if len(behaviors) == 0 {
		return nil
	}

	if _, err := collection.InsertMany(ctx, behaviors); err != nil {
		return fmt.Errorf("插入 reader_behaviors 失败: %w", err)
	}

	fmt.Printf("  创建了 %d 条读者行为记录\n", len(behaviors))
	return nil
}

func (s *StatsSeeder) appendRetentionAnchors(
	behaviors *[]interface{},
	readers []statsUserSeed,
	book statsBookSeed,
	chapters []statsChapterSeed,
	now time.Time,
) {
	if len(readers) == 0 || len(chapters) == 0 {
		return
	}
	anchorDays := []int{30, 7, 1, 0}
	chapter := chapters[0]
	limit := minInt(len(readers), 6)

	for i := 0; i < limit; i++ {
		reader := readers[i]
		for _, dayOffset := range anchorDays {
			readAt := now.AddDate(0, 0, -dayOffset).Add(time.Duration(9+i) * time.Hour)
			*behaviors = append(*behaviors, statsModel.ReaderBehavior{
				ID:            primitive.NewObjectID().Hex(),
				UserID:        reader.ID,
				BookID:        book.ID,
				ChapterID:     chapter.ID,
				BehaviorType:  statsModel.BehaviorTypeView,
				StartPosition: 0,
				EndPosition:   maxInt(1, chapter.WordCount/2),
				Progress:      0.5,
				ReadDuration:  360 + i*30,
				ReadAt:        readAt,
				DeviceType:    statsModel.DeviceTypeMobile,
				ClientIP:      s.randomIP(),
				Source:        statsModel.SourceBookshelf,
				Referrer:      "/bookshelf",
				CreatedAt:     readAt,
			})
		}
	}
}

func (s *StatsSeeder) buildStatsDocuments(
	books []statsBookSeed,
	chaptersByBook map[string][]statsChapterSeed,
	behaviors []statsModel.ReaderBehavior,
	subscriptionEvents map[string]map[string]int64,
	subscriptionCounts map[string]int64,
	commentCounts map[string]int64,
	likeCounts map[string]int64,
	collectionCounts map[string]int64,
	bookmarkCounts map[string]int64,
) ([]statsModel.BookStatsDaily, []statsModel.ChapterStats, []statsModel.BookStats, []statsModel.ReaderRetention) {
	now := time.Now()
	startDate := now.AddDate(0, 0, -29).Truncate(24 * time.Hour)
	endDate := now.Truncate(24 * time.Hour)

	bookDaily := make(map[string]map[string]*bookDailyAccumulator)
	bookReadersByDay := make(map[string]map[string]map[string]struct{})
	bookAllReaders := make(map[string]map[string]struct{})
	chapterAgg := make(map[string]*chapterBehaviorAccumulator)
	bookDurationSum := make(map[string]int64)
	bookDurationCount := make(map[string]int64)
	bookFirstSeen := make(map[string]map[string]string)

	for _, behavior := range behaviors {
		dateKey := dayKey(behavior.ReadAt)
		if _, ok := bookDaily[behavior.BookID]; !ok {
			bookDaily[behavior.BookID] = make(map[string]*bookDailyAccumulator)
		}
		if _, ok := bookReadersByDay[behavior.BookID]; !ok {
			bookReadersByDay[behavior.BookID] = make(map[string]map[string]struct{})
		}
		if _, ok := bookReadersByDay[behavior.BookID][dateKey]; !ok {
			bookReadersByDay[behavior.BookID][dateKey] = make(map[string]struct{})
		}
		if _, ok := bookAllReaders[behavior.BookID]; !ok {
			bookAllReaders[behavior.BookID] = make(map[string]struct{})
		}
		if _, ok := bookFirstSeen[behavior.BookID]; !ok {
			bookFirstSeen[behavior.BookID] = make(map[string]string)
		}
		if _, ok := chapterAgg[behavior.ChapterID]; !ok {
			chapterAgg[behavior.ChapterID] = &chapterBehaviorAccumulator{
				uniqueReaders: make(map[string]struct{}),
			}
		}
		if _, ok := bookDaily[behavior.BookID][dateKey]; !ok {
			bookDaily[behavior.BookID][dateKey] = &bookDailyAccumulator{
				newReaders: make(map[string]struct{}),
			}
		}

		bookAllReaders[behavior.BookID][behavior.UserID] = struct{}{}
		bookReadersByDay[behavior.BookID][dateKey][behavior.UserID] = struct{}{}

		if firstDate, exists := bookFirstSeen[behavior.BookID][behavior.UserID]; !exists || dateKey < firstDate {
			bookFirstSeen[behavior.BookID][behavior.UserID] = dateKey
		}

		agg := chapterAgg[behavior.ChapterID]
		switch behavior.BehaviorType {
		case statsModel.BehaviorTypeView:
			agg.viewCount++
			agg.uniqueReaders[behavior.UserID] = struct{}{}
			agg.totalReadTime += int64(behavior.ReadDuration)
			agg.readTimeEvents++
			bookDaily[behavior.BookID][dateKey].views++
			bookDurationSum[behavior.BookID] += int64(behavior.ReadDuration)
			bookDurationCount[behavior.BookID]++
		case statsModel.BehaviorTypeComplete:
			agg.completeCount++
			agg.uniqueReaders[behavior.UserID] = struct{}{}
		case statsModel.BehaviorTypeDropOff:
			agg.dropOffCount++
			agg.uniqueReaders[behavior.UserID] = struct{}{}
		}
	}

	for bookID, readers := range bookFirstSeen {
		for userID, firstDate := range readers {
			if dayStats, ok := bookDaily[bookID][firstDate]; ok {
				dayStats.newReaders[userID] = struct{}{}
			}
		}
	}

	for bookID, dailyCounts := range subscriptionEvents {
		if _, ok := bookDaily[bookID]; !ok {
			bookDaily[bookID] = make(map[string]*bookDailyAccumulator)
		}
		for dateKey, count := range dailyCounts {
			if _, ok := bookDaily[bookID][dateKey]; !ok {
				bookDaily[bookID][dateKey] = &bookDailyAccumulator{
					newReaders: make(map[string]struct{}),
				}
			}
			bookDaily[bookID][dateKey].subscribers += count
		}
	}

	dailyDocs := make([]statsModel.BookStatsDaily, 0, len(books)*30)
	chapterDocs := make([]statsModel.ChapterStats, 0)
	bookDocs := make([]statsModel.BookStats, 0, len(books))
	retentionDocs := make([]statsModel.ReaderRetention, 0, len(books))

	for _, book := range books {
		chapters := chaptersByBook[book.ID]
		if len(chapters) == 0 {
			continue
		}

		totalSubscribers := subscriptionCounts[book.ID]
		totalComments := commentCounts[book.ID]
		totalLikes := likeCounts[book.ID]
		totalBookmarks := collectionCounts[book.ID] + bookmarkCounts[book.ID]
		totalShares := totalLikes / 12
		if totalLikes > 0 && totalShares == 0 {
			totalShares = 1
		}

		totalViews := int64(0)
		totalDropOffs := int64(0)
		totalChapterRevenue := 0.0
		totalChapterViews := float64(0)
		totalCompletionRate := 0.0
		totalDropOffRate := 0.0
		highestDropOffRate := -1.0
		dropOffChapterTitle := ""

		for _, chapter := range chapters {
			agg := chapterAgg[chapter.ID]
			if agg == nil {
				agg = &chapterBehaviorAccumulator{uniqueReaders: make(map[string]struct{})}
			}

			avgReadTime := 0.0
			if agg.readTimeEvents > 0 {
				avgReadTime = float64(agg.totalReadTime) / float64(agg.readTimeEvents)
			}

			completionRate := 0.0
			if agg.viewCount > 0 {
				completionRate = float64(agg.completeCount) / float64(agg.viewCount)
			}

			dropOffRate := 0.0
			if agg.viewCount > 0 {
				dropOffRate = float64(agg.dropOffCount) / float64(agg.viewCount)
			}

			chapterSubscriberCount := int64(0)
			if !chapter.IsFree {
				chapterSubscriberCount = minInt64(totalSubscribers, maxInt64(agg.viewCount/2, 1))
			}

			revenue := 0.0
			if !chapter.IsFree && chapterSubscriberCount > 0 {
				revenue = float64(chapterSubscriberCount) * math.Max(chapter.Price, 0.05)
			}

			chapterDocs = append(chapterDocs, statsModel.ChapterStats{
				ID:             primitive.NewObjectID().Hex(),
				BookID:         book.ID,
				ChapterID:      chapter.ID,
				Title:          chapter.Title,
				WordCount:      chapter.WordCount,
				ViewCount:      agg.viewCount,
				UniqueViewers:  int64(len(agg.uniqueReaders)),
				AvgReadTime:    avgReadTime,
				CompletionRate: roundFloat(completionRate, 4),
				DropOffCount:   agg.dropOffCount,
				DropOffRate:    roundFloat(dropOffRate, 4),
				CommentCount:   0,
				LikeCount:      agg.viewCount / 6,
				BookmarkCount:  agg.viewCount / 10,
				SubscribeCount: chapterSubscriberCount,
				Revenue:        roundFloat(revenue, 2),
				StatDate:       endDate,
				CreatedAt:      now,
				UpdatedAt:      now,
			})

			totalViews += agg.viewCount
			totalDropOffs += agg.dropOffCount
			totalChapterRevenue += revenue
			totalChapterViews += float64(agg.viewCount)
			totalCompletionRate += completionRate
			totalDropOffRate += dropOffRate
			if dropOffRate > highestDropOffRate {
				highestDropOffRate = dropOffRate
				dropOffChapterTitle = chapter.Title
			}
		}

		bookDailySeries := make([]statsModel.BookStatsDaily, 0, 30)
		for current := startDate; !current.After(endDate); current = current.Add(24 * time.Hour) {
			dateKey := dayKey(current)
			acc := bookDaily[book.ID][dateKey]
			if acc == nil {
				acc = &bookDailyAccumulator{newReaders: make(map[string]struct{})}
			}

			dailyRevenue := float64(acc.subscribers) * s.estimateDailyRevenue(chapters)
			acc.dailyRevenue = dailyRevenue

			bookDailySeries = append(bookDailySeries, statsModel.BookStatsDaily{
				ID:               primitive.NewObjectID().Hex(),
				BookID:           book.ID,
				Date:             current,
				DailyViews:       acc.views,
				DailyNewReaders:  int64(len(acc.newReaders)),
				DailyRevenue:     roundFloat(acc.dailyRevenue, 2),
				DailySubscribers: acc.subscribers,
				CreatedAt:        now,
				UpdatedAt:        now,
			})
		}
		dailyDocs = append(dailyDocs, bookDailySeries...)

		totalWords := int64(0)
		for _, chapter := range chapters {
			totalWords += int64(chapter.WordCount)
		}
		if totalWords == 0 {
			totalWords = book.WordCount
		}

		avgChapterViews := 0.0
		if len(chapters) > 0 {
			avgChapterViews = totalChapterViews / float64(len(chapters))
		}

		avgCompletionRate := 0.0
		avgDropOffRate := 0.0
		if len(chapters) > 0 {
			avgCompletionRate = totalCompletionRate / float64(len(chapters))
			avgDropOffRate = totalDropOffRate / float64(len(chapters))
		}

		avgReadingDuration := 0.0
		if bookDurationCount[book.ID] > 0 {
			avgReadingDuration = float64(bookDurationSum[book.ID]) / float64(bookDurationCount[book.ID])
		}

		rewardRevenue := float64(totalLikes) * 0.08
		totalRevenue := totalChapterRevenue + rewardRevenue
		avgRevenuePerUser := 0.0
		if len(bookAllReaders[book.ID]) > 0 {
			avgRevenuePerUser = totalRevenue / float64(len(bookAllReaders[book.ID]))
		}

		day1Retention := s.calculateRetention(bookReadersByDay[book.ID], 1)
		day3Retention := s.calculateRetention(bookReadersByDay[book.ID], 3)
		day7Retention := s.calculateRetention(bookReadersByDay[book.ID], 7)
		day30Retention := s.calculateRetention(bookReadersByDay[book.ID], 30)

		viewTrend := calculateTrend(bookDailySeries, func(item statsModel.BookStatsDaily) float64 {
			return float64(item.DailyViews)
		})
		revenueTrend := calculateTrend(bookDailySeries, func(item statsModel.BookStatsDaily) float64 {
			return item.DailyRevenue
		})

		bookDocs = append(bookDocs, statsModel.BookStats{
			ID:                 primitive.NewObjectID().Hex(),
			BookID:             book.ID,
			Title:              book.Title,
			AuthorID:           book.AuthorID,
			TotalChapter:       len(chapters),
			TotalWords:         totalWords,
			TotalViews:         totalViews,
			UniqueReaders:      int64(len(bookAllReaders[book.ID])),
			AvgChapterViews:    roundFloat(avgChapterViews, 2),
			AvgCompletionRate:  roundFloat(avgCompletionRate, 4),
			AvgReadingDuration: roundFloat(avgReadingDuration, 2),
			TotalDropOffs:      totalDropOffs,
			AvgDropOffRate:     roundFloat(avgDropOffRate, 4),
			DropOffChapter:     dropOffChapterTitle,
			TotalComments:      totalComments,
			TotalLikes:         totalLikes,
			TotalBookmarks:     totalBookmarks,
			TotalShares:        totalShares,
			TotalSubscribers:   totalSubscribers,
			AvgSubscribeRate:   safeRate(totalSubscribers, int64(len(bookAllReaders[book.ID]))),
			TotalRevenue:       roundFloat(totalRevenue, 2),
			ChapterRevenue:     roundFloat(totalChapterRevenue, 2),
			SubscribeRevenue:   roundFloat(totalChapterRevenue, 2),
			RewardRevenue:      roundFloat(rewardRevenue, 2),
			AvgRevenuePerUser:  roundFloat(avgRevenuePerUser, 2),
			Day1Retention:      day1Retention,
			Day7Retention:      day7Retention,
			Day30Retention:     day30Retention,
			ViewTrend:          viewTrend,
			RevenueTrend:       revenueTrend,
			StatDate:           endDate,
			CreatedAt:          now,
			UpdatedAt:          now,
		})

		retentionDocs = append(retentionDocs, statsModel.ReaderRetention{
			BookID:         book.ID,
			Day1Retention:  day1Retention,
			Day3Retention:  day3Retention,
			Day7Retention:  day7Retention,
			Day30Retention: day30Retention,
			NewReaders:     sumNewReaders(bookDailySeries),
			ActiveReaders:  int64(len(bookReadersByDay[book.ID][dayKey(endDate)])),
			StatDate:       endDate,
			CreatedAt:      now,
		})
	}

	return dailyDocs, chapterDocs, bookDocs, retentionDocs
}

func (s *StatsSeeder) estimateDailyRevenue(chapters []statsChapterSeed) float64 {
	total := 0.0
	for _, chapter := range chapters {
		if chapter.IsFree {
			continue
		}
		total += math.Max(chapter.Price, 0.05)
	}
	if total == 0 {
		return 0
	}
	return roundFloat(total/float64(len(chapters)+1), 2)
}

func (s *StatsSeeder) calculateRetention(readersByDay map[string]map[string]struct{}, days int) float64 {
	if len(readersByDay) == 0 {
		return 0
	}
	baseDay := dayKey(time.Now().AddDate(0, 0, -days))
	today := dayKey(time.Now())
	baseReaders := readersByDay[baseDay]
	if len(baseReaders) == 0 {
		return 0
	}
	todayReaders := readersByDay[today]
	active := int64(0)
	for userID := range baseReaders {
		if _, ok := todayReaders[userID]; ok {
			active++
		}
	}
	return safeRate(active, int64(len(baseReaders)))
}

func (s *StatsSeeder) pickReaderCount(total int) int {
	if total <= 0 {
		return 0
	}
	scaleCfg := config.GetScaleConfig(s.config.Scale)
	base := maxInt(8, scaleCfg.Users/25)
	count := base + s.rng.Intn(maxInt(2, base/2))
	return minInt(total, maxInt(6, count))
}

func (s *StatsSeeder) pickUsers(users []statsUserSeed, count int) []statsUserSeed {
	if count >= len(users) {
		cloned := make([]statsUserSeed, len(users))
		copy(cloned, users)
		return cloned
	}

	indices := s.rng.Perm(len(users))
	selected := make([]statsUserSeed, 0, count)
	for i := 0; i < count; i++ {
		selected = append(selected, users[indices[i]])
	}
	return selected
}

func (s *StatsSeeder) randomRecentTime(now time.Time, maxDays int) time.Time {
	dayOffset := s.rng.Intn(maxDays + 1)
	base := now.AddDate(0, 0, -dayOffset).Truncate(24 * time.Hour)
	hour := 7 + s.rng.Intn(16)
	minute := s.rng.Intn(60)
	second := s.rng.Intn(60)
	return base.Add(time.Duration(hour)*time.Hour + time.Duration(minute)*time.Minute + time.Duration(second)*time.Second)
}

func (s *StatsSeeder) randomDeviceType() string {
	devices := []string{statsModel.DeviceTypeMobile, statsModel.DeviceTypeDesktop, statsModel.DeviceTypeTablet}
	return devices[s.rng.Intn(len(devices))]
}

func (s *StatsSeeder) randomSource() string {
	sources := []string{
		statsModel.SourceRecommendation,
		statsModel.SourceSearch,
		statsModel.SourceBookshelf,
		statsModel.SourceRanking,
		statsModel.SourceCategory,
	}
	return sources[s.rng.Intn(len(sources))]
}

func (s *StatsSeeder) randomIP() string {
	return fmt.Sprintf("10.%d.%d.%d", s.rng.Intn(255), s.rng.Intn(255), s.rng.Intn(255))
}

func (s *StatsSeeder) getBooks(ctx context.Context) ([]statsBookSeed, error) {
	cursor, err := s.db.Collection("books").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []bson.M
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	books := make([]statsBookSeed, 0, len(docs))
	for _, doc := range docs {
		books = append(books, statsBookSeed{
			ID:          normalizeID(doc["_id"]),
			Title:       stringValue(doc["title"]),
			AuthorID:    normalizeID(doc["author_id"]),
			PublishedAt: timeValue(doc["published_at"]),
			WordCount:   int64Value(doc["word_count"]),
		})
	}
	return books, nil
}

func (s *StatsSeeder) getChaptersByBook(ctx context.Context) (map[string][]statsChapterSeed, error) {
	cursor, err := s.db.Collection("chapters").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []bson.M
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	chaptersByBook := make(map[string][]statsChapterSeed)
	for _, doc := range docs {
		bookID := normalizeID(doc["book_id"])
		if bookID == "" {
			continue
		}
		chaptersByBook[bookID] = append(chaptersByBook[bookID], statsChapterSeed{
			ID:         normalizeID(doc["_id"]),
			BookID:     bookID,
			Title:      stringValue(doc["title"]),
			ChapterNum: intValue(doc["chapter_num"]),
			WordCount:  intValue(doc["word_count"]),
			Price:      floatValue(doc["price"]),
			IsFree:     boolValue(doc["is_free"]),
		})
	}

	for bookID := range chaptersByBook {
		sort.SliceStable(chaptersByBook[bookID], func(i, j int) bool {
			return chaptersByBook[bookID][i].ChapterNum < chaptersByBook[bookID][j].ChapterNum
		})
	}

	return chaptersByBook, nil
}

func (s *StatsSeeder) getReaders(ctx context.Context) ([]statsUserSeed, error) {
	cursor, err := s.db.Collection("users").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []bson.M
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	readers := make([]statsUserSeed, 0, len(docs))
	for _, doc := range docs {
		roles := stringSliceValue(doc["roles"])
		if len(roles) == 0 {
			if role := stringValue(doc["role"]); role != "" {
				roles = []string{role}
			}
		}
		if !containsRole(roles, "reader") && !containsRole(roles, "author") && !containsRole(roles, "admin") {
			continue
		}
		readers = append(readers, statsUserSeed{
			ID:    normalizeID(doc["_id"]),
			Roles: roles,
		})
	}
	return readers, nil
}

func (s *StatsSeeder) getSubscriptions(ctx context.Context) (map[string]map[string]int64, map[string]int64, error) {
	cursor, err := s.db.Collection("subscriptions").Find(ctx, bson.M{})
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	var docs []bson.M
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, nil, err
	}

	daily := make(map[string]map[string]int64)
	totals := make(map[string]int64)
	for _, doc := range docs {
		bookID := normalizeID(doc["book_id"])
		if bookID == "" {
			continue
		}
		dateKey := dayKey(timeValue(doc["subscribed_at"]))
		if _, ok := daily[bookID]; !ok {
			daily[bookID] = make(map[string]int64)
		}
		daily[bookID][dateKey]++
		totals[bookID]++
	}
	return daily, totals, nil
}

func (s *StatsSeeder) getReaderBehaviors(ctx context.Context) ([]statsModel.ReaderBehavior, error) {
	cursor, err := s.db.Collection("reader_behaviors").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var behaviors []statsModel.ReaderBehavior
	if err := cursor.All(ctx, &behaviors); err != nil {
		return nil, err
	}
	return behaviors, nil
}

func (s *StatsSeeder) getGroupedCounts(ctx context.Context, collName, field string, filter bson.M) (map[string]int64, error) {
	pipeline := mongoPipelineForGroupedCount(field, filter)
	cursor, err := s.db.Collection(collName).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rows []bson.M
	if err := cursor.All(ctx, &rows); err != nil {
		return nil, err
	}

	result := make(map[string]int64, len(rows))
	for _, row := range rows {
		result[normalizeID(row["_id"])] = int64Value(row["count"])
	}
	return result, nil
}

func mongoPipelineForGroupedCount(field string, filter bson.M) mongo.Pipeline {
	pipeline := mongo.Pipeline{}
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: filter}})
	}
	pipeline = append(pipeline, bson.D{{Key: "$group", Value: bson.M{
		"_id":   "$" + field,
		"count": bson.M{"$sum": 1},
	}}})
	return pipeline
}

// Clean 清空统计数据
func (s *StatsSeeder) Clean() error {
	ctx := context.Background()
	collections := []string{"book_stats", "book_stats_daily", "chapter_stats", "reader_behaviors", "reader_retentions"}

	for _, collName := range collections {
		if _, err := s.db.Collection(collName).DeleteMany(ctx, bson.M{}); err != nil {
			return fmt.Errorf("清空 %s 集合失败: %w", collName, err)
		}
	}

	fmt.Println("  已清空统计数据集合")
	return nil
}

func toDocSlice[T any](items []T) []interface{} {
	docs := make([]interface{}, 0, len(items))
	for _, item := range items {
		docs = append(docs, item)
	}
	return docs
}

func sumNewReaders(items []statsModel.BookStatsDaily) int64 {
	total := int64(0)
	for _, item := range items {
		total += item.DailyNewReaders
	}
	return total
}

func calculateTrend(items []statsModel.BookStatsDaily, valueFn func(statsModel.BookStatsDaily) float64) string {
	if len(items) < 2 {
		return statsModel.TrendStable
	}
	mid := len(items) / 2
	firstHalf := averageDaily(items[:mid], valueFn)
	secondHalf := averageDaily(items[mid:], valueFn)
	if secondHalf > firstHalf*1.1 {
		return statsModel.TrendUp
	}
	if secondHalf < firstHalf*0.9 {
		return statsModel.TrendDown
	}
	return statsModel.TrendStable
}

func averageDaily(items []statsModel.BookStatsDaily, valueFn func(statsModel.BookStatsDaily) float64) float64 {
	if len(items) == 0 {
		return 0
	}
	total := 0.0
	for _, item := range items {
		total += valueFn(item)
	}
	return total / float64(len(items))
}

func safeRate(numerator, denominator int64) float64 {
	if denominator == 0 {
		return 0
	}
	return roundFloat(float64(numerator)/float64(denominator), 4)
}

func roundFloat(value float64, precision int) float64 {
	pow := math.Pow10(precision)
	return math.Round(value*pow) / pow
}

func dayKey(t time.Time) string {
	return t.Truncate(24 * time.Hour).Format("2006-01-02")
}

func containsRole(roles []string, target string) bool {
	for _, role := range roles {
		if role == target {
			return true
		}
	}
	return false
}

func normalizeID(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case primitive.ObjectID:
		return v.Hex()
	case *primitive.ObjectID:
		if v == nil {
			return ""
		}
		return v.Hex()
	case map[string]interface{}:
		if raw, ok := v["$oid"]; ok {
			return stringValue(raw)
		}
	case bson.M:
		if raw, ok := v["$oid"]; ok {
			return stringValue(raw)
		}
	}
	return ""
}

func stringValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	}
	return ""
}

func stringSliceValue(value interface{}) []string {
	switch v := value.(type) {
	case []string:
		return v
	case bson.A:
		result := make([]string, 0, len(v))
		for _, item := range v {
			if s := stringValue(item); s != "" {
				result = append(result, s)
			}
		}
		return result
	case []interface{}:
		result := make([]string, 0, len(v))
		for _, item := range v {
			if s := stringValue(item); s != "" {
				result = append(result, s)
			}
		}
		return result
	}
	return nil
}

func timeValue(value interface{}) time.Time {
	switch v := value.(type) {
	case time.Time:
		return v
	case *time.Time:
		if v != nil {
			return *v
		}
	case primitive.DateTime:
		return v.Time()
	}
	return time.Now()
}

func intValue(value interface{}) int {
	return int(int64Value(value))
}

func int64Value(value interface{}) int64 {
	switch v := value.(type) {
	case int:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	case float32:
		return int64(v)
	case float64:
		return int64(v)
	}
	return 0
}

func floatValue(value interface{}) float64 {
	switch v := value.(type) {
	case float32:
		return float64(v)
	case float64:
		return v
	case int:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	}
	return 0
}

func boolValue(value interface{}) bool {
	switch v := value.(type) {
	case bool:
		return v
	}
	return false
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
