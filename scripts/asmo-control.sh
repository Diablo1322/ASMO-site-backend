#!/bin/bash

case "$1" in
    start)
        echo "🚀 Starting ASMO full stack..."

        # Start backend
        cd /opt/asmo-backend
        docker-compose -f docker-compose.vps.yml up -d

        # Start frontend
        /opt/asmo-backend/scripts/frontend-control.sh start

        echo "✅ All services started"
        ;;
    stop)
        echo "🛑 Stopping ASMO full stack..."

        # Stop backend
        cd /opt/asmo-backend
        docker-compose -f docker-compose.vps.yml down

        # Stop frontend
        /opt/asmo-backend/scripts/frontend-control.sh stop

        echo "✅ All services stopped"
        ;;
    restart)
        echo "🔄 Restarting ASMO full stack..."
        $0 stop
        sleep 3
        $0 start
        ;;
    status)
        echo "📊 ASMO Services Status:"
        echo ""
        echo "Backend Services:"
        cd /opt/asmo-backend
        docker-compose -f docker-compose.vps.yml ps

        echo ""
        echo "Frontend:"
        /opt/asmo-backend/scripts/frontend-control.sh status

        echo ""
        echo "Network:"
        echo "  Frontend: https://your-domain.com"
        echo "  API:      https://your-domain.com/api/health"
        echo "  Grafana:  https://your-domain.com/grafana"
        ;;
    update)
        echo "🔄 Updating ASMO full stack..."

        # Update backend
        cd /opt/asmo-backend
        git pull origin main
        docker-compose -f docker-compose.vps.yml up -d --build

        # Update frontend
        /opt/asmo-backend/scripts/frontend-control.sh update

        echo "✅ All services updated"
        ;;
    backup)
        echo "💾 Creating backup..."
        cd /opt/asmo-backend
        ./scripts/backup-db.sh
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|status|update|backup}"
        exit 1
esac