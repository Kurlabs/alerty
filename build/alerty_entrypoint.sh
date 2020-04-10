#!/usr/bin/env bash
set -e

if [ "$1" = 'brain' ]; then
    go run cmd/brain/brain.go
fi

if [ "$1" = 'runner_monitor_websites' ]; then
    go run cmd/websites-cron/run.go
fi

if [ "$1" = 'runner_monitor_sockets' ]; then
    go run cmd/sockets-cron/run.go
fi

if [ "$1" = 'runner_monitor_robots' ]; then
    go run cmd/robots-cron/run.go
fi

if [ "$1" = 'controller' ]; then
    go run cmd/controller/controller.go
fi

exec "$@"