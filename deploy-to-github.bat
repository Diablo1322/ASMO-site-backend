@echo off
echo ========================================
echo    ASMO Backend - GitHub Deploy (HTTPS)
echo ========================================

:: Проверяем инициализацию Git
if not exist ".git" (
    echo Initializing Git repository...
    git init
)

:: Добавляем все файлы
echo.
echo Adding files to Git...
git add .

:: Проверяем статус
echo.
echo Checking Git status...
git status

:: Создаем коммит
echo.
set /p commit_msg="Enter commit message (or press Enter for default): "
if "%commit_msg%"=="" (
    set commit_msg="Deploy: %date% %time% - ASMO backend updates"
)

echo Committing changes: %commit_msg%
git commit -m %commit_msg%

:: Проверяем наличие remote
echo.
git remote -v | findstr origin >nul
if errorlevel 1 (
    echo Configuring remote repository...
    git remote add origin https://github.com/ASMO-team/ASMO-site-backend.git
    echo Remote origin set to: https://github.com/ASMO-team/ASMO-site-backend.git
)

:: Пушим изменения
echo.
echo Pushing to GitHub via HTTPS...
git push -u origin main

if %errorlevel% equ 0 (
    echo.
    echo ========================================
    echo    ✅ Successfully deployed to GitHub!
    echo ========================================
    echo.
    echo Repository: https://github.com/ASMO-team/ASMO-site-backend.git
    echo.
    echo You may need to enter your GitHub credentials.
    echo.
) else (
    echo.
    echo ========================================
    echo    ❌ Failed to push to GitHub!
    echo ========================================
    echo.
    echo Possible solutions:
    echo 1. Check if repository exists: https://github.com/ASMO-team/ASMO-site-backend
    echo 2. Verify GitHub credentials
    echo 3. Try: git push --force-with-lease origin main
    echo 4. Use Personal Access Token if 2FA is enabled
    echo.
)

pause