@echo off
echo 🚀 构建MCP服务器项目...

REM 检查Go是否安装
go version >nul 2>&1
if errorlevel 1 (
    echo ❌ Go未安装，请先安装Go 1.19+
    pause
    exit /b 1
)

REM 构建Hello MCP服务器
echo 📦 构建Hello MCP服务器...
go build -o sayhi-server.exe cmd/sayhi_server/main.go
if errorlevel 1 (
    echo ❌ 构建Hello MCP服务器失败
    pause
    exit /b 1
)

REM 构建数据库MCP服务器
echo 📦 构建数据库MCP服务器...
go build -o database-server.exe cmd/database_server/main.go
if errorlevel 1 (
    echo ❌ 构建数据库MCP服务器失败
    pause
    exit /b 1
)

REM 构建Redis MCP服务器
echo 📦 构建Redis MCP服务器...
go build -o redis-server.exe cmd/redis_server/main.go
if errorlevel 1 (
    echo ❌ 构建Redis MCP服务器失败
    pause
    exit /b 1
)

REM 构建测试客户端
echo 📦 构建测试客户端...
go build -o test-client.exe examples/test_client.go
if errorlevel 1 (
    echo ❌ 构建测试客户端失败
    pause
    exit /b 1
)

echo ✅ 构建完成！
echo.
echo 📋 可执行文件：
echo   - sayhi-server.exe (Hello MCP服务器)
echo   - database-server.exe (数据库MCP服务器)
echo   - redis-server.exe (Redis MCP服务器)
echo   - test-client.exe (测试客户端)
echo.
echo 🎯 运行Hello服务器：
echo   sayhi-server.exe
echo.
echo 🗄️ 运行数据库服务器：
echo   database-server.exe
echo.
echo 🔴 运行Redis服务器：
echo   redis-server.exe
echo.
echo 🧪 运行测试：
echo   test-client.exe
echo.
pause 