#!/bin/bash
set -e

# Start API in background
/go-service &
SERVICE_PID=$!

# Trap SIGTERM/SIGINT to kill child processes
# trap "echo 'Shutting down...'; kill $API_PID $CRON_PID; exit" SIGTERM SIGINT

# Wait for both processes to finish (they shouldn't unless killed)
# wait $API_PID $CRON_PID

shutdown() {
    echo "Shutting down..."
    kill $SERVICE_PID 2>/dev/null
    wait $SERVICE_PID 2>/dev/null
    exit 0
}

trap shutdown SIGTERM SIGINT
wait