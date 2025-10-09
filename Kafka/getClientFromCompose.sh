#!/bin/bash

# This script helps you get a shell inside the producer or consumer container
# Usage: ./getClientFromCompose.sh producer-1
#        ./getClientFromCompose.sh consumer-1

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 {producer-1|consumer-1}"
    exit 1
fi

SERVICE_NAME=$1

if [ "$SERVICE_NAME" = "producer-1" ]; then
    EXEC_NAME="./producer"
elif [ "$SERVICE_NAME" = "consumer-1" ]; then
    EXEC_NAME="./consumer"
else
    echo "Unknown service: $SERVICE_NAME"
    exit 1
fi

# Check if compose is up and running, if not, start it
if ! docker compose ps | grep -q "$SERVICE_NAME"; then
    echo "Starting Docker Compose services..."
    docker compose up -d
fi

docker compose exec "$SERVICE_NAME" sh -c "$EXEC_NAME"