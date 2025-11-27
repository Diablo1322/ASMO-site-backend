#!/bin/bash

echo "🚀 Deploying ASMO Full Stack to VPS..."

# Create project directories
mkdir -p /opt/asmo-backend
mkdir -p /opt/asmo-frontend

# Deploy Backend
cd /opt/asmo-backend

# Clone backend repository
if [ -d ".git" ]; then
    git pull origin main
else
    git clone https://github.com/ASMO-team/ASMO-site-backend.git .
fi

# Create environment file
cat > .env << EOF
# Database
DB_PASSWORD=your_secure_password_here

# Monitoring
GRAFANA_PASSWORD=your_grafana_password_here

# Redis
REDIS_PASSWORD=your_redis_password_here

# CORS - фронтенд на этом же домене
ALLOWED_ORIGINS=https://your-domain.com,https://www.your-domain.com
EOF

# Get SSL certificates
certbot certonly --nginx -d your-domain.com -d www.your-domain.com

# Create SSL symlinks
mkdir -p ssl/live/your-domain.com
ln -s /etc/letsencrypt/live/your-domain.com/fullchain.pem ssl/live/your-domain.com/
ln -s /etc/letsencrypt/live/your-domain.com/privkey.pem ssl/live/your-domain.com/

# Start backend services
docker-compose -f docker-compose.vps.yml up -d --build

# Deploy Frontend
cd /opt/asmo-frontend
chmod +x /opt/asmo-backend/scripts/deploy-frontend.sh
/opt/asmo-backend/scripts/deploy-frontend.sh

# Setup backup cron and health checks
(crontab -l 2>/dev/null; echo "0 2 * * * /opt/asmo-backend/scripts/backup-db.sh") | crontab -
(crontab -l 2>/dev/null; echo "*/5 * * * * /opt/asmo-backend/scripts/health-check.sh") | crontab -

echo "✅ Full stack deployment completed!"
echo "📊 Access your services:"
echo "   Frontend:    https://your-domain.com"
echo "   Backend API: https://your-domain.com/api/health"
echo "   Grafana:     https://your-domain.com/grafana"