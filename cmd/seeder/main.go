// Package main provides the test data seeder for Qingyu writing platform
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"Qingyu_backend/cmd/seeder/config"
	"Qingyu_backend/cmd/seeder/utils"
	"Qingyu_backend/cmd/seeder/validator"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	// 命令行标志
	cfgFile string
	scale   string
	clean   bool
	createDemo    bool
	forceCleanDemo bool

	// 配置
	cfg *config.Config

	// 根命令
	rootCmd = &cobra.Command{
		Use:   "seeder",
		Short: "青羽写作平台测试数据填充工具",
		Long: `seeder 是一个用于生成青羽写作平台测试数据的命令行工具。
支持生成用户、书籍、订阅关系等多种测试数据，并提供数据验证功能。`,
	}

	// allCmd 执行所有填充操作
	allCmd = &cobra.Command{
		Use:   "all",
		Short: "执行核心联调基线数据填充",
		Run:   runAll,
	}

	// baselineCmd 构建联调基线数据
	baselineCmd = &cobra.Command{
		Use:   "baseline",
		Short: "构建联调基线数据（用户/书城/章节/社交/阅读/统计）",
		Run:   runBaseline,
	}

	// fullCmd 构建完整测试数据
	fullCmd = &cobra.Command{
		Use:   "full",
		Short: "构建完整测试数据（基线 + 钱包/通知/消息/财务/AI配额）",
		Run:   runFull,
	}

	// usersCmd 只填充用户数据
	usersCmd = &cobra.Command{
		Use:   "users",
		Short: "只填充用户数据",
		Run:   runUsers,
	}

	// categoriesCmd 只填充分类数据
	categoriesCmd = &cobra.Command{
		Use:   "categories",
		Short: "填充标准分类数据",
		Run:   runCategories,
	}

	// bookstoreCmd 只填充书籍数据
	bookstoreCmd = &cobra.Command{
		Use:   "bookstore",
		Short: "只填充书籍数据",
		Run:   runBookstore,
	}

	// showcaseCmd 只填充精选演示书籍数据
	showcaseCmd = &cobra.Command{
		Use:   "showcase",
		Short: "填充精选演示书籍数据",
		Run:   runShowcase,
	}

	// subscriptionsCmd 只填充书籍订阅关系
	subscriptionsCmd = &cobra.Command{
		Use:   "subscriptions",
		Short: "填充书籍订阅关系",
		Run:   runSubscriptions,
	}

	// cleanCmd 清空所有数据
	cleanCmd = &cobra.Command{
		Use:   "clean",
		Short: "清空所有测试数据",
		Run:   runClean,
	}

	// verifyCmd 验证数据完整性
	verifyCmd = &cobra.Command{
		Use:   "verify",
		Short: "验证数据完整性",
		Run:   runVerify,
	}

	// testCmd 填充E2E测试所需的数据
	testCmd = &cobra.Command{
		Use:   "test",
		Short: "填充E2E测试所需的特定数据",
		Run:   runTestData,
	}

	// publicationFlowCmd 填充完整发布流程测试数据
	publicationFlowCmd = &cobra.Command{
		Use:   "publication-flow",
		Short: "通过完整发布流程创建测试数据（作者→项目→文档→发布→审核→书城）",
		Run:   runPublicationFlow,
	}

	// chaptersCmd 填充章节数据
	chaptersCmd = &cobra.Command{
		Use:   "chapters",
		Short: "填充章节数据（需要先有书籍）",
		Run:   runChapters,
	}

	// socialCmd 填充社交数据
	socialCmd = &cobra.Command{
		Use:   "social",
		Short: "填充社交数据（评论、点赞、收藏、关注）",
		Run:   runSocial,
	}

	// walletsCmd 填充钱包数据
	walletsCmd = &cobra.Command{
		Use:   "wallets",
		Short: "填充钱包和交易数据",
		Run:   runWallets,
	}

	// rankingsCmd 填充榜单数据
	rankingsCmd = &cobra.Command{
		Use:   "rankings",
		Short: "填充榜单数据",
		Run:   runRankings,
	}

	// aiQuotaCmd 激活AI配额
	aiQuotaCmd = &cobra.Command{
		Use:   "ai-quota",
		Short: "激活用户AI配额",
		Run:   runAIQuota,
	}

	// importCmd 导入小说数据
	importCmd = &cobra.Command{
		Use:   "import",
		Short: "从JSON文件导入小说数据",
		Run:   runImport,
	}

	// readerCmd 填充阅读数据
	readerCmd = &cobra.Command{
		Use:   "reader",
		Short: "填充阅读数据（阅读历史、书架、订阅）",
		Run:   runReader,
	}

	// notificationsCmd 填充通知数据
	notificationsCmd = &cobra.Command{
		Use:   "notifications",
		Short: "填充通知数据",
		Run:   runNotifications,
	}

	// messagingCmd 填充消息数据
	messagingCmd = &cobra.Command{
		Use:   "messaging",
		Short: "填充消息数据（对话、消息、公告）",
		Run:   runMessaging,
	}

	// statsCmd 填充统计数据
	statsCmd = &cobra.Command{
		Use:   "stats",
		Short: "填充统计数据（书籍统计、章节统计）",
		Run:   runStats,
	}

	// financeCmd 填充财务数据
	financeCmd = &cobra.Command{
		Use:   "finance",
		Short: "填充财务数据（作者收入、会员）",
		Run:   runFinance,
	}

	// booklistsCmd 填充书单数据
	booklistsCmd = &cobra.Command{
		Use:   "booklists",
		Short: "填充书单数据（书单、点赞）",
		Run:   runBooklists,
	}

	// collectionFoldersCmd 填充收藏夹数据
	collectionFoldersCmd = &cobra.Command{
		Use:   "collection-folders",
		Short: "填充收藏夹数据",
		Run:   runCollectionFolders,
	}

	// withdrawalsCmd 填充提现申请数据
	withdrawalsCmd = &cobra.Command{
		Use:   "withdrawals",
		Short: "填充提现申请数据",
		Run:   runWithdrawals,
	}

	// auditReaderCmd 审查读者视角数据关联
	auditReaderCmd = &cobra.Command{
		Use:   "audit-reader",
		Short: "审查读者视角数据关联完整性",
		Run:   runAuditReader,
	}

	// auditAuthorCmd 审查作者视角数据关联
	auditAuthorCmd = &cobra.Command{
		Use:   "audit-author",
		Short: "审查作者视角数据关联完整性",
		Run:   runAuthorAuditCmd,
	}

	// settingsCmd 填充用户设置数据
	settingsCmd = &cobra.Command{
		Use:   "settings",
		Short: "填充用户设置数据",
		Run:   runSettings,
	}

	// e2eCmd E2E测试专用一键填充命令
	e2eCmd = &cobra.Command{
		Use:   "e2e",
		Short: "E2E测试专用：一键填充所有测试数据（清空+重建）",
		Long: `E2E测试专用命令，一键填充所有测试数据。
包含：用户、分类、书籍、章节、社交、阅读、钱包等完整数据。
适用于Playwright E2E多链路联调测试。`,
		Run: runE2E,
	}

	// demoCmd 演示数据填充命令
	demoCmd = &cobra.Command{
		Use:   "demo",
		Short: "创建演示数据（科幻故事：星际觉醒）",
		Long: `创建完整的演示数据项目"星际觉醒"。
包含：1个项目、4卷24章、12个角色、30+关系、8道具、6地点、3时间线。
适用于产品演示和功能测试。`,
		Run: runDemo,
	}
)

// init 初始化命令
func init() {
	cobra.OnInitialize(initConfig)

	// 根命令标志
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件路径 (默认为 $HOME/.seeder.yaml)")
	rootCmd.PersistentFlags().StringVarP(&scale, "scale", "s", "medium", "数据规模: small, medium, large")
	rootCmd.PersistentFlags().BoolVarP(&clean, "clean", "c", false, "填充前清空现有数据")

	// 添加子命令
	rootCmd.AddCommand(allCmd)
	rootCmd.AddCommand(baselineCmd)
	rootCmd.AddCommand(fullCmd)
	rootCmd.AddCommand(usersCmd)
	rootCmd.AddCommand(categoriesCmd)
	rootCmd.AddCommand(bookstoreCmd)
	rootCmd.AddCommand(showcaseCmd)
	rootCmd.AddCommand(subscriptionsCmd)
	rootCmd.AddCommand(chaptersCmd)
	rootCmd.AddCommand(socialCmd)
	rootCmd.AddCommand(walletsCmd)
	rootCmd.AddCommand(rankingsCmd)
	rootCmd.AddCommand(aiQuotaCmd)
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(readerCmd)
	rootCmd.AddCommand(notificationsCmd)
	rootCmd.AddCommand(messagingCmd)
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(financeCmd)
	rootCmd.AddCommand(booklistsCmd)
	rootCmd.AddCommand(collectionFoldersCmd)
	rootCmd.AddCommand(withdrawalsCmd)
	rootCmd.AddCommand(auditReaderCmd)
	rootCmd.AddCommand(auditAuthorCmd)
	rootCmd.AddCommand(settingsCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(verifyCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(publicationFlowCmd)
	rootCmd.AddCommand(e2eCmd)
	rootCmd.AddCommand(demoCmd)

	// demoCmd 专属标志
	demoCmd.Flags().BoolVar(&forceCleanDemo, "force-clean", false, "强制清空后重新创建演示数据")
}

// initConfig 初始化配置
func initConfig() {
	cfgCopy := *config.DefaultConfig
	cfg = &cfgCopy

	// Allow CI/test workflows to inject isolated MongoDB settings without
	// requiring a seeder-specific config file.
	if uri := strings.TrimSpace(os.Getenv("QINGYU_DATABASE_PRIMARY_MONGODB_URI")); uri != "" {
		cfg.MongoDB.URI = uri
	} else if uri := strings.TrimSpace(os.Getenv("MONGODB_URI")); uri != "" {
		cfg.MongoDB.URI = uri
	}

	if database := strings.TrimSpace(os.Getenv("QINGYU_DATABASE_PRIMARY_MONGODB_DATABASE")); database != "" {
		cfg.MongoDB.Database = database
	} else if database := strings.TrimSpace(os.Getenv("MONGODB_DATABASE")); database != "" {
		cfg.MongoDB.Database = database
	}

	if scale != "" {
		cfg.Scale = scale
	}
	cfg.Clean = clean
}

// main 主函数
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

// getDatabase 连接 MongoDB 数据库
func getDatabase() (*utils.Database, error) {
	db, err := utils.ConnectDB(cfg.MongoDB.URI, cfg.MongoDB.Database)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}
	return db, nil
}

// runAll 执行所有填充操作
func runAll(cmd *cobra.Command, args []string) {
	runBaseline(cmd, args)
}

func runBaseline(cmd *cobra.Command, args []string) {
	fmt.Println("开始填充联调基线测试数据...")
	runSeedWorkflow(seedWorkflowOptions{
		cleanMessage: "清空现有基线数据...",
		cleanFn:      cleanAllData,
		extraSteps:   nil,
	})
}

func runFull(cmd *cobra.Command, args []string) {
	fmt.Println("开始填充完整测试数据...")
	runSeedWorkflow(seedWorkflowOptions{
		cleanMessage: "清空现有完整测试数据...",
		cleanFn:      cleanAllData,
		extraSteps: []seedStep{
			{title: "填充钱包数据...", run: seedWalletsData},
			{title: "填充通知数据...", run: seedNotificationsData},
			{title: "填充消息数据...", run: seedMessagingData},
			{title: "填充财务数据...", run: seedFinanceData},
			{title: "激活AI配额...", run: seedAIQuotaData},
			{title: "填充书单数据...", run: seedBooklistsData},
			{title: "填充收藏夹数据...", run: seedCollectionFoldersData},
			{title: "填充提现申请数据...", run: seedWithdrawalsData},
		},
	})
}

// runUsers 只填充用户数据
func runUsers(cmd *cobra.Command, args []string) {
	fmt.Println("开始填充用户数据...")
	fmt.Printf("数据规模: %s\n", cfg.Scale)

	db, err := getDatabase()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	if cfg.Clean {
		fmt.Println("\n清空用户数据...")
		if err := seedUsersClean(db); err != nil {
			fmt.Printf("清空用户数据失败: %v\n", err)
			os.Exit(1)
		}
	}

	if err := seedUsers(db); err != nil {
		fmt.Printf("填充用户数据失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n用户数据填充完成!")
}

// runBookstore 只填充书籍数据
func runBookstore(cmd *cobra.Command, args []string) {
	fmt.Println("开始填充书籍数据...")
	fmt.Printf("数据规模: %s\n", cfg.Scale)

	db, err := getDatabase()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	// 确保有足够的author用户
	if err := ensureAuthorUsers(db); err != nil {
		fmt.Printf("确保author用户存在失败: %v\n", err)
		os.Exit(1)
	}

	if cfg.Clean {
		fmt.Println("\n清空书籍数据...")
		if err := seedBooksClean(db); err != nil {
			fmt.Printf("清空书籍数据失败: %v\n", err)
			os.Exit(1)
		}
		if err := seedCategoriesClean(db); err != nil {
			fmt.Printf("清空分类数据失败: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("\n填充分类数据...")
	if err := seedCategories(db); err != nil {
		fmt.Printf("填充分类数据失败: %v\n", err)
		os.Exit(1)
	}

	if err := seedBooks(db); err != nil {
		fmt.Printf("填充书籍数据失败: %v\n", err)
		os.Exit(1)
	}

	// 自动生成榜单数据
	fmt.Println("\n生成榜单数据...")
	if err := seedRankings(db); err != nil {
		fmt.Printf("生成榜单数据失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n书籍和榜单数据填充完成!")
}

// runShowcase 只填充精选演示书籍和榜单
func runShowcase(cmd *cobra.Command, args []string) {
	fmt.Println("开始填充精选演示书籍数据...")

	db, err := getDatabase()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	if err := ensureAuthorUsers(db); err != nil {
		fmt.Printf("确保author用户存在失败: %v\n", err)
		os.Exit(1)
	}

	if cfg.Clean {
		fmt.Println("\n清空书籍和分类数据...")
		if err := seedBooksClean(db); err != nil {
			fmt.Printf("清空书籍数据失败: %v\n", err)
			os.Exit(1)
		}
		if err := seedCategoriesClean(db); err != nil {
			fmt.Printf("清空分类数据失败: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("\n填充分类数据...")
	if err := seedCategories(db); err != nil {
		fmt.Printf("填充分类数据失败: %v\n", err)
		os.Exit(1)
	}

	if err := seedShowcaseBooks(db); err != nil {
		fmt.Printf("填充精选演示书籍失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n生成精选 banner...")
	if err := seedShowcaseBanners(db); err != nil {
		fmt.Printf("生成精选 banner 失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n生成精选章节和内容...")
	if err := seedShowcaseChapters(db); err != nil {
		fmt.Printf("生成精选章节失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n生成榜单数据...")
	if err := seedRankings(db); err != nil {
		fmt.Printf("生成榜单数据失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n精选演示书籍和榜单数据填充完成!")
}

// runClean 清空所有数据
func runClean(cmd *cobra.Command, args []string) {
	fmt.Println("警告: 此操作将清空所有测试数据!")
	fmt.Print("请输入 'YES' 确认: ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("读取输入失败: %v\n", err)
		os.Exit(1)
	}

	input = strings.TrimSpace(input)
	if input != "YES" {
		fmt.Println("操作已取消")
		return
	}

	fmt.Println("\n清空所有测试数据...")

	db, err := getDatabase()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	if err := cleanAllData(db); err != nil {
		fmt.Printf("清空数据失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("数据清空完成!")
}

// runVerify 验证数据完整性
func runVerify(cmd *cobra.Command, args []string) {
	fmt.Println("验证数据完整性...")
	fmt.Println()

	db, err := getDatabase()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	// 使用新的验证器
	if err := validateData(db); err != nil {
		fmt.Printf("验证失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n数据验证完成!")
}

// cleanAllData 清空所有数据
func cleanAllData(db *utils.Database) error {
	if err := seedUsersClean(db); err != nil {
		return err
	}
	if err := seedSettingsClean(db); err != nil {
		return err
	}
	if err := seedCategoriesClean(db); err != nil {
		return err
	}
	if err := seedBooksClean(db); err != nil {
		return err
	}
	if err := seedRankingsClean(db); err != nil {
		return err
	}
	if err := seedChaptersClean(db); err != nil {
		return err
	}
	if err := seedSubscriptionsClean(db); err != nil {
		return err
	}
	if err := seedSocialClean(db); err != nil {
		return err
	}
	if err := seedReaderClean(db); err != nil {
		return err
	}
	if err := seedStatsClean(db); err != nil {
		return err
	}
	return nil
}

type seedStep struct {
	title string
	run   func(db *utils.Database) error
}

type seedWorkflowOptions struct {
	cleanMessage string
	cleanFn      func(db *utils.Database) error
	extraSteps   []seedStep
}

func runSeedWorkflow(opts seedWorkflowOptions) {
	fmt.Printf("数据规模: %s\n", cfg.Scale)

	db, err := getDatabase()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	if cfg.Clean && opts.cleanFn != nil {
		fmt.Printf("\n%s\n", opts.cleanMessage)
		if err := opts.cleanFn(db); err != nil {
			fmt.Printf("清空数据失败: %v\n", err)
			os.Exit(1)
		}
	}

	coreSteps := []seedStep{
		{title: "填充用户数据...", run: seedUsers},
		{title: "填充用户设置数据...", run: seedSettings},
		{title: "填充分类数据...", run: seedCategories},
		{title: "填充书籍数据...", run: seedBooks},
		{title: "生成榜单数据...", run: seedRankings},
		{title: "填充章节和章节内容数据...", run: seedChaptersData},
		{title: "填充订阅关系...", run: seedSubscriptions},
		{title: "填充社交数据...", run: seedSocialData},
		{title: "填充阅读数据...", run: seedReaderData},
		{title: "填充统计数据...", run: seedStatsData},
	}

	steps := append(coreSteps, opts.extraSteps...)
	for _, step := range steps {
		fmt.Printf("\n%s\n", step.title)
		if err := step.run(db); err != nil {
			fmt.Printf("%s失败: %v\n", strings.TrimSuffix(step.title, "..."), err)
			os.Exit(1)
		}
	}

	fmt.Println("\n验证数据完整性...")
	if err := validateData(db); err != nil {
		fmt.Printf("数据验证失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n数据填充完成!")
}

// seedUsers 填充用户数据
func seedUsers(db *utils.Database) error {
	seeder := NewUserSeeder(db, cfg)

	// 填充真实测试用户
	if err := seeder.SeedRealUsers(); err != nil {
		return err
	}

	// 填充生成的用户
	if err := seeder.SeedGeneratedUsers(); err != nil {
		return err
	}

	// 统计用户数量
	count, err := seeder.Count()
	if err != nil {
		return err
	}

	fmt.Printf("用户填充完成，共 %d 个用户\n", count)
	return nil
}

// seedUsersClean 清空用户数据
func seedUsersClean(db *utils.Database) error {
	seeder := NewUserSeeder(db, cfg)
	return seeder.Clean()
}

// seedCategories 填充分类数据
func seedCategories(db *utils.Database) error {
	seeder := NewCategorySeeder(db, cfg)
	return seeder.SeedCategories()
}

// seedCategoriesClean 清空分类数据
func seedCategoriesClean(db *utils.Database) error {
	seeder := NewCategorySeeder(db, cfg)
	return seeder.Clean()
}

// seedBooks 填充书籍数据
func seedBooks(db *utils.Database) error {
	seeder := NewBookstoreSeeder(db, cfg)

	// 填充书籍数据
	if err := seeder.SeedGeneratedBooks(); err != nil {
		return err
	}

	// 填充 banner 数据
	if err := seeder.SeedBanners(); err != nil {
		return err
	}

	// 统计书籍数量
	count, err := seeder.Count()
	if err != nil {
		return err
	}

	fmt.Printf("书籍填充完成，共 %d 本书\n", count)
	return nil
}

func seedShowcaseBooks(db *utils.Database) error {
	seeder := NewBookstoreSeeder(db, cfg)
	return seeder.SeedShowcaseBooks()
}

func seedShowcaseBanners(db *utils.Database) error {
	seeder := NewBookstoreSeeder(db, cfg)
	return seeder.SeedBanners()
}

// seedBooksClean 清空书籍数据
func seedBooksClean(db *utils.Database) error {
	seeder := NewBookstoreSeeder(db, cfg)
	return seeder.Clean()
}

// seedSubscriptions 填充用户订阅书籍关系
func seedSubscriptions(db *utils.Database) error {
	seeder := NewSubscriptionSeeder(db, cfg)
	return seeder.SeedSubscriptions()
}

// seedSubscriptionsClean 清空订阅关系
func seedSubscriptionsClean(db *utils.Database) error {
	seeder := NewSubscriptionSeeder(db, cfg)
	return seeder.Clean()
}

func seedChaptersData(db *utils.Database) error {
	seeder := NewChapterSeeder(db, cfg)
	if err := seeder.SeedChapters(); err != nil {
		return err
	}
	return seeder.SeedChapterContents()
}

func seedShowcaseChapters(db *utils.Database) error {
	seeder := NewChapterSeeder(db, cfg)
	if err := seeder.SeedShowcaseChapters(); err != nil {
		return err
	}
	return seeder.SeedShowcaseChapterContents()
}

func seedChaptersClean(db *utils.Database) error {
	seeder := NewChapterSeeder(db, cfg)
	return seeder.Clean()
}

func seedSocialData(db *utils.Database) error {
	seeder := NewSocialSeeder(db, cfg)
	return seeder.SeedSocialData()
}

func seedSocialClean(db *utils.Database) error {
	seeder := NewSocialSeeder(db, cfg)
	return seeder.Clean()
}

func seedReaderData(db *utils.Database) error {
	seeder := NewReaderSeeder(db, cfg)
	return seeder.SeedReadingData()
}

func seedReaderClean(db *utils.Database) error {
	seeder := NewReaderSeeder(db, cfg)
	return seeder.Clean()
}

func seedStatsData(db *utils.Database) error {
	seeder := NewStatsSeeder(db, cfg)
	return seeder.SeedStats()
}

func seedStatsClean(db *utils.Database) error {
	seeder := NewStatsSeeder(db, cfg)
	return seeder.Clean()
}

func seedWalletsData(db *utils.Database) error {
	seeder := NewWalletSeeder(db, cfg)
	return seeder.SeedWallets()
}

func seedNotificationsData(db *utils.Database) error {
	seeder := NewNotificationSeeder(db, cfg)
	return seeder.SeedNotifications()
}

func seedMessagingData(db *utils.Database) error {
	seeder := NewMessagingSeeder(db, cfg)
	return seeder.SeedMessagingData()
}

func seedFinanceData(db *utils.Database) error {
	seeder := NewFinanceSeeder(db, cfg)
	return seeder.SeedFinanceData()
}

func seedAIQuotaData(db *utils.Database) error {
	seeder := NewAIQuotaSeeder(db, cfg)
	return seeder.SeedAIQuota()
}

func seedBooklistsData(db *utils.Database) error {
	seeder := NewBooklistSeeder(db, cfg)
	return seeder.SeedBooklists()
}

func seedCollectionFoldersData(db *utils.Database) error {
	seeder := NewCollectionFolderSeeder(db, cfg)
	return seeder.SeedCollectionFolders()
}

func seedWithdrawalsData(db *utils.Database) error {
	seeder := NewWithdrawalSeeder(db, cfg)
	return seeder.SeedWithdrawals()
}

func runBooklists(cmd *cobra.Command, args []string) {
	fmt.Println("开始填充书单数据...")
	db, err := getDatabase()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()
	seeder := NewBooklistSeeder(db, cfg)
	if err := seeder.SeedBooklists(); err != nil {
		fmt.Printf("填充书单数据失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("\n书单数据填充完成!")
}

func runCollectionFolders(cmd *cobra.Command, args []string) {
	fmt.Println("开始填充收藏夹数据...")
	db, err := getDatabase()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()
	seeder := NewCollectionFolderSeeder(db, cfg)
	if err := seeder.SeedCollectionFolders(); err != nil {
		fmt.Printf("填充收藏夹数据失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("\n收藏夹数据填充完成!")
}

func runWithdrawals(cmd *cobra.Command, args []string) {
	fmt.Println("开始填充提现申请数据...")
	db, err := getDatabase()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()
	seeder := NewWithdrawalSeeder(db, cfg)
	if err := seeder.SeedWithdrawals(); err != nil {
		fmt.Printf("填充提现申请数据失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("\n提现申请数据填充完成!")
}

// runTestData 填充E2E测试数据
func runTestData(cmd *cobra.Command, args []string) {
	fmt.Println("开始填充E2E测试数据...")

	db, err := getDatabase()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	// 如果需要清空数据
	if cfg.Clean {
		fmt.Println("\n清空现有测试数据...")
		seeder := NewTestDataSeeder(db)
		if err := seeder.Clean(); err != nil {
			fmt.Printf("清空测试数据失败: %v\n", err)
			os.Exit(1)
		}
	}

	seeder := NewTestDataSeeder(db)

	if err := seeder.SeedTestData(); err != nil {
		fmt.Printf("填充测试数据失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✅ E2E测试数据填充完成!")
	fmt.Println("\n测试账号: testuser / 123456")
	fmt.Println("测试书籍: 修仙世界、修仙归来、万古修仙等")
	fmt.Println("测试分类: 玄幻、修仙")
}

// runPublicationFlow 通过完整发布流程创建测试数据
func runPublicationFlow(cmd *cobra.Command, args []string) {
	fmt.Println("开始通过完整发布流程创建测试数据...")

	db, err := getDatabase()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	seeder := NewPublicationFlowSeeder(db, cfg)

	// 如果需要清空数据
	if cfg.Clean {
		fmt.Println("\n清空现有测试数据...")
		if err := seeder.Clean(); err != nil {
			fmt.Printf("清空数据失败: %v\n", err)
			os.Exit(1)
		}
	}

	// 获取参数（默认每个作者2本书，每本书3个章节）
	booksPerAuthor := 2
	chaptersPerBook := 3

	if err := seeder.SeedPublicationFlow(booksPerAuthor, chaptersPerBook); err != nil {
		fmt.Printf("\n❌ 发布流程数据创建失败: %v\n", err)
		os.Exit(1)
	}
}

// runChapters 填充章节数据
func runChapters(cmd *cobra.Command, args []string) {
	fmt.Println("开始填充章节数据...")

	db, err := getDatabase()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	if cfg.Clean {
		fmt.Println("\n清空章节数据...")
		seeder := NewChapterSeeder(db, cfg)
		if err := seeder.Clean(); err != nil {
			fmt.Printf("清空章节数据失败: %v\n", err)
			os.Exit(1)
		}
	}

	seeder := NewChapterSeeder(db, cfg)

	fmt.Println("\n生成章节...")
	if err := seeder.SeedChapters(); err != nil {
		fmt.Printf("生成章节失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n生成章节内容...")
	if err := seeder.SeedChapterContents(); err != nil {
		fmt.Printf("生成章节内容失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n章节数据填充完成!")
}

// runCategories 填充分类数据
func runCategories(cmd *cobra.Command, args []string) {
	fmt.Println("开始填充分类数据...")

	db, err := getDatabase()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	if cfg.Clean {
		fmt.Println("\n清空分类数据...")
		if err := seedCategoriesClean(db); err != nil {
			fmt.Printf("清空分类数据失败: %v\n", err)
			os.Exit(1)
		}
	}

	if err := seedCategories(db); err != nil {
		fmt.Printf("填充分类数据失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n分类数据填充完成!")
}

func runSubscriptions(cmd *cobra.Command, args []string) {
	fmt.Println("开始填充书籍订阅关系...")

	db, err := getDatabase()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	if cfg.Clean {
		fmt.Println("\n清空订阅关系数据...")
		if err := seedSubscriptionsClean(db); err != nil {
			fmt.Printf("清空订阅关系失败: %v\n", err)
			os.Exit(1)
		}
	}

	if err := seedSubscriptions(db); err != nil {
		fmt.Printf("填充订阅关系失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n订阅关系填充完成!")
}

// runSocial 填充社交数据
func runSocial(cmd *cobra.Command, args []string) {
	fmt.Println("开始填充社交数据...")

	db, err := getDatabase()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	seeder := NewSocialSeeder(db, cfg)

	if cfg.Clean {
		fmt.Println("\n清空社交数据...")
		if err := seeder.Clean(); err != nil {
			fmt.Printf("清空社交数据失败: %v\n", err)
			os.Exit(1)
		}
	}

	if err := seeder.SeedSocialData(); err != nil {
		fmt.Printf("填充社交数据失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n社交数据填充完成!")
}

// runWallets 填充钱包数据
func runWallets(cmd *cobra.Command, args []string) {
	fmt.Println("开始填充钱包数据...")

	db, err := getDatabase()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	seeder := NewWalletSeeder(db, cfg)

	if err := seeder.SeedWallets(); err != nil {
		fmt.Printf("填充钱包数据失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n钱包数据填充完成!")
}

// runRankings 填充榜单数据
func runRankings(cmd *cobra.Command, args []string) {
	fmt.Println("开始填充榜单数据...")

	db, err := getDatabase()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	seeder := NewRankingSeeder(db, cfg)

	if cfg.Clean {
		fmt.Println("\n清空榜单数据...")
		if err := seeder.Clean(); err != nil {
			fmt.Printf("清空榜单数据失败: %v\n", err)
			os.Exit(1)
		}
	}

	if err := seeder.SeedRankings(); err != nil {
		fmt.Printf("填充榜单数据失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n榜单数据填充完成!")
}

// runAIQuota 激活AI配额
func runAIQuota(cmd *cobra.Command, args []string) {
	fmt.Println("开始激活AI配额...")

	db, err := getDatabase()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	seeder := NewAIQuotaSeeder(db, cfg)

	if err := seeder.SeedAIQuota(); err != nil {
		fmt.Printf("激活AI配额失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nAI配额激活完成!")
}

// runImport 导入小说数据
func runImport(cmd *cobra.Command, args []string) {
	fmt.Println("开始导入小说数据...")

	db, err := getDatabase()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	seeder := NewImportSeeder(db, cfg)

	if err := seeder.SeedFromJSON("data/novels_100.json"); err != nil {
		fmt.Printf("导入小说数据失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n小说数据导入完成!")
}

// seedRankings 填充榜单数据
func seedRankings(db *utils.Database) error {
	seeder := NewRankingSeeder(db, cfg)
	return seeder.SeedRankings()
}

// seedRankingsClean 清空榜单数据
func seedRankingsClean(db *utils.Database) error {
	seeder := NewRankingSeeder(db, cfg)
	return seeder.Clean()
}

// runReader 填充阅读数据
func runReader(cmd *cobra.Command, args []string) {
	fmt.Println("开始填充阅读数据...")

	db, err := getDatabase()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	if cfg.Clean {
		fmt.Println("\n清空阅读数据...")
		seeder := NewReaderSeeder(db, cfg)
		if err := seeder.Clean(); err != nil {
			fmt.Printf("清空阅读数据失败: %v\n", err)
			os.Exit(1)
		}
	}

	seeder := NewReaderSeeder(db, cfg)

	if err := seeder.SeedReadingData(); err != nil {
		fmt.Printf("填充阅读数据失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n阅读数据填充完成!")
}

// runNotifications 填充通知数据
func runNotifications(cmd *cobra.Command, args []string) {
	fmt.Println("开始填充通知数据...")

	db, err := getDatabase()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	if cfg.Clean {
		fmt.Println("\n清空通知数据...")
		seeder := NewNotificationSeeder(db, cfg)
		if err := seeder.Clean(); err != nil {
			fmt.Printf("清空通知数据失败: %v\n", err)
			os.Exit(1)
		}
	}

	seeder := NewNotificationSeeder(db, cfg)

	if err := seeder.SeedNotifications(); err != nil {
		fmt.Printf("填充通知数据失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n通知数据填充完成!")
}

// runMessaging 填充消息数据
func runMessaging(cmd *cobra.Command, args []string) {
	fmt.Println("开始填充消息数据...")

	db, err := getDatabase()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	if cfg.Clean {
		fmt.Println("\n清空消息数据...")
		seeder := NewMessagingSeeder(db, cfg)
		if err := seeder.Clean(); err != nil {
			fmt.Printf("清空消息数据失败: %v\n", err)
			os.Exit(1)
		}
	}

	seeder := NewMessagingSeeder(db, cfg)
	if err := seeder.SeedMessagingData(); err != nil {
		fmt.Printf("填充消息数据失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n消息数据填充完成!")
}

// runStats 填充统计数据
func runStats(cmd *cobra.Command, args []string) {
	fmt.Println("开始填充统计数据...")

	db, err := getDatabase()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	if cfg.Clean {
		fmt.Println("\n清空统计数据...")
		seeder := NewStatsSeeder(db, cfg)
		if err := seeder.Clean(); err != nil {
			fmt.Printf("清空统计数据失败: %v\n", err)
			os.Exit(1)
		}
	}

	seeder := NewStatsSeeder(db, cfg)
	if err := seeder.SeedStats(); err != nil {
		fmt.Printf("填充统计数据失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n统计数据填充完成!")
}

// runFinance 填充财务数据
func runFinance(cmd *cobra.Command, args []string) {
	fmt.Println("开始填充财务数据...")

	db, err := getDatabase()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	if cfg.Clean {
		fmt.Println("\n清空财务数据...")
		seeder := NewFinanceSeeder(db, cfg)
		if err := seeder.Clean(); err != nil {
			fmt.Printf("清空财务数据失败: %v\n", err)
			os.Exit(1)
		}
	}

	seeder := NewFinanceSeeder(db, cfg)
	if err := seeder.SeedFinanceData(); err != nil {
		fmt.Printf("填充财务数据失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n财务数据填充完成!")
}

// runAuditReader 审查读者视角数据关联
func runAuditReader(cmd *cobra.Command, args []string) {
	fmt.Println("🔍 开始读者视角数据关联审查...")

	db, err := getDatabase()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	if err := RunReaderAudit(db.Database); err != nil {
		fmt.Printf("审查执行失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✅ 读者视角数据关联审查完成!")
	fmt.Println("📄 报告已输出到控制台")
}

// runAuthorAuditCmd 审查作者视角数据关联
func runAuthorAuditCmd(cmd *cobra.Command, args []string) {
	fmt.Println("🚀 开始执行作者视角数据关联审查...")
	fmt.Println()

	db, err := getDatabase()
	if err != nil {
		fmt.Printf("❌ 数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	report, err := RunAuthorAudit(db)
	if err != nil {
		fmt.Printf("❌ 审查失败: %v\n", err)
		os.Exit(1)
	}

	report.PrintReport()

	// 根据评分返回退出码
	if report.TotalScore == "差 (D)" {
		os.Exit(1)
	}
}

// validateData 验证数据完整性
func validateData(db *utils.Database) error {
	v := validator.NewDataValidator(db.Database)

	report, err := v.ValidateRelationships(context.Background())
	if err != nil {
		return fmt.Errorf("验证执行失败: %w", err)
	}

	// 打印验证结果
	fmt.Println("========== 数据验证报告 ==========")
	fmt.Printf("验证时间: %s\n", report.ValidatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("验证结果: ")
	if report.IsValid {
		fmt.Println("✅ 通过")
	} else {
		fmt.Println("❌ 失败")
	}

	// 打印集合统计信息
	if len(report.CollectionStats) > 0 {
		fmt.Println("\n集合统计:")
		for coll, count := range report.CollectionStats {
			fmt.Printf("  - %s: %d 条记录\n", coll, count)
		}
	}

	// 打印孤儿记录详情
	if report.TotalOrphanedRecords > 0 {
		fmt.Printf("\n孤儿记录 (共 %d 条):\n", report.TotalOrphanedRecords)
		for _, detail := range report.OrphanDetails {
			fmt.Printf("  - %s\n", detail)
		}
	}

	// 打印ID格式不一致问题
	if len(report.InconsistentRecords) > 0 {
		fmt.Printf("\nID格式不一致 (共 %d 条):\n", len(report.InconsistentRecords))
		for _, issue := range report.InconsistentRecords {
			fmt.Printf("  - %s\n", issue)
		}
	}

	// 打印错误信息
	if len(report.Errors) > 0 {
		fmt.Printf("\n错误 (共 %d 个):\n", len(report.Errors))
		for _, err := range report.Errors {
			fmt.Printf("  - %v\n", err)
		}
	}

	// 打印摘要
	fmt.Printf("\n%s\n", report.Summary)
	fmt.Println("=================================")

	// 如果验证失败，返回错误
	if !report.IsValid {
		return fmt.Errorf("数据验证失败: %s", report.Summary)
	}

	return nil
}

// ensureAuthorUsers 确保有足够的author用户
func ensureAuthorUsers(db *utils.Database) error {
	ctx := context.Background()

	// 检查现有author数量
	count, err := db.Collection("users").CountDocuments(ctx, bson.M{"role": "author"})
	if err != nil {
		return fmt.Errorf("检查author用户失败: %w", err)
	}

	// 获取配置中期望的author数量
	scale := config.GetScaleConfig(cfg.Scale)
	minAuthors := int64(scale.Authors)

	if count >= minAuthors {
		fmt.Printf("✓ 已有 %d 个author用户\n", count)
		return nil
	}

	// 不足则生成
	needed := minAuthors - count
	fmt.Printf("需要生成 %d 个author用户...\n", needed)

	seeder := NewUserSeeder(db, cfg)

	// 只生成author角色用户
	authors := seeder.GetGenerator().GenerateUsers(int(needed), "author")
	if err := seeder.GetInserter().InsertMany(ctx, authors); err != nil {
		return fmt.Errorf("生成author用户失败: %w", err)
	}

	fmt.Printf("✓ 成功生成 %d 个author用户\n", needed)
	return nil
}

// runSettings 填充用户设置数据
func runSettings(cmd *cobra.Command, args []string) {
	fmt.Println("开始填充用户设置数据...")

	db, err := getDatabase()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	if cfg.Clean {
		fmt.Println("\n清空用户设置数据...")
		seeder := NewSettingsSeeder(db, cfg)
		if err := seeder.Clean(); err != nil {
			fmt.Printf("清空用户设置数据失败: %v\n", err)
			os.Exit(1)
		}
	}

	seeder := NewSettingsSeeder(db, cfg)
	if err := seeder.SeedUserSettings(); err != nil {
		fmt.Printf("填充用户设置数据失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n用户设置数据填充完成!")
}

// seedSettings 填充用户设置数据
func seedSettings(db *utils.Database) error {
	seeder := NewSettingsSeeder(db, cfg)
	return seeder.SeedUserSettings()
}

// seedSettingsClean 清空用户设置数据
func seedSettingsClean(db *utils.Database) error {
	seeder := NewSettingsSeeder(db, cfg)
	return seeder.Clean()
}

// runE2E E2E测试专用一键填充命令
func runE2E(cmd *cobra.Command, args []string) {
	fmt.Println("🚀 开始E2E测试数据一键填充...")
	fmt.Println("=" + strings.Repeat("=", 50))
	fmt.Println()

	db, err := getDatabase()
	if err != nil {
		fmt.Printf("❌ 数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	// 步骤1: 清空所有数据
	fmt.Println("📋 步骤1: 清空现有数据...")
	if err := cleanAllData(db); err != nil {
		fmt.Printf("❌ 清空数据失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("   ✓ 数据清空完成")

	// 步骤2: 填充核心数据
	steps := []struct {
		name string
		fn   func(*utils.Database) error
	}{
		{"用户数据", seedUsers},
		{"用户设置", seedSettings},
		{"分类数据", seedCategories},
		{"书籍数据", seedBooks},
		{"榜单数据", seedRankings},
		{"章节数据", seedChaptersData},
		{"订阅关系", seedSubscriptions},
		{"社交数据", seedSocialData},
		{"阅读数据", seedReaderData},
		{"统计数据", seedStatsData},
		{"钱包数据", seedWalletsData},
		{"消息数据", seedMessagingData},
	}

	for i, step := range steps {
		fmt.Printf("📋 步骤%d/%d: 填充%s...\n", i+2, len(steps)+1, step.name)
		if err := step.fn(db); err != nil {
			fmt.Printf("❌ 填充%s失败: %v\n", step.name, err)
			os.Exit(1)
		}
		fmt.Printf("   ✓ %s完成\n", step.name)
	}

	// 步骤3: 验证数据
	fmt.Println("\n📋 步骤" + fmt.Sprintf("%d", len(steps)+2) + ": 验证数据完整性...")
	if err := validateData(db); err != nil {
		fmt.Printf("⚠️  数据验证警告: %v\n", err)
	}

	// 输出测试账号信息
	fmt.Println("\n" + strings.Repeat("=", 52))
	fmt.Println("✅ E2E测试数据填充完成!")
	fmt.Println(strings.Repeat("=", 52))
	fmt.Println("\n📝 测试账号信息:")
	fmt.Println("   读者账号: testuser001 / password")
	fmt.Println("   读者账号: testuser002 / password")
	fmt.Println("   读者账号: testuser003 / password")
	fmt.Println("   读者账号: testuser004 / password")
	fmt.Println("   作者账号: testauthor001 / password")
	fmt.Println("   管理账号: testadmin001 / password")
	fmt.Println("\n💡 使用方式:")
	fmt.Println("   cd Qingyu_backend && go run cmd/seeder/main.go e2e")
}

// runDemo 演示数据填充命令
func runDemo(cmd *cobra.Command, args []string) {
	fmt.Println("🚀 开始创建演示数据...")
	fmt.Println("=" + strings.Repeat("=", 50))
	fmt.Println()
	fmt.Println("📖 项目名称: 星际觉醒")
	fmt.Println("🎭 类型: 科幻小说")
	fmt.Println("📊 规模: 中型")
	fmt.Println()

	db, err := getDatabase()
	if err != nil {
		fmt.Printf("❌ 数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	seeder := NewDemoDataSeeder(db)

	// 如果指定了 force-clean，先清理
	if forceCleanDemo {
		fmt.Println("⚠️  强制清空模式，将删除已存在的演示数据...")
		if err := seeder.Clean(); err != nil {
			fmt.Printf("❌ 清空数据失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Println()
	}

	// 创建演示数据
	if err := seeder.CreateDemoData(); err != nil {
		fmt.Printf("❌ 创建演示数据失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n" + strings.Repeat("=", 52))
	fmt.Println("✅ 演示数据创建完成!")
	fmt.Println(strings.Repeat("=", 52))
	fmt.Println("\n📝 登录信息:")
	fmt.Println("   管理员: admin / qingyu2024")
	fmt.Println("\n📊 数据统计:")
	fmt.Println("   项目: 1个（星际觉醒）")
	fmt.Println("   文档: 28个（4卷 + 24章）")
	fmt.Println("   大纲: 21个（1级1个 + 2级4个 + 3级16个）")
	fmt.Println("   角色: 12个")
	fmt.Println("   关系: 30+条边")
	fmt.Println("   道具: 8个")
	fmt.Println("   地点: 6个")
	fmt.Println("   时间线: 3条")
	fmt.Println("\n💡 使用方式:")
	fmt.Println("   cd Qingyu_backend && go run cmd/seeder/main.go demo")
}
