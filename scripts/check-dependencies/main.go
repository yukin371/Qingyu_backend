package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Violation 表示依赖违规
type Violation struct {
	File     string
	Line     int
	Import   string
	Rule     string
	Severity string // "error", "warning", or "deprecated"
}

// Rule 表示依赖规则
type Rule struct {
	Pattern    string
	TargetPath string
	Reason     string
	Severity   string
}

var rules = []Rule{
	// {
	// 	Pattern:    `^api/v1/`,
	// 	TargetPath: `Qingyu_backend/service/shared`,
	// 	Reason:     "API层不应该直接依赖shared模块实现",
	// 	Severity:   "warning",
	// },
	// {
	// 	Pattern:    `^service/(user|writer|reader|ai)/`,
	// 	TargetPath: `Qingyu_backend/service/shared`,
	// 	Reason:     "业务服务不应该直接依赖shared模块实现，应该使用Port接口",
	// 	Severity:   "error",
	// },
	// 可以添加更多规则...
}

// 定义禁止的导入模式
var forbiddenImports = map[string]string{
	// 业务服务不应该直接导入shared的具体实现
	`service/user`:     `不应该直接导入shared模块，请使用service/interfaces/shared中的Port接口`,
	`service/writer`:   `不应该直接导入shared模块，请使用service/interfaces/shared中的Port接口`,
	`service/reader`:   `不应该直接导入shared模块，请使用service/interfaces/shared中的Port接口`,
	`service/ai`:       `不应该直接导入shared模块，请使用service/interfaces/shared中的Port接口`,
}

// 允许直接导入shared的模块
var allowedSharedImporters = map[string]bool{
	`service/container`:          true,
	`service/interfaces/shared`: true,
	`router/shared`:             true,
	`api/v1/auth`:               true,
	`api/v1/shared`:             true,
	`realtime/websocket`:        true,
	`middleware`:                true,
}

// 已废弃的导入路径（auth模块迁移）
var deprecatedImports = map[string]string{
	`Qingyu_backend/service/shared/auth`: `auth模块已迁移到service/auth，请使用新的导入路径`,
}

// 允许使用新auth路径的模块（用于向后兼容检查）
var newAuthAllowedImporters = map[string]bool{
	`service/container`:          true,
	`service/interfaces/shared`: true,
	`service/auth`:               true,
	`api/v1/auth`:                true,
	`api/v1/shared`:              true,
	`router/shared`:             true,
	`middleware`:                true,
}

func main() {
	fmt.Println("🔍 检查代码依赖关系...")
	fmt.Println()

	violations := []Violation{}

	// 遍历所有Go文件
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过vendor和隐藏目录
		if strings.Contains(path, "vendor") || strings.Contains(path, ".git") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// 只处理.go文件，跳过测试文件
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// 检查文件的依赖
		fileViolations := checkFile(path)
		violations = append(violations, fileViolations...)

		return nil
	})

	if err != nil {
		fmt.Printf("❌ 遍历文件失败: %v\n", err)
		os.Exit(1)
	}

	// 输出结果
	if len(violations) == 0 {
		fmt.Println("✅ 未发现依赖违规")
		os.Exit(0)
	}

	// 统计违规
	errorCount := 0
	warningCount := 0
	deprecatedCount := 0
	for _, v := range violations {
		if v.Severity == "error" {
			errorCount++
		} else if v.Severity == "warning" {
			warningCount++
		} else if v.Severity == "deprecated" {
			deprecatedCount++
		}
	}

	// 输出统计信息
	if errorCount > 0 || warningCount > 0 {
		if deprecatedCount > 0 {
			fmt.Printf("❌ 发现 %d 个错误, %d 个警告, %d 个废弃提示\n\n", errorCount, warningCount, deprecatedCount)
		} else {
			fmt.Printf("❌ 发现 %d 个错误, %d 个警告\n\n", errorCount, warningCount)
		}
	} else if deprecatedCount > 0 {
		fmt.Printf("ℹ️  发现 %d 个使用废弃路径的导入（仅供参考）\n\n", deprecatedCount)
	}

	// 输出详细违规信息
	for i, v := range violations {
		severityIcon := "ℹ️ "
		if v.Severity == "error" {
			severityIcon = "❌"
		} else if v.Severity == "warning" {
			severityIcon = "⚠️ "
		}
		fmt.Printf("%s [%d] %s:%d\n", severityIcon, i+1, v.File, v.Line)
		fmt.Printf("   导入: %s\n", v.Import)
		fmt.Printf("   规则: %s\n", v.Rule)
		fmt.Println()
	}

	// 输出修复建议
	fmt.Println("💡 修复建议:")
	if deprecatedCount > 0 {
		fmt.Println("   废弃路径迁移:")
		fmt.Println("   - 将 Qingyu_backend/service/shared/auth 改为 Qingyu_backend/service/auth")
		fmt.Println("   - 兼容层会继续工作，但建议尽快迁移")
		fmt.Println("   - 迁移指南: docs/migration/auth-module-migration.md")
	}
	if errorCount > 0 || warningCount > 0 {
		fmt.Println("   依赖规范:")
		fmt.Println("   1. 使用service/interfaces/shared中定义的Port接口")
		fmt.Println("   2. 通过依赖注入而非直接导入")
		fmt.Println("   3. 参考文档: docs/architecture/dependency-rules.md")
	}

	// 根据违规类型决定退出码
	if errorCount > 0 {
		os.Exit(1)
	} else if warningCount > 0 {
		os.Exit(1)
	}
	// deprecated不影响CI通过
	os.Exit(0)
}

// checkFile 检查单个文件的依赖
func checkFile(filePath string) []Violation {
	violations := []Violation{}

	// 获取相对路径
	relPath, err := filepath.Rel(".", filePath)
	if err != nil {
		return violations
	}

	// 转换为Unix风格的路径
	relPath = filepath.ToSlash(relPath)

	// 解析Go文件
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ImportsOnly)
	if err != nil {
		return violations
	}

	// 获取文件所在目录
	fileDir := filepath.Dir(relPath)

	// 检查每个import
	ast.Inspect(node, func(n ast.Node) bool {
		importSpec, ok := n.(*ast.ImportSpec)
		if !ok {
			return true
		}

		if importSpec.Path == nil {
			return true
		}

		importPath := strings.Trim(importSpec.Path.Value, `"`)
		violation := checkImport(relPath, fileDir, importPath, fset.Position(importSpec.Pos()).Line)
		if violation != nil {
			violations = append(violations, *violation)
		}

		return true
	})

	return violations
}

// checkImport 检查单个import是否违规
func checkImport(filePath, fileDir, importPath string, line int) *Violation {
	// 只检查项目内部的导入
	if !strings.HasPrefix(importPath, "Qingyu_backend") {
		return nil
	}

	// 获取导入的模块路径
	importModule := strings.TrimPrefix(importPath, "Qingyu_backend/")

	// 规则0: 检查是否使用了已废弃的auth路径
	if deprecationMsg, isDeprecated := deprecatedImports[importPath]; isDeprecated {
		// 检查是否是测试文件
		isTestFile := strings.HasSuffix(filePath, "_test.go")
		if isTestFile {
			// 测试文件给出info级别提示
			return &Violation{
				File:     filePath,
				Line:     line,
				Import:   importPath,
				Rule:     fmt.Sprintf("%s (测试文件可以继续使用，但建议迁移)", deprecationMsg),
				Severity: "deprecated",
			}
		}

		// 生产代码给出warning
		return &Violation{
			File:     filePath,
			Line:     line,
			Import:   importPath,
			Rule:     deprecationMsg,
			Severity: "warning",
		}
	}

	// 规则1: 检查业务服务是否直接导入shared
	if strings.HasPrefix(importModule, "service/shared/") {
		// 检查是否在允许列表中
		if allowedSharedImporters[fileDir] {
			return nil
		}

		// 检查文件是否在业务服务目录下
		for businessDir := range forbiddenImports {
			if strings.HasPrefix(fileDir, businessDir) {
				return &Violation{
					File:     filePath,
					Line:     line,
					Import:   importPath,
					Rule:     forbiddenImports[businessDir],
					Severity: "error",
				}
			}
		}
	}

	return nil
}

// checkCircularDependencies 检查循环依赖
func checkCircularDependencies() error {
	// TODO: 实现循环依赖检测
	return nil
}

// parseImportDecl 解析import声明
func parseImportDecl(line string) string {
	re := regexp.MustCompile(`import\s+(?P<import>"[^"]+"|` + "`[^`]+`)")
	matches := re.FindStringSubmatch(line)
	if len(matches) > 1 {
		return strings.Trim(matches[1], "\"`")
	}
	return ""
}

// readImportsFromFile 从文件中读取所有import
func readImportsFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	imports := []string{}
	scanner := bufio.NewScanner(file)
	inImport := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 检测import块开始
		if strings.HasPrefix(line, "import (") {
			inImport = true
			continue
		}

		// 检测import块结束
		if line == ")" {
			inImport = false
			continue
		}

		// 解析import行
		if strings.HasPrefix(line, "import ") || inImport {
			importPath := parseImportDecl(line)
			if importPath != "" {
				imports = append(imports, importPath)
			}
		}
	}

	return imports, scanner.Err()
}
