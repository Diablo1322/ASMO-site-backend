#!/bin/bash

if [ -z "$1" ]; then
    echo "Usage: $0 <backup_file.sql.gz>"
    exit 1
fi

BACKUP_FILE=$1

echo "Restoring database from $BACKUP_FILE..."

# Stop backend to prevent connections
docker-compose stop backend

# Restore database
gunzip -c $BACKUP_FILE | docker-compose exec -T postgres psql -U asmo_prod_user asmo_production

# Start backend
docker-compose start backend

echo "Database restore completed!"