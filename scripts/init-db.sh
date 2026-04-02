#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
CONFIG_PATH="${GITIMPACT_CONFIG:-$ROOT_DIR/backend/config.yaml}"
DB_TYPE="${1:-auto}"

if [[ "$DB_TYPE" == "auto" ]]; then
  if [[ ! -f "$CONFIG_PATH" ]]; then
    echo "[init-db] config not found: $CONFIG_PATH"
    echo "[init-db] pass database type explicitly: ./scripts/init-db.sh mysql|dameng"
    exit 1
  fi
  DB_TYPE="$(
    awk -F: '/^[[:space:]]*type[[:space:]]*:/ {gsub(/ /, "", $2); print tolower($2); exit}' "$CONFIG_PATH" \
    | tr -d '\r\n"' \
    | sed "s/'//g"
  )"
  if [[ -z "$DB_TYPE" ]]; then
    DB_TYPE="mysql"
  fi
fi

case "$DB_TYPE" in
  mysql)
    SQL_FILE="$ROOT_DIR/sql/mysql/init.sql"
    ;;
  dameng)
    SQL_FILE="$ROOT_DIR/sql/dameng/init.sql"
    ;;
  *)
    echo "[init-db] unsupported database type: $DB_TYPE"
    echo "[init-db] allowed: mysql | dameng"
    exit 1
    ;;
esac

echo "[init-db] database.type = $DB_TYPE"
echo "[init-db] sql file      = $SQL_FILE"

echo "[init-db] Manual initialization examples:"
if [[ "$DB_TYPE" == "mysql" ]]; then
  echo "  mysql -h<host> -P<port> -u<user> -p < $SQL_FILE"
else
  echo "  Use Dameng client (disql/manager) to execute: $SQL_FILE"
fi

echo "[init-db] GitImpact startup checks core tables and will fail clearly if schema is missing."
