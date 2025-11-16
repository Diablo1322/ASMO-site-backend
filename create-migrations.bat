@echo off
echo ========================================
echo    Creating Migration Files
echo ========================================

echo Creating migrations folder...
mkdir backend\migrations 2>nul

echo Creating migration files...

REM Create web projects table migration
echo CREATE TABLE web_projects ( > backend\migrations\001_create_web_projects_table.up.sql
echo     id SERIAL PRIMARY KEY, >> backend\migrations\001_create_web_projects_table.up.sql
echo     name VARCHAR(100) NOT NULL, >> backend\migrations\001_create_web_projects_table.up.sql
echo     description TEXT NOT NULL, >> backend\migrations\001_create_web_projects_table.up.sql
echo     img TEXT NOT NULL, >> backend\migrations\001_create_web_projects_table.up.sql
echo     price DECIMAL(10,2) NOT NULL, >> backend\migrations\001_create_web_projects_table.up.sql
echo     time_develop INTEGER NOT NULL, >> backend\migrations\001_create_web_projects_table.up.sql
echo     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, >> backend\migrations\001_create_web_projects_table.up.sql
echo     update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP >> backend\migrations\001_create_web_projects_table.up.sql
echo ); >> backend\migrations\001_create_web_projects_table.up.sql

echo DROP TABLE IF EXISTS web_projects; > backend\migrations\001_create_web_projects_table.down.sql

REM Create mobile projects table migration
echo CREATE TABLE mobile_projects ( > backend\migrations\002_create_mobile_projects_table.up.sql
echo     id SERIAL PRIMARY KEY, >> backend\migrations\002_create_mobile_projects_table.up.sql
echo     name VARCHAR(100) NOT NULL, >> backend\migrations\002_create_mobile_projects_table.up.sql
echo     description TEXT NOT NULL, >> backend\migrations\002_create_mobile_projects_table.up.sql
echo     img TEXT NOT NULL, >> backend\migrations\002_create_mobile_projects_table.up.sql
echo     price DECIMAL(10,2) NOT NULL, >> backend\migrations\002_create_mobile_projects_table.up.sql
echo     time_develop INTEGER NOT NULL, >> backend\migrations\002_create_mobile_projects_table.up.sql
echo     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, >> backend\migrations\002_create_mobile_projects_table.up.sql
echo     update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP >> backend\migrations\002_create_mobile_projects_table.up.sql
echo ); >> backend\migrations\002_create_mobile_projects_table.up.sql

echo DROP TABLE IF EXISTS mobile_projects; > backend\migrations\002_create_mobile_projects_table.down.sql

REM Create bots projects table migration
echo CREATE TABLE bots_projects ( > backend\migrations\003_create_bots_projects_table.up.sql
echo     id SERIAL PRIMARY KEY, >> backend\migrations\003_create_bots_projects_table.up.sql
echo     name VARCHAR(100) NOT NULL, >> backend\migrations\003_create_bots_projects_table.up.sql
echo     description TEXT NOT NULL, >> backend\migrations\003_create_bots_projects_table.up.sql
echo     img TEXT NOT NULL, >> backend\migrations\003_create_bots_projects_table.up.sql
echo     price DECIMAL(10,2) NOT NULL, >> backend\migrations\003_create_bots_projects_table.up.sql
echo     time_develop INTEGER NOT NULL, >> backend\migrations\003_create_bots_projects_table.up.sql
echo     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, >> backend\migrations\003_create_bots_projects_table.up.sql
echo     update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP >> backend\migrations\003_create_bots_projects_table.up.sql
echo ); >> backend\migrations\003_create_bots_projects_table.up.sql

echo DROP TABLE IF EXISTS bots_projects; > backend\migrations\003_create_bots_projects_table.down.sql

echo.
echo ✅ Migration files created successfully!
echo Location: backend\migrations\
dir backend\migrations\*.sql /b

echo.
pause