// Package main provides the test data seeder for Qingyu writing platform
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"Qingyu_backend/cmd/seeder/config"
	"Qingyu_backend/cmd/seeder/utils"

	"github.com/spf13/cobra"
)

var (
	// 命令行标志
	cfgFile string
	scale   string
	clean   bool

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
		Short: "执行所有数据填充",
		Run:   runAll,
	}

	// usersCmd 只填充用户数据
	usersCmd = &cobra.Command{
		Use:   "users",
		Short: "只填充用户数据",
		Run:   runUsers,
	}

	// bookstoreCmd 只填充书籍数据
	bookstoreCmd = &cobra.Command{
		Use:   "bookstore",
		Short: "只填充书籍数据",
		Run:   runBookstore,
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
	rootCmd.AddCommand(usersCmd)
	rootCmd.AddCommand(bookstoreCmd)
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
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(verifyCmd)
	rootCmd.AddCommand(testCmd)
}

// initConfig 初始化配置
func initConfig() {
	cfg = config.DefaultConfig
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
	fmt.Println("开始填充所有测试数据...")
	fmt.Printf("数据规模: %s\n", cfg.Scale)

	db, err := getDatabase()
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	// 如果需要，先清空数据
	if cfg.Clean {
		fmt.Println("\n清空现有数据...")
		if err := cleanAllData(db); err != nil {
			fmt.Printf("清空数据失败: %v\n", err)
			os.Exit(1)
		}
	}

	// 填充用户
	fmt.Println("\n填充用户数据...")
	if err := seedUsers(db); err != nil {
		fmt.Printf("填充用户数据失败: %v\n", err)
		os.Exit(1)
	}

	// 填充书籍
	fmt.Println("\n填充书籍数据...")
	if err := seedBooks(db); err != nil {
		fmt.Printf("填充书籍数据失败: %v\n", err)
		os.Exit(1)
	}

	// 生成榜单数据
	fmt.Println("\n生成榜单数据...")
	if err := seedRankings(db); err != nil {
		fmt.Printf("生成榜单数据失败: %v\n", err)
		os.Exit(1)
	}

	// 填充订阅关系
	fmt.Println("\n填充订阅关系...")
	if err := seedSubscriptions(db); err != nil {
		fmt.Printf("填充订阅关系失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n数据填充完成!")
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

	if cfg.Clean {
		fmt.Println("\n清空书籍数据...")
		if err := seedBooksClean(db); err != nil {
			fmt.Printf("清空书籍数据失败: %v\n", err)
			os.Exit(1)
		}
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

	reports := utils.VerifyData(db)

	// 显示验证结果
	for _, report := range reports {
		if report.Passed {
			fmt.Printf("✅ %s: 通过\n", report.Category)
		} else {
			fmt.Printf("❌ %s: 失败\n", report.Category)
		}

		for _, issue := range report.Issues {
			fmt.Printf("   - %s\n", issue)
		}
		fmt.Println()
	}

	// 统计结果
	passedCount := 0
	for _, report := range reports {
		if report.Passed {
			passedCount++
		}
	}

	fmt.Printf("总计: %d/%d 验证通过\n", passedCount, len(reports))

	// 如果有验证失败，返回错误码
	if passedCount < len(reports) {
		os.Exit(1)
	}
}

// cleanAllData 清空所有数据
func cleanAllData(db *utils.Database) error {
	if err := seedUsersClean(db); err != nil {
		return err
	}
	if err := seedBooksClean(db); err != nil {
		return err
	}
	if err := seedSubscriptionsClean(db); err != nil {
		return err
	}
	return nil
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

// seedBooksClean 清空书籍数据
func seedBooksClean(db *utils.Database) error {
	seeder := NewBookstoreSeeder(db, cfg)
	return seeder.Clean()
}

// seedSubscriptions 填充订阅关系 (TODO: 实现)
func seedSubscriptions(db *utils.Database) error {
	fmt.Println("订阅关系填充功能待实现...")
	return nil
}

// seedSubscriptionsClean 清空订阅关系 (TODO: 实现)
func seedSubscriptionsClean(db *utils.Database) error {
	fmt.Println("订阅关系清空功能待实现...")
	return nil
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
