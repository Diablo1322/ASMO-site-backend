@echo off
echo ========================================
echo    ASMO Backend - GitHub Deploy
echo ========================================

:: Проверяем инициализацию Git
if not exist ".git" (
    echo Initializing Git repository...
    git init
)

:: Проверяем существование remote origin
git remote get-url origin >nul 2>&1
if errorlevel 1 (
    echo Configuring remote repository...
    git remote add origin https://github.com/ASMO-team/ASMO-site-backend.git
    echo Remote origin set to: https://github.com/ASMO-team/ASMO-site-backend.git
)

:: Проверяем текущую ветку
for /f "tokens=*" %%i in ('git branch --show-current') do set CURRENT_BRANCH=%%i

if "%CURRENT_BRANCH%"=="" (
    echo Creating initial commit on main branch...
    git checkout -b main
) else (
    echo Current branch: %CURRENT_BRANCH%
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
set /p commit_msg="Enter commit message: "
if "%commit_msg%"=="" (
    set commit_msg="Deploy: %date% %time% - ASMO backend updates"
)

echo Committing changes: %commit_msg%
git commit -m %commit_msg%

:: Пушим изменения
echo.
echo Pushing to GitHub...
git push -u origin %CURRENT_BRANCH%

if %errorlevel% equ 0 (
    echo.
    echo ========================================
    echo    ✅ Successfully deployed to GitHub!
    echo ========================================
    echo.
    echo Repository: https://github.com/ASMO-team/ASMO-site-backend
    echo Branch: %CURRENT_BRANCH%
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
    echo 3. Try: git push --force-with-lease origin %CURRENT_BRANCH%
    echo 4. Use Personal Access Token if 2FA is enabled
    echo 5. Check internet connection
    echo.
)

pause