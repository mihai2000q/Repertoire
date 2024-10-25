DIR_SCRIPT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

docker-compose -f ${DIR_SCRIPT}/repertoire.server/docker-compose.yml -d up
