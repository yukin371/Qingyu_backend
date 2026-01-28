//go:build ignore
// +build ignore

package main

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// MiddlewareConfig 中间件配置结构
type MiddlewareConfig struct {
	Middleware struct {
		RequestID map[string]interface{} `yaml:"request_id"`
		Recovery  map[string]interface{} `yaml:"recovery"`
		Security  map[string]interface{} `yaml:"security"`
		Logger    map[string]interface{} `yaml:"logger"`
		CORS      map[string]interface{} `yaml:"cors"`
		Compression map[string]interface{} `yaml:"compression"`
	} `yaml:"middleware"`
	PriorityOverrides map[string]int `yaml:"priority_overrides"`
}

// MiddlewareInfo 中间件信息
type MiddlewareInfo struct {
	Name           string
	DefaultPriority int
	OverridePriority *int
	Description     string
}

func main() {
	// 定义命令行参数
	configPath := flag.String("config", "configs/middleware.yaml", "中间件配置文件路径")
	outputFormat := flag.String("format", "text", "输出格式：text, json, markdown")
	flag.Parse()

	// 检查配置文件是否存在
	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		fmt.Printf("错误: 配置文件不存在: %s\n", *configPath)
		os.Exit(1)
	}

	// 读取配置文件
	config, err := loadConfig(*configPath)
	if err != nil {
		fmt.Printf("错误: 无法加载配置: %v\n", err)
		os.Exit(1)
	}

	// 定义所有中间件的默认信息
	middlewares := defineMiddlewares()

	// 应用优先级覆盖
	applyOverrides(middlewares, config.PriorityOverrides)

	// 按优先级排序
	sortedMiddlewares := sortMiddlewares(middlewares)

	// 检测冲突
	conflicts := detectConflicts(sortedMiddlewares)

	// 生成报告
	generateReport(sortedMiddlewares, conflicts, *outputFormat)
}

// loadConfig 加载配置文件
func loadConfig(path string) (*MiddlewareConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config MiddlewareConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// defineMiddlewares 定义所有中间件
func defineMiddlewares() map[string]*MiddlewareInfo {
	return map[string]*MiddlewareInfo{
		"request_id": {
			Name:            "request_id",
			DefaultPriority: 1,
			Description:     "请求ID中间件 - 为每个请求生成唯一标识",
		},
		"recovery": {
			Name:            "recovery",
			DefaultPriority: 2,
			Description:     "异常恢复中间件 - 捕获panic，恢复服务",
		},
		"error_handler": {
			Name:            "error_handler",
			DefaultPriority: 2,
			Description:     "错误处理中间件 - 捕获中间件链中的错误",
		},
		"security": {
			Name:            "security",
			DefaultPriority: 3,
			Description:     "安全头中间件 - 添加安全响应头",
		},
		"cors": {
			Name:            "cors",
			DefaultPriority: 4,
			Description:     "CORS中间件 - 处理跨域请求",
		},
		"timeout": {
			Name:            "timeout",
			DefaultPriority: 6,
			Description:     "超时中间件 - 设置请求超时",
		},
		"logger": {
			Name:            "logger",
			DefaultPriority: 7,
			Description:     "日志中间件 - 记录请求和响应日志",
		},
		"metrics": {
			Name:            "metrics",
			DefaultPriority: 7,
			Description:     "指标中间件 - 收集性能指标",
		},
		"rate_limit": {
			Name:            "rate_limit",
			DefaultPriority: 8,
			Description:     "限流中间件 - 限制请求频率",
		},
		"auth": {
			Name:            "auth",
			DefaultPriority: 9,
			Description:     "认证中间件 - 验证用户身份",
		},
		"permission": {
			Name:            "permission",
			DefaultPriority: 10,
			Description:     "权限中间件 - 检查用户权限",
		},
		"validation": {
			Name:            "validation",
			DefaultPriority: 11,
			Description:     "验证中间件 - 验证请求数据",
		},
		"compression": {
			Name:            "compression",
			DefaultPriority: 12,
			Description:     "压缩中间件 - gzip压缩响应",
		},
	}
}

// applyOverrides 应用优先级覆盖
func applyOverrides(middlewares map[string]*MiddlewareInfo, overrides map[string]int) {
	for name, priority := range overrides {
		if mw, exists := middlewares[name]; exists {
			mw.OverridePriority = &priority
		}
	}
}

// sortMiddlewares 按优先级排序中间件
func sortMiddlewares(middlewares map[string]*MiddlewareInfo) []*MiddlewareInfo {
	sorted := make([]*MiddlewareInfo, 0, len(middlewares))
	for _, mw := range middlewares {
		sorted = append(sorted, mw)
	}

	// 简单的冒泡排序
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			priorityI := getPriority(sorted[i])
			priorityJ := getPriority(sorted[j])
			if priorityJ < priorityI {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	return sorted
}

// getPriority 获取中间件的有效优先级
func getPriority(mw *MiddlewareInfo) int {
	if mw.OverridePriority != nil {
		return *mw.OverridePriority
	}
	return mw.DefaultPriority
}

// detectConflicts 检测优先级冲突
func detectConflicts(middlewares []*MiddlewareInfo) [][2]string {
	conflicts := [][2]string{}
	priorityMap := make(map[int][]string)

	// 按优先级分组
	for _, mw := range middlewares {
		priority := getPriority(mw)
		priorityMap[priority] = append(priorityMap[priority], mw.Name)
	}

	// 检测冲突
	for _, names := range priorityMap {
		if len(names) > 1 {
			// 有冲突
			for i := 0; i < len(names)-1; i++ {
				for j := i + 1; j < len(names); j++ {
					conflicts = append(conflicts, [2]string{names[i], names[j]})
				}
			}
		}
	}

	return conflicts
}

// generateReport 生成报告
func generateReport(middlewares []*MiddlewareInfo, conflicts [][2]string, format string) {
	switch format {
	case "json":
		generateJSONReport(middlewares, conflicts)
	case "markdown":
		generateMarkdownReport(middlewares, conflicts)
	default:
		generateTextReport(middlewares, conflicts)
	}
}

// generateTextReport 生成文本格式报告
func generateTextReport(middlewares []*MiddlewareInfo, conflicts [][2]string) {
	fmt.Println("===========================================")
	fmt.Println("      中间件执行顺序验证报告")
	fmt.Println("===========================================")
	fmt.Println()

	fmt.Printf("中间件总数: %d\n\n", len(middlewares))

	fmt.Println("========== 执行顺序 ==========")
	for i, mw := range middlewares {
		priority := getPriority(mw)
		overrideStatus := ""
		if mw.OverridePriority != nil {
			overrideStatus = fmt.Sprintf(" (覆盖自 %d)", mw.DefaultPriority)
		}

		fmt.Printf("%2d. %-20s Priority: %2d%s\n", i+1, mw.Name, priority, overrideStatus)
		fmt.Printf("    %s\n", mw.Description)
		fmt.Println()
	}

	if len(conflicts) > 0 {
		fmt.Println("========== 警告: 优先级冲突 ==========")
		for _, conflict := range conflicts {
			fmt.Printf("冲突: %s 和 %s 有相同的优先级\n", conflict[0], conflict[1])
		}
		fmt.Println()
	} else {
		fmt.Println("========== 优先级检查 ==========")
		fmt.Println("未发现优先级冲突")
		fmt.Println()
	}

	fmt.Println("===========================================")
	fmt.Println("              报告结束")
	fmt.Println("===========================================")
}

// generateJSONReport 生成JSON格式报告
func generateJSONReport(middlewares []*MiddlewareInfo, conflicts [][2]string) {
	fmt.Println("{")
	fmt.Println(`  "middlewares": [`)

	for i, mw := range middlewares {
		priority := getPriority(mw)
		comma := ","
		if i == len(middlewares)-1 {
			comma = ""
		}

		fmt.Printf("    {\n")
		fmt.Printf("      \"name\": \"%s\",\n", mw.Name)
		fmt.Printf("      \"priority\": %d,\n", priority)
		fmt.Printf("      \"default_priority\": %d,\n", mw.DefaultPriority)
		fmt.Printf("      \"description\": \"%s\"\n", mw.Description)
		fmt.Printf("    }%s\n", comma)
	}

	fmt.Println(`  ],`)
	fmt.Println(`  "conflicts": [`)

	for i, conflict := range conflicts {
		comma := ","
		if i == len(conflicts)-1 {
			comma = ""
		}
		fmt.Printf("    [\"%s\", \"%s\"]%s\n", conflict[0], conflict[1], comma)
	}

	fmt.Println("  ]")
	fmt.Println("}")
}

// generateMarkdownReport 生成Markdown格式报告
func generateMarkdownReport(middlewares []*MiddlewareInfo, conflicts [][2]string) {
	fmt.Println("# 中间件执行顺序验证报告")
	fmt.Println()
	fmt.Printf("**中间件总数**: %d\n\n", len(middlewares))

	fmt.Println("## 执行顺序")
	fmt.Println()
	fmt.Println("| 序号 | 中间件名称 | 优先级 | 默认优先级 | 描述 |")
	fmt.Println("|------|----------|--------|-----------|------|")

	for i, mw := range middlewares {
		priority := getPriority(mw)
		overrideMark := ""
		if mw.OverridePriority != nil {
			overrideMark = fmt.Sprintf(" (覆盖自 %d)", mw.DefaultPriority)
		}

		fmt.Printf("| %d | %s | %d%s | %d | %s |\n",
			i+1, mw.Name, priority, overrideMark, mw.DefaultPriority, mw.Description)
	}

	fmt.Println()

	if len(conflicts) > 0 {
		fmt.Println("## 警告: 优先级冲突")
		fmt.Println()
		for _, conflict := range conflicts {
			fmt.Printf("- **冲突**: %s 和 %s 有相同的优先级\n", conflict[0], conflict[1])
		}
		fmt.Println()
	} else {
		fmt.Println("## 优先级检查")
		fmt.Println()
		fmt.Println("未发现优先级冲突 ✓")
		fmt.Println()
	}
}
