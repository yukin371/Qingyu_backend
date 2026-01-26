package search

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"

	"Qingyu_backend/pkg/logger"
	"Qingyu_backend/service/search/engine"
)

// IndexVersion 索引版本信息
type IndexVersion struct {
	Name      string      `json:"name"`       // 索引名称（如: books_v1, books_v2）
	Alias     string      `json:"alias"`      // Alias 名称（如: books_search）
	Mapping   interface{} `json:"mapping"`    // Mapping 定义
	Status    string      `json:"status"`     // 状态: creating, active, inactive, deleting
	CreatedAt time.Time   `json:"created_at"` // 创建时间
	UpdatedAt time.Time   `json:"updated_at"` // 更新时间
}

// IndexVersionStatus 索引版本状态
const (
	StatusCreating = "creating"
	StatusActive   = "active"
	StatusInactive = "inactive"
	StatusDeleting = "deleting"
)

// IndexVersionManager 索引版本管理器接口
type IndexVersionManager interface {
	// CreateIndexVersion 创建新索引版本
	CreateIndexVersion(ctx context.Context, alias string, mapping interface{}) (*IndexVersion, error)

	// SwitchAlias 切换 alias 到新版本
	SwitchAlias(ctx context.Context, alias string, newIndex string) error

	// GetActiveVersion 获取当前激活的版本
	GetActiveVersion(ctx context.Context, alias string) (*IndexVersion, error)

	// ListVersions 列出所有版本
	ListVersions(ctx context.Context, alias string) ([]*IndexVersion, error)

	// RollbackVersion 回滚到指定版本
	RollbackVersion(ctx context.Context, alias string, targetVersion string) error

	// CleanupOldVersions 清理旧版本索引
	CleanupOldVersions(ctx context.Context, alias string, keep int) error
}

// elasticsearchIndexVersionManager Elasticsearch 索引版本管理器实现
type elasticsearchIndexVersionManager struct {
	client *elastic.Client
	engine engine.Engine
	logger *logger.Logger
}

// NewElasticsearchIndexVersionManager 创建 Elasticsearch 索引版本管理器
func NewElasticsearchIndexVersionManager(client *elastic.Client, eng engine.Engine) (IndexVersionManager, error) {
	if client == nil {
		return nil, fmt.Errorf("elasticsearch client cannot be nil")
	}
	if eng == nil {
		return nil, fmt.Errorf("search engine cannot be nil")
	}

	return &elasticsearchIndexVersionManager{
		client: client,
		engine: eng,
		logger: logger.Get().WithModule("index-version-manager"),
	}, nil
}

// CreateIndexVersion 创建新索引版本
func (m *elasticsearchIndexVersionManager) CreateIndexVersion(ctx context.Context, alias string, mapping interface{}) (*IndexVersion, error) {
	startTime := time.Now()

	// 1. 查找当前最新版本号
	versions, err := m.ListVersions(ctx, alias)
	if err != nil {
		m.logger.Error("Failed to list versions",
			zap.Error(err),
			zap.String("alias", alias),
		)
		return nil, fmt.Errorf("failed to list versions: %w", err)
	}

	// 2. 生成新版本号
	newVersionNumber := 1
	if len(versions) > 0 {
		// 获取最新版本号并递增
		latestVersion := versions[len(versions)-1]
		versionNum, err := extractVersionNumber(latestVersion.Name)
		if err != nil {
			m.logger.Error("Failed to extract version number",
				zap.Error(err),
				zap.String("index", latestVersion.Name),
			)
			return nil, fmt.Errorf("failed to extract version number: %w", err)
		}
		newVersionNumber = versionNum + 1
	}

	// 3. 生成新索引名称
	newIndexName := fmt.Sprintf("%s_v%d", alias, newVersionNumber)

	// 4. 创建新索引
	if err := m.engine.CreateIndex(ctx, newIndexName, mapping); err != nil {
		m.logger.Error("Failed to create index",
			zap.Error(err),
			zap.String("index", newIndexName),
			zap.String("alias", alias),
		)
		return nil, fmt.Errorf("failed to create index: %w", err)
	}

	// 5. 创建 IndexVersion 对象
	indexVersion := &IndexVersion{
		Name:      newIndexName,
		Alias:     alias,
		Mapping:   mapping,
		Status:    StatusInactive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	took := time.Since(startTime)

	m.logger.Info("Index version created",
		zap.String("alias", alias),
		zap.String("index", newIndexName),
		zap.Int("version", newVersionNumber),
		zap.Duration("took", took),
	)

	return indexVersion, nil
}

// SwitchAlias 切换 alias 到新版本
func (m *elasticsearchIndexVersionManager) SwitchAlias(ctx context.Context, alias string, newIndex string) error {
	startTime := time.Now()

	// 1. 验证新索引存在
	exists, err := m.client.IndexExists(newIndex).Do(ctx)
	if err != nil {
		m.logger.Error("Failed to check index existence",
			zap.Error(err),
			zap.String("index", newIndex),
		)
		return fmt.Errorf("failed to check index existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("index does not exist: %s", newIndex)
	}

	// 2. 获取当前 alias 指向的索引
	aliasResult, err := m.client.Aliases().Alias(alias).Do(ctx)
	if err != nil {
		m.logger.Error("Failed to get alias",
			zap.Error(err),
			zap.String("alias", alias),
		)
		return fmt.Errorf("failed to get alias: %w", err)
	}

	// 3. 构建 alias 切换操作
	aliasService := m.client.Alias()

	// 如果 alias 已存在，先移除旧索引的 alias
	if aliasResult != nil && len(aliasResult.Indices) > 0 {
		for indexName := range aliasResult.Indices {
			aliasService = aliasService.Remove(indexName, alias)
			m.logger.Info("Removing alias from old index",
				zap.String("alias", alias),
				zap.String("old_index", indexName),
			)
		}
	}

	// 4. 添加新索引的 alias
	aliasService = aliasService.Add(newIndex, alias)

	// 5. 执行 alias 切换
	result, err := aliasService.Do(ctx)
	if err != nil {
		m.logger.Error("Failed to switch alias",
			zap.Error(err),
			zap.String("alias", alias),
			zap.String("new_index", newIndex),
		)
		return fmt.Errorf("failed to switch alias: %w", err)
	}

	if !result.Acknowledged {
		return fmt.Errorf("alias switch not acknowledged")
	}

	took := time.Since(startTime)

	m.logger.Info("Alias switched successfully",
		zap.String("alias", alias),
		zap.String("new_index", newIndex),
		zap.Bool("acknowledged", result.Acknowledged),
		zap.Duration("took", took),
	)

	return nil
}

// GetActiveVersion 获取当前激活的版本
func (m *elasticsearchIndexVersionManager) GetActiveVersion(ctx context.Context, alias string) (*IndexVersion, error) {
	// 1. 查询 alias 指向的索引
	aliasResult, err := m.client.Aliases().Alias(alias).Do(ctx)
	if err != nil {
		m.logger.Error("Failed to get alias",
			zap.Error(err),
			zap.String("alias", alias),
		)
		return nil, fmt.Errorf("failed to get alias: %w", err)
	}

	// 2. 检查 alias 是否存在
	if aliasResult == nil || len(aliasResult.Indices) == 0 {
		return nil, fmt.Errorf("alias not found: %s", alias)
	}

	// 3. 获取指向的索引名称（alias 应该只指向一个索引）
	var activeIndex string
	for indexName := range aliasResult.Indices {
		activeIndex = indexName
		break
	}

	// 4. 获取索引设置
	indexService, err := m.client.IndexGetSettings(activeIndex).Do(ctx)
	if err != nil {
		m.logger.Error("Failed to get index settings",
			zap.Error(err),
			zap.String("index", activeIndex),
		)
		return nil, fmt.Errorf("failed to get index settings: %w", err)
	}

	// 5. 构建 IndexVersion 对象
	indexVersion := &IndexVersion{
		Name:      activeIndex,
		Alias:     alias,
		Status:    StatusActive,
		CreatedAt: time.Now(),   // 实际应从索引 settings 获取
		UpdatedAt: time.Now(),   // 实际应从索引 settings 获取
	}

	// 从索引设置中获取创建时间
	if settings, ok := indexService[activeIndex]; ok && settings.Settings != nil {
		if creationDateStr, ok := settings.Settings["index.creation_date"].(string); ok {
			if timestamp, err := strconv.ParseInt(creationDateStr, 10, 64); err == nil {
				indexVersion.CreatedAt = time.Unix(0, timestamp*int64(time.Millisecond))
			}
		}
	}

	m.logger.Info("Active version retrieved",
		zap.String("alias", alias),
		zap.String("index", activeIndex),
	)

	return indexVersion, nil
}

// ListVersions 列出所有版本
func (m *elasticsearchIndexVersionManager) ListVersions(ctx context.Context, alias string) ([]*IndexVersion, error) {
	// 1. 查询所有索引
	catIndicesResult, err := m.client.CatIndices().Do(ctx)
	if err != nil {
		m.logger.Error("Failed to list indices",
			zap.Error(err),
			zap.String("alias", alias),
		)
		return nil, fmt.Errorf("failed to list indices: %w", err)
	}

	// 2. 获取当前 alias 指向的索引
	catAliasesResult, err := m.client.CatAliases().Alias(alias).Do(ctx)
	var activeIndex string
	if err == nil && len(catAliasesResult) > 0 {
		// CatAliasesResponse 是 []CatAliasesResponseRow，第一个元素就是 alias 指向的索引
		activeIndex = catAliasesResult[0].Index
	}

	// 3. 过滤出匹配 alias 的版本索引
	versionPattern := regexp.MustCompile(fmt.Sprintf("^%s_v\\d+$", regexp.QuoteMeta(alias)))
	var versions []*IndexVersion

	for _, indexInfo := range catIndicesResult {
		if versionPattern.MatchString(indexInfo.Index) {
			// 判断索引状态
			status := StatusInactive
			if activeIndex == indexInfo.Index {
				status = StatusActive
			}

			// 解析创建时间
			var createdAt time.Time
			if indexInfo.CreationDate > 0 {
				// Elasticsearch 返回的时间戳是毫秒级
				createdAt = time.Unix(0, indexInfo.CreationDate*int64(time.Millisecond))
			}

			versions = append(versions, &IndexVersion{
				Name:      indexInfo.Index,
				Alias:     alias,
				Status:    status,
				CreatedAt: createdAt,
				UpdatedAt: createdAt,
			})
		}
	}

	// 4. 按版本号排序
	sortVersions(versions)

	m.logger.Info("Versions listed",
		zap.String("alias", alias),
		zap.Int("count", len(versions)),
	)

	return versions, nil
}

// RollbackVersion 回滚到指定版本
func (m *elasticsearchIndexVersionManager) RollbackVersion(ctx context.Context, alias string, targetVersion string) error {
	startTime := time.Now()

	// 1. 验证目标版本索引存在
	exists, err := m.client.IndexExists(targetVersion).Do(ctx)
	if err != nil {
		m.logger.Error("Failed to check index existence",
			zap.Error(err),
			zap.String("index", targetVersion),
		)
		return fmt.Errorf("failed to check index existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("target version index does not exist: %s", targetVersion)
	}

	// 2. 验证目标版本是否匹配 alias
	expectedPrefix := alias + "_v"
	if !strings.HasPrefix(targetVersion, expectedPrefix) {
		return fmt.Errorf("target version %s does not match alias %s", targetVersion, alias)
	}

	// 3. 切换 alias 到目标版本
	if err := m.SwitchAlias(ctx, alias, targetVersion); err != nil {
		m.logger.Error("Failed to rollback alias",
			zap.Error(err),
			zap.String("alias", alias),
			zap.String("target_version", targetVersion),
		)
		return fmt.Errorf("failed to rollback alias: %w", err)
	}

	took := time.Since(startTime)

	m.logger.Info("Alias rolled back successfully",
		zap.String("alias", alias),
		zap.String("target_version", targetVersion),
		zap.Duration("took", took),
	)

	return nil
}

// CleanupOldVersions 清理旧版本索引
func (m *elasticsearchIndexVersionManager) CleanupOldVersions(ctx context.Context, alias string, keep int) error {
	startTime := time.Now()

	if keep < 1 {
		return fmt.Errorf("keep must be at least 1")
	}

	// 1. 列出所有版本
	versions, err := m.ListVersions(ctx, alias)
	if err != nil {
		m.logger.Error("Failed to list versions",
			zap.Error(err),
			zap.String("alias", alias),
		)
		return fmt.Errorf("failed to list versions: %w", err)
	}

	// 2. 如果版本数量不超过保留数量，无需清理
	if len(versions) <= keep {
		m.logger.Info("No versions to cleanup",
			zap.String("alias", alias),
			zap.Int("total", len(versions)),
			zap.Int("keep", keep),
		)
		return nil
	}

	// 3. 找出需要删除的旧版本（保留最新的 keep 个版本）
	var toDelete []*IndexVersion
	for i := 0; i < len(versions)-keep; i++ {
		// 只删除 inactive 状态的索引
		if versions[i].Status == StatusInactive {
			toDelete = append(toDelete, versions[i])
		}
	}

	// 4. 删除旧版本索引
	deletedCount := 0
	for _, version := range toDelete {
		// 确认不是 active 状态
		if version.Status == StatusActive {
			m.logger.Warn("Skipping active index",
				zap.String("index", version.Name),
				zap.String("status", version.Status),
			)
			continue
		}

		// 删除索引
		deleteResult, err := m.client.DeleteIndex(version.Name).Do(ctx)
		if err != nil {
			m.logger.Error("Failed to delete index",
				zap.Error(err),
				zap.String("index", version.Name),
			)
			continue
		}

		if deleteResult.Acknowledged {
			deletedCount++
			m.logger.Info("Index version deleted",
				zap.String("index", version.Name),
				zap.String("alias", alias),
			)
		}
	}

	took := time.Since(startTime)

	m.logger.Info("Old versions cleanup completed",
		zap.String("alias", alias),
		zap.Int("deleted", deletedCount),
		zap.Int("kept", len(versions)-deletedCount),
		zap.Duration("took", took),
	)

	return nil
}

// extractVersionNumber 从索引名称中提取版本号
// 例如: books_v1 -> 1, books_v2 -> 2
func extractVersionNumber(indexName string) (int, error) {
	// 查找最后一个下划线位置
	lastUnderscore := strings.LastIndex(indexName, "_")
	if lastUnderscore == -1 {
		return 0, fmt.Errorf("invalid index name format: %s", indexName)
	}

	// 提取版本号部分
	versionStr := indexName[lastUnderscore+1:]
	if !strings.HasPrefix(versionStr, "v") {
		return 0, fmt.Errorf("version prefix not found: %s", indexName)
	}

	// 去掉 'v' 前缀
	versionNum, err := strconv.Atoi(versionStr[1:])
	if err != nil {
		return 0, fmt.Errorf("failed to parse version number: %w", err)
	}

	return versionNum, nil
}

// sortVersions 按版本号排序索引版本列表
func sortVersions(versions []*IndexVersion) {
	// 使用简单的冒泡排序按版本号升序排列
	n := len(versions)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			v1, err1 := extractVersionNumber(versions[j].Name)
			v2, err2 := extractVersionNumber(versions[j+1].Name)

			// 如果提取失败，保持原顺序
			if err1 != nil || err2 != nil {
				continue
			}

			if v1 > v2 {
				versions[j], versions[j+1] = versions[j+1], versions[j]
			}
		}
	}
}
