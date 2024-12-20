DIR_SCRIPT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

source ${DIR_SCRIPT}/.env

export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING=postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}
export GOOSE_MIGRATION_DIR=migrations/

goose up

