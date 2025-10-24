@echo off
REM 阅读端功能自动化测试脚本 (Windows)

setlocal enabledelayedexpansion

REM 配置
set BASE_URL=http://localhost:8080
set API_PREFIX=/api/v1
set OUTPUT_DIR=test_results

REM 创建输出目录
if not exist "%OUTPUT_DIR%" mkdir "%OUTPUT_DIR%"

REM 测试结果统计
set /a TOTAL_TESTS=0
set /a PASSED_TESTS=0
set /a FAILED_TESTS=0

echo ========================================
echo 青羽阅读端功能自动化测试
echo ========================================
echo 测试服务器: %BASE_URL%
echo 测试时间: %date% %time%
echo.

REM ========================================
REM 一、书城浏览功能测试
REM ========================================
echo ========================================
echo 1. 书城浏览功能测试
echo ========================================

call :test_api "获取书籍列表" "/bookstore/books?page=1&pageSize=10"
call :test_api "分页测试-第2页" "/bookstore/books?page=2&pageSize=5"
call :test_api "按字数排序" "/bookstore/books?sortBy=word_count&sortOrder=desc&limit=10"
call :test_api "按章节数排序" "/bookstore/books?sortBy=chapter_count&sortOrder=desc&limit=10"
call :test_api "玄幻分类筛选" "/bookstore/books?category=玄幻&limit=10"
call :test_api "仙侠分类筛选" "/bookstore/books?category=仙侠&limit=10"
call :test_api "都市分类筛选" "/bookstore/books?category=都市&limit=10"
call :test_api "搜索功能测试" "/bookstore/books/search?keyword=书"

REM ========================================
REM 二、榜单功能测试
REM ========================================
echo.
echo ========================================
echo 2. 榜单功能测试
echo ========================================

call :test_api "热门榜-按字数" "/bookstore/books?sortBy=word_count&sortOrder=desc&limit=20"
call :test_api "热门榜-按章节" "/bookstore/books?sortBy=chapter_count&sortOrder=desc&limit=20"
call :test_api "热门书籍标记" "/bookstore/books?is_hot=true&limit=20"
call :test_api "推荐书籍" "/bookstore/books?is_recommended=true&limit=20"
call :test_api "精选书籍" "/bookstore/books?is_featured=true&limit=20"
call :test_api "新书榜" "/bookstore/books?sortBy=created_at&sortOrder=desc&limit=20"
call :test_api "最近更新榜" "/bookstore/books?sortBy=updated_at&sortOrder=desc&limit=20"

REM 分类榜单
call :test_api "玄幻热门榜" "/bookstore/books?category=玄幻&sortBy=word_count&sortOrder=desc&limit=20"
call :test_api "仙侠热门榜" "/bookstore/books?category=仙侠&sortBy=word_count&sortOrder=desc&limit=20"
call :test_api "都市热门榜" "/bookstore/books?category=都市&sortBy=word_count&sortOrder=desc&limit=20"

REM ========================================
REM 测试总结
REM ========================================
echo.
echo ========================================
echo 测试总结
echo ========================================
echo.
echo 总测试项: %TOTAL_TESTS%
echo 通过: %PASSED_TESTS%
echo 失败: %FAILED_TESTS%

if %FAILED_TESTS%==0 (
    echo [OK] 所有测试通过！
    set PASS_RATE=100
) else (
    set /a PASS_RATE=PASSED_TESTS*100/TOTAL_TESTS
    echo 通过率: !PASS_RATE!%%
)

REM 生成测试报告
set REPORT_FILE=%OUTPUT_DIR%\test_report_%date:~0,4%%date:~5,2%%date:~8,2%_%time:~0,2%%time:~3,2%%time:~6,2%.txt
(
    echo ========================================
    echo 青羽阅读端功能测试报告
    echo ========================================
    echo 测试时间: %date% %time%
    echo 测试服务器: %BASE_URL%
    echo.
    echo 测试结果统计:
    echo - 总测试项: %TOTAL_TESTS%
    echo - 通过项: %PASSED_TESTS%
    echo - 失败项: %FAILED_TESTS%
    echo - 通过率: !PASS_RATE!%%
    echo.
    echo 详细结果请查看: %OUTPUT_DIR% 目录
) > "%REPORT_FILE%"

echo.
echo 测试报告已保存: %REPORT_FILE%
echo 详细响应数据保存在: %OUTPUT_DIR% 目录
echo.

pause
goto :eof

REM ========================================
REM 测试函数
REM ========================================
:test_api
set TEST_NAME=%~1
set URL=%~2
set /a TOTAL_TESTS+=1

curl -s -o "%OUTPUT_DIR%\%TEST_NAME%.json" -w "%%{http_code}" "%BASE_URL%%API_PREFIX%%URL%" > "%OUTPUT_DIR%\http_code.txt" 2>nul

set /p HTTP_CODE=<"%OUTPUT_DIR%\http_code.txt"

if "%HTTP_CODE:~0,1%"=="2" (
    echo [OK] %TEST_NAME%
    set /a PASSED_TESTS+=1
) else (
    echo [FAIL] %TEST_NAME% ^(HTTP %HTTP_CODE%^)
    set /a FAILED_TESTS+=1
)

goto :eof

