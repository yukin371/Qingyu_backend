package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"time"
)

// MigrationResult 迁移结果
type MigrationResult struct {
	Filename           string   `json:"filename"`
	SuccessCount       int      `json:"success_count"`
	ErrorCount         int      `json:"error_count"`
	ValidationCount    int      `json:"validation_count"`
	CallsMigrated      int      `json:"calls_migrated"`
	ImportsRemoved     []string `json:"imports_removed"`
	ImportsAdded       []string `json:"imports_added"`
	BackupFile         string   `json:"backup_file"`
	DryRun             bool     `json:"dry_run"`
	MigrationTime      float64  `json:"migration_time"`
}

// MigrateFile 迁移单个文件
func MigrateFile(filePath string, dryRun bool, backup bool) (*MigrationResult, error) {
	startTime := time.Now()

	result := &MigrationResult{
		Filename:       filePath,
		DryRun:         dryRun,
		ImportsRemoved: []string{},
		ImportsAdded:   []string{},
	}

	// 读取文件
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	// 创建备份
	if backup && !dryRun {
		backupFile := filePath + ".bak"
		err = os.WriteFile(backupFile, content, 0644)
		if err != nil {
			return nil, fmt.Errorf("创建备份失败: %w", err)
		}
		result.BackupFile = backupFile
	}

	// 解析AST
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("解析文件失败: %w", err)
	}

	// 创建AST重写器
	rewriter := &migrationRewriter{
		fset:    fset,
		result:  result,
		dryRun:  dryRun,
		origContent: string(content),
	}

	// 重写AST
	ast.Inspect(node, rewriter.Rewrite)

	// 如果不是dry run，写入文件
	if !dryRun {
		// 格式化并写入
		var buf bytes.Buffer
		err = format.Node(&buf, fset, node)
		if err != nil {
			return nil, fmt.Errorf("格式化代码失败: %w", err)
		}

		err = os.WriteFile(filePath, buf.Bytes(), 0644)
		if err != nil {
			return nil, fmt.Errorf("写入文件失败: %w", err)
		}
	}

	result.MigrationTime = time.Since(startTime).Seconds()

	return result, nil
}

// migrationRewriter AST重写器
type migrationRewriter struct {
	fset        *token.FileSet
	result      *MigrationResult
	dryRun      bool
	origContent string
	importsAdded map[string]bool
	importsRemoved map[string]bool
}

// Rewrite 重写AST节点
func (r *migrationRewriter) Rewrite(n ast.Node) bool {
	call, ok := n.(*ast.CallExpr)
	if !ok {
		return true
	}

	// 检查是否是 selector expression (如 shared.Error)
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return true
	}

	// 获取包名
	pkgName, ok := sel.X.(*ast.Ident)
	if !ok {
		return true
	}

	// 获取函数名
	funcName := sel.Sel.Name

	// 处理shared包的调用
	if pkgName.Name == "shared" {
		r.migrateSharedCall(call, sel, funcName)
	}

	return true
}

// migrateSharedCall 迁移shared调用
func (r *migrationRewriter) migrateSharedCall(call *ast.CallExpr, sel *ast.SelectorExpr, funcName string) {
	switch funcName {
	case "Success":
		r.migrateSuccessCall(call, sel)
	case "Error":
		r.migrateErrorCall(call, sel)
	case "ValidationError":
		r.migrateValidationErrorCall(call, sel)
	}
}

// migrateSuccessCall 迁移Success调用
func (r *migrationRewriter) migrateSuccessCall(call *ast.CallExpr, sel *ast.SelectorExpr) {
	// shared.Success(c, http.StatusCreated, "创建成功", data)
	// → response.Created(c, data)

	if len(call.Args) < 4 {
		return
	}

	// 检查第二个参数是否是HTTP状态码
	statusCode, ok := r.extractStatusCode(call.Args[1])
	if !ok {
		return
	}

	// 确定使用哪个response函数
	var newFuncName string
	if statusCode == 201 {
		newFuncName = "Created"
	} else {
		newFuncName = "Success"
	}

	// 更新调用
	if ident, ok := sel.X.(*ast.Ident); ok {
		ident.Name = "response"
	}
	sel.Sel.Name = newFuncName

	// 移除HTTP状态码参数和消息参数
	if newFuncName == "Success" {
		// Success(c, data)
		call.Args = []ast.Expr{call.Args[0], call.Args[3]}
	} else {
		// Created(c, data)
		call.Args = []ast.Expr{call.Args[0], call.Args[3]}
	}

	r.result.SuccessCount++
	r.result.CallsMigrated++
	r.recordImportChange("response", true)
}

// migrateErrorCall 迁移Error调用
func (r *migrationRewriter) migrateErrorCall(call *ast.CallExpr, sel *ast.SelectorExpr) {
	// shared.Error(c, http.StatusBadRequest, "参数错误", "details")
	// → response.BadRequest(c, "参数错误", "details")

	if len(call.Args) < 4 {
		return
	}

	// 检查第二个参数是否是HTTP状态码
	statusCode, ok := r.extractStatusCode(call.Args[1])
	if !ok {
		return
	}

	// 确定使用哪个response函数
	var newFuncName string
	switch statusCode {
	case 400:
		newFuncName = "BadRequest"
	case 401:
		newFuncName = "Unauthorized"
	case 403:
		newFuncName = "Forbidden"
	case 404:
		newFuncName = "NotFound"
	case 409:
		newFuncName = "Conflict"
	default:
		if statusCode >= 500 {
			newFuncName = "InternalError"
		} else {
			// 其他状态码不迁移
			return
		}
	}

	// 更新调用
	if ident, ok := sel.X.(*ast.Ident); ok {
		ident.Name = "response"
	}
	sel.Sel.Name = newFuncName

	// 移除HTTP状态码参数
	if newFuncName == "InternalError" {
		// InternalError(c, err)
		call.Args = []ast.Expr{call.Args[0], call.Args[3]}
	} else {
		// BadRequest(c, "参数错误", "details")
		call.Args = []ast.Expr{call.Args[0], call.Args[2], call.Args[3]}
	}

	r.result.ErrorCount++
	r.result.CallsMigrated++
	r.recordImportChange("response", true)
}

// migrateValidationErrorCall 迁移ValidationError调用
func (r *migrationRewriter) migrateValidationErrorCall(call *ast.CallExpr, sel *ast.SelectorExpr) {
	// shared.ValidationError(c, err)
	// → response.BadRequest(c, "参数错误", err.Error())

	if len(call.Args) < 2 {
		return
	}

	// 更新调用
	if ident, ok := sel.X.(*ast.Ident); ok {
		ident.Name = "response"
	}
	sel.Sel.Name = "BadRequest"

	// 添加消息参数
	// BadRequest(c, "参数错误", err.Error())
	// TODO: 需要添加 "参数错误" 字符串和调用err.Error()

	r.result.ValidationCount++
	r.result.CallsMigrated++
	r.recordImportChange("response", true)
}

// extractStatusCode 提取HTTP状态码
func (r *migrationRewriter) extractStatusCode(expr ast.Expr) (int, bool) {
	// 检查是否是 selector expression (如 http.StatusOK)
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return 0, false
	}

	// 检查包名是否是http
	pkg, ok := sel.X.(*ast.Ident)
	if !ok || pkg.Name != "http" {
		return 0, false
	}

	// 获取状态码
	switch sel.Sel.Name {
	case "StatusOK":
		return 200, true
	case "StatusCreated":
		return 201, true
	case "StatusNoContent":
		return 204, true
	case "StatusBadRequest":
		return 400, true
	case "StatusUnauthorized":
		return 401, true
	case "StatusForbidden":
		return 403, true
	case "StatusNotFound":
		return 404, true
	case "StatusConflict":
		return 409, true
	case "StatusInternalServerError":
		return 500, true
	default:
		return 0, false
	}
}

// recordImportChange 记录导入变更
func (r *migrationRewriter) recordImportChange(importPath string, added bool) {
	if r.importsAdded == nil {
		r.importsAdded = make(map[string]bool)
		r.importsRemoved = make(map[string]bool)
	}

	if added {
		if !r.importsAdded[importPath] {
			r.importsAdded[importPath] = true
			r.result.ImportsAdded = append(r.result.ImportsAdded, importPath)
		}
	} else {
		if !r.importsRemoved[importPath] {
			r.importsRemoved[importPath] = true
			r.result.ImportsRemoved = append(r.result.ImportsRemoved, importPath)
		}
	}
}

// Print 打印迁移结果
func (r *MigrationResult) Print() {
	fmt.Printf("\n=== 迁移结果: %s ===\n", r.Filename)
	if r.DryRun {
		fmt.Printf("模式: DRY RUN (未实际修改文件)\n")
	}
	fmt.Printf("Success调用迁移: %d\n", r.SuccessCount)
	fmt.Printf("Error调用迁移: %d\n", r.ErrorCount)
	fmt.Printf("ValidationError迁移: %d\n", r.ValidationCount)
	fmt.Printf("总调用迁移: %d\n", r.CallsMigrated)
	if len(r.ImportsAdded) > 0 {
		fmt.Printf("添加导入: %v\n", r.ImportsAdded)
	}
	if len(r.ImportsRemoved) > 0 {
		fmt.Printf("移除导入: %v\n", r.ImportsRemoved)
	}
	if r.BackupFile != "" {
		fmt.Printf("备份文件: %s\n", r.BackupFile)
	}
	fmt.Printf("耗时: %.2f秒\n", r.MigrationTime)
}

// GenerateMigrationReport 生成迁移报告
func GenerateMigrationReport(results []*MigrationResult) string {
	var buf bytes.Buffer

	buf.WriteString("\n=== 迁移报告 ===\n")
	buf.WriteString(fmt.Sprintf("文件总数: %d\n", len(results)))

	totalCalls := 0
	totalTime := 0.0
	for _, result := range results {
		totalCalls += result.CallsMigrated
		totalTime += result.MigrationTime
	}

	buf.WriteString(fmt.Sprintf("总调用迁移: %d\n", totalCalls))
	buf.WriteString(fmt.Sprintf("总耗时: %.2f秒\n", totalTime))

	buf.WriteString("\n详细结果:\n")
	for _, result := range results {
		buf.WriteString(fmt.Sprintf("\n%s:\n", result.Filename))
		buf.WriteString(fmt.Sprintf("  迁移调用: %d\n", result.CallsMigrated))
		buf.WriteString(fmt.Sprintf("  耗时: %.2f秒\n", result.MigrationTime))
		if result.BackupFile != "" {
			buf.WriteString(fmt.Sprintf("  备份: %s\n", result.BackupFile))
		}
	}

	return buf.String()
}

// SimpleMigrate 简单迁移（使用字符串替换，用于快速原型）
func SimpleMigrate(filePath string, dryRun bool) (*MigrationResult, error) {
	startTime := time.Now()

	result := &MigrationResult{
		Filename:   filePath,
		DryRun:     dryRun,
	}

	// 读取文件
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	originalContent := string(content)
	newContent := originalContent

	// 定义替换规则
	replacements := []struct{
		from string
		to   string
		count *int
	}{
		// shared.Success → response.Success
		{
			from: "shared.Success(c, http.StatusOK,",
			to:   "response.Success(c,",
			count: &result.SuccessCount,
		},
		{
			from: "shared.Success(c, http.StatusCreated,",
			to:   "response.Created(c,",
			count: &result.SuccessCount,
		},
		// shared.Error → response.*
		{
			from: "shared.Error(c, http.StatusBadRequest,",
			to:   "response.BadRequest(c,",
			count: &result.ErrorCount,
		},
		{
			from: "shared.Error(c, http.StatusUnauthorized,",
			to:   "response.Unauthorized(c,",
			count: &result.ErrorCount,
		},
		{
			from: "shared.Error(c, http.StatusForbidden,",
			to:   "response.Forbidden(c,",
			count: &result.ErrorCount,
		},
		{
			from: "shared.Error(c, http.StatusNotFound,",
			to:   "response.NotFound(c,",
			count: &result.ErrorCount,
		},
		{
			from: "shared.Error(c, http.StatusConflict,",
			to:   "response.Conflict(c,",
			count: &result.ErrorCount,
		},
		{
			from: "shared.Error(c, http.StatusInternalServerError,",
			to:   "response.InternalError(c,",
			count: &result.ErrorCount,
		},
		// shared.ValidationError → response.BadRequest
		{
			from: "shared.ValidationError(c,",
			to:   "response.BadRequest(c, \"参数错误\",",
			count: &result.ValidationCount,
		},
	}

	// 执行替换
	for _, repl := range replacements {
		if strings.Contains(newContent, repl.from) {
			count := strings.Count(newContent, repl.from)
			newContent = strings.ReplaceAll(newContent, repl.from, repl.to)
			*repl.count += count
		}
	}

	// 计算总迁移数
	result.CallsMigrated = result.SuccessCount + result.ErrorCount + result.ValidationCount

	// 如果内容有变化且不是dry run，写入文件
	if newContent != originalContent && !dryRun {
		// 创建备份
		backupFile := filePath + ".bak"
		err = os.WriteFile(backupFile, content, 0644)
		if err != nil {
			return nil, fmt.Errorf("创建备份失败: %w", err)
		}
		result.BackupFile = backupFile

		// 写入新内容
		err = os.WriteFile(filePath, []byte(newContent), 0644)
		if err != nil {
			return nil, fmt.Errorf("写入文件失败: %w", err)
		}
	}

	result.MigrationTime = time.Since(startTime).Seconds()

	return result, nil
}
