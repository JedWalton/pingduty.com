#!/bin/bash

# Load environment variables from .env file in project root
set -a
source ./.env.local.production
set +a

# Use a more robust way to extract URL components
PROTO="$(echo $POSTGRES_URL | sed -e's,^\(.*://\).*,\1,g')"
URL_NO_PROTO="$(echo ${POSTGRES_URL/$PROTO/})"
USER="$(echo $URL_NO_PROTO | sed -e's,^\([^:]*\):.*,\1,g')"
PASS="$(echo $URL_NO_PROTO | sed -e's,^[^:]*:\([^@]*\)@.*,\1,g')"
HOST_PORT_PATH="$(echo $URL_NO_PROTO | sed -e's,^[^@]*@,,g')"
HOST="$(echo $HOST_PORT_PATH | cut -d':' -f1)"
PORT="$(echo $HOST_PORT_PATH | sed -e's,[^:]*:\([^/]*\)/.*,\1,g')"
DB="$(echo $HOST_PORT_PATH | sed -e's,.*\/\(.*\),\1,g' | cut -d'?' -f1)"

# Prompt for password to avoid showing it in the command
export PGPASSWORD=$PASS

# Set the date format and backup directory
BACKUP_DATE=$(date +%Y%m%d%H%M%S)
mkdir -p ./database_backups
BACKUP_DIR="./database_backups"  # Ensure this directory exists

# Execute pg_dump
pg_dump -U "$USER" -h "$HOST" -p "$PORT" -d "$DB" > "$BACKUP_DIR/db-backup-$BACKUP_DATE.sql"

# Clear the password variable
unset PGPASSWORD

