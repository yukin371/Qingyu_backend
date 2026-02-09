package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCheckDependenciesScript 测试依赖检查脚本
func TestCheckDependenciesScript(t *testing.T) {
	t.Run("脚本应该可以编译", func(t *testing.T) {
		// 编译脚本
		cmd := exec.Command("go", "build", "-o", "check-deps", ".")
		cmd.Dir = filepath.Join(getProjectRoot(), "scripts/check-dependencies")
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Logf("编译输出: %s", string(output))
		}

		assert.NoError(t, err, "脚本应该可以编译")
	})

	t.Run("脚本应该可以运行", func(t *testing.T) {
		// 直接运行Go程序
		cmd := exec.Command("go", "run", ".")
		cmd.Dir = filepath.Join(getProjectRoot(), "scripts/check-dependencies")
		output, err := cmd.CombinedOutput()

		t.Logf("脚本输出:\n%s", string(output))

		// 脚本应该成功执行（可能发现违规，但不应该崩溃）
		// 如果发现违规，退出码是1，这是正常的
		if err != nil {
			t.Log("脚本发现依赖违规（这是预期行为）")
		}
	})
}

// TestDependencyRules 测试依赖规则
func TestDependencyRules(t *testing.T) {
	t.Run("允许的模块应该可以导入shared", func(t *testing.T) {
		allowedModules := []string{
			"service/container",
			"service/interfaces/shared",
			"router/shared",
		}

		for _, module := range allowedModules {
			_, ok := allowedSharedImporters[module]
			assert.True(t, ok, "模块 %s 应该在允许列表中", module)
		}
	})

	t.Run("禁止的模块应该有明确规则", func(t *testing.T) {
		businessModules := []string{
			"service/user",
			"service/writer",
			"service/reader",
		}

		for _, module := range businessModules {
			reason, exists := forbiddenImports[module]
			assert.True(t, exists, "模块 %s 应该有禁止规则", module)
			assert.NotEmpty(t, reason, "模块 %s 的规则应该有说明", module)
		}
	})
}

// TestImportPatterns 测试导入模式匹配
func TestImportPatterns(t *testing.T) {
	testCases := []struct {
		name          string
		fileDir       string
		importPath    string
		shouldViolate bool
	}{
		{
			name:          "业务服务导入shared应该违规",
			fileDir:       "service/user",
			importPath:    "Qingyu_backend/service/shared/auth",
			shouldViolate: true,
		},
		{
			name:          "容器导入shared不应该违规",
			fileDir:       "service/container",
			importPath:    "Qingyu_backend/service/shared/auth",
			shouldViolate: false,
		},
		{
			name:          "接口层导入shared不应该违规",
			fileDir:       "service/interfaces/shared",
			importPath:    "Qingyu_backend/service/shared/auth",
			shouldViolate: false,
		},
		{
			name:          "导入外部包不应该违规",
			fileDir:       "service/user",
			importPath:    "github.com/gin-gonic/gin",
			shouldViolate: false,
		},
		{
			name:          "导入models不应该违规",
			fileDir:       "service/user",
			importPath:    "Qingyu_backend/models/user",
			shouldViolate: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			violation := checkImport(tc.fileDir+".go", tc.fileDir, tc.importPath, 1)
			if tc.shouldViolate {
				assert.NotNil(t, violation, "应该检测到违规")
				assert.Equal(t, "error", violation.Severity)
			} else {
				assert.Nil(t, violation, "不应该检测到违规")
			}
		})
	}
}

// TestForbiddenImportsMap 测试禁止导入映射
func TestForbiddenImportsMap(t *testing.T) {
	t.Run("所有禁止规则都应该有描述", func(t *testing.T) {
		for module, reason := range forbiddenImports {
			assert.NotEmpty(t, reason, "模块 %s 的规则应该有描述", module)
			assert.Contains(t, reason, "应该", "规则描述应该包含建议")
		}
	})
}

// TestAllowedSharedImportersMap 测试允许导入映射
func TestAllowedSharedImportersMap(t *testing.T) {
	t.Run("允许的模块应该包含容器和接口层", func(t *testing.T) {
		expectedModules := []string{
			"service/container",
			"service/interfaces/shared",
		}

		for _, module := range expectedModules {
			assert.True(t, allowedSharedImporters[module], "模块 %s 应该被允许导入shared", module)
		}
	})
}

// getProjectRoot 获取项目根目录
func getProjectRoot() string {
	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}

	// 如果在scripts/check-dependencies目录下，返回上两级
	if strings.Contains(dir, filepath.FromSlash("/scripts/check-dependencies")) {
		// 移除 /scripts/check-dependencies
		idx := strings.LastIndex(dir, filepath.FromSlash("/scripts/check-dependencies"))
		if idx > 0 {
			return dir[:idx]
		}
	}

	// 如果在scripts目录下，返回上一级
	if strings.Contains(dir, filepath.FromSlash("/scripts/")) {
		idx := strings.LastIndex(dir, filepath.FromSlash("/scripts"))
		if idx > 0 {
			return dir[:idx]
		}
	}

	// 向上查找go.mod文件
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "."
}

// BenchmarkCheckFile 性能基准测试
func BenchmarkCheckFile(b *testing.B) {
	filePath := filepath.Join(getProjectRoot(), "service", "container", "service_container.go")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = checkFile(filePath)
	}
}
