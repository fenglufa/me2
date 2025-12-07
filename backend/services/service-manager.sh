#!/bin/bash

SERVICES_DIR="$(cd "$(dirname "$0")" && pwd)"
PID_DIR="$SERVICES_DIR/.pids"
LOG_DIR="$SERVICES_DIR/.logs"

SERVICES=(action ai avatar diary event oss scheduler sms user world gateway)

mkdir -p "$PID_DIR" "$LOG_DIR"

get_port() {
    case $1 in
        user) echo 8003 ;;
        avatar) echo 8004 ;;
        world) echo 8005 ;;
        event) echo 8006 ;;
        action) echo 8007 ;;
        diary) echo 8008 ;;
        ai) echo 8009 ;;
        scheduler) echo 8010 ;;
        oss) echo 8011 ;;
        sms) echo 8012 ;;
        gateway) echo 8888 ;;
        *) echo "" ;;
    esac
}

start_service() {
    local service=$1
    local pid_file="$PID_DIR/$service.pid"
    local port=$(get_port $service)

    if [ -f "$pid_file" ] && kill -0 $(cat "$pid_file") 2>/dev/null; then
        echo "[$service] already running (PID: $(cat $pid_file))"
        return
    fi

    cd "$SERVICES_DIR/$service"
    nohup make run > "$LOG_DIR/$service.log" 2>&1 &
    sleep 2
    local actual_pid=$(lsof -ti:$port 2>/dev/null)
    echo ${actual_pid} > "$pid_file"
    echo "[$service] started (PID: $actual_pid, Port: $port)"
}

stop_service() {
    local service=$1
    local pid_file="$PID_DIR/$service.pid"

    if [ ! -f "$pid_file" ]; then
        echo "[$service] not running"
        return
    fi

    local pid=$(cat "$pid_file")
    if kill -0 $pid 2>/dev/null; then
        kill $pid
        echo "[$service] stopped (PID: $pid)"
    else
        echo "[$service] not running"
    fi
    rm -f "$pid_file"
}

status_service() {
    local service=$1
    local pid_file="$PID_DIR/$service.pid"

    if [ -f "$pid_file" ] && kill -0 $(cat "$pid_file") 2>/dev/null; then
        echo "[$service] running (PID: $(cat $pid_file))"
    else
        echo "[$service] stopped"
    fi
}

case "$1" in
    start)
        if [ -z "$2" ]; then
            for service in "${SERVICES[@]}"; do
                start_service "$service"
            done
        else
            start_service "$2"
        fi
        ;;
    stop)
        if [ -z "$2" ]; then
            for service in "${SERVICES[@]}"; do
                stop_service "$service"
            done
        else
            stop_service "$2"
        fi
        ;;
    restart)
        if [ -z "$2" ]; then
            for service in "${SERVICES[@]}"; do
                stop_service "$service"
                start_service "$service"
            done
        else
            stop_service "$2"
            start_service "$2"
        fi
        ;;
    status)
        if [ -z "$2" ]; then
            for service in "${SERVICES[@]}"; do
                status_service "$service"
            done
        else
            status_service "$2"
        fi
        ;;
    logs)
        if [ -z "$2" ]; then
            echo "Usage: $0 logs <service>"
            exit 1
        fi
        tail -f "$LOG_DIR/$2.log"
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|status|logs} [service]"
        echo "Services: ${SERVICES[*]}"
        exit 1
        ;;
esac
