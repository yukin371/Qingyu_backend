package middleware

import (
	"context"
	"hash/crc32"
	"math/rand"
	"sync"
	"time"

	"go.uber.org/zap"

	"Qingyu_backend/service/search"
)

// GrayScaleDecision 灰度决策器接口
type GrayScaleDecision interface {
	// ShouldUseES 判断是否应该使用 ES
	ShouldUseES(ctx context.Context, searchType string, userID string) bool

	// RecordUsage 记录引擎使用情况
	RecordUsage(engine string, took time.Duration)

	// GetMetrics 获取灰度指标
	GetMetrics() GrayScaleMetrics
}

// GrayScaleMetrics 灰度指标
type GrayScaleMetrics struct {
	ESCount       int64         // ES 使用次数
	MongoDBCount  int64         // MongoDB 使用次数
	ESAvgTook     time.Duration // ES 平均耗时
	MongoDBAvgTook time.Duration // MongoDB 平均耗时
	ESTotalTook   time.Duration // ES 总耗时
	MongoDBTotalTook time.Duration // MongoDB 总耗时
}

// grayscaleMetrics 灰度指标收集器
type grayscaleMetrics struct {
	mu            sync.RWMutex
	esCount       int64
	mongoCount    int64
	esTotalTook   time.Duration
	mongoTotalTook time.Duration
}

// grayscaleDecision 灰度决策器实现
type grayscaleDecision struct {
	config  *search.GrayScaleConfig
	metrics *grayscaleMetrics
	logger  *zap.Logger
}

// NewGrayScaleDecision 创建灰度决策器
func NewGrayScaleDecision(config *search.GrayScaleConfig, logger *zap.Logger) GrayScaleDecision {
	return &grayscaleDecision{
		config:  config,
		metrics: &grayscaleMetrics{},
		logger:  logger,
	}
}

// ShouldUseES 判断是否应该使用 ES
func (g *grayscaleDecision) ShouldUseES(ctx context.Context, searchType string, userID string) bool {
	// 如果灰度未启用，使用 MongoDB
	if g.config == nil || !g.config.Enabled {
		g.logger.Debug("灰度未启用，使用 MongoDB",
			zap.String("user_id", userID),
			zap.String("search_type", searchType),
		)
		return false
	}

	percent := g.config.Percent

	// 边界检查
	if percent <= 0 {
		g.logger.Debug("灰度百分比为0，使用 MongoDB",
			zap.String("user_id", userID),
			zap.Int("percent", percent),
		)
		return false
	}
	if percent >= 100 {
		g.logger.Debug("灰度百分比为100，使用 ES",
			zap.String("user_id", userID),
			zap.Int("percent", percent),
		)
		return true
	}

	// 使用用户 ID 哈希策略（保证同一用户的一致性）
	useES := hashBasedGrayScale(userID, percent)

	engine := "MongoDB"
	if useES {
		engine = "ES"
	}

	g.logger.Info("灰度决策完成",
		zap.String("user_id", userID),
		zap.String("search_type", searchType),
		zap.String("engine", engine),
		zap.Int("percent", percent),
		zap.Bool("use_es", useES),
	)

	return useES
}

// RecordUsage 记录引擎使用情况
func (g *grayscaleDecision) RecordUsage(engine string, took time.Duration) {
	g.metrics.mu.Lock()
	defer g.metrics.mu.Unlock()

	if engine == "es" || engine == "ES" {
		g.metrics.esCount++
		g.metrics.esTotalTook += took
	} else if engine == "mongodb" || engine == "MongoDB" {
		g.metrics.mongoCount++
		g.metrics.mongoTotalTook += took
	}

	g.logger.Debug("引擎使用情况已记录",
		zap.String("engine", engine),
		zap.Duration("took", took),
	)
}

// GetMetrics 获取灰度指标
func (g *grayscaleDecision) GetMetrics() GrayScaleMetrics {
	g.metrics.mu.RLock()
	defer g.metrics.mu.RUnlock()

	esAvgTook := time.Duration(0)
	if g.metrics.esCount > 0 {
		esAvgTook = g.metrics.esTotalTook / time.Duration(g.metrics.esCount)
	}

	mongoAvgTook := time.Duration(0)
	if g.metrics.mongoCount > 0 {
		mongoAvgTook = g.metrics.mongoTotalTook / time.Duration(g.metrics.mongoCount)
	}

	return GrayScaleMetrics{
		ESCount:        g.metrics.esCount,
		MongoDBCount:   g.metrics.mongoCount,
		ESAvgTook:      esAvgTook,
		MongoDBAvgTook: mongoAvgTook,
		ESTotalTook:    g.metrics.esTotalTook,
		MongoDBTotalTook: g.metrics.mongoTotalTook,
	}
}

// hashBasedGrayScale 基于用户 ID 哈希的灰度策略
func hashBasedGrayScale(userID string, percent int) bool {
	if userID == "" {
		// 如果没有用户 ID，使用随机策略
		return randomGrayScale(percent)
	}

	hash := crc32.ChecksumIEEE([]byte(userID))
	return int(hash%100) < percent
}

// randomGrayScale 基于随机数的灰度策略
func randomGrayScale(percent int) bool {
	return rand.Intn(100) < percent
}

// GetTrafficDistribution 获取流量分配比例
func (m *GrayScaleMetrics) GetTrafficDistribution() (esPercent, mongoPercent float64) {
	total := m.ESCount + m.MongoDBCount
	if total == 0 {
		return 0, 0
	}

	esPercent = float64(m.ESCount) / float64(total) * 100
	mongoPercent = float64(m.MongoDBCount) / float64(total) * 100

	return esPercent, mongoPercent
}
