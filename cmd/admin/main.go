package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		printUsage(os.Stdout)
		return nil
	}

	switch args[0] {
	case "bootstrap":
		return runBootstrap(args[1:])
	case "-h", "--help", "help":
		printUsage(os.Stdout)
		return nil
	default:
		return fmt.Errorf("未知命令: %s", args[0])
	}
}

func printUsage(out *os.File) {
	_, _ = fmt.Fprintln(out, "Qingyu 后端管理员工具")
	_, _ = fmt.Fprintln(out, "")
	_, _ = fmt.Fprintln(out, "用法:")
	_, _ = fmt.Fprintln(out, "  admin bootstrap --username <用户名> --email <邮箱> [--config <配置文件>] [--force-reset-password] [--dry-run]")
	_, _ = fmt.Fprintln(out, "")
	_, _ = fmt.Fprintln(out, "说明:")
	_, _ = fmt.Fprintln(out, "  - 默认从环境变量 QINGYU_BOOTSTRAP_ADMIN_PASSWORD 读取初始密码")
	_, _ = fmt.Fprintln(out, "  - 推荐通过环境变量注入密码，避免命令历史泄漏")
}

func newBootstrapFlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("bootstrap", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	return fs
}
