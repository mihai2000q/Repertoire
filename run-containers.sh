DIR_SCRIPT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

docker-compose -f ${DIR_SCRIPT}/repertoire.server/docker-compose.yaml up -d
docker-compose -f ${DIR_SCRIPT}/repertoire.storage/docker-compose.yaml up -d
docker-compose -f ${DIR_SCRIPT}/repertoire.web/docker-compose.yaml up -d
