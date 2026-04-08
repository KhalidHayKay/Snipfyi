#!/bin/sh
if [ "$APP_ENV" = "development" ]; then
    exec air
else
    exec ./main
fi