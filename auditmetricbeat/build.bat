@echo off
REM AuditMetricbeat Windows Build Script

echo =========================================
echo   AuditMetricbeat Build Script (Windows)
echo =========================================
echo.

REM 检查Go环境
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Go compiler not found
    echo Please install Go: https://golang.org/dl/
    pause
    exit /b 1
)

echo [OK] Go version:
go version
echo.

REM 设置输出文件名
set OUTPUT=auditmetricbeat.exe
if not "%1"=="" set OUTPUT=%1

REM 构建
echo Starting build...
echo.

go build -o %OUTPUT% -ldflags="-s -w" .\main.go

if %ERRORLEVEL% EQU 0 (
    echo [SUCCESS] Build completed: %OUTPUT%
    dir %OUTPUT% | findstr %OUTPUT%
) else (
    echo [FAILED] Build failed
    pause
    exit /b 1
)

echo.
echo =========================================
echo   Build Complete
echo =========================================
pause
