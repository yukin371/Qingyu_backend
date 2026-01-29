package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// AnalysisResult 分析结果
type AnalysisResult struct {
	Module         string            `json:"module"`
	TotalFiles     int               `json:"total_files"`
	TotalCalls     int               `json:"total_calls"`
	Breakdown      CallBreakdown     `json:"breakdown"`
	Files          []FileAnalysis    `json:"files"`
	Complexity     string            `json:"complexity"`
	Recommendation string            `json:"recommendation"`
}

// CallBreakdown 调用分类统计
type CallBreakdown struct {
	SharedSuccess        int `json:"shared_success"`
	SharedError          int `json:"shared_error"`
	SharedValidationError int `json:"shared_validation_error"`
	ResponseSuccess      int `json:"response_success"`
	ResponseBadRequest   int `json:"response_bad_request"`
	ResponseUnauthorized int `json:"response_unauthorized"`
	ResponseForbidden    int `json:"response_forbidden"`
	ResponseNotFound     int `json:"response_not_found"`
	ResponseConflict     int `json:"response_conflict"`
	ResponseInternal     int `json:"response_internal_error"`
	ResponseCreated      int `json:"response_created"`
	ResponseNoContent    int `json:"response_no_content"`
	ResponsePaginated    int `json:"response_paginated"`
}

// FileAnalysis 单个文件分析
type FileAnalysis struct {
	Filename             string            `json:"filename"`
	FilePath             string            `json:"file_path"`
	Functions            int               `json:"functions"`
	TotalCalls           int               `json:"total_calls"`
	SharedCalls          int               `json:"shared_calls"`
	ResponseCalls        int               `json:"response_calls"`
	CallBreakdown        CallBreakdown     `json:"call_breakdown"`
	HasWebSocket         bool              `json:"has_websocket"`
	HasFileDownload      bool              `json:"has_file_download"`
	Complexity           string            `json:"complexity"`
	RiskLevel            string            `json:"risk_level"`
	EstimatedMinutes     int               `json:"estimated_minutes"`
}

// AnalyzePath 分析指定路径
func AnalyzePath(path string, verbose bool) (*AnalysisResult, error) {
	result := &AnalysisResult{
		Module: filepath.Base(path),
		Files:  []FileAnalysis{},
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

	// 分析每个文件
	for _, file := range files {
		fileAnalysis, err := analyzeFile(file, verbose)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to analyze %s: %v\n", file, err)
			continue
		}
		result.Files = append(result.Files, *fileAnalysis)
		result.TotalFiles++
		result.TotalCalls += fileAnalysis.TotalCalls
		result.Breakdown.SharedSuccess += fileAnalysis.CallBreakdown.SharedSuccess
		result.Breakdown.SharedError += fileAnalysis.CallBreakdown.SharedError
		result.Breakdown.SharedValidationError += fileAnalysis.CallBreakdown.SharedValidationError
		result.Breakdown.ResponseSuccess += fileAnalysis.CallBreakdown.ResponseSuccess
		result.Breakdown.ResponseBadRequest += fileAnalysis.CallBreakdown.ResponseBadRequest
		result.Breakdown.ResponseUnauthorized += fileAnalysis.CallBreakdown.ResponseUnauthorized
		result.Breakdown.ResponseForbidden += fileAnalysis.CallBreakdown.ResponseForbidden
		result.Breakdown.ResponseNotFound += fileAnalysis.CallBreakdown.ResponseNotFound
		result.Breakdown.ResponseConflict += fileAnalysis.CallBreakdown.ResponseConflict
		result.Breakdown.ResponseInternal += fileAnalysis.CallBreakdown.ResponseInternal
		result.Breakdown.ResponseCreated += fileAnalysis.CallBreakdown.ResponseCreated
		result.Breakdown.ResponseNoContent += fileAnalysis.CallBreakdown.ResponseNoContent
		result.Breakdown.ResponsePaginated += fileAnalysis.CallBreakdown.ResponsePaginated
	}

	// 计算整体复杂度和建议
	result.Complexity = calculateOverallComplexity(result)
	result.Recommendation = generateRecommendation(result)

	return result, nil
}

// analyzeFile 分析单个文件
func analyzeFile(filePath string, verbose bool) (*FileAnalysis, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	analysis := &FileAnalysis{
		Filename: filepath.Base(filePath),
		FilePath: filePath,
		CallBreakdown: CallBreakdown{},
	}

	// 遍历AST
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			// 统计函数数量
			if x.Recv == nil || len(x.Recv.List) == 0 || x.Recv.List[0].Names == nil {
				// 这是一个函数或方法
				analysis.Functions++
			}
		case *ast.CallExpr:
			// 分析函数调用
			analyzeCallExpr(x, analysis, fset)
		}
		return true
	})

	// 计算总调用数
	analysis.TotalCalls = analysis.SharedCalls + analysis.ResponseCalls

	// 检测特殊场景
	// TODO: 实现WebSocket和FileDownload检测
	analysis.HasWebSocket = false
	analysis.HasFileDownload = false

	// 评估复杂度和风险
	analysis.Complexity, analysis.RiskLevel = evaluateComplexity(analysis)
	analysis.EstimatedMinutes = estimateTime(analysis)

	if verbose {
		fmt.Printf("  File: %s\n", analysis.Filename)
		fmt.Printf("    Functions: %d\n", analysis.Functions)
		fmt.Printf("    Shared calls: %d\n", analysis.SharedCalls)
		fmt.Printf("    Response calls: %d\n", analysis.ResponseCalls)
		fmt.Printf("    Complexity: %s\n", analysis.Complexity)
		fmt.Printf("    Risk level: %s\n", analysis.RiskLevel)
	}

	return analysis, nil
}

// analyzeCallExpr 分析函数调用表达式
func analyzeCallExpr(call *ast.CallExpr, analysis *FileAnalysis, fset *token.FileSet) {
	// 检查是否是 selector expression (如 shared.Error)
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	// 获取包名
	pkgName, ok := sel.X.(*ast.Ident)
	if !ok {
		return
	}

	// 获取函数名
	funcName := sel.Sel.Name

	switch pkgName.Name {
	case "shared":
		switch funcName {
		case "Success":
			analysis.CallBreakdown.SharedSuccess++
			analysis.SharedCalls++
		case "Error":
			analysis.CallBreakdown.SharedError++
			analysis.SharedCalls++
		case "ValidationError":
			analysis.CallBreakdown.SharedValidationError++
			analysis.SharedCalls++
		}
	case "response":
		switch funcName {
		case "Success":
			analysis.CallBreakdown.ResponseSuccess++
			analysis.ResponseCalls++
		case "Created":
			analysis.CallBreakdown.ResponseCreated++
			analysis.ResponseCalls++
		case "NoContent":
			analysis.CallBreakdown.ResponseNoContent++
			analysis.ResponseCalls++
		case "BadRequest":
			analysis.CallBreakdown.ResponseBadRequest++
			analysis.ResponseCalls++
		case "Unauthorized":
			analysis.CallBreakdown.ResponseUnauthorized++
			analysis.ResponseCalls++
		case "Forbidden":
			analysis.CallBreakdown.ResponseForbidden++
			analysis.ResponseCalls++
		case "NotFound":
			analysis.CallBreakdown.ResponseNotFound++
			analysis.ResponseCalls++
		case "Conflict":
			analysis.CallBreakdown.ResponseConflict++
			analysis.ResponseCalls++
		case "InternalError":
			analysis.CallBreakdown.ResponseInternal++
			analysis.ResponseCalls++
		case "Paginated":
			analysis.CallBreakdown.ResponsePaginated++
			analysis.ResponseCalls++
		}
	}
}

// hasWebSocketImport 检查是否导入了WebSocket
func hasWebSocketImport(node *ast.Node) bool {
	// TODO: 实现WebSocket检测
	return false
}

// hasFileAttachment 检查是否使用了文件附件
func hasFileAttachment(node *ast.Node) bool {
	// TODO: 实现FileAttachment检测
	return false
}

// evaluateComplexity 评估复杂度和风险等级
func evaluateComplexity(analysis *FileAnalysis) (complexity, risk string) {
	// 基于调用次数和函数数量评估
	totalCalls := analysis.SharedCalls + analysis.ResponseCalls

	if totalCalls > 50 {
		complexity = "high"
	} else if totalCalls > 20 {
		complexity = "medium"
	} else {
		complexity = "low"
	}

	// 风险评估
	if analysis.HasWebSocket || analysis.HasFileDownload {
		risk = "high"
	} else if totalCalls > 30 {
		risk = "medium"
	} else {
		risk = "low"
	}

	return
}

// estimateTime 估算迁移时间（分钟）
func estimateTime(analysis *FileAnalysis) int {
	// 基础时间：每个shared调用2分钟
	baseTime := analysis.SharedCalls * 2

	// 复杂度加成
	switch analysis.Complexity {
	case "high":
		baseTime = int(float64(baseTime) * 1.5)
	case "medium":
		baseTime = int(float64(baseTime) * 1.2)
	}

	// 特殊场景加成
	if analysis.HasWebSocket || analysis.HasFileDownload {
		baseTime += 15
	}

	// 最少10分钟
	if baseTime < 10 {
		baseTime = 10
	}

	return baseTime
}

// calculateOverallComplexity 计算整体复杂度
func calculateOverallComplexity(result *AnalysisResult) string {
	sharedCalls := result.Breakdown.SharedSuccess + result.Breakdown.SharedError + result.Breakdown.SharedValidationError

	if sharedCalls > 100 {
		return "high"
	} else if sharedCalls > 50 {
		return "medium"
	}
	return "low"
}

// generateRecommendation 生成迁移建议
func generateRecommendation(result *AnalysisResult) string {
	sharedCalls := result.Breakdown.SharedSuccess + result.Breakdown.SharedError + result.Breakdown.SharedValidationError

	if sharedCalls == 0 {
		return "所有文件已完成迁移"
	}

	if sharedCalls > 100 {
		return "建议分批迁移，优先处理高风险文件"
	}

	return "建议按复杂度从低到高逐步迁移"
}

// SaveToFile 保存分析结果到文件
func (r *AnalysisResult) SaveToFile(filename string) error {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// Print 打印分析结果
func (r *AnalysisResult) Print() {
	fmt.Printf("\n=== 分析结果: %s ===\n", r.Module)
	fmt.Printf("文件总数: %d\n", r.TotalFiles)
	fmt.Printf("总响应调用: %d\n", r.TotalCalls)
	fmt.Printf("\n调用分类:\n")
	fmt.Printf("  shared.Success:        %d\n", r.Breakdown.SharedSuccess)
	fmt.Printf("  shared.Error:          %d\n", r.Breakdown.SharedError)
	fmt.Printf("  shared.ValidationError: %d\n", r.Breakdown.SharedValidationError)
	fmt.Printf("  response.*:            %d\n", r.TotalCalls-r.Breakdown.SharedSuccess-r.Breakdown.SharedError-r.Breakdown.SharedValidationError)
	fmt.Printf("\n复杂度: %s\n", r.Complexity)
	fmt.Printf("建议: %s\n", r.Recommendation)

	if len(r.Files) > 0 {
		fmt.Printf("\n文件详情:\n")
		for _, file := range r.Files {
			fmt.Printf("  %s: %d次shared调用, %d次response调用, 复杂度=%s, 风险=%s\n",
				file.Filename, file.SharedCalls, file.ResponseCalls, file.Complexity, file.RiskLevel)
		}
	}
}
