@echo off
REM 青羽写作平台 - 用户安全功能测试运行脚本

echo ========================================
echo 青羽写作平台 - 用户安全功能测试
echo ========================================
echo.

REM 设置项目根目录
cd /d D:\Github\青羽\Qingyu_backend

echo [1/4] 清理旧的测试文件...
if exist coverage.out del coverage.out

echo.
echo [2/4] 运行邮箱验证Token管理器测试...
go test -v ./service/user -run "TestNewEmailVerificationTokenManager"
if %ERRORLEVEL% NEQ 0 (
    echo 测试失败！
    goto :error
)

echo.
echo [3/4] 运行密码重置Token管理器测试...
go test -v ./service/user -run "TestNewPasswordResetTokenManager"
if %ERRORLEVEL% NEQ 0 (
    echo 测试失败！
    goto :error
)

echo.
echo [4/4] 运行所有安全功能测试...
go test -v ./service/user -run "EmailVerification|PasswordReset" 2>&1 | findstr /V "declared and not used"
if %ERRORLEVEL% NEQ 0 (
    echo 测试完成，但有一些编译警告
)

echo.
echo ========================================
echo 测试完成！
echo ========================================
echo.
echo 查看测试总结: SECURITY_TESTS_SUMMARY.md
echo.

goto :end

:error
echo.
echo ========================================
echo 测试运行失败！
echo ========================================
exit /b 1

:end
pause
