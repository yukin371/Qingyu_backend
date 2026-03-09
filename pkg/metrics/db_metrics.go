package metrics

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type commandTrace struct {
	commandName  string
	collection   string
	startedAt    time.Time
	command      string
	connectionID string
}

// DbMetricsCollector 数据库指标收集器
type DbMetricsCollector struct {
	metrics         *Metrics
	databaseName    string
	slowThreshold   time.Duration
	profilingLevel  int
	inFlightQueries sync.Map
}

// NewDbMetricsCollector 创建数据库指标收集器
func NewDbMetricsCollector(databaseName string, slowThreshold time.Duration, profilingLevel int) *DbMetricsCollector {
	return &DbMetricsCollector{
		metrics:        GetMetrics(),
		databaseName:   databaseName,
		slowThreshold:  slowThreshold,
		profilingLevel: profilingLevel,
	}
}

// GetMonitorCommand 监控命令事件
func (c *DbMetricsCollector) GetMonitorCommand() *event.CommandMonitor {
	return &event.CommandMonitor{
		Started: func(_ context.Context, started *event.CommandStartedEvent) {
			c.inFlightQueries.Store(started.RequestID, commandTrace{
				commandName:  started.CommandName,
				collection:   collectionNameFromCommand(started.Command, started.CommandName),
				startedAt:    time.Now(),
				command:      truncateCommand(started.Command.String()),
				connectionID: started.ConnectionID,
			})
		},
		Succeeded: func(_ context.Context, succeeded *event.CommandSucceededEvent) {
			c.finishQuery(succeeded.RequestID, succeeded.CommandName, succeeded.Duration, true, "")
		},
		Failed: func(_ context.Context, failed *event.CommandFailedEvent) {
			c.finishQuery(failed.RequestID, failed.CommandName, failed.Duration, false, failed.Failure)
		},
	}
}

func (c *DbMetricsCollector) finishQuery(requestID int64, commandName string, duration time.Duration, success bool, failure string) {
	trace := commandTrace{
		commandName: commandName,
		startedAt:   time.Now().Add(-duration),
	}

	if stored, ok := c.inFlightQueries.LoadAndDelete(requestID); ok {
		trace = stored.(commandTrace)
	}

	operation := commandName
	if trace.commandName != "" {
		operation = trace.commandName
	}

	c.metrics.RecordDbQuery(c.databaseName, operation, duration, success)

	if !success {
		zap.L().Warn("MongoDB query failed",
			zap.String("database", c.databaseName),
			zap.String("operation", operation),
			zap.String("collection", trace.collection),
			zap.Duration("duration", duration),
			zap.String("connection_id", trace.connectionID),
			zap.String("command", trace.command),
			zap.String("failure", failure),
		)
		return
	}

	if c.shouldLogQuery(duration) {
		zap.L().Warn("MongoDB slow query detected",
			zap.String("database", c.databaseName),
			zap.String("operation", operation),
			zap.String("collection", trace.collection),
			zap.Duration("duration", duration),
			zap.Duration("threshold", c.slowThreshold),
			zap.String("connection_id", trace.connectionID),
			zap.String("command", trace.command),
		)
	}
}

func (c *DbMetricsCollector) shouldLogQuery(duration time.Duration) bool {
	if c.profilingLevel == 2 {
		return true
	}
	if c.profilingLevel <= 0 {
		return false
	}
	return duration >= c.slowThreshold
}

func collectionNameFromCommand(command bson.Raw, commandName string) string {
	value := command.Lookup(commandName)
	if value.Type != bson.TypeString {
		return ""
	}
	return value.StringValue()
}

func truncateCommand(command string) string {
	const maxLen = 512
	if len(command) <= maxLen {
		return command
	}
	return command[:maxLen] + "...(truncated)"
}

// GetPoolMonitor 获取连接池监控
func (c *DbMetricsCollector) GetPoolMonitor() *event.PoolMonitor {
	return &event.PoolMonitor{
		Event: func(evt *event.PoolEvent) {
			switch evt.Type {
			case event.ConnectionCreated:
				// 连接创建
			case event.ConnectionClosed:
				// 连接关闭
			case event.GetSucceeded:
				// 获取连接成功
			case event.ConnectionReturned:
				// 连接归还
			case event.PoolCleared:
				// 连接池清理
			}
		},
	}
}

// UpdatePoolMetrics 更新连接池指标
func (c *DbMetricsCollector) UpdatePoolMetrics(client *mongo.Client) error {
	// TODO: 从MongoDB客户端获取连接池统计信息
	// MongoDB Go Driver 1.8+ 提供了连接池统计信息
	// 需要通过类型断言访问内部统计信息

	return nil
}

// StartPoolMonitoring 启动连接池监控
// 定期收集连接池指标
func (c *DbMetricsCollector) StartPoolMonitoring(ctx context.Context, client *mongo.Client, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.UpdatePoolMetrics(client)
		}
	}
}

// DbStats 数据库统计信息
type DbStats struct {
	// 连接池信息
	PoolConnections     int
	PoolIdleConnections int
	PoolWaitCount       int64
	PoolWaitDuration    time.Duration

	// 查询统计
	QueryTotal       int64
	QueryErrors      int64
	QueryDuration    time.Duration
	AvgQueryDuration time.Duration

	// 数据库信息
	DatabaseName    string
	DatabaseSize    int64
	CollectionCount int
}

// GetDbStats 获取数据库统计信息
func (c *DbMetricsCollector) GetDbStats(client *mongo.Client) (*DbStats, error) {
	stats := &DbStats{
		DatabaseName: c.databaseName,
	}

	// TODO: 从MongoDB获取实际的统计信息
	// 这需要运行dbStats命令

	return stats, nil
}
