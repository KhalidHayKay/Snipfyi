#!/bin/sh
if [ "$APP_ENV" = "development" ]; then
    exec air
else
    exec ./cmd/web/main
fi