#!/bin/bash

# weatherChatBot 프로그램 경로
BIN_PATH="./weatherChatBot"
# 로그 파일 경로
LOG_FILE="./weatherChatBot.log"
# PID 파일 경로
PID_FILE="./weatherChatBot.pid"

start() {
    echo "Starting weatherChatBot ..."
    nohup $BIN_PATH > $LOG_FILE 2>&1 &
    echo $! > $PID_FILE
    echo "weatherChatBot started with PID $(cat $PID_FILE)"
}

stop() {
    if [ -f $PID_FILE ]; then
        PID=$(cat $PID_FILE)
        echo "Stopping weatherChatBot with PID $PID..."
        kill $PID
        rm $PID_FILE
        echo "weatherChatBot stopped."
    else
        echo "PID file not found. Is the weatherChatBot running?"
    fi
}

status() {
    if [ -f $PID_FILE ]; then
        PID=$(cat $PID_FILE)
        if ps -p $PID > /dev/null; then
            echo "weatherChatBot is running with PID $PID."
        else
            echo "weatherChatBot is not running, but PID file exists."
        fi
    else
        echo "weatherChatBot is not running."
    fi
}

case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    status)
        status
        ;;
    *)
        echo "Usage: $0 {start|stop|status}"
        exit 1
        ;;
esac
