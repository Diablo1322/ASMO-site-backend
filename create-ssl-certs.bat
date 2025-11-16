@echo off
echo ========================================
echo    Creating SSL Certificates
echo ========================================

echo Creating SSL directory...
mkdir ssl 2>nul

echo Generating private key...
openssl genrsa -out ssl/asmo-backend.key 2048

echo Generating certificate signing request...
openssl req -new -key ssl/asmo-backend.key -out ssl/asmo-backend.csr -subj "/C=RU/ST=Moscow/L=Moscow/O=ASMO/CN=localhost"

echo Generating self-signed certificate...
openssl x509 -req -days 365 -in ssl/asmo-backend.csr -signkey ssl/asmo-backend.key -out ssl/asmo-backend.crt

echo Cleaning up...
del ssl\asmo-backend.csr

echo.
echo ✅ SSL certificates created:
echo   - ssl/asmo-backend.key (private key)
echo   - ssl/asmo-backend.crt (certificate)
echo.
echo For production, replace with certificates from Let's Encrypt.
pause