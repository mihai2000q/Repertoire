DIRECTORY_SCRIPT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

source ${DIRECTORY_SCRIPT}/.env

export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING=postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}
export GOOSE_MIGRATION_DIR=${DIRECTORY_SCRIPT}/migrations/

goose up

