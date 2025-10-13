#!/bin/bash

set -e

./migrate -database "${APP_DATABASE_DSN}" -path /app/cmd/server/migrations up &&
./main
