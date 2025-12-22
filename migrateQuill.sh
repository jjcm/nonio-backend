#!/usr/bin/env bash
set -euo pipefail

# Runs the Quill -> Markdown backfill using the same local env vars as localRun.sh.
# Usage:
#   ./migrateQuill.sh --dry-run
#   ./migrateQuill.sh            # performs updates
#   ./migrateQuill.sh --limit 500

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

cd "${SCRIPT_DIR}/cmd"
go build -o ../dist/socid

export APP_KEY="asdfa323faefjifajwiefawef"
export WEB_HOST="http://localhost:4200"
export DB_HOST="127.0.0.1"
export DB_PORT="3306"
export DB_DATABASE="socidb"
export DB_USER="dbuser"
export DB_PASSWORD="password"
export APP_PORT="4201"
export ADMIN_EMAIL="test@example.com"
export ADMIN_EMAIL_PASSWORD="password"
export EMAIL_ACCESS_TOKEN=""
export EMAIL_REFRESH_TOKEN=""
export EMAIL_CLIENT_ID=""
export EMAIL_CLIENT_SECRET=""
export STRIPE_SECRET_KEY="sk_test_51EpA4oH4gvdXgbs5rBv4JI29C38uWuNEGuB8Agt5hfya1fjgVGOQePyfj7x6ANDPE7hyYNZEMRWwkP93NAa7QTCl00GPr79F0w"
export STRIPE_PUBLISHABLE_KEY="pk_test_51EpA4oH4gvdXgbs5r0aq0i3U6IzOwbWRVYaBYXMFLLHvihVHGHotHPAi2EJ7Km9JqudFZyLE30kt2YQSUOSK88Xx00Q6eEqxmS"
export WEBHOOK_ENDPOINT_SECRET=""

# Dev-only: keep consistent with localRun.sh
export DEV_TOOLS_ENABLED="true"
export DEV_SUBSCRIPTION_PAYOUTS="true"
export PAYOUT_CYCLE_DAYS="1"

# Ensure schema is up-to-date before running the backfill
cd "${SCRIPT_DIR}/migrations"
goose mysql "${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_DATABASE}" up

cd "${SCRIPT_DIR}"
./dist/socid migrate-quill-to-markdown "$@"


