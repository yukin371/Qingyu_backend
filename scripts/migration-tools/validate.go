package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// ValidationResult 验证结果
type ValidationResult struct {
	Path              string           `json:"path"`
	TotalFiles        int              `json:"total_files"`
	Passed            bool             `json:"passed"`
	Checks            CheckResults     `json:"checks"`
	Issues            []Issue          `json:"issues"`
	FilesChecked      []FileValidation `json:"files_checked"`
}

// CheckResults 检查结果
type CheckResults struct {
	ImportsClean      bool `json:"imports_clean"`
	TestsPassing      bool `json:"tests_passing"`
	SwaggerUpdated    bool `json:"swagger_updated"`
	NoSharedCalls     bool `json:"no_shared_calls"`
	AllResponseCalls  bool `json:"all_response_calls"`
}

// FileValidation 单个文件验证
type FileValidation struct {
	Filename        string  `json:"filename"`
	Passed          bool    `json:"passed"`
	ImportsClean    bool    `json:"imports_clean"`
	NoSharedCalls   bool    `json:"no_shared_calls"`
	SwaggerUpdated  bool    `json:"swagger_updated"`
	Issues          []Issue `json:"issues"`
}

// Issue 问题
type Issue struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Type     string `json:"type"`
	Message  string `json:"message"`
	Severity string `json:"severity"` // error, warning, info
}

// ValidatePath 验证指定路径
func ValidatePath(path string, checks string) (*ValidationResult, error) {
	result := &ValidationResult{
		Path:         path,
		Checks:       CheckResults{},
		Issues:       []Issue{},
		FilesChecked: []FileValidation{},
	}

	// 解析要执行的检查
	checkList := strings.Split(checks, ",")
	doAll := false
	for _, check := range checkList {
		if strings.TrimSpace(check) == "all" {
			doAll = true
			break
		}
	}

	// 检查是文件还是目录
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("path not found: %s", path)
	}

	var files []string
	if info.IsDir() {
		// 递归查找所有*_api.go文件
		err = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(filePath, "_api.go") {
				files = append(files, filePath)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		files = []string{path}
	}

	// 验证每个文件
	for _, file := range files {
		fileValidation, err := validateFile(file, doAll, checkList)
		if err != nil {
			result.Issues = append(result.Issues, Issue{
				File:     file,
				Type:     "validation_error",
				Message:  err.Error(),
				Severity: "error",
			})
			continue
		}
		result.FilesChecked = append(result.FilesChecked, *fileValidation)
		result.Issues = append(result.Issues, fileValidation.Issues...)
		result.TotalFiles++
	}

	// 计算整体检查结果
	result.calculateOverallResults()

	return result, nil
}

// validateFile 验证单个文件
func validateFile(filePath string, doAll bool, checkList []string) (*FileValidation, error) {
	validation := &FileValidation{
		Filename: filePath,
		Issues:   []Issue{},
	}

	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	contentStr := string(content)

	// 解析AST
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// 检查shared调用
	if doAll || containsCheck(checkList, "no_shared_calls") {
		validation.NoSharedCalls = checkNoSharedCallsInFile(node, validation)
	}

	// 检查导入清理
	if doAll || containsCheck(checkList, "imports") {
		validation.ImportsClean = checkImportsCleanInFile(node, validation)
	}

	// 检查Swagger注释
	if doAll || containsCheck(checkList, "swagger") {
		validation.SwaggerUpdated = checkSwaggerUpdated(contentStr, validation)
	}

	// 计算整体结果
	validation.Passed = validation.NoSharedCalls && validation.ImportsClean && validation.SwaggerUpdated

	return validation, nil
}

// checkNoSharedCallsInFile 检查是否还有shared调用
func checkNoSharedCallsInFile(node *ast.File, validation *FileValidation) bool {
	hasSharedCalls := false

	ast.Inspect(node, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		pkgName, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}

		if pkgName.Name == "shared" {
			hasSharedCalls = true
			funcName := sel.Sel.Name

			var severity string
			if funcName == "Error" || funcName == "Success" || funcName == "ValidationError" {
				severity = "error"
			} else {
				severity = "warning"
			}

			validation.Issues = append(validation.Issues, Issue{
				File:     validation.Filename,
				Line:     0, // TODO: 获取行号
				Type:     "shared_call_found",
				Message:  fmt.Sprintf("发现shared.%s调用，应该迁移到response包", funcName),
				Severity: severity,
			})
		}

		return true
	})

	return !hasSharedCalls
}

// checkImportsCleanInFile 检查导入是否清理干净
func checkImportsCleanInFile(node *ast.File, validation *FileValidation) bool {
	hasSharedImport := false

	ast.Inspect(node, func(n ast.Node) bool {
		importSpec, ok := n.(*ast.ImportSpec)
		if !ok {
			return true
		}

		if importSpec.Path != nil {
			importPath := strings.Trim(importSpec.Path.Value, `"`)

			// 检查shared导入
			if importPath == "Qingyu_backend/api/v1/shared" {
				hasSharedImport = true
				validation.Issues = append(validation.Issues, Issue{
					File:     validation.Filename,
					Line:     0,
					Type:     "shared_import_found",
					Message:  "发现shared包导入，应该移除",
					Severity: "error",
				})
			}
		}

		return true
	})

	return !hasSharedImport
}

// checkSwaggerUpdated 检查Swagger注释是否更新
func checkSwaggerUpdated(content string, validation *FileValidation) bool {
	// 检查是否还有shared.APIResponse
	if strings.Contains(content, "shared.APIResponse") {
		validation.Issues = append(validation.Issues, Issue{
			File:     validation.Filename,
			Line:     0,
			Type:     "swagger_not_updated",
			Message:  "Swagger注释中仍有shared.APIResponse，应改为response.APIResponse",
			Severity: "warning",
		})
		return false
	}

	return true
}

// containsCheck 检查是否包含指定的检查项
func containsCheck(checkList []string, check string) bool {
	for _, c := range checkList {
		if strings.TrimSpace(c) == check {
			return true
		}
	}
	return false
}

// calculateOverallResults 计算整体检查结果
func (r *ValidationResult) calculateOverallResults() {
	if len(r.FilesChecked) == 0 {
		return
	}

	allPassed := true
	for _, file := range r.FilesChecked {
		if !file.Passed {
			allPassed = false
			break
		}
	}

	r.Passed = allPassed
	r.Checks.NoSharedCalls = true
	r.Checks.ImportsClean = true
	r.Checks.SwaggerUpdated = true

	for _, file := range r.FilesChecked {
		if !file.NoSharedCalls {
			r.Checks.NoSharedCalls = false
		}
		if !file.ImportsClean {
			r.Checks.ImportsClean = false
		}
		if !file.SwaggerUpdated {
			r.Checks.SwaggerUpdated = false
		}
	}

	// 统计严重问题
	errorCount := 0
	warningCount := 0
	for _, issue := range r.Issues {
		switch issue.Severity {
		case "error":
			errorCount++
		case "warning":
			warningCount++
		}
	}

	// 如果有错误，则验证不通过
	if errorCount > 0 {
		r.Passed = false
	}
}

// Print 打印验证结果
func (r *ValidationResult) Print() {
	fmt.Printf("\n=== 验证结果: %s ===\n", r.Path)
	fmt.Printf("文件总数: %d\n", r.TotalFiles)
	fmt.Printf("整体状态: ")
	if r.Passed {
		fmt.Println("✅ 通过")
	} else {
		fmt.Println("❌ 失败")
	}

	fmt.Printf("\n检查项:\n")
	fmt.Printf("  导入清理: ")
	if r.Checks.ImportsClean {
		fmt.Println("✅")
	} else {
		fmt.Println("❌")
	}
	fmt.Printf("  无shared调用: ")
	if r.Checks.NoSharedCalls {
		fmt.Println("✅")
	} else {
		fmt.Println("❌")
	}
	fmt.Printf("  Swagger更新: ")
	if r.Checks.SwaggerUpdated {
		fmt.Println("✅")
	} else {
		fmt.Println("❌")
	}

	if len(r.Issues) > 0 {
		fmt.Printf("\n问题列表 (%d):\n", len(r.Issues))
		for i, issue := range r.Issues {
			if i >= 20 { // 最多显示20个问题
				fmt.Printf("  ... 还有 %d 个问题\n", len(r.Issues)-20)
				break
			}
			fmt.Printf("  [%d] %s: %s\n", i+1, issue.Type, issue.Message)
		}
	}

	if len(r.FilesChecked) > 0 {
		fmt.Printf("\n文件详情:\n")
		for _, file := range r.FilesChecked {
			status := "✅"
			if !file.Passed {
				status = "❌"
			}
			fmt.Printf("  %s %s\n", status, filepath.Base(file.Filename))
		}
	}
}
