package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"

	"Qingyu_backend/pkg/logger"
	"Qingyu_backend/service/search"
	engine "Qingyu_backend/service/search/engine"
)

var (
	// 全局标志
	esURL      = flag.String("es-url", "http://localhost:9200", "Elasticsearch URL")
	configPath = flag.String("config", "config/search_indices.yaml", "Config file path")
	verbose    = flag.Bool("verbose", false, "Verbose output")
)

// elasticLogger 适配器，使 zap.Logger 实现 elastic.Logger 接口
type elasticLogger struct {
	*zap.Logger
}

// Printf 实现 elastic.Logger 接口
func (l *elasticLogger) Printf(format string, v ...interface{}) {
	l.Info(fmt.Sprintf(format, v...))
}


func main() {
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("Usage: search_index <command> [options]")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  init <index>       Initialize an index")
		fmt.Println("  migrate <index>    Create a new index version and migrate data")
		fmt.Println("  switch <index> <version>   Switch alias to a new index version")
		fmt.Println("  rollback <index> <version> Rollback to a previous index version")
		fmt.Println("  status <index>     Show index version status")
		fmt.Println("  cleanup <index> [keep]   Clean up old index versions")
		fmt.Println()
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// 初始化日志
	if *verbose {
		// 注意: logger.Logger 不支持 SetLevel, 这里只是记录要使用 verbose 模式
		logger.Get().Debug("Verbose mode enabled")
	}

	command := flag.Arg(0)
	args := flag.Args()[1:]

	var err error
	switch command {
	case "init":
		err = runInit(args)
	case "migrate":
		err = runMigrate(args)
	case "switch":
		err = runSwitch(args)
	case "rollback":
		err = runRollback(args)
	case "status":
		err = runStatus(args)
	case "cleanup":
		err = runCleanup(args)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// createESClient 创建 Elasticsearch 客户端
func createESClient() (*elastic.Client, error) {
	// 创建 elastic logger 适配器
	elasticLogger := &elasticLogger{Logger: logger.Get().Logger}

	client, err := elastic.NewClient(
		elastic.SetURL(*esURL),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetErrorLog(elasticLogger),
		elastic.SetInfoLog(elasticLogger),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %w", err)
	}

	return client, nil
}

// loadConfig 加载索引配置
func loadConfig() (*search.SearchIndicesConfig, error) {
	config, err := search.LoadSearchIndicesConfig(*configPath)
	if err != nil {
		// 如果配置文件不存在，使用默认配置
		if os.IsNotExist(err) {
			logger.Get().Warn("Config file not found, using default config",
				zap.String("path", *configPath),
			)
			return search.GetDefaultIndicesConfig(), nil
		}
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return config, nil
}

// getIndexConfig 获取索引配置
func getIndexConfig(config *search.SearchIndicesConfig, indexName string) (*search.IndexConfig, error) {
	indexConfig, ok := config.Indices[indexName]
	if !ok {
		return nil, fmt.Errorf("index '%s' not found in config", indexName)
	}

	return &indexConfig, nil
}

// runInit 执行初始化命令
func runInit(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("init command requires index name")
	}

	indexName := args[0]
	ctx := context.Background()

	fmt.Printf("Initializing index: %s\n", indexName)

	// 1. 创建 ES 客户端
	client, err := createESClient()
	if err != nil {
		return err
	}

	// 2. 加载配置
	config, err := loadConfig()
	if err != nil {
		return err
	}

	// 3. 获取索引配置
	indexConfig, err := getIndexConfig(config, indexName)
	if err != nil {
		return err
	}

	// 4. 创建搜索引擎
	esEngine, err := engine.NewElasticsearchEngine(client)
	if err != nil {
		return fmt.Errorf("failed to create search engine: %w", err)
	}

	// 5. 创建索引版本管理器
	versionManager, err := search.NewElasticsearchIndexVersionManager(client, esEngine)
	if err != nil {
		return fmt.Errorf("failed to create version manager: %w", err)
	}

	// 6. 检查 alias 是否已存在
	exists, err := client.Aliases().Alias(indexConfig.Alias).Do(ctx)
	if err == nil && exists != nil && len(exists.Indices) > 0 {
		return fmt.Errorf("alias '%s' already exists, use 'migrate' to create new version", indexConfig.Alias)
	}

	// 7. 创建初始索引 (v1)
	mapping := indexConfig.BuildIndexMapping(config.Settings)
	indexVersion, err := versionManager.CreateIndexVersion(ctx, indexConfig.Alias, mapping)
	if err != nil {
		return fmt.Errorf("failed to create index version: %w", err)
	}

	fmt.Printf("Created index: %s\n", indexVersion.Name)

	// 8. 创建 alias 指向新索引
	if err := versionManager.SwitchAlias(ctx, indexConfig.Alias, indexVersion.Name); err != nil {
		return fmt.Errorf("failed to create alias: %w", err)
	}

	fmt.Printf("Created alias: %s -> %s\n", indexConfig.Alias, indexVersion.Name)
	fmt.Printf("Index '%s' initialized successfully!\n", indexName)

	return nil
}

// runMigrate 执行迁移命令
func runMigrate(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("migrate command requires index name")
	}

	indexName := args[0]
	ctx := context.Background()

	fmt.Printf("Migrating index: %s\n", indexName)

	// 1. 创建 ES 客户端
	client, err := createESClient()
	if err != nil {
		return err
	}

	// 2. 加载配置
	config, err := loadConfig()
	if err != nil {
		return err
	}

	// 3. 获取索引配置
	indexConfig, err := getIndexConfig(config, indexName)
	if err != nil {
		return err
	}

	// 4. 创建搜索引擎
	esEngine, err := engine.NewElasticsearchEngine(client)
	if err != nil {
		return fmt.Errorf("failed to create search engine: %w", err)
	}

	// 5. 创建索引版本管理器
	versionManager, err := search.NewElasticsearchIndexVersionManager(client, esEngine)
	if err != nil {
		return fmt.Errorf("failed to create version manager: %w", err)
	}

	// 6. 获取当前激活的版本
	activeVersion, err := versionManager.GetActiveVersion(ctx, indexConfig.Alias)
	if err != nil {
		return fmt.Errorf("failed to get active version: %w", err)
	}

	fmt.Printf("Current active version: %s\n", activeVersion.Name)

	// 7. 创建新版本索引
	mapping := indexConfig.BuildIndexMapping(config.Settings)
	newIndexVersion, err := versionManager.CreateIndexVersion(ctx, indexConfig.Alias, mapping)
	if err != nil {
		return fmt.Errorf("failed to create new index version: %w", err)
	}

	fmt.Printf("Created new index: %s\n", newIndexVersion.Name)

	// 8. 使用 reindex API 迁移数据
	reindexResult, err := client.Reindex().
		Source(elastic.NewReindexSource().Index(activeVersion.Name)).
		DestinationIndex(newIndexVersion.Name).
		Refresh("true").
		WaitForCompletion(true).
		Do(ctx)

	if err != nil {
		// 如果 reindex 失败，删除新创建的索引
		_, _ = client.DeleteIndex(newIndexVersion.Name).Do(ctx)
		return fmt.Errorf("failed to reindex data: %w", err)
	}

	fmt.Printf("Migrated %d documents from %s to %s\n",
		reindexResult.Created, activeVersion.Name, newIndexVersion.Name)

	fmt.Printf("Migration completed successfully!\n")
	fmt.Printf("New index '%s' is ready. Use 'switch' command to activate it.\n", newIndexVersion.Name)

	return nil
}

// runSwitch 执行切换命令
func runSwitch(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("switch command requires index name and version number")
	}

	indexName := args[0]
	var version int
	_, err := fmt.Sscanf(args[1], "%d", &version)
	if err != nil {
		return fmt.Errorf("invalid version number: %s", args[1])
	}

	ctx := context.Background()

	targetIndex := fmt.Sprintf("%s_v%d", indexName, version)

	fmt.Printf("Switching alias for index: %s\n", indexName)
	fmt.Printf("Target version: %d (index: %s)\n", version, targetIndex)

	// 1. 创建 ES 客户端
	client, err := createESClient()
	if err != nil {
		return err
	}

	// 2. 加载配置
	config, err := loadConfig()
	if err != nil {
		return err
	}

	// 3. 获取索引配置
	indexConfig, err := getIndexConfig(config, indexName)
	if err != nil {
		return err
	}

	// 4. 创建搜索引擎
	esEngine, err := engine.NewElasticsearchEngine(client)
	if err != nil {
		return fmt.Errorf("failed to create search engine: %w", err)
	}

	// 5. 创建索引版本管理器
	versionManager, err := search.NewElasticsearchIndexVersionManager(client, esEngine)
	if err != nil {
		return fmt.Errorf("failed to create version manager: %w", err)
	}

	// 6. 检查目标索引是否存在
	indexExists, err := client.IndexExists(targetIndex).Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to check index existence: %w", err)
	}
	if !indexExists {
		return fmt.Errorf("target index does not exist: %s", targetIndex)
	}

	// 7. 获取当前激活的版本
	oldActiveVersion, err := versionManager.GetActiveVersion(ctx, indexConfig.Alias)
	if err != nil {
		fmt.Printf("Warning: failed to get current active version: %v\n", err)
	} else {
		fmt.Printf("Current active version: %s\n", oldActiveVersion.Name)
	}

	// 8. 切换 alias
	if err := versionManager.SwitchAlias(ctx, indexConfig.Alias, targetIndex); err != nil {
		return fmt.Errorf("failed to switch alias: %w", err)
	}

	fmt.Printf("Alias switched successfully: %s -> %s\n", indexConfig.Alias, targetIndex)
	fmt.Printf("Index '%s' is now using version %d\n", indexName, version)

	return nil
}

// runRollback 执行回滚命令
func runRollback(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("rollback command requires index name and version number")
	}

	indexName := args[0]
	var version int
	_, err := fmt.Sscanf(args[1], "%d", &version)
	if err != nil {
		return fmt.Errorf("invalid version number: %s", args[1])
	}

	ctx := context.Background()

	targetIndex := fmt.Sprintf("%s_v%d", indexName, version)

	fmt.Printf("Rolling back index: %s\n", indexName)
	fmt.Printf("Target version: %d (index: %s)\n", version, targetIndex)

	// 1. 创建 ES 客户端
	client, err := createESClient()
	if err != nil {
		return err
	}

	// 2. 加载配置
	config, err := loadConfig()
	if err != nil {
		return err
	}

	// 3. 获取索引配置
	indexConfig, err := getIndexConfig(config, indexName)
	if err != nil {
		return err
	}

	// 4. 创建搜索引擎
	esEngine, err := engine.NewElasticsearchEngine(client)
	if err != nil {
		return fmt.Errorf("failed to create search engine: %w", err)
	}

	// 5. 创建索引版本管理器
	versionManager, err := search.NewElasticsearchIndexVersionManager(client, esEngine)
	if err != nil {
		return fmt.Errorf("failed to create version manager: %w", err)
	}

	// 6. 执行回滚
	if err := versionManager.RollbackVersion(ctx, indexConfig.Alias, targetIndex); err != nil {
		return fmt.Errorf("failed to rollback: %w", err)
	}

	fmt.Printf("Rollback completed successfully!\n")
	fmt.Printf("Index '%s' is now using version %d\n", indexName, version)

	return nil
}

// runStatus 执行状态命令
func runStatus(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("status command requires index name")
	}

	indexName := args[0]
	ctx := context.Background()

	fmt.Printf("Index status for: %s\n", indexName)
	fmt.Println()

	// 1. 创建 ES 客户端
	client, err := createESClient()
	if err != nil {
		return err
	}

	// 2. 加载配置
	config, err := loadConfig()
	if err != nil {
		return err
	}

	// 3. 获取索引配置
	indexConfig, err := getIndexConfig(config, indexName)
	if err != nil {
		return err
	}

	// 4. 创建搜索引擎
	esEngine, err := engine.NewElasticsearchEngine(client)
	if err != nil {
		return fmt.Errorf("failed to create search engine: %w", err)
	}

	// 5. 创建索引版本管理器
	versionManager, err := search.NewElasticsearchIndexVersionManager(client, esEngine)
	if err != nil {
		return fmt.Errorf("failed to create version manager: %w", err)
	}

	// 6. 获取所有版本
	versions, err := versionManager.ListVersions(ctx, indexConfig.Alias)
	if err != nil {
		return fmt.Errorf("failed to list versions: %w", err)
	}

	if len(versions) == 0 {
		fmt.Printf("No versions found for index '%s'\n", indexName)
		fmt.Printf("Use 'init' command to create the first version.\n")
		return nil
	}

	// 7. 显示版本列表
	fmt.Println("Version List:")
	fmt.Println("-------------")
	for i, v := range versions {
		statusIndicator := "  "
		if v.Status == search.StatusActive {
			statusIndicator = "* " // 标记激活版本
		}

		fmt.Printf("%sVersion: %s\n", statusIndicator, v.Name)
		fmt.Printf("   Status: %s\n", v.Status)
		fmt.Printf("   Alias:  %s\n", v.Alias)
		fmt.Printf("   Created: %s\n", v.CreatedAt.Format("2006-01-02 15:04:05"))

		if i < len(versions)-1 {
			fmt.Println()
		}
	}

	fmt.Println()
	fmt.Printf("Total versions: %d\n", len(versions))

	return nil
}

// runCleanup 执行清理命令
func runCleanup(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("cleanup command requires index name")
	}

	indexName := args[0]
	keep := 2 // 默认保留最近 2 个版本
	if len(args) >= 2 {
		_, err := fmt.Sscanf(args[1], "%d", &keep)
		if err != nil {
			return fmt.Errorf("invalid keep number: %s", args[1])
		}
	}

	ctx := context.Background()

	fmt.Printf("Cleaning up old versions for index: %s\n", indexName)
	fmt.Printf("Keeping the most recent %d versions\n", keep)
	fmt.Println()

	// 1. 创建 ES 客户端
	client, err := createESClient()
	if err != nil {
		return err
	}

	// 2. 加载配置
	config, err := loadConfig()
	if err != nil {
		return err
	}

	// 3. 获取索引配置
	indexConfig, err := getIndexConfig(config, indexName)
	if err != nil {
		return err
	}

	// 4. 创建搜索引擎
	esEngine, err := engine.NewElasticsearchEngine(client)
	if err != nil {
		return fmt.Errorf("failed to create search engine: %w", err)
	}

	// 5. 创建索引版本管理器
	versionManager, err := search.NewElasticsearchIndexVersionManager(client, esEngine)
	if err != nil {
		return fmt.Errorf("failed to create version manager: %w", err)
	}

	// 6. 获取所有版本
	versions, err := versionManager.ListVersions(ctx, indexConfig.Alias)
	if err != nil {
		return fmt.Errorf("failed to list versions: %w", err)
	}

	if len(versions) <= keep {
		fmt.Printf("Total versions: %d (keep: %d)\n", len(versions), keep)
		fmt.Println("No cleanup needed. Total versions <= keep count.")
		return nil
	}

	// 7. 显示将要删除的版本
	fmt.Printf("Current versions: %d\n", len(versions))
	fmt.Printf("Versions to keep: %d\n", keep)
	fmt.Printf("Versions to delete: %d\n", len(versions)-keep)
	fmt.Println()

	if len(versions) > keep {
		fmt.Println("Versions to be deleted:")
		for i := 0; i < len(versions)-keep; i++ {
			if versions[i].Status == search.StatusInactive {
				fmt.Printf("  - %s (created: %s)\n",
					versions[i].Name,
					versions[i].CreatedAt.Format("2006-01-02 15:04:05"))
			}
		}
		fmt.Println()
	}

	// 8. 执行清理
	if err := versionManager.CleanupOldVersions(ctx, indexConfig.Alias, keep); err != nil {
		return fmt.Errorf("failed to cleanup old versions: %w", err)
	}

	fmt.Println("Cleanup completed successfully!")

	return nil
}
