#!/bin/bash
set -e

# Start Flask app in the background
python main.py &
FLASK_PID=$!

# Graceful shutdown function
shutdown() {
    echo "🛑 Shutting down Flask app..."
    # Send SIGTERM to Flask process
    kill -TERM "$FLASK_PID" 2>/dev/null
    # Wait up to 10 seconds for graceful exit
    wait "$FLASK_PID" 2>/dev/null &
    WAIT_PID=$!
    sleep 10
    # If still running, force kill
    if kill -0 "$FLASK_PID" 2>/dev/null; then
        echo "⚠️  Force killing Flask process..."
        kill -KILL "$FLASK_PID" 2>/dev/null
    fi
    wait "$WAIT_PID" 2>/dev/null || true
    echo "✅ Shutdown complete."
    exit 0
}

# Trap SIGTERM and SIGINT (Ctrl+C)
trap shutdown SIGTERM SIGINT

# Wait for the Flask process to finish (it won't unless killed)
wait "$FLASK_PID"