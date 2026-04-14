#!/bin/sh
if [ "$APP_ENV" = "development" ]; then
    exec air -c .air.worker.toml
else
    exec ./cmd/worker/main
fi