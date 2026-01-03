@echo off
REM 青羽写作平台AI辅助功能测试脚本
REM 运行内容总结、文本校对、敏感词检测服务的单元测试

echo ======================================
echo 青羽写作平台 AI辅助功能测试
echo ======================================
echo.

echo [1/5] 编译检查...
cd /d "D:\Github\青羽\Qingyu_backend"
go build ./service/ai/summarize_service.go
go build ./service/ai/proofread_service.go
go build ./service/ai/sensitive_words_service.go
if errorlevel 1 (
    echo 编译失败！
    pause
    exit /b 1
)
echo 编译通过 ✓
echo.

echo [2/5] 运行 Mock 适配器测试...
go test -v ./service/ai/mocks -run TestMockAIAdapter 2>&1
if errorlevel 1 (
    echo Mock 测试失败！
)
echo.

echo [3/5] 运行内容总结服务测试...
go test -v ./service/ai -run "TestSummarizeService.*" -timeout 30s 2>&1
if errorlevel 1 (
    echo 内容总结服务测试失败！
)
echo.

echo [4/5] 运行文本校对服务测试...
go test -v ./service/ai -run "TestProofreadService.*" -timeout 30s 2>&1
if errorlevel 1 (
    echo 文本校对服务测试失败！
)
echo.

echo [5/5] 运行敏感词检测服务测试...
go test -v ./service/ai -run "TestSensitiveWordsService.*" -timeout 30s 2>&1
if errorlevel 1 (
    echo 敏感词检测服务测试失败！
)
echo.

echo [6/6] 运行 API 层测试...
go test -v ./api/v1/ai -run "TestWritingAssistantApi.*" -timeout 30s 2>&1
if errorlevel 1 (
    echo API 测试失败！
)
echo.

echo ======================================
echo 测试完成！
echo ======================================
echo.
echo 测试文件位置：
echo - 服务测试：D:\Github\青羽\Qingyu_backend\service\ai\summarize_service_test.go
echo - 服务测试：D:\Github\青羽\Qingyu_backend\service\ai\proofread_service_test.go
echo - 服务测试：D:\Github\青羽\Qingyu_backend\service\ai\sensitive_words_service_test.go
echo - Mock测试：D:\Github\青羽\Qingyu_backend\service\ai\mocks\ai_adapter_mock.go
echo - API测试：D:\Github\青羽\Qingyu_backend\api\v1\ai\writing_assistant_api_test.go
echo.

pause
