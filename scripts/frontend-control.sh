#!/bin/bash

FRONTEND_DIR="/opt/asmo-frontend"
PORT=3001

case "$1" in
    start)
        echo "🚀 Starting frontend..."
        cd $FRONTEND_DIR
        nohup npm run start -- -p $PORT > frontend.log 2>&1 &
        echo "✅ Frontend started on port $PORT"
        ;;
    stop)
        echo "🛑 Stopping frontend..."
        pkill -f "next start"
        echo "✅ Frontend stopped"
        ;;
    restart)
        echo "🔄 Restarting frontend..."
        $0 stop
        sleep 2
        $0 start
        ;;
    status)
        echo "📊 Frontend status:"
        if pgrep -f "next start" > /dev/null; then
            echo "✅ Running on port $PORT"
            echo "📄 Log file: $FRONTEND_DIR/frontend.log"
        else
            echo "❌ Not running"
        fi
        ;;
    logs)
        echo "📄 Frontend logs:"
        tail -f $FRONTEND_DIR/frontend.log
        ;;
    update)
        echo "📥 Updating frontend..."
        cd $FRONTEND_DIR
        git pull origin main
        npm install
        npm run build
        $0 restart
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|status|logs|update}"
        exit 1
esac