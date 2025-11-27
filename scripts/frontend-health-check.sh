#!/bin/bash

echo "🏥 Checking frontend health..."

MAX_RETRIES=10
RETRY_COUNT=0
FRONTEND_URL="https://your-domain.com"

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" $FRONTEND_URL || true)

    if [ "$RESPONSE" = "200" ]; then
        echo "✅ Frontend health check passed"
        exit 0
    fi

    RETRY_COUNT=$((RETRY_COUNT + 1))
    echo "⏳ Waiting for frontend to be ready... ($RETRY_COUNT/$MAX_RETRIES)"
    sleep 5
done

echo "❌ Frontend health check failed after $MAX_RETRIES attempts"
exit 1