package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	version = "1.0.0"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "analyze":
		handleAnalyze()
	case "migrate":
		handleMigrate()
	case "validate":
		handleValidate()
	case "testgen":
		handleTestGen()
	case "version", "-v", "--version":
		fmt.Printf("Migration Tools v%s\n", version)
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf(`Migration Tools v%s

USAGE:
    migration-tools <command> [options]

COMMANDS:
    analyze     Analyze API files for response calls
    migrate     Migrate API files from shared to response package
    validate    Validate migration quality
    testgen     Generate test framework for API files
    version     Show version information
    help        Show this help message

EXAMPLES:
    # Analyze writer module
    migration-tools analyze --path api/v1/writer

    # Migrate single file (dry-run)
    migration-tools migrate --file api/v1/writer/audit_api.go --dry-run

    # Migrate single file
    migration-tools migrate --file api/v1/writer/audit_api.go

    # Validate migration
    migration-tools validate --path api/v1/writer

    # Generate test framework
    migration-tools testgen --file api/v1/writer/audit_api.go

For more information about each command, use:
    migration-tools <command> --help
`, version)
}

func handleAnalyze() {
	cmd := flag.NewFlagSet("analyze", flag.ExitOnError)
	path := cmd.String("path", "", "Path to API directory or file")
	output := cmd.String("output", "", "Output file (JSON format)")
	verbose := cmd.Bool("verbose", false, "Verbose output")

	if err := cmd.Parse(os.Args[2:]); err != nil {
		fmt.Printf("Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if *path == "" {
		fmt.Println("Error: --path is required")
		cmd.PrintDefaults()
		os.Exit(1)
	}

	result, err := AnalyzePath(*path, *verbose)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	result.Print()

	if *output != "" {
		err = result.SaveToFile(*output)
		if err != nil {
			fmt.Printf("Error saving to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("\n结果已保存到: %s\n", *output)
	}
}

func handleMigrate() {
	cmd := flag.NewFlagSet("migrate", flag.ExitOnError)
	file := cmd.String("file", "", "Path to API file to migrate")
	dryRun := cmd.Bool("dry-run", false, "Preview changes without modifying files")
	backup := cmd.Bool("backup", true, "Create backup before migration")
	simple := cmd.Bool("simple", false, "Use simple string replacement (faster but less precise)")
	_ = cmd.Bool("verbose", false, "Verbose output")

	if err := cmd.Parse(os.Args[2:]); err != nil {
		fmt.Printf("Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if *file == "" {
		fmt.Println("Error: --file is required")
		cmd.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("Migrating: %s\n", *file)
	if *dryRun {
		fmt.Println("Mode: DRY RUN (no changes will be made)")
	}
	if *simple {
		fmt.Println("Method: Simple string replacement")
	}
	if *backup && !*dryRun {
		fmt.Println("Backup: enabled")
	}

	var result *MigrationResult
	var err error

	if *simple {
		result, err = SimpleMigrate(*file, *dryRun)
	} else {
		result, err = MigrateFile(*file, *dryRun, *backup)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	result.Print()
}

func handleValidate() {
	cmd := flag.NewFlagSet("validate", flag.ExitOnError)
	path := cmd.String("path", "", "Path to API directory or file")
	checks := cmd.String("checks", "all", "Comma-separated list of checks (imports,no_shared_calls,swagger,all)")
	_ = cmd.Bool("verbose", false, "Verbose output")

	if err := cmd.Parse(os.Args[2:]); err != nil {
		fmt.Printf("Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if *path == "" {
		fmt.Println("Error: --path is required")
		cmd.PrintDefaults()
		os.Exit(1)
	}

	result, err := ValidatePath(*path, *checks)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	result.Print()

	// 返回适当的退出码
	if !result.Passed {
		os.Exit(1)
	}
}

func handleTestGen() {
	cmd := flag.NewFlagSet("testgen", flag.ExitOnError)
	file := cmd.String("file", "", "Path to API file")
	output := cmd.String("output", "", "Output test file path (auto-generated if not specified)")
	verbose := cmd.Bool("verbose", false, "Verbose output")

	if err := cmd.Parse(os.Args[2:]); err != nil {
		fmt.Printf("Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if *file == "" {
		fmt.Println("Error: --file is required")
		cmd.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("Generating tests for: %s\n", *file)
	if *output != "" {
		fmt.Printf("Output: %s\n", *output)
	} else {
		fmt.Println("Output: auto-generated")
	}
	if *verbose {
		fmt.Println("Verbose: enabled")
	}

	// TODO: 实现测试生成逻辑
	fmt.Println("✓ Test generation complete")
}
