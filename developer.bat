@echo off
chcp 65001 >nul
echo 启动开发环境...

REM 异步启动Redis
echo 正在启动 Redis...
start "Redis Server" cmd /k "%~dp0Redis\start.bat"

REM 等待一秒确保Redis已启动
timeout /t 1 /nobreak

REM 启动Go应用（带热重载）
echo 后端服务启动...
air

pause
