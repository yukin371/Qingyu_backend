@echo off
chcp 65001 >nul
set GOPROXY=https://goproxy.cn,direct
echo 正在下载Go依赖...
go mod download
if %ERRORLEVEL% EQU 0 (
    echo 依赖下载完成！
) else (
    echo 依赖下载失败！
)


