#!/bin/bash

# Configuration
BACKUP_DIR="/opt/asmo/backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="asmo_backup_$DATE.sql"
RETENTION_DAYS=7

# Create backup directory
mkdir -p $BACKUP_DIR

# Backup database
docker-compose exec -T postgres pg_dump -U asmo_prod_user asmo_production > $BACKUP_DIR/$BACKUP_FILE

# Compress backup
gzip $BACKUP_DIR/$BACKUP_FILE

# Remove old backups
find $BACKUP_DIR -name "asmo_backup_*.sql.gz" -mtime +$RETENTION_DAYS -delete

# Sync with remote storage (optional)
# rclone copy $BACKUP_DIR/asmo_backup_$DATE.sql.gz remote:backups/

echo "Backup completed: $BACKUP_DIR/$BACKUP_FILE.gz"