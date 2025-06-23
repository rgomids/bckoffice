#!/bin/sh
set -e

FILE="$(date +"%Y-%m-%d_%H%M").sql.gz"

echo "[BackupAgent] starting pg_dump to /backup/$FILE"
pg_dump -Fc -Z9 -f "/backup/$FILE"

echo "[BackupAgent] upload to bucket $MINIO_BUCKET"
mc cp "/backup/$FILE" "minio/$MINIO_BUCKET/"

# Remove dumps older than 7 days
find /backup -name '*.sql.gz' -mtime +7 -delete
