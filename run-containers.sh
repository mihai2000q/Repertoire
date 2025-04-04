DIR_SCRIPT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

docker-compose -f ${DIR_SCRIPT}/repertoire.server/docker-compose.yaml up -d
sh ${DIR_SCRIPT}/repertoire.server/apply-database-migrations.sh
sh ${DIR_SCRIPT}/repertoire.server/apply-search-migrations.sh
docker-compose -f ${DIR_SCRIPT}/repertoire.auth/docker-compose.yaml up -d
docker-compose -f ${DIR_SCRIPT}/repertoire.storage/docker-compose.yaml up -d
docker-compose -f ${DIR_SCRIPT}/repertoire.web/docker-compose.yaml up -d
