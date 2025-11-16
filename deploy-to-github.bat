@echo off
echo ========================================
echo    ASMO Backend - GitHub Deploy
echo ========================================

:: Проверяем инициализацию Git
if not exist ".git" (
    echo Initializing Git repository...
    git init
)

:: Добавляем все файлы
echo Adding files to Git...
git add .

:: Проверяем статус
echo Checking Git status...
git status

:: Создаем коммит
set /p commit_msg="Enter commit message (or press Enter for default): "
if "%commit_msg%"=="" (
    set commit_msg="Auto deploy: %date% %time% - ASMO backend updates"
)

echo Committing changes...
git commit -m %commit_msg%

:: Проверяем наличие remote
git remote -v | findstr origin >nul
if errorlevel 1 (
    echo.
    echo ERROR: Remote repository not configured!
    echo.
    echo Please set up remote repository first:
    echo git remote add origin https://github.com/Diablo1322/ASMO-site-backend.git
    echo.
    echo Or create repository on GitHub first.
    echo.
    pause
    exit /b 1
)

:: Пушим изменения
echo Pushing to GitHub...
git push -u origin main

if %errorlevel% equ 0 (
    echo.
    echo ========================================
    echo    Successfully deployed to GitHub!
    echo ========================================
    echo.
) else (
    echo.
    echo ERROR: Failed to push to GitHub.
    echo Check your credentials and remote URL.
    echo.
)

pause